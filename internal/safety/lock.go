package safety

// for preventing multiple gitresolve processes from running at the same time on the same repo,
//  We create a lock file when we start, and delete it when we're done.
//  If another process tries to start while the lock file exists, it will fail with an error.
import (
	"fmt"
	"os"
	"time"

	gserrors "github.com/jhanvi857/gitresolve/pkg/errors"
)

// lockfile : file we create to signal that gitresolve is running.
// If it exists, we assume another process is running and exit with an error.
const LockFile = ".gitresolve.lock"

// AcquireLock creates the lock file. If it already exists, it returns an error.
func AcquireLock(repoPath string) error {
	lockPath := repoPath + "/" + LockFile
	// os.O_CREATE|os.O_EXCL : create this file but fail if already exists
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("AcquireLock : %w", gserrors.ErrLockAcquired)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "locked at %s\n", time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("AcquireLock: writing timestamp: %w", err)
	}

	return nil
}

// ReleaseLock deletes the lock file when done.
func ReleaseLock(repoPath string) error {
	lockPath := repoPath + "/" + LockFile
	err := os.Remove(lockPath)
	if err != nil {
		return fmt.Errorf("ReleaseLock %w", err)
	}
	return nil
}
