package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/app"
	"github.com/dotcommander/zai/internal/config"
)

func TestVideoImageURLValidation(t *testing.T) {
	t.Parallel()
	for _, input := range []string{"photo.jpg", "/tmp/photo.jpg", "file:///tmp/photo.jpg", "ftp://example.com/a.jpg", "https:///a.jpg", "//example.com/a.jpg"} {
		if err := validateVideoImageURLs([]string{input}); err == nil {
			t.Errorf("accepted %q", input)
		}
	}
	if err := validateVideoImageURLs([]string{"https://example.com/first.jpg", "http://example.com/last.jpg"}); err != nil {
		t.Fatal(err)
	}
}

func TestVideoLocalInputFailsBeforeClient(t *testing.T) {
	// Serial: command flags are process-global state.
	original := videoImageURLs
	t.Cleanup(func() { videoImageURLs = original })
	videoImageURLs = []string{"local.jpg"}
	if err := runVideoGeneration("prompt"); err == nil || !strings.Contains(err.Error(), "HTTP(S)") {
		t.Fatalf("error = %v", err)
	}
}

func TestMediaDownloadLimitPrecedence(t *testing.T) {
	// Serial: Viper and environment are process-global state.
	viper.Reset()
	t.Cleanup(viper.Reset)
	config.SetDefaults()
	cmd := &cobra.Command{}
	cmd.Flags().Int64("max-download-bytes", app.DefaultMaxDownloadBytes, "")
	assertLimit := func(want int64) {
		t.Helper()
		got, err := mediaDownloadLimit(cmd)
		if err != nil || got != want {
			t.Fatalf("limit = %d, %v; want %d", got, err, want)
		}
	}
	assertLimit(app.DefaultMaxDownloadBytes)
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(strings.NewReader("media:\n  max_download_bytes: 12\n")); err != nil {
		t.Fatal(err)
	}
	assertLimit(12)
	viper.SetEnvPrefix("ZAI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	t.Setenv("ZAI_MEDIA_MAX_DOWNLOAD_BYTES", "15")
	assertLimit(15)
	if err := cmd.Flags().Set("max-download-bytes", "20"); err != nil {
		t.Fatal(err)
	}
	assertLimit(20)
	t.Setenv("ZAI_MEDIA_MAX_DOWNLOAD_BYTES", "0")
	assertLimit(20) // An explicit valid flag overrides an invalid environment value.
	if err := cmd.Flags().Set("max-download-bytes", "0"); err != nil {
		t.Fatal(err)
	}
	if _, err := mediaDownloadLimit(cmd); err == nil {
		t.Fatal("accepted zero limit")
	}
}
