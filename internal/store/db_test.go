package store

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenCorruptDB(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "corrupt.db")

	// Write garbage to the file
	err := os.WriteFile(dbPath, []byte("NOT A SQLITE DATABASE FILE"), 0644)
	if err != nil {
		t.Fatalf("failed to write garbage: %v", err)
	}

	// Attempt to open it
	_, err = Open(dbPath)
	if err == nil {
		t.Fatal("expected error when opening corrupt DB, got nil")
	}

	if !errors.Is(err, ErrDBCorrupt) {
		t.Errorf("expected ErrDBCorrupt, got %v", err)
	}
}
