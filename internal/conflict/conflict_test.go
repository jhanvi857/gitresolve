package conflict

import (
	"strings"
	"testing"
)

func TestClassifierAndAutoResolve_Whitespace(t *testing.T) {
	c := &ConflictBlock{
		FilePath: "main.go",
		OursLines: []string{
			"func main() {",
			"    fmt.Println(\"hello\")",
			"}",
		},
		TheirsLines: []string{
			"func main() {",
			"\tfmt.Println(\"hello\")",
			"}",
		},
	}

	Classify(c)

	if c.Type != TypeWhitespace {
		t.Errorf("expected type Whitespace, got %v", c.Type)
	}
	if !c.CanAutoResolve {
		t.Error("expected CanAutoResolve to be true for whitespace")
	}

	resolved := AutoResolve(c, Options{})
	if !resolved {
		t.Error("AutoResolve failed to resolve whitespace conflict")
	}

	expectedResolution := "func main() {\n    fmt.Println(\"hello\")\n}"
	if c.Resolution != expectedResolution {
		t.Errorf("expected `%s`, got `%s`", expectedResolution, c.Resolution)
	}
}

func TestClassifierAndAutoResolve_Imports(t *testing.T) {
	c := &ConflictBlock{
		FilePath: "api.go",
		OursLines: []string{
			"import \"fmt\"",
			"import \"net/http\"",
		},
		TheirsLines: []string{
			"import \"fmt\"",
			"import \"os\"",
		},
	}

	Classify(c)

	if c.Type != TypeImport {
		t.Errorf("expected type Import, got %v", c.Type)
	}
	if !c.CanAutoResolve {
		t.Error("expected CanAutoResolve to be true for imports")
	}

	resolved := AutoResolve(c, Options{})
	if !resolved {
		t.Error("AutoResolve failed")
	}

	if !strings.HasPrefix(c.Resolution, "import (") ||
		!strings.Contains(c.Resolution, "\"fmt\"") ||
		!strings.Contains(c.Resolution, "\"net/http\"") ||
		!strings.Contains(c.Resolution, "\"os\"") {
		t.Errorf("Imports did not merge correctly. Got: %v", c.Resolution)
	}
}

func TestClassifierAndAutoResolve_IdenticalBothSides(t *testing.T) {
	// Simulate what git WOULD produce if it didn't auto-resolve.
	// When both sides have identical changes, this should be classified as TypeIdentical
	// and auto-resolved with the reason code "strategy.identical_both_sides".
	c := &ConflictBlock{
		FilePath: "package/main.go",
		OursLines: []string{
			"func hello() string {",
			"    return \"hello world\"",
			"}",
		},
		TheirsLines: []string{
			"func hello() string {",
			"    return \"hello world\"",
			"}",
		},
	}

	Classify(c)

	if c.Type != TypeIdentical {
		t.Errorf("expected type Identical, got %v", c.Type)
	}
	if c.Severity != SeverityTrivial {
		t.Errorf("expected severity Trivial, got %v", c.Severity)
	}
	if !c.CanAutoResolve {
		t.Error("expected CanAutoResolve to be true for identical changes")
	}
	if c.Confidence < 0.95 {
		t.Errorf("expected high confidence (>0.95), got %f", c.Confidence)
	}

	resolved := AutoResolve(c, Options{})
	if !resolved {
		t.Error("AutoResolve failed to resolve identical conflict")
	}

	expectedResolution := "func hello() string {\n    return \"hello world\"\n}"
	if c.Resolution != expectedResolution {
		t.Errorf("expected resolution `%s`, got `%s`", expectedResolution, c.Resolution)
	}
}

func TestClassifier_LogicConflict(t *testing.T) {
	c := &ConflictBlock{
		FilePath: "auth/login.go",
		OursLines: []string{
			"if user.Valid() {",
			"    login(user)",
			"}",
		},
		TheirsLines: []string{
			"if user.IsActive() {",
			"    generateToken(user)",
			"}",
		},
	}

	Classify(c)

	if c.Type != TypeLogic {
		t.Errorf("expected logic conflict, got %v", c.Type)
	}
	if c.Severity != SeverityCritical {
		t.Errorf("expected critical severity due to sensitive auth path, got %v", c.Severity)
	}
	if c.CanAutoResolve {
		t.Error("logic conflicts shouldn't be auto resolved")
	}
}

func TestParser_UnbalancedBraceNoSwallow(t *testing.T) {
	content := `<<<<<<< OURS
func test() {
=======
func test() { {
>>>>>>> THEIRS
some unrelated code
<<<<<<< OURS
var x = 1
=======
var x = 2
>>>>>>> THEIRS`

	blocks := ParseFile("test.go", []byte(content))
	if len(blocks) != 2 {
		t.Fatalf("expected 2 conflict blocks, got %d. Subsequent conflict block was swallowed!", len(blocks))
	}
	if len(blocks[0].TheirsLines) != 2 || blocks[0].TheirsLines[0] != "func test() { {" || blocks[0].TheirsLines[1] != "some unrelated code" {
		t.Errorf("expected 2 lines in first TheirsLines, got %v", blocks[0].TheirsLines)
	}
	if len(blocks[1].OursLines) != 1 || blocks[1].OursLines[0] != "var x = 1" {
		t.Errorf("expected var x = 1 in second OursLines, got %v", blocks[1].OursLines)
	}
}
