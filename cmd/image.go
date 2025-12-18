package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"zai/internal/app"
)

var (
	imagePrompt    string
	imageQuality   string
	imageSize      string
	imageOutput    string
	imageShow      bool
	imageCopy      bool
	imageModel     string
	imageUserID    string
	imageEnhance   bool
	imageNoEnhance bool
)

var imageCmd = &cobra.Command{
	Use:   "image \"description\"",
	Short: "Generate images using Z.AI's image generation API",
	Long: `Generate images with Z.AI's CogView-4 model.

Examples:
  zai image "a cat wearing sunglasses"
  zai image "sunset on mars" --quality hd --size 1024x1024
  zai image "abstract art" --output my-art.png
  zai image "logo" --copy --size 512x512
  zai image "sunset" --no-enhance    # Skip prompt enhancement`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runImageGeneration(args[0])
	},
}

var imageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available image generation models",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runImageModelList()
	},
}

func init() {
	// Main image command
	imageCmd.Flags().StringVarP(&imageQuality, "quality", "q", "hd", "Image quality: hd or standard (default: hd)")
	imageCmd.Flags().StringVarP(&imageSize, "size", "s", "1024x1024", "Image size: 1024x1024, 1024x768, 768x1024, or 512x512 (default: 1024x1024)")
	imageCmd.Flags().StringVarP(&imageOutput, "output", "o", "", "Save image to file path")
	imageCmd.Flags().BoolVarP(&imageShow, "show", "S", false, "Open image with default viewer after generation")
	imageCmd.Flags().BoolVarP(&imageCopy, "copy", "c", false, "Copy image to clipboard (macOS only)")
	imageCmd.Flags().StringVarP(&imageModel, "model", "m", "", "Override default image model")
	imageCmd.Flags().StringVar(&imageUserID, "user-id", "", "User ID for analytics")
	imageCmd.Flags().BoolVarP(&imageEnhance, "enhance", "e", true, "Enhance prompt with AI before generation")
	imageCmd.Flags().BoolVar(&imageNoEnhance, "no-enhance", false, "Disable prompt enhancement")

	// Add subcommands
	imageCmd.AddCommand(imageListCmd)

	// Register with root
	rootCmd.AddCommand(imageCmd)
}

// shouldEnhancePrompt determines if prompt enhancement should be used.
// --no-enhance explicitly disables, otherwise --enhance (default true) controls.
func shouldEnhancePrompt() bool {
	if imageNoEnhance {
		return false
	}
	return imageEnhance
}

// enhanceImagePrompt uses the chat LLM to transform a simple description
// into a professional, detailed image generation prompt.
func enhanceImagePrompt(client *app.Client, prompt string) (string, error) {
	systemPrompt := `You are an expert image prompt engineer. Transform simple descriptions into ultra-detailed, professional image generation prompts.

Include relevant aspects from:
- Photography: lens type, aperture, focal length, depth of field, exposure
- Lighting: natural/studio, direction, quality, golden hour, rim light, volumetric
- Style: photorealistic, cinematic, artistic, illustration, 3D render
- Mood: atmosphere, emotion, color palette, tone
- Technical: resolution (8K, 4K), detail level, sharpness, clarity
- Composition: rule of thirds, leading lines, framing, perspective
- Film/Camera: specific camera models, film stock, color grading

Output ONLY the enhanced prompt, nothing else. Keep it under 500 characters.`

	opts := app.ChatOptions{
		Temperature: app.Float64Ptr(0.7),
		MaxTokens:   app.IntPtr(200),
		Context: []app.Message{
			{Role: "system", Content: systemPrompt},
		},
	}

	enhanced, err := client.Chat(context.Background(), "Enhance this image prompt: "+prompt, opts)
	if err != nil {
		return prompt, err // Fall back to original on error
	}
	return strings.TrimSpace(enhanced), nil
}

