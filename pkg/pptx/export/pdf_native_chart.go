//nolint:mnd // Native chart title rendering uses fixed visual offsets from PPT defaults.
package export

import (
	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type chartRect struct{ x, y, w, h float64 }

//nolint:funlen // Chart rendering enumerates all supported chart fields explicitly for deterministic output.
func renderNativePDFSlideCharts(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	if slide.Chart != nil {
		renderBarLike(
			pdf,
			slide.Chart.Title,
			chartRectFromLength(slide.Chart.X.Emu(), slide.Chart.Y.Emu(), slide.Chart.CX.Emu(), slide.Chart.CY.Emu()),
			slide.Chart.Values,
			slide.Chart.Categories,
			false,
		)
	}
	if slide.BarHorizontal != nil {
		renderBarLike(
			pdf,
			slide.BarHorizontal.Title,
			chartRectFromLength(
				slide.BarHorizontal.X.Emu(),
				slide.BarHorizontal.Y.Emu(),
				slide.BarHorizontal.CX.Emu(),
				slide.BarHorizontal.CY.Emu(),
			),
			slide.BarHorizontal.Values,
			slide.BarHorizontal.Categories,
			true,
		)
	}
	if slide.BarStacked != nil {
		renderBarLike(
			pdf,
			slide.BarStacked.Title,
			chartRectFromLength(
				slide.BarStacked.X.Emu(),
				slide.BarStacked.Y.Emu(),
				slide.BarStacked.CX.Emu(),
				slide.BarStacked.CY.Emu(),
			),
			slide.BarStacked.Values,
			slide.BarStacked.Categories,
			true,
		)
	}
	if slide.BarStacked100 != nil {
		renderBarLike(
			pdf,
			slide.BarStacked100.Title,
			chartRectFromLength(
				slide.BarStacked100.X.Emu(),
				slide.BarStacked100.Y.Emu(),
				slide.BarStacked100.CX.Emu(),
				slide.BarStacked100.CY.Emu(),
			),
			normalizePercentSeries(slide.BarStacked100.Values),
			slide.BarStacked100.Categories,
			true,
		)
	}
	if slide.Line != nil {
		renderLineLike(
			pdf,
			slide.Line.Title,
			chartRectFromLength(slide.Line.X.Emu(), slide.Line.Y.Emu(), slide.Line.CX.Emu(), slide.Line.CY.Emu()),
			slide.Line.Values,
			slide.Line.Categories,
			false,
		)
	}
	if slide.LineMarkers != nil {
		renderLineLike(
			pdf,
			slide.LineMarkers.Title,
			chartRectFromLength(
				slide.LineMarkers.X.Emu(),
				slide.LineMarkers.Y.Emu(),
				slide.LineMarkers.CX.Emu(),
				slide.LineMarkers.CY.Emu(),
			),
			slide.LineMarkers.Values,
			slide.LineMarkers.Categories,
			true,
		)
	}
	if slide.LineStacked != nil {
		renderLineLike(
			pdf,
			slide.LineStacked.Title,
			chartRectFromLength(
				slide.LineStacked.X.Emu(),
				slide.LineStacked.Y.Emu(),
				slide.LineStacked.CX.Emu(),
				slide.LineStacked.CY.Emu(),
			),
			slide.LineStacked.Values,
			slide.LineStacked.Categories,
			true,
		)
	}
	if slide.Area != nil {
		renderAreaLike(
			pdf,
			slide.Area.Title,
			chartRectFromLength(slide.Area.X.Emu(), slide.Area.Y.Emu(), slide.Area.CX.Emu(), slide.Area.CY.Emu()),
			slide.Area.Values,
			slide.Area.Categories,
		)
	}
	if slide.AreaStacked != nil {
		renderAreaLike(
			pdf,
			slide.AreaStacked.Title,
			chartRectFromLength(
				slide.AreaStacked.X.Emu(),
				slide.AreaStacked.Y.Emu(),
				slide.AreaStacked.CX.Emu(),
				slide.AreaStacked.CY.Emu(),
			),
			slide.AreaStacked.Values,
			slide.AreaStacked.Categories,
		)
	}
	if slide.AreaStacked100 != nil {
		renderAreaLike(
			pdf,
			slide.AreaStacked100.Title,
			chartRectFromLength(
				slide.AreaStacked100.X.Emu(),
				slide.AreaStacked100.Y.Emu(),
				slide.AreaStacked100.CX.Emu(),
				slide.AreaStacked100.CY.Emu(),
			),
			normalizePercentSeries(slide.AreaStacked100.Values),
			slide.AreaStacked100.Categories,
		)
	}
	if slide.Pie != nil {
		renderPieLike(
			pdf,
			slide.Pie.Title,
			chartRectFromLength(slide.Pie.X.Emu(), slide.Pie.Y.Emu(), slide.Pie.CX.Emu(), slide.Pie.CY.Emu()),
			slide.Pie.Values,
			false,
		)
	}
	if slide.Doughnut != nil {
		renderPieLike(
			pdf,
			slide.Doughnut.Title,
			chartRectFromLength(
				slide.Doughnut.X.Emu(),
				slide.Doughnut.Y.Emu(),
				slide.Doughnut.CX.Emu(),
				slide.Doughnut.CY.Emu(),
			),
			slide.Doughnut.Values,
			true,
		)
	}
	if slide.Scatter != nil {
		renderScatterLike(
			pdf,
			slide.Scatter.Title,
			chartRectFromLength(
				slide.Scatter.X.Emu(),
				slide.Scatter.Y.Emu(),
				slide.Scatter.CX.Emu(),
				slide.Scatter.CY.Emu(),
			),
			slide.Scatter.XValues,
			slide.Scatter.YValues,
			nil,
		)
	}
	if slide.Bubble != nil {
		renderScatterLike(
			pdf,
			slide.Bubble.Title,
			chartRectFromLength(slide.Bubble.X, slide.Bubble.Y, slide.Bubble.CX, slide.Bubble.CY),
			slide.Bubble.XValues,
			slide.Bubble.YValues,
			slide.Bubble.BubbleSizes,
		)
	}
	if slide.Radar != nil {
		renderRadarLike(
			pdf,
			slide.Radar.Title,
			chartRectFromLength(slide.Radar.X.Emu(), slide.Radar.Y.Emu(), slide.Radar.CX.Emu(), slide.Radar.CY.Emu()),
			slide.Radar.Values,
			false,
		)
	}
	if slide.RadarFilled != nil {
		renderRadarLike(
			pdf,
			slide.RadarFilled.Title,
			chartRectFromLength(
				slide.RadarFilled.X.Emu(),
				slide.RadarFilled.Y.Emu(),
				slide.RadarFilled.CX.Emu(),
				slide.RadarFilled.CY.Emu(),
			),
			slide.RadarFilled.Values,
			true,
		)
	}
	if slide.StockHLC != nil {
		renderStockLike(
			pdf,
			slide.StockHLC.Title,
			chartRectFromLength(
				slide.StockHLC.X.Emu(),
				slide.StockHLC.Y.Emu(),
				slide.StockHLC.CX.Emu(),
				slide.StockHLC.CY.Emu(),
			),
			nil,
			slide.StockHLC.HighValues,
			slide.StockHLC.LowValues,
			slide.StockHLC.CloseValues,
		)
	}
	if slide.StockOHLC != nil {
		renderStockLike(
			pdf,
			slide.StockOHLC.Title,
			chartRectFromLength(
				slide.StockOHLC.X.Emu(),
				slide.StockOHLC.Y.Emu(),
				slide.StockOHLC.CX.Emu(),
				slide.StockOHLC.CY.Emu(),
			),
			slide.StockOHLC.OpenValues,
			slide.StockOHLC.HighValues,
			slide.StockOHLC.LowValues,
			slide.StockOHLC.CloseValues,
		)
	}
	if slide.Combo != nil {
		renderComboLike(
			pdf,
			slide.Combo.Title,
			chartRectFromLength(slide.Combo.X.Emu(), slide.Combo.Y.Emu(), slide.Combo.CX.Emu(), slide.Combo.CY.Emu()),
			slide.Combo.BarSeries,
			slide.Combo.LineSeries,
			slide.Combo.Categories,
		)
	}
}

func chartRectFromLength(x, y, w, h int64) chartRect {
	return chartRect{emuToPt(x), emuToPt(y), emuToPt(w), emuToPt(h)}
}

func renderChartTitle(pdf *gopdf.GoPdf, title string, r chartRect) {
	if title == "" {
		return
	}
	pdf.SetTextColor(40, 40, 40)
	pdf.SetX(r.x + 6)
	pdf.SetY(r.y + 4)
	_ = pdf.Cell(nil, title)
}
