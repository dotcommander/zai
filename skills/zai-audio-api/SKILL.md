---
name: zai-audio-api
description: Transcribes audio files to text using Z.AI GLM-ASR-2512 model. Use when user mentions audio transcription, speech-to-text, or converting audio files to text.
license: MIT
---

# Z.AI Audio Transcriptions API

## Quick Reference

| Intent | Action |
|--------|--------|
| "Transcribe audio file" | POST to `/api/paas/v4/audio/transcriptions` |
| "Streaming transcription" | Set `stream: true` in request |
| "Add domain vocabulary" | Use `hotwords` array (max 100 items) |
| "Context-aware transcription" | Use `prompt` for prior context |
| "Preprocess audio" | Convert to 16kHz mono PCM WAV with ffmpeg |

## Endpoint

```
POST https://api.z.ai/api/paas/v4/audio/transcriptions
```

**Authentication:** `Authorization: Bearer <token>` header required.

## Request Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file` | file | Yes* | Audio file, max 25MB, 30 seconds |
| `file_base64` | string | Yes* | Base64 encoded audio (use instead of file) |
| `model` | enum | No | Model ID, default: `glm-asr-2512` |
| `prompt` | string | No | Context from prior transcriptions, max 8000 chars |
| `hotwords` | string[] | No | Domain vocabulary, max 100 items |
| `stream` | boolean | No | Enable streaming via Event Stream, default: `false` |
| `request_id` | string | No | Client-provided unique identifier |
| `user_id` | string | No | End user ID, 6-128 characters |

*Either `file` or `file_base64` is required.

## Response Format

**Content-Type:** `application/json`

```json
{
  "id": "<string>",
  "created": 1234567890,
  "request_id": "<string>",
  "model": "<string>",
  "text": "<transcribed text>"
}
```

## Configuration Options

### Supported Audio Formats
- `.wav` - WAV audio
- `.mp3` - MP3 audio
- `.mp4` - MP4 video (audio extracted)
- `.m4a` - M4A audio
- `.flac` - FLAC audio
- `.aac` - AAC audio
- `.ogg` - OGG audio

### Constraints
- Maximum file size: 25MB
- Maximum audio duration: 30 seconds

### Response Formats
Only `application/json` response format is supported.

## Audio Preprocessing Best Practices

For optimal transcription quality, preprocess audio using ffmpeg:

### Optimal Format Settings

| Setting | Value | Rationale |
|---------|-------|-----------|
| Sample Rate | 16000 Hz | Optimal for speech recognition |
| Codec | PCM 16-bit LE (`pcm_s16le`) | Uncompressed, widely compatible |
| Channels | 1 (Mono) | Speech is mono; reduces file size |
| Format | WAV | Preserves quality |

### FFmpeg Conversion Command

```bash
ffmpeg -hide_banner -loglevel error -y -i input.mp3 \
  -vn -acodec pcm_s16le -ar 16000 -ac 1 output.wav
```

### Voice Activity Detection (VAD)

Remove silence to reduce API costs and improve accuracy:

```bash
ffmpeg -hide_banner -loglevel error -y -i input.mp3 \
  -vn -acodec pcm_s16le -ar 16000 -ac 1 \
  -af "silenceremove=start_periods=1:start_duration=1:start_threshold=-50dB:detection=peak,aformat=dblp,areverse,silenceremove=start_periods=1:start_duration=1:start_threshold=-50dB:detection=peak,aformat=dblp,areverse" \
  output.wav
```

### Handling Large Files

For files >25MB or >30 seconds:

1. **Split by duration**: Use 25-second segments
2. **Transcribe each chunk**: Process sequentially
3. **Combine results**: Join with newlines

```bash
# Split into 25-second chunks
ffmpeg -hide_banner -loglevel error -i long_audio.mp3 \
  -f segment -segment_time 25 -c copy chunk-%03d.wav

# Transcribe each chunk and combine
for chunk in chunk-*.wav; do
  zai audio "$chunk"
done
```

## CLI Usage Examples

