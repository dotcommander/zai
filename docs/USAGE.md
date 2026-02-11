# Usage Guide

Comprehensive usage documentation for zai CLI commands, configuration, and advanced patterns.

## Table of Contents

- [Commands](#commands)
  - [Search Command](#search-command)
  - [Reader Command](#reader-command)
- [Configuration](#configuration)
  - [Rate Limiting](#rate-limiting)
  - [Web Search](#web-search)
- [Advanced Usage](#advanced-usage)
  - [Piping and Chaining](#piping-and-chaining)
  - [Shell Script Integration](#shell-script-integration)
  - [File Processing Pipelines](#file-processing-pipelines)
  - [Automation Patterns](#automation-patterns)

---

## Commands

### Search Command

Search the web using the Z.AI search engine with results optimized for LLM consumption.

#### Usage

```bash
zai search [query] [flags]
```

The search query can be provided as a command-line argument or piped via stdin.

#### Options

##### `-c, --count <number>`

Number of search results to return (1-50). If not specified, uses the default from configuration (`web_search.default_count`).

```bash
zai search "golang patterns" -c 5
```

##### `-r, --recency <filter>`

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

##### `-d, --domain <domain>`

Limit search results to a specific domain. Useful for searching within a particular website.

```bash
zai search "kubernetes tutorials" -d kubernetes.io
zai search "github.com" -d github.com
```

##### `-o, --format <format>`

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

#### Examples

**Basic Search:**
```bash
zai search "golang best practices"
```

**Search from Stdin:**
```bash
echo "rust ownership patterns" | zai search
```

**Limit Results and Timeframe:**
```bash
zai search "latest AI developments" -c 10 -r oneWeek
```

**Domain-Specific Search:**
```bash
# Search only GitHub
zai search "fiber web framework" -d github.com

# Search only a specific site
zai search "kubernetes deployment" -d kubernetes.io
```

**JSON Output for Scripting:**
```bash
zai search "docker compose" -o json
```

Output:
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

**Verbose Output:**
```bash
zai search "golang generics" -v
```

Shows search duration, result count, and content previews.

#### Search in Chat Mode

Web search is also available within interactive chat sessions:

```bash
zai chat
you> search "latest golang release" -c 5 -r oneWeek
you> /search "rust async patterns" -r oneMonth
```

---

### Reader Command

Fetch and display web content from URLs using Z.AI's web reader API. The reader extracts the main content from web pages, filtering out navigation, ads, and other clutter.

#### Usage

```bash
zai reader <url> [flags]
```

The URL argument is required and must be a valid HTTP/HTTPS URL.

#### Options

##### `--format <format>`

Return format for the fetched content.

| Format | Description |
|--------|-------------|
| `markdown` | GitHub Flavored Markdown (default) |
| `text` | Plain text without formatting |

```bash
zai reader https://example.com --format markdown
zai reader https://example.com --format text
```

##### `--timeout <seconds>`

Request timeout in seconds. Must be a positive integer. Default: `20`.

```bash
zai reader https://example.com --timeout 30
```

##### `--no-cache`

Disable caching for this request. By default, the reader caches responses to speed up repeated requests.

```bash
zai reader https://example.com --no-cache
```

##### `--no-gfm`

Disable GitHub Flavored Markdown formatting. Only applies when `--format markdown` is used.

```bash
zai reader https://example.com --no-gfm
```

##### `--keep-img-data-url`

Preserve image data URLs in the output. By default, image data URLs are stripped.

```bash
zai reader https://example.com --keep-img-data-url
```

##### `--with-images-summary`

Include a summary of images found in the content.

```bash
zai reader https://example.com --with-images-summary
```

##### `--with-links-summary`

Include a summary of links found in the content.

```bash
zai reader https://example.com --with-links-summary
```

##### `--no-retain-images`

Do not retain images in the processed content. Default is to retain images.

```bash
zai reader https://example.com --no-retain-images
```

##### `--json`

Output the fetched content in JSON format instead of human-readable text. Useful for scripting and programmatic processing.

```bash
zai reader https://example.com --json
```

Output:
```json
{
  "url": "https://example.com",
  "title": "Example Domain",
  "description": "This domain is for use in illustrative examples...",
  "content": "# Example Domain\n\nThis domain is...",
  "metadata": {
    "author": "Example Author",
    "publishedTime": "2024-01-15T10:00:00Z"
  },
  "external_resources": {
    "images": 5,
    "links": 12
  },
  "timestamp": "2026-01-18T22:00:00Z"
}
```

#### Examples

**Basic Web Content:**
```bash
zai reader https://example.com
```

Fetches and displays the content in Markdown format with default settings.

**Plain Text Output:**
```bash
zai reader https://example.com --format text
```

Returns the content as plain text without Markdown formatting.

**Custom Timeout:**
```bash
zai reader https://example.com --timeout 60
```

Increases the timeout to 60 seconds for slow-loading pages.

**Bypass Cache:**
```bash
zai reader https://example.com --no-cache
```

Forces a fresh fetch, bypassing any cached response.

**Include Summaries:**
```bash
zai reader https://blog.example.com --with-images-summary --with-links-summary
```

Fetches content and includes summaries of images and links found in the page.

**JSON Output for Scripting:**
```bash
zai reader https://example.com --json | jq '.title'
```

Parses the JSON output and extracts the title using `jq`.

**Disable Image Retention:**
```bash
zai reader https://example.com --no-retain-images
```

Fetches content without preserving embedded images.

#### Output Format

By default, the reader outputs content in a human-readable format:

```
Title: Example Page Title
URL: https://example.com
Description: A brief description of the page content

Content:
# Main Heading

The extracted content appears here in Markdown format...

Metadata:
  author: John Doe
  publishedTime: 2024-01-15T10:00:00Z

External Resources:
  images: 5
  links: 12
```

#### Auto-Detection in Chat

URLs in chat prompts are automatically fetched and included:

```bash
zai chat
you> Summarize https://example.com
```

The content is automatically fetched and included in the prompt context.

#### File Flag Support

The `-f` flag in other commands also supports URLs:

```bash
zai chat -f https://example.com "What is this about?"
```

---

## Configuration

### Rate Limiting

The zai CLI tool implements rate limiting for all API calls to prevent overwhelming the Z.AI API and ensure fair usage.

#### Configuration

Rate limiting is configured in the `api.rate_limit` section of your config file (`~/.config/zai/config.yaml`):

```yaml
api:
  rate_limit:
    requests_per_second: 10  # Maximum requests per second
    burst: 5                 # Maximum burst requests
```

#### Default Values

- `requests_per_second`: 10
- `burst`: 5

#### Disabling Rate Limiting

To disable rate limiting, set `requests_per_second` to 0:

```yaml
api:
  rate_limit:
    requests_per_second: 0
    burst: 0
```

#### How It Works

The rate limiter uses a token bucket algorithm implemented with `golang.org/x/time/rate`:

1. **Token Bucket**: A bucket with tokens that refills at a constant rate
2. **Request Consumption**: Each API request consumes one token
3. **Wait Behavior**: If no tokens are available, requests wait until tokens are available
4. **Burst Handling**: The burst size allows short bursts of requests above the sustained rate

#### Behavior Examples

**Default Configuration (10 req/sec, burst 5):**
- Requests 1-5: Execute immediately (burst capacity)
- Requests 6+: Wait 0.1 seconds between requests (rate limiting)
- Concurrent requests: Properly synchronized across all API calls

**Disabled (0 req/sec):**
- All execute immediately with no waiting

#### Logging

When rate limiting is active, you'll see debug log messages when requests are delayed:

```
rate limit exceeded: rate limit exceeded
```

This helps identify when the rate limiting is affecting your workflow.

#### API Coverage

Rate limiting applies to all API calls:
- Chat completions
- Image generation
- Web search
- Web content fetching
- Audio transcription
- Video generation
- Vision analysis
- Model listing

#### Best Practices

1. **For most users**: Keep the default settings
2. **For batch processing**: Consider increasing the rate limit temporarily
3. **For scripts**: Handle rate limiting errors appropriately
4. **For development**: Use `verbose` mode to see rate limiting activity

---

### Web Search

Web search behavior is configured via the `web_search` section in `~/.config/zai/config.yaml`.

#### Configuration Keys

##### `enabled`

**Type:** `boolean`
**Default:** `true`

Enable or disable web search functionality globally.

```yaml
web_search:
  enabled: true
```

When disabled, all search commands (`zai search`, `/search` in chat) will fail with an error message.

##### `default_count`

**Type:** `integer`
**Default:** `10`
**Range:** `1-50`

Default number of search results to return per search query.

```yaml
web_search:
  default_count: 10
```

This value is used when the `-c, --count` flag is not specified in a search command. Higher values provide more comprehensive results but increase API usage and processing time.

##### `default_recency`

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

##### `timeout`

**Type:** `integer`
**Default:** `30`
**Unit:** seconds

Maximum duration to wait for search API responses.

```yaml
web_search:
  timeout: 30
```

Increase this value if you experience timeouts on slower network connections.

#### Caching Configuration

##### `cache_enabled`

**Type:** `boolean`
**Default:** `true`

Enable or disable search result caching.

```yaml
web_search:
  cache_enabled: true
```

When enabled, search results are cached locally to avoid redundant API calls for identical queries. The cache key includes the query string, count, recency filter, and domain filter.

##### `cache_dir`

**Type:** `string`
**Default:** `~/.config/zai/search_cache`

Directory where cached search results are stored.

```yaml
web_search:
  cache_dir: "~/.config/zai/search_cache"
```

Cache files are stored as JSON files named by the SHA256 hash of the search parameters. The path is expanded for `~` (home directory) and environment variables.

##### `cache_ttl`

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

#### Complete Example

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

#### Performance Considerations

- **Higher `default_count`** increases API usage and processing time
- **Shorter `cache_ttl`** increases API usage (more frequent re-fetches)
- **`cache_enabled: false`** significantly increases API costs for repeated queries
- **`timeout`** too low may cause failures on slow connections; too high may cause long waits on API failures

#### Web Reader Configuration

Reader behavior can be configured via `~/.config/zai/config.yaml`:

```yaml
web_reader:
  enabled: true           # Enable/disable web reader
  timeout: 20            # Default timeout in seconds
  cache_enabled: true    # Enable response caching
  return_format: markdown # Default format (markdown or text)
  auto_detect: true      # Auto-detect URLs in chat
  max_content_length: 50000 # Max characters to include
```

---

## Advanced Usage

### Piping and Chaining

#### Chain Multiple Commands

Pipe output between commands for multi-step processing:

```bash
# Chain zai with other tools
cat largefile.log | zai "extract errors" | grep ERROR | sort | uniq -c

# Process git output
git diff HEAD~1 | zai "summarize changes" | tee summary.txt

# Combine with jq for JSON processing
curl -s https://api.example.com/data | zai "analyze trends" | jq '.trends'
```

#### Conditional Execution

Use zai output for conditional logic:

```bash
# Only proceed if zai approves
if zai "check if this code is safe: $(cat script.sh)" | grep -q "safe"; then
    chmod +x script.sh
    ./script.sh
fi

# Retry with different prompts
for prompt in "explain" "simplify" "summarize"; do
    cat data.txt | zai "$prompt" | head -20
done
```

---

### Shell Script Integration

#### Wrapper Scripts

Create reusable shell functions:

```bash
#!/bin/bash
# zai-repo-helper.sh

# Function to review commits
review_commit() {
    local commit=$1
    git show "$commit" | zai "code review this commit. Focus on bugs and security issues."
}

# Function to analyze logs
analyze_logs() {
    local logfile=$1
    local pattern=$2
    grep "$pattern" "$logfile" | zai "analyze these log entries for trends"
}

# Function to generate docs
generate_docs() {
    local source=$1
    cat "$source" | zai "generate markdown documentation" > "${source%.go}.md"
}

# Main script
case "$1" in
    review)
        review_commit "$2"
        ;;
    logs)
        analyze_logs "$2" "$3"
        ;;
    docs)
        generate_docs "$2"
        ;;
    *)
        echo "Usage: $0 {review|logs|docs} [args...]"
        exit 1
        ;;
esac
```

#### Interactive Menus

Build interactive tools with zai:

```bash
#!/bin/bash
# zai-menu.sh

show_menu() {
    echo "=== ZAI Helper ==="
    echo "1) Code Review"
    echo "2) Generate Tests"
    echo "3) Write Documentation"
    echo "4) Optimize Code"
    echo "5) Exit"
    read -p "Choose: " choice
}

process_choice() {
    case $choice in
        1)
            read -p "File: " file
            cat "$file" | zai "code review. Focus on bugs, security, and performance."
            ;;
        2)
            read -p "File: " file
            cat "$file" | zai "generate comprehensive unit tests using table-driven tests"
            ;;
        3)
            read -p "File: " file
            cat "$file" | zai "generate markdown documentation with examples"
            ;;
        4)
            read -p "File: " file
            cat "$file" | zai "optimize for performance. Explain changes."
            ;;
        5)
            exit 0
            ;;
    esac
}

while true; do
    show_menu
    process_choice
    echo ""
    read -p "Press Enter to continue..."
    clear
done
```

---

### File Processing Pipelines

#### Bulk Operations

Process entire directories:

```bash
# Generate summaries for all markdown files
find docs -name "*.md" -exec sh -c '
    cat "$1" | zai "summarize in 3 bullet points" > "$1.summary"
' _ {} \;

# Convert all .txt to .md with formatting
for file in *.txt; do
    cat "$file" | zai "convert to markdown with proper formatting" > "${file%.txt}.md"
done

# Extract and categorize code snippets
find . -name "*.go" -exec sh -c '
    cat "$1" | zai "extract and categorize all functions"
' _ {} \; > function-catalog.md
```

#### Multi-Stage Pipelines

Complex processing workflows:

```bash
# Stage 1: Extract, Stage 2: Analyze, Stage 3: Report
cat large-dataset.json | \
    zai "extract all user records" | \
    zai "analyze for patterns" | \
    zai "generate markdown report" > analysis-report.md

# Documentation generation pipeline
find src -name "*.go" | while read file; do
    echo "## $(basename $file)" >> docs/API.md
    cat "$file" | zai "extract function signatures and doc comments" >> docs/API.md
    echo "" >> docs/API.md
done
```

---

### Automation Patterns

#### Cron Jobs

Scheduled tasks with zai:

```bash
# Daily log analysis (add to crontab: 0 2 * * *)
#!/bin/bash
LOG_FILE="/var/log/app.log"
REPORT_DIR="/var/reports"
DATE=$(date +%Y-%m-%d)

tail -n 10000 "$LOG_FILE" | \
    zai "analyze errors and warnings, generate summary report" \
    > "$REPORT_DIR/daily-$DATE.md"

# Weekly code review
#!/bin/bash
git log --since="1 week ago" --pretty=format:"%h %s" | \
    zai "review commits for quality issues" \
    > /tmp/weekly-review.txt

# Hourly health check
#!/bin/bash
curl -s http://localhost:8080/health | \
    zai "check if health metrics are within normal range" | \
    mail -s "Health Check Report" admin@example.com
```

#### Git Hooks

Integrate zai into git workflow:

```bash
# .git/hooks/pre-commit
#!/bin/bash
# Run zai on staged files

STAGED_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')

if [ -n "$STAGED_FILES" ]; then
    echo "Running zai code review on staged files..."
    for file in $STAGED_FILES; do
        REVIEW=$(git show ":$file" | zai "quick code review. Check only for critical bugs.")
        if echo "$REVIEW" | grep -qi "critical\|security\|vulnerability"; then
            echo "⚠️  Critical issues found in $file:"
            echo "$REVIEW"
            echo ""
            read -p "Commit anyway? (y/N) " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                exit 1
            fi
        fi
    done
fi
```

#### CI/CD Integration

```bash
# .github/workflows/zai-review.yml
name: ZAI Code Review
on: [pull_request]

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Review changed files
        run: |
          git diff origin/main --name-only | grep '\.go$' | \
          xargs -I {} sh -c 'cat {} | zai "review for bugs" >> review.txt'
      - name: Comment on PR
        run: |
          gh pr comment ${{ github.event.number }} --body-file review.txt
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

### Tips and Best Practices

#### Performance Optimization

```bash
# Batch requests instead of individual calls
# ❌ Slow
for file in *.go; do zai "summarize $file" < "$file"; done

# ✅ Fast (parallel)
find . -name "*.go" | parallel -j 4 'cat {} | zai "summarize"'
```

#### Error Handling

```bash
# Always check exit status
if ! cat file.go | zai "review" > review.txt; then
    echo "ZAI failed - check API key or connection"
    exit 1
fi

# Retry logic
MAX_RETRIES=3
for i in $(seq 1 $MAX_RETRIES); do
    if cat data.txt | zai "process" > output.txt; then
        break
    fi
    echo "Retry $i/$MAX_RETRIES"
    sleep 5
done
```

#### Caching

```bash
# Cache results to avoid redundant API calls
cache_zai() {
    local prompt=$1
    local cache_key=$(echo "$prompt" | md5sum | cut -d' ' -f1)
    local cache_file="$HOME/.cache/zai/$cache_key"

    if [ -f "$cache_file" ]; then
        cat "$cache_file"
        return
    fi

    zai "$prompt" | tee "$cache_file"
}
```

#### Memory Management

```bash
# Process large files in chunks to avoid memory issues
split -l 1000 largefile.txt chunk_
for chunk in chunk_*; do
    cat "$chunk" | zai "process" >> output.txt
    rm "$chunk"
done
```
