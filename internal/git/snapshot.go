package git

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jhanvi857/gitresolve/internal/safepath"
)

// purpose : if user hitted ctrl +c in between then have to return back to the curernt state thatt's why need to store the head in disk and not in variable.
func storeHead(root *os.Root, head string) error {
	headFilePath := filepath.Join(".git", "gitresolve_head")
	// safepath: CWE-22 hardened
	err := safepath.SafeWrite(root, headFilePath, []byte(head), 0644)
	if err != nil {
		return fmt.Errorf("snapshot : storeHead %w", err)
	}
	return nil
}

func StoreHead(root *os.Root, head string) error {
	return storeHead(root, head)
}

func getStoredHead(root *os.Root) (string, error) {
	headFilePath := filepath.Join(".git", "gitresolve_head")
	// safepath: CWE-22 hardened
	f, err := root.Open(headFilePath)
	if err != nil {
		return "", fmt.Errorf("snapshot : getStoredHead %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("snapshot : reading head %w", err)
	}
	return string(data), nil
}

func GetStoredHead(root *os.Root) (string, error) {
	return getStoredHead(root)
}

func clearStoredHead(root *os.Root) error {
	headFilePath := filepath.Join(".git", "gitresolve_head")
	// safepath: CWE-22 hardened
	err := root.Remove(headFilePath)
	if err != nil {
		return fmt.Errorf("snapshot : clearstoredHead %w", err)
	}
	return nil
}

func ClearStoredHead(root *os.Root) error {
	return clearStoredHead(root)
}
