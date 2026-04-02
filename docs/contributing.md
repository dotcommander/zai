Contributing to zai starts here. This guide covers development setup, code style, testing expectations, and the pull request process. Follow these guidelines to keep reviews fast and merges clean.

- [Development Setup](#development-setup)
    - [Building](#building)
    - [Running Tests](#running-tests)
- [Project Structure](#project-structure)
- [Code Style](#code-style)
    - [Naming](#naming)
    - [Error Handling](#error-handling)
- [Testing](#testing)
    - [Writing Tests](#writing-tests)
    - [Coverage](#coverage)
- [Pull Request Process](#pull-request-process)
    - [Branch Naming](#branch-naming)
    - [Before Submitting](#before-submitting)
- [Security](#security)

<a name="development-setup"></a>
## Development Setup

You need Go 1.26 or later and git. Fork the repository, clone your fork, and add the upstream remote:

```bash
git clone https://github.com/your-username/zai.git
cd zai
git remote add upstream https://github.com/dotcommander/zai.git
```

<a name="building"></a>
### Building

```bash
go build -o bin/zai .
```

To install locally for testing:

```bash
go build -o bin/zai . && ln -sf $(pwd)/bin/zai ~/go/bin/zai
```

<a name="running-tests"></a>
### Running Tests

```bash
go test ./...
go vet ./...
golangci-lint run ./...
```

All three must pass before submitting a PR.

<a name="project-structure"></a>
## Project Structure

```
zai/
├── cmd/                        # CLI commands (thin cobra wrappers)
│   ├── root.go                # Main command, stdin handling, one-shot mode
│   ├── chat.go                # Interactive REPL
│   ├── search.go              # Web search
│   ├── web.go                 # Web reader
│   ├── image.go               # Image generation
│   ├── vision.go              # Vision analysis
│   ├── audio.go               # Audio transcription
│   ├── tts.go                 # Text-to-speech
│   ├── embed.go               # Text embeddings
│   ├── video.go               # Video generation
│   └── model.go               # Model management
├── internal/
│   └── app/
│       ├── client.go          # Core HTTP client, interfaces, shared helpers
│       ├── client_chat.go     # Chat and streaming
│       ├── client_image.go    # Image generation
│       ├── client_vision.go   # Vision analysis
│       ├── client_web.go      # Web reader and search
│       ├── client_audio.go    # Audio transcription
│       ├── client_video.go    # Video generation
│       ├── client_tts.go      # Text-to-speech
│       ├── client_embedding.go # Embeddings
│       ├── types.go           # Request/response types
│       ├── history.go         # History storage
│       ├── cache.go           # Search caching
│       └── utils.go           # URL detection, formatting
├── internal/config/
│   └── config.go              # Viper defaults and loading
├── docs/                      # User-facing documentation
└── CLAUDE.md                  # Project reference
```

The client is split by domain -- each `client_*.go` file handles one API surface. Interfaces are segregated per domain (ChatClient, VisionClient, etc.) and defined in `client.go`.

<a name="code-style"></a>
## Code Style

Follow [Effective Go](https://go.dev/doc/effective_go) and the project's `CLAUDE.md` for detailed conventions. The short version:

- Format with `gofmt` and `goimports`
- Run `golangci-lint run ./...` before committing
- Functions under 80 lines, nesting under 4 levels
- Error path first: `if err != nil { return }`

<a name="naming"></a>
### Naming

- **Packages**: lowercase, single word
- **Exports**: PascalCase (`ChatClient`)
- **Private**: camelCase (`buildContent`)
- **Interfaces**: end with behavior or capability (`HTTPDoer`, `FileReader`)

<a name="error-handling"></a>
### Error Handling

Wrap errors with context using `fmt.Errorf("operation: %w", err)`. Use sentinel errors for expected cases that callers need to branch on.

<a name="testing"></a>
## Testing

<a name="writing-tests"></a>
### Writing Tests

Every test and subtest must call `t.Parallel()`. Use table-driven tests for multiple cases. Mock external dependencies through interfaces -- no real I/O in unit tests.

```go
func TestFunctionName(t *testing.T) {
    t.Parallel()
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {"happy path", validInput, expectedOutput, false},
        {"error case", invalidInput, OutputType{}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            got, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

<a name="coverage"></a>
### Coverage

Aim for 80% coverage on new code:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

<a name="pull-request-process"></a>
## Pull Request Process

<a name="branch-naming"></a>
### Branch Naming

- `feat/feature-name` -- new features
- `fix/bug-description` -- bug fixes
- `docs/update-name` -- documentation
- `refactor/component-name` -- refactoring

Use conventional commit messages: `feat(scope): description`, `fix(scope): description`.

<a name="before-submitting"></a>
### Before Submitting

- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] `golangci-lint run ./...` clean
- [ ] New tests added for new features
- [ ] Documentation updated if behavior changed
- [ ] Commit messages follow conventional format
- [ ] No merge conflicts with upstream/main

Push your branch and open a PR on GitHub. Describe the change, link related issues, and list any breaking changes.

<a name="security"></a>
## Security

If you discover a security vulnerability, email the maintainers privately rather than opening an issue.

Never commit API keys. Use environment variables or config files. Validate user inputs at system boundaries. Review HTTP request construction for injection risks. Keep dependencies current with `go get -u ./... && go mod tidy`.
