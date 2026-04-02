package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/dotcommander/zai/internal/app"
)

var (
	ttsVoice  string
	ttsModel  string
	ttsSpeed  int
	ttsFormat string
	ttsOutput string
)

var ttsCmd = &cobra.Command{
	Use:   "tts [text]",
	Short: "Convert text to speech using Z.AI TTS",
	Long: `Convert text to speech using Z.AI's GLM-TTS model.

Supports multiple voices and output formats.

Examples:
  zai tts "Hello, world!"
  zai tts "Welcome to Z.AI" --voice xiaochen
  echo "text" | zai tts --output greeting.wav
  zai tts "Fast speech" --speed 2 --format wav`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTTS,
}

func registerTTSCmd() {
	rootCmd.AddCommand(ttsCmd)

	ttsCmd.Flags().StringVar(&ttsVoice, "voice", "", "Voice: tongtong, xiaochen, chuichui, jam, kazi, douji, luodo (default: config tts.voice)")
	ttsCmd.Flags().StringVarP(&ttsModel, "model", "m", "", "Override TTS model (default: config api.tts_model)")
	ttsCmd.Flags().IntVar(&ttsSpeed, "speed", 0, "Speech speed multiplier (default: 1)")
	ttsCmd.Flags().StringVar(&ttsFormat, "format", "", "Output format: wav or pcm (default: config tts.response_format)")
	ttsCmd.Flags().StringVarP(&ttsOutput, "output", "o", "", "Output file path (default: zai-tts-{timestamp}.wav)")
}

func runTTS(cmd *cobra.Command, args []string) error {
	var text string
	switch {
	case len(args) > 0:
		text = args[0]
	case hasStdinData():
		stdinText, err := readStdin()
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		if stdinText == "" {
			return fmt.Errorf("empty text from stdin")
		}
		text = stdinText
	default:
		return fmt.Errorf("text argument or stdin input required")
	}

	voice := ttsVoice
	if voice == "" {
		voice = getModelWithDefault("tts.voice", "tongtong")
	}

	format := ttsFormat
	if format == "" {
		format = getModelWithDefault("tts.response_format", "wav")
	}

	model := ttsModel
	if model == "" {
		model = getModelWithDefault("api.tts_model", "glm-tts")
	}

	opts := app.TTSOptions{
		Model:          model,
		Voice:          voice,
		ResponseFormat: format,
	}
	if cmd.Flags().Changed("speed") {
		opts.Speed = app.IntPtr(ttsSpeed)
	}

	outputPath := ttsOutput
	if outputPath == "" {
		timestamp := time.Now().Format("20060102-150405")
		outputPath = fmt.Sprintf("zai-tts-%s.%s", timestamp, format)
	}

	client := newClient()
	ctx, cancel := createContext(2 * time.Minute)
	defer cancel()

	fmt.Fprintf(os.Stderr, "Generating speech: voice=%s, format=%s\n", voice, format)

	audioData, err := client.TextToSpeech(ctx, text, opts)
	if err != nil {
		return fmt.Errorf("TTS failed: %w", err)
	}

	if err := os.WriteFile(outputPath, audioData, 0o644); err != nil {
		return fmt.Errorf("failed to write audio file: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Saved to: %s (%d bytes)\n", outputPath, len(audioData))

	historyStore := app.NewFileHistoryStore("")
	entry := app.NewTTSHistoryEntry(text, voice, opts.Model)
	if err := historyStore.Save(entry); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save to history: %v\n", err)
	}

	return nil
}
