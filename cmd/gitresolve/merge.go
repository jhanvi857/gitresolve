package gitresolve

import (
	"errors"
	"fmt"
	"os"

	"github.com/jhanvi857/gitresolve/internal/conflict"
	"github.com/jhanvi857/gitresolve/internal/git"
	"github.com/jhanvi857/gitresolve/internal/safety"
	"github.com/jhanvi857/gitresolve/internal/store"
	gserrors "github.com/jhanvi857/gitresolve/pkg/errors"
	"github.com/spf13/cobra"
)

var dryRun bool
var noAutoStructured bool

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Run smart merge on current conflicts",
	Long:  `Analyzes and auto-resolves smart merge conflicts using deterministic rule-based algorithms. Escapes complex semantic or structural discrepancies to manual review securely.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Engine Bootup: Initializing gitresolve in directory '.' (DryRun: %v, NoAutoStructured: %v)\n", dryRun, noAutoStructured)

		repoPath := "."
		r, err := git.Open(".")
		if err != nil {
			fmt.Println("Fatal: Failed to open git repository: ", err)
			return
		}
		defer git.Close(r)
		HandleSignals(r)

		db, dbErr := openStore(repoPath)
		if dbErr == nil {
			defer db.Close()
		}

		if !dryRun && dbErr == nil {
			head, headErr := r.HeadCommit()
			if headErr == nil {
				_ = db.SaveSession(repoPath, "merge", head)
				_ = git.StoreHead(repoPath, head)
			}
		}

		files, err := git.ConflictedFiles(r)
		if err != nil {
			fmt.Println("Status check:", err)
			return
		}

		fmt.Printf("Scanning index. Found %d unmerged conflicts...\n", len(files))
		writer := safety.NewWriter(dryRun)
		hadHardFailure := false

		for _, file := range files {
			fmt.Printf("\n--- Processing %s ---\n", file)

			if !dryRun {
				if err := safety.PreserveOriginal(file); err != nil {
					fmt.Println("Warning: Could not create backup:", err)
					continue
				}
			}

			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}

			conflicts := conflict.ParseFile(file, content)
			var autoResolvedCount int

			for _, c := range conflicts {
				conflict.Classify(c)
				if conflict.ShouldAutoApply(c) {
					resolved := conflict.AutoResolve(c, conflict.Options{
						NoAutoStructured: noAutoStructured,
					})
					if resolved {
						autoResolvedCount++
						if dbErr == nil {
							_ = db.SaveConflict(store.ConflictRecord{
								RepoPath:     repoPath,
								FilePath:     file,
								ConflictType: typeLabel(c.Type),
								Severity:     severityLabel(c.Severity),
								Strategy:     "auto",
							})
						}
					} else {
						fmt.Printf(" > Escalating conflict [Severity %d] %v\n", c.Severity, c.Type)
						if c.ManualReason != "" {
							fmt.Printf("   reason: %s\n", c.ManualReason)
						}
						if c.SuggestHint != "" {
							fmt.Printf("   hint: %s\n", c.SuggestHint)
						}
						if dbErr == nil {
							_ = db.SaveConflict(store.ConflictRecord{
								RepoPath:     repoPath,
								FilePath:     file,
								ConflictType: typeLabel(c.Type),
								Severity:     severityLabel(c.Severity),
								Strategy:     "manual-required",
							})
						}
					}
				} else {
					if conflict.NeedsGuidedChoice(c) {
						fmt.Printf(" > Guided choice needed [severity=%d confidence=%.2f type=%v] in %s. Suggested strategies: ours/theirs/both.\n", c.Severity, c.Confidence, c.Type, file)
					} else {
						fmt.Printf(" > Escalating conflict [Severity %d] %v\n", c.Severity, c.Type)
					}
					if c.ManualReason != "" {
						fmt.Printf("   reason: %s\n", c.ManualReason)
					}
					if c.SuggestHint != "" {
						fmt.Printf("   hint: %s\n", c.SuggestHint)
					}
					if dbErr == nil {
						_ = db.SaveConflict(store.ConflictRecord{
							RepoPath:     repoPath,
							FilePath:     file,
							ConflictType: typeLabel(c.Type),
							Severity:     severityLabel(c.Severity),
							Strategy:     "manual-required",
						})
					}
				}
			}

			if autoResolvedCount > 0 {
				newContent := conflict.CompileResolution(content, conflicts)
				if err := conflict.Verify(file, newContent); err != nil {
					fmt.Println("Error: Verification failed (hard stop for file):", err)
					hadHardFailure = true
					continue
				}

				err := writer.Write(file, []byte(newContent))
				if err != nil {
					if dryRun && errors.Is(err, gserrors.ErrDryRun) {
						fmt.Printf(" > [dry-run] would apply auto-resolution to %s (%d/%d blocks).\n", file, autoResolvedCount, len(conflicts))
					} else {
						fmt.Println("Error: Atomic write failed:", err)
						hadHardFailure = true
						continue
					}
				}

				if autoResolvedCount == len(conflicts) && !dryRun {
					git.MarkResolved(r, file)
					fmt.Printf(" > Successfully auto-resolved 100%% of conflicts in %s and staged.\n", file)
				} else {
					fmt.Printf(" > Auto-resolved %d of %d conflicts in %s. Manual review still required for remainder.\n", autoResolvedCount, len(conflicts), file)
				}
			} else {
				fmt.Printf(" > No safe resolutions could be applied to %s.\n", file)
			}
		}

		fmt.Println("\nMerge scan complete.")
		if hadHardFailure {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would happen without writing")
	mergeCmd.Flags().BoolVar(&noAutoStructured, "no-auto-structured", false, "disable auto-resolution for structured files (JSON/YAML/TOML)")
}
