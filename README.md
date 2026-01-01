# ZAI CLI

A powerful command-line interface for Z.AI GLM models with chat, web search, image generation, and more.

## Quick Start

```bash
# Set API key
export ZAI_API_KEY="your_api_key"

# Start chatting
zai "Explain quantum computing"
```

## Installation

### Go Install
```bash
go install github.com/garyblankenship/zai@latest
```

### From Source
```bash
git clone https://github.com/garyblankenship/zai.git
cd zai
go build -o bin/zai . && ln -sf $(pwd)/bin/zai ~/go/bin/zai
```

### Homebrew (macOS)
```bash
brew install ffmpeg yt-dlp  # Optional: for audio transcription
```

## Configuration

### Environment Variables (Recommended)
```bash
export ZAI_API_KEY="your_api_key"
export ZAI_API_BASE_URL="https://api.z.ai/api/paas/v4"  # optional
export ZAI_API_MODEL="glm-4.7"  # optional
```

### Config File: `~/.config/zai/config.yaml`
```yaml
api:
  key: "your-api-key"
  base_url: "https://api.z.ai/api/paas/v4"  # default
  model: "glm-4.7"                           # default

# Web reader configuration
web_reader:
  enabled: true                # Enable web content fetching
  timeout: 20                 # Timeout in seconds
  cache_enabled: true         # Enable response caching
  return_format: markdown     # Default format (markdown/text)
  auto_detect: true          # Auto-detect URLs in chat
  max_content_length: 50000  # Max characters to include

# Web search configuration
web_search:
  enabled: true               # Enable web search
  default_count: 10          # Default number of results
  default_recency: "noLimit" # Time filter (oneDay, oneWeek, oneMonth, noLimit)
  timeout: 30               # Default timeout in seconds
  cache_enabled: true       # Enable search caching
  cache_dir: "~/.config/zai/search_cache"
  cache_ttl: 24h            # Cache duration
```

## Usage

### One-shot Mode
```bash
zai "Explain quantum computing"
zai -f main.go "Explain this code"
zai -f https://example.com "Summarize"   # -f supports URLs too!
zai -v "What is 2+2?"                    # verbose with debug info
zai --think "Analyze this problem"       # enable reasoning mode
zai --json "What is 2+2?"               # output as JSON
zai --search "What happened in AI today?" # web search augmented
```

### Piped Input
```bash
pbpaste | zai "explain this"
cat file.txt | zai "summarize"
echo "Hello" | zai "translate to French"
```

### Interactive Chat (REPL)
```bash
zai chat                                 # Start interactive session
zai chat -f main.go                     # With file context
zai chat --think                        # With reasoning mode
```

**REPL Commands:**
- `help` - Show available commands
- `history` - Show conversation history
- `context` - Show current context window
- `clear` - Clear conversation history
- `exit` or `Ctrl+D` - Exit chat

### Web Reader
```bash
# Fetch and display web content
zai reader https://example.com
zai reader https://example.com --format text
zai reader https://example.com --no-cache
zai reader https://example.com --timeout 30
zai reader https://example.com --json       # output as JSON

# Auto-detect URLs in prompts
zai "Summarize https://example.com/article"
```

### Web Search
```bash
# Basic search
zai search "golang best practices"
zai search "AI news" -c 5 -r oneWeek      # With count and recency
zai search "github.com" -d github.com     # Domain filter
zai search "query" -o json               # JSON output (format flag)
zai search "query" --json                # JSON output (global flag)

# Search in chat
zai chat
you> search "latest AI news" -c 3
you> /search "golang patterns" -r oneMonth
```

### Image Generation
```bash
zai image "a wizard"                        # AI-enhanced prompt + auto-save
zai image "sunset" --size 1024x1024         # Custom size
zai image "cat" --no-enhance                # Skip AI prompt enhancement
zai image "logo" -o logo.png                # Custom output path
zai image "art" --show                      # Open in viewer after generation
zai image "product" --copy                  # Copy URL to clipboard (macOS)
```

