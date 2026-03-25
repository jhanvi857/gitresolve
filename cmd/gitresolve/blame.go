package gitresolve

import (
	"fmt"
	"github.com/spf13/cobra"
)

var blameFileName string

var blameCmd = &cobra.Command{
	Use:   "blame",
	Short: "Show conflict history for the current repo",
	Long:  `Queries the SQLite session log to output the history of conflicts and resolutions in this repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		if blameFileName != "" {
			fmt.Printf("Displaying conflict history for %s...\n", blameFileName)
		} else {
			fmt.Println("Displaying conflict history for the entire repository...")
		}
	},
}

func init() {
	rootCmd.AddCommand(blameCmd)
	blameCmd.Flags().StringVar(&blameFileName, "file", "", "show history for a specific file")
}
