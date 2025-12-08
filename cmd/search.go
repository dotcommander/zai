package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"zai/internal/app"
	"zai/internal/config"
)

var (
	searchCount     int
	searchRecency   string
	searchDomain    string
	searchFormat    string
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search the web using Z.AI search engine",
	Long: `Search the web and return results optimized for LLM consumption.

The search query can be provided as an argument or piped via stdin.

Examples:
  zai search "golang best practices"
  echo "golang best practices" | zai search
  zai search "latest AI news" -c 5 -r oneWeek
  zai search "site:github.com golang" -d github.com`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().IntVarP(&searchCount, "count", "c", 0, "Number of results (1-50)")
	searchCmd.Flags().StringVarP(&searchRecency, "recency", "r", "", "Time filter: oneDay, oneWeek, oneMonth, oneYear, noLimit")
	searchCmd.Flags().StringVarP(&searchDomain, "domain", "d", "", "Limit to specific domain")
	searchCmd.Flags().StringVarP(&searchFormat, "format", "o", "table", "Output format: table, detailed, json")
}

func runSearch(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if web search is enabled
	if !cfg.WebSearch.Enabled {
		return fmt.Errorf("web search is disabled in configuration")
	}

	// Get query from args or stdin
	var query string
	if len(args) > 0 {
		query = args[0]
	} else {
		// Read from stdin if piped
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Stdin is piped
			data, err := os.ReadFile(os.Stdin.Name())
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			query = strings.TrimSpace(string(data))
			if query == "" {
				return fmt.Errorf("empty query from stdin")
			}
		} else {
			return fmt.Errorf("search query is required")
		}
	}

	// Validate format
	validFormats := map[string]bool{
		"table": true, "detailed": true, "json": true,
	}
	if !validFormats[searchFormat] {
		return fmt.Errorf("invalid format: %s (must be table, detailed, or json)", searchFormat)
	}

	// Prepare search options
	opts := app.SearchOptions{
		Count:         searchCount,
		DomainFilter:  searchDomain,
		RecencyFilter: searchRecency,
	}

	// Use defaults if not specified
	if opts.Count == 0 {
		opts.Count = cfg.WebSearch.DefaultCount
	}
	if opts.RecencyFilter == "" {
		opts.RecencyFilter = cfg.WebSearch.DefaultRecency
	}

	// Create client
	logger := &app.StderrLogger{Verbose: viper.GetBool("verbose")}
	history := app.NewFileHistoryStore("")
	client := app.NewClient(app.ClientConfig{
		APIKey:  cfg.API.Key,
		BaseURL: cfg.API.BaseURL,
		Model:   cfg.API.Model,
		Timeout: time.Duration(cfg.WebSearch.Timeout) * time.Second,
		Verbose: viper.GetBool("verbose"),
	}, logger, history)

	// Set context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.WebSearch.Timeout)*time.Second)
	defer cancel()

	// Perform search
	start := time.Now()
	resp, err := client.SearchWeb(ctx, query, opts)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	duration := time.Since(start)

	// Format and display results
	output, err := formatSearchOutput(resp.SearchResult, searchFormat, query, duration, viper.GetBool("verbose"))
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	fmt.Print(output)

	return nil
}

// formatSearchOutput formats search results according to the specified format
func formatSearchOutput(results []app.SearchResult, format, query string, duration time.Duration, verbose bool) (string, error) {
	switch format {
	case "json":
		return formatSearchJSON(results, query, duration)
	case "detailed":
		return formatSearchDetailed(results, query, duration)
	default: // table
		return formatSearchTable(results, query, duration, verbose)
	}
}

