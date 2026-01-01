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
	"github.com/garyblankenship/zai/internal/app"
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
	imageCmd.Flags().BoolVarP(&imageCopy, "copy", "c", false, "Copy image to clipboard (macOS, Linux, Windows)")
	imageCmd.Flags().StringVarP(&imageModel, "model", "m", "", "Override default image model")
	imageCmd.Flags().StringVar(&imageUserID, "user-id", "", "User ID for analytics")
	imageCmd.Flags().BoolVarP(&imageEnhance, "enhance", "e", true, "Enhance prompt with AI before generation")
	imageCmd.Flags().BoolVar(&imageNoEnhance, "no-enhance", false, "Disable prompt enhancement")

	// Mark mutually exclusive flags
	imageCmd.MarkFlagsMutuallyExclusive("enhance", "no-enhance")

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

// enhanceImagePromptWithCtx is the context-aware version of enhanceImagePrompt.
func enhanceImagePromptWithCtx(ctx context.Context, client *app.Client, prompt string) (string, error) {
	systemPrompt := `You are an expert at creating detailed, evocative prompts for AI image generation.

## YOUR TASK
Transform the user's simple prompt into a rich, detailed image generation prompt.

## STYLE GUIDE
Use this framework: [Subject + Action] + [Environment/Setting] + [Lighting/Atmosphere] + [Technical Specs]

## EXAMPLES
Input: "a wizard"
Output: "An elderly wizard with a flowing silver beard stands atop a crystalline tower, arms raised to channel crackling blue energy, surrounded by ancient runes glowing softly on stone walls. Dramatic moonlight beams through arched windows, illuminating floating dust motes. Cinematic wide shot, photorealistic, 8K resolution, volumetric lighting."

Input: "cat in a garden"
Output: "A fluffy ginger cat lounging on a weathered wooden bench in a sun-drenched cottage garden, surrounded by cascading ivy and blooming roses. Golden hour sunlight filters through oak leaves, creating dappled patterns on fur. Macro photography style, shallow depth of field, warm color palette."

## OUTPUT RULES
- Write as vivid natural language sentences, NOT keyword lists
- Be specific and evocative - paint a picture with words
- 150-300 characters ideal for CogView
- Output ONLY the enhanced prompt - no explanations, no quotes, no prefixes`

	opts := app.ChatOptions{
		Temperature: app.Float64Ptr(0.8),
		MaxTokens:   app.IntPtr(250),
		Context: []app.Message{
			{Role: "system", Content: systemPrompt},
		},
	}

	userPrompt := fmt.Sprintf("Transform this into a cinematic image prompt: %s", prompt)
	enhanced, err := client.Chat(ctx, userPrompt, opts)
	if err != nil {
		return "", err // Return error, let caller handle fallback
	}

	// Clean up any quotes or prefixes the model might add
	result := strings.TrimSpace(enhanced)
	result = strings.Trim(result, "\"'")
	result = strings.TrimPrefix(result, "Enhanced prompt: ")
	result = strings.TrimPrefix(result, "Prompt: ")

	return result, nil
}

func enhanceImagePrompt(client *app.Client, prompt string) (string, error) {
	ctx, cancel := createContext(2 * time.Minute)
	defer cancel()
	return enhanceImagePromptWithCtx(ctx, client, prompt)
}

func runImageGeneration(prompt string) error {
	client := newClient()
	ctx, cancel := createContext(5 * time.Minute)
	defer cancel()

	// Build options and enhance prompt
	opts := buildImageOptions()
	finalPrompt := buildFinalPrompt(client, prompt)

	// Generate image
	fmt.Printf("\n🖼️  Generating image...\n")
	response, err := client.GenerateImage(ctx, finalPrompt, opts)
	if err != nil {
		return fmt.Errorf("failed to generate image: %w", err)
	}

	imageData := response.Data[0]

	// Save to history (non-blocking)
	saveToHistory(prompt, imageData, opts.Model)

	// Display and handle the result
	return displayImageResult(imageData, finalPrompt, imageSize)
}

// buildImageOptions creates image options from command line flags and config.
func buildImageOptions() app.ImageOptions {
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

	return opts
}

