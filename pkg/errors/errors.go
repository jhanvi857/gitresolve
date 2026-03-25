package errors

import "errors"

var (
	ErrNoRepo             = errors.New("not a git repository")
	ErrLockAcquired       = errors.New("another gitresolve process is running")
	ErrNoConflicts        = errors.New("no conflicts found in file")
	ErrSyntaxInvalid      = errors.New("file has syntax errors after resolution")
	ErrDryRun             = errors.New("dry-run mode: no changes written")
	ErrMergeBaseNotFound  = errors.New("could not find common ancestor commit")
	ErrRenameDetected     = errors.New("rename conflict detected")
	ErrVerificationFailed = errors.New("resolved file failed verification")
	ErrSeverityCritical   = errors.New("conflict severity is critical, manual review required")
	ErrUnsupportedLang    = errors.New("language not supported for AST parsing")
)
