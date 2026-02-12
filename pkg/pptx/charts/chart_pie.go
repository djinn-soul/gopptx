package charts

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// PieChart is a simple categorical pie chart.
type PieChart struct {
	Title        string
	TitleOverlay bool
	Categories   []string
	Values       []float64
	X            styling.Length
	Y            styling.Length
	CX           styling.Length
	CY           styling.Length

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
		Title:      "Chart",
		Categories: cats,
		Values:     vals,
		X:          styling.Emu(685800),
		Y:          styling.Emu(1800000),
		CX:         styling.Emu(7772400),
		CY:         styling.Emu(4114800),

		SeriesName:     "Series 1",
		ShowLegend:     false,
		LegendPosition: LegendPositionRight,
		ShowDataLabels: false,
	}
}

// Position sets chart position in EMU.
func (c PieChart) Position(x styling.Length, y styling.Length) PieChart {
	c.X = x
	c.Y = y
	return c
}

// Size sets chart size in EMU.
func (c PieChart) Size(cx styling.Length, cy styling.Length) PieChart {
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
		Kind:         pptxxml.ChartKindPie,
		Title:        c.Title,
		TitleOverlay: c.TitleOverlay,
		Categories:   CopyStringSlice(c.Categories),
		Values:       CopyFloat64Slice(c.Values),
		X:            c.X.Emu(),
		Y:            c.Y.Emu(),
		CX:           c.CX.Emu(),
		CY:           c.CY.Emu(),

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
		false,
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
