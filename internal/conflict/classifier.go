package conflict

import (
	"regexp"
	"strings"

	"github.com/jhanvi857/gitresolve/internal/analysis"
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
	if analysis.IsStructuredFile(c.FilePath) {
		c.Type = TypeStructured
		if analysis.IsCriticalFile(c.FilePath) {
			c.Severity = SeverityHigh
			c.CanAutoResolve = false
		} else {
			c.Severity = SeverityLow
			c.CanAutoResolve = true
		}
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
	if isSignatureChange(c.FilePath, c.OurLines, c.TheirLines) {
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

	// rule 7: check for critical files (go.mod, etc. which might not be structured)
	if analysis.IsCriticalFile(c.FilePath) {
		c.Type = TypeLogic
		c.Severity = SeverityHigh
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
	// Refined heuristics with Regex to avoid false positives on normal strings
	goImport := regexp.MustCompile(`^(import\s*(?:\(\s*)?|"[a-zA-Z0-9_\-\./]+"|\s+"[a-zA-Z0-9_\-\./]+")`)
	jsImport := regexp.MustCompile(`^(import\s+.*from\s+['"].*['"]|require\(['"].*['"]\))`)
	pyImport := regexp.MustCompile(`^(import\s+[a-zA-Z0-9_\.]+|from\s+[a-zA-Z0-9_\.]+\s+import)`)
	javaImport := regexp.MustCompile(`^import\s+[a-zA-Z0-9_\.]+;*`)

	return goImport.MatchString(line) ||
		jsImport.MatchString(line) ||
		pyImport.MatchString(line) ||
		javaImport.MatchString(line)
}


func isDeleteModify(ours, theirs []string) bool {
	// one side has zero lines = deletion
	// other side has lines = modification
	ourEmpty := len(ours) == 0 || (len(ours) == 1 && strings.TrimSpace(ours[0]) == "")
	theirEmpty := len(theirs) == 0 || (len(theirs) == 1 && strings.TrimSpace(theirs[0]) == "")
	return ourEmpty != theirEmpty
}

func isSignatureChange(filePath string, ours, theirs []string) bool {
	// Fallback fast heuristic
	ourHasFunc := containsFuncDef(ours)
	theirHasFunc := containsFuncDef(theirs)

	if ourHasFunc && theirHasFunc {
		// Verify using AST if it really contains function declarations
		ourAST, err1 := analysis.ParseFile(filePath, []byte(strings.Join(ours, "\n")))
		theirAST, err2 := analysis.ParseFile(filePath, []byte(strings.Join(theirs, "\n")))

		if err1 == nil && err2 == nil && ourAST != nil && theirAST != nil {
			ourNodes := analysis.FindChangedNodes(ourAST, 0, len(ours)+1)
			theirNodes := analysis.FindChangedNodes(theirAST, 0, len(theirs)+1)

			hasFuncChange := func(nodes []*analysis.Node) bool {
				for _, n := range nodes {
					if n.Type == "function_declaration" || n.Type == "method_declaration" || n.Type == "arrow_function" || n.Type == "lexical_declaration" || n.Type == "ERROR" {
						return true
					}
				}
				return false
			}

			if hasFuncChange(ourNodes) && hasFuncChange(theirNodes) {
				return true
			}
		}
		
		// If AST fails or nodes not parsed perfectly cleanly (due to incomplete snippet), return true since heuristic passed
		return true
	}
	
	return false
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
			strings.HasPrefix(trimmed, "type ") ||
			strings.HasPrefix(trimmed, "interface ") ||
			strings.HasPrefix(trimmed, "class ") ||
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
