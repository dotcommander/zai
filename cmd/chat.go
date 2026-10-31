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
	if value := strings.TrimSpace(viper.GetString("system")); value != "" {
		resolved, _, err := resolveSystemPrompt(value)
		if err != nil {
			return err
		}
		baseOpts.SystemPrompt = resolved
	}

	if !term.IsTerminal(os.Stdout.Fd()) || !term.IsTerminal(os.Stdin.Fd()) {
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
	client.EnableSessionFileReuse(baseOpts.FilePath)
	baseOpts.Model = strings.TrimSpace(chatModelFlag)
	if baseOpts.Model == "" {
		baseOpts.Model = viper.GetString("api.model")
	}
	baseOpts.Think = viper.GetBool("think")

	systemVal := strings.TrimSpace(viper.GetString("system"))
	if systemVal != "" {
		resolvedSystem, _, err := resolveSystemPrompt(systemVal)
		if err == nil {
			baseOpts.SystemPrompt = resolvedSystem
		}
	}

	searchEnabled := viper.GetBool("search")
	return client, baseOpts, searchEnabled
}

// runChatREPLRaw is the non-TTY fallback: bare scanner, raw token streaming.
func runChatREPLRaw(client *app.Client, baseOpts app.ChatOptions, searchEnabled bool) error {

	var chatCtx []app.Message
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 64*1024), MaxStdinSize)
	var sessionHist []string

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" || input == "/exit" || input == "/quit" {
			return nil
		}

		local := strings.TrimPrefix(input, "/")
		switch local {
		case "clear":
			chatCtx = nil
			sessionHist = nil
			continue
		case "help", "?":
			fmt.Println((&chatModel{}).helpText())
			continue
		case "history":
			fmt.Println((&chatModel{sessionHist: sessionHist}).historyText())
			continue
		case "context":
			fmt.Println((&chatModel{chatCtx: chatCtx}).contextText())
			continue
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		sessionHist = append(sessionHist, input)
		if isSearchCommand(input) {
			query, options := parseSearchCommand(stripCommandPrefix(input, "search"))
			results, err := client.SearchWeb(ctx, query, options)
			stop()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			text := app.FormatSearchResultsForChat(results.SearchResult, query)
			fmt.Println(text)
			chatCtx = append(chatCtx, app.Message{Role: "user", Content: input}, app.Message{Role: "assistant", Content: text})
			continue
		}
		if isWebCommand(input) {
			url := strings.TrimSpace(stripCommandPrefix(input, "web"))
			response, err := client.FetchWebContent(ctx, url, &app.WebReaderOptions{ReturnFormat: "markdown"})
			stop()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			text := app.FormatWebContent(url, response.ReaderResult.Title, response.ReaderResult.Content)
			fmt.Println(text)
			chatCtx = append(chatCtx, app.Message{Role: "user", Content: input}, app.Message{Role: "assistant", Content: text})
			continue
		}
		opts := baseOpts
		opts.Context = chatCtx

		messageToSend := augmentWithSearch(ctx, client, input, searchEnabled)

		reader, err := client.StreamChat(ctx, messageToSend, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			stop()
			continue
		}

		streamErr := readStreamToStdout(reader)
		fmt.Println()
		reader.Close() //nolint:errcheck // best-effort

		stop()
		if streamErr != nil {
			continue
		}
		response := reader.FullContent()
		client.SaveStreamToHistory(messageToSend, response, reader.StreamUsage(), reader.Model())
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

func augmentWithSearch(ctx context.Context, client *app.Client, input string, searchEnabled bool) string {
	if !searchEnabled {
		return input
	}
	searchOpts := app.SearchOptions{
		Count:         5,
		RecencyFilter: "oneWeek",
		Language:      viper.GetString("web_search.language"),
	}
	results, err := client.SearchWeb(ctx, input, searchOpts)
	if err == nil && results != nil && len(results.SearchResult) > 0 {
		searchContext := app.FormatSearchForContext(results.SearchResult)
		return searchContext + "\n\nUser question: " + input
	}
	return input
}

func readStreamToStdout(reader *app.StreamReader) error {
	for {
		token, err := reader.Next()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Fprintf(os.Stderr, "stream error: %v\n", err)
			}
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		fmt.Print(token)
	}
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
	if maxLen <= 0 {
		return ""
	}
	if len(r) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return string(r[:maxLen])
	}
	return string(r[:maxLen-3]) + "..."
}
