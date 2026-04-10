package conflict

import "fmt"

const (
	ReasonParserMalformedNestedMarker = "parser.malformed_nested_marker"
	ReasonParserMissingDivider        = "parser.missing_divider"
	ReasonSemanticUnsupportedLanguage = "semantic.unsupported_language"
	ReasonSemanticParseFailed         = "semantic.parse_failed"
	ReasonSafetyIncompleteStructure   = "safety.incomplete_structure"
	ReasonStrategyBothBlockedRisk     = "strategy.both_blocked_high_risk"
	ReasonValidationSyntaxFailed      = "validation.syntax_failed"
	ReasonStructuredAutoDisabled      = "structured.auto_disabled"
	ReasonStructuredOverlap           = "structured.overlap"
	ReasonStructuredParseFailed       = "structured.parse_failed"
	ReasonImportOverlapCritical       = "import.overlap_critical"
	ReasonImportMergeFailed           = "import.merge_failed"
	ReasonImportParseFailed           = "import.parse_failed"
	ReasonDecisionUnknown             = "decision.unknown"
	ReasonShadowDiff                  = "decision.shadow_diff"
)

var stableReasonCodeSet = map[string]struct{}{
	ReasonParserMalformedNestedMarker: {},
	ReasonParserMissingDivider:        {},
	ReasonSemanticUnsupportedLanguage: {},
	ReasonSemanticParseFailed:         {},
	ReasonSafetyIncompleteStructure:   {},
	ReasonStrategyBothBlockedRisk:     {},
	ReasonValidationSyntaxFailed:      {},
	ReasonStructuredAutoDisabled:      {},
	ReasonStructuredOverlap:           {},
	ReasonStructuredParseFailed:       {},
	ReasonImportOverlapCritical:       {},
	ReasonImportMergeFailed:           {},
	ReasonImportParseFailed:           {},
	ReasonDecisionUnknown:             {},
	ReasonShadowDiff:                  {},
}

func IsStableReasonCode(code string) bool {
	_, ok := stableReasonCodeSet[code]
	return ok
}

func SetManualEscalation(c *ConflictBlock, code, reason, hint string) {
	if c == nil {
		return
	}
	if code != "" && IsStableReasonCode(code) {
		c.ManualReasonCode = code
	} else if code != "" {
		// Keep the human reason useful even if caller used an unknown code.
		reason = fmt.Sprintf("%s (invalid reason code: %s)", reason, code)
	}
	if reason != "" {
		c.ManualReason = reason
	}
	if hint != "" {
		c.SuggestHint = hint
	}
}
