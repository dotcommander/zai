package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/dotcommander/zai/internal/app"
	"github.com/dotcommander/zai/internal/app/utils"
)

var (
	videoPrompt      string
	videoQuality     string
	videoSize        string
	videoFPS         int
	videoDuration    int
	videoWithAudio   bool
	videoOutput      string
	videoShow        bool
	videoModel       string
	videoUserID      string
	videoRequestID   string
	videoImageURLs   []string
	videoPollTimeout time.Duration
)

var videoCmd = &cobra.Command{
	Use:   "video \"description\"",
	Short: "Generate videos using Z.AI's CogVideoX-3 API",
	Long: `Generate videos with Z.AI's CogVideoX-3 model.

Text-to-Video:
  zai video "A cat playing with a ball"

Image-to-Video:
  zai video -f https://example.com/image.jpg "Make the picture move"

First & Last Frame:
  zai video -f first.jpg -f last.jpg "Smooth transition"

Examples:
  zai video "a sunset over the ocean" --quality quality --size 1920x1080
  zai video "prompt" --fps 60 --duration 10 --with-audio
  zai video "prompt" --output my-video.mp4 --show`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runVideoGeneration(args[0])
	},
}

func init() {
	// Main video command
	videoCmd.Flags().StringVarP(&videoQuality, "quality", "q", "speed", "Quality mode: speed (fast) or quality (higher quality)")
	videoCmd.Flags().StringVarP(&videoSize, "size", "s", "1920x1080", "Video size: 1280x720, 1024x1024, 1920x1080, 3840x2160, etc.")
	videoCmd.Flags().IntVar(&videoFPS, "fps", 30, "Frame rate: 30 or 60")
	videoCmd.Flags().IntVar(&videoDuration, "duration", 5, "Duration: 5 or 10 seconds")
	videoCmd.Flags().BoolVar(&videoWithAudio, "with-audio", false, "Generate AI sound effects")
	videoCmd.Flags().StringVarP(&videoOutput, "output", "o", "", "Save video to file path")
	videoCmd.Flags().BoolVarP(&videoShow, "show", "S", false, "Open video with default player after generation")
	videoCmd.Flags().StringVarP(&videoModel, "model", "m", "", "Override default video model")
	videoCmd.Flags().StringVar(&videoUserID, "user-id", "", "User ID for analytics")
	videoCmd.Flags().StringVar(&videoRequestID, "request-id", "", "Unique request ID")
	videoCmd.Flags().StringArrayVarP(&videoImageURLs, "file", "f", []string{}, "Image URL(s) for image-to-video or first/last frame mode (can specify 1 or 2)")
	videoCmd.Flags().DurationVar(&videoPollTimeout, "poll-timeout", 3*time.Minute, "Maximum time to wait for video generation")

	// Register with root
	rootCmd.AddCommand(videoCmd)
}

func runVideoGeneration(prompt string) error {
	client := newClient()
	ctx, cancel := context.WithTimeout(context.Background(), videoPollTimeout)
	defer cancel()

	// Build options
	opts := buildVideoOptions()

	// Start video generation
	fmt.Printf("\n🎬 Starting video generation...\n")
	fmt.Printf("📝 Prompt: %s\n", prompt)
	if len(videoImageURLs) > 0 {
		fmt.Printf("🖼️  Image URLs: %d provided\n", len(videoImageURLs))
	}
	fmt.Printf("⚙️  Quality: %s, Size: %s, FPS: %d, Duration: %ds\n", opts.Quality, opts.Size, opts.FPS, opts.Duration)
	if opts.WithAudio {
		fmt.Printf("🔊 Audio: enabled\n")
	}

	response, err := client.GenerateVideo(ctx, prompt, opts)
	if err != nil {
		return fmt.Errorf("failed to start video generation: %w", err)
	}

	// Poll for result
	fmt.Printf("📋 Task ID: %s\n", response.ID)
	fmt.Printf("⏳ Polling for result (this may take 1-3 minutes)...\n")

	result, err := pollForResult(ctx, client, response.ID)
	if err != nil {
		return err
	}

	// Display and handle the result
	return displayVideoResult(result, prompt)
}

