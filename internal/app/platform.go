//go:build !windows

package app

import (
	"context"
	"errors"
	"os/exec"
	"runtime"
)

// OpenWith opens a file or URL with the system's default handler.
// On macOS: uses `open`, Linux: `xdg-open`.
func OpenWith(ctx context.Context, target string) error {
	cmd, err := buildOpenCommand(ctx, target)
	if err != nil {
		return err
	}
	return cmd.Run()
}

// buildOpenCommand creates the platform-specific command to open a file/URL.
// The target argument comes from CLI user input (file paths, URLs) which is intentional.
func buildOpenCommand(ctx context.Context, target string) (*exec.Cmd, error) {
	switch runtime.GOOS {
	case "darwin":
		return exec.CommandContext(ctx, "open", target), nil //nolint:gosec // G204: the CLI intentionally passes the requested target to the fixed OS opener.
	case "linux":
		return exec.CommandContext(ctx, "xdg-open", target), nil //nolint:gosec // G204: the CLI intentionally passes the requested target to the fixed OS opener.
	default:
		return nil, errors.New("no platform opener available")
	}
}
