package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"zai/internal/app"
)

// Help styles (reuse colors from chat.go where applicable)
var (
	helpTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	helpSectionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7D56F4"))

	helpCommandStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#73F59F"))

	helpFlagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00D4FF"))

	helpDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA"))

	helpExampleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFD700")).
				Italic(true)

	helpDividerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#444444"))
)

var (
	cfgFile    string
	verbose    bool
	filePath   string
	think      bool
	jsonOutput bool
	search     bool
)

var rootCmd = &cobra.Command{
	Use:   "zai [prompt]",
	Short: "Chat with Z.AI GLM model",
	Long: `ZAI is a CLI tool for chatting with Z.AI models.

One-shot mode:
  zai "Explain quantum computing"
  zai -f main.go "Explain this code"

Piped input:
  pbpaste | zai "explain this"
  cat file.txt | zai "summarize"
  echo "Hello" | zai

Interactive REPL:
  zai chat
  zai chat -f main.go

History:
  zai history`,
	Args: cobra.ArbitraryArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config init for commands that don't need API
		if cmd.Name() == "history" || cmd.Name() == "completion" || cmd.Name() == "help" {
			return nil
		}
		return initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var prompt string
		var stdinData string

		// Check for stdin data (piped input)
		if hasStdinData() {
			data, err := readStdin()
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			stdinData = data
		}

		// Build prompt from args
		if len(args) > 0 {
			prompt = strings.Join(args, " ")
		}

		// Combine: prompt + stdin (stdin as context with XML tags)
		if stdinData != "" {
			if prompt != "" {
				prompt = prompt + "\n\n<stdin>\n" + stdinData + "\n</stdin>"
			} else {
				prompt = stdinData
			}
		}

		// Require some input
		if prompt == "" {
			return cmd.Help()
		}

		return runOneShot(prompt)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpFunc(styledHelp)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.config/zai/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "include file contents in prompt")
	rootCmd.PersistentFlags().BoolVar(&think, "think", false, "enable thinking/reasoning mode")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&search, "search", false, "augment prompt with web search results")

	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("file", rootCmd.PersistentFlags().Lookup("file"))
	_ = viper.BindPFlag("think", rootCmd.PersistentFlags().Lookup("think"))
	_ = viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
	_ = viper.BindPFlag("search", rootCmd.PersistentFlags().Lookup("search"))
}

// styledHelp displays the custom styled help output.
func styledHelp(cmd *cobra.Command, args []string) {
	// Title
	fmt.Println()
	fmt.Println(helpTitleStyle.Render(" ZAI ") + " " + helpDescStyle.Render("Chat with Z.AI models"))
	fmt.Println()

	// Examples section
	fmt.Println(helpSectionStyle.Render("Examples"))
	fmt.Println(helpDividerStyle.Render(strings.Repeat("-", 50)))
	examples := []string{
		`zai "Explain quantum computing"`,
		`zai -f main.go "Review this code"`,
		`zai --search "Latest AI news"`,
		`echo "text" | zai "summarize"`,
	}
	for _, ex := range examples {
		fmt.Printf("  %s\n", helpExampleStyle.Render(ex))
	}
	fmt.Println()

	// Commands section
	fmt.Println(helpSectionStyle.Render("Commands"))
	fmt.Println(helpDividerStyle.Render(strings.Repeat("-", 50)))
	commands := [][]string{
		{"chat", "Interactive chat session (REPL)"},
		{"search", "Search the web"},
		{"web", "Fetch web content"},
		{"image", "Generate images with AI enhancement"},
		{"history", "View chat history"},
		{"model", "Model management"},
	}
	for _, c := range commands {
		fmt.Printf("  %s  %s\n",
			helpCommandStyle.Render(fmt.Sprintf("%-10s", c[0])),
			helpDescStyle.Render(c[1]))
	}
	fmt.Println()

	// Flags section
	fmt.Println(helpSectionStyle.Render("Flags"))
	fmt.Println(helpDividerStyle.Render(strings.Repeat("-", 50)))
	flags := [][]string{
		{"-f, --file <path>", "Include file or URL in prompt"},
		{"--search", "Augment with web search results"},
		{"--think", "Enable reasoning mode"},
		{"--json", "Output as JSON"},
		{"-v, --verbose", "Show debug info"},
		{"-h, --help", "Show this help"},
	}
	for _, f := range flags {
		fmt.Printf("  %s  %s\n",
			helpFlagStyle.Render(fmt.Sprintf("%-18s", f[0])),
			helpDescStyle.Render(f[1]))
	}
	fmt.Println()

	// Footer
	fmt.Println(helpDescStyle.Render("Use \"zai <command> --help\" for command details"))
	fmt.Println()
}

