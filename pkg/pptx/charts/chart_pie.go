package charts

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

// PieChart is a simple categorical pie chart.
type PieChart struct {
	Title          string
	TitleOverlay   bool
	Categories     []string
	Values         []float64
	X              int64
	Y              int64
	CX             int64
	CY             int64
	SeriesName     string
	ShowLegend     bool
	LegendPosition string
	LegendOverlay  bool
	ShowDataLabels bool
}

// NewPieChart creates a pie chart with default layout and style.
func NewPieChart(categories []string, values []float64) PieChart {
	cats, vals := copyChartData(categories, values)
	return PieChart{
		Title:          "Chart",
		Categories:     cats,
		Values:         vals,
		X:              685800,
		Y:              1800000,
		CX:             7772400,
		CY:             4114800,
		SeriesName:     "Series 1",
		ShowLegend:     false,
		LegendPosition: LegendPositionRight,
		ShowDataLabels: false,
	}
}

// Position sets chart position in EMU.
func (c PieChart) Position(x int64, y int64) PieChart {
	c.X = x
	c.Y = y
	return c
}

// Size sets chart size in EMU.
func (c PieChart) Size(cx int64, cy int64) PieChart {
	c.CX = cx
	c.CY = cy
	return c
}

// WithTitle sets the chart title.
func (c PieChart) WithTitle(title string) PieChart {
	c.Title = title
	return c
}

// ToChartSpec converts PieChart to internal XML spec.
func (c PieChart) ToChartSpec() *pptxxml.ChartSpec {
	return &pptxxml.ChartSpec{
		Kind:           pptxxml.ChartKindPie,
		Title:          c.Title,
		TitleOverlay:   c.TitleOverlay,
		Categories:     CopyStringSlice(c.Categories),
		Values:         CopyFloat64Slice(c.Values),
		X:              c.X,
		Y:              c.Y,
		CX:             c.CX,
		CY:             c.CY,
		SeriesName:     c.SeriesName,
		ShowLegend:     c.ShowLegend,
		LegendPosition: c.LegendPosition,
		LegendOverlay:  c.LegendOverlay,
		ShowDataLabels: c.ShowDataLabels,
	}
}

// Validate checks the pie chart for consistency.
func (c PieChart) Validate(slideIndex int) error {
	if err := validateChartCore(
		slideIndex,
		c.Title,
		c.Categories,
		c.Values,
		c.X,
		c.Y,
		c.CX,
		c.CY,
	); err != nil {
		return err
	}
	if strings.TrimSpace(c.SeriesName) == "" {
		return fmt.Errorf("slide %d pie chart series name cannot be empty", slideIndex)
	}
	if !IsLegendPosition(c.LegendPosition) {
		return fmt.Errorf("slide %d pie chart legend position must be one of r,l,t,b", slideIndex)
	}
	return nil
}
