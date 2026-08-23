package conflict

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type VerificationError struct {
	File        string
	Reason      string
	Output      string // the invalid content, for debugging
	IsMarkerErr bool   // true if failure was due to remaining markers
}

func (e *VerificationError) Error() string {
	return fmt.Sprintf("verification failed for %s: %s", e.File, e.Reason)
}

func HasConflictMarkers(content string) bool {
	return strings.Contains(content, "<<<<<<<") ||
		strings.Contains(content, "=======") ||
		strings.Contains(content, ">>>>>>>") ||
		strings.Contains(content, "|||||||")
}

func EnsureNoConflictMarkers(filePath, content string) error {
	return checkNoMarkers(filePath, content)
}

// Verify checks that a resolved file is actually valid
func Verify(filePath, content string) error {
	if err := checkNoMarkers(filePath, content); err != nil {
		return err
	}

	if strings.HasSuffix(filePath, ".json") {
		if err := verifyJSON(filePath, content); err != nil {
			return err
		}
	}
	if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml") {
		if err := verifyYAML(filePath, content); err != nil {
			return err
		}
	}
	if strings.HasSuffix(filePath, ".toml") {
		if err := verifyTOML(filePath, content); err != nil {
			return err
		}
	}
	if strings.HasSuffix(filePath, ".go") {
		if err := verifyGo(filePath, content); err != nil {
			return err
		}
	}

	return nil
}

func checkNoMarkers(filePath, content string) error {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "<<<<<<<") ||
			strings.HasPrefix(line, ">>>>>>>") ||
			strings.HasPrefix(line, "|||||||") ||
			line == "=======" {
			return &VerificationError{
				File:        filePath,
				Reason:      fmt.Sprintf("conflict marker found on line %d", i+1),
				Output:      content,
				IsMarkerErr: true,
			}
		}
	}
	return nil
}

func verifyJSON(filePath, content string) error {
	var v interface{}
	if err := json.Unmarshal([]byte(content), &v); err != nil {
		// Try wrapping as an object snippet (trimming any trailing comma)
		trimmed := strings.TrimSpace(content)
		trimmedObj := strings.TrimSuffix(trimmed, ",")
		wrappedObj := "{" + trimmedObj + "}"
		if err2 := json.Unmarshal([]byte(wrappedObj), &v); err2 == nil {
			return nil
		}
		// Try wrapping as an array snippet
		wrappedArr := "[" + trimmedObj + "]"
		if err3 := json.Unmarshal([]byte(wrappedArr), &v); err3 == nil {
			return nil
		}
		return &VerificationError{
			File:   filePath,
			Reason: fmt.Sprintf("invalid JSON: %v", err),
			Output: content,
		}
	}
	return nil
}

func verifyYAML(filePath, content string) error {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(content), &root); err != nil {
		return &VerificationError{
			File:   filePath,
			Reason: fmt.Sprintf("invalid YAML: %v", err),
			Output: content,
		}
	}

	if err := detectYAMLDuplicateKeys(&root); err != nil {
		return &VerificationError{
			File:   filePath,
			Reason: err.Error(),
			Output: content,
		}
	}
	return nil
}

func verifyTOML(filePath, content string) error {
	var v map[string]interface{}
	if err := toml.Unmarshal([]byte(content), &v); err != nil {
		return &VerificationError{
			File:   filePath,
			Reason: fmt.Sprintf("invalid TOML: %v", err),
			Output: content,
		}
	}
	return nil
}

func verifyGo(filePath, content string) error {
	if err := ValidateGoSyntax(filePath, content); err != nil {
		return &VerificationError{
			File:   filePath,
			Reason: fmt.Sprintf("invalid Go syntax: %v", err),
			Output: content,
		}
	}
	return nil
}

func ValidateGoSyntax(filePath, content string) error {
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, filePath, content, parser.AllErrors)
	if err != nil {
		if strings.Contains(err.Error(), "expected 'package'") {
			// Try wrapping as a package-level declaration snippet
			wrappedDecl := "package p\n" + content
			_, errDecl := parser.ParseFile(fset, filePath, wrappedDecl, parser.AllErrors)
			if errDecl == nil {
				return nil
			}
			// If it's a statement snippet (e.g. inside a function body), try wrapping inside a function
			wrappedStmt := "package p\nfunc _() {\n" + content + "\n}"
			_, errStmt := parser.ParseFile(fset, filePath, wrappedStmt, parser.AllErrors)
			if errStmt == nil {
				return nil
			}
			return errDecl
		}
		return err
	}

	tmp, err := os.CreateTemp("", "gitresolve-vet-*.go")
	if err != nil {
		return fmt.Errorf("vet: temp file: %w", err)
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()
	if _, err := tmp.Write([]byte(content)); err != nil {
		return fmt.Errorf("vet: write: %w", err)
	}
	tmp.Close()
	cmd := exec.Command("go", "vet", tmp.Name()) // #nosec G204 -- fixed command, temp file path is generated locally
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go vet failed: %s", out)
	}
	return nil
}

func detectYAMLDuplicateKeys(node *yaml.Node) error {
	if node == nil {
		return nil
	}

	if node.Kind == yaml.MappingNode {
		seen := make(map[string]struct{})
		for i := 0; i+1 < len(node.Content); i += 2 {
			k := node.Content[i]
			v := node.Content[i+1]
			if _, ok := seen[k.Value]; ok {
				return fmt.Errorf("duplicate YAML key %q at line %d", k.Value, k.Line)
			}
			seen[k.Value] = struct{}{}
			if err := detectYAMLDuplicateKeys(v); err != nil {
				return err
			}
		}
		return nil
	}

	for _, child := range node.Content {
		if err := detectYAMLDuplicateKeys(child); err != nil {
			return err
		}
	}

	return nil
}
