package conflict

import (
	"strings"
	"testing"
)

func TestRegression_test_m2_StrictMarkerFailure(t *testing.T) {
	content := "line 1\n<<<<<<< ours\nA\n=======\nB\n>>>>>>> theirs\nline 2\n"
	if err := Verify("main.go", content); err == nil {
		t.Fatal("expected Verify to fail when conflict markers remain")
	}
}

func TestRegression_test_h1_PartialStructuredFallback(t *testing.T) {
	c := &Conflict{
		FilePath:       "config.json",
		Type:           TypeStructured,
		CanAutoResolve: true,
		BaseLines: []string{
			"{\"service\":\"api\",\"replicas\":1}",
		},
		OurLines: []string{
			"{\"service\":\"api\",\"replicas\":2}",
		},
		TheirLines: []string{
			"{\"service\":\"api\",\"replicas\":3}",
		},
	}

	if AutoResolve(c, Options{}) {
		t.Fatal("expected structured overlap to fallback to manual")
	}
	if c.ManualReason == "" {
		t.Fatal("expected manual reason for structured fallback")
	}
	if c.SuggestHint == "" {
		t.Fatal("expected suggested strategy hint for structured fallback")
	}
}

func TestRegression_test_h5_UnresolvedMarkersBlockedInCompiledOutput(t *testing.T) {
	content := []byte(strings.Join([]string{
		"start",
		"<<<<<<< ours",
		"A",
		"=======",
		"B",
		">>>>>>> theirs",
		"end",
	}, "\n"))
	conflicts := ParseFile("config.yaml", content)
	if len(conflicts) != 1 {
		t.Fatalf("expected one conflict, got %d", len(conflicts))
	}

	// Leave unresolved on purpose to simulate partial resolution output.
	compiled := CompileResolution(content, conflicts)
	if err := Verify("config.yaml", compiled); err == nil {
		t.Fatal("expected Verify to fail when compiled output still has markers")
	}
}

func TestRegression_test_e3_ImportParseGate(t *testing.T) {
	c := &Conflict{
		FilePath:       "main.go",
		Type:           TypeImport,
		CanAutoResolve: true,
		OurLines: []string{
			"import (",
			"\t\"fmt\"",
			")",
		},
		TheirLines: []string{
			"import )",
			"\t\"os\"",
			")",
		},
	}

	if AutoResolve(c, Options{}) {
		t.Fatal("expected import auto-merge to fallback when parse-check fails")
	}
	if !strings.Contains(c.ManualReason, "parse-check") && !strings.Contains(c.ManualReason, "valid merged import block") {
		t.Fatalf("expected parse safety reason, got: %s", c.ManualReason)
	}
}

func TestRegression_test_e3_GoImportMergeSuccess(t *testing.T) {
	c := &Conflict{
		FilePath:       "main.go",
		Type:           TypeImport,
		CanAutoResolve: true,
		OurLines: []string{
			"import (",
			"\t\"fmt\"",
			"\t\"os\"",
			")",
		},
		TheirLines: []string{
			"import (",
			"\t\"fmt\"",
			"\t\"net/http\"",
			")",
		},
	}

	if !AutoResolve(c, Options{}) {
		t.Fatal("expected go import blocks to merge")
	}
	if !strings.Contains(c.Resolution, "\"os\"") || !strings.Contains(c.Resolution, "\"net/http\"") {
		t.Fatalf("expected merged imports, got: %s", c.Resolution)
	}
	if err := Verify("main.go", "package main\n"+c.Resolution+"\nfunc main(){}\n"); err != nil {
		t.Fatalf("expected merged import block to remain parseable: %v", err)
	}
}

func TestRegression_test_e3_GoSingleLineImportNormalizedToBlock(t *testing.T) {
	c := &Conflict{
		FilePath:       "main.go",
		Type:           TypeImport,
		CanAutoResolve: true,
		OurLines:       []string{"import \"fmt\""},
		TheirLines:     []string{"import \"net/http\""},
	}

	if !AutoResolve(c, Options{}) {
		t.Fatal("expected single-line go imports to merge")
	}
	if !strings.HasPrefix(c.Resolution, "import (") {
		t.Fatalf("expected normalized import block, got: %s", c.Resolution)
	}
}

