package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// GenerateVideo creates a video using Z.AI's CogVideoX-3 API (async).
func (c *Client) GenerateVideo(ctx context.Context, prompt string, opts VideoOptions) (*VideoGenerationResponse, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	model := resolveModel(opts.Model, "", "cogvideox-3")

	// Validate options
	if err := validateVideoOptions(opts); err != nil {
		return nil, fmt.Errorf("invalid video options: %w", err)
	}

	// Build request
	reqData := VideoGenerationRequest{
		Model:     model,
		Prompt:    prompt,
		ImageURL:  opts.ImageURLs,
		Quality:   opts.Quality,
		WithAudio: opts.WithAudio,
		Size:      opts.Size,
		FPS:       opts.FPS,
		Duration:  opts.Duration,
		RequestID: opts.RequestID,
		UserID:    opts.UserID,
	}

	// Set defaults
	if reqData.Quality == "" {
		reqData.Quality = "speed"
	}
	if reqData.Size == "" {
		reqData.Size = "1920x1080"
	}
	if reqData.FPS == 0 {
		reqData.FPS = 30
	}
	if reqData.Duration == 0 {
		reqData.Duration = 5
	}

	var videoResp VideoGenerationResponse
	body, err := c.executeJSONRequest(ctx, "videos/generations", reqData)
	if err != nil {
		return nil, fmt.Errorf("video generation API error: %w", err)
	}
	if err := json.Unmarshal(body, &videoResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal video response: %w", err)
	}

	c.logger.DebugContext(ctx, "video generation task created", "id", videoResp.ID, "status", videoResp.TaskStatus)

	return &videoResp, nil
}

// RetrieveVideoResult polls for async video generation result.
func (c *Client) RetrieveVideoResult(ctx context.Context, taskID string) (*VideoResultResponse, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	// Validate task ID
	if taskID == "" {
		return nil, errors.New("task ID is required")
	}

	var resultResp VideoResultResponse
	endpoint := fmt.Sprintf("async-result/%s", taskID)
	body, err := c.executeGetRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("retrieve video result API error: %w", err)
	}
	if err := json.Unmarshal(body, &resultResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal video result response: %w", err)
	}

	c.logger.DebugContext(ctx, "video result retrieved", "id", taskID, "status", resultResp.TaskStatus)

	return &resultResp, nil
}

// validateVideoOptions checks if video options are valid.
func validateVideoOptions(opts VideoOptions) error {
	// Validate quality
	if opts.Quality != "" && opts.Quality != "quality" && opts.Quality != "speed" {
		return fmt.Errorf("invalid quality: %s (must be 'quality' or 'speed')", opts.Quality)
	}

	// Validate size format
	if opts.Size != "" {
		supportedSizes := map[string]bool{
			"1280x720": true, "720x1280": true,
			"1024x1024": true,
			"1920x1080": true, "1080x1920": true,
			"2048x1080": true,
			"3840x2160": true,
		}
		if !supportedSizes[opts.Size] {
			return fmt.Errorf("invalid size: %s (supported: 1280x720, 720x1280, 1024x1024, 1920x1080, 1080x1920, 2048x1080, 3840x2160)", opts.Size)
		}
	}

	// Validate FPS
	if opts.FPS != 0 && opts.FPS != 30 && opts.FPS != 60 {
		return fmt.Errorf("invalid fps: %d (must be 30 or 60)", opts.FPS)
	}

	// Validate duration
	if opts.Duration != 0 && opts.Duration != 5 && opts.Duration != 10 {
		return fmt.Errorf("invalid duration: %d (must be 5 or 10 seconds)", opts.Duration)
	}

	// Validate image URLs (max 2 for first/last frame mode)
	if len(opts.ImageURLs) > 2 {
		return errors.New("too many image URLs (max 2 for first/last frame mode)")
	}

	return nil
}
