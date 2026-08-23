package conflict

// file for : defining type of conflicts and severity + strcut after including all things with resolution.
// ConflictType tells you what kind of conflict this is
type ConflictType int

// types of conflicts :
// only whitespace differs
// import block changes
// function or variable renamed
// function signature changed
// logic inside function changed
// JSON/YAML/TOML key conflict
// one side deleted, other modified
// classifier could not determine
const (
	TypeWhitespace ConflictType = iota
	TypeImport
	TypeIdentical
	TypeRename
	TypeSignature
	TypeLogic
	TypeStructured
	TypeDeleteModify
	TypeScalar
	TypeUnknown
)

// dangerous conflict or not if yes how much
type Severity int

const (
	SeverityTrivial Severity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

type ConflictBlock struct {
	FilePath         string
	StartLine        int
	EndLine          int
	StartIndex       int
	EndIndex         int
	OursLines        []string
	TheirsLines      []string
	BaseLines        []string
	PreLines         []string
	PostLines        []string
	Type             ConflictType
	Severity         Severity
	Confidence       float64
	CanAutoResolve   bool
	Resolution       string
	ManualReasonCode string
	ManualReason     string
	SuggestHint      string
}
