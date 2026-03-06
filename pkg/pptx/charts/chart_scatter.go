package charts

import (
	"fmt"
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	ScatterStyleMarker       = "marker"
	ScatterStyleLineMarker   = "lineMarker"
	ScatterStyleSmoothMarker = "smoothMarker"
)

// ScatterChart is a simple XY scatter chart.
type ScatterChart struct {
	Title        string
	TitleOverlay bool
	XValues      []float64
	YValues      []float64
	X            styling.Length
	Y            styling.Length
	CX           styling.Length
	CY           styling.Length

	LineColor                  string
	SeriesName                 string
	ScatterStyle               string
	ShowLegend                 bool
	LegendPosition             string
	LegendOverlay              bool
	ShowDataLabels             bool
	DataLabels                 DataLabelSettings
	ShowMajorGridlines         bool
	ShowCategoryMajorGridlines bool
	CategoryAxisTitle          string
	ValueAxisTitle             string
	CategoryTickLabelPosition  string
	ValueTickLabelPosition     string
	CategoryAxisCrosses        string
	ValueAxisCrosses           string
	ValueFormat                string
	ValueAxisCrossBetween      string
	MinValue                   *float64
	MaxValue                   *float64

	// Accessibility
	AltText      string
	IsDecorative bool
}

// NewScatterChart creates a scatter chart with default layout and style.
func NewScatterChart(xValues []float64, yValues []float64) ScatterChart {
	xs := make([]float64, len(xValues))
	copy(xs, xValues)
	ys := make([]float64, len(yValues))
	copy(ys, yValues)
	return ScatterChart{
		Title:   "Chart",
		XValues: xs,
		YValues: ys,
		X:       styling.Emu(defaultChartX),
		Y:       styling.Emu(defaultChartY),
		CX:      styling.Emu(defaultChartCX),
		CY:      styling.Emu(defaultChartCY),

		LineColor:                  "4F81BD",
		SeriesName:                 "Series 1",
		ScatterStyle:               ScatterStyleMarker,
		ShowLegend:                 false,
		LegendPosition:             LegendPositionRight,
		ShowDataLabels:             false,
		ShowMajorGridlines:         true,
		ShowCategoryMajorGridlines: false,
		CategoryTickLabelPosition:  AxisTickLabelPositionNextTo,
		ValueTickLabelPosition:     AxisTickLabelPositionNextTo,
		CategoryAxisCrosses:        AxisCrossesAutoZero,
		ValueAxisCrosses:           AxisCrossesAutoZero,
		ValueFormat:                "General",
		ValueAxisCrossBetween:      ValueAxisCrossBetweenBetween,
	}
}

// WithAltText sets the alternative text for accessibility.
func (c ScatterChart) WithAltText(text string) ScatterChart {
	c.AltText = text
	return c
}

// WithDecorative marks the chart as decorative (ignored by screen readers).
func (c ScatterChart) WithDecorative(enabled bool) ScatterChart {
	c.IsDecorative = enabled
	return c
}

// Position sets chart position in EMU.
func (c ScatterChart) Position(x styling.Length, y styling.Length) ScatterChart {
	c.X = x
	c.Y = y
	return c
}

// Size sets chart size in EMU.
func (c ScatterChart) Size(cx styling.Length, cy styling.Length) ScatterChart {
	c.CX = cx
	c.CY = cy
	return c
}

// WithTitle sets the chart title.
func (c ScatterChart) WithTitle(title string) ScatterChart {
	c.Title = title
	return c
}

// WithLineColor sets the scatter line/marker color using RGB hex.
func (c ScatterChart) WithLineColor(color string) ScatterChart {
	c.LineColor = NormalizeHexColor(color)
	return c
}

// ToChartSpec converts ScatterChart to internal XML spec.
func (c ScatterChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := &pptxxml.ChartSpec{
		Kind:         pptxxml.ChartKindScatter,
		Title:        c.Title,
		TitleOverlay: c.TitleOverlay,
		XValues:      CopyFloat64Slice(c.XValues),
		Values:       CopyFloat64Slice(c.YValues),
		X:            c.X.Emu(),
		Y:            c.Y.Emu(),
		CX:           c.CX.Emu(),
		CY:           c.CY.Emu(),

		Color:                      NormalizeHexColor(c.LineColor),
		SeriesName:                 c.SeriesName,
		ScatterStyle:               c.ScatterStyle,
		ShowLegend:                 c.ShowLegend,
		LegendPosition:             c.LegendPosition,
		LegendOverlay:              c.LegendOverlay,
		ShowDataLabels:             c.ShowDataLabels,
		ShowMajorGridlines:         c.ShowMajorGridlines,
		ShowCategoryMajorGridlines: c.ShowCategoryMajorGridlines,
		CategoryAxisTitle:          c.CategoryAxisTitle,
		ValueAxisTitle:             c.ValueAxisTitle,
		CategoryTickLabelPosition:  c.CategoryTickLabelPosition,
		ValueTickLabelPosition:     c.ValueTickLabelPosition,
		CategoryAxisCrosses:        c.CategoryAxisCrosses,
		ValueAxisCrosses:           c.ValueAxisCrosses,
		ValueFormat:                c.ValueFormat,
		ValueAxisCrossBetween:      c.ValueAxisCrossBetween,
		MinValue:                   CopyFloat64Pointer(c.MinValue),
		MaxValue:                   CopyFloat64Pointer(c.MaxValue),
		AltText:                    c.AltText,
		IsDecorative:               c.IsDecorative,
	}
	applyDataLabelSettings(spec, c.DataLabels)
	return spec
}

