package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mime/multipart"
	"golang.org/x/sync/errgroup"
)

// ClientConfig holds all configuration for the ZAI client.
// Injected at construction time - no global state.
type ClientConfig struct {
	APIKey     string
	BaseURL    string
	Model      string
	Timeout    time.Duration
	Verbose    bool
	RetryConfig RetryConfig
}

// DefaultChatOptions returns sensible defaults for CLI usage.
func DefaultChatOptions() ChatOptions {
	return ChatOptions{
		Temperature:  Float64Ptr(0.6),
		MaxTokens:    IntPtr(8192),
		TopP:         Float64Ptr(0.9),
		WebEnabled:   BoolPtr(true),
		WebTimeout:   IntPtr(20),
		Think:        false, // Legacy field default
	}
}

// Helper functions for creating pointers to literals (exported for use in cmd package)
func Float64Ptr(v float64) *float64 { return &v }
func IntPtr(v int) *int             { return &v }
func BoolPtr(v bool) *bool          { return &v }

// Logger interface for output control (ISP compliance).
type Logger interface {
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
}

// StderrLogger writes to stderr with emoji prefixes.
type StderrLogger struct {
	Verbose bool
}

func (l *StderrLogger) Info(format string, args ...any) {
	if l.Verbose {
		fmt.Fprintf(os.Stderr, "📡 "+format+"\n", args...)
	}
}

func (l *StderrLogger) Warn(format string, args ...any) {
	if l.Verbose {
		fmt.Fprintf(os.Stderr, "⚠️  "+format+"\n", args...)
	}
}

func (l *StderrLogger) Error(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "❌ "+format+"\n", args...)
}

// ChatClient interface for testability (ISP compliance).
type ChatClient interface {
	Chat(ctx context.Context, prompt string, opts ChatOptions) (string, error)
	Vision(ctx context.Context, prompt string, imageBase64 string, opts VisionOptions) (string, error)
	ListModels(ctx context.Context) ([]Model, error)
	GenerateImage(ctx context.Context, prompt string, opts ImageOptions) (*ImageResponse, error)
	FetchWebContent(ctx context.Context, url string, opts *WebReaderOptions) (*WebReaderResponse, error)
	SearchWeb(ctx context.Context, query string, opts SearchOptions) (*WebSearchResponse, error)
	TranscribeAudio(ctx context.Context, audioPath string, opts TranscriptionOptions) (*TranscriptionResponse, error)
}

// HistoryStore interface for storage abstraction (ISP compliance).
type HistoryStore interface {
	Save(entry HistoryEntry) error
	GetRecent(limit int) ([]HistoryEntry, error)
}

// HTTPDoer interface for HTTP operations (DIP compliance, enables testing).
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// FileReader interface for file operations (DIP compliance, enables testing).
type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

// OSFileReader implements FileReader using os.ReadFile.
type OSFileReader struct{}

// ReadFile reads a file from the filesystem.
func (r OSFileReader) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// Client implements ChatClient with Z.AI API.
type Client struct {
	config     ClientConfig
	httpClient HTTPDoer
	logger     Logger
	history    HistoryStore
	fileReader FileReader
}

// ClientDeps holds optional dependencies for NewClient.
// Zero values mean "use default implementation".
type ClientDeps struct {
	HTTPClient HTTPDoer
	FileReader FileReader
}

// NewClient creates a client with injected dependencies.
// If deps is nil or fields are nil, default implementations are used.
func NewClient(cfg ClientConfig, logger Logger, history HistoryStore, httpClient HTTPDoer) *Client {
	return NewClientWithDeps(cfg, logger, history, &ClientDeps{HTTPClient: httpClient})
}

// NewClientWithDeps creates a client with full dependency injection.
// Allows injection of all dependencies for testing.
func NewClientWithDeps(cfg ClientConfig, logger Logger, history HistoryStore, deps *ClientDeps) *Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	var httpClient HTTPDoer
	var fileReader FileReader

	if deps != nil {
		httpClient = deps.HTTPClient
		fileReader = deps.FileReader
	}

	if httpClient == nil {
		httpClient = &http.Client{Timeout: timeout}
	}
	if fileReader == nil {
		fileReader = OSFileReader{}
	}

	return &Client{
		config:     cfg,
		httpClient: httpClient,
		logger:     logger,
		history:    history,
		fileReader: fileReader,
	}
}

