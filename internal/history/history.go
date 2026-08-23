// Package history provides a lightweight in-memory index of git history
// facts (co-changes, file authors, import edges) for use in rule-based
// escalation decisions. The index is built per-run from git log output
// and optionally persisted to SQLite for incremental updates.
//
// This package deliberately avoids any ML, LLM, or external service
// dependency — all data is derived from the repo's own commit history
// and AST structure.
package history

import (
	"bufio"
	"fmt"
	"math"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// SymbolID uniquely identifies a symbol within the codebase.
type SymbolID struct {
	File string
	Name string
}

// CallerRef records one call-site that references a symbol.
type CallerRef struct {
	File string
	Line int
	Name string
}

// AuthorContribution tracks a single author's weighted contribution to a file.
type AuthorContribution struct {
	Email       string
	Weight      float64   // recency-weighted contribution (0–1)
	LastTouched time.Time // most recent commit touching this file by this author
}

// CoChange records the co-change coupling between two files.
type CoChange struct {
	FileA    string
	FileB    string
	Count    int
	Strength float64 // normalized 0–1
}

// Cycle is an ordered list of packages forming an import cycle.
type Cycle []string

// Index is the in-memory per-run history index.
// It is built from git log output and optionally from existing AST data.
type Index struct {
	SymbolCallers map[SymbolID][]CallerRef
	FileCoChanges map[string]map[string]int // file -> co-changed file -> raw count
	FileAuthors   map[string][]AuthorContribution
	ImportEdges   map[string][]string // package/file -> imported packages
	totalCommits  int
	maxCommits    int
	newHeadSHA    string // the HEAD SHA after building (for sync state persistence)
}

// NewIndex creates an empty Index with the given max commit limit.
func NewIndex(maxCommits int) *Index {
	if maxCommits <= 0 {
		maxCommits = 500
	}
	return &Index{
		SymbolCallers: make(map[SymbolID][]CallerRef),
		FileCoChanges: make(map[string]map[string]int),
		FileAuthors:   make(map[string][]AuthorContribution),
		ImportEdges:   make(map[string][]string),
		maxCommits:    maxCommits,
	}
}

// NewHeadSHA returns the HEAD SHA at the time the index was built.
func (idx *Index) NewHeadSHA() string { return idx.newHeadSHA }

// TotalCommits returns the number of commits processed during index build.
func (idx *Index) TotalCommits() int { return idx.totalCommits }

// gitLogRecord is a parsed record from git log --numstat --format output.
type gitLogRecord struct {
	sha         string
	authorEmail string
	timestamp   time.Time
	files       []string
}

// BuildFromGitLog shells out to `git log` and fills FileCoChanges and
// FileAuthors. If lastProcessedSHA is non-empty, only commits after that
// SHA are processed (incremental update).
func (idx *Index) BuildFromGitLog(repoPath string, lastProcessedSHA string) error {
	records, err := idx.runGitLog(repoPath, lastProcessedSHA)
	if err != nil {
		return fmt.Errorf("BuildFromGitLog: %w", err)
	}

	idx.processRecords(records)
	return nil
}

// BuildFromParsedLog populates the index from pre-parsed records.
// Used in tests to avoid shelling out to git.
func (idx *Index) BuildFromParsedLog(records []gitLogRecord) {
	idx.processRecords(records)
}

func (idx *Index) processRecords(records []gitLogRecord) {
	now := time.Now()
	// Track per-file author stats for recency weighting
	type authorStat struct {
		linesChanged int
		lastTouched  time.Time
		commitCount  int
	}
	fileAuthorStats := make(map[string]map[string]*authorStat)

	for _, rec := range records {
		idx.totalCommits++

		// Co-change: every pair of files in this commit
		for i := 0; i < len(rec.files); i++ {
			for j := i + 1; j < len(rec.files); j++ {
				a, b := rec.files[i], rec.files[j]
				if a > b {
					a, b = b, a // canonical order
				}
				if idx.FileCoChanges[a] == nil {
					idx.FileCoChanges[a] = make(map[string]int)
				}
				idx.FileCoChanges[a][b]++
			}
		}

		// Author contributions
		for _, file := range rec.files {
			if fileAuthorStats[file] == nil {
				fileAuthorStats[file] = make(map[string]*authorStat)
			}
			stat := fileAuthorStats[file][rec.authorEmail]
			if stat == nil {
				stat = &authorStat{}
				fileAuthorStats[file][rec.authorEmail] = stat
			}
			stat.commitCount++
			if rec.timestamp.After(stat.lastTouched) {
				stat.lastTouched = rec.timestamp
			}
		}
	}

	// Convert author stats to recency-weighted contributions
	for file, authors := range fileAuthorStats {
		var contribs []AuthorContribution
		for email, stat := range authors {
			// Recency weight: exponential decay over 90 days
			daysSince := now.Sub(stat.lastTouched).Hours() / 24
			recencyWeight := math.Exp(-daysSince / 90.0)
			weight := float64(stat.commitCount) * recencyWeight

			contribs = append(contribs, AuthorContribution{
				Email:       email,
				Weight:      weight,
				LastTouched: stat.lastTouched,
			})
		}
		// Sort by weight descending
		sort.Slice(contribs, func(i, j int) bool {
			return contribs[i].Weight > contribs[j].Weight
		})
		idx.FileAuthors[file] = contribs
	}
}

// runGitLog executes git log and parses output into records.
func (idx *Index) runGitLog(repoPath string, lastProcessedSHA string) ([]gitLogRecord, error) {
	args := []string{
		"log",
		"--format=COMMIT:%H|%ae|%at",
		"--name-only",
		fmt.Sprintf("-%d", idx.maxCommits),
	}

	if lastProcessedSHA != "" {
		args = append(args, lastProcessedSHA+"..HEAD")
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath

	out, err := cmd.Output()
	if err != nil {
		// If the range is invalid (e.g., first run), fall back to full log
		if lastProcessedSHA != "" {
			args = args[:len(args)-1] // remove the range
			cmd = exec.Command("git", args...)
			cmd.Dir = repoPath
			out, err = cmd.Output()
			if err != nil {
				return nil, fmt.Errorf("git log: %w", err)
			}
		} else {
			return nil, fmt.Errorf("git log: %w", err)
		}
	}

	records := ParseGitLogOutput(string(out))

	// Record the newest SHA for sync state
	if len(records) > 0 {
		idx.newHeadSHA = records[0].sha
	}

	return records, nil
}

// ParseGitLogOutput parses the output of git log --format=COMMIT:%H|%ae|%at
// --name-only into structured records. Exported for testing.
func ParseGitLogOutput(output string) []gitLogRecord {
	var records []gitLogRecord
	var current *gitLogRecord

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "COMMIT:") {
			if current != nil && current.sha != "" {
				records = append(records, *current)
			}
			current = &gitLogRecord{}
			parts := strings.SplitN(strings.TrimPrefix(line, "COMMIT:"), "|", 3)
			if len(parts) >= 1 {
				current.sha = parts[0]
			}
			if len(parts) >= 2 {
				current.authorEmail = parts[1]
			}
			if len(parts) >= 3 {
				if ts, err := strconv.ParseInt(parts[2], 10, 64); err == nil {
					current.timestamp = time.Unix(ts, 0)
				}
			}
		} else if current != nil {
			// This is a filename from --name-only
			current.files = append(current.files, line)
		}
	}

	// Don't forget the last record
	if current != nil && current.sha != "" {
		records = append(records, *current)
	}

	return records
}

