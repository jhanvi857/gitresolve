package safety

import (
	"fmt"
	"os"

	"github.com/jhanvi857/gitresolve/internal/safepath"
	gserrors "github.com/jhanvi857/gitresolve/pkg/errors"
	"github.com/jhanvi857/gitresolve/pkg/logger"
)

// true = dry run mode, false = real writes
type DryRunWriter struct {
	Enabled bool
	root    *os.Root
}

// Write either performs the real write or just logs what would happen
func (d DryRunWriter) Write(targetPath string, content []byte) error {
	if d.Enabled {
		// just log and return a special error so the caller knows nothing was written
		logger.Info().
			Str("path", targetPath).
			Int("bytes", len(content)).
			Msg("dry-run: would write file")
		return fmt.Errorf("Write: %w", gserrors.ErrDryRun)
	}
	// safepath: CWE-22 hardened
	return safepath.SafeWrite(d.root, targetPath, content, 0644)
}

// NewWriter creates a DryRunWriter based on the --dry-run flag
func NewWriter(dryRun bool, root *os.Root) DryRunWriter {
	return DryRunWriter{Enabled: dryRun, root: root}
}
