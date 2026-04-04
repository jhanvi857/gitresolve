package conflict

import (
	"fmt"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jhanvi857/gitresolve/internal/analysis"
)

type Options struct {
	NoAutoStructured bool
}

func AutoResolve(c *Conflict, opts Options) bool {
	if !c.CanAutoResolve {
		return false
	}

	switch c.Type {
	case TypeStructured:
		if opts.NoAutoStructured {
			c.ManualReason = "structured auto-merge disabled by --no-auto-structured"
			c.SuggestHint = "try --strategy ours|theirs|both on this file"
			return false
		}
		ext := filepath.Ext(c.FilePath)
		baseBytes := []byte(strings.Join(c.BaseLines, "\n"))
		ourBytes := []byte(strings.Join(c.OurLines, "\n"))
		theirBytes := []byte(strings.Join(c.TheirLines, "\n"))

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
			c.ManualReason = fmt.Sprintf("structured merge parse failed: %v", err)
		} else {
			c.ManualReason = "structured merge has overlapping key/block edits"
		}
		c.SuggestHint = "resolve manually with --strategy ours|theirs|both"

	case TypeWhitespace:
		c.Resolution = strings.Join(c.OurLines, "\n")
		return true

	case TypeImport:
		if analysis.IsCriticalFile(c.FilePath) {
			if unsafe, reason := hasCriticalImportOverlap(c); unsafe {
				c.ManualReason = reason
				c.SuggestHint = "critical file overlap detected; choose ours/theirs manually"
				return false
			}
		}

		merged := mergeImports(c.OurLines, c.TheirLines)
		if !importBlockParses(c.FilePath, merged) {
			c.ManualReason = "merged import block failed syntax parse-check"
			c.SuggestHint = "retry with --strategy ours|theirs|both"
			return false
		}
		c.Resolution = strings.Join(merged, "\n")
		return true

	case TypeIdentical:
		c.Resolution = strings.Join(c.OurLines, "\n")
		return true
	}

	return false
}

func hasCriticalImportOverlap(c *Conflict) (bool, string) {
	if !strings.EqualFold(filepath.Base(c.FilePath), "go.mod") {
		return false, ""
	}

	base := parseGoModRequires(c.BaseLines)
	ours := parseGoModRequires(c.OurLines)
	theirs := parseGoModRequires(c.TheirLines)

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
