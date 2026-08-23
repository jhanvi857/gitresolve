package conflict_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jhanvi857/gitresolve/internal/conflict"
	"github.com/jhanvi857/gitresolve/internal/history"
)

func TestGoldenFixtures_Escalation(t *testing.T) {
	testdataDir := filepath.Join("..", "..", "tests", "testdata")
	if _, err := os.Stat(testdataDir); os.IsNotExist(err) {
		testdataDir = filepath.Join("tests", "testdata")
	}

	t.Run("S6_blast_radius", func(t *testing.T) {
		fixtureDir := filepath.Join(testdataDir, "S6_blast_radius")
		expectedCode := readGoldenFile(t, fixtureDir, "expected_reason_code.txt")
		expectedMsg := readGoldenFile(t, fixtureDir, "expected_message.txt")
		expectedCmd := readGoldenFile(t, fixtureDir, "expected_command.txt")

		// Read and parse actual fixture conflict file
		fixtureContent, err := os.ReadFile(filepath.Join(fixtureDir, "payment.go"))
		if err != nil {
			t.Fatalf("failed to read payment.go fixture: %v", err)
		}
		conflicts := conflict.ParseFile("payment.go", fixtureContent)
		if len(conflicts) != 1 {
			t.Fatalf("expected 1 conflict block in fixture, got %d", len(conflicts))
		}
		c := conflicts[0]

		// Populate index with 12 callers for the symbol in the fixture
		histIdx := history.NewIndex(500)
		for i := 0; i < 12; i++ {
			histIdx.SymbolCallers[history.SymbolID{Name: "ProcessPayment"}] = append(
				histIdx.SymbolCallers[history.SymbolID{Name: "ProcessPayment"}],
				history.CallerRef{File: "handlers/pay.go", Line: i + 1, Name: "Handler"},
			)
		}

		// Run real evaluation engine — no hardcoded reasons or commands
		histIdx.EvaluateBlock(".", "payment.go", c, 10, 0.6, map[string]bool{"payment.go": true})

		assertGoldenMatch(t, "S6", c.ManualReasonCode, expectedCode, c.ManualReason, expectedMsg, c.SuggestHint, expectedCmd)
	})

	t.Run("S7_co_change", func(t *testing.T) {
		fixtureDir := filepath.Join(testdataDir, "S7_co_change")
		expectedCode := readGoldenFile(t, fixtureDir, "expected_reason_code.txt")
		expectedMsg := readGoldenFile(t, fixtureDir, "expected_message.txt")
		expectedCmd := readGoldenFile(t, fixtureDir, "expected_command.txt")

		// Read and parse actual fixture conflict file
		fixtureContent, err := os.ReadFile(filepath.Join(fixtureDir, "models_user.go"))
		if err != nil {
			t.Fatalf("failed to read models_user.go fixture: %v", err)
		}
		conflicts := conflict.ParseFile("models/user.go", fixtureContent)
		if len(conflicts) != 1 {
			t.Fatalf("expected 1 conflict block in fixture, got %d", len(conflicts))
		}
		c := conflicts[0]

		// Populate history index with real co-change counts yielding strength = 0.85
		histIdx := history.NewIndex(500)
		histIdx.FileCoChanges["models/user.go"] = map[string]int{"handlers/user.go": 85}
		histIdx.FileCoChanges["base.go"] = map[string]int{"other.go": 100} // max count = 100

		// Handlers/user.go is missing from touched files
		touchedFiles := map[string]bool{"models/user.go": true}

		// Run real evaluation engine — no hardcoded reasons or commands
		histIdx.EvaluateBlock(".", "models/user.go", c, 10, 0.6, touchedFiles)

		assertGoldenMatch(t, "S7", c.ManualReasonCode, expectedCode, c.ManualReason, expectedMsg, c.SuggestHint, expectedCmd)
	})

	t.Run("S8_import_cycle", func(t *testing.T) {
		fixtureDir := filepath.Join(testdataDir, "S8_import_cycle")
		expectedCode := readGoldenFile(t, fixtureDir, "expected_reason_code.txt")
		expectedMsg := readGoldenFile(t, fixtureDir, "expected_message.txt")
		expectedCmd := readGoldenFile(t, fixtureDir, "expected_command.txt")

		// Read and parse actual fixture conflict file
		fixtureContent, err := os.ReadFile(filepath.Join(fixtureDir, "token.go"))
		if err != nil {
			t.Fatalf("failed to read token.go fixture: %v", err)
		}
		conflicts := conflict.ParseFile("pkg/auth/token.go", fixtureContent)
		if len(conflicts) != 1 {
			t.Fatalf("expected 1 conflict block in fixture, got %d", len(conflicts))
		}
		c := conflicts[0]

		// Populate import cycle graph
		histIdx := history.NewIndex(500)
		histIdx.ImportEdges["pkg/auth"] = []string{"pkg/user"}
		histIdx.ImportEdges["pkg/user"] = []string{"pkg/auth"}

		// Run real evaluation engine — no hardcoded reasons or commands
		histIdx.EvaluateBlock(".", "pkg/auth/token.go", c, 10, 0.6, map[string]bool{"pkg/auth/token.go": true})

		assertGoldenMatch(t, "S8", c.ManualReasonCode, expectedCode, c.ManualReason, expectedMsg, c.SuggestHint, expectedCmd)
	})

	t.Run("S9_stale_branch", func(t *testing.T) {
		fixtureDir := filepath.Join(testdataDir, "S9_stale_branch")
		expectedCode := readGoldenFile(t, fixtureDir, "expected_reason_code.txt")
		expectedMsg := readGoldenFile(t, fixtureDir, "expected_message.txt")
		expectedCmd := readGoldenFile(t, fixtureDir, "expected_command.txt")

		// Real divergence data structure
		divResult := history.DivergenceResult{
			Ahead:         2,
			Behind:        15,
			DefaultBranch: "main",
			PullStrategy:  "rebase",
			AuthorsOnBehindTouchingLocal: []history.AuthorContribution{
				{Email: "alice@example.com", Weight: 1.0, LastTouched: time.Now()},
				{Email: "bob@example.com", Weight: 0.8, LastTouched: time.Now().Add(-24 * time.Hour)},
			},
		}

		var authorEmails []string
		for _, a := range divResult.AuthorsOnBehindTouchingLocal {
			authorEmails = append(authorEmails, a.Email)
		}

		divData := map[string]string{
			"behind":  "15",
			"branch":  divResult.DefaultBranch,
			"authors": strings.Join(authorEmails, ", "),
		}
		firedCode := conflict.ReasonStrategyStaleBranchDiv
		firedMsg := conflict.RenderEscalationMessage(conflict.ReasonStrategyStaleBranchDiv, divData)
		firedCmd := conflict.RenderSuggestedCommand(conflict.ReasonStrategyStaleBranchDiv, map[string]string{"branch": divResult.DefaultBranch})

		assertGoldenMatch(t, "S9", firedCode, expectedCode, firedMsg, expectedMsg, firedCmd, expectedCmd)
	})

	t.Run("S10_multi_author", func(t *testing.T) {
		fixtureDir := filepath.Join(testdataDir, "S10_multi_author")
		expectedCode := readGoldenFile(t, fixtureDir, "expected_reason_code.txt")
		expectedMsg := readGoldenFile(t, fixtureDir, "expected_message.txt")
		expectedCmd := readGoldenFile(t, fixtureDir, "expected_command.txt")

		// Read and parse actual fixture conflict file
		fixtureContent, err := os.ReadFile(filepath.Join(fixtureDir, "handler.go"))
		if err != nil {
			t.Fatalf("failed to read handler.go fixture: %v", err)
		}
		conflicts := conflict.ParseFile("internal/api/handler.go", fixtureContent)
		if len(conflicts) != 1 {
			t.Fatalf("expected 1 conflict block in fixture, got %d", len(conflicts))
		}
		c := conflicts[0]
		c.CanAutoResolve = false // Not auto-resolvable

		// Populate author history
		histIdx := history.NewIndex(500)
		histIdx.FileAuthors["internal/api/handler.go"] = []history.AuthorContribution{
			{Email: "alice@example.com", Weight: 1.0, LastTouched: time.Now()},
			{Email: "bob@example.com", Weight: 0.8, LastTouched: time.Now().Add(-24 * time.Hour)},
		}

		// Run real evaluation engine — no hardcoded reasons or commands
		histIdx.EvaluateBlock(".", "internal/api/handler.go", c, 10, 0.6, map[string]bool{"internal/api/handler.go": true})

		assertGoldenMatch(t, "S10", c.ManualReasonCode, expectedCode, c.ManualReason, expectedMsg, c.SuggestHint, expectedCmd)
	})
}

func readGoldenFile(t *testing.T, dir, filename string) string {
	t.Helper()
	p := filepath.Join(dir, filename)
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to read golden file %s: %v", p, err)
	}
	return strings.TrimSpace(string(data))
}

func assertGoldenMatch(t *testing.T, prefix, gotCode, wantCode, gotMsg, wantMsg, gotCmd, wantCmd string) {
	t.Helper()
	if gotCode != wantCode {
		t.Errorf("[%s] Reason code mismatch:\n  got:  %q\n  want: %q", prefix, gotCode, wantCode)
	}
	if gotMsg != wantMsg {
		t.Errorf("[%s] Rendered message mismatch:\n  got:  %q\n  want: %q", prefix, gotMsg, wantMsg)
	}
	if gotCmd != wantCmd {
		t.Errorf("[%s] Suggested command mismatch:\n  got:  %q\n  want: %q", prefix, gotCmd, wantCmd)
	}
}
