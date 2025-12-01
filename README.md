# ZAI CLI

Command-line interface for Z.AI GLM models.

## Installation

```bash
go build -o zai .
```

## Configuration

```bash
# Environment variable (recommended)
export ZAI_API_KEY="your_api_key"

# Or config file: ~/.config/zai/config.yaml
api:
  key: "your_api_key"
  base_url: "https://api.z.ai/api/paas/v4"
  model: "glm-4.6"
```

## Usage

### One-shot Mode
```bash
zai "Explain quantum computing"
zai -f main.go "Explain this code"
zai -v "What is 2+2?"  # verbose
```

### Interactive REPL
```bash
zai chat
zai chat -f main.go  # with file context
```

REPL commands: `help`, `history`, `context`, `clear`, `exit`

### History
```bash
zai history
zai history -l 0  # show all
```

## Flags

| Flag | Description |
|------|-------------|
| `-f, --file` | Include file contents in prompt |
| `-v, --verbose` | Show debug info and token usage |
| `--config` | Custom config file path |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ZAI_API_KEY` | API key (required) |
| `ZAI_API_BASE_URL` | Override base URL |
| `ZAI_API_MODEL` | Override model |

## Architecture

```
zai/
├── cmd/
│   ├── root.go     # CLI setup, one-shot mode
│   ├── chat.go     # Interactive REPL
│   └── history.go  # History command
├── internal/app/
│   ├── client.go   # API client (DI, interfaces)
│   ├── history.go  # JSONL history store
│   └── types.go    # Request/Response types
└── main.go
```

## Requirements

- Go 1.21+
- Z.AI API key
