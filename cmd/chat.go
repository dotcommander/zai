package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/garyblankenship/zai/internal/app"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start interactive chat session (REPL)",
	Long: `Start an interactive chat session with Z.AI.

The -f flag loads a file into context for the entire session.

Examples:
  zai chat                    # Start REPL
  zai chat -f main.go         # Start REPL with file in context`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runChatREPL()
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}

// animateThinking displays an animated spinner while waiting for API response.
// Accepts io.Writer for testability (pass nil to use os.Stdout).
func animateThinking(w io.Writer, stop *atomic.Bool) {
	if w == nil {
		w = os.Stdout
	}
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerStyle := theme.SpinnerStyle()
	i := 0
	for !stop.Load() {
		fmt.Fprintf(w, "\r%s %s", spinnerStyle.Render(frames[i%len(frames)]), theme.Dim.Render("Thinking..."))
		time.Sleep(80 * time.Millisecond)
		i++
	}
	fmt.Fprint(w, "\r\033[K") // Clear line
}

// printWelcomeBanner displays the styled welcome message.
func printWelcomeBanner(filePath string, searchEnabled bool) {
	fmt.Println()
	fmt.Println(theme.Title.Render(" Z.AI Chat "))
	fmt.Println()

	if filePath != "" {
		fmt.Println(theme.Info.Render("  File: ") + theme.Dim.Render(filePath))
	}
	if searchEnabled {
		fmt.Println(theme.Info.Render("  Search: ") + theme.Dim.Render("enabled (answers include web search)"))
	}

	fmt.Println()
	fmt.Println(theme.HelpText.Render("  Commands: help, history, clear, search <query>, exit"))
	fmt.Println(theme.Divider.Render(strings.Repeat("─", 50)))
	fmt.Println()
}

// printStyledHelp displays the help text with styling.
func printStyledHelp() {
	fmt.Println()
	fmt.Println(theme.Section.Render("Commands"))
	fmt.Println(theme.Divider.Render(strings.Repeat("─", 40)))

	commands := []struct {
		cmd  string
		desc string
	}{
		{"help, ?", "Show this help"},
		{"history", "Show session history"},
		{"context", "Show conversation context"},
		{"clear", "Clear conversation and screen"},
		{"search <query>", "Search the web"},
		{"web <url>", "Fetch and display web page"},
		{"exit, quit", "Exit chat"},
	}

	for _, c := range commands {
		fmt.Printf("  %s  %s\n",
			theme.Info.Render(fmt.Sprintf("%-16s", c.cmd)),
			theme.Dim.Render(c.desc))
	}

	fmt.Println()
	fmt.Println(theme.Section.Render("Search Options"))
	fmt.Println(theme.Divider.Render(strings.Repeat("─", 40)))
	fmt.Println(theme.Dim.Render(`  search "golang" -c 5 -r oneWeek
  /search "AI news" -d github.com`))

	fmt.Println()
	fmt.Println(theme.Section.Render("Tips"))
	fmt.Println(theme.Divider.Render(strings.Repeat("─", 40)))
	tips := []string{
		"Previous messages are used as context",
		"URLs in messages are auto-fetched",
		"Web/search results are added to context",
		"Use --search flag to auto-search every message",
	}
	for _, tip := range tips {
		fmt.Println(theme.Dim.Render("  - " + tip))
	}
	fmt.Println()
}

