package gitresolve

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jhanvi857/gitresolve/internal/safepath"
	"github.com/jhanvi857/gitresolve/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var verbose bool
var force bool
var logLevel string

var rootCmd = &cobra.Command{
	Use:   "gitresolve",
	Short: "A fully local, privacy-first Git conflict resolution engine",
	Long: `gitresolve classifies merge conflicts using Abstract Syntax Tree analysis 
and deterministic rule-based reasoning, auto-resolves safe conflict types, 
detects cross-file semantic breakages after merge, and predicts conflicts 
before they happen.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		lvl := zerolog.WarnLevel
		if verbose {
			lvl = zerolog.InfoLevel
		} else {
			var err error
			lvl, err = zerolog.ParseLevel(logLevel)
			if err != nil {
				return fmt.Errorf("invalid log level: %w", err)
			}
		}

		logger.InitWithLevel(lvl)
		safepath.SetForceAllowUnsupported(force)
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
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "warn", "set log level (error, warn, info, debug, trace)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable info logging (shorthand for --log-level info)")
	rootCmd.PersistentFlags().BoolVar(&force, "force", false, "allow unsafe fallback on unsupported platforms (plan9/js)")
}
