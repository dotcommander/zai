# Search Command

Search the web using the Z.AI search engine with results optimized for LLM consumption.

## Usage

```bash
zai search [query] [flags]
```

The search query can be provided as a command-line argument or piped via stdin.

## Options

### `-c, --count <number>`

Number of search results to return (1-50). If not specified, uses the default from configuration (`web_search.default_count`).

```bash
zai search "golang patterns" -c 5
```

### `-r, --recency <filter>`

Time filter for search results. Limits results to content published within the specified timeframe.

| Value | Description |
|-------|-------------|
| `oneDay` | Results from the last 24 hours |
| `oneWeek` | Results from the last 7 days |
| `oneMonth` | Results from the last 30 days |
| `oneYear` | Results from the last 365 days |
| `noLimit` | No time restriction (default) |

```bash
zai search "AI news" -r oneWeek
zai search "golang release notes" -r oneMonth
```

### `-d, --domain <domain>`

Limit search results to a specific domain. Useful for searching within a particular website.

```bash
zai search "kubernetes tutorials" -d kubernetes.io
zai search "github.com" -d github.com
```

### `-o, --format <format>`

Output format for search results.

| Format | Description |
|--------|-------------|
| `table` | Compact table with title, domain, and URL (default) |
| `detailed` | Full results with content preview, media, and publish dates |
| `json` | Machine-readable JSON output |

```bash
zai search "rust vs go" -o detailed
zai search "machine learning" -o json
```

## Examples

### Basic Search

```bash
zai search "golang best practices"
```

### Search from Stdin

```bash
echo "rust ownership patterns" | zai search
```

### Limit Results and Timeframe

```bash
zai search "latest AI developments" -c 10 -r oneWeek
```

### Domain-Specific Search

```bash
# Search only GitHub
zai search "fiber web framework" -d github.com

# Search only a specific site
zai search "kubernetes deployment" -d kubernetes.io
```

### JSON Output for Scripting

```bash
zai search "docker compose" -o json
```

**Output:**

```json
{
  "query": "docker compose",
  "duration": "1.234s",
  "count": 10,
  "timestamp": "2026-01-18T22:00:00Z",
  "results": [
    {
      "title": "Docker Compose overview",
      "link": "https://docs.docker.com/compose/",
      "content": "Docker Compose is a tool for defining and running multi-container Docker applications...",
      "media": "docs.docker.com",
      "publishDate": "2024-12-15"
    }
  ]
}
```

### Verbose Output

```bash
zai search "golang generics" -v
```

Shows search duration, result count, and content previews.

## Configuration

Search behavior can be configured via `~/.config/zai/config.yaml`:

```yaml
web_search:
  enabled: true           # Enable/disable web search
  default_count: 10       # Default number of results
  default_recency: "noLimit"  # Default time filter
  timeout: 30            # Request timeout in seconds
  cache_enabled: true    # Enable response caching
  cache_dir: "~/.config/zai/search_cache"
  cache_ttl: 24h         # Cache duration
```

## Search in Chat Mode

Web search is also available within interactive chat sessions:

```bash
zai chat
you> search "latest golang release" -c 5 -r oneWeek
you> /search "rust async patterns" -r oneMonth
```

## See Also

- [Web Reader Command](./reader.md) - Fetch and parse web content
- [Chat Command](./chat.md) - Interactive AI chat with search integration
