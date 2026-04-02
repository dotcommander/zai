Chatting with models is zai's primary function. You send a prompt, and zai streams back the response token by token -- you start reading before the model finishes thinking. Zai supports one-shot prompts, a persistent REPL with conversation memory, file and URL context, custom system prompts, reasoning mode, and a coding-optimized endpoint.

- [One-Shot Mode](#one-shot-mode)
    - [Stdin Input](#stdin-input)
- [Interactive REPL](#interactive-repl)
    - [REPL Commands](#repl-commands)
- [Streaming vs JSON Output](#streaming-vs-json-output)
- [File Context](#file-context)
- [Custom System Prompts](#custom-system-prompts)
- [Reasoning Mode](#reasoning-mode)
- [Coding API](#coding-api)
- [Pipe Detection](#pipe-detection)
- [The --search Flag vs the search Command](#search-flag-vs-search-command)
- [History](#history)
- [Model Management](#model-management)

<a name="one-shot-mode"></a>
## One-Shot Mode

Pass a prompt as a positional argument. Zai streams the response to stdout, then exits:

```bash
zai "What is the capital of France?"
```

You may attach a file for context with `-f`:

```bash
zai -f main.go "Explain what this code does"
```

<a name="stdin-input"></a>
### Stdin Input

You may pipe input from stdin. Zai wraps the piped data as context and combines it with any positional prompt:

```bash
echo "func main() { fmt.Println(42) }" | zai "What does this print?"
```

When stdin provides context and you also pass a positional argument, both are sent together:

```bash
cat server.go | zai "Find security issues"
```

For more on piping and automation, see the [scripting documentation](/docs/scripting).

<a name="interactive-repl"></a>
## Interactive REPL

Start a persistent chat session where zai retains previous messages as context:

```bash
zai chat
```

The REPL keeps the last 10 exchanges (20 messages) as conversation context, so follow-up questions reference earlier parts of the discussion naturally.

You may load a file into the session context on startup:

```bash
zai chat -f main.go
```

You may override the model for a session:

```bash
zai chat --model glm-4.7
```

<a name="repl-commands"></a>
### REPL Commands

| Command | Description |
|---------|-------------|
| `help` or `?` | Show available commands |
| `history` | Show messages sent this session |
| `context` | Show current conversation context |
| `clear` | Clear context and screen |
| `search <query>` | Run a web search (results added to context) |
| `web <url>` | Fetch a web page (content added to context) |
| `exit` or `quit` | End the session |

All commands work with or without a `/` prefix: `search "query"` and `/search "query"` are equivalent.

<a name="streaming-vs-json-output"></a>
## Streaming vs JSON Output

By default, all chat responses stream -- tokens appear one at a time as they are generated. This is faster for interactive use because you start reading immediately.

When you need a complete, structured response for scripts or further processing, use `--json`:

```bash
zai "What is 2+2?" --json
```

Zai returns a JSON object containing the prompt, response, model, and metadata. JSON mode disables streaming because the entire response must be assembled before it can be serialized.

<a name="file-context"></a>
## File Context

The `-f` flag attaches file contents (or web page contents) to your prompt. It works in both one-shot and REPL modes:

```bash
zai -f utils.go "Are there any bugs here?"
zai -f https://example.com/docs "Summarize this page"
```

When `-f` receives an HTTP or HTTPS URL, zai fetches the page through the web reader API automatically. For more on web content fetching, see the [reader documentation](/docs/reader).

<a name="custom-system-prompts"></a>
## Custom System Prompts

Override the default system prompt with `--system`:

```bash
zai --system "You are a senior Go developer. Be concise." "Review this function"
```

The value can be a literal string or a path to a file:

```bash
zai --system ~/prompts/reviewer.txt -f main.go "Review this"
```

If the path does not resolve to a file, zai falls back to its built-in default prompt.

<a name="reasoning-mode"></a>
## Reasoning Mode

Enable extended reasoning with `--think`. The model takes longer but produces more thorough analysis:

```bash
zai --think "Prove that the square root of 2 is irrational"
```

This works in both one-shot and REPL modes:

```bash
zai chat --think
```

<a name="coding-api"></a>
## Coding API

Use the `-C` (or `--coding`) flag to route requests through the coding-optimized endpoint:

```bash
zai -C "Write a binary search in Go"
```

The coding endpoint uses a separate base URL configured as `api.coding_base_url` in your config. For details, see the [configuration documentation](/docs/configuration#api-settings).

<a name="pipe-detection"></a>
## Pipe Detection

When stdout is not a terminal (pipe or redirect), zai outputs raw text with no styling. Tokens stream through the pipe as they arrive, making zai composable with other tools. For pipe chaining patterns and shell integration, see the [scripting documentation](/docs/scripting).

<a name="search-flag-vs-search-command"></a>
## The --search Flag vs the search Command

These two features both involve web search, but they serve different purposes.

The `--search` flag augments your chat prompt with search results automatically. You ask a question; zai searches for context behind the scenes and includes it in the prompt:

```bash
zai --search "What happened at GopherCon 2026?"
```

You never see the raw search results. They are injected as context for the model to reference.

The `search` command performs a standalone web search and displays the results directly to you:

```bash
zai search "GopherCon 2026"
```

Use `--search` when you want the model to answer using current information. Use `search` when you want to browse results yourself. For the full search command reference, see the [search documentation](/docs/search).

<a name="history"></a>
## History

View your command history across all session types:

```bash
zai history
```

Zai displays a table of recent prompts with timestamps, types, models, and truncated responses.

Show more entries with `-l`:

```bash
zai history -l 50
zai history -l 0
```

Passing `-l 0` shows all entries. Add `--json` for machine-readable output.

History is stored at `~/.config/zai/history.jsonl`.

<a name="model-management"></a>
## Model Management

List available models from the Z.AI API:

```bash
zai model list
```

Zai displays each model name and creation date. Add `--json` for machine-readable output.

Each command uses a default model from your configuration, and any command accepts `-m` to override the model for that invocation. For model configuration, see the [configuration documentation](/docs/configuration#api-settings).
