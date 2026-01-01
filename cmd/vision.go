package cmd

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/garyblankenship/zai/internal/app"
)

var (
	visionFile   string
	visionPrompt string
	visionModel  string
	visionTemp   float64
)

var visionCmd = &cobra.Command{
	Use:   "vision [prompt]",
	Short: "Analyze images with AI vision (glm-4.6v)",
	Long: `Analyze images using Z.AI's GLM-4.6V vision model.

Supports local files and HTTP/HTTPS URLs via -f flag.

Examples:
  zai vision -f photo.jpg                     # Describe image
  zai vision -f screenshot.png "What text?"   # Extract text
  zai vision -f https://example.com/img.jpg   # Analyze URL
  zai vision -f chart.png -p "Explain trends" # With prompt flag`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if visionFile == "" {
			return fmt.Errorf("image required: use -f <image-path-or-url>")
		}
		prompt := ""
		if len(args) > 0 {
			prompt = args[0]
		}
		return runVision(visionFile, prompt)
	},
}

func init() {
	visionCmd.Flags().StringVarP(&visionFile, "file", "f", "", "Image file path or URL (required)")
	visionCmd.Flags().StringVarP(&visionPrompt, "prompt", "p", "", "Analysis prompt (default: describe the image)")
	visionCmd.Flags().StringVarP(&visionModel, "model", "m", "", "Override vision model (default: glm-4.6v)")
	visionCmd.Flags().Float64VarP(&visionTemp, "temperature", "t", 0.3, "Temperature (0.0-1.0, default: 0.3)")

	// Register with root
	rootCmd.AddCommand(visionCmd)
}

func runVision(imageSource, prompt string) error {
	client := newClient()

	ctx, cancel := createContext(5 * time.Minute)
	defer cancel()

	// Use flag prompt if provided, otherwise use argument
	if visionPrompt != "" {
		prompt = visionPrompt
	}

	// Default prompt
	if prompt == "" {
		prompt = "What do you see in this image? Please provide a detailed description."
	}

	// Determine if imageSource is a URL or local file
	var imageBase64 string
	var err error

	if strings.HasPrefix(imageSource, "http://") || strings.HasPrefix(imageSource, "https://") {
		// For URLs, pass directly
		imageBase64 = imageSource
		fmt.Printf("🌐 Fetching image from URL: %s\n", imageSource)
	} else {
		// Local file - read and encode to base64
		imageBase64, err = encodeImageToBase64(imageSource)
		if err != nil {
			return fmt.Errorf("failed to read image: %w", err)
		}
		fmt.Printf("📁 Analyzing image: %s\n", imageSource)
	}

	// Build options
	opts := app.VisionOptions{
		Model:       visionModel,
		Temperature: app.Float64Ptr(visionTemp),
	}

	fmt.Printf("🔍 Analyzing with prompt: %s\n", prompt)
	fmt.Println()

	// Call vision API
	response, err := client.Vision(ctx, prompt, imageBase64, opts)
	if err != nil {
		return fmt.Errorf("vision analysis failed: %w", err)
	}

	// Output response
	fmt.Println("📝 Analysis:")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println(response)
	fmt.Println(strings.Repeat("─", 50))

	return nil
}

// encodeImageToBase64 reads an image file and returns it as a data URI.
func encodeImageToBase64(imagePath string) (string, error) {
	// Check if file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", imagePath)
	}

	// Read file
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Determine MIME type from extension
	ext := strings.ToLower(filepath.Ext(imagePath))
	var mimeType string
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	default:
		return "", fmt.Errorf("unsupported image format: %s (supported: jpg, jpeg, png, gif, webp)", ext)
	}

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(data)

	// Return as data URI
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

// downloadImage downloads an image from URL and returns base64 data URI.
// Uses provided HTTPDoer for connection pooling.
// Not currently used but available for future enhancement (download + local processing).
func downloadImage(client app.HTTPDoer, url string) (string, error) {
	// Create request using shared HTTP client
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Detect content type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg" // default
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", contentType, encoded), nil
}
