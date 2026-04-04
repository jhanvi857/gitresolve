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
	if !strings.Contains(c.ManualReason, "parse-check") {
		t.Fatalf("expected parse-check reason, got: %s", c.ManualReason)
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
