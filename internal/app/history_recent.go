package app

import (
	"bytes"
	"encoding/json"
	"os"
	"slices"
	"strings"
)

// readRecentHistoryEntries reads backwards so a small limit does not decode or
// retain the entire JSONL file. A single partial record is capped at the same
// size as the forward reader, including its newline when present.
func readRecentHistoryEntries(file *os.File, limit int) ([]HistoryEntry, error) {
	const blockSize = 64 * 1024
	const maxRecordBytes = 1024 * 1024
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	var recent []HistoryEntry // newest first until return
	var suffix []byte
	oversized := false
	terminated := false // Whether the current record ends in a newline.
	block := make([]byte, blockSize)
	for pos := info.Size(); pos > 0 && len(recent) < limit; {
		start := pos - blockSize
		if start < 0 {
			start = 0
		}
		chunk := block[:int(pos-start)]
		if _, err := file.ReadAt(chunk, start); err != nil {
			return nil, err
		}
		pos = start
		end := len(chunk)
		for i := len(chunk) - 1; i >= 0; i-- {
			if chunk[i] != '\n' {
				continue
			}
			// The bytes after this delimiter belong to the next record.
			recordSize := len(chunk[i+1:end]) + len(suffix)
			if terminated {
				recordSize++
			}
			if !oversized && recordSize <= maxRecordBytes {
				recent = appendRecentHistoryRecord(recent, string(chunk[i+1:end])+string(suffix))
			}
			suffix = nil
			oversized = false
			terminated = true
			end = i
			if len(recent) == limit {
				break
			}
		}
		if len(recent) == limit {
			break
		}
		// Prefix is an incomplete record; keep only bounded bytes across blocks.
		if !oversized {
			recordSize := len(chunk[:end]) + len(suffix)
			if terminated {
				recordSize++
			}
			if recordSize > maxRecordBytes {
				suffix = nil
				oversized = true
			} else {
				suffix = append(bytes.Clone(chunk[:end]), suffix...)
			}
		}
	}
	// At byte zero, there may be a record without a preceding newline.
	if len(recent) < limit && !oversized && len(suffix) > 0 {
		recent = appendRecentHistoryRecord(recent, string(suffix))
	}
	slices.Reverse(recent)
	if recent == nil {
		return []HistoryEntry{}, nil
	}
	return recent, nil
}

func appendRecentHistoryRecord(entries []HistoryEntry, record string) []HistoryEntry {
	line := strings.TrimSpace(record)
	var entry HistoryEntry
	if line != "" && json.Unmarshal([]byte(line), &entry) == nil {
		return append(entries, entry)
	}
	return entries
}
