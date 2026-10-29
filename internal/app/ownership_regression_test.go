package app

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestRepairChatMessagesPreserveContextStorage(t *testing.T) {
	t.Parallel()
	storage := []Message{{Role: "user", Content: "prior"}, {Role: "assistant", Content: "reply"}, {Role: "user", Content: "sentinel"}, {Role: "assistant", Content: "sentinel"}}
	before := append([]Message(nil), storage...)
	opts := ChatOptions{SystemPrompt: "system", Context: storage[:2]}
	client := &Client{}
	want := []Message{{Role: "system", Content: "system"}, before[0], before[1], {Role: "user", Content: "current"}}
	for _, got := range [][]Message{client.buildMessages("current", opts), client.buildMessagesWithContext("current", opts)} {
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("messages = %#v; want %#v", got, want)
		}
	}
	if !reflect.DeepEqual(storage, before) {
		t.Fatalf("context backing storage changed: %#v", storage)
	}
}

func TestRepairHistorySkipsOversizedRecords(t *testing.T) {
	t.Parallel()
	first, _ := json.Marshal(HistoryEntry{ID: "first"})
	last, _ := json.Marshal(HistoryEntry{ID: "last"})
	input := string(first) + "\n" + strings.Repeat("x", 1024*1024+1) + "\nmalformed\n" + string(last)
	entries, err := readHistoryEntries(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 || entries[0].ID != "first" || entries[1].ID != "last" {
		t.Fatalf("entries = %#v", entries)
	}
}

type failingHistoryReader struct{ err error }

func (r failingHistoryReader) Read([]byte) (int, error) { return 0, r.err }

func TestRepairHistoryPropagatesReadError(t *testing.T) {
	t.Parallel()
	expected := errors.New("read failed")
	_, err := readHistoryEntries(io.MultiReader(strings.NewReader("{}\n"), failingHistoryReader{expected}))
	if !errors.Is(err, expected) {
		t.Fatalf("error = %v", err)
	}
}

type embeddingModelDoer struct {
	t    *testing.T
	want string
}

func (d embeddingModelDoer) Do(req *http.Request) (*http.Response, error) {
	var body EmbeddingRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		d.t.Fatal(err)
	}
	if body.Model != d.want {
		d.t.Errorf("model = %q; want %q", body.Model, d.want)
	}
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{}"))}, nil
}

func TestRepairEmbeddingModelPrecedence(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, option, configured, want string }{
		{name: "fallback", want: "embedding-3"},
		{name: "configuration", configured: "configured", want: "configured"},
		{name: "option", option: "explicit", configured: "configured", want: "explicit"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := NewClient(ClientConfig{APIKey: "test", BaseURL: "https://example.invalid", EmbeddingModel: tc.configured}, DiscardLogger(), nil, embeddingModelDoer{t, tc.want})
			if _, err := client.CreateEmbedding(context.Background(), []string{"text"}, EmbeddingOptions{Model: tc.option}); err != nil {
				t.Fatal(err)
			}
		})
	}
}
