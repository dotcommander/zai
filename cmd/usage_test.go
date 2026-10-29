package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/dotcommander/zai/internal/app"
)

type stubUsage func(context.Context) ([]app.QuotaWindow, error)

func (f stubUsage) Usage(ctx context.Context) ([]app.QuotaWindow, error) { return f(ctx) }
func TestUsageOutput(t *testing.T) {
	t.Parallel()
	windows := []app.QuotaWindow{{Type: "CREDIT_LIMIT", Percentage: 0}}
	for _, jsonMode := range []bool{false, true} {
		var out bytes.Buffer
		if err := writeUsage(context.Background(), &out, jsonMode, stubUsage(func(context.Context) ([]app.QuotaWindow, error) { return windows, nil })); err != nil {
			t.Fatal(err)
		}
		if jsonMode {
			var result struct {
				Windows []app.QuotaWindow `json:"windows"`
			}
			if err := json.Unmarshal(out.Bytes(), &result); err != nil || len(result.Windows) != 1 || result.Windows[0].Percentage != 0 {
				t.Fatalf("json: %s %v", out.String(), err)
			}
		} else if !strings.Contains(out.String(), "CREDIT_LIMIT") || !strings.Contains(out.String(), "0%") {
			t.Fatalf("text %q", out.String())
		}
		out.Reset()
		err := writeUsage(context.Background(), &out, jsonMode, stubUsage(func(context.Context) ([]app.QuotaWindow, error) { return nil, errors.New("secret-key provider body") }))
		if err == nil || strings.Contains(err.Error(), "secret-key") || strings.Contains(out.String(), "secret-key") {
			t.Fatalf("error leak: %v %s", err, out.String())
		}
		if jsonMode {
			var failure struct {
				Error struct {
					Reason string `json:"reason"`
				} `json:"error"`
			}
			if json.Unmarshal(out.Bytes(), &failure) != nil || failure.Error.Reason != "request_failed" {
				t.Fatalf("failure: %s", out.String())
			}
		}
	}
}