// HTTPClient returns the underlying HTTP client for connection reuse.
func (c *Client) HTTPClient() HTTPDoer {
	return c.httpClient
}

// Chat sends a prompt and returns the response.
// Orchestrates content building, URL enrichment, and request execution.
func (c *Client) Chat(ctx context.Context, prompt string, opts ChatOptions) (string, error) {
	if c.config.APIKey == "" {
		return "", fmt.Errorf("API key is not configured")
	}

	// Build message content (with optional file)
	content, err := c.buildContent(ctx, prompt, opts.FilePath)
	if err != nil {
		return "", err
	}

	// Enrich content with web URLs if enabled
	content = c.enrichWithURLContent(ctx, prompt, content, opts)

	// Build messages array with context
	messages := c.buildMessagesWithContext(content, opts)

	// Handle legacy Think field
	if opts.Think && opts.Thinking == nil {
		opts.Thinking = &opts.Think
	}

	// Execute request with retry
	response, usage, err := c.doRequestWithRetry(ctx, messages, opts)
	if err != nil {
		return "", err
	}

	// Save to history (non-blocking, log errors)
	c.saveToHistory(prompt, response, usage)

	return response, nil
}

// enrichWithURLContent fetches web content for URLs in the prompt if web is enabled.
// Uses concurrent fetching with errgroup for improved performance.
func (c *Client) enrichWithURLContent(ctx context.Context, prompt, content string, opts ChatOptions) string {
	if !c.isWebEnabled(opts) {
		return content
	}

	urls := ExtractURLs(prompt)
	if len(urls) == 0 {
		return content
	}

	webOpts := c.defaultWebReaderOptions(opts.WebTimeout)

	// Use errgroup for concurrent URL fetching
	g, ctx := errgroup.WithContext(ctx)
	results := make([]struct {
		url   string
		title string
		body  string
	}, len(urls))

	// Fetch all URLs concurrently
	for i, url := range urls {
		i, url := i, url // capture loop variables
		g.Go(func() error {
			webResp, err := c.FetchWebContent(ctx, url, webOpts)
			if err != nil {
				c.logger.Warn("Failed to fetch web content from %s: %v", url, err)
				return nil // Don't fail entire group for single URL error
			}
			results[i].url = url
			results[i].title = webResp.ReaderResult.Title
			results[i].body = webResp.ReaderResult.Content
			return nil
		})
	}

	// Wait for all fetches to complete
	if err := g.Wait(); err != nil {
		c.logger.Warn("Error fetching web content: %v", err)
	}

	// Append results in original order
	for _, r := range results {
		if r.url != "" { // Only append successful fetches
			content += "\n\n" + FormatWebContent(r.url, r.title, r.body)
		}
	}

	return content
}

// isWebEnabled checks if web content fetching is enabled.
func (c *Client) isWebEnabled(opts ChatOptions) bool {
	if opts.WebEnabled != nil {
		return *opts.WebEnabled
	}
	return true // Default to enabled
}

// defaultWebReaderOptions creates default options for web content fetching.
func (c *Client) defaultWebReaderOptions(timeout *int) *WebReaderOptions {
	trueVal := true
	falseVal := false
	return &WebReaderOptions{
		Timeout:          timeout,
		ReturnFormat:     "markdown",
		RetainImages:     &trueVal,
		NoCache:          &falseVal,
		NoGFM:            &falseVal,
		KeepImgDataURL:   &falseVal,
		WithImagesSummary: &falseVal,
		WithLinksSummary:  &falseVal,
	}
}

// buildMessagesWithContext constructs messages array including conversation context.
func (c *Client) buildMessagesWithContext(content string, opts ChatOptions) []Message {
	messages := c.buildMessages(content, opts)

	// Prepend context messages if provided
	if len(opts.Context) > 0 {
		messages = append(opts.Context, messages...)
	}

	return messages
}

