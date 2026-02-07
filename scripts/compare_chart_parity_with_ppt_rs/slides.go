package main

import "github.com/vegito/goppt/pkg/pptx"

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
