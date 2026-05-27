package ownership

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePolicy_BinarySearch(t *testing.T) {
	// We need to bypass LoadPolicyConfig because it reads from disk.
	// But ResolvePolicy calls it.
	// We can use a temp dir.
	tmpDir := t.TempDir()
	dotGitResolve := filepath.Join(tmpDir, ".gitresolve")
	if err := os.MkdirAll(dotGitResolve, 0o755); err != nil {
		t.Fatalf("mkdir .gitresolve: %v", err)
	}

	policyContent := `{
		"default": "balanced",
		"path_profiles": {
			"a/": "strict",
			"a/b/": "aggressive",
			"x/": "auto"
		}
	}`
	if err := os.WriteFile(filepath.Join(dotGitResolve, "policy.json"), []byte(policyContent), 0o644); err != nil {
		t.Fatalf("write policy.json: %v", err)
	}

	tests := []struct {
		filePath string
		expected string
	}{
		{"a/file.go", "strict"},
		{"a/b/file.go", "aggressive"},
		{"a/c/file.go", "strict"},
		{"x/file.go", "auto"}, // auto is a valid profile in config
		{"y/file.go", "balanced"},
	}

	root, err := os.OpenRoot(tmpDir)
	if err != nil {
		t.Fatalf("open root: %v", err)
	}
	defer root.Close()

	for _, tt := range tests {
		res, err := ResolvePolicy(root, tt.filePath, "auto")
		if err != nil {
			t.Errorf("unexpected error for %s: %v", tt.filePath, err)
			continue
		}
		if res.ResolvedProfile != tt.expected {
			t.Errorf("for %s: expected %s, got %s", tt.filePath, tt.expected, res.ResolvedProfile)
		}
	}
}