// saveToHistory persists the chat exchange to history storage.
func (c *Client) saveToHistory(prompt, response string, usage Usage) {
	if c.history == nil {
		return
	}
	entry := NewChatHistoryEntry(time.Now(), prompt, response, c.config.Model, usage)
	if err := c.history.Save(entry); err != nil {
		c.logger.Warn("Failed to save to history: %v", err)
	}
}

// buildContent combines prompt with optional file contents or URL content.
func (c *Client) buildContent(ctx context.Context, prompt, filePath string) (string, error) {
	if filePath == "" {
		return prompt, nil
	}

	// Check if it's a URL
	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		// Fetch web content
		webOpts := &WebReaderOptions{
			ReturnFormat: "markdown",
		}
		resp, err := c.FetchWebContent(ctx, filePath, webOpts)
		if err != nil {
			return "", fmt.Errorf("failed to fetch URL %s: %w", filePath, err)
		}
		return fmt.Sprintf("%s\n\n<web_content url=\"%s\" title=\"%s\">\n%s\n</web_content>",
			prompt, filePath, resp.ReaderResult.Title, resp.ReaderResult.Content), nil
	}

	// Local file
	data, err := c.fileReader.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return fmt.Sprintf("%s\n\nFile contents (%s):\n```\n%s\n```", prompt, filePath, string(data)), nil
}

// buildMessages constructs the messages array for the API.
func (c *Client) buildMessages(content string, opts ChatOptions) []Message {
	var messages []Message

	// Add system prompt
	messages = append(messages, Message{
		Role:    "system",
		Content: "Be concise and direct. Answer briefly and to the point.",
	})

	// Add current user message
	messages = append(messages, Message{
		Role:    "user",
		Content: content,
	})

	return messages
}

// isRetryableError checks if an error should trigger a retry.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Network errors: timeout, connection refused, etc.
	var netErr interface{ Timeout() bool }
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	// Check for specific error patterns
	errStr := err.Error()
	retryablePatterns := []string{
		"connection refused",
		"connection reset",
		"temporary failure",
		"timeout",
		"503", // Service Unavailable
		"502", // Bad Gateway
		"504", // Gateway Timeout
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(strings.ToLower(errStr), pattern) {
			return true
		}
	}

	return false
}

// calculateBackoff calculates exponential backoff with jitter.
func calculateBackoff(attempt int, initialBackoff, maxBackoff time.Duration) time.Duration {
	// Exponential backoff: initial * 2^(attempt-1)
	backoff := initialBackoff * time.Duration(1<<uint(attempt-1))

	// Cap at max backoff
	if backoff > maxBackoff {
		backoff = maxBackoff
	}

	// Add jitter (±12.5%, centered - so jitter can add or subtract up to 12.5%)
	// This ensures we never go below the base value by more than 12.5%
	jitterRange := float64(backoff) * 0.125
	jitter := time.Duration(jitterRange * (2.0*rand.Float64() - 1.0))

	return backoff + jitter
}

