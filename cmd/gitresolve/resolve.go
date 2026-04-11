package gitresolve

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jhanvi857/gitresolve/internal/conflict"
	"github.com/jhanvi857/gitresolve/internal/git"
	"github.com/jhanvi857/gitresolve/internal/ownership"
	"github.com/jhanvi857/gitresolve/internal/safety"
	"github.com/jhanvi857/gitresolve/internal/store"
	gserrors "github.com/jhanvi857/gitresolve/pkg/errors"
	"github.com/jhanvi857/gitresolve/pkg/logger"
	"github.com/spf13/cobra"
)

var resolveFileName string
var resolveStrategy string
var resolveDryRun bool
var resolveNonInteractive bool
var resolveTimeout time.Duration
var resolveShadow bool
var resolveEnforceGates bool
var resolveManualRateGate float64
var resolvePolicyProfile string

var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Interactively resolve remaining conflicts",
	Long:  `Fall back to manual or interactive resolution for conflicts that could not be auto-resolved.`,
	Run: func(cmd *cobra.Command, args []string) {
		strategy, err := parseStrategy(resolveStrategy)
		if err != nil {
			fmt.Println("Resolve failed:", err)
			return
		}

		r, err := git.Open(".")
		if err != nil {
			fmt.Println("Fatal: Failed to open git repository:", err)
			return
		}
		defer git.Close(r)
		HandleSignals(r)

		repoPath := "."
		files, err := git.ConflictedFiles(r)
		if err != nil {
			fmt.Println("No unmerged files in index. Scanning for mis-staged markers...")
			files, _ = git.ScanForMarkers(repoPath)
		}

		if len(files) == 0 {
			fmt.Println("No conflicts found (index or content).")
			return
		}

		writer := safety.NewWriter(resolveDryRun)

		db, dbErr := openStore(repoPath)
		if dbErr == nil {
			defer db.Close()
		}

		if !resolveDryRun && dbErr == nil {
			head, headErr := r.HeadCommit()
			if headErr == nil {
				_ = db.SaveSession(repoPath, "resolve", head)
				_ = git.StoreHead(repoPath, head)
			}
		}

		autoResolved := 0
		interactiveResolved := 0
		validationFailed := 0
		filesUpdated := 0
		totalDecisions := 0
		manualEscalations := 0
		var failedFiles []string

		for _, file := range files {
			if resolveFileName != "" && file != resolveFileName {
				continue
			}

			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading %s: %v\n", file, err)
				continue
			}

			conflicts := conflict.ParseFile(file, content)
			logger.Debug(fmt.Sprintf("parsed %d conflict block(s) in %s", len(conflicts), file))
			if len(conflicts) == 0 {
				continue
			}

			if !resolveDryRun {
				if err := safety.PreserveOriginal(file); err != nil {
					fmt.Printf("Warning: backup failed for %s: %v\n", file, err)
					continue
				}
			}

			fileValidationFailed := false
			fileSkipped := false
			for _, c := range conflicts {
				totalDecisions++
				logger.Debug(fmt.Sprintf("conflict block parsed: file=%s start=%d end=%d ours=%d theirs=%d", file, c.StartLine, c.EndLine, len(c.OursLines), len(c.TheirsLines)))
				conflict.Classify(c)

				resolvedPolicy, policyErr := ownership.ResolvePolicyProfile(repoPath, file, resolvePolicyProfile)
				if policyErr != nil {
					fmt.Printf("Warning: policy resolution failed for %s: %v (falling back to balanced)\n", file, policyErr)
					resolvedPolicy = ownership.PolicyBalanced
				}
				if strategy == conflict.StrategyBoth && policyBlocksBothForFile(resolvedPolicy, file) {
					manualEscalations++
					conflict.SetManualEscalation(c, conflict.ReasonStrategyBothBlockedRisk, "BOTH disabled by strict policy profile for source file", "use ours/theirs/manual under strict policy")
					if dbErr == nil {
						_ = db.SaveDecision(store.DecisionRecord{
							RepoPath:     repoPath,
							FilePath:     file,
							Operation:    "resolve",
							ConflictType: typeLabel(c.Type),
							Severity:     severityLabel(c.Severity),
							Action:       "manual-escalate",
							ReasonCode:   reasonCodeOrUnknown(c),
							Reason:       c.ManualReason,
							Confidence:   c.Confidence,
							Shadow:       resolveShadow,
						})
					}
					fileSkipped = true
					continue
				}
				isAuto := false
				if shouldAutoApplyWithPolicy(c, resolvedPolicy) {
					isAuto = true
				} else {
					if conflict.NeedsGuidedChoice(c) {
						fmt.Printf("Guided choice: %s L%d-%d [confidence=%.2f]. Suggested strategy: ours|theirs|both\n", file, c.StartLine, c.EndLine, c.Confidence)
					}
					if c.ManualReason != "" {
						fmt.Printf("  reason: %s\n", c.ManualReason)
					}
					if c.SuggestHint != "" {
						fmt.Printf("  hint: %s\n", c.SuggestHint)
					}
				}

				opts := conflict.ResolveOptions{
					NonInteractive: resolveNonInteractive,
					Timeout:        resolveTimeout,
				}
				result, resolveErr := conflict.Resolve(c, strategy, opts)
				if resolveErr != nil {
					fmt.Printf("Resolve failed for %s: %v\n", file, resolveErr)
					logger.Debug(fmt.Sprintf("resolution failure: file=%s start=%d end=%d err=%v", file, c.StartLine, c.EndLine, resolveErr))
					manualEscalations++
					if dbErr == nil {
						_ = db.SaveDecision(store.DecisionRecord{
							RepoPath:     repoPath,
							FilePath:     file,
							Operation:    "resolve",
							ConflictType: typeLabel(c.Type),
							Severity:     severityLabel(c.Severity),
							Action:       "manual-escalate",
							ReasonCode:   reasonCodeOrUnknown(c),
							Reason:       c.ManualReason,
							Confidence:   c.Confidence,
							Shadow:       resolveShadow,
						})
					}
					validationFailed++
					failedFiles = append(failedFiles, file)
					fileValidationFailed = true
					if resolveNonInteractive {
						os.Exit(1)
					}
					break
				}

				logger.Debug(fmt.Sprintf("resolution selected: file=%s start=%d end=%d choice=%s applied=%v", file, c.StartLine, c.EndLine, result.SelectedLabel, result.Applied))
				if !result.Applied {
					manualEscalations++
					if dbErr == nil {
						_ = db.SaveDecision(store.DecisionRecord{
							RepoPath:     repoPath,
							FilePath:     file,
							Operation:    "resolve",
							ConflictType: typeLabel(c.Type),
							Severity:     severityLabel(c.Severity),
							Action:       "manual",
							ReasonCode:   reasonCodeOrUnknown(c),
							Reason:       c.ManualReason,
							Confidence:   c.Confidence,
							Shadow:       resolveShadow,
						})
					}
					fileSkipped = true
					continue
				}

				if dbErr == nil {
					action := "resolve"
					if isAuto {
						action = "auto-resolve"
					}
					_ = db.SaveDecision(store.DecisionRecord{
						RepoPath:     repoPath,
						FilePath:     file,
						Operation:    "resolve",
						ConflictType: typeLabel(c.Type),
						Severity:     severityLabel(c.Severity),
						Action:       action,
						ReasonCode:   reasonCodeOrUnknown(c),
						Reason:       c.ManualReason,
						Confidence:   c.Confidence,
						Shadow:       resolveShadow,
					})
				}

				if isAuto {
					autoResolved++
				} else if strategy == conflict.StrategyInteractive {
					interactiveResolved++
				}
			}

			if fileValidationFailed {
				continue
			}
			if fileSkipped {
				fmt.Printf("Skipped unresolved blocks in %s; leaving file unchanged.\n", file)
				continue
			}

			newContent := conflict.CompileResolution(content, conflicts)
			if resolveShadow {
				if dbErr == nil {
					_ = db.SaveDecision(store.DecisionRecord{
						RepoPath:      repoPath,
						FilePath:      file,
						Operation:     "resolve",
						ConflictType:  "file",
						Severity:      "info",
						Action:        "shadow-diff",
						ReasonCode:    conflict.ReasonShadowDiff,
						Reason:        "shadow simulation recorded",
						Confidence:    1,
						Shadow:        true,
						OriginalHash:  hashContent(content),
						SimulatedHash: hashContent([]byte(newContent)),
					})
				}
				fmt.Printf("[shadow] simulated resolution for %s (no write)\n", file)
				continue
			}
			if strings.HasSuffix(file, ".go") {
				if err := conflict.ValidateGoSyntax(file, newContent); err != nil {
					reason := "reconstructed output failed Go syntax validation"
					fmt.Printf("Escalating %s to manual: %s (%v)\n", file, reason, err)
					for _, c := range conflicts {
						conflict.SetManualEscalation(c, conflict.ReasonValidationSyntaxFailed, reason, "resolve manually with --strategy ours|theirs")
					}
					manualEscalations += len(conflicts)
					validationFailed++
					failedFiles = append(failedFiles, file)
					if resolveNonInteractive {
						os.Exit(1)
					}
					continue
				}
			}
			if err := conflict.EnsureNoConflictMarkers(file, newContent); err != nil {
				fmt.Printf("Safety check failed for %s: %v\n", file, err)
				validationFailed++
				failedFiles = append(failedFiles, file)
				logger.Debug(fmt.Sprintf("marker cleanup failed: file=%s err=%v", file, err))
				if resolveNonInteractive {
					os.Exit(1)
				}
				continue
			}
			if err := conflict.Verify(file, newContent); err != nil {
				fmt.Printf("Verification failed for %s: %v\n", file, err)
				logger.Debug(fmt.Sprintf("validation failure: file=%s err=%v", file, err))
				validationFailed++
				failedFiles = append(failedFiles, file)
				if resolveNonInteractive {
					os.Exit(1)
				}
				continue
			}

			if err := writer.Write(file, []byte(newContent)); err != nil {
				if resolveDryRun && errors.Is(err, gserrors.ErrDryRun) {
					fmt.Printf("[dry-run] would resolve %s using strategy '%s'\n", file, resolveStrategy)
					filesUpdated++
					continue
				}
				fmt.Printf("Error writing %s: %v\n", file, err)
				logger.Debug(fmt.Sprintf("write failure: file=%s err=%v", file, err))
				validationFailed++
				failedFiles = append(failedFiles, file)
				if resolveNonInteractive {
					os.Exit(1)
				}
				continue
			}

			if !resolveDryRun {
				if err := git.MarkResolved(r, file); err != nil {
					fmt.Printf("Warning: failed to stage %s: %v\n", file, err)
				}
			}

			if dbErr == nil {
				for _, c := range conflicts {
					_ = db.SaveConflict(storeConflict(repoPath, file, c, resolveStrategy))
				}
			}

			fmt.Printf("Resolved %s using strategy '%s'\n", file, resolveStrategy)
			filesUpdated++
		}

		if resolveFileName != "" && filesUpdated == 0 && validationFailed == 0 {
			fmt.Printf("No conflicted file matched '%s'.\n", resolveFileName)
		}

		fmt.Printf("\nResolve complete. Summary:\n")
		fmt.Printf("  auto_resolved: %d\n", autoResolved)
		fmt.Printf("  interactive_resolved: %d\n", interactiveResolved)
		fmt.Printf("  manual_escalations: %d\n", manualEscalations)
		fmt.Printf("  total_decisions: %d\n", totalDecisions)
		fmt.Printf("  validation_failed: %d\n", validationFailed)
		fmt.Printf("  files_updated: %d\n", filesUpdated)
		if totalDecisions > 0 {
			manualRate := (float64(manualEscalations) / float64(totalDecisions)) * 100
			fmt.Printf("  manual_escalation_rate: %.2f%%\n", manualRate)
			if resolveEnforceGates && manualRate > resolveManualRateGate {
				fmt.Printf("Release gate failed: manual escalation rate %.2f%% exceeds threshold %.2f%%\n", manualRate, resolveManualRateGate)
				os.Exit(1)
			}
		}

		if validationFailed > 0 {
			fmt.Println("\nFiles with validation failures:")
			for _, f := range failedFiles {
				fmt.Printf("  - %s\n", f)
			}
			os.Exit(1)
		}
	},
}

