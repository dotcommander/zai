package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"mime/multipart"

	"golang.org/x/time/rate"

	"golang.org/x/sync/errgroup"

	"github.com/dotcommander/zai/internal/app/utils"
	"github.com/dotcommander/zai/internal/config"
)

const (
	maxAudioFileSize = 25 * 1024 * 1024 // 25MB
)

// ClientConfig holds all configuration for the ZAI client.
// Injected at construction time - no global state.
type ClientConfig struct {
	APIKey         string
	BaseURL        string
	CodingBaseURL  string
	Model          string
	Timeout        time.Duration
	Verbose        bool
	RateLimit      RateLimitConfig
	RetryConfig    RetryConfig
	CircuitBreaker config.CircuitBreakerConfig
}

// RateLimitConfig holds rate limiting configuration.
type RateLimitConfig struct {
	RequestsPerSecond int
	Burst             int
}

// DefaultChatOptions returns sensible defaults for CLI usage.
func DefaultChatOptions() ChatOptions {
	return ChatOptions{
		Temperature: Float64Ptr(0.6),
		MaxTokens:   IntPtr(8192),
		TopP:        Float64Ptr(0.9),
		WebEnabled:  BoolPtr(true),
		WebTimeout:  IntPtr(20),
		Think:       false, // Legacy field default
	}
}

// Helper functions for creating pointers to literals (exported for use in cmd package)
func Float64Ptr(v float64) *float64 { return &v }
func IntPtr(v int) *int             { return &v }
func BoolPtr(v bool) *bool          { return &v }

// NewLogger creates a slog.Logger for the application.
// If verbose is true, logs at Debug level; otherwise Info level.
func NewLogger(verbose bool) *slog.Logger {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{Level: level}
	return slog.New(slog.NewTextHandler(os.Stderr, opts))
}

// DiscardLogger returns a logger that discards all output (for testing).
func DiscardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// ChatClient interface for testability (ISP compliance).
// Provides the main chat functionality.
type ChatClient interface {
	Chat(ctx context.Context, prompt string, opts ChatOptions) (string, error)
}

// VisionClient interface for image analysis (ISP compliance).
type VisionClient interface {
	Vision(ctx context.Context, prompt string, imageBase64 string, opts VisionOptions) (string, error)
}

// ImageClient interface for image generation (ISP compliance).
type ImageClient interface {
	GenerateImage(ctx context.Context, prompt string, opts ImageOptions) (*ImageResponse, error)
}

// ModelClient interface for model listing (ISP compliance).
type ModelClient interface {
	ListModels(ctx context.Context) ([]Model, error)
}

// WebReaderClient interface for web content fetching (ISP compliance).
type WebReaderClient interface {
	FetchWebContent(ctx context.Context, url string, opts *WebReaderOptions) (*WebReaderResponse, error)
}

// WebSearchClient interface for web searching (ISP compliance).
type WebSearchClient interface {
	SearchWeb(ctx context.Context, query string, opts SearchOptions) (*WebSearchResponse, error)
}

// AudioClient interface for audio transcription (ISP compliance).
type AudioClient interface {
	TranscribeAudio(ctx context.Context, audioPath string, opts TranscriptionOptions) (*TranscriptionResponse, error)
}

// VideoClient interface for video generation (ISP compliance).
type VideoClient interface {
	GenerateVideo(ctx context.Context, prompt string, opts VideoOptions) (*VideoGenerationResponse, error)
	RetrieveVideoResult(ctx context.Context, taskID string) (*VideoResultResponse, error)
}

// FullClient composes all client interfaces into one (backward compatibility).
type FullClient interface {
	ChatClient
	VisionClient
	ImageClient
	ModelClient
	WebReaderClient
	WebSearchClient
	AudioClient
	VideoClient
}

// Client implements all client interfaces with Z.AI API.