// doRequest executes the HTTP request to Z.AI API.
// Single place for all HTTP logic (DRY compliance).
func (c *Client) doRequest(ctx context.Context, messages []Message, opts ChatOptions) (string, Usage, error) {
	reqData := ChatRequest{
		Model:       c.config.Model,
		Messages:    messages,
		Stream:      false,
		Thinking:    &Thinking{Type: "disabled"},
	}

	// Apply optional overrides
	if opts.Temperature != nil {
		reqData.Temperature = *opts.Temperature
	} else {
		reqData.Temperature = 0.6 // default
	}

	if opts.MaxTokens != nil {
		reqData.MaxTokens = *opts.MaxTokens
	} else {
		reqData.MaxTokens = 8192 // default
	}

	if opts.TopP != nil {
		reqData.TopP = *opts.TopP
	} else {
		reqData.TopP = 0.9 // default
	}

	// Apply model override if provided
	if opts.Model != "" {
		reqData.Model = opts.Model
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", Usage{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", Usage{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Accept-Language", "en-US,en")

	c.logger.Info("Sending request to: %s", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", Usage{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", Usage{}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", Usage{}, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", Usage{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", Usage{}, fmt.Errorf("no choices in response")
	}

	c.logger.Info("Usage: %d tokens (prompt: %d, completion: %d)",
		chatResp.Usage.TotalTokens,
		chatResp.Usage.PromptTokens,
		chatResp.Usage.CompletionTokens)

	return chatResp.Choices[0].Message.Content, chatResp.Usage, nil
}

// doRequestWithRetry executes doRequest with exponential backoff retry logic.
func (c *Client) doRequestWithRetry(ctx context.Context, messages []Message, opts ChatOptions) (string, Usage, error) {
	var lastErr error

	// Apply defaults for zero values
	maxAttempts := c.config.RetryConfig.MaxAttempts
	if maxAttempts < 1 {
		maxAttempts = 1 // No retry if not configured
	}

	initialBackoff := c.config.RetryConfig.InitialBackoff
	if initialBackoff < 1 {
		initialBackoff = 1 * time.Second
	}

	maxBackoff := c.config.RetryConfig.MaxBackoff
	if maxBackoff < 1 {
		maxBackoff = 30 * time.Second
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Check context before attempting
		select {
		case <-ctx.Done():
			return "", Usage{}, ctx.Err()
		default:
		}

		// On retry (not first attempt), log and wait
		if attempt > 1 {
			backoff := calculateBackoff(attempt, initialBackoff, maxBackoff)
			c.logger.Info("Retrying request (attempt %d/%d) after %v: %v", attempt, maxAttempts, backoff, lastErr)

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return "", Usage{}, ctx.Err()
			}
		}

		// Execute request
		response, usage, err := c.doRequest(ctx, messages, opts)
		if err == nil {
			return response, usage, nil
		}

		lastErr = err

		// Don't retry if error is not retryable or this was the last attempt
		if !isRetryableError(err) || attempt == maxAttempts {
			break
		}
	}

	return "", Usage{}, fmt.Errorf("request failed after %d attempts: %w", maxAttempts, lastErr)
}

// ListModels fetches available models from the API.
func (c *Client) ListModels(ctx context.Context) ([]Model, error) {
	if c.config.APIKey == "" {
		return nil, fmt.Errorf("API key is not configured")
	}

	url := fmt.Sprintf("%s/models", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	c.logger.Info("Fetching models from: %s", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var modelsResp ModelsResponse
	if err := json.Unmarshal(body, &modelsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return modelsResp.Data, nil
}

// GenerateImage creates an image using the Z.AI image generation API.
func (c *Client) GenerateImage(ctx context.Context, prompt string, opts ImageOptions) (*ImageResponse, error) {
	if c.config.APIKey == "" {
		return nil, fmt.Errorf("API key is not configured")
	}

	// Validate options
	if err := validateImageOptions(opts); err != nil {
		return nil, fmt.Errorf("invalid image options: %w", err)
	}

	// Build request
	model := opts.Model
	if model == "" {
		model = "cogview-4-250304" // Default image model
	}

	reqData := ImageGenerationRequest{
		Model:   model,
		Prompt:  prompt,
		Quality: opts.Quality,
		Size:    opts.Size,
		UserID:  opts.UserID,
	}

	// Set defaults
	if reqData.Quality == "" {
		reqData.Quality = "hd"
	}
	if reqData.Size == "" {
		reqData.Size = "1024x1024"
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal image request: %w", err)
	}

	url := fmt.Sprintf("%s/images/generations", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create image request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Accept-Language", "en-US,en")

	c.logger.Info("Sending image generation request to: %s", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send image request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("image generation API error: %d - %s", resp.StatusCode, string(body))
	}

	var imageResp ImageResponse
	if err := json.Unmarshal(body, &imageResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal image response: %w", err)
	}

	if len(imageResp.Data) == 0 {
		return nil, fmt.Errorf("no images in response")
	}

	c.logger.Info("Generated image: %s (%dx%d)", imageResp.Data[0].URL, imageResp.Data[0].Width, imageResp.Data[0].Height)

	return &imageResp, nil
}

// FetchWebContent retrieves and processes web content from a URL.
func (c *Client) FetchWebContent(ctx context.Context, url string, opts *WebReaderOptions) (*WebReaderResponse, error) {
	if c.config.APIKey == "" {
		return nil, fmt.Errorf("API key is not configured")
	}

	// Validate URL
	if url == "" {
		return nil, fmt.Errorf("URL is required")
	}

	// Build request with defaults
	req := WebReaderRequest{
		URL:          url,
		ReturnFormat: "markdown",
	}
	trueVal := true
	req.RetainImages = &trueVal

	// Apply options
	if opts != nil {
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

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal web reader request: %w", err)
	}

	// Use /paas/v4/reader endpoint
	readerURL := fmt.Sprintf("%s/reader", c.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", readerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create web reader request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	httpReq.Header.Set("Accept-Language", "en-US,en")

	c.logger.Info("Fetching web content from: %s", url)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch web content: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read web reader response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("web reader API error: %d - %s", resp.StatusCode, string(body))
	}

	var webResp WebReaderResponse
	if err := json.Unmarshal(body, &webResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal web reader response: %w", err)
	}

	c.logger.Info("Successfully fetched web content: %s (title: %s)",
		webResp.ReaderResult.URL, webResp.ReaderResult.Title)

	return &webResp, nil
}

// validateImageOptions checks if image options are valid.
func validateImageOptions(opts ImageOptions) error {
	// Validate quality
	if opts.Quality != "" && opts.Quality != "hd" && opts.Quality != "standard" {
		return fmt.Errorf("invalid quality: %s (must be 'hd' or 'standard')", opts.Quality)
	}

	// Validate size format
	if opts.Size != "" {
		supportedSizes := map[string]bool{
			"1024x1024": true,
			"1024x768":  true,
			"768x1024":  true,
			"512x512":   true,
		}
		if !supportedSizes[opts.Size] {
			return fmt.Errorf("invalid size: %s (supported: 1024x1024, 1024x768, 768x1024, 512x512)", opts.Size)
		}
	}

	return nil
}

// SearchWeb performs a web search using Z.AI's search API.
func (c *Client) SearchWeb(ctx context.Context, query string, opts SearchOptions) (*WebSearchResponse, error) {
	if c.config.APIKey == "" {
		return nil, fmt.Errorf("API key is not configured")
	}

	// Validate query
	if query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	// Validate count
	if opts.Count < 1 || opts.Count > 50 {
		return nil, fmt.Errorf("count must be between 1 and 50")
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
	reqData := WebSearchRequest{
		SearchEngine: "search-prime",
		SearchQuery:  query,
		Count:        &opts.Count,
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

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search request: %w", err)
	}

	// Use /paas/v4/web_search endpoint
	url := fmt.Sprintf("%s/web_search", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Accept-Language", "en-US,en")

	c.logger.Info("Searching web for: %s", query)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send search request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read search response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Try to parse error response
		var apiErr struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != "" {
			return nil, fmt.Errorf("search API error: %s - %s", apiErr.Error, apiErr.Message)
		}
		return nil, fmt.Errorf("search API error: %d - %s", resp.StatusCode, string(body))
	}

	var searchResp WebSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	c.logger.Info("Found %d search results for query: %s", len(searchResp.SearchResult), query)

	// Save to history (non-blocking, log errors)
	if c.history != nil {
		entry := NewSearchHistoryEntry(time.Now(), query, &searchResp)
		if err := c.history.Save(entry); err != nil {
			c.logger.Warn("Failed to save search to history: %v", err)
		}
	}

	return &searchResp, nil
}

// Vision analyzes an image using Z.AI's vision model (glm-4.6v).
// imageBase64 should be a data URI like "data:image/jpeg;base64,<base64-data>" or a raw base64 string.
func (c *Client) Vision(ctx context.Context, prompt string, imageBase64 string, opts VisionOptions) (string, error) {
	if c.config.APIKey == "" {
		return "", fmt.Errorf("API key is not configured")
	}

	// Validate prompt
	if prompt == "" {
		prompt = "What do you see in this image? Please describe it in detail."
	}

	// Validate image input
	if imageBase64 == "" {
		return "", fmt.Errorf("image data is required")
	}

	// Build vision model
	model := opts.Model
	if model == "" {
		model = "glm-4.6v" // Default vision model
	}

	// Build multimodal messages
	messages := []VisionMessage{
		{
			Role: "user",
			Content: []ContentPart{
				{
					Type: "text",
					Text: prompt,
				},
				{
					Type: "image_url",
					ImageURL: &ImageURLContent{
						URL: imageBase64,
					},
				},
			},
		},
	}

	// Build request
	reqData := VisionRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	// Apply optional overrides
	if opts.Temperature != nil {
		reqData.Temperature = *opts.Temperature
	} else {
		reqData.Temperature = 0.3 // Lower temp for vision
	}

	if opts.MaxTokens != nil {
		reqData.MaxTokens = *opts.MaxTokens
	} else {
		reqData.MaxTokens = 4096
	}

	if opts.TopP != nil {
		reqData.TopP = *opts.TopP
	} else {
		reqData.TopP = 0.9
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal vision request: %w", err)
	}

	// Use chat/completions endpoint for vision
	url := fmt.Sprintf("%s/chat/completions", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create vision request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Accept-Language", "en-US,en")

	c.logger.Info("Sending vision request to: %s", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send vision request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read vision response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vision API error: %d - %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal vision response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in vision response")
	}

	c.logger.Info("Vision complete: %d tokens (prompt: %d, completion: %d)",
		chatResp.Usage.TotalTokens,
		chatResp.Usage.PromptTokens,
		chatResp.Usage.CompletionTokens)

	return chatResp.Choices[0].Message.Content, nil
}

// TranscribeAudio transcribes an audio file using Z.AI's ASR model.
func (c *Client) TranscribeAudio(ctx context.Context, audioPath string, opts TranscriptionOptions) (*TranscriptionResponse, error) {
	if c.config.APIKey == "" {
		return nil, fmt.Errorf("API key is not configured")
	}

	// Validate audio file
	if audioPath == "" {
		return nil, fmt.Errorf("audio file path is required")
	}

	// Open audio file
	file, err := os.Open(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}
	defer file.Close()

	// Check file size
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	const maxFileSize = 25 * 1024 * 1024 // 25MB
	if stat.Size() > maxFileSize {
		return nil, fmt.Errorf("audio file too large: %d bytes (max: %d MB)", stat.Size(), maxFileSize/1024/1024)
	}

	// Build model
	model := opts.Model
	if model == "" {
		model = "glm-asr-2512"
	}

	// Build multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Add model
	writer.WriteField("model", model)

	// Add optional fields
	if opts.Prompt != "" {
		writer.WriteField("prompt", opts.Prompt)
	}
	if opts.Stream {
		writer.WriteField("stream", "true")
	}
	if opts.UserID != "" {
		writer.WriteField("user_id", opts.UserID)
	}
	if opts.RequestID != "" {
		writer.WriteField("request_id", opts.RequestID)
	}
	if len(opts.Hotwords) > 0 {
		hotwordsJSON, err := json.Marshal(opts.Hotwords)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal hotwords: %w", err)
		}
		writer.WriteField("hotwords", string(hotwordsJSON))
	}

	writer.Close()

	// Build request
	url := fmt.Sprintf("%s/audio/transcriptions", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Accept-Language", "en-US,en")

	c.logger.Info("Sending audio transcription request to: %s", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("transcription API error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var transcriptionResp TranscriptionResponse
	if err := json.Unmarshal(bodyBytes, &transcriptionResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.logger.Info("Transcription complete: %d chars, model: %s", len(transcriptionResp.Text), transcriptionResp.Model)

	return &transcriptionResp, nil
}
