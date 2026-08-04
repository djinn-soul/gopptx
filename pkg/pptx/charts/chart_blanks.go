package charts

import "github.com/djinn-soul/gopptx/internal/pptxxml"

// How a chart draws a category that has no value.
const (
	// DisplayBlanksAsGap leaves a hole in the series, PowerPoint's default.
	DisplayBlanksAsGap = pptxxml.DisplayBlanksAsGap
	// DisplayBlanksAsZero plots the blank at zero.
	DisplayBlanksAsZero = pptxxml.DisplayBlanksAsZero
	// DisplayBlanksAsSpan joins the neighbouring points across the blank.
	DisplayBlanksAsSpan = pptxxml.DisplayBlanksAsSpan
)

// BlankValue is the value that marks a data point as missing rather than zero.
// A series holding it writes no <c:pt> for that category, so PowerPoint draws a
// gap instead of a reading of zero (upstream issue #968).
func BlankValue() float64 {
	return pptxxml.BlankValue()
}

// IsBlankValue reports whether a value stands for a missing data point.
func IsBlankValue(value float64) bool {
	return pptxxml.IsBlankValue(value)
}

// WithDisplayBlanksAs sets how blanks in this chart are drawn: gap, zero or
// span. An unrecognised token falls back to gap.
func (c LineChart) WithDisplayBlanksAs(mode string) LineChart {
	c.DisplayBlanksAs = mode
	return c
}

// WithDisplayBlanksAs sets how blanks in this chart are drawn.
func (c BarChart) WithDisplayBlanksAs(mode string) BarChart {
	c.DisplayBlanksAs = mode
	return c
}

// WithDisplayBlanksAs sets how blanks in this chart are drawn.
func (c AreaChart) WithDisplayBlanksAs(mode string) AreaChart {
	c.DisplayBlanksAs = mode
	return c
}
