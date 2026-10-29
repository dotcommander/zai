package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAudioStreamFrames(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, body, want, err string }{
		{"deltas and final", "data: {\"type\":\"transcript.text.delta\",\"delta\":\"hello \"}\n\ndata: {\"type\":\"transcript.text.delta\",\"delta\":\"world\"}\n\ndata: {\"type\":\"transcript.text.done\",\"delta\":\"hello world\",\"model\":\"asr\"}\n\n", "hello world", ""},
		{"done marker", "data: {\"type\":\"transcript.text.delta\",\"delta\":\"hi\"}\n\ndata: [DONE]\n\n", "hi", ""},
		{"event field", "event: transcript.text.delta\ndata: {\"delta\":\"hi\"}\n\nevent: transcript.text.done\ndata: {}\n\n", "hi", ""},
		{"malformed", "data: invalid\n\n", "", "malformed"},
		{"truncated", "data: {\"type\":\"transcript.text.delta\",\"delta\":\"hi\"}\n\n", "", "before completion"},
		{"empty", "", "", "before completion"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var printed strings.Builder
			resp, err := readAudioStream(t.Context(), strings.NewReader(tc.body), func(s string) { printed.WriteString(s) })
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("error = %v, want %s", err, tc.err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if resp.Text != tc.want || printed.String() != tc.want {
				t.Fatalf("result=%q printed=%q want=%q", resp.Text, printed.String(), tc.want)
			}
		})
	}
}

func TestAudioStreamIncrementalAndCancellation(t *testing.T) {
	t.Parallel()
	for _, cancelStream := range []bool{false, true} {
		t.Run(fmt.Sprint(cancelStream), func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			t.Cleanup(cancel)
			deltaSeen := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseMultipartForm(1024); err != nil { //nolint:gosec // G120: this local test server receives only the fixture request constructed by this test.
					t.Error(err)
					return
				}
				if r.FormValue("stream") != "true" {
					t.Error("missing stream flag")
				}
				w.Header().Set("Content-Type", "text/event-stream")
				if _, err := fmt.Fprint(w, "data: {\"type\":\"transcript.text.delta\",\"delta\":\"hello\"}\n\n"); err != nil {
					t.Error(err)
					return
				}
				w.(http.Flusher).Flush()
				select {
				case <-deltaSeen:
				case <-r.Context().Done():
					return
				}
				if cancelStream {
					<-r.Context().Done()
					return
				}
				if _, err := fmt.Fprint(w, "data: {\"type\":\"transcript.text.done\",\"delta\":\"hello\"}\n\n"); err != nil {
					t.Error(err)
					return
				}
			}))
			t.Cleanup(server.Close)
			file := filepath.Join(t.TempDir(), "audio.wav")
			if err := os.WriteFile(file, []byte("audio"), 0600); err != nil {
				t.Fatal(err)
			}
			client := NewClientWithDeps(ClientConfig{APIKey: "test", BaseURL: server.URL}, DiscardLogger(), nil, &ClientDeps{StreamHTTPClient: server.Client()})
			resp, err := client.TranscribeAudioStream(ctx, file, TranscriptionOptions{}, func(delta string) {
				if delta != "hello" {
					t.Errorf("delta=%q", delta)
				}
				close(deltaSeen)
				if cancelStream {
					cancel()
				}
			})
			if cancelStream {
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("error=%v", err)
				}
				return
			}
			if err != nil || resp.Text != "hello" {
				t.Fatalf("response=%v error=%v", resp, err)
			}
		})
	}
}

func TestAudioStreamHTTPError(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "denied", http.StatusUnauthorized) }))
	t.Cleanup(server.Close)
	file := filepath.Join(t.TempDir(), "audio.wav")
	if err := os.WriteFile(file, []byte("audio"), 0600); err != nil {
		t.Fatal(err)
	}
	client := NewClientWithDeps(ClientConfig{APIKey: "test", BaseURL: server.URL}, DiscardLogger(), nil, &ClientDeps{StreamHTTPClient: server.Client()})
	_, err := client.TranscribeAudioStream(t.Context(), file, TranscriptionOptions{}, nil)
	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("error=%v", err)
	}
}
