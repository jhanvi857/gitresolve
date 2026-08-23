package history

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// DivergenceResult holds the outcome of a pre-resolve divergence check.
type DivergenceResult struct {
	Ahead         int    // commits HEAD is ahead of origin
	Behind        int    // commits HEAD is behind origin
	DefaultBranch string // the remote default branch name
	// Authors on the behind-side who touched files also modified locally
	AuthorsOnBehindTouchingLocal []AuthorContribution
	// Detected pull strategy from git config
	PullStrategy string // "rebase" or "merge"
}

// CheckDivergence runs `git rev-list --left-right --count` to determine
// how far the current branch has diverged from the remote default branch.
// It cross-references FileAuthors from the index to identify which authors
// on the behind-side touched files that are also locally modified.
func CheckDivergence(repoPath string, localFiles []string, idx *Index) (*DivergenceResult, error) {
	defaultBranch, err := detectDefaultBranch(repoPath)
	if err != nil {
		return nil, fmt.Errorf("CheckDivergence: %w", err)
	}

	ahead, behind, err := revListCount(repoPath, defaultBranch)
	if err != nil {
		return nil, fmt.Errorf("CheckDivergence: %w", err)
	}

	pullStrategy := detectPullStrategy(repoPath)

	// Cross-reference: find authors who committed to files on the behind
	// side that overlap with locally modified files
	localSet := make(map[string]bool, len(localFiles))
	for _, f := range localFiles {
		localSet[f] = true
	}

	seenAuthors := make(map[string]bool)
	var overlappingAuthors []AuthorContribution
	for file, authors := range idx.FileAuthors {
		if !localSet[file] {
			continue
		}
		for _, a := range authors {
			if !seenAuthors[a.Email] {
				seenAuthors[a.Email] = true
				overlappingAuthors = append(overlappingAuthors, a)
			}
		}
	}

	return &DivergenceResult{
		Ahead:                        ahead,
		Behind:                       behind,
		DefaultBranch:                defaultBranch,
		AuthorsOnBehindTouchingLocal: overlappingAuthors,
		PullStrategy:                 pullStrategy,
	}, nil
}

// detectDefaultBranch tries to determine the remote's default branch.
// Tries `git symbolic-ref refs/remotes/origin/HEAD` first, then falls
// back to checking for main/master.
func detectDefaultBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err == nil {
		ref := strings.TrimSpace(string(out))
		// refs/remotes/origin/main -> main
		parts := strings.Split(ref, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1], nil
		}
	}

	// Fallback: check if main or master exists
	for _, branch := range []string{"main", "master"} {
		// #nosec G204 -- branch is restricted to fixed list "main" / "master"
		cmd = exec.Command("git", "rev-parse", "--verify", "origin/"+branch)
		cmd.Dir = repoPath
		if err := cmd.Run(); err == nil {
			return branch, nil
		}
	}

	return "main", nil // safe default
}

// revListCount returns ahead/behind counts relative to origin/<branch>.
func revListCount(repoPath, branch string) (ahead, behind int, err error) {
	// #nosec G204 -- branch is detected from git symbolic-ref or sanitized default branch
	cmd := exec.Command("git", "rev-list", "--left-right", "--count",
		fmt.Sprintf("HEAD...origin/%s", branch))
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("rev-list: %w", err)
	}

	parts := strings.Fields(strings.TrimSpace(string(out)))
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected rev-list output: %q", string(out))
	}

	ahead, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("parsing ahead count: %w", err)
	}
	behind, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("parsing behind count: %w", err)
	}

	return ahead, behind, nil
}

// detectPullStrategy reads `git config pull.rebase` to determine the
// configured pull strategy.
func detectPullStrategy(repoPath string) string {
	cmd := exec.Command("git", "config", "pull.rebase")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return "merge" // default
	}

	val := strings.TrimSpace(string(out))
	if val == "true" || val == "1" {
		return "rebase"
	}
	return "merge"
}
