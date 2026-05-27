//go:build plan9 || js

package safepath

import (
	"fmt"
	"os"
	"path/filepath"
)

func RepoRoot(repoDir string) (*os.Root, error) {
	if !IsForceAllowed() {
		return nil, UnsupportedPlatformErr()
	}
	return nil, nil
}

func SafeOpen(root *os.Root, relPath string) (*os.File, error) {
	if !filepath.IsLocal(relPath) {
		return nil, ErrUnsafePath
	}
	if !IsForceAllowed() {
		return nil, UnsupportedPlatformErr()
	}
	f, err := os.Open(relPath)
	if err != nil {
		return nil, fmt.Errorf("SafeOpen(force): %w", err)
	}
	return f, nil
}

func SafeWrite(root *os.Root, relPath string, data []byte, perm os.FileMode) error {
	if !filepath.IsLocal(relPath) {
		return ErrUnsafePath
	}
	if !IsForceAllowed() {
		return UnsupportedPlatformErr()
	}
	return os.WriteFile(relPath, data, perm)
}
