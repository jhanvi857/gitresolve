package lock

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestAcquireLockContention(t *testing.T) {
	repoDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repoDir, ".gitresolve"), 0o755); err != nil {
		t.Fatalf("mkdir .gitresolve: %v", err)
	}
	root, err := os.OpenRoot(repoDir)
	if err != nil {
		t.Fatalf("open root: %v", err)
	}
	defer root.Close()

	// First acquire should succeed
	l1, err := Acquire(root)
	if err != nil {
		t.Fatalf("Expected first acquire to succeed, got %v", err)
	}
	defer func() { _ = l1.Release() }()

	// Second acquire should fail with ErrLockContention
	l2, err := Acquire(root)
	if err != ErrLockContention {
		t.Fatalf("Expected ErrLockContention, got %v", err)
	}
	if l2 != nil {
		_ = l2.Release()
		t.Fatal("Expected second lock to be nil")
	}

	// Release first lock
	if err := l1.Release(); err != nil {
		t.Fatalf("Failed to release first lock: %v", err)
	}

	// Third acquire should now succeed
	l3, err := Acquire(root)
	if err != nil {
		t.Fatalf("Expected third acquire to succeed after release, got %v", err)
	}
	defer func() { _ = l3.Release() }()
}

func TestAcquireConcurrency(t *testing.T) {
	repoDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repoDir, ".gitresolve"), 0o755); err != nil {
		t.Fatalf("mkdir .gitresolve: %v", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	failureCount := 0

	var locks []*RepoLock

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			root, _ := os.OpenRoot(repoDir)
			if root == nil {
				return
			}
			defer root.Close()
			l, err := Acquire(root)
			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				successCount++
				locks = append(locks, l)
			} else if err == ErrLockContention {
				failureCount++
			}
		}()
	}

	wg.Wait()

	for _, l := range locks {
		_ = l.Release()
	}

	if successCount != 1 {
		t.Errorf("Expected exactly 1 success, got %d", successCount)
	}
	if failureCount != 9 {
		t.Errorf("Expected exactly 9 failures, got %d", failureCount)
	}
}