// Validate checks the scatter chart for consistency.
func (c ScatterChart) Validate(slideIndex int) error {
	if !c.IsDecorative && len(c.AltText) > common.MaxAltTextLength {
		return fmt.Errorf("slide %d scatter chart alt text exceeds %d characters", slideIndex, common.MaxAltTextLength)
	}
	if err := c.validateLayout(slideIndex); err != nil {
		return err
	}
	if err := c.validateMetadata(slideIndex); err != nil {
		return err
	}
	if err := c.validatePoints(slideIndex); err != nil {
		return err
	}
	return validateValueRange(c.MinValue, c.MaxValue, slideIndex)
}

func (c ScatterChart) validateLayout(slideIndex int) error {
	if c.X < 0 || c.Y < 0 {
		return fmt.Errorf("slide %d scatter chart position cannot be negative", slideIndex)
	}
	if c.CX <= 0 || c.CY <= 0 {
		return fmt.Errorf("slide %d scatter chart size must be > 0", slideIndex)
	}
	return nil
}

func (c ScatterChart) validateMetadata(slideIndex int) error {
	if strings.TrimSpace(c.Title) == "" {
		return fmt.Errorf("slide %d scatter chart title cannot be empty", slideIndex)
	}
	if !IsHexColor(c.LineColor) {
		return fmt.Errorf("slide %d scatter chart color must be 6-digit RGB hex", slideIndex)
	}
	if strings.TrimSpace(c.SeriesName) == "" {
		return fmt.Errorf("slide %d scatter chart series name cannot be empty", slideIndex)
	}
	if !IsLegendPosition(c.LegendPosition) {
		return fmt.Errorf("slide %d scatter chart legend position must be one of r,l,t,b", slideIndex)
	}
	if !isScatterStyle(c.ScatterStyle) {
		return fmt.Errorf("slide %d scatter style must be one of marker,lineMarker,smoothMarker", slideIndex)
	}
	if !IsDataLabelPosition(c.DataLabels.Position) {
		return fmt.Errorf(
			"slide %d scatter chart data-label position must be ctr,inEnd,inBase,outEnd,bestFit,l,r,t,or b",
			slideIndex,
		)
	}
	if strings.TrimSpace(c.ValueFormat) == "" {
		return fmt.Errorf("slide %d scatter chart value format cannot be empty", slideIndex)
	}
	if !IsAxisTickLabelPosition(c.CategoryTickLabelPosition) {
		return fmt.Errorf(
			"slide %d scatter chart x-axis tick label position must be nextTo, low, high, or none",
			slideIndex,
		)
	}
	if !IsAxisTickLabelPosition(c.ValueTickLabelPosition) {
		return fmt.Errorf(
			"slide %d scatter chart y-axis tick label position must be nextTo, low, high, or none",
			slideIndex,
		)
	}
	if !IsAxisCrosses(c.CategoryAxisCrosses) {
		return fmt.Errorf("slide %d scatter chart x-axis crosses must be autoZero, min, or max", slideIndex)
	}
	if !IsAxisCrosses(c.ValueAxisCrosses) {
		return fmt.Errorf("slide %d scatter chart y-axis crosses must be autoZero, min, or max", slideIndex)
	}
	if !IsValueAxisCrossBetween(c.ValueAxisCrossBetween) {
		return fmt.Errorf("slide %d scatter chart value-axis crossBetween must be between or midCat", slideIndex)
	}
	return nil
}

func (c ScatterChart) validatePoints(slideIndex int) error {
	if len(c.XValues) == 0 {
		return fmt.Errorf("slide %d scatter chart must define at least one point", slideIndex)
	}
	if len(c.XValues) != len(c.YValues) {
		return fmt.Errorf("slide %d scatter chart x/y length mismatch (%d vs %d)",
			slideIndex, len(c.XValues), len(c.YValues))
	}
	for i := range c.XValues {
		if math.IsNaN(c.XValues[i]) || math.IsInf(c.XValues[i], 0) {
			return fmt.Errorf("slide %d scatter x value %d must be finite", slideIndex, i+1)
		}
		if math.IsNaN(c.YValues[i]) || math.IsInf(c.YValues[i], 0) {
			return fmt.Errorf("slide %d scatter y value %d must be finite", slideIndex, i+1)
		}
	}
	return nil
}

func (c ScatterChart) GetCategories() []string {
	cats := make([]string, len(c.XValues))
	for i, v := range c.XValues {
		cats[i] = fmt.Sprintf("%g", v)
	}
	return cats
}

func (c ScatterChart) GetValues() []float64 {
	return c.YValues
}

func isScatterStyle(style string) bool {
	switch strings.TrimSpace(style) {
	case ScatterStyleMarker, ScatterStyleLineMarker, ScatterStyleSmoothMarker:
		return true
	default:
		return false
	}
}
