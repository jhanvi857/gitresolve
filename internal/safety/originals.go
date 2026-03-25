package safety

import (
	"fmt"
	"os"
)

func PreserveOriginal(filepath string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("PreserveOriginal : reading file %w", err)
	}
	backupPath := filepath + ".gitresolve-orig"
	err = os.WriteFile(backupPath, content, 0644)
	if err != nil {
		return fmt.Errorf("preserveOriginal : writing backup %w", err)
	}
	return nil
}

// restore original.
func RestoreOriginal(filePath string) error {
	backupPath := filePath + ".gitresolve-orig"
	conent, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("RestoreOriginal : reading file %w", err)
	}
	// using atomicwrite to restore original content.
	err = writeAtomic(filePath, conent)
	if err != nil {
		return fmt.Errorf("RestoreOriginal : restoring original %w", err)
	}
	return nil
}
