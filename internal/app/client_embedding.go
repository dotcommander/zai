package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// CreateEmbedding generates embeddings for the given texts.
func (c *Client) CreateEmbedding(ctx context.Context, texts []string, opts EmbeddingOptions) (*EmbeddingResponse, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}
	if len(texts) == 0 {
		return nil, errors.New("at least one text input is required")
	}

	model := resolveModel(opts.Model, c.config.EmbeddingModel, "embedding-3")

	reqData := EmbeddingRequest{
		Model: model,
		Input: texts,
	}

	var embResp EmbeddingResponse
	body, err := c.executeJSONRequest(ctx, "embeddings", reqData)
	if err != nil {
		return nil, fmt.Errorf("embedding API error: %w", err)
	}
	if err := json.Unmarshal(body, &embResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal embedding response: %w", err)
	}

	c.logger.DebugContext(ctx, "embedding complete",
		"model", embResp.Model,
		"results", len(embResp.Data),
		"prompt_tokens", embResp.Usage.PromptTokens)

	return &embResp, nil
}
