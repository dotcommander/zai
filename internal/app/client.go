package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ClientConfig holds all configuration for the ZAI client.
// Injected at construction time - no global state.
type ClientConfig struct {
	APIKey  string
	BaseURL string
	Model   string
	Timeout time.Duration
	Verbose bool
}

// DefaultChatOptions returns sensible defaults for CLI usage.
func DefaultChatOptions() ChatOptions {
	return ChatOptions{
		Temperature:  float64Ptr(0.6),
		MaxTokens:    intPtr(8192),
		TopP:         float64Ptr(0.9),
		WebEnabled:   boolPtr(true),
		WebTimeout:   intPtr(20),
		Think:        false, // Legacy field default
	}
}

// Helper functions for creating pointers to literals
func float64Ptr(v float64) *float64 { return &v }
func intPtr(v int) *int { return &v }
func boolPtr(v bool) *bool { return &v }

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
	ListModels(ctx context.Context) ([]Model, error)
	GenerateImage(ctx context.Context, prompt string, opts ImageOptions) (*ImageResponse, error)
	FetchWebContent(ctx context.Context, url string, opts *WebReaderOptions) (*WebReaderResponse, error)
	SearchWeb(ctx context.Context, query string, opts SearchOptions) (*WebSearchResponse, error)
}

// HistoryStore interface for storage abstraction (ISP compliance).
type HistoryStore interface {
	Save(entry HistoryEntry) error
	GetRecent(limit int) ([]HistoryEntry, error)
}

// Client implements ChatClient with Z.AI API.
type Client struct {
	config     ClientConfig
	httpClient *http.Client
	logger     Logger
	history    HistoryStore
}

// NewClient creates a client with injected dependencies.
func NewClient(cfg ClientConfig, logger Logger, history HistoryStore) *Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	return &Client{
		config:     cfg,
		httpClient: &http.Client{Timeout: timeout},
		logger:     logger,
		history:    history,
	}
}

// Chat sends a prompt and returns the response.
// Single method handles all cases: simple, with file, with context, with web content.
func (c *Client) Chat(ctx context.Context, prompt string, opts ChatOptions) (string, error) {
	if c.config.APIKey == "" {
		return "", fmt.Errorf("API key is not configured")
	}

	// Build message content (with optional file)
	content, err := c.buildContent(prompt, opts.FilePath)
	if err != nil {
		return "", err
	}

	// Check for URLs and fetch web content if enabled
	webEnabled := false
	if opts.WebEnabled != nil {
		webEnabled = *opts.WebEnabled
	} else {
		// Check configuration for default
		// For now, default to true if web reader is configured
		webEnabled = true
	}

	if webEnabled {
		urls := ExtractURLs(prompt)
		if len(urls) > 0 {
			webCtx := ctx
			webOpts := &WebReaderOptions{
				Timeout:       opts.WebTimeout,
				ReturnFormat:  "markdown", // Default format for chat
			}
			// Set default values
			trueVal := true
			falseVal := false
			webOpts.RetainImages = &trueVal
			webOpts.NoCache = &falseVal
			webOpts.NoGFM = &falseVal
			webOpts.KeepImgDataURL = &falseVal
			webOpts.WithImagesSummary = &falseVal
			webOpts.WithLinksSummary = &falseVal

			for _, url := range urls {
				webResp, err := c.FetchWebContent(webCtx, url, webOpts)
				if err != nil {
					c.logger.Warn("Failed to fetch web content from %s: %v", url, err)
					continue
				}
				// Append formatted web content to the prompt
				content += "\n\n" + FormatWebContent(url, webResp.ReaderResult.Title, webResp.ReaderResult.Content)
			}
		}
	}

	// Build messages array
	messages := c.buildMessages(content, opts)

	// Add context messages if provided (legacy support)
	if len(opts.Context) > 0 {
		// Prepend context messages
		allMessages := append(opts.Context, messages...)
		messages = allMessages
	}

	// Handle legacy Think field
	if opts.Think && opts.Thinking == nil {
		opts.Thinking = &opts.Think
	}

	// Execute request
	response, usage, err := c.doRequest(ctx, messages, opts)
	if err != nil {
		return "", err
	}

	// Save to history (non-blocking, log errors)
	if c.history != nil {
		entry := NewChatHistoryEntry(time.Now(), prompt, response, c.config.Model, usage)
		if err := c.history.Save(entry); err != nil {
			c.logger.Warn("Failed to save to history: %v", err)
		}
	}

	return response, nil
}

// buildContent combines prompt with optional file contents.
func (c *Client) buildContent(prompt, filePath string) (string, error) {
	if filePath == "" {
		return prompt, nil
	}

	data, err := os.ReadFile(filePath)
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
