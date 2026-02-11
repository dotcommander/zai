# Advanced Usage Patterns

This guide covers advanced workflows including piping, scripting, automation, and shell integration patterns.

## Table of Contents

- [Piping and Chaining](#piping-and-chaining)
- [Stdin Processing](#stdin-processing)
- [Shell Script Integration](#shell-script-integration)
- [File Processing Pipelines](#file-processing-pipelines)
- [Automation Patterns](#automation-patterns)
- [Output Capture and Parsing](#output-capture-and-parsing)

---

## Piping and Chaining

### Chain Multiple Commands

Pipe output between commands for multi-step processing:

```bash
# Chain zai with other tools
cat largefile.log | zai "extract errors" | grep ERROR | sort | uniq -c

# Process git output
git diff HEAD~1 | zai "summarize changes" | tee summary.txt

# Combine with jq for JSON processing
curl -s https://api.example.com/data | zai "analyze trends" | jq '.trends'
```

### Conditional Execution

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

## Stdin Processing

### File Processing Loops

Process multiple files through zai:

```bash
# Process all Go files in current directory
for file in *.go; do
    echo "=== $file ===" >> analysis.md
    cat "$file" | zai "code review this file" >> analysis.md
    echo "" >> analysis.md
done

# Parallel processing with xargs
find . -name "*.md" -print0 | xargs -0 -P4 -I{} sh -c 'cat {} | zai "summarize" > {}.summary'

# Batch processing with parallel
ls *.json | parallel -j 4 'cat {} | zai "extract key metrics" > {.}-metrics.txt'
```

### Stream Processing

Process streaming data with zai:

```bash
# Process log file in chunks
tail -f /var/log/app.log | while read line; do
    echo "$line" | zai "categorize log entry"
done

# Process CSV row by row
cat data.csv | zai "format as table" | split -l 1000 - batch-

# Continuous monitoring
watch -n 10 'ps aux | zai "detect anomalies" | head -5'
```

---

## Shell Script Integration

### Wrapper Scripts

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

### Interactive Menus

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

## File Processing Pipelines

### Bulk Operations

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

### Multi-Stage Pipelines

Complex processing workflows:

```bash
# Stage 1: Extract, Stage 2: Analyze, Stage 3: Report
cat large-dataset.json | \
    zai "extract all user records" | \
    zai "analyze for patterns" | \
    zai "generate markdown report" > analysis-report.md

# Image batch processing with video generation
ls images/*.png | while read img; do
    base=$(basename "$img" .png)
    zai image "$(cat "$img" | zai "describe this image for video generation")" \
        -o "$base.mp4"
done

# Documentation generation pipeline
find src -name "*.go" | while read file; do
    echo "## $(basename $file)" >> docs/API.md
    cat "$file" | zai "extract function signatures and doc comments" >> docs/API.md
    echo "" >> docs/API.md
done
```

---

## Automation Patterns

### Cron Jobs

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

### Git Hooks

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

### CI/CD Integration

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

## Output Capture and Parsing

### JSON Processing

Parse structured zai output:

```bash
# Request JSON output and process
cat data.json | zai "analyze and return results as JSON list" | \
    jq -r '.[] | .category + ": " + .description'

# Store results for later processing
RESULT=$(cat file.go | zai "return function list as JSON array")
echo "$RESULT" | jq '.[] | select(.complexity > 10)'

# Merge with existing JSON
jq -s '.[0] + .[1]' existing.json <(cat data.txt | zai "convert to JSON")
```

### Output Formatting

Transform zai output for different tools:

```bash
# Format for Jira tickets
cat bug-report.txt | zai "format as Jira ticket description" | \
    sed 's/^$/\n---\n/' > ticket.md

# Format for email
cat summary.txt | zai "format as professional email" | \
    mail -s "Weekly Report" team@example.com

# Format for Slack
cat update.txt | zai "format as Slack message with emojis" | \
    curl -X POST -d @- $SLACK_WEBHOOK_URL
```

### Logging and Auditing

Track zai usage:

```bash
# Log all zai queries with timestamps
zai_wrapper() {
    local prompt="$1"
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local log_file="$HOME/.zai/history/$timestamp.log"

    {
        echo "=== Query ==="
        echo "$prompt"
        echo ""
        echo "=== Response ==="
        zai "$prompt"
        echo ""
        echo "=== Metadata ==="
        echo "Time: $(date)"
        echo "User: $(whoami)"
        echo "PWD: $(pwd)"
    } | tee -a "$log_file"
}

# Audit trail for critical operations
cat deploy-script.sh | zai "check for security issues" | \
    tee >(cat >> /var/log/security-audit.log) | \
    grep -i "critical" && mail -s "Security Alert" security@example.com
```

---

## Tips and Best Practices

### Performance Optimization

```bash
# Batch requests instead of individual calls
# ❌ Slow
for file in *.go; do zai "summarize $file" < "$file"; done

# ✅ Fast (parallel)
find . -name "*.go" | parallel -j 4 'cat {} | zai "summarize"'
```

### Error Handling

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

### Caching

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

### Memory Management

```bash
# Process large files in chunks to avoid memory issues
split -l 1000 largefile.txt chunk_
for chunk in chunk_*; do
    cat "$chunk" | zai "process" >> output.txt
    rm "$chunk"
done
```
