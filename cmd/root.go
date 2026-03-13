package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/app"
)

// Constants for input size limits
const (
	MaxStdinSize  = 10 * 1024 * 1024  // 10MB
	MaxFileSize   = 100 * 1024 * 1024 // 100MB
	AudioFileSize = 25 * 1024 * 1024  // 25MB (audio specific)
)

// Flag variables for Cobra binding (required for PersistentFlags).
var (
	cfgFile    string
	verbose    bool
	filePath   string
	think      bool
	jsonOutput bool
	search     bool
	coding     bool
	system     string
)

// RunConfig holds runtime configuration collected from flags and config file.
// Passed to functions instead of accessing globals directly.
type RunConfig struct {
	FilePath   string
	Think      bool
	JSONOutput bool
	Search     bool
	Verbose    bool
	System     string
}

// NewRunConfig creates RunConfig from viper settings (collected after flag parsing).
func NewRunConfig() RunConfig {
	return RunConfig{
		FilePath:   viper.GetString("file"),
		Think:      viper.GetBool("think"),
		JSONOutput: viper.GetBool("json"),
		Search:     viper.GetBool("search"),
		Verbose:    viper.GetBool("verbose"),
		System:     viper.GetString("system"),
	}
}

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
		if cmd.Name() == "history" || cmd.Name() == "completion" || cmd.Name() == "help" || cmd.Name() == "version" {
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

		// Handle --system flag: "-", "/dev/stdin", or file paths
		systemVal := viper.GetString("system")
		stdinUsedForSystem := false
		if systemVal == "-" || systemVal == "/dev/stdin" {
			if stdinData == "" {
				return fmt.Errorf("--system %q requires stdin input", systemVal)
			}
			viper.Set("system", stdinData)
			stdinUsedForSystem = true
		} else if systemVal != "" {
			resolvedSystem, _, err := resolveSystemPrompt(systemVal)
			if err != nil {
				return fmt.Errorf("failed to resolve system prompt: %w", err)
			}
			viper.Set("system", resolvedSystem)
		}

		// If stdin wasn't used for system prompt, prepend it to user prompt as context
		if stdinData != "" && !stdinUsedForSystem {
			var b strings.Builder
			b.WriteString("<stdin>\n")
			b.WriteString(stdinData)
			b.WriteString("\n</stdin>\n\n")
			b.WriteString(prompt)
			prompt = b.String()
		}

		// Build prompt from args
		if len(args) > 0 {
			var b strings.Builder
			if prompt != "" {
				b.WriteString(prompt)
				b.WriteString(" ")
			}
			b.WriteString(strings.Join(args, " "))
			prompt = b.String()
		}

		// Require some input
		if prompt == "" {
			return cmd.Help()
		}

		return runOneShot(prompt)
	},
}

// Execute is the entry point for the CLI. It registers all subcommands and runs the root command.
func Execute() {
	registerRootCmd()
	registerAudioCmd()
	registerChatCmd()
	registerHistoryCmd()
	registerImageCmd()
	registerModelCmd()
	registerReaderCmd()
	registerSearchCmd()
	registerVersionCmd()
	registerVideoCmd()
	registerVisionCmd()

	if err := rootCmd.Execute(); err != nil {
		printStyledError(err)
		os.Exit(1)
	}
}

// printStyledError displays an error with lipgloss styling.
// Detects usage errors and conditionally shows help hint.
func printStyledError(err error) {
	errMsg := err.Error()

	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "%s %s\n",
		theme.ErrorText.Render("Error:"),
		theme.Description.Render(errMsg))

	// Detect usage errors and show help hint
	if isUsageError(errMsg) {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "%s\n",
			theme.HelpText.Render("Run 'zai --help' for usage information"))
	}
	fmt.Fprintln(os.Stderr)
}

// isUsageError detects if an error is a usage/flag error.
// Pattern from charmbracelet/fang for detecting flag parsing errors.
func isUsageError(errMsg string) bool {
	usageErrorPrefixes := []string{
		"unknown flag:",
		"unknown shorthand flag:",
		"flag needs an argument:",
		"invalid argument",
		"accepts",
		"unknown command",
		"required flag",
	}

	for _, prefix := range usageErrorPrefixes {
		if strings.HasPrefix(errMsg, prefix) || strings.Contains(errMsg, prefix) {
			return true
		}
	}
	return false
}

