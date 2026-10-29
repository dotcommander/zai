All configuration lives in `~/.config/zai/config.yaml`. Every setting has a sensible default -- you only need to set your API key to get started. Individual flags override config values per-invocation when you need something different for a single command.

- [API Settings](#api-settings)
- [Rate Limiting](#rate-limiting)
- [Retry Behavior](#retry-behavior)
- [Circuit Breaker](#circuit-breaker)
- [Web Search Settings](#web-search-settings)
- [Web Reader Settings](#web-reader-settings)
- [TTS Settings](#tts-settings)
- [Media Downloads](#media-downloads)
- [Account Quota](#account-quota)
- [Environment Variables](#environment-variables)

<a name="api-settings"></a>
## API Settings

```yaml
api:
  key: "your-api-key"
  base_url: "https://api.z.ai/api/paas/v4"
  coding_base_url: "https://api.z.ai/api/coding/paas/v4"
  coding_plan: false
  model: "glm-4.7"
  image_model: "cogview-4-250304"
  video_model: "cogvideox-3"
  vision_model: "glm-4.6v"
  audio_model: "glm-asr-2512"
  tts_model: "glm-tts"
  embedding_model: "embedding-3"
```

Each model can be overridden per-command with the `-m` flag. The `base_url` is the primary API endpoint for chat, search, reader, and most operations. The `coding_base_url` is used exclusively when the `-C` flag is active. See the [chat documentation](/docs/chat#coding-api) for details on the coding endpoint.

Ordinary API requests use a total request timeout. Chat and audio streaming requests have no total timeout; their default transport retains a 30-second dial timeout, a 10-second TLS handshake timeout, and a 60-second response-header timeout. Streams still stop when their caller context is canceled or its deadline expires.

<a name="rate-limiting"></a>
## Rate Limiting

```yaml
api:
  rate_limit:
    requests_per_second: 10
    burst: 5
```

Rate limiting applies to all API calls. The `burst` value allows short spikes above the sustained rate. When the request rate is positive, `burst: 0` is treated as `burst: 1`. Streaming and ordinary API requests share the rate limit. Set `requests_per_second: 0` to disable rate limiting entirely.

For most users, the defaults are fine. If you are running batch scripts, you may want to increase the limits temporarily.

<a name="retry-behavior"></a>
## Retry Behavior

```yaml
api:
  retry:
    max_attempts: 3
    initial_backoff: 1s
    max_backoff: 30s
```

Failed requests are retried with exponential backoff and jitter. The backoff starts at `initial_backoff` and doubles on each attempt, capping at `max_backoff`. Jitter prevents thundering herd when multiple requests fail simultaneously.

> [!WARNING]
> Streaming requests are not retried because they are not idempotent mid-stream. If a streaming request fails partway through, you receive a partial response and an error.

<a name="circuit-breaker"></a>
## Circuit Breaker

```yaml
api:
  circuit_breaker:
    enabled: true
    failure_threshold: 5
    success_threshold: 2
    timeout: 60s
```

The circuit breaker prevents cascading failures when the API is unhealthy. It operates in three states:

**Closed** (normal operation): requests flow through normally. Each failure increments a counter. When the counter reaches `failure_threshold`, the circuit opens.

**Open** (failing fast): all requests fail immediately without hitting the API. After `timeout` elapses, the circuit transitions to half-open.

**Half-Open** (testing recovery): one probe request is allowed through at a time. If `success_threshold` consecutive requests succeed, the circuit closes and normal operation resumes. If any request fails, the circuit reopens.

HTTP 5xx responses and transport failures, including transport timeouts, count as failures. HTTP 4xx responses and caller cancellation or deadline expiration do not change success or failure counters; an excluded half-open result releases the probe slot. Concurrent closed-state requests can overlap, and completions from an older state cannot reverse a newer transition.

Set `enabled: false` to disable the circuit breaker entirely.

<a name="web-search-settings"></a>
## Web Search Settings

```yaml
web_search:
  enabled: true
  default_count: 10
  default_recency: "noLimit"
  language: "en"
  timeout: 30
  search_engine: "search_std"
  content_size: "medium"
  cache_enabled: true
  cache_dir: "~/.config/zai/search_cache"
  cache_ttl: 24h
```

Search results are cached locally by default. Identical queries within the TTL window are served from cache without an API call. The cache key is a SHA256 hash of the query parameters.

The `default_count`, `default_recency`, `search_engine`, and `content_size` values set defaults for the `search` command. All of them can be overridden per-invocation with flags. For the full flag reference, see the [search documentation](/docs/search).

Set `enabled: false` to disable web search entirely.

<a name="web-reader-settings"></a>
## Web Reader Settings

```yaml
web_reader:
  enabled: true
  timeout: 20
  cache_enabled: true
  return_format: markdown
  auto_detect: true
  max_content_length: 50000
```

When `auto_detect` is `true`, URLs in chat prompts and in the `-f` flag are automatically fetched and included as context. Set it to `false` if you want URLs treated as literal text.

The `max_content_length` caps how much content is extracted from a page. Increase it for dense reference pages; decrease it if you find context windows filling up too fast.

For the full reader command reference, see the [reader documentation](/docs/reader).

<a name="tts-settings"></a>
## TTS Settings

```yaml
tts:
  voice: "tongtong"
  response_format: "wav"
```

These set the default voice and format for the `tts` command. Both can be overridden per-invocation with `--voice` and `--format`. For available voices, see the [TTS documentation](/docs/tts#choosing-a-voice).

<a name="media-downloads"></a>
## Media Downloads

```yaml
media:
  max_download_bytes: 536870912
```

Image and video downloads are limited to 512 MiB by default. Set a positive byte limit with `media.max_download_bytes`, `ZAI_MEDIA_MAX_DOWNLOAD_BYTES`, or the image/video `--max-download-bytes` flag. An explicit flag overrides the environment, which overrides the configuration file. Downloads exceeding the limit fail without replacing an existing output file; temporary output is removed on failure.

<a name="account-quota"></a>
## Account Quota

`zai usage` prints account quota windows, and `zai usage --json` returns a `windows` array. The command reads the monitoring endpoint without a model-generation call. It requires both configured standard and coding API URLs to be HTTPS endpoints at the same approved origin (`api.z.ai`, `open.bigmodel.cn`, or `dev.bigmodel.cn`). A custom proxy or mixed region returns an unavailable error before sending credentials. The monitoring request has no redirects or automatic retries.

```yaml
quota:
  timeout: 10s
  max_response_bytes: 262144
```

The default config is created privately on first real command execution if absent. An existing custom config is never rewritten; missing quota settings receive defaults in memory. `--config` does not create a file, and help/version does not bootstrap configuration. A monitoring endpoint may be unavailable for a given account or key; no live account entitlement is assumed. JSON failures return an `error.reason` category and a nonzero exit code; unknown values and unsupported windows are not reported as zero.

<a name="environment-variables"></a>
## Environment Variables

| Variable | Effect |
|----------|--------|
| `ZAI_API_KEY` | Overrides `api.key` in the config file |
| `ZAI_MEDIA_MAX_DOWNLOAD_BYTES` | Overrides `media.max_download_bytes` (positive bytes) |

The environment variable takes precedence when both are set. For initial setup, see the [installation documentation](/docs/installation#configuring-your-api-key).
