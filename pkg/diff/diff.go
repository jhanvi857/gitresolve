package diff

import (
	"fmt"
	"strings"
)

type DiffLineType int

const (
	LineUnchanged DiffLineType = iota
	LineAdded
	LineDeleted
)

type DiffLine struct {
	Type    DiffLineType
	Content string
}

func (dl DiffLine) String() string {
	switch dl.Type {
	case LineAdded:
		return "+ " + dl.Content
	case LineDeleted:
		return "- " + dl.Content
	default:
		return "  " + dl.Content
	}
}

type Hunk struct {
	StartLineBase int
	StartLineNew  int
	Lines         []DiffLine
}

func (h Hunk) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("@@ -%d +%d @@\n", h.StartLineBase, h.StartLineNew))
	for _, l := range h.Lines {
		sb.WriteString(l.String() + "\n")
	}
	return sb.String()
}