// runChatREPL starts the interactive chat session.
func runChatREPL() error {
	// Set up signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := newClient()
	baseOpts := app.DefaultChatOptions()
	baseOpts.FilePath = viper.GetString("file")
	baseOpts.Think = viper.GetBool("think")
	searchEnabled := viper.GetBool("search")

	// Track conversation context
	var conversationContext []app.Message

	// Show welcome
	printWelcomeBanner(baseOpts.FilePath, searchEnabled)

	scanner := bufio.NewScanner(os.Stdin)
	var sessionHistory []string

	for {
		// Check if context was cancelled (Ctrl+C)
		select {
		case <-ctx.Done():
			fmt.Println()
			fmt.Println(theme.Dim.Render("Goodbye!"))
			fmt.Println()
			return nil
		default:
		}

		fmt.Print(theme.Prompt.Render("you> "))
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Check for search command
		if strings.HasPrefix(input, "/search ") || strings.HasPrefix(input, "search ") {
			query := strings.TrimSpace(input[len("/search "):])
			if strings.HasPrefix(input, "search ") {
				query = strings.TrimSpace(input[len("search "):])
			}

			// Parse search options
			query, opts := parseSearchCommand(query)

			// Perform search with spinner
			fmt.Println()
			fmt.Println(theme.Info.Render("  Searching: ") + theme.Dim.Render(query))

			var stop atomic.Bool
			go animateThinking(nil, &stop)

			start := time.Now()
			resp, err := client.SearchWeb(ctx, query, opts)
			stop.Store(true)
			time.Sleep(100 * time.Millisecond) // Let spinner clear

			if err != nil {
				fmt.Println(theme.ErrorText.Render("  Error: ") + theme.Dim.Render(err.Error()))
				fmt.Println()
				continue
			}

			duration := time.Since(start)
			fmt.Println(theme.Dim.Render(fmt.Sprintf("  Found %d results in %v", len(resp.SearchResult), duration.Round(time.Millisecond))))
			fmt.Println()

			// Format and display results
			for i, result := range resp.SearchResult {
				fmt.Printf("  %s %s\n",
					theme.Dim.Render(fmt.Sprintf("%d.", i+1)),
					theme.ResultTitle.Render(result.Title))
				fmt.Printf("     %s\n", theme.ResultLink.Render(result.Link))
				if result.PublishDate != "" {
					fmt.Printf("     %s\n", theme.ResultDate.Render(result.PublishDate))
				}
				if result.Content != "" {
					content := result.Content
					if len(content) > 200 {
						content = content[:200] + "..."
					}
					fmt.Printf("     %s\n", theme.Dim.Render(content))
				}
				fmt.Println()
			}

			// Add formatted search results to conversation
			searchFormatted := app.FormatSearchResultsForChat(resp.SearchResult, query)
			conversationContext = append(conversationContext,
				app.Message{Role: "user", Content: fmt.Sprintf("Search: %s", query)},
				app.Message{Role: "assistant", Content: searchFormatted},
			)
			if len(conversationContext) > 20 {
				conversationContext = conversationContext[2:]
			}

			sessionHistory = append(sessionHistory, input)
			continue
		}

		// Check for web command
		if strings.HasPrefix(input, "/web ") || strings.HasPrefix(input, "web ") {
			url := strings.TrimSpace(input[len("/web "):])
			if strings.HasPrefix(input, "web ") {
				url = strings.TrimSpace(input[len("web "):])
			}

			if url == "" {
				fmt.Println(theme.ErrorText.Render("  Usage: /web <url>"))
				fmt.Println()
				continue
			}

			// Fetch web content with spinner
			fmt.Println()
			fmt.Println(theme.Info.Render("  Fetching: ") + theme.ResultLink.Render(url))

			var stop atomic.Bool
			go animateThinking(nil, &stop)

			webOpts := &app.WebReaderOptions{
				ReturnFormat: "markdown",
			}
			resp, err := client.FetchWebContent(ctx, url, webOpts)
			stop.Store(true)
			time.Sleep(100 * time.Millisecond) // Let spinner clear

			if err != nil {
				fmt.Println(theme.ErrorText.Render("  Error: ") + theme.Dim.Render(err.Error()))
				fmt.Println()
				continue
			}

			// Display content
			fmt.Println()
			fmt.Println(theme.Section.Render("  " + resp.ReaderResult.Title))
			fmt.Println(theme.ResultLink.Render("  " + resp.ReaderResult.URL))
			fmt.Println()

			// Truncate content for display
			content := resp.ReaderResult.Content
			if len(content) > 2000 {
				content = content[:2000] + "\n\n" + theme.Dim.Render("[Content truncated - full content added to context]")
			}
			fmt.Println(theme.Dim.Render(content))
			fmt.Println()

			// Add to conversation context
			formattedContent := app.FormatWebContent(url, resp.ReaderResult.Title, resp.ReaderResult.Content)
			conversationContext = append(conversationContext,
				app.Message{Role: "user", Content: fmt.Sprintf("Fetched web page: %s", url)},
				app.Message{Role: "assistant", Content: formattedContent},
			)
			if len(conversationContext) > 20 {
				conversationContext = conversationContext[2:]
			}

			sessionHistory = append(sessionHistory, input)
			continue
		}

		// Handle commands
		switch strings.ToLower(input) {
		case "exit", "quit", "/exit", "/quit":
			fmt.Println()
			fmt.Println(theme.Dim.Render("Goodbye!"))
			fmt.Println()
			return nil

		case "help", "/help", "?":
			printStyledHelp()
			continue

		case "history", "/history":
			printSessionHistoryStyled(sessionHistory)
			continue

		case "clear", "/clear":
			conversationContext = nil
			sessionHistory = nil
			fmt.Print("\033[2J\033[H") // Clear screen
			printWelcomeBanner("", false)
			continue

		case "context", "/context":
			printContextStyled(conversationContext)
			continue
		}

		// Add to session history
		sessionHistory = append(sessionHistory, input)

		// Build options with current context
		opts := baseOpts
		opts.Context = conversationContext

		// Only include file on first message or if explicitly requested
		if len(conversationContext) > 0 {
			opts.FilePath = ""
		}

		// Augment with web search if enabled
		messageToSend := input
		if searchEnabled {
			searchOpts := app.SearchOptions{
				Count:         5,
				RecencyFilter: "oneWeek",
			}
			results, err := client.SearchWeb(ctx, input, searchOpts)
			if err == nil && len(results.SearchResult) > 0 {
				searchContext := app.FormatSearchForContext(results.SearchResult)
				messageToSend = searchContext + "\n\nUser question: " + input
			}
		}

		// Send to API with spinner
		var stop atomic.Bool
		go animateThinking(nil, &stop)

		response, err := client.Chat(ctx, messageToSend, opts)
		stop.Store(true)
		time.Sleep(100 * time.Millisecond) // Let spinner clear

		if err != nil {
			fmt.Println(theme.ErrorText.Render("Error: ") + theme.Dim.Render(err.Error()))
			fmt.Println()
			continue
		}

		// Update conversation context (keep last 10 exchanges = 20 messages)
		conversationContext = append(conversationContext,
			app.Message{Role: "user", Content: input},
			app.Message{Role: "assistant", Content: response},
		)
		if len(conversationContext) > 20 {
			conversationContext = conversationContext[2:]
		}

		// Display response with styling
		fmt.Println()
		fmt.Printf("%s %s\n", theme.AILabel.Render("AI>"), response)
		fmt.Println()
	}

	return nil
}

