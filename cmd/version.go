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
		fmt.Println()
		fmt.Println(theme.Title.Render(" ZAI ") + " " + theme.Info.Render(version.String()))
		fmt.Println()
		if version.Build != "unknown" {
			fmt.Printf("  %s %s\n", theme.Dim.Render("Built:"), theme.Description.Render(version.Build))
		}
		if version.Commit != "unknown" {
			fmt.Printf("  %s %s\n", theme.Dim.Render("Commit:"), theme.Description.Render(version.Commit))
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
