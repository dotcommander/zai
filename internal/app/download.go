package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// closeBodyResponse closes the response body and logs any error.
func closeBodyResponse(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to close response body: %v\n", err)
	}
}

// DefaultMaxDownloadBytes is the default bound for media transfers.
const DefaultMaxDownloadBytes int64 = 536870912

// MediaDownloader handles downloading media files with DI support.
type MediaDownloader struct {
	httpClient HTTPDoer
	maxBytes   int64
}

// NewMediaDownloader creates a MediaDownloader with the provided HTTP client.
// If httpClient is nil, a default http.Client is used.
func NewMediaDownloader(httpClient HTTPDoer, maxBytes ...int64) *MediaDownloader {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	limit := DefaultMaxDownloadBytes
	if len(maxBytes) > 0 {
		limit = maxBytes[0]
	}
	return &MediaDownloader{httpClient: httpClient, maxBytes: limit}
}

// DownloadResult contains the result of a download operation.
type DownloadResult struct {
	FilePath string
	Size     int64
	Error    error
}

// Download fetches a URL and saves to file with directory creation.
func (d *MediaDownloader) Download(ctx context.Context, url, filePath string) *DownloadResult {
	if d.maxBytes <= 0 {
		return &DownloadResult{FilePath: filePath, Error: fmt.Errorf("max download bytes must be positive")}
	}
	if err := ensureDir(filePath); err != nil {
		return &DownloadResult{FilePath: filePath, Error: err}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return &DownloadResult{FilePath: filePath, Error: fmt.Errorf("create request: %w", err)}
	}
	req.Header.Set("User-Agent", "zai/1.0")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return &DownloadResult{FilePath: filePath, Error: fmt.Errorf("download: %w", err)}
	}
	defer closeBodyResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return &DownloadResult{FilePath: filePath, Error: fmt.Errorf("download failed: status %d", resp.StatusCode)}
	}

	if resp.ContentLength > d.maxBytes {
		return &DownloadResult{FilePath: filePath, Error: fmt.Errorf("download exceeds maximum of %d bytes", d.maxBytes)}
	}

	size, err := writeToFile(filePath, resp.Body, d.maxBytes)
	if err != nil {
		return &DownloadResult{FilePath: filePath, Error: err}
	}

	return &DownloadResult{FilePath: filePath, Size: size, Error: nil}
}

// ensureDir creates the parent directory for a file path if needed.
func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("create directory: %w", err)
		}
	}
	return nil
}

// writeToFile publishes a completed, bounded transfer atomically.
func writeToFile(filePath string, r io.Reader, maxBytes int64) (int64, error) {
	out, err := os.CreateTemp(filepath.Dir(filePath), "."+filepath.Base(filePath)+"-*")
	if err != nil {
		return 0, fmt.Errorf("create temporary file: %w", err)
	}
	defer os.Remove(out.Name()) //nolint:errcheck // Best-effort cleanup; a successful rename already removes the temporary path.
	defer out.Close()           //nolint:errcheck // Failure paths preserve their primary error; the successful path explicitly checks Close below.

	size, err := io.Copy(out, io.LimitReader(r, maxBytes))
	if err != nil {
		return 0, fmt.Errorf("write file: %w", err)
	}
	// Read one extra byte separately to avoid overflowing maxBytes + 1.
	var extra [1]byte
	n, err := io.ReadFull(r, extra[:])
	if n > 0 {
		return 0, fmt.Errorf("download exceeds maximum of %d bytes", maxBytes)
	}
	if err != nil && err != io.EOF {
		return 0, fmt.Errorf("read download: %w", err)
	}
	if err := out.Close(); err != nil {
		return 0, fmt.Errorf("close file: %w", err)
	}
	if err := os.Rename(out.Name(), filePath); err != nil {
		return 0, fmt.Errorf("publish file: %w", err)
	}
	return size, nil
}
