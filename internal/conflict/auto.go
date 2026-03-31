package conflict

import (
	"path/filepath"
	"strings"

	"github.com/jhanvi857/gitresolve/internal/analysis"
)

type Options struct {
	NoAutoStructured bool
}

func AutoResolve(c *Conflict, opts Options) bool {
	if !c.CanAutoResolve {
		return false
	}

	switch c.Type {
	case TypeStructured:
		if opts.NoAutoStructured {
			return false
		}
		ext := filepath.Ext(c.FilePath)
		baseBytes := []byte(strings.Join(c.BaseLines, "\n"))
		ourBytes := []byte(strings.Join(c.OurLines, "\n"))
		theirBytes := []byte(strings.Join(c.TheirLines, "\n"))

		var res analysis.StructuredMergeResult
		var err error

		if ext == ".json" {
			res, err = analysis.MergeJSON(baseBytes, ourBytes, theirBytes)
		} else if ext == ".yaml" || ext == ".yml" {
			res, err = analysis.MergeYAML(baseBytes, ourBytes, theirBytes)
		} else if ext == ".toml" {
			res, err = analysis.MergeTOML(baseBytes, ourBytes, theirBytes)
		}

		if err == nil && !res.HasConflicts {
			c.Resolution = res.Content
			return true
		}

	case TypeWhitespace:
		c.Resolution = strings.Join(c.OurLines, "\n")
		return true

	case TypeImport:
		merged := mergeImports(c.OurLines, c.TheirLines)
		c.Resolution = strings.Join(merged, "\n")
		return true

	case TypeIdentical:
		c.Resolution = strings.Join(c.OurLines, "\n")
		return true
	}

	return false
}

func mergeImports(ours, theirs []string) []string {
	seen := make(map[string]bool)
	var merged []string

	for _, line := range ours {
		if !seen[line] {
			seen[line] = true
			merged = append(merged, line)
		}
	}

	for _, line := range theirs {
		if !seen[line] {
			seen[line] = true
			merged = append(merged, line)
		}
	}

	return merged
}
