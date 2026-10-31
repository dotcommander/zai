package fileutil

import (
	"os"
	"sync"
)

// SessionFileReader reuses one local file while its observed metadata is unchanged.
// The injected reader remains the source of bytes, including when stat is unavailable.
type SessionFileReader struct {
	reader FileReader
	path   string
	mu     sync.Mutex
	info   os.FileInfo
	data   []byte
}

func NewSessionFileReader(reader FileReader, path string) *SessionFileReader {
	return &SessionFileReader{reader: reader, path: path}
}

func sameSessionFile(a, b os.FileInfo) bool {
	return a != nil && b != nil && os.SameFile(a, b) && a.Size() == b.Size() &&
		a.ModTime().Equal(b.ModTime()) && a.Mode() == b.Mode()
}

func (r *SessionFileReader) ReadFile(name string) ([]byte, error) {
	if name != r.path {
		return r.reader.ReadFile(name)
	}

	// Serialize fills to avoid duplicating disk reads when requests overlap.
	r.mu.Lock()
	defer r.mu.Unlock()

	before, statErr := os.Stat(name)
	if statErr == nil && sameSessionFile(r.info, before) {
		return r.data, nil
	}
	r.info, r.data = nil, nil

	data, err := r.reader.ReadFile(name)
	if err != nil {
		return nil, err
	}
	if statErr == nil {
		after, err := os.Stat(name)
		if err == nil && sameSessionFile(before, after) {
			r.info, r.data = after, data
		}
	}
	return data, nil
}
