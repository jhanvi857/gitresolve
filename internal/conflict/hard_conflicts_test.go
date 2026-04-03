package conflict

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestStructuredConflict_PackageJSON(t *testing.T) {
	c := &Conflict{
		FilePath: "package.json",
		BaseLines: []string{
			"{",
			"  \"name\": \"myapp\",",
			"  \"version\": \"1.0.0\",",
			"  \"dependencies\": {",
			"    \"lodash\": \"4.17.20\"",
			"  }",
			"}",
		},
		OurLines: []string{
			"{",
			"  \"name\": \"myapp\",",
			"  \"version\": \"1.0.1\",",
			"  \"dependencies\": {",
			"    \"lodash\": \"4.17.21\"",
			"  }",
			"}",
		},
		TheirLines: []string{
			"{",
			"  \"name\": \"myapp\",",
			"  \"version\": \"1.0.0\",",
			"  \"dependencies\": {",
			"    \"lodash\": \"4.17.20\",",
			"    \"axios\": \"0.21.1\"",
			"  }",
			"}",
		},
	}

	Classify(c)

	if c.Type != TypeStructured {
		t.Errorf("expected Structured, got %v", c.Type)
	}
	if c.Severity != SeverityHigh {
		t.Errorf("expected High severity for package.json, got %v", c.Severity)
	}
	if !c.CanAutoResolve {
		t.Error("package.json should be allowed to attempt auto-resolution")
	}

	// Try auto-resolve
	resolved := AutoResolve(c, Options{})
	if !resolved {
		t.Error("AutoResolve should succeed for non-overlapping package.json changes")
	}
	if !strings.Contains(c.Resolution, "\"version\": \"1.0.1\"") || !strings.Contains(c.Resolution, "\"axios\": \"0.21.1\"") {
		t.Errorf("AutoResolve failed to merge package.json correctly: %s", c.Resolution)
	}

	// Test the case where we force it but it has conflicting array edits (YAML example)
}

func TestStructuredConflict_ArrayAmbiguity(t *testing.T) {
	c := &Conflict{
		FilePath: "config.yaml",
		BaseLines: []string{
			"servers:",
			"  - host: db1",
		},
		OurLines: []string{
			"servers:",
			"  - host: db1",
			"  - host: db2",
		},
		TheirLines: []string{
			"servers:",
			"  - host: db1",
			"  - host: db3",
		},
	}

	Classify(c)
	// YAML arrays are now merged seamlessly 
	resolved := AutoResolve(c, Options{})
	if !resolved {
		t.Error("should auto-resolve additive array merge")
	}
	if !strings.Contains(c.Resolution, "host: db2") || !strings.Contains(c.Resolution, "host: db3") {
		t.Errorf("array additive merge failed: %s", c.Resolution)
	}
}

func TestSignatureChange_Go(t *testing.T) {
	c := &Conflict{
		FilePath: "main.go",
		OurLines: []string{
			"func Process(ctx context.Context, data string) error {",
		},
		TheirLines: []string{
			"func Process(data string, timeout int) {",
		},
	}

	Classify(c)

	if c.Type != TypeSignature {
		t.Errorf("expected Signature conflict, got %v", c.Type)
	}
	if c.CanAutoResolve {
		t.Error("signature changes must not be auto-resolved")
	}
}

func TestNoAutoStructuredFlag(t *testing.T) {
	c := &Conflict{
		FilePath: "data.json",
		BaseLines: []string{"{\"key\": \"base\"}"},
		OurLines: []string{"{\"key\": \"ours\"}"},
		TheirLines: []string{"{\"key\": \"base\"}"},
		CanAutoResolve: true, // Manually set for test
		Type: TypeStructured,
	}

	// Should resolve if flag is false
	if !AutoResolve(c, Options{NoAutoStructured: false}) {
		t.Error("should auto-resolve simple JSON if flag is false")
	}

	// Should NOT resolve if flag is true
	c.Resolution = ""
	if AutoResolve(c, Options{NoAutoStructured: true}) {
		t.Error("should NOT auto-resolve if NoAutoStructured is true")
	}
}

func TestValuesEqual(t *testing.T) {
	// This tests the underlying structured logic via internal knowledge
	// but we'll just test if json parsing etc works for comparison
	var a interface{} = map[string]interface{}{"a": 1, "b": []int{1, 2}}
	var b interface{} = map[string]interface{}{"b": []int{1, 2}, "a": 1}
	
	// They are semantically equal but might have different keys order in JSON if not careful
	// But json.Marshal for maps is deterministic? Yes, it sorts by key.
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	if string(aj) != string(bj) {
		t.Errorf("JSON marshal not deterministic for maps? %s != %s", aj, bj)
	}
}
func TestTSXConflict(t *testing.T) {
	c := &Conflict{
		FilePath: "Component.tsx",
		OurLines: []string{"const App = () => <div className='foo'>{count}</div>;"},
		TheirLines: []string{"const App = () => <div className='bar'>{total}</div>;"},
	}
	Classify(c)
	// JSX/TSX changes should be caught by logic/signature rules
	if c.CanAutoResolve {
		t.Error("TSX logic changes should not be auto-resolved")
	}
}

