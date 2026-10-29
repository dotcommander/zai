package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/dotcommander/zai/internal/app"
)

func TestReadRequestObject(t *testing.T) {
	for _, input := range []string{"null", "[]", "42", "{} {}"} {
		if _, err := readRequestObject("-", strings.NewReader(input)); err == nil {
			t.Errorf("accepted %q", input)
		}
	}
	body, err := readRequestObject("-", strings.NewReader(`{"large":9007199254740993,"nested":{"tools":[{"name":"x"}]}}`))
	if err != nil {
		t.Fatal(err)
	}
	if string(body["large"]) != "9007199254740993" {
		t.Fatalf("number changed: %s", body["large"])
	}
}

func TestAdvancedBodyExplicitOverrides(t *testing.T) {
	cmd := newAdvancedCommand("image", "images")
	cmd.Flags().String("quality", "hd", "")
	cmd.Flags().Bool("unused", false, "")
	cmd.Flags().String("model", "default-model", "")
	cmd.SetIn(strings.NewReader(`{"quality":"standard","model":"request-model","unknown":{"n":9007199254740993}}`))
	if err := cmd.Flags().Set("request-json", "-"); err != nil {
		t.Fatal(err)
	}
	data, err := advancedBody(cmd, nil, "images")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte(`"quality":"standard"`)) || !bytes.Contains(data, []byte(`"model":"request-model"`)) {
		t.Fatalf("defaults replaced request: %s", data)
	}
	// Explicit flags replace only their top-level field; nested unknowns survive.
	dir := t.TempDir()
	path := filepath.Join(dir, "request.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("request-json", path); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("quality", "hd"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("model", "explicit"); err != nil {
		t.Fatal(err)
	}
	data, err = advancedBody(cmd, nil, "images")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte(`"quality":"hd"`)) || !bytes.Contains(data, []byte(`"model":"explicit"`)) || !bytes.Contains(data, []byte(`9007199254740993`)) {
		t.Fatalf("unexpected merge: %s", data)
	}
}

