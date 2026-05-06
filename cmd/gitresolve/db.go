package gitresolve

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jhanvi857/gitresolve/internal/store"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage the gitresolve database",
}

var dbRepairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Check database integrity and repair if corrupt",
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot, err := ResolveRepoRoot()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		dbPath := dbPathForRepo(repoRoot)
		fmt.Printf("Checking database at: %s\n", dbPath)

		db, err := store.Open(dbPath)
		if err != nil {
			if errors.Is(err, store.ErrDBCorrupt) {
				fmt.Printf("Database corruption detected: %v\n", err)
				repairDB(dbPath)
				return
			}
			fmt.Printf("Failed to open database: %v\n", err)
			return
		}
		db.Close()

		fmt.Println("Database integrity check passed.")
	},
}

func repairDB(dbPath string) {
	timestamp := time.Now().Unix()
	backupPath := fmt.Sprintf("%s.corrupt.%d", dbPath, timestamp)

	fmt.Printf("Moving corrupt database to: %s\n", backupPath)
	if err := os.Rename(dbPath, backupPath); err != nil {
		fmt.Printf("Failed to backup corrupt database: %v\n", err)
		return
	}

	fmt.Println("Initializing fresh database...")
	db, err := store.Open(dbPath)
	if err != nil {
		fmt.Printf("Failed to initialize fresh database: %v\n", err)
		return
	}
	db.Close()
	fmt.Println("Database successfully repaired (re-initialized).")
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbRepairCmd)
}
