package conflict

import (
	"strings"
)

func AutoResolve(c *Conflict) bool {
	if !c.CanAutoResolve {
		return false
	}

	switch c.Type {
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
