package merge

import (
	"fmt"
	"strings"

	"github.com/jhanvi857/gitresolve/internal/git"
)

type MergeResult struct {
	FilePath string
	// final merged content if no conflicts
	Content      string
	HasConflicts bool
	Hunks        []git.Hunk
}

// MergeFile attempts to merge three versions of a file : returns clean content if auto-mergeable conflict markers if not
func MergeFile(filePath string, base, ours, theirs []byte) (MergeResult, error) {
	hunks, err := git.ThreeWayDiff(base, ours, theirs)
	if err != nil {
		return MergeResult{}, fmt.Errorf("MergeFile: %w", err)
	}

	result := MergeResult{
		FilePath: filePath,
		Hunks:    hunks,
	}

	baseLines := splitContent(base)
	var output []string
	i := 0

	for _, hunk := range hunks {
		// add unchanged lines before this hunk
		for i < hunk.StartLine {
			output = append(output, baseLines[i])
			i++
		}

		switch hunk.Kind {
		case git.HunkOursOnly:
			// only our side changed : take ours
			output = append(output, hunk.OurLines...)

		case git.HunkTheirsOnly:
			// only their side changed : take theirs
			output = append(output, hunk.TheirLines...)

		case git.HunkIdentical:
			// both sides made same change : take either
			output = append(output, hunk.OurLines...)

		case git.HunkConflict:
			// both sides changed differently : write conflict markers
			result.HasConflicts = true
			output = append(output, "<<<<<<< ours")
			output = append(output, hunk.OurLines...)
			output = append(output, "=======")
			output = append(output, hunk.TheirLines...)
			output = append(output, ">>>>>>> theirs")
		}

		i = hunk.EndLine + 1
	}

	// add any remaining unchanged lines after last hunk
	for i < len(baseLines) {
		output = append(output, baseLines[i])
		i++
	}

	result.Content = strings.Join(output, "\n")
	return result, nil
}

func splitContent(content []byte) []string {
	if len(content) == 0 {
		return []string{}
	}
	return strings.Split(string(content), "\n")
}
