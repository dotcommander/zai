package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

//go:embed quota-defaults.yaml
var quotaDefaults []byte

// QuotaConfig holds operator-adjustable bounds and presentation text.
type QuotaConfig struct {
	Timeout          time.Duration     `mapstructure:"timeout"`
	MaxResponseBytes int64             `mapstructure:"max_response_bytes"`
	Route            string            `mapstructure:"route"`
	ApprovedHosts    []string          `mapstructure:"approved_hosts"`
	CommandHelp      string            `mapstructure:"command_help"`
	Heading          string            `mapstructure:"heading"`
	Unavailable      string            `mapstructure:"unavailable"`
	Labels           map[string]string `mapstructure:"labels"`
}

// QuotaSafety returns the pinned credential destination policy from the embedded asset.
func QuotaSafety() (string, []string, error) {
	defaults, err := quotaAsset()
	if err != nil {
		return "", nil, err
	}
	route := defaults.GetString("quota.route")
	hosts := defaults.GetStringSlice("quota.approved_hosts")
	if route == "" || len(hosts) == 0 {
		return "", nil, fmt.Errorf("invalid quota destination policy")
	}
	return route, hosts, nil
}

func quotaAsset() (*viper.Viper, error) {
	defaults := viper.New()
	defaults.SetConfigType("yaml")
	if err := defaults.ReadConfig(bytes.NewReader(quotaDefaults)); err != nil {
		return nil, fmt.Errorf("invalid bundled defaults: %w", err)
	}
	return defaults, nil
}

// SetQuotaDefaults supplies missing settings without modifying custom config files.
func SetQuotaDefaults() {
	defaults, err := quotaAsset()
	if err != nil {
		return
	}
	for _, key := range []string{"timeout", "max_response_bytes", "command_help", "heading", "unavailable", "labels"} {
		viper.SetDefault("quota."+key, defaults.Get("quota."+key))
	}
	for _, key := range []string{"badge", "bar_width", "filled", "empty", "percent_mark", "group_separator", "duration_hour", "duration_minute", "reset_format", "indent", "gap", "types.time_limit", "types.tokens_limit", "types.credit_limit"} {
		viper.SetDefault("quota.display."+key, defaults.Get("quota.display."+key))
	}
}

// BootstrapDefault creates a private default config only when absent.
func BootstrapDefault(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if os.IsExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("create config: %w", err)
	}
	_, writeErr := f.Write(quotaDefaults)
	closeErr := f.Close()
	if writeErr != nil {
		return fmt.Errorf("write config: %w", writeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close config: %w", closeErr)
	}
	return nil
}
