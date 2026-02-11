# Contributing to zai

Thank you for your interest in contributing to zai! This document provides guidelines and instructions for contributing effectively.

## Table of Contents

- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Code Style Standards](#code-style-standards)
- [Pull Request Process](#pull-request-process)
- [Testing](#testing)
- [Security](#security)

## Development Setup

### Prerequisites

- **Go**: 1.23 or later
- **Git**: For version control
- **make**: For build automation (optional, scripts are provided)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/your-username/zai.git
   cd zai
   ```

3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/dotcommander/zai.git
   ```

### Build the Project

```bash
go build -o bin/zai .
```

### Run Tests

```bash
go test ./...
```

### Install Locally

To test your changes locally:

```bash
go build -o bin/zai . && ln -sf $(pwd)/bin/zai ~/go/bin/zai
```

## Project Structure

```
zai/
├── cmd/                 # CLI commands
│   ├── root.go         # Main command and entry point
│   ├── chat.go         # Interactive REPL
│   ├── search.go       # Web search command
│   ├── web.go          # Web reader command
│   ├── image.go        # Image generation command
│   ├── vision.go       # Vision analysis command
│   ├── audio.go        # Audio transcription command
│   ├── video.go        # Video generation command
│   └── model.go        # Model management command
├── internal/
│   └── app/            # Application logic
│       ├── cache.go    # Search caching
│       ├── client.go   # HTTP client and API calls
│       ├── types.go    # Request/response types
│       ├── history.go  # History storage
│       └── utils.go    # Utility functions
├── docs/               # Documentation
└── CLAUDE.md          # Project documentation
```

### Architecture Principles

- **SOLID Compliance**: All components follow SOLID principles with dependency injection
- **Interface Segregation**: Client interfaces are segregated by functionality (ChatClient, VisionClient, etc.)
- **No Global State**: Configuration is injected at construction time
- **Repository Pattern**: No raw I/O outside of abstraction layers

## Code Style Standards

### Go Conventions

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Run `golint` and `go vet` before committing

### Naming

- **Packages**: lowercase, single word when possible
- **Exports**: PascalCase (e.g., `ChatClient`)
- **Private**: camelCase (e.g., `buildContent`)
- **Interfaces**:通常 end with behavior or capability (e.g., `HTTPDoer`, `FileReader`)

### Error Handling

- Always check errors
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Use sentinel errors for expected cases

### Dependencies

- Prefer standard library over external packages
- Keep dependencies minimal
- Document why non-standard dependencies are needed

### Documentation

- Exported functions must have godoc comments
- Include usage examples in godoc
- Update CLAUDE.md for architectural changes

## Pull Request Process

### Branch Naming

Use descriptive branch names:

- `feat/feature-name` - New features
- `fix/bug-description` - Bug fixes
- `docs/update-name` - Documentation updates
- `refactor/component-name` - Refactoring

### Making Changes

1. Create a branch from `main`:
   ```bash
   git checkout -b feat/your-feature-name
   ```

2. Make your changes following the code style standards

3. Test your changes:
   ```bash
   go test ./...
   go build -o bin/zai .
   ```

4. Commit with conventional commits:
   ```
   feat(scope): description
   fix(scope): description
   docs(scope): description
   ```

### Before Submitting

- [ ] Code builds successfully
- [ ] All tests pass
- [ ] New tests added for new features
- [ ] Documentation updated
- [ ] Commit messages follow conventional format
- [ ] No merge conflicts with upstream/main

### Submitting

1. Push your branch:
   ```bash
   git push origin feat/your-feature-name
   ```

2. Create a pull request on GitHub

3. Fill out the PR template:
   - Describe the change
   - Link related issues
   - List breaking changes (if any)
   - Include screenshots (if applicable)

4. Wait for review and address feedback

### Review Process

- Maintainers will review your PR
- Address all review comments
- Keep the conversation focused and constructive
- Squash commits if requested before merge

## Testing

### Writing Tests

- Test files should be named `*_test.go`
- Use table-driven tests for multiple cases
- Test both happy path and edge cases
- Mock external dependencies (HTTP clients, file readers)

### Example Test Structure

```go
func TestFunctionName(t *testing.T) {
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

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestFunctionName ./internal/app

# Verbose output
go test -v ./...
```

### Coverage

Aim for >80% test coverage for new code. View coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Security

### Reporting Vulnerabilities

If you discover a security vulnerability, please email it to the maintainers privately rather than opening an issue.

### Security Best Practices

- **API Keys**: Never commit API keys. Use environment variables or config files
- **Input Validation**: Always validate user inputs at system boundaries
- **Dependencies**: Keep dependencies updated; run `go mod tidy` regularly
- **Secrets**: Use `.env` files (gitignored) for local development secrets

### Code Review for Security

- Check for hardcoded secrets
- Validate file paths (prevent path traversal)
- Sanitize user inputs before command execution
- Review HTTP request construction for injection risks

## Getting Help

- Open an issue for bugs or feature requests
- Check existing documentation in `docs/` and `CLAUDE.md`
- Review existing pull requests for context

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

---

Thank you for contributing to zai!
