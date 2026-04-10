package gitresolve

import (
	"errors"
	"fmt"
	"os"
	"strings"

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

		autoResolved := 0
		interactiveResolved := 0 // remain 0 for merge command
		validationFailed := 0
		filesUpdated := 0
		var failedFiles []string

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
			var fileAutoResolved int

			for _, c := range conflicts {
				conflict.Classify(c)
				if conflict.ShouldAutoApply(c) {
					resolved := conflict.AutoResolve(c, conflict.Options{
						NoAutoStructured: noAutoStructured,
					})
					if resolved {
						fileAutoResolved++
						autoResolved++
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
			}

			if fileAutoResolved > 0 {
				newContent := conflict.CompileResolution(content, conflicts)
				if strings.HasSuffix(file, ".go") {
					if err := conflict.ValidateGoSyntax(file, newContent); err != nil {
						reason := "reconstructed output failed Go syntax validation"
						fmt.Printf("Escalating %s to manual: %s (%v)\n", file, reason, err)
						validationFailed++
						failedFiles = append(failedFiles, file)
						continue
					}
				}
				if err := conflict.Verify(file, newContent); err != nil {
					fmt.Println("Error: Verification failed:", err)
					validationFailed++
					failedFiles = append(failedFiles, file)
					continue
				}

				if err := writer.Write(file, []byte(newContent)); err != nil {
					if dryRun && errors.Is(err, gserrors.ErrDryRun) {
						fmt.Printf(" > [dry-run] would apply auto-resolution to %s\n", file)
						filesUpdated++
						continue
					}
					fmt.Println("Error: Write failed:", err)
					validationFailed++
					failedFiles = append(failedFiles, file)
					continue
				}

				if fileAutoResolved == len(conflicts) && !dryRun {
					git.MarkResolved(r, file)
				}
				filesUpdated++
			}
		}

		fmt.Printf("\nMerge complete. Summary:\n")
		fmt.Printf("  auto_resolved: %d\n", autoResolved)
		fmt.Printf("  interactive_resolved: %d\n", interactiveResolved)
		fmt.Printf("  validation_failed: %d\n", validationFailed)
		fmt.Printf("  files_updated: %d\n", filesUpdated)

		if validationFailed > 0 {
			fmt.Println("\nFiles with validation failures:")
			for _, f := range failedFiles {
				fmt.Printf("  - %s\n", f)
			}
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would happen without writing")
	mergeCmd.Flags().BoolVar(&dryRun, "dryrun", false, "alias for --dry-run")
	mergeCmd.Flags().BoolVar(&noAutoStructured, "no-auto-structured", false, "disable auto-resolution for structured files (JSON/YAML/TOML)")
}
