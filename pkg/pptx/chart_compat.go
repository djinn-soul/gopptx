package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

// BarChart and related types are chart aliases for backward compatibility.
type (
	// BarChart is an alias for charts.BarChart.
	BarChart = charts.BarChart
	// BarHorizontalChart is an alias for charts.BarHorizontalChart.
	BarHorizontalChart = charts.BarHorizontalChart
	// BarStackedChart is an alias for charts.BarStackedChart.
	BarStackedChart = charts.BarStackedChart
	// BarStacked100Chart is an alias for charts.BarStacked100Chart.
	BarStacked100Chart = charts.BarStacked100Chart
	// LineChart is an alias for charts.LineChart.
	LineChart = charts.LineChart
	// LineMarkersChart is an alias for charts.LineMarkersChart.
	LineMarkersChart = charts.LineMarkersChart
	// LineStackedChart is an alias for charts.LineStackedChart.
	LineStackedChart = charts.LineStackedChart
	// ScatterChart is an alias for charts.ScatterChart.
	ScatterChart = charts.ScatterChart
	// AreaChart is an alias for charts.AreaChart.
	AreaChart = charts.AreaChart
	// AreaStackedChart is an alias for charts.AreaStackedChart.
	AreaStackedChart = charts.AreaStackedChart
	// AreaStacked100Chart is an alias for charts.AreaStacked100Chart.
	AreaStacked100Chart = charts.AreaStacked100Chart
	// PieChart is an alias for charts.PieChart.
	PieChart = charts.PieChart
	// DoughnutChart is an alias for charts.DoughnutChart.
	DoughnutChart = charts.DoughnutChart
	// BubbleChart is an alias for charts.BubbleChart.
	BubbleChart = charts.BubbleChart
	// RadarChart is an alias for charts.RadarChart.
	RadarChart = charts.RadarChart
	// RadarFilledChart is an alias for charts.RadarFilledChart.
	RadarFilledChart = charts.RadarFilledChart
	// StockHLCChart is an alias for charts.StockHLCChart.
	StockHLCChart = charts.StockHLCChart
	// StockOHLCChart is an alias for charts.StockOHLCChart.
	StockOHLCChart = charts.StockOHLCChart
	// ComboChart is an alias for charts.ComboChart.
	ComboChart = charts.ComboChart
	// Series is an alias for charts.Series.
	Series = charts.Series
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

func NewStockHLCChart(categories []string, high, low, closeValues []float64) StockHLCChart {
	return charts.NewStockHLCChart(categories, high, low, closeValues)
}

func NewStockOHLCChart(categories []string, open, high, low, closeValues []float64) StockOHLCChart {
	return charts.NewStockOHLCChart(categories, open, high, low, closeValues)
}

func NewComboChart(categories []string, barSeries []Series, lineSeries []Series) ComboChart {
	return charts.NewComboChart(categories, barSeries, lineSeries)
}

// Constants.
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
