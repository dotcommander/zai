package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/mcp"
)

type fakeMCPSession struct {
	method   string
	params   json.RawMessage
	result   json.RawMessage
	closed   bool
	deadline bool
	callErr  error
	closeErr error
}

func (s *fakeMCPSession) Call(ctx context.Context, method string, params json.RawMessage) (json.RawMessage, error) {
	s.method, s.params = method, params
	_, s.deadline = ctx.Deadline()
	return s.result, s.callErr
}
func (s *fakeMCPSession) Close() error { s.closed = true; return s.closeErr }

func TestMCPCommands(t *testing.T) {
	oldFactory, oldKey := mcpFactory, viper.Get("api.key")
	defer func() { mcpFactory = oldFactory; viper.Set("api.key", oldKey) }()
	viper.Set("api.key", "fixture-key")
	for _, product := range []string{"vision", "search", "reader", "zread"} {
		t.Run(product, func(t *testing.T) {
			fake := &fakeMCPSession{result: json.RawMessage(`{"content":[{"type":"image","data":"bytes"}],"structuredContent":{"large":9007199254740993}}`)}
			var options mcp.Options
			mcpFactory = func(ctx context.Context, opts mcp.Options) (mcpSession, error) { options = opts; return fake, nil }
			cmd := newMCPCommand()
			var output bytes.Buffer
			cmd.SetOut(&output)
			cmd.SetErr(&output)
			cmd.SetIn(strings.NewReader(`{"nested":{"large":9007199254740993}}`))
			cmd.SetArgs([]string{product, "call", "tool-from-server", "--request-json", "-"})
			if err := cmd.ExecuteContext(context.Background()); err != nil {
				t.Fatal(err)
			}
			if fake.method != "tools/call" || !strings.Contains(string(fake.params), `9007199254740993`) || !strings.Contains(string(fake.params), `"name":"tool-from-server"`) {
				t.Fatalf("request: %s %s", fake.method, fake.params)
			}
			if !fake.closed || !fake.deadline || !strings.Contains(output.String(), `9007199254740993`) {
				t.Fatalf("result/lifecycle: %s %#v", output.String(), fake)
			}
			if options.APIKey != "fixture-key" {
				t.Fatal("key missing")
			}
			if product == "vision" {
				if options.Command != "npx" || len(options.Args) != 2 || options.Args[1] != "@z_ai/mcp-server" || !strings.Contains(strings.Join(options.Env, "\n"), "Z_AI_API_KEY=fixture-key") {
					t.Fatalf("vision options: command=%s args=%v", options.Command, options.Args)
				}
			} else if options.URL != mcpEndpoints[product] || options.HTTPClient == nil {
				t.Fatal("remote options missing")
			}
		})
	}
}

func TestMCPValidationAndPagination(t *testing.T) {
	oldFactory, oldKey := mcpFactory, viper.Get("api.key")
	defer func() { mcpFactory = oldFactory; viper.Set("api.key", oldKey) }()
	viper.Set("api.key", "fixture-key")
	calls := 0
	fake := &fakeMCPSession{result: json.RawMessage(`{"tools":[],"nextCursor":"next"}`)}
	mcpFactory = func(ctx context.Context, opts mcp.Options) (mcpSession, error) { calls++; return fake, nil }
	for _, args := range [][]string{{"search", "call", "x", "--request-json", "-"}, {"search", "tools", "--timeout", "0s"}, {"search", "tools", "--command", "x"}, {"vision", "tools", "--server-url", "https://example.test"}, {"vision", "tools", "--command-arg", "x"}} {
		cmd := newMCPCommand()
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})
		cmd.SetIn(strings.NewReader(`[]`))
		cmd.SetArgs(args)
		if err := cmd.ExecuteContext(context.Background()); err == nil {
			t.Errorf("accepted %v", args)
		}
	}
	if calls != 0 {
		t.Fatal("invalid requests started a server")
	}
	cmd := newMCPCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetArgs([]string{"reader", "tools", "--cursor", "opaque", "--server-url", "http://127.0.0.1:1234/mcp"})
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if fake.method != "tools/list" || string(fake.params) != `{"cursor":"opaque"}` || !strings.Contains(output.String(), `"nextCursor":"next"`) {
		t.Fatalf("pagination: %s %s", fake.params, output.String())
	}
}

func TestMCPToolErrorAndCancellation(t *testing.T) {
	oldFactory, oldKey := mcpFactory, viper.Get("api.key")
	defer func() { mcpFactory = oldFactory; viper.Set("api.key", oldKey) }()
	viper.Set("api.key", "fixture-key")
	fake := &fakeMCPSession{result: json.RawMessage(`{"isError":true,"content":[{"type":"text","text":"tool failed"}]}`)}
	mcpFactory = func(ctx context.Context, opts mcp.Options) (mcpSession, error) { return fake, nil }
	cmd := newMCPCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"search", "call", "x"})
	if err := cmd.ExecuteContext(context.Background()); err == nil || !strings.Contains(output.String(), `"isError":true`) {
		t.Fatalf("lost tool error: %v %s", err, output.String())
	}
	mcpFactory = func(ctx context.Context, opts mcp.Options) (mcpSession, error) { <-ctx.Done(); return nil, ctx.Err() }
	cmd = newMCPCommand()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"search", "tools", "--timeout", "1ms"})
	start := time.Now()
	if err := cmd.ExecuteContext(context.Background()); !errors.Is(err, context.DeadlineExceeded) || time.Since(start) > time.Second {
		t.Fatalf("deadline: %v", err)
	}
}

func TestMCPCloseErrorsPreserveCallFailure(t *testing.T) {
	oldFactory, oldKey := mcpFactory, viper.Get("api.key")
	defer func() { mcpFactory = oldFactory; viper.Set("api.key", oldKey) }()
	viper.Set("api.key", "fixture-key")
	closeErr := errors.New("session cleanup failed")
	callErr := errors.New("request failed")
	for _, requestErr := range []error{nil, callErr} {
		fake := &fakeMCPSession{result: json.RawMessage(`{"tools":[]}`), callErr: requestErr, closeErr: closeErr}
		mcpFactory = func(context.Context, mcp.Options) (mcpSession, error) { return fake, nil }
		cmd := newMCPCommand()
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})
		cmd.SetArgs([]string{"search", "tools"})
		err := cmd.ExecuteContext(context.Background())
		if !errors.Is(err, closeErr) || (requestErr != nil && !errors.Is(err, requestErr)) {
			t.Fatalf("lost call or cleanup error: %v", err)
		}
	}
}
