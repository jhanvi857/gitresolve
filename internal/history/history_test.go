package history

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestParseGitLogOutput(t *testing.T) {
	output := `COMMIT:abc123def|alice@example.com|1692000000
models/user.go
handlers/user.go
COMMIT:def456ghi|bob@example.com|1692100000
models/user.go
services/auth.go
COMMIT:ghi789jkl|alice@example.com|1692200000
models/user.go
handlers/user.go
services/auth.go
`

	records := ParseGitLogOutput(output)

	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}

	if records[0].sha != "abc123def" {
		t.Errorf("record[0].sha = %q, want %q", records[0].sha, "abc123def")
	}
	if records[0].authorEmail != "alice@example.com" {
		t.Errorf("record[0].authorEmail = %q, want %q", records[0].authorEmail, "alice@example.com")
	}
	if len(records[0].files) != 2 {
		t.Errorf("record[0].files = %v, want 2 files", records[0].files)
	}

	if records[2].sha != "ghi789jkl" {
		t.Errorf("record[2].sha = %q, want %q", records[2].sha, "ghi789jkl")
	}
	if len(records[2].files) != 3 {
		t.Errorf("record[2].files = %v, want 3 files", records[2].files)
	}
}

func TestCoChangeNormalization(t *testing.T) {
	idx := NewIndex(500)
	// Simulate: commits where models/user.go and handlers/user.go always change together
	records := []gitLogRecord{
		{sha: "a", authorEmail: "a@x.com", timestamp: time.Now(), files: []string{"models/user.go", "handlers/user.go"}},
		{sha: "b", authorEmail: "a@x.com", timestamp: time.Now(), files: []string{"models/user.go", "handlers/user.go"}},
		{sha: "c", authorEmail: "a@x.com", timestamp: time.Now(), files: []string{"models/user.go", "handlers/user.go"}},
		{sha: "d", authorEmail: "b@x.com", timestamp: time.Now(), files: []string{"models/user.go", "services/auth.go"}},
	}
	idx.BuildFromParsedLog(records)

	// models/user.go ↔ handlers/user.go should have count=3
	// models/user.go ↔ services/auth.go should have count=1
	cochanges := idx.NormalizedCoChanges()

	var userHandlerCC, userAuthCC *CoChange
	for i := range cochanges {
		if (cochanges[i].FileA == "handlers/user.go" && cochanges[i].FileB == "models/user.go") ||
			(cochanges[i].FileA == "models/user.go" && cochanges[i].FileB == "handlers/user.go") {
			userHandlerCC = &cochanges[i]
		}
		if (cochanges[i].FileA == "models/user.go" && cochanges[i].FileB == "services/auth.go") ||
			(cochanges[i].FileA == "services/auth.go" && cochanges[i].FileB == "models/user.go") {
			userAuthCC = &cochanges[i]
		}
	}

	if userHandlerCC == nil {
		t.Fatal("expected co-change entry for user.go ↔ handler/user.go")
	}
	if userHandlerCC.Count != 3 {
		t.Errorf("user-handler count = %d, want 3", userHandlerCC.Count)
	}
	if userHandlerCC.Strength != 1.0 {
		t.Errorf("user-handler strength = %f, want 1.0 (max)", userHandlerCC.Strength)
	}

	if userAuthCC == nil {
		t.Fatal("expected co-change entry for user.go ↔ auth.go")
	}
	if userAuthCC.Count != 1 {
		t.Errorf("user-auth count = %d, want 1", userAuthCC.Count)
	}
	// Strength should be 1/3 ≈ 0.333
	expectedStrength := 1.0 / 3.0
	if diff := userAuthCC.Strength - expectedStrength; diff > 0.01 || diff < -0.01 {
		t.Errorf("user-auth strength = %f, want ~%f", userAuthCC.Strength, expectedStrength)
	}
}

