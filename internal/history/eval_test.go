package history_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/jhanvi857/gitresolve/internal/conflict"
	"github.com/jhanvi857/gitresolve/internal/history"
)

func TestEvaluateBlock_BlastRadius(t *testing.T) {
	idx := history.NewIndex(500)
	// Add 15 callers for symbol "CalculateTax"
	for i := 0; i < 15; i++ {
		idx.SymbolCallers[history.SymbolID{Name: "CalculateTax"}] = append(
			idx.SymbolCallers[history.SymbolID{Name: "CalculateTax"}],
			history.CallerRef{File: "service/invoice.go", Line: i + 1, Name: "GenerateInvoice"},
		)
	}

	c := &conflict.ConflictBlock{
		FilePath: "finance/tax.go",
		OursLines: []string{
			"func CalculateTax(amount float64) float64 {",
			"    return amount * 0.15",
			"}",
		},
		TheirsLines: []string{
			"func CalculateTax(amount float64) float64 {",
			"    return amount * 0.20",
			"}",
		},
		CanAutoResolve: true,
	}

	idx.EvaluateBlock(".", "finance/tax.go", c, 10, 0.6, map[string]bool{"finance/tax.go": true})

	if c.CanAutoResolve {
		t.Error("expected CanAutoResolve to be false for high blast radius")
	}
	if c.ManualReasonCode != conflict.ReasonSemanticHighBlastRadius {
		t.Errorf("got reason code %q, want %q", c.ManualReasonCode, conflict.ReasonSemanticHighBlastRadius)
	}
	expectedMsg := "CalculateTax is called from 15 other locations — escalating for manual review"
	if c.ManualReason != expectedMsg {
		t.Errorf("got manual reason %q, want %q", c.ManualReason, expectedMsg)
	}
	expectedCmd := "go test ./... (run full suite before committing)"
	if c.SuggestHint != expectedCmd {
		t.Errorf("got suggest hint %q, want %q", c.SuggestHint, expectedCmd)
	}
}

func TestEvaluateBlock_BlastRadius_BelowThreshold(t *testing.T) {
	idx := history.NewIndex(500)
	// Add only 3 callers (below threshold 10)
	for i := 0; i < 3; i++ {
		idx.SymbolCallers[history.SymbolID{Name: "CalculateTax"}] = append(
			idx.SymbolCallers[history.SymbolID{Name: "CalculateTax"}],
			history.CallerRef{File: "service/invoice.go", Line: i + 1, Name: "GenerateInvoice"},
		)
	}

	c := &conflict.ConflictBlock{
		FilePath: "finance/tax.go",
		OursLines: []string{
			"func CalculateTax(amount float64) float64 {",
			"    return amount * 0.15",
			"}",
		},
		TheirsLines: []string{
			"func CalculateTax(amount float64) float64 {",
			"    return amount * 0.20",
			"}",
		},
		CanAutoResolve: true,
	}

	idx.EvaluateBlock(".", "finance/tax.go", c, 10, 0.6, map[string]bool{"finance/tax.go": true})

	// Should not escalate on blast radius
	if c.ManualReasonCode == conflict.ReasonSemanticHighBlastRadius {
		t.Error("should not escalate when callers <= maxCallers")
	}
}

func TestEvaluateBlock_MissingCoupledFile(t *testing.T) {
	idx := history.NewIndex(500)
	// fileA and fileB have co-change strength 0.90
	idx.FileCoChanges["schema/db.sql"] = map[string]int{"models/schema.go": 90}
	idx.FileCoChanges["x.go"] = map[string]int{"y.go": 100} // max count = 100

	c := &conflict.ConflictBlock{
		FilePath: "schema/db.sql",
		OursLines: []string{
			"ALTER TABLE users ADD COLUMN age INT;",
		},
		TheirsLines: []string{
			"ALTER TABLE users ADD COLUMN birthday DATE;",
		},
		CanAutoResolve: true,
	}

	// schema/db.sql touched, but models/schema.go NOT touched
	touchedFiles := map[string]bool{"schema/db.sql": true}

	idx.EvaluateBlock(".", "schema/db.sql", c, 10, 0.6, touchedFiles)

	if c.CanAutoResolve {
		t.Error("expected CanAutoResolve to be false for missing coupled file")
	}
	if c.ManualReasonCode != conflict.ReasonSemanticMissingCoupledFile {
		t.Errorf("got reason code %q, want %q", c.ManualReasonCode, conflict.ReasonSemanticMissingCoupledFile)
	}
	expectedMsg := "schema/db.sql is frequently changed alongside models/schema.go (strength 0.90) but models/schema.go was not touched in this branch"
	if c.ManualReason != expectedMsg {
		t.Errorf("got manual reason %q, want %q", c.ManualReason, expectedMsg)
	}
	expectedCmd := "gitresolve status models/schema.go"
	if c.SuggestHint != expectedCmd {
		t.Errorf("got suggest hint %q, want %q", c.SuggestHint, expectedCmd)
	}
}

