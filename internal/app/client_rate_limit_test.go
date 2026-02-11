package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// MockHTTPClient is a simple mock HTTP client for testing
type MockHTTPClient struct {
	callCount atomic.Int64
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	count := m.callCount.Add(1)
	fmt.Printf("MockHTTPClient call #%d at time: %s\n", count, time.Now().Format(time.RFC3339Nano))

	// Simulate a successful response
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("test response"))),
	}, nil
}

func TestRateLimitedClient(t *testing.T) {
	// Create a mock client
	mockClient := &MockHTTPClient{}

	// Create rate limiter: 2 requests per second, burst of 3
	rateLimitConfig := RateLimitConfig{
		RequestsPerSecond: 2,
		Burst:             3,
	}

	logger := DiscardLogger()

	// Create rate-limited client
	rateLimitedClient := NewRateLimitedClient(mockClient, rateLimitConfig, logger)

	// Make 5 requests concurrently
	var wg sync.WaitGroup
	var times []time.Time
	var mu sync.Mutex

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com", nil)
			resp, err := rateLimitedClient.Do(req)
			if err != nil {
				t.Errorf("Request %d failed: %v", id, err)
				return
			}
			_ = resp.Body.Close()

			mu.Lock()
			times = append(times, time.Now())
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	// All requests should complete (rate limiting doesn't fail, just waits)
	if callCount := mockClient.callCount.Load(); callCount != 5 {
		t.Errorf("Expected 5 calls, got %d", callCount)
	}

	// Check that requests were rate limited (not all at the exact same time)
	// With 2 req/sec, the last requests should be spaced out
	if len(times) >= 2 {
		durationBetweenLastTwo := times[len(times)-1].Sub(times[len(times)-2])
		if durationBetweenLastTwo < time.Millisecond*100 {
			t.Errorf("Expected rate limiting to space out requests, but last two were only %v apart", durationBetweenLastTwo)
		}
	}
}

func TestRateLimitingDisabled(t *testing.T) {
	// Create a mock client
	mockClient := &MockHTTPClient{}

	// Create rate limiter with 0 (disabled)
	rateLimitConfig := RateLimitConfig{
		RequestsPerSecond: 0,
		Burst:             0,
	}

	logger := DiscardLogger()

	// Create rate-limited client
	rateLimitedClient := NewRateLimitedClient(mockClient, rateLimitConfig, logger)

	// Make 5 requests
	var start time.Time
	for i := 0; i < 5; i++ {
		if i == 0 {
			start = time.Now()
		}
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com", nil)
		resp, err := rateLimitedClient.Do(req)
		if err != nil {
			t.Errorf("Request %d failed: %v", i, err)
			return
		}
		_ = resp.Body.Close()
	}
	elapsed := time.Since(start)

	// All requests should complete immediately
	if callCount := mockClient.callCount.Load(); callCount != 5 {
		t.Errorf("Expected 5 calls, got %d", callCount)
	}

	// When rate limiting is disabled, all requests should complete almost immediately
	if elapsed > time.Millisecond*100 {
		t.Errorf("Expected requests to complete quickly when rate limiting is disabled, took %v", elapsed)
	}
}
