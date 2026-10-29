package app

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// SearchCache interface for search caching (ISP compliance).
type SearchCache interface {
	Get(query string, opts SearchOptions) ([]SearchResult, bool)
	Set(query string, opts SearchOptions, results []SearchResult, ttl time.Duration) error
	Clear() error
	Cleanup() error // Remove expired entries
}

// FileSearchCache implements persistent file-based caching.
type FileSearchCache struct {
	dir            string
	mutex          sync.RWMutex
	cleanupRunning sync.Mutex // Guards async cleanup to prevent storms
}

// NewFileSearchCache creates a new file-based search cache.
func NewFileSearchCache(dir string) *FileSearchCache {
	return &FileSearchCache{
		dir: dir,
	}
}

// Get retrieves cached search results.
func (fsc *FileSearchCache) Get(query string, opts SearchOptions) ([]SearchResult, bool) {
	fsc.mutex.RLock()
	defer fsc.mutex.RUnlock()

	cacheKey := generateCacheKey(query, opts)
	filename := filepath.Join(fsc.dir, cacheKey+".json")

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, false
	}

	// Read and parse cache entry
	data, err := os.ReadFile(filename) //nolint:gosec // G304: filename is constructed internally, not from user input
	if err != nil {
		return nil, false
	}

	var entry SearchCacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		// Corrupted cache entry - trigger cleanup to remove it
		go fsc.tryCleanup()
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		// Trigger async cleanup to remove expired entries
		go fsc.tryCleanup()
		return nil, false
	}

	return entry.Results, true
}

// Set stores search results in cache.
func (fsc *FileSearchCache) Set(query string, opts SearchOptions, results []SearchResult, ttl time.Duration) error {
	fsc.mutex.Lock()
	defer fsc.mutex.Unlock()

	// Ensure cache directory exists
	if err := os.MkdirAll(fsc.dir, 0750); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	cacheKey := generateCacheKey(query, opts)
	filename := filepath.Join(fsc.dir, cacheKey+".json")

	entry := SearchCacheEntry{
		Query:     query,
		Results:   results,
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(ttl),
		Hash:      cacheKey,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	tmp, err := os.CreateTemp(fsc.dir, ".cache-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name()) //nolint:errcheck // Best-effort cleanup; a successful rename already removes the temporary path.
	if _, err := tmp.Write(data); err != nil {
		tmp.Close() //nolint:errcheck // Preserve the write error; closing the failed temporary file is best effort.
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp.Name(), filename); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// Clear removes all cached entries.
func (fsc *FileSearchCache) Clear() error {
	fsc.mutex.Lock()
	defer fsc.mutex.Unlock()

	entries, err := os.ReadDir(fsc.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".json") {
			if err := os.Remove(filepath.Join(fsc.dir, entry.Name())); err != nil {
				return fmt.Errorf("failed to remove cache file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// tryCleanup attempts to run cleanup if not already running.
// Uses TryLock to avoid blocking and prevent cleanup storms.
func (fsc *FileSearchCache) tryCleanup() {
	if !fsc.cleanupRunning.TryLock() {
		return // Another cleanup is already running
	}
	defer fsc.cleanupRunning.Unlock()
	_ = fsc.Cleanup() // Ignore error in async cleanup
}

// Cleanup removes expired entries.
func (fsc *FileSearchCache) Cleanup() error {

	entries, err := os.ReadDir(fsc.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		filename := filepath.Join(fsc.dir, entry.Name())
		data, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		var cacheEntry SearchCacheEntry
		if err := json.Unmarshal(data, &cacheEntry); err != nil {
			// Corrupted entry, remove it
			fsc.removeUnchanged(filename, data)
			continue
		}

		if time.Now().After(cacheEntry.ExpiresAt) {
			fsc.removeUnchanged(filename, data)
		}
	}

	return nil
}

// generateCacheKey creates a unique hash for the query and options.
func generateCacheKey(query string, opts SearchOptions) string {
	data, _ := json.Marshal(struct {
		Query   string
		Options SearchOptions
	}{query, opts})
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// CacheStats provides statistics about the cache.
type CacheStats struct {
	TotalEntries   int    `json:"total_entries"`
	ExpiredEntries int    `json:"expired_entries"`
	SizeBytes      int64  `json:"size_bytes"`
	CacheDir       string `json:"cache_dir"`
}

// Stats returns cache statistics.
func (fsc *FileSearchCache) Stats() (*CacheStats, error) {
	fsc.mutex.RLock()
	defer fsc.mutex.RUnlock()

	entries, err := os.ReadDir(fsc.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return &CacheStats{CacheDir: fsc.dir}, nil
		}
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	stats := &CacheStats{
		CacheDir: fsc.dir,
	}

	now := time.Now()

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		stats.TotalEntries++

		// Get file size
		info, err := entry.Info()
		if err == nil {
			stats.SizeBytes += info.Size()
		}

		// Check if expired
		filename := filepath.Join(fsc.dir, entry.Name())
		data, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		var cacheEntry SearchCacheEntry
		if err := json.Unmarshal(data, &cacheEntry); err != nil {
			continue
		}

		if now.After(cacheEntry.ExpiresAt) {
			stats.ExpiredEntries++
		}
	}

	return stats, nil
}

func (fsc *FileSearchCache) removeUnchanged(filename string, original []byte) {
	fsc.mutex.Lock()
	defer fsc.mutex.Unlock()
	current, err := os.ReadFile(filename)
	if err == nil && string(current) == string(original) {
		_ = os.Remove(filename)
	}
}
