package charts

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	RadarStyleMarker = "marker"
	RadarStyleFilled = "filled"
)

// RadarChart is a radar chart using marker style.
type RadarChart struct {
	Title        string
	TitleOverlay bool
	Categories   []string
	Values       []float64
	X            styling.Length
	Y            styling.Length
	CX           styling.Length
	CY           styling.Length

	LineColor             string
	SeriesName            string
	ShowLegend            bool
	LegendPosition        string
	LegendOverlay         bool
	ShowDataLabels        bool
	DataLabels            DataLabelSettings
	ShowMajorGridlines    bool
	CategoryAxisTitle     string
	ValueAxisTitle        string
	ValueFormat           string
	ValueAxisCrossBetween string
	MinValue              *float64
	MaxValue              *float64
	RadarStyle            string

	// Accessibility
	AltText      string
	IsDecorative bool
}

func NewRadarChart(categories []string, values []float64) RadarChart {
	cats, vals := copyChartData(categories, values)
	return RadarChart{
		Title:      defaultChartTitle,
		Categories: cats,
		Values:     vals,
		X:          styling.Emu(defaultChartX),
		Y:          styling.Emu(defaultChartY),
		CX:         styling.Emu(defaultChartCX),
		CY:         styling.Emu(defaultChartCY),

		LineColor:             defaultChartSeriesColor,
		SeriesName:            defaultChartSeriesName,
		ShowLegend:            false,
		LegendPosition:        LegendPositionRight,
		ShowDataLabels:        false,
		ShowMajorGridlines:    true,
		ValueFormat:           defaultChartValueFormat,
		ValueAxisCrossBetween: ValueAxisCrossBetweenBetween,
		RadarStyle:            RadarStyleMarker,
	}
}

// WithAltText sets the alternative text for accessibility.
func (c RadarChart) WithAltText(text string) RadarChart {
	c.AltText = text
	return c
}

// WithDecorative marks the chart as decorative (ignored by screen readers).
func (c RadarChart) WithDecorative(enabled bool) RadarChart {
	c.IsDecorative = enabled
	return c
}

func (c RadarChart) Position(x styling.Length, y styling.Length) RadarChart {
	c.X = x
	c.Y = y
	return c
}

func (c RadarChart) Size(cx styling.Length, cy styling.Length) RadarChart {
	c.CX = cx
	c.CY = cy
	return c
}

func (c RadarChart) WithTitle(title string) RadarChart {
	c.Title = title
	return c
}

func (c RadarChart) WithLineColor(color string) RadarChart {
	c.LineColor = NormalizeHexColor(color)
	return c
}

// ToChartSpec converts RadarChart to internal XML spec.
func (c RadarChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := &pptxxml.ChartSpec{
		Kind:         pptxxml.ChartKindRadar,
		Title:        c.Title,
		TitleOverlay: c.TitleOverlay,
		Categories:   CopyStringSlice(c.Categories),
		Values:       CopyFloat64Slice(c.Values),
		X:            c.X.Emu(),
		Y:            c.Y.Emu(),
		CX:           c.CX.Emu(),
		CY:           c.CY.Emu(),

		Color:                 NormalizeHexColor(c.LineColor),
		SeriesName:            c.SeriesName,
		ShowLegend:            c.ShowLegend,
		LegendPosition:        c.LegendPosition,
		LegendOverlay:         c.LegendOverlay,
		ShowDataLabels:        c.ShowDataLabels,
		ShowMajorGridlines:    c.ShowMajorGridlines,
		CategoryAxisTitle:     c.CategoryAxisTitle,
		ValueAxisTitle:        c.ValueAxisTitle,
		ValueFormat:           c.ValueFormat,
		ValueAxisCrossBetween: c.ValueAxisCrossBetween,
		MinValue:              CopyFloat64Pointer(c.MinValue),
		MaxValue:              CopyFloat64Pointer(c.MaxValue),
		RadarStyle:            c.RadarStyle,
		AltText:               c.AltText,
		IsDecorative:          c.IsDecorative,
	}
	applyDataLabelSettings(spec, c.DataLabels)
	return spec
}

// Validate checks the radar chart for consistency.
func (c RadarChart) Validate(slideIndex int) error {
	if !c.IsDecorative && len(c.AltText) > common.MaxAltTextLength {
		return fmt.Errorf("slide %d radar chart alt text exceeds %d characters", slideIndex, common.MaxAltTextLength)
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
	if !IsHexColor(c.LineColor) {
		return fmt.Errorf("slide %d radar chart color must be 6-digit RGB hex", slideIndex)
	}
	if strings.TrimSpace(c.SeriesName) == "" {
		return fmt.Errorf("slide %d radar chart series name cannot be empty", slideIndex)
	}
	if !IsLegendPosition(c.LegendPosition) {
		return fmt.Errorf("slide %d radar chart legend position must be one of r,l,t,b", slideIndex)
	}
	if !IsDataLabelPosition(c.DataLabels.Position) {
		return fmt.Errorf(
			"slide %d radar chart data-label position must be ctr,inEnd,inBase,outEnd,bestFit,l,r,t,or b",
			slideIndex,
		)
	}
	if strings.TrimSpace(c.ValueFormat) == "" {
		return fmt.Errorf("slide %d radar chart value format cannot be empty", slideIndex)
	}
	if !IsValueAxisCrossBetween(c.ValueAxisCrossBetween) {
		return fmt.Errorf("slide %d radar chart value-axis crossBetween must be between or midCat", slideIndex)
	}
	if c.RadarStyle != RadarStyleMarker && c.RadarStyle != RadarStyleFilled {
		return fmt.Errorf("slide %d radar chart style must be marker or filled", slideIndex)
	}
	if err := validateValueRange(c.MinValue, c.MaxValue, slideIndex); err != nil {
		return err
	}
	return nil
}

func (c RadarChart) GetCategories() []string {
	return c.Categories
}

func (c RadarChart) GetValues() []float64 {
	return c.Values
}

// RadarFilledChart is a radar chart using filled style.
type RadarFilledChart struct {
	RadarChart
}

func NewRadarFilledChart(categories []string, values []float64) RadarFilledChart {
	base := NewRadarChart(categories, values)
	base.RadarStyle = RadarStyleFilled
	return RadarFilledChart{RadarChart: base}
}

// ToChartSpec converts RadarFilledChart to internal XML spec.
func (c RadarFilledChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.RadarChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindRadarFilled
	spec.RadarStyle = RadarStyleFilled
	return spec
}

// WithAltText sets the alternative text for accessibility.
func (c RadarFilledChart) WithAltText(text string) RadarFilledChart {
	c.AltText = text
	return c
}

// WithDecorative marks the chart as decorative (ignored by screen readers).
func (c RadarFilledChart) WithDecorative(enabled bool) RadarFilledChart {
	c.IsDecorative = enabled
	return c
}

// Validate checks the radar chart for consistency.
func (c RadarFilledChart) Validate(slideIndex int) error {
	return c.RadarChart.Validate(slideIndex)
}
