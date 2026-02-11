package app

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// urlRegex matches HTTP/HTTPS URLs
var urlRegex = regexp.MustCompile(`https?://[^\s<>"]+|www\.[^\s<>"]+`)

// ExtractURLs extracts all URLs from the given text.
func ExtractURLs(text string) []string {
	matches := urlRegex.FindAllString(text, -1)
	if len(matches) == 0 {
		return nil
	}

	// Deduplicate and normalize URLs with pre-allocated capacity
	seen := make(map[string]bool, len(matches))
	urls := make([]string, 0, len(matches))

	for _, match := range matches {
		// Normalize URL
		normalized := normalizeURL(match)
		if normalized == "" {
			continue
		}

		// Skip duplicates
		if seen[normalized] {
			continue
		}
		seen[normalized] = true

		urls = append(urls, normalized)
	}

	return urls
}

// normalizeURL ensures URL has proper scheme and is valid.
func normalizeURL(raw string) string {
	// Trim trailing punctuation
	raw = strings.TrimRight(raw, ".,!?;:)]}")

	// Add scheme if missing
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		if strings.HasPrefix(raw, "www.") {
			raw = "https://" + raw
		} else {
			return "" // Not a valid URL pattern
		}
	}

	// Validate URL structure
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	// Only allow http/https schemes
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ""
	}

	// Ensure host exists
	if parsed.Host == "" {
		return ""
	}

	// Return the cleaned URL
	return parsed.String()
}

// IsWebContentPrompt checks if the prompt is asking for web content.
func IsWebContentPrompt(text string) bool {
	textLower := strings.ToLower(text)

	// Keywords that indicate web content request
	webKeywords := []string{
		"fetch", "read", "summarize", "analyze", "extract", "get content",
		"what's on", "check out", "look at", "visit", "browse",
	}

	// Check if any web keywords are present
	for _, keyword := range webKeywords {
		if strings.Contains(textLower, keyword) {
			// Check if there's also a URL
			if len(ExtractURLs(text)) > 0 {
				return true
			}
		}
	}

	// If URLs are present, assume web content is desired
	return len(ExtractURLs(text)) > 0
}

// FormatWebContent formats web content for inclusion in chat prompts.
func FormatWebContent(url, title, content string) string {
	const webContentFormat = `<web_content>
<source_url>%s</source_url>
<title>%s</title>
<content>
%s
</content>
</web_content>`

	return fmt.Sprintf(webContentFormat, url, title, content)
}

// IsValidWebURL checks if a string is a valid web URL.
func IsValidWebURL(raw string) bool {
	normalized := normalizeURL(raw)
	return normalized != ""
}

// IsSearchPrompt checks if the prompt is requesting a web search.
func IsSearchPrompt(text string) bool {
	textLower := strings.ToLower(text)

	// Exclude if it contains a URL (likely a web fetch request, not search)
	if len(ExtractURLs(text)) > 0 {
		return false
	}

	// Keywords that indicate search request
	searchKeywords := []string{
		"search for", "find information about", "look up", "search",
		"what is", "who is", "when was", "where is", "why is", "how to",
		"latest news about", "recent", "current", "define", "explain",
	}

	// Check if any search keywords are present
	for _, keyword := range searchKeywords {
		if strings.Contains(textLower, keyword) {
			return true
		}
	}

	// Check for question patterns
	questionPatterns := []string{
		"?", "what are", "who are", "when were", "where are",
		"why are", "how do", "how can", "how would",
	}

	for _, pattern := range questionPatterns {
		if strings.Contains(textLower, pattern) {
			return true
		}
	}

	return false
}

// ExtractSearchQuery extracts the search query from a prompt.
func ExtractSearchQuery(text string) string {
	// Remove common prefixes
	prefixes := []string{
		"search for", "find information about", "look up", "search",
		"what is", "who is", "when was", "where is", "why is", "how to",
		"tell me about", "define", "explain",
	}

	textLower := strings.ToLower(text)
	for _, prefix := range prefixes {
		if strings.HasPrefix(textLower, prefix) {
			// Preserve original case for the query
			return strings.TrimSpace(text[len(prefix):])
		}
	}

	// Remove trailing question marks
	text = strings.TrimSuffix(text, "?")
	return strings.TrimSpace(text)
}

// FormatSearchResultsForChat formats search results for inclusion in chat prompts.
func FormatSearchResultsForChat(results []SearchResult, query string) string {
	if len(results) == 0 {
		var b strings.Builder
		b.WriteString("No search results found for: ")
		b.WriteString(query)
		return b.String()
	}

	var sb strings.Builder
	sb.WriteString(`<search_results query="`)
	sb.WriteString(query)
	sb.WriteString(`">
`)

	resultTemplate := `<result index="%d">
<title>%s</title>
<link>%s</link>
`

	for i, result := range results {
		sb.WriteString(fmt.Sprintf(resultTemplate, i+1, result.Title, result.Link))

		if result.Content != "" {
			// Truncate content if too long
			content := result.Content
			if len(content) > 500 {
				content = content[:500] + "..."
			}
			sb.WriteString("<summary>")
			sb.WriteString(content)
			sb.WriteString("</summary>\n")
		}

		if result.PublishDate != "" {
			sb.WriteString("<publish_date>")
			sb.WriteString(result.PublishDate)
			sb.WriteString("</publish_date>\n")
		}

		if result.Media != "" {
			sb.WriteString("<source>")
			sb.WriteString(result.Media)
			sb.WriteString("</source>\n")
		}

		sb.WriteString("</result>\n\n")
	}

	sb.WriteString("</search_results>")
	return sb.String()
}

// DetectSearchIntent checks if a prompt has search intent and returns the query.
func DetectSearchIntent(text string) (hasIntent bool, query string) {
	if IsSearchPrompt(text) {
		return true, ExtractSearchQuery(text)
	}
	return false, ""
}

// FormatSearchForContext formats search results as XML context for prompt augmentation.
// This is used by the --search flag to prepend search results to prompts.
func FormatSearchForContext(results []SearchResult) string {
	if len(results) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("<web_search_results>\n")

	resultTemplate := `<result>
<title>%s</title>
<url>%s</url>
`

	for _, result := range results {
		sb.WriteString(fmt.Sprintf(resultTemplate, result.Title, result.Link))

		if result.Content != "" {
			// Truncate very long content to keep context manageable
			content := result.Content
			if len(content) > 1000 {
				content = content[:1000] + "..."
			}
			sb.WriteString("<content>")
			sb.WriteString(content)
			sb.WriteString("</content>\n")
		}

		if result.PublishDate != "" {
			sb.WriteString("<date>")
			sb.WriteString(result.PublishDate)
			sb.WriteString("</date>\n")
		}

		sb.WriteString("</result>\n")
	}

	sb.WriteString("</web_search_results>")
	return sb.String()
}
