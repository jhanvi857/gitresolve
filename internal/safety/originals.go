package safety

import (
	"fmt"
	"io"
	"os"

	"github.com/jhanvi857/gitresolve/internal/safepath"
)

func PreserveOriginal(root *os.Root, filePath string) error {
	f, err := safepath.SafeOpen(root, filePath)
	if err != nil {
		return fmt.Errorf("PreserveOriginal : reading file %w", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("PreserveOriginal : reading file %w", err)
	}

	backupPath := filePath + ".gitresolve-orig"
	// safepath: CWE-22 hardened
	err = safepath.SafeWrite(root, backupPath, content, 0644)
	if err != nil {
		return fmt.Errorf("preserveOriginal : writing backup %w", err)
	}
	return nil
}

// restore original.
func RestoreOriginal(root *os.Root, filePath string) error {
	backupPath := filePath + ".gitresolve-orig"
	f, err := safepath.SafeOpen(root, backupPath)
	if err != nil {
		return fmt.Errorf("RestoreOriginal : reading file %w", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("RestoreOriginal : reading file %w", err)
	}

	// safepath: CWE-22 hardened
	err = safepath.SafeWrite(root, filePath, content, 0644)
	if err != nil {
		return fmt.Errorf("RestoreOriginal : restoring original %w", err)
	}
	return nil
}
