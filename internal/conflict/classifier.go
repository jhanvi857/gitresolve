package conflict

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jhanvi857/gitresolve/internal/analysis"
)

const (
	AutoResolveConfidenceThreshold = 0.80
)

func Classify(c *ConflictBlock) {
	if hasEmbeddedConflictMarkers(c) || strings.Contains(c.ManualReason, "malformed conflict markers") {
		c.Type = TypeUnknown
		c.Severity = SeverityCritical
		c.Confidence = 0.05
		c.CanAutoResolve = false
		if c.ManualReason == "" || c.ManualReasonCode == "" {
			SetManualEscalation(c, ReasonParserMalformedNestedMarker, "malformed conflict markers detected", "prefer ours/theirs or manual edit for nested/irregular markers")
		}
		return
	}

	// rule 1: both sides made identical changes
	// this happens when two devs independently fix the same bug
	if linesIdentical(c.OursLines, c.TheirsLines) {
		c.Type = TypeIdentical
		c.Severity = SeverityTrivial
		c.Confidence = 0.99
		c.CanAutoResolve = true
		return
	}

	// rule 2: whitespace only
	// strip all whitespace from both sides and compare
	// if they are identical after stripping = whitespace conflict
	if isWhitespaceOnly(c.OursLines, c.TheirsLines) {
		c.Type = TypeWhitespace
		c.Severity = SeverityTrivial
		c.Confidence = 0.99
		c.CanAutoResolve = true
		return
	}

	// rule 3: import block conflict
	// all changed lines on both sides are import statements
	if isImportConflict(c.FilePath, c.OursLines, c.TheirsLines) {
		c.Type = TypeImport
		c.Severity = SeverityLow
		c.Confidence = 0.84
		if analysis.IsCriticalFile(c.FilePath) {
			c.Severity = SeverityHigh
			c.Confidence = 0.82
		}

		if containsComplexImports(c.FilePath, c.OursLines, c.TheirsLines) {
			c.Severity = SeverityMedium
			c.Confidence = 0.48
			c.CanAutoResolve = false // Fallback to manual for complex python/java imports
		} else {
			c.CanAutoResolve = true
		}
		return
	}

	// rule 4: structured file conflict
	// JSON/YAML/TOML : handled by structured.go not line diff
	if analysis.IsStructuredFile(c.FilePath) {
		c.Type = TypeStructured
		if analysis.IsCriticalFile(c.FilePath) {
			c.Severity = SeverityHigh
			c.Confidence = 0.83
			c.CanAutoResolve = true // Allow structured merger to attempt safe resolution
		} else {
			c.Severity = SeverityLow
			c.Confidence = 0.82
			c.CanAutoResolve = true
		}
		return
	}

	// rule 5: delete vs modify
	// one side has no lines (deletion) other side has lines (modification)
	// this is dangerous : someone deleted something the other person was using
	if isDeleteModify(c.OursLines, c.TheirsLines) {
		c.Type = TypeDeleteModify
		c.Severity = SeverityCritical
		c.Confidence = 0.10
		c.CanAutoResolve = false
		return
	}

	// rule 6: function signature change
	// check if lines contain function definition keywords
	if isSignatureChange(c.FilePath, c.OursLines, c.TheirsLines) {
		c.Type = TypeSignature
		c.Severity = SeverityHigh
		c.Confidence = 0.20
		c.CanAutoResolve = false
		return
	}

	// rule 7: check file path for sensitive areas
	// auth, security, crypto, migration files get elevated severity
	if isSensitivePath(c.FilePath) {
		c.Type = TypeLogic
		c.Severity = SeverityCritical
		c.Confidence = 0.12
		c.CanAutoResolve = false
		return
	}

	// rule 8: check for critical files (go.mod, etc. which might not be structured)
	if analysis.IsCriticalFile(c.FilePath) {
		c.Type = TypeLogic
		c.Severity = SeverityHigh
		c.Confidence = 0.58
		c.CanAutoResolve = false
		return
	}

	if isSourceLikeFile(c.FilePath) && !hasSemanticResolverCoverage(c.FilePath) {
		c.Type = TypeUnknown
		c.Severity = SeverityHigh
		c.Confidence = 0.35
		c.CanAutoResolve = false
		if c.ManualReason == "" || c.ManualReasonCode == "" {
			SetManualEscalation(c, ReasonSemanticUnsupportedLanguage, "language-specific semantic resolver not available for this file type", "use ours/theirs/manual and run language-native checks after merge")
		}
		return
	}

	if isSourceLikeFile(c.FilePath) && hasSemanticResolverCoverage(c.FilePath) && !semanticParserAvailable(c.FilePath) {
		c.Type = TypeUnknown
		c.Severity = SeverityHigh
		c.Confidence = 0.30
		c.CanAutoResolve = false
		SetManualEscalation(c, ReasonSemanticParseFailed, "semantic parser unavailable for this environment", "install parser/runtime support or resolve manually with ours/theirs")
		return
	}

	// rule 8: scalar change (single line, non-critical, non-signature)
	if isScalarChange(c.OursLines, c.TheirsLines) {
		c.Type = TypeScalar
		c.Severity = SeverityMedium
		c.Confidence = 0.58
		c.CanAutoResolve = false
		return
	}

	// default: logic conflict, medium severity, needs human review
	c.Type = TypeLogic
	c.Severity = SeverityMedium
	c.Confidence = 0.50
	c.CanAutoResolve = false
}

