package conflict

import "testing"

func applyStrategy(content []byte, strategy Strategy) (string, error) {
	conflicts := ParseFile("idempotency.go", content)
	for _, c := range conflicts {
		Classify(c)
		if _, err := Resolve(c, strategy, ResolveOptions{NonInteractive: true}); err != nil {
			return "", err
		}
	}
	return CompileResolution(content, conflicts), nil
}

func TestIdempotencyPerStrategy(t *testing.T) {
	base := []byte("package main\n<<<<<<< HEAD\nconst X = 1\n=======\nconst X = 2\n>>>>>>> branch\n")
	strategies := []Strategy{StrategyOurs, StrategyTheirs, StrategyBoth}

	for _, s := range strategies {
		first, err := applyStrategy(base, s)
		if err != nil {
			t.Fatalf("first pass failed for strategy %d: %v", s, err)
		}
		second, err := applyStrategy([]byte(first), s)
		if err != nil {
			t.Fatalf("second pass failed for strategy %d: %v", s, err)
		}
		if first != second {
			t.Fatalf("idempotency failed for strategy %d", s)
		}
	}
}

func TestStrategyConsistencyIsolation(t *testing.T) {
	base := []byte("package main\n<<<<<<< HEAD\nconst Y = 10\n=======\nconst Y = 20\n>>>>>>> branch\n")

	oursA, err := applyStrategy(base, StrategyOurs)
	if err != nil {
		t.Fatalf("ours A failed: %v", err)
	}
	theirs, err := applyStrategy(base, StrategyTheirs)
	if err != nil {
		t.Fatalf("theirs failed: %v", err)
	}
	oursB, err := applyStrategy(base, StrategyOurs)
	if err != nil {
		t.Fatalf("ours B failed: %v", err)
	}

	if oursA != oursB {
		t.Fatal("strategy contamination detected: repeated ours changed output")
	}
	if oursA == theirs {
		t.Fatal("expected different outputs for ours and theirs strategies")
	}
}
