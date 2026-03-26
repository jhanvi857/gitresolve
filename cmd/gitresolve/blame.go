package gitresolve

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var blameFileName string

var blameCmd = &cobra.Command{
	Use:   "blame",
	Short: "Show conflict history for the current repo",
	Long:  `Queries the SQLite session log to output the history of conflicts and resolutions in this repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := openStore(".")
		if err != nil {
			fmt.Println("Blame failed: could not open history DB:", err)
			return
		}
		defer db.Close()

		records, err := db.GetHistory(".")
		if err != nil {
			fmt.Println("Blame failed:", err)
			return
		}

		if len(records) == 0 {
			fmt.Println("No conflict history found yet.")
			return
		}

		fmt.Println("Conflict History:")
		fmt.Println("  TYPE            SEVERITY   STRATEGY   FILE")
		shown := 0
		for _, r := range records {
			if blameFileName != "" && !strings.EqualFold(r.FilePath, blameFileName) {
				continue
			}
			fmt.Printf("  %-14s  %-8s  %-9s  %s\n", r.ConflictType, r.Severity, r.Strategy, r.FilePath)
			shown++
		}

		if shown == 0 {
			fmt.Printf("No history found for file '%s'.\n", blameFileName)
		}
	},
}

func init() {
	rootCmd.AddCommand(blameCmd)
	blameCmd.Flags().StringVar(&blameFileName, "file", "", "show history for a specific file")
}
