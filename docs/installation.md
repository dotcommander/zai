Zai installs from source with a single `go build` and needs only an API key to start working. The entire setup takes under a minute.

- [Building From Source](#building-from-source)
- [Configuring Your API Key](#configuring-your-api-key)
- [Your First Command](#your-first-command)
- [Global Flags](#global-flags)

<a name="building-from-source"></a>
## Building From Source

Build the binary and symlink it into your Go bin directory so it is available everywhere:

```bash
go build -o bin/zai .
ln -sf $(pwd)/bin/zai ~/go/bin/zai
```

Verify the installation by checking the version:

```bash
zai --version
```

Zai prints its current version and exits.

<a name="configuring-your-api-key"></a>
## Configuring Your API Key

You need a Z.AI API key before zai does anything useful. There are two ways to provide it.

**Environment variable** (recommended for scripts and CI):

```bash
export ZAI_API_KEY="your-api-key"
```

**Config file** (recommended for daily use):

Create `~/.config/zai/config.yaml` with your key:

```yaml
api:
  key: "your-api-key"
```

The environment variable takes precedence when both are set. For the full configuration reference, see the [configuration documentation](/docs/configuration).

<a name="your-first-command"></a>
## Your First Command

Send a one-shot prompt to confirm everything works:

```bash
zai "Explain what a goroutine is in two sentences"
```

Zai streams the response to your terminal as tokens are generated. You start reading immediately -- there is no wait for the full response to complete.

<a name="global-flags"></a>
## Global Flags

These flags are available on every command:

| Flag | Short | Description |
|------|-------|-------------|
| `--file <path>` | `-f` | Include file or URL contents in prompt |
| `--search` | | Augment prompt with web search results |
| `--think` | | Enable reasoning mode |
| `--coding` | `-C` | Use coding API endpoint |
| `--json` | | Output as JSON |
| `--system <prompt>` | | Custom system prompt (string or file path) |
| `--verbose` | `-v` | Show debug information |
| `--config <path>` | | Use a specific config file |
| `--help` | `-h` | Show help |
| `--version` | | Show version |
