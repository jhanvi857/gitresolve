package gitresolve

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var targetBranch string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Predict conflicts before push",
	Long:  `Scans the current branch against the target branch (or current tracking branch by default) to predict conflicts before they happen.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Scanning for potential conflicts against %s...\n", targetBranch)

		if _, err := runGit("rev-parse", "--verify", targetBranch); err != nil {
			fmt.Printf("Scan failed: target ref '%s' not found (%v)\n", targetBranch, err)
			return
		}

		baseSHA, err := runGit("merge-base", "HEAD", targetBranch)
		if err != nil {
			fmt.Println("Scan failed: unable to compute merge-base:", err)
			return
		}

		mergeTreeOut, err := runGit("merge-tree", strings.TrimSpace(baseSHA), "HEAD", targetBranch)
		if err != nil {
			fmt.Println("Scan failed: unable to run merge-tree:", err)
			return
		}

		conflictCount := strings.Count(mergeTreeOut, "<<<<<<<")
		if conflictCount == 0 {
			fmt.Println("No potential conflicts found.")
			return
		}

		fmt.Printf("Potential conflicts detected: %d block(s).\n", conflictCount)
		fmt.Println("Conflict hints:")
		for _, line := range strings.Split(mergeTreeOut, "\n") {
			trim := strings.TrimSpace(line)
			if strings.HasPrefix(trim, "changed in both") || strings.HasPrefix(trim, "Auto-merging") {
				fmt.Println(" -", trim)
			}
		}
	},
}

func runGit(args ...string) (string, error) {
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVar(&targetBranch, "target", "main", "scan against a specific branch")
}
