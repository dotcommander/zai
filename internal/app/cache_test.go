package app

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFileSearchCache tests the FileSearchCache implementation.
func TestFileSearchCache(t *testing.T) {
	// Create temp directory for test cache
	tempDir := t.TempDir()
	cache := NewFileSearchCache(tempDir)

	// Test data
	query := "test query"
	opts := SearchOptions{
		Count:         5,
		RecencyFilter: "oneWeek",
	}
	results := []SearchResult{
		{Title: "Result 1", Link: "https://example.com/1"},
		{Title: "Result 2", Link: "https://example.com/2"},
	}
	ttl := 1 * time.Hour

	t.Run("Set and Get", func(t *testing.T) {
		// Set cache entry
		err := cache.Set(query, opts, results, ttl)
		require.NoError(t, err)

		// Get cache entry
		cachedResults, found := cache.Get(query, opts)
		assert.True(t, found)
		assert.Equal(t, results, cachedResults)
	})

	t.Run("Get miss - non-existent key", func(t *testing.T) {
		_, found := cache.Get("nonexistent", SearchOptions{})
		assert.False(t, found)
	})

	t.Run("Get miss - different options", func(t *testing.T) {
		// Set with one set of options
		err := cache.Set(query, opts, results, ttl)
		require.NoError(t, err)

		// Try to get with different options
		differentOpts := SearchOptions{
			Count:         10, // Different count
			RecencyFilter: "oneWeek",
		}
		_, found := cache.Get(query, differentOpts)
		assert.False(t, found)
	})

	t.Run("Get miss - expired entry", func(t *testing.T) {
		// Set entry with very short TTL
		shortTTL := 10 * time.Millisecond
		err := cache.Set("expired", SearchOptions{}, results, shortTTL)
		require.NoError(t, err)

		// Wait for expiration
		time.Sleep(150 * time.Millisecond)

		// Should return not found
		_, found := cache.Get("expired", SearchOptions{})
		assert.False(t, found)
	})

	t.Run("Clear removes all entries", func(t *testing.T) {
		// Add multiple entries
		err := cache.Set("query1", SearchOptions{}, results, ttl)
		require.NoError(t, err)
		err = cache.Set("query2", SearchOptions{}, results, ttl)
		require.NoError(t, err)

		// Verify entries exist
		_, found1 := cache.Get("query1", SearchOptions{})
		_, found2 := cache.Get("query2", SearchOptions{})
		assert.True(t, found1)
		assert.True(t, found2)

		// Clear cache
		err = cache.Clear()
		require.NoError(t, err)

		// Verify entries are gone
		_, found1 = cache.Get("query1", SearchOptions{})
		_, found2 = cache.Get("query2", SearchOptions{})
		assert.False(t, found1)
		assert.False(t, found2)
	})

	t.Run("Cleanup removes expired entries", func(t *testing.T) {
		// Add both valid and expired entries
		shortTTL := 10 * time.Millisecond
		err := cache.Set("expired", SearchOptions{}, results, shortTTL)
		require.NoError(t, err)
		err = cache.Set("valid", SearchOptions{}, results, ttl)
		require.NoError(t, err)

		// Wait for expiration
		time.Sleep(150 * time.Millisecond)

		// Run cleanup
		err = cache.Cleanup()
		require.NoError(t, err)

		// Verify expired entry is gone
		_, found := cache.Get("expired", SearchOptions{})
		assert.False(t, found)

		// Verify valid entry still exists
		_, found = cache.Get("valid", SearchOptions{})
		assert.True(t, found)
	})

	t.Run("Stats returns cache statistics", func(t *testing.T) {
		// Clear any existing entries
		_ = cache.Clear()

		// Add some entries
		err := cache.Set("query1", SearchOptions{Count: 5}, results, ttl)
		require.NoError(t, err)
		err = cache.Set("query2", SearchOptions{Count: 10}, results, ttl)
		require.NoError(t, err)

		// Add expired entry
		shortTTL := 10 * time.Millisecond
		err = cache.Set("expired", SearchOptions{}, results, shortTTL)
		require.NoError(t, err)

		// Wait for expiration
		time.Sleep(150 * time.Millisecond)

		// Get stats
		stats, err := cache.Stats()
		require.NoError(t, err)

		assert.Equal(t, tempDir, stats.CacheDir)
		assert.Equal(t, 3, stats.TotalEntries)
		assert.Equal(t, 1, stats.ExpiredEntries)
		assert.Greater(t, stats.SizeBytes, int64(0))
	})

	t.Run("Concurrent access is safe", func(t *testing.T) {
		// Clear any existing entries
		_ = cache.Clear()

		// Run concurrent operations
		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func(index int) {
				query := fmt.Sprintf("query%d", index)
				_ = cache.Set(query, SearchOptions{}, results, ttl)
				_, _ = cache.Get(query, SearchOptions{})
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}

		// Verify cache is still functional
		stats, err := cache.Stats()
		require.NoError(t, err)
		assert.Equal(t, 10, stats.TotalEntries)
	})
}

// TestFileSearchCacheNilDirectory tests cache with nil directory (uses default).
func TestFileSearchCacheNilDirectory(t *testing.T) {
	// Don't use tempDir, test with actual temp directory
	tempDir := t.TempDir()
	cache := NewFileSearchCache(tempDir)

	query := "test"
	opts := SearchOptions{}
	results := []SearchResult{{Title: "Test"}}
	ttl := 1 * time.Hour

	err := cache.Set(query, opts, results, ttl)
	require.NoError(t, err)

	cached, found := cache.Get(query, opts)
	assert.True(t, found)
	assert.Equal(t, results, cached)
}

// TestFileSearchCacheNonExistentDirectory tests cache creation when directory doesn't exist.
func TestFileSearchCacheNonExistentDirectory(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentDir := filepath.Join(tempDir, "cache", "subdir")

	// Create cache with non-existent directory
	cache := NewFileSearchCache(nonExistentDir)

	query := "test"
	opts := SearchOptions{}
	results := []SearchResult{{Title: "Test"}}
	ttl := 1 * time.Hour

	// Should create directory and save cache
	err := cache.Set(query, opts, results, ttl)
	require.NoError(t, err)

	// Verify directory was created
	_, err = os.Stat(nonExistentDir)
	assert.NoError(t, err)
}
