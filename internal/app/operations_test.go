package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func operationTestClient(server *httptest.Server) *Client {
	return NewClientWithDeps(ClientConfig{APIKey: "test-key", BaseURL: server.URL + "/api/paas/v4"}, DiscardLogger(), nil, &ClientDeps{HTTPClient: server.Client(), StreamHTTPClient: server.Client()})
}

func TestOperationContracts(t *testing.T) {
	for _, op := range Operations() {
		t.Run(op.Name, func(t *testing.T) {
			response := `{"id":"task-1","unknown":{"number":9007199254740993,"nested":[true,null]}}`
			input := json.RawMessage(`{"model":"chosen","extension":{"integer":9007199254740993,"array":[{"video_url":{"url":"https://example.com/video"}}]},"stream":false}`)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				family := map[string]string{"paas": "/api/paas/v4/", "v1": "/api/v1/", "anthropic": "/api/anthropic/"}[op.Family]
				want := family + strings.ReplaceAll(op.Route, "{id}", "task-1")
				if r.URL.Path != want || r.Method != op.Method {
					t.Errorf("route %s %s, want %s %s", r.Method, r.URL.Path, op.Method, want)
				}
				if r.Header.Get("Authorization") != "Bearer test-key" {
					t.Error("missing bearer authentication")
				}
				if op.Auth == "anthropic" && (r.Header.Get("x-api-key") != "test-key" || r.Header.Get("anthropic-version") == "") {
					t.Error("missing Messages headers")
				}
				if op.RequestMode == "json" {
					data, _ := io.ReadAll(r.Body)
					if string(data) != string(input) {
						t.Errorf("lossy request: %s", data)
					}
				}
				if op.RequestMode == "multipart" {
					if err := r.ParseMultipartForm(1 << 20); err != nil { //nolint:gosec // G120: this local test server receives only the fixture request constructed by this test.
						t.Error(err)
						return
					}
					defer func() {
						if err := r.MultipartForm.RemoveAll(); err != nil {
							t.Error(err)
						}
					}()
					if r.FormValue("extension") != `{"integer":9007199254740993,"array":[{"video_url":{"url":"https://example.com/video"}}]}` {
						t.Errorf("multipart lost field: %s", r.FormValue("extension"))
					}
					f, header, err := r.FormFile("file")
					if err != nil {
						t.Error(err)
						return
					}
					defer func() {
						if err := f.Close(); err != nil {
							t.Error(err)
						}
					}()
					data, _ := io.ReadAll(f)
					if string(data) != "file contents" || header.Filename == "" {
						t.Error("bad upload")
					}
				}
				w.Header().Set("Content-Type", "application/json")
				if _, err := fmt.Fprint(w, response); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			req := OperationRequest{Operation: op.Name, PathParams: map[string]string{"id": "task-1"}, JSON: input}
			if op.RequestMode == "multipart" {
				ext := ".pdf"
				if op.Name == "audio" {
					ext = ".wav"
				}
				req.UploadPath = filepath.Join(t.TempDir(), "input"+ext)
				if err := os.WriteFile(req.UploadPath, []byte("file contents"), 0600); err != nil {
					t.Fatal(err)
				}
				req.Fields = map[string]string{"purpose": "user_data"}
			}
			got, err := operationTestClient(server).ExecuteOperation(context.Background(), req)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != response {
				t.Errorf("lossy response %s", got)
			}
		})
	}
}

func TestOperationInventory(t *testing.T) {
	data, err := os.ReadFile("../../docs/openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	var spec struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err = json.Unmarshal(data, &spec); err != nil {
		t.Fatal(err)
	}
	count := 0
	for path, methods := range spec.Paths {
		for method := range methods {
			if method != "post" && method != "get" {
				continue
			}
			count++
			found := false
			for _, op := range Operations() {
				prefix := map[string]string{"paas": "/paas/v4/", "v1": "/v1/"}[op.Family]
				if prefix+op.Route == path && strings.EqualFold(op.Method, method) {
					found = true
				}
			}
			if !found {
				t.Errorf("uncovered operation %s %s", method, path)
			}
		}
	}
	if count != 14 {
		t.Fatalf("inventory changed: %d operations", count)
	}
}

