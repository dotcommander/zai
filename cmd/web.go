package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"zai/internal/app"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web <url>",
	Short: "Fetch and display web content",
	Long: `Fetch and display web content from a URL using Z.AI's web reader API.

Examples:
  zai web https://example.com
  zai web https://example.com --format text
  zai web https://example.com --no-cache
  zai web https://example.com --timeout 30
  zai web https://example.com --with-links-summary`,
	Args: cobra.ExactArgs(1),
	RunE: runWeb,
}

var (
	webFormat          string
	webTimeout         int
	webNoCache         bool
	webNoGFM           bool
	webKeepImgDataURL  bool
	webWithImagesSum   bool
	webWithLinksSum    bool
	webNoRetainImages  bool
	webJSON            bool
)

func runWeb(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	url := args[0]

	// Validate API key
	if viper.GetString("api.key") == "" {
		return fmt.Errorf("API key required: set ZAI_API_KEY or configure in ~/.config/zai/config.yaml")
	}

	// Create logger
	logger := &app.StderrLogger{Verbose: viper.GetBool("verbose")}

	// Create client
	clientConfig := app.ClientConfig{
		APIKey:  viper.GetString("api.key"),
		BaseURL: viper.GetString("api.base_url"),
		Model:   viper.GetString("api.model"),
		Verbose: viper.GetBool("verbose"),
		Timeout: time.Duration(webTimeout) * time.Second,
	}

	client := app.NewClient(clientConfig, logger, nil)

	// Build web reader options
	opts := &app.WebReaderOptions{
		ReturnFormat:      webFormat,
		Timeout:           &webTimeout,
		NoCache:           &webNoCache,
		NoGFM:             &webNoGFM,
		KeepImgDataURL:    &webKeepImgDataURL,
		WithImagesSummary: &webWithImagesSum,
		WithLinksSummary:  &webWithLinksSum,
	}

	// Set retain images (default true)
	retainImages := !webNoRetainImages
	opts.RetainImages = &retainImages

	// Validate format
	if webFormat != "markdown" && webFormat != "text" {
		return fmt.Errorf("invalid format: %s (must be 'markdown' or 'text')", webFormat)
	}

	// Validate timeout
	if webTimeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	// Fetch web content
	resp, err := client.FetchWebContent(ctx, url, opts)
	if err != nil {
		return fmt.Errorf("failed to fetch web content: %w", err)
	}

	// Output results
	if webJSON {
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
	rootCmd.AddCommand(webCmd)

	// Web reader flags
	webCmd.Flags().StringVar(&webFormat, "format", "markdown", "Return format (markdown or text)")
	webCmd.Flags().IntVar(&webTimeout, "timeout", 20, "Request timeout in seconds")
	webCmd.Flags().BoolVar(&webNoCache, "no-cache", false, "Disable caching")
	webCmd.Flags().BoolVar(&webNoGFM, "no-gfm", false, "Disable GitHub Flavored Markdown")
	webCmd.Flags().BoolVar(&webKeepImgDataURL, "keep-img-data-url", false, "Keep image data URLs")
	webCmd.Flags().BoolVar(&webWithImagesSum, "with-images-summary", false, "Include image summary")
	webCmd.Flags().BoolVar(&webWithLinksSum, "with-links-summary", false, "Include links summary")
	webCmd.Flags().BoolVar(&webNoRetainImages, "no-retain-images", false, "Do not retain images")
	webCmd.Flags().BoolVar(&webJSON, "json", false, "Output in JSON format")
}