func ShouldAutoApply(c *ConflictBlock) bool {
	return c.CanAutoResolve && c.Confidence >= AutoResolveConfidenceThreshold
}

func isScalarChange(ours, theirs []string) bool {
	// A scalar change is a very small (1-line) modification to an existing line
	// that doesn't trigger signature detection, import detection, etc.
	// It's still safer to have human review, but we mark it as Scalar for better UX.
	return len(ours) == 1 && len(theirs) == 1
}

func containsComplexImports(filePath string, ours, theirs []string) bool {
	// Go imports will be handled correctly by go/ast in auto.go, so we don't block them here
	if strings.HasSuffix(filePath, ".go") {
		return false
	}

	allLines := append(ours, theirs...)
	for _, line := range allLines {
		trimmed := strings.TrimSpace(line)

		// Python: relative import or alias
		if strings.HasPrefix(trimmed, "from .") || strings.Contains(trimmed, " as ") {
			return true
		}

		// Java: wildcard import
		if strings.HasPrefix(trimmed, "import ") && strings.Contains(trimmed, "*;") {
			return true
		}
	}

	return false
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

func isImportConflict(filePath string, ours, theirs []string) bool {
	// every line on both sides must look like an import statement
	// we check for common import patterns across languages
	allLines := append(ours, theirs...)
	for _, line := range allLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !isImportLine(filePath, line) {
			return false
		}
	}
	return len(allLines) > 0
}

func isImportLine(filePath, line string) bool {
	// Refined heuristics with Regex to avoid false positives on normal strings
	goImport := regexp.MustCompile(`^(import\s*(?:\(\s*)?|(?:[a-zA-Z0-9_.]+\s+)?"[a-zA-Z0-9_\-\./]+"|\))`)
	jsImport := regexp.MustCompile(`^(import\s+.*from\s+['"].*['"]|require\(['"].*['"]\))`)
	pyImport := regexp.MustCompile(`^(import\s+[a-zA-Z0-9_\.]+|from\s+[a-zA-Z0-9_\.]+\s+import)`)
	javaImport := regexp.MustCompile(`^import\s+[a-zA-Z0-9_\.]+;*`)

	// go.mod support
	if strings.HasSuffix(strings.ToLower(filePath), "go.mod") {
		goMod := regexp.MustCompile(`^(require|module|go|retract|exclude|replace)(\s+|$)`)
		trimmed := strings.TrimSpace(line)
		if goMod.MatchString(line) || trimmed == "(" || trimmed == ")" || strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") || strings.Contains(line, " v") {
			return true
		}
	}

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
	sensitivePatterns := []string{
		"auth", "security", "crypto", "password",
		"token", "secret", "migration", "payment",
		"billing", "admin",
	}
	normalized := strings.ToLower(filepath.ToSlash(filePath))
	segments := strings.Split(normalized, "/")
	for _, seg := range segments {
		// strip extension from the final segment for matching
		seg = strings.TrimSuffix(seg, filepath.Ext(seg))
		for _, pattern := range sensitivePatterns {
			if seg == pattern {
				return true
			}
		}
	}
	return false
}

func hasEmbeddedConflictMarkers(c *ConflictBlock) bool {
	all := append([]string{}, c.OursLines...)
	all = append(all, c.BaseLines...)
	all = append(all, c.TheirsLines...)
	for _, line := range all {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "<<<<<<<") || strings.HasPrefix(trimmed, "=======") || strings.HasPrefix(trimmed, ">>>>>>>") || strings.HasPrefix(trimmed, "|||||||") {
			return true
		}
	}
	return false
}

func isSourceLikeFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go", ".js", ".jsx", ".ts", ".tsx", ".py", ".java", ".kt", ".rb", ".php", ".rs", ".c", ".cc", ".cpp", ".h", ".hpp", ".cs", ".swift":
		return true
	default:
		return false
	}
}

func hasSemanticResolverCoverage(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".go" || ext == ".js" || ext == ".jsx" || ext == ".ts" || ext == ".tsx"
}

func semanticParserAvailable(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".go" {
		return true
	}
	if ext == ".js" || ext == ".jsx" || ext == ".ts" || ext == ".tsx" {
		_, err := analysis.ParseFile("probe"+ext, []byte("const __gitresolve_probe = 1;"))
		return err == nil
	}
	return false
}
