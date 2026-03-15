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
	 zai chat -f main.go         # Start REPL with file in context
	 zai chat --model glm-4.7    # Override model for this session`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runChatREPL()
	},
}

var chatModel string

// chatSession holds the mutable state for a single chat REPL session.
type chatSession struct {
	client        *app.Client
	searchEnabled bool
	context       []app.Message
	history       []string
}

func registerChatCmd() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().StringVarP(&chatModel, "model", "m", "", "Override model for this chat session")
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
func runChatREPL() error { //nolint:gocognit,gocyclo,revive,funlen // REPL loop - signal handling, goroutine, and loop are tightly coupled
	// Set up signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize client and options
	client, baseOpts, searchEnabled := initializeChatOptions()

	// Track conversation context and history
	sess := &chatSession{
		client:        client,
		searchEnabled: searchEnabled,
	}

	// Show welcome
	printWelcomeBanner(baseOpts.FilePath, searchEnabled)

	// Read input on a separate goroutine so Ctrl-C can exit immediately.
	inputCh := make(chan string)
	inputErrCh := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inputCh <- strings.TrimSpace(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			inputErrCh <- err
		}
		close(inputCh)
	}()

	// Main REPL loop
	for {
		fmt.Print(theme.Prompt.Render("you> "))

		var input string
		select {
		case <-ctx.Done():
			fmt.Println()
			fmt.Println(theme.Dim.Render("Goodbye!"))
			fmt.Println()
			return nil
		case err := <-inputErrCh:
			return fmt.Errorf("failed to read chat input: %w", err)
		case line, ok := <-inputCh:
			if !ok {
				fmt.Println()
				fmt.Println(theme.Dim.Render("Goodbye!"))
				fmt.Println()
				return nil
			}
			input = line
		}

		if input == "" {
			continue
		}

		// Handle special commands
		handled, shouldExit := handleSpecialCommands(input, sess)
		if handled {
			if shouldExit {
				fmt.Println()
				fmt.Println(theme.Dim.Render("Goodbye!"))
				fmt.Println()
				return nil
			}
			continue
		}

		// Handle search command
		if isSearchCommand(input) {
			if err := handleSearchCommand(ctx, sess, input); err != nil {
				fmt.Println(theme.ErrorText.Render("Error: ") + theme.Dim.Render(err.Error()))
				fmt.Println()
			}
			continue
		}

		// Handle web command
		if isWebCommand(input) {
			if err := handleWebCommand(ctx, sess, input); err != nil {
				fmt.Println(theme.ErrorText.Render("Error: ") + theme.Dim.Render(err.Error()))
				fmt.Println()
			}
			continue
		}

		// Handle regular chat message
		if err := handleRegularChat(ctx, sess, baseOpts, input); err != nil {
			fmt.Println(theme.ErrorText.Render("Error: ") + theme.Dim.Render(err.Error()))
			fmt.Println()
		}
	}
}

