package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/garyblankenship/zai/internal/app"
)

var (
	historyLimit int
	historyJSON  bool
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show chat history",
	Long:  `Display your chat history with timestamps and model information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showHistory()
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.Flags().IntVarP(&historyLimit, "limit", "l", 10, "number of entries (0 for all)")
	historyCmd.Flags().BoolVar(&historyJSON, "json", false, "Output in JSON format")
}

func showHistory() error {
	store := app.NewFileHistoryStore("")
	entries, err := store.GetRecent(historyLimit)
	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No chat history found.")
		return nil
	}

	if historyJSON {
		// Create structured JSON output
		output := map[string]interface{}{
			"history":   entries,
			"count":     len(entries),
			"limit":     historyLimit,
			"timestamp": time.Now().Format(time.RFC3339),
		}

		data, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
	} else {
		// Display human-readable table format
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIME\tTYPE\tMODEL\tPROMPT\tRESPONSE")
		fmt.Fprintln(w, "────\t────\t─────\t──────\t────────")

		for _, entry := range entries {
			// Determine type display
			typeDisplay := entry.Type
			if typeDisplay == "" {
				typeDisplay = "chat" // Default for backward compatibility
			}

			// Special handling for image entries
			var responseDisplay string
			if entry.Type == "image" {
				responseDisplay = fmt.Sprintf("🖼️ %s", entry.ImageSize)
			} else if entry.Type == "web" {
				responseDisplay = "🌐 web content"
			} else {
				// Handle Response as interface{}
				if respStr, ok := entry.Response.(string); ok {
					responseDisplay = truncate(respStr, 30)
				} else {
					responseDisplay = "📝 complex response"
				}
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				entry.Timestamp.Format("01-02 15:04"),
				typeDisplay,
				entry.Model,
				truncate(entry.Prompt, 30),
				responseDisplay,
			)
		}
		w.Flush()

		if historyLimit > 0 && len(entries) >= historyLimit {
			fmt.Printf("\nShowing %d most recent. Use -l 0 for all.\n", historyLimit)
		}
	}

	return nil
}
