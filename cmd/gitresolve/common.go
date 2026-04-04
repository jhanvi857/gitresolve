package gitresolve

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/jhanvi857/gitresolve/internal/conflict"
	"github.com/jhanvi857/gitresolve/internal/git"
	"github.com/jhanvi857/gitresolve/internal/store"
	"github.com/jhanvi857/gitresolve/pkg/logger"
)

func HandleSignals(r *git.Repository) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if r != nil {
			fmt.Printf("\nInterrupted. Releasing lock on %s...\n", r.Path)
			_ = git.Close(r)
		}
		os.Exit(1)
	}()
}

func dbPathForRepo(repoPath string) string {
	return filepath.Join(repoPath, ".git", "gitresolve.db")
}

func openStore(repoPath string) (*store.DB, error) {
	db, err := store.Open(dbPathForRepo(repoPath))
	if err != nil {
		return nil, fmt.Errorf("openStore: %w", err)
	}
	return db, nil
}

func ResolveRepoRoot() (string, error) {
	out, err := runGit("rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("not inside a git repository: %w", err)
	}
	root := strings.TrimSpace(out)
	root = filepath.FromSlash(root)
	return root, nil
}
func ValidatePath(repoRoot, filePath string) error {
	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return fmt.Errorf("resolving repo root: %w", err)
	}

	absFile, err := filepath.Abs(filepath.Join(repoRoot, filePath))
	if err != nil {
		return fmt.Errorf("resolving file path: %w", err)
	}

	if !strings.HasPrefix(absFile, absRoot+string(filepath.Separator)) && absFile != absRoot {
		return fmt.Errorf("path %q escapes repository root", filePath)
	}

	realRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		logger.Debug("symlink eval on root failed (non-fatal): " + err.Error())
		realRoot = absRoot
	}

	info, statErr := os.Lstat(absFile)
	if statErr == nil && info.Mode()&os.ModeSymlink != 0 {
		realFile, err := filepath.EvalSymlinks(absFile)
		if err == nil && !strings.HasPrefix(realFile, realRoot+string(filepath.Separator)) {
			return fmt.Errorf("symlink %q points outside repository", filePath)
		}
	}

	return nil
}

func severityLabel(s conflict.Severity) string {
	switch s {
	case conflict.SeverityTrivial:
		return "trivial"
	case conflict.SeverityLow:
		return "low"
	case conflict.SeverityMedium:
		return "medium"
	case conflict.SeverityHigh:
		return "high"
	case conflict.SeverityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

func typeLabel(t conflict.ConflictType) string {
	switch t {
	case conflict.TypeWhitespace:
		return "whitespace"
	case conflict.TypeImport:
		return "import"
	case conflict.TypeIdentical:
		return "identical"
	case conflict.TypeRename:
		return "rename"
	case conflict.TypeSignature:
		return "signature"
	case conflict.TypeLogic:
		return "logic"
	case conflict.TypeStructured:
		return "structured"
	case conflict.TypeDeleteModify:
		return "delete-modify"
	case conflict.TypeUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

func hasConflictMarkers(content string) bool {
	return strings.Contains(content, "<<<<<<<") &&
		strings.Contains(content, "=======") &&
		strings.Contains(content, ">>>>>>>")
}
