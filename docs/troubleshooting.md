When something goes wrong, this page covers the most common issues and how to fix them. Each entry lists the symptom, cause, and solution.

- [API Authentication](#api-authentication)
- [Rate Limiting](#rate-limiting)
- [Network Timeouts](#network-timeouts)
- [Invalid File Path](#invalid-file-path)
- [Audio Transcription Failures](#audio-transcription-failures)
- [Missing Dependencies](#missing-dependencies)
- [Web Content Not Fetched](#web-content-not-fetched)
- [Video Generation Timeout](#video-generation-timeout)
- [Circuit Breaker Tripped](#circuit-breaker-tripped)
- [TTS and Embedding Errors](#tts-and-embedding-errors)
- [History File Corruption](#history-file-corruption)

<a name="api-authentication"></a>
## API Authentication

**Symptom**: zai reports that the API key is not configured.

**Cause**: The key is missing from both the environment and the config file.

**Solution**: Set the `ZAI_API_KEY` environment variable or add `api.key` to `~/.config/zai/config.yaml`. For details, see the [installation documentation](/docs/installation#configuring-your-api-key).

<a name="rate-limiting"></a>
## Rate Limiting

**Symptom**: Requests fail with a 429 status code.

**Cause**: Too many API requests in a short period.

**Solution**: Wait a few minutes and retry. If you are running batch scripts, increase the rate limit in your config or add delays between calls. For rate limit configuration, see the [configuration documentation](/docs/configuration#rate-limiting).

<a name="network-timeouts"></a>
## Network Timeouts

**Symptom**: Requests fail with "context deadline exceeded."

**Cause**: Slow network or the API is taking longer than expected to respond.

**Solution**: Check your internet connection. Increase the timeout in your config under `web_reader.timeout` or `web_search.timeout`. For timeout settings, see the [configuration documentation](/docs/configuration#web-reader-settings).

<a name="invalid-file-path"></a>
## Invalid File Path

**Symptom**: zai reports "no such file or directory."

**Cause**: The file specified with `-f` does not exist or the path is wrong.

**Solution**: Verify the file exists. Use absolute paths or confirm you are in the correct directory. For URLs, ensure they start with `http://` or `https://`.

<a name="audio-transcription-failures"></a>
## Audio Transcription Failures

**Symptom**: Audio transcription fails with a format or size error.

**Cause**: The file is in an unsupported format, exceeds 25MB, or preprocessing failed.

**Solution**: Convert to a supported format (WAV, MP3, MP4, M4A, FLAC, AAC, OGG). Ensure the file is under 25MB -- zai splits longer files automatically, but the individual file must fit within the limit. If preprocessing is failing, install ffmpeg. For supported formats, see the [audio documentation](/docs/audio#supported-formats).

<a name="missing-dependencies"></a>
## Missing Dependencies

**Symptom**: Audio or YouTube features fail with "command not found."

**Cause**: Optional external tools are not installed.

**Solution**: The audio command requires [ffmpeg](https://ffmpeg.org/) for preprocessing (enabled by default) and [yt-dlp](https://github.com/yt-dlp/yt-dlp) for YouTube transcription. Install them with your package manager:

```bash
brew install ffmpeg yt-dlp
```

Alternatively, disable preprocessing with `--preprocess=false` if you only need ffmpeg for that purpose.

<a name="web-content-not-fetched"></a>
## Web Content Not Fetched

**Symptom**: URLs in chat prompts are not being fetched automatically.

**Cause**: The web reader is disabled or auto-detection is turned off.

**Solution**: Ensure `web_reader.enabled` and `web_reader.auto_detect` are both `true` in your config. You may also fetch content manually with `zai reader <url>`. For reader settings, see the [configuration documentation](/docs/configuration#web-reader-settings).

<a name="video-generation-timeout"></a>
## Video Generation Timeout

**Symptom**: Video generation times out before completing.

**Cause**: Video generation is asynchronous and typically takes 1-3 minutes. The default poll timeout may be too short for complex prompts or high-quality settings.

**Solution**: Increase the poll timeout:

```bash
zai video "prompt" --poll-timeout 5m
```

If the process was interrupted, check whether the video file was already downloaded by looking for `zai-video-*.mp4` in the current directory.

<a name="circuit-breaker-tripped"></a>
## Circuit Breaker Tripped

**Symptom**: All requests fail immediately without reaching the API.

**Cause**: The circuit breaker opened after too many consecutive failures. This protects both your client and the API from cascading errors.

**Solution**: Wait for the circuit breaker timeout to elapse (default: 60 seconds), then retry. The circuit transitions to half-open and allows a few test requests through. If those succeed, normal operation resumes. You may also disable the circuit breaker temporarily in your config. For circuit breaker settings, see the [configuration documentation](/docs/configuration#circuit-breaker).

<a name="tts-and-embedding-errors"></a>
## TTS and Embedding Errors

**Symptom**: TTS or embedding requests fail with a model or format error.

**Cause**: The model name is incorrect, the input exceeds the model's limits, or the API key lacks access to these endpoints.

**Solution**: Verify your config uses the correct model names (`glm-tts` for TTS, `embedding-3` for embeddings). For TTS, ensure the voice name is one of the supported options. For embeddings, check that your input text is not empty. For model configuration, see the [configuration documentation](/docs/configuration#api-settings).

<a name="history-file-corruption"></a>
## History File Corruption

**Symptom**: The `zai history` command fails or shows garbled output.

**Cause**: The history file at `~/.config/zai/history.jsonl` contains malformed lines, possibly from a crash during a write or from very long responses exceeding the scanner buffer.

**Solution**: Back up the file, then remove or edit the corrupted lines. Each line is an independent JSON object, so removing a single bad line does not affect the others:

```bash
cp ~/.config/zai/history.jsonl ~/.config/zai/history.jsonl.bak
```

If the file is too damaged to salvage, delete it. Zai creates a new one on the next command.
