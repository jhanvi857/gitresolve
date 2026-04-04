package gitresolve

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jhanvi857/gitresolve/internal/conflict"
	"github.com/jhanvi857/gitresolve/internal/git"
	"github.com/jhanvi857/gitresolve/internal/safety"
	"github.com/jhanvi857/gitresolve/internal/store"
	gserrors "github.com/jhanvi857/gitresolve/pkg/errors"
	"github.com/spf13/cobra"
)

var resolveFileName string
var resolveStrategy string
var resolveDryRun bool
var resolveNonInteractive bool
var resolveTimeout time.Duration

var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Interactively resolve remaining conflicts",
	Long:  `Fall back to manual or interactive resolution for conflicts that could not be auto-resolved.`,
	Run: func(cmd *cobra.Command, args []string) {
		strategy, err := parseStrategy(resolveStrategy)
		if err != nil {
			fmt.Println("Resolve failed:", err)
			return
		}

		r, err := git.Open(".")
		if err != nil {
			fmt.Println("Fatal: Failed to open git repository:", err)
			return
		}
		defer git.Close(r)
		HandleSignals(r)

		files, err := git.ConflictedFiles(r)
		if err != nil {
			fmt.Println("Resolve check:", err)
			return
		}

		writer := safety.NewWriter(resolveDryRun)

		repoPath := "."
		db, dbErr := openStore(repoPath)
		if dbErr == nil {
			defer db.Close()
		}

		if !resolveDryRun && dbErr == nil {
			head, headErr := r.HeadCommit()
			if headErr == nil {
				_ = db.SaveSession(repoPath, "resolve", head)
				_ = git.StoreHead(repoPath, head)
			}
		}

		resolvedFiles := 0
		autoResolved := 0
		manualRequired := 0
		criticalConflicts := 0
		hadHardFailure := false

		for _, file := range files {
			if resolveFileName != "" && file != resolveFileName {
				continue
			}

			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading %s: %v\n", file, err)
				continue
			}

			if !hasConflictMarkers(string(content)) {
				continue
			}

			conflicts := conflict.ParseFile(file, content)
			if len(conflicts) == 0 {
				continue
			}

			if !resolveDryRun {
				if err := safety.PreserveOriginal(file); err != nil {
					fmt.Printf("Warning: backup failed for %s: %v\n", file, err)
					continue
				}
			}

			for _, c := range conflicts {
				conflict.Classify(c)
				if conflict.ShouldAutoApply(c) {
					autoResolved++
				} else {
					manualRequired++
					if conflict.NeedsGuidedChoice(c) {
						fmt.Printf("Guided choice: %s L%d-%d [confidence=%.2f]. Suggested strategy: ours|theirs|both\n", file, c.StartLine, c.EndLine, c.Confidence)
					}
					if c.ManualReason != "" {
						fmt.Printf("  reason: %s\n", c.ManualReason)
					}
					if c.SuggestHint != "" {
						fmt.Printf("  hint: %s\n", c.SuggestHint)
					}
					if c.Severity == conflict.SeverityCritical {
						criticalConflicts++
					}
				}

				opts := conflict.ResolveOptions{
					NonInteractive: resolveNonInteractive,
					Timeout:        resolveTimeout,
				}
				if err := conflict.Resolve(c, strategy, opts); err != nil {
					fmt.Printf("Resolve failed for %s: %v\n", file, err)
					if resolveNonInteractive {
						os.Exit(1)
					}
					continue
				}
			}

			newContent := conflict.CompileResolution(content, conflicts)
			if err := conflict.Verify(file, newContent); err != nil {
				fmt.Printf("Verification failed for %s: %v\n", file, err)
				hadHardFailure = true
				if resolveNonInteractive {
					os.Exit(1)
				}
				continue
			}

			if err := writer.Write(file, []byte(newContent)); err != nil {
				if resolveDryRun && errors.Is(err, gserrors.ErrDryRun) {
					fmt.Printf("[dry-run] would resolve %s using strategy '%s'\n", file, resolveStrategy)
					resolvedFiles++
					continue
				}
				fmt.Printf("Error writing %s: %v\n", file, err)
				hadHardFailure = true
				if resolveNonInteractive {
					os.Exit(1)
				}
				continue
			}

			if !resolveDryRun {
				if err := git.MarkResolved(r, file); err != nil {
					fmt.Printf("Warning: failed to stage %s: %v\n", file, err)
				}
			}

			if dbErr == nil {
				for _, c := range conflicts {
					_ = db.SaveConflict(storeConflict(repoPath, file, c, resolveStrategy))
				}
			}

			fmt.Printf("Resolved %s using strategy '%s'\n", file, resolveStrategy)
			resolvedFiles++
		}

		if resolveFileName != "" && resolvedFiles == 0 {
			fmt.Printf("No conflicted file matched '%s'.\n", resolveFileName)
		}

		if resolveDryRun {
			fmt.Printf("\nDry-run complete. Files that would be resolved: %d\n", resolvedFiles)
		} else {
			fmt.Printf("\nResolve complete. Summary:\n")
			fmt.Printf("  %d blocks auto-resolved\n", autoResolved)
			fmt.Printf("  %d blocks required manual input\n", manualRequired)
			if criticalConflicts > 0 {
				fmt.Printf("  %d critical conflicts handled\n", criticalConflicts)
			}
			fmt.Printf("Total files updated: %d\n", resolvedFiles)
		}

		if hadHardFailure {
			os.Exit(1)
		}
	},
}

func parseStrategy(v string) (conflict.Strategy, error) {
	switch v {
	case "ours":
		return conflict.StrategyOurs, nil
	case "theirs":
		return conflict.StrategyTheirs, nil
	case "both":
		return conflict.StrategyBoth, nil
	case "interactive":
		return conflict.StrategyInteractive, nil
	default:
		return 0, fmt.Errorf("unknown strategy '%s' (use interactive|ours|theirs|both)", v)
	}
}

func storeConflict(repoPath, file string, c *conflict.Conflict, strategy string) store.ConflictRecord {
	return store.ConflictRecord{
		RepoPath:     repoPath,
		FilePath:     file,
		ConflictType: typeLabel(c.Type),
		Severity:     severityLabel(c.Severity),
		Strategy:     strategy,
	}
}

func init() {
	rootCmd.AddCommand(resolveCmd)
	resolveCmd.Flags().StringVar(&resolveFileName, "file", "", "resolve a specific file")
	resolveCmd.Flags().StringVar(&resolveStrategy, "strategy", "interactive", "resolve strategy: interactive|ours|theirs|both")
	resolveCmd.Flags().BoolVar(&resolveDryRun, "dry-run", false, "show what would happen without writing")
	resolveCmd.Flags().BoolVar(&resolveNonInteractive, "non-interactive", false, "fail on conflicts requiring manual resolution instead of prompting")
	resolveCmd.Flags().DurationVar(&resolveTimeout, "timeout", 0, "timeout for interactive prompt (e.g. 30s). Auto-selects theirs if reached.")
}
