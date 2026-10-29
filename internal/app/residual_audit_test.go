package app

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dotcommander/zai/internal/config"
)

func TestResidualRetryAndBackoff(t *testing.T) {
	for _, err := range []error{context.Canceled, context.DeadlineExceeded, &APIError{StatusCode: 400, Body: "timeout 429 503"}, errors.New("customer text 503")} {
		if isRetryableError(err) {
			t.Fatalf("spurious retry: %v", err)
		}
	}
	for _, attempt := range []int{-1, 0, 35, 63, 1000000} {
		d := calculateBackoff(attempt, time.Second, 30*time.Second)
		if d <= 0 || d > 34*time.Second {
			t.Fatalf("attempt %d: %v", attempt, d)
		}
	}
}

func TestResidualCacheOptionIdentity(t *testing.T) {
	base := generateCacheKey("a", SearchOptions{SearchEngine: "standard", ContentSize: "medium"})
	for _, options := range []SearchOptions{{SearchEngine: "pro", ContentSize: "medium"}, {SearchEngine: "standard", ContentSize: "high"}} {
		if generateCacheKey("a", options) == base {
			t.Fatal("cache option collision")
		}
	}
	if generateCacheKey("adomain:b", SearchOptions{}) == generateCacheKey("a", SearchOptions{DomainFilter: "b"}) {
		t.Fatal("delimiter collision")
	}
}

func TestResidualChatBreaker(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++; w.WriteHeader(503) }))
	defer server.Close()
	client := NewClient(ClientConfig{APIKey: "test", BaseURL: server.URL}, DiscardLogger(), nil, nil)
	client.circuitBreakers["chat"] = resilienceBreaker()
	_, _, _ = client.doRequest(context.Background(), nil, ChatOptions{})
	_, _, err := client.doRequest(context.Background(), nil, ChatOptions{})
	if err == nil || calls != 1 {
		t.Fatalf("breaker bypass: calls=%d err=%v", calls, err)
	}
}

func TestResidualActiveModelHistory(t *testing.T) {
	history := NewFileHistoryStore(filepath.Join(t.TempDir(), "history.jsonl"))
	client := NewClient(ClientConfig{Model: "configured"}, DiscardLogger(), history, nil)
	client.SaveStreamToHistory("prompt", "answer", Usage{}, "active")
	entries, err := history.GetRecent(1)
	if err != nil || len(entries) != 1 || entries[0].Model != "active" {
		t.Fatalf("history=%#v err=%v", entries, err)
	}
}

func TestResidualCacheConcurrentCleanup(t *testing.T) {
	cache := NewFileSearchCache(t.TempDir())
	var tasks sync.WaitGroup
	for i := 0; i < 3; i++ {
		tasks.Go(func() {
			for j := 0; j < 10; j++ {
				if err := cache.Set("query", SearchOptions{}, []SearchResult{{Title: "complete"}}, time.Hour); err != nil {
					t.Error(err)
				}
				if _, ok := cache.Get("query", SearchOptions{}); !ok {
					t.Error("lost live entry")
				}
				if err := cache.Cleanup(); err != nil {
					t.Error(err)
				}
			}
		})
	}
	tasks.Wait()
}

func TestResidualConvenienceProviderErrors(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("x-request-id", "provider-request")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = io.WriteString(w, strings.Repeat("x", maxDiagnostic+100))
	}))
	defer server.Close()
	audio := filepath.Join(t.TempDir(), "sample.wav")
	if err := os.WriteFile(audio, []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	client := NewClient(ClientConfig{APIKey: "test", BaseURL: server.URL, RetryConfig: RetryConfig{MaxAttempts: 3}}, DiscardLogger(), nil, nil)
	cases := []struct {
		name    string
		request func() error
	}{
		{"JSON", func() error {
			_, err := client.executeJSONRequestInternal(context.Background(), "images/generations", map[string]string{})
			return err
		}},
		{"GET", func() error { _, err := client.executeGetRequestInternal(context.Background(), "models"); return err }},
		{"chat", func() error { _, err := client.Chat(context.Background(), "prompt", ChatOptions{}); return err }},
		{"chat-compat", func() error {
			_, _, err := client.doRequestWithRetry(context.Background(), nil, ChatOptions{})
			return err
		}},
		{"chat-stream", func() error { _, err := client.StreamChat(context.Background(), "prompt", ChatOptions{}); return err }},
		{"audio", func() error {
			_, err := client.TranscribeAudio(context.Background(), audio, TranscriptionOptions{})
			return err
		}},
		{"audio-stream", func() error {
			_, err := client.TranscribeAudioStream(context.Background(), audio, TranscriptionOptions{}, nil)
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Isolate each request from the deliberately repeated provider failures.
			client.circuitBreakers["chat"] = NewCircuitBreaker("chat", config.CircuitBreakerConfig{}, DiscardLogger())
			before := calls
			err := tc.request()
			var apiErr *APIError
			if !errors.As(err, &apiErr) || apiErr.StatusCode != 503 || apiErr.RequestID != "provider-request" {
				t.Fatalf("provider error=%#v (%v)", apiErr, err)
			}
			if len(apiErr.Body) > maxDiagnostic+30 || !strings.HasSuffix(apiErr.Body, " [truncated]") {
				t.Fatalf("unbounded diagnostic length %d", len(apiErr.Body))
			}
			if calls != before+1 {
				t.Fatalf("generation replayed: %d requests", calls-before)
			}
		})
	}
}
