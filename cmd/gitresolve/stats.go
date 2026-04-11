package gitresolve

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var statsOperation string
var statsJSON bool
var statsTop int

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show decision and escalation metrics",
	Long:  `Displays aggregated decision logs for operational observability and CI gate monitoring.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := openStore(".")
		if err != nil {
			fmt.Println("Stats failed: could not open history DB:", err)
			return
		}
		defer db.Close()

		actions, err := db.GetDecisionActionCounts(".", statsOperation)
		if err != nil {
			fmt.Println("Stats failed to read action counts:", err)
			return
		}
		reasons, err := db.GetTopDecisionReasons(".", statsOperation, statsTop)
		if err != nil {
			fmt.Println("Stats failed to read top reason codes:", err)
			return
		}

		total := 0
		manual := 0
		for _, a := range actions {
			total += a.Count
			if a.Action == "manual" || a.Action == "manual-escalate" {
				manual += a.Count
			}
		}
		manualRate := 0.0
		if total > 0 {
			manualRate = (float64(manual) / float64(total)) * 100
		}

		if statsJSON {
			payload := map[string]interface{}{
				"operation":              statsOperation,
				"total_decisions":        total,
				"manual_decisions":       manual,
				"manual_escalation_rate": manualRate,
				"actions":                actions,
				"top_reason_codes":       reasons,
			}
			enc, _ := json.MarshalIndent(payload, "", "  ")
			fmt.Println(string(enc))
			return
		}

		fmt.Println("Decision Stats")
		fmt.Printf("  operation: %s\n", statsOperation)
		fmt.Printf("  total_decisions: %d\n", total)
		fmt.Printf("  manual_decisions: %d\n", manual)
		fmt.Printf("  manual_escalation_rate: %.2f%%\n", manualRate)

		fmt.Println("\nAction Counts:")
		for _, a := range actions {
			fmt.Printf("  %-18s %d\n", a.Action, a.Count)
		}

		fmt.Println("\nTop Reason Codes:")
		for _, r := range reasons {
			fmt.Printf("  %-36s %d\n", r.ReasonCode, r.Count)
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.Flags().StringVar(&statsOperation, "operation", "all", "filter by operation: all|resolve|merge")
	statsCmd.Flags().BoolVar(&statsJSON, "json", false, "emit stats in JSON format")
	statsCmd.Flags().IntVar(&statsTop, "top", 8, "number of top reason codes to display")
}
