package conflict

import (
	"strings"
)

func ParseFile(filePath string, content []byte) []*Conflict {
	lines := strings.Split(string(content), "\n")
	var conflicts []*Conflict

	var inConflict bool
	var currentConflict *Conflict
	var side int
	for i, line := range lines {
		if strings.HasPrefix(line, "<<<<<<<") {
			inConflict = true
			side = 1
			currentConflict = &Conflict{
				FilePath:  filePath,
				StartLine: i + 1,
			}
			continue
		} else if strings.HasPrefix(line, "|||||||") {
			side = 2
			continue
		} else if strings.HasPrefix(line, "=======") {
			if inConflict {
				side = 3
			}
			continue
		} else if strings.HasPrefix(line, ">>>>>>>") {
			inConflict = false
			side = 0
			if currentConflict != nil {
				currentConflict.EndLine = i + 1
				conflicts = append(conflicts, currentConflict)
				currentConflict = nil
			}
			continue
		}

		if inConflict && currentConflict != nil {
			if side == 1 {
				currentConflict.OurLines = append(currentConflict.OurLines, line)
			} else if side == 2 {
				currentConflict.BaseLines = append(currentConflict.BaseLines, line)
			} else if side == 3 {
				currentConflict.TheirLines = append(currentConflict.TheirLines, line)
			}
		}
	}

	return conflicts
}

func CompileResolution(content []byte, conflicts []*Conflict) string {
	lines := strings.Split(string(content), "\n")
	var result []string

	cIdx := 0
	inConflict := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if strings.HasPrefix(line, "<<<<<<<") {
			inConflict = true
			if cIdx < len(conflicts) {
				if conflicts[cIdx].Resolution != "" {
					result = append(result, conflicts[cIdx].Resolution)
				} else {
					result = append(result, line)
				}
			}
			continue
		} else if strings.HasPrefix(line, ">>>>>>>") {
			inConflict = false
			if cIdx < len(conflicts) && conflicts[cIdx].Resolution == "" {
				result = append(result, line)
			}
			cIdx++
			continue
		}

		if inConflict {
			if cIdx < len(conflicts) && conflicts[cIdx].Resolution == "" {
				result = append(result, line)
			}
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
