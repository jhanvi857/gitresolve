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
