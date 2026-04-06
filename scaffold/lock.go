package scaffold

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

func acquireLock(dir string) (*os.File, error) {
	lockPath := path.Join(dir, ".boilerplate-go.lock")

	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err == nil {
		fmt.Fprint(f, strconv.Itoa(os.Getpid()))
		return f, nil
	}

	if !os.IsExist(err) {
		return nil, fmt.Errorf("create lock file: %w", err)
	}

	if tryRemoveStaleLock(lockPath) {
		return acquireLock(dir)
	}

	content, readErr := os.ReadFile(lockPath)
	if readErr == nil && len(content) > 0 {
		return nil, fmt.Errorf("another boilerplate-go process (PID %s) is already running in this directory", string(content))
	}

	return nil, fmt.Errorf("another boilerplate-go process is already running in this directory (if this is stale, remove %s)", lockPath)
}

// tryRemoveStaleLock reads the lock file PID and removes it if the process is dead.
func tryRemoveStaleLock(lockPath string) bool {
	content, err := os.ReadFile(lockPath)
	if err != nil || len(content) == 0 {
		return false
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		return false
	}

	if isProcessAlive(pid) {
		return false
	}

	os.Remove(lockPath)
	return true
}

func releaseLock(f *os.File) {
	name := f.Name()
	f.Close()
	os.Remove(name)
}
