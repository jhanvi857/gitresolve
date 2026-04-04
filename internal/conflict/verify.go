package conflict

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// Verifying checks that a resolved file is actually valid
func Verify(filePath, content string) error {
	if err := checkNoMarkers(content); err != nil {
		return err
	}
	if strings.HasSuffix(filePath, ".json") {
		return verifyJSON(content)
	}
	if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml") {
		return verifyYAML(content)
	}
	if strings.HasSuffix(filePath, ".toml") {
		return verifyTOML(content)
	}
	if strings.HasSuffix(filePath, ".go") {
		return verifyGo(content)
	}
	return nil
}

func hasMarkers(content string) bool {
	return strings.Contains(content, "<<<<<<<") ||
		strings.Contains(content, "=======") ||
		strings.Contains(content, ">>>>>>>")
}

func checkNoMarkers(content string) error {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "<<<<<<<") ||
			strings.HasPrefix(line, "=======") ||
			strings.HasPrefix(line, ">>>>>>>") {
			return fmt.Errorf("verify: conflict marker found at line %d", i+1)
		}
	}
	return nil
}

// json.Unmarshal : Go's standard JSON parser. It takes a string and tries to parse it into a Go value. If the string is not valid JSON it returns an error describing exactly what is wrong.
// var v interface{} : why interface? coz we don't care what the JSON contains. we just want to know if it parses successfully. interface{} accepts any valid JSON structure - object, array, string, number, anything.
func verifyJSON(content string) error {
	var v interface{}
	if err := json.Unmarshal([]byte(content), &v); err != nil {
		return fmt.Errorf("verify: invalid JSON after resolution: %w", err)
	}
	return nil
}

// same as json. but yaml is more indentation sensitive
func verifyYAML(content string) error {
	var v interface{}
	if err := yaml.Unmarshal([]byte(content), &v); err != nil {
		return fmt.Errorf("verify: invalid YAML after resolution: %w", err)
	}
	return nil
}

func verifyTOML(content string) error {
	var v map[string]interface{}
	if err := toml.Unmarshal([]byte(content), &v); err != nil {
		return fmt.Errorf("verify: invalid TOML after resolution: %w", err)
	}
	return nil
}

func verifyGo(content string) error {
	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, "resolved.go", content, parser.AllErrors); err != nil {
		return fmt.Errorf("verify: invalid Go syntax after resolution: %w", err)
	}
	return nil
}
