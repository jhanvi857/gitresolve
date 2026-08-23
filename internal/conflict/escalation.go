package conflict

import "strings"

// EscalationTemplates maps each history-aware reason code to a fixed
// plain-English template with named {placeholders}. Placeholders are
// filled from real computed data only — no free-text generation.
var EscalationTemplates = map[string]string{
	ReasonSemanticHighBlastRadius:     "{symbol} is called from {count} other locations — escalating for manual review",
	ReasonSemanticMissingCoupledFile:  "{file} is frequently changed alongside {coupled_file} (strength {strength}) but {coupled_file} was not touched in this branch",
	ReasonSemanticImportCycle:         "conflict block in {file} sits inside an import cycle: {cycle}",
	ReasonStrategyStaleBranchDiv:      "branch is {behind} commits behind {branch} — {authors} authored changes touching files you also modified",
	ReasonStrategyMultiAuthorConflict: "both sides of the conflict in {file} were authored by different people: {author_ours} vs {author_theirs}",
}

// SuggestedCommandTemplates maps each history-aware reason code to a
// concrete command suggestion with optional {placeholders}.
var SuggestedCommandTemplates = map[string]string{
	ReasonSemanticHighBlastRadius:     "go test ./... (run full suite before committing)",
	ReasonSemanticMissingCoupledFile:  "gitresolve status {coupled_file}",
	ReasonSemanticImportCycle:         "review import dependencies and break the cycle",
	ReasonStrategyStaleBranchDiv:      "git fetch && git rebase origin/{branch}",
	ReasonStrategyMultiAuthorConflict: "coordinate with {author_theirs} before resolving",
}

// RenderEscalationMessage fills the template for the given reason code
// using the provided data map. Returns empty string if the code has no
// template registered.
func RenderEscalationMessage(code string, data map[string]string) string {
	tmpl, ok := EscalationTemplates[code]
	if !ok {
		return ""
	}
	return renderTemplate(tmpl, data)
}

// RenderSuggestedCommand fills the command template for the given reason
// code using the provided data map. Returns the fallback command if the
// code has no template registered.
func RenderSuggestedCommand(code string, data map[string]string) string {
	tmpl, ok := SuggestedCommandTemplates[code]
	if !ok {
		return FallbackSuggestedCommand()
	}
	return renderTemplate(tmpl, data)
}

// FallbackSuggestedCommand returns a safe default command for any
// unhandled escalation scenario.
func FallbackSuggestedCommand() string {
	return "git merge --abort"
}

// renderTemplate replaces every {key} in tmpl with data[key].
// Unknown keys are left as-is (defensive — never panic on missing data).
func renderTemplate(tmpl string, data map[string]string) string {
	result := tmpl
	for key, val := range data {
		result = strings.ReplaceAll(result, "{"+key+"}", val)
	}
	return result
}
