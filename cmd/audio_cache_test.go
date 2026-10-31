package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEfficiencyAudioCacheHashAndLimit(t *testing.T) {
	// getCachePath uses HOME to choose its directory, so this test cannot run in parallel.
	home := t.TempDir()
	t.Setenv("HOME", home)
	source := filepath.Join(t.TempDir(), "source.wav")
	content := []byte("audio source content")
	if err := os.WriteFile(source, content, 0600); err != nil {
		t.Fatal(err)
	}
	path, err := getCachePath(source)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(content)
	want := filepath.Join(home, ".cache", "zai", "audio", hex.EncodeToString(digest[:8])+".json")
	if path != want {
		t.Fatalf("cache path = %q, want %q", path, want)
	}
	large := filepath.Join(t.TempDir(), "large.wav")
	file, err := os.Create(large)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = file.Close() })
	if err := file.Truncate(MaxFileSize + 1); err != nil {
		t.Fatal(err)
	}
	if _, err := getCachePath(large); err == nil {
		t.Fatal("oversized source accepted")
	}
}

func TestEfficiencyAudioJournalRecoveryAndCompaction(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "audio.json")
	legacy := AudioCache{Chunks: map[int]string{0: "first"}}
	data, err := json.Marshal(legacy)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	if err := appendAudioChunk(path, 1, "second"); err != nil {
		t.Fatal(err)
	}
	// A partial final line is ignored and the next append is isolated from it.
	journal := audioJournalPath(path)
	file, err := os.OpenFile(journal, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString("partial"); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if err := appendAudioChunk(path, 2, "third"); err != nil {
		t.Fatal(err)
	}
	cache, err := loadCache(path)
	if err != nil {
		t.Fatal(err)
	}
	for index, want := range []string{"first", "second", "third"} {
		if cache.Chunks[index] != want {
			t.Fatalf("chunk %d = %q, want %q", index, cache.Chunks[index], want)
		}
	}
	if err := compactAudioCache(path, cache); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(journal); !os.IsNotExist(err) {
		t.Fatalf("journal remains after compaction: %v", err)
	}
	cache, err = loadCache(path)
	if err != nil || len(cache.Chunks) != 3 {
		t.Fatalf("compacted cache = %#v, error = %v", cache, err)
	}
	if err := clearAudioCache(path); err != nil {
		t.Fatal(err)
	}
	cache, err = loadCache(path)
	if err != nil || len(cache.Chunks) != 0 {
		t.Fatalf("cleared cache = %#v, error = %v", cache, err)
	}
}

func TestEfficiencyAudioJournalIncompleteFinalRecord(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "audio.json")
	first, err := json.Marshal(audioChunkRecord{Index: 0, Text: "one"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := json.Marshal(audioChunkRecord{Index: 1, Text: "unfinished"})
	if err != nil {
		t.Fatal(err)
	}
	// Simulate a crash halfway through the last encoded record.
	incomplete := append(append(append([]byte(nil), first...), '\n'), second[:len(second)/2]...)
	if err := os.WriteFile(audioJournalPath(path), incomplete, 0600); err != nil {
		t.Fatal(err)
	}
	cache, err := loadCache(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cache.Chunks) != 1 || cache.Chunks[0] != "one" {
		t.Fatalf("incomplete journal replayed: %#v", cache.Chunks)
	}
}