func initConfig() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configDir := filepath.Join(home, ".config", "zai")
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return err
			}
		}
	}

	viper.SetEnvPrefix("ZAI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if viper.GetString("api.key") == "" {
		return fmt.Errorf("API key required: set ZAI_API_KEY or configure in ~/.config/zai/config.yaml")
	}

	return nil
}

// buildClientConfig creates ClientConfig from viper settings.
func buildClientConfig() app.ClientConfig {
	return app.ClientConfig{
		APIKey:  viper.GetString("api.key"),
		BaseURL: viper.GetString("api.base_url"),
		Model:   viper.GetString("api.model"),
		Verbose: viper.GetBool("verbose"),
	}
}

// newClient creates a fully configured client with dependencies.
func newClient() *app.Client {
	cfg := buildClientConfig()
	logger := &app.StderrLogger{Verbose: cfg.Verbose}
	history := app.NewFileHistoryStore("")
	return app.NewClient(cfg, logger, history)
}

// hasStdinData detects if stdin has piped/redirected data.
func hasStdinData() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// readStdin reads all data from stdin.
func readStdin() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// runOneShot executes a single prompt and exits.
func runOneShot(prompt string) error {
	client := newClient()
	opts := app.DefaultChatOptions()
	opts.FilePath = viper.GetString("file")
	opts.Think = viper.GetBool("think")

	if viper.GetBool("verbose") {
		fmt.Fprintf(os.Stderr, "Prompt: %s\n", prompt)
		if opts.FilePath != "" {
			fmt.Fprintf(os.Stderr, "File: %s\n", opts.FilePath)
		}
	}

	// Augment prompt with web search results if --search flag is set
	if viper.GetBool("search") {
		if viper.GetBool("verbose") {
			fmt.Fprintf(os.Stderr, "Searching web for: %s\n", prompt)
		}

		searchOpts := app.SearchOptions{
			Count:         5,
			RecencyFilter: "oneWeek",
		}
		results, err := client.SearchWeb(context.Background(), prompt, searchOpts)
		if err != nil {
			if viper.GetBool("verbose") {
				fmt.Fprintf(os.Stderr, "Search failed (continuing without): %v\n", err)
			}
		} else if len(results.SearchResult) > 0 {
			searchContext := app.FormatSearchForContext(results.SearchResult)
			prompt = searchContext + "\n\nUser question: " + prompt

			if viper.GetBool("verbose") {
				fmt.Fprintf(os.Stderr, "Found %d search results\n", len(results.SearchResult))
			}
		}
	}

	response, err := client.Chat(context.Background(), prompt, opts)
	if err != nil {
		return fmt.Errorf("failed to get response: %w", err)
	}

	if viper.GetBool("json") {
		// Create structured JSON output
		output := map[string]interface{}{
			"prompt":      prompt,
			"response":    response,
			"model":       viper.GetString("api.model"),
			"file":        opts.FilePath,
			"think":       opts.Think,
			"search":      viper.GetBool("search"),
			"timestamp":   time.Now().Format(time.RFC3339),
		}

		data, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
	} else {
		fmt.Println(response)
	}

	return nil
}