**Features:**
- **AI Prompt Enhancement**: Simple prompts are automatically enhanced with professional photography/lighting/style details (disable with `--no-enhance`)
- **Auto-save**: Images automatically saved to `zai-image-{timestamp}.png`
- **Professional prompts**: Enhancement adds cinematic lighting, camera specs, composition, mood

### Audio Transcription
```bash
zai audio recording.wav                      # Transcribe audio file
zai audio speech.mp3 --model glm-asr-2512   # Specify ASR model
zai audio interview.wav --prompt "Context"   # Add context
zai audio lecture.wav --hotwords "k8s,docker" # Domain vocabulary
zai audio --video https://youtu.be/abc123    # YouTube support
zai audio recording.wav --vad                # Remove silence (reduces cost)
zai audio recording.wav --resume             # Resume partial transcription
zai audio recording.wav --clear-cache        # Start fresh
zai audio recording.wav --json               # JSON output
```

**Features:**
- **YouTube Support**: Transcribe YouTube videos directly with `--video`
- **VAD (Voice Activity Detection)**: Remove silence to reduce API costs
- **Auto-chunking**: Large files (>25MB) automatically split into chunks
- **Resume Support**: Continue from where you left off on failures
- **Preprocessing**: Auto-converts to optimal 16kHz mono WAV
- **Hotwords**: Domain vocabulary for better accuracy (up to 100 terms)

**Supported formats**: .wav, .mp3, .mp4, .m4a, .flac, .aac, .ogg

### Vision Analysis
```bash
zai vision photo.jpg                         # Describe image
zai vision screenshot.png "What text?"        # Extract text
zai vision https://example.com/img.jpg       # Analyze URL
zai vision chart.png -p "Explain this chart" # Custom prompt
zai vision diagram.png -t 0.1                # Lower temperature
```

**Supported formats**: .jpg, .jpeg, .png, .gif, .webp

### Model Management
```bash
zai model list                    # List available models
zai model list --json             # List models as JSON
zai model current                 # Show current model
zai model set glm-4.7            # Switch model
```

### History Management
```bash
zai history                       # Show chat history
zai history -l 0                 # Show all history
zai history -n 10                # Show last 10 entries
zai history --search "quantum"    # Search history
zai history --json               # Output history as JSON
```

### Version
```bash
zai version                       # Show version, build time, commit
```

## Command Reference

### Global Flags
| Flag | Description |
|------|-------------|
| `-f, --file string` | Include file or URL contents in prompt |
| `-v, --verbose` | Show debug info and token usage |
| `--config string` | Custom config file path (default `~/.config/zai/config.yaml`) |
| `--think` | Enable thinking/reasoning mode |
| `--search` | Augment prompt with web search results |
| `--json` | Output in JSON format (structured data with metadata) |
| `-h, --help` | Help for zai |

### Commands
| Command | Description |
|---------|-------------|
| `chat` | Start interactive chat session (REPL) |
| `search` | Search the web using Z.AI search engine |
| `reader` | Fetch and display web content |
| `image` | Generate images using Z.AI's image generation API |
| `model` | Model management commands |
| `history` | Show chat history |
| `audio` | Transcribe audio files to text |
| `vision` | Analyze images with AI vision |
| `version` | Show version information |
| `completion` | Generate shell autocompletion script |

### Shell Autocompletion
```bash
# Bash
zai completion bash > ~/.local/share/bash-completion/completions/zai
source ~/.bashrc

# Zsh
zai completion zsh > "${fpath[1]}/_zai"
compinit

# Fish
zai completion fish > ~/.config/fish/completions/zai.fish
```

## JSON Output

The `--json` flag provides structured output for programmatic use and integration with other tools.

### Supported Commands
- **Root (one-shot)**: `zai "prompt" --json`
- **Search**: `zai search "query" --json`
- **Reader**: `zai reader https://example.com --json`
- **Model List**: `zai model list --json`
- **History**: `zai history --json`

### JSON Structure
Each JSON output includes:
- **Main data**: Results, content, or response
- **Metadata**: Timestamp, count, model information
- **Format**: Consistent 2-space indentation

