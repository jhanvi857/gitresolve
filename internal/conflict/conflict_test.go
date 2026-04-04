package conflict

import (
	"strings"
	"testing"
)

func TestClassifierAndAutoResolve_Whitespace(t *testing.T) {
	c := &Conflict{
		FilePath: "main.go",
		OurLines: []string{
			"func main() {",
			"    fmt.Println(\"hello\")",
			"}",
		},
		TheirLines: []string{
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
	c := &Conflict{
		FilePath: "api.go",
		OurLines: []string{
			"import \"fmt\"",
			"import \"net/http\"",
		},
		TheirLines: []string{
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

func TestClassifier_LogicConflict(t *testing.T) {
	c := &Conflict{
		FilePath: "auth/login.go",
		OurLines: []string{
			"if user.Valid() {",
			"    login(user)",
			"}",
		},
		TheirLines: []string{
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
