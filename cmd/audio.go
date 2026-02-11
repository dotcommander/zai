package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/dotcommander/zai/internal/app"
)

// closeFile closes a file and logs any error.
func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to close file: %v\n", err)
	}
} //nolint:errcheck // error is already handled in the function

var (
	audioFile     string
	audioModel    string
	audioPrompt   string
	audioLanguage string
	audioHotwords string
	audioStream   bool
	audioJSON     bool
	audioUserID   string
	// Preprocessing options
	audioVAD        bool   // Voice Activity Detection - remove silence
	audioVideo      string // YouTube video URL to transcribe
	audioPreprocess bool   // Auto-convert to optimal format (16kHz mono WAV)
	// Cache options
	audioResume     bool // Resume from previous partial transcription
	audioClearCache bool // Clear cached transcription and start fresh
)

var audioCmd = &cobra.Command{
	Use:   "audio",
	Short: "Transcribe audio files to text",
	Long: `Transcribe audio files to text using Z.AI's GLM-ASR-2512 model.

Examples:
  zai audio -f recording.wav
  zai audio -f speech.mp3 --model glm-asr-2512
  zai audio -f interview.wav --prompt "Previous context"
  zai audio -f lecture.wav --hotwords "kubernetes,docker"
  zai audio --video https://youtu.be/abc123  # YouTube support
  zai audio -f recording.wav --vad  # Remove silence
  zai audio -f recording.wav --resume  # Resume partial transcription
  cat audio.wav | zai audio  # From stdin

Supported formats: .wav, .mp3, .mp4, .m4a, .flac, .aac, .ogg
Maximum file size: 25MB
Maximum duration: 30 seconds per chunk`,
	Args: cobra.NoArgs,
	RunE: runAudioTranscription,
}

func init() {
	rootCmd.AddCommand(audioCmd)

	audioCmd.Flags().StringVarP(&audioFile, "file", "f", "", "Audio file path")
	audioCmd.Flags().StringVarP(&audioModel, "model", "m", "glm-asr-2512", "ASR model to use")
	audioCmd.Flags().StringVarP(&audioPrompt, "prompt", "p", "", "Context from prior transcriptions (max 8000 chars)")
	audioCmd.Flags().StringVarP(&audioLanguage, "language", "l", "", "Language code (e.g., en, zh, ja)")
	audioCmd.Flags().StringVar(&audioHotwords, "hotwords", "", "Comma-separated domain vocabulary (max 100 items)")
	audioCmd.Flags().BoolVar(&audioStream, "stream", false, "Enable streaming transcription")
	audioCmd.Flags().BoolVar(&audioJSON, "json", false, "Output in JSON format")
	audioCmd.Flags().StringVar(&audioUserID, "user-id", "", "User ID for analytics (6-128 characters)")
	// Preprocessing flags
	audioCmd.Flags().BoolVar(&audioVAD, "vad", false, "Apply Voice Activity Detection to remove silence (reduces API costs)")
	audioCmd.Flags().StringVar(&audioVideo, "video", "", "YouTube video URL to transcribe")
	audioCmd.Flags().BoolVar(&audioPreprocess, "preprocess", true, "Auto-convert audio to optimal format (16kHz mono WAV)")
	// Cache flags
	audioCmd.Flags().BoolVar(&audioResume, "resume", false, "Resume from previous partial transcription")
	audioCmd.Flags().BoolVar(&audioClearCache, "clear-cache", false, "Clear cached transcription and start fresh")
}

// sanitizePath validates and cleans a file path to prevent command injection.
func sanitizePath(path string) (string, error) {
	// Clean the path to remove any . and .. components
	cleanPath := filepath.Clean(path)

	// Check if the path exists and is a regular file
	info, err := os.Stat(cleanPath)
	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("path is not a regular file: %s", cleanPath)
	}

	// Check for suspicious patterns that could indicate command injection
	const suspiciousChars = "|;&$()`><{}*?'\"\\"
	if strings.ContainsAny(cleanPath, suspiciousChars) {
		return "", fmt.Errorf("path contains invalid characters: %s", cleanPath)
	}

	return cleanPath, nil
}

// checkFFmpeg verifies ffmpeg is installed before audio processing.
func checkFFmpeg() error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg required for audio processing\n  Install: brew install ffmpeg (macOS) | apt install ffmpeg (Linux) | choco install ffmpeg (Windows)")
	}
	return nil
}

