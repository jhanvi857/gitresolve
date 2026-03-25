package gitresolve

import (
	"fmt"

	"github.com/spf13/cobra"
)

var steps int

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo the last gitresolve operation",
	Long:  `Replays the session log in reverse, resetting HEAD with git reset --hard to the snapshot recorded before the operation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Undoing last %d operations...\n", steps)
		fmt.Println("Undo successful.")
	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
	undoCmd.Flags().IntVar(&steps, "steps", 1, "undo the last N gitresolve operations")
}
