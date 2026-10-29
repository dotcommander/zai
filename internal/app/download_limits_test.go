package app

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type downloadDoerFunc func(*http.Request) (*http.Response, error)

func (f downloadDoerFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

type brokenDownloadReader struct{}

func (brokenDownloadReader) Read([]byte) (int, error) { return 0, errors.New("broken transfer") }

func TestMediaDownloadLimitAndAtomicPublication(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name      string
		body      string
		length    int64
		broken    bool
		wantError bool
	}{
		{"exact known", "abcd", 4, false, false},
		{"exact unknown", "abcd", -1, false, false},
		{"oversized known", "abcde", 5, false, true},
		{"oversized unknown", "abcde", -1, false, true},
		{"misleading length", "abcde", 2, false, true},
		{"read failure", "", -1, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			destination := filepath.Join(dir, "image.png")
			if err := os.WriteFile(destination, []byte("original"), 0600); err != nil {
				t.Fatal(err)
			}
			doer := downloadDoerFunc(func(r *http.Request) (*http.Response, error) {
				var body io.Reader = strings.NewReader(tc.body)
				if tc.broken {
					body = brokenDownloadReader{}
				}
				return &http.Response{StatusCode: http.StatusOK, ContentLength: tc.length, Body: io.NopCloser(body)}, nil
			})
			result := NewMediaDownloader(doer, 4).Download(context.Background(), "https://example.com/image", destination)
			if (result.Error != nil) != tc.wantError {
				t.Fatalf("error = %v", result.Error)
			}
			data, err := os.ReadFile(destination)
			if err != nil {
				t.Fatal(err)
			}
			want := tc.body
			if tc.wantError {
				want = "original"
			} else if result.Size != 4 {
				t.Fatalf("size = %d", result.Size)
			}
			if string(data) != want {
				t.Fatalf("destination = %q, want %q", data, want)
			}
			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 1 {
				t.Fatalf("partial files left: %v", entries)
			}
		})
	}
}

func TestMediaDownloadRejectsNonpositiveLimitBeforeHTTP(t *testing.T) {
	t.Parallel()
	doer := downloadDoerFunc(func(*http.Request) (*http.Response, error) { t.Fatal("unexpected request"); return nil, nil })
	for _, limit := range []int64{0, -1} {
		result := NewMediaDownloader(doer, limit).Download(context.Background(), "https://example.com", filepath.Join(t.TempDir(), "output"))
		if result.Error == nil {
			t.Fatalf("limit %d accepted", limit)
		}
	}
}
