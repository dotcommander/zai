// Package mcp implements the session-based MCP transports used by Z.AI's
// Coding Plan products. It does not execute tools requested by a server.
package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

const maxResponse = 64 << 20
const maxLine = 8 << 20

// HTTPDoer permits reuse of a caller's transport. Custom implementations must
// not forward Authorization across origins when following redirects.
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type Options struct {
	HTTPClient HTTPDoer
	URL        string
	APIKey     string
	Command    string
	Args       []string
	// Env entries override inherited environment variables for the local server.
	Env []string
}

// RPCError preserves the complete server error object without decoding numbers.
type RPCError struct {
	Raw json.RawMessage
	key string
}

func (e *RPCError) Error() string {
	text := string(e.Raw)
	if len(text) > 65536 {
		text = text[:65536] + " [truncated]"
	}
	if e.key != "" {
		text = strings.ReplaceAll(text, e.key, "[redacted]")
	}
	return "MCP server returned a JSON-RPC error: " + text
}

// HTTPError retains bounded diagnostics. Credentials are redacted.
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string { return fmt.Sprintf("MCP HTTP %d: %s", e.StatusCode, e.Body) }

type Client struct {
	mu                              sync.Mutex
	ctx                             context.Context
	cancel                          context.CancelFunc
	http                            HTTPDoer
	endpoint, key, session, version string
	id                              uint64
	closed                          bool
	command                         *exec.Cmd
	input                           io.WriteCloser
	output                          io.ReadCloser
	scanner                         *bufio.Scanner
	stopPipes                       func() bool
}

type envelope struct {
	Version string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Result  json.RawMessage `json:"result"`
	Error   json.RawMessage `json:"error"`
}

// New starts and initializes one server session. Exactly one of URL or Command
// must be supplied. No requests, including tool calls, are automatically retried.
func New(ctx context.Context, options Options) (*Client, error) {
	if ctx == nil {
		return nil, errors.New("MCP requires a context")
	}
	if (options.URL == "") == (options.Command == "") {
		return nil, errors.New("MCP requires exactly one HTTP URL or local command")
	}
	cctx, cancel := context.WithCancel(ctx)
	c := &Client{ctx: cctx, cancel: cancel, key: options.APIKey, version: "2025-03-26"}
	fail := func(err error) (*Client, error) {
		_ = c.Close() // Initialization errors take precedence over cleanup errors.
		return nil, err
	}
	if options.URL != "" {
		endpoint, err := url.Parse(options.URL)
		if err != nil || endpoint.Host == "" || (endpoint.Scheme != "http" && endpoint.Scheme != "https") || endpoint.User != nil || endpoint.Fragment != "" || endpoint.RawQuery != "" {
			return fail(errors.New("MCP server URL must be absolute HTTP(S) without credentials, query, or fragment"))
		}
		c.endpoint = endpoint.String()
		c.http = options.HTTPClient
		if c.http == nil {
			c.http = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
		}
		if h, ok := c.http.(*http.Client); ok {
			clone := *h
			clone.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
			c.http = &clone
		}
	} else {
		command := exec.CommandContext(cctx, options.Command, options.Args...) // #nosec G204 -- Callers explicitly select the local MCP executable and arguments; no shell is used.
		configureProcess(command)
		command.Env = append(os.Environ(), options.Env...)
		command.Stderr = io.Discard
		command.WaitDelay = 2 * time.Second
		input, err := command.StdinPipe()
		if err != nil {
			return fail(errors.New("MCP could not open server stdin"))
		}
		output, err := command.StdoutPipe()
		if err != nil {
			_ = input.Close() // Best-effort cleanup after setup failed.
			return fail(errors.New("MCP could not open server stdout"))
		}
		c.command = command
		c.input = input
		c.output = output
		if err = command.Start(); err != nil {
			_ = input.Close()  // Best-effort cleanup after setup failed.
			_ = output.Close() // Best-effort cleanup after setup failed.
			c.command = nil
			return fail(errors.New("MCP could not start local server; check the executable and runtime requirements"))
		}
		c.stopPipes = context.AfterFunc(cctx, func() {
			// Cancellation closes pipes to unblock I/O; they may already be closed.
			_ = input.Close()
			_ = output.Close()
		})
		c.scanner = bufio.NewScanner(output)
		c.scanner.Buffer(make([]byte, 4096), maxLine)
	}
	init := json.RawMessage(`{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"zai","version":"1"}}`)
	raw, err := c.Call(ctx, "initialize", init)
	if err != nil {
		return fail(err)
	}
	var result struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	if err = json.Unmarshal(raw, &result); err != nil {
		return fail(errors.New("MCP initialize result is not an object"))
	}
	switch result.ProtocolVersion {
	case "2024-11-05", "2025-03-26", "2025-06-18", "2025-11-25":
		c.version = result.ProtocolVersion
	default:
		return fail(fmt.Errorf("MCP server selected unsupported protocol version %q", result.ProtocolVersion))
	}
	if err = c.notify(ctx, "notifications/initialized"); err != nil {
		return fail(err)
	}
	return c, nil
}

