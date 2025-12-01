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

// ChatOptions configures individual chat requests.
type ChatOptions struct {
	Temperature  float64
	MaxTokens    int
	TopP         float64
	SystemPrompt string
	FilePath     string        // Optional file to include in context
	Context      []Message     // Previous messages for context
}

// DefaultChatOptions returns sensible defaults for CLI usage.
func DefaultChatOptions() ChatOptions {
	return ChatOptions{
		Temperature:  0.2,
		MaxTokens:    2000,
		TopP:         0.9,
		SystemPrompt: "Be concise and direct. Answer briefly and to the point.",
	}
}

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
// Single method handles all cases: simple, with file, with context.
func (c *Client) Chat(ctx context.Context, prompt string, opts ChatOptions) (string, error) {
	if c.config.APIKey == "" {
		return "", fmt.Errorf("API key is not configured")
	}

	// Build message content (with optional file)
	content, err := c.buildContent(prompt, opts.FilePath)
	if err != nil {
		return "", err
	}

	// Build messages array
	messages := c.buildMessages(content, opts)

	// Execute request
	response, usage, err := c.doRequest(ctx, messages, opts)
	if err != nil {
		return "", err
	}

	// Save to history (non-blocking, log errors)
	if c.history != nil {
		entry := HistoryEntry{
			Timestamp:  time.Now(),
			Prompt:     prompt,
			Response:   response,
			Model:      c.config.Model,
			TokenUsage: usage,
		}
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

	// Add context messages first (previous conversation)
	messages = append(messages, opts.Context...)

	// Add system prompt if no context (single-shot mode)
	if len(opts.Context) == 0 && opts.SystemPrompt != "" {
		messages = append(messages, Message{
			Role:    "system",
			Content: opts.SystemPrompt,
		})
	}

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
		Temperature: opts.Temperature,
		MaxTokens:   opts.MaxTokens,
		TopP:        opts.TopP,
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