func TestOperationErrorsAndNoReplay(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("x-request-id", "request-42")
		w.WriteHeader(503)
		if _, err := fmt.Fprint(w, strings.Repeat("x", maxDiagnostic+100)); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	c := operationTestClient(server)
	c.config.RetryConfig.MaxAttempts = 3
	_, err := c.ExecuteOperation(context.Background(), OperationRequest{Operation: "images", JSON: json.RawMessage(`{}`)})
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 503 || apiErr.RequestID != "request-42" || len(apiErr.Body) > maxDiagnostic+20 {
		t.Fatalf("bad API error: %#v", err)
	}
	if calls.Load() != 1 {
		t.Fatal("billable request replayed")
	}
	calls.Store(0)
	_, err = c.ExecuteOperation(context.Background(), OperationRequest{Operation: "task_get", PathParams: map[string]string{"id": "1"}})
	if err == nil || calls.Load() != 3 {
		t.Fatalf("safe retrieval retry calls=%d error=%v", calls.Load(), err)
	}
}

func TestOperationRouting(t *testing.T) {
	c := NewClient(ClientConfig{BaseURL: "https://example.test/custom/paas", CodingBaseURL: "https://coding.test/api/coding/paas/v4", UseCoding: true}, DiscardLogger(), nil, nil)
	for _, tt := range []struct{ name, want string }{{"chat", "https://coding.test/api/coding/paas/v4/chat/completions"}, {"images", "https://example.test/custom/paas/images/generations"}, {"responses", "https://example.test/api/v1/responses"}, {"messages", "https://example.test/api/anthropic/v1/messages"}} {
		op, _ := lookupOperation(tt.name)
		got, err := c.operationURL(op, nil)
		if err != nil || got != tt.want {
			t.Errorf("%s: %q %v", tt.name, got, err)
		}
	}
	op, _ := lookupOperation("task_get")
	got, _ := c.operationURL(op, map[string]string{"id": "x/y?z"})
	if !strings.HasSuffix(got, "x%2Fy%3Fz") {
		t.Errorf("unescaped id: %s", got)
	}
}

// Small reads split SSE fields, UTF-8, and JSON independently of network packets.
type fragmentedReader struct{ data string }

func (r *fragmentedReader) Read(p []byte) (int, error) {
	if r.data == "" {
		return 0, io.EOF
	}
	p[0] = r.data[0]
	r.data = r.data[1:]
	return 1, nil
}
func TestOperationEvents(t *testing.T) {
	wire := ": ping\r\nevent: reasoning\r\nid: 42\r\ndata: {\"text\":\"思考\"}\r\n\r\nevent: tool\ndata: {\"tool_calls\":[\ndata: {\"id\":\"call-1\"}]}\n\ndata: {\"choices\":[{\"delta\":{\"content\":\"answer\"}}]}\n\ndata: [DONE]\n\n"
	var events []Event
	err := readOperationEvents(context.Background(), &fragmentedReader{wire}, "chat", func(e Event) error { events = append(events, e); return nil })
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 4 || events[0].Event != "reasoning" || events[0].ID != "42" || events[1].Event != "tool" || !strings.Contains(events[1].Data, "\n") || events[3].Data != "[DONE]" {
		t.Fatalf("lost event data: %#v", events)
	}
	for _, tt := range []struct{ operation, wire string }{{"audio", "event: transcript.text.done\ndata: {\"type\":\"transcript.text.done\"}\n\n"}, {"messages", "event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n"}, {"responses", "event: response.completed\ndata: {\"type\":\"response.completed\"}\n\n"}, {"agent_run", "data: {\"choices\":[{\"finish_reason\":\"stop\"}]}\n\n"}} {
		if err = readOperationEvents(context.Background(), strings.NewReader(tt.wire), tt.operation, func(Event) error { return nil }); err != nil {
			t.Errorf("%s: %v", tt.operation, err)
		}
	}
	if err = readOperationEvents(context.Background(), strings.NewReader("data: {}\n\n"), "chat", func(Event) error { return nil }); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("premature termination: %v", err)
	}
	if err = readOperationEvents(context.Background(), strings.NewReader("event: error\ndata: {\"error\":\"bad\"}\n\n"), "chat", func(Event) error { return nil }); err == nil {
		t.Error("missing SSE error")
	}
}

func TestOperationStreamCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		if _, err := fmt.Fprint(w, "data: {\"delta\":\"first\"}\n\n"); err != nil {
			t.Error(err)
		}
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := operationTestClient(server).StreamOperation(ctx, OperationRequest{Operation: "chat", JSON: json.RawMessage(`{"messages":[]}`)}, func(Event) error { cancel(); return nil })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel error: %v", err)
	}
}

