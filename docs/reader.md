The reader command fetches a web page, strips navigation and ads, and returns the main content in clean markdown or plain text. It is the same mechanism zai uses behind the scenes when you include a URL in a chat prompt or pass one to the `-f` flag.

- [Fetching a Page](#fetching-a-page)
- [Output Options](#output-options)
- [Additional Flags](#additional-flags)
- [Auto-Detection in Chat](#auto-detection-in-chat)

<a name="fetching-a-page"></a>
## Fetching a Page

```bash
zai reader https://go.dev/blog/go1.22
```

Zai fetches the page, extracts the title, URL, description, and primary content, then prints it to your terminal. The output defaults to GitHub Flavored Markdown.

<a name="output-options"></a>
## Output Options

Change the format between markdown (default) and plain text:

```bash
zai reader https://example.com --format text
```

Increase the timeout for slow-loading pages:

```bash
zai reader https://example.com --timeout 60
```

Force a fresh fetch, bypassing the cache:

```bash
zai reader https://example.com --no-cache
```

Include summaries of images and links found on the page:

```bash
zai reader https://example.com --with-images-summary --with-links-summary
```

Get JSON output for programmatic use:

```bash
zai reader https://example.com --json
```

<a name="additional-flags"></a>
## Additional Flags

| Flag | Effect |
|------|--------|
| `--no-gfm` | Disable GitHub Flavored Markdown |
| `--keep-img-data-url` | Preserve image data URLs |
| `--no-retain-images` | Strip images from output |

<a name="auto-detection-in-chat"></a>
## Auto-Detection in Chat

When `web_reader.auto_detect` is enabled in your configuration (it is by default), URLs in chat prompts and in the `-f` flag are automatically fetched and included as context:

```bash
zai "Summarize https://go.dev/blog/go1.22"
```

```bash
zai chat -f https://go.dev/blog/go1.22 "What changed?"
```

Inside the REPL, you can also type a URL directly in your message and zai fetches it transparently. The `web <url>` REPL command does the same thing explicitly. For web reader configuration, see the [configuration documentation](/docs/configuration#web-reader-settings).
