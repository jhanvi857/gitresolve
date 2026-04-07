package conflict

import (
	"strings"
)

func countBraces(line string) int {
	depth := 0
	for _, char := range line {
		if char == '{' {
			depth++
		} else if char == '}' {
			depth--
		}
	}
	return depth
}

func ParseFile(filePath string, content []byte) []*ConflictBlock {
	lines := strings.Split(string(content), "\n")
	var conflicts []*ConflictBlock

	var currentDepth int
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "<<<<<<<") {
			cb := &ConflictBlock{
				FilePath:  filePath,
				StartLine: i + 1,
				PreLines:  append([]string{}, lines[:i]...), // Deep copy suggested
			}
			baselineDepth := currentDepth

			// Track depth for the marker itself if it contains braces? (unlikely but safe)
			// Actually, just ignore markers for depth.
			
			i++
			side := 1 // 1=ours, 2=base, 3=theirs
			for i < len(lines) {
				currLine := lines[i]
				if strings.HasPrefix(currLine, "|||||||") {
					side = 2
					i++
					continue
				} else if strings.HasPrefix(currLine, "=======") {
					side = 3
					i++
					continue
				} else if strings.HasPrefix(currLine, ">>>>>>>") {
					// Don't break yet, need to check depth after consuming THEIRS
					break
				}
				
				if side == 1 {
					cb.OursLines = append(cb.OursLines, currLine)
				} else if side == 2 {
					cb.BaseLines = append(cb.BaseLines, currLine)
				} else if side == 3 {
					cb.TheirsLines = append(cb.TheirsLines, currLine)
				}
				
				// We track depth for ALL lines in the conflict to know the balance
				currentDepth += countBraces(currLine)
				i++
			}
			
			// We hit >>>>>>> at index i
			if i < len(lines) && strings.HasPrefix(lines[i], ">>>>>>>") {
				i++
				// Continue consuming lines until depth matches baselineDepth
				for i < len(lines) {
					if currentDepth == baselineDepth {
						break
					}
					cb.TheirsLines = append(cb.TheirsLines, lines[i])
					currentDepth += countBraces(lines[i])
					i++
				}
				cb.EndLine = i 
			}
			
			cb.PostLines = append([]string{}, lines[i:]...)
			conflicts = append(conflicts, cb)
			
			// Stay at index i-1 for the next loop iteration because i is advanced?
            // Actually i is at the first line of PostLines. The loop increment will do i++.
            // So we should do i--.
            i--
		} else {
			currentDepth += countBraces(line)
		}
	}

	return conflicts
}

func CompileResolution(content []byte, conflicts []*ConflictBlock) string {
	// If the user wants the formula: output = PreLines + OursLines + PostLines
	// And if multiple conflicts are present, they assume they are resolved sequentially?
	
	// If we have one conflict, we can just use the first one's Resolution if it's the full file.
	// But usually CompileResolution joins them.
	
	// Let's stick to a simple join if Resolution is just the block resolution.
    // BUT the prompt says "[O]urs: output = PreLines + OursLines + PostLines".
    // This implies c.Resolution is the FULL FILE content.
    
    // If c.Resolution is the full file content, then CompileResolution is just returning that.
    if len(conflicts) > 0 && conflicts[len(conflicts)-1].Resolution != "" {
        // If they were resolved sequentially, the last one has the final state.
        // Wait, that's not safe.
    }
    
    // Let's re-read Resolve logic.
    
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
                    // Split the resolution back into lines if it's a full file?
                    // No, usually Resolution is just the replacement part.
                    // BUT the prompt says resolution should produce FULL FILE.
                    
                    // Actually, if I follow the prompt's Resolve formula, c.Resolution will be the full file.
                    // If so, CompileResolution could just return c.Resolution.
                    
					result = append(result, conflicts[cIdx].Resolution)
                    // Skip till EndLine
                    i = conflicts[cIdx].EndLine - 1
                    inConflict = false
                    cIdx++
                    continue
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
