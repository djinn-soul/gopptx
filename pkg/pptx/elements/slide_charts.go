package elements

import "github.com/djinn-soul/gopptx/pkg/pptx/charts"

// WithBarChart sets one bar chart for the slide.
func (s SlideContent) WithBarChart(chart charts.BarChart) SlideContent {
	clearCharts(&s)
	s.Chart = &chart
	return s
}

func (s SlideContent) WithBarHorizontalChart(chart charts.BarHorizontalChart) SlideContent {
	clearCharts(&s)
	s.BarHorizontal = &chart
	return s
}

func (s SlideContent) WithBarStackedChart(chart charts.BarStackedChart) SlideContent {
	clearCharts(&s)
	s.BarStacked = &chart
	return s
}

func (s SlideContent) WithBarStacked100Chart(chart charts.BarStacked100Chart) SlideContent {
	clearCharts(&s)
	s.BarStacked100 = &chart
	return s
}

// WithLineChart sets one line chart for the slide.
func (s SlideContent) WithLineChart(chart charts.LineChart) SlideContent {
	clearCharts(&s)
	s.Line = &chart
	return s
}

func (s SlideContent) WithLineMarkersChart(chart charts.LineMarkersChart) SlideContent {
	clearCharts(&s)
	s.LineMarkers = &chart
	return s
}

func (s SlideContent) WithLineStackedChart(chart charts.LineStackedChart) SlideContent {
	clearCharts(&s)
	s.LineStacked = &chart
	return s
}

// WithScatterChart sets one scatter chart for the slide.
func (s SlideContent) WithScatterChart(chart charts.ScatterChart) SlideContent {
	clearCharts(&s)
	s.Scatter = &chart
	return s
}

// WithAreaChart sets one area chart for the slide.
func (s SlideContent) WithAreaChart(chart charts.AreaChart) SlideContent {
	clearCharts(&s)
	s.Area = &chart
	return s
}

func (s SlideContent) WithAreaStackedChart(chart charts.AreaStackedChart) SlideContent {
	clearCharts(&s)
	s.AreaStacked = &chart
	return s
}

func (s SlideContent) WithAreaStacked100Chart(chart charts.AreaStacked100Chart) SlideContent {
	clearCharts(&s)
	s.AreaStacked100 = &chart
	return s
}

// WithPieChart sets one pie chart for the slide.
func (s SlideContent) WithPieChart(chart charts.PieChart) SlideContent {
	clearCharts(&s)
	s.Pie = &chart
	return s
}

// WithDoughnutChart sets one doughnut chart for the slide.
func (s SlideContent) WithDoughnutChart(chart charts.DoughnutChart) SlideContent {
	clearCharts(&s)
	s.Doughnut = &chart
	return s
}

func (s SlideContent) WithBubbleChart(chart charts.BubbleChart) SlideContent {
	clearCharts(&s)
	s.Bubble = &chart
	return s
}

func (s SlideContent) WithRadarChart(chart charts.RadarChart) SlideContent {
	clearCharts(&s)
	s.Radar = &chart
	return s
}

func (s SlideContent) WithRadarFilledChart(chart charts.RadarFilledChart) SlideContent {
	clearCharts(&s)
	s.RadarFilled = &chart
	return s
}

func (s SlideContent) WithStockHLCChart(chart charts.StockHLCChart) SlideContent {
	clearCharts(&s)
	s.StockHLC = &chart
	return s
}

func (s SlideContent) WithStockOHLCChart(chart charts.StockOHLCChart) SlideContent {
	clearCharts(&s)
	s.StockOHLC = &chart
	return s
}

func (s SlideContent) WithComboChart(chart charts.ComboChart) SlideContent {
	clearCharts(&s)
	s.Combo = &chart
	return s
}

func clearCharts(s *SlideContent) {
	s.Chart = nil
	s.BarHorizontal = nil
	s.BarStacked = nil
	s.BarStacked100 = nil
	s.Line = nil
	s.LineMarkers = nil
	s.LineStacked = nil
	s.Scatter = nil
	s.Area = nil
	s.AreaStacked = nil
	s.AreaStacked100 = nil
	s.Pie = nil
	s.Doughnut = nil
	s.Bubble = nil
	s.Radar = nil
	s.RadarFilled = nil
	s.StockHLC = nil
	s.StockOHLC = nil
	s.Combo = nil
}
