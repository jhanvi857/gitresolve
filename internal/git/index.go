package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	gserrors "github.com/jhanvi857/gitresolve/pkg/errors"
)

// when two devs push conflicting changes git marks those files in staging area
// this file finds those files, checks individual ones, and marks them resolved

func ConflictedFiles(r *Repository) ([]string, error) {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("ConflictedFiles: getting worktree: %w", err)
	}

	// status gives us state of every file git knows about
	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("ConflictedFiles: reading status: %w", err)
	}

	var conflicted []string

	for filePath, fileStatus := range status {
		// go-git marks conflicted paths as UpdatedButUnmerged on either side.
		if isStatusConflicted(fileStatus) {
			conflicted = append(conflicted, filePath)
		}
	}

	if len(conflicted) == 0 {
		return nil, fmt.Errorf("ConflictedFiles: %w", gserrors.ErrNoConflicts)
	}

	return conflicted, nil
}

func IsConflicted(r *Repository, filePath string) (bool, error) {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("IsConflicted: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return false, fmt.Errorf("IsConflicted: %w", err)
	}

	// get status of this specific file
	fileStatus := status.File(filePath)

	isConflicted := isStatusConflicted(fileStatus)

	return isConflicted, nil
}

func isStatusConflicted(fileStatus *gogit.FileStatus) bool {
	if fileStatus == nil {
		return false
	}

	return fileStatus.Staging == gogit.UpdatedButUnmerged ||
		fileStatus.Worktree == gogit.UpdatedButUnmerged
}

func MarkResolved(r *Repository, filePath string) error {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("MarkResolved: %w", err)
	}

	// git add : moves file from conflicted to staged
	_, err = worktree.Add(filePath)
	if err != nil {
		return fmt.Errorf("MarkResolved: staging %s: %w", filePath, err)
	}

	return nil
}
