//go:build unix

package mcp

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

// Own a fresh process group so npm's server child does not survive cancellation.
// The negative PID targets only the group created for this invocation.
func configureProcess(command *exec.Cmd) {
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	command.Cancel = func() error {
		if command.Process == nil {
			return os.ErrProcessDone
		}
		err := syscall.Kill(-command.Process.Pid, syscall.SIGKILL)
		if errors.Is(err, syscall.ESRCH) {
			return os.ErrProcessDone
		}
		return err
	}
}
