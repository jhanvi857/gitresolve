package git

// imports :
// 1. os : read .git/objects directory
// 2. filepath : safe file handling
// 3. compress/zlib : decompressing git objects
// 4. io : reading streams
// 5. bytes : parsing binary content
// 6. strings : parsing commit or tree text
// 7. fmt : error wrapping

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Commit represents a single git commit object
// uppercase fields = public, other packages can read them
type Commit struct {
	SHA     string
	Tree    string
	Parent  string
	Author  string
	Message string
	Time    string
}

// TreeEntry : one item inside a git tree object
// can be either a file (blob) or a directory (another tree)
// "100644" = file, "040000" = directory, "100755" = executable
type TreeEntry struct {
	Mode string
	Name string
	SHA  string
	Type string
}

// readObject reads and decompresses any git object from .git/objects/
// returns the object type commit/tree/blob and the raw content bytes
func readObject(repoPath string, sha string) (string, []byte, error) {
	// git splits SHA into two parts for storage on disk
	// first 2 chars become the directory name
	// remaining 38 chars become the filename
	// example: "a3f2c1d9e4b87654321abcdef1234567890abcd"
	// stored at: .git/objects/a3/f2c1d9e4b87654321abcdef1234567890abcd
	dir := sha[:2]
	file := sha[2:]

	// on windows uses \ and on linux/mac uses /
	objectPath := filepath.Join(repoPath, ".git", "objects", dir, file)
	f, err := os.Open(objectPath)
	if err != nil {
		return "", nil, fmt.Errorf("readObject: opening %s: %w", sha, err)
	}
	// defer : run this when the function returns no matter what
	// guarantees file gets closed even if we return early due to error
	defer f.Close()

	// git compresses every object with zlib before saving to disk
	// zlib.NewReader wraps file and decompresses as we read
	// without this we get meaningless compressed bytes
	zlibReader, err := zlib.NewReader(f)
	if err != nil {
		return "", nil, fmt.Errorf("readObject: creating zlib reader for %s: %w", sha, err)
	}
	defer zlibReader.Close()
	rawData, err := io.ReadAll(zlibReader)
	if err != nil {
		return "", nil, fmt.Errorf("readObject: reading decompressed data for %s: %w", sha, err)
	}
	// the \0 is a null byte : used to separate header from content
	nullIndex := bytes.IndexByte(rawData, 0)
	if nullIndex == -1 {
		return "", nil, fmt.Errorf("readObject: malformed object %s, no null byte found", sha)
	}

	header := string(rawData[:nullIndex])
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return "", nil, fmt.Errorf("readObject: malformed header in object %s: %q", sha, header)
	}

	// parts[0] is the type: commit / tree / blob
	objectType := parts[0]
	content := rawData[nullIndex+1:]
	return objectType, content, nil
}

func parseCommit(sha string, content []byte) (Commit, error) {
	text := string(content)
	lines := strings.Split(text, "\n")
	commit := Commit{SHA: sha}
	pastHeaders := false
	var messageLines []string

	for _, line := range lines {
		if line == "" {
			pastHeaders = true
			continue
		}
		if pastHeaders {
			messageLines = append(messageLines, line)
			continue
		}
		if strings.HasPrefix(line, "tree ") {
			commit.Tree = strings.TrimPrefix(line, "tree ")
		}

		if strings.HasPrefix(line, "parent ") {
			commit.Parent = strings.TrimPrefix(line, "parent ")
		}

		if strings.HasPrefix(line, "author ") {
			authorFull := strings.TrimPrefix(line, "author ")
			authorParts := strings.Split(authorFull, " ")

			if len(authorParts) >= 3 {
				commit.Author = strings.Join(authorParts[:len(authorParts)-2], " ")
				commit.Time = authorParts[len(authorParts)-2]
			}
		}
	}
	commit.Message = strings.TrimSpace(strings.Join(messageLines, "\n"))

	return commit, nil
}

