package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/config"
)

func TestUsageInitializationFailure(t *testing.T) {
	// Serial: command configuration uses process-global Viper and cfgFile.
	previous := cfgFile
	t.Cleanup(func() { cfgFile = previous; viper.Reset(); config.SetDefaults() })
	viper.Reset()
	config.SetDefaults()
	viper.Set("json", true)
	cfgFile = filepath.Join(t.TempDir(), "invalid.yaml")
	if err := os.WriteFile(cfgFile, []byte("api: [\n"), 0600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	command := &cobra.Command{Use: "usage"}
	command.SetOut(&out)
	err := initializeCommandConfig(command, nil)
	if err == nil || !strings.Contains(err.Error(), "quota unavailable") || strings.Contains(err.Error(), cfgFile) {
		t.Fatalf("unsafe initialization error: %v", err)
	}
	var failure struct {
		Error struct {
			Reason string `json:"reason"`
		} `json:"error"`
	}
	if err := json.Unmarshal(out.Bytes(), &failure); err != nil || failure.Error.Reason != "invalid_response" {
		t.Fatalf("missing safe JSON error: %s (%v)", out.String(), err)
	}
}

func TestUsageConfigBootstrapWaitsForRootInput(t *testing.T) {
	// Serial: HOME, command configuration and Viper are process-global.
	previous := cfgFile
	t.Cleanup(func() { cfgFile = previous; viper.Reset(); config.SetDefaults() })
	home := t.TempDir()
	t.Setenv("HOME", home)
	viper.Reset()
	config.SetDefaults()
	cfgFile = ""
	command := &cobra.Command{Use: rootCmd.Name()}
	if err := initializeCommandConfig(command, nil); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(home, ".config", "zai", "config.yaml")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("root metadata created configuration: %v", err)
	}
	viper.Set("api.key", "offline-placeholder")
	if err := initializeRootInputConfig(command, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("actual root input did not bootstrap configuration: %v", err)
	}
}
