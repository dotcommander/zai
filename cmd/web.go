package cmd

import (
	"context"
	"encoding/json"
	"fmt"
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
	readerFormat          string
	readerTimeout         int
	readerNoCache         bool
	readerNoGFM           bool
	readerKeepImgDataURL  bool
	readerWithImagesSum   bool
	readerWithLinksSum    bool
	readerNoRetainImages  bool
	readerJSON            bool
)

func runReader(cmd *cobra.Command, args []string) error {
	var ctx context.Context
	ctx, cancel := createContext(2 * time.Minute)
	defer cancel()

	url := args[0]

	// Validate API key
	if viper.GetString("api.key") == "" {
		return fmt.Errorf("API key required: set ZAI_API_KEY or configure in ~/.config/zai/config.yaml")
	}

	// Create client using factory with custom timeout (no history needed)
	clientConfig := app.ClientConfig{
		APIKey:  viper.GetString("api.key"),
		BaseURL: viper.GetString("api.base_url"),
		Model:   viper.GetString("api.model"),
		Verbose: viper.GetBool("verbose"),
		Timeout: time.Duration(readerTimeout) * time.Second,
	}
	logger := &app.StderrLogger{Verbose: clientConfig.Verbose}
	client := app.NewClient(clientConfig, logger, nil, nil)

	// Build web reader options
	opts := &app.WebReaderOptions{
		ReturnFormat:      readerFormat,
		Timeout:           &readerTimeout,
		NoCache:           &readerNoCache,
		NoGFM:             &readerNoGFM,
		KeepImgDataURL:    &readerKeepImgDataURL,
		WithImagesSummary: &readerWithImagesSum,
		WithLinksSummary:  &readerWithLinksSum,
	}

	// Set retain images (default true)
	retainImages := !readerNoRetainImages
	opts.RetainImages = &retainImages

	// Validate format
	if readerFormat != "markdown" && readerFormat != "text" {
		return fmt.Errorf("invalid format: %s (must be 'markdown' or 'text')", readerFormat)
	}

	// Validate timeout
	if readerTimeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	// Fetch web content
	resp, err := client.FetchWebContent(ctx, url, opts)
	if err != nil {
		return fmt.Errorf("failed to fetch web content: %w", err)
	}

	// Output results
	if readerJSON {
		// Create structured JSON output
		output := map[string]interface{}{
			"url":               resp.ReaderResult.URL,
			"title":             resp.ReaderResult.Title,
			"description":       resp.ReaderResult.Description,
			"content":           resp.ReaderResult.Content,
			"metadata":          resp.ReaderResult.Metadata,
			"external_resources": resp.ReaderResult.ExternalResources,
			"timestamp":         time.Now().Format(time.RFC3339),
		}

		data, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
	} else {
		// Display human-readable results
		fmt.Printf("Title: %s\n", resp.ReaderResult.Title)
		fmt.Printf("URL: %s\n", resp.ReaderResult.URL)
		if resp.ReaderResult.Description != "" {
			fmt.Printf("Description: %s\n", resp.ReaderResult.Description)
		}
		fmt.Printf("\nContent:\n%s\n", resp.ReaderResult.Content)

		// Display metadata if available
		if len(resp.ReaderResult.Metadata) > 0 {
			fmt.Printf("\nMetadata:\n")
			for k, v := range resp.ReaderResult.Metadata {
				fmt.Printf("  %s: %v\n", k, v)
			}
		}

		// Display external resources if available
		if len(resp.ReaderResult.ExternalResources) > 0 {
			fmt.Printf("\nExternal Resources:\n")
			for k, v := range resp.ReaderResult.ExternalResources {
				fmt.Printf("  %s: %v\n", k, v)
			}
		}
	}

	// Save to history (using default location)
	history := app.NewFileHistoryStore("")

	// Create a history entry for web content
	entry := app.NewWebHistoryEntry(
		resp.ID,
		fmt.Sprintf("Fetch web content: %s", url),
		resp,
		[]string{url},
	)
	if err := history.Save(entry); err != nil {
		logger.Warn("Failed to save to history: %v", err)
	}

	return nil
}

func init() {
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