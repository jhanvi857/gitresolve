package conflict

import (
	"bufio"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jhanvi857/gitresolve/internal/analysis"
	"github.com/jhanvi857/gitresolve/pkg/logger"
)

type Strategy int

type ResolveOptions struct {
	NonInteractive bool
	Timeout        time.Duration
}

type ResolveResult struct {
	Applied         bool
	Selected        Strategy
	SelectedLabel   string
	FailureHint     string
	BothAllowedNext bool
}

const (
	StrategyOurs Strategy = iota
	StrategyTheirs
	StrategyBoth
	StrategyInteractive
)

func Resolve(c *ConflictBlock, strategy Strategy, opts ResolveOptions) (ResolveResult, error) {
	applySelection := func(strat Strategy) (ResolveResult, error) {
		res := ResolveResult{Applied: true, Selected: strat, BothAllowedNext: true}

		switch strat {
		case StrategyOurs:
			res.SelectedLabel = "ours"
			c.Resolution = strings.Join(c.OursLines, "\n")
			return res, nil
		case StrategyTheirs:
			res.SelectedLabel = "theirs"
			c.Resolution = strings.Join(c.TheirsLines, "\n")
			return res, nil
		case StrategyBoth:
			res.SelectedLabel = "both"
			if bothUnsafeReason := riskReasonForBoth(c); bothUnsafeReason != "" {
				SetManualEscalation(c, ReasonStrategyBothBlockedRisk, bothUnsafeReason, "choose ours/theirs or edit manually")
				res.BothAllowedNext = false
				res.FailureHint = bothUnsafeReason
				return res, errors.New(bothUnsafeReason)
			}
			merged, err := buildBothResolution(c)
			if err != nil {
				res.BothAllowedNext = false
				res.FailureHint = err.Error()
				return res, err
			}
			c.Resolution = merged
			return res, nil
		default:
			return ResolveResult{}, fmt.Errorf("unknown selection strategy")
		}
	}

	switch strategy {
	case StrategyOurs, StrategyTheirs, StrategyBoth:
		res, err := applySelection(strategy)
		return res, err
	case StrategyInteractive:
		if opts.NonInteractive {
			return ResolveResult{}, fmt.Errorf("conflict in %s requires manual resolution, but --non-interactive is set", c.FilePath)
		}

		printConflictForChoice(c)
		if c.Confidence > 0 {
			fmt.Printf("Suggested option: %s (confidence %.2f)\n", suggestedChoice(c), c.Confidence)
		}
		if c.ManualReason != "" {
			fmt.Printf("Reason: %s\n", c.ManualReason)
		}
		if c.SuggestHint != "" {
			fmt.Printf("Hint: %s\n", c.SuggestHint)
		}

		bothAllowed := true
		attempts := 0
		const maxAttempts = 8

		for attempts < maxAttempts {
			attempts++
			choice, timedOut := readChoice(opts.Timeout, bothAllowed)
			if timedOut {
				fmt.Printf("\nTimeout reached (%s). Auto-selecting [T]heirs.\n", opts.Timeout.String())
				res, err := applySelection(StrategyTheirs)
				return res, err
			}

			switch choice {
			case "O":
				return applySelection(StrategyOurs)
			case "T":
				return applySelection(StrategyTheirs)
			case "B":
				if !bothAllowed {
					fmt.Println("[B]oth disabled for this conflict after previous validation failure.")
					continue
				}
				res, err := applySelection(StrategyBoth)
				if err != nil {
					fmt.Printf("BOTH failed: %v\n", err)
					fmt.Println("Choose [O], [T], [M], or [S].")
					bothAllowed = false
					continue
				}
				return res, nil
			case "M":
				return ResolveResult{Applied: false, SelectedLabel: "manual", FailureHint: "manual edit requested", BothAllowedNext: bothAllowed}, nil
			case "S":
				return ResolveResult{Applied: false, SelectedLabel: "skip", FailureHint: "skipped by user", BothAllowedNext: bothAllowed}, nil
			default:
				fmt.Println("Invalid option. Choose [O]urs, [T]heirs, [B]oth, [M]anual edit, or [S]kip.")
			}
		}

		fmt.Println("Too many invalid attempts. Auto-selecting [T]heirs to avoid infinite retry.")
		res, err := applySelection(StrategyTheirs)
		return res, err
	default:
		return ResolveResult{}, fmt.Errorf("Resolve: unknown strategy %d", strategy)
	}
}

func printConflictForChoice(c *ConflictBlock) {
	if c.Type == TypeScalar {
		fmt.Printf("\n[Scalar] %s (L%d-%d)\n", c.FilePath, c.StartLine, c.EndLine)
		fmt.Printf(" [O]urs:   %s\n", strings.Join(c.OursLines, " "))
		fmt.Printf(" [T]heirs: %s\n", strings.Join(c.TheirsLines, " "))
		fmt.Println(" Options: [O]urs [T]heirs [B]oth [M]anual edit [S]kip")
		return
	}

	fmt.Printf("\n--- Conflict in %s ---\n", c.FilePath)
	fmt.Println("<<<<<<< OURS")
	fmt.Println(strings.Join(c.OursLines, "\n"))
	fmt.Println("=======")
	fmt.Println(strings.Join(c.TheirsLines, "\n"))
	fmt.Println(">>>>>>> THEIRS")
	fmt.Println("Options: [O]urs [T]heirs [B]oth [M]anual edit [S]kip")
}