func TestCoChangeStrengthQuery(t *testing.T) {
	idx := NewIndex(500)
	records := []gitLogRecord{
		{sha: "a", authorEmail: "a@x.com", timestamp: time.Now(), files: []string{"a.go", "b.go"}},
		{sha: "b", authorEmail: "a@x.com", timestamp: time.Now(), files: []string{"a.go", "b.go"}},
		{sha: "c", authorEmail: "a@x.com", timestamp: time.Now(), files: []string{"a.go", "c.go"}},
	}
	idx.BuildFromParsedLog(records)

	// a.go ↔ b.go: count=2 (max)
	// a.go ↔ c.go: count=1
	s := idx.CoChangeStrength("a.go", "b.go")
	if s != 1.0 {
		t.Errorf("CoChangeStrength(a, b) = %f, want 1.0", s)
	}

	s = idx.CoChangeStrength("b.go", "a.go") // order shouldn't matter
	if s != 1.0 {
		t.Errorf("CoChangeStrength(b, a) = %f, want 1.0", s)
	}

	s = idx.CoChangeStrength("a.go", "c.go")
	if s != 0.5 {
		t.Errorf("CoChangeStrength(a, c) = %f, want 0.5", s)
	}

	s = idx.CoChangeStrength("x.go", "y.go")
	if s != 0 {
		t.Errorf("CoChangeStrength(x, y) = %f, want 0", s)
	}
}

func TestAuthorRecencyWeighting(t *testing.T) {
	idx := NewIndex(500)
	now := time.Now()
	records := []gitLogRecord{
		// alice: recent commit
		{sha: "a", authorEmail: "alice@x.com", timestamp: now.Add(-24 * time.Hour), files: []string{"main.go"}},
		// bob: old commit (180 days ago)
		{sha: "b", authorEmail: "bob@x.com", timestamp: now.Add(-180 * 24 * time.Hour), files: []string{"main.go"}},
		// alice: another recent commit
		{sha: "c", authorEmail: "alice@x.com", timestamp: now.Add(-2 * 24 * time.Hour), files: []string{"main.go"}},
	}
	idx.BuildFromParsedLog(records)

	authors := idx.AuthorsForFile("main.go")
	if len(authors) != 2 {
		t.Fatalf("expected 2 authors, got %d", len(authors))
	}

	// alice should be ranked higher due to recency + 2 commits
	if authors[0].Email != "alice@x.com" {
		t.Errorf("top author = %q, want alice@x.com", authors[0].Email)
	}
	if authors[0].Weight <= authors[1].Weight {
		t.Errorf("alice weight (%f) should be > bob weight (%f)", authors[0].Weight, authors[1].Weight)
	}
}

func TestCoupledFilesMissing(t *testing.T) {
	idx := NewIndex(500)
	records := []gitLogRecord{
		{sha: "a", authorEmail: "a@x.com", timestamp: time.Now(), files: []string{"models/user.go", "handlers/user.go"}},
		{sha: "b", authorEmail: "a@x.com", timestamp: time.Now(), files: []string{"models/user.go", "handlers/user.go"}},
		{sha: "c", authorEmail: "a@x.com", timestamp: time.Now(), files: []string{"models/user.go", "handlers/user.go"}},
	}
	idx.BuildFromParsedLog(records)

	// Touch models/user.go but NOT handlers/user.go
	touched := map[string]bool{"models/user.go": true}
	missing := idx.CoupledFilesMissing("models/user.go", touched, 0.5)

	if len(missing) != 1 {
		t.Fatalf("expected 1 missing coupled file, got %d", len(missing))
	}

	found := false
	for _, cc := range missing {
		if cc.FileA == "handlers/user.go" || cc.FileB == "handlers/user.go" {
			found = true
		}
	}
	if !found {
		t.Error("expected handlers/user.go to be in missing coupled files")
	}
}

func TestTarjanCycleDetection(t *testing.T) {
	tests := []struct {
		name     string
		graph    map[string][]string
		wantLen  int
		wantCyc  []string // sorted members of expected cycle (or nil)
	}{
		{
			name: "simple cycle A→B→C→A",
			graph: map[string][]string{
				"pkg/a": {"pkg/b"},
				"pkg/b": {"pkg/c"},
				"pkg/c": {"pkg/a"},
			},
			wantLen: 1,
			wantCyc: []string{"pkg/a", "pkg/b", "pkg/c"},
		},
		{
			name: "no cycle (DAG)",
			graph: map[string][]string{
				"pkg/a": {"pkg/b"},
				"pkg/b": {"pkg/c"},
			},
			wantLen: 0,
		},
		{
			name: "self-loop",
			graph: map[string][]string{
				"pkg/a": {"pkg/a"},
			},
			wantLen: 1,
			wantCyc: []string{"pkg/a"},
		},
		{
			name: "two separate cycles",
			graph: map[string][]string{
				"a": {"b"},
				"b": {"a"},
				"c": {"d"},
				"d": {"c"},
			},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cycles := tarjanSCC(tt.graph)
			if len(cycles) != tt.wantLen {
				t.Errorf("got %d cycles, want %d: %v", len(cycles), tt.wantLen, cycles)
			}
			if tt.wantCyc != nil && len(cycles) > 0 {
				got := strings.Join([]string(cycles[0]), ", ")
				want := strings.Join(tt.wantCyc, ", ")
				if got != want {
					t.Errorf("cycle = [%s], want [%s]", got, want)
				}
			}
		})
	}
}

