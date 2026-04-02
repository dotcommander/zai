package app

import (
	"context"
	"errors"
	"fmt"
)

// TextToSpeech synthesizes text to audio using Z.AI's TTS API.
// Returns raw audio bytes (WAV or PCM format).
func (c *Client) TextToSpeech(ctx context.Context, text string, opts TTSOptions) ([]byte, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}
	if text == "" {
		return nil, errors.New("text is required for TTS")
	}

	model := resolveModel(opts.Model, c.config.TTSModel, "glm-tts")

	voice := opts.Voice
	if voice == "" {
		voice = "tongtong"
	}

	responseFormat := opts.ResponseFormat
	if responseFormat == "" {
		responseFormat = "wav"
	}

	reqData := TTSSpeechRequest{
		Model:          model,
		Input:          text,
		Voice:          voice,
		Speed:          opts.Speed,
		Volume:         opts.Volume,
		ResponseFormat: responseFormat,
	}

	body, err := c.executeJSONRequest(ctx, "audio/speech", reqData)
	if err != nil {
		return nil, fmt.Errorf("TTS API error: %w", err)
	}

	c.logger.DebugContext(ctx, "TTS complete", "bytes", len(body))

	return body, nil
}
