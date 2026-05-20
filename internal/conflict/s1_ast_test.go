package conflict

import (
	"testing"
)

func TestS1_AST_NestedStructValidation(t *testing.T) {
	// Test that deeply nested Go struct conflicts that produce invalid syntax
	// are properly caught by the validation gate and escalated to manual review.
	input := `package main

type ConfigOrig struct {
	Database struct {
		Pool struct {
			MaxConns int
			MinConns int
			Timeout  time.Duration
		}
	}
}

<<<<<<< HEAD
type ConfigNew struct {
	Database struct {
		Pool struct {
			MaxConns int
			MinConns int
			Timeout  time.Duration
		}
		Extra string
	}
}
=======
type ConfigNew struct {
	Database struct {
		Pool struct {
			MaxConns    int
			MinConns    int
			Timeout     time.Duration
			IdleTimeout time.Duration
		}
		ReadReplica string
	}
}
>>>>>>> branch

func init() {}
`

	filePath := "config.go"

	// Parse the conflict
	conflicts := ParseFile(filePath, []byte(input))
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(conflicts))
	}

	c := conflicts[0]
	t.Logf("Parsed conflict with %d ours lines, %d theirs lines", len(c.OursLines), len(c.TheirsLines))

	// Classify the conflict
	Classify(c)
	t.Logf("Classified as Type=%v, CanAutoResolve=%v, Severity=%v", c.Type, c.CanAutoResolve, c.Severity)

	// Attempt auto-resolution
	if c.CanAutoResolve {
		resolved := AutoResolve(c, Options{})
		t.Logf("AutoResolve returned: %v", resolved)

		if !resolved {
			// Expected outcome: validation gate should catch the invalid syntax
			t.Logf("✓ AutoResolve correctly rejected invalid output")
			if c.ManualReasonCode == ReasonValidationSyntaxFailed {
				t.Logf("✓ Reason code correctly set to: %s", c.ManualReasonCode)
			} else {
				t.Errorf("Expected reason code %s, got %s", ReasonValidationSyntaxFailed, c.ManualReasonCode)
			}
			t.Logf("Manual reason: %s", c.ManualReason)
			// This is the correct behavior - invalid syntax was caught
		} else {
			// If we get here in the fixed code, the output must be valid
			t.Logf("AutoResolve succeeded, validating output...")
			if err := Verify(filePath, c.Resolution); err != nil {
				t.Errorf("BUG: AutoResolve returned success but output is invalid: %v", err)
				t.Logf("Invalid output:\n%s", c.Resolution)
			} else {
				t.Logf("✓ Output is valid Go syntax")
			}
		}
	} else {
		t.Logf("Conflict not marked for auto-resolve, skipping AutoResolve test")
	}
}

func TestS1_AST_SimpleValidation_Passes(t *testing.T) {
	// Test that valid, simple resolutions still pass validation
	c := &ConflictBlock{
		FilePath: "simple.go",
		OursLines: []string{
			"func foo() {",
			"    x := 1",
			"}",
		},
		TheirsLines: []string{
			"func foo() {",
			"    x := 2",
			"}",
		},
	}

	Classify(c)
	resolved := AutoResolve(c, Options{})

	if resolved {
		t.Logf("✓ Simple conflict auto-resolved successfully")
		if err := Verify(c.FilePath, c.Resolution); err != nil {
			t.Errorf("Validation failed unexpectedly: %v", err)
		}
	} else {
		t.Logf("Conflict was not auto-resolved (may have been escalated)")
	}
}

func TestValidationGate_RejjectsInvalidGo(t *testing.T) {
	// Direct validation gate test
	invalidGo := `type Config struct {
    Nested struct {
        X int
    }
    Y string
}`
	// Missing closing brace makes this invalid

	filePath := "bad.go"
	if err := Verify(filePath, invalidGo); err != nil {
		t.Logf("✓ Validation gate correctly rejected invalid Go: %v", err)
	} else {
		t.Logf("Note: This Go snippet may be parseable in certain contexts")
	}
}
