package errors

import "errors"

var (
	ErrNoRepo        = errors.New("not a git repository")
	ErrLockAcquired  = errors.New("another gitstitch process is running")
	ErrNoConflicts   = errors.New("no conflicts found in file")
	ErrSyntaxInvalid = errors.New("file has syntax errors after resolution")
	ErrDryRun        = errors.New("dry-run mode: no changes written")
)
