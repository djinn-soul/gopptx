package mermaid

// Mermaid edge tokens shared across the flowchart, class and state parsers.
const (
	arrowSolid      = "-->"
	arrowThick      = "==>"
	arrowDotted     = "-.->"
	arrowOpen       = "---"
	arrowThickOpen  = "==="
	arrowDottedOpen = "-.-"
	arrowCircle     = "--o"
	arrowCross      = "--x"
	arrowInherit    = "<|--"
	arrowRealize    = "<|.."
	stateTypeNormal = "normal"
)

// flowArrowTokens are the edge tokens a flowchart line may be split on, longest
// first so "-.->"" is never mistaken for the "-.-" inside it, and "-->" never
// for the "--o"/"--x" it shares a prefix with.
func flowArrowTokens() []string {
	return []string{
		arrowDotted, arrowThick, arrowSolid,
		arrowCircle, arrowCross,
		arrowOpen, arrowThickOpen, arrowDottedOpen,
		"->",
	}
}

// backgroundWhite is the default slide background used by the mermaid themes.
const backgroundWhite = "FFFFFF"
