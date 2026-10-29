package app

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

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

	// Execute billable generation once; callers decide whether to submit again.
	response, usage, err := c.doRequest(ctx, messages, opts)
	if err != nil {
		return "", err
	}

	// Save to history (non-blocking, log errors)
	c.saveToHistory(prompt, response, usage, c.buildChatRequest(nil, opts, false).Model)

	return response, nil
}

// StreamReader reads SSE chunks from a streaming API response.
type StreamReader struct {
	scanner     *bufio.Scanner
	body        io.ReadCloser
	usage       Usage
	fullContent strings.Builder
	done        bool
	model       string
}

// Next reads the next token from the stream.
// Returns the delta content string and nil on success.
// Returns "" and io.EOF when the stream is complete.
func (sr *StreamReader) Next() (string, error) {
	if sr.done {
		return "", io.EOF
	}

	for sr.scanner.Scan() {
		line := sr.scanner.Text()

		if line == "" {
			continue
		}

		if line == "data: [DONE]" {
			sr.done = true
			return "", io.EOF
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if chunk.Usage != nil {
			sr.usage = *chunk.Usage
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		content := chunk.Choices[0].Delta.Content
		if content != "" {
			sr.fullContent.WriteString(content)
			return content, nil
		}
	}

	if err := sr.scanner.Err(); err != nil {
		return "", fmt.Errorf("stream read error: %w", err)
	}

	sr.done = true
	return "", io.EOF
}

// Close closes the underlying response body.
func (sr *StreamReader) Close() error {
	return sr.body.Close()
}

// FullContent returns all accumulated content from the stream.
func (sr *StreamReader) Model() string { return sr.model }

func (sr *StreamReader) FullContent() string {
	return sr.fullContent.String()
}

// StreamUsage returns the token usage from the stream (available after completion).
func (sr *StreamReader) StreamUsage() Usage {
	return sr.usage
}

// StreamChat sends a streaming chat request and returns a StreamReader.
// The caller must call Close() on the returned StreamReader when done.
func (c *Client) StreamChat(ctx context.Context, prompt string, opts ChatOptions) (*StreamReader, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	content, err := c.buildContent(ctx, prompt, opts.FilePath)
	if err != nil {
		return nil, err
	}

	content = c.enrichWithURLContent(ctx, prompt, content, opts)
	messages := c.buildMessagesWithContext(content, opts)

	return c.doStreamRequest(ctx, messages, opts)
}

// SaveStreamToHistory saves a completed stream exchange to history.
func (c *Client) SaveStreamToHistory(prompt, response string, usage Usage, models ...string) {
	c.saveToHistory(prompt, response, usage, models...)
}

// buildChatRequest constructs a ChatRequest from messages and options.
// Handles the legacy Think→Thinking bridge and applies all defaults.
func (c *Client) buildChatRequest(messages []Message, opts ChatOptions, stream bool) ChatRequest {
	// Handle legacy Think field
	if opts.Think && opts.Thinking == nil {
		opts.Thinking = &opts.Think
	}

	var thinking *Thinking
	if opts.Thinking != nil && *opts.Thinking {
		thinking = &Thinking{Type: "enabled"}
	} else {
		thinking = &Thinking{Type: "disabled"}
	}

	req := ChatRequest{
		Model:    resolveModel(opts.Model, c.config.Model, "glm-4.7"),
		Messages: messages,
		Stream:   stream,
		Thinking: thinking,
	}

	if opts.Temperature != nil {
		req.Temperature = *opts.Temperature
	} else {
		req.Temperature = 0.6
	}

	if opts.MaxTokens != nil {
		req.MaxTokens = *opts.MaxTokens
	} else {
		req.MaxTokens = 8192
	}

	if opts.TopP != nil {
		req.TopP = *opts.TopP
	} else {
		req.TopP = 0.9
	}

	return req
}

// doStreamRequest executes the streaming HTTP request to Z.AI API.
func (c *Client) doStreamRequest(ctx context.Context, messages []Message, opts ChatOptions) (*StreamReader, error) {
	reqData := c.buildChatRequest(messages, opts, true)

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	baseURL := c.config.BaseURL
	if c.config.UseCoding {
		baseURL = c.config.CodingBaseURL
	}

	url := fmt.Sprintf("%s/chat/completions", baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Accept-Language", "en-US,en")

	c.logger.DebugContext(ctx, "sending streaming request", "url", url)

	resp, err := c.doChatHTTP(req, true)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer closeBody(resp)
		return nil, providerError(resp)
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	return &StreamReader{
		scanner: scanner,
		body:    resp.Body,
		model:   reqData.Model,
	}, nil
}

// doRequest executes the HTTP request to Z.AI API.
// Single place for all HTTP logic (DRY compliance).
func (c *Client) doRequest(ctx context.Context, messages []Message, opts ChatOptions) (string, Usage, error) {
	reqData := c.buildChatRequest(messages, opts, false)

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", Usage{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	baseURL := c.config.BaseURL
	if c.config.UseCoding {
		baseURL = c.config.CodingBaseURL
	}

	url := fmt.Sprintf("%s/chat/completions", baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", Usage{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Accept-Language", "en-US,en")

	c.logger.DebugContext(ctx, "sending request", "url", url)

	resp, err := c.doChatHTTP(req, false)
	if err != nil {
		return "", Usage{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		return "", Usage{}, providerError(resp)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", Usage{}, fmt.Errorf("failed to read response: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", Usage{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", Usage{}, errors.New("no choices in response")
	}

	c.logger.DebugContext(ctx, "usage",
		"total_tokens", chatResp.Usage.TotalTokens,
		"prompt_tokens", chatResp.Usage.PromptTokens,
		"completion_tokens", chatResp.Usage.CompletionTokens)

	return chatResp.Choices[0].Message.Content, chatResp.Usage, nil
}

// doRequestWithRetry retains its internal compatibility name but submits generation once.
// An ambiguous provider failure must not automatically replay billable work.
func (c *Client) doRequestWithRetry(ctx context.Context, messages []Message, opts ChatOptions) (string, Usage, error) {
	return c.doRequest(ctx, messages, opts)
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
		g.Go(func() error {
			webResp, err := c.FetchWebContent(ctx, url, webOpts)
			if err != nil {
				c.logger.WarnContext(ctx, "failed to fetch web content", "url", url, "error", err)
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
		c.logger.WarnContext(ctx, "error fetching web content", "error", err)
	}

	// Append results in original order
	var sb strings.Builder
	sb.WriteString(content)
	for _, r := range results {
		if r.url != "" { // Only append successful fetches
			sb.WriteString("\n\n")
			sb.WriteString(FormatWebContent(r.url, r.title, r.body))
		}
	}
	content = sb.String()

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
	return c.buildMessages(content, opts)
}

// saveToHistory persists the chat exchange to history storage.
func (c *Client) saveToHistory(prompt, response string, usage Usage, models ...string) {
	if c.history == nil {
		return
	}
	model := c.config.Model
	if len(models) > 0 && models[0] != "" {
		model = models[0]
	}
	entry := NewChatHistoryEntry(time.Now(), prompt, response, model, usage)
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
	messages := make([]Message, 0, len(opts.Context)+2)

	// Add system prompt (custom or default)
	systemPrompt := opts.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = "Be concise and direct. Answer briefly and to the point."
	}
	messages = append(messages, Message{
		Role:    "system",
		Content: systemPrompt,
	})

	messages = append(messages, opts.Context...)

	// Add current user message
	messages = append(messages, Message{
		Role:    "user",
		Content: content,
	})

	return messages
}

// doChatHTTP applies the configured breaker to chat request and response headers.
func (c *Client) doChatHTTP(req *http.Request, stream bool) (*http.Response, error) {
	var resp *http.Response
	run := func() error {
		client := c.httpClient
		if stream {
			client = c.streamHTTPClient
		}
		var err error
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return &APIError{StatusCode: resp.StatusCode}
		}
		return nil
	}
	var err error
	if breaker := c.getCircuitBreaker("chat"); breaker != nil {
		err = breaker.Execute(run)
	} else {
		err = run()
	}
	if resp != nil {
		return resp, nil
	}
	return nil, err
}
