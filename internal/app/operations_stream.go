package app

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Event retains SSE data verbatim (including non-JSON payloads such as [DONE]).
// Encoding this struct as NDJSON preserves multiline data without losing names.
type Event struct {
	Event string `json:"event"`
	Data  string `json:"data"`
	ID    string `json:"id,omitempty"`
}

// StreamOperation executes once and passes each complete SSE event to emit.
// The callback owns output and may stop the request by returning an error.
func (c *Client) StreamOperation(ctx context.Context, input OperationRequest, emit func(Event) error) error {
	if emit == nil {
		return errors.New("event callback is required")
	}
	if err := c.requireAPIKey(); err != nil {
		return err
	}
	op, err := lookupOperation(input.Operation)
	if err != nil {
		return err
	}
	supported := false
	for _, mode := range op.ResponseModes {
		if mode == "sse" {
			supported = true
		}
	}
	if !supported {
		return fmt.Errorf("%s does not support SSE", op.Name)
	}
	obj, err := JSONObject(input.JSON)
	if len(input.JSON) == 0 {
		obj = make(map[string]json.RawMessage)
		err = nil
	}
	if err != nil {
		return err
	}
	obj["stream"] = json.RawMessage(`true`)
	input.JSON, err = json.Marshal(obj)
	if err != nil {
		return err
	}
	return c.runOperation(ctx, op, func() error {
		req, cleanup, err := c.operationRequest(ctx, op, input, true)
		if err != nil {
			return err
		}
		defer cleanup()
		resp, err := c.streamHTTPClient.Do(req)
		if err != nil {
			return fmt.Errorf("%s stream: %w", op.Name, err)
		}
		defer resp.Body.Close() //nolint:errcheck // Response-body cleanup must not replace the stream or callback result.
		if !operationSuccess(op, resp.StatusCode) {
			return providerError(resp)
		}
		if !strings.HasPrefix(strings.ToLower(resp.Header.Get("Content-Type")), "text/event-stream") {
			return fmt.Errorf("%s expected text/event-stream, got %q", op.Name, resp.Header.Get("Content-Type"))
		}
		return readOperationEvents(ctx, resp.Body, op.Name, emit)
	})
}

func readOperationEvents(ctx context.Context, reader io.Reader, operation string, emit func(Event) error) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 4096), 8<<20)
	event := Event{}
	var data []string
	hasData := false
	terminal := false
	dispatch := func() error {
		if !hasData {
			event = Event{}
			return nil
		}
		if event.Event == "" {
			event.Event = "message"
		}
		event.Data = strings.Join(data, "\n")
		if err := emit(event); err != nil {
			return err
		}

		var payload map[string]json.RawMessage
		_ = json.Unmarshal([]byte(event.Data), &payload)
		var kind string
		_ = json.Unmarshal(payload["type"], &kind)
		if event.Event == "error" || event.Event == "response.failed" || event.Event == "response.incomplete" || kind == "error" || kind == "response.failed" || kind == "response.incomplete" {
			return eventAPIError(event)
		}
		if raw, ok := payload["error"]; ok && string(raw) != "null" {
			return eventAPIError(event)
		}
		terminal = operationTerminal(operation, event, payload)
		event = Event{}
		data = nil
		hasData = false
		return nil
	}
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return err
		}
		line := scanner.Text()
		if line == "" {
			if err := dispatch(); err != nil {
				return err
			}
			if terminal {
				return nil
			}
			continue
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		field, value, ok := strings.Cut(line, ":")
		if !ok {
			value = ""
		}
		value = strings.TrimPrefix(value, " ")
		switch field {
		case "event":
			event.Event = value
		case "data":
			data = append(data, value)
			hasData = true
		case "id":
			if !strings.ContainsRune(value, '\x00') {
				event.ID = value
			}
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read SSE: %w", err)
	}
	// A final event without a blank separator is still useful and must be retained.
	if err := dispatch(); err != nil {
		return err
	}
	if terminal {
		return nil
	}
	return io.ErrUnexpectedEOF
}

func operationTerminal(operation string, event Event, payload map[string]json.RawMessage) bool {
	if event.Data == "[DONE]" {
		return true
	}
	var kind string
	_ = json.Unmarshal(payload["type"], &kind)
	if operation == "audio" {
		return event.Event == "transcript.text.done" || kind == "transcript.text.done"
	}
	if operation == "messages" {
		return event.Event == "message_stop" || kind == "message_stop"
	}
	if operation == "responses" {
		for _, name := range []string{"response.completed", "response.failed", "response.incomplete"} {
			if event.Event == name || kind == name {
				return true
			}
		}
	}
	if operation == "agent_run" {
		var status string
		_ = json.Unmarshal(payload["status"], &status)
		if status == "completed" || status == "success" || status == "SUCCESS" || event.Event == "done" || event.Event == "finish" {
			return true
		}
		var choices []struct {
			FinishReason *string `json:"finish_reason"`
		}
		_ = json.Unmarshal(payload["choices"], &choices)
		for _, choice := range choices {
			if choice.FinishReason != nil && *choice.FinishReason != "" {
				return true
			}
		}
	}
	return false
}

func eventAPIError(event Event) error {
	body := event.Data
	if len(body) > maxDiagnostic {
		body = body[:maxDiagnostic] + " [truncated]"
	}
	return &APIError{StatusCode: 200, Body: body, RequestID: event.ID}
}
