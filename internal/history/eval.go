package history

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jhanvi857/gitresolve/internal/conflict"
)

// EvaluateBlock evaluates history-aware escalation rules for a single conflict block.
// It checks in priority order:
// 1. High blast radius: symbol modified has > maxCallers callers across the repository.
// 2. Missing coupled file: a file with co-change strength >= coChangeMinStrength was not touched.
// 3. Import cycle: the conflict block file sits inside a detected package import cycle.
// 4. Multi-author conflict: both sides were authored by different people (informational).
func (idx *Index) EvaluateBlock(
	repoPath string,
	file string,
	c *conflict.ConflictBlock,
	maxCallers int,
	coChangeMinStrength float64,
	touchedFiles map[string]bool,
) {
	if c == nil || idx == nil {
		return
	}

	if maxCallers <= 0 {
		maxCallers = 10
	}
	if coChangeMinStrength <= 0 {
		coChangeMinStrength = 0.6
	}

	// 1. High blast radius: symbol modified has > maxCallers callers
	symbols := ExtractSymbolsFromLines(append(c.OursLines, c.TheirsLines...))
	if enc := ExtractEnclosingSymbol(repoPath, file, c.StartLine); enc != "" {
		symbols = append(symbols, enc)
	}

	for _, sym := range symbols {
		callers := idx.CallersOfSymbol(sym)
		if len(callers) > maxCallers {
			data := map[string]string{
				"symbol": sym,
				"count":  strconv.Itoa(len(callers)),
			}
			msg := conflict.RenderEscalationMessage(conflict.ReasonSemanticHighBlastRadius, data)
			hint := conflict.RenderSuggestedCommand(conflict.ReasonSemanticHighBlastRadius, data)
			conflict.SetManualEscalation(c, conflict.ReasonSemanticHighBlastRadius, msg, hint)
			c.CanAutoResolve = false
			c.Severity = conflict.SeverityHigh
			c.Confidence = 0.20
			return
		}
	}

	// 2. Missing coupled file: strongly coupled file not touched in this branch
	missingCoupled := idx.CoupledFilesMissing(file, touchedFiles, coChangeMinStrength)
	if len(missingCoupled) > 0 {
		top := missingCoupled[0]
		coupled := top.FileA
		if coupled == file {
			coupled = top.FileB
		}
		data := map[string]string{
			"file":         file,
			"coupled_file": coupled,
			"strength":     fmt.Sprintf("%.2f", top.Strength),
		}
		msg := conflict.RenderEscalationMessage(conflict.ReasonSemanticMissingCoupledFile, data)
		hint := conflict.RenderSuggestedCommand(conflict.ReasonSemanticMissingCoupledFile, map[string]string{"coupled_file": coupled})
		conflict.SetManualEscalation(c, conflict.ReasonSemanticMissingCoupledFile, msg, hint)
		c.CanAutoResolve = false
		c.Severity = conflict.SeverityHigh
		c.Confidence = 0.25
		return
	}

	// 3. Import cycle: conflict sits inside an import cycle
	if inCycle, cycle := idx.FileInImportCycle(file); inCycle {
		cycleStr := strings.Join([]string(cycle), " → ") + " → " + cycle[0]
		data := map[string]string{
			"file":  file,
			"cycle": cycleStr,
		}
		msg := conflict.RenderEscalationMessage(conflict.ReasonSemanticImportCycle, data)
		hint := conflict.RenderSuggestedCommand(conflict.ReasonSemanticImportCycle, nil)
		conflict.SetManualEscalation(c, conflict.ReasonSemanticImportCycle, msg, hint)
		c.CanAutoResolve = false
		c.Severity = conflict.SeverityCritical
		c.Confidence = 0.15
		return
	}

	// 4. Multi-author conflict (informational / for decision log)
	authors := idx.AuthorsForFile(file)
	if len(authors) >= 2 && c.ManualReasonCode == "" && !c.CanAutoResolve {
		data := map[string]string{
			"file":          file,
			"author_ours":   authors[0].Email,
			"author_theirs": authors[1].Email,
		}
		msg := conflict.RenderEscalationMessage(conflict.ReasonStrategyMultiAuthorConflict, data)
		hint := conflict.RenderSuggestedCommand(conflict.ReasonStrategyMultiAuthorConflict, map[string]string{"author_theirs": authors[1].Email})
		conflict.SetManualEscalation(c, conflict.ReasonStrategyMultiAuthorConflict, msg, hint)
	}
}