func parseStrategy(v string) (conflict.Strategy, error) {
	switch v {
	case "ours":
		return conflict.StrategyOurs, nil
	case "theirs":
		return conflict.StrategyTheirs, nil
	case "both":
		return conflict.StrategyBoth, nil
	case "interactive":
		return conflict.StrategyInteractive, nil
	default:
		return 0, fmt.Errorf("unknown strategy '%s' (use interactive|ours|theirs|both)", v)
	}
}

func storeConflict(repoPath, file string, c *conflict.ConflictBlock, strategy string) store.ConflictRecord {
	return store.ConflictRecord{
		RepoPath:     repoPath,
		FilePath:     file,
		ConflictType: typeLabel(c.Type),
		Severity:     severityLabel(c.Severity),
		Strategy:     strategy,
	}
}

func init() {
	rootCmd.AddCommand(resolveCmd)
	resolveCmd.Flags().StringVar(&resolveFileName, "file", "", "resolve a specific file")
	resolveCmd.Flags().StringVar(&resolveStrategy, "strategy", "interactive", "resolve strategy: interactive|ours|theirs|both")
	resolveCmd.Flags().BoolVar(&resolveDryRun, "dry-run", false, "show what would happen without writing")
	resolveCmd.Flags().BoolVar(&resolveDryRun, "dryrun", false, "alias for --dry-run")
	resolveCmd.Flags().BoolVar(&resolveShadow, "shadow", false, "simulate resolution and record hash-only diff decisions without writing")
	resolveCmd.Flags().StringVar(&resolvePolicyProfile, "policy-profile", ownership.PolicyAuto, "policy profile: auto|strict|balanced|aggressive")
	resolveCmd.Flags().BoolVar(&resolveNonInteractive, "non-interactive", false, "fail on conflicts requiring manual resolution instead of prompting")
	resolveCmd.Flags().BoolVar(&resolveEnforceGates, "enforce-gates", false, "enforce release gate thresholds (manual rate and validation failures)")
	resolveCmd.Flags().Float64Var(&resolveManualRateGate, "manual-rate-gate", 60, "maximum allowed manual escalation rate percentage when --enforce-gates is set")
	resolveCmd.Flags().DurationVar(&resolveTimeout, "timeout", 0, "timeout for interactive prompt (e.g. 30s). Auto-selects theirs if reached.")
}