// pollForResult polls for video generation completion with spinner.
func pollForResult(ctx context.Context, client *app.Client, taskID string) (*app.VideoResultResponse, error) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	spinnerIdx := 0
	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("video generation timeout")
		case <-ticker.C:
			elapsed := time.Since(startTime)
			result, err := client.RetrieveVideoResult(ctx, taskID)
			if err != nil {
				return nil, err
			}

			// Update spinner
			spinner := SpinnerFrames[spinnerIdx%len(SpinnerFrames)]
			spinnerIdx++

			switch result.TaskStatus {
			case "SUCCESS":
				fmt.Printf("\r%s ✅ Video generation complete! (%.1fs elapsed)\n", spinner, elapsed.Seconds())
				return result, nil
			case "FAIL":
				return nil, fmt.Errorf("video generation failed on server")
			case "PROCESSING":
				fmt.Printf("\r%s ⏳ Processing... (%.1fs elapsed)   ", spinner, elapsed.Seconds())
			}
		}
	}
}

// displayVideoResult handles displaying, saving, and opening the generated video.
func displayVideoResult(result *app.VideoResultResponse, prompt string) error {
	if len(result.VideoResult) == 0 {
		return fmt.Errorf("no video in result")
	}

	videoData := result.VideoResult[0]

	fmt.Printf("\n✅ Video generated successfully!\n")
	fmt.Printf("🔗 URL: %s\n", videoData.URL)
	if videoData.CoverImageURL != "" {
		fmt.Printf("🖼️  Cover: %s\n", videoData.CoverImageURL)
	}

	// Determine output path
	outputPath := videoOutput
	if outputPath == "" {
		timestamp := time.Now().Format("20060102-150405")
		ext := ".mp4"
		outputPath = fmt.Sprintf("zai-video-%s%s", timestamp, ext)
	}

	// Save video to disk
	fmt.Printf("💾 Downloading to: %s\n", outputPath)
	downloader := app.NewMediaDownloader(nil)
	downloadResult := downloader.Download(videoData.URL, outputPath)
	if downloadResult.Error != nil {
		return fmt.Errorf("failed to save video: %w", downloadResult.Error)
	}

	fmt.Printf("📊 Size: %.2f MB\n", float64(downloadResult.Size)/(1024*1024))
	fmt.Printf("✅ Saved to: %s\n", outputPath)

	// Open in player
	if videoShow {
		if err := openVideoPlayer(outputPath); err != nil {
			fmt.Printf("⚠️  Warning: Failed to open video player: %v\n", err)
		}
	}

	return nil
}

// openVideoPlayer opens video file with default player.
func openVideoPlayer(filePath string) error {
	fmt.Printf("🎬 Opening video player...\n")
	return app.OpenWith(filePath)
}

// buildVideoOptions creates video options from command line flags and config.
func buildVideoOptions() app.VideoOptions {
	opts := app.VideoOptions{
		Quality:   videoQuality,
		Size:      videoSize,
		FPS:       videoFPS,
		Duration:  videoDuration,
		WithAudio: videoWithAudio,
		ImageURLs: videoImageURLs,
		UserID:    videoUserID,
		RequestID: videoRequestID,
		Model:     videoModel,
	}

	// Use configured model if not overridden
	if opts.Model == "" {
		opts.Model = getModelWithDefault("api.video_model", "cogvideox-3")
	}

	return opts
}

// enhanceVideoPrompt enhances the prompt using AI (optional feature, similar to image).
// This can be added later if needed - for now keeping it simple.
func enhanceVideoPrompt(client *app.Client, prompt string) (string, error) {
	// For future enhancement: use LLM to expand simple prompts into detailed video descriptions
	return prompt, nil
}

// readImageAndEncode reads an image file and returns base64 data URI.
// Returns URL unchanged if it's already an HTTP(S) URL.
func readImageAndEncode(imagePath string) (string, error) {
	// Pass through URLs unchanged
	if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
		return imagePath, nil
	}

	// Read local file using utils
	fileReader := utils.OSFileReader{}
	data, err := fileReader.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("read image: %w", err)
	}

	// Check file size (5MB max per API)
	const maxImageSize = 5 * 1024 * 1024
	if len(data) > maxImageSize {
		return "", fmt.Errorf("image file too large: %d bytes (max: 5MB)", len(data))
	}

	// Detect MIME type using utils (supports more formats: jpg, png, gif, webp)
	mimeType, err := utils.DetectImageMimeType(imagePath)
	if err != nil {
		return "", err
	}

	// Encode as data URI using utils
	return utils.EncodeBytesToDataURI(data, mimeType), nil
}
