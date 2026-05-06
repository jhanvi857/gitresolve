package gitresolve

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jhanvi857/gitresolve/internal/conflict"
	"github.com/jhanvi857/gitresolve/internal/safepath"
)

func TestReadConflictFileWithLimit_SkipsAndLogsLargeFile(t *testing.T) {
	repo := t.TempDir()
	fileName := "huge-conflict.txt"
	filePath := filepath.Join(repo, fileName)

	// 11MB synthetic file to exceed the 10MB default gate.
	data := make([]byte, 11*1024*1024)
	for i := range data {
		data[i] = 'a'
	}
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	root, err := safepath.RepoRoot(repo)
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	defer root.Close()

	cfg := conflict.ResolverConfig{MaxFileBytes: conflict.DefaultMaxConflictFileBytes}
	logged := false
	var loggedSize int64

	content, skipped, size, err := readConflictFileWithLimit(root, fileName, cfg, func(file string, gotSize int64, gotCfg conflict.ResolverConfig) {
		logged = true
		loggedSize = gotSize
		if file != fileName {
			t.Fatalf("unexpected file in callback: %s", file)
		}
		if gotCfg.MaxFileBytes != conflict.DefaultMaxConflictFileBytes {
			t.Fatalf("unexpected max-file-bytes in callback: %d", gotCfg.MaxFileBytes)
		}
	})
	if err != nil {
		t.Fatalf("readConflictFileWithLimit: %v", err)
	}

	if !skipped {
		t.Fatal("expected file to be skipped by size gate")
	}
	if !logged {
		t.Fatal("expected oversized file to be logged via callback")
	}
	if loggedSize != size {
		t.Fatalf("logged size mismatch: logged=%d returned=%d", loggedSize, size)
	}
	if content != nil {
		t.Fatalf("expected no content for skipped file, got %d bytes", len(content))
	}
}
