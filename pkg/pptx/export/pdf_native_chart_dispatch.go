package export

import (
	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func renderBarAndLineCharts(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	if slide.Chart != nil {
		renderBarChart(pdf, slide.Chart)
	}
	if slide.BarHorizontal != nil {
		renderBarHorizontalChart(pdf, slide.BarHorizontal)
	}
	if slide.BarStacked != nil {
		renderBarStackedChart(pdf, slide.BarStacked)
	}
	if slide.BarStacked100 != nil {
		renderBarStacked100Chart(pdf, slide.BarStacked100)
	}
	if slide.Line != nil {
		renderLineChart(pdf, slide.Line)
	}
	if slide.LineMarkers != nil {
		renderLineMarkersChart(pdf, slide.LineMarkers)
	}
	if slide.LineStacked != nil {
		renderLineStackedChart(pdf, slide.LineStacked)
	}
	if slide.Area != nil {
		renderAreaChart(pdf, slide.Area)
	}
	if slide.AreaStacked != nil {
		renderAreaStackedChart(pdf, slide.AreaStacked)
	}
	if slide.AreaStacked100 != nil {
		renderAreaStacked100Chart(pdf, slide.AreaStacked100)
	}
}

func renderOtherCharts(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	if slide.Pie != nil {
		renderPieChart(pdf, slide.Pie)
	}
	if slide.Doughnut != nil {
		renderDoughnutChart(pdf, slide.Doughnut)
	}
	if slide.Scatter != nil {
		renderScatterChart(pdf, slide.Scatter)
	}
	if slide.Bubble != nil {
		renderBubbleChart(pdf, slide.Bubble)
	}
	if slide.Radar != nil {
		renderRadarChart(pdf, slide.Radar)
	}
	if slide.RadarFilled != nil {
		renderRadarFilledChart(pdf, slide.RadarFilled)
	}
	if slide.StockHLC != nil {
		renderStockHLCChart(pdf, slide.StockHLC)
	}
	if slide.StockOHLC != nil {
		renderStockOHLCChart(pdf, slide.StockOHLC)
	}
	if slide.Combo != nil {
		renderComboChart(pdf, slide.Combo)
	}
}