func runAudioTranscription(cmd *cobra.Command, args []string) error { //nolint:gocognit,gocyclo // TODO: decompose into smaller functions
	// Use extended timeout for large audio files (10 min for long recordings)
	ctx, cancel := createContext(10 * time.Minute)
	defer cancel()

	// Setup temporary file management
	tempMgr := &TempFileManager{}
	defer tempMgr.Cleanup()

	// Determine audio source and get audio path
	audioPath, err := determineAudioSource()
	if err != nil {
		return err
	}
	tempMgr.Add(audioPath)

	// Validate file exists and get original source path
	originalSource, err := validateAndGetAudioPath(audioPath)
	if err != nil {
		return err
	}

	// Preprocess audio if needed
	audioPath, err = preprocessAudioIfNeeded(audioPath, tempMgr)
	if err != nil {
		return err
	}

	// Handle large files by chunking
	if shouldChunkFile(audioPath) {
		return handleLargeAudioFile(ctx, audioPath, originalSource, tempMgr)
	}

	// Perform regular transcription for normal-sized files
	return performRegularTranscription(ctx, audioPath, originalSource)
}

// determineAudioSource determines the audio source (YouTube, file, or stdin) and returns the path.
func determineAudioSource() (string, error) {
	switch {
	case audioVideo != "":
		// YouTube source
		ytPath, err := downloadYouTubeAudio(audioVideo)
		if err != nil {
			return "", fmt.Errorf("YouTube download failed: %w", err)
		}
		return ytPath, nil
	case audioFile == "-" || (audioFile == "" && hasStdinData()):
		// Explicit -f - or auto-detected stdin
		// Read from stdin and write to temp file
		stdinPath, _, err := createTempAudioFile()
		if err != nil {
			return "", fmt.Errorf("failed to create temp file: %w", err)
		}
		return stdinPath, nil
	default:
		return "", fmt.Errorf("audio file required: use -f <file> or --video <youtube_url>, or pipe via stdin")
	}
}

// validateAndGetAudioPath validates the audio file exists and returns the original source path.
func validateAndGetAudioPath(audioPath string) (string, error) {
	// Validate file exists
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		return "", fmt.Errorf("audio file not found: %s", audioPath)
	}

	// Save original source path for cache key (before preprocessing)
	return audioPath, nil
}

// preprocessAudioIfNeeded preprocesses audio if needed and returns the final audio path.
func preprocessAudioIfNeeded(audioPath string, tempMgr *TempFileManager) (string, error) {
	// Check ffmpeg before any processing that requires it
	needsFFmpeg := audioPreprocess || audioVAD
	if needsFFmpeg {
		if err := checkFFmpeg(); err != nil {
			return "", err
		}
	}

	// Preprocessing: convert to optimal format if needed
	if audioPreprocess || audioVAD {
		processedPath, err := preprocessAudio(audioPath, audioVAD)
		if err != nil {
			return "", fmt.Errorf("audio preprocessing failed: %w", err)
		}
		if processedPath != audioPath {
			tempMgr.Add(processedPath)
			return processedPath, nil
		}
	}

	return audioPath, nil
}

// shouldChunkFile checks if the audio file should be chunked based on size.
func shouldChunkFile(audioPath string) bool {
	info, err := os.Stat(audioPath)
	if err != nil {
		return false
	}
	const maxFileSize = 25 * 1024 * 1024
	return info.Size() > maxFileSize
}

// handleLargeAudioFile handles large audio files by chunking them.
func handleLargeAudioFile(ctx context.Context, audioPath, originalSource string, tempMgr *TempFileManager) error {
	// Check ffmpeg for splitting (required even if preprocessing was skipped)
	if err := checkFFmpeg(); err != nil {
		return err
	}

	// Try to chunk the file
	info, err := os.Stat(audioPath)
	if err != nil {
		return fmt.Errorf("failed to access audio file: %w", err)
	}
	fmt.Fprintf(os.Stderr, "File too large (%d MB), splitting into chunks...\n", info.Size()/1024/1024)
	chunks, err := splitAudio(audioPath, 25) // 25-second chunks (API limit 30s)
	if err != nil {
		return fmt.Errorf("failed to chunk audio: %w", err)
	}
	tempMgr.AddAll(chunks)

	// Create client once for all chunk processing
	client := newClientWithoutHistory()

	// Transcribe each chunk and combine
	return transcribeChunks(ctx, client, chunks, originalSource, audioPath)
}

