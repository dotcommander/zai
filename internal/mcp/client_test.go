package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func writeFixture(t *testing.T, w io.Writer, text string) {
	t.Helper()
	if _, err := io.WriteString(w, text); err != nil {
		t.Errorf("write fixture: %v", err)
	}
}

func closeClient(t *testing.T, c *Client) {
	t.Helper()
	if err := c.Close(); err != nil {
		t.Errorf("close client: %v", err)
	}
}

func rpcReply(t *testing.T, w http.ResponseWriter, id json.RawMessage, result string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	writeFixture(t, w, fmt.Sprintf(`{"jsonrpc":"2.0","id":%s,"result":%s}`, id, result))
}
func TestHTTPJSONSession(t *testing.T) {
	var initialized atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer fixture-secret" {
			t.Error("missing auth")
		}
		if r.Header.Get("Accept") != "application/json, text/event-stream" {
			t.Error("missing accept")
		}
		var req struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
			Params json.RawMessage `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
			return
		}
		if req.Method != "initialize" && (r.Header.Get("Mcp-Session-Id") != "fixture-session" || r.Header.Get("MCP-Protocol-Version") != "2025-06-18") {
			t.Error("missing negotiated session headers")
		}
		switch req.Method {
		case "initialize":
			w.Header().Set("Mcp-Session-Id", "fixture-session")
			rpcReply(t, w, req.ID, `{"protocolVersion":"2025-06-18","capabilities":{"tools":{}}}`)
		case "notifications/initialized":
			initialized.Store(true)
			w.WriteHeader(http.StatusAccepted)
		case "tools/list":
			if !initialized.Load() {
				t.Error("not initialized")
			}
			rpcReply(t, w, req.ID, `{"tools":[{"name":"arbitrary","inputSchema":{"type":"object"}}],"nextCursor":"opaque+/="}`)
		case "tools/call":
			if !strings.Contains(string(req.Params), "9007199254740993") {
				t.Error("argument number lost")
			}
			rpcReply(t, w, req.ID, `{"content":[{"type":"text","text":"ok"}],"large":9007199254740993}`)
		default:
			t.Error("unexpected method", req.Method)
		}
	}))
	defer server.Close()
	client, err := New(context.Background(), Options{URL: server.URL, APIKey: "fixture-secret"})
	if err != nil {
		t.Fatal(err)
	}
	defer closeClient(t, client)
	result, err := client.Call(context.Background(), "tools/list", json.RawMessage(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(result), `opaque+/=`) {
		t.Fatal(string(result))
	}
	result, err = client.Call(context.Background(), "tools/call", json.RawMessage(`{"name":"arbitrary","arguments":{"n":9007199254740993}}`))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(result), "9007199254740993") {
		t.Fatal(string(result))
	}
}

func TestHTTPSSEAndErrors(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
			return
		}
		switch req.Method {
		case "initialize":
			rpcReply(t, w, req.ID, `{"protocolVersion":"2025-03-26"}`)
		case "notifications/initialized":
			w.WriteHeader(202)
		case "tools/list":
			w.Header().Set("Content-Type", "text/event-stream")
			writeFixture(t, w, ": keepalive\n\ndata: {\"jsonrpc\":\"2.0\",\"method\":\"notifications/progress\",\"params\":{}}\n\n")
			writeFixture(t, w, fmt.Sprintf("event: message\ndata: {\"jsonrpc\":\"2.0\",\"id\":%s,\ndata: \"result\":{\"tools\":[{\"name\":\"工具\"}]}}\n\n", req.ID))
		case "tools/call":
			calls.Add(1)
			w.Header().Set("Content-Type", "application/json")
			writeFixture(t, w, fmt.Sprintf(`{"jsonrpc":"2.0","id":%s,"error":{"code":-32000,"message":"fixture-secret rejected","extra":9007199254740993}}`, req.ID))
		}
	}))
	defer server.Close()
	c, err := New(context.Background(), Options{URL: server.URL, APIKey: "fixture-secret"})
	if err != nil {
		t.Fatal(err)
	}
	defer closeClient(t, c)
	result, err := c.Call(context.Background(), "tools/list", nil)
	if err != nil || !strings.Contains(string(result), "工具") {
		t.Fatalf("%s %v", result, err)
	}
	_, err = c.Call(context.Background(), "tools/call", json.RawMessage(`{}`))
	var rpcErr *RPCError
	if !errors.As(err, &rpcErr) {
		t.Fatal(err)
	}
	if strings.Contains(err.Error(), "fixture-secret") || !strings.Contains(string(rpcErr.Raw), "9007199254740993") {
		t.Fatal(err)
	}
	if calls.Load() != 1 {
		t.Fatal("call replayed")
	}
}

func TestHTTPFailureAndCancellation(t *testing.T) {
	for _, mode := range []string{"status", "id", "eof", "version", "cancel", "redirect"} {
		t.Run(mode, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req struct {
					ID     json.RawMessage `json:"id"`
					Method string          `json:"method"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("decode request: %v", err)
					return
				}
				if req.Method == "initialize" {
					version := "2025-03-26"
					if mode == "version" {
						version = "unsupported"
					}
					rpcReply(t, w, req.ID, `{"protocolVersion":"`+version+`"}`)
					return
				}
				if req.Method == "notifications/initialized" {
					w.WriteHeader(202)
					return
				}
				switch mode {
				case "status":
					w.WriteHeader(429)
					writeFixture(t, w, strings.Repeat("x", 70000)+"fixture-secret")
				case "id":
					rpcReply(t, w, json.RawMessage(`999`), `{}`)
				case "eof":
					w.Header().Set("Content-Type", "text/event-stream")
					writeFixture(t, w, "data: {\"jsonrpc\":\"2.0\",\"method\":\"notifications/progress\"}\n\n")
				case "cancel":
					<-r.Context().Done()
				case "redirect":
					w.Header().Set("Location", "http://example.invalid/")
					w.WriteHeader(302)
				}
			}))
			defer server.Close()
			c, err := New(context.Background(), Options{URL: server.URL, APIKey: "fixture-secret"})
			if mode == "version" {
				if err == nil {
					closeClient(t, c)
					t.Fatal("accepted unknown version")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			defer closeClient(t, c)
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()
			_, err = c.Call(ctx, "tools/list", nil)
			if err == nil {
				t.Fatal("expected error")
			}
			if mode == "cancel" && !errors.Is(err, context.DeadlineExceeded) {
				t.Fatal(err)
			}
			if mode == "status" {
				var httpErr *HTTPError
				if !errors.As(err, &httpErr) || httpErr.StatusCode != 429 || len(httpErr.Body) > 65600 {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestStdioHelperProcess(t *testing.T) {
	if os.Getenv("ZAI_MCP_TEST_HELPER") != "1" {
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var req struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			t.Fatal(err)
		}
		if len(req.ID) == 0 {
			continue
		}
		switch req.Method {
		case "initialize":
			fmt.Printf("{\"jsonrpc\":\"2.0\",\"id\":%s,\"result\":{\"protocolVersion\":\"2024-11-05\"}}\n", req.ID)
		case "tools/list":
			fmt.Printf("{\"jsonrpc\":\"2.0\",\"id\":%s,\"result\":{\"tools\":[{\"name\":\"vision\"}]}}\n", req.ID)
		case "tools/call":
			if os.Getenv("ZAI_MCP_TEST_BLOCK") == "1" {
				time.Sleep(time.Hour)
			}
			fmt.Printf("{\"jsonrpc\":\"2.0\",\"id\":%s,\"result\":{\"content\":[{\"type\":\"text\",\"text\":\"ok\"}]}}\n", req.ID)
		}
	}
	os.Exit(0)
}
func TestStdioLifecycle(t *testing.T) {
	for _, block := range []bool{false, true} {
		t.Run(fmt.Sprint(block), func(t *testing.T) {
			env := []string{"ZAI_MCP_TEST_HELPER=1"}
			if block {
				env = append(env, "ZAI_MCP_TEST_BLOCK=1")
			}
			c, err := New(context.Background(), Options{Command: os.Args[0], Args: []string{"-test.run=^TestStdioHelperProcess$"}, Env: env})
			if err != nil {
				t.Fatal(err)
			}
			result, err := c.Call(context.Background(), "tools/list", nil)
			if err != nil || !strings.Contains(string(result), "vision") {
				t.Fatalf("%s %v", result, err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			_, err = c.Call(ctx, "tools/call", json.RawMessage(`{"name":"vision","arguments":{}}`))
			if block && !errors.Is(err, context.DeadlineExceeded) {
				t.Fatal(err)
			}
			if !block && err != nil {
				t.Fatal(err)
			}
			if err := c.Close(); err != nil {
				t.Fatal(err)
			}
			if c.command.ProcessState == nil {
				t.Fatal("process not reaped")
			}
			if err := c.Close(); err != nil {
				t.Fatal(err)
			}
			if _, err := c.Call(context.Background(), "tools/list", nil); err == nil {
				t.Fatal("closed client accepted call")
			}
		})
	}
}

func TestInputValidation(t *testing.T) {
	for _, opts := range []Options{{}, {URL: "https://example.test/?key=x"}, {URL: "https://u:p@example.test"}, {URL: "file:///tmp/server"}, {URL: "https://example.test", Command: "server"}, {Command: "/nonexistent/zai-mcp-test"}} { // #nosec G101 -- Invalid URL credentials are a validation fixture, not usable secrets.
		if c, err := New(context.Background(), opts); err == nil {
			closeClient(t, c)
			t.Fatal("accepted invalid options")
		}
	}
}
