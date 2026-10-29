package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/app"
)

// mediaDownloadLimit keeps explicit command flags above environment/config values.
func mediaDownloadLimit(cmd *cobra.Command, overrides ...int64) (int64, error) {
	limit := app.DefaultMaxDownloadBytes
	if viper.IsSet("media.max_download_bytes") {
		limit = viper.GetInt64("media.max_download_bytes")
	}
	if cmd != nil && cmd.Flags().Changed("max-download-bytes") {
		var err error
		limit, err = cmd.Flags().GetInt64("max-download-bytes")
		if err != nil {
			return 0, err
		}
	}
	if len(overrides) > 0 {
		limit = overrides[0]
	}
	if limit <= 0 {
		return 0, fmt.Errorf("media.max_download_bytes / --max-download-bytes must be positive")
	}
	return limit, nil
}