func TestAdvancedStdinConflict(t *testing.T) {
	cmd := newAdvancedCommand("chat", "chat")
	cmd.Flags().String("system", "", "")
	if err := cmd.Flags().Set("system", "-"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("request-json", "-"); err != nil {
		t.Fatal(err)
	}
	if _, err := advancedBody(cmd, nil, "chat"); err == nil {
		t.Fatal("expected stdin conflict")
	}
}

type fakeAdvancedAPI struct {
	req      app.OperationRequest
	deadline bool
	waited   bool
	waitErr  error
}

func (f *fakeAdvancedAPI) ExecuteOperation(ctx context.Context, req app.OperationRequest) (json.RawMessage, error) {
	f.req = req
	_, f.deadline = ctx.Deadline()
	return json.RawMessage(`{"id":"task-1","unknown":{"value":9007199254740993}}`), nil
}
func (f *fakeAdvancedAPI) StreamOperation(ctx context.Context, req app.OperationRequest, emit func(app.Event) error) error {
	f.req = req
	return emit(app.Event{Event: "tool.delta", ID: "event-1", Data: `{"tool_calls":[{"index":0,"arguments":"{"}]}`})
}
func (f *fakeAdvancedAPI) WaitTask(ctx context.Context, id string, interval time.Duration) (json.RawMessage, error) {
	f.waited = true
	if f.waitErr != nil {
		return nil, f.waitErr
	}
	return json.RawMessage(`{"id":"task-1","task_status":"SUCCESS"}`), nil
}

func (f *fakeAdvancedAPI) WaitAgent(ctx context.Context, agentID, asyncID string, interval time.Duration) (json.RawMessage, error) {
	f.waited = true
	return json.RawMessage(`{"status":"SUCCESS"}`), nil
}

func TestAdvancedDispatchAndEvents(t *testing.T) {
	fake := &fakeAdvancedAPI{}
	previous := advancedClientFactory
	advancedClientFactory = func() advancedAPI { return fake }
	defer func() { advancedClientFactory = previous }()
	cmd := newAdvancedCommand("responses", "responses")
	cmd.SetContext(context.Background())
	cmd.SetIn(strings.NewReader(`{"input":[{"type":"message"}],"extra":true}`))
	var output bytes.Buffer
	cmd.SetOut(&output)
	if err := cmd.Flags().Set("request-json", "-"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("events", "true"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	if fake.req.Operation != "responses" || !bytes.Contains(fake.req.JSON, []byte(`"stream":true`)) {
		t.Fatalf("wrong request: %+v", fake.req)
	}
	var event app.Event
	if err := json.Unmarshal(output.Bytes(), &event); err != nil {
		t.Fatal(err)
	}
	if event.Event != "tool.delta" || event.ID != "event-1" || !strings.Contains(event.Data, "tool_calls") {
		t.Fatalf("event lost fields: %+v", event)
	}
}

func TestAdvancedWaitRequiresDeadline(t *testing.T) {
	fake := &fakeAdvancedAPI{}
	previous := advancedClientFactory
	advancedClientFactory = func() advancedAPI { return fake }
	defer func() { advancedClientFactory = previous }()
	cmd := newAdvancedCommand("image", "image_async")
	cmd.SetContext(context.Background())
	cmd.SetOut(&bytes.Buffer{})
	if err := cmd.Flags().Set("wait", "true"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.RunE(cmd, nil); err == nil {
		t.Fatal("expected missing caller deadline error")
	}
	if err := cmd.Flags().Set("wait-timeout", "1m"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	if !fake.deadline || !fake.waited {
		t.Fatal("waiting did not receive caller deadline")
	}
}

func TestAttachAdvancedRelaxesOnlyAdvancedArgs(t *testing.T) {
	cmd := &cobra.Command{Use: "image", Args: cobra.ExactArgs(1), RunE: func(*cobra.Command, []string) error { return nil }}
	attachAdvanced(cmd, "images")
	if err := cmd.Args(cmd, nil); err == nil {
		t.Fatal("convenience argument check lost")
	}
	if err := cmd.Flags().Set("request-json", "-"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Args(cmd, nil); err != nil {
		t.Fatalf("raw request rejected: %v", err)
	}
}

func TestAdvancedPositionalRouting(t *testing.T) {
	for _, tc := range []struct{ op, field, value string }{{"agent_run", "agent_id", "slides_glm_agent"}, {"agent_result", "async_id", "async-1"}, {"agent_conversation", "conversation_id", "conversation-1"}} {
		cmd := newAdvancedCommand("test", tc.op)
		body, err := advancedBody(cmd, []string{tc.value}, tc.op)
		if err != nil {
			t.Fatal(err)
		}
		var object map[string]json.RawMessage
		if err := json.Unmarshal(body, &object); err != nil {
			t.Fatal(err)
		}
		var value string
		if err := json.Unmarshal(object[tc.field], &value); err != nil {
			t.Fatal(err)
		}
		if value != tc.value {
			t.Fatalf("%s missing positional value: %s", tc.op, body)
		}
	}
	fake := &fakeAdvancedAPI{}
	previous := advancedClientFactory
	advancedClientFactory = func() advancedAPI { return fake }
	defer func() { advancedClientFactory = previous }()
	cmd := newAdvancedCommand("audio", "audio")
	cmd.SetContext(context.Background())
	cmd.SetOut(&bytes.Buffer{})
	if err := executeAdvanced(cmd, []string{"audio.wav"}, "audio", nil); err != nil {
		t.Fatal(err)
	}
	if fake.req.UploadPath != "audio.wav" || bytes.Contains(fake.req.JSON, []byte(`"prompt"`)) {
		t.Fatalf("audio path encoded as prompt: %+v", fake.req)
	}
}

func TestAdvancedChatOverridesPreserveMessages(t *testing.T) {
	cmd := newAdvancedCommand("chat", "chat")
	cmd.Flags().Bool("think", false, "")
	cmd.Flags().String("system", "", "")
	if err := cmd.Flags().Set("think", "false"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("system", "new system"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("request-json", "-"); err != nil {
		t.Fatal(err)
	}
	cmd.SetIn(strings.NewReader(`{"thinking":{"type":"enabled","extra":9007199254740993},"messages":[{"role":"system","content":"old"},{"role":"assistant","tool_calls":[{"id":"call-1"}]}]}`))
	body, err := advancedBody(cmd, nil, "chat")
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{`"type":"disabled"`, `"extra":9007199254740993`, `"content":"new system"`, `"tool_calls"`} {
		if !bytes.Contains(body, []byte(text)) {
			t.Fatalf("lost override or field %s: %s", text, body)
		}
	}
	if bytes.Contains(body, []byte(`"old"`)) {
		t.Fatalf("old system survived: %s", body)
	}
}

func TestAdvancedStreamJSONAndWaitGuard(t *testing.T) {
	fake := &fakeAdvancedAPI{}
	previous := advancedClientFactory
	advancedClientFactory = func() advancedAPI { return fake }
	defer func() { advancedClientFactory = previous }()
	cmd := newAdvancedCommand("responses", "responses")
	cmd.SetContext(context.Background())
	cmd.SetOut(&bytes.Buffer{})
	if err := cmd.Flags().Set("request-json", "-"); err != nil {
		t.Fatal(err)
	}
	cmd.SetIn(strings.NewReader(`{"stream":true}`))
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	if fake.req.Operation != "responses" {
		t.Fatal("stream JSON did not dispatch")
	}
	guarded := newAdvancedCommand("image", "images")
	if err := guarded.Flags().Set("wait", "true"); err != nil {
		t.Fatal(err)
	}
	if err := guarded.Flags().Set("wait-timeout", "1m"); err != nil {
		t.Fatal(err)
	}
	fake.req = app.OperationRequest{}
	if err := guarded.RunE(guarded, nil); err == nil {
		t.Fatal("expected sync operation wait guard")
	}
	if fake.req.Operation != "" {
		t.Fatal("wait guard submitted synchronous request")
	}
}

func TestAdvancedWaitFailureRetainsSubmission(t *testing.T) {
	fake := &fakeAdvancedAPI{waitErr: errors.New("poll failed")}
	previous := advancedClientFactory
	advancedClientFactory = func() advancedAPI { return fake }
	defer func() { advancedClientFactory = previous }()
	cmd := newAdvancedCommand("image", "image_async")
	cmd.SetContext(context.Background())
	var stderr bytes.Buffer
	cmd.SetErr(&stderr)
	if err := cmd.Flags().Set("wait", "true"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("wait-timeout", "1m"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.RunE(cmd, nil); err == nil {
		t.Fatal("expected polling failure")
	}
	if !strings.Contains(stderr.String(), `"id":"task-1"`) {
		t.Fatalf("submission id lost: %s", stderr.String())
	}
}

func TestAdvancedResultViews(t *testing.T) {
	cmd := newAdvancedCommand("chat", "chat")
	cmd.Flags().Bool("json", false, "")
	text := json.RawMessage(`{"choices":[{"message":{"content":"hello"}}],"usage":{"total_tokens":1}}`)
	if got := advancedResultText(cmd, "chat", text); got != "hello" {
		t.Fatalf("text view: %s", got)
	}
	tool := json.RawMessage(`{"choices":[{"message":{"content":null,"tool_calls":[{"id":"call-1"}]}}]}`)
	if got := advancedResultText(cmd, "chat", tool); got != string(tool) {
		t.Fatalf("tool-only response lost: %s", got)
	}
	if err := cmd.Flags().Set("json", "true"); err != nil {
		t.Fatal(err)
	}
	if got := advancedResultText(cmd, "chat", text); got != string(text) {
		t.Fatalf("JSON response lost: %s", got)
	}
}

func TestAdvancedTextWithToolsRetained(t *testing.T) {
	cmd := newAdvancedCommand("chat", "chat")
	result := json.RawMessage(`{"choices":[{"message":{"content":"calling tool","tool_calls":[{"id":"1"}]}}]}`)
	if got := advancedResultText(cmd, "chat", result); got != string(result) {
		t.Fatalf("tool call lost: %s", got)
	}
}

func TestAdvancedWaitJSONRetainsBothResponses(t *testing.T) {
	previous := advancedClientFactory
	advancedClientFactory = func() advancedAPI { return &fakeAdvancedAPI{} }
	defer func() { advancedClientFactory = previous }()
	cmd := newAdvancedCommand("image", "image_async")
	cmd.Flags().Bool("json", true, "")
	if err := cmd.Flags().Set("wait", "true"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("wait-timeout", "1m"); err != nil {
		t.Fatal(err)
	}
	cmd.SetContext(context.Background())
	var out bytes.Buffer
	cmd.SetOut(&out)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(out.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(envelope["submission"], []byte(`"id":"task-1"`)) || len(envelope["result"]) == 0 {
		t.Fatalf("lost submission: %s", out.String())
	}
}

func TestAdvancedAgentStreamMode(t *testing.T) {
	for _, tc := range []struct {
		body   string
		events bool
		want   string
	}{{`{}`, false, `false`}, {`{"stream":true}`, false, `true`}, {`{"stream":false}`, true, `true`}} {
		cmd := newAdvancedCommand("run", "agent_run")
		cmd.SetIn(strings.NewReader(tc.body))
		if err := cmd.Flags().Set("request-json", "-"); err != nil {
			t.Fatal(err)
		}
		if tc.events {
			if err := cmd.Flags().Set("events", "true"); err != nil {
				t.Fatal(err)
			}
		}
		body, err := advancedBody(cmd, nil, "agent_run")
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(body, &fields); err != nil {
			t.Fatal(err)
		}
		if string(fields["stream"]) != tc.want {
			t.Fatalf("stream = %s, want %s", fields["stream"], tc.want)
		}
	}
}

func TestAdvancedAudioHotwordsArray(t *testing.T) {
	for _, tc := range []struct {
		name, flag string
		changed    bool
		want       string
	}{
		{"preserve JSON", "", false, `["from-json"]`},
		{"explicit comma list", " alpha, beta , ,gamma ", true, `["alpha","beta","gamma"]`},
		{"explicit empty list", "", true, `[]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newAdvancedCommand("audio", "audio")
			cmd.Flags().String("hotwords", "", "")
			if err := cmd.Flags().Set("request-json", "-"); err != nil {
				t.Fatal(err)
			}
			cmd.SetIn(strings.NewReader(`{"hotwords":["from-json"]}`))
			if tc.changed {
				if err := cmd.Flags().Set("hotwords", tc.flag); err != nil {
					t.Fatal(err)
				}
			}
			body, err := advancedBody(cmd, nil, "audio")
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(body, &fields); err != nil {
				t.Fatal(err)
			}
			if string(fields["hotwords"]) != tc.want {
				t.Fatalf("hotwords = %s, want %s", fields["hotwords"], tc.want)
			}
		})
	}
}

func TestAdvancedToolStreamFlagDispatch(t *testing.T) {
	fake := &fakeAdvancedAPI{}
	previous := advancedClientFactory
	advancedClientFactory = func() advancedAPI { return fake }
	defer func() { advancedClientFactory = previous }()
	convenience := false
	cmd := &cobra.Command{Use: "chat", RunE: func(*cobra.Command, []string) error { convenience = true; return nil }}
	attachAdvanced(cmd, "chat")
	cmd.SetContext(context.Background())
	var output bytes.Buffer
	cmd.SetOut(&output)
	if err := cmd.Flags().Set("tool-stream", "true"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.RunE(cmd, []string{"Call a tool"}); err != nil {
		t.Fatal(err)
	}
	if convenience || !bytes.Contains(fake.req.JSON, []byte(`"tool_stream":true`)) || !bytes.Contains(fake.req.JSON, []byte(`"stream":true`)) {
		t.Fatalf("tool streaming did not use advanced SSE: convenience=%v request=%s", convenience, fake.req.JSON)
	}
	var event app.Event
	if err := json.Unmarshal(output.Bytes(), &event); err != nil {
		t.Fatal(err)
	}
	if event.Event != "tool.delta" || !strings.Contains(event.Data, "tool_calls") {
		t.Fatalf("tool event lost: %+v", event)
	}
}

func TestAdvancedToolStreamExplicitPrecedence(t *testing.T) {
	for _, tc := range []struct {
		changed bool
		want    string
	}{{false, "true"}, {true, "false"}} {
		cmd := newAdvancedCommand("chat", "chat")
		if err := cmd.Flags().Set("request-json", "-"); err != nil {
			t.Fatal(err)
		}
		cmd.SetIn(strings.NewReader(`{"stream":true,"tool_stream":true}`))
		if tc.changed {
			if err := cmd.Flags().Set("tool-stream", "false"); err != nil {
				t.Fatal(err)
			}
		}
		body, err := advancedBody(cmd, nil, "chat")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(body, []byte(`"tool_stream":`+tc.want)) || !bytes.Contains(body, []byte(`"stream":true`)) {
			t.Fatalf("unexpected merge: %s", body)
		}
	}
}

type capabilityAPI struct {
	*fakeAdvancedAPI
}

func (f *capabilityAPI) ExecuteOperation(ctx context.Context, req app.OperationRequest) (json.RawMessage, error) {
	f.req = req
	return json.RawMessage(`{"choices":[{"message":{"content":"Found it","tool_calls":[{"id":"mcp-1","type":"mcp","result":{"document":"report"}}]}}],"usage":{"prompt_tokens_details":{"cached_tokens":123}}}`), nil
}

func TestAdvancedMCPAndCacheFields(t *testing.T) {
	fake := &capabilityAPI{&fakeAdvancedAPI{}}
	previous := advancedClientFactory
	advancedClientFactory = func() advancedAPI { return fake }
	defer func() { advancedClientFactory = previous }()
	cmd := newAdvancedCommand("chat", "chat")
	cmd.Flags().Bool("json", false, "")
	if err := cmd.Flags().Set("json", "true"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("request-json", "-"); err != nil {
		t.Fatal(err)
	}
	cmd.SetContext(context.Background())
	cmd.SetIn(strings.NewReader(`{"messages":[{"role":"user","content":"Read report"}],"tools":[{"type":"mcp","mcp":{"server_label":"documents","server_url":"https://example.com/mcp","transport_type":"streamable-http","headers":{"X-Custom":"value"},"allowed_tools":["read_document"]}}],"tool_choice":"auto"}`))
	var output bytes.Buffer
	cmd.SetOut(&output)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{`"type":"mcp"`, `"server_label":"documents"`, `"transport_type":"streamable-http"`, `"headers":{"X-Custom":"value"}`, `"allowed_tools":["read_document"]`, `"tool_choice":"auto"`} {
		if !bytes.Contains(fake.req.JSON, []byte(field)) {
			t.Fatalf("request lost %s: %s", field, fake.req.JSON)
		}
	}
	for _, field := range []string{`"tool_calls"`, `"cached_tokens":123`} {
		if !strings.Contains(output.String(), field) {
			t.Fatalf("response lost %s: %s", field, output.String())
		}
	}
}
