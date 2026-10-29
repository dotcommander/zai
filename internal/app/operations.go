package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Operation describes a finite, supported provider operation. Routes are relative
// to their family, never supplied by untrusted request JSON.
type Operation struct {
	Name            string   `json:"name"`
	Method          string   `json:"method"`
	Family          string   `json:"family"`
	Route           string   `json:"route"`
	Auth            string   `json:"auth"`
	RequestMode     string   `json:"request_mode"`
	ResponseModes   []string `json:"response_modes"`
	SuccessStatuses []int    `json:"success_statuses"`
}

var operationRegistry = []Operation{
	{"chat", "POST", "paas", "chat/completions", "bearer", "json", []string{"json", "sse"}, []int{200}},
	{"videos", "POST", "paas", "videos/generations", "bearer", "json", []string{"json"}, []int{200}},
	{"task_get", "GET", "paas", "async-result/{id}", "bearer", "none", []string{"json"}, []int{200}},
	{"images", "POST", "paas", "images/generations", "bearer", "json", []string{"json"}, []int{200}},
	{"image_async", "POST", "paas", "async/images/generations", "bearer", "json", []string{"json"}, []int{200}},
	{"audio", "POST", "paas", "audio/transcriptions", "bearer", "multipart", []string{"json", "sse"}, []int{200}},
	{"tokenize", "POST", "paas", "tokenizer", "bearer", "json", []string{"json"}, []int{200}},
	{"layout", "POST", "paas", "layout_parsing", "bearer", "json", []string{"json"}, []int{200}},
	{"search", "POST", "paas", "web_search", "bearer", "json", []string{"json"}, []int{200}},
	{"reader", "POST", "paas", "reader", "bearer", "json", []string{"json"}, []int{200}},
	{"files_upload", "POST", "paas", "files", "bearer", "multipart", []string{"json"}, []int{200}},
	{"agent_run", "POST", "v1", "agents", "bearer", "json", []string{"json", "sse"}, []int{200}},
	{"agent_result", "POST", "v1", "agents/async-result", "bearer", "json", []string{"json"}, []int{200}},
	{"agent_conversation", "POST", "v1", "agents/conversation", "bearer", "json", []string{"json"}, []int{200}},
	// Compatibility operations are documented by integration guides, outside OpenAPI.
	{"responses", "POST", "v1", "responses", "bearer", "json", []string{"json", "sse"}, []int{200}},
	{"messages", "POST", "anthropic", "v1/messages", "anthropic", "json", []string{"json", "sse"}, []int{200}},
}

// Operations returns an independent inventory suitable for discovery/help.
func Operations() []Operation {
	out := append([]Operation(nil), operationRegistry...)
	for i := range out {
		out[i].ResponseModes = append([]string(nil), out[i].ResponseModes...)
		out[i].SuccessStatuses = append([]int(nil), out[i].SuccessStatuses...)
	}
	return out
}

// OperationRequest preserves all JSON fields and nested values. UploadPath is a
// local file streamed as the multipart file part; Fields overrides JSON parts.
type OperationRequest struct {
	AcceptLanguage string
	Operation      string
	PathParams     map[string]string
	JSON           json.RawMessage
	UploadPath     string
	Fields         map[string]string
}

func lookupOperation(name string) (Operation, error) {
	for _, op := range operationRegistry {
		if op.Name == name {
			return op, nil
		}
	}
	return Operation{}, fmt.Errorf("unknown operation %q", name)
}

// JSONObject rejects arrays, null, trailing input, and malformed JSON without
// coercing numbers to floats. Callers may mutate individual RawMessage fields.
func JSONObject(data []byte) (map[string]json.RawMessage, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("request must be a JSON object: %w", err)
	}
	if obj == nil {
		return nil, errors.New("request must be a JSON object")
	}
	return obj, nil
}

func (c *Client) operationURL(op Operation, params map[string]string) (string, error) {
	base := strings.TrimRight(c.config.BaseURL, "/")
	if op.Name == "chat" && c.config.UseCoding && c.config.CodingBaseURL != "" {
		base = strings.TrimRight(c.config.CodingBaseURL, "/")
	}
	u, err := url.Parse(base)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return "", errors.New("API base URL must be an absolute HTTP(S) URL")
	}
	if u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return "", errors.New("API base URL cannot contain credentials, query, or fragment")
	}
	if op.Family != "paas" {
		u.Path = "/api/" + op.Family
		u.RawPath = ""
	} else if u.Path == "" || u.Path == "/" || u.Path == "/api" {
		u.Path = "/api/paas/v4"
	}
	route := op.Route
	if strings.Contains(route, "{id}") {
		id := params["id"]
		if strings.TrimSpace(id) == "" {
			return "", errors.New("task id is required")
		}
		route = strings.ReplaceAll(route, "{id}", url.PathEscape(id))
	}
	// Append escaped route to the serialized URL to preserve encoded path segments.
	return strings.TrimRight(u.String(), "/") + "/" + route, nil
}