// performRegularTranscription performs transcription for normal-sized audio files.
func performRegularTranscription(ctx context.Context, audioPath, originalSource string) error {
	// Create client
	client := newClientWithoutHistory()

	// Build transcription options
	opts := buildTranscriptionOptions()

	// Perform transcription
	resp, err := client.TranscribeAudio(ctx, audioPath, opts)
	if err != nil {
		return fmt.Errorf("transcription failed: %w", err)
	}

	// Output results
	outputTranscriptionResult(resp)

	// Save to history (non-blocking)
	saveAudioToHistory(resp)

	return nil
}

// buildTranscriptionOptions builds the transcription options from command flags.
func buildTranscriptionOptions() app.TranscriptionOptions {
	opts := app.TranscriptionOptions{
		Model:    audioModel,
		Prompt:   audioPrompt,
		Stream:   audioStream,
		UserID:   audioUserID,
		Hotwords: parseHotwords(audioHotwords),
	}

	// Handle language via prompt if provided
	if audioLanguage != "" {
		if opts.Prompt != "" {
			opts.Prompt = "Language: " + audioLanguage + ". " + opts.Prompt
		} else {
			opts.Prompt = "Language: " + audioLanguage
		}
	}

	return opts
}

// outputTranscriptionResult outputs the transcription result in the requested format.
func outputTranscriptionResult(resp *app.TranscriptionResponse) {
	if audioJSON {
		output := map[string]interface{}{
			"id":      resp.ID,
			"model":   resp.Model,
			"text":    resp.Text,
			"created": resp.Created,
		}
		data, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to marshal JSON: %v\n", err)
			return
		}
		fmt.Println(string(data))
	} else {
		fmt.Println(resp.Text)
	}
}

// saveAudioToHistory saves the transcription result to history.
func saveAudioToHistory(resp *app.TranscriptionResponse) {
	history := app.NewFileHistoryStore("")
	entry := app.NewAudioHistoryEntry(resp.Text, resp.Model)
	if err := history.Save(entry); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to save to history: %v\n", err)
	}
}

// AudioCache stores partial transcription results for resume support.
type AudioCache struct {
	Chunks map[int]string `json:"chunks"` // chunk index -> transcribed text
}

// getCachePath returns the cache file path for a given source file.
func getCachePath(sourcePath string) (string, error) {
	file, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer closeFile(file)

	limitedReader := io.LimitReader(file, MaxFileSize)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", err
	}

	// Check if file exceeds our limit
	fileInfo, err := os.Stat(sourcePath)
	if err != nil {
		return "", err
	}
	if fileInfo.Size() > MaxFileSize {
		return "", fmt.Errorf("audio cache file exceeds maximum size of %d bytes", MaxFileSize)
	}

	hash := sha256.Sum256(data)
	hashStr := fmt.Sprintf("%x", hash[:8])

	cacheDir := filepath.Join(os.Getenv("HOME"), ".cache", "zai", "audio")
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, hashStr+".json"), nil
}

// loadCache loads cached transcription results.
func loadCache(cachePath string) (*AudioCache, error) {
	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &AudioCache{Chunks: make(map[int]string)}, nil
		}
		return nil, err
	}

	var cache AudioCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	if cache.Chunks == nil {
		cache.Chunks = make(map[int]string)
	}
	return &cache, nil
}

// saveCache saves transcription cache to disk.
func saveCache(cachePath string, cache *AudioCache) error {
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cachePath, data, 0600)
}

// chunkResult holds the result of transcribing a single chunk.
type chunkResult struct {
	index int
	text  string
	err   error
}

