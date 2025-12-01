package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

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

	// Track conversation context
	var conversationContext []app.Message

	// Show welcome
	fmt.Println("🤖 Z.AI Interactive Chat")
	fmt.Println("Commands: help, history, clear, exit")
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

func printChatHelp() {
	fmt.Print(`
Commands:
  help, ?      Show this help
  history      Show session history
  context      Show conversation context
  clear        Clear conversation and screen
  exit, quit   Exit chat

Tips:
  - Previous messages are used as context
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