func registerRootCmd() {
	// Enable custom styled error output
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.SetHelpFunc(styledHelp)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.config/zai/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "include file contents in prompt")
	rootCmd.PersistentFlags().BoolVar(&think, "think", false, "enable thinking/reasoning mode")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&search, "search", false, "augment prompt with web search results")
	rootCmd.PersistentFlags().BoolVarP(&coding, "coding", "C", false, "use coding API endpoint")
	rootCmd.PersistentFlags().StringVar(&system, "system", "", "custom system prompt")

	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("file", rootCmd.PersistentFlags().Lookup("file"))
	_ = viper.BindPFlag("think", rootCmd.PersistentFlags().Lookup("think"))
	_ = viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
	_ = viper.BindPFlag("search", rootCmd.PersistentFlags().Lookup("search"))
	_ = viper.BindPFlag("coding", rootCmd.PersistentFlags().Lookup("coding"))
	_ = viper.BindPFlag("system", rootCmd.PersistentFlags().Lookup("system"))
}

// styledHelp displays the custom styled help output.
// For subcommands, delegates to default cobra help to show command-specific usage.
func styledHelp(cmd *cobra.Command, args []string) {
	// If this is a subcommand (not root), use default cobra help
	if cmd != rootCmd {
		rootCmd.SetHelpFunc(nil)        // Temporarily unset to use default
		cmd.Help()                      //nolint:errcheck // help output
		rootCmd.SetHelpFunc(styledHelp) // Restore custom help
		return
	}

	// Title
	fmt.Println()
	fmt.Println(theme.Title.Render(" ZAI ") + " " + theme.Description.Render("Chat with Z.AI models"))
	fmt.Println()

	// Examples section
	fmt.Println(theme.Section.Render("Examples"))
	fmt.Println(theme.Divider.Render(strings.Repeat("-", 50)))
	examples := []string{
		`zai "Explain quantum computing"`,
		`zai -f main.go "Review this code"`,
		`zai --search "Latest AI news"`,
		`echo "text" | zai "summarize"`,
	}
	for _, ex := range examples {
		fmt.Printf("  %s\n", theme.Example.Render(ex))
	}
	fmt.Println()

	// Commands section
	fmt.Println(theme.Section.Render("Commands"))
	fmt.Println(theme.Divider.Render(strings.Repeat("-", 50)))
	commands := [][]string{
		{"chat", "Interactive chat session (REPL)"},
		{"search", "Search the web"},
		{"reader", "Fetch web content"},
		{"image", "Generate images with AI enhancement"},
		{"vision", "Analyze images with AI vision"},
		{"audio", "Transcribe audio to text"},
		{"video", "Generate videos with AI"},
		{"history", "View chat history"},
		{"model", "Model management"},
		{"version", "Show version information"},
	}
	for _, c := range commands {
		fmt.Printf("  %s  %s\n",
			theme.Command.Render(fmt.Sprintf("%-10s", c[0])),
			theme.Description.Render(c[1]))
	}
	fmt.Println()

	// Flags section
	fmt.Println(theme.Section.Render("Flags"))
	fmt.Println(theme.Divider.Render(strings.Repeat("-", 50)))
	flags := [][]string{
		{"-f, --file <path>", "Include file or URL in prompt"},
		{"--search", "Augment with web search results"},
		{"--think", "Enable reasoning mode"},
		{"-C, --coding", "Use coding API endpoint"},
		{"--json", "Output as JSON"},
		{"-v, --verbose", "Show debug info"},
		{"-h, --help", "Show this help"},
	}
	for _, f := range flags {
		fmt.Printf("  %s  %s\n",
			theme.Flag.Render(fmt.Sprintf("%-18s", f[0])),
			theme.Description.Render(f[1]))
	}
	fmt.Println()

	// Footer
	fmt.Println(theme.Description.Render("Use \"zai <command> --help\" for command details"))
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
	}

	// Always read config file (moved outside if/else to fix --config flag bug)
	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return err
		}
	}

	viper.SetEnvPrefix("ZAI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if viper.GetString("api.key") == "" {
		return errors.New("API key required: set ZAI_API_KEY or configure in ~/.config/zai/config.yaml")
	}

	return nil
}

