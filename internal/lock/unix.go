//go:build linux || darwin

package lock

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func platformAcquire(f *os.File) error {
	err := unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB)
	if err != nil {
		if err == unix.EWOULDBLOCK || err == syscall.EAGAIN {
			return ErrLockContention
		}
		return err
	}
	return nil
}

func platformRelease(f *os.File) {
	_ = unix.Flock(int(f.Fd()), unix.LOCK_UN)
}