// Call returns the complete result. The client serializes calls; use separate
// sessions for concurrent operations. params must be a JSON object or omitted.
func (c *Client) Call(ctx context.Context, method string, params json.RawMessage) (json.RawMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil, errors.New("MCP session is closed")
	}
	if method == "" {
		return nil, errors.New("MCP method is required")
	}
	if len(params) > 0 {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(params, &object); err != nil || object == nil {
			return nil, errors.New("MCP parameters must be a JSON object")
		}
	}
	requestCtx, cancel := context.WithCancel(ctx)
	stop := context.AfterFunc(c.ctx, cancel)
	defer cancel()
	defer stop()
	if err := requestCtx.Err(); err != nil {
		return nil, err
	}
	c.id++
	id := strconv.FormatUint(c.id, 10)
	message := map[string]any{"jsonrpc": "2.0", "id": c.id, "method": method}
	if len(params) > 0 {
		message["params"] = params
	}
	data, _ := json.Marshal(message)
	raw, err := c.exchange(requestCtx, data, id)
	if err != nil && requestCtx.Err() != nil {
		return nil, requestCtx.Err()
	}
	return raw, err
}

func (c *Client) notify(ctx context.Context, method string) error {
	data, _ := json.Marshal(map[string]string{"jsonrpc": "2.0", "method": method})
	_, err := c.exchange(ctx, data, "")
	return err
}

func (c *Client) exchange(ctx context.Context, data []byte, id string) (json.RawMessage, error) {
	if c.command != nil {
		stop := context.AfterFunc(ctx, func() {
			c.cancel()
			// Cancellation closes pipes to unblock I/O; they may already be closed.
			_ = c.input.Close()
			_ = c.output.Close()
		})
		defer stop()
		if _, err := c.input.Write(append(data, '\n')); err != nil {
			return nil, errors.New("MCP local server input closed")
		}
		if id == "" {
			return nil, nil
		}
		total := 0
		for c.scanner.Scan() {
			line := c.scanner.Bytes()
			total += len(line)
			if total > maxResponse {
				return nil, errors.New("MCP response exceeds 64 MiB")
			}
			if len(bytes.TrimSpace(line)) == 0 {
				continue
			}
			result, matched, err := c.message(ctx, line, id)
			if err != nil || matched {
				return result, err
			}
		}
		if err := c.scanner.Err(); err != nil {
			return nil, errors.New("MCP local server output failed or exceeded its line limit")
		}
		return nil, io.ErrUnexpectedEOF
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	c.headers(req)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errors.New("MCP HTTP transport failed")
	}
	defer func() {
		// Response read errors determine the request result; closing is cleanup.
		_ = resp.Body.Close()
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, c.httpError(resp)
	}
	if value := resp.Header.Get("Mcp-Session-Id"); value != "" {
		c.session = value
	}
	if id == "" {
		if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return nil, fmt.Errorf("MCP notification returned HTTP %d", resp.StatusCode)
		}
		return nil, nil
	}
	media, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if media == "text/event-stream" {
		return c.readSSE(ctx, resp.Body, id)
	}
	if media != "application/json" {
		return nil, fmt.Errorf("MCP expected JSON or SSE response, got %q", media)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponse+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxResponse {
		return nil, errors.New("MCP response exceeds 64 MiB")
	}
	result, matched, err := c.message(ctx, body, id)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errors.New("MCP response did not match request ID")
	}
	return result, nil
}

