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

// HistoryEntry represents a single chat, image generation, or web reader history entry.
type HistoryEntry struct {
	ID         string      `json:"id,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	Prompt     string      `json:"prompt"`
	Response   interface{} `json:"response"` // Support string or complex response
	Model      string      `json:"model"`
	TokenUsage Usage       `json:"token_usage,omitempty"`

	// Image generation fields
	ImageURL    string `json:"image_url,omitempty"`
	ImageSize   string `json:"image_size,omitempty"`
	ImageFormat string `json:"image_format,omitempty"`
	Type        string `json:"type"` // "chat", "image", or "web"

	// Web reader fields
	WebSources []string `json:"web_sources,omitempty"`
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

	// Handle response conversion for compatibility
	if _, ok := entry.Response.(string); ok {
		// Response is already a string, no conversion needed
	} else {
		// Response is complex type, convert to JSON string for storage
		data, err := json.Marshal(entry.Response)
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}
		entry.Response = string(data)
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

// NewImageHistoryEntry creates a history entry for image generation.
func NewImageHistoryEntry(prompt string, imageData ImageData, model string) HistoryEntry {
	return HistoryEntry{
		Timestamp:   time.Now(),
		Prompt:      prompt,
		Response:    fmt.Sprintf("Generated image: %s", imageData.URL),
		Model:       model,
		ImageURL:    imageData.URL,
		ImageSize:   fmt.Sprintf("%dx%d", imageData.Width, imageData.Height),
		ImageFormat: imageData.Format,
		Type:        "image",
	}
}

// NewChatHistoryEntry creates a history entry for chat (sets type to "chat").
func NewChatHistoryEntry(timestamp time.Time, prompt, response, model string, usage Usage) HistoryEntry {
	return HistoryEntry{
		Timestamp:  timestamp,
		Prompt:     prompt,
		Response:   response,
		Model:      model,
		TokenUsage: usage,
		Type:       "chat",
	}
}

// NewWebHistoryEntry creates a history entry for web content fetching.
func NewWebHistoryEntry(id, prompt string, resp *WebReaderResponse, sources []string) HistoryEntry {
	return HistoryEntry{
		ID:        id,
		Timestamp: time.Now(),
		Prompt:    prompt,
		Response: map[string]interface{}{
			"url":         resp.ReaderResult.URL,
			"title":       resp.ReaderResult.Title,
			"description": resp.ReaderResult.Description,
			"content":     resp.ReaderResult.Content,
		},
		Model:      "web-reader",
		Type:       "web",
		WebSources: sources,
	}
}

// NewSearchHistoryEntry creates a history entry for web search.
func NewSearchHistoryEntry(timestamp time.Time, query string, resp *WebSearchResponse) HistoryEntry {
	return HistoryEntry{
		Timestamp: timestamp,
		Prompt:    fmt.Sprintf("search: %s", query),
		Response: map[string]interface{}{
			"query":        query,
			"result_count": len(resp.SearchResult),
			"results":      resp.SearchResult,
		},
		Model: "web-search",
		Type:  "web_search",
	}
}

// NewAudioHistoryEntry creates a history entry for audio transcription.
func NewAudioHistoryEntry(text string, model string) HistoryEntry {
	return HistoryEntry{
		Timestamp: time.Now(),
		Prompt:    "audio transcription",
		Response:  text,
		Model:     model,
		Type:      "audio",
	}
}