// transcribeChunks transcribes multiple audio chunks with caching, resume, and parallel processing.
func transcribeChunks(ctx context.Context, client *app.Client, chunks []string, cacheSourcePath, audioPath string) error { //nolint:gocognit,gocyclo // TODO: decompose into smaller functions
	// Get cache path using original source file for consistent cache keys
	cachePath, err := getCachePath(cacheSourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not determine cache path: %v\n", err)
	}

	var cache *AudioCache
	if cachePath != "" && !audioClearCache {
		cache, err = loadCache(cachePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not load cache: %v\n", err)
			cache = &AudioCache{Chunks: make(map[int]string)}
		}
	} else {
		cache = &AudioCache{Chunks: make(map[int]string)}
	}

	// Clear cache if requested
	if audioClearCache && cachePath != "" {
		if err := os.Remove(cachePath); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Warning: Could not clear cache: %v\n", err)
		}
		fmt.Fprintf(os.Stderr, "Cache cleared.\n")
		cache = &AudioCache{Chunks: make(map[int]string)}
	}

	// Find chunks that need transcription (resume support)
	pending := []int{}
	for i := range chunks {
		if _, ok := cache.Chunks[i]; !ok {
			pending = append(pending, i)
		}
	}

	allDone := len(pending) == 0
	if allDone {
		fmt.Fprintf(os.Stderr, "All %d chunks already transcribed (from cache)\n", len(chunks))
	} else {
		fmt.Fprintf(os.Stderr, "Processing %d chunks in parallel...\n", len(pending))
	}

	// Process pending chunks in parallel
	if !allDone { //nolint:nestif // TODO: reduce nesting
		results := transcribeParallel(ctx, client, chunks, pending)
		for res := range results {
			if res.err != nil {
				if cachePath != "" {
					_ = saveCache(cachePath, cache) // Best effort save on error
				}
				return fmt.Errorf("chunk %d failed: %w", res.index+1, res.err)
			}
			cache.Chunks[res.index] = res.text
			if cachePath != "" {
				if err := saveCache(cachePath, cache); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Could not save cache: %v\n", err)
				}
			}
		}
	}

	// Assemble final text in order
	var fullText string
	for i := range chunks {
		if text, ok := cache.Chunks[i]; ok {
			if fullText != "" {
				fullText += "\n"
			}
			fullText += text
		}
	}

	// Output results
	if audioJSON {
		output := map[string]interface{}{
			"model": audioModel,
			"text":  fullText,
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println(fullText)
	}

	return nil
}

// transcribeParallel processes chunks concurrently using a worker pool.
// Client is shared across workers for connection pooling.
func transcribeParallel(ctx context.Context, client *app.Client, chunks []string, pendingIndices []int) <-chan chunkResult { //nolint:gocognit // TODO: decompose into smaller functions
	numWorkers := 5
	results := make(chan chunkResult, len(pendingIndices))
	jobs := make(chan int, len(pendingIndices))

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			opts := app.TranscriptionOptions{Model: audioModel, Prompt: audioPrompt}

			for idx := range jobs {
				var resp *app.TranscriptionResponse
				var err error

				// Retry with exponential backoff + jitter (matches Chat pattern)
				for attempt := 1; attempt <= 3; attempt++ {
					resp, err = client.TranscribeAudio(ctx, chunks[idx], opts)
					if err == nil {
						break
					}
					if attempt < 3 {
						// Exponential backoff: 1s, 2s, 4s
						backoff := time.Second * time.Duration(1<<uint(attempt-1)) //nolint:gosec // G115: attempt count is small, overflow impossible
						// Add jitter ±12.5%
						jitter := time.Duration(float64(backoff) * 0.125 * (2*rand.Float64() - 1)) //nolint:gosec // G404: jitter doesn't need crypto-grade randomness
						time.Sleep(backoff + jitter)
					}
				}

				if err != nil {
					results <- chunkResult{index: idx, err: err}
				} else {
					results <- chunkResult{index: idx, text: resp.Text}
				}
			}
		}(w)
	}

	go func() {
		for _, idx := range pendingIndices {
			jobs <- idx
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// preprocessAudio converts audio to optimal format and optionally applies VAD.
func preprocessAudio(inputPath string, applyVAD bool) (string, error) {
	// Sanitize input path to prevent command injection
	sanitizedPath, err := sanitizePath(inputPath)
	if err != nil {
		return "", fmt.Errorf("input path validation failed: %w", err)
	}

	// Check if already optimal WAV
	ext := strings.ToLower(filepath.Ext(sanitizedPath))
	if ext == ".wav" && !applyVAD {
		return sanitizedPath, nil
	}

	tempDir := os.TempDir()
	outputPath := filepath.Join(tempDir, fmt.Sprintf("zai-audio-processed-%d.wav", time.Now().UnixNano()))

	args := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-y", "-i", inputPath,
		"-vn",                  // No video
		"-acodec", "pcm_s16le", // 16-bit PCM
		"-ar", "16000", // 16kHz sample rate (optimal for speech)
		"-ac", "1", // Mono
	}

	// Apply VAD filter if requested
	if applyVAD {
		args = append(args, "-af", "silenceremove=start_periods=1:start_duration=1:start_threshold=-50dB:detection=peak,aformat=dblp,areverse,silenceremove=start_periods=1:start_duration=1:start_threshold=-50dB:detection=peak,aformat=dblp,areverse")
	}

	args = append(args, outputPath)

	cmd := exec.Command("ffmpeg", args...) //nolint:gosec // G204: ffmpeg binary is hardcoded, args are controlled
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg failed: %w (is ffmpeg installed?)", err)
	}

	return outputPath, nil
}

