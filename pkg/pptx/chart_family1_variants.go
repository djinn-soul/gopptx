package pptx

import "github.com/djinn-soul/gopptx/internal/pptxxml"

// BarHorizontalChart is a horizontal clustered bar chart variant.
type BarHorizontalChart struct {
	BarChart
}

func NewBarHorizontalChart(categories []string, values []float64) BarHorizontalChart {
	return BarHorizontalChart{BarChart: NewBarChart(categories, values)}
}

// ToChartSpec converts BarHorizontalChart to internal XML spec.
func (c BarHorizontalChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.BarChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindBarHorizontal
	spec.BarDir = "bar"
	return spec
}

// Validate checks the bar chart for consistency.
func (c BarHorizontalChart) Validate(slideIndex int) error {
	return c.BarChart.Validate(slideIndex)
}

// BarStackedChart is a stacked bar chart variant.
type BarStackedChart struct {
	BarChart
}

func NewBarStackedChart(categories []string, values []float64) BarStackedChart {
	return BarStackedChart{BarChart: NewBarChart(categories, values)}
}

// ToChartSpec converts BarStackedChart to internal XML spec.
func (c BarStackedChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.BarChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindBarStacked
	spec.Grouping = "stacked"
	return spec
}

// Validate checks the bar chart for consistency.
func (c BarStackedChart) Validate(slideIndex int) error {
	return c.BarChart.Validate(slideIndex)
}

// BarStacked100Chart is a 100%% stacked bar chart variant.
type BarStacked100Chart struct {
	BarChart
}

func NewBarStacked100Chart(categories []string, values []float64) BarStacked100Chart {
	return BarStacked100Chart{BarChart: NewBarChart(categories, values)}
}

// ToChartSpec converts BarStacked100Chart to internal XML spec.
func (c BarStacked100Chart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.BarChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindBarStacked100
	spec.Grouping = "percentStacked"
	return spec
}

// Validate checks the bar chart for consistency.
func (c BarStacked100Chart) Validate(slideIndex int) error {
	return c.BarChart.Validate(slideIndex)
}

// LineMarkersChart is a line-with-markers chart variant.
type LineMarkersChart struct {
	LineChart
}

func NewLineMarkersChart(categories []string, values []float64) LineMarkersChart {
	return LineMarkersChart{LineChart: NewLineChart(categories, values)}
}

// ToChartSpec converts LineMarkersChart to internal XML spec.
func (c LineMarkersChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.LineChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindLineMarkers
	spec.ShowMarkers = true
	return spec
}

// Validate checks the line chart for consistency.
func (c LineMarkersChart) Validate(slideIndex int) error {
	return c.LineChart.Validate(slideIndex)
}

// LineStackedChart is a stacked line chart variant.
type LineStackedChart struct {
	LineChart
}

func NewLineStackedChart(categories []string, values []float64) LineStackedChart {
	return LineStackedChart{LineChart: NewLineChart(categories, values)}
}

// ToChartSpec converts LineStackedChart to internal XML spec.
func (c LineStackedChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.LineChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindLineStacked
	spec.Grouping = "stacked"
	return spec
}

// Validate checks the line chart for consistency.
func (c LineStackedChart) Validate(slideIndex int) error {
	return c.LineChart.Validate(slideIndex)
}

// AreaStackedChart is a stacked area chart variant.
type AreaStackedChart struct {
	AreaChart
}

func NewAreaStackedChart(categories []string, values []float64) AreaStackedChart {
	return AreaStackedChart{AreaChart: NewAreaChart(categories, values)}
}

// ToChartSpec converts AreaStackedChart to internal XML spec.
func (c AreaStackedChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.AreaChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindAreaStacked
	spec.Grouping = "stacked"
	return spec
}

// Validate checks the area chart for consistency.
func (c AreaStackedChart) Validate(slideIndex int) error {
	return c.AreaChart.Validate(slideIndex)
}

// AreaStacked100Chart is a 100%% stacked area chart variant.
type AreaStacked100Chart struct {
	AreaChart
}

func NewAreaStacked100Chart(categories []string, values []float64) AreaStacked100Chart {
	return AreaStacked100Chart{AreaChart: NewAreaChart(categories, values)}
}

// ToChartSpec converts AreaStacked100Chart to internal XML spec.
func (c AreaStacked100Chart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.AreaChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindAreaStacked100
	spec.Grouping = "percentStacked"
	return spec
}

// Validate checks the area chart for consistency.
func (c AreaStacked100Chart) Validate(slideIndex int) error {
	return c.AreaChart.Validate(slideIndex)
}
