//go:build !plan9 && !js

package safepath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func RepoRoot(repoDir string) (*os.Root, error) {
	root, err := os.OpenRoot(repoDir)
	if err != nil {
		return nil, fmt.Errorf("RepoRoot: %w", err)
	}
	return root, nil
}

func SafeOpen(root *os.Root, relPath string) (*os.File, error) {
	if root == nil {
		return nil, fmt.Errorf("SafeOpen: nil root")
	}
	if !filepath.IsLocal(relPath) {
		return nil, ErrUnsafePath
	}

	f, err := root.Open(relPath)
	if err != nil {
		if isUnsafeRootError(err) {
			return nil, fmt.Errorf("%w: %v", ErrUnsafePath, err)
		}
		return nil, err
	}
	return f, nil
}

func SafeWrite(root *os.Root, relPath string, data []byte, perm os.FileMode) error {
	if root == nil {
		return fmt.Errorf("SafeWrite: nil root")
	}
	if !filepath.IsLocal(relPath) {
		return ErrUnsafePath
	}

	dir := filepath.Dir(relPath)
	if dir == "." {
		dir = ""
	}
	base := filepath.Base(relPath)
	tmpBase := ".gitresolve-tmp-" + strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + base
	tmpPath := tmpBase
	if dir != "" {
		tmpPath = filepath.Join(dir, tmpBase)
	}

	f, err := root.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		if isUnsafeRootError(err) {
			return fmt.Errorf("%w: %v", ErrUnsafePath, err)
		}
		return fmt.Errorf("SafeWrite: create temp: %w", err)
	}

	cleanup := true
	defer func() {
		if cleanup {
			if err := root.Remove(tmpPath); err != nil {
				_ = err
			}
		}
	}()

	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return fmt.Errorf("SafeWrite: write temp: %w", err)
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return fmt.Errorf("SafeWrite: sync temp: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("SafeWrite: close temp: %w", err)
	}

	if err := root.Rename(tmpPath, relPath); err != nil {
		if isUnsafeRootError(err) {
			return fmt.Errorf("%w: %v", ErrUnsafePath, err)
		}
		return fmt.Errorf("SafeWrite: rename temp: %w", err)
	}

	cleanup = false
	return nil
}

func isUnsafeRootError(err error) bool {
	if errors.Is(err, os.ErrPermission) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "escape") || strings.Contains(msg, "outside") || strings.Contains(msg, "travers") || strings.Contains(msg, "symlink")
}
