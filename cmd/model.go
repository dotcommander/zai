package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/app"
)

var modelCmd = &cobra.Command{
	Use:   "model <subcommand>",
	Short: "Model management commands",
}

var (
	modelJSON bool
)

var modelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available models from the API",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runModelList()
	},
}

func registerModelCmd() {
	rootCmd.AddCommand(modelCmd)
	modelCmd.AddCommand(modelListCmd)

	modelListCmd.Flags().BoolVar(&modelJSON, "json", false, "Output in JSON format")
}

// configuredModel represents a model set via config, not discovered from the API.
type configuredModel struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func runModelList() error {
	client := newClient()

	ctx, cancel := createContext(30 * time.Second)
	defer cancel()

	models, err := client.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	configured := collectConfiguredModels()

	if modelJSON {
		return printModelListJSON(models, configured)
	}
	printModelListTable(models, configured)
	return nil
}

func collectConfiguredModels() []configuredModel {
	type entry struct {
		key, label string
	}
	sources := []entry{
		{"api.image_model", "image"},
		{"api.video_model", "video"},
		{"api.vision_model", "vision"},
		{"api.audio_model", "audio"},
		{"api.tts_model", "tts"},
		{"api.embedding_model", "embedding"},
	}

	seen := make(map[string]bool)
	var out []configuredModel
	for _, s := range sources {
		id := viper.GetString(s.key)
		if id == "" || seen[id+s.label] {
			continue
		}
		seen[id+s.label] = true
		out = append(out, configuredModel{ID: id, Type: s.label})
	}
	return out
}

func printModelListTable(apiModels []app.Model, configured []configuredModel) {
	// Chat models from API
	fmt.Println("Chat Models (from API):")
	fmt.Println("───────────────────────")
	if len(apiModels) == 0 {
		fmt.Println("  (none)")
	} else {
		sort.Slice(apiModels, func(i, j int) bool {
			return apiModels[i].ID < apiModels[j].ID
		})
		for _, m := range apiModels {
			created := time.Unix(m.Created, 0).Format("2006-01-02")
			fmt.Printf("  %-28s  %s\n", m.ID, created)
		}
	}

	// Configured models by type
	if len(configured) > 0 {
		fmt.Println()
		fmt.Println("Configured Models:")
		fmt.Println("──────────────────")
		for _, m := range configured {
			fmt.Printf("  %-28s  %s\n", m.ID, m.Type)
		}
	}
}

func printModelListJSON(apiModels interface{}, configured []configuredModel) error {
	output := map[string]interface{}{
		"chat_models":       apiModels,
		"configured_models": configured,
		"timestamp":         time.Now().Format(time.RFC3339),
	}
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
