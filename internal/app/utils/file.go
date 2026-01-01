package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileReader interface for reading files (ISP compliance, enables testing).
type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

// OSFileReader implements FileReader using os.ReadFile.
type OSFileReader struct{}

// ReadFile reads a file from the filesystem.
func (r OSFileReader) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// DetectImageMimeType determines the MIME type from file extension.
func DetectImageMimeType(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg", nil
	case ".png":
		return "image/png", nil
	case ".gif":
		return "image/gif", nil
	case ".webp":
		return "image/webp", nil
	default:
		return "", fmt.Errorf("unsupported image format: %s (supported: jpg, jpeg, png, gif, webp)", ext)
	}
}

// EncodeBytesToDataURI encodes bytes to base64 data URI.
func EncodeBytesToDataURI(data []byte, mimeType string) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)
}

// ValidateImageFile checks if an image file exists and is readable.
func ValidateImageFile(filePath string, fileReader FileReader) error {
	if _, err := fileReader.ReadFile(filePath); err != nil {
		return fmt.Errorf("image file validation failed: %w", err)
	}
	return nil
}