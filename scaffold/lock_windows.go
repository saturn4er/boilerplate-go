//go:build windows

package scaffold

import (
	"os"
)

func isProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// On Windows, FindProcess always succeeds. Sending signal 0 is not
	// supported, so we attempt to release the process handle — if it
	// fails the process likely doesn't exist.
	err = proc.Release()
	return err == nil
}
