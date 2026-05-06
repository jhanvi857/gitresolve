package lock

import (
	"errors"
	"os"
)

type RepoLock struct {
	f *os.File
}

var ErrLockContention = errors.New("repository is locked by another gitresolve process")

const LockFile = ".gitresolve/repo.lock"

func Acquire(root *os.Root) (*RepoLock, error) {
	// safepath: CWE-22 hardened
	f, err := root.OpenFile(LockFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

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

	platformRelease(l.f)

	path := l.f.Name()
	l.f.Close()
	l.f = nil

	_ = os.Remove(path)
	return nil
}
