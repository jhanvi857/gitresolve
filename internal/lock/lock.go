package lock

import (
	"errors"
	"fmt"
	"os"
)

type RepoLock struct {
	f *os.File
}

var ErrLockContention = errors.New("repository is locked by another gitresolve process")

const LockFile = ".gitresolve/repo.lock"

func Acquire(root *os.Root) (*RepoLock, error) {
	f, err := root.OpenFile(LockFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if err != nil {
		if os.IsExist(err) {
			return nil, ErrLockContention
		}
		return nil, err
	}
	fmt.Fprintf(f, "%d\n", os.Getpid())

	if err := platformAcquire(f); err != nil {
		f.Close()
		return nil, err
	}

	return &RepoLock{f: f}, nil
}

func (l *RepoLock) Release() error {
	if l.f == nil {
		return nil
	}

	var releaseErr error
	if err := platformRelease(l.f); err != nil {
		releaseErr = err
	}

	path := l.f.Name()
	if err := l.f.Close(); err != nil && releaseErr == nil {
		releaseErr = err
	}
	l.f = nil

	if err := os.Remove(path); err != nil && releaseErr == nil {
		releaseErr = err
	}
	return releaseErr
}