func TestOperationAsync(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if n == 1 {
			if _, err := fmt.Fprint(w, `{"id":"1","task_status":"PROCESSING"}`); err != nil {
				t.Error(err)
			}
		} else {
			if _, err := fmt.Fprint(w, `{"id":"1","task_status":"SUCCESS","video_result":[{"url":"https://example.com/a"}]}`); err != nil {
				t.Error(err)
			}
		}
	}))
	defer server.Close()
	result, err := operationTestClient(server).WaitTask(context.Background(), "1", time.Millisecond)
	if err != nil || !strings.Contains(string(result), "SUCCESS") {
		t.Fatalf("poll: %s %v", result, err)
	}
	for _, status := range []string{"FAIL", "UNKNOWN"} {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := fmt.Fprintf(w, `{"task_status":%q}`, status); err != nil {
				t.Error(err)
			}
		}))
		if _, err = operationTestClient(s).WaitTask(context.Background(), "1", time.Millisecond); err == nil {
			t.Errorf("status %s accepted", status)
		}
		s.Close()
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{"task_status":"PROCESSING"}`); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if _, err = operationTestClient(s).WaitTask(ctx, "1", time.Second); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("deadline: %v", err)
	}
}

func TestOperationUploadConstraints(t *testing.T) {
	for _, tt := range []struct {
		op, purpose, path string
		size              int64
		valid             bool
	}{
		{"files_upload", "agent", "a.txt", 100 << 20, true}, {"files_upload", "agent", "a.txt", (100 << 20) + 1, false},
		{"files_upload", "user_data", "a.pdf", 101 << 20, true}, {"files_upload", "user_data", "a.pdf", 1 << 40, true}, {"files_upload", "user_data", "a.txt", 1, false},
		{"audio", "", "a.wav", 25 << 20, true}, {"audio", "", "a.mp3", (25 << 20) + 1, false},
	} {
		if err := validateOperationUpload(tt.op, tt.purpose, tt.path, tt.size); (err == nil) != tt.valid {
			t.Errorf("%+v: %v", tt, err)
		}
	}
	for _, bad := range []string{`null`, `[]`, `{} {}`, `42`} {
		if _, err := JSONObject([]byte(bad)); err == nil {
			t.Errorf("accepted %s", bad)
		}
	}
	original := Operations()
	copy := Operations()
	copy[0].ResponseModes[0] = "bad"
	if reflect.DeepEqual(copy, Operations()) || original[0].ResponseModes[0] != "json" {
		t.Error("registry mutable")
	}
}

// These fixtures cover known framing conventions, not undocumented vendor
// guarantees. An unrecognized EOF must retain events and fail visibly.
func TestOperationAgentEvents(t *testing.T) {
	for _, terminal := range []string{
		"data: [DONE]\n\n",
		"event: done\ndata: {}\n\n",
		"event: finish\ndata: {}\n\n",
		"data: {\"status\":\"success\"}\n\n",
		"data: {\"choices\":[{\"finish_reason\":\"stop\"}]}\n\n",
	} {
		wire := "event: translation.delta\nid: request-1\ndata: {\"text\":\"翻译\",\"extra\":9007199254740993}\n\n" + terminal
		var events []Event
		err := readOperationEvents(context.Background(), &fragmentedReader{wire}, "agent_run", func(e Event) error { events = append(events, e); return nil })
		if err != nil || len(events) != 2 || events[0].Event != "translation.delta" || events[0].ID != "request-1" || !strings.Contains(events[0].Data, "9007199254740993") {
			t.Fatalf("agent events: %#v %v", events, err)
		}
	}
	emitted := 0
	err := readOperationEvents(context.Background(), strings.NewReader("event: progress\ndata: {\"status\":\"pending\"}\n\n"), "agent_run", func(Event) error { emitted++; return nil })
	if !errors.Is(err, io.ErrUnexpectedEOF) || emitted != 1 {
		t.Fatalf("unknown terminal accepted or lost: emitted=%d err=%v", emitted, err)
	}
	err = readOperationEvents(context.Background(), strings.NewReader("event: error\nid: agent-error\ndata: {\"error\":{\"code\":\"failed\"}}\n\n"), "agent_run", func(Event) error { return nil })
	var apiError *APIError
	if !errors.As(err, &apiError) || apiError.RequestID != "agent-error" {
		t.Fatalf("agent error: %v", err)
	}
}
