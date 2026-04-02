package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// FetchWebContent retrieves and processes web content from a URL.
func (c *Client) FetchWebContent(ctx context.Context, url string, opts *WebReaderOptions) (*WebReaderResponse, error) { //nolint:gocognit
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	// Validate URL
	if err := c.validateWebContentURL(url); err != nil {
		return nil, err
	}

	// Build request with defaults and options
	req := c.buildWebReaderRequest(url, opts)

	// Execute API request and parse response
	webResp, err := c.executeWebReaderRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	c.logger.DebugContext(ctx, "fetched web content",
		"url", webResp.ReaderResult.URL,
		"title", webResp.ReaderResult.Title)

	return &webResp, nil
}

// validateWebContentURL validates the URL parameter for web content fetching.
func (c *Client) validateWebContentURL(url string) error {
	if url == "" {
		return errors.New("URL is required")
	}
	return nil
}

// buildWebReaderRequest builds a WebReaderRequest with defaults and applies options.
func (c *Client) buildWebReaderRequest(url string, opts *WebReaderOptions) WebReaderRequest {
	// Build request with defaults
	req := WebReaderRequest{
		URL:          url,
		ReturnFormat: "markdown",
	}
	trueVal := true
	req.RetainImages = &trueVal

	// Apply options
	if opts != nil {
		c.applyWebReaderOptions(&req, opts)
	}

	return req
}

// applyWebReaderOptions applies WebReaderOptions to the request.
func (c *Client) applyWebReaderOptions(req *WebReaderRequest, opts *WebReaderOptions) {
	if opts.Timeout != nil {
		req.Timeout = opts.Timeout
	}
	if opts.NoCache != nil {
		req.NoCache = opts.NoCache
	}
	if opts.ReturnFormat != "" {
		req.ReturnFormat = opts.ReturnFormat
	}
	if opts.RetainImages != nil {
		req.RetainImages = opts.RetainImages
	}
	if opts.NoGFM != nil {
		req.NoGFM = opts.NoGFM
	}
	if opts.KeepImgDataURL != nil {
		req.KeepImgDataURL = opts.KeepImgDataURL
	}
	if opts.WithImagesSummary != nil {
		req.WithImagesSummary = opts.WithImagesSummary
	}
	if opts.WithLinksSummary != nil {
		req.WithLinksSummary = opts.WithLinksSummary
	}
}

// executeWebReaderRequest executes the web reader API call and parses the response.
func (c *Client) executeWebReaderRequest(ctx context.Context, req WebReaderRequest) (WebReaderResponse, error) {
	var webResp WebReaderResponse
	body, err := c.executeJSONRequest(ctx, "reader", req)
	if err != nil {
		return WebReaderResponse{}, fmt.Errorf("web reader API error: %w", err)
	}
	if err := json.Unmarshal(body, &webResp); err != nil {
		return WebReaderResponse{}, fmt.Errorf("failed to unmarshal web reader response: %w", err)
	}
	return webResp, nil
}

