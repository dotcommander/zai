package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"

	"github.com/dotcommander/zai/internal/app"
)

// chatMessage holds a single message in the conversation display.
type chatMessage struct {
	role     string // "user", "assistant", "system"
	content  string // raw text
	rendered string // glamour-rendered (assistant only)
}

// chatModel is the bubbletea model for the interactive chat REPL.
type chatModel struct {
	// UI
	viewport viewport.Model
	input    textarea.Model
	ready    bool

	// Chat state
	client        *app.Client
	opts          app.ChatOptions
	chatCtx       []app.Message // API conversation context
	sessionHist   []string      // raw input history
	searchEnabled bool
	messages      []chatMessage

	// Streaming
	streaming bool
	sc        streamChans
	streamBuf strings.Builder
	streamCtx context.Context //nolint:containedctx // bubbletea model needs cancel handle
	cancelFn  context.CancelFunc

	// Glamour
	glamRenderer *glamour.TermRenderer

	// Layout
	width  int
	height int
}

func newChatModel(client *app.Client, opts app.ChatOptions, searchEnabled bool) *chatModel {
	ti := textarea.New()
	ti.Placeholder = "Type a message..."
	ti.Prompt = theme.Prompt.Render("you> ")
	ti.ShowLineNumbers = false
	ti.SetHeight(1)
	ti.CharLimit = 0                          // unlimited
	ti.KeyMap.InsertNewline.SetEnabled(false) // Enter sends
	_ = ti.Focus()                            // set focused state before program starts

	vp := viewport.New()
	vp.SoftWrap = true

	width := 80
	if w, _, err := term.GetSize(os.Stdout.Fd()); err == nil && w > 0 {
		width = w
	}

	gr, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width-4),
	)

	m := &chatModel{
		viewport:      vp,
		input:         ti,
		client:        client,
		opts:          opts,
		searchEnabled: searchEnabled,
		glamRenderer:  gr,
		width:         width,
	}

	// Add welcome banner as first system message
	m.messages = append(m.messages, chatMessage{
		role:    "system",
		content: m.welcomeText(),
	})

	return m
}

func (m *chatModel) welcomeText() string {
	var b strings.Builder
	b.WriteString(theme.Title.Render(" Z.AI Chat "))
	b.WriteString("\n\n")
	if m.opts.FilePath != "" {
		b.WriteString(theme.Info.Render("  File: ") + theme.Dim.Render(m.opts.FilePath) + "\n")
	}
	if m.searchEnabled {
		b.WriteString(theme.Info.Render("  Search: ") + theme.Dim.Render("enabled") + "\n")
	}
	b.WriteString("\n")
	b.WriteString(theme.HelpText.Render("  Commands: help, history, clear, search <query>, exit"))
	b.WriteString("\n")
	b.WriteString(theme.Divider.Render(strings.Repeat("─", 50)))
	return b.String()
}

