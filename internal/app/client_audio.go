package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

// TranscribeAudio transcribes an audio file using Z.AI's ASR model.
func (c *Client) TranscribeAudio(ctx context.Context, audioPath string, opts TranscriptionOptions) (*TranscriptionResponse, error) {
	if opts.Stream {
		return c.TranscribeAudioStream(ctx, audioPath, opts, nil)
	}
	req, err := c.audioRequest(ctx, audioPath, opts)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer closeBody(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, audioResponseError(resp)
	}
	var result TranscriptionResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 10*1024*1024)).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &result, nil
}

func audioResponseError(resp *http.Response) error {
	return providerError(resp)
}

func (c *Client) audioRequest(ctx context.Context, audioPath string, opts TranscriptionOptions) (*http.Request, error) { //nolint:gocyclo,funlen
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	// Validate audio file
	if audioPath == "" {
		return nil, errors.New("audio file path is required")
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

	model := resolveModel(opts.Model, c.config.AudioModel, "glm-asr-2512")

	// Build only multipart framing in memory; retain the injected reader's bytes
	// as the file part so requests remain replayable without a second full copy.
	var framing bytes.Buffer
	writer := multipart.NewWriter(&framing)
	if _, err := writer.CreateFormFile("file", filepath.Base(audioPath)); err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	prefix := bytes.Clone(framing.Bytes())
	framing.Reset()

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
		var hotwordsJSON []byte
		hotwordsJSON, err = json.Marshal(opts.Hotwords)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal hotwords: %w", err)
		}
		writer.WriteField("hotwords", string(hotwordsJSON)) //nolint:errcheck // multipart field write
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart form: %w", err)
	}
	suffix := bytes.Clone(framing.Bytes())
	newBody := func() io.ReadCloser {
		return io.NopCloser(io.MultiReader(bytes.NewReader(prefix), bytes.NewReader(data), bytes.NewReader(suffix)))
	}

	// Build request
	url := fmt.Sprintf("%s/audio/transcriptions", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, newBody())
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.ContentLength = int64(len(prefix) + len(data) + len(suffix))
	req.GetBody = func() (io.ReadCloser, error) { return newBody(), nil }
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Accept-Language", "en-US,en")

	return req, nil
}
