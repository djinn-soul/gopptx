package charts

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
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
	DataLabels     DataLabelSettings

	// Accessibility
	AltText      string
	IsDecorative bool
}

// NewPieChart creates a pie chart with default layout and style.
func NewPieChart(categories []string, values []float64) PieChart {
	cats, vals := copyChartData(categories, values)
	return PieChart{
		Title:      defaultChartTitle,
		Categories: cats,
		Values:     vals,
		X:          styling.Emu(defaultChartX),
		Y:          styling.Emu(defaultChartY),
		CX:         styling.Emu(defaultChartCX),
		CY:         styling.Emu(defaultChartCY),

		SeriesName:     defaultChartSeriesName,
		ShowLegend:     false,
		LegendPosition: LegendPositionRight,
		ShowDataLabels: false,
	}
}

// WithAltText sets the alternative text for accessibility.
func (c PieChart) WithAltText(text string) PieChart {
	c.AltText = text
	return c
}

// WithDecorative marks the chart as decorative (ignored by screen readers).
func (c PieChart) WithDecorative(enabled bool) PieChart {
	c.IsDecorative = enabled
	return c
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
	spec := &pptxxml.ChartSpec{
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
		AltText:        c.AltText,
		IsDecorative:   c.IsDecorative,
	}
	applyDataLabelSettings(spec, c.DataLabels)
	return spec
}

// Validate checks the pie chart for consistency.
func (c PieChart) Validate(slideIndex int) error {
	if !c.IsDecorative && len(c.AltText) > common.MaxAltTextLength {
		return fmt.Errorf("slide %d pie chart alt text exceeds %d characters", slideIndex, common.MaxAltTextLength)
	}
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
	if !IsDataLabelPosition(c.DataLabels.Position) {
		return fmt.Errorf(
			"slide %d pie chart data-label position must be ctr,inEnd,inBase,outEnd,bestFit,l,r,t,or b",
			slideIndex,
		)
	}
	return nil
}

func (c PieChart) GetCategories() []string {
	return c.Categories
}

func (c PieChart) GetValues() []float64 {
	return c.Values
}
