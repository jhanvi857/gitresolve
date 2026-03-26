package conflict

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Strategy int

const (
	StrategyOurs Strategy = iota
	StrategyTheirs
	StrategyBoth
	StrategyInteractive
)

func Resolve(c *Conflict, strategy Strategy) error {
	switch strategy {
	case StrategyOurs:
		c.Resolution = strings.Join(c.OurLines, "\n")

	case StrategyTheirs:
		c.Resolution = strings.Join(c.TheirLines, "\n")

	case StrategyBoth:
		both := append(c.OurLines, c.TheirLines...)
		c.Resolution = strings.Join(both, "\n")

	case StrategyInteractive:
		fmt.Printf("\n--- Conflict in %s ---\n", c.FilePath)
		fmt.Println("<<<<<<< OURS")
		fmt.Println(strings.Join(c.OurLines, "\n"))
		fmt.Println("=======")
		fmt.Println(strings.Join(c.TheirLines, "\n"))
		fmt.Println(">>>>>>> THEIRS")
		
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Select resolution [O]urs, [T]heirs, [B]oth : ")
			input, _ := reader.ReadString('\n')
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
			}
		}

	default:
		return fmt.Errorf("Resolve: unknown strategy %d", strategy)
	}

	return nil
}
