# Web Search Configuration

Web search behavior is configured via the `web_search` section in `~/.config/zai/config.yaml`.

## Configuration Keys

### `enabled`

**Type:** `boolean`
**Default:** `true`

Enable or disable web search functionality globally.

```yaml
web_search:
  enabled: true
```

When disabled, all search commands (`zai search`, `/search` in chat) will fail with an error message.

### `default_count`

**Type:** `integer`
**Default:** `10`
**Range:** `1-50`

Default number of search results to return per search query.

```yaml
web_search:
  default_count: 10
```

This value is used when the `-c, --count` flag is not specified in a search command. Higher values provide more comprehensive results but increase API usage and processing time.

### `default_recency`

**Type:** `string`
**Default:** `"noLimit"`
**Valid Values:** `"oneDay"`, `"oneWeek"`, `"oneMonth"`, `"oneYear"`, `"noLimit"`

Default time filter for search results. Limits results to content published within the specified timeframe.

```yaml
web_search:
  default_recency: "noLimit"
```

**Recency Filter Values:**

| Value | Description |
|-------|-------------|
| `oneDay` | Results from the last 24 hours |
| `oneWeek` | Results from the last 7 days |
| `oneMonth` | Results from the last 30 days |
| `oneYear` | Results from the last 365 days |
| `noLimit` | No time restriction |

This value is used when the `-r, --recency` flag is not specified in a search command.

### `timeout`

**Type:** `integer`
**Default:** `30`
**Unit:** seconds

Maximum duration to wait for search API responses.

```yaml
web_search:
  timeout: 30
```

Increase this value if you experience timeouts on slower network connections.

## Caching Configuration

### `cache_enabled`

**Type:** `boolean`
**Default:** `true`

Enable or disable search result caching.

```yaml
web_search:
  cache_enabled: true
```

When enabled, search results are cached locally to avoid redundant API calls for identical queries. The cache key includes the query string, count, recency filter, and domain filter.

### `cache_dir`

**Type:** `string`
**Default:** `~/.config/zai/search_cache`

Directory where cached search results are stored.

```yaml
web_search:
  cache_dir: "~/.config/zai/search_cache"
```

Cache files are stored as JSON files named by the SHA256 hash of the search parameters. The path is expanded for `~` (home directory) and environment variables.

### `cache_ttl`

**Type:** `duration`
**Default:** `"24h"`
**Format:** Go duration string (e.g., `"1h"`, `"30m"`, `"24h"`, `"7d"`)

Maximum time to retain cached search results before re-fetching.

```yaml
web_search:
  cache_ttl: 24h
```

**Duration Format Examples:**

| Value | Description |
|-------|-------------|
| `"1h"` | 1 hour |
| `"30m"` | 30 minutes |
| `"24h"` | 24 hours (1 day) |
| `"168h"` | 168 hours (1 week) |
| `"7d"` | 7 days (supported for readability) |

Cached results older than the TTL are automatically ignored and re-fetched from the API.

## Complete Example

```yaml
web_search:
  enabled: true
  default_count: 10
  default_recency: "noLimit"
  timeout: 30
  cache_enabled: true
  cache_dir: "~/.config/zai/search_cache"
  cache_ttl: 24h
```

## Performance Considerations

- **Higher `default_count`** increases API usage and processing time
- **Shorter `cache_ttl`** increases API usage (more frequent re-fetches)
- **`cache_enabled: false`** significantly increases API costs for repeated queries
- **`timeout`** too low may cause failures on slow connections; too high may cause long waits on API failures

## Related Documentation

- [Search Command](../commands/search.md) - Command-line usage and flags
- [Chat Command](../commands/chat.md) - Interactive chat with search integration
- [Web Reader Configuration](./web-reader.md) - Web content fetching configuration
