package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var modelCmd = &cobra.Command{
	Use:   "model <subcommand>",
	Short: "Model management commands",
	Long:  `Commands for managing and listing available models.`,
}

var modelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available models from the API",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runModelList()
	},
}

func init() {
	rootCmd.AddCommand(modelCmd)
	modelCmd.AddCommand(modelListCmd)
}

func runModelList() error {
	client := newClient()

	models, err := client.ListModels(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	fmt.Println("Available Models:")
	fmt.Println("─────────────────")
	for _, m := range models {
		created := time.Unix(m.Created, 0).Format("2006-01-02")
		fmt.Printf("  %s  (created: %s)\n", m.ID, created)
	}

	return nil
}
