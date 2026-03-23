package git

import (
	"fmt"
	"os"
	"path/filepath"
)

// purpose : if user hitted ctrl +c in between then have to return back to the curernt state thatt's why need to store the head in disk and not in variable.
func storeHead(repoPath string, head string) error {

	dir := filepath.Join(repoPath, ".git", "gitresolve_head")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("storeHead: creating directory: %w", err)
	}
	headFilePath := filepath.Join(repoPath, ".git", "gitresolve_head")
	err = os.WriteFile(headFilePath, []byte(head), 0644)
	if err != nil {
		return fmt.Errorf("snapshot : storeHead %w", err)
	}
	return nil
}

func getStoredHead(repoPath string) (string, error) {
	headFilePath := filepath.Join(repoPath, ".git", "gitresolve_head")
	data, err := os.ReadFile(headFilePath)
	if err != nil {
		return "", fmt.Errorf("snapshot : getStoredHead %w", err)
	}
	return string(data), nil
}
func clearStoredHead(repoPath string) error {
	headFilePath := filepath.Join(repoPath, ".git", "gitresolve_head")
	err := os.Remove(headFilePath)
	if err != nil {
		return fmt.Errorf("snapshot : clearstoredHead %w", err)
	}
	return nil
}
