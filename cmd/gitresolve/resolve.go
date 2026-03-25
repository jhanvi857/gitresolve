package gitresolve

import (
	"fmt"
	"github.com/spf13/cobra"
)

var resolveFileName string

var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Interactively resolve remaining conflicts",
	Long:  `Fall back to manual or interactive resolution for conflicts that could not be auto-resolved.`,
	Run: func(cmd *cobra.Command, args []string) {
		if resolveFileName != "" {
			fmt.Printf("Resolving conflicts for %s...\n", resolveFileName)
		} else {
			fmt.Println("Resolving all remaining conflicts interactively...")
		}
	},
}

func init() {
	rootCmd.AddCommand(resolveCmd)
	resolveCmd.Flags().StringVar(&resolveFileName, "file", "", "resolve a specific file")
}
