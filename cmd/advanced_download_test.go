package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestAdvancedArtifactDownload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("artifact")); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	target := filepath.Join(t.TempDir(), "report.pdf")
	cmd := newAdvancedCommand("conversation", "agent_conversation")
	_ = cmd.Flags().Set("download", target)
	result, _ := json.Marshal(map[string]any{"results": []any{map[string]string{"pdf_url": server.URL + "/report"}}})
	if err := validateAdvancedDownload(cmd, "agent_conversation", false, false); err != nil {
		t.Fatal(err)
	}
	if err := downloadAdvancedResult(context.Background(), cmd, &fakeAdvancedAPI{}, result); err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(target); err != nil || string(got) != "artifact" {
		t.Fatalf("artifact: %q %v", got, err)
	}
	_ = cmd.Flags().Set("max-download-bytes", "2")
	if err := downloadAdvancedResult(context.Background(), cmd, &fakeAdvancedAPI{}, result); err == nil {
		t.Fatal("unbounded download")
	}
	if got, _ := os.ReadFile(target); string(got) != "artifact" {
		t.Fatal("failed download replaced prior artifact")
	}
	if err := validateAdvancedDownload(cmd, "videos", false, false); err == nil {
		t.Fatal("submitted download without wait")
	}
	if err := validateAdvancedDownload(cmd, "agent_conversation", false, true); err == nil {
		t.Fatal("event download accepted")
	}
}
