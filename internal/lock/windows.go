//go:build windows

package lock

import (
	"golang.org/x/sys/windows"
	"os"
)

func platformAcquire(f *os.File) error {
	// LockFileEx flags:
	// LOCKFILE_EXCLUSIVE_LOCK = 2
	// LOCKFILE_FAIL_IMMEDIATELY = 1
	var flags uint32 = windows.LOCKFILE_EXCLUSIVE_LOCK | windows.LOCKFILE_FAIL_IMMEDIATELY

	// Overlapped structure is required but can be zeroed for simple locking
	overlapped := &windows.Overlapped{}

	// We lock the first byte of the file.
	err := windows.LockFileEx(windows.Handle(f.Fd()), flags, 0, 1, 0, overlapped)
	if err != nil {
		// ERROR_LOCK_VIOLATION = 33
		if err == windows.ERROR_LOCK_VIOLATION || err == windows.ERROR_IO_PENDING {
			return ErrLockContention
		}
		return err
	}
	return nil
}

func platformRelease(f *os.File) {
	overlapped := &windows.Overlapped{}
	_ = windows.UnlockFileEx(windows.Handle(f.Fd()), 0, 1, 0, overlapped)
}
