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
func (c *Client) TranscribeAudio(ctx context.Context, audioPath string, opts TranscriptionOptions) (*TranscriptionResponse, error) { //nolint:gocyclo,funlen
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

	// Build multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file from memory
	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, copyErr := io.Copy(part, bytes.NewReader(data)); copyErr != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", copyErr)
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
		var hotwordsJSON []byte
		hotwordsJSON, err = json.Marshal(opts.Hotwords)
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

	c.logger.DebugContext(ctx, "sending audio transcription request", "url", url)

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

	c.logger.DebugContext(ctx, "transcription complete", "chars", len(transcriptionResp.Text), "model", transcriptionResp.Model)

	return &transcriptionResp, nil
}
