package merge

import (
	"fmt"

	"github.com/jhanvi857/gitresolve/internal/git"
)

func findMergeBase(r *git.Repository, ourSHA, theirSHA string) (string, error) {
	// collect all ancestors of our branch into a set
	ourAncestors := make(map[string]bool)
	current := ourSHA
	for current != "" {
		ourAncestors[current] = true
		commit, err := git.GetCommit(r, current)
		if err != nil {
			break
		}
		current = commit.Parent
	}
	// walk their branch backwards until we hit a commit that exists in ourAncestors
	// first match is the merge base
	current = theirSHA
	for current != "" {
		if ourAncestors[current] {
			return current, nil
		}
		commit, err := git.GetCommit(r, current)
		if err != nil {
			break
		}
		current = commit.Parent
	}

	return "", fmt.Errorf("FindMergeBase: %w", fmt.Errorf("no common ancestor found"))
}