func readChoice(timeout time.Duration, bothAllowed bool) (string, bool) {
	prompt := "Select [O/T/B/M/S]: "
	if !bothAllowed {
		prompt = "Select [O/T/M/S]: "
	}

	inputChan := make(chan string, 1)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		inputChan <- input
	}()

	if timeout > 0 {
		select {
		case input := <-inputChan:
			return normalizeChoice(input), false
		case <-time.After(timeout):
			return "", true
		}
	}

	input := <-inputChan
	return normalizeChoice(input), false
}

func normalizeChoice(input string) string {
	u := strings.TrimSpace(strings.ToUpper(input))
	switch u {
	case "O", "OURS":
		return "O"
	case "T", "THEIRS":
		return "T"
	case "B", "BOTH":
		return "B"
	case "M", "MANUAL":
		return "M"
	case "S", "SKIP":
		return "S"
	default:
		return ""
	}
}

func suggestedChoice(c *ConflictBlock) string {
	switch c.Type {
	case TypeWhitespace, TypeIdentical:
		return "[O]urs"
	case TypeImport, TypeStructured:
		return "[B]oth"
	case TypeDeleteModify:
		return "[T]heirs"
	default:
		if c.Confidence >= 0.75 {
			return "[B]oth"
		}
		return "[T]heirs"
	}
}

func buildBothResolution(c *ConflictBlock) (string, error) {
	ext := strings.ToLower(filepath.Ext(c.FilePath))
	if ext == ".yaml" || ext == ".yml" {
		merged, err := mergeYAMLBoth(c)
		if err == nil {
			return merged, nil
		}
		logger.Debug("yaml both merge fallback: " + err.Error())
	}

	oursDepth := braceDepth(c.OursLines)
	theirsDepth := braceDepth(c.TheirsLines)
	if oursDepth > 0 || theirsDepth > 0 {
		combined := make([]string, 0, len(c.OursLines)+len(c.TheirsLines)+3)
		combined = append(combined, c.OursLines...)
		if oursDepth > 0 {
			combined = append(combined, "}")
		}
		combined = append(combined, "")
		combined = append(combined, c.TheirsLines...)
		if theirsDepth > 0 {
			combined = append(combined, "}")
		}
		return strings.Join(combined, "\n"), nil
	}

	combined := combineWithIndent(c.OursLines, c.TheirsLines)
	if ext == ".go" {
		if err := validateGoFragment(combined); err != nil {
			return "", fmt.Errorf("invalid Go construct from BOTH: %w", err)
		}
	}

	return strings.Join(combined, "\n"), nil
}

func combineWithIndent(ours, theirs []string) []string {
	combined := make([]string, 0, len(ours)+len(theirs)+1)
	combined = append(combined, ours...)
	if len(ours) > 0 && len(theirs) > 0 && strings.TrimSpace(ours[len(ours)-1]) != "" && strings.TrimSpace(theirs[0]) != "" {
		combined = append(combined, "")
	}
	combined = append(combined, theirs...)
	return combined
}

func mergeYAMLBoth(c *ConflictBlock) (string, error) {
	res, err := analysis.MergeYAML([]byte(strings.Join(c.BaseLines, "\n")), []byte(strings.Join(c.OursLines, "\n")), []byte(strings.Join(c.TheirsLines, "\n")))
	if err != nil {
		return "", err
	}
	if res.HasConflicts {
		return "", fmt.Errorf("YAML map merge has overlapping key edits")
	}
	return strings.TrimRight(res.Content, "\n"), nil
}

func validateGoFragment(lines []string) error {
	content := strings.Join(lines, "\n")
	if strings.Count(content, "{") != strings.Count(content, "}") {
		return fmt.Errorf("unbalanced braces in merged fragment")
	}

	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, "merged_fragment.go", "package p\n"+content+"\n", parser.AllErrors); err == nil {
		return nil
	}

	wrapped := "package p\nfunc _gitresolve_probe() {\n" + content + "\n}\n"
	if _, err := parser.ParseFile(fset, "merged_wrapped.go", wrapped, parser.AllErrors); err == nil {
		return nil
	} else {
		return err
	}
}

func riskReasonForBoth(c *ConflictBlock) string {
	if !isCodeFileForBothRisk(c.FilePath) {
		return ""
	}
	if c.Type == TypeDeleteModify {
		return "BOTH disabled for delete-vs-modify conflicts on source code; choose ours/theirs/manual"
	}
	if c.Type == TypeSignature {
		return "BOTH disabled for function signature conflicts on source code; choose ours/theirs/manual"
	}
	if c.Type == TypeLogic && c.Severity >= SeverityHigh {
		return "BOTH disabled for high-risk semantic logic conflicts; choose ours/theirs/manual"
	}
	return ""
}

func isCodeFileForBothRisk(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go", ".js", ".jsx", ".ts", ".tsx", ".py", ".java", ".kt", ".rb", ".php", ".rs", ".c", ".cc", ".cpp", ".h", ".hpp", ".cs", ".swift":
		return true
	default:
		return false
	}
}
