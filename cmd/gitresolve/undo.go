package gitresolve

import (
	"fmt"

	"github.com/jhanvi857/gitresolve/internal/git"
	"github.com/jhanvi857/gitresolve/internal/safepath"
	"github.com/jhanvi857/gitresolve/pkg/logger"
	"github.com/spf13/cobra"
)

var steps int

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo the last gitresolve operation",
	Long:  `Replays the session log in reverse, resetting HEAD with git reset --hard to the snapshot recorded before the operation.`,
	Run: func(cmd *cobra.Command, args []string) {
		if steps < 1 {
			fmt.Println("Undo failed: --steps must be >= 1")
			return
		}

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

		db, err := openStore(".")
		if err != nil {
			fmt.Println("Undo failed: could not open session DB:", err)
			return
		}
		defer db.Close()

		sessions, err := db.GetRecentSessions(".", steps)
		if err != nil {
			fmt.Println("Undo failed: could not read sessions:", err)
			return
		}

		if len(sessions) < steps {
			fmt.Printf("Undo failed: requested %d step(s) but only %d session(s) available.\n", steps, len(sessions))
			return
		}

		target := sessions[steps-1].SnapshotSHA
		fmt.Printf("Undoing last %d operation(s) -> resetting to %s\n", steps, target)

		if err := r.ResetHardTo(target); err != nil {
			fmt.Println("Undo failed: reset error:", err)
			return
		}

		if err := db.DeleteRecentSessions(".", steps); err != nil {
			fmt.Println("Warning: reset succeeded but could not prune session history:", err)
		}

		if err := git.ClearStoredHead(root); err != nil {
			logger.Debug().Err(err).Msg("failed to clear stored head")
		}
		fmt.Println("Undo successful.")
	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
	undoCmd.Flags().IntVar(&steps, "steps", 1, "undo the last N gitresolve operations")
}