func TestFileInImportCycle(t *testing.T) {
	idx := NewIndex(500)
	idx.ImportEdges = map[string][]string{
		"pkg/auth":  {"pkg/user"},
		"pkg/user":  {"pkg/auth"},
		"pkg/utils": {"fmt"},
	}

	inCycle, cycle := idx.FileInImportCycle("pkg/auth/token.go")
	if !inCycle {
		t.Error("expected pkg/auth/token.go to be in a cycle")
	}
	if len(cycle) != 2 {
		t.Errorf("cycle length = %d, want 2", len(cycle))
	}

	inCycle, _ = idx.FileInImportCycle("pkg/utils/helper.go")
	if inCycle {
		t.Error("expected pkg/utils/helper.go NOT to be in a cycle")
	}
}

func TestBuildFromParsedLogEmpty(t *testing.T) {
	idx := NewIndex(500)
	idx.BuildFromParsedLog(nil)

	if idx.TotalCommits() != 0 {
		t.Errorf("TotalCommits = %d, want 0", idx.TotalCommits())
	}
	if len(idx.NormalizedCoChanges()) != 0 {
		t.Error("expected no co-changes from empty log")
	}
}

func TestParseGitLogOutputEmpty(t *testing.T) {
	records := ParseGitLogOutput("")
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}
}

func TestParseGitLogOutputMalformed(t *testing.T) {
	// Malformed COMMIT line — missing fields
	output := `COMMIT:abc123
somefile.go
`
	records := ParseGitLogOutput(output)
	// Should still parse with partial data (sha only, no author/timestamp)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].sha != "abc123" {
		t.Errorf("sha = %q, want %q", records[0].sha, "abc123")
	}
	if records[0].authorEmail != "" {
		t.Errorf("authorEmail = %q, want empty", records[0].authorEmail)
	}
	if len(records[0].files) != 1 || records[0].files[0] != "somefile.go" {
		t.Errorf("files = %v, want [somefile.go]", records[0].files)
	}
}

func TestIndexTotalCommits(t *testing.T) {
	idx := NewIndex(500)
	records := make([]gitLogRecord, 42)
	for i := range records {
		records[i] = gitLogRecord{
			sha:         fmt.Sprintf("sha%d", i),
			authorEmail: "dev@x.com",
			timestamp:   time.Now(),
			files:       []string{"file.go"},
		}
	}
	idx.BuildFromParsedLog(records)

	if idx.TotalCommits() != 42 {
		t.Errorf("TotalCommits = %d, want 42", idx.TotalCommits())
	}
}

func TestExtractSymbolsFromLines(t *testing.T) {
	lines := []string{
		"func ProcessOrder(order Order) error {",
		"	total := calculateTotal(order)",
		"	return nil",
		"}",
		"func (s *Service) RefundPayment(id string) error {",
	}

	symbols := ExtractSymbolsFromLines(lines)
	if len(symbols) != 2 {
		t.Fatalf("expected 2 symbols, got %d: %v", len(symbols), symbols)
	}
	if symbols[0] != "ProcessOrder" || symbols[1] != "RefundPayment" {
		t.Errorf("got symbols %v, want [ProcessOrder RefundPayment]", symbols)
	}
}

func TestCallersOfSymbol(t *testing.T) {
	idx := NewIndex(500)
	idx.SymbolCallers[SymbolID{Name: "ProcessPayment"}] = []CallerRef{
		{File: "handlers/checkout.go", Line: 42, Name: "CheckoutHandler"},
		{File: "services/billing.go", Line: 88, Name: "BillCustomer"},
	}

	callers := idx.CallersOfSymbol("ProcessPayment")
	if len(callers) != 2 {
		t.Fatalf("expected 2 callers, got %d", len(callers))
	}
	if callers[0].Name != "CheckoutHandler" || callers[1].Name != "BillCustomer" {
		t.Errorf("unexpected callers: %v", callers)
	}
}

