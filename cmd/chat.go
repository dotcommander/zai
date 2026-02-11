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

	"golang.org/x/sync/errgroup"

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
	spinnerStyle := theme.SpinnerStyle()
	i := 0
	for !stop.Load() {
		fmt.Fprintf(w, "\r%s %s", spinnerStyle.Render(SpinnerFrames[i%len(SpinnerFrames)]), theme.Dim.Render("Thinking...")) //nolint:errcheck // terminal output
		time.Sleep(80 * time.Millisecond)
		i++
	}
	fmt.Fprint(w, "\r\033[K") //nolint:errcheck // terminal output
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
func runChatREPL() error { //nolint:gocognit,gocyclo // TODO: decompose REPL into smaller functions
	// Set up signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize client and options
	client, baseOpts, searchEnabled := initializeChatOptions()

	// Track conversation context and history
	var conversationContext []app.Message
	var sessionHistory []string

	// Show welcome
	printWelcomeBanner(baseOpts.FilePath, searchEnabled)

	// Main REPL loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if shouldExitREPL(ctx) {
			break
		}

		input := readUserInput(scanner)
		if input == "" {
			continue
		}

		// Handle special commands
		if handled, err := handleSpecialCommands(input, &conversationContext, &sessionHistory); handled {
			if err != nil {
				fmt.Println(theme.ErrorText.Render("Error: ") + theme.Dim.Render(err.Error()))
				fmt.Println()
			}
			continue
		}

		// Handle search command
		if isSearchCommand(input) {
			if err := handleSearchCommand(ctx, client, input, &conversationContext, &sessionHistory); err != nil {
				fmt.Println(theme.ErrorText.Render("Error: ") + theme.Dim.Render(err.Error()))
				fmt.Println()
			}
			continue
		}

		// Handle web command
		if isWebCommand(input) {
			if err := handleWebCommand(ctx, client, input, &conversationContext, &sessionHistory); err != nil {
				fmt.Println(theme.ErrorText.Render("Error: ") + theme.Dim.Render(err.Error()))
				fmt.Println()
			}
			continue
		}

		// Handle regular chat message
		if err := handleRegularChat(ctx, client, baseOpts, input, searchEnabled, &conversationContext, &sessionHistory); err != nil {
			fmt.Println(theme.ErrorText.Render("Error: ") + theme.Dim.Render(err.Error()))
			fmt.Println()
			continue
		}
	}

	return nil
}

// initializeChatOptions sets up the client and base options for the chat session.
func initializeChatOptions() (*app.Client, app.ChatOptions, bool) {
	client := newClient()
	baseOpts := app.DefaultChatOptions()
	baseOpts.FilePath = viper.GetString("file")
	baseOpts.Think = viper.GetBool("think")
	baseOpts.SystemPrompt = viper.GetString("system")
	searchEnabled := viper.GetBool("search")
	return client, baseOpts, searchEnabled
}

// shouldExitREPL checks if the REPL should exit due to context cancellation.
func shouldExitREPL(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		fmt.Println()
		fmt.Println(theme.Dim.Render("Goodbye!"))
		fmt.Println()
		return true
	default:
		return false
	}
}

// readUserInput reads user input from the scanner.
func readUserInput(scanner *bufio.Scanner) string {
	fmt.Print(theme.Prompt.Render("you> "))
	if !scanner.Scan() {
		return ""
	}
	return strings.TrimSpace(scanner.Text())
}

// handleSpecialCommands handles built-in commands like exit, help, clear, etc.
func handleSpecialCommands(input string, conversationContext *[]app.Message, sessionHistory *[]string) (bool, error) {
	switch strings.ToLower(input) {
	case "exit", "quit", "/exit", "/quit":
		fmt.Println()
		fmt.Println(theme.Dim.Render("Goodbye!"))
		fmt.Println()
		return true, nil

	case "help", "/help", "?":
		printStyledHelp()
		return true, nil

	case "history", "/history":
		printSessionHistoryStyled(*sessionHistory)
		return true, nil

	case "clear", "/clear":
		*conversationContext = nil
		*sessionHistory = nil
		fmt.Print("\033[2J\033[H") // Clear screen
		printWelcomeBanner("", false)
		return true, nil

	case "context", "/context":
		printContextStyled(*conversationContext)
		return true, nil
	}
	return false, nil
}

// isSearchCommand checks if the input is a search command.
func isSearchCommand(input string) bool {
	return strings.HasPrefix(input, "/search ") || strings.HasPrefix(input, "search ")
}

// isWebCommand checks if the input is a web command.
func isWebCommand(input string) bool {
	return strings.HasPrefix(input, "/web ") || strings.HasPrefix(input, "web ")
}

