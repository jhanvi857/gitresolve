package gitresolve

import (
	"fmt"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current conflict state with severity scores",
	Long:  `Displays current unresolved and resolved conflicts sorted by severity. Evaluates and predicts conflict severity without performing auto-resolution.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Conflict Status Check:")
		fmt.Println("  SCORE  TYPE          FILE")
		// TODO: Print from internal/conflict/classifier results
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
