//nolint:mnd // Chart text drawing uses fixed typographic and stroke constants.
package export

import (
	"github.com/signintech/gopdf"
)

// Chart text sizes, measured off PowerPoint's own render of a generated chart
// whose XML carries no c:txPr and no sz, so PowerPoint supplies its defaults.
//
// Axis and legend text came out at a 6.00pt cap height; Calibri's cap height is
// 1294/2048 em, which puts the size at 9.5pt -> 10pt. The chart title measured
// 53.25pt wide for "Bubble", which Calibri reproduces at 18pt.
//
// These matter because nothing in the chart renderers used to set a size at all:
// chart text inherited whatever the previously drawn element happened to leave
// on the canvas, so the same chart rendered differently depending on what sat
// next to it on the slide.
const (
	chartLabelFontSize = 10
	chartTitleFontSize = 18

	// chartAxisLabelGapPt is the clearance between an axis line and its labels.
	chartAxisLabelGapPt = 4.0

	// chartTickMarkPt is the length of a major tick mark. OOXML's default for
	// c:majorTickMark is "out", so PowerPoint draws a short stub outside the plot
	// at every major tick even when the chart XML says nothing about ticks.
	chartTickMarkPt = 3.0

	// chartGridlineGrey is the grey PowerPoint uses for major gridlines.
	chartGridlineGrey = 90

	// chartMinorTickMarkPt is the length of a minor tick, and
	// chartMinorTicksPerMajor how many intervals each major step is divided
	// into. PowerPoint's default minor unit is a fifth of the major one.
	chartMinorTickMarkPt    = 1.5
	chartMinorTicksPerMajor = 5

	// Legend marker geometry, sized to sit against 10pt text.
	chartLegendMarkerWPt = 8.0
	chartLegendMarkerHPt = 6.0
	chartLegendRowPt     = 14.0
)

// chartTextAlign selects how a chart label sits relative to the anchor point it
// is drawn at.
type chartTextAlign int

const (
	chartTextLeft chartTextAlign = iota
	chartTextCenter
	chartTextRight
)

// setChartFont selects the chart text face at size. Charts always use the sans
// alias: OOXML chart parts carry their own text properties rather than
// inheriting the shape's typeface, and gopptx does not emit any.
func setChartFont(pdf *gopdf.GoPdf, size int) {
	setPDFTextFontWithHint(pdf, size, false, false, "")
}

// drawChartLabel draws one piece of chart furniture at (x, y), where y is the
// vertical centre of the text rather than its top. Alignment is applied
// horizontally about x.
//
// Callers previously positioned this text with len(text)*3 style estimates,
// which assume a monospaced font and a fixed size; both are wrong here, so the
// real advance width is measured instead.
func drawChartLabel(pdf *gopdf.GoPdf, text string, x, y float64, size int, align chartTextAlign) {
	if text == "" {
		return
	}
	setChartFont(pdf, size)
	width := measuredWidth(pdf, text)
	drawX := x
	switch align {
	case chartTextCenter:
		drawX = x - width/2
	case chartTextRight:
		drawX = x - width
	case chartTextLeft:
	}
	pdf.SetX(drawX)
	pdf.SetY(y - pdfLineHeight(size)/2 + fontBaselineShift(pdf, "", size))
	_ = pdf.Cell(nil, text)
}

// chartAxisOrientation selects which side of the plot a tick mark sticks out on.
type chartAxisOrientation int

const (
	// chartAxisVertical is the value axis down the left of the plot; its ticks
	// point left.
	chartAxisVertical chartAxisOrientation = iota
	// chartAxisHorizontal is the axis along the bottom; its ticks point down.
	chartAxisHorizontal
)

// drawChartAxisTick draws one major tick mark outside the plot area at (x, y),
// which is the point where the tick meets the axis line.
func drawChartAxisTick(pdf *gopdf.GoPdf, x, y float64, orientation chartAxisOrientation) {
	drawChartAxisTickOfLength(pdf, x, y, orientation, chartTickMarkPt)
}

// drawChartMinorTicks draws the minor ticks between two neighbouring major ticks
// at axis positions from and to, which are measured along the axis.
func drawChartMinorTicks(pdf *gopdf.GoPdf, from, to, cross float64, orientation chartAxisOrientation) {
	step := (to - from) / chartMinorTicksPerMajor
	for i := 1; i < chartMinorTicksPerMajor; i++ {
		at := from + step*float64(i)
		if orientation == chartAxisVertical {
			drawChartAxisTickOfLength(pdf, cross, at, orientation, chartMinorTickMarkPt)
			continue
		}
		drawChartAxisTickOfLength(pdf, at, cross, orientation, chartMinorTickMarkPt)
	}
}

func drawChartAxisTickOfLength(
	pdf *gopdf.GoPdf,
	x, y float64,
	orientation chartAxisOrientation,
	length float64,
) {
	pdf.SetStrokeColor(30, 30, 30)
	if orientation == chartAxisVertical {
		pdf.Line(x-length, y, x, y)
		return
	}
	pdf.Line(x, y, x, y+length)
}

// chartLabelWidth is the advance width of text at the chart label size, for
// callers that need to reserve space before drawing.
func chartLabelWidth(pdf *gopdf.GoPdf, text string, size int) float64 {
	setChartFont(pdf, size)
	return measuredWidth(pdf, text)
}
