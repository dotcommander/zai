# CLAUDE.md

CLI tool for chatting with Z.AI GLM models.

## Build & Run

```bash
go build -o bin/zai .           # Build
./bin/zai "prompt"              # One-shot
./bin/zai chat                  # Interactive REPL (styled with lipgloss)
echo "text" | ./bin/zai         # Stdin pipe
./bin/zai -f file.go "explain"  # With file context
./bin/zai -f https://url "sum"  # -f supports URLs too
./bin/zai web <url>             # Fetch web content
./bin/zai search "query"        # Web search
./bin/zai --search "news"       # Search-augmented generation
./bin/zai image "wizard"        # AI-enhanced image generation
```

## Install Globally

```bash
ln -sf /Users/vampire/go/src/zai/bin/zai ~/go/bin/zai
```

## Configuration

Config file: `~/.config/zai/config.yaml`

```yaml
api:
  key: "your-api-key"
  base_url: "https://api.z.ai/api/paas/v4"  # default
  model: "glm-4.7"                           # default
web_reader:
  enabled: true           # Enable web content fetching
  timeout: 20            # Default timeout in seconds
  cache_enabled: true    # Enable response caching
  return_format: markdown # Default format
  auto_detect: true      # Auto-detect URLs in chat
  max_content_length: 50000 # Max characters to include
web_search:
  enabled: true           # Enable web search
  default_count: 10       # Default number of results
  default_recency: "noLimit" # Time filter
  timeout: 30            # Default timeout in seconds
  cache_enabled: true    # Enable search caching
  cache_dir: "~/.config/zai/search_cache"
  cache_ttl: 24h         # Cache duration
```

Environment: `ZAI_API_KEY` overrides config file.

## Commands

### Chat
```bash
zai chat                    # Interactive REPL
zai chat -f file.go         # Chat with file context
zai chat --think            # Enable reasoning mode
```

### Search
```bash
zai search "golang best practices"              # Basic search
zai search "AI news" -c 5 -r oneWeek              # With options
zai search "github.com" -d github.com             # Domain filter
zai search "query" -o json                        # JSON output
echo "golang" | zai search                         # From stdin
```

**Search in Chat:**
```bash
zai chat
you> search "latest AI news" -c 3
you> /search "golang patterns" -r oneMonth
```

### Web Reader
```bash
zai web https://example.com                 # Fetch and display
zai web https://example.com --format text   # Plain text format
zai web https://example.com --no-cache      # Disable caching
zai web https://example.com --timeout 30    # Custom timeout
```

### Auto Web Content in Chat
URLs in prompts are automatically fetched and included:
```bash
zai "Summarize https://example.com"  # Fetches and includes content
```

### Search-Augmented Generation
```bash
zai --search "What's happening in AI today?"  # Searches web, adds context
zai chat --search                              # Enable for entire chat session
```

### Image Generation
```bash
zai image "a wizard"              # AI-enhanced prompt + auto-download
zai image "sunset" -s 1024x768    # Custom size
zai image "cat" --no-enhance      # Skip prompt enhancement
zai image "logo" -o output.png    # Custom filename
```

**Features:**
- **AI prompt enhancement** (default on): Transforms simple prompts into professional image generation prompts with lighting, composition, style
- **Auto-download**: Images saved to `zai-image-{timestamp}.png`
- Combines original + enhanced prompt for best results

## Architecture

```
cmd/
  root.go     # Main command, stdin handling, one-shot mode
  chat.go     # Interactive REPL with conversation context and search
  history.go  # History viewing command
  search.go   # Web search command
  web.go      # Web reader command
  image.go    # Image generation command
  model.go    # Model management command
internal/
  app/
    cache.go    # File-based search caching
    client.go   # HTTP client, API calls (DI, interfaces)
    types.go    # Request/response types
    history.go  # File-based history storage
    utils.go    # URL detection, web content and search formatting
  config/
    config.go   # Viper defaults and loading (includes web_search config)
```

**Design**: SOLID-compliant with dependency injection. Client takes `Logger` and `HistoryStore` interfaces.

## Key Patterns

- **Stdin detection**: `(stat.Mode() & os.ModeCharDevice) == 0`
- **Stdin + prompt**: Combines as `prompt + <stdin>data</stdin>`
- **History**: JSONL file at `~/.config/zai/history.jsonl` (includes search history)
- **Context**: REPL keeps last 20 messages (10 exchanges)
- **Web Content**: Auto-detects URLs, fetches via `/paas/v4/reader` API
- **Web Content Format**: Wrapped in `<web_content>` XML tags with metadata
- **Web Search**: Queries `/paas/v4/web_search` API with caching and filtering
- **Search Cache**: File-based with SHA256 keys and TTL expiration
- **Search Augmentation**: `--search` flag prepends web results as `<web_search_results>` context
- **File flag URLs**: `-f` detects http/https and routes to web reader
- **Image Enhancement**: LLM transforms simple prompts using professional image engineering framework
- **Chat TUI**: Styled with Charmbracelet lipgloss (colors, spinner, styled output)