func (m *chatModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:gocyclo,gocognit // bubbletea Update handles many msg types
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		inputH := m.input.Height() + 1
		m.viewport.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - inputH)
		m.input.SetWidth(msg.Width)
		if !m.ready {
			m.ready = true
			m.rebuildViewport()
		}
		return m, nil

	case tea.KeyPressMsg:
		return m.handleKey(msg)

	case tokenMsg:
		m.streamBuf.WriteString(string(msg))
		m.rebuildViewport()
		m.viewport.GotoBottom()
		return m, waitForToken(m.sc)

	case streamDoneMsg:
		return m.finishStream(msg)

	case streamErrMsg:
		m.streaming = false
		m.appendSystemMessage(theme.ErrorText.Render("Error: ") + msg.err.Error())
		m.rebuildViewport()
		m.viewport.GotoBottom()
		return m, m.input.Focus()

	case searchResultMsg:
		return m.handleSearchResult(msg)

	case webContentMsg:
		return m.handleWebContent(msg)
	}

	// Forward to child bubbles
	if m.streaming {
		m.viewport, _ = m.viewport.Update(msg)
	} else {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *chatModel) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "ctrl+c":
		if m.streaming {
			// Cancel the active stream
			if m.cancelFn != nil {
				m.cancelFn()
			}
			m.streaming = false
			// Finalize whatever we got
			raw := m.streamBuf.String()
			if raw != "" {
				m.appendAssistantMessage(raw)
			}
			m.streamBuf.Reset()
			m.rebuildViewport()
			m.viewport.GotoBottom()
			return m, m.input.Focus()
		}
		return m, tea.Quit

	case "enter":
		if m.streaming {
			return m, nil
		}
		input := strings.TrimSpace(m.input.Value())
		if input == "" {
			return m, nil
		}
		m.input.Reset()
		return m.handleUserInput(input)

	case "pgup", "pgdown":
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	// Forward to input when not streaming
	if !m.streaming {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	// Forward to viewport when streaming (for scrolling)
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *chatModel) handleUserInput(input string) (tea.Model, tea.Cmd) {
	lower := strings.ToLower(input)

	// Exit
	switch lower {
	case "exit", "quit", "/exit", "/quit":
		return m, tea.Quit

	case "help", "/help", "?":
		m.appendSystemMessage(m.helpText())
		m.rebuildViewport()
		m.viewport.GotoBottom()
		return m, nil

	case "history", "/history":
		m.appendSystemMessage(m.historyText())
		m.rebuildViewport()
		m.viewport.GotoBottom()
		return m, nil

	case "clear", "/clear":
		m.chatCtx = nil
		m.sessionHist = nil
		m.messages = m.messages[:0]
		m.messages = append(m.messages, chatMessage{role: "system", content: m.welcomeText()})
		m.rebuildViewport()
		m.viewport.GotoBottom()
		return m, nil

	case "context", "/context":
		m.appendSystemMessage(m.contextText())
		m.rebuildViewport()
		m.viewport.GotoBottom()
		return m, nil
	}

	// Search command
	if isSearchCommand(input) {
		query := strings.TrimSpace(stripCommandPrefix(input, "search"))
		query, searchOpts := parseSearchCommand(query)
		m.appendUserMessage(input)
		m.sessionHist = append(m.sessionHist, input)
		m.appendSystemMessage(theme.Info.Render("Searching: ") + query)
		m.rebuildViewport()
		m.viewport.GotoBottom()
		return m, doSearch(context.Background(), m.client, query, searchOpts)
	}

	// Web command
	if isWebCommand(input) {
		url := strings.TrimSpace(stripCommandPrefix(input, "web"))
		if url == "" {
			m.appendSystemMessage(theme.ErrorText.Render("Usage: /web <url>"))
			m.rebuildViewport()
			m.viewport.GotoBottom()
			return m, nil
		}
		m.appendUserMessage(input)
		m.sessionHist = append(m.sessionHist, input)
		m.appendSystemMessage(theme.Info.Render("Fetching: ") + url)
		m.rebuildViewport()
		m.viewport.GotoBottom()
		return m, doWebFetch(context.Background(), m.client, url)
	}

	// Regular chat message
	m.appendUserMessage(input)
	m.sessionHist = append(m.sessionHist, input)

	// Build opts with context
	opts := m.opts
	opts.Context = m.chatCtx
	if len(m.chatCtx) > 0 {
		opts.FilePath = ""
	}

	m.streaming = true
	m.streamBuf.Reset()
	m.input.Blur()

	// Append a placeholder assistant message (will be rendered from streamBuf)
	m.rebuildViewport()
	m.viewport.GotoBottom()

	m.streamCtx, m.cancelFn = context.WithCancel(context.Background())
	m.sc = newStreamChans()

	if m.searchEnabled {
		return m, tea.Batch(
			launchSearchAugmentedStream(m.streamCtx, m.client, input, opts, m.sc),
			waitForToken(m.sc),
		)
	}

	return m, tea.Batch(
		launchStream(m.streamCtx, m.client, input, opts, m.sc),
		waitForToken(m.sc),
	)
}

