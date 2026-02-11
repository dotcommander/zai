# Reader Command

Fetch and display web content from URLs using Z.AI's web reader API. The reader extracts the main content from web pages, filtering out navigation, ads, and other clutter.

## Usage

```bash
zai reader <url> [flags]
```

The URL argument is required and must be a valid HTTP/HTTPS URL.

## Options

### `--format <format>`

Return format for the fetched content.

| Format | Description |
|--------|-------------|
| `markdown` | GitHub Flavored Markdown (default) |
| `text` | Plain text without formatting |

```bash
zai reader https://example.com --format markdown
zai reader https://example.com --format text
```

### `--timeout <seconds>`

Request timeout in seconds. Must be a positive integer. Default: `20`.

```bash
zai reader https://example.com --timeout 30
```

### `--no-cache`

Disable caching for this request. By default, the reader caches responses to speed up repeated requests.

```bash
zai reader https://example.com --no-cache
```

### `--no-gfm`

Disable GitHub Flavored Markdown formatting. Only applies when `--format markdown` is used.

```bash
zai reader https://example.com --no-gfm
```

### `--keep-img-data-url`

Preserve image data URLs in the output. By default, image data URLs are stripped.

```bash
zai reader https://example.com --keep-img-data-url
```

### `--with-images-summary`

Include a summary of images found in the content.

```bash
zai reader https://example.com --with-images-summary
```

### `--with-links-summary`

Include a summary of links found in the content.

```bash
zai reader https://example.com --with-links-summary
```

### `--no-retain-images`

Do not retain images in the processed content. Default is to retain images.

```bash
zai reader https://example.com --no-retain-images
```

### `--json`

Output the fetched content in JSON format instead of human-readable text. Useful for scripting and programmatic processing.

```bash
zai reader https://example.com --json
```

**Output:**

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

## Examples

### Basic Web Content

```bash
zai reader https://example.com
```

Fetches and displays the content in Markdown format with default settings.

### Plain Text Output

```bash
zai reader https://example.com --format text
```

Returns the content as plain text without Markdown formatting.

### Custom Timeout

```bash
zai reader https://example.com --timeout 60
```

Increases the timeout to 60 seconds for slow-loading pages.

### Bypass Cache

```bash
zai reader https://example.com --no-cache
```

Forces a fresh fetch, bypassing any cached response.

### Include Summaries

```bash
zai reader https://blog.example.com --with-images-summary --with-links-summary
```

Fetches content and includes summaries of images and links found in the page.

### JSON Output for Scripting

```bash
zai reader https://example.com --json | jq '.title'
```

Parses the JSON output and extracts the title using `jq`.

### Disable Image Retention

```bash
zai reader https://example.com --no-retain-images
```

Fetches content without preserving embedded images.

## Output Format

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

## Configuration

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

## Auto-Detection in Chat

URLs in chat prompts are automatically fetched and included:

```bash
zai chat
you> Summarize https://example.com
```

The content is automatically fetched and included in the prompt context.

## File Flag Support

The `-f` flag in other commands also supports URLs:

```bash
zai chat -f https://example.com "What is this about?"
```

## See Also

- [Search Command](./search.md) - Search the web
- [Chat Command](./chat.md) - Interactive AI chat
- [Vision Command](./vision.md) - Analyze images from URLs
