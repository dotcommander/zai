package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/dotcommander/zai/internal/config"
)

func resilienceBreaker() *CircuitBreaker {
	return NewCircuitBreaker("test", config.CircuitBreakerConfig{FailureThreshold: 1, SuccessThreshold: 1, Timeout: time.Hour}, DiscardLogger())
}

func awaitResilience(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("request did not progress")
	}
}

func TestCircuitBreakerClosedRequestsOverlap(t *testing.T) {
	cb := resilienceBreaker()
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	finished := make(chan struct{}, 2)
	for range 2 {
		go func() {
			_ = cb.Execute(func() error { entered <- struct{}{}; <-release; return nil })
			finished <- struct{}{}
		}()
	}
	awaitResilience(t, entered)
	awaitResilience(t, entered)
	close(release)
	awaitResilience(t, finished)
	awaitResilience(t, finished)
}

func TestCircuitBreakerHalfOpenSingleProbeAndExcludedRelease(t *testing.T) {
	cb := resilienceBreaker()
	cb.state = Open
	cb.lastStateChange = time.Now().Add(-2 * time.Hour)
	entered := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	go func() {
		_ = cb.Execute(func() error { close(entered); <-release; return &APIError{StatusCode: 429} })
		close(done)
	}()
	awaitResilience(t, entered)
	called := false
	if err := cb.Execute(func() error { called = true; return nil }); err == nil || called {
		t.Fatal("concurrent half-open probe admitted")
	}
	close(release)
	awaitResilience(t, done)
	if cb.state != HalfOpen || cb.failureCount != 0 || cb.successCount != 0 {
		t.Fatalf("excluded result changed state/counters: %+v", cb)
	}
	if err := cb.Execute(func() error { return nil }); err != nil || cb.state != Closed {
		t.Fatal("excluded result did not release probe", err)
	}
}

func TestCircuitBreakerStaleCompletions(t *testing.T) {
	for _, staleErr := range []error{nil, errors.New("old failure")} {
		t.Run(fmt.Sprint(staleErr), func(t *testing.T) {
			cb := resilienceBreaker()
			entered := make(chan struct{})
			release := make(chan struct{})
			done := make(chan struct{})
			go func() { _ = cb.Execute(func() error { close(entered); <-release; return staleErr }); close(done) }()
			awaitResilience(t, entered)
			if staleErr == nil {
				_ = cb.Execute(func() error { return errors.New("new failure") })
			} else {
				cb.Reset()
			}
			close(release)
			awaitResilience(t, done)
			want := Open
			if staleErr != nil {
				want = Closed
			}
			if cb.state != want || (want == Closed && cb.failureCount != 0) {
				t.Fatalf("stale result changed newer state: %s", cb.state)
			}
		})
	}
}

func TestCircuitBreakerErrorClassification(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		cancel   bool
		wantOpen bool
	}{
		{"server", &APIError{StatusCode: 503}, false, true},
		{"client", &APIError{StatusCode: 400}, false, false},
		{"rate", &APIError{StatusCode: 429}, false, false},
		{"transport", &url.Error{Op: "Get", URL: "https://example.test", Err: io.ErrUnexpectedEOF}, false, true},
		{"transport timeout", &url.Error{Op: "Get", URL: "https://example.test", Err: context.DeadlineExceeded}, false, true},
		{"net timeout wrapping deadline", &net.OpError{Op: "read", Net: "tcp", Err: context.DeadlineExceeded}, false, true},
		{"caller timeout", context.DeadlineExceeded, false, false},
		{"caller cancellation", context.Canceled, false, false},
		{"expired context", &url.Error{Op: "Get", Err: context.DeadlineExceeded}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := resilienceBreaker()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tt.cancel {
				cancel()
			}
			_ = cb.executeWithContext(ctx, func() error { return fmt.Errorf("wrapped: %w", tt.err) })
			if (cb.state == Open) != tt.wantOpen {
				t.Fatalf("state = %s, wantOpen=%v", cb.state, tt.wantOpen)
			}
		})
	}
}

