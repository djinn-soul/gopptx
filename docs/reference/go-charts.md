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

## `Series` type

Used by `NewComboChart` and wherever a named data series is needed:

```go
type Series struct {
    Name   string
    Values []float64
}
```

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

## Chart data builders

Use these to update or replace chart data on an existing `Presentation`.

Source file: `pkg/pptx/presentation_chart_data_builder.go`

### `NewCategoryChartData(categories []string) *CategoryChartData`

Build a replacement dataset for bar, line, area, pie, and similar category-based charts.

- `AddSeries(name string, values []float64) *CategoryChartData`
- `AddCategoryLevel(categories []string) *CategoryChartData` — multi-level category axis

### `NewXyChartData() *XyChartData`

Build a replacement dataset for scatter / XY charts.

- `AddSeries(name string, xValues, yValues []float64) *XyChartData`

### `NewBubbleChartData() *BubbleChartData`

Build a replacement dataset for bubble charts.

- `AddSeries(name string, xValues, yValues, sizeValues []float64) *BubbleChartData`

## Presentation runtime chart API

Methods on `*Presentation` for reading and mutating charts in an opened file.

Source file: `pkg/pptx/presentation_chart_api.go`

### Listing

- `ListSlideCharts(slideIndex int) ([]SlideChartRef, error)` — return all chart references on a slide

### Updating data

- `UpdateChartData(slideIndex int, selector ChartSelector, data ChartDataUpdate) error`
- `UpdateChartDataByIndex(slideIndex, chartIndex int, data ChartDataUpdate) error`
- `UpdateChartDataByRelID(slideIndex int, relID string, data ChartDataUpdate) error`
- `UpdateChartDataByIndexFromBuilder(slideIndex, chartIndex int, builder ChartDataBuilder) error`
- `UpdateChartDataByRelIDFromBuilder(slideIndex int, relID string, builder ChartDataBuilder) error`

### Updating formatting

- `UpdateChartFormatting(slideIndex int, selector ChartSelector, format ChartFormatUpdate) error`
- `UpdateChartFormattingByIndex(slideIndex, chartIndex int, format ChartFormatUpdate) error`
- `UpdateChartFormattingByRelID(slideIndex int, relID string, format ChartFormatUpdate) error`

### Replacing data (convenience)

- `ReplaceChartData(slideIndex, chartIndex int, categories []string, values []float64) error`
- `ReplaceChartDataByRelID(slideIndex int, relID string, categories []string, values []float64) error`

### Types

- `ChartSelector` — identifies a chart by index (`Index *int`) or relationship ID (`RelID string`)
- `ChartDataUpdate` — complete chart replacement payload
- `ChartFormatUpdate` — partial formatting patch
- `SlideChartRef` — chart reference returned by `ListSlideCharts`
- `ChartDataBuilder` — interface implemented by `*CategoryChartData`, `*XyChartData`, `*BubbleChartData`

See also:

- [Go API Reference](go-api.md)
