package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"zai/internal/app"
)

var (
	cfgFile  string
	verbose  bool
	filePath string
)

var rootCmd = &cobra.Command{
	Use:   "zai [prompt]",
	Short: "Chat with Z.AI GLM model",
	Long: `ZAI is a CLI tool for chatting with Z.AI models.

One-shot mode:
  zai "Explain quantum computing"
  zai -f main.go "Explain this code"

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
		if len(args) == 0 {
			return cmd.Help()
		}
		// One-shot mode: zai "prompt"
		prompt := strings.Join(args, " ")
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.config/zai/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "include file contents in prompt")

	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("file", rootCmd.PersistentFlags().Lookup("file"))
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

// runOneShot executes a single prompt and exits.
func runOneShot(prompt string) error {
	client := newClient()
	opts := app.DefaultChatOptions()
	opts.FilePath = viper.GetString("file")

	if viper.GetBool("verbose") {
		fmt.Fprintf(os.Stderr, "🤖 Prompt: %s\n", prompt)
		if opts.FilePath != "" {
			fmt.Fprintf(os.Stderr, "📄 File: %s\n", opts.FilePath)
		}
	}

	response, err := client.Chat(context.Background(), prompt, opts)
	if err != nil {
		return fmt.Errorf("failed to get response: %w", err)
	}

	fmt.Println(response)
	return nil
}
