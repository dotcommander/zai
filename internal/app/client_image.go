package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// GenerateImage creates an image using the Z.AI image generation API.
func (c *Client) GenerateImage(ctx context.Context, prompt string, opts ImageOptions) (*ImageResponse, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	// Validate options
	if err := validateImageOptions(opts); err != nil {
		return nil, fmt.Errorf("invalid image options: %w", err)
	}

	model := resolveModel(opts.Model, c.config.ImageModel, "cogview-4-250304")

	reqData := ImageGenerationRequest{
		Model:   model,
		Prompt:  prompt,
		Quality: opts.Quality,
		Size:    opts.Size,
		UserID:  opts.UserID,
	}

	// Set defaults
	if reqData.Quality == "" {
		reqData.Quality = "hd"
	}
	if reqData.Size == "" {
		reqData.Size = "1024x1024"
	}

	var imageResp ImageResponse
	body, err := c.executeJSONRequest(ctx, "images/generations", reqData)
	if err != nil {
		return nil, fmt.Errorf("image generation API error: %w", err)
	}
	if err := json.Unmarshal(body, &imageResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal image response: %w", err)
	}

	if len(imageResp.Data) == 0 {
		return nil, errors.New("no images in response")
	}

	c.logger.DebugContext(ctx, "generated image",
		"url", imageResp.Data[0].URL,
		"width", imageResp.Data[0].Width,
		"height", imageResp.Data[0].Height)

	return &imageResp, nil
}

// validateImageOptions checks if image options are valid.
func validateImageOptions(opts ImageOptions) error {
	// Validate quality
	if opts.Quality != "" && opts.Quality != "hd" && opts.Quality != "standard" {
		return fmt.Errorf("invalid quality: %s (must be 'hd' or 'standard')", opts.Quality)
	}

	// Validate size format
	if opts.Size != "" {
		supportedSizes := map[string]bool{
			"1024x1024": true,
			"1024x768":  true,
			"768x1024":  true,
			"512x512":   true,
			"768x1344":  true,
			"864x1152":  true,
			"1344x768":  true,
			"1152x864":  true,
			"1440x720":  true,
			"720x1440":  true,
		}
		if !supportedSizes[opts.Size] {
			return fmt.Errorf("invalid size: %s (supported: 1024x1024, 1024x768, 768x1024, 512x512, 768x1344, 864x1152, 1344x768, 1152x864, 1440x720, 720x1440)", opts.Size)
		}
	}

	return nil
}
