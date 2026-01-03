package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/dotcommander/zai/internal/app"
	"github.com/dotcommander/zai/internal/app/utils"
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

// ImageSource represents the type of image source (URL or local file)
type ImageSource int

const (
	ImageSourceURL ImageSource = iota
	ImageSourceFile
)

// detectImageSource determines if the image source is a URL or local file
func detectImageSource(source string) ImageSource {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return ImageSourceURL
	}
	return ImageSourceFile
}

// buildVisionPrompt builds the vision prompt from user inputs
func buildVisionPrompt(userPrompt, flagPrompt, defaultPrompt string) string {
	if flagPrompt != "" {
		return flagPrompt
	}
	if userPrompt != "" {
		return userPrompt
	}
	return defaultPrompt
}

// encodeLocalImage reads and encodes a local image file to base64 data URI
func encodeLocalImage(imagePath string, fileReader utils.FileReader) (string, error) {
	data, err := fileReader.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	mimeType, err := utils.DetectImageMimeType(imagePath)
	if err != nil {
		return "", err
	}

	return utils.EncodeBytesToDataURI(data, mimeType), nil
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

	// Build the prompt using pure function
	prompt = buildVisionPrompt(prompt, visionPrompt, "What do you see in this image? Please provide a detailed description.")

	// Determine image source type and handle accordingly
	imageBase64, err := processImageSource(imageSource, client)
	if err != nil {
		return fmt.Errorf("failed to process image: %w", err)
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

// processImageSource handles URL and local image sources appropriately
func processImageSource(imageSource string, client *app.Client) (string, error) {
	sourceType := detectImageSource(imageSource)

	switch sourceType {
	case ImageSourceURL:
		fmt.Printf("🌐 Fetching image from URL: %s\n", imageSource)
		return imageSource, nil
	case ImageSourceFile:
		fmt.Printf("📁 Analyzing image: %s\n", imageSource)
		fileReader := utils.OSFileReader{}
		return encodeLocalImage(imageSource, fileReader)
	default:
		return "", fmt.Errorf("unsupported image source: %s", imageSource)
	}
}
