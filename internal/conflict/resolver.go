package conflict

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Strategy int

type ResolveOptions struct {
	NonInteractive bool
	Timeout        time.Duration
}

const (
	StrategyOurs Strategy = iota
	StrategyTheirs
	StrategyBoth
	StrategyInteractive
)

func Resolve(c *ConflictBlock, strategy Strategy, opts ResolveOptions) error {
	applySelection := func(strat Strategy) (string, error) {
		pre := strings.Join(c.PreLines, "\n")
		if pre != "" {
			pre += "\n"
		}
		post := strings.Join(c.PostLines, "\n")
		if post != "" {
			post = "\n" + post
		}

		switch strat {
		case StrategyOurs:
			c.Resolution = strings.Join(c.OursLines, "\n")
			return pre + c.Resolution + post, nil

		case StrategyTheirs:
			c.Resolution = strings.Join(c.TheirsLines, "\n")
			return pre + c.Resolution + post, nil

		case StrategyBoth:
			// check for closing braces if it's a Go file and we have function logic
			if strings.HasSuffix(c.FilePath, ".go") && (c.Type == TypeLogic || c.Type == TypeSignature) {
				ourHasBrace := false
				if len(c.OursLines) > 0 {
					ourHasBrace = strings.Contains(c.OursLines[len(c.OursLines)-1], "}")
				}
				theirHasBrace := false
				if len(c.TheirsLines) > 0 {
					theirHasBrace = strings.Contains(c.TheirsLines[len(c.TheirsLines)-1], "}")
				}

				if !ourHasBrace || !theirHasBrace {
					c.CanAutoResolve = false
					c.ManualReason = "incomplete function block detected, manual edit required"
					return "", fmt.Errorf("Both: %s", c.ManualReason)
				}
			}

			both := append(c.OursLines, c.TheirsLines...)
			c.Resolution = strings.Join(both, "\n")
			return pre + c.Resolution + post, nil
		}
		return "", fmt.Errorf("unknown selection strategy")
	}

	switch strategy {
	case StrategyOurs, StrategyTheirs, StrategyBoth:
		_, err := applySelection(strategy)
		return err

	case StrategyInteractive:
		if opts.NonInteractive {
			return fmt.Errorf("conflict in %s requires manual resolution, but --non-interactive is set", c.FilePath)
		}

		if c.Type == TypeScalar {
			fmt.Printf("\n[Scalar] %s (L%d-%d)\n", c.FilePath, c.StartLine, c.EndLine)
			fmt.Printf(" [O]urs:   %s\n", strings.Join(c.OursLines, " "))
			fmt.Printf(" [T]heirs: %s\n", strings.Join(c.TheirsLines, " "))
		} else {
			fmt.Printf("\n--- Conflict in %s ---\n", c.FilePath)
			fmt.Println("<<<<<<< OURS")
			fmt.Println(strings.Join(c.OursLines, "\n"))
			fmt.Println("=======")
			fmt.Println(strings.Join(c.TheirsLines, "\n"))
			fmt.Println(">>>>>>> THEIRS")
		}

		inputChan := make(chan string)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			if c.Type == TypeScalar {
				fmt.Print("Resolve [O|T|B]: ")
			} else {
				fmt.Print("Select resolution [O]urs, [T]heirs, [B]oth : ")
			}
			input, _ := reader.ReadString('\n')
			inputChan <- input
		}()

		var input string
		if opts.Timeout > 0 {
			select {
			case input = <-inputChan:
				// Proceed
			case <-time.After(opts.Timeout):
				fmt.Printf("\nTimeout reached (%s). Auto-selecting [T]heirs.\n", opts.Timeout.String())
				_, err := applySelection(StrategyTheirs)
				return err
			}
		} else {
			input = <-inputChan
		}

		for {
			input = strings.TrimSpace(strings.ToUpper(input))
			var strat Strategy
			valid := true
			if input == "O" || input == "OURS" {
				strat = StrategyOurs
			} else if input == "T" || input == "THEIRS" {
				strat = StrategyTheirs
			} else if input == "B" || input == "BOTH" {
				strat = StrategyBoth
			} else {
				valid = false
			}

			if valid {
				output, err := applySelection(strat)
				if err != nil {
					// Both check might fail
					fmt.Printf("ERROR: %v\n", err)
				} else {
					// Run validation BEFORE write (well, the write is done in the caller, but we validate here as requested)
					if err := Verify(c.FilePath, output); err != nil {
						vErr, ok := err.(*VerificationError)
						if ok && vErr.IsMarkerErr {
							// For files with multiple conflicts, the intermediate state will still have markers.
							// We allow this to continue so other blocks can be resolved.
							return nil
						}
						// Real syntax error (parsers failed)
						fmt.Println("Resolution produced invalid syntax. File left unchanged.")
						fmt.Printf("Run: gitresolve resolve --file %s to retry this file.\n", c.FilePath)
						return err
					}
					return nil
				}
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Select resolution [O]urs, [T]heirs, [B]oth : ")
			input, _ = reader.ReadString('\n')
		}

	default:
		return fmt.Errorf("Resolve: unknown strategy %d", strategy)
	}

	// return nil
}
