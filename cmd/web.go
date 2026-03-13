package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/app"
)

// readerCmd represents the reader command
var readerCmd = &cobra.Command{
	Use:   "reader <url>",
	Short: "Fetch and display web content",
	Long: `Fetch and display web content from a URL using Z.AI's web reader API.

Examples:
  zai reader https://example.com
  zai reader https://example.com --format text
  zai reader https://example.com --no-cache
  zai reader https://example.com --timeout 30
  zai reader https://example.com --with-links-summary`,
	Args: cobra.ExactArgs(1),
	RunE: runReader,
}

var (
	readerFormat         string
	readerTimeout        int
	readerNoCache        bool
	readerNoGFM          bool
	readerKeepImgDataURL bool
	readerWithImagesSum  bool
	readerWithLinksSum   bool
	readerNoRetainImages bool
	readerJSON           bool
)

func runReader(cmd *cobra.Command, args []string) error {
	ctx, cancel := createContext(2 * time.Minute)
	defer cancel()

	url := args[0]

	if err := validateReaderFlags(); err != nil {
		return err
	}

	client, logger := buildReaderClient()

	opts := buildReaderOpts()

	resp, err := client.FetchWebContent(ctx, url, opts)
	if err != nil {
		return fmt.Errorf("failed to fetch web content: %w", err)
	}

	if err := outputReaderResult(resp); err != nil {
		return err
	}

	saveReaderHistory(logger, url, resp)
	return nil
}

// validateReaderFlags validates the reader command flags.
func validateReaderFlags() error {
	if readerFormat != "markdown" && readerFormat != "text" {
		return fmt.Errorf("invalid format: %s (must be 'markdown' or 'text')", readerFormat)
	}
	if readerTimeout <= 0 {
		return errors.New("timeout must be positive")
	}
	return nil
}

// buildReaderClient creates the HTTP client and logger for the reader command.
func buildReaderClient() (*app.Client, *slog.Logger) {
	clientConfig := app.ClientConfig{
		APIKey:  viper.GetString("api.key"),
		BaseURL: viper.GetString("api.base_url"),
		Model:   viper.GetString("api.model"),
		Verbose: viper.GetBool("verbose"),
		Timeout: time.Duration(readerTimeout) * time.Second,
	}
	logger := app.NewLogger(clientConfig.Verbose)
	return app.NewClient(clientConfig, logger, nil, nil), logger
}

// buildReaderOpts creates the web reader options from command flags.
func buildReaderOpts() *app.WebReaderOptions {
	retainImages := !readerNoRetainImages
	return &app.WebReaderOptions{
		ReturnFormat:      readerFormat,
		Timeout:           &readerTimeout,
		NoCache:           &readerNoCache,
		NoGFM:             &readerNoGFM,
		KeepImgDataURL:    &readerKeepImgDataURL,
		WithImagesSummary: &readerWithImagesSum,
		WithLinksSummary:  &readerWithLinksSum,
		RetainImages:      &retainImages,
	}
}

// outputReaderResult prints the web reader response to stdout.
func outputReaderResult(resp *app.WebReaderResponse) error {
	if readerJSON {
		return outputReaderJSON(resp)
	}
	outputReaderText(resp)
	return nil
}

// outputReaderJSON prints the response as JSON.
func outputReaderJSON(resp *app.WebReaderResponse) error {
	output := map[string]interface{}{
		"url":                resp.ReaderResult.URL,
		"title":              resp.ReaderResult.Title,
		"description":        resp.ReaderResult.Description,
		"content":            resp.ReaderResult.Content,
		"metadata":           resp.ReaderResult.Metadata,
		"external_resources": resp.ReaderResult.ExternalResources,
		"timestamp":          time.Now().Format(time.RFC3339),
	}
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// outputReaderText prints the response in human-readable format.
func outputReaderText(resp *app.WebReaderResponse) {
	fmt.Printf("Title: %s\n", resp.ReaderResult.Title)
	fmt.Printf("URL: %s\n", resp.ReaderResult.URL)
	if resp.ReaderResult.Description != "" {
		fmt.Printf("Description: %s\n", resp.ReaderResult.Description)
	}
	fmt.Printf("\nContent:\n%s\n", resp.ReaderResult.Content)

	if len(resp.ReaderResult.Metadata) > 0 {
		fmt.Printf("\nMetadata:\n")
		for k, v := range resp.ReaderResult.Metadata {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}

	if len(resp.ReaderResult.ExternalResources) > 0 {
		fmt.Printf("\nExternal Resources:\n")
		for k, v := range resp.ReaderResult.ExternalResources {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}
}

// saveReaderHistory saves the web reader result to history.
func saveReaderHistory(logger *slog.Logger, url string, resp *app.WebReaderResponse) {
	history := app.NewFileHistoryStore("")
	entry := app.NewWebHistoryEntry(
		resp.ID,
		fmt.Sprintf("Fetch web content: %s", url),
		resp,
		[]string{url},
	)
	if err := history.Save(entry); err != nil {
		logger.Warn("failed to save to history", "error", err)
	}
}

func registerReaderCmd() {
	rootCmd.AddCommand(readerCmd)

	// Web reader flags
	readerCmd.Flags().StringVar(&readerFormat, "format", "markdown", "Return format (markdown or text)")
	readerCmd.Flags().IntVar(&readerTimeout, "timeout", 20, "Request timeout in seconds")
	readerCmd.Flags().BoolVar(&readerNoCache, "no-cache", false, "Disable caching")
	readerCmd.Flags().BoolVar(&readerNoGFM, "no-gfm", false, "Disable GitHub Flavored Markdown")
	readerCmd.Flags().BoolVar(&readerKeepImgDataURL, "keep-img-data-url", false, "Keep image data URLs")
	readerCmd.Flags().BoolVar(&readerWithImagesSum, "with-images-summary", false, "Include image summary")
	readerCmd.Flags().BoolVar(&readerWithLinksSum, "with-links-summary", false, "Include links summary")
	readerCmd.Flags().BoolVar(&readerNoRetainImages, "no-retain-images", false, "Do not retain images")
	readerCmd.Flags().BoolVar(&readerJSON, "json", false, "Output in JSON format")
}
