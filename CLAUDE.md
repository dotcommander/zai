# CLAUDE.md

CLI tool for chatting with Z.AI GLM models.

## Build & Install

```bash
go build -o bin/zai .                                     # Build
go build -o bin/zai . && ln -sf $(pwd)/bin/zai ~/go/bin/zai  # Install globally
```

## Usage

```bash
zai "prompt"                      # One-shot (streaming)
zai chat                          # Interactive REPL (streaming)
echo "text" | zai                 # Stdin pipe
zai -f file.go "explain"          # With file context
zai --search "query"              # Search-augmented generation
zai "prompt" --json               # Non-streaming JSON output
zai "prompt" --system "Be brief"  # Custom system prompt
zai "prompt" -C                   # Use coding API endpoint
zai version                       # Show version info
```

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--file` | `-f` | Include file or URL in prompt |
| `--search` | | Augment with web search results |
| `--think` | | Enable reasoning mode |
| `--coding` | `-C` | Use coding API endpoint |
| `--json` | | Output as JSON |
| `--system` | | Custom system prompt (literal or file path) |
| `--config` | | Config file path (default `~/.config/zai/config.yaml`) |
| `--verbose` | `-v` | Show debug info |

### Pipe Chaining

When stdout is not a TTY, zai outputs raw text — no styling, no wrappers. This enables chaining:

```bash
cat code.go | zai "find bugs" | zai "suggest fixes"
git diff | zai "summarize changes" > summary.txt
```

## Configuration

Config file: `~/.config/zai/config.yaml`

```yaml
api:
  key: "your-api-key"
  base_url: "https://api.z.ai/api/paas/v4"
  coding_base_url: "https://api.z.ai/api/coding/paas/v4"
  embedding_base_url: ""  # Set to "https://open.bigmodel.cn/api/paas/v4" for dedicated embedding-3 model
  coding_plan: false
  model: "glm-4.7"
  image_model: "cogview-4-250304"
  video_model: "cogvideox-3"
  vision_model: "glm-4.6v"
  audio_model: "glm-asr-2512"
  tts_model: "glm-tts"
  embedding_model: "glm-4.7"
  rate_limit:
    requests_per_second: 10
    burst: 5
  retry:
    max_attempts: 3
    initial_backoff: 1s
    max_backoff: 30s
  circuit_breaker:
    enabled: true
    failure_threshold: 5
    success_threshold: 2
    timeout: 60s

web_reader:
  enabled: true
  timeout: 20
  cache_enabled: true
  return_format: markdown
  auto_detect: true
  max_content_length: 50000

web_search:
  enabled: true
  default_count: 10
  default_recency: "noLimit"
  language: "en"
  timeout: 30
  cache_enabled: true
  cache_dir: "~/.config/zai/search_cache"
  cache_ttl: 24h
  search_engine: "search_std"
  content_size: "medium"

tts:
  voice: "tongtong"
  response_format: "wav"
```

Environment: `ZAI_API_KEY` overrides config file.

## Commands

### Chat
```bash
zai chat                          # Interactive REPL with lipgloss styling
zai chat -f file.go               # With file context
zai chat --think                  # Enable reasoning mode
zai chat --model glm-4.7          # Override model for session
```

### Search
```bash
zai search "query"                                        # Web search
zai search "query" -c 5 -r oneWeek -d github.com          # With filters
zai search "query" --engine search_pro --content-size high --sources  # Enhanced
zai search "query" -o json                                 # JSON output
zai search "query" -l zh --require-search                  # Language + force search
zai chat --search                                          # Search-augmented chat
```

Engines: `search_std`, `search_pro`, `search_pro_sogou`, `search_pro_quark`. Output formats: `table`, `detailed`, `json`.

In chat: `search "query"` or `/search "query"`

### Reader
```bash
zai reader https://example.com                      # Fetch web content
zai reader https://example.com --format text         # Plain text
zai reader https://example.com --with-links-summary  # Include link summary
zai reader https://example.com --no-cache            # Skip cache
zai "Summarize https://example.com"                  # Auto-fetch URLs in prompts
```

### Image
```bash
zai image "wizard"                          # AI-enhanced prompt + auto-download
zai image "sunset" -s 1024x768 --no-enhance -o output.png
zai image "landscape" -s 768x1344           # CogView-4 tall size
zai image "logo" --copy -S                  # Copy to clipboard + open viewer
zai image "art" -m cogview-4-250304 -q standard
zai image list                              # List available image models
```

Auto-downloads to `zai-image-{timestamp}.png`. Uses CogView-4 by default. AI enhancement transforms prompts with lighting/composition/style. Quality: `hd` (default), `standard`.

### Vision
```bash
zai vision -f photo.jpg "What text?"           # Analyze image (local or URL)
zai vision -f chart.png -p "Explain trends"    # Custom prompt
zai vision -f img.jpg -m glm-4v-flash -t 0.5  # Override model + temperature
```

### Audio
```bash
zai audio -f recording.wav                              # Transcribe audio
zai audio -f speech.mp3 --hotwords "kubernetes,docker"  # Domain vocabulary
zai audio --video https://youtu.be/abc123 --vad         # YouTube with VAD
zai audio -f lecture.wav --stream                        # Streaming transcription
zai audio -f long.wav --resume                           # Resume partial transcription
zai audio -f speech.wav -l en --prompt "Prior context"   # Language + context
cat audio.wav | zai audio                                # From stdin
```

Supports: .wav, .mp3, .mp4, .m4a, .flac, .aac, .ogg (max 25MB). Auto-splits long files into 30s chunks. `--preprocess` (on by default) converts to optimal 16kHz mono WAV.

### Video
```bash
zai video "A cat playing"                                  # Text-to-video
zai video -f img.jpg "Animate this"                        # Image-to-video
zai video -f first.jpg -f last.jpg "transition"            # First/last frame
zai video "prompt" -q quality -s 1920x1080 -S              # High quality + open
zai video "prompt" --fps 60 --duration 10 --with-audio     # 60fps, 10s, AI audio
```

Auto-downloads to `zai-video-{timestamp}.mp4`. Async polling (1-3 min). Defaults: `speed` quality, `1920x1080`, 30fps, 5s. Sizes: 1280x720, 1024x1024, 1920x1080, 3840x2160.

### TTS
```bash
zai tts "Hello, world!"                    # Text-to-speech (default voice: tongtong)
zai tts "Welcome" --voice xiaochen         # Different voice
echo "text" | zai tts --output out.wav     # From stdin
zai tts "Fast speech" --speed 2            # Speed multiplier
zai tts "Hello" --format pcm              # PCM output format
```

Voices: tongtong (default), xiaochen, chuichui, jam, kazi, douji, luodo. Formats: `wav` (default), `pcm`. Saves to `zai-tts-{timestamp}.wav`.

### Embed
```bash
zai embed "Hello, world!"                  # Single embedding (JSON output)
zai embed "sentence one" "sentence two"    # Multiple embeddings
echo "text" | zai embed                    # From stdin
zai embed "text" -m glm-5                  # Override model
```

Always outputs JSON (embedding vectors are numeric arrays). Uses GLM-4.7 by default. For dedicated `embedding-3` model, set `embedding_base_url` to `https://open.bigmodel.cn/api/paas/v4` (requires China API key).

