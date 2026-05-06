package safepath

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSafeOpen_NormalPath(t *testing.T) {
	repo := t.TempDir()
	file := filepath.Join(repo, "nested", "ok.txt")
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(file, []byte("ok"), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	root, err := RepoRoot(repo)
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	defer root.Close()

	f, err := SafeOpen(root, filepath.Join("nested", "ok.txt"))
	if err != nil {
		t.Fatalf("SafeOpen normal path: %v", err)
	}
	defer f.Close()

	got, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(got) != "ok" {
		t.Fatalf("unexpected content: %q", string(got))
	}
}

func TestSafeOpen_RejectsParentTraversal(t *testing.T) {
	root, err := RepoRoot(t.TempDir())
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	defer root.Close()

	_, err = SafeOpen(root, filepath.Join("..", "escape.txt"))
	if !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("expected ErrUnsafePath, got %v", err)
	}
}

func TestSafeOpen_RejectsAbsolutePath(t *testing.T) {
	root, err := RepoRoot(t.TempDir())
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	defer root.Close()

	abs := "/tmp/escape.txt"
	if runtime.GOOS == "windows" {
		abs = `C:\\tmp\\escape.txt`
	}

	_, err = SafeOpen(root, abs)
	if !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("expected ErrUnsafePath, got %v", err)
	}
}

func TestSafeOpen_RejectsSymlinkEscape(t *testing.T) {
	repo := t.TempDir()
	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "secret.txt")
	if err := os.WriteFile(outsideFile, []byte("secret"), 0o644); err != nil {
		t.Fatalf("seed outside file: %v", err)
	}

	linkPath := filepath.Join(repo, "link-out.txt")
	if err := os.Symlink(outsideFile, linkPath); err != nil {
		if runtime.GOOS == "windows" && strings.Contains(strings.ToLower(err.Error()), "privilege") {
			t.Skipf("symlink creation requires privilege on this Windows environment: %v", err)
		}
		t.Fatalf("create symlink: %v", err)
	}

	root, err := RepoRoot(repo)
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	defer root.Close()

	_, err = SafeOpen(root, "link-out.txt")
	if !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("expected ErrUnsafePath for symlink escape, got %v", err)
	}
}

func TestSafeWrite_SymlinkSwapCannotEscapeRoot(t *testing.T) {
	repo := t.TempDir()
	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "outside.txt")
	if err := os.WriteFile(outsideFile, []byte("outside-original"), 0o644); err != nil {
		t.Fatalf("seed outside file: %v", err)
	}

	linkPath := filepath.Join(repo, "victim.txt")
	if err := os.Symlink(outsideFile, linkPath); err != nil {
		if runtime.GOOS == "windows" && strings.Contains(strings.ToLower(err.Error()), "privilege") {
			t.Skipf("symlink creation requires privilege on this Windows environment: %v", err)
		}
		t.Fatalf("create symlink: %v", err)
	}

	root, err := RepoRoot(repo)
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	defer root.Close()

	if err := SafeWrite(root, "victim.txt", []byte("inside-new"), 0o644); err != nil {
		t.Fatalf("SafeWrite through symlink path: %v", err)
	}

	outsideData, err := os.ReadFile(outsideFile)
	if err != nil {
		t.Fatalf("read outside file: %v", err)
	}
	if string(outsideData) != "outside-original" {
		t.Fatalf("outside file was modified unexpectedly: %q", outsideData)
	}

	insidePath := filepath.Join(repo, "victim.txt")
	insideData, err := os.ReadFile(insidePath)
	if err != nil {
		t.Fatalf("read inside file: %v", err)
	}
	if string(insideData) != "inside-new" {
		t.Fatalf("inside file mismatch: %q", insideData)
	}

	info, err := os.Lstat(insidePath)
	if err != nil {
		t.Fatalf("lstat inside file: %v", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Fatal("inside path remains a symlink; expected regular file after safe atomic rename")
	}
}
