package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	API APIConfig `mapstructure:"api"`
}

type APIConfig struct {
	Key     string `mapstructure:"key"`
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
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
	viper.SetDefault("api.model", "glm-4.6")
}