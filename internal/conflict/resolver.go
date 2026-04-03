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

func Resolve(c *Conflict, strategy Strategy, opts ResolveOptions) error {
	switch strategy {
	case StrategyOurs:
		c.Resolution = strings.Join(c.OurLines, "\n")

	case StrategyTheirs:
		c.Resolution = strings.Join(c.TheirLines, "\n")

	case StrategyBoth:
		both := append(c.OurLines, c.TheirLines...)
		c.Resolution = strings.Join(both, "\n")

	case StrategyInteractive:
		if opts.NonInteractive {
			return fmt.Errorf("conflict in %s requires manual resolution, but --non-interactive is set", c.FilePath)
		}

		if c.Type == TypeScalar {
			fmt.Printf("\n[Scalar] %s (L%d-%d)\n", c.FilePath, c.StartLine, c.EndLine)
			fmt.Printf(" [O]urs:   %s\n", strings.Join(c.OurLines, " "))
			fmt.Printf(" [T]heirs: %s\n", strings.Join(c.TheirLines, " "))
		} else {
			fmt.Printf("\n--- Conflict in %s ---\n", c.FilePath)
			fmt.Println("<<<<<<< OURS")
			fmt.Println(strings.Join(c.OurLines, "\n"))
			fmt.Println("=======")
			fmt.Println(strings.Join(c.TheirLines, "\n"))
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
				c.Resolution = strings.Join(c.TheirLines, "\n")
				return nil
			}
		} else {
			input = <-inputChan
		}

		for {
			input = strings.TrimSpace(strings.ToUpper(input))
			
			if input == "O" || input == "OURS" {
				c.Resolution = strings.Join(c.OurLines, "\n")
				break
			} else if input == "T" || input == "THEIRS" {
				c.Resolution = strings.Join(c.TheirLines, "\n")
				break
			} else if input == "B" || input == "BOTH" {
				both := append(c.OurLines, c.TheirLines...)
				c.Resolution = strings.Join(both, "\n")
				break
			} else {
				fmt.Println("Invalid option. Please press O, T, or B.")
				// We don't loop correctly with the channel setup for retry on invalid input if timeout is used natively like this.
				// For real use, we need the inner retry loop without a channel if no timeout, or just a simple block. Just block again:
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Select resolution [O]urs, [T]heirs, [B]oth : ")
				input, _ = reader.ReadString('\n')
			}
		}

	default:
		return fmt.Errorf("Resolve: unknown strategy %d", strategy)
	}

	return nil
}
