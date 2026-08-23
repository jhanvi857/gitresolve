package conflict

import (
	"strings"
	"testing"
)

func TestClassifier_MalformedMarkersEscalatesManual(t *testing.T) {
	c := &ConflictBlock{
		FilePath: "main.go",
		OursLines: []string{
			"func ok() {",
			"<<<<<<< nested",
			"}",
		},
		TheirsLines: []string{"func ok() {}"},
	}

	Classify(c)

	if c.Type != TypeUnknown {
		t.Fatalf("expected TypeUnknown for malformed marker content, got %v", c.Type)
	}
	if c.Severity != SeverityCritical {
		t.Fatalf("expected SeverityCritical, got %v", c.Severity)
	}
	if c.CanAutoResolve {
		t.Fatal("expected malformed marker conflicts to disable auto resolve")
	}
	if !strings.Contains(c.ManualReason, "malformed conflict markers") {
		t.Fatalf("expected malformed marker manual reason, got %q", c.ManualReason)
	}
}

func TestClassifier_UnsupportedSemanticLanguageEscalatesManual(t *testing.T) {
	c := &ConflictBlock{
		FilePath: "engine.rs",
		OursLines: []string{
			"fn run(input: i32) -> i32 {",
			"    input + 1",
			"}",
		},
		TheirsLines: []string{
			"fn run(input: i32) -> i32 {",
			"    input + 2",
			"}",
		},
	}

	Classify(c)

	if c.Type != TypeUnknown {
		t.Fatalf("expected TypeUnknown for unsupported semantic language, got %v", c.Type)
	}
	if c.Severity != SeverityHigh {
		t.Fatalf("expected SeverityHigh, got %v", c.Severity)
	}
	if c.CanAutoResolve {
		t.Fatal("expected unsupported semantic language conflicts to disable auto resolve")
	}
	if !strings.Contains(c.ManualReason, "semantic resolver not available") {
		t.Fatalf("expected semantic coverage manual reason, got %q", c.ManualReason)
	}
}

func TestResolveBoth_HighRiskSemanticConflictBlocked(t *testing.T) {
	c := &ConflictBlock{
		FilePath: "auth/login.go",
		Type:     TypeSignature,
		Severity: SeverityHigh,
		OursLines: []string{
			"func Login(user string) error {",
			"    return nil",
			"}",
		},
		TheirsLines: []string{
			"func Login(user string, traceID string) error {",
			"    return nil",
			"}",
		},
	}

	_, err := Resolve(c, StrategyBoth, ResolveOptions{})
	if err == nil {
		t.Fatal("expected BOTH to be blocked for high-risk signature conflicts")
	}
	if !strings.Contains(err.Error(), "BOTH disabled") {
		t.Fatalf("expected BOTH disabled error, got %v", err)
	}
}
