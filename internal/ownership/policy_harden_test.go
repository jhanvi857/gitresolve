package ownership

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadPolicyConfig_Hardening(t *testing.T) {
	tmpDir := t.TempDir()
	dotGitResolve := filepath.Join(tmpDir, ".gitresolve")
	if err := os.MkdirAll(dotGitResolve, 0o755); err != nil {
		t.Fatalf("mkdir .gitresolve: %v", err)
	}

	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{
			name:    "oversized file",
			content: strings.Repeat(" ", 1024*1024+1),
			wantErr: ErrPolicyTooLarge.Error(),
		},
		{
			name:    "unknown key",
			content: `{"default": "strict", "unknown_key": "val"}`,
			wantErr: ErrPolicyUnknownKey("unknown_key").Error(),
		},
		{
			name: "too many path entries",
			content: func() string {
				c := `{"path_profiles": {`
				for i := 0; i < 501; i++ {
					c += fmt.Sprintf(`"p%d": "auto",`, i)
				}
				return c[:len(c)-1] + `}}`
			}(),
			wantErr: ErrPolicyTooManyEntries.Error(),
		},
		{
			name: "too many team entries",
			content: func() string {
				c := `{"team_profiles": {`
				for i := 0; i < 101; i++ {
					c += fmt.Sprintf(`"t%d": "auto",`, i)
				}
				return c[:len(c)-1] + `}}`
			}(),
			wantErr: ErrPolicyTooManyEntries.Error(),
		},
		{
			name:    "invalid profile value",
			content: `{"default": "super_strict"}`,
			wantErr: ErrPolicyInvalidProfile("super_strict").Error(),
		},
		{
			name: "valid complex policy",
			content: `{
				"default": "balanced",
				"path_profiles": {
					"src/": "strict",
					"tests/": "aggressive"
				},
				"team_profiles": {
					"core": "strict"
				}
			}`,
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dotGitResolve, "policy.json")
			if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("write policy.json: %v", err)
			}

			root, err := os.OpenRoot(tmpDir)
			if err != nil {
				t.Fatalf("open root: %v", err)
			}
			defer root.Close()
			_, loadErr := LoadPolicyConfig(root)
			if tt.wantErr != "" {
				if loadErr == nil {
					t.Fatal("expected error, got nil")
				}
				if loadErr.Error() != tt.wantErr {
					t.Errorf("expected error %q, got %q", tt.wantErr, loadErr.Error())
				}
			} else if loadErr != nil {
				t.Errorf("unexpected error: %v", loadErr)
			}
		})
	}
}