func parseSearchCommand(input string) (query string, opts app.SearchOptions) {
	// Default options
	opts = app.SearchOptions{
		Count:         10,
		RecencyFilter: "noLimit",
	}

	// Parse flags
	flagRegex := regexp.MustCompile(`-(\w+)\s*(\S+)`)
	matches := flagRegex.FindAllStringSubmatch(input, -1)

	// Remove flags from query
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
		}

		// Remove this flag from query
		cleanQuery = strings.ReplaceAll(cleanQuery, match[0], "")
	}

	query = strings.TrimSpace(cleanQuery)
	return query, opts
}

func printSessionHistoryStyled(history []string) {
	fmt.Println()
	if len(history) == 0 {
		fmt.Println(theme.Dim.Render("  No messages yet."))
		fmt.Println()
		return
	}

	fmt.Println(theme.Section.Render(fmt.Sprintf("Session History (%d messages)", len(history))))
	fmt.Println(theme.Divider.Render(strings.Repeat("─", 40)))

	for i, msg := range history {
		fmt.Printf("  %s %s\n",
			theme.Dim.Render(fmt.Sprintf("%2d.", i+1)),
			truncate(msg, 60))
	}
	fmt.Println()
}

func printContextStyled(ctx []app.Message) {
	fmt.Println()
	if len(ctx) == 0 {
		fmt.Println(theme.Dim.Render("  No context yet."))
		fmt.Println()
		return
	}

	fmt.Println(theme.Section.Render(fmt.Sprintf("Conversation Context (%d messages)", len(ctx))))
	fmt.Println(theme.Divider.Render(strings.Repeat("─", 40)))

	for _, msg := range ctx {
		var roleName string
		var styledRole string
		if msg.Role == "user" {
			roleName = "You"
			styledRole = theme.Prompt.Render(fmt.Sprintf("[%s]", roleName))
		} else {
			roleName = "AI"
			styledRole = theme.AILabel.Render(fmt.Sprintf("[%s]", roleName))
		}
		fmt.Printf("  %s %s\n",
			styledRole,
			theme.Dim.Render(truncate(msg.Content, 50)))
	}
	fmt.Println()
}

func truncate(s string, maxLen int) string {
	// Remove newlines for display
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