// HistoryStore interface for storage abstraction (ISP compliance).
type HistoryStore interface {
	Save(entry HistoryEntry) error
	GetRecent(limit int) ([]HistoryEntry, error)
}

// HTTPDoer interface for HTTP operations (DIP compliance, enables testing).
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// CircuitBreaker states
type CircuitBreakerState int

const (
	Closed CircuitBreakerState = iota
	Open
	HalfOpen
)

func (s CircuitBreakerState) String() string {
	switch s {
	case Closed:
		return "closed"
	case Open:
		return "open"
	case HalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreaker implements a circuit breaker pattern for API calls.
type CircuitBreaker struct {
	name            string
	config          config.CircuitBreakerConfig
	logger          *slog.Logger
	state           CircuitBreakerState
	mu              sync.Mutex
	failureCount    int
	successCount    int
	lastStateChange time.Time
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(name string, config config.CircuitBreakerConfig, logger *slog.Logger) *CircuitBreaker {
	return &CircuitBreaker{
		name:   name,
		config: config,
		logger: logger,
		state:  Closed,
	}
}

// Execute wraps a function call with circuit breaker protection.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Check if circuit breaker is open
	if cb.state == Open {
		// Check if timeout has passed
		if time.Since(cb.lastStateChange) < cb.config.Timeout {
			return fmt.Errorf("circuit breaker '%s' is open (timeout: %v)", cb.name, cb.config.Timeout)
		}
		// Move to half-open state
		cb.state = HalfOpen
		cb.successCount = 0
		cb.logger.Info("circuit breaker state change",
			"name", cb.name,
			"from", "open",
			"to", "half-open")
	}

	// Execute the function
	err := fn()

	// Record the result
	cb.recordResult(err)

	return err
}

// Reset manually resets the circuit breaker to closed state.
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = Closed
	cb.failureCount = 0
	cb.successCount = 0
	cb.lastStateChange = time.Now()

	cb.logger.Info("circuit breaker reset",
		"name", cb.name,
		"action", "manual reset")
}

// recordResult records the success/failure and updates state accordingly.
func (cb *CircuitBreaker) recordResult(err error) {
	switch cb.state {
	case Closed:
		if err != nil {
			cb.failureCount++
			if cb.failureCount >= cb.config.FailureThreshold {
				cb.setState(Open, err)
			}
		} else {
			cb.failureCount = 0
			cb.successCount = 0
		}

	case HalfOpen:
		if err == nil {
			cb.successCount++
			if cb.successCount >= cb.config.SuccessThreshold {
				cb.setState(Closed, nil)
			}
		} else {
			cb.setState(Open, err)
		}

	case Open:
		// In open state, do nothing until timeout
	}
}

// setState changes the circuit breaker state and logs the transition.
func (cb *CircuitBreaker) setState(newState CircuitBreakerState, err error) {
	if cb.state != newState {
		oldState := cb.state
		cb.state = newState
		cb.lastStateChange = time.Now()

		// Reset counters when transitioning to closed
		if newState == Closed {
			cb.failureCount = 0
			cb.successCount = 0
		}

		var message string
		if err != nil {
			message = fmt.Sprintf("error: %v", err)
		}

		cb.logger.Info("circuit breaker state change",
			"name", cb.name,
			"from", oldState.String(),
			"to", newState.String(),
			"reason", message)
	}
}

// RateLimitedClient implements HTTPDoer with rate limiting.
type RateLimitedClient struct {
	client  HTTPDoer
	limiter *rate.Limiter
	logger  *slog.Logger
}

// NewRateLimitedClient creates a new rate-limited HTTP client.
func NewRateLimitedClient(client HTTPDoer, rateLimitConfig RateLimitConfig, logger *slog.Logger) HTTPDoer {
	if rateLimitConfig.RequestsPerSecond <= 0 {
		// Rate limiting disabled, return original client
		return client
	}

	limiter := rate.NewLimiter(rate.Limit(rateLimitConfig.RequestsPerSecond), rateLimitConfig.Burst)
	return &RateLimitedClient{
		client:  client,
		limiter: limiter,
		logger:  logger,
	}
}

// Do implements HTTPDoer interface with rate limiting.
func (c *RateLimitedClient) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	// Wait for token from rate limiter
	err := c.limiter.Wait(ctx)
	if err != nil {
		c.logger.Error("rate limit exceeded", "error", err)
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	return c.client.Do(req)
}

// FileReader interface for file operations (DIP compliance, enables testing).
// Deprecated: Use utils.FileReader instead. Kept for backward compatibility.
type FileReader = utils.FileReader

// OSFileReader implements FileReader using os.ReadFile.
// Deprecated: Use utils.OSFileReader instead. Kept for backward compatibility.
type OSFileReader = utils.OSFileReader

// Client implements ChatClient with Z.AI API.
type Client struct {
	config          ClientConfig
	httpClient      HTTPDoer
	logger          *slog.Logger
	history         HistoryStore
	fileReader      FileReader
	circuitBreakers map[string]*CircuitBreaker
	mu              sync.RWMutex
}

// ClientDeps holds optional dependencies for NewClient.
// Zero values mean "use default implementation".
type ClientDeps struct {
	HTTPClient HTTPDoer
	FileReader FileReader
}

// NewClient creates a client with injected dependencies.
// If deps is nil or fields are nil, default implementations are used.
func NewClient(cfg ClientConfig, logger *slog.Logger, history HistoryStore, httpClient HTTPDoer) *Client {
	return NewClientWithDeps(cfg, logger, history, &ClientDeps{HTTPClient: httpClient})
}

// NewClientWithDeps creates a client with full dependency injection.
// Allows injection of all dependencies for testing.
func NewClientWithDeps(cfg ClientConfig, logger *slog.Logger, history HistoryStore, deps *ClientDeps) *Client {
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

	// Wrap HTTP client with rate limiting
	httpClient = NewRateLimitedClient(httpClient, cfg.RateLimit, logger)

	client := &Client{
		config:          cfg,
		httpClient:      httpClient,
		logger:          logger,
		history:         history,
		fileReader:      fileReader,
		circuitBreakers: make(map[string]*CircuitBreaker),
	}

	// Initialize circuit breakers
	if cfg.CircuitBreaker.Enabled {
		client.initCircuitBreakers()
	}

	return client
}

// HTTPClient returns the underlying HTTP client for connection reuse.
func (c *Client) HTTPClient() HTTPDoer {
	return c.httpClient
}

// initCircuitBreakers initializes circuit breakers for different API endpoints.
func (c *Client) initCircuitBreakers() {
	c.circuitBreakers["chat"] = NewCircuitBreaker("chat", c.config.CircuitBreaker, c.logger)
	c.circuitBreakers["web_search"] = NewCircuitBreaker("web_search", c.config.CircuitBreaker, c.logger)
	c.circuitBreakers["reader"] = NewCircuitBreaker("reader", c.config.CircuitBreaker, c.logger)
	c.circuitBreakers["models"] = NewCircuitBreaker("models", c.config.CircuitBreaker, c.logger)
	c.circuitBreakers["images"] = NewCircuitBreaker("images", c.config.CircuitBreaker, c.logger)
	c.circuitBreakers["videos"] = NewCircuitBreaker("videos", c.config.CircuitBreaker, c.logger)
}

// getCircuitBreaker returns the appropriate circuit breaker for an endpoint.
func (c *Client) getCircuitBreaker(endpoint string) *CircuitBreaker {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.circuitBreakers[endpoint]
}

// requireAPIKey validates the API key is configured.
// Returns an error with helpful message if not set.
func (c *Client) requireAPIKey() error {
	if c.config.APIKey == "" {
		return fmt.Errorf("API key is not configured. Set ZAI_API_KEY or configure in ~/.config/zai/config.yaml")
	}
	return nil
}

// Chat sends a prompt and returns the response.
// Orchestrates content building, URL enrichment, and request execution.
func (c *Client) Chat(ctx context.Context, prompt string, opts ChatOptions) (string, error) {
	if err := c.requireAPIKey(); err != nil {
		return "", err
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
				c.logger.Warn("failed to fetch web content", "url", url, "error", err)
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
		c.logger.Warn("error fetching web content", "error", err)
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
		Timeout:           timeout,
		ReturnFormat:      "markdown",
		RetainImages:      &trueVal,
		NoCache:           &falseVal,
		NoGFM:             &falseVal,
		KeepImgDataURL:    &falseVal,
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
		c.logger.Warn("failed to save to history", "error", err)
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

	// Add system prompt (custom or default)
	systemPrompt := opts.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = "Be concise and direct. Answer briefly and to the point."
	}
	messages = append(messages, Message{
		Role:    "system",
		Content: systemPrompt,
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
		"429", // Too Many Requests (rate limit)
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
	// Cap attempt to prevent overflow (2^62 would overflow time.Duration)
	if attempt > 62 {
		attempt = 62
	}

	// Exponential backoff: initial * 2^(attempt-1)
	backoff := initialBackoff * time.Duration(1<<uint(attempt-1)) //nolint:gosec // G115: attempt count is small, overflow impossible

	// Cap at max backoff
	if backoff > maxBackoff {
		backoff = maxBackoff
	}

	// Add jitter (±12.5%, centered - so jitter can add or subtract up to 12.5%)
	// This ensures we never go below the base value by more than 12.5%
	jitterRange := float64(backoff) * 0.125
	jitter := time.Duration(jitterRange * (2.0*rand.Float64() - 1.0)) //nolint:gosec // G404: jitter doesn't need crypto-grade randomness

	return backoff + jitter
}

// buildJSONRequest creates an HTTP POST request with JSON data.
func buildJSONRequest(baseURL, apiKey string, ctx context.Context, endpoint string, data interface{}) (*http.Request, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/%s", baseURL, endpoint)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	setJSONHeaders(req, apiKey)
	return req, nil
}

// buildGetRequest creates an HTTP GET request.
func buildGetRequest(baseURL, apiKey string, ctx context.Context, endpoint string) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", baseURL, endpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	return req, nil
}

// setJSONHeaders sets common headers for JSON requests.
func setJSONHeaders(req *http.Request, apiKey string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Accept-Language", "en-US,en")
}

// extractEndpointName extracts a standardized name from endpoint path.
func extractEndpointName(endpoint string) string {
	switch {
	case strings.Contains(endpoint, "chat"):
		return "chat"
	case strings.Contains(endpoint, "web_search"):
		return "web_search"
	case strings.Contains(endpoint, "reader"):
		return "reader"
	case strings.Contains(endpoint, "models"):
		return "models"
	case strings.Contains(endpoint, "images"):
		return "images"
	case strings.Contains(endpoint, "videos"):
		return "videos"
	case strings.Contains(endpoint, "audio"):
		return "audio"
	default:
		return "default"
	}
}

// closeBody closes the response body and logs any error.
func closeBody(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to close response body: %v\n", err)
	}
}

// executeJSONRequest executes a JSON POST request using HTTPDoer interface.
func (c *Client) executeJSONRequest(ctx context.Context, endpoint string, reqData interface{}) ([]byte, error) {
	if c.config.CircuitBreaker.Enabled {
		cb := c.getCircuitBreaker(extractEndpointName(endpoint))
		if cb != nil {
			var result []byte
			var internalErr error
			err := cb.Execute(func() error {
				result, internalErr = c.executeJSONRequestInternal(ctx, endpoint, reqData)
				return internalErr
			})
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	}

	return c.executeJSONRequestInternal(ctx, endpoint, reqData)
}

// executeJSONRequestInternal is the internal implementation without circuit breaker.
func (c *Client) executeJSONRequestInternal(ctx context.Context, endpoint string, reqData interface{}) ([]byte, error) {
	req, err := buildJSONRequest(c.config.BaseURL, c.config.APIKey, ctx, endpoint, reqData)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("sending request", "url", req.URL)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer closeBody(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	return body, nil
}

// executeGetRequest executes a GET request using HTTPDoer interface.
func (c *Client) executeGetRequest(ctx context.Context, endpoint string) ([]byte, error) {
	if c.config.CircuitBreaker.Enabled {
		cb := c.getCircuitBreaker(extractEndpointName(endpoint))
		if cb != nil {
			var result []byte
			var internalErr error
			err := cb.Execute(func() error {
				result, internalErr = c.executeGetRequestInternal(ctx, endpoint)
				return internalErr
			})
			return result, err
		}
	}

	return c.executeGetRequestInternal(ctx, endpoint)
}

// executeGetRequestInternal is the internal implementation without circuit breaker.
func (c *Client) executeGetRequestInternal(ctx context.Context, endpoint string) ([]byte, error) {
	req, err := buildGetRequest(c.config.BaseURL, c.config.APIKey, ctx, endpoint)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("sending request", "url", req.URL)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer closeBody(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	return body, nil
}

// doRequest executes the HTTP request to Z.AI API.
// Single place for all HTTP logic (DRY compliance).
func (c *Client) doRequest(ctx context.Context, messages []Message, opts ChatOptions) (string, Usage, error) {
	// Use opts.Thinking (bool pointer) to build the API request structure
	var thinking *Thinking
	if opts.Thinking != nil && *opts.Thinking {
		thinking = &Thinking{Type: "enabled"}
	} else {
		thinking = &Thinking{Type: "disabled"}
	}

	reqData := ChatRequest{
		Model:    c.config.Model,
		Messages: messages,
		Stream:   false,
		Thinking: thinking,
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

	c.logger.Debug("sending request", "url", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", Usage{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer closeBody(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", Usage{}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", Usage{}, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", Usage{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", Usage{}, fmt.Errorf("no choices in response")
	}

	c.logger.Debug("usage",
		"total_tokens", chatResp.Usage.TotalTokens,
		"prompt_tokens", chatResp.Usage.PromptTokens,
		"completion_tokens", chatResp.Usage.CompletionTokens)

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
			c.logger.Debug("retrying request",
				"attempt", attempt,
				"max_attempts", maxAttempts,
				"backoff", backoff,
				"error", lastErr)

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
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	var modelsResp ModelsResponse
	body, err := c.executeGetRequest(ctx, "models")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &modelsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal models response: %w", err)
	}

	return modelsResp.Data, nil
}

// GenerateImage creates an image using the Z.AI image generation API.
func (c *Client) GenerateImage(ctx context.Context, prompt string, opts ImageOptions) (*ImageResponse, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	// Validate options
	if err := validateImageOptions(opts); err != nil {
		return nil, fmt.Errorf("invalid image options: %w", err)
	}

	// Build request with defaults
	model := opts.Model
	if model == "" {
		model = "glm-image" // Default image model
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

	var imageResp ImageResponse
	body, err := c.executeJSONRequest(ctx, "images/generations", reqData)
	if err != nil {
		return nil, fmt.Errorf("image generation API error: %w", err)
	}
	if err := json.Unmarshal(body, &imageResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal image response: %w", err)
	}

	if len(imageResp.Data) == 0 {
		return nil, fmt.Errorf("no images in response")
	}

	c.logger.Debug("generated image",
		"url", imageResp.Data[0].URL,
		"width", imageResp.Data[0].Width,
		"height", imageResp.Data[0].Height)

	return &imageResp, nil
}

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

	c.logger.Debug("fetched web content",
		"url", webResp.ReaderResult.URL,
		"title", webResp.ReaderResult.Title)

	return &webResp, nil
}

// validateWebContentURL validates the URL parameter for web content fetching.
func (c *Client) validateWebContentURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL is required")
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
func (c *Client) SearchWeb(ctx context.Context, query string, opts SearchOptions) (*WebSearchResponse, error) { //nolint:gocognit,gocyclo
	if err := c.requireAPIKey(); err != nil {
		return nil, err
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

	c.logger.Debug("search complete", "results", len(searchResp.SearchResult), "query", query)

	// Save to history (non-blocking, log errors)
	if c.history != nil {
		entry := NewSearchHistoryEntry(time.Now(), query, &searchResp)
		if err := c.history.Save(entry); err != nil {
			c.logger.Warn("failed to save search to history", "error", err)
		}
	}

	return &searchResp, nil
}

// Vision analyzes an image using Z.AI's vision model (glm-4.6v).
// imageBase64 should be a data URI like "data:image/jpeg;base64,<base64-data>" or a raw base64 string.
func (c *Client) Vision(ctx context.Context, prompt string, imageBase64 string, opts VisionOptions) (string, error) {
	if err := c.requireAPIKey(); err != nil {
		return "", err
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

	var chatResp ChatResponse
	body, err := c.executeJSONRequest(ctx, "chat/completions", reqData)
	if err != nil {
		return "", fmt.Errorf("vision API error: %w", err)
	}
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal vision response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in vision response")
	}

	c.logger.Debug("vision complete",
		"total_tokens", chatResp.Usage.TotalTokens,
		"prompt_tokens", chatResp.Usage.PromptTokens,
		"completion_tokens", chatResp.Usage.CompletionTokens)

	return chatResp.Choices[0].Message.Content, nil
}

// TranscribeAudio transcribes an audio file using Z.AI's ASR model.
func (c *Client) TranscribeAudio(ctx context.Context, audioPath string, opts TranscriptionOptions) (*TranscriptionResponse, error) { //nolint:gocyclo,funlen
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	// Validate audio file
	if audioPath == "" {
		return nil, fmt.Errorf("audio file path is required")
	}

	// Read audio file using injected FileReader
	data, err := c.fileReader.ReadFile(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio file: %w", err)
	}

	// Check file size
	if len(data) > maxAudioFileSize {
		return nil, fmt.Errorf("audio file too large: %d bytes (max: %d MB)", len(data), maxAudioFileSize/1024/1024)
	}

	// Build model
	model := opts.Model
	if model == "" {
		model = "glm-asr-2512"
	}

	// Build multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file from memory
	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(data)); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	// Add model
	writer.WriteField("model", model) //nolint:errcheck // multipart field write

	// Add optional fields
	if opts.Prompt != "" {
		writer.WriteField("prompt", opts.Prompt) //nolint:errcheck // multipart field write
	}
	if opts.Stream {
		writer.WriteField("stream", "true") //nolint:errcheck // multipart field write
	}
	if opts.UserID != "" {
		writer.WriteField("user_id", opts.UserID) //nolint:errcheck // multipart field write
	}
	if opts.RequestID != "" {
		writer.WriteField("request_id", opts.RequestID) //nolint:errcheck // multipart field write
	}
	if len(opts.Hotwords) > 0 {
		hotwordsJSON, err := json.Marshal(opts.Hotwords)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal hotwords: %w", err)
		}
		writer.WriteField("hotwords", string(hotwordsJSON)) //nolint:errcheck // multipart field write
	}

	writer.Close() //nolint:errcheck // multipart writer close

	// Build request
	url := fmt.Sprintf("%s/audio/transcriptions", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Accept-Language", "en-US,en")

	c.logger.Debug("sending audio transcription request", "url", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer closeBody(resp)

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

	c.logger.Debug("transcription complete", "chars", len(transcriptionResp.Text), "model", transcriptionResp.Model)

	return &transcriptionResp, nil
}

// GenerateVideo creates a video using Z.AI's CogVideoX-3 API (async).
func (c *Client) GenerateVideo(ctx context.Context, prompt string, opts VideoOptions) (*VideoGenerationResponse, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	// Build model
	model := opts.Model
	if model == "" {
		model = "cogvideox-3" // Default video model
	}

	// Validate options
	if err := validateVideoOptions(opts); err != nil {
		return nil, fmt.Errorf("invalid video options: %w", err)
	}

	// Build request
	reqData := VideoGenerationRequest{
		Model:     model,
		Prompt:    prompt,
		ImageURL:  opts.ImageURLs,
		Quality:   opts.Quality,
		WithAudio: opts.WithAudio,
		Size:      opts.Size,
		FPS:       opts.FPS,
		Duration:  opts.Duration,
		RequestID: opts.RequestID,
		UserID:    opts.UserID,
	}

	// Set defaults
	if reqData.Quality == "" {
		reqData.Quality = "speed"
	}
	if reqData.Size == "" {
		reqData.Size = "1920x1080"
	}
	if reqData.FPS == 0 {
		reqData.FPS = 30
	}
	if reqData.Duration == 0 {
		reqData.Duration = 5
	}

	var videoResp VideoGenerationResponse
	body, err := c.executeJSONRequest(ctx, "videos/generations", reqData)
	if err != nil {
		return nil, fmt.Errorf("video generation API error: %w", err)
	}
	if err := json.Unmarshal(body, &videoResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal video response: %w", err)
	}

	c.logger.Debug("video generation task created", "id", videoResp.ID, "status", videoResp.TaskStatus)

	return &videoResp, nil
}

// RetrieveVideoResult polls for async video generation result.
func (c *Client) RetrieveVideoResult(ctx context.Context, taskID string) (*VideoResultResponse, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	// Validate task ID
	if taskID == "" {
		return nil, fmt.Errorf("task ID is required")
	}

	var resultResp VideoResultResponse
	endpoint := fmt.Sprintf("async-result/%s", taskID)
	body, err := c.executeGetRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("retrieve video result API error: %w", err)
	}
	if err := json.Unmarshal(body, &resultResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal video result response: %w", err)
	}

	c.logger.Debug("video result retrieved", "id", taskID, "status", resultResp.TaskStatus)

	return &resultResp, nil
}

// validateVideoOptions checks if video options are valid.
func validateVideoOptions(opts VideoOptions) error {
	// Validate quality
	if opts.Quality != "" && opts.Quality != "quality" && opts.Quality != "speed" {
		return fmt.Errorf("invalid quality: %s (must be 'quality' or 'speed')", opts.Quality)
	}

	// Validate size format
	if opts.Size != "" {
		supportedSizes := map[string]bool{
			"1280x720": true, "720x1280": true,
			"1024x1024": true,
			"1920x1080": true, "1080x1920": true,
			"2048x1080": true,
			"3840x2160": true,
		}
		if !supportedSizes[opts.Size] {
			return fmt.Errorf("invalid size: %s (supported: 1280x720, 720x1280, 1024x1024, 1920x1080, 1080x1920, 2048x1080, 3840x2160)", opts.Size)
		}
	}

	// Validate FPS
	if opts.FPS != 0 && opts.FPS != 30 && opts.FPS != 60 {
		return fmt.Errorf("invalid fps: %d (must be 30 or 60)", opts.FPS)
	}

	// Validate duration
	if opts.Duration != 0 && opts.Duration != 5 && opts.Duration != 10 {
		return fmt.Errorf("invalid duration: %d (must be 5 or 10 seconds)", opts.Duration)
	}

	// Validate image URLs (max 2 for first/last frame mode)
	if len(opts.ImageURLs) > 2 {
		return fmt.Errorf("too many image URLs (max 2 for first/last frame mode)")
	}

	return nil
}
