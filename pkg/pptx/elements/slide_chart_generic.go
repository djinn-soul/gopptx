package elements

import "github.com/djinn-soul/gopptx/pkg/pptx/charts"

// WithChart sets any chart on the slide, replacing a chart already set. It is
// the type-generic counterpart to the WithBarChart / WithPieChart / ... family
// and dispatches to the same fields, so the two styles are interchangeable.
//
// A nil chart leaves the slide unchanged.
func (s SlideContent) WithChart(chart charts.Chart) SlideContent {
	switch c := chart.(type) {
	case charts.BarChart:
		return s.WithBarChart(c)
	case charts.BarHorizontalChart:
		return s.WithBarHorizontalChart(c)
	case charts.BarStackedChart:
		return s.WithBarStackedChart(c)
	case charts.BarStacked100Chart:
		return s.WithBarStacked100Chart(c)
	case charts.LineChart:
		return s.WithLineChart(c)
	case charts.LineMarkersChart:
		return s.WithLineMarkersChart(c)
	case charts.LineStackedChart:
		return s.WithLineStackedChart(c)
	case charts.ScatterChart:
		return s.WithScatterChart(c)
	case charts.AreaChart:
		return s.WithAreaChart(c)
	case charts.AreaStackedChart:
		return s.WithAreaStackedChart(c)
	case charts.AreaStacked100Chart:
		return s.WithAreaStacked100Chart(c)
	case charts.PieChart:
		return s.WithPieChart(c)
	case charts.DoughnutChart:
		return s.WithDoughnutChart(c)
	case charts.BubbleChart:
		return s.WithBubbleChart(c)
	case charts.RadarChart:
		return s.WithRadarChart(c)
	case charts.RadarFilledChart:
		return s.WithRadarFilledChart(c)
	case charts.StockHLCChart:
		return s.WithStockHLCChart(c)
	case charts.StockOHLCChart:
		return s.WithStockOHLCChart(c)
	case charts.ComboChart:
		return s.WithComboChart(c)
	}
	return s
}
