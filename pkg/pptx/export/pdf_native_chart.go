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
	// Centre the title horizontally over the chart rect.
	titleX := r.x + r.w/2 - float64(len(title))*3.5
	pdf.SetX(titleX)
	pdf.SetY(r.y + 4)
	_ = pdf.Cell(nil, title)
}
