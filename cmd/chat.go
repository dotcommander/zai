package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/app"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start interactive chat session (REPL)",
	Long: `Start an interactive chat session with Z.AI.

The -f flag loads a file into context for the entire session.

Examples:
	 zai chat                    # Start REPL
	 zai chat -f main.go         # Start REPL with file in context
	 zai chat --model glm-4.7    # Override model for this session`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runChatREPL()
	},
}

var chatModelFlag string

func registerChatCmd() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().StringVarP(&chatModelFlag, "model", "m", "", "Override model for this chat session")
}

// runChatREPL starts the interactive chat session. Uses bubbletea for TTY,
// falls back to raw streaming for non-TTY (pipes).
func runChatREPL() error {
	client, baseOpts, searchEnabled := initializeChatOptions()

	if !term.IsTerminal(os.Stdout.Fd()) {
		return runChatREPLRaw(client, baseOpts, searchEnabled)
	}

	m := newChatModel(client, baseOpts, searchEnabled)
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

// initializeChatOptions sets up the client and base options for the chat session.
func initializeChatOptions() (*app.Client, app.ChatOptions, bool) {
	client := newClient()
	baseOpts := app.DefaultChatOptions()
	baseOpts.FilePath = viper.GetString("file")
	baseOpts.Model = strings.TrimSpace(chatModelFlag)
	baseOpts.Think = viper.GetBool("think")

	systemVal := strings.TrimSpace(viper.GetString("system"))
	if systemVal != "" {
		resolvedSystem, fromFile, err := resolveSystemPrompt(systemVal)
		switch {
		case err != nil:
			fmt.Fprintf(os.Stderr, "warning: failed to resolve system prompt (%v); using default concise prompt\n", err)
			baseOpts.SystemPrompt = ""
		case !fromFile && looksLikeSystemPromptPath(systemVal):
			fmt.Fprintf(os.Stderr, "warning: system prompt %q looks like a file path but was not found; using default concise prompt\n", systemVal)
			baseOpts.SystemPrompt = ""
		default:
			baseOpts.SystemPrompt = resolvedSystem
		}
	}

	searchEnabled := viper.GetBool("search")
	return client, baseOpts, searchEnabled
}

// runChatREPLRaw is the non-TTY fallback: bare scanner, raw token streaming.
func runChatREPLRaw(client *app.Client, baseOpts app.ChatOptions, searchEnabled bool) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var chatCtx []app.Message
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" {
			return nil
		}

		opts := baseOpts
		opts.Context = chatCtx
		if len(chatCtx) > 0 {
			opts.FilePath = ""
		}

		// Augment with search if enabled
		messageToSend := input
		if searchEnabled {
			searchOpts := app.SearchOptions{
				Count:         5,
				RecencyFilter: "oneWeek",
				Language:      viper.GetString("web_search.language"),
			}
			results, err := client.SearchWeb(ctx, input, searchOpts)
			if err == nil && results != nil && len(results.SearchResult) > 0 {
				searchContext := app.FormatSearchForContext(results.SearchResult)
				messageToSend = searchContext + "\n\nUser question: " + input
			}
		}

		reader, err := client.StreamChat(ctx, messageToSend, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			continue
		}

		for {
			token, err := reader.Next()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Fprintf(os.Stderr, "stream error: %v\n", err)
				break
			}
			fmt.Print(token)
		}
		fmt.Println()
		reader.Close() //nolint:errcheck // best-effort

		response := reader.FullContent()
		client.SaveStreamToHistory(messageToSend, response, reader.StreamUsage())
		chatCtx = append(chatCtx,
			app.Message{Role: "user", Content: messageToSend},
			app.Message{Role: "assistant", Content: response},
		)
		if len(chatCtx) > 20 {
			chatCtx = chatCtx[2:]
		}
	}

	return scanner.Err()
}

// --- Shared utilities used by both bubbletea model and raw fallback ---

// isSearchCommand checks if the input is a search command.
func isSearchCommand(input string) bool {
	return strings.HasPrefix(input, "/search ") || strings.HasPrefix(input, "search ")
}

// isWebCommand checks if the input is a web command.
func isWebCommand(input string) bool {
	return strings.HasPrefix(input, "/web ") || strings.HasPrefix(input, "web ")
}

// stripCommandPrefix removes the leading "/cmd " or "cmd " prefix from input.
func stripCommandPrefix(input, cmd string) string {
	if after, ok := strings.CutPrefix(input, "/"+cmd+" "); ok {
		return after
	}
	if after, ok := strings.CutPrefix(input, cmd+" "); ok {
		return after
	}
	return input
}

func parseSearchCommand(input string) (query string, opts app.SearchOptions) {
	opts = app.SearchOptions{
		Count:         10,
		RecencyFilter: "noLimit",
		Language:      viper.GetString("web_search.language"),
	}

	matches := searchFlagRE.FindAllStringSubmatch(input, -1)
	cleanQuery := input
	for _, match := range matches {
		flag := match[1]
		value := match[2]

		switch flag {
		case "c", "count":
			if count, err := strconv.Atoi(value); err == nil && count > 0 && count <= 50 {
				opts.Count = count
			}
		case "r", "recency":
			opts.RecencyFilter = value
		case "d", "domain":
			opts.DomainFilter = value
		case "l", "lang":
			opts.Language = value
		}

		cleanQuery = strings.ReplaceAll(cleanQuery, match[0], "")
	}

	query = strings.TrimSpace(cleanQuery)
	return query, opts
}

// searchFlagRE matches flag/value pairs like "-c 5" or "-r oneWeek" in search commands.
var searchFlagRE = regexp.MustCompile(`-(\w+)\s*(\S+)`)

func truncate(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	r := []rune(s)
	if len(r) <= maxLen {
		return s
	}
	return string(r[:maxLen-3]) + "..."
}
