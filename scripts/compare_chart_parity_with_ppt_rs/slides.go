package main

import "github.com/djinn09/goppt/pkg/pptx"

func barSlide() pptx.SlideContent {
	chart := pptx.NewBarChart(
		[]string{"Q1", "Q2", "Q3"},
		[]float64{12, 18, 24},
	).WithTitle("Bar")
	return pptx.NewSlide("Bar").WithBarChart(chart)
}

func barHorizontalSlide() pptx.SlideContent {
	chart := pptx.NewBarHorizontalChart(
		[]string{"Q1", "Q2", "Q3"},
		[]float64{12, 18, 24},
	).WithTitle("Bar Horizontal")
	return pptx.NewSlide("Bar Horizontal").WithBarHorizontalChart(chart)
}

func barStackedSlide() pptx.SlideContent {
	chart := pptx.NewBarStackedChart(
		[]string{"Q1", "Q2", "Q3"},
		[]float64{12, 18, 24},
	).WithTitle("Bar Stacked")
	return pptx.NewSlide("Bar Stacked").WithBarStackedChart(chart)
}

func barStacked100Slide() pptx.SlideContent {
	chart := pptx.NewBarStacked100Chart(
		[]string{"Q1", "Q2", "Q3"},
		[]float64{12, 18, 24},
	).WithTitle("Bar Stacked 100")
	return pptx.NewSlide("Bar Stacked 100").WithBarStacked100Chart(chart)
}

func lineSlide() pptx.SlideContent {
	chart := pptx.NewLineChart(
		[]string{"Q1", "Q2", "Q3"},
		[]float64{10, 16, 22},
	).WithTitle("Line")
	return pptx.NewSlide("Line").WithLineChart(chart)
}

func lineMarkersSlide() pptx.SlideContent {
	chart := pptx.NewLineMarkersChart(
		[]string{"Q1", "Q2", "Q3"},
		[]float64{10, 16, 22},
	).WithTitle("Line Markers")
	return pptx.NewSlide("Line Markers").WithLineMarkersChart(chart)
}

func lineStackedSlide() pptx.SlideContent {
	chart := pptx.NewLineStackedChart(
		[]string{"Q1", "Q2", "Q3"},
		[]float64{10, 16, 22},
	).WithTitle("Line Stacked")
	return pptx.NewSlide("Line Stacked").WithLineStackedChart(chart)
}

func areaSlide() pptx.SlideContent {
	chart := pptx.NewAreaChart(
		[]string{"Q1", "Q2", "Q3"},
		[]float64{14, 17, 23},
	).WithTitle("Area")
	return pptx.NewSlide("Area").WithAreaChart(chart)
}

func areaStackedSlide() pptx.SlideContent {
	chart := pptx.NewAreaStackedChart(
		[]string{"Q1", "Q2", "Q3"},
		[]float64{14, 17, 23},
	).WithTitle("Area Stacked")
	return pptx.NewSlide("Area Stacked").WithAreaStackedChart(chart)
}

func areaStacked100Slide() pptx.SlideContent {
	chart := pptx.NewAreaStacked100Chart(
		[]string{"Q1", "Q2", "Q3"},
		[]float64{14, 17, 23},
	).WithTitle("Area Stacked 100")
	return pptx.NewSlide("Area Stacked 100").WithAreaStacked100Chart(chart)
}

func pieSlide() pptx.SlideContent {
	chart := pptx.NewPieChart(
		[]string{"A", "B", "C"},
		[]float64{30, 45, 25},
	).WithTitle("Pie")
	return pptx.NewSlide("Pie").WithPieChart(chart)
}

func doughnutSlide() pptx.SlideContent {
	chart := pptx.NewDoughnutChart(
		[]string{"A", "B", "C"},
		[]float64{30, 45, 25},
	).WithTitle("Doughnut")
	return pptx.NewSlide("Doughnut").WithDoughnutChart(chart)
}

func scatterMarkerSlide() pptx.SlideContent {
	chart := pptx.NewScatterChart(
		[]float64{1, 2, 3},
		[]float64{10, 15, 20},
	).WithTitle("Scatter Marker").WithScatterStyle(pptx.ScatterStyleMarker)
	return pptx.NewSlide("Scatter Marker").WithScatterChart(chart)
}

func scatterLinesSlide() pptx.SlideContent {
	chart := pptx.NewScatterChart(
		[]float64{1, 2, 3},
		[]float64{10, 15, 20},
	).WithTitle("Scatter Lines").WithScatterStyle(pptx.ScatterStyleLineMarker)
	return pptx.NewSlide("Scatter Lines").WithScatterChart(chart)
}

func scatterSmoothSlide() pptx.SlideContent {
	chart := pptx.NewScatterChart(
		[]float64{1, 2, 3},
		[]float64{10, 15, 20},
	).WithTitle("Scatter Smooth").WithScatterStyle(pptx.ScatterStyleSmoothMarker)
	return pptx.NewSlide("Scatter Smooth").WithScatterChart(chart)
}

func bubbleSlide() pptx.SlideContent {
	chart := pptx.NewBubbleChart(
		[]float64{1, 2, 3},
		[]float64{10, 20, 30},
		[]float64{10, 20, 30},
	).WithTitle("Bubble").WithSeriesName("Series 1").WithBubbleScale(100)
	return pptx.NewSlide("Bubble").WithBubbleChart(chart)
}

func radarSlide() pptx.SlideContent {
	chart := pptx.NewRadarChart(
		[]string{"A", "B", "C"},
		[]float64{2, 3, 4},
	).WithTitle("Radar")
	return pptx.NewSlide("Radar").WithRadarChart(chart)
}

func radarFilledSlide() pptx.SlideContent {
	chart := pptx.NewRadarFilledChart(
		[]string{"A", "B", "C"},
		[]float64{3, 4, 5},
	).WithTitle("Radar Filled")
	return pptx.NewSlide("Radar Filled").WithRadarFilledChart(chart)
}

func stockHLCSlide() pptx.SlideContent {
	chart := pptx.NewStockHLCChart(
		[]string{"D1", "D2", "D3"},
		[]float64{12, 13, 14},
		[]float64{8, 9, 10},
		[]float64{10, 11, 12},
	).WithTitle("StockHLC")
	return pptx.NewSlide("StockHLC").WithStockHLCChart(chart)
}

func stockOHLCSlide() pptx.SlideContent {
	chart := pptx.NewStockOHLCChart(
		[]string{"D1", "D2", "D3"},
		[]float64{9, 10, 11},
		[]float64{12, 13, 14},
		[]float64{8, 9, 10},
		[]float64{10, 11, 12},
	).WithTitle("StockOHLC")
	return pptx.NewSlide("StockOHLC").WithStockOHLCChart(chart)
}

func comboSlide() pptx.SlideContent {
	chart := pptx.NewComboChart(
		[]string{"Q1", "Q2", "Q3"},
		[]pptx.Series{{Name: "Bar A", Values: []float64{1, 2, 3}}},
		[]pptx.Series{{Name: "Line A", Values: []float64{2, 3, 4}}},
	).WithTitle("Combo")
	return pptx.NewSlide("Combo").WithComboChart(chart)
}
