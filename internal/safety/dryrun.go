package safety

import (
	"fmt"

	gserrors "github.com/jhanvi857/gitresolve/pkg/errors"
	"github.com/jhanvi857/gitresolve/pkg/logger"
)

// true = dry run mode, false = real writes
type DryRunWriter struct {
	Enabled bool
}

// Write either performs the real write or just logs what would happen
func (d DryRunWriter) Write(targetPath string, content []byte) error {
	if d.Enabled {
		// just log and return a special error so the caller knows nothing was written
		logger.Infof("dry-run: would write file", map[string]any{
			"path":  targetPath,
			"bytes": len(content),
		})
		return fmt.Errorf("Write: %w", gserrors.ErrDryRun)
	}
	return writeAtomic(targetPath, content)
}

// NewWriter creates a DryRunWriter based on the --dry-run flag
func NewWriter(dryRun bool) DryRunWriter {
	return DryRunWriter{Enabled: dryRun}
}
