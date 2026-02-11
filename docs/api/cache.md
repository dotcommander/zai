# Cache Interface API Reference

The `SearchCache` interface provides persistent file-based caching for web search results. It implements the Interface Segregation Principle (ISP) with a minimal interface focused on search result caching.

## Interface Definition

```go
type SearchCache interface {
    Get(query string, opts SearchOptions) ([]SearchResult, bool)
    Set(query string, opts SearchOptions, results []SearchResult, ttl time.Duration) error
    Clear() error
    Cleanup() error
}
```

---

## Constructor

### `NewFileSearchCache`

```go
func NewFileSearchCache(dir string) *FileSearchCache
```

Creates a new file-based search cache instance.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `dir` | `string` | Cache directory path (created if missing) |

**Returns:** `*FileSearchCache` - Cache instance

**Example:**

```go
cache := app.NewFileSearchCache("/path/to/cache")
```

---

## Methods

### `Get`

```go
func (fsc *FileSearchCache) Get(query string, opts SearchOptions) ([]SearchResult, bool)
```

Retrieves cached search results for a query and options.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `query` | `string` | Search query string |
| `opts` | `SearchOptions` | Search options (affects cache key) |

**Returns:** `([]SearchResult, bool)` - Cached results and true if found, nil and false if not found/expired

**Behavior:**

- Returns `false` if cache entry does not exist
- Returns `false` if cache entry has expired (based on TTL)
- Returns `false` if cache file is corrupted (triggers async cleanup)
- Cache key includes query + domain filter + recency filter + count

**Example:**

```go
results, found := cache.Get("golang patterns", app.SearchOptions{
    Count:         10,
    RecencyFilter: "oneWeek",
})
if found {
    fmt.Println("Cache hit:", len(results), "results")
}
```

---

### `Set`

```go
func (fsc *FileSearchCache) Set(query string, opts SearchOptions, results []SearchResult, ttl time.Duration) error
```

Stores search results in cache with a time-to-live (TTL).

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `query` | `string` | Search query string |
| `opts` | `SearchOptions` | Search options (affects cache key) |
| `results` | `[]SearchResult` | Search results to cache |
| `ttl` | `time.Duration` | Time-to-live before expiration |

**Returns:** `error` - nil on success, error if directory creation or file write fails

**Behavior:**

- Creates cache directory if it doesn't exist
- Overwrites existing cache entry with same key
- Each entry stores query, results, cached timestamp, expiration timestamp, and hash

**Example:**

```go
results := []app.SearchResult{
    {Title: "Go Patterns", Link: "https://example.com/go"},
}
err := cache.Set("golang patterns", app.SearchOptions{Count: 10}, results, 24*time.Hour)
```

---

### `Clear`

```go
func (fsc *FileSearchCache) Clear() error
```

Removes all cached entries from the cache directory.

**Returns:** `error` - nil on success, error if directory read fails

**Behavior:**

- Only removes files ending in `.json`
- Returns nil if cache directory doesn't exist (no-op)
- Thread-safe (uses exclusive lock)

**Example:**

```go
err := cache.Clear()
if err != nil {
    log.Printf("Failed to clear cache: %v", err)
}
```

---

### `Cleanup`

```go
func (fsc *FileSearchCache) Cleanup() error
```

Removes expired and corrupted cache entries.

**Returns:** `error` - nil on success, error if directory read fails

**Behavior:**

- Removes entries where `ExpiresAt < time.Now()`
- Removes entries that fail to unmarshal (corrupted)
- Returns nil if cache directory doesn't exist (no-op)
- Thread-safe (uses exclusive lock)

**Example:**

```go
err := cache.Cleanup()
if err != nil {
    log.Printf("Cleanup failed: %v", err)
}
```

---

## Cache Key Generation

Cache keys are generated using SHA256 hash of the query and relevant search options:

```go
func generateCacheKey(query string, opts SearchOptions) string
```

**Key Components:**

