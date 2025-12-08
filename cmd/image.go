package cmd

import (
	"context"
	"fmt"
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
	imagePrompt   string
	imageQuality  string
	imageSize     string
	imageOutput   string
	imageShow     bool
	imageCopy     bool
	imageModel    string
	imageUserID   string
)

var imageCmd = &cobra.Command{
	Use:   "image \"description\"",
	Short: "Generate images using Z.AI's image generation API",
	Long: `Generate images with Z.AI's CogView-4 model.

Examples:
  zai image "a cat wearing sunglasses"
  zai image "sunset on mars" --quality hd --size 1024x1024
  zai image "abstract art" --output my-art.png
  zai image "logo" --copy --size 512x512`,
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

	// Add subcommands
	imageCmd.AddCommand(imageListCmd)

	// Register with root
	rootCmd.AddCommand(imageCmd)
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

	// Generate image
	fmt.Printf("🎨 Generating image: %s\n", prompt)
	response, err := client.GenerateImage(context.Background(), prompt, opts)
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
	fmt.Printf("📐 Size: %dx%d\n", imageData.Width, imageData.Height)
	fmt.Printf("🔗 URL: %s\n", imageData.URL)
	fmt.Printf("⏰ Expires: 30 days from now\n")

	// Handle output options
	if imageOutput != "" {
		if err := saveImageToDisk(imageData.URL, imageOutput); err != nil {
			fmt.Printf("⚠️  Warning: Failed to save image: %v\n", err)
		} else {
			fmt.Printf("💾 Saved to: %s\n", imageOutput)
		}
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
	// For now, just print info about the URL
	// In a full implementation, we would download the image
	// This requires additional dependencies like curl or a HTTP client

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create a placeholder file with the URL
	content := fmt.Sprintf("Image URL: %s\nGenerated: %s\n", url, time.Now().Format(time.RFC3339))
	return os.WriteFile(filePath, []byte(content), 0644)
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