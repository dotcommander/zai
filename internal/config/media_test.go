package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestMediaAndEmbeddingDefaults(t *testing.T) {
	// Serial: Viper is process-global state.
	viper.Reset()
	t.Cleanup(viper.Reset)
	SetDefaults()
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Media.MaxDownloadBytes != 536870912 {
		t.Fatalf("media limit = %d", cfg.Media.MaxDownloadBytes)
	}
	if cfg.API.EmbeddingModel != "embedding-3" {
		t.Fatalf("embedding model = %q", cfg.API.EmbeddingModel)
	}
	viper.Set("media.max_download_bytes", 0)
	if _, err := Load(); err == nil {
		t.Fatal("accepted zero media limit")
	}
	viper.Set("media.max_download_bytes", -1)
	if _, err := Load(); err == nil {
		t.Fatal("accepted negative media limit")
	}
}
