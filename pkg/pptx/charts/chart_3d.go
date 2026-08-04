package charts

import (
	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Column3DChart is a clustered column chart drawn with 3-D perspective.
type Column3DChart struct {
	BarChart
}

// NewColumn3DChart creates a 3-D column chart with default layout and style.
func NewColumn3DChart(categories []string, values []float64) Column3DChart {
	return Column3DChart{BarChart: NewBarChart(categories, values)}
}

// ToChartSpec converts Column3DChart to the internal 3-D column XML spec.
func (c Column3DChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.BarChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindThreeDColumn
	return spec
}

// ChartKind reports the internal 3-D column chart kind.
func (Column3DChart) ChartKind() string { return pptxxml.ChartKindThreeDColumn }

// Validate checks the chart for consistency.
func (c Column3DChart) Validate(slideIndex int) error {
	return c.BarChart.Validate(slideIndex)
}

// Position sets chart position in EMU.
func (c Column3DChart) Position(x styling.Length, y styling.Length) Column3DChart {
	c.BarChart = c.BarChart.Position(x, y)
	return c
}

// Size sets chart size in EMU.
func (c Column3DChart) Size(cx styling.Length, cy styling.Length) Column3DChart {
	c.BarChart = c.BarChart.Size(cx, cy)
	return c
}

// WithTitle sets the chart title.
func (c Column3DChart) WithTitle(title string) Column3DChart {
	c.BarChart = c.BarChart.WithTitle(title)
	return c
}

// WithLegend toggles the legend.
func (c Column3DChart) WithLegend(show bool) Column3DChart {
	c.BarChart = c.BarChart.WithLegend(show)
	return c
}

// WithSeriesName sets the series name.
func (c Column3DChart) WithSeriesName(name string) Column3DChart {
	c.BarChart = c.BarChart.WithSeriesName(name)
	return c
}

// Bar3DChart is a clustered horizontal bar chart drawn with 3-D perspective.
type Bar3DChart struct {
	BarChart
}

// NewBar3DChart creates a 3-D horizontal bar chart with default layout and style.
func NewBar3DChart(categories []string, values []float64) Bar3DChart {
	return Bar3DChart{BarChart: NewBarChart(categories, values)}
}

// ToChartSpec converts Bar3DChart to the internal 3-D bar XML spec.
func (c Bar3DChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.BarChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindThreeDBar
	spec.BarDir = chartBarDirection
	return spec
}

// ChartKind reports the internal 3-D bar chart kind.
func (Bar3DChart) ChartKind() string { return pptxxml.ChartKindThreeDBar }

// Validate checks the chart for consistency.
func (c Bar3DChart) Validate(slideIndex int) error {
	return c.BarChart.Validate(slideIndex)
}

// Position sets chart position in EMU.
func (c Bar3DChart) Position(x styling.Length, y styling.Length) Bar3DChart {
	c.BarChart = c.BarChart.Position(x, y)
	return c
}

// Size sets chart size in EMU.
func (c Bar3DChart) Size(cx styling.Length, cy styling.Length) Bar3DChart {
	c.BarChart = c.BarChart.Size(cx, cy)
	return c
}

// WithTitle sets the chart title.
func (c Bar3DChart) WithTitle(title string) Bar3DChart {
	c.BarChart = c.BarChart.WithTitle(title)
	return c
}

// WithLegend toggles the legend.
func (c Bar3DChart) WithLegend(show bool) Bar3DChart {
	c.BarChart = c.BarChart.WithLegend(show)
	return c
}

// WithSeriesName sets the series name.
func (c Bar3DChart) WithSeriesName(name string) Bar3DChart {
	c.BarChart = c.BarChart.WithSeriesName(name)
	return c
}

// Line3DChart is a categorical line chart drawn with 3-D perspective, each
// series sitting on its own depth row.
type Line3DChart struct {
	LineChart
}

// NewLine3DChart creates a 3-D line chart with default layout and style.
func NewLine3DChart(categories []string, values []float64) Line3DChart {
	return Line3DChart{LineChart: NewLineChart(categories, values)}
}

// ToChartSpec converts Line3DChart to the internal 3-D line XML spec.
func (c Line3DChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.LineChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindThreeDLine
	return spec
}

// ChartKind reports the internal 3-D line chart kind.
func (Line3DChart) ChartKind() string { return pptxxml.ChartKindThreeDLine }

// Validate checks the chart for consistency.
func (c Line3DChart) Validate(slideIndex int) error {
	return c.LineChart.Validate(slideIndex)
}

// Position sets chart position in EMU.
func (c Line3DChart) Position(x styling.Length, y styling.Length) Line3DChart {
	c.LineChart = c.LineChart.Position(x, y)
	return c
}

// Size sets chart size in EMU.
func (c Line3DChart) Size(cx styling.Length, cy styling.Length) Line3DChart {
	c.LineChart = c.LineChart.Size(cx, cy)
	return c
}

// WithTitle sets the chart title.
func (c Line3DChart) WithTitle(title string) Line3DChart {
	c.LineChart = c.LineChart.WithTitle(title)
	return c
}

// WithLegend toggles the legend.
func (c Line3DChart) WithLegend(show bool) Line3DChart {
	c.LineChart = c.LineChart.WithLegend(show)
	return c
}

// WithSeriesName sets the series name.
func (c Line3DChart) WithSeriesName(name string) Line3DChart {
	c.LineChart = c.LineChart.WithSeriesName(name)
	return c
}

// Area3DChart is a categorical area chart drawn with 3-D perspective.
type Area3DChart struct {
	AreaChart
}

// NewArea3DChart creates a 3-D area chart with default layout and style.
func NewArea3DChart(categories []string, values []float64) Area3DChart {
	return Area3DChart{AreaChart: NewAreaChart(categories, values)}
}

// ToChartSpec converts Area3DChart to the internal 3-D area XML spec.
func (c Area3DChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.AreaChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindThreeDArea
	return spec
}

// ChartKind reports the internal 3-D area chart kind.
func (Area3DChart) ChartKind() string { return pptxxml.ChartKindThreeDArea }

// Validate checks the chart for consistency.
func (c Area3DChart) Validate(slideIndex int) error {
	return c.AreaChart.Validate(slideIndex)
}

// Position sets chart position in EMU.
func (c Area3DChart) Position(x styling.Length, y styling.Length) Area3DChart {
	c.AreaChart = c.AreaChart.Position(x, y)
	return c
}

// Size sets chart size in EMU.
func (c Area3DChart) Size(cx styling.Length, cy styling.Length) Area3DChart {
	c.AreaChart = c.AreaChart.Size(cx, cy)
	return c
}

// WithTitle sets the chart title.
func (c Area3DChart) WithTitle(title string) Area3DChart {
	c.AreaChart = c.AreaChart.WithTitle(title)
	return c
}

// WithLegend toggles the legend.
func (c Area3DChart) WithLegend(show bool) Area3DChart {
	c.AreaChart = c.AreaChart.WithLegend(show)
	return c
}

// WithSeriesName sets the series name.
func (c Area3DChart) WithSeriesName(name string) Area3DChart {
	c.AreaChart = c.AreaChart.WithSeriesName(name)
	return c
}
