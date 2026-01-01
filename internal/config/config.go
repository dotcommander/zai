package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	API       APIConfig       `mapstructure:"api"`
	WebReader WebReaderConfig `mapstructure:"web_reader"`
	WebSearch WebSearchConfig `mapstructure:"web_search"`
}

type APIConfig struct {
	Key        string      `mapstructure:"key"`
	BaseURL    string      `mapstructure:"base_url"`
	Model      string      `mapstructure:"model"`
	ImageModel string      `mapstructure:"image_model"`
	Retry      RetryConfig `mapstructure:"retry"`
}

type RetryConfig struct {
	MaxAttempts    int           `mapstructure:"max_attempts"`
	InitialBackoff time.Duration `mapstructure:"initial_backoff"`
	MaxBackoff     time.Duration `mapstructure:"max_backoff"`
}

type WebReaderConfig struct {
	Enabled          bool   `mapstructure:"enabled"`
	Timeout          int    `mapstructure:"timeout"`
	CacheEnabled     bool   `mapstructure:"cache_enabled"`
	ReturnFormat     string `mapstructure:"return_format"`
	AutoDetect       bool   `mapstructure:"auto_detect"`
	MaxContentLength int    `mapstructure:"max_content_length"`
}

type WebSearchConfig struct {
	Enabled          bool          `mapstructure:"enabled"`
	DefaultCount     int           `mapstructure:"default_count"`
	DefaultRecency   string        `mapstructure:"default_recency"`
	Timeout          int           `mapstructure:"timeout"`
	CacheEnabled     bool          `mapstructure:"cache_enabled"`
	CacheDir         string        `mapstructure:"cache_dir"`
	CacheTTL         time.Duration `mapstructure:"cache_ttl"`
}

// Load unmarshals viper config into struct
func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}
	return &cfg, nil
}

// SetDefaults sets default values
func SetDefaults() {
	viper.SetDefault("api.base_url", "https://api.z.ai/api/paas/v4")
	viper.SetDefault("api.model", "glm-4.7")
	viper.SetDefault("api.image_model", "cogview-4-250304")

	// Retry defaults
	viper.SetDefault("api.retry.max_attempts", 3)
	viper.SetDefault("api.retry.initial_backoff", "1s")
	viper.SetDefault("api.retry.max_backoff", "30s")

	// Web reader defaults
	viper.SetDefault("web_reader.enabled", true)
	viper.SetDefault("web_reader.timeout", 20)
	viper.SetDefault("web_reader.cache_enabled", true)
	viper.SetDefault("web_reader.return_format", "markdown")
	viper.SetDefault("web_reader.auto_detect", true)
	viper.SetDefault("web_reader.max_content_length", 50000)

	// Web search defaults
	home, _ := os.UserHomeDir()
	viper.SetDefault("web_search.enabled", true)
	viper.SetDefault("web_search.default_count", 10)
	viper.SetDefault("web_search.default_recency", "noLimit")
	viper.SetDefault("web_search.timeout", 30)
	viper.SetDefault("web_search.cache_enabled", true)
	viper.SetDefault("web_search.cache_dir", filepath.Join(home, ".config", "zai", "search_cache"))
	viper.SetDefault("web_search.cache_ttl", "24h")
}