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
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/dotcommander/zai/internal/app/fileutil"
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
	UseCoding      bool
	Model          string
	ImageModel     string
	VisionModel    string
	AudioModel     string
	TTSModel       string
	EmbeddingModel string
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

// Float64Ptr returns a pointer to the given float64 value.
func Float64Ptr(v float64) *float64 { return &v }

// IntPtr returns a pointer to the given int value.
func IntPtr(v int) *int { return &v }

// BoolPtr returns a pointer to the given bool value.
func BoolPtr(v bool) *bool { return &v }

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

// TTSClient interface for text-to-speech (ISP compliance).
type TTSClient interface {
	TextToSpeech(ctx context.Context, text string, opts TTSOptions) ([]byte, error)
}

// EmbeddingClient interface for text embeddings (ISP compliance).
type EmbeddingClient interface {
	CreateEmbedding(ctx context.Context, texts []string, opts EmbeddingOptions) (*EmbeddingResponse, error)
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

// CircuitBreakerState represents the state of the circuit breaker.
type CircuitBreakerState int

// Circuit breaker state constants.
const (
	// Closed is the normal operating state; requests pass through.
	Closed CircuitBreakerState = iota
	// Open is the tripped state; requests are rejected.
	Open
	// HalfOpen is the recovery state; a limited number of requests are allowed.
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
//
// Deprecated: Use fileutil.FileReader instead. Kept for backward compatibility.
type FileReader = fileutil.FileReader

// OSFileReader implements FileReader using os.ReadFile.
//
// Deprecated: Use fileutil.OSFileReader instead. Kept for backward compatibility.
type OSFileReader = fileutil.OSFileReader

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
	c.circuitBreakers["tts"] = NewCircuitBreaker("tts", c.config.CircuitBreaker, c.logger)
	c.circuitBreakers["embeddings"] = NewCircuitBreaker("embeddings", c.config.CircuitBreaker, c.logger)
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
		return errors.New("API key is not configured. Set ZAI_API_KEY or configure in ~/.config/zai/config.yaml")
	}
	return nil
}

// closeBody closes the response body and logs any error.
func closeBody(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to close response body: %v\n", err)
	}
}

// setJSONHeaders sets common headers for JSON requests.
func setJSONHeaders(req *http.Request, apiKey string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Accept-Language", "en-US,en")
}

// buildJSONRequest creates an HTTP POST request with JSON data.
func buildJSONRequest(ctx context.Context, baseURL, apiKey, endpoint string, data interface{}) (*http.Request, error) {
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
func buildGetRequest(ctx context.Context, baseURL, apiKey, endpoint string) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", baseURL, endpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	return req, nil
}

// withCircuitBreaker wraps fn with the circuit breaker for the given endpoint,
// if circuit breaking is enabled and a breaker exists. Otherwise calls fn directly.
func (c *Client) withCircuitBreaker(endpoint string, fn func() ([]byte, error)) ([]byte, error) {
	if !c.config.CircuitBreaker.Enabled {
		return fn()
	}
	cb := c.getCircuitBreaker(extractEndpointName(endpoint))
	if cb == nil {
		return fn()
	}
	var result []byte
	err := cb.Execute(func() error {
		var innerErr error
		result, innerErr = fn()
		return innerErr
	})
	return result, err
}

// executeJSONRequest executes a JSON POST request using HTTPDoer interface.
func (c *Client) executeJSONRequest(ctx context.Context, endpoint string, reqData interface{}) ([]byte, error) {
	return c.withCircuitBreaker(endpoint, func() ([]byte, error) {
		return c.executeJSONRequestInternal(ctx, endpoint, reqData)
	})
}

// executeJSONRequestInternal is the internal implementation without circuit breaker.
func (c *Client) executeJSONRequestInternal(ctx context.Context, endpoint string, reqData interface{}) ([]byte, error) {
	baseURL := c.config.BaseURL
	if c.config.UseCoding && endpoint == "chat/completions" {
		baseURL = c.config.CodingBaseURL
	}
	req, err := buildJSONRequest(ctx, baseURL, c.config.APIKey, endpoint, reqData)
	if err != nil {
		return nil, err
	}

	c.logger.DebugContext(ctx, "sending request", "url", req.URL)
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
	return c.withCircuitBreaker(endpoint, func() ([]byte, error) {
		return c.executeGetRequestInternal(ctx, endpoint)
	})
}

// executeGetRequestInternal is the internal implementation without circuit breaker.
func (c *Client) executeGetRequestInternal(ctx context.Context, endpoint string) ([]byte, error) {
	req, err := buildGetRequest(ctx, c.config.BaseURL, c.config.APIKey, endpoint)
	if err != nil {
		return nil, err
	}

	c.logger.DebugContext(ctx, "sending request", "url", req.URL)
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

// resolveModel returns the first non-empty value among override, configDefault, fallback.
func resolveModel(override, configDefault, fallback string) string {
	if override != "" {
		return override
	}
	if configDefault != "" {
		return configDefault
	}
	return fallback
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
	case strings.Contains(endpoint, "audio/speech"), strings.Contains(endpoint, "tts"), strings.Contains(endpoint, "voice"):
		return "tts"
	case strings.Contains(endpoint, "audio"):
		return "audio"
	case strings.Contains(endpoint, "embeddings"):
		return "embeddings"
	default:
		return "default"
	}
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
