package charts

import (
	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Pie3DChart is a categorical pie chart rendered with 3D perspective.
type Pie3DChart struct {
	PieChart
}

// NewPie3DChart creates a 3D pie chart with default layout and style.
func NewPie3DChart(categories []string, values []float64) Pie3DChart {
	return Pie3DChart{PieChart: NewPieChart(categories, values)}
}

// WithAltText sets alternative text for accessibility.
func (c Pie3DChart) WithAltText(text string) Pie3DChart {
	c.PieChart = c.PieChart.WithAltText(text)
	return c
}

// WithDecorative marks the chart as decorative.
func (c Pie3DChart) WithDecorative(enabled bool) Pie3DChart {
	c.PieChart = c.PieChart.WithDecorative(enabled)
	return c
}

// WithTitle sets the chart title.
func (c Pie3DChart) WithTitle(title string) Pie3DChart {
	c.PieChart = c.PieChart.WithTitle(title)
	return c
}

// Position sets chart position in EMU.
func (c Pie3DChart) Position(x styling.Length, y styling.Length) Pie3DChart {
	c.PieChart = c.PieChart.Position(x, y)
	return c
}

// Size sets chart size in EMU.
func (c Pie3DChart) Size(cx styling.Length, cy styling.Length) Pie3DChart {
	c.PieChart = c.PieChart.Size(cx, cy)
	return c
}

// ToChartSpec converts Pie3DChart to the internal 3D-pie XML spec.
func (c Pie3DChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.PieChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindThreeDPie
	return spec
}

// ChartKind reports the internal 3D-pie chart kind.
func (Pie3DChart) ChartKind() string { return pptxxml.ChartKindThreeDPie }

func (Pie3DChart) sealedChart() {}
