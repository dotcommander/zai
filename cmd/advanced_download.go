package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dotcommander/zai/internal/app"
)

func addAdvancedDownloadFlags(cmd *cobra.Command) {
	cmd.Flags().String("download", "", "Download the first result artifact atomically to this file")
	if cmd.Flag("max-download-bytes") == nil {
		cmd.Flags().Int64("max-download-bytes", app.DefaultMaxDownloadBytes, "Maximum artifact download size")
	}
}

func advancedDownloadTarget(cmd *cobra.Command) string {
	target, _ := cmd.Flags().GetString("download")
	if target == "" {
		if f := cmd.Flag("output"); f != nil && f.Changed {
			target = f.Value.String()
		}
	}
	return target
}

func validateAdvancedDownload(cmd *cobra.Command, operation string, wait, events bool) error {
	if advancedDownloadTarget(cmd) == "" {
		return nil
	}
	switch operation {
	case "images", "image_async", "videos", "task_get", "agent_run", "agent_result", "agent_conversation":
	default:
		return fmt.Errorf("artifact download is unavailable for %s", operation)
	}
	if events {
		return fmt.Errorf("artifact download requires a completed JSON result; omit --events")
	}
	if (operation == "videos" || operation == "image_async" || operation == "agent_run") && !wait {
		return fmt.Errorf("artifact download of a submitted task requires --wait and --wait-timeout")
	}
	limit, _ := cmd.Flags().GetInt64("max-download-bytes")
	if limit <= 0 {
		return fmt.Errorf("--max-download-bytes must be positive")
	}
	return nil
}

func downloadAdvancedResult(ctx context.Context, cmd *cobra.Command, client advancedAPI, result json.RawMessage) error {
	target := advancedDownloadTarget(cmd)
	if target == "" {
		return nil
	}
	var value any
	if err := json.Unmarshal(result, &value); err != nil {
		return err
	}
	artifact := firstArtifactURL(value)
	if artifact == "" {
		return fmt.Errorf("completed response has no HTTP(S) artifact URL")
	}
	limit, _ := cmd.Flags().GetInt64("max-download-bytes")
	var httpClient app.HTTPDoer
	if transport, ok := client.(interface{ HTTPClient() app.HTTPDoer }); ok {
		httpClient = transport.HTTPClient()
	}
	download := app.NewMediaDownloader(httpClient, limit).Download(ctx, artifact, target)
	if download.Error != nil {
		return fmt.Errorf("save result artifact: %w", download.Error)
	}
	return nil
}

func firstArtifactURL(value any) string {
	switch value := value.(type) {
	case []any:
		for _, child := range value {
			if found := firstArtifactURL(child); found != "" {
				return found
			}
		}
	case map[string]any:
		for _, key := range []string{"download_url", "file_url", "video_url", "image_url", "pptx_url", "pdf_url", "url"} {
			if candidate, ok := value[key].(string); ok {
				if validArtifactURL(candidate) {
					return candidate
				}
			}
		}
		keys := make([]string, 0, len(value))
		for key := range value {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if found := firstArtifactURL(value[key]); found != "" {
				return found
			}
		}
	case string:
		if validArtifactURL(value) {
			return value
		}
	}
	return ""
}
func validArtifactURL(value string) bool {
	if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		return false
	}
	parsed, err := url.Parse(value)
	return err == nil && parsed.Host != "" && parsed.User == nil
}
