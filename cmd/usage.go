package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/term"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/app"
	"github.com/dotcommander/zai/internal/config"
)

type usageReader interface {
	Usage(context.Context) ([]app.QuotaWindow, error)
}

func newUsageCommand() *cobra.Command {
	return &cobra.Command{Use: "usage", Short: viper.GetString("quota.command_help"), Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		settings, err := config.Load()
		if err != nil {
			return usageFailure(cmd.OutOrStdout(), viper.GetBool("json"), "invalid_response")
		}
		client := &app.QuotaClient{APIKey: viper.GetString("api.key"), BaseURL: viper.GetString("api.base_url"), CodingBaseURL: viper.GetString("api.coding_base_url"), Settings: settings.Quota}
		return writeUsage(cmd.Context(), cmd.OutOrStdout(), viper.GetBool("json"), client)
	}}
}

func usageFailure(out io.Writer, jsonMode bool, reason string) error {
	if jsonMode {
		payload := struct {
			Error struct {
				Reason  string `json:"reason"`
				Message string `json:"message"`
			} `json:"error"`
		}{}
		payload.Error.Reason = reason
		payload.Error.Message = viper.GetString("quota.unavailable")
		if payload.Error.Message == "" {
			payload.Error.Message = reason
		}
		_ = json.NewEncoder(out).Encode(payload)
	}
	return fmt.Errorf("quota unavailable: %s", reason)
}

func writeUsage(ctx context.Context, out io.Writer, jsonMode bool, client usageReader) error {
	windows, err := client.Usage(ctx)
	if err != nil || len(windows) == 0 {
		var unavailable *app.QuotaUnavailable
		reason := "request_failed"
		if errors.As(err, &unavailable) {
			reason = unavailable.Reason
		}
		if err == nil {
			reason = "invalid_response"
		}
		return usageFailure(out, jsonMode, reason)
	}
	if jsonMode {
		return json.NewEncoder(out).Encode(struct {
			Windows []app.QuotaWindow `json:"windows"`
		}{windows})
	}
	if destination, ok := out.(*os.File); ok && term.IsTerminal(destination.Fd()) {
		width, _, sizeErr := term.GetSize(destination.Fd())
		if sizeErr != nil {
			width = 0
		}
		_, noColor := os.LookupEnv("NO_COLOR")
		return writeUsagePretty(out, windows, currentUsageDisplay(), width, !noColor)
	}
	heading := viper.GetString("quota.heading")
	if heading != "" {
		if _, err := fmt.Fprintln(out, heading); err != nil {
			return err
		}
	}
	for _, window := range windows {
		details := []string{fmt.Sprintf("%s: %g%%", window.Type, window.Percentage)}
		if window.Used != nil {
			details = append(details, fmt.Sprintf("%s %g", viper.GetString("quota.labels.used"), *window.Used))
		}
		if window.Limit != nil {
			details = append(details, fmt.Sprintf("%s %g", viper.GetString("quota.labels.limit"), *window.Limit))
		}
		if window.Remaining != nil {
			details = append(details, fmt.Sprintf("%s %g", viper.GetString("quota.labels.remaining"), *window.Remaining))
		}
		if window.DurationMinutes != nil {
			details = append(details, fmt.Sprintf("%s %dm", viper.GetString("quota.labels.duration"), *window.DurationMinutes))
		}
		if window.ResetsAt != nil {
			details = append(details, viper.GetString("quota.labels.reset")+" "+window.ResetsAt.Format("2006-01-02T15:04:05Z07:00"))
		}
		if _, err := fmt.Fprintln(out, strings.Join(details, " | ")); err != nil {
			return err
		}
	}
	return nil
}
