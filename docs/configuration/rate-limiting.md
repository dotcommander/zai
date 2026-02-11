# Rate Limiting Configuration

The zai CLI tool implements rate limiting for all API calls to prevent overwhelming the Z.AI API and ensure fair usage.

## Configuration

Rate limiting is configured in the `api.rate_limit` section of your config file (`~/.config/zai/config.yaml`):

```yaml
api:
  rate_limit:
    requests_per_second: 10  # Maximum requests per second
    burst: 5                 # Maximum burst requests
```

### Default Values

- `requests_per_second`: 10
- `burst`: 5

### Disabling Rate Limiting

To disable rate limiting, set `requests_per_second` to 0:

```yaml
api:
  rate_limit:
    requests_per_second: 0
    burst: 0
```

## How It Works

The rate limiter uses a token bucket algorithm implemented with `golang.org/x/time/rate`:

1. **Token Bucket**: A bucket with tokens that refills at a constant rate
2. **Request Consumption**: Each API request consumes one token
3. **Wait Behavior**: If no tokens are available, requests wait until tokens are available
4. **Burst Handling**: The burst size allows short bursts of requests above the sustained rate

## Behavior Examples

### Default Configuration (10 req/sec, burst 5)

- Requests 1-5: Execute immediately (burst capacity)
- Requests 6+: Wait 0.1 seconds between requests (rate limiting)
- Concurrent requests: Properly synchronized across all API calls

### Disabled (0 req/sec)

- All execute immediately with no waiting

## Logging

When rate limiting is active, you'll see debug log messages when requests are delayed:

```
rate limit exceeded: rate limit exceeded
```

This helps identify when the rate limiting is affecting your workflow.

## API Coverage

Rate limiting applies to all API calls:
- Chat completions
- Image generation
- Web search
- Web content fetching
- Audio transcription
- Video generation
- Vision analysis
- Model listing

## Best Practices

1. **For most users**: Keep the default settings
2. **For batch processing**: Consider increasing the rate limit temporarily
3. **For scripts**: Handle rate limiting errors appropriately
4. **For development**: Use `verbose` mode to see rate limiting activity