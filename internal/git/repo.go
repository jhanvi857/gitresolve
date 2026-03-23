package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	"github.com/jhanvi857/gitresolve/internal/safety"
	gserrors "github.com/jhanvi857/gitresolve/pkg/errors"
)

type Repository struct {
	Path string
	repo *gogit.Repository
}

// opening git repo at given path and aquire lock.
func Open(path string) (*Repository, error) {
	r, err := gogit.PlainOpenWithOptions(path, &gogit.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, fmt.Errorf("Open: %w", gserrors.ErrNoRepo)
	}

	err = safety.AcquireLock(path)
	if err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	}

	return &Repository{
		Path: path,
		repo: r,
	}, nil
}

// closing repo and release lock.
func Close(r *Repository) error {
	err := safety.ReleaseLock(r.Path)
	if err != nil {
		return fmt.Errorf("Close: %w", err)
	}
	return nil
}

// checking if the working directory has no uncommitted changes
// useful to warn the user before gitresolve starts modifying files
func (r *Repository) IsClean() (bool, error) {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("IsClean: getting worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return false, fmt.Errorf("IsClean: reading status: %w", err)
	}

	return status.IsClean(), nil
}

// fetching sha of head commit
func (r *Repository) HeadCommit() (string, error) {
	ref, err := r.repo.Head()
	if err != nil {
		return "", fmt.Errorf("HeadCommit: %w", err)
	}

	return ref.Hash().String(), nil
}