| Component | Condition | Format |
|-----------|-----------|--------|
| Query | Always | Raw query string |
| Domain filter | If non-empty | `domain:{value}` |
| Recency filter | If not "noLimit" | `recency:{value}` |
| Count | If > 0 | `count:{value}` |

**Examples:**

```
"hello world" → sha256("hello world")
"golang" with domain filter → sha256("golangdomain:github.com")
"AI news" with recency filter → sha256("AI newsrecency:oneWeekcount:10")
```

**Notes:**

- Same query + different options = different cache keys
- Omitted options don't affect the cache key
- Hex-encoded SHA256 produces 64-character filenames

---

## Statistics

### `Stats`

```go
func (fsc *FileSearchCache) Stats() (*CacheStats, error)
```

Returns statistics about the cache state.

**Returns:** `(*CacheStats, error)` - Cache statistics or error

**Statistics Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `TotalEntries` | `int` | Total number of cache entries |
| `ExpiredEntries` | `int` | Number of expired entries |
| `SizeBytes` | `int64` | Total size of cache files in bytes |
| `CacheDir` | `string` | Cache directory path |

**Example:**

```go
stats, err := cache.Stats()
if err == nil {
    fmt.Printf("Cache: %d entries (%d expired), %.2f MB\n",
        stats.TotalEntries,
        stats.ExpiredEntries,
        float64(stats.SizeBytes)/(1024*1024),
    )
}
```

---

## Thread Safety

`FileSearchCache` is safe for concurrent use:

- `Get` uses read locks (multiple concurrent reads allowed)
- `Set`, `Clear`, `Cleanup`, `Stats` use exclusive locks
- Async cleanup triggered by expired/corrupted entries uses `TryLock` to prevent cleanup storms

---

## Cache Entry Format

Each cache file is stored as JSON:

```go
type SearchCacheEntry struct {
    Query     string         `json:"query"`
    Results   []SearchResult `json:"results"`
    CachedAt  time.Time      `json:"cached_at"`
    ExpiresAt time.Time      `json:"expires_at"`
    Hash      string         `json:"hash"`
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `Query` | `string` | Original search query |
| `Results` | `[]SearchResult` | Cached search results |
| `CachedAt` | `time.Time` | When the entry was cached |
| `ExpiresAt` | `time.Time` | When the entry expires |
| `Hash` | `string` | SHA256 hash used as cache key |

**File Naming:** `{sha256hash}.json`

**Example File:**

```json
{
  "query": "golang patterns",
  "results": [
    {
      "title": "Go Patterns",
      "link": "https://example.com/go",
      "content": "...",
      "media": "",
      "icon": "",
      "refer": "",
      "publish_date": ""
    }
  ],
  "cached_at": "2026-01-19T00:00:00Z",
  "expires_at": "2026-01-20T00:00:00Z",
  "hash": "a1b2c3d4..."
}
```

---

## Expiration Behavior

Cache entries expire based on the TTL provided to `Set`:

- Get returns `false` for expired entries (lazy expiration)
- Expired entries are removed by `Cleanup()` (manual) or `tryCleanup()` (async)
- Expiration check compares `time.Now()` > `entry.ExpiresAt`

**Common TTL Values:**

| Duration | Use Case |
|----------|----------|
| `1 * time.Hour` | Fast-changing topics (news, stocks) |
| `24 * time.Hour` | Default for general searches |
| `7 * 24 * time.Hour` | Stable reference content |

---

## Configuration

Search caching is configured via `config.yaml`:

```yaml
web_search:
  cache_enabled: true    # Enable caching
  cache_dir: "~/.config/zai/search_cache"
  cache_ttl: 24h         # Default TTL
```

**Environment Override:** `ZAI_SEARCH_CACHE_DIR` can override the cache directory.

---

## See Also

- [Client Interface](client.md) - Client API documentation
- [Types Reference](types.md) - Request/response type definitions
- [History Interface](history.md) - Chat history storage
