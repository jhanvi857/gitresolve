package conflict

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestSymmetricBraceRecovery(t *testing.T) {
	content := []byte(strings.Join([]string{
		"package auth",
		"",
		"import \"time\"",
		"",
		"<<<<<<< HEAD",
		"const TokenExpiry = 48 * time.Hour",
		"=======",
		"const TokenExpiry = 72 * time.Hour",
		">>>>>>> feature/performance-updates",
		"",
		"func ValidateToken(token string) bool {",
		"\treturn len(token) > 0",
		"}",
		"",
		"<<<<<<< HEAD",
		"func GenerateToken(userID string) string {",
		"\treturn userID + \"-token\"",
		"=======",
		"func RevokeToken(token string) error {",
		"\treturn nil",
		">>>>>>> feature/performance-updates",
		"}",
	}, "\n"))

	conflicts := ParseFile("auth.go", content)
	if len(conflicts) != 2 {
		t.Fatalf("expected 2 conflict blocks, got %d", len(conflicts))
	}

	if _, err := Resolve(conflicts[0], StrategyTheirs, ResolveOptions{}); err != nil {
		t.Fatalf("expected first conflict (theirs) to resolve, got %v", err)
	}
	if _, err := Resolve(conflicts[1], StrategyBoth, ResolveOptions{}); err != nil {
		t.Fatalf("expected second conflict (both) to resolve, got %v", err)
	}

	output := CompileResolution(content, conflicts)

	if !strings.Contains(output, "const TokenExpiry = 72 * time.Hour") {
		t.Fatal("expected output to contain theirs const TokenExpiry value")
	}
	if !strings.Contains(output, "func GenerateToken(userID string) string {") {
		t.Fatal("expected output to contain GenerateToken")
	}
	if !strings.Contains(output, "func RevokeToken(token string) error {") {
		t.Fatal("expected output to contain RevokeToken")
	}
	if strings.Contains(output, "<<<<<<<") || strings.Contains(output, "=======") || strings.Contains(output, ">>>>>>>") {
		t.Fatal("expected output to contain zero conflict markers")
	}

	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, "auth.go", output, parser.AllErrors); err != nil {
		t.Fatalf("expected reconstructed output to pass go/parser.ParseFile, got %v", err)
	}
}

func TestRegression_test_m2_StrictMarkerFailure(t *testing.T) {
	content := "line 1\n<<<<<<< ours\nA\n=======\nB\n>>>>>>> theirs\nline 2\n"
	if err := Verify("main.go", content); err == nil {
		t.Fatal("expected Verify to fail when conflict markers remain")
	}
}

func TestRegression_test_h1_PartialStructuredFallback(t *testing.T) {
	c := &ConflictBlock{
		FilePath:       "config.json",
		Type:           TypeStructured,
		CanAutoResolve: true,
		BaseLines: []string{
			"{\"service\":\"api\",\"replicas\":1}",
		},
		OursLines: []string{
			"{\"service\":\"api\",\"replicas\":2}",
		},
		TheirsLines: []string{
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

func TestBothSelectionClosingBrace(t *testing.T) {
	// MALFORMED INPUT: closing brace belongs to THEIRS but is outside marker.
	// Brace-aware parser should consume it.
	content := []byte(strings.Join([]string{
		"package main",
		"<<<<<<< HEAD",
		"func GenerateToken(userID string) string {",
		"	return userID + \"-token\"",
		"}",
		"=======",
		"func RevokeToken(token string) error {",
		"	return nil",
		">>>>>>> feature/performance-updates",
		"}",
	}, "\n"))

	conflicts := ParseFile("jwt.go", content)
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict block, got %d", len(conflicts))
	}

	c := conflicts[0]
	if len(c.TheirsLines) != 3 {
		t.Fatalf("expected TheirsLines to have 3 lines, got %d", len(c.TheirsLines))
	}
	if !strings.Contains(c.TheirsLines[2], "}") {
		t.Fatal("expected TheirsLines to include closing brace consumed after marker")
	}

	// Simulate user selecting [B]oth and expect success.
	_, err := Resolve(c, StrategyBoth, ResolveOptions{})
	if err != nil {
		t.Fatalf("expected Resolve(Both) to succeed, got %v", err)
	}

	output := CompileResolution(content, conflicts)
	if strings.Contains(output, "<<<<<<<") || strings.Contains(output, ">>>>>>>") {
		t.Fatal("output still contains conflict markers")
	}

	if err := Verify("jwt.go", output); err != nil {
		t.Fatalf("expected output file to pass verification, got: %v", err)
	}
}

func TestTheirsSelectionDanglingBrace(t *testing.T) {
	content := []byte(strings.Join([]string{
		"package main",
		"<<<<<<< HEAD",
		"func GenerateToken(userID string) string {",
		"	return userID + \"-token\"",
		"}",
		"======= ",
		"func RevokeToken(token string) error {",
		"	return nil",
		">>>>>>> feature/performance-updates",
		"}",
	}, "\n"))

	conflicts := ParseFile("jwt.go", content)
	c := conflicts[0]
	if len(c.TheirsLines) != 3 {
		t.Fatalf("expected TheirsLines to have 3 lines, got %d", len(c.TheirsLines))
	}

	// Simulate user selecting [T]heirs
	_, err := Resolve(c, StrategyTheirs, ResolveOptions{})
	if err != nil {
		t.Fatalf("expected Resolve(Theirs) to succeed, got %v", err)
	}

	output := CompileResolution(content, conflicts)
	output = strings.TrimSpace(output)

	if err := Verify("jwt.go", output); err != nil {
		t.Fatalf("expected output file to pass verification, got: %v", err)
	}
	if strings.Contains(output, "func GenerateToken") {
		t.Fatal("output should not contain GenerateToken")
	}
	if !strings.Contains(output, "func RevokeToken") {
		t.Fatal("output missing RevokeToken")
	}
}

func TestRegression_test_e3_ImportParseGate(t *testing.T) {
	c := &ConflictBlock{
		FilePath:       "main.go",
		Type:           TypeImport,
		CanAutoResolve: true,
		OursLines: []string{
			"import (",
			"\t\"fmt\"",
			")",
		},
		TheirsLines: []string{
			"import )",
			"\t\"os\"",
			")",
		},
	}

	if AutoResolve(c, Options{}) {
		t.Fatal("expected import auto-merge to fallback when parse-check fails")
	}
}

func TestRegression_test_h2_m3_m4_ConfidenceThresholdGuidance(t *testing.T) {
	medium := &ConflictBlock{
		FilePath:    "feature.ts",
		OursLines:   []string{"const x = 1"},
		TheirsLines: []string{"const x = 2"},
	}
	Classify(medium)
	if ShouldAutoApply(medium) {
		t.Fatal("expected medium-confidence scalar conflict to skip auto-apply")
	}
}
