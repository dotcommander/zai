package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestQuotaDefaults(t *testing.T) {
	// Serial: Viper is process-global state.
	viper.Reset()
	t.Cleanup(viper.Reset)
	SetDefaults()
	viper.Set("api.key", "custom")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.API.Key != "custom" || cfg.Quota.Timeout != 10*time.Second || cfg.Quota.MaxResponseBytes != 262144 {
		t.Fatalf("quota defaults: %+v", cfg)
	}
}

func TestQuotaDisplayDefaultsMerge(t *testing.T) {
	// Serial: Viper is process-global state.
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(strings.NewReader("quota:\n  heading: My quota\n  labels:\n    used: consumed\n  display:\n    types:\n      time_limit: Custom time\n")); err != nil {
		t.Fatal(err)
	}
	SetDefaults()
	if viper.GetString("quota.heading") != "My quota" || viper.GetString("quota.labels.used") != "consumed" || viper.GetString("quota.display.types.time_limit") != "Custom time" {
		t.Fatal("custom quota settings overwritten")
	}
	for _, key := range []string{"quota.display.badge", "quota.display.types.tokens_limit", "quota.display.types.credit_limit", "quota.display.bar_width", "quota.display.reset_format"} {
		if viper.GetString(key) == "" {
			t.Errorf("missing default %s", key)
		}
	}
	cfg, err := Load()
	if err != nil || cfg.Quota.Timeout != 10*time.Second {
		t.Fatalf("quota defaults missing alongside custom settings: %+v: %v", cfg, err)
	}
}

func TestQuotaBootstrap(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := BootstrapDefault(path); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("permissions %v", info.Mode())
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := BootstrapDefault(path); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("existing config overwritten")
	}
}
