//nolint:mnd // Native chart title rendering uses fixed visual offsets from PPT defaults.
package export

import (
	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type chartRect struct{ x, y, w, h float64 }

func renderNativePDFSlideCharts(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	renderBarAndLineCharts(pdf, slide)
	renderOtherCharts(pdf, slide)
}

func chartRectFromLength(x, y, w, h int64) chartRect {
	return chartRect{emuToPt(x), emuToPt(y), emuToPt(w), emuToPt(h)}
}

func renderChartTitle(pdf *gopdf.GoPdf, title string, r chartRect) {
	if title == "" {
		return
	}
	pdf.SetTextColor(40, 40, 40)
	// PowerPoint wraps a chart title that is wider than the chart rather than
	// letting it run past the edge, so measure and wrap here too.
	setChartFont(pdf, chartTitleFontSize)
	lines := wrapPDFTextWithMetrics(pdf, title, r.w)
	lineHeight := pdfLineHeight(chartTitleFontSize)
	// The anchor Y is each line's centre, so start half a line below the top.
	y := r.y + 4 + lineHeight/2
	for _, line := range lines {
		drawChartLabel(pdf, line, r.x+r.w/2, y, chartTitleFontSize, chartTextCenter)
		y += lineHeight
	}
	pdf.SetTextColor(0, 0, 0)
}

// chartTitleHeight is the vertical space a wrapped title occupies, so the plot
// area can be pushed below it instead of being drawn over it.
func chartTitleHeight(pdf *gopdf.GoPdf, title string, width float64) float64 {
	if title == "" {
		return 0
	}
	setChartFont(pdf, chartTitleFontSize)
	lines := wrapPDFTextWithMetrics(pdf, title, width)
	return float64(len(lines)) * pdfLineHeight(chartTitleFontSize)
}
