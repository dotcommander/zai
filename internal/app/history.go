package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// HistoryEntry represents a single chat history entry.
type HistoryEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	Prompt     string    `json:"prompt"`
	Response   string    `json:"response"`
	Model      string    `json:"model"`
	TokenUsage Usage     `json:"token_usage,omitempty"`
}

// FileHistoryStore implements HistoryStore with JSONL file storage.
type FileHistoryStore struct {
	path string
}

// NewFileHistoryStore creates a history store at the given path.
// If path is empty, uses default XDG location.
func NewFileHistoryStore(path string) *FileHistoryStore {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			path = "history.jsonl"
		} else {
			path = filepath.Join(home, ".config", "zai", "history.jsonl")
		}
	}
	return &FileHistoryStore{path: path}
}

// Save appends an entry to the history file.
func (h *FileHistoryStore) Save(entry HistoryEntry) error {
	// Ensure directory exists
	dir := filepath.Dir(h.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create history directory: %w", err)
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal history entry: %w", err)
	}

	file, err := os.OpenFile(h.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(string(data) + "\n"); err != nil {
		return fmt.Errorf("failed to write history entry: %w", err)
	}

	return nil
}

// GetRecent returns the most recent history entries.
func (h *FileHistoryStore) GetRecent(limit int) ([]HistoryEntry, error) {
	file, err := os.Open(h.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []HistoryEntry{}, nil
		}
		return nil, fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	var entries []HistoryEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var entry HistoryEntry
		if err := json.Unmarshal([]byte(line), &entry); err == nil {
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading history file: %w", err)
	}

	// Return most recent entries
	if limit > 0 && len(entries) > limit {
		entries = entries[len(entries)-limit:]
	}

	return entries, nil
}

// Path returns the history file path.
func (h *FileHistoryStore) Path() string {
	return h.path
}
