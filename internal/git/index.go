package git

import (
	"fmt"
	"os/exec"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	gserrors "github.com/jhanvi857/gitresolve/pkg/errors"
)

// when two devs push conflicting changes git marks those files in staging area
// this file finds those files, checks individual ones, and marks them resolved

func ConflictedFiles(r *Repository) ([]string, error) {
	out, err := exec.Command("git", "diff", "--name-only", "--diff-filter=U").Output()
	if err != nil {
		return nil, fmt.Errorf("ConflictedFiles: calling git diff: %w", err)
	}

	lines := strings.Split(string(out), "\n")
	var conflicted []string

	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			conflicted = append(conflicted, l)
		}
	}

	if len(conflicted) == 0 {
		return nil, fmt.Errorf("ConflictedFiles: %w", gserrors.ErrNoConflicts)
	}

	return conflicted, nil
}

func IsConflicted(r *Repository, filePath string) (bool, error) {
	conflicted, err := ConflictedFiles(r)
	if err != nil {
		if strings.Contains(err.Error(), "no conflicts") {
			return false, nil
		}
		return false, err
	}
	
	for _, f := range conflicted {
		if strings.TrimSpace(f) == strings.TrimSpace(filePath) {
			return true, nil
		}
	}
	
	return false, nil
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
