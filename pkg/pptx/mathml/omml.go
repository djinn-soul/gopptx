// Package mathml renders a small LaTeX subset as Office MathML (OMML), the
// markup PowerPoint's own equation editor uses. Writing an equation natively
// means it stays editable and scales with the slide, instead of being pasted in
// as a picture (upstream issue #126).
package mathml

import (
	"fmt"
	"strings"
)

// Namespaces used by an inline equation inside a drawing text body.
const (
	// MathNS is the OMML namespace.
	MathNS = "http://schemas.openxmlformats.org/officeDocument/2006/math"
	// A14NS is the DrawingML 2010 extension namespace that carries an equation
	// inside a shape's text body.
	A14NS = "http://schemas.microsoft.com/office/drawing/2010/main"
	// MCNS is the markup-compatibility namespace.
	MCNS = "http://schemas.openxmlformats.org/markup-compatibility/2006"
)

// greekLetters maps the LaTeX command name to its Unicode character. OMML holds
// the literal character, not the command.
//
//nolint:gochecknoglobals // immutable lookup table
var greekLetters = map[string]string{
	"alpha": "α", "beta": "β", "gamma": "γ", "delta": "δ",
	"epsilon": "ε", "zeta": "ζ", "eta": "η", "theta": "θ",
	"iota": "ι", "kappa": "κ", "lambda": "λ", "mu": "μ",
	"nu": "ν", "xi": "ξ", "pi": "π", "rho": "ρ",
	"sigma": "σ", "tau": "τ", "phi": "φ", "chi": "χ",
	"psi": "ψ", "omega": "ω",
	"Gamma": "Γ", "Delta": "Δ", "Theta": "Θ", "Lambda": "Λ",
	"Xi": "Ξ", "Pi": "Π", "Sigma": "Σ", "Phi": "Φ",
	"Psi": "Ψ", "Omega": "Ω",
}

// symbols maps other supported LaTeX commands to their Unicode characters.
//
//nolint:gochecknoglobals // immutable lookup table
var symbols = map[string]string{
	"times": "×", "div": "÷", "pm": "±", "mp": "∓",
	"leq": "≤", "geq": "≥", "neq": "≠", "approx": "≈",
	"cdot": "⋅", "infty": "∞", "partial": "∂",
	"rightarrow": "→", "leftarrow": "←", "to": "→",
	"in": "∈", "forall": "∀", "exists": "∃",
}

// naryOperators maps the n-ary LaTeX commands to their operator characters.
//
//nolint:gochecknoglobals // immutable lookup table
var naryOperators = map[string]string{
	"sum": "∑", "prod": "∏", "int": "∫",
}

// ToOMML renders a LaTeX fragment as an <m:oMath> element.
//
// The supported subset is: literal text and numbers, Greek letters and the
// symbols above, ^ and _ scripts, \frac{a}{b}, \sqrt{x}, \sqrt[n]{x},
// \sum/\prod/\int with optional _{lower} and ^{upper}, and {} grouping.
// Anything outside it is an error rather than silently dropped markup.
func ToOMML(latex string) (string, error) {
	parser := &parser{src: []rune(latex)}
	body, err := parser.parseSequence(0)
	if err != nil {
		return "", err
	}
	if parser.pos < len(parser.src) {
		return "", fmt.Errorf("unexpected %q at position %d", string(parser.src[parser.pos]), parser.pos)
	}
	return `<m:oMath xmlns:m="` + MathNS + `">` + body + `</m:oMath>`, nil
}

// ParagraphXML wraps an equation in the mc:AlternateContent a shape's text body
// needs, leaving the text size to the theme. See ParagraphXMLSized.
func ParagraphXML(latex string) (string, error) {
	return ParagraphXMLSized(latex, 0)
}

// ParagraphXMLSized wraps an equation in the mc:AlternateContent a shape's text
// body needs: the mc:Choice holds the real equation for readers that understand
// the a14 extension, and the mc:Fallback holds the source text so the paragraph
// is never empty for readers that do not.
//
// sizeHundredths is the run size in hundredths of a point; zero leaves it to the
// theme. The size has to land on the math runs themselves -- run properties on a
// following paragraph do not flow backwards into an equation.
func ParagraphXMLSized(latex string, sizeHundredths int) (string, error) {
	omml, err := ToOMML(latex)
	if err != nil {
		return "", err
	}
	runProps := runSizeProps(sizeHundredths)
	var b strings.Builder
	b.WriteString(`<a:p><mc:AlternateContent xmlns:mc="` + MCNS + `">`)
	b.WriteString(`<mc:Choice xmlns:a14="` + A14NS + `" Requires="a14"><a14:m>`)
	b.WriteString(`<m:oMathPara xmlns:m="` + MathNS + `">`)
	body := strings.Replace(omml, `<m:oMath xmlns:m="`+MathNS+`">`, `<m:oMath>`, 1)
	b.WriteString(sizeMathRuns(body, runProps))
	b.WriteString(`</m:oMathPara></a14:m></mc:Choice>`)
	b.WriteString(`<mc:Fallback><a:r>` + runProps + `<a:t>`)
	b.WriteString(escapeXML(latex))
	b.WriteString(`</a:t></a:r></mc:Fallback>`)
	b.WriteString(`</mc:AlternateContent></a:p>`)
	return b.String(), nil
}

// runSizeProps renders the a:rPr every run in the equation carries.
func runSizeProps(sizeHundredths int) string {
	if sizeHundredths <= 0 {
		return `<a:rPr lang="en-US" dirty="0"/>`
	}
	return fmt.Sprintf(`<a:rPr lang="en-US" sz="%d" dirty="0"/>`, sizeHundredths)
}

// sizeMathRuns puts the run properties on every math run. runXML is the only
// producer of <m:r>, and it always writes the two tags adjacent, so this reaches
// each run exactly once.
func sizeMathRuns(omml, runProps string) string {
	return strings.ReplaceAll(omml, `<m:r><m:t>`, `<m:r>`+runProps+`<m:t>`)
}
