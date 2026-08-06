package export

import (
	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// addChartPaintItems queues every chart the slide carries against its position
// in the shape tree.
//
// Each chart kind lives in its own SlideContent field, so the kinds are listed
// here rather than iterated; the order of the list is what a slide with no
// recovered tree order falls back to, and it matches the order the renderer
// used before paint order existed.
func addChartPaintItems(list *paintList, pdf *gopdf.GoPdf, slide elements.SlideContent, order slidePaintOrder) {
	for _, entry := range chartPaintEntries(slide) {
		draw := entry.draw
		list.add(order.chartZ(entry.kind), func() error { draw(pdf); return nil })
	}
}

// chartPaintEntry pairs the reader's chart kind with the call that draws it.
type chartPaintEntry struct {
	kind string
	draw func(*gopdf.GoPdf)
}

func chartPaintEntries(slide elements.SlideContent) []chartPaintEntry {
	var out []chartPaintEntry
	add := func(kind string, present bool, draw func(*gopdf.GoPdf)) {
		if present {
			out = append(out, chartPaintEntry{kind: kind, draw: draw})
		}
	}
	add(chartKindBar, slide.Chart != nil, func(p *gopdf.GoPdf) { renderBarChart(p, slide.Chart) })
	add(chartKindBarHoriz, slide.BarHorizontal != nil, func(p *gopdf.GoPdf) {
		renderBarHorizontalChart(p, slide.BarHorizontal)
	})
	add(chartKindBarStacked, slide.BarStacked != nil, func(p *gopdf.GoPdf) {
		renderBarStackedChart(p, slide.BarStacked)
	})
	add(chartKindBar100, slide.BarStacked100 != nil, func(p *gopdf.GoPdf) {
		renderBarStacked100Chart(p, slide.BarStacked100)
	})
	add(chartKindLine, slide.Line != nil, func(p *gopdf.GoPdf) { renderLineChart(p, slide.Line) })
	add(chartKindLineMarkers, slide.LineMarkers != nil, func(p *gopdf.GoPdf) {
		renderLineMarkersChart(p, slide.LineMarkers)
	})
	add(chartKindLineStacked, slide.LineStacked != nil, func(p *gopdf.GoPdf) {
		renderLineStackedChart(p, slide.LineStacked)
	})
	add(chartKindArea, slide.Area != nil, func(p *gopdf.GoPdf) { renderAreaChart(p, slide.Area) })
	add(chartKindAreaStacked, slide.AreaStacked != nil, func(p *gopdf.GoPdf) {
		renderAreaStackedChart(p, slide.AreaStacked)
	})
	add(chartKindArea100, slide.AreaStacked100 != nil, func(p *gopdf.GoPdf) {
		renderAreaStacked100Chart(p, slide.AreaStacked100)
	})
	add(chartKindPie, slide.Pie != nil, func(p *gopdf.GoPdf) { renderPieChart(p, slide.Pie) })
	add(chartKindDoughnut, slide.Doughnut != nil, func(p *gopdf.GoPdf) { renderDoughnutChart(p, slide.Doughnut) })
	add(chartKindScatter, slide.Scatter != nil, func(p *gopdf.GoPdf) { renderScatterChart(p, slide.Scatter) })
	add(chartKindBubble, slide.Bubble != nil, func(p *gopdf.GoPdf) { renderBubbleChart(p, slide.Bubble) })
	add(chartKindRadar, slide.Radar != nil, func(p *gopdf.GoPdf) { renderRadarChart(p, slide.Radar) })
	add(chartKindRadarFilled, slide.RadarFilled != nil, func(p *gopdf.GoPdf) {
		renderRadarFilledChart(p, slide.RadarFilled)
	})
	add(chartKindStockHLC, slide.StockHLC != nil, func(p *gopdf.GoPdf) { renderStockHLCChart(p, slide.StockHLC) })
	add(chartKindStockOHLC, slide.StockOHLC != nil, func(p *gopdf.GoPdf) { renderStockOHLCChart(p, slide.StockOHLC) })
	add(chartKindCombo, slide.Combo != nil, func(p *gopdf.GoPdf) { renderComboChart(p, slide.Combo) })
	return out
}
