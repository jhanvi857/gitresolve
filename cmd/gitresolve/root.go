package gitresolve

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jhanvi857/gitresolve/pkg/logger"
	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "gitresolve",
	Short: "A fully local, privacy-first Git conflict resolution engine",
	Long: `gitresolve classifies merge conflicts using Abstract Syntax Tree analysis 
and deterministic rule-based reasoning, auto-resolves safe conflict types, 
detects cross-file semantic breakages after merge, and predicts conflicts 
before they happen.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger.Init(verbose)
		if err := preflightChecks(); err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// preflightChecks verifies required external dependencies are available.
func preflightChecks() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("preflight: git is not installed or not in PATH")
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable debug logging")
}