// SearchWeb performs a web search using Z.AI's search API.
func (c *Client) SearchWeb(ctx context.Context, query string, opts SearchOptions) (*WebSearchResponse, error) { //nolint:gocognit,gocyclo,revive // cognitive complexity is inherent to search validation
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	// Validate query
	if query == "" {
		return nil, errors.New("search query is required")
	}

	// Over-fetch when filtering by non-Chinese language
	requestedCount := opts.Count
	if opts.Language != "" && opts.Language != "zh" {
		opts.Count = min(opts.Count*3, 50)
	}

	// Validate count
	if opts.Count < 1 || opts.Count > 50 {
		return nil, errors.New("count must be between 1 and 50")
	}

	// Validate recency filter
	validRecencyFilters := map[string]bool{
		"": true, "noLimit": true,
		"oneDay": true, "oneWeek": true,
		"oneMonth": true, "oneYear": true,
	}
	if !validRecencyFilters[opts.RecencyFilter] {
		return nil, fmt.Errorf("invalid recency filter: %s (must be one of: oneDay, oneWeek, oneMonth, oneYear, noLimit)", opts.RecencyFilter)
	}

	// Build request
	searchEngine := opts.SearchEngine
	if searchEngine == "" {
		searchEngine = "search_std"
	}
	reqData := WebSearchRequest{
		SearchEngine:   searchEngine,
		SearchQuery:    query,
		Count:          &opts.Count,
		ContentSize:    opts.ContentSize,
		SearchPrompt:   opts.SearchPrompt,
		ResultSequence: opts.ResultSequence,
	}

	// Add optional parameters
	if opts.DomainFilter != "" {
		reqData.SearchDomainFilter = &opts.DomainFilter
	}
	if opts.RecencyFilter != "" && opts.RecencyFilter != "noLimit" {
		reqData.SearchRecencyFilter = &opts.RecencyFilter
	}
	if opts.RequestID != "" {
		reqData.RequestID = &opts.RequestID
	}
	if opts.UserID != "" {
		reqData.UserID = &opts.UserID
	}
	if opts.SearchResult {
		reqData.SearchResult = &opts.SearchResult
	}
	if opts.RequireSearch {
		reqData.RequireSearch = &opts.RequireSearch
	}

	var searchResp WebSearchResponse
	body, err := c.executeJSONRequest(ctx, "web_search", reqData)
	if err != nil {
		// Try to extract structured error from API response
		var apiError *APIError
		if errors.As(err, &apiError) {
			// Try to parse JSON error body
			var jsonErr struct {
				Error   string `json:"error"`
				Message string `json:"message"`
			}
			if json.Unmarshal([]byte(apiError.Body), &jsonErr) == nil && jsonErr.Error != "" {
				return nil, fmt.Errorf("search API error: %s - %s", jsonErr.Error, jsonErr.Message)
			}
		}
		return nil, fmt.Errorf("search API error: %w", err)
	}
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	// Filter results by language if specified
	if opts.Language != "" && opts.Language != "zh" {
		filtered := make([]SearchResult, 0, len(searchResp.SearchResult))
		for _, r := range searchResp.SearchResult {
			// Check both title and content for language match
			if isLikelyLanguage(r.Title, opts.Language) && isLikelyLanguage(r.Content, opts.Language) {
				filtered = append(filtered, r)
			}
		}
		if requestedCount > 0 && len(filtered) > requestedCount {
			filtered = filtered[:requestedCount]
		}
		searchResp.SearchResult = filtered
	}

	c.logger.DebugContext(ctx, "search complete", "results", len(searchResp.SearchResult), "query", query)

	// Save to history (non-blocking, log errors)
	if c.history != nil {
		entry := NewSearchHistoryEntry(time.Now(), query, &searchResp)
		if err := c.history.Save(entry); err != nil {
			c.logger.WarnContext(ctx, "failed to save search to history", "error", err)
		}
	}

	return &searchResp, nil
}

// isLikelyLanguage checks if text is likely written in the target language.
// Uses character-class heuristics: Latin scripts for "en"/"es"/"fr"/"de"/"pt"/"it",
// CJK detection for "zh"/"ja"/"ko", Cyrillic for "ru", Arabic for "ar".
func isLikelyLanguage(text, lang string) bool {
	if text == "" || lang == "" {
		return true
	}

	runes := []rune(text)
	if len(runes) == 0 {
		return true
	}

	// Count character classes
	var latin, cjk, cyrillic, arabic int
	for _, r := range runes {
		switch {
		case (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z'):
			latin++
		case r >= 0x4E00 && r <= 0x9FFF || r >= 0x3400 && r <= 0x4DBF || r >= 0x3040 && r <= 0x30FF || r >= 0xAC00 && r <= 0xD7AF:
			cjk++
		case r >= 0x0400 && r <= 0x04FF:
			cyrillic++
		case r >= 0x0600 && r <= 0x06FF:
			arabic++
		}
	}

	total := latin + cjk + cyrillic + arabic
	if total == 0 {
		return true // all punctuation/numbers — can't determine
	}

	switch lang {
	case "en", "es", "fr", "de", "pt", "it":
		// Reject if any CJK characters present — Latin-language content shouldn't contain them
		if cjk > 0 {
			return false
		}
		return float64(latin)/float64(total) > 0.5
	case "zh", "ja", "ko":
		return float64(cjk)/float64(total) > 0.3
	case "ru":
		return float64(cyrillic)/float64(total) > 0.3
	case "ar":
		return float64(arabic)/float64(total) > 0.3
	default:
		return true // unknown language, don't filter
	}
}
