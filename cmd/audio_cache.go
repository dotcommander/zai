package cmd

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// AudioCache stores partial transcription results for resume support.
type AudioCache struct {
	Chunks map[int]string `json:"chunks"`
}

type audioChunkRecord struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
}

const maxAudioJournalRecord = 10 * 1024 * 1024

// The cache identity remains the first eight bytes of the source SHA-256.
func getCachePath(sourcePath string) (string, error) {
	file, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer closeFile(file)
	info, err := file.Stat()
	if err != nil {
		return "", err
	}
	if info.Size() > MaxFileSize {
		return "", fmt.Errorf("audio cache file exceeds maximum size of %d bytes", MaxFileSize)
	}
	hash := sha256.New()
	n, err := io.Copy(hash, io.LimitReader(file, MaxFileSize+1))
	if err != nil {
		return "", err
	}
	if n > MaxFileSize {
		return "", fmt.Errorf("audio cache file exceeds maximum size of %d bytes", MaxFileSize)
	}
	cacheDir := filepath.Join(os.Getenv("HOME"), ".cache", "zai", "audio")
	if err := os.MkdirAll(cacheDir, 0700); err != nil { //nolint:gosec // G703: HOME determines cache location
		return "", err
	}
	return filepath.Join(cacheDir, hex.EncodeToString(hash.Sum(nil)[:8])+".json"), nil
}

func audioJournalPath(cachePath string) string { return cachePath + ".jsonl" }

// loadCache overlays completed journal records on the older snapshot format.
func loadCache(cachePath string) (*AudioCache, error) {
	cache := &AudioCache{Chunks: make(map[int]string)}
	data, err := os.ReadFile(cachePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err == nil {
		if err := json.Unmarshal(data, cache); err != nil {
			return nil, err
		}
		if cache.Chunks == nil {
			cache.Chunks = make(map[int]string)
		}
	}
	file, err := os.Open(audioJournalPath(cachePath))
	if os.IsNotExist(err) {
		return cache, nil
	}
	if err != nil {
		return nil, err
	}
	defer closeFile(file)
	reader := bufio.NewReader(file)
	for {
		line, readErr := reader.ReadSlice('\n')
		if readErr == bufio.ErrBufferFull {
			// A record may exceed the reader's buffer; cap it before decoding.
			line, readErr = readAudioJournalRecord(reader, line)
		}
		if readErr == io.EOF {
			return cache, nil // Ignore an incomplete final append.
		}
		if readErr != nil {
			return nil, readErr
		}
		if len(line) > maxAudioJournalRecord {
			continue
		}
		var record audioChunkRecord
		if json.Unmarshal(line, &record) == nil && record.Index >= 0 {
			cache.Chunks[record.Index] = record.Text
		}
	}
}

// Consume the rest of an oversized record without retaining it.
func readAudioJournalRecord(reader *bufio.Reader, first []byte) ([]byte, error) {
	line := append([]byte(nil), first...)
	oversized := false
	for {
		fragment, err := reader.ReadSlice('\n')
		if !oversized {
			if len(line)+len(fragment) <= maxAudioJournalRecord {
				line = append(line, fragment...)
			} else {
				line = nil
				oversized = true
			}
		}
		if err == bufio.ErrBufferFull {
			continue
		}
		return line, err
	}
}

func appendAudioChunk(cachePath string, index int, value string) error {
	data, err := json.Marshal(audioChunkRecord{Index: index, Text: value})
	if err != nil {
		return err
	}
	if len(data)+1 > maxAudioJournalRecord {
		return fmt.Errorf("audio chunk cache record exceeds %d bytes", maxAudioJournalRecord)
	}
	file, err := os.OpenFile(audioJournalPath(cachePath), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer closeFile(file)
	info, err := file.Stat()
	if err != nil {
		return err
	}
	prefix := []byte(nil)
	if info.Size() > 0 {
		var last [1]byte
		if _, err := file.ReadAt(last[:], info.Size()-1); err != nil {
			return err
		}
		if last[0] != '\n' {
			prefix = []byte{'\n'} // Isolate a torn prior append from the next record.
		}
	}
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	entry := make([]byte, 0, len(prefix)+len(data)+1)
	entry = append(entry, prefix...)
	entry = append(entry, data...)
	entry = append(entry, '\n')
	if _, err := file.Write(entry); err != nil {
		return err
	}
	return file.Sync()
}

// Publish a complete snapshot before discarding the journal. A failed rename or
// removal leaves the journal available for replay on the next invocation.
func compactAudioCache(cachePath string, cache *AudioCache) error {
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(cachePath), ".audio-cache-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name()) //nolint:errcheck // renamed files have no temporary path
	if err := file.Chmod(0600); err != nil {
		return errors.Join(err, file.Close())
	}
	if _, err := file.Write(data); err != nil {
		return errors.Join(err, file.Close())
	}
	if err := file.Sync(); err != nil {
		return errors.Join(err, file.Close())
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Rename(file.Name(), cachePath); err != nil {
		return err
	}
	if err := os.Remove(audioJournalPath(cachePath)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func clearAudioCache(cachePath string) error {
	var errs []error
	for _, path := range []string{cachePath, audioJournalPath(cachePath)} {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
