package gitresolve

import (
	"fmt"
	"io"

	"github.com/jhanvi857/gitresolve/internal/conflict"
	"github.com/jhanvi857/gitresolve/internal/git"
	"github.com/jhanvi857/gitresolve/internal/safepath"
	"github.com/jhanvi857/gitresolve/pkg/logger"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current conflict state with severity scores",
	Long:  `Displays current unresolved and resolved conflicts sorted by severity. Evaluates and predicts conflict severity without performing auto-resolution.`,
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot, err := ResolveRepoRoot()
		if err != nil {
			fmt.Println("Fatal: failed to resolve repository root:", err)
			return
		}

		root, err := safepath.RepoRoot(repoRoot)
		if err != nil {
			fmt.Println("Fatal: failed to open repository sandbox:", err)
			return
		}
		defer root.Close()

		r, err := git.Open(".", root)
		if err != nil {
			fmt.Println("Fatal: Failed to open git repository:", err)
			return
		}
		defer func() {
			if rec := recover(); rec != nil {
				_ = git.Close(r)
				panic(rec)
			}
		}()
		defer git.Close(r)

		files, err := git.ConflictedFiles(r)
		if err != nil {
			fmt.Println("Status check:", err)
			return
		}

		fmt.Println("Conflict Status Check:")
		fmt.Println("  SCORE  TYPE            AUTO  FILE")

		var total int
		for _, file := range files {
			f, err := safepath.SafeOpen(root, file)
			if err != nil {
				fmt.Printf("  --     read-error      --    %s (%v)\n", file, err)
				continue
			}

			content, err := io.ReadAll(f)
			if closeErr := f.Close(); closeErr != nil {
				logger.Debug().Err(closeErr).Str("file", file).Msg("failed to close file")
			}
			if err != nil {
				fmt.Printf("  --     read-error      --    %s (%v)\n", file, err)
				continue
			}

			parsed := conflict.ParseFile(file, content)
			if len(parsed) == 0 {
				fmt.Printf("  --     unparsed         --    %s\n", file)
				continue
			}

			for _, c := range parsed {
				conflict.Classify(c)
				auto := "no"
				if c.CanAutoResolve {
					auto = "yes"
				}
				fmt.Printf("  %-5d  %-14s %-4s  %s\n", c.Severity, typeLabel(c.Type), auto, file)
				total++
			}
		}

		fmt.Printf("\nTotal conflict blocks: %d\n", total)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
