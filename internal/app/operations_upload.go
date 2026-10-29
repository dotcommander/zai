package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// operationMultipart assembles framing in memory, but reads the file directly
// from disk. No producer goroutine can outlive a rejected/canceled request.
func (c *Client) operationMultipart(ctx context.Context, op Operation, input OperationRequest) (io.Reader, string, func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, "", nil, err
	}
	fields := make(map[string]string)
	if len(input.JSON) > 0 {
		obj, err := JSONObject(input.JSON)
		if err != nil {
			return nil, "", nil, err
		}
		for key, value := range obj {
			var str string
			if json.Unmarshal(value, &str) == nil {
				fields[key] = str
			} else {
				fields[key] = string(value)
			}
		}
	}
	for key, value := range input.Fields {
		fields[key] = value
	}
	if _, ok := fields["file"]; ok {
		return nil, "", nil, errors.New("binary file must be supplied as UploadPath, not a JSON field")
	}
	if op.Name == "files_upload" && fields["purpose"] == "" {
		fields["purpose"] = "user_data"
	}
	if op.Name == "files_upload" && input.UploadPath == "" {
		return nil, "", nil, errors.New("file upload requires a local file")
	}
	if op.Name == "audio" && input.UploadPath == "" && fields["file_base64"] == "" {
		return nil, "", nil, errors.New("transcription requires a file or file_base64")
	}
	var file *os.File
	cleanup := func() {
		if file != nil {
			_ = file.Close()
		}
	}
	if input.UploadPath != "" {
		var err error
		file, err = os.Open(input.UploadPath)
		if err != nil {
			return nil, "", nil, fmt.Errorf("open upload: %w", err)
		}
		info, err := file.Stat()
		if err != nil {
			cleanup()
			return nil, "", nil, err
		}
		if !info.Mode().IsRegular() {
			cleanup()
			return nil, "", nil, errors.New("upload must be a regular file")
		}
		if err = validateOperationUpload(op.Name, fields["purpose"], input.UploadPath, info.Size()); err != nil {
			cleanup()
			return nil, "", nil, err
		}
	}
	var framing bytes.Buffer
	writer := multipart.NewWriter(&framing)
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if err := writer.WriteField(key, fields[key]); err != nil {
			cleanup()
			return nil, "", nil, err
		}
	}
	if file == nil {
		if err := writer.Close(); err != nil {
			cleanup()
			return nil, "", nil, err
		}
		return bytes.NewReader(framing.Bytes()), writer.FormDataContentType(), cleanup, nil
	}
	if _, err := writer.CreateFormFile("file", filepath.Base(input.UploadPath)); err != nil {
		cleanup()
		return nil, "", nil, err
	}
	split := framing.Len()
	if err := writer.Close(); err != nil {
		cleanup()
		return nil, "", nil, err
	}
	data := framing.Bytes()
	return io.MultiReader(bytes.NewReader(data[:split]), file, bytes.NewReader(data[split:])), writer.FormDataContentType(), cleanup, nil
}

func validateOperationUpload(operation, purpose, path string, size int64) error {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	allowed := []string{}
	var limit int64
	switch operation {
	case "audio":
		allowed = []string{"wav", "mp3"}
		limit = 25 << 20
	case "files_upload":
		switch purpose {
		case "agent":
			allowed = []string{"pdf", "doc", "xlsx", "ppt", "txt", "jpg", "png"}
			limit = 100 << 20
		case "user_data":
			allowed = []string{"pptx", "ppt", "docx", "doc", "xlsx", "xls", "pdf"}
			limit = 1 << 40
		default:
			return fmt.Errorf("unsupported upload purpose %q (use agent or user_data)", purpose)
		}
	}
	if limit > 0 && size > limit {
		return fmt.Errorf("%s upload exceeds purpose-specific limit of %d bytes", operation, limit)
	}
	for _, candidate := range allowed {
		if ext == candidate {
			return nil
		}
	}
	return fmt.Errorf("unsupported %s file extension %q for purpose %q", operation, ext, purpose)
}
