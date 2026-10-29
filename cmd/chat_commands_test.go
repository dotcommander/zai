package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"github.com/dotcommander/zai/internal/app"
)

func awaitStreamTerminal(t *testing.T, sc streamChans) tea.Msg {
	t.Helper()
	result := make(chan tea.Msg, 1)
	go func() {
		for {
			msg := waitForToken(sc)()
			if _, ok := msg.(tokenMsg); !ok {
				result <- msg
				return
			}
		}
	}()
	select {
	case msg := <-result:
		if len(sc.done) != 0 {
			t.Fatal("more than one terminal result")
		}
		return msg
	case <-time.After(time.Second):
		t.Fatal("stream waiter blocked")
		return nil
	}
}

func TestRepairChatStreamTerminal(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"success", "request-error", "cancel"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if mode == "request-error" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "text/event-stream")
				w.WriteHeader(http.StatusOK)
				w.(http.Flusher).Flush()
				if mode == "cancel" {
					<-r.Context().Done()
					return
				}
				if _, err := fmt.Fprintln(w, "data: [DONE]"); err != nil {
					t.Error(err)
				}
			}))
			t.Cleanup(server.Close)
			client := app.NewClient(app.ClientConfig{APIKey: "test", BaseURL: server.URL}, app.DiscardLogger(), nil, nil)
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			sc := newStreamChans()
			if msg := launchStream(ctx, client, "hello", app.ChatOptions{}, sc)(); msg != nil {
				t.Fatalf("launch returned duplicate terminal: %#v", msg)
			}
			if mode == "cancel" {
				cancel()
			}
			msg := awaitStreamTerminal(t, sc)
			if mode == "success" {
				if _, ok := msg.(streamDoneMsg); !ok {
					t.Fatalf("terminal = %#v", msg)
				}
			} else {
				failure, ok := msg.(streamErrMsg)
				if !ok {
					t.Fatalf("terminal = %#v", msg)
				}
				if mode == "cancel" && !errors.Is(failure.err, context.Canceled) {
					t.Fatalf("error = %v", failure.err)
				}
			}
		})
	}
}

func TestRepairChatCancellationPresentation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name         string
		streaming    bool
		err          error
		wantMessages int
	}{
		{name: "already canceled", err: fmt.Errorf("request: %w", context.Canceled)},
		{name: "active cancellation", streaming: true, err: context.Canceled, wantMessages: 1},
		{name: "genuine error", err: errors.New("request failed"), wantMessages: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := &chatModel{input: textarea.New(), viewport: viewport.New(), width: 80, streaming: tc.streaming}
			m.Update(streamErrMsg{err: tc.err})
			if len(m.messages) != tc.wantMessages {
				t.Fatalf("messages = %#v", m.messages)
			}
			if m.streaming {
				t.Fatal("stream remains active")
			}
		})
	}
}
