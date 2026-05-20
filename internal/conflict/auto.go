package conflict

import (
	"fmt"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jhanvi857/gitresolve/internal/analysis"
)

type Options struct {
	NoAutoStructured bool
}

func AutoResolve(c *ConflictBlock, opts Options) bool {
	if !c.CanAutoResolve {
		return false
	}

	switch c.Type {
	case TypeStructured:
		if opts.NoAutoStructured {
			SetManualEscalation(c, ReasonStructuredAutoDisabled, "structured auto-merge disabled by --no-auto-structured", "try --strategy ours|theirs|both on this file")
			return false
		}
		ext := filepath.Ext(c.FilePath)
		baseBytes := []byte(strings.Join(c.BaseLines, "\n"))
		ourBytes := []byte(strings.Join(c.OursLines, "\n"))
		theirBytes := []byte(strings.Join(c.TheirsLines, "\n"))

		var res analysis.StructuredMergeResult
		var err error

		if ext == ".json" {
			res, err = analysis.MergeJSON(baseBytes, ourBytes, theirBytes)
		} else if ext == ".yaml" || ext == ".yml" {
			res, err = analysis.MergeYAML(baseBytes, ourBytes, theirBytes)
		} else if ext == ".toml" {
			res, err = analysis.MergeTOML(baseBytes, ourBytes, theirBytes)
		}

		if err == nil && !res.HasConflicts {
			c.Resolution = res.Content
			return true
		}
		if err != nil {
			SetManualEscalation(c, ReasonStructuredParseFailed, fmt.Sprintf("structured merge parse failed: %v", err), "resolve manually with --strategy ours|theirs|both")
		} else {
			SetManualEscalation(c, ReasonStructuredOverlap, "structured merge has overlapping key/block edits", "resolve manually with --strategy ours|theirs|both")
		}

	case TypeWhitespace:
		c.Resolution = strings.Join(c.OursLines, "\n")
		// Validate the resolution before returning
		if err := Verify(c.FilePath, c.Resolution); err != nil {
			SetManualEscalation(c, ReasonValidationSyntaxFailed, fmt.Sprintf("auto-resolved content failed validation: %v", err), "manually review and apply resolution strategy")
			return false
		}
		return true

	case TypeImport:
		if analysis.IsCriticalFile(c.FilePath) {
			if unsafe, reason := hasCriticalImportOverlap(c); unsafe {
				SetManualEscalation(c, ReasonImportOverlapCritical, reason, "critical file overlap detected; choose ours/theirs manually")
				return false
			}
		}

		var merged []string
		if strings.EqualFold(filepath.Ext(c.FilePath), ".go") {
			merged = mergeGoImports(c.OursLines, c.TheirsLines)
		} else {
			merged = mergeImports(c.OursLines, c.TheirsLines)
		}
		if len(merged) == 0 {
			SetManualEscalation(c, ReasonImportMergeFailed, "could not produce a valid merged import block", "retry with --strategy ours|theirs|both")
			return false
		}
		if !importBlockParses(c.FilePath, merged) {
			SetManualEscalation(c, ReasonImportParseFailed, "merged import block failed syntax parse-check", "retry with --strategy ours|theirs|both")
			return false
		}
		c.Resolution = strings.Join(merged, "\n")
		// Validate the resolution before returning
		if err := Verify(c.FilePath, c.Resolution); err != nil {
			SetManualEscalation(c, ReasonValidationSyntaxFailed, fmt.Sprintf("auto-resolved content failed validation: %v", err), "manually review and apply resolution strategy")
			return false
		}
		return true

	case TypeIdentical:
		c.Resolution = strings.Join(c.OursLines, "\n")
		// Validate the resolution before returning
		if err := Verify(c.FilePath, c.Resolution); err != nil {
			SetManualEscalation(c, ReasonValidationSyntaxFailed, fmt.Sprintf("auto-resolved content failed validation: %v", err), "manually review and apply resolution strategy")
			return false
		}
		return true
	}

	return false
}

func hasCriticalImportOverlap(c *ConflictBlock) (bool, string) {
	if !strings.EqualFold(filepath.Base(c.FilePath), "go.mod") {
		return false, ""
	}

	base := parseGoModRequires(c.BaseLines)
	ours := parseGoModRequires(c.OursLines)
	theirs := parseGoModRequires(c.TheirsLines)

	keys := make(map[string]struct{})
	for k := range ours {
		keys[k] = struct{}{}
	}
	for k := range theirs {
		keys[k] = struct{}{}
	}

	for mod := range keys {
		baseV := base[mod]
		ourV, ourTouched := ours[mod]
		theirV, theirTouched := theirs[mod]
		if !ourTouched || !theirTouched {
			continue
		}
		oursChanged := ourV != baseV
		theirsChanged := theirV != baseV
		if oursChanged && theirsChanged && ourV != theirV {
			return true, fmt.Sprintf("critical go.mod entry %q changed on both sides (%s vs %s)", mod, ourV, theirV)
		}
	}

	return false, ""
}