// initializeChatOptions sets up the client and base options for the chat session.
func initializeChatOptions() (*app.Client, app.ChatOptions, bool) {
	client := newClient()
	baseOpts := app.DefaultChatOptions()
	baseOpts.FilePath = viper.GetString("file")
	baseOpts.Model = strings.TrimSpace(chatModel)
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

// handleSpecialCommands handles built-in commands like exit, help, clear, etc.
// Returns (handled, shouldExit).
func handleSpecialCommands(input string, sess *chatSession) (handled bool, shouldExit bool) {
	switch strings.ToLower(input) {
	case "exit", "quit", "/exit", "/quit":
		return true, true

	case "help", "/help", "?":
		printStyledHelp()
		return true, false

	case "history", "/history":
		printSessionHistoryStyled(sess.history)
		return true, false

	case "clear", "/clear":
		sess.context = nil
		sess.history = nil
		fmt.Print("\033[2J\033[H") // Clear screen
		printWelcomeBanner("", false)
		return true, false

	case "context", "/context":
		printContextStyled(sess.context)
		return true, false
	}
	return false, false
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
func handleSearchCommand(ctx context.Context, sess *chatSession, input string) error {
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
	resp, err := sess.client.SearchWeb(ctx, query, opts)
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
	sess.context = append(sess.context,
		app.Message{Role: "user", Content: fmt.Sprintf("Search: %s", query)},
		app.Message{Role: "assistant", Content: searchFormatted},
	)
	if len(sess.context) > 20 {
		sess.context = sess.context[2:]
	}

	sess.history = append(sess.history, input)
	return nil
}

// handleWebCommand processes web commands and displays fetched content.
func handleWebCommand(ctx context.Context, sess *chatSession, input string) error {
	url := strings.TrimSpace(input[len("/web "):])
	if strings.HasPrefix(input, "web ") {
		url = strings.TrimSpace(input[len("web "):])
	}

	if url == "" {
		fmt.Println(theme.ErrorText.Render("  Usage: /web <url>"))
		fmt.Println()
		return errors.New("usage: /web <url>")
	}

	// Fetch web content with spinner
	fmt.Println()
	fmt.Println(theme.Info.Render("  Fetching: ") + theme.ResultLink.Render(url))

	var stop atomic.Bool
	go animateThinking(nil, &stop)

	webOpts := &app.WebReaderOptions{
		ReturnFormat: "markdown",
	}
	resp, err := sess.client.FetchWebContent(ctx, url, webOpts)
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
	sess.context = append(sess.context,
		app.Message{Role: "user", Content: userMsg},
		app.Message{Role: "assistant", Content: formattedContent},
	)
	if len(sess.context) > 20 {
		sess.context = sess.context[2:]
	}

	sess.history = append(sess.history, input)
	return nil
}

// handleRegularChat processes regular chat messages.
func handleRegularChat(ctx context.Context, sess *chatSession, baseOpts app.ChatOptions, input string) error {
	// Add to session history
	sess.history = append(sess.history, input)

	// Build options with current context
	opts := baseOpts
	opts.Context = sess.context

	// Only include file on first message or if explicitly requested
	if len(sess.context) > 0 {
		opts.FilePath = ""
	}

	// If search is not enabled, proceed with regular chat
	if !sess.searchEnabled {
		return sendChatMessage(ctx, sess, input, opts)
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
			Language:      viper.GetString("web_search.language"),
		}
		results, err := sess.client.SearchWeb(ctx, input, searchOpts)
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
	return sendChatMessage(ctx, sess, messageToSend, opts)
}

// sendChatMessage handles the actual chat API call with streaming output
func sendChatMessage(ctx context.Context, sess *chatSession, messageToSend string, opts app.ChatOptions) error {
	reader, err := sess.client.StreamChat(ctx, messageToSend, opts)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Print AI label, then stream tokens inline
	fmt.Println()
	fmt.Printf("%s ", theme.AILabel.Render("AI>"))

	for {
		token, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("stream error: %w", err)
		}
		fmt.Print(token)
	}
	fmt.Println()
	fmt.Println()

	// Get accumulated response for context and history
	response := reader.FullContent()

	// Save to history
	sess.client.SaveStreamToHistory(messageToSend, response, reader.StreamUsage())

	// Update conversation context (keep last 10 exchanges = 20 messages)
	sess.context = append(sess.context,
		app.Message{Role: "user", Content: messageToSend},
		app.Message{Role: "assistant", Content: response},
	)
	if len(sess.context) > 20 {
		sess.context = sess.context[2:]
	}

	return nil
}

func parseSearchCommand(input string) (query string, opts app.SearchOptions) {
	// Default options
	opts = app.SearchOptions{
		Count:         10,
		RecencyFilter: "noLimit",
		Language:      viper.GetString("web_search.language"),
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
		case "l", "lang":
			opts.Language = value
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
