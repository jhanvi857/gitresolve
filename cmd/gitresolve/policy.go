package gitresolve

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	securejoin "github.com/cyphar/filepath-securejoin"
	"github.com/jhanvi857/gitresolve/internal/ownership"
	"github.com/jhanvi857/gitresolve/internal/safepath"
	"github.com/spf13/cobra"
)

var policyCheckProfile string
var policyCheckJSON bool

const policySchemaVersion = "1.0"

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Inspect policy profile resolution",
	Long:  `Inspect policy profile resolution for a specific file path using explicit, path, team, and default policy chain.`,
}

var policyCheckCmd = &cobra.Command{
	Use:   "check <file>",
	Short: "Show the profile that applies to a file",
	Long:  `Resolves policy profile for a given file and reports whether the profile came from explicit flag, path rule, team rule, or default.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot, err := ResolveRepoRoot()
		if err != nil {
			fmt.Println("Policy check failed:", err)
			return
		}

		root, err := safepath.RepoRoot(repoRoot)
		if err != nil {
			fmt.Println("Policy check failed:", err)
			return
		}
		defer root.Close()

		inputPath := args[0]
		absPath, err := securejoin.SecureJoin(repoRoot, inputPath)
		if err != nil {
			fmt.Println("Policy check failed: path traversal rejected:", err)
			return
		}

		relPath, err := filepath.Rel(repoRoot, absPath)
		if err != nil {
			fmt.Println("Policy check failed: could not resolve path relative to repo root:", err)
			return
		}
		relPath = filepath.ToSlash(relPath)

		if err := ValidatePath(repoRoot, relPath); err != nil {
			fmt.Println("Policy check failed:", err)
			return
		}

		resolution, err := ownership.ResolvePolicy(root, relPath, policyCheckProfile)
		if err != nil {
			fmt.Println("Policy check failed:", err)
			return
		}

		strictBlocksBoth := policyBlocksBothForFile(resolution.ResolvedProfile, relPath)
		if policyCheckJSON {
			payload := map[string]interface{}{
				"schema_version":                policySchemaVersion,
				"file":                          relPath,
				"policy_resolution":             resolution,
				"strict_blocks_both_for_file":   strictBlocksBoth,
				"strict_source_like_extensions": strictPolicySourceLikeExtensions(),
			}
			enc, err := json.MarshalIndent(payload, "", "  ")
			if err != nil {
				fmt.Println("Error encoding policy check response:", err)
				return
			}
			fmt.Println(string(enc))
			return
		}

		fmt.Println("Policy Check")
		fmt.Printf("  file: %s\n", relPath)
		fmt.Printf("  requested_profile: %s\n", resolution.RequestedProfile)
		fmt.Printf("  resolved_profile: %s\n", resolution.ResolvedProfile)
		fmt.Printf("  source: %s\n", resolution.Source)
		if resolution.MatchedPath != "" {
			fmt.Printf("  matched_path: %s\n", resolution.MatchedPath)
		}
		if resolution.MatchedTeam != "" {
			fmt.Printf("  matched_team: %s\n", resolution.MatchedTeam)
		}
		fmt.Printf("  strict_blocks_both_for_file: %t\n", strictBlocksBoth)
		if strictBlocksBoth {
			fmt.Printf("  strict_source_like_extensions: %s\n", strings.Join(strictPolicySourceLikeExtensions(), ", "))
		}
	},
}

func init() {
	rootCmd.AddCommand(policyCmd)
	policyCmd.AddCommand(policyCheckCmd)
	policyCheckCmd.Flags().StringVar(&policyCheckProfile, "policy-profile", ownership.PolicyAuto, "policy profile: auto|strict|balanced|aggressive")
	policyCheckCmd.Flags().BoolVar(&policyCheckJSON, "json", false, "emit policy check result in JSON format")
}