// buildFinalPrompt creates the final prompt by optionally enhancing the original.
func buildFinalPrompt(client *app.Client, originalPrompt string) string {
	if !shouldEnhancePrompt() {
		fmt.Printf("🎨 Generating image: %s\n", originalPrompt)
		return originalPrompt
	}

	fmt.Printf("🎨 Original: %s\n", originalPrompt)
	fmt.Printf("✨ Enhancing prompt...\n")

	enhanced, err := enhanceImagePrompt(client, originalPrompt)
	if err != nil {
		fmt.Printf("⚠️  Enhancement failed, using original: %v\n", err)
		return originalPrompt
	}

	// Combine original + enhanced for best results
	finalPrompt := originalPrompt + ". " + enhanced
	fmt.Printf("✨ Enhanced: %s\n", enhanced)
	return finalPrompt
}

// ImageResult represents the structured result of image generation.
type ImageResult struct {
	Data       app.ImageData
	Prompt     string
	Size       string
	OutputPath string
	SaveError  error
}

// ImageOutputHandler handles output operations for image results.
type ImageOutputHandler interface {
	PrintSuccess(result *ImageResult)
	PrintSaveError(err error)
	PrintCopyError(err error)
	PrintViewerError(err error)
	PrintSaveSuccess(path string)
	PrintCopySuccess()
}

// DefaultImageOutputHandler prints to stdout/stderr.
type DefaultImageOutputHandler struct{}

func (h *DefaultImageOutputHandler) PrintSuccess(result *ImageResult) {
	fmt.Printf("\n✅ Image generated successfully!\n")
	if result.Data.Width > 0 && result.Data.Height > 0 {
		fmt.Printf("📐 Size: %dx%d\n", result.Data.Width, result.Data.Height)
	} else {
		fmt.Printf("📐 Size: %s\n", result.Size)
	}
	fmt.Printf("🔗 URL: %s\n", result.Data.URL)
	fmt.Printf("⏰ Expires: 30 days from now\n")
}

func (h *DefaultImageOutputHandler) PrintSaveError(err error) {
	fmt.Printf("⚠️  Warning: Failed to save image: %v\n", err)
}

func (h *DefaultImageOutputHandler) PrintCopyError(err error) {
	fmt.Printf("⚠️  Warning: Failed to copy to clipboard: %v\n", err)
}

func (h *DefaultImageOutputHandler) PrintViewerError(err error) {
	fmt.Printf("⚠️  Warning: Failed to open image viewer: %v\n", err)
}

func (h *DefaultImageOutputHandler) PrintSaveSuccess(path string) {
	fmt.Printf("💾 Saved to: %s\n", path)
}

func (h *DefaultImageOutputHandler) PrintCopySuccess() {
	fmt.Printf("📋 Copied URL to clipboard\n")
}

// ImageOutputConfig holds configuration for image output operations.
type ImageOutputConfig struct {
	Copy   bool
	Show   bool
	Output string
}

// ProcessImageResult processes the image result and handles all output operations.
func ProcessImageResult(result *ImageResult, cfg ImageOutputConfig, handler ImageOutputHandler, saver *ImageSaver) error {
	// Print success message
	handler.PrintSuccess(result)

	// Determine output path
	outputPath := cfg.Output
	if outputPath == "" {
		timestamp := time.Now().Format("20060102-150405")
		outputPath = fmt.Sprintf("zai-image-%s.png", timestamp)
	}

	// Save to disk
	saveResult := saver.Save(result.Data.URL, outputPath)
	if saveResult.Error != nil {
		handler.PrintSaveError(saveResult.Error)
	} else {
		handler.PrintSaveSuccess(outputPath)
	}

	// Copy to clipboard
	if cfg.Copy {
		if err := copyToClipboard(result.Data.URL); err != nil {
			handler.PrintCopyError(err)
		} else {
			handler.PrintCopySuccess()
		}
	}

	// Open in viewer
	if cfg.Show {
		if err := openImageViewer(result.Data.URL); err != nil {
			handler.PrintViewerError(err)
		}
	}

	return nil
}

