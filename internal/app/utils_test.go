package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExtractURLs tests the ExtractURLs function with table-driven tests.
func TestExtractURLs(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
		firstURL      string // If expectedCount > 0
	}{
		{
			name:          "empty string",
			input:         "",
			expectedCount: 0,
		},
		{
			name:          "no URLs",
			input:         "This is just plain text with no URLs",
			expectedCount: 0,
		},
		{
			name:          "single HTTP URL",
			input:         "Check out https://example.com",
			expectedCount: 1,
			firstURL:      "https://example.com",
		},
		{
			name:          "single HTTPS URL",
			input:         "Visit http://test.org for more info",
			expectedCount: 1,
			firstURL:      "http://test.org",
		},
		{
			name:          "multiple URLs",
			input:         "Go to https://example.com and http://test.org",
			expectedCount: 2,
			firstURL:      "https://example.com",
		},
		{
			name:          "URL with path",
			input:         "See https://example.com/path/to/resource",
			expectedCount: 1,
			firstURL:      "https://example.com/path/to/resource",
		},
		{
			name:          "URL with query params",
			input:         "Link: https://example.com?foo=bar&baz=qux",
			expectedCount: 1,
			firstURL:      "https://example.com?foo=bar&baz=qux",
		},
		{
			name:          "URL with www prefix",
			input:         "Go to www.example.com",
			expectedCount: 1,
			firstURL:      "https://www.example.com", // Auto-adds https://
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractURLs(tt.input)
			assert.Len(t, result, tt.expectedCount)
			if tt.expectedCount > 0 {
				assert.Equal(t, tt.firstURL, result[0])
			}
		})
	}
}

// TestNormalizeURL tests the normalizeURL function with table-driven tests.
func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "valid URL",
			input: "https://example.com",
		},
		{
			name:  "trailing slash",
			input: "https://example.com/",
		},
		{
			name:  "with path",
			input: "https://example.com/path",
		},
		{
			name:  "URL with query params",
			input: "https://example.com?foo=bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeURL(tt.input)
			// Just verify it produces a valid URL without error
			assert.NotEmpty(t, result)
		})
	}
}

// TestFormatWebContent tests the FormatWebContent function.
func TestFormatWebContent(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		title   string
		content string
	}{
		{
			name:    "basic formatting",
			url:     "https://example.com",
			title:   "Example Site",
			content: "This is the content",
		},
		{
			name:    "empty content",
			url:     "https://example.com",
			title:   "Example Site",
			content: "",
		},
		{
			name:    "multi-line content",
			url:     "https://example.com",
			title:   "Example Site",
			content: "Line 1\nLine 2\nLine 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatWebContent(tt.url, tt.title, tt.content)
			// Just verify it contains the key elements
			assert.Contains(t, result, "<web_content>")
			assert.Contains(t, result, tt.url)
			assert.Contains(t, result, tt.title)
			assert.Contains(t, result, "</web_content>")
		})
	}
}

// TestFormatSearchForContext tests the FormatSearchForContext function.
func TestFormatSearchForContext(t *testing.T) {
	results := []SearchResult{
		{
			Title:   "First Result",
			Link:    "https://example.com/1",
			Content: "Content of first result",
		},
		{
			Title:   "Second Result",
			Link:    "https://example.com/2",
			Content: "Content of second result",
		},
	}

	result := FormatSearchForContext(results)

	assert.Contains(t, result, "<web_search_results>")
	assert.Contains(t, result, "First Result")
	assert.Contains(t, result, "https://example.com/1")
	assert.Contains(t, result, "Content of first result")
	assert.Contains(t, result, "Second Result")
	assert.Contains(t, result, "</web_search_results>")
}

// TestDefaultChatOptions tests the DefaultChatOptions function.
func TestDefaultChatOptions(t *testing.T) {
	opts := DefaultChatOptions()

	assert.NotNil(t, opts)
	assert.Equal(t, 0.6, *opts.Temperature)
	assert.Equal(t, 8192, *opts.MaxTokens)
	assert.Equal(t, 0.9, *opts.TopP)
}