func TestNestedYAMLMerge(t *testing.T) {
	c := &Conflict{
		FilePath: "deploy.yaml",
		BaseLines: []string{
			"services:",
			"  api:",
			"    image: v1",
			"    env:",
			"      DEBUG: \"false\"",
		},
		OurLines: []string{
			"services:",
			"  api:",
			"    image: v2",
			"    env:",
			"      DEBUG: \"false\"",
		},
		TheirLines: []string{
			"services:",
			"  api:",
			"    image: v1",
			"    env:",
			"      DEBUG: \"true\"",
			"      LOG_LEVEL: \"info\"",
		},
	}
	Classify(c)
	resolved := AutoResolve(c, Options{})
	if !resolved {
		t.Error("Nested YAML merge should succeed when changes are in different keys")
	}
	if !strings.Contains(c.Resolution, "image: v2") || !strings.Contains(c.Resolution, "DEBUG: \"true\"") {
		t.Errorf("Resolution failed to merge nested keys: %s", c.Resolution)
	}
}

func TestGoModConflict(t *testing.T) {
	c := &Conflict{
		FilePath: "go.mod",
		BaseLines: []string{
			"require (",
			"    github.com/gin-gonic/gin v1.7.0",
			")",
		},
		OurLines: []string{
			"require (",
			"    github.com/gin-gonic/gin v1.8.0",
			")",
		},
		TheirLines: []string{
			"require (",
			"    github.com/gin-gonic/gin v1.7.0",
			"    github.com/spf13/cobra v1.4.0",
			")",
		},
	}
	Classify(c)
	if c.Severity != SeverityHigh || !c.CanAutoResolve {
		t.Error("go.mod should be High severity and allowed to attempt auto-resolve")
	}
	
	resolved := AutoResolve(c, Options{})
	if !resolved {
		t.Error("go.mod non-overlapping changes should be auto-merged by import deduplication logic")
	}
}

func TestAliasedImports(t *testing.T) {
	c := &Conflict{
		FilePath: "main.go",
		OurLines: []string{
			"import g \"github.com/go-git/go-git/v5\"",
			"import \"fmt\"",
		},
		TheirLines: []string{
			"import g \"github.com/go-git/go-git/v5\"",
			"import \"os\"",
		},
	}
	Classify(c)
	resolved := AutoResolve(c, Options{})
	if !resolved {
		t.Error("Aliased imports should be merged correctly")
	}
	if !strings.Contains(c.Resolution, "import g \"github.com/go-git/go-git/v5\"") {
		t.Error("Alias was lost in merge")
	}
}

func TestDeleteModifySensitive(t *testing.T) {
	c := &Conflict{
		FilePath: "internal/auth/provider.go",
		OurLines: []string{}, // deleted
		TheirLines: []string{
			"func VerifyToken(token string) bool {",
			"    return true",
			"}",
		},
	}
	Classify(c)
	if c.Type != TypeDeleteModify || c.Severity != SeverityCritical {
		t.Errorf("Expected Critical DeleteModify for auth file, got %v with severity %v", c.Type, c.Severity)
	}
}

func TestTOMLNestedMerge(t *testing.T) {
	c := &Conflict{
		FilePath: "Cargo.toml",
		BaseLines: []string{
			"[package]",
			"name = \"foo\"",
			"version = \"0.1.0\"",
		},
		OurLines: []string{
			"[package]",
			"name = \"foo\"",
			"version = \"0.1.1\"",
		},
		TheirLines: []string{
			"[package]",
			"name = \"foo\"",
			"version = \"0.1.0\"",
			"authors = [\"me\"]",
		},
	}
	Classify(c)
	// Cargo.toml is critical. 
	if c.Severity != SeverityHigh {
		t.Error("Cargo.toml should be high severity")
	}
}

func TestIndentationWhitespace(t *testing.T) {
	c := &Conflict{
		FilePath: "style.css",
		OurLines: []string{
			".btn {",
			"  color: red;",
			"}",
		},
		TheirLines: []string{
			".btn {",
			"\tcolor: red;",
			"}",
		},
	}
	Classify(c)
	if c.Type != TypeWhitespace {
		t.Errorf("Expected TypeWhitespace for indentation change, got %v", c.Type)
	}
}

func TestLogicConflict_Renames(t *testing.T) {
	c := &Conflict{
		FilePath: "util.js",
		OurLines: []string{"const calculateTotal = (price, tax) => price * tax;"},
		TheirLines: []string{"const getFullAmount = (val, rate) => val * rate;"},
	}
	Classify(c)
	if c.CanAutoResolve {
		t.Error("Significant logic/naming changes should not be auto-resolved")
	}
}

func TestConflictedStructuredMerge(t *testing.T) {
	c := &Conflict{
		FilePath: "settings.json",
		BaseLines: []string{"{\"theme\": \"light\"}"},
		OurLines: []string{"{\"theme\": \"dark\"}"},
		TheirLines: []string{"{\"theme\": \"high-contrast\"}"},
	}
	Classify(c)
	resolved := AutoResolve(c, Options{})
	if resolved {
		t.Error("Simultaneous scalar edits to same key should NOT auto-resolve")
	}
}

func TestGoInterfaceChange(t *testing.T) {
	c := &Conflict{
		FilePath: "service.go",
		OurLines: []string{
			"type Store interface {",
			"    Get(id string) (Item, error)",
			"}",
		},
		TheirLines: []string{
			"type Store interface {",
			"    Get(id string) (Item, error)",
			"    Save(item Item) error",
			"}",
		},
	}
	Classify(c)
	if c.Type == TypeUnknown {
		t.Error("Interface change should be classified")
	}
}