func (m *chatModel) finishStream(msg streamDoneMsg) (tea.Model, tea.Cmd) {
	m.streaming = false
	raw := m.streamBuf.String()
	m.streamBuf.Reset()

	if raw != "" {
		m.appendAssistantMessage(raw)

		// Save to history and update context
		m.client.SaveStreamToHistory(msg.input, raw, msg.usage)
		m.chatCtx = append(m.chatCtx,
			app.Message{Role: "user", Content: msg.input},
			app.Message{Role: "assistant", Content: raw},
		)
		m.trimContext()
	}

	m.rebuildViewport()
	m.viewport.GotoBottom()
	return m, m.input.Focus()
}

func (m *chatModel) handleSearchResult(msg searchResultMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.appendSystemMessage(theme.ErrorText.Render("Search error: ") + msg.err.Error())
		m.rebuildViewport()
		m.viewport.GotoBottom()
		return m, nil
	}

	var b strings.Builder
	b.WriteString(theme.Dim.Render(fmt.Sprintf("Found %d results", len(msg.results.SearchResult))))
	b.WriteString("\n\n")
	for i, result := range msg.results.SearchResult {
		b.WriteString(fmt.Sprintf("  %s %s\n",
			theme.Dim.Render(fmt.Sprintf("%d.", i+1)),
			theme.ResultTitle.Render(result.Title)))
		b.WriteString(fmt.Sprintf("     %s\n", theme.ResultLink.Render(result.Link)))
		if result.Content != "" {
			b.WriteString(fmt.Sprintf("     %s\n", theme.Dim.Render(truncate(result.Content, 200))))
		}
		b.WriteString("\n")
	}
	m.appendSystemMessage(b.String())

	// Add to conversation context
	formatted := app.FormatSearchResultsForChat(msg.results.SearchResult, msg.query)
	m.chatCtx = append(m.chatCtx,
		app.Message{Role: "user", Content: fmt.Sprintf("Search: %s", msg.query)},
		app.Message{Role: "assistant", Content: formatted},
	)
	m.trimContext()

	m.rebuildViewport()
	m.viewport.GotoBottom()
	return m, nil
}

func (m *chatModel) handleWebContent(msg webContentMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.appendSystemMessage(theme.ErrorText.Render("Fetch error: ") + msg.err.Error())
		m.rebuildViewport()
		m.viewport.GotoBottom()
		return m, nil
	}

	var b strings.Builder
	b.WriteString(theme.Section.Render(msg.resp.ReaderResult.Title))
	b.WriteString("\n")
	b.WriteString(theme.ResultLink.Render(msg.resp.ReaderResult.URL))
	b.WriteString("\n\n")
	content := msg.resp.ReaderResult.Content
	if len(content) > 2000 {
		content = content[:2000] + "\n\n" + theme.Dim.Render("[truncated — full content in context]")
	}
	b.WriteString(theme.Dim.Render(content))
	m.appendSystemMessage(b.String())

	// Add to conversation context
	formatted := app.FormatWebContent(msg.url, msg.resp.ReaderResult.Title, msg.resp.ReaderResult.Content)
	m.chatCtx = append(m.chatCtx,
		app.Message{Role: "user", Content: fmt.Sprintf("Fetched web page: %s", msg.url)},
		app.Message{Role: "assistant", Content: formatted},
	)
	m.trimContext()

	m.rebuildViewport()
	m.viewport.GotoBottom()
	return m, nil
}

func (m *chatModel) View() tea.View {
	if !m.ready {
		return tea.NewView(theme.Dim.Render("Initializing..."))
	}

	vpView := m.viewport.View()
	content := vpView + "\n" + m.input.View()

	v := tea.NewView(content)
	v.AltScreen = true
	if !m.streaming {
		c := m.input.Cursor()
		if c != nil {
			c.Position.Y += lipgloss.Height(vpView) + 1
		}
		v.Cursor = c
	}
	return v
}

// rebuildViewport re-renders all messages into the viewport content.
func (m *chatModel) rebuildViewport() {
	var b strings.Builder
	for _, msg := range m.messages {
		switch msg.role {
		case "user":
			b.WriteString(theme.Prompt.Render("you> ") + msg.content + "\n\n")
		case "assistant":
			b.WriteString(theme.AILabel.Render("AI>") + "\n")
			b.WriteString(msg.rendered + "\n\n")
		case "system":
			b.WriteString(msg.content + "\n\n")
		}
	}

	// Streaming in-progress content
	if m.streaming {
		raw := m.streamBuf.String()
		b.WriteString(theme.AILabel.Render("AI>") + "\n")
		if raw != "" {
			b.WriteString(m.renderStreaming(raw))
		} else {
			b.WriteString(theme.Dim.Render("Thinking..."))
		}
		b.WriteString("\n")
	}

	m.viewport.SetContent(b.String())
}