```bash
# Basic transcription
zai audio recording.wav

# With model specification
zai audio speech.mp3 --model glm-asr-2512

# With context prompt
zai audio interview.wav --prompt "Previous context about the interview topic"

# With domain vocabulary
zai audio lecture.wav --hotwords "kubernetes,docker,microservices"

# YouTube video transcription
zai audio --video https://youtu.be/abc123

# Remove silence with VAD
zai audio recording.wav --vad

# JSON output
zai audio recording.wav --json
```

## Example Usage

### cURL Request
```bash
curl --request POST \
  --url https://api.z.ai/api/paas/v4/audio/transcriptions \
  --header 'Authorization: Bearer ZAI_API_KEY' \
  --header 'Content-Type: multipart/form-data' \
  --form 'model=glm-asr-2512' \
  --form 'stream=false' \
  --form 'file=@/path/to/audio.wav'
```

### Go HTTP Client
```go
func TranscribeAudio(ctx context.Context, client *http.Client, apiKey string, audioPath string) (*TranscriptionResponse, error) {
    file, err := os.Open(audioPath)
    if err != nil {
        return nil, fmt.Errorf("open audio file: %w", err)
    }
    defer file.Close()

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
    if err != nil {
        return nil, fmt.Errorf("create form file: %w", err)
    }
    if _, err := io.Copy(part, file); err != nil {
        return nil, fmt.Errorf("copy file: %w", err)
    }
    writer.WriteField("model", "glm-asr-2512")
    writer.Close()

    req, err := http.NewRequestWithContext(ctx, "POST",
        "https://api.z.ai/api/paas/v4/audio/transcriptions", body)
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("api request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("api error: %s", string(bodyBytes))
    }

    var result TranscriptionResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("decode response: %w", err)
    }
    return &result, nil
}
```

### With Hotwords
```bash
curl --request POST \
  --url https://api.z.ai/api/paas/v4/audio/transcriptions \
  --header 'Authorization: Bearer ZAI_API_KEY' \
  --header 'Content-Type: multipart/form-data' \
  --form 'model=glm-asr-2512' \
  --form 'hotwords=["kubernetes","microservices","api gateway"]' \
  --form 'file=@/path/to/audio.wav'
```

### Streaming Response
```bash
curl --request POST \
  --url https://api.z.ai/api/paas/v4/audio/transcriptions \
  --header 'Authorization: Bearer ZAI_API_KEY' \
  --header 'Content-Type: multipart/form-data' \
  --form 'model=glm-asr-2512' \
  --form 'stream=true' \
  --form 'file=@/path/to/audio.wav'
```

## Anti-Patterns

| Anti-Pattern | Problem | Fix |
|--------------|---------|-----|
| Sending file >25MB | Request fails with size error | Split audio or compress first |
| Audio >30 seconds | Transcription incomplete | Trim audio to 30s segments |
| Missing Authorization header | 401 Unauthorized error | Add `Bearer <token>` header |
| Wrong Content-Type | Request parsing fails | Use `multipart/form-data` for file upload |
| Base64 without file param | Duplicate parameter error | Use ONLY one of file or file_base64 |
| High sample rate audio | Larger file size, no quality gain | Downsample to 16kHz |
| Stereo audio | Unnecessary file size | Convert to mono |
| Audio with long silences | Wasted API quota | Apply VAD filter |

## Success Criteria

- [ ] Authorization header properly set with Bearer token
- [ ] Either file or file_base64 provided (not both)
- [ ] Audio file within 25MB size limit
- [ ] Audio duration within 30 seconds limit
- [ ] Response parsed as JSON with text field
- [ ] Hotwords array contains at most 100 items
- [ ] Prompt context does not exceed 8000 characters
- [ ] Audio preprocessed to optimal format (16kHz mono WAV)

## Integration Notes

- Use environment variable `ZAI_API_KEY` for authentication
- The GLM-ASR-2512 model supports multilingual transcription
- Hotwords improve recognition of domain-specific terminology
- Prompt helps with context from previous transcriptions for consistency
- Streaming mode uses Event Stream format for real-time results
- For production use, always preprocess audio with ffmpeg
- VAD filtering recommended for long recordings with pauses