### History
```bash
zai history                 # Show last 10 entries
zai history -l 50           # Show last 50 entries
zai history --json          # JSON output
```

### Model
```bash
zai model list              # List available models from API
```

### Version
```bash
zai version                 # Show version, build time, commit hash
```

## Architecture

```
cmd/
  root.go     # Main command, stdin handling, one-shot mode
  chat.go     # Interactive REPL with conversation context
  history.go  # History viewing
  search.go   # Web search (enhanced: engine, content-size, sources)
  web.go      # Web reader (reader subcommand)
  image.go    # Image generation (CogView-4)
  vision.go   # Vision analysis
  audio.go    # Audio transcription
  tts.go      # Text-to-speech (GLM-TTS)
  embed.go    # Text embeddings (Embedding-3)
  video.go    # Video generation
  model.go    # Model management
  version.go  # Version info
  theme.go    # Shared lipgloss styling/colors
internal/
  app/
    cache.go    # File-based search caching
    client.go   # HTTP client, API calls (DI, interfaces)
    types.go    # Request/response types
    history.go  # File-based history storage
    utils.go    # URL detection, web content/search formatting
  config/
    config.go   # Viper defaults and loading
```

**Design**: SOLID-compliant with dependency injection. Client uses `Logger` and `HistoryStore` interfaces.

## Key Patterns

- **Streaming**: SSE streaming is the default for chat — tokens render as they arrive via `StreamChat()`/`StreamReader`. JSON mode (`--json`) falls back to non-streaming `Chat()`
- **Pipe detection**: `isTerminal(os.Stdout)` checks `os.ModeCharDevice`. Non-TTY stdout outputs raw streamed text for pipe chaining
- **Stdin detection**: `(stat.Mode() & os.ModeCharDevice) == 0`
- **Stdin + prompt**: Combines as `prompt + <stdin>data</stdin>`
- **History**: JSONL at `~/.config/zai/history.jsonl`
- **Context**: REPL keeps last 20 messages (10 exchanges)
- **Web Content**: Auto-detects URLs, fetches via `/paas/v4/reader` API, wraps in `<web_content>` XML tags
- **Web Search**: `/paas/v4/web_search` API with SHA256-keyed file cache. Supports `search_std`, `search_pro`, `search_pro_sogou`, `search_pro_quark` engines. Configurable content size, source details, and result ordering
- **Search Augmentation**: `--search` flag prepends `<web_search_results>` context
- **File flag URLs**: `-f` detects http/https and routes to web reader
- **Image Enhancement**: LLM transforms prompts using professional image engineering framework
- **Retry/Circuit Breaker**: Exponential backoff with jitter; circuit breaker per endpoint (Closed → Open → Half-Open). Streaming requests do NOT retry (not idempotent mid-stream)

## Operational Context

- In `cmd/chat.go`, `exit`/`/exit` must return from the REPL loop (not just print Goodbye), otherwise chat continues running.
- Chat input is read in a goroutine and consumed via channels so `SIGINT`/Ctrl-C can terminate even while waiting on stdin.
- `zai chat` supports command-local model override via `--model`/`-m`, wired to `app.ChatOptions.Model` in `initializeChatOptions`.
- `system` can be a literal prompt or file path; in chat init, path-like values that do not resolve fall back to the built-in concise default prompt.