// createContext creates a context with timeout for CLI operations.
// If timeout is 0, returns a cancelable context without timeout.
func createContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(context.Background(), timeout)
	}
	return context.WithCancel(context.Background())
}

// getModelWithDefault returns the model from config key or uses the fallback.
// Simplifies the pattern: if flag empty -> check config -> use default.
func getModelWithDefault(configKey, fallback string) string {
	if model := viper.GetString(configKey); model != "" {
		return model
	}
	return fallback
}

// buildClientConfig creates ClientConfig from viper settings.
func buildClientConfig() app.ClientConfig {
	// Load retry config from viper
	retryCfg := app.RetryConfig{
		MaxAttempts:    viper.GetInt("api.retry.max_attempts"),
		InitialBackoff: viper.GetDuration("api.retry.initial_backoff"),
		MaxBackoff:     viper.GetDuration("api.retry.max_backoff"),
	}

	// Load rate limit config from viper
	rateLimitCfg := app.RateLimitConfig{
		RequestsPerSecond: viper.GetInt("api.rate_limit.requests_per_second"),
		Burst:             viper.GetInt("api.rate_limit.burst"),
	}

	baseURL := viper.GetString("api.base_url")
	codingBaseURL := viper.GetString("api.coding_base_url")

	// Swap to coding API if --coding flag or api.coding_plan config is set
	if viper.GetBool("coding") || viper.GetBool("api.coding_plan") {
		baseURL = codingBaseURL
	}

	return app.ClientConfig{
		APIKey:        viper.GetString("api.key"),
		BaseURL:       baseURL,
		CodingBaseURL: codingBaseURL,
		Model:         viper.GetString("api.model"),
		Verbose:       viper.GetBool("verbose"),
		RateLimit:     rateLimitCfg,
		RetryConfig:   retryCfg,
	}
}

// newClient creates a fully configured client with dependencies.
// Uses default http.Client by passing nil for httpClient.
func newClient() *app.Client {
	cfg := buildClientConfig()
	logger := app.NewLogger(cfg.Verbose)
	history := app.NewFileHistoryStore("")
	return app.NewClient(cfg, logger, history, nil)
}

// newClientWithoutHistory creates a client without history storage.
// Used for commands that don't need history (e.g., web fetch).
func newClientWithoutHistory() *app.Client {
	cfg := buildClientConfig()
	logger := app.NewLogger(cfg.Verbose)
	return app.NewClient(cfg, logger, nil, nil)
}

// newClientWithConfig creates a client with custom config.
// Used when command-specific config overrides are needed.
func newClientWithConfig(cfg app.ClientConfig) *app.Client {
	logger := app.NewLogger(cfg.Verbose)
	history := app.NewFileHistoryStore("")
	return app.NewClient(cfg, logger, history, nil)
}

// hasStdinData detects if stdin has piped/redirected data.
func hasStdinData() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// readStdin reads all data from stdin with a size limit.
func readStdin() (string, error) {
	limitedReader := io.LimitReader(os.Stdin, MaxStdinSize)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", err
	}
	if len(data) == MaxStdinSize {
		return "", fmt.Errorf("stdin exceeds maximum size of %d bytes", MaxStdinSize)
	}
	dataStr := string(data)
	return strings.TrimSpace(dataStr), nil
}

// resolveSystemPrompt returns file contents if value points to an existing file.
// Returns (value, false, nil) when the path doesn't exist so callers can decide fallback behavior.
func resolveSystemPrompt(value string) (string, bool, error) {
	if value == "" {
		return "", false, nil
	}

	cleanPath := filepath.Clean(value)
	if strings.HasPrefix(cleanPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", false, fmt.Errorf("failed to get home directory: %w", err)
		}
		if cleanPath == "~" {
			cleanPath = home
		} else {
			cleanPath = filepath.Join(home, strings.TrimPrefix(cleanPath, "~/"))
		}
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return value, false, nil
		}
		return "", false, err
	}

	if info.IsDir() {
		return "", false, fmt.Errorf("system prompt path is a directory: %s", cleanPath)
	}

	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return "", false, fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), true, nil
}

