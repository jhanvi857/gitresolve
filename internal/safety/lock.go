package safety

// for preventing multiple gitresolve processes from running at the same time on the same repo,
//  We create a lock file when we start, and delete it when we're done.
//  If another process tries to start while the lock file exists, it will fail with an error.
import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	gserrors "github.com/jhanvi857/gitresolve/pkg/errors"
)

// lockfile : file we create to signal that gitresolve is running.
// Stored inside .git/ to avoid polluting the working tree and git status.
const LockFile = ".git/gitresolve.lock"

// AcquireLock creates the lock file. If it already exists, it checks if the process is still alive.
func AcquireLock(repoPath string) error {
	lockPath := filepath.Join(repoPath, LockFile)

	// Check if lock file exist
	if data, err := os.ReadFile(lockPath); err == nil {
		// Read PID from file
		lines := strings.Split(string(data), "\n")
		if len(lines) > 0 {
			parts := strings.Fields(lines[0])
			if len(parts) >= 2 && parts[0] == "pid" {
				oldPid, ferr := strconv.Atoi(parts[1])
				if ferr == nil && isProcessRunning(oldPid) {
					return fmt.Errorf("AcquireLock : %w (pid %d is still running)", gserrors.ErrLockAcquired, oldPid)
				}
			}
		}
		// If process is not running, we can remove the stale lock
		_ = os.Remove(lockPath)
	}

	// os.O_CREATE|os.O_EXCL : create this file but fail if already exists
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("AcquireLock : %w", gserrors.ErrLockAcquired)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "pid %d\nlocked at %s\n", os.Getpid(), time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("AcquireLock: writing PID: %w", err)
	}

	return nil
}

// ReleaseLock deletes the lock file when done.
func ReleaseLock(repoPath string) error {
	lockPath := filepath.Join(repoPath, LockFile)
	err := os.Remove(lockPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("ReleaseLock %w", err)
	}
	return nil
}

// isProcessRunning checks if a process with the given PID is currently active.
func isProcessRunning(pid int) bool {
	if pid <= 0 {
		return false
	}
	// On Windows, os.FindProcess always returns a process handle even if it doesn't exist.
	// We need to use Signals or other checks.
	// On Unix, p.Signal(0) is perfect.
	// On Windows, we can use the result of FindProcess and then check its status.
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	return checkProcessAlive(p)
}
