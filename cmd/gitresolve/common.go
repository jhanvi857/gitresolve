package gitresolve

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/jhanvi857/gitresolve/internal/conflict"
	"github.com/jhanvi857/gitresolve/internal/git"
	"github.com/jhanvi857/gitresolve/internal/ownership"
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
			if err := git.Close(r); err != nil {
				logger.Debug().Err(err).Msg("failed to close git repository")
			}
		}
		os.Exit(1)
	}()
}

func dbPathForRepo(repoPath string) string {
	if envPath := strings.TrimSpace(os.Getenv("GITRESOLVE_DB_PATH")); envPath != "" {
		if filepath.IsAbs(envPath) {
			return envPath
		}
		return filepath.Join(repoPath, envPath)
	}
	return filepath.Join(repoPath, ".gitresolve", "conflicts.db")
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
		logger.Debug().Msg("symlink eval on root failed (non-fatal): " + err.Error())
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

func hashContent(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func reasonCodeOrUnknown(c *conflict.ConflictBlock) string {
	if c != nil && conflict.IsStableReasonCode(c.ManualReasonCode) {
		return c.ManualReasonCode
	}
	return conflict.ReasonDecisionUnknown
}

func shouldAutoApplyWithPolicy(c *conflict.ConflictBlock, policy string) bool {
	switch strings.ToLower(policy) {
	case ownership.PolicyStrict:
		if c.Type != conflict.TypeWhitespace && c.Type != conflict.TypeIdentical {
			return false
		}
		return c.CanAutoResolve && c.Confidence >= conflict.AutoResolveConfidenceThreshold
	case ownership.PolicyAggressive:
		return c.CanAutoResolve && c.Confidence >= 0.70
	default:
		return conflict.ShouldAutoApply(c)
	}
}

func policyBlocksBothForFile(policy, filePath string) bool {
	if strings.ToLower(policy) != ownership.PolicyStrict {
		return false
	}
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, sourceExt := range strictPolicySourceLikeExtensions() {
		if ext == sourceExt {
			return true
		}
	}
	return false
}

func strictPolicySourceLikeExtensions() []string {
	return []string{".go", ".js", ".jsx", ".ts", ".tsx", ".py", ".java", ".kt", ".rb", ".php", ".rs", ".c", ".cc", ".cpp", ".h", ".hpp", ".cs", ".swift"}
}
