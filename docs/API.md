# API Reference

This document provides comprehensive API documentation for the zai CLI client interfaces and caching layer.

## Table of Contents

- [Client Interface](#client-interface)
  - [Constructor](#constructor)
  - [Chat API](#chat-api)
  - [Vision API](#vision-api)
  - [Image Generation API](#image-generation-api)
  - [Model Listing API](#model-listing-api)
  - [Web Reader API](#web-reader-api)
  - [Web Search API](#web-search-api)
  - [Audio Transcription API](#audio-transcription-api)
  - [Video Generation API](#video-generation-api)
- [Cache Interface](#cache-interface)
  - [Constructor](#cache-constructor)
  - [Methods](#cache-methods)
  - [Statistics](#statistics)
- [Configuration Types](#configuration-types)
- [Error Handling](#error-handling)

---

## Client Interface

The `Client` is the main interface for interacting with the Z.AI API. It implements multiple client interfaces following the Interface Segregation Principle (ISP).

### Constructor

#### `NewClient`

```go
func NewClient(cfg ClientConfig, logger *slog.Logger, history HistoryStore, httpClient HTTPDoer) *Client
```

Creates a new client with injected dependencies.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `cfg` | `ClientConfig` | Client configuration (API key, base URL, model, timeout) |
| `logger` | `*slog.Logger` | Structured logger for debug/info output |
| `history` | `HistoryStore` | Optional history storage (can be nil) |
| `httpClient` | `HTTPDoer` | HTTP client (uses `http.Client` if nil) |

**Returns:** `*Client` - Configured client instance

**Example:**

```go
cfg := app.ClientConfig{
    APIKey:  "your-api-key",
    BaseURL: "https://api.z.ai/api/paas/v4",
    Model:   "glm-4.7",
    Timeout: 60 * time.Second,
}
logger := app.NewLogger(false)
client := app.NewClient(cfg, logger, nil, nil)
```

#### `NewClientWithDeps`

```go
func NewClientWithDeps(cfg ClientConfig, logger *slog.Logger, history HistoryStore, deps *ClientDeps) *Client
```

Creates a client with full dependency injection support for testing.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `cfg` | `ClientConfig` | Client configuration |
| `logger` | `*slog.Logger` | Structured logger |
| `history` | `HistoryStore` | Optional history storage |
| `deps` | `*ClientDeps` | Optional dependencies (HTTPClient, FileReader) |

**Returns:** `*Client` - Configured client instance

---

### Chat API

#### `Chat`

```go
func (c *Client) Chat(ctx context.Context, prompt string, opts ChatOptions) (string, error)
```

Sends a prompt to the chat completion API and returns the response.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | Context for cancellation/timeout |
| `prompt` | `string` | User prompt/question |
| `opts` | `ChatOptions` | Optional configuration overrides |

**Returns:** `(string, error)` - Response text or error

**Errors:**

- `*APIError` - API returned non-200 status (check `StatusCode` and `Body`)
- `fmt.Errorf` - Missing API key, file read failure, URL fetch failure
- `context.Canceled` - Context was canceled
- `context.DeadlineExceeded` - Context deadline exceeded

**Example:**

```go
response, err := client.Chat(ctx, "Explain Go interfaces", app.ChatOptions{
    Temperature: app.Float64Ptr(0.7),
    MaxTokens:   app.IntPtr(4096),
})
```

---

### Vision API

#### `Vision`

```go
func (c *Client) Vision(ctx context.Context, prompt string, imageBase64 string, opts VisionOptions) (string, error)
```

Analyzes an image using the vision model (glm-4.6v).

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | Context for cancellation/timeout |
| `prompt` | `string` | Analysis prompt (defaults to "What do you see...") |
| `imageBase64` | `string` | Base64-encoded image (data URI or raw) |
| `opts` | `VisionOptions` | Optional configuration overrides |

**Returns:** `(string, error)` - Analysis text or error

**Errors:**

- `*APIError` - API returned non-200 status
- `fmt.Errorf` - Missing API key, empty image data, no choices in response

**Example:**

```go
response, err := client.Vision(ctx, "What text is in this image?", base64Data, app.VisionOptions{
    Model: "glm-4.6v",
})
```

---

### Image Generation API

#### `GenerateImage`

```go
func (c *Client) GenerateImage(ctx context.Context, prompt string, opts ImageOptions) (*ImageResponse, error)
```

Generates an image using the GLM-Image model.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | Context for cancellation/timeout |
| `prompt` | `string` | Image generation prompt |
| `opts` | `ImageOptions` | Optional configuration (quality, size, model) |

**Returns:** `(*ImageResponse, error)` - Image generation response or error

**Errors:**

- `*APIError` - API returned non-200 status
- `fmt.Errorf` - Invalid image options (quality/size), no images in response

**Validation:**

- `Quality`: Must be "hd" or "standard"
- `Size`: Must be one of 1024x1024, 1024x768, 768x1024, 512x512

**Example:**

```go
response, err := client.GenerateImage(ctx, "a sunset over mountains", app.ImageOptions{
    Quality: "hd",
    Size:    "1024x1024",
})
```

---

### Model Listing API

#### `ListModels`

```go
func (c *Client) ListModels(ctx context.Context) ([]Model, error)
```

Fetches available models from the API.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | Context for cancellation/timeout |

**Returns:** `([]Model, error)` - List of available models or error

**Errors:**

- `*APIError` - API returned non-200 status
- `fmt.Errorf` - Unmarshal failure

**Example:**

```go
models, err := client.ListModels(ctx)
for _, model := range models {
    fmt.Println(model.ID, model.OwnedBy)
}
```

---

### Web Reader API

#### `FetchWebContent`

```go
func (c *Client) FetchWebContent(ctx context.Context, url string, opts *WebReaderOptions) (*WebReaderResponse, error)
```

Fetches and processes web content from a URL.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | Context for cancellation/timeout |
| `url` | `string` | URL to fetch (required) |
| `opts` | `*WebReaderOptions` | Optional configuration (can be nil) |

**Returns:** `(*WebReaderResponse, error)` - Web content response or error

**Errors:**

- `*APIError` - API returned non-200 status
- `fmt.Errorf` - Empty URL, unmarshal failure

**Defaults:**

- `ReturnFormat`: "markdown"
- `RetainImages`: true
- `Timeout`: 20 seconds

**Example:**

```go
response, err := client.FetchWebContent(ctx, "https://example.com", &app.WebReaderOptions{
    Timeout:      app.IntPtr(30),
    ReturnFormat: "text",
})
```

---

### Web Search API

#### `SearchWeb`

```go
func (c *Client) SearchWeb(ctx context.Context, query string, opts SearchOptions) (*WebSearchResponse, error)
```

Performs a web search using Z.AI's search API.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | Context for cancellation/timeout |
| `query` | `string` | Search query (required) |
| `opts` | `SearchOptions` | Search configuration |

**Returns:** `(*WebSearchResponse, error)` - Search results or error

**Errors:**

- `*APIError` - API returned non-200 status
- `fmt.Errorf` - Empty query, invalid count (1-50), invalid recency filter

**Validation:**

- `Count`: Must be between 1 and 50
- `RecencyFilter`: Must be one of "", "noLimit", "oneDay", "oneWeek", "oneMonth", "oneYear"

**Example:**

```go
response, err := client.SearchWeb(ctx, "golang best practices", app.SearchOptions{
    Count:         10,
    RecencyFilter: "oneWeek",
})
```

---

### Audio Transcription API

#### `TranscribeAudio`

```go
func (c *Client) TranscribeAudio(ctx context.Context, audioPath string, opts TranscriptionOptions) (*TranscriptionResponse, error)
```

Transcribes an audio file using the ASR model.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | Context for cancellation/timeout |
| `audioPath` | `string` | Path to audio file (required) |
| `opts` | `TranscriptionOptions` | Transcription configuration |

**Returns:** `(*TranscriptionResponse, error)` - Transcription result or error

**Errors:**

- `*APIError` - API returned non-200 status
- `fmt.Errorf` - Missing audio path, file read failure, file too large (>25MB)

**Constraints:**

- Max file size: 25MB
- Prompt max length: 8000 characters
- Hotwords max items: 100

**Example:**

```go
response, err := client.TranscribeAudio(ctx, "recording.wav", app.TranscriptionOptions{
    Model:    "glm-asr-2512",
    Hotwords: []string{"kubernetes", "docker"},
})
```

---

### Video Generation API

#### `GenerateVideo`

```go
func (c *Client) GenerateVideo(ctx context.Context, prompt string, opts VideoOptions) (*VideoGenerationResponse, error)
```

Creates a video using the CogVideoX-3 model (async operation).

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | Context for cancellation/timeout |
| `prompt` | `string` | Video description (max 512 chars) |
| `opts` | `VideoOptions` | Video configuration |

**Returns:** `(*VideoGenerationResponse, error)` - Task response with task ID or error

**Errors:**

- `*APIError` - API returned non-200 status
- `fmt.Errorf` - Invalid video options (quality, size, FPS, duration, too many images)

**Validation:**

- `Quality`: "quality" or "speed"
- `Size`: One of 1280x720, 720x1280, 1024x1024, 1920x1080, 1080x1920, 2048x1080, 3840x2160
- `FPS`: 30 or 60
- `Duration`: 5 or 10 seconds
- `ImageURLs`: Max 2 images (for first/last frame mode)

**Defaults:**

- `Quality`: "speed"
- `Size`: "1920x1080"
- `FPS`: 30
- `Duration`: 5 seconds

**Example:**

```go
response, err := client.GenerateVideo(ctx, "a cat playing with a ball", app.VideoOptions{
    Quality:  "quality",
    Size:     "1920x1080",
    Duration: 5,
})
// Use response.ID with RetrieveVideoResult to poll for completion
```

#### `RetrieveVideoResult`

```go
func (c *Client) RetrieveVideoResult(ctx context.Context, taskID string) (*VideoResultResponse, error)
```

Retrieves the result of an async video generation task.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | Context for cancellation/timeout |
| `taskID` | `string` | Task ID from GenerateVideo (required) |

**Returns:** `(*VideoResultResponse, error)` - Video result or error

**Errors:**

- `*APIError` - API returned non-200 status
- `fmt.Errorf` - Empty task ID, unmarshal failure

**Task Status:**

- `PROCESSING` - Video still being generated
- `SUCCESS` - Video ready (check `VideoResult` array)
- `FAIL` - Generation failed

**Example:**

```go
result, err := client.RetrieveVideoResult(ctx, taskID)
if result.TaskStatus == "SUCCESS" && len(result.VideoResult) > 0 {
    videoURL := result.VideoResult[0].URL
    fmt.Println("Video ready:", videoURL)
}
```

---

### Error Handling

#### `APIError`

```go
type APIError struct {
    StatusCode int
    Body       string
}
```

Use `errors.As` to extract `*APIError` from wrapped errors:

```go
var apiErr *app.APIError
if errors.As(err, &apiErr) {
    log.Printf("API error: %d - %s", apiErr.StatusCode, apiErr.Body)
}
```

**Common Status Codes:**

- `401` - Invalid or missing API key
- `429` - Rate limit exceeded (auto-retry with backoff)
- `500` - Internal server error (auto-retry with backoff)
- `502` - Bad gateway (auto-retry with backoff)
- `503` - Service unavailable (auto-retry with backoff)
- `504` - Gateway timeout (auto-retry with backoff)

---

### Interfaces

The client implements multiple interfaces following ISP:

```go
type ChatClient interface {
    Chat(ctx context.Context, prompt string, opts ChatOptions) (string, error)
}

type VisionClient interface {
    Vision(ctx context.Context, prompt string, imageBase64 string, opts VisionOptions) (string, error)
}

type ImageClient interface {
    GenerateImage(ctx context.Context, prompt string, opts ImageOptions) (*ImageResponse, error)
}

type ModelClient interface {
    ListModels(ctx context.Context) ([]Model, error)
}

type WebReaderClient interface {
    FetchWebContent(ctx context.Context, url string, opts *WebReaderOptions) (*WebReaderResponse, error)
}

type WebSearchClient interface {
    SearchWeb(ctx context.Context, query string, opts SearchOptions) (*WebSearchResponse, error)
}

type AudioClient interface {
    TranscribeAudio(ctx context.Context, audioPath string, opts TranscriptionOptions) (*TranscriptionResponse, error)
}

type VideoClient interface {
    GenerateVideo(ctx context.Context, prompt string, opts VideoOptions) (*VideoGenerationResponse, error)
    RetrieveVideoResult(ctx context.Context, taskID string) (*VideoResultResponse, error)
}
```

Use these interfaces for dependency injection in your code.

---

## Cache Interface

The `SearchCache` interface provides persistent file-based caching for web search results. It implements the Interface Segregation Principle (ISP) with a minimal interface focused on search result caching.

### Interface Definition

```go
type SearchCache interface {
    Get(query string, opts SearchOptions) ([]SearchResult, bool)
    Set(query string, opts SearchOptions, results []SearchResult, ttl time.Duration) error
    Clear() error
    Cleanup() error
}
```

---

### Cache Constructor

#### `NewFileSearchCache`

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

### Cache Methods

#### `Get`

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

#### `Set`

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

#### `Clear`

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

#### `Cleanup`

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

### Statistics

#### `Stats`

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

### Cache Key Generation

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

### Thread Safety

`FileSearchCache` is safe for concurrent use:

- `Get` uses read locks (multiple concurrent reads allowed)
- `Set`, `Clear`, `Cleanup`, `Stats` use exclusive locks
- Async cleanup triggered by expired/corrupted entries uses `TryLock` to prevent cleanup storms

---

### Cache Entry Format

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

### Expiration Behavior

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

## Configuration Types

### `ClientConfig`

```go
type ClientConfig struct {
    APIKey        string
    BaseURL       string
    CodingBaseURL string
    Model         string
    Timeout       time.Duration
    Verbose       bool
    RetryConfig   RetryConfig
}
```

### `ChatOptions`

```go
type ChatOptions struct {
    Model       string
    Temperature *float64
    MaxTokens   *int
    TopP        *float64
    Thinking    *bool
    WebEnabled  *bool
    WebTimeout  *int
    FilePath    string   // Legacy
    Context     []Message // Legacy
    Think       bool      // Legacy
}
```

### `VisionOptions`

```go
type VisionOptions struct {
    Model       string
    Temperature *float64
    MaxTokens   *int
    TopP        *float64
}
```

### `ImageOptions`

```go
type ImageOptions struct {
    Quality string
    Size    string
    UserID  string
    Model   string
}
```

### `WebReaderOptions`

```go
type WebReaderOptions struct {
    Timeout           *int
    NoCache           *bool
    ReturnFormat      string
    RetainImages      *bool
    NoGFM             *bool
    KeepImgDataURL    *bool
    WithImagesSummary *bool
    WithLinksSummary  *bool
}
```

### `SearchOptions`

```go
type SearchOptions struct {
    Count         int
    DomainFilter  string
    RecencyFilter string
    RequestID     string
    UserID        string
}
```

### `TranscriptionOptions`

```go
type TranscriptionOptions struct {
    Model     string
    Prompt    string
    Hotwords  []string
    Stream    bool
    UserID    string
    RequestID string
}
```

### `VideoOptions`

```go
type VideoOptions struct {
    Model     string
    Quality   string
    Size      string
    FPS       int
    Duration  int
    WithAudio bool
    ImageURLs []string
    UserID    string
    RequestID string
}
```

---

## Helper Functions

```go
func Float64Ptr(v float64) *float64
func IntPtr(v int) *int
func BoolPtr(v bool) *bool
```

Use these to create pointers to literal values for option structs.
