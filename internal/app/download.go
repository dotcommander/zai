package app

import (
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

// MediaDownloader handles downloading media files with DI support.
type MediaDownloader struct {
	httpClient HTTPDoer
}

// NewMediaDownloader creates a MediaDownloader with the provided HTTP client.
// If httpClient is nil, a default http.Client is used.
func NewMediaDownloader(httpClient HTTPDoer) *MediaDownloader {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &MediaDownloader{httpClient: httpClient}
}

// DownloadResult contains the result of a download operation.
type DownloadResult struct {
	FilePath string
	Size     int64
	Error    error
}

// Download fetches a URL and saves to file with directory creation.
func (d *MediaDownloader) Download(url, filePath string) *DownloadResult {
	if err := ensureDir(filePath); err != nil {
		return &DownloadResult{FilePath: filePath, Error: err}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &DownloadResult{FilePath: filePath, Error: fmt.Errorf("create request: %w", err)}
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return &DownloadResult{FilePath: filePath, Error: fmt.Errorf("download: %w", err)}
	}
	defer closeBodyResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return &DownloadResult{FilePath: filePath, Error: fmt.Errorf("download failed: status %d", resp.StatusCode)}
	}

	size, err := writeToFile(filePath, resp.Body)
	if err != nil {
		return &DownloadResult{FilePath: filePath, Error: err}
	}

	return &DownloadResult{FilePath: filePath, Size: size, Error: nil}
}

// ensureDir creates the parent directory for a file path if needed.
func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory: %w", err)
		}
	}
	return nil
}

// writeToFile writes reader content to a file and returns bytes written.
func writeToFile(filePath string, r io.Reader) (int64, error) {
	out, err := os.Create(filePath)
	if err != nil {
		return 0, fmt.Errorf("create file: %w", err)
	}
	defer closeFile(out)

	size, err := io.Copy(out, r)
	if err != nil {
		return 0, fmt.Errorf("write file: %w", err)
	}

	return size, nil
}