// handleSearchCommand processes search commands and displays results.
func handleSearchCommand(ctx context.Context, client *app.Client, input string, conversationContext *[]app.Message, sessionHistory *[]string) error {
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
		return err
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
	*conversationContext = append(*conversationContext,
		app.Message{Role: "user", Content: fmt.Sprintf("Search: %s", query)},
		app.Message{Role: "assistant", Content: searchFormatted},
	)
	if len(*conversationContext) > 20 {
		*conversationContext = (*conversationContext)[2:]
	}

	*sessionHistory = append(*sessionHistory, input)
	return nil
}

// handleWebCommand processes web commands and displays fetched content.
func handleWebCommand(ctx context.Context, client *app.Client, input string, conversationContext *[]app.Message, sessionHistory *[]string) error {
	url := strings.TrimSpace(input[len("/web "):])
	if strings.HasPrefix(input, "web ") {
		url = strings.TrimSpace(input[len("web "):])
	}

	if url == "" {
		fmt.Println(theme.ErrorText.Render("  Usage: /web <url>"))
		fmt.Println()
		return fmt.Errorf("usage: /web <url>")
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
		return err
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
	userMsg := fmt.Sprintf("Fetched web page: %s", url)
	*conversationContext = append(*conversationContext,
		app.Message{Role: "user", Content: userMsg},
		app.Message{Role: "assistant", Content: formattedContent},
	)
	if len(*conversationContext) > 20 {
		*conversationContext = (*conversationContext)[2:]
	}

	*sessionHistory = append(*sessionHistory, input)
	return nil
}

// handleRegularChat processes regular chat messages.
func handleRegularChat(ctx context.Context, client *app.Client, baseOpts app.ChatOptions, input string, searchEnabled bool, conversationContext *[]app.Message, sessionHistory *[]string) error {
	// Add to session history
	*sessionHistory = append(*sessionHistory, input)

	// Build options with current context
	opts := baseOpts
	opts.Context = *conversationContext

	// Only include file on first message or if explicitly requested
	if len(*conversationContext) > 0 {
		opts.FilePath = ""
	}

	// If search is not enabled, proceed with regular chat
	if !searchEnabled {
		return sendChatMessage(ctx, client, input, opts, conversationContext)
	}

	// Run search and chat in parallel using errgroup
	g, ctx := errgroup.WithContext(ctx)

	// Channel for search results and error
	type searchResult struct {
		results *app.WebSearchResponse
		err     error
	}
	searchChan := make(chan searchResult, 1)

	// Start search in goroutine
	g.Go(func() error {
		searchOpts := app.SearchOptions{
			Count:         5,
			RecencyFilter: "oneWeek",
		}
		results, err := client.SearchWeb(ctx, input, searchOpts)
		searchChan <- searchResult{results: results, err: err}
		return nil
	})

	// Wait for search to complete or context to be cancelled
	var searchContext string
	var searchErr error

	select {
	case result := <-searchChan:
		searchErr = result.err
		if result.err == nil && result.results != nil && len(result.results.SearchResult) > 0 {
			searchContext = app.FormatSearchForContext(result.results.SearchResult)
		}
	case <-ctx.Done():
		return fmt.Errorf("search cancelled: %w", ctx.Err())
	}

	// If search failed, proceed with regular chat (no search context)
	messageToSend := input
	if searchErr == nil && searchContext != "" {
		messageToSend = searchContext + "\n\nUser question: " + input
	}

	// Send chat message
	return sendChatMessage(ctx, client, messageToSend, opts, conversationContext)
}

// sendChatMessage handles the actual chat API call with spinner animation
func sendChatMessage(ctx context.Context, client *app.Client, messageToSend string, opts app.ChatOptions, conversationContext *[]app.Message) error {
	// Send to API with spinner
	var stop atomic.Bool
	go animateThinking(nil, &stop)

	response, err := client.Chat(ctx, messageToSend, opts)
	stop.Store(true)
	time.Sleep(100 * time.Millisecond) // Let spinner clear

	if err != nil {
		return err
	}

	// Update conversation context (keep last 10 exchanges = 20 messages)
	*conversationContext = append(*conversationContext,
		app.Message{Role: "user", Content: messageToSend},
		app.Message{Role: "assistant", Content: response},
	)
	if len(*conversationContext) > 20 {
		*conversationContext = (*conversationContext)[2:]
	}

	// Display response with styling
	fmt.Println()
	fmt.Printf("%s %s\n", theme.AILabel.Render("AI>"), response)
	fmt.Println()

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

	query = cleanQuery
	query = strings.TrimSpace(query)
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
	if strings.Contains(s, "\n") {
		s = strings.ReplaceAll(s, "\n", " ")
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
