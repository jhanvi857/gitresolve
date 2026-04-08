package conflict

import (
	"sort"
	"strings"
)

func ParseFile(filePath string, content []byte) []*ConflictBlock {
	lines := strings.Split(string(content), "\n")
	var conflicts []*ConflictBlock

	for i := 0; i < len(lines); i++ {
		if !strings.HasPrefix(lines[i], "<<<<<<<") {
			continue
		}

		start := i
		cb := &ConflictBlock{
			FilePath:    filePath,
			StartLine:   start + 1,
			StartIndex:  start,
			PreLines:    append([]string{}, lines[:start]...),
			OursLines:   []string{},
			BaseLines:   []string{},
			TheirsLines: []string{},
		}

		i++
		side := 1 // 1=ours, 2=base, 3=theirs
		foundEnd := false
		for i < len(lines) {
			currLine := lines[i]
			if strings.HasPrefix(currLine, "|||||||") {
				side = 2
				i++
				continue
			}
			if strings.HasPrefix(currLine, "=======") {
				side = 3
				i++
				continue
			}
			if strings.HasPrefix(currLine, ">>>>>>>") {
				cb.EndIndex = i
				cb.EndLine = i + 1
				foundEnd = true
				break
			}

			switch side {
			case 1:
				cb.OursLines = append(cb.OursLines, currLine)
			case 2:
				cb.BaseLines = append(cb.BaseLines, currLine)
			case 3:
				cb.TheirsLines = append(cb.TheirsLines, currLine)
			}
			i++
		}

		if !foundEnd {
			continue
		}

		postStart := cb.EndIndex + 1
		if postStart < 0 || postStart > len(lines) {
			postStart = len(lines)
		}
		cb.PostLines = append([]string{}, lines[postStart:]...)
		conflicts = append(conflicts, cb)
	}

	return conflicts
}

func CompileResolution(content []byte, conflicts []*ConflictBlock) string {
	if len(conflicts) == 0 {
		return string(content)
	}

	lines := strings.Split(string(content), "\n")
	sortedConflicts := append([]*ConflictBlock{}, conflicts...)
	sort.Slice(sortedConflicts, func(i, j int) bool {
		return sortedConflicts[i].StartIndex < sortedConflicts[j].StartIndex
	})

	result := make([]string, 0, len(lines))
	cursor := 0
	for _, c := range sortedConflicts {
		if c.StartIndex < cursor || c.StartIndex >= len(lines) {
			continue
		}
		if c.EndIndex < c.StartIndex || c.EndIndex >= len(lines) {
			continue
		}

		result = append(result, lines[cursor:c.StartIndex]...)
		if c.Resolution != "" {
			result = append(result, strings.Split(c.Resolution, "\n")...)
		} else {
			result = append(result, lines[c.StartIndex:c.EndIndex+1]...)
		}
		cursor = c.EndIndex + 1
	}

	if cursor < len(lines) {
		result = append(result, lines[cursor:]...)
	}

	return strings.Join(result, "\n")
}
