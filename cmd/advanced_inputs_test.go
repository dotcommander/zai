package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestAdvancedChatFileInputs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "context.txt")
	if err := os.WriteFile(path, []byte("retained context"), 0600); err != nil {
		t.Fatal(err)
	}
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("file", "", "")
	cmd.Flags().String("request-json", "", "")
	cmd.Flags().String("system", "", "")
	cmd.Flags().Bool("think", false, "")
	_ = cmd.Flags().Set("file", path)
	body := map[string]json.RawMessage{"messages": json.RawMessage(`[{"role":"assistant","tool_calls":[{"id":"tool-1","arguments":9007199254740993}]}]`)}
	if err := advancedChatInputs(cmd, []string{"explain"}, body); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body["messages"]), "retained context") || !strings.Contains(string(body["messages"]), "9007199254740993") {
		t.Fatalf("lost context: %s", body["messages"])
	}
	_ = cmd.Flags().Set("file", "-")
	_ = cmd.Flags().Set("system", "-")
	if err := advancedChatInputs(cmd, nil, body); err == nil {
		t.Error("shared stdin accepted")
	}
}

func TestAdvancedVisionJSONInputs(t *testing.T) {
	old := visionCmd
	defer func() { visionCmd = old }()
	visionCmd = &cobra.Command{Use: "vision"}
	visionCmd.Flags().String("file", "", "")
	visionCmd.Flags().String("request-json", "", "")
	visionCmd.Flags().String("prompt", "", "")
	_ = visionCmd.Flags().Set("file", "https://example.com/photo.png")
	body := map[string]json.RawMessage{}
	if err := advancedChatInputs(visionCmd, []string{"describe"}, body); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body["messages"]), "image_url") || !strings.Contains(string(body["messages"]), "https://example.com/photo.png") {
		t.Fatalf("no image: %s", body["messages"])
	}
}
