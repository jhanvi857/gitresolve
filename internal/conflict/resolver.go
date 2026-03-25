package conflict

import (
	"fmt"
	"strings"
)

type Strategy int

const (
	StrategyOurs Strategy = iota
	StrategyTheirs
	StrategyBoth
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

	default:
		return fmt.Errorf("Resolve: unknown strategy %d", strategy)
	}

	return nil
}
