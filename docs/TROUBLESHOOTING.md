# Troubleshooting Guide

Common issues and their solutions when using the `zai` CLI tool.

## API Authentication Issues

### Symptom
```
Error: API key not configured. Set ZAI_API_KEY environment variable or configure in ~/.config/zai/config.yaml
```

### Cause
The API key is missing from both the environment variables and the configuration file.

### Solution
1. Set the environment variable:
   ```bash
   export ZAI_API_KEY="your-api-key"
   ```
2. Or add to `~/.config/zai/config.yaml`:
   ```yaml
   api:
     key: "your-api-key"
   ```

---

## Rate Limit Errors

### Symptom
```
Error: 429 Too Many Requests - Rate limit exceeded
```

### Cause
Too many API requests have been made in a short period. The API has rate limits to prevent abuse.

### Solution
1. Wait a few minutes before retrying
2. Implement exponential backoff in your scripts
3. Configure retry settings in `~/.config/zai/config.yaml`:
   ```yaml
   api:
     retry:
       max_attempts: 3
       initial_backoff: 1s
       max_backoff: 30s
   ```

---

## Network Timeout Issues

### Symptom
```
Error: context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

### Cause
Slow network connection or the API server is taking longer than expected to respond.

### Solution
1. Check your internet connection
2. Increase the timeout in your config:
   ```yaml
   web_reader:
     timeout: 30  # Increase from default 20s
   web_search:
     timeout: 60  # Increase from default 30s
   ```

---

## Invalid File Path

### Symptom
```
Error: open file.go: no such file or directory
```

### Cause
The file specified with `-f` flag does not exist or the path is incorrect.

### Solution
1. Verify the file exists:
   ```bash
   ls -la /path/to/file.go
   ```
2. Use absolute paths or ensure you're in the correct directory
3. For URLs, ensure they start with `http://` or `https://`

---

## Audio Transcription Fails

### Symptom
```
Error: audio file format not supported or file too large
```

### Cause
The audio file is in an unsupported format, exceeds 25MB, or is longer than 30 seconds per chunk.

### Solution
1. Check file size is under 25MB
2. Convert to supported format (.wav, .mp3, .mp4, .m4a, .flac, .aac, .ogg):
   ```bash
   ffmpeg -i input.mov -acodec mp3 output.mp3
   ```
3. For longer files, split into chunks:
   ```bash
   ffmpeg -i long_audio.mp3 -f segment -segment_time 30 -c copy chunk_%03d.mp3
   ```
4. Use `--vad` flag to remove silence and reduce costs

---

## Web Content Not Fetched

### Symptom
URLs in chat are not being automatically fetched and included in the context.

### Cause
Web reader is disabled in configuration or the URL is not being detected.

### Solution
1. Ensure web reader is enabled in `~/.config/zai/config.yaml`:
   ```yaml
   web_reader:
     enabled: true
     auto_detect: true
   ```
2. Verify the URL starts with `http://` or `https://`
3. Try using the `reader` command directly:
   ```bash
   zai reader https://example.com
   ```

---

## Video Generation Timeout

### Symptom
Video generation appears to hang or times out after a few minutes.

### Cause
Video generation is an async process that typically takes 1-3 minutes. The default poll timeout may be too short.

### Solution
1. Increase the poll timeout:
   ```bash
   zai video "prompt" --poll-timeout 5m
   ```
2. Wait for the animated spinner to complete
3. Check if the video file was created: `ls -la zai-video-*.mp4`
