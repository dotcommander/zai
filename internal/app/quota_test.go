package app

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dotcommander/zai/internal/config"
)

type quotaDoer func(*http.Request) (*http.Response, error)

func (f quotaDoer) Do(r *http.Request) (*http.Response, error) { return f(r) }
func quotaJSON(t *testing.T, v any) string {
	t.Helper()
	b, e := json.Marshal(v)
	if e != nil {
		t.Fatal(e)
	}
	return string(b)
}
func quotaSettings() config.QuotaConfig {
	return config.QuotaConfig{Timeout: time.Second, MaxResponseBytes: 262144}
}
func quotaClient(doer HTTPDoer) *QuotaClient {
	return &QuotaClient{APIKey: "private-key", BaseURL: "https://api.z.ai/api/paas/v4", CodingBaseURL: "https://api.z.ai/api/coding/paas/v4", Settings: quotaSettings(), Doer: doer}
}

func TestQuotaWindows(t *testing.T) {
	t.Parallel()
	body := quotaJSON(t, map[string]any{"data": map[string]any{"limits": []any{
		map[string]any{"type": "CREDIT_LIMIT", "percentage": 0, "currentValue": 0, "usage": 100, "remaining": 100, "nextResetTime": int64(1730000000000), "unit": 3, "number": 5},
		map[string]any{"type": "TOKENS_LIMIT", "currentValue": 2, "usage": 4},
		map[string]any{"type": "TIME_LIMIT", "percentage": 25},
		map[string]any{"type": "UNKNOWN", "percentage": 1},
	}}})
	windows, err := parseQuotaWindows([]byte(body))
	if err != nil {
		t.Fatal(err)
	}
	if len(windows) != 3 || windows[0].Percentage != 0 || windows[0].Used == nil || *windows[0].Used != 0 || windows[0].DurationMinutes == nil || *windows[0].DurationMinutes != 300 || windows[0].ResetsAt == nil || windows[1].Percentage != 50 || windows[1].Remaining == nil || *windows[1].Remaining != 2 || windows[2].Used != nil || windows[2].Remaining != nil {
		t.Fatalf("windows: %+v", windows)
	}
	for _, v := range []any{map[string]any{"data": map[string]any{"limits": []any{map[string]any{"type": "UNKNOWN", "percentage": 0}}}}, map[string]any{"data": map[string]any{"limits": []any{map[string]any{"type": "CREDIT_LIMIT", "percentage": -1}}}}} {
		if _, err := parseQuotaWindows([]byte(quotaJSON(t, v))); err == nil {
			t.Fatal("accepted invalid limits")
		}
	}
}

func TestQuotaClientSafety(t *testing.T) {
	t.Parallel()
	calls := 0
	client := quotaClient(quotaDoer(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != "GET" || r.URL.String() != "https://api.z.ai/api/monitor/usage/quota/limit" || r.Header.Get("Authorization") != "private-key" {
			t.Fatalf("unsafe request: %s %s", r.Method, r.URL)
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(quotaJSON(t, map[string]any{"data": map[string]any{"limits": []any{map[string]any{"type": "CREDIT_LIMIT", "percentage": 0}}}})))}, nil
	}))
	if _, err := client.Usage(context.Background()); err != nil || calls != 1 {
		t.Fatalf("usage: %v calls %d", err, calls)
	}
	for _, urls := range [][2]string{{"https://proxy.example/api", "https://api.z.ai/api/coding/paas/v4"}, {"https://api.z.ai:443/api", "https://api.z.ai/api"}, {"http://api.z.ai/api", "http://api.z.ai/api"}, {"https://api.z.ai@proxy.example/api", "https://api.z.ai/api"}, {"https://api.z.ai/api", "https://dev.bigmodel.cn/api"}} {
		bad := *client
		bad.BaseURL = urls[0]
		bad.CodingBaseURL = urls[1]
		if _, err := bad.Usage(context.Background()); quotaReason(err) != "unsupported_endpoint" {
			t.Fatalf("urls %v: %v", urls, err)
		}
	}
	if calls != 1 {
		t.Fatalf("unsafe call count %d", calls)
	}
	client.Doer = quotaDoer(func(*http.Request) (*http.Response, error) { return nil, errors.New("private-key secret transport") })
	if _, err := client.Usage(context.Background()); strings.Contains(err.Error(), "private-key") || quotaReason(err) != "request_failed" {
		t.Fatalf("unsafe error %v", err)
	}
	client.Doer = quotaDoer(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 302, Header: http.Header{"Location": []string{"https://evil.example/private-key"}}, Body: io.NopCloser(strings.NewReader("private-key"))}, nil
	})
	if _, err := client.Usage(context.Background()); quotaReason(err) != "unsupported_endpoint" || strings.Contains(err.Error(), "private-key") {
		t.Fatalf("redirect %v", err)
	}
}
