package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestEfficiencyHistoryRecentSkipsInvalidAndOversized(t *testing.T) {
	t.Parallel()
	entry := func(id string) string {
		t.Helper()
		data, err := json.Marshal(HistoryEntry{ID: id, Prompt: strings.Repeat(id, 70000)})
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	path := filepath.Join(t.TempDir(), "history.jsonl")
	// The oversized record crosses multiple blocks; neither it nor the invalid
	// record should displace one of the last two valid records.
	content := entry("old") + "\n" + strings.Repeat("x", 1024*1024+1) + "\n" +
		entry("middle") + "\nmalformed\n" + entry("latest")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	store := NewFileHistoryStore(path)
	got, err := store.GetRecent(2)
	if err != nil {
		t.Fatal(err)
	}
	ids := []string{}
	for _, v := range got {
		ids = append(ids, v.ID)
	}
	if want := []string{"middle", "latest"}; !reflect.DeepEqual(ids, want) {
		t.Fatalf("recent IDs = %v, want %v", ids, want)
	}
	all, err := store.GetRecent(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 3 {
		t.Fatalf("unlimited history contains %d entries, want 3", len(all))
	}
}

func TestEfficiencyHistoryRecentRecordBoundary(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "history.jsonl")
	entry := func(prompt string) []byte {
		t.Helper()
		data, err := json.Marshal(HistoryEntry{ID: "boundary", Prompt: prompt})
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	base := entry("")
	valid := entry(strings.Repeat("x", 1024*1024-len(base)-1))
	oversized := entry(strings.Repeat("x", 1024*1024-len(base)))
	if len(valid)+1 != 1024*1024 || len(oversized)+1 != 1024*1024+1 {
		t.Fatal("fixture record lengths are incorrect")
	}
	if err := os.WriteFile(path, append(append(valid, '\n'), append(oversized, '\n')...), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := NewFileHistoryStore(path).GetRecent(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || len(got[0].Prompt) != 1024*1024-len(base)-1 {
		t.Fatalf("boundary entry not returned: count=%d", len(got))
	}
}
