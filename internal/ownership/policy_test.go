package ownership

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePolicyProfile_DefaultBalanced(t *testing.T) {
	tmp := t.TempDir()
	profile, err := ResolvePolicyProfile(tmp, "internal/foo/bar.go", PolicyAuto)
	if err != nil {
		t.Fatalf("ResolvePolicyProfile failed: %v", err)
	}
	if profile != PolicyBalanced {
		t.Fatalf("expected default balanced profile, got %q", profile)
	}
}

func TestResolvePolicyProfile_PathAndTeam(t *testing.T) {
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".gitresolve"), 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	ownersJSON := `{"owners":{"internal/auth/":"security"}}`
	if err := os.WriteFile(filepath.Join(tmp, ".gitresolve", "owners.json"), []byte(ownersJSON), 0644); err != nil {
		t.Fatalf("write owners failed: %v", err)
	}
	policyJSON := `{
  "default":"balanced",
  "path_profiles":{"internal/auth/":"strict"},
  "team_profiles":{"security":"aggressive"}
}`
	if err := os.WriteFile(filepath.Join(tmp, ".gitresolve", "policy.json"), []byte(policyJSON), 0644); err != nil {
		t.Fatalf("write policy failed: %v", err)
	}

	profile, err := ResolvePolicyProfile(tmp, "internal/auth/token.go", PolicyAuto)
	if err != nil {
		t.Fatalf("ResolvePolicyProfile failed: %v", err)
	}
	if profile != PolicyStrict {
		t.Fatalf("expected path policy strict to win, got %q", profile)
	}

	profile, err = ResolvePolicyProfile(tmp, "internal/other/file.go", PolicyAuto)
	if err != nil {
		t.Fatalf("ResolvePolicyProfile failed: %v", err)
	}
	if profile != PolicyBalanced {
		t.Fatalf("expected fallback default profile, got %q", profile)
	}

	explicit, err := ResolvePolicyProfile(tmp, "internal/auth/token.go", PolicyAggressive)
	if err != nil {
		t.Fatalf("ResolvePolicyProfile explicit failed: %v", err)
	}
	if explicit != PolicyAggressive {
		t.Fatalf("expected explicit aggressive policy, got %q", explicit)
	}
}
