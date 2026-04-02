package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/app"
	"github.com/dotcommander/zai/internal/config"
)

var (
	searchCount          int
	searchRecency        string
	searchDomain         string
	searchLang           string
	searchFormat         string
	searchEngine         string
	searchContentSize    string
	searchShowSources    bool
	searchForceSearch    bool
	searchCustomPrompt   string
	searchResultSequence string
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

func registerSearchCmd() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().IntVarP(&searchCount, "count", "c", 0, "Number of results (1-50)")
	searchCmd.Flags().StringVarP(&searchRecency, "recency", "r", "", "Time filter: oneDay, oneWeek, oneMonth, oneYear, noLimit")
	searchCmd.Flags().StringVarP(&searchDomain, "domain", "d", "", "Limit to specific domain")
	searchCmd.Flags().StringVarP(&searchLang, "lang", "l", "", "Search language (e.g., en, zh). Default: config value")
	searchCmd.Flags().StringVarP(&searchFormat, "format", "o", "table", "Output format: table, detailed, json")
	searchCmd.Flags().StringVar(&searchEngine, "engine", "", "Search engine: search_std, search_pro, search_pro_sogou, search_pro_quark")
	searchCmd.Flags().StringVar(&searchContentSize, "content-size", "", "Content detail: medium (400-600 chars) or high (2500 chars)")
	searchCmd.Flags().BoolVar(&searchShowSources, "sources", false, "Include detailed source information")
	searchCmd.Flags().BoolVar(&searchForceSearch, "require-search", false, "Force response based on search results")
	searchCmd.Flags().StringVar(&searchCustomPrompt, "search-prompt", "", "Custom prompt for processing results (supports {{current_date}})")
	searchCmd.Flags().StringVar(&searchResultSequence, "result-sequence", "", "Result ordering: before or after (default: after)")
}

