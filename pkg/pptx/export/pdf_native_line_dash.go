package export

import (
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// A shape, connector or rich line can ask for a dashed outline through
// a:prstDash, and the PPTX reader has always carried that value through — but
// the native renderer stroked every outline solid, so a dashed callout border
// or a dotted connector printed as a plain line.
//
// OOXML states a preset dash as a repeating run of dash and gap lengths in
// multiples of the line width, which is exactly the shape of a PDF dash array,
// so the mapping is a scale by the stroke width rather than a translation.
//
// The multiples below are the ones ECMA-376 lists for each preset.

// Keys are lower-cased: a dash reaches the renderer either as one of the
// package's canonical constants or, for the values those constants do not
// cover, straight off the a:prstDash attribute.
//
//nolint:gochecknoglobals // Immutable preset table shared by every stroke.
var ooxmlDashPatterns = map[string][]float64{
	shapes.LineDashDot:  {1, 3},
	shapes.LineDashDash: {4, 3},
	"lgdash":            {8, 3},
	"dashdot":           {4, 3, 1, 3},
	"lgdashdot":         {8, 3, 1, 3},
	"dashdotdot":        {8, 3, 1, 3, 1, 3},
	"lgdashdotdot":      {8, 3, 1, 3, 1, 3},
	"sysdash":           {3, 1},
	"sysdot":            {1, 1},
	"sysdashdot":        {3, 1, 1, 1},
	"sysdashdotdot":     {3, 1, 1, 1, 1, 1},
}

// minDashUnitPt keeps a hairline's dashes visible. A 0.75pt line dotted at one
// times its width would produce 0.75pt marks, which is fine, but a line thinner
// than this rounds up so the pattern does not collapse into a solid rule.
const minDashUnitPt = 0.5

// pdfDashPattern returns the PDF dash array for one OOXML preset dash at the
// given stroke width, or nil for a solid line.
func pdfDashPattern(dash string, widthPt float64) []float64 {
	normalized := strings.ToLower(strings.TrimSpace(dash))
	if normalized == "" || normalized == shapes.LineDashSolid {
		return nil
	}
	multiples, ok := ooxmlDashPatterns[normalized]
	if !ok {
		return nil
	}
	unit := max(widthPt, minDashUnitPt)
	pattern := make([]float64, len(multiples))
	for i, m := range multiples {
		pattern[i] = m * unit
	}
	return pattern
}

// applyPDFLineDash switches the stroke to the requested preset dash. A solid or
// unrecognised dash leaves the current state alone, so callers clear the dash
// themselves once the outline is drawn: gopdf keeps it as document-wide state
// that would otherwise leak onto the next shape.
func applyPDFLineDash(pdf *gopdf.GoPdf, dash string, widthPt float64) {
	pattern := pdfDashPattern(dash, widthPt)
	if pattern == nil {
		return
	}
	pdf.SetCustomLineType(pattern, 0)
}

// clearPDFLineDash restores solid stroking.
func clearPDFLineDash(pdf *gopdf.GoPdf) {
	pdf.SetLineType("")
}

// shapeLineDash is the dash a shape asks for. A rich line (a:ln read in full)
// states it separately from the simple line, and wins when both are present
// because it is the more complete reading of the same element.
func shapeLineDash(s shapes.Shape) string {
	if s.RichLine != nil && s.RichLine.DashStyle != "" {
		return string(s.RichLine.DashStyle)
	}
	if s.Line != nil {
		return s.Line.Dash
	}
	return ""
}
