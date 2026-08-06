package export

import (
	"math"

	"github.com/signintech/gopdf"
)

// Chart layout is circular: the plot area is inset by the axis labels, but how
// many labels an axis draws depends on how long the axis is. PowerPoint resolves
// this by fitting as many "nice" ticks as the axis has room for, so a wide plot
// gets 0, 0.5, 1, 1.5 … where a narrow one gets 0, 1, 2.
//
// solveChartLayout breaks the cycle with two passes: size the plot from a
// first-guess tick set, then re-derive the ticks from that size and re-inset.
// The densities it returns are the ones its final labels were measured from, so
// the frame drawers put ticks exactly where the insets reserved room for them.

// chartAxisSpec describes one axis for layout purposes. An axis is either a
// value axis (numeric ticks across MinV..MaxV) or a category axis (one label per
// category, no tick fitting).
type chartAxisSpec struct {
	MinV, MaxV  float64
	ValueFormat string
	// Categories, when non-empty, makes this a category axis.
	Categories []string
}

// labels returns what this axis draws at the given density.
func (a chartAxisSpec) labels(density chartAxisTickDensity) []string {
	if len(a.Categories) > 0 {
		return a.Categories
	}
	return chartAxisTickLabels(a.MinV, a.MaxV, a.ValueFormat, density)
}

// chartPlotLayout is a resolved plot area plus the tick densities the frame
// drawers must use to stay consistent with it.
type chartPlotLayout struct {
	X, Y, W, H float64
	Vertical   chartAxisTickDensity
	Horizontal chartAxisTickDensity
}

// chartLabelClearancePt is the gap kept between neighbouring axis labels before
// dropping to a coarser tick interval.
//
// It used to be 60pt — six times the label size — as a way of forcing coarse
// intervals while the axis *maximum* was still wrong: more gridlines only
// multiplied that misalignment. With the maximum and the step now solved
// together against PowerPoint's own ten-interval rule (see niceAxisMax), the
// clearance is back to what it says it is: enough room that two labels do not
// touch. Measured against PowerPoint, a 10pt vertical axis labels at a 34.75pt
// pitch — roughly a label height plus this gap.
const chartLabelClearancePt = 8.0

// solveChartLayout resolves the plot area and the tick density of both axes.
func solveChartLayout(
	pdf *gopdf.GoPdf,
	r chartRect,
	titleOverlay bool,
	title string,
	vertical, horizontal chartAxisSpec,
) chartPlotLayout {
	// Pass 1: no density known, so the tick intervals fall back to niceStep and
	// give a first estimate of how much room the labels need.
	titleHeight := chartTitleHeight(pdf, title, r.w)
	_, _, w, h := chartPlotInsets(
		pdf, r, titleOverlay, titleHeight,
		vertical.labels(chartAxisTickDensity{}),
		horizontal.labels(chartAxisTickDensity{}),
	)

	// Pass 2: with a provisional plot size, work out how many labels each axis
	// can actually carry, then re-inset for the labels that produces.
	verticalDensity := vertical.density(pdf, h, true)
	horizontalDensity := horizontal.density(pdf, w, false)
	px, py, pw, ph := chartPlotInsets(
		pdf, r, titleOverlay, titleHeight,
		vertical.labels(verticalDensity),
		horizontal.labels(horizontalDensity),
	)

	return chartPlotLayout{
		X: px, Y: py, W: pw, H: ph,
		Vertical:   verticalDensity,
		Horizontal: horizontalDensity,
	}
}

// density reports how many labels fit along an axis of the given length. A
// category axis returns the zero value: its ticks are fixed by the data.
func (a chartAxisSpec) density(pdf *gopdf.GoPdf, axisLengthPt float64, isVertical bool) chartAxisTickDensity {
	if len(a.Categories) > 0 || pdf == nil || axisLengthPt <= 0 {
		return chartAxisTickDensity{}
	}
	extent := pdfLineHeight(chartLabelFontSize) + chartLabelClearancePt
	if !isVertical {
		// Along a horizontal axis it is the label's width that limits how many
		// fit, so measure the widest one the axis might draw.
		widest := widestChartLabel(pdf, a.labels(chartAxisTickDensity{}))
		extent = widest + chartLabelClearancePt
	}
	return chartAxisTickDensity{AxisLengthPt: axisLengthPt, LabelExtentPt: extent}
}

// chartPlotInsets sizes the plot area from the labels that must sit around it.
// The widest vertical-axis label sets the left inset; the horizontal-axis labels
// set the bottom inset plus a right inset wide enough that the last label is not
// clipped.
//
// pdf may be nil, in which case fixed fallbacks are used.
func chartPlotInsets(
	pdf *gopdf.GoPdf,
	r chartRect,
	titleOverlay bool,
	titleHeight float64,
	verticalAxisLabels []string,
	horizontalAxisLabels []string,
) (float64, float64, float64, float64) {
	leftPad := chartFallbackLeftPadPt
	rightPad := chartFallbackRightPadPt
	bottomPad := chartFallbackBottomPadPt

	if pdf != nil {
		if widest := widestChartLabel(pdf, verticalAxisLabels); widest > 0 {
			// Floor it: PowerPoint keeps a minimum gutter beside the plot even
			// when the labels are a single digit, so sizing purely to the text
			// pulls the plot too far left.
			leftPad = math.Max(
				chartMinLeftPadPt,
				widest+chartAxisLabelGapPt+chartTickMarkPt,
			)
		}
		if widest := widestChartLabel(pdf, horizontalAxisLabels); widest > 0 {
			// Half the last label overhangs the axis end, so reserve that much.
			rightPad = math.Max(chartFallbackRightPadPt, widest/2+chartAxisLabelGapPt)
			bottomPad = pdfLineHeight(chartLabelFontSize) +
				chartAxisLabelGapPt + chartTickMarkPt + chartCategoryLabelPadPt
		}
	}

	topPad := chartMinTopPadPt
	if titleHeight > 0 {
		// Reserve the title's real height plus the margin PowerPoint leaves under
		// it, so a title that wrapped to two lines is not drawn over the plot.
		// Measured: a 3.75in chart with a one-line 18pt title puts its plot top
		// 39.75pt below the chart top, which is the 21.6pt title plus ~18pt.
		topPad = math.Max(topPad, titleHeight+chartTitleMarginPt)
	}
	if titleOverlay {
		topPad = chartOverlayTopPadPt
	}
	return r.x + leftPad, r.y + topPad, r.w - leftPad - rightPad, r.h - topPad - bottomPad
}

const (
	// chartMinTopPadPt is the gap PowerPoint leaves above the plot area.
	// Measured at 40pt on both a 360pt chart with a one-line title and a 468pt
	// chart with none, which is why it is a constant rather than the fraction of
	// the chart height this used to scale by: at 12% a full-slide chart reserved
	// 56pt and pushed its plot 16pt too low.
	chartMinTopPadPt = 40.0

	// chartTitleMarginPt is the gap PowerPoint leaves between the chart title and
	// the top of the plot area.
	chartTitleMarginPt = 8.0

	// chartCategoryLabelPadPt is the clearance PowerPoint keeps below the
	// category labels. Measured: a 25pt bottom inset against the 19pt the label
	// line box, tick and gap account for.
	chartCategoryLabelPadPt = 6.0
)
