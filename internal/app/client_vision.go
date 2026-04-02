package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// Vision analyzes an image using Z.AI's vision model (glm-4.6v).
// imageBase64 should be a data URI like "data:image/jpeg;base64,<base64-data>" or a raw base64 string.
func (c *Client) Vision(ctx context.Context, prompt string, imageBase64 string, opts VisionOptions) (string, error) {
	if err := c.requireAPIKey(); err != nil {
		return "", err
	}

	// Validate prompt
	if prompt == "" {
		prompt = "What do you see in this image? Please describe it in detail."
	}

	// Validate image input
	if imageBase64 == "" {
		return "", errors.New("image data is required")
	}

	model := resolveModel(opts.Model, c.config.VisionModel, "glm-4.6v")

	// Build multimodal messages
	messages := []VisionMessage{
		{
			Role: "user",
			Content: []ContentPart{
				{
					Type: "text",
					Text: prompt,
				},
				{
					Type: "image_url",
					ImageURL: &ImageURLContent{
						URL: imageBase64,
					},
				},
			},
		},
	}

	// Build request
	reqData := VisionRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	// Apply optional overrides
	if opts.Temperature != nil {
		reqData.Temperature = *opts.Temperature
	} else {
		reqData.Temperature = 0.3 // Lower temp for vision
	}

	if opts.MaxTokens != nil {
		reqData.MaxTokens = *opts.MaxTokens
	} else {
		reqData.MaxTokens = 4096
	}

	if opts.TopP != nil {
		reqData.TopP = *opts.TopP
	} else {
		reqData.TopP = 0.9
	}

	var chatResp ChatResponse
	body, err := c.executeJSONRequest(ctx, "chat/completions", reqData)
	if err != nil {
		return "", fmt.Errorf("vision API error: %w", err)
	}
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal vision response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", errors.New("no choices in vision response")
	}

	c.logger.DebugContext(ctx, "vision complete",
		"total_tokens", chatResp.Usage.TotalTokens,
		"prompt_tokens", chatResp.Usage.PromptTokens,
		"completion_tokens", chatResp.Usage.CompletionTokens)

	return chatResp.Choices[0].Message.Content, nil
}