// formatSearchTable formats results as a table
func formatSearchTable(results []app.SearchResult, query string, duration time.Duration, verbose bool) (string, error) {
	var sb strings.Builder

	// Header
	if verbose {
		sb.WriteString(fmt.Sprintf("🔍 Search results for: %s\n", query))
		sb.WriteString(fmt.Sprintf("⏱️  Duration: %v\n", duration))
		sb.WriteString(fmt.Sprintf("📊 Results: %d\n\n", len(results)))
	}

	if len(results) == 0 {
		sb.WriteString("No results found.\n")
		return sb.String(), nil
	}

	// Find max widths for alignment
	maxTitleLen := 0
	maxDomainLen := 0
	for _, result := range results {
		if len(result.Title) > maxTitleLen {
			maxTitleLen = len(result.Title)
		}
		domain := extractDomain(result.Link)
		if len(domain) > maxDomainLen {
			maxDomainLen = len(domain)
		}
	}

	// Limit max width for readability
	if maxTitleLen > 60 {
		maxTitleLen = 60
	}
	if maxDomainLen > 20 {
		maxDomainLen = 20
	}

	// Table header
	sb.WriteString(fmt.Sprintf("%-*s  %-*s  %s\n", maxTitleLen, "Title", maxDomainLen, "Domain", "URL"))
	sb.WriteString(strings.Repeat("-", maxTitleLen+maxDomainLen+50) + "\n")

	// Table rows
	for i, result := range results {
		title := result.Title
		if len(title) > maxTitleLen {
			title = title[:maxTitleLen-3] + "..."
		}

		domain := extractDomain(result.Link)
		if len(domain) > maxDomainLen {
			domain = domain[:maxDomainLen-3] + "..."
		}

		sb.WriteString(fmt.Sprintf("%-*s  %-*s  %s\n", maxTitleLen, title, maxDomainLen, domain, result.Link))

		// Add content preview for first few results in verbose mode
		if verbose && i < 3 {
			content := result.Content
			if len(content) > 100 {
				content = content[:100] + "..."
			}
			sb.WriteString(fmt.Sprintf("   %s\n\n", content))
		}
	}

	return sb.String(), nil
}

// formatSearchDetailed formats results with full details
func formatSearchDetailed(results []app.SearchResult, query string, duration time.Duration) (string, error) {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("🔍 Search results for: %s\n", query))
	sb.WriteString(fmt.Sprintf("⏱️  Duration: %v\n", duration))
	sb.WriteString(fmt.Sprintf("📊 Results: %d\n\n", len(results)))

	if len(results) == 0 {
		sb.WriteString("No results found.\n")
		return sb.String(), nil
	}

	// Detailed results
	for i, result := range results {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, result.Title))
		sb.WriteString(fmt.Sprintf("   URL: %s\n", result.Link))
		if result.Media != "" {
			sb.WriteString(fmt.Sprintf("   Media: %s\n", result.Media))
		}
		if result.PublishDate != "" {
			sb.WriteString(fmt.Sprintf("   Published: %s\n", result.PublishDate))
		}
		sb.WriteString("\n")

		// Content
		content := strings.ReplaceAll(result.Content, "\n", " ")
		if len(content) > 300 {
			content = content[:300] + "..."
		}
		sb.WriteString(fmt.Sprintf("   %s\n\n", content))
		sb.WriteString(strings.Repeat("-", 80) + "\n\n")
	}

	return sb.String(), nil
}

// formatSearchJSON formats results as JSON
func formatSearchJSON(results []app.SearchResult, query string, duration time.Duration) (string, error) {
	// Create a structured output
	output := map[string]interface{}{
		"query":     query,
		"duration":  duration.String(),
		"count":     len(results),
		"results":   results,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Convert to JSON
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// extractDomain extracts domain from URL
func extractDomain(url string) string {
	if strings.HasPrefix(url, "http://") {
		url = url[7:]
	} else if strings.HasPrefix(url, "https://") {
		url = url[8:]
	}

	// Remove path
	if idx := strings.Index(url, "/"); idx != -1 {
		url = url[:idx]
	}

	// Remove www prefix
	if strings.HasPrefix(url, "www.") {
		url = url[4:]
	}

	return url
}