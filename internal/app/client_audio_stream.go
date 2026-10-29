package app

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TranscribeAudioStream emits text deltas as they arrive and returns the complete transcript.
// A completion event or [DONE] is required; EOF alone is an incomplete transcript.
func (c *Client) TranscribeAudioStream(ctx context.Context, audioPath string, opts TranscriptionOptions, onDelta func(string)) (*TranscriptionResponse, error) {
	opts.Stream = true
	req, err := c.audioRequest(ctx, audioPath, opts)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/event-stream")
	resp, err := c.streamHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send streaming transcription request: %w", err)
	}
	defer closeBody(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, audioResponseError(resp)
	}
	return readAudioStream(ctx, resp.Body, onDelta)
}

func readAudioStream(ctx context.Context, body io.Reader, onDelta func(string)) (*TranscriptionResponse, error) {
	result := &TranscriptionResponse{}
	var transcript strings.Builder
	var data strings.Builder
	var eventName string
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		line := scanner.Text()
		if line == "" {
			done, err := consumeAudioStreamFrame(&data, eventName, result, &transcript, onDelta)
			eventName = ""
			if err != nil {
				return nil, err
			}
			if done {
				result.Text = transcript.String()
				return result, nil
			}
			continue
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		field, value, _ := strings.Cut(line, ":")
		value = strings.TrimPrefix(value, " ")
		switch field {
		case "event":
			eventName = value
		case "data":
			if data.Len()+len(value) > 1024*1024 {
				return nil, errors.New("audio stream frame exceeds 1 MiB")
			}
			data.WriteString(value)
			data.WriteByte('\n')
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read audio stream: %w", err)
	}
	if data.Len() != 0 {
		done, err := consumeAudioStreamFrame(&data, eventName, result, &transcript, onDelta)
		if err != nil {
			return nil, err
		}
		if done {
			result.Text = transcript.String()
			return result, nil
		}
	}
	return nil, errors.New("audio stream terminated before completion")
}

func consumeAudioStreamFrame(data *strings.Builder, eventName string, result *TranscriptionResponse, transcript *strings.Builder, onDelta func(string)) (bool, error) {
	payload := strings.TrimSuffix(data.String(), "\n")
	data.Reset()
	if payload == "" {
		return false, nil
	}
	if payload == "[DONE]" {
		return true, nil
	}
	var event struct {
		TranscriptionResponse
		Type  string `json:"type"`
		Delta string `json:"delta"`
	}
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		return false, fmt.Errorf("malformed audio stream frame: %w", err)
	}
	if event.Type == "" {
		event.Type = eventName
	}
	if event.ID != "" {
		result.ID = event.ID
	}
	if event.Model != "" {
		result.Model = event.Model
	}
	if event.Created != 0 {
		result.Created = event.Created
	}
	if event.RequestID != "" {
		result.RequestID = event.RequestID
	}
	switch event.Type {
	case "transcript.text.delta":
		transcript.WriteString(event.Delta)
		if onDelta != nil && event.Delta != "" {
			onDelta(event.Delta)
		}
	case "transcript.text.done":
		// Some providers repeat the whole transcript on completion.
		// Deltas already emitted must not be printed or appended again.
		if transcript.Len() == 0 {
			text := event.Text
			if text == "" {
				text = event.Delta
			}
			transcript.WriteString(text)
			if onDelta != nil && text != "" {
				onDelta(text)
			}
		}
		return true, nil
	default:
		return false, fmt.Errorf("unexpected audio stream event %q", event.Type)
	}
	return false, nil
}