func TestCircuitBreakerExpiredCallerDeadline(t *testing.T) {
	cb := resilienceBreaker()
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	_ = cb.executeWithContext(ctx, func() error { return &net.OpError{Op: "read", Err: context.DeadlineExceeded} })
	if cb.state != Closed || cb.failureCount != 0 {
		t.Fatal("expired caller deadline counted as transport failure")
	}
}

func TestStreamingHTTPPolicyAndCancellation(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ordinary" {
			<-r.Context().Done()
			return
		}
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		select {
		case <-release:
			_, _ = io.WriteString(w, "stream survived")
		case <-r.Context().Done():
		}
	}))
	defer server.Close()
	client := NewClient(ClientConfig{Timeout: 20 * time.Millisecond}, DiscardLogger(), nil, nil)
	streamClient := client.StreamHTTPClient().(*http.Client)
	transport := streamClient.Transport.(*http.Transport)
	if streamClient.Timeout != 0 || transport.ResponseHeaderTimeout != time.Minute || transport.TLSHandshakeTimeout == 0 || transport.DialContext == nil {
		t.Fatal("stream transport policy missing")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/stream", nil)
	response, err := streamClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := response.Body.Close(); err != nil {
			t.Error(err)
		}
	})
	ordinaryReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/ordinary", nil)
	if _, err := client.HTTPClient().Do(ordinaryReq); err == nil {
		t.Fatal("ordinary timeout did not fire")
	}
	close(release)
	body, err := io.ReadAll(response.Body)
	if err != nil || string(body) != "stream survived" {
		t.Fatalf("stream: %q, %v", body, err)
	}

	// A separate server holds the body open until cancellation.
	held := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer held.Close()
	cancelCtx, stop := context.WithCancel(context.Background())
	cancelReq, _ := http.NewRequestWithContext(cancelCtx, http.MethodGet, held.URL, nil)
	cancelResp, err := streamClient.Do(cancelReq)
	if err != nil {
		stop()
		t.Fatal(err)
	}
	stop()
	_, err = io.ReadAll(cancelResp.Body)
	_ = cancelResp.Body.Close()
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("body cancellation: %v", err)
	}
}

type resilienceDoer func(*http.Request) (*http.Response, error)

func (d resilienceDoer) Do(r *http.Request) (*http.Response, error) { return d(r) }

func TestZeroBurstRateLimitAdmitsRequest(t *testing.T) {
	called := false
	doer := NewRateLimitedClient(resilienceDoer(func(*http.Request) (*http.Response, error) {
		called = true
		return &http.Response{Body: io.NopCloser(strings.NewReader(""))}, nil
	}), RateLimitConfig{RequestsPerSecond: 1}, DiscardLogger())
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	resp, err := doer.Do(req)
	if err != nil || !called {
		t.Fatalf("zero burst blocked request: %v", err)
	}
	_ = resp.Body.Close()
}

func TestStreamingHTTPInjectionAndSharedRateLimit(t *testing.T) {
	normal := resilienceDoer(func(*http.Request) (*http.Response, error) { return nil, errors.New("normal") })
	stream := resilienceDoer(func(*http.Request) (*http.Response, error) { return nil, errors.New("stream") })
	client := NewClientWithDeps(ClientConfig{RateLimit: RateLimitConfig{RequestsPerSecond: 1, Burst: 1}}, DiscardLogger(), nil, &ClientDeps{HTTPClient: normal, StreamHTTPClient: stream})
	if client.httpClient.(*RateLimitedClient).limiter != client.streamHTTPClient.(*RateLimitedClient).limiter {
		t.Fatal("stream bypasses shared rate limit")
	}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	// Inspect the injected doer beneath rate limiting without consuming tokens.
	_, err := client.streamHTTPClient.(*RateLimitedClient).client.Do(req)
	if err == nil || err.Error() != "stream" {
		t.Fatal("stream injection ignored")
	}
}