func runSearch(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if !cfg.WebSearch.Enabled {
		return errors.New("web search is disabled in configuration")
	}

	query, err := parseSearchQuery(args)
	if err != nil {
		return err
	}

	opts, err := buildSearchOptions(cfg)
	if err != nil {
		return err
	}

	client := newClientWithConfig(app.ClientConfig{
		APIKey:  cfg.API.Key,
		BaseURL: cfg.API.BaseURL,
		Model:   cfg.API.Model,
		Timeout: time.Duration(cfg.WebSearch.Timeout) * time.Second,
		Verbose: viper.GetBool("verbose"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.WebSearch.Timeout)*time.Second)
	defer cancel()

	start := time.Now()
	resp, err := client.SearchWeb(ctx, query, opts)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	format := searchFormat
	if viper.GetBool("json") {
		format = "json"
	}

	output, err := formatSearchOutput(resp.SearchResult, format, query, time.Since(start), viper.GetBool("verbose"))
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	fmt.Print(output)
	return nil
}

// parseSearchQuery extracts query from args or stdin.
func parseSearchQuery(args []string) (string, error) {
	switch {
	case len(args) > 0:
		return args[0], nil
	case hasStdinData():
		query, err := readStdin()
		if err != nil {
			return "", fmt.Errorf("failed to read from stdin: %w", err)
		}
		if query == "" {
			return "", errors.New("empty query from stdin")
		}
		return query, nil
	default:
		return "", errors.New("search query is required")
	}
}

var validSearchFormats = map[string]bool{"table": true, "detailed": true, "json": true}

// buildSearchOptions creates SearchOptions from flags with config defaults.
func buildSearchOptions(cfg *config.Config) (app.SearchOptions, error) {
	if !validSearchFormats[searchFormat] {
		return app.SearchOptions{}, fmt.Errorf("invalid format: %s (must be table, detailed, or json)", searchFormat)
	}

	opts := app.SearchOptions{
		Count:          searchCount,
		DomainFilter:   searchDomain,
		RecencyFilter:  searchRecency,
		Language:       searchLang,
		SearchEngine:   searchEngine,
		ContentSize:    searchContentSize,
		SearchResult:   searchShowSources,
		RequireSearch:  searchForceSearch,
		SearchPrompt:   searchCustomPrompt,
		ResultSequence: searchResultSequence,
	}

	if opts.Count == 0 {
		opts.Count = cfg.WebSearch.DefaultCount
	}
	if opts.RecencyFilter == "" {
		opts.RecencyFilter = cfg.WebSearch.DefaultRecency
	}
	if opts.Language == "" {
		opts.Language = cfg.WebSearch.Language
	}
	if opts.SearchEngine == "" {
		opts.SearchEngine = cfg.WebSearch.SearchEngine
	}
	if opts.ContentSize == "" {
		opts.ContentSize = cfg.WebSearch.ContentSize
	}

	return opts, nil
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
		fmt.Fprintf(&sb, "🔍 Search results for: %s\n", query)
		fmt.Fprintf(&sb, "⏱️  Duration: %v\n", duration)
		fmt.Fprintf(&sb, "📊 Results: %d\n\n", len(results))
	}

	if len(results) == 0 {
		sb.WriteString("No results found.\n")
		return sb.String(), nil
	}

	// Pre-compute domains once to avoid redundant URL parsing
	domains := make([]string, len(results))
	maxTitleLen := 0
	maxDomainLen := 0
	for i, result := range results {
		domains[i] = extractDomain(result.Link)
		if len(result.Title) > maxTitleLen {
			maxTitleLen = len(result.Title)
		}
		if len(domains[i]) > maxDomainLen {
			maxDomainLen = len(domains[i])
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
	fmt.Fprintf(&sb, "%-*s  %-*s  %s\n", maxTitleLen, "Title", maxDomainLen, "Domain", "URL")
	sb.WriteString(strings.Repeat("-", maxTitleLen+maxDomainLen+50) + "\n")

	// Table rows
	for i, result := range results {
		title := result.Title
		if len(title) > maxTitleLen {
			title = title[:maxTitleLen-3] + "..."
		}

		domain := domains[i]
		if len(domain) > maxDomainLen {
			domain = domain[:maxDomainLen-3] + "..."
		}

		fmt.Fprintf(&sb, "%-*s  %-*s  %s\n", maxTitleLen, title, maxDomainLen, domain, result.Link)

		// Add content preview for first few results in verbose mode
		if verbose && i < 3 {
			content := result.Content
			if len(content) > 100 {
				content = content[:100] + "..."
			}
			fmt.Fprintf(&sb, "   %s\n\n", content)
		}
	}

	return sb.String(), nil
}

// formatSearchDetailed formats results with full details
func formatSearchDetailed(results []app.SearchResult, query string, duration time.Duration) (string, error) {
	var sb strings.Builder

	// Header
	fmt.Fprintf(&sb, "🔍 Search results for: %s\n", query)
	fmt.Fprintf(&sb, "⏱️  Duration: %v\n", duration)
	fmt.Fprintf(&sb, "📊 Results: %d\n\n", len(results))

	if len(results) == 0 {
		sb.WriteString("No results found.\n")
		return sb.String(), nil
	}

	// Detailed results
	for i, result := range results {
		fmt.Fprintf(&sb, "%d. %s\n", i+1, result.Title)
		fmt.Fprintf(&sb, "   URL: %s\n", result.Link)
		if result.Media != "" {
			fmt.Fprintf(&sb, "   Media: %s\n", result.Media)
		}
		if result.PublishDate != "" {
			fmt.Fprintf(&sb, "   Published: %s\n", result.PublishDate)
		}
		sb.WriteString("\n")

		// Content
		content := strings.ReplaceAll(result.Content, "\n", " ")
		if len(content) > 300 {
			content = content[:300] + "..."
		}
		sb.WriteString("   ")
		sb.WriteString(content)
		sb.WriteString("\n\n")
		sb.WriteString(strings.Repeat("-", 80))
		sb.WriteString("\n\n")
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

// extractDomain extracts domain from URL using net/url stdlib.
// Handles edge cases like ports, IPv6, and malformed URLs.
func extractDomain(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL // Return as-is if parsing fails
	}

	host := u.Host

	// Handle empty host (relative URLs, etc.)
	if host == "" {
		return rawURL
	}

	// Remove port if present (handles both IPv4:port and [IPv6]:port)
	if colonIdx := strings.LastIndex(host, ":"); colonIdx != -1 {
		// Check if it's not an IPv6 address without brackets
		if !strings.Contains(host, "[") || strings.Contains(host, "]:") {
			host = host[:colonIdx]
		}
	}

	// Remove brackets from IPv6
	host = strings.TrimPrefix(host, "[")
	host = strings.TrimSuffix(host, "]")

	// Remove www prefix
	host = strings.TrimPrefix(host, "www.")

	return host
}
