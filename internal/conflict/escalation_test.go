package conflict

import "testing"

func TestRenderEscalationMessage(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		data     map[string]string
		expected string
	}{
		{
			name: "high blast radius",
			code: ReasonSemanticHighBlastRadius,
			data: map[string]string{
				"symbol": "ProcessPayment",
				"count":  "14",
			},
			expected: "ProcessPayment is called from 14 other locations — escalating for manual review",
		},
		{
			name: "missing coupled file",
			code: ReasonSemanticMissingCoupledFile,
			data: map[string]string{
				"file":         "models/user.go",
				"coupled_file": "handlers/user.go",
				"strength":     "0.82",
			},
			expected: "models/user.go is frequently changed alongside handlers/user.go (strength 0.82) but handlers/user.go was not touched in this branch",
		},
		{
			name: "import cycle",
			code: ReasonSemanticImportCycle,
			data: map[string]string{
				"file":  "pkg/auth/token.go",
				"cycle": "pkg/auth → pkg/user → pkg/auth",
			},
			expected: "conflict block in pkg/auth/token.go sits inside an import cycle: pkg/auth → pkg/user → pkg/auth",
		},
		{
			name: "stale branch divergence",
			code: ReasonStrategyStaleBranchDiv,
			data: map[string]string{
				"behind":  "23",
				"branch":  "main",
				"authors": "alice@co.com, bob@co.com",
			},
			expected: "branch is 23 commits behind main — alice@co.com, bob@co.com authored changes touching files you also modified",
		},
		{
			name: "multi author conflict",
			code: ReasonStrategyMultiAuthorConflict,
			data: map[string]string{
				"file":          "internal/api/handler.go",
				"author_ours":   "alice@co.com",
				"author_theirs": "bob@co.com",
			},
			expected: "both sides of the conflict in internal/api/handler.go were authored by different people: alice@co.com vs bob@co.com",
		},
		{
			name:     "unknown code returns empty",
			code:     "nonexistent.code",
			data:     map[string]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderEscalationMessage(tt.code, tt.data)
			if got != tt.expected {
				t.Errorf("RenderEscalationMessage(%q) =\n  %q\nwant:\n  %q", tt.code, got, tt.expected)
			}
		})
	}
}

func TestRenderSuggestedCommand(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		data     map[string]string
		expected string
	}{
		{
			name:     "high blast radius",
			code:     ReasonSemanticHighBlastRadius,
			data:     map[string]string{},
			expected: "go test ./... (run full suite before committing)",
		},
		{
			name: "missing coupled file",
			code: ReasonSemanticMissingCoupledFile,
			data: map[string]string{
				"coupled_file": "handlers/user.go",
			},
			expected: "gitresolve status handlers/user.go",
		},
		{
			name: "stale branch divergence rebase",
			code: ReasonStrategyStaleBranchDiv,
			data: map[string]string{
				"branch": "main",
			},
			expected: "git fetch && git rebase origin/main",
		},
		{
			name: "multi author conflict",
			code: ReasonStrategyMultiAuthorConflict,
			data: map[string]string{
				"author_theirs": "bob@co.com",
			},
			expected: "coordinate with bob@co.com before resolving",
		},
		{
			name:     "unknown code returns fallback",
			code:     "nonexistent.code",
			data:     map[string]string{},
			expected: "git merge --abort",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderSuggestedCommand(tt.code, tt.data)
			if got != tt.expected {
				t.Errorf("RenderSuggestedCommand(%q) =\n  %q\nwant:\n  %q", tt.code, got, tt.expected)
			}
		})
	}
}

func TestFallbackSuggestedCommand(t *testing.T) {
	got := FallbackSuggestedCommand()
	expected := "git merge --abort"
	if got != expected {
		t.Errorf("FallbackSuggestedCommand() = %q, want %q", got, expected)
	}
}
