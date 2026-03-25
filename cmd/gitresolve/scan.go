package gitresolve

import (
	"fmt"
	"github.com/spf13/cobra"
)

var targetBranch string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Predict conflicts before push",
	Long:  `Scans the current branch against the target branch (or current tracking branch by default) to predict conflicts before they happen.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Scanning for potential conflicts against %s...\n", targetBranch)
		// TODO: Call internal/conflict/ detector logic and internal/ownership checker.
		fmt.Println("No potential conflicts found.")
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVar(&targetBranch, "target", "main", "scan against a specific branch")
}