// looksLikeSystemPromptPath checks whether a value appears to be a filesystem path.
func looksLikeSystemPromptPath(value string) bool {
	if value == "" {
		return false
	}

	if strings.HasPrefix(value, "~") || strings.HasPrefix(value, "/") || strings.HasPrefix(value, "./") || strings.HasPrefix(value, "../") {
		return true
	}

	if strings.Contains(value, string(filepath.Separator)) {
		return true
	}

	ext := filepath.Ext(value)
	return ext != ""
}

// runOneShot executes a single prompt and exits.
func runOneShot(prompt string) error {
	cfg := NewRunConfig()
	client, opts := setupOneShotConfig(cfg)
	logConfigDetails(cfg, opts, prompt)

	ctx, cancel := createContext(5 * time.Minute)
	defer cancel()

	prompt = augmentWithWebSearch(ctx, client, cfg, prompt)
	response, err := callChatAPI(ctx, client, prompt, opts)
	if err != nil {
		return fmt.Errorf("failed to get response: %w", err)
	}

	formatOutput(response, cfg, prompt, opts)

	return nil
}

// setupOneShotConfig initializes configuration and creates client with options
func setupOneShotConfig(cfg RunConfig) (*app.Client, app.ChatOptions) {
	client := newClient()
	opts := app.DefaultChatOptions()
	opts.FilePath = cfg.FilePath
	opts.Think = cfg.Think
	opts.SystemPrompt = cfg.System
	return client, opts
}

// logConfigDetails logs configuration details if verbose mode is enabled
func logConfigDetails(cfg RunConfig, opts app.ChatOptions, prompt string) {
	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Prompt: %s\n", prompt)
		if opts.FilePath != "" {
			fmt.Fprintf(os.Stderr, "File: %s\n", opts.FilePath)
		}
		if opts.SystemPrompt != "" {
			fmt.Fprintf(os.Stderr, "System prompt: %s\n", opts.SystemPrompt)
		}
	}
}

// augmentWithWebSearch augments the prompt with web search results if --search flag is set
func augmentWithWebSearch(ctx context.Context, client *app.Client, cfg RunConfig, prompt string) string {
	if !cfg.Search {
		return prompt
	}

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Searching web for: %s\n", prompt)
	}

	searchOpts := app.SearchOptions{
		Count:         5,
		RecencyFilter: "oneWeek",
		Language:      viper.GetString("web_search.language"),
	}
	results, err := client.SearchWeb(ctx, prompt, searchOpts)
	if err != nil {
		if cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Search failed (continuing without): %v\n", err)
		}
		return prompt
	}

	if len(results.SearchResult) > 0 {
		searchContext := app.FormatSearchForContext(results.SearchResult)
		var b strings.Builder
		b.WriteString(searchContext)
		b.WriteString("\n\nUser question: ")
		b.WriteString(prompt)
		augmentedPrompt := b.String()

		if cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Found %d search results\n", len(results.SearchResult))
		}

		return augmentedPrompt
	}

	return prompt
}

// callChatAPI makes the chat API call and returns the response
func callChatAPI(ctx context.Context, client *app.Client, prompt string, opts app.ChatOptions) (string, error) {
	return client.Chat(ctx, prompt, opts)
}

// formatOutput formats and prints the response according to configuration
func formatOutput(response string, cfg RunConfig, prompt string, opts app.ChatOptions) {
	if cfg.JSONOutput {
		output := map[string]interface{}{
			"prompt":    prompt,
			"response":  response,
			"model":     viper.GetString("api.model"),
			"file":      opts.FilePath,
			"think":     opts.Think,
			"search":    cfg.Search,
			"timestamp": time.Now().Format(time.RFC3339),
		}

		data, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to marshal JSON: %v\n", err)
			return
		}
		fmt.Println(string(data))
	} else {
		fmt.Println(response)
	}
}
