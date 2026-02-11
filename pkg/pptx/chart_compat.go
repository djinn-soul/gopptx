package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

// Chart aliases for backward compatibility
type (
	BarChart            = charts.BarChart
	BarHorizontalChart  = charts.BarHorizontalChart
	BarStackedChart     = charts.BarStackedChart
	BarStacked100Chart  = charts.BarStacked100Chart
	LineChart           = charts.LineChart
	LineMarkersChart    = charts.LineMarkersChart
	LineStackedChart    = charts.LineStackedChart
	ScatterChart        = charts.ScatterChart
	AreaChart           = charts.AreaChart
	AreaStackedChart    = charts.AreaStackedChart
	AreaStacked100Chart = charts.AreaStacked100Chart
	PieChart            = charts.PieChart
	DoughnutChart       = charts.DoughnutChart
	BubbleChart         = charts.BubbleChart
	RadarChart          = charts.RadarChart
	RadarFilledChart    = charts.RadarFilledChart
	StockHLCChart       = charts.StockHLCChart
	StockOHLCChart      = charts.StockOHLCChart
	ComboChart          = charts.ComboChart
	Series              = charts.Series
)

type (
	ChartDefinition = charts.ChartDefinition
)

func NewBarChart(categories []string, values []float64) BarChart {
	return charts.NewBarChart(categories, values)
}

func NewBarHorizontalChart(categories []string, values []float64) BarHorizontalChart {
	return charts.NewBarHorizontalChart(categories, values)
}

func NewBarStackedChart(categories []string, values []float64) BarStackedChart {
	return charts.NewBarStackedChart(categories, values)
}

func NewBarStacked100Chart(categories []string, values []float64) BarStacked100Chart {
	return charts.NewBarStacked100Chart(categories, values)
}

func NewLineChart(categories []string, values []float64) LineChart {
	return charts.NewLineChart(categories, values)
}

func NewLineMarkersChart(categories []string, values []float64) LineMarkersChart {
	return charts.NewLineMarkersChart(categories, values)
}

func NewLineStackedChart(categories []string, values []float64) LineStackedChart {
	return charts.NewLineStackedChart(categories, values)
}

func NewScatterChart(xValues, yValues []float64) ScatterChart {
	return charts.NewScatterChart(xValues, yValues)
}

func NewAreaChart(categories []string, values []float64) AreaChart {
	return charts.NewAreaChart(categories, values)
}

func NewAreaStackedChart(categories []string, values []float64) AreaStackedChart {
	return charts.NewAreaStackedChart(categories, values)
}

func NewAreaStacked100Chart(categories []string, values []float64) AreaStacked100Chart {
	return charts.NewAreaStacked100Chart(categories, values)
}

func NewPieChart(categories []string, values []float64) PieChart {
	return charts.NewPieChart(categories, values)
}

func NewDoughnutChart(categories []string, values []float64) DoughnutChart {
	return charts.NewDoughnutChart(categories, values)
}

func NewBubbleChart(xValues, yValues, sizeValues []float64) BubbleChart {
	return charts.NewBubbleChart(xValues, yValues, sizeValues)
}

func NewRadarChart(categories []string, values []float64) RadarChart {
	return charts.NewRadarChart(categories, values)
}

func NewRadarFilledChart(categories []string, values []float64) RadarFilledChart {
	return charts.NewRadarFilledChart(categories, values)
}

func NewStockHLCChart(categories []string, high, low, close []float64) StockHLCChart {
	return charts.NewStockHLCChart(categories, high, low, close)
}

func NewStockOHLCChart(categories []string, open, high, low, close []float64) StockOHLCChart {
	return charts.NewStockOHLCChart(categories, open, high, low, close)
}

func NewComboChart(categories []string, barSeries []Series, lineSeries []Series) ComboChart {
	return charts.NewComboChart(categories, barSeries, lineSeries)
}

// Constants
const (
	LegendPositionRight  = charts.LegendPositionRight
	LegendPositionLeft   = charts.LegendPositionLeft
	LegendPositionTop    = charts.LegendPositionTop
	LegendPositionBottom = charts.LegendPositionBottom

	ValueAxisCrossBetweenBetween     = charts.ValueAxisCrossBetweenBetween
	ValueAxisCrossBetweenMidCategory = charts.ValueAxisCrossBetweenMidCategory

	RadarStyleFilled = charts.RadarStyleFilled

	ScatterStyleMarker       = charts.ScatterStyleMarker
	ScatterStyleLineMarker   = charts.ScatterStyleLineMarker
	ScatterStyleSmoothMarker = charts.ScatterStyleSmoothMarker
)
