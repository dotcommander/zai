package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/dotcommander/zai/internal/app"
)

func sampleUsageDisplay() usageDisplay {
	return usageDisplay{
		Badge: " ZAI ", Heading: "Account quota", Types: map[string]string{"TIME_LIMIT": "Time quota", "TOKENS_LIMIT": "Token quota", "CREDIT_LIMIT": "Credit quota"},
		BarWidth: 16, Filled: "━", Empty: "─", PercentMark: "%", GroupSeparator: ",", DurationHour: "h", DurationMinute: "m", Indent: "  ", Gap: "  ",
		ResetFormat: "Mon, 02 Jan 2006 15:04 UTC", Labels: map[string]string{"used": "used", "limit": "limit", "remaining": "remaining", "duration": "duration", "reset": "reset"},
	}
}

func TestUsagePrettySample(t *testing.T) {
	t.Parallel()
	used, limit, remaining := 56.0, 1000.0, 944.0
	duration := int64(300)
	resetTime := time.Date(2026, 11, 26, 19, 17, 35, 0, time.UTC)
	resetTokens := time.Date(2026, 10, 29, 5, 38, 54, 0, time.UTC)
	windows := []app.QuotaWindow{
		{Type: "TIME_LIMIT", Percentage: 5, Used: &used, Limit: &limit, Remaining: &remaining, ResetsAt: &resetTime},
		{Type: "TOKENS_LIMIT", Percentage: 1, DurationMinutes: &duration, ResetsAt: &resetTokens},
	}
	var out bytes.Buffer
	if err := writeUsagePretty(&out, windows, sampleUsageDisplay(), 80, false); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	t.Log("\n" + got)
	if groupUsageNumber(1000000, ",") != "1,000,000" {
		t.Error("large counts are not grouped")
	}
	for _, want := range []string{"ZAI", "Account quota", "Time quota", "Token quota", "5% used", "1% used", "used 56", "limit 1,000", "remaining 944", "duration 5h", "reset Thu, 26 Nov 2026 19:17 UTC", "reset Thu, 29 Oct 2026 05:38 UTC", "━"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in %q", want, got)
		}
	}
	if strings.Count(got, "━") != 2 || !strings.Contains(got, "used 56  limit 1,000  remaining 944") {
		t.Errorf("low-usage bars or spacing: %q", got)
	}
	if strings.Index(got, "Time quota") > strings.Index(got, "Token quota") || strings.Contains(got, "300m") || strings.Contains(got, "\x1b[") {
		t.Errorf("layout: %q", got)
	}
}

func TestUsagePrettyMissingZeroAndNarrow(t *testing.T) {
	t.Parallel()
	zero := 0.0
	windows := []app.QuotaWindow{{Type: "CREDIT_LIMIT", Percentage: 0, Used: &zero}, {Type: "FUTURE_LIMIT", Percentage: 42}}
	var out bytes.Buffer
	if err := writeUsagePretty(&out, windows, sampleUsageDisplay(), 12, false); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{"Credit quota", "0% used", "used 0", "FUTURE_LIMIT", "42% used"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in %q", want, got)
		}
	}
	for _, absent := range []string{"limit 0", "remaining 0", "duration 0", "reset unknown"} {
		if strings.Contains(got, absent) {
			t.Errorf("unexpected %q in %q", absent, got)
		}
	}
}

func TestUsagePrettyOutputError(t *testing.T) {
	t.Parallel()
	if err := writeUsagePretty(failingUsageWriter{}, []app.QuotaWindow{{Type: "TIME_LIMIT"}}, sampleUsageDisplay(), 80, false); err == nil {
		t.Fatal("expected writer failure")
	}
}

type failingUsageWriter struct{}

func (failingUsageWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

func TestUsageJSONAndPlainUnchanged(t *testing.T) {
	t.Parallel()
	windows := []app.QuotaWindow{{Type: "TIME_LIMIT", Percentage: 5}}
	reader := stubUsage(func(context.Context) ([]app.QuotaWindow, error) { return windows, nil })
	var plain, encoded bytes.Buffer
	if err := writeUsage(context.Background(), &plain, false, reader); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(plain.String(), "TIME_LIMIT: 5%\n") || strings.Contains(plain.String(), "\x1b[") || strings.Contains(plain.String(), "Time quota") {
		t.Fatalf("plain: %q", plain.String())
	}
	if err := writeUsage(context.Background(), &encoded, true, reader); err != nil {
		t.Fatal(err)
	}
	var result struct {
		Windows []app.QuotaWindow `json:"windows"`
	}
	if err := json.Unmarshal(encoded.Bytes(), &result); err != nil || len(result.Windows) != 1 || result.Windows[0].Type != "TIME_LIMIT" {
		t.Fatalf("json: %q: %v", encoded.String(), err)
	}
	if strings.Contains(encoded.String(), "\x1b[") {
		t.Fatalf("ANSI in JSON: %q", encoded.String())
	}
}
