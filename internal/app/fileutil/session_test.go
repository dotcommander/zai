package fileutil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type sessionCountingReader struct {
	calls int
	err   error
}

func (r *sessionCountingReader) ReadFile(path string) ([]byte, error) {
	r.calls++
	if r.err != nil {
		return nil, r.err
	}
	return os.ReadFile(path)
}

func TestEfficiencySessionFileReuseAndRefresh(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "context.txt")
	if err := os.WriteFile(path, []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	underlying := &sessionCountingReader{}
	reader := NewSessionFileReader(underlying, path)
	for range 2 {
		got, err := reader.ReadFile(path)
		if err != nil || string(got) != "old" {
			t.Fatalf("first content = %q, %v", got, err)
		}
	}
	if underlying.calls != 1 {
		t.Fatalf("unchanged file read %d times, want 1", underlying.calls)
	}

	// Explicit modtime avoids depending on filesystem timestamp resolution.
	if err := os.WriteFile(path, []byte("new"), 0600); err != nil {
		t.Fatal(err)
	}
	stamp := time.Unix(1_000_000, 0)
	if err := os.Chtimes(path, stamp, stamp); err != nil {
		t.Fatal(err)
	}
	got, err := reader.ReadFile(path)
	if err != nil || string(got) != "new" || underlying.calls != 2 {
		t.Fatalf("refreshed content = %q, err %v, calls %d", got, err, underlying.calls)
	}

	other := filepath.Join(t.TempDir(), "replacement")
	if err := os.WriteFile(other, []byte("replaced"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(other, path); err != nil {
		t.Fatal(err)
	}
	got, err = reader.ReadFile(path)
	if err != nil || string(got) != "replaced" || underlying.calls != 3 {
		t.Fatalf("replacement = %q, err %v, calls %d", got, err, underlying.calls)
	}
}

func TestEfficiencySessionFileReadErrorsNotCached(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "context.txt")
	if err := os.WriteFile(path, []byte("ok"), 0600); err != nil {
		t.Fatal(err)
	}
	underlying := &sessionCountingReader{err: errors.New("temporary failure")}
	reader := NewSessionFileReader(underlying, path)
	if _, err := reader.ReadFile(path); err == nil {
		t.Fatal("expected read failure")
	}
	underlying.err = nil
	got, err := reader.ReadFile(path)
	if err != nil || string(got) != "ok" || underlying.calls != 2 {
		t.Fatalf("retry = %q, err %v, calls %d", got, err, underlying.calls)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if _, err := reader.ReadFile(path); err == nil {
		t.Fatal("stale cached content returned for missing file")
	}
}
