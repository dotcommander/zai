package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"zai/internal/app"
)

var historyLimit int

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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIME\tMODEL\tPROMPT\tRESPONSE")
	fmt.Fprintln(w, "────\t─────\t──────\t────────")

	for _, entry := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			entry.Timestamp.Format("01-02 15:04"),
			entry.Model,
			truncate(entry.Prompt, 30),
			truncate(entry.Response, 30),
		)
	}
	w.Flush()

	if historyLimit > 0 && len(entries) >= historyLimit {
		fmt.Printf("\nShowing %d most recent. Use -l 0 for all.\n", historyLimit)
	}

	return nil
}