// splitAudio splits an audio file into chunks using ffmpeg.
func splitAudio(inputPath string, chunkDuration int) ([]string, error) {
	// Sanitize input path to prevent command injection
	sanitizedPath, err := sanitizePath(inputPath)
	if err != nil {
		return nil, fmt.Errorf("input path validation failed: %w", err)
	}

	tempDir := os.TempDir()
	chunkPattern := filepath.Join(tempDir, fmt.Sprintf("zai-chunk-%d-%%03d.wav", os.Getpid()))

	args := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-i", sanitizedPath,
		"-f", "segment",
		"-segment_time", fmt.Sprintf("%d", chunkDuration),
		"-c", "copy",
		chunkPattern,
	}

	cmd := exec.Command("ffmpeg", args...) //nolint:gosec // G204: ffmpeg binary is hardcoded, args are controlled
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to split audio: %w", err)
	}

	// Find generated chunks
	basePattern := strings.Replace(chunkPattern, "%03d", "*", 1)
	chunks, err := filepath.Glob(basePattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find chunks: %w", err)
	}

	if len(chunks) == 0 {
		return nil, fmt.Errorf("no chunks generated")
	}

	return chunks, nil
}

// downloadYouTubeAudio downloads audio from a YouTube video using yt-dlp.
func downloadYouTubeAudio(url string) (string, error) {
	// Check if yt-dlp is available
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return "", fmt.Errorf("yt-dlp not found (required for YouTube): %w", err)
	}

	tempDir := os.TempDir()
	outputPath := filepath.Join(tempDir, fmt.Sprintf("zai-yt-audio-%d.%%(ext)s", time.Now().UnixNano()))

	args := []string{
		"-x",                    // Extract audio
		"--audio-format", "wav", // Convert to WAV
		"--audio-quality", "0", // Best quality
		"-o", outputPath,
		url,
	}

	cmd := exec.Command("yt-dlp", args...) //nolint:gosec // G204: yt-dlp binary is hardcoded, args are controlled
	cmd.Stdout = os.Stderr                 // yt-dlp progress to stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("yt-dlp failed: %w", err)
	}

	// Find the downloaded file (replace %(ext)s with actual extension)
	globPattern := strings.Replace(outputPath, "%(ext)s", "*", 1)
	matches, err := filepath.Glob(globPattern)
	if err != nil || len(matches) == 0 {
		return "", fmt.Errorf("yt-dlp did not produce any audio file")
	}

	return matches[0], nil
}

// createTempAudioFile reads stdin and writes to a temp file.
func createTempAudioFile() (string, func(), error) {
	limitedReader := io.LimitReader(os.Stdin, AudioFileSize)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read stdin: %w", err)
	}
	if len(data) == 0 {
		return "", nil, fmt.Errorf("no audio data in stdin")
	}
	if len(data) == AudioFileSize {
		return "", nil, fmt.Errorf("audio data exceeds maximum size of %d bytes", AudioFileSize)
	}

	// Create temp file with .wav extension
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, fmt.Sprintf("zai-audio-stdin-%d.wav", time.Now().UnixNano()))

	if err := os.WriteFile(tempFile, data, 0600); err != nil {
		return "", nil, fmt.Errorf("failed to write temp file: %w", err)
	}

	cleanup := func() {
		_ = os.Remove(tempFile)
	}

	return tempFile, cleanup, nil
}

// TempFileManager tracks temporary files for cleanup.
type TempFileManager struct {
	files []string
}

// Add registers a file for cleanup.
func (m *TempFileManager) Add(path string) {
	m.files = append(m.files, path)
}

// AddAll registers multiple files for cleanup.
func (m *TempFileManager) AddAll(paths []string) {
	m.files = append(m.files, paths...)
}

// Cleanup removes all registered files.
func (m *TempFileManager) Cleanup() {
	for _, f := range m.files {
		_ = os.Remove(f)
	}
}

// parseHotwords parses comma-separated hotwords into a slice.
func parseHotwords(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	hotwords := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			hotwords = append(hotwords, p)
		}
	}
	// Limit to 100 items
	if len(hotwords) > 100 {
		hotwords = hotwords[:100]
	}
	return hotwords
}
