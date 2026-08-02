package pptxxml

import (
	"math"
	"strconv"
	"strings"
)

// BlankValue marks a data point that has no value at all, as opposed to one
// worth zero. A chart written with blanks leaves the <c:pt> out entirely, which
// is how PowerPoint tells "no reading" from "a reading of zero" -- the
// distinction upstream issue #968 is about, where None came out as 0.
func BlankValue() float64 {
	return math.NaN()
}

// IsBlankValue reports whether a value stands for a missing data point.
func IsBlankValue(value float64) bool {
	return math.IsNaN(value)
}

// writeNumericPoints writes the <c:pt> children of a numeric cache, skipping
// blanks. ptCount still counts every category so the points stay aligned with
// the category axis.
func writeNumericPoints(b *strings.Builder, values []float64) {
	b.WriteString(`
<c:ptCount val="`)
	b.WriteString(strconv.Itoa(len(values)))
	b.WriteString(`"/>`)
	for i, value := range values {
		if IsBlankValue(value) {
			continue
		}
		b.WriteString(`
<c:pt idx="`)
		b.WriteString(strconv.Itoa(i))
		b.WriteString(`"><c:v>`)
		b.WriteString(strconv.FormatFloat(value, 'f', 6, 64))
		b.WriteString(`</c:v></c:pt>`)
	}
}

// plotVisOnlyElement is the element <c:dispBlanksAs> has to follow.
const plotVisOnlyElement = `<c:plotVisOnly val="1"/>`

// DisplayBlanksAs values control how a chart draws a gap in its data.
const (
	// DisplayBlanksAsGap leaves a hole, which is PowerPoint's own default.
	DisplayBlanksAsGap = "gap"
	// DisplayBlanksAsZero plots the blank as zero.
	DisplayBlanksAsZero = "zero"
	// DisplayBlanksAsSpan joins the neighbouring points with a straight line.
	DisplayBlanksAsSpan = "span"
)

// normalizedDisplayBlanksAs maps a caller's token to a schema-valid one,
// falling back to "gap" so a typo cannot produce a package PowerPoint rejects.
func normalizedDisplayBlanksAs(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case DisplayBlanksAsZero:
		return DisplayBlanksAsZero
	case DisplayBlanksAsSpan:
		return DisplayBlanksAsSpan
	default:
		return DisplayBlanksAsGap
	}
}

// displayBlanksAsXML renders <c:dispBlanksAs>, which sits after
// <c:plotVisOnly> in a chart (ECMA-376 §21.2.2.42).
func displayBlanksAsXML(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return `
<c:dispBlanksAs val="` + normalizedDisplayBlanksAs(value) + `"/>`
}
