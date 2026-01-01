package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/dotcommander/zai/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Display version information including version, build time, and commit hash.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("zai version %s\n", version.Version)
		fmt.Printf("Build:   %s\n", version.Build)
		fmt.Printf("Commit:  %s\n", version.Commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
