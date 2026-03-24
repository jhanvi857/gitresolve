package git

import "fmt"

// HunkKind describes what kind of change a hunk represents
type HunkKind int

const (
	HunkIdentical HunkKind = iota // both sides unchanged
	HunkOursOnly
	HunkTheirsOnly
	HunkConflict
)

// Hunk is one region of change between file versions
type Hunk struct {
	StartLine  int
	EndLine    int
	Kind       HunkKind
	BaseLines  []string
	OurLines   []string
	TheirLines []string
}

// ThreeWayDiff compares base/ours/theirs and returns all changed regions
// this function splits content into lines and orchestrates the comparison
func ThreeWayDiff(base, ours, theirs []byte) ([]Hunk, error) {
	baseLines := splitLines(base)
	ourLines := splitLines(ours)
	theirLines := splitLines(theirs)

	if len(baseLines) == 0 {
		return nil, fmt.Errorf("ThreeWayDiff: base content is empty")
	}

	// compare base against ours and base against theirs separately
	// then combine results to find true conflicts
	ourChanges := diffLines(baseLines, ourLines)
	theirChanges := diffLines(baseLines, theirLines)

	hunks := combineChanges(baseLines, ourLines, theirLines, ourChanges, theirChanges)

	return hunks, nil
}

// FileDiff is a simple two-way diff between any two file versions
func FileDiff(a, b []byte) ([]Hunk, error) {
	aLines := splitLines(a)
	bLines := splitLines(b)

	changes := diffLines(aLines, bLines)

	var hunks []Hunk
	for _, c := range changes {
		hunks = append(hunks, Hunk{
			StartLine: c.start,
			EndLine:   c.end,
			Kind:      HunkOursOnly,
			OurLines:  c.lines,
		})
	}

	return hunks, nil
}

// splitLines converts raw file bytes into a slice of strings
// each element is one line without the newline character
func splitLines(content []byte) []string {
	if len(content) == 0 {
		return []string{}
	}

	var lines []string
	start := 0

	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			lines = append(lines, string(content[start:i]))
			start = i + 1
		}
	}

	// handle last line if file does not end with newline
	if start < len(content) {
		lines = append(lines, string(content[start:]))
	}

	return lines
}

// change represents one set of modified lines from diffLines
type change struct {
	start int
	end   int
	lines []string
}

// diffLines finds which line ranges changed between a and b
// this is a simplified diff that myers.go will replace with full LCS algorithm
// returns list of changes describing what is different in b compared to a
func diffLines(a, b []string) []change {
	var changes []change

	// built a map of line content to line numbers in a
	// lets quickly check if a line from b exists in a
	lineIndex := make(map[string]bool)
	for _, line := range a {
		lineIndex[line] = true
	}

	i := 0
	for i < len(b) {
		// if this line from b does not exist in a it is a change
		if !lineIndex[b[i]] {
			start := i
			var changedLines []string

			// collect consecutive changed lines
			for i < len(b) && !lineIndex[b[i]] {
				changedLines = append(changedLines, b[i])
				i++
			}

			changes = append(changes, change{
				start: start,
				end:   i - 1,
				lines: changedLines,
			})
			continue
		}
		i++
	}

	return changes
}

// combineChanges takes our changes and their changes and produces final hunks
// if only our side changed a region = HunkOursOnly
// if only their side changed a region = HunkTheirsOnly
// if both sides changed same region differently = HunkConflict
func combineChanges(base, ours, theirs []string, ourChanges, theirChanges []change) []Hunk {
	var hunks []Hunk

	// build sets of line numbers changed on each side
	ourChanged := make(map[int]bool)
	theirChanged := make(map[int]bool)

	for _, c := range ourChanges {
		for i := c.start; i <= c.end; i++ {
			ourChanged[i] = true
		}
	}

	for _, c := range theirChanges {
		for i := c.start; i <= c.end; i++ {
			theirChanged[i] = true
		}
	}

	// walk through base lines and classify each region
	i := 0
	for i < len(base) {
		inOurs := ourChanged[i]
		inTheirs := theirChanged[i]

		if !inOurs && !inTheirs {
			i++
			continue
		}

		// collecting full region
		start := i
		var baseRegion, ourRegion, theirRegion []string

		for i < len(base) && (ourChanged[i] || theirChanged[i]) {
			baseRegion = append(baseRegion, base[i])

			if i < len(ours) {
				ourRegion = append(ourRegion, ours[i])
			}
			if i < len(theirs) {
				theirRegion = append(theirRegion, theirs[i])
			}
			i++
		}

		// classifying this region
		kind := classifyHunk(inOurs, inTheirs, ourRegion, theirRegion)

		hunks = append(hunks, Hunk{
			StartLine:  start,
			EndLine:    i - 1,
			Kind:       kind,
			BaseLines:  baseRegion,
			OurLines:   ourRegion,
			TheirLines: theirRegion,
		})
	}

	return hunks
}

// classifyHunk decides what kind of conflict a region is
func classifyHunk(inOurs, inTheirs bool, ourLines, theirLines []string) HunkKind {
	if inOurs && !inTheirs {
		return HunkOursOnly
	}

	if inTheirs && !inOurs {
		return HunkTheirsOnly
	}

	// both sides changed this region
	// check if they ended up with the same content
	if linesEqual(ourLines, theirLines) {
		return HunkIdentical
	}

	return HunkConflict
}

// linesEqual checks if two line slices have identical content
func linesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
