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

	if url == "" {
		return nil, errors.New("URL is required")
	}

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
func (c *Client) SearchWeb(ctx context.Context, query string, opts SearchOptions) (*WebSearchResponse, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}
	if query == "" {
		return nil, errors.New("search query is required")
	}

	if err := validateSearchOpts(opts); err != nil {
		return nil, err
	}

	// Over-fetch when filtering by non-Chinese language so we have enough after filtering
	requestedCount := opts.Count
	if opts.Language != "" && opts.Language != "zh" {
		opts.Count = min(opts.Count*2, 50)
	}

	reqData := buildSearchRequest(query, opts)

	body, err := c.executeJSONRequest(ctx, "web_search", reqData)
	if err != nil {
		return nil, unwrapSearchError(err)
	}

	var searchResp WebSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	filterSearchByLanguage(&searchResp, opts.Language, requestedCount)

	c.logger.DebugContext(ctx, "search complete", "results", len(searchResp.SearchResult), "query", query)

	if c.history != nil {
		entry := NewSearchHistoryEntry(time.Now(), query, &searchResp)
		if err := c.history.Save(entry); err != nil {
			c.logger.WarnContext(ctx, "failed to save search to history", "error", err)
		}
	}

	return &searchResp, nil
}

// validRecencyFilters is the set of accepted recency filter values.
var validRecencyFilters = map[string]bool{
	"": true, "noLimit": true, "oneDay": true,
	"oneWeek": true, "oneMonth": true, "oneYear": true,
}

// validateSearchOpts validates search options without mutating them.
func validateSearchOpts(opts SearchOptions) error {
	if opts.Count < 1 || opts.Count > 50 {
		return errors.New("count must be between 1 and 50")
	}

	if !validRecencyFilters[opts.RecencyFilter] {
		return fmt.Errorf("invalid recency filter: %s (must be one of: oneDay, oneWeek, oneMonth, oneYear, noLimit)", opts.RecencyFilter)
	}
	return nil
}

// buildSearchRequest creates a WebSearchRequest from query and options.
func buildSearchRequest(query string, opts SearchOptions) WebSearchRequest {
	req := WebSearchRequest{
		SearchEngine:   opts.SearchEngine,
		SearchQuery:    query,
		Count:          &opts.Count,
		ContentSize:    opts.ContentSize,
		SearchPrompt:   opts.SearchPrompt,
		ResultSequence: opts.ResultSequence,
	}

	if opts.DomainFilter != "" {
		req.SearchDomainFilter = &opts.DomainFilter
	}
	if opts.RecencyFilter != "" && opts.RecencyFilter != "noLimit" {
		req.SearchRecencyFilter = &opts.RecencyFilter
	}
	if opts.RequestID != "" {
		req.RequestID = &opts.RequestID
	}
	if opts.UserID != "" {
		req.UserID = &opts.UserID
	}
	if opts.SearchResult {
		req.SearchResult = &opts.SearchResult
	}
	if opts.RequireSearch {
		req.RequireSearch = &opts.RequireSearch
	}
	return req
}

// unwrapSearchError extracts a structured error message from an API error.
func unwrapSearchError(err error) error {
	var apiError *APIError
	if errors.As(err, &apiError) {
		var jsonErr struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		if json.Unmarshal([]byte(apiError.Body), &jsonErr) == nil && jsonErr.Error != "" {
			return fmt.Errorf("search API error: %s - %s", jsonErr.Error, jsonErr.Message)
		}
	}
	return fmt.Errorf("search API error: %w", err)
}

// filterSearchByLanguage filters results by language if a non-Chinese language is specified.
func filterSearchByLanguage(resp *WebSearchResponse, lang string, requestedCount int) {
	if lang == "" || lang == "zh" {
		return
	}

	filtered := make([]SearchResult, 0, len(resp.SearchResult))
	for _, r := range resp.SearchResult {
		if isLikelyLanguage(r.Title, lang) && isLikelyLanguage(r.Content, lang) {
			filtered = append(filtered, r)
		}
	}
	if requestedCount > 0 && len(filtered) > requestedCount {
		filtered = filtered[:requestedCount]
	}
	resp.SearchResult = filtered
}

type charCounts struct {
	latin, cjk, cyrillic, arabic int
}

func (c charCounts) total() int { return c.latin + c.cjk + c.cyrillic + c.arabic }

func isLatin(r rune) bool { return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') }

// isCJK covers CJK Unified Ideographs, Extension A, Hiragana/Katakana, and Hangul Syllables.
func isCJK(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) || (r >= 0x3400 && r <= 0x4DBF) ||
		(r >= 0x3040 && r <= 0x30FF) || (r >= 0xAC00 && r <= 0xD7AF)
}

// maxLangSampleRunes caps character scanning — language is determinable from a small sample.
const maxLangSampleRunes = 300

func countCharClasses(text string) charCounts {
	var c charCounts
	n := 0
	for _, r := range text {
		if n >= maxLangSampleRunes {
			break
		}
		n++
		switch {
		case isLatin(r):
			c.latin++
		case isCJK(r):
			c.cjk++
		case r >= 0x0400 && r <= 0x04FF:
			c.cyrillic++
		case r >= 0x0600 && r <= 0x06FF:
			c.arabic++
		}
	}
	return c
}

// isLikelyLanguage checks if text is likely written in the target language.
// Uses character-class heuristics: Latin for "en"/"es"/"fr"/"de"/"pt"/"it",
// CJK for "zh"/"ja"/"ko", Cyrillic for "ru", Arabic for "ar".
func isLikelyLanguage(text, lang string) bool {
	if text == "" || lang == "" {
		return true
	}

	cc := countCharClasses(text)
	total := cc.total()
	if total == 0 {
		return true // all punctuation/numbers — can't determine
	}
	ftotal := float64(total)

	switch lang {
	case "en", "es", "fr", "de", "pt", "it":
		return cc.cjk == 0 && float64(cc.latin)/ftotal > 0.5
	case "zh", "ja", "ko":
		return float64(cc.cjk)/ftotal > 0.3
	case "ru":
		return float64(cc.cyrillic)/ftotal > 0.3
	case "ar":
		return float64(cc.arabic)/ftotal > 0.3
	default:
		return true
	}
}
