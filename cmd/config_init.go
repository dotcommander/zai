package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/config"
)

func initializeCommandConfig(cmd *cobra.Command, args []string) error {
	switch cmd.Name() {
	case "history", "completion", "help", "version":
		return nil
	}
	pendingRootInput := rootInputPending(cmd, args)
	noAPI := cmd.Name() == "usage" || pendingRootInput || (cmd.Name() == "model" && len(args) == 0)
	err := readCommandConfig(noAPI, !pendingRootInput)
	if err != nil && cmd.Name() == "usage" {
		return usageFailure(cmd.OutOrStdout(), viper.GetBool("json"), "invalid_response")
	}
	return err
}

func rootInputPending(cmd *cobra.Command, args []string) bool {
	return cmd.Parent() == nil && cmd.Name() == "zai" && len(args) == 0
}

func initializeRootInputConfig(cmd *cobra.Command, args []string) error {
	if rootInputPending(cmd, args) {
		return initConfig()
	}
	return nil
}

func initConfig(noAPI ...bool) error {
	return readCommandConfig(len(noAPI) > 0 && noAPI[0], true)
}

func readCommandConfig(noAPI, bootstrap bool) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configDir := filepath.Join(home, ".config", "zai")
		if bootstrap {
			if err := config.BootstrapDefault(filepath.Join(configDir, "config.yaml")); err != nil {
				return err
			}
		}
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}
	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return err
		}
	}
	viper.SetEnvPrefix("ZAI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if viper.GetString("api.key") == "" && !noAPI {
		return errors.New("API key required: set ZAI_API_KEY or configure in ~/.config/zai/config.yaml")
	}
	return nil
}
