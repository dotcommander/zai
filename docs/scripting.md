When stdout is not a terminal, zai outputs raw text with no styling. Tokens stream through the pipe as they arrive, making zai composable with standard Unix tools and other zai invocations. This page covers pipe chaining, shell integration, and automation patterns.

- [Pipe Chaining](#pipe-chaining)
- [Shell Integration](#shell-integration)
- [Automation Patterns](#automation-patterns)
    - [Conditional Execution](#conditional-execution)
    - [Batch Processing](#batch-processing)
    - [Git Integration](#git-integration)
    - [Error Handling](#error-handling)

<a name="pipe-chaining"></a>
## Pipe Chaining

Chain multiple zai calls. Each step streams into the next:

```bash
cat code.go | zai "find bugs" | zai "suggest fixes for each bug"
```

Mix zai with standard Unix tools:

```bash
git diff HEAD~1 | zai "summarize these changes" | tee summary.txt
```

Redirect streaming output to a file:

```bash
zai "Write a haiku about concurrency" > haiku.txt
```

Zai detects whether stdout is a TTY and adjusts its output accordingly. In a terminal, you get styled output with lipgloss formatting. In a pipe or redirect, you get plain text -- no ANSI codes, no decorations.

<a name="shell-integration"></a>
## Shell Integration

Create reusable shell functions by adding them to your `.bashrc` or `.zshrc`:

```bash
review_commit() {
    git show "$1" | zai "Code review this commit. Focus on bugs and security."
}

explain_file() {
    zai -f "$1" "Explain what this code does"
}
```

Then use them like any other command:

```bash
review_commit HEAD
explain_file internal/app/client.go
```

<a name="automation-patterns"></a>
## Automation Patterns

<a name="conditional-execution"></a>
### Conditional Execution

Branch on zai's output to drive decisions in scripts:

```bash
if zai "Is this code safe? $(cat script.sh)" | grep -qi "safe"; then
    chmod +x script.sh
fi
```

<a name="batch-processing"></a>
### Batch Processing

Process multiple files in parallel:

```bash
find . -name "*.go" | xargs -P 4 -I {} sh -c 'cat {} | zai "summarize" > {}.summary'
```

<a name="git-integration"></a>
### Git Integration

Generate weekly development summaries:

```bash
git log --since="1 week ago" --pretty=format:"%h %s" | \
    zai "Summarize this week's commits" > weekly-summary.md
```

<a name="error-handling"></a>
### Error Handling

Check zai's exit code in scripts to handle failures gracefully:

```bash
if ! zai -f broken.go "find bugs" > review.txt 2>&1; then
    echo "ZAI failed -- check API key or network connection"
    exit 1
fi
```

> [!NOTE]
> Zai exits with a non-zero status when it encounters an error (network failure, invalid API key, etc.). Always check the exit code in automated pipelines.

For information on configuring timeouts, retries, and rate limits for batch workloads, see the [configuration documentation](/docs/configuration).
