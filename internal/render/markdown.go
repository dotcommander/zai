package render

import (
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/x/term"
	"github.com/muesli/termenv"
)

// MarkdownStream renders streaming markdown progressively in a terminal.
// Accumulates tokens and re-renders with glamour when content changes
// meaningfully (newline or 40+ bytes of growth). Uses ANSI escape codes
// to overwrite previous output between renders.
type MarkdownStream struct {
	buf            strings.Builder
	renderer       *glamour.TermRenderer
	output         *termenv.Output
	lastLen        int
	lastNL         int // total newlines seen so far (incremental)
	lastRenderedNL int // lastNL value at last render trigger
	lastLines      int // visible line count of last rendered output
	isTTY          bool
}

// New creates a MarkdownStream. If stdout is not a TTY, markdown rendering
// is disabled and tokens pass through raw.
func New() *MarkdownStream {
	ms := &MarkdownStream{
		isTTY:  term.IsTerminal(os.Stdout.Fd()),
		output: termenv.NewOutput(os.Stdout),
	}
	if !ms.isTTY {
		return ms
	}

	width, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil || width < 40 {
		width = 80
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width-4),
	)
	if err != nil {
		ms.isTTY = false // fall back to raw
		return ms
	}
	ms.renderer = r
	return ms
}

// Token accepts a streaming token. In TTY mode it accumulates and
// periodically re-renders; in raw mode it returns the token unchanged
// for immediate printing.
func (ms *MarkdownStream) Token(token string) string {
	if !ms.isTTY {
		return token
	}

	ms.buf.WriteString(token)
	ms.lastNL += strings.Count(token, "\n")
	newLen := ms.buf.Len()

	// Render immediately on first content (eliminates pause after AI> label),
	// then on each newline or every ~40 bytes of growth for responsive streaming.
	firstContent := ms.lastLen == 0 && newLen > 0
	if firstContent || ms.lastNL > ms.lastRenderedNL || newLen-ms.lastLen >= 40 {
		ms.rerender(ms.buf.String())
		ms.lastLen = newLen
		ms.lastRenderedNL = ms.lastNL
	}
	return ""
}

// Flush performs a final glamour render of the complete accumulated buffer.
// Must be called after the last Token to ensure the full response is rendered.
// No-op in non-TTY mode (tokens already passed through raw).
func (ms *MarkdownStream) Flush() {
	if !ms.isTTY || ms.buf.Len() == 0 {
		return
	}
	ms.rerender(ms.buf.String())
}

// Content returns the raw accumulated text (for history/context storage).
func (ms *MarkdownStream) Content() string {
	return ms.buf.String()
}

func (ms *MarkdownStream) rerender(text string) {
	rendered, err := ms.renderer.Render(text)
	if err != nil {
		return // keep last good render
	}

	// Glamour wraps output in a leading \n and multiple trailing \n.
	// Strip both: the leading \n would drift content downward on re-renders
	// (since ClearLines repositions to where content started, not to the
	// blank separator line). Keep exactly one trailing \n so the cursor
	// ends on a fresh line after content.
	rendered = strings.Trim(rendered, "\n") + "\n"

	// Clear previous output. After writing lastLines lines the cursor sits on
	// the blank line below; ClearLines(lastLines) erases it plus the lastLines
	// content lines above, leaving the cursor on the line before our output.
	if ms.lastLines > 0 {
		ms.output.ClearLines(ms.lastLines)
	}

	ms.lastLines = strings.Count(rendered, "\n")
	_, _ = os.Stdout.WriteString(rendered)
}