func TestRegression_test_h2_m3_m4_ConfidenceThresholdGuidance(t *testing.T) {
	medium := &Conflict{
		FilePath:   "feature.ts",
		OurLines:   []string{"const x = 1"},
		TheirLines: []string{"const x = 2"},
	}
	Classify(medium)
	if ShouldAutoApply(medium) {
		t.Fatal("expected medium-confidence scalar conflict to skip auto-apply")
	}
	if !NeedsGuidedChoice(medium) {
		t.Fatal("expected guided choice for medium-confidence conflict")
	}

	high := &Conflict{
		FilePath:   "main.go",
		OurLines:   []string{"\tfmt.Println(\"x\")"},
		TheirLines: []string{"    fmt.Println(\"x\")"},
	}
	Classify(high)
	if !ShouldAutoApply(high) {
		t.Fatal("expected high-confidence whitespace conflict to auto-apply")
	}
}

func TestRegression_test_m2_CriticalGoModOverlapRejected(t *testing.T) {
	c := &Conflict{
		FilePath:       "go.mod",
		Type:           TypeImport,
		CanAutoResolve: true,
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
			"    github.com/gin-gonic/gin v1.9.0",
			")",
		},
	}

	if AutoResolve(c, Options{}) {
		t.Fatal("expected go.mod overlapping dependency change to require manual resolution")
	}
	if !strings.Contains(c.ManualReason, "go.mod entry") {
		t.Fatalf("expected critical overlap reason, got: %s", c.ManualReason)
	}
}

func TestRegression_SyntaxAssertions_JSON_YAML_TOML_Go(t *testing.T) {
	cases := []struct {
		name    string
		file    string
		content string
	}{
		{name: "json", file: "a.json", content: "{\"name\":\"ok\"}"},
		{name: "yaml", file: "a.yaml", content: "name: ok\ncount: 1\n"},
		{name: "toml", file: "a.toml", content: "name = \"ok\"\ncount = 1\n"},
		{name: "go", file: "a.go", content: "package main\nfunc main() {}\n"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := Verify(tc.file, tc.content); err != nil {
				t.Fatalf("expected valid syntax for %s: %v", tc.file, err)
			}
		})
	}
}

func TestRegression_test_m2_YAMLSequenceUnion(t *testing.T) {
	c := &Conflict{
		FilePath:       "items.yaml",
		Type:           TypeStructured,
		CanAutoResolve: true,
		BaseLines:      []string{"- apple", "- banana"},
		OurLines:       []string{"- apple", "- banana", "- cherry"},
		TheirLines:     []string{"- apple", "- banana", "- date"},
	}

	if !AutoResolve(c, Options{}) {
		t.Fatal("expected yaml sequence overlap to auto-merge")
	}
	if !strings.Contains(c.Resolution, "- cherry") || !strings.Contains(c.Resolution, "- date") {
		t.Fatalf("expected union list merge, got: %s", c.Resolution)
	}
	if err := Verify("items.yaml", c.Resolution); err != nil {
		t.Fatalf("expected merged yaml to be valid: %v", err)
	}
}

func TestRegression_test_m4_TOMLSnippetMerge(t *testing.T) {
	c := &Conflict{
		FilePath:       "app.toml",
		Type:           TypeStructured,
		CanAutoResolve: true,
		BaseLines: []string{
			"enabled = false",
			"timeout = 30",
		},
		OurLines: []string{
			"enabled = false",
			"timeout = 30",
		},
		TheirLines: []string{
			"enabled = true",
			"retries = 3",
		},
	}

	if !AutoResolve(c, Options{}) {
		t.Fatal("expected toml key-value snippet merge to succeed")
	}
	if !strings.Contains(c.Resolution, "enabled") || !strings.Contains(c.Resolution, "retries") {
		t.Fatalf("expected merged toml keys, got: %s", c.Resolution)
	}
	if err := Verify("app.toml", c.Resolution); err != nil {
		t.Fatalf("expected merged toml to be valid: %v", err)
	}
}
