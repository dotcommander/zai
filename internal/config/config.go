package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config holds the complete application configuration.
type Config struct {
	API       APIConfig       `mapstructure:"api"`
	WebReader WebReaderConfig `mapstructure:"web_reader"`
	WebSearch WebSearchConfig `mapstructure:"web_search"`
	TTS       TTSConfig       `mapstructure:"tts"`
}

// APIConfig holds API connection settings.
type APIConfig struct {
	Key            string               `mapstructure:"key"`
	BaseURL        string               `mapstructure:"base_url"`
	CodingBaseURL  string               `mapstructure:"coding_base_url"`
	CodingPlan     bool                 `mapstructure:"coding_plan"`
	Model          string               `mapstructure:"model"`
	ImageModel     string               `mapstructure:"image_model"`
	VideoModel     string               `mapstructure:"video_model"`
	VisionModel    string               `mapstructure:"vision_model"`
	AudioModel     string               `mapstructure:"audio_model"`
	TTSModel       string               `mapstructure:"tts_model"`
	EmbeddingModel string               `mapstructure:"embedding_model"`
	RateLimit      RateLimitConfig      `mapstructure:"rate_limit"`
	Retry          RetryConfig          `mapstructure:"retry"`
	CircuitBreaker CircuitBreakerConfig `mapstructure:"circuit_breaker"`
}

// RateLimitConfig holds rate limiting settings.
type RateLimitConfig struct {
	RequestsPerSecond int `mapstructure:"requests_per_second"`
	Burst             int `mapstructure:"burst"`
}

// RetryConfig holds retry behavior settings.
type RetryConfig struct {
	MaxAttempts    int           `mapstructure:"max_attempts"`
	InitialBackoff time.Duration `mapstructure:"initial_backoff"`
	MaxBackoff     time.Duration `mapstructure:"max_backoff"`
}

// CircuitBreakerConfig holds circuit breaker settings.
type CircuitBreakerConfig struct {
	Enabled          bool          `mapstructure:"enabled"`
	FailureThreshold int           `mapstructure:"failure_threshold"`
	SuccessThreshold int           `mapstructure:"success_threshold"`
	Timeout          time.Duration `mapstructure:"timeout"`
}

// WebReaderConfig holds web content fetching settings.
type WebReaderConfig struct {
	Enabled          bool   `mapstructure:"enabled"`
	Timeout          int    `mapstructure:"timeout"`
	CacheEnabled     bool   `mapstructure:"cache_enabled"`
	ReturnFormat     string `mapstructure:"return_format"`
	AutoDetect       bool   `mapstructure:"auto_detect"`
	MaxContentLength int    `mapstructure:"max_content_length"`
}

// WebSearchConfig holds web search settings.
type WebSearchConfig struct {
	Enabled        bool          `mapstructure:"enabled"`
	DefaultCount   int           `mapstructure:"default_count"`
	DefaultRecency string        `mapstructure:"default_recency"`
	Language       string        `mapstructure:"language"`
	Timeout        int           `mapstructure:"timeout"`
	CacheEnabled   bool          `mapstructure:"cache_enabled"`
	CacheDir       string        `mapstructure:"cache_dir"`
	CacheTTL       time.Duration `mapstructure:"cache_ttl"`
	SearchEngine   string        `mapstructure:"search_engine"`
	ContentSize    string        `mapstructure:"content_size"`
}

// TTSConfig holds text-to-speech settings.
type TTSConfig struct {
	Voice          string `mapstructure:"voice"`
	ResponseFormat string `mapstructure:"response_format"`
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
	viper.SetDefault("api.coding_base_url", "https://api.z.ai/api/coding/paas/v4")
	viper.SetDefault("api.coding_plan", false)
	viper.SetDefault("api.model", "glm-4.7")
	viper.SetDefault("api.image_model", "cogview-4-250304")
	viper.SetDefault("api.video_model", "cogvideox-3")
	viper.SetDefault("api.vision_model", "glm-4.6v")
	viper.SetDefault("api.audio_model", "glm-asr-2512")
	viper.SetDefault("api.tts_model", "glm-tts")
	viper.SetDefault("api.embedding_model", "embedding-3")

	// Rate limit defaults
	viper.SetDefault("api.rate_limit.requests_per_second", 10)
	viper.SetDefault("api.rate_limit.burst", 5)

	// Retry defaults
	viper.SetDefault("api.retry.max_attempts", 3)
	viper.SetDefault("api.retry.initial_backoff", "1s")
	viper.SetDefault("api.retry.max_backoff", "30s")

	// Circuit breaker defaults
	viper.SetDefault("api.circuit_breaker.enabled", true)
	viper.SetDefault("api.circuit_breaker.failure_threshold", 5)
	viper.SetDefault("api.circuit_breaker.success_threshold", 2)
	viper.SetDefault("api.circuit_breaker.timeout", "60s")

	// Web reader defaults
	viper.SetDefault("web_reader.enabled", true)
	viper.SetDefault("web_reader.timeout", 20)
	viper.SetDefault("web_reader.cache_enabled", true)
	viper.SetDefault("web_reader.return_format", "markdown")
	viper.SetDefault("web_reader.auto_detect", true)
	viper.SetDefault("web_reader.max_content_length", 50000)

	// Web search defaults
	home, err := os.UserHomeDir()
	if err != nil {
		home = "/tmp" // Fallback for cache dir
	}
	viper.SetDefault("web_search.enabled", true)
	viper.SetDefault("web_search.default_count", 10)
	viper.SetDefault("web_search.default_recency", "noLimit")
	viper.SetDefault("web_search.language", "en")
	viper.SetDefault("web_search.timeout", 30)
	viper.SetDefault("web_search.cache_enabled", true)
	viper.SetDefault("web_search.cache_dir", filepath.Join(home, ".config", "zai", "search_cache"))
	viper.SetDefault("web_search.cache_ttl", "24h")
	viper.SetDefault("web_search.search_engine", "search_std")
	viper.SetDefault("web_search.content_size", "medium")

	// TTS defaults
	viper.SetDefault("tts.voice", "tongtong")
	viper.SetDefault("tts.response_format", "wav")
}