// parseTree turns raw tree content bytes into a slice of TreeEntry
// tree objects are BINARY not text : each entry is packed as:
// [mode] [space] [filename] [null byte] [20 raw binary SHA bytes]
// there are no newlines : we must read byte by byte
func parseTree(content []byte) ([]TreeEntry, error) {
	var entries []TreeEntry
	i := 0

	for i < len(content) {

		spaceIndex := bytes.IndexByte(content[i:], ' ')
		if spaceIndex == -1 {
			return nil, fmt.Errorf("parseTree: no space found at position %d", i)
		}
		mode := string(content[i : i+spaceIndex])

		// move past the space character
		i += spaceIndex + 1
		nullIndex := bytes.IndexByte(content[i:], 0)
		if nullIndex == -1 {
			return nil, fmt.Errorf("parseTree: no null byte found at position %d", i)
		}
		name := string(content[i : i+nullIndex])
		i += nullIndex + 1

		if i+20 > len(content) {
			return nil, fmt.Errorf("parseTree: not enough bytes for SHA at position %d", i)
		}

		rawSHA := content[i : i+20]

		// %x converts each byte to its two-char hex representation
		// []byte{0xa3, 0xf2} becomes the string "a3f2"
		// this gives us the familiar 40-character hex SHA string
		hexSHA := fmt.Sprintf("%x", rawSHA)

		// move past the 20 SHA bytes :  ready for next entry
		i += 20
		// "40000" or "040000" = directory = this entry points to another tree
		// everything else = file = this entry points to a blob
		entryType := "blob"
		if mode == "40000" || mode == "040000" {
			entryType = "tree"
		}

		entries = append(entries, TreeEntry{
			Mode: mode,
			Name: name,
			SHA:  hexSHA,
			Type: entryType,
		})
	}

	return entries, nil
}

// GetCommit fetches and parses a commit object by its SHA
// every other package in gitstitch calls this to read commit history
func GetCommit(r *Repository, sha string) (Commit, error) {
	objType, content, err := readObject(r.Path, sha)
	if err != nil {
		return Commit{}, fmt.Errorf("GetCommit: %w", err)
	}
	if objType != "commit" {
		return Commit{}, fmt.Errorf("GetCommit: expected commit got %s for SHA %s", objType, sha)
	}

	commit, err := parseCommit(sha, content)
	if err != nil {
		return Commit{}, fmt.Errorf("GetCommit: %w", err)
	}

	return commit, nil
}

// GetParentCommit fetches the commit that came before the given commit
// merge/base.go uses this to walk backwards through history to find the merge base
func GetParentCommit(r *Repository, sha string) (Commit, error) {
	// get the commit we were given
	commit, err := GetCommit(r, sha)
	if err != nil {
		return Commit{}, fmt.Errorf("GetParentCommit: %w", err)
	}

	// the very first commit in a repo has no parent
	// Parent field will be empty string
	if commit.Parent == "" {
		return Commit{}, fmt.Errorf("GetParentCommit: commit %s is the root commit, no parent exists", sha)
	}

	// fetch and return the parent commit using its SHA
	parent, err := GetCommit(r, commit.Parent)
	if err != nil {
		return Commit{}, fmt.Errorf("GetParentCommit: %w", err)
	}

	return parent, nil
}

