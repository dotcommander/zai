package app

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestQuotaFailures(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{"authentication", 401, "private-key", "authentication_failed"},
		{"too large", 200, strings.Repeat("x", 262145), "invalid_response"},
		{"invalid schema", 200, "private-key", "invalid_response"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := quotaClient(quotaDoer(func(r *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: tc.status, Body: io.NopCloser(strings.NewReader(tc.body))}, nil
			}))
			_, err := client.Usage(context.Background())
			if quotaReason(err) != tc.want || strings.Contains(err.Error(), "private-key") {
				t.Fatalf("failure: %v", err)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client := quotaClient(quotaDoer(func(r *http.Request) (*http.Response, error) {
		if r.Context().Err() == nil {
			t.Fatal("cancellation not propagated")
		}
		return nil, r.Context().Err()
	}))
	if _, err := client.Usage(ctx); quotaReason(err) != "request_failed" {
		t.Fatalf("cancel: %v", err)
	}
}
