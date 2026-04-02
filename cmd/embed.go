package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/dotcommander/zai/internal/app"
)

var embedModel string

var embedCmd = &cobra.Command{
	Use:   "embed [text]",
	Short: "Generate text embeddings using Z.AI Embedding-3",
	Long: `Generate vector embeddings for text using Z.AI's Embedding-3 model.

Always outputs JSON (embedding vectors are numeric arrays).

Examples:
  zai embed "Hello, world!"
  echo "text" | zai embed
  zai embed "sentence one" "sentence two"`,
	Args: cobra.ArbitraryArgs,
	RunE: runEmbed,
}

func registerEmbedCmd() {
	rootCmd.AddCommand(embedCmd)
	embedCmd.Flags().StringVarP(&embedModel, "model", "m", "", "Override embedding model (default: config api.embedding_model)")
}

func runEmbed(cmd *cobra.Command, args []string) error {
	texts := args

	if hasStdinData() {
		stdinText, err := readStdin()
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		if stdinText != "" {
			texts = append(texts, stdinText)
		}
	}

	if len(texts) == 0 {
		return fmt.Errorf("at least one text input required")
	}

	model := embedModel
	if model == "" {
		model = getModelWithDefault("api.embedding_model", "embedding-3")
	}

	client := newClient()
	ctx, cancel := createContext(30 * time.Second)
	defer cancel()

	opts := app.EmbeddingOptions{Model: model}

	resp, err := client.CreateEmbedding(ctx, texts, opts)
	if err != nil {
		return fmt.Errorf("embedding failed: %w", err)
	}

	output := struct {
		Model      string              `json:"model"`
		InputCount int                 `json:"input_count"`
		Usage      app.EmbeddingUsage  `json:"usage"`
		Data       []app.EmbeddingData `json:"data"`
		Timestamp  string              `json:"timestamp"`
	}{
		Model:      resp.Model,
		InputCount: len(texts),
		Usage:      resp.Usage,
		Data:       resp.Data,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	out, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}
	fmt.Println(string(out))

	return nil
}
