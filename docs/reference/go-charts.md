# Go Charts Reference

This page documents the chart constructors exposed by `pkg/pptx`.

Primary source files:

- `pkg/pptx/chart_compat.go`
- `pkg/pptx/charts/chart.go`
- `pkg/pptx/charts/chart_family1_variants.go`
- `pkg/pptx/charts/chart_scatter.go`
- `pkg/pptx/charts/chart_bubble.go`
- `pkg/pptx/charts/chart_doughnut.go`
- `pkg/pptx/charts/chart_area_variant_options.go`
- `pkg/pptx/charts/chart_bar_variant_options.go`

## Constructors

- `NewBarChart(categories []string, values []float64) BarChart`
- `NewBarHorizontalChart(categories []string, values []float64) BarHorizontalChart`
- `NewBarStackedChart(categories []string, values []float64) BarStackedChart`
- `NewBarStacked100Chart(categories []string, values []float64) BarStacked100Chart`
- `NewLineChart(categories []string, values []float64) LineChart`
- `NewLineMarkersChart(categories []string, values []float64) LineMarkersChart`
- `NewLineStackedChart(categories []string, values []float64) LineStackedChart`
- `NewScatterChart(xValues, yValues []float64) ScatterChart`
- `NewAreaChart(categories []string, values []float64) AreaChart`
- `NewAreaStackedChart(categories []string, values []float64) AreaStackedChart`
- `NewAreaStacked100Chart(categories []string, values []float64) AreaStacked100Chart`
- `NewPieChart(categories []string, values []float64) PieChart`
- `NewDoughnutChart(categories []string, values []float64) DoughnutChart`
- `NewBubbleChart(xValues, yValues, sizeValues []float64) BubbleChart`
- `NewRadarChart(categories []string, values []float64) RadarChart`
- `NewRadarFilledChart(categories []string, values []float64) RadarFilledChart`
- `NewStockHLCChart(categories []string, high, low, closeValues []float64) StockHLCChart`
- `NewStockOHLCChart(categories []string, open, high, low, closeValues []float64) StockOHLCChart`
- `NewComboChart(categories []string, barSeries []Series, lineSeries []Series) ComboChart`

## Fluent methods

Most chart types share a `With*` fluent API for title, legend, labels, position, and size:

- `WithTitle(...)`
- `WithAltText(...)`
- `WithDecorative(...)`
- `WithLegend(...)`
- `WithLegendPosition(...)`
- `WithDataLabels(...)`
- `WithAxisTitles(...)`
- `WithMajorGridlines(...)`
- `WithValueFormat(...)`
- `Position(...)`
- `Size(...)`

## Constants

- `LegendPositionRight`
- `LegendPositionLeft`
- `LegendPositionTop`
- `LegendPositionBottom`
- `ValueAxisCrossBetweenBetween`
- `ValueAxisCrossBetweenMidCategory`
- `RadarStyleFilled`
- `ScatterStyleMarker`
- `ScatterStyleLineMarker`
- `ScatterStyleSmoothMarker`

See also:

- [Go API Reference](go-api.md)
