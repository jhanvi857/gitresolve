package history

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// BenchmarkHistoryIndex measures the time to build a history index from
// a synthetic git log representing ~5000 commits / ~2000 files.
// Target: < 200ms on commodity hardware.
func BenchmarkHistoryIndex(b *testing.B) {
	records := generateSyntheticLog(5000, 2000)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		idx := NewIndex(5000)
		idx.BuildFromParsedLog(records)
		_ = idx.NormalizedCoChanges()
		// Import cycle detection on a synthetic graph
		for j := 0; j < 200; j++ {
			pkg := fmt.Sprintf("pkg/%d", j)
			dep := fmt.Sprintf("pkg/%d", (j+1)%200)
			idx.ImportEdges[pkg] = append(idx.ImportEdges[pkg], dep)
		}
		_ = idx.ImportCycles()
	}
}

// generateSyntheticLog creates a realistic synthetic git log.
// Each commit touches 1–8 files drawn from a pool of numFiles files.
func generateSyntheticLog(numCommits, numFiles int) []gitLogRecord {
	// #nosec G404 -- deterministic pseudo-random seed used strictly for benchmarking
	rng := rand.New(rand.NewSource(42))
	files := make([]string, numFiles)
	for i := 0; i < numFiles; i++ {
		files[i] = fmt.Sprintf("pkg/%d/file_%d.go", i%50, i)
	}

	authors := []string{
		"alice@example.com", "bob@example.com", "carol@example.com",
		"dave@example.com", "eve@example.com",
	}

	records := make([]gitLogRecord, numCommits)
	baseTime := time.Now().Add(-365 * 24 * time.Hour)

	for i := 0; i < numCommits; i++ {
		numTouched := 1 + rng.Intn(8)
		touched := make([]string, 0, numTouched)
		seen := make(map[int]bool)
		for j := 0; j < numTouched; j++ {
			idx := rng.Intn(numFiles)
			if !seen[idx] {
				seen[idx] = true
				touched = append(touched, files[idx])
			}
		}

		records[i] = gitLogRecord{
			sha:         fmt.Sprintf("%040x", rng.Int63()),
			authorEmail: authors[rng.Intn(len(authors))],
			timestamp:   baseTime.Add(time.Duration(i) * time.Hour),
			files:       touched,
		}
	}

	return records
}
