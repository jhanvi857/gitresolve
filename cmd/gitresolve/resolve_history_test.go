package gitresolve

import (
	"testing"
	"time"

	"github.com/jhanvi857/gitresolve/internal/conflict"
	"github.com/jhanvi857/gitresolve/internal/history"
	"github.com/jhanvi857/gitresolve/internal/ownership"
)

func TestEvaluateHistoryEscalation_BlastRadius(t *testing.T) {
	idx := history.NewIndex(500)
	for i := 0; i < 15; i++ {
		idx.SymbolCallers[history.SymbolID{Name: "ExecuteTrade"}] = append(
			idx.SymbolCallers[history.SymbolID{Name: "ExecuteTrade"}],
			history.CallerRef{File: "caller.go", Line: i + 1, Name: "Main"},
		)
	}

	c := &conflict.ConflictBlock{
		FilePath:       "trade.go",
		OursLines:      []string{"func ExecuteTrade(t Trade) error {"},
		TheirsLines:    []string{"func ExecuteTrade(t Trade, ctx Context) error {"},
		CanAutoResolve: true,
	}

	cfg := &ownership.PolicyConfig{MaxCallers: 10, CoChangeMinStrength: 0.6}
	evaluateHistoryEscalation(".", "trade.go", c, idx, cfg, map[string]bool{"trade.go": true})

	if c.ManualReasonCode != conflict.ReasonSemanticHighBlastRadius {
		t.Errorf("expected %s, got %s", conflict.ReasonSemanticHighBlastRadius, c.ManualReasonCode)
	}
	if c.CanAutoResolve {
		t.Error("expected CanAutoResolve to be false")
	}
	if c.SuggestHint != "go test ./... (run full suite before committing)" {
		t.Errorf("unexpected suggestion: %s", c.SuggestHint)
	}
}

func TestEvaluateHistoryEscalation_MissingCoupledFile(t *testing.T) {
	idx := history.NewIndex(500)
	idx.FileCoChanges["api/routes.go"] = map[string]int{"api/handlers.go": 90}
	idx.FileCoChanges["root.go"] = map[string]int{"base.go": 100}

	c := &conflict.ConflictBlock{
		FilePath:       "api/routes.go",
		OursLines:      []string{"r.GET(\"/users\", GetUsers)"},
		TheirsLines:    []string{"r.GET(\"/v2/users\", GetUsersV2)"},
		CanAutoResolve: true,
	}

	cfg := &ownership.PolicyConfig{MaxCallers: 10, CoChangeMinStrength: 0.6}
	// routes.go touched, handlers.go NOT touched
	evaluateHistoryEscalation(".", "api/routes.go", c, idx, cfg, map[string]bool{"api/routes.go": true})

	if c.ManualReasonCode != conflict.ReasonSemanticMissingCoupledFile {
		t.Errorf("expected %s, got %s", conflict.ReasonSemanticMissingCoupledFile, c.ManualReasonCode)
	}
	if c.SuggestHint != "gitresolve status api/handlers.go" {
		t.Errorf("unexpected suggestion: %s", c.SuggestHint)
	}
}

func TestEvaluateHistoryEscalation_ImportCycle(t *testing.T) {
	idx := history.NewIndex(500)
	idx.ImportEdges["pkg/x"] = []string{"pkg/y"}
	idx.ImportEdges["pkg/y"] = []string{"pkg/x"}

	c := &conflict.ConflictBlock{
		FilePath:       "pkg/x/mod.go",
		OursLines:      []string{"var X = 1"},
		TheirsLines:    []string{"var X = 2"},
		CanAutoResolve: true,
	}

	cfg := &ownership.PolicyConfig{MaxCallers: 10, CoChangeMinStrength: 0.6}
	evaluateHistoryEscalation(".", "pkg/x/mod.go", c, idx, cfg, map[string]bool{"pkg/x/mod.go": true})

	if c.ManualReasonCode != conflict.ReasonSemanticImportCycle {
		t.Errorf("expected %s, got %s", conflict.ReasonSemanticImportCycle, c.ManualReasonCode)
	}
	if c.SuggestHint != "review import dependencies and break the cycle" {
		t.Errorf("unexpected suggestion: %s", c.SuggestHint)
	}
}

func TestEvaluateHistoryEscalation_MultiAuthor(t *testing.T) {
	idx := history.NewIndex(500)
	idx.FileAuthors["shared.go"] = []history.AuthorContribution{
		{Email: "dev1@corp.com", Weight: 2.0, LastTouched: time.Now()},
		{Email: "dev2@corp.com", Weight: 1.5, LastTouched: time.Now().Add(-12 * time.Hour)},
	}

	c := &conflict.ConflictBlock{
		FilePath:       "shared.go",
		OursLines:      []string{"const Version = \"1.0\""},
		TheirsLines:    []string{"const Version = \"2.0\""},
		CanAutoResolve: false,
	}

	cfg := &ownership.PolicyConfig{MaxCallers: 10, CoChangeMinStrength: 0.6}
	evaluateHistoryEscalation(".", "shared.go", c, idx, cfg, map[string]bool{"shared.go": true})

	if c.ManualReasonCode != conflict.ReasonStrategyMultiAuthorConflict {
		t.Errorf("expected %s, got %s", conflict.ReasonStrategyMultiAuthorConflict, c.ManualReasonCode)
	}
	if c.SuggestHint != "coordinate with dev2@corp.com before resolving" {
		t.Errorf("unexpected suggestion: %s", c.SuggestHint)
	}
}

func TestSkipSyncCheckFlag(t *testing.T) {
	flag := resolveCmd.Flags().Lookup("skip-sync-check")
	if flag == nil {
		t.Fatal("expected --skip-sync-check flag to be registered on resolveCmd")
	}
	if flag.DefValue != "false" {
		t.Errorf("expected default value 'false', got %q", flag.DefValue)
	}
}

func TestPrintVerboseEvidence_NoError(t *testing.T) {
	idx := history.NewIndex(500)
	idx.FileAuthors["main.go"] = []history.AuthorContribution{
		{Email: "author@test.com", Weight: 1.0, LastTouched: time.Now()},
	}
	idx.FileCoChanges["main.go"] = map[string]int{"util.go": 10}

	c := &conflict.ConflictBlock{
		FilePath:    "main.go",
		OursLines:   []string{"fmt.Println(\"ours\")"},
		TheirsLines: []string{"fmt.Println(\"theirs\")"},
		Type:        conflict.TypeLogic,
		Severity:    conflict.SeverityMedium,
		Confidence:  0.5,
	}

	// Should execute and print without panicking
	printVerboseEvidence("main.go", c, idx)
}
