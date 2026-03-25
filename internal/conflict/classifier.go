package conflict

import (
	"path/filepath"
	"strings"
)

func Classify(c *Conflict) {
	// rule 1: whitespace only
	// strip all whitespace from both sides and compare
	// if they are identical after stripping = whitespace conflict
	if isWhitespaceOnly(c.OurLines, c.TheirLines) {
		c.Type = TypeWhitespace
		c.Severity = SeverityTrivial
		c.CanAutoResolve = true
		return
	}

	// rule 2: both sides made identical changes
	// this happens when two devs independently fix the same bug
	if linesIdentical(c.OurLines, c.TheirLines) {
		c.Type = TypeIdentical
		c.Severity = SeverityTrivial
		c.CanAutoResolve = true
		return
	}

	// rule 3: import block conflict
	// all changed lines on both sides are import statements
	if isImportConflict(c.OurLines, c.TheirLines) {
		c.Type = TypeImport
		c.Severity = SeverityLow
		c.CanAutoResolve = true
		return
	}

	// rule 4: structured file conflict
	// JSON/YAML/TOML : handled by structured.go not line diff
	if isStructuredFile(c.FilePath) {
		c.Type = TypeStructured
		c.Severity = SeverityLow
		c.CanAutoResolve = false
		return
	}

	// rule 5: delete vs modify
	// one side has no lines (deletion) other side has lines (modification)
	// this is dangerous : someone deleted something the other person was using
	if isDeleteModify(c.OurLines, c.TheirLines) {
		c.Type = TypeDeleteModify
		c.Severity = SeverityCritical
		c.CanAutoResolve = false
		return
	}

	// rule 6: function signature change
	// check if lines contain function definition keywords
	if isSignatureChange(c.OurLines, c.TheirLines) {
		c.Type = TypeSignature
		c.Severity = SeverityHigh
		c.CanAutoResolve = false
		return
	}

	// rule 7: check file path for sensitive areas
	// auth, security, crypto, migration files get elevated severity
	if isSensitivePath(c.FilePath) {
		c.Type = TypeLogic
		c.Severity = SeverityCritical
		c.CanAutoResolve = false
		return
	}

	// default: logic conflict, medium severity, needs human review
	c.Type = TypeLogic
	c.Severity = SeverityMedium
	c.CanAutoResolve = false
}

func isWhitespaceOnly(ours, theirs []string) bool {
	// strip all whitespace from every line on both sides
	// if the remaining content is identical = only whitespace differs
	ourStripped := stripWhitespace(ours)
	theirStripped := stripWhitespace(theirs)
	return ourStripped == theirStripped
}

// stripWhitespace removes all spaces and tabs from lines
// joins everything into one string for easy comparison
func stripWhitespace(lines []string) string {
	var result strings.Builder
	for _, line := range lines {
		// strings.ReplaceAll replaces every occurrence of " " with ""
		line = strings.ReplaceAll(line, " ", "")
		line = strings.ReplaceAll(line, "\t", "")
		result.WriteString(line)
	}
	return result.String()
}

func linesIdentical(ours, theirs []string) bool {
	if len(ours) != len(theirs) {
		return false
	}
	for i := range ours {
		if ours[i] != theirs[i] {
			return false
		}
	}
	return true
}

func isImportConflict(ours, theirs []string) bool {
	// every line on both sides must look like an import statement
	// we check for common import patterns across languages
	allLines := append(ours, theirs...)
	for _, line := range allLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !isImportLine(trimmed) {
			return false
		}
	}
	return len(allLines) > 0
}

func isImportLine(line string) bool {
	// Go:         import "fmt"  or  "github.com/..."
	// JavaScript: import React from 'react'
	// Python:     import os  or  from os import path
	// Java:       import java.util.List;
	return strings.HasPrefix(line, "import ") ||
		strings.HasPrefix(line, "from ") ||
		strings.HasPrefix(line, "\"") ||
		strings.HasPrefix(line, "'")
}

func isStructuredFile(filePath string) bool {
	ext := filepath.Ext(filePath)
	return ext == ".json" || ext == ".yaml" || ext == ".yml" || ext == ".toml"
}

func isDeleteModify(ours, theirs []string) bool {
	// one side has zero lines = deletion
	// other side has lines = modification
	ourEmpty := len(ours) == 0 || (len(ours) == 1 && strings.TrimSpace(ours[0]) == "")
	theirEmpty := len(theirs) == 0 || (len(theirs) == 1 && strings.TrimSpace(theirs[0]) == "")
	return ourEmpty != theirEmpty
}

func isSignatureChange(ours, theirs []string) bool {
	// check if any line looks like a function signature
	// if both sides changed lines that contain function definitions
	// that is a signature conflict
	ourHasFunc := containsFuncDef(ours)
	theirHasFunc := containsFuncDef(theirs)
	return ourHasFunc && theirHasFunc
}

func containsFuncDef(lines []string) bool {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Go:         func FunctionName(
		// JavaScript: function name( or const name = ( or ) =>
		// Python:     def function_name(
		// Java/C#:    public/private/protected ... (
		if strings.HasPrefix(trimmed, "func ") ||
			strings.HasPrefix(trimmed, "function ") ||
			strings.HasPrefix(trimmed, "def ") ||
			strings.HasPrefix(trimmed, "public ") ||
			strings.HasPrefix(trimmed, "private ") ||
			strings.Contains(trimmed, "=> {") {
			return true
		}
	}
	return false
}

func isSensitivePath(filePath string) bool {
	// file paths containing these words get critical severity
	// because bugs in these areas have serious consequences
	sensitivePatterns := []string{
		"auth", "security", "crypto", "password",
		"token", "secret", "migration", "payment",
		"billing", "admin",
	}

	lowerPath := strings.ToLower(filePath)
	for _, pattern := range sensitivePatterns {
		if strings.Contains(lowerPath, pattern) {
			return true
		}
	}
	return false
}
