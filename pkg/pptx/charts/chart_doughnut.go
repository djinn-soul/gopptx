package charts

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

// DoughnutChart is a simple categorical doughnut chart.
type DoughnutChart struct {
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
	HoleSize       int
}

// NewDoughnutChart creates a doughnut chart with default layout and style.
func NewDoughnutChart(categories []string, values []float64) DoughnutChart {
	cats, vals := copyChartData(categories, values)
	return DoughnutChart{
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
		HoleSize:       50,
	}
}

// Position sets chart position in EMU.
func (c DoughnutChart) Position(x int64, y int64) DoughnutChart {
	c.X = x
	c.Y = y
	return c
}

// Size sets chart size in EMU.
func (c DoughnutChart) Size(cx int64, cy int64) DoughnutChart {
	c.CX = cx
	c.CY = cy
	return c
}

// WithTitle sets the chart title.
func (c DoughnutChart) WithTitle(title string) DoughnutChart {
	c.Title = title
	return c
}

// ToChartSpec converts DoughnutChart to internal XML spec.
func (c DoughnutChart) ToChartSpec() *pptxxml.ChartSpec {
	return &pptxxml.ChartSpec{
		Kind:           pptxxml.ChartKindDoughnut,
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
		HoleSize:       c.HoleSize,
	}
}

// Validate checks the doughnut chart for consistency.
func (c DoughnutChart) Validate(slideIndex int) error {
	if err := validateChartCore(
		slideIndex,
		c.Title,
		c.Categories,
		c.Values,
		c.X,
		c.Y,
		c.CX,
		c.CY,
		false,
	); err != nil {
		return err
	}
	if strings.TrimSpace(c.SeriesName) == "" {
		return fmt.Errorf("slide %d doughnut chart series name cannot be empty", slideIndex)
	}
	if !IsLegendPosition(c.LegendPosition) {
		return fmt.Errorf("slide %d doughnut chart legend position must be one of r,l,t,b", slideIndex)
	}
	if c.HoleSize < 10 || c.HoleSize > 90 {
		return fmt.Errorf("slide %d doughnut chart hole size must be between 10 and 90", slideIndex)
	}
	return nil
}