// displayImageResult handles displaying, saving, and opening the generated image.
func displayImageResult(imageData app.ImageData, prompt, size string) error {
	result := &ImageResult{
		Data:   imageData,
		Prompt: prompt,
		Size:   size,
	}

	cfg := ImageOutputConfig{
		Copy:   imageCopy,
		Show:   imageShow,
		Output: imageOutput,
	}

	handler := &DefaultImageOutputHandler{}
	saver := NewImageSaver(nil)

	return ProcessImageResult(result, cfg, handler, saver)
}

// saveToHistory saves the image to history store.
func saveToHistory(prompt string, imageData app.ImageData, model string) {
	historyStore := app.NewFileHistoryStore("")
	historyEntry := app.NewImageHistoryEntry(prompt, imageData, model)
	if err := historyStore.Save(historyEntry); err != nil {
		fmt.Printf("⚠️  Warning: Failed to save to history: %v\n", err)
	}
}

func runImageModelList() error {
	client := newClient()

	ctx, cancel := createContext(30 * time.Second)
	defer cancel()

	// Note: Using the same ListModels method as chat for now
	// In a real implementation, this might be a separate endpoint
	models, err := client.ListModels(ctx)
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

// ImageSaver handles saving images to disk.
type ImageSaver struct {
	httpClient app.HTTPDoer
}

// NewImageSaver creates an ImageSaver with the provided HTTP client.
func NewImageSaver(httpClient app.HTTPDoer) *ImageSaver {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &ImageSaver{httpClient: httpClient}
}

// ImageSaveResult contains the result of saving an image.
type ImageSaveResult struct {
	FilePath string
	URL      string
	Size     int64
	Error    error
}

// Save downloads an image from URL and saves to file.
func (s *ImageSaver) Save(url, filePath string) *ImageSaveResult {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return &ImageSaveResult{FilePath: filePath, URL: url, Error: fmt.Errorf("failed to create directory: %w", err)}
		}
	}

	// Download the image
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &ImageSaveResult{FilePath: filePath, URL: url, Error: fmt.Errorf("failed to create request: %w", err)}
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return &ImageSaveResult{FilePath: filePath, URL: url, Error: fmt.Errorf("failed to download image: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &ImageSaveResult{FilePath: filePath, URL: url, Error: fmt.Errorf("failed to download image: status %d", resp.StatusCode)}
	}

	// Create output file
	out, err := os.Create(filePath)
	if err != nil {
		return &ImageSaveResult{FilePath: filePath, URL: url, Error: fmt.Errorf("failed to create file: %w", err)}
	}
	defer out.Close()

	// Copy response body to file
	size, err := io.Copy(out, resp.Body)
	if err != nil {
		return &ImageSaveResult{FilePath: filePath, URL: url, Error: fmt.Errorf("failed to write image: %w", err)}
	}

	return &ImageSaveResult{FilePath: filePath, URL: url, Size: size, Error: nil}
}

// saveImageToDisk downloads an image from URL and saves to file.
// Uses provided HTTPDoer for connection pooling.
func saveImageToDisk(client app.HTTPDoer, url, filePath string) error {
	saver := NewImageSaver(client)
	result := saver.Save(url, filePath)
	return result.Error
}

// copyToClipboard copies URL to clipboard (macOS, Linux, Windows)
func copyToClipboard(url string) error {
	var cmd *exec.Cmd
	var err error

	// Platform-specific commands
	if _, err = exec.LookPath("pbcopy"); err == nil { // macOS
		cmd = exec.Command("pbcopy")
	} else if _, err = exec.LookPath("xclip"); err == nil { // Linux
		cmd = exec.Command("xclip", "-selection", "clipboard")
	} else if _, err = exec.LookPath("xsel"); err == nil { // Linux (alternative)
		cmd = exec.Command("xsel", "--clipboard", "--input")
	} else if _, err = exec.LookPath("clip"); err == nil { // Windows
		cmd = exec.Command("clip")
	} else {
		return fmt.Errorf("no suitable clipboard tool found (requires: pbcopy/macOS, xclip/xsel/Linux, or clip/Windows)")
	}

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