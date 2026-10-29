//go:build windows

package app

import (
	"context"
	"fmt"

	"golang.org/x/sys/windows"
)

// OpenWith opens a file or URL with the system's default handler.
// ShellExecute passes the target as a file argument without invoking a command shell.
// Cancellation is checked before launch; ShellExecute cannot cancel an opened handler.
func OpenWith(ctx context.Context, target string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	file, err := windows.UTF16PtrFromString(target)
	if err != nil {
		return fmt.Errorf("encode opener target: %w", err)
	}
	return windows.ShellExecute(0, nil, file, nil, nil, windows.SW_SHOWNORMAL)
}
