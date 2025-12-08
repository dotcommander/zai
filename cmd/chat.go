package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"zai/internal/app"
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

// runChatREPL starts the interactive chat session.
func runChatREPL() error {
	client := newClient()
	baseOpts := app.DefaultChatOptions()
	baseOpts.FilePath = viper.GetString("file")
	baseOpts.Think = viper.GetBool("think")

	// Track conversation context
	var conversationContext []app.Message

	// Show welcome
	fmt.Println("🤖 Z.AI Interactive Chat")
	fmt.Println("Commands: help, history, clear, search <query>, exit")
	if baseOpts.FilePath != "" {
		fmt.Printf("📄 File loaded: %s\n", baseOpts.FilePath)
	}
	fmt.Println(strings.Repeat("─", 40))

	scanner := bufio.NewScanner(os.Stdin)
	var sessionHistory []string

	for {
		fmt.Print("you> ")
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

			// Perform search
			fmt.Printf("🔍 Searching for: %s\n", query)
			start := time.Now()

			resp, err := client.SearchWeb(context.Background(), query, opts)
			if err != nil {
				fmt.Printf("❌ Search failed: %v\n", err)
				continue
			}

			duration := time.Since(start)
			fmt.Printf("⏱️  Found %d results in %v\n\n", len(resp.SearchResult), duration)

			// Format and display results
			for i, result := range resp.SearchResult {
				fmt.Printf("%d. %s\n", i+1, result.Title)
				fmt.Printf("   %s\n", result.Link)
				if result.PublishDate != "" {
					fmt.Printf("   📅 %s\n", result.PublishDate)
				}
				if result.Content != "" {
					content := result.Content
					if len(content) > 200 {
						content = content[:200] + "..."
					}
					fmt.Printf("   %s\n\n", content)
				}
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

		// Handle commands
		switch strings.ToLower(input) {
		case "exit", "quit", "/exit", "/quit":
			fmt.Println("👋 Goodbye!")
			return nil

		case "help", "/help", "?":
			printChatHelp()
			continue

		case "history", "/history":
			printSessionHistory(sessionHistory)
			continue

		case "clear", "/clear":
			conversationContext = nil
			sessionHistory = nil
			fmt.Print("\033[2J\033[H") // Clear screen
			fmt.Println("🤖 Z.AI Interactive Chat (cleared)")
			fmt.Println(strings.Repeat("─", 40))
			continue

		case "context", "/context":
			printContext(conversationContext)
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

		// Send to API
		response, err := client.Chat(context.Background(), input, opts)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
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

		fmt.Printf("🤖 %s\n\n", response)
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

func printChatHelp() {
	fmt.Print(`
Commands:
  help, ?           Show this help
  history           Show session history
  context           Show conversation context
  clear             Clear conversation and screen
  search <query>    Search the web
  /search <query>   Search the web (alternative format)
  exit, quit        Exit chat

Search options:
  search "golang" -c 5 -r oneWeek
  /search "AI news" -d github.com

Tips:
  - Previous messages are used as context
  - Search results are automatically added to conversation
  - Use -f flag to include a file in context
  - Arrow keys for input history (if supported)
`)
}

func printSessionHistory(history []string) {
	if len(history) == 0 {
		fmt.Println("No messages yet.")
		return
	}
	fmt.Printf("\n📜 Session History (%d messages):\n", len(history))
	fmt.Println(strings.Repeat("─", 40))
	for i, msg := range history {
		fmt.Printf("%2d: %s\n", i+1, truncate(msg, 60))
	}
	fmt.Println()
}

func printContext(ctx []app.Message) {
	if len(ctx) == 0 {
		fmt.Println("No context yet.")
		return
	}
	fmt.Printf("\n💬 Conversation Context (%d messages):\n", len(ctx))
	fmt.Println(strings.Repeat("─", 40))
	for _, msg := range ctx {
		role := "You"
		if msg.Role == "assistant" {
			role = "AI"
		}
		fmt.Printf("[%s] %s\n", role, truncate(msg.Content, 50))
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
