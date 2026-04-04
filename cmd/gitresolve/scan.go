package gitresolve

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var targetBranch string

// safeBranchNameRe allows only safe characters in branch names.
var safeBranchNameRe = regexp.MustCompile(`^[a-zA-Z0-9_./-]+$`)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Predict conflicts before push",
	Long:  `Scans the current branch against the target branch (or current tracking branch by default) to predict conflicts before they happen.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !safeBranchNameRe.MatchString(targetBranch) {
			fmt.Printf("Scan failed: invalid branch name %q (only alphanumeric, '.', '/', '-', '_' allowed)\n", targetBranch)
			return
		}

		fmt.Printf("Scanning for potential conflicts against %s...\n", targetBranch)

		// Use "--" to prevent flag injection via branch names
		if _, err := runGit("rev-parse", "--verify", "--", targetBranch); err != nil {
			fmt.Printf("Scan failed: target ref '%s' not found (%v)\n", targetBranch, err)
			return
		}

		// Use modern merge-tree which is more accurate and handles renames better
		mergeTreeOut, err := runGit("merge-tree", "HEAD", targetBranch)
		if err != nil {
			fmt.Println("Scan failed: unable to run merge-tree:", err)
			return
		}

		// merge-tree output:
		// <tree-hash>
		// Conflict ...
		// <other info>
		lines := strings.Split(mergeTreeOut, "\n")
		var conflicts []string
		for _, line := range lines {
			if strings.HasPrefix(line, "Conflict") || strings.Contains(line, "CONFLICT") {
				conflicts = append(conflicts, line)
			}
		}

		if len(conflicts) == 0 {
			fmt.Println("No potential conflicts found.")
			fmt.Println("\nNote: scan may miss conflicts in binary files and low-similarity renames.")
			fmt.Println("Always run `gitresolve status` after merging.")
			return
		}

		fmt.Printf("Potential conflicts detected: %d files/blocks.\n", len(conflicts))
		fmt.Println("Conflict hints:")
		for _, c := range conflicts {
			fmt.Println(" -", c)
		}

		fmt.Println("\nNote: scan may miss conflicts in binary files and low-similarity renames.")
		fmt.Println("Always run `gitresolve status` after merging.")
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