const maxOperationResponse = 64 << 20
const maxDiagnostic = 64 << 10

func providerError(resp *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxDiagnostic+1))
	if err != nil {
		return fmt.Errorf("read API error: %w", err)
	}
	if len(body) > maxDiagnostic {
		body = append(body[:maxDiagnostic], []byte(" [truncated]")...)
	}
	id := resp.Header.Get("x-request-id")
	if id == "" {
		id = resp.Header.Get("request-id")
	}
	if id == "" {
		var obj struct {
			RequestID string `json:"request_id"`
		}
		_ = json.Unmarshal(body, &obj)
		id = obj.RequestID
	}
	return &APIError{StatusCode: resp.StatusCode, Body: string(body), RequestID: id}
}

func (c *Client) operationRequest(ctx context.Context, op Operation, input OperationRequest, stream bool) (*http.Request, func(), error) {
	endpoint, err := c.operationURL(op, input.PathParams)
	if err != nil {
		return nil, nil, err
	}
	var body io.Reader
	contentType := "application/json"
	cleanup := func() {}
	switch op.RequestMode {
	case "multipart":
		body, contentType, cleanup, err = c.operationMultipart(ctx, op, input)
		if err != nil {
			return nil, nil, err
		}
	case "json":
		data := input.JSON
		if len(data) == 0 {
			data = json.RawMessage(`{}`)
		}
		if _, err = JSONObject(data); err != nil {
			return nil, nil, err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, op.Method, endpoint, body)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("Accept-Language", "en-US,en")
	if input.AcceptLanguage != "" {
		req.Header.Set("Accept-Language", input.AcceptLanguage)
	}
	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}
	if stream {
		req.Header.Set("Accept", "text/event-stream")
	} else {
		req.Header.Set("Accept", "application/json")
	}
	if op.Auth == "anthropic" {
		req.Header.Set("x-api-key", c.config.APIKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	}
	return req, cleanup, nil
}

// operationBreaker creates a breaker for every registry family under the same
// policy as convenience requests. Existing endpoint breakers remain shared.
func (c *Client) operationBreaker(op Operation) *CircuitBreaker {
	if !c.config.CircuitBreaker.Enabled {
		return nil
	}
	key := extractEndpointName(op.Route)
	if key == "default" {
		key = op.Family + ":" + op.Name
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	cb := c.circuitBreakers[key]
	if cb == nil {
		cb = NewCircuitBreaker(key, c.config.CircuitBreaker, c.logger)
		c.circuitBreakers[key] = cb
	}
	return cb
}

func (c *Client) runOperation(ctx context.Context, op Operation, fn func() error) error {
	if cb := c.operationBreaker(op); cb != nil {
		return cb.executeWithContext(ctx, fn)
	}
	return fn()
}

// ExecuteOperation executes exactly once for generation and upload operations.
// Only safe GET retrievals are eligible for configured transient retries.
func (c *Client) ExecuteOperation(ctx context.Context, input OperationRequest) (json.RawMessage, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}
	op, err := lookupOperation(input.Operation)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	attempt := func() error {
		return c.runOperation(ctx, op, func() error {
			req, cleanup, err := c.operationRequest(ctx, op, input, false)
			if err != nil {
				return err
			}
			defer cleanup()
			resp, err := c.httpClient.Do(req)
			if err != nil {
				return fmt.Errorf("%s request: %w", op.Name, err)
			}
			defer resp.Body.Close() //nolint:errcheck // Response-body cleanup must not replace the request or decoding result.
			if !operationSuccess(op, resp.StatusCode) {
				return providerError(resp)
			}
			data, err := io.ReadAll(io.LimitReader(resp.Body, maxOperationResponse+1))
			if err != nil {
				return fmt.Errorf("read %s response: %w", op.Name, err)
			}
			if len(data) > maxOperationResponse {
				return errors.New("API response exceeds 64 MiB")
			}
			if !json.Valid(data) {
				return fmt.Errorf("%s returned invalid JSON", op.Name)
			}
			result = append(json.RawMessage(nil), data...)
			return nil
		})
	}
	err = attempt()
	if op.Method == http.MethodGet {
		for n := 0; err != nil && n+1 < c.config.RetryConfig.MaxAttempts && isRetryableError(err); n++ {
			delay := calculateBackoff(n+1, c.config.RetryConfig.InitialBackoff, c.config.RetryConfig.MaxBackoff)
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
			}
			err = attempt()
		}
	}
	return result, err
}

func operationSuccess(op Operation, status int) bool {
	for _, allowed := range op.SuccessStatuses {
		if status == allowed {
			return true
		}
	}
	return false
}