func (c *Client) headers(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if c.key != "" {
		req.Header.Set("Authorization", "Bearer "+c.key)
	}
	if c.session != "" {
		req.Header.Set("Mcp-Session-Id", c.session)
	}
	req.Header.Set("MCP-Protocol-Version", c.version)
}

func (c *Client) httpError(resp *http.Response) error {
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 65537))
	text := string(raw)
	if len(raw) > 65536 {
		text = string(raw[:65536]) + " [truncated]"
	}
	if c.key != "" {
		text = strings.ReplaceAll(text, c.key, "[redacted]")
	}
	return &HTTPError{StatusCode: resp.StatusCode, Body: text}
}

func (c *Client) message(ctx context.Context, raw []byte, id string) (json.RawMessage, bool, error) {
	var message envelope
	if err := json.Unmarshal(raw, &message); err != nil || message.Version != "2.0" {
		return nil, false, errors.New("MCP received invalid JSON-RPC message")
	}
	if message.Method != "" {
		if len(message.ID) == 0 {
			return nil, false, nil
		}
		reply := map[string]any{"jsonrpc": "2.0", "id": message.ID}
		if message.Method == "ping" {
			reply["result"] = map[string]any{}
		} else {
			reply["error"] = map[string]any{"code": -32601, "message": "Client method not supported"}
		}
		data, _ := json.Marshal(reply)
		if _, err := c.exchange(ctx, data, ""); err != nil {
			return nil, false, err
		}
		return nil, false, nil
	}
	if string(message.ID) != id {
		return nil, false, nil
	}
	if len(message.Error) > 0 && string(message.Error) != "null" {
		raw := append(json.RawMessage(nil), message.Error...)
		return nil, true, &RPCError{Raw: raw, key: c.key}
	}
	if len(message.Result) == 0 {
		return nil, true, errors.New("MCP response has neither result nor error")
	}
	return append(json.RawMessage(nil), message.Result...), true, nil
}

func (c *Client) readSSE(ctx context.Context, reader io.Reader, id string) (json.RawMessage, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 4096), maxLine)
	var data []string
	total := 0
	dispatch := func() (json.RawMessage, bool, error) {
		if len(data) == 0 {
			return nil, false, nil
		}
		raw := strings.Join(data, "\n")
		data = nil
		return c.message(ctx, []byte(raw), id)
	}
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		line := scanner.Text()
		total += len(line)
		if total > maxResponse {
			return nil, errors.New("MCP response exceeds 64 MiB")
		}
		if line == "" {
			result, matched, err := dispatch()
			if matched || err != nil {
				return result, err
			}
			continue
		}
		if strings.HasPrefix(line, "data:") {
			data = append(data, strings.TrimPrefix(strings.TrimPrefix(line, "data:"), " "))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	result, matched, err := dispatch()
	if matched || err != nil {
		return result, err
	}
	return nil, io.ErrUnexpectedEOF
}

// Close cancels in-flight operations and reaps an owned local server. Remote
// sessions are left to server expiration; MCP session DELETE is optional.
func (c *Client) Close() error {
	c.cancel()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	if c.stopPipes != nil {
		c.stopPipes()
	}
	if c.input != nil {
		_ = c.input.Close() // Cancellation may already have closed the pipe.
	}
	if c.output != nil {
		_ = c.output.Close() // Cancellation may already have closed the pipe.
	}
	if c.command != nil {
		_ = c.command.Wait()
	}
	return nil
}
