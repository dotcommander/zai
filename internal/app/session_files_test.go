package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type sessionFileFake struct{ reads int }

func (f *sessionFileFake) ReadFile(name string) ([]byte, error) {
	f.reads++
	return os.ReadFile(name)
}

func TestEfficiencySessionClientUsesInjectedReader(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "context.txt")
	if err := os.WriteFile(path, []byte("context data"), 0600); err != nil {
		t.Fatal(err)
	}
	fake := &sessionFileFake{}
	client := NewClientWithDeps(ClientConfig{}, DiscardLogger(), nil, &ClientDeps{FileReader: fake})
	client.EnableSessionFileReuse(path)
	for range 2 {
		content, err := client.buildContent(context.Background(), "prompt", path)
		if err != nil || !strings.Contains(content, "context data") {
			t.Fatalf("content = %q, err = %v", content, err)
		}
	}
	if fake.reads != 1 {
		t.Fatalf("injected reader used %d times, want 1", fake.reads)
	}
}