// renderStreaming renders completed blocks with glamour, trailing block raw.
func (m *chatModel) renderStreaming(raw string) string {
	if m.glamRenderer == nil {
		return raw
	}

	lastBound := strings.LastIndex(raw, "\n\n")
	if lastBound < 0 {
		return raw // all in-progress, show raw
	}

	completed := raw[:lastBound+2]
	tail := raw[lastBound+2:]

	rendered, err := m.glamRenderer.Render(completed)
	if err != nil {
		return raw
	}

	return strings.TrimLeft(rendered, "\n") + tail
}

// renderMarkdown renders full markdown content with glamour.
func (m *chatModel) renderMarkdown(content string) string {
	if m.glamRenderer == nil {
		return content
	}
	rendered, err := m.glamRenderer.Render(content)
	if err != nil {
		return content
	}
	return strings.Trim(rendered, "\n")
}

func (m *chatModel) appendUserMessage(content string) {
	m.messages = append(m.messages, chatMessage{role: "user", content: content})
}

func (m *chatModel) appendAssistantMessage(content string) {
	m.messages = append(m.messages, chatMessage{
		role:     "assistant",
		content:  content,
		rendered: m.renderMarkdown(content),
	})
}

func (m *chatModel) appendSystemMessage(content string) {
	m.messages = append(m.messages, chatMessage{role: "system", content: content})
}

func (m *chatModel) trimContext() {
	if len(m.chatCtx) > 20 {
		m.chatCtx = m.chatCtx[2:]
	}
}

func (m *chatModel) helpText() string {
	var b strings.Builder
	b.WriteString(theme.Section.Render("Commands") + "\n")
	b.WriteString(theme.Divider.Render(strings.Repeat("─", 40)) + "\n")
	cmds := []struct{ cmd, desc string }{
		{"help, ?", "Show this help"},
		{"history", "Show session history"},
		{"context", "Show conversation context"},
		{"clear", "Clear conversation and screen"},
		{"search <query>", "Search the web"},
		{"web <url>", "Fetch and display web page"},
		{"exit, quit", "Exit chat"},
	}
	for _, c := range cmds {
		b.WriteString(fmt.Sprintf("  %s  %s\n",
			theme.Info.Render(fmt.Sprintf("%-16s", c.cmd)),
			theme.Dim.Render(c.desc)))
	}
	return b.String()
}

func (m *chatModel) historyText() string {
	if len(m.sessionHist) == 0 {
		return theme.Dim.Render("  No messages yet.")
	}
	var b strings.Builder
	b.WriteString(theme.Section.Render(fmt.Sprintf("Session History (%d messages)", len(m.sessionHist))) + "\n")
	b.WriteString(theme.Divider.Render(strings.Repeat("─", 40)) + "\n")
	for i, msg := range m.sessionHist {
		b.WriteString(fmt.Sprintf("  %s %s\n",
			theme.Dim.Render(fmt.Sprintf("%2d.", i+1)),
			truncate(msg, 60)))
	}
	return b.String()
}

func (m *chatModel) contextText() string {
	if len(m.chatCtx) == 0 {
		return theme.Dim.Render("  No context yet.")
	}
	var b strings.Builder
	b.WriteString(theme.Section.Render(fmt.Sprintf("Conversation Context (%d messages)", len(m.chatCtx))) + "\n")
	b.WriteString(theme.Divider.Render(strings.Repeat("─", 40)) + "\n")
	for _, msg := range m.chatCtx {
		var styledRole string
		if msg.Role == "user" {
			styledRole = theme.Prompt.Render("[You]")
		} else {
			styledRole = theme.AILabel.Render("[AI]")
		}
		b.WriteString(fmt.Sprintf("  %s %s\n", styledRole, theme.Dim.Render(truncate(msg.Content, 50))))
	}
	return b.String()
}