func runImageGeneration(prompt string) error {
	client := newClient()

	// Build image options
	opts := app.ImageOptions{
		Quality: imageQuality,
		Size:    imageSize,
		UserID:  imageUserID,
		Model:   imageModel,
	}

	// Use configured model if not overridden
	if opts.Model == "" {
		opts.Model = viper.GetString("api.image_model")
		if opts.Model == "" {
			opts.Model = "cogview-4-250304"
		}
	}

	// Enhance prompt if enabled
	finalPrompt := prompt
	if shouldEnhancePrompt() {
		fmt.Printf("🎨 Original: %s\n", prompt)
		fmt.Printf("✨ Enhancing prompt...\n")
		enhanced, err := enhanceImagePrompt(client, prompt)
		if err != nil {
			fmt.Printf("⚠️  Enhancement failed, using original: %v\n", err)
		} else {
			// Combine original + enhanced for best results
			finalPrompt = prompt + ". " + enhanced
			fmt.Printf("✨ Enhanced: %s\n", enhanced)
		}
	} else {
		fmt.Printf("🎨 Generating image: %s\n", prompt)
	}

	// Generate image
	fmt.Printf("\n🖼️  Generating image...\n")
	response, err := client.GenerateImage(context.Background(), finalPrompt, opts)
	if err != nil {
		return fmt.Errorf("failed to generate image: %w", err)
	}

	imageData := response.Data[0]

	// Save to history (non-blocking, log errors)
	historyStore := app.NewFileHistoryStore("")
	historyEntry := app.NewImageHistoryEntry(prompt, imageData, opts.Model)
	if err := historyStore.Save(historyEntry); err != nil {
		fmt.Printf("⚠️  Warning: Failed to save to history: %v\n", err)
	}

	// Display basic info
	fmt.Printf("\n✅ Image generated successfully!\n")
	if imageData.Width > 0 && imageData.Height > 0 {
		fmt.Printf("📐 Size: %dx%d\n", imageData.Width, imageData.Height)
	} else {
		fmt.Printf("📐 Size: %s\n", imageSize) // Use the requested size
	}
	fmt.Printf("🔗 URL: %s\n", imageData.URL)
	fmt.Printf("⏰ Expires: 30 days from now\n")

	// Auto-save to current directory if no output specified
	outputPath := imageOutput
	if outputPath == "" {
		// Generate filename from timestamp
		timestamp := time.Now().Format("20060102-150405")
		outputPath = fmt.Sprintf("zai-image-%s.png", timestamp)
	}

	if err := saveImageToDisk(imageData.URL, outputPath); err != nil {
		fmt.Printf("⚠️  Warning: Failed to save image: %v\n", err)
	} else {
		fmt.Printf("💾 Saved to: %s\n", outputPath)
	}

	if imageCopy {
		if err := copyToClipboard(imageData.URL); err != nil {
			fmt.Printf("⚠️  Warning: Failed to copy to clipboard: %v\n", err)
		} else {
			fmt.Printf("📋 Copied URL to clipboard\n")
		}
	}

	if imageShow {
		if err := openImageViewer(imageData.URL); err != nil {
			fmt.Printf("⚠️  Warning: Failed to open image viewer: %v\n", err)
		}
	}

	return nil
}

func runImageModelList() error {
	client := newClient()

	// Note: Using the same ListModels method as chat for now
	// In a real implementation, this might be a separate endpoint
	models, err := client.ListModels(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	fmt.Println("Available Models:")
	fmt.Println("─────────────────")

	// Show the image generation model specifically
	imageModel := viper.GetString("api.image_model")
	if imageModel == "" {
		imageModel = "cogview-4-250304"
	}

	fmt.Printf("  %s  (image generation)\n", imageModel)

	// Show other models that might support images
	for _, m := range models {
		if strings.Contains(strings.ToLower(m.ID), "image") ||
		   strings.Contains(strings.ToLower(m.ID), "cogview") ||
		   strings.Contains(strings.ToLower(m.ID), "dall-e") {
			fmt.Printf("  %s  (image capable)\n", m.ID)
		}
	}

	return nil
}

// saveImageToDisk downloads an image from URL and saves to file
func saveImageToDisk(url, filePath string) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Download the image
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	// Create output file
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy response body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write image: %w", err)
	}

	return nil
}

// copyToClipboard copies URL to clipboard (macOS only)
func copyToClipboard(url string) error {
	if _, err := exec.LookPath("pbcopy"); err != nil {
		return fmt.Errorf("pbcopy not found (macOS only)")
	}

	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(url)
	return cmd.Run()
}

// openImageViewer opens URL with default viewer
func openImageViewer(url string) error {
	var cmd *exec.Cmd
	var err error

	// Platform-specific commands
	if _, err = exec.LookPath("open"); err == nil { // macOS
		cmd = exec.Command("open", url)
	} else if _, err = exec.LookPath("xdg-open"); err == nil { // Linux
		cmd = exec.Command("xdg-open", url)
	} else if _, err = exec.LookPath("start"); err == nil { // Windows
		cmd = exec.Command("cmd", "/c", "start", url)
	} else {
		return fmt.Errorf("no suitable viewer found for this platform")
	}

	return cmd.Start()
}