var goModRequireLine = regexp.MustCompile(`^([A-Za-z0-9_\-\./]+)\s+v\S+`)

func parseGoModRequires(lines []string) map[string]string {
	result := make(map[string]string)
	inRequireBlock := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}
		if strings.HasPrefix(trimmed, "require (") {
			inRequireBlock = true
			continue
		}
		if inRequireBlock && trimmed == ")" {
			inRequireBlock = false
			continue
		}

		candidate := trimmed
		if strings.HasPrefix(trimmed, "require ") {
			candidate = strings.TrimSpace(strings.TrimPrefix(trimmed, "require"))
		}

		if !inRequireBlock && !strings.HasPrefix(trimmed, "require ") {
			continue
		}

		m := goModRequireLine.FindStringSubmatch(candidate)
		if len(m) < 2 {
			continue
		}
		parts := strings.Fields(candidate)
		if len(parts) >= 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

func importBlockParses(filePath string, merged []string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	block := strings.Join(merged, "\n")

	if ext == ".go" {
		snippet := "package p\n" + block + "\nvar _ = 1\n"
		_, err := parser.ParseFile(token.NewFileSet(), "snippet.go", snippet, parser.AllErrors)
		return err == nil
	}

	if ext == ".js" || ext == ".jsx" || ext == ".ts" || ext == ".tsx" {
		snippet := block + "\nconst __gitresolve_probe = 1;\n"
		_, err := analysis.ParseFile("probe"+ext, []byte(snippet))
		if err != nil {
			if strings.Contains(err.Error(), "requires cgo-enabled build") {
				return true
			}
			return false
		}
	}

	return true
}

func mergeGoImports(ours, theirs []string) []string {
	specs := make(map[string]struct{})
	malformed := false

	extract := func(lines []string) {
		inBlock := false
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "//") {
				continue
			}

			if strings.HasPrefix(trimmed, "import ") {
				rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "import"))
				if rest == "(" {
					inBlock = true
					continue
				}
				if rest == ")" {
					malformed = true
					continue
				}
				if spec := normalizeGoImportSpec(rest); spec != "" {
					specs[spec] = struct{}{}
				} else {
					malformed = true
				}
				continue
			}

			if inBlock {
				if trimmed == ")" {
					inBlock = false
					continue
				}
				if spec := normalizeGoImportSpec(trimmed); spec != "" {
					specs[spec] = struct{}{}
				} else {
					malformed = true
				}
				continue
			}

			if strings.Contains(trimmed, "\"") {
				// Quoted import-like spec outside import context indicates malformed fragment.
				malformed = true
			}
		}
	}

	extract(ours)
	extract(theirs)

	if malformed || len(specs) == 0 {
		return nil
	}

	ordered := make([]string, 0, len(specs))
	for s := range specs {
		ordered = append(ordered, s)
	}
	sort.Strings(ordered)

	result := []string{"import ("}
	for _, spec := range ordered {
		result = append(result, "\t"+spec)
	}
	result = append(result, ")")
	return result
}

func normalizeGoImportSpec(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if idx := strings.Index(raw, "//"); idx >= 0 {
		raw = strings.TrimSpace(raw[:idx])
	}
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "(" || raw == ")" {
		return ""
	}

	firstQuote := strings.Index(raw, "\"")
	lastQuote := strings.LastIndex(raw, "\"")
	if firstQuote < 0 || lastQuote <= firstQuote {
		return ""
	}

	alias := strings.TrimSpace(raw[:firstQuote])
	path := strings.TrimSpace(raw[firstQuote : lastQuote+1])
	if alias == "" {
		return path
	}
	return alias + " " + path
}

func mergeImports(ours, theirs []string) []string {
	seen := make(map[string]bool)
	var merged []string
	var hasOpenParen bool
	var hasCloseParen bool

	normalize := func(line string) string {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "import ")
		line = strings.Trim(line, "()\"' ")
		return line
	}

	process := func(lines []string) {
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			if trimmed == "(" || strings.HasSuffix(trimmed, "(") {
				hasOpenParen = true
				continue
			}
			if trimmed == ")" {
				hasCloseParen = true
				continue
			}

			norm := normalize(line)
			if !seen[norm] {
				seen[norm] = true
				merged = append(merged, line)
			}
		}
	}

	process(ours)
	process(theirs)

	var result []string
	if hasOpenParen {
		prefix := "("
		for _, line := range append(ours, theirs...) {
			if strings.Contains(line, "(") && !strings.HasPrefix(strings.TrimSpace(line), "(") {
				prefix = line
				break
			}
		}
		result = append(result, prefix)
	}

	result = append(result, merged...)

	if hasCloseParen {
		result = append(result, ")")
	}

	return result
}
