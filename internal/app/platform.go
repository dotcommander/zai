package app

import (
	"fmt"
	"os/exec"
)

// OpenWith opens a file or URL with the system's default handler.
// On macOS: uses `open`, Linux: `xdg-open`, Windows: `cmd /c start`.
func OpenWith(target string) error {
	cmd, err := buildOpenCommand(target)
	if err != nil {
		return err
	}
	return cmd.Start()
}

// buildOpenCommand creates the platform-specific command to open a file/URL.
func buildOpenCommand(target string) (*exec.Cmd, error) {
	// macOS
	if _, err := exec.LookPath("open"); err == nil {
		return exec.Command("open", target), nil
	}
	// Linux
	if _, err := exec.LookPath("xdg-open"); err == nil {
		return exec.Command("xdg-open", target), nil
	}
	// Windows - note the empty string argument is required for proper parsing
	if _, err := exec.LookPath("start"); err == nil {
		return exec.Command("cmd", "/c", "start", "", target), nil
	}
	return nil, fmt.Errorf("no platform opener available (need: open, xdg-open, or start)")
}
