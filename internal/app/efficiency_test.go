package app

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type efficiencyFileReader struct{ data []byte }

func (r efficiencyFileReader) ReadFile(string) ([]byte, error) { return r.data, nil }

func TestEfficiencyAudioMultipartFramingAndReplay(t *testing.T) {
	t.Parallel()
	data := []byte("audio-payload")
	client := &Client{config: ClientConfig{APIKey: "test", BaseURL: "https://example.test"}, fileReader: efficiencyFileReader{data}}
	req, err := client.audioRequest(context.Background(), "fixture.wav", TranscriptionOptions{Model: "sample", Prompt: "context", UserID: "u", RequestID: "r", Hotwords: []string{"word"}})
	if err != nil {
		t.Fatal(err)
	}
	original, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	replay, err := req.GetBody()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = replay.Close() })
	replayed, err := io.ReadAll(replay)
	if err != nil || string(replayed) != string(original) {
		t.Fatalf("replayed multipart differs: %v", err)
	}
	if req.ContentLength != int64(len(original)) {
		t.Fatalf("content length %d != %d", req.ContentLength, len(original))
	}
	parsed := multipart.NewReader(strings.NewReader(string(original)), strings.TrimPrefix(req.Header.Get("Content-Type"), "multipart/form-data; boundary="))
	values := map[string]string{}
	for {
		part, err := parsed.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(part)
		if err != nil {
			t.Fatal(err)
		}
		values[part.FormName()] = string(body)
	}
	for key, want := range map[string]string{"file": string(data), "model": "sample", "prompt": "context", "user_id": "u", "request_id": "r"} {
		if values[key] != want {
			t.Fatalf("%s = %q, want %q", key, values[key], want)
		}
	}
	var hotwords []string
	if err := json.Unmarshal([]byte(values["hotwords"]), &hotwords); err != nil || len(hotwords) != 1 || hotwords[0] != "word" {
		t.Fatalf("hotwords = %v (%v)", hotwords, err)
	}
}

func TestEfficiencySearchEvictsOnlyFailedEntry(t *testing.T) {
	t.Parallel()
	cache := NewFileSearchCache(t.TempDir())
	options := SearchOptions{}
	if err := cache.Set("expired", options, nil, -time.Second); err != nil {
		t.Fatal(err)
	}
	if err := cache.Set("other", options, nil, -time.Second); err != nil {
		t.Fatal(err)
	}
	if err := cache.Set("live", options, []SearchResult{{Title: "retained"}}, time.Hour); err != nil {
		t.Fatal(err)
	}
	if _, ok := cache.Get("expired", options); ok {
		t.Fatal("expired cache hit")
	}
	if _, err := os.Stat(filepath.Join(cache.dir, generateCacheKey("expired", options)+".json")); !os.IsNotExist(err) {
		t.Fatalf("expired entry not evicted: %v", err)
	}
	if _, err := os.Stat(filepath.Join(cache.dir, generateCacheKey("other", options)+".json")); err != nil {
		t.Fatalf("unrelated entry removed: %v", err)
	}
	corrupt := filepath.Join(cache.dir, generateCacheKey("corrupt", options)+".json")
	if err := os.WriteFile(corrupt, []byte("invalid"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, ok := cache.Get("corrupt", options); ok {
		t.Fatal("corrupt cache hit")
	}
	if _, err := os.Stat(corrupt); !os.IsNotExist(err) {
		t.Fatalf("corrupt entry not evicted: %v", err)
	}
	if results, ok := cache.Get("live", options); !ok || len(results) != 1 || results[0].Title != "retained" {
		t.Fatalf("live result lost: %v (%v)", results, ok)
	}
}