// NormalizedCoChanges converts raw co-change counts into normalized 0–1
// strength scores. Strength = count / max_count across all pairs.
func (idx *Index) NormalizedCoChanges() []CoChange {
	maxCount := 0
	for _, neighbors := range idx.FileCoChanges {
		for _, count := range neighbors {
			if count > maxCount {
				maxCount = count
			}
		}
	}

	if maxCount == 0 {
		return nil
	}

	var result []CoChange
	for fileA, neighbors := range idx.FileCoChanges {
		for fileB, count := range neighbors {
			result = append(result, CoChange{
				FileA:    fileA,
				FileB:    fileB,
				Count:    count,
				Strength: float64(count) / float64(maxCount),
			})
		}
	}

	// Sort for deterministic output
	sort.Slice(result, func(i, j int) bool {
		if result[i].FileA != result[j].FileA {
			return result[i].FileA < result[j].FileA
		}
		return result[i].FileB < result[j].FileB
	})

	return result
}

// CoChangeStrength returns the normalized co-change strength between two files.
// Returns 0 if no coupling exists.
func (idx *Index) CoChangeStrength(fileA, fileB string) float64 {
	if fileA > fileB {
		fileA, fileB = fileB, fileA
	}

	count := 0
	if neighbors, ok := idx.FileCoChanges[fileA]; ok {
		count = neighbors[fileB]
	}

	if count == 0 {
		return 0
	}

	// Compute max for normalization
	maxCount := 0
	for _, neighbors := range idx.FileCoChanges {
		for _, c := range neighbors {
			if c > maxCount {
				maxCount = c
			}
		}
	}

	if maxCount == 0 {
		return 0
	}

	return float64(count) / float64(maxCount)
}

// AuthorsForFile returns the recency-weighted author contributions for a file.
func (idx *Index) AuthorsForFile(file string) []AuthorContribution {
	return idx.FileAuthors[file]
}

// CallersOf returns all call-sites referencing the given symbol.
func (idx *Index) CallersOf(sym SymbolID) []CallerRef {
	return idx.SymbolCallers[sym]
}

// CoupledFilesMissing checks if any file strongly coupled with touchedFile
// (strength >= minStrength) is absent from touchedFiles set.
func (idx *Index) CoupledFilesMissing(touchedFile string, touchedFiles map[string]bool, minStrength float64) []CoChange {
	var missing []CoChange
	cochanges := idx.NormalizedCoChanges()

	for _, cc := range cochanges {
		if cc.Strength < minStrength {
			continue
		}

		var coupled string
		if cc.FileA == touchedFile {
			coupled = cc.FileB
		} else if cc.FileB == touchedFile {
			coupled = cc.FileA
		} else {
			continue
		}

		if !touchedFiles[coupled] {
			missing = append(missing, cc)
		}
	}

	return missing
}
