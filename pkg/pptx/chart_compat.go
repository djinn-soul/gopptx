package pptx

import "github.com/djinn-soul/gopptx/pkg/pptx/charts"

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
)

// Other type aliases
type (
	Series          = charts.Series
	ChartDefinition = charts.ChartDefinition
)

// Function aliases (if any were public)
var (
	NewBarChart            = charts.NewBarChart
	NewBarHorizontalChart  = charts.NewBarHorizontalChart
	NewBarStackedChart     = charts.NewBarStackedChart
	NewBarStacked100Chart  = charts.NewBarStacked100Chart
	NewLineChart           = charts.NewLineChart
	NewLineMarkersChart    = charts.NewLineMarkersChart
	NewLineStackedChart    = charts.NewLineStackedChart
	NewScatterChart        = charts.NewScatterChart
	NewAreaChart           = charts.NewAreaChart
	NewAreaStackedChart    = charts.NewAreaStackedChart
	NewAreaStacked100Chart = charts.NewAreaStacked100Chart
	NewPieChart            = charts.NewPieChart
	NewDoughnutChart       = charts.NewDoughnutChart
	NewBubbleChart         = charts.NewBubbleChart
	NewRadarChart          = charts.NewRadarChart
	NewRadarFilledChart    = charts.NewRadarFilledChart
	NewStockHLCChart       = charts.NewStockHLCChart
	NewStockOHLCChart      = charts.NewStockOHLCChart
	NewComboChart          = charts.NewComboChart
)

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

// Helper aliases
var (
	copyStringSlice         = charts.CopyStringSlice
	copyFloat64Slice        = charts.CopyFloat64Slice
	copyFloat64Pointer      = charts.CopyFloat64Pointer
	normalizeHexColor       = charts.NormalizeHexColor
	toXMLSeries             = charts.ToXMLSeries
	isHexColor              = charts.IsHexColor
	isLegendPosition        = charts.IsLegendPosition
	isValueAxisCrossBetween = charts.IsValueAxisCrossBetween
	validateSeriesList      = charts.ValidateSeriesList
	copySeriesList          = charts.CopySeriesList
)