func TestEvaluateBlock_CoupledFile_AlreadyTouched(t *testing.T) {
	idx := history.NewIndex(500)
	idx.FileCoChanges["schema/db.sql"] = map[string]int{"models/schema.go": 90}
	idx.FileCoChanges["x.go"] = map[string]int{"y.go": 100}

	c := &conflict.ConflictBlock{
		FilePath: "schema/db.sql",
		OursLines: []string{
			"ALTER TABLE users ADD COLUMN age INT;",
		},
		TheirsLines: []string{
			"ALTER TABLE users ADD COLUMN birthday DATE;",
		},
		CanAutoResolve: true,
	}

	// Both files ARE touched in this branch
	touchedFiles := map[string]bool{"schema/db.sql": true, "models/schema.go": true}

	idx.EvaluateBlock(".", "schema/db.sql", c, 10, 0.6, touchedFiles)

	if c.ManualReasonCode == conflict.ReasonSemanticMissingCoupledFile {
		t.Error("should not escalate when coupled file is already touched")
	}
}

func TestEvaluateBlock_ImportCycle(t *testing.T) {
	idx := history.NewIndex(500)
	idx.ImportEdges["pkg/a"] = []string{"pkg/b"}
	idx.ImportEdges["pkg/b"] = []string{"pkg/c"}
	idx.ImportEdges["pkg/c"] = []string{"pkg/a"}

	c := &conflict.ConflictBlock{
		FilePath: "pkg/b/service.go",
		OursLines: []string{
			"func Run() { c.Do() }",
		},
		TheirsLines: []string{
			"func Run() { c.Execute() }",
		},
		CanAutoResolve: true,
	}

	idx.EvaluateBlock(".", "pkg/b/service.go", c, 10, 0.6, map[string]bool{"pkg/b/service.go": true})

	if c.CanAutoResolve {
		t.Error("expected CanAutoResolve to be false for import cycle")
	}
	if c.ManualReasonCode != conflict.ReasonSemanticImportCycle {
		t.Errorf("got reason code %q, want %q", c.ManualReasonCode, conflict.ReasonSemanticImportCycle)
	}
	expectedMsg := "conflict block in pkg/b/service.go sits inside an import cycle: pkg/a → pkg/b → pkg/c → pkg/a"
	if c.ManualReason != expectedMsg {
		t.Errorf("got manual reason %q, want %q", c.ManualReason, expectedMsg)
	}
	expectedCmd := "review import dependencies and break the cycle"
	if c.SuggestHint != expectedCmd {
		t.Errorf("got suggest hint %q, want %q", c.SuggestHint, expectedCmd)
	}
}

func TestEvaluateBlock_MultiAuthor(t *testing.T) {
	idx := history.NewIndex(500)
	idx.FileAuthors["cmd/server.go"] = []history.AuthorContribution{
		{Email: "lead@corp.com", Weight: 2.5, LastTouched: time.Now()},
		{Email: "contractor@corp.com", Weight: 1.2, LastTouched: time.Now().Add(-48 * time.Hour)},
	}

	c := &conflict.ConflictBlock{
		FilePath: "cmd/server.go",
		OursLines: []string{
			"func Start() { listen(8080) }",
		},
		TheirsLines: []string{
			"func Start() { listen(9090) }",
		},
		CanAutoResolve: false,
	}

	idx.EvaluateBlock(".", "cmd/server.go", c, 10, 0.6, map[string]bool{"cmd/server.go": true})

	if c.ManualReasonCode != conflict.ReasonStrategyMultiAuthorConflict {
		t.Errorf("got reason code %q, want %q", c.ManualReasonCode, conflict.ReasonStrategyMultiAuthorConflict)
	}
	expectedMsg := "both sides of the conflict in cmd/server.go were authored by different people: lead@corp.com vs contractor@corp.com"
	if c.ManualReason != expectedMsg {
		t.Errorf("got manual reason %q, want %q", c.ManualReason, expectedMsg)
	}
	expectedCmd := "coordinate with contractor@corp.com before resolving"
	if c.SuggestHint != expectedCmd {
		t.Errorf("got suggest hint %q, want %q", c.SuggestHint, expectedCmd)
	}
}

func TestEvaluateBlock_PriorityOrder(t *testing.T) {
	// If both blast radius and missing coupled file apply, blast radius (higher risk) takes priority
	idx := history.NewIndex(500)
	for i := 0; i < 20; i++ {
		idx.SymbolCallers[history.SymbolID{Name: "CoreAPI"}] = append(
			idx.SymbolCallers[history.SymbolID{Name: "CoreAPI"}],
			history.CallerRef{File: "caller.go", Line: i + 1, Name: "Call"},
		)
	}
	idx.FileCoChanges["core.go"] = map[string]int{"other.go": 95}
	idx.FileCoChanges["a.go"] = map[string]int{"b.go": 100}

	c := &conflict.ConflictBlock{
		FilePath: "core.go",
		OursLines: []string{
			"func CoreAPI() { doV1() }",
		},
		TheirsLines: []string{
			"func CoreAPI() { doV2() }",
		},
		CanAutoResolve: true,
	}

	idx.EvaluateBlock(".", "core.go", c, 10, 0.6, map[string]bool{"core.go": true})

	if c.ManualReasonCode != conflict.ReasonSemanticHighBlastRadius {
		t.Errorf("expected blast radius priority, got %q", c.ManualReasonCode)
	}
}

func TestDivergenceRendering_MergeStrategy(t *testing.T) {
	data := map[string]string{
		"behind":  strconv.Itoa(18),
		"branch":  "develop",
		"authors": "carol@example.com",
	}
	msg := conflict.RenderEscalationMessage(conflict.ReasonStrategyStaleBranchDiv, data)
	expectedMsg := "branch is 18 commits behind develop — carol@example.com authored changes touching files you also modified"
	if msg != expectedMsg {
		t.Errorf("got msg %q, want %q", msg, expectedMsg)
	}
}