// GetFileAtCommit : returns the raw content of a file as it existed at a specific commit
// this is the most important function in this file for conflict resolution
// conflict/classifier.go calls this to get the three versions of a conflicted file:
// base version, our version, their version
func GetFileAtCommit(r *Repository, commitSHA string, filePath string) ([]byte, error) {
	// step 1 :  get the commit to find its root tree SHA
	commit, err := GetCommit(r, commitSHA)
	if err != nil {
		return nil, fmt.Errorf("GetFileAtCommit: %w", err)
	}

	// step 2 : read and parse the root tree object
	// the tree tells us all files and directories at this commit
	objType, treeContent, err := readObject(r.Path, commit.Tree)
	if err != nil {
		return nil, fmt.Errorf("GetFileAtCommit: reading root tree: %w", err)
	}
	if objType != "tree" {
		return nil, fmt.Errorf("GetFileAtCommit: expected tree got %s", objType)
	}

	entries, err := parseTree(treeContent)
	if err != nil {
		return nil, fmt.Errorf("GetFileAtCommit: parsing root tree: %w", err)
	}

	// step 3 :  walk the file path one directory at a time
	// "src/config/app.js" splits into ["src", "config", "app.js"]
	// filepath.ToSlash makes sure we always use / even on Windows
	parts := strings.Split(filepath.ToSlash(filePath), "/")

	// currentEntries starts as root tree entries
	// as we walk into subdirectories we replace it with subtree entries
	currentEntries := entries

	for partIndex, part := range parts {
		found := false

		for _, entry := range currentEntries {
			// does this tree entry match the current path component
			if entry.Name != part {
				continue
			}

			found = true
			if partIndex == len(parts)-1 {
				if entry.Type != "blob" {
					return nil, fmt.Errorf("GetFileAtCommit: %s is a directory not a file", filePath)
				}

				// read the blob :  this contains the raw file content
				blobType, blobContent, err := readObject(r.Path, entry.SHA)
				if err != nil {
					return nil, fmt.Errorf("GetFileAtCommit: reading blob: %w", err)
				}
				if blobType != "blob" {
					return nil, fmt.Errorf("GetFileAtCommit: expected blob got %s", blobType)
				}
				return blobContent, nil
			}
			if entry.Type != "tree" {
				return nil, fmt.Errorf("GetFileAtCommit: %s is not a directory", part)
			}
			_, subTreeContent, err := readObject(r.Path, entry.SHA)
			if err != nil {
				return nil, fmt.Errorf("GetFileAtCommit: reading subtree: %w", err)
			}

			currentEntries, err = parseTree(subTreeContent)
			if err != nil {
				return nil, fmt.Errorf("GetFileAtCommit: parsing subtree: %w", err)
			}
			break
		}

		if !found {
			return nil, fmt.Errorf("GetFileAtCommit: %q not found in tree", part)
		}
	}

	return nil, fmt.Errorf("GetFileAtCommit: %s not found in commit %s", filePath, commitSHA)
}

// ListChangedFiles returns which files were modified in a given commit
// compared to its parent commit
// index.go uses this to know which files need conflict checking
func ListChangedFiles(r *Repository, commitSHA string) ([]string, error) {
	commit, err := GetCommit(r, commitSHA)
	if err != nil {
		return nil, fmt.Errorf("ListChangedFiles: %w", err)
	}

	// get parent commit :  if this is the root commit return empty list
	parent, err := GetParentCommit(r, commitSHA)
	if err != nil {
		// root commit :  no parent exists, return empty not an error
		return []string{}, nil
	}

	// read and parse current commit tree
	_, currentContent, err := readObject(r.Path, commit.Tree)
	if err != nil {
		return nil, fmt.Errorf("ListChangedFiles: reading current tree: %w", err)
	}
	currentEntries, err := parseTree(currentContent)
	if err != nil {
		return nil, fmt.Errorf("ListChangedFiles: parsing current tree: %w", err)
	}

	// read and parse parent commit tree
	_, parentContent, err := readObject(r.Path, parent.Tree)
	if err != nil {
		return nil, fmt.Errorf("ListChangedFiles: reading parent tree: %w", err)
	}
	parentEntries, err := parseTree(parentContent)
	if err != nil {
		return nil, fmt.Errorf("ListChangedFiles: parsing parent tree: %w", err)
	}

	// build a map of filename -> SHA for the parent tree
	// map lookup is O(1) :  much faster than nested loops
	parentMap := make(map[string]string)
	for _, entry := range parentEntries {
		parentMap[entry.Name] = entry.SHA
	}

	// compare every file in current tree against parent tree
	var changedFiles []string
	for _, entry := range currentEntries {
		// skip directories :  only care about actual files
		if entry.Type != "blob" {
			continue
		}

		parentSHA, existedBefore := parentMap[entry.Name]

		if !existedBefore {
			// file is new :  did not exist in parent commit
			changedFiles = append(changedFiles, entry.Name)
			continue
		}

		if entry.SHA != parentSHA {
			// file existed but SHA changed :  content was modified
			// git generates SHA from file content
			// same content always produces same SHA
			// different SHA always means different content
			changedFiles = append(changedFiles, entry.Name)
		}
	}

	return changedFiles, nil
}