### Example JSON Output
```json
{
  "query": "golang best practices",
  "count": 10,
  "duration": "1.234s",
  "results": [
    {
      "title": "Go Best Practices",
      "link": "https://example.com/go-practices",
      "content": "Comprehensive guide to Go development..."
    }
  ],
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Use Cases
- **CI/CD pipelines**: Parse results for automation
- **Scripts**: Process data programmatically
- **Integration**: Feed into other tools
- **Logging**: Structured output for analysis

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ZAI_API_KEY` | API key (required) | - |
| `ZAI_API_BASE_URL` | Override base URL | `https://api.z.ai/api/paas/v4` |
| `ZAI_API_MODEL` | Override model | `glm-4.7` |

## Architecture

The ZAI CLI follows SOLID principles with dependency injection:

```
zai/
├── cmd/                    # Command definitions
│   ├── root.go           # Main command, one-shot mode, stdin handling
│   ├── chat.go           # Interactive REPL with conversation context
│   ├── history.go        # History viewing command
│   ├── search.go         # Web search command
│   ├── web.go            # Web reader command
│   ├── image.go          # Image generation command
│   ├── audio.go          # Audio transcription command
│   ├── vision.go         # Vision analysis command
│   ├── model.go          # Model management command
│   └── version.go        # Version info command
├── internal/
│   ├── app/
│   │   ├── client.go     # HTTP client, API calls (DI, interfaces)
│   │   ├── types.go      # Request/response types
│   │   ├── history.go    # File-based history storage
│   │   ├── cache.go      # File-based search caching
│   │   └── utils.go      # URL detection, web content formatting
│   ├── config/
│   │   └── config.go     # Viper defaults and loading
│   └── version/
│       └── version.go    # Build-time version info
├── bin/                   # Built binaries
└── main.go
```

### Key Features
- **SOLID Architecture**: Dependency injection, interface-driven design
- **Context Management**: REPL maintains conversation context (last 20 messages)
- **Web Integration**: Auto-detect URLs and fetch content via Z.AI reader API
- **Search Capabilities**: Built-in web search with caching and filtering
- **Audio Transcription**: Speech-to-text with YouTube support and chunking
- **Vision Analysis**: Image understanding with GLM-4.6V
- **History Storage**: JSONL file at `~/.config/zai/history.jsonl`
- **Caching**: Intelligent caching for web content and search results
- **stdin Support**: Detects piped input automatically

## Requirements

- **Go**: 1.21+
- **OS**: Cross-platform (Linux, macOS, Windows)
- **API**: Z.AI API key (required)

**Optional** (for audio transcription):
- **ffmpeg**: Audio preprocessing and chunking
- **yt-dlp**: YouTube audio extraction

## Examples

### Code Analysis
```bash
zai -f main.go "Review this code for bugs"
zai -f script.py "Convert this to JavaScript"
```

### Research Assistant
```bash
zai search "Rust vs Go performance 2024"
zai reader https://arxiv.org/abs/2301.07041 "Summarize this paper"
```

### Audio Transcription
```bash
zai audio meeting.wav                           # Transcribe recording
zai audio podcast.mp3 --vad                     # Remove silence first
zai audio lecture.wav --hotwords "API,REST,GraphQL"  # Domain terms
zai audio --video https://youtu.be/dQw4w9WgXcQ  # Transcribe YouTube
zai audio long-interview.wav --resume           # Resume if interrupted
cat audio.wav | zai audio                       # From stdin
```

### Vision Analysis
```bash
zai vision photo.jpg                            # Describe image
zai vision screenshot.png "What text is here?"  # OCR/text extraction
zai vision chart.png "Explain the trends"       # Chart analysis
zai vision error.png "What's wrong?"            # Debug screenshots
zai vision https://example.com/diagram.jpg      # Analyze from URL
```

### Daily Use
```bash
zai "Write a professional email about the meeting delay"
zai "Explain machine learning like I'm 10"
zai "Create a git commit message for these changes"
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Follow the existing SOLID architecture
4. Add tests for new features
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
