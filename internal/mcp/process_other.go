//go:build !unix

package mcp

import "os/exec"

// CommandContext owns and terminates the launched process on other platforms.
func configureProcess(command *exec.Cmd) {}
