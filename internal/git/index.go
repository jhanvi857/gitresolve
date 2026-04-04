package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	gserrors "github.com/jhanvi857/gitresolve/pkg/errors"
)

// when two devs push conflicting changes git marks those files in staging area
// this file finds those files, checks individual ones, and marks them resolved

func ConflictedFiles(r *Repository) ([]string, error) {
	// Run git diff from the repository's working directory, not CWD
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	cmd.Dir = r.Path
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Warning: Could not get conflicted files from index.")
	}

	lines := strings.Split(string(out), "\n")
	var conflicted []string

	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		// Validate each path is within the repo (prevents path traversal)
		if err := validateRepoPath(r.Path, l); err != nil {
			fmt.Printf("Warning: skipping unsafe path %q: %v\n", l, err)
			continue
		}
		conflicted = append(conflicted, l)
	}

	if len(conflicted) == 0 {
		return nil, fmt.Errorf("ConflictedFiles: %w", gserrors.ErrNoConflicts)
	}

	return conflicted, nil
}

// validateRepoPath ensures a relative path does not escape the repository root.
func validateRepoPath(repoRoot, relativePath string) error {
	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return fmt.Errorf("resolving repo root: %w", err)
	}

	absFile, err := filepath.Abs(filepath.Join(repoRoot, relativePath))
	if err != nil {
		return fmt.Errorf("resolving file path: %w", err)
	}

	if !strings.HasPrefix(absFile, absRoot+string(filepath.Separator)) && absFile != absRoot {
		return fmt.Errorf("path escapes repository root")
	}

	return nil
}

func IsConflicted(r *Repository, filePath string) (bool, error) {
	conflicted, err := ConflictedFiles(r)
	if err != nil {
		if strings.Contains(err.Error(), "no conflicts") {
			return false, nil
		}
		return false, err
	}

	for _, f := range conflicted {
		if strings.TrimSpace(f) == strings.TrimSpace(filePath) {
			return true, nil
		}
	}

	return false, nil
}

func isStatusConflicted(fileStatus *gogit.FileStatus) bool {
	if fileStatus == nil {
		return false
	}

	return fileStatus.Staging == gogit.UpdatedButUnmerged ||
		fileStatus.Worktree == gogit.UpdatedButUnmerged
}

func MarkResolved(r *Repository, filePath string) error {
	// Validate path before passing to git add
	if err := validateRepoPath(r.Path, filePath); err != nil {
		return fmt.Errorf("MarkResolved: unsafe path: %w", err)
	}

	// calling git add directly is way more reliable on windows with unmerged entries
	// Use "--" to separate flags from file paths (prevents flag injection)
	cmd := exec.Command("git", "add", "--", filePath)
	cmd.Dir = r.Path
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("MarkResolved: git add: %w (output: %s)", err, string(out))
	}
	return nil
}

func ScanForMarkers(repoPath string) ([]string, error) {
	var results []string
	absRoot, _ := filepath.Abs(repoPath)

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "node_modules" || info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		// check for markers in text files
		if isLikelyText(path) {
			content, err := os.ReadFile(path)
			if err == nil {
				if strings.Contains(string(content), "<<<<<<<") &&
					strings.Contains(string(content), "=======") &&
					strings.Contains(string(content), ">>>>>>>") {
					rel, err := filepath.Rel(repoPath, path)
					if err == nil {
						// Validate the path doesn't escape the repo
						absFile, _ := filepath.Abs(path)
						if strings.HasPrefix(absFile, absRoot+string(filepath.Separator)) {
							results = append(results, rel)
						}
					}
				}
			}
		}
		return nil
	})
	return results, err
}

func isLikelyText(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go", ".js", ".ts", ".json", ".yaml", ".yml", ".toml", ".txt", ".md", ".py", ".java", ".c", ".cpp", ".h", ".sh", ".rs", ".rb", ".php", ".css", ".html", ".xml", ".jsx", ".tsx", ".swift", ".kt":
		return true
	}
	return false
}

