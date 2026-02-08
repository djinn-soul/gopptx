package pptx

// WithBarChart sets one bar chart for the slide.
func (s SlideContent) WithBarChart(chart BarChart) SlideContent {
	s.clearCharts()
	s.Chart = &chart
	return s
}

func (s SlideContent) WithBarHorizontalChart(chart BarHorizontalChart) SlideContent {
	s.clearCharts()
	s.BarHorizontal = &chart
	return s
}

func (s SlideContent) WithBarStackedChart(chart BarStackedChart) SlideContent {
	s.clearCharts()
	s.BarStacked = &chart
	return s
}

func (s SlideContent) WithBarStacked100Chart(chart BarStacked100Chart) SlideContent {
	s.clearCharts()
	s.BarStacked100 = &chart
	return s
}

// WithLineChart sets one line chart for the slide.
func (s SlideContent) WithLineChart(chart LineChart) SlideContent {
	s.clearCharts()
	s.Line = &chart
	return s
}

func (s SlideContent) WithLineMarkersChart(chart LineMarkersChart) SlideContent {
	s.clearCharts()
	s.LineMarkers = &chart
	return s
}

func (s SlideContent) WithLineStackedChart(chart LineStackedChart) SlideContent {
	s.clearCharts()
	s.LineStacked = &chart
	return s
}

// WithScatterChart sets one scatter chart for the slide.
func (s SlideContent) WithScatterChart(chart ScatterChart) SlideContent {
	s.clearCharts()
	s.Scatter = &chart
	return s
}

// WithAreaChart sets one area chart for the slide.
func (s SlideContent) WithAreaChart(chart AreaChart) SlideContent {
	s.clearCharts()
	s.Area = &chart
	return s
}

func (s SlideContent) WithAreaStackedChart(chart AreaStackedChart) SlideContent {
	s.clearCharts()
	s.AreaStacked = &chart
	return s
}

func (s SlideContent) WithAreaStacked100Chart(chart AreaStacked100Chart) SlideContent {
	s.clearCharts()
	s.AreaStacked100 = &chart
	return s
}

// WithPieChart sets one pie chart for the slide.
func (s SlideContent) WithPieChart(chart PieChart) SlideContent {
	s.clearCharts()
	s.Pie = &chart
	return s
}

// WithDoughnutChart sets one doughnut chart for the slide.
func (s SlideContent) WithDoughnutChart(chart DoughnutChart) SlideContent {
	s.clearCharts()
	s.Dough = &chart
	return s
}

func (s SlideContent) WithBubbleChart(chart BubbleChart) SlideContent {
	s.clearCharts()
	s.Bubble = &chart
	return s
}

func (s SlideContent) WithRadarChart(chart RadarChart) SlideContent {
	s.clearCharts()
	s.Radar = &chart
	return s
}

func (s SlideContent) WithRadarFilledChart(chart RadarFilledChart) SlideContent {
	s.clearCharts()
	s.RadarFilled = &chart
	return s
}

func (s SlideContent) WithStockHLCChart(chart StockHLCChart) SlideContent {
	s.clearCharts()
	s.StockHLC = &chart
	return s
}

func (s SlideContent) WithStockOHLCChart(chart StockOHLCChart) SlideContent {
	s.clearCharts()
	s.StockOHLC = &chart
	return s
}

func (s SlideContent) WithComboChart(chart ComboChart) SlideContent {
	s.clearCharts()
	s.Combo = &chart
	return s
}

func (s *SlideContent) clearCharts() {
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
	s.Dough = nil
	s.Bubble = nil
	s.Radar = nil
	s.RadarFilled = nil
	s.StockHLC = nil
	s.StockOHLC = nil
	s.Combo = nil
}
