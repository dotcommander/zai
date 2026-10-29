package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"

	"github.com/dotcommander/zai/internal/app"
)

func TestResidualStdinBoundary(t *testing.T) {
	previous := os.Stdin
	defer func() { os.Stdin = previous }()
	for _, size := range []int{MaxStdinSize, MaxStdinSize + 1} {
		file, err := os.CreateTemp(t.TempDir(), "stdin")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = file.WriteString(strings.Repeat("x", size)); err != nil {
			t.Fatal(err)
		}
		_, _ = file.Seek(0, 0)
		os.Stdin = file
		_, err = readStdin()
		if closeErr := file.Close(); closeErr != nil {
			t.Fatal(closeErr)
		}
		if (err != nil) != (size > MaxStdinSize) {
			t.Fatalf("size=%d err=%v", size, err)
		}
	}
}

func TestResidualMissingSystemAndUTF8(t *testing.T) {
	if _, _, err := resolveSystemPrompt(filepath.Join(t.TempDir(), "missing.txt")); err == nil {
		t.Fatal("missing path accepted")
	}
	for _, limit := range []int{0, 1, 2, 3, 4, 8} {
		if text := truncate("日本語の文章です", limit); !utf8.ValidString(text) {
			t.Fatalf("invalid UTF8: %q", text)
		}
	}
}

func TestResidualStaleStream(t *testing.T) {
	m := &chatModel{turn: 2, streaming: true, input: textarea.New(), viewport: viewport.New()}
	m.streamBuf.WriteString("current")
	m.Update(turnTokenMsg{turn: 1, token: "stale"})
	m.Update(streamDoneMsg{turn: 1})
	m.Update(streamErrMsg{turn: 1})
	if !m.streaming || m.streamBuf.String() != "current" {
		t.Fatal("stale turn changed active stream")
	}
}

func TestResidualRawCommandsAndFileContext(t *testing.T) {
	previous := os.Stdin
	defer func() { os.Stdin = previous }()
	input, err := os.CreateTemp(t.TempDir(), "input")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = input.WriteString("/help\n/clear\nfirst\nsecond\n/quit\n")
	_, _ = input.Seek(0, 0)
	t.Cleanup(func() {
		if err := input.Close(); err != nil {
			t.Error(err)
		}
	})
	os.Stdin = input
	file := filepath.Join(t.TempDir(), "context.txt")
	if err := os.WriteFile(file, []byte("retained-file-marker"), 0600); err != nil {
		t.Fatal(err)
	}
	var requests []app.ChatRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request app.ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
		}
		requests = append(requests, request)
		if _, err := fmt.Fprintln(w, `data: {"choices":[{"delta":{"content":"ok"}}]}`); err != nil {
			t.Error(err)
		}
		if _, err := fmt.Fprintln(w, "data: [DONE]"); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	client := app.NewClient(app.ClientConfig{APIKey: "test", BaseURL: server.URL}, app.DiscardLogger(), nil, nil)
	if err := runChatREPLRaw(client, app.ChatOptions{FilePath: file, Model: "override"}, false); err != nil {
		t.Fatal(err)
	}
	if len(requests) != 2 {
		t.Fatalf("slash commands made paid requests: %d", len(requests))
	}
	for _, request := range requests {
		if request.Model != "override" || !strings.Contains(request.Messages[len(request.Messages)-1].Content, "retained-file-marker") {
			t.Fatalf("lost model or file: %#v", request)
		}
	}
}
