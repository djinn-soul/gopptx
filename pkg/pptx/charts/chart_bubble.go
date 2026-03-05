package charts

import (
	"fmt"
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

const defaultBubbleScale = 100

// BubbleChart is a bubble chart using x/y coordinates and bubble sizes.
type BubbleChart struct {
	Title                      string
	TitleOverlay               bool
	XValues                    []float64
	YValues                    []float64
	BubbleSizes                []float64
	X                          int64
	Y                          int64
	CX                         int64
	CY                         int64
	LineColor                  string
	SeriesName                 string
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
	BubbleScale                int

	// Accessibility
	AltText      string
	IsDecorative bool
}

func NewBubbleChart(xValues []float64, yValues []float64, bubbleSizes []float64) BubbleChart {
	xs := make([]float64, len(xValues))
	copy(xs, xValues)
	ys := make([]float64, len(yValues))
	copy(ys, yValues)
	bs := make([]float64, len(bubbleSizes))
	copy(bs, bubbleSizes)
	return BubbleChart{
		Title:                      "Chart",
		XValues:                    xs,
		YValues:                    ys,
		BubbleSizes:                bs,
		X:                          defaultChartX,
		Y:                          defaultChartY,
		CX:                         defaultChartCX,
		CY:                         defaultChartCY,
		LineColor:                  "4F81BD",
		SeriesName:                 "Series 1",
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
		BubbleScale:                defaultBubbleScale,
	}
}

// WithAltText sets the alternative text for accessibility.
func (c BubbleChart) WithAltText(text string) BubbleChart {
	c.AltText = text
	return c
}

// WithDecorative marks the chart as decorative (ignored by screen readers).
func (c BubbleChart) WithDecorative(enabled bool) BubbleChart {
	c.IsDecorative = enabled
	return c
}

func (c BubbleChart) Position(x int64, y int64) BubbleChart {
	c.X = x
	c.Y = y
	return c
}

func (c BubbleChart) Size(cx int64, cy int64) BubbleChart {
	c.CX = cx
	c.CY = cy
	return c
}

func (c BubbleChart) WithTitle(title string) BubbleChart {
	c.Title = title
	return c
}

func (c BubbleChart) WithLineColor(color string) BubbleChart {
	c.LineColor = NormalizeHexColor(color)
	return c
}

// ToChartSpec converts BubbleChart to internal XML spec.
func (c BubbleChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := &pptxxml.ChartSpec{
		Kind:                       pptxxml.ChartKindBubble,
		Title:                      c.Title,
		TitleOverlay:               c.TitleOverlay,
		XValues:                    CopyFloat64Slice(c.XValues),
		Values:                     CopyFloat64Slice(c.YValues),
		BubbleSizes:                CopyFloat64Slice(c.BubbleSizes),
		X:                          c.X,
		Y:                          c.Y,
		CX:                         c.CX,
		CY:                         c.CY,
		Color:                      NormalizeHexColor(c.LineColor),
		SeriesName:                 c.SeriesName,
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
		BubbleScale:                c.BubbleScale,
		AltText:                    c.AltText,
		IsDecorative:               c.IsDecorative,
	}
	applyDataLabelSettings(spec, c.DataLabels)
	return spec
}

// Validate checks the bubble chart for consistency.
func (c BubbleChart) Validate(slideIndex int) error {
	if !c.IsDecorative && len(c.AltText) > common.MaxAltTextLength {
		return fmt.Errorf("slide %d bubble chart alt text exceeds %d characters", slideIndex, common.MaxAltTextLength)
	}
	if err := c.validateCoordinates(slideIndex); err != nil {
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

func (c BubbleChart) validateCoordinates(slideIndex int) error {
	if c.X < 0 || c.Y < 0 {
		return fmt.Errorf("slide %d bubble chart position cannot be negative", slideIndex)
	}
	if c.CX <= 0 || c.CY <= 0 {
		return fmt.Errorf("slide %d bubble chart size must be > 0", slideIndex)
	}
	return nil
}

func (c BubbleChart) validateMetadata(slideIndex int) error {
	if strings.TrimSpace(c.Title) == "" {
		return fmt.Errorf("slide %d bubble chart title cannot be empty", slideIndex)
	}
	if !IsHexColor(c.LineColor) {
		return fmt.Errorf("slide %d bubble chart color must be 6-digit RGB hex", slideIndex)
	}
	if strings.TrimSpace(c.SeriesName) == "" {
		return fmt.Errorf("slide %d bubble chart series name cannot be empty", slideIndex)
	}
	if !IsLegendPosition(c.LegendPosition) {
		return fmt.Errorf("slide %d bubble chart legend position must be one of r,l,t,b", slideIndex)
	}
	if !IsDataLabelPosition(c.DataLabels.Position) {
		return fmt.Errorf(
			"slide %d bubble chart data-label position must be ctr,inEnd,inBase,outEnd,bestFit,l,r,t,or b",
			slideIndex,
		)
	}
	if strings.TrimSpace(c.ValueFormat) == "" {
		return fmt.Errorf("slide %d bubble chart value format cannot be empty", slideIndex)
	}
	if !IsAxisTickLabelPosition(c.CategoryTickLabelPosition) {
		return fmt.Errorf(
			"slide %d bubble chart x-axis tick label position must be nextTo, low, high, or none",
			slideIndex,
		)
	}
	if !IsAxisTickLabelPosition(c.ValueTickLabelPosition) {
		return fmt.Errorf(
			"slide %d bubble chart y-axis tick label position must be nextTo, low, high, or none",
			slideIndex,
		)
	}
	if !IsAxisCrosses(c.CategoryAxisCrosses) {
		return fmt.Errorf("slide %d bubble chart x-axis crosses must be autoZero, min, or max", slideIndex)
	}
	if !IsAxisCrosses(c.ValueAxisCrosses) {
		return fmt.Errorf("slide %d bubble chart y-axis crosses must be autoZero, min, or max", slideIndex)
	}
	if !IsValueAxisCrossBetween(c.ValueAxisCrossBetween) {
		return fmt.Errorf("slide %d bubble chart value-axis crossBetween must be between or midCat", slideIndex)
	}
	if c.BubbleScale < 1 || c.BubbleScale > 300 {
		return fmt.Errorf("slide %d bubble chart scale must be between 1 and 300", slideIndex)
	}
	return nil
}

func (c BubbleChart) validatePoints(slideIndex int) error {
	if len(c.XValues) == 0 {
		return fmt.Errorf("slide %d bubble chart must define at least one point", slideIndex)
	}
	if len(c.XValues) != len(c.YValues) || len(c.XValues) != len(c.BubbleSizes) {
		return fmt.Errorf("slide %d bubble chart x/y/size lengths must match", slideIndex)
	}
	for i := range c.XValues {
		if err := c.validatePoint(slideIndex, i); err != nil {
			return err
		}
	}
	return nil
}

func (c BubbleChart) validatePoint(slideIndex, i int) error {
	if math.IsNaN(c.XValues[i]) || math.IsInf(c.XValues[i], 0) {
		return fmt.Errorf("slide %d bubble x value %d must be finite", slideIndex, i+1)
	}
	if math.IsNaN(c.YValues[i]) || math.IsInf(c.YValues[i], 0) {
		return fmt.Errorf("slide %d bubble y value %d must be finite", slideIndex, i+1)
	}
	if math.IsNaN(c.BubbleSizes[i]) || math.IsInf(c.BubbleSizes[i], 0) {
		return fmt.Errorf("slide %d bubble size %d must be finite", slideIndex, i+1)
	}
	if c.BubbleSizes[i] < 0 {
		return fmt.Errorf("slide %d bubble size %d cannot be negative", slideIndex, i+1)
	}
	return nil
}

func (c BubbleChart) GetCategories() []string {
	// Browse compliance: Convert XValues to string for categories
	cats := make([]string, len(c.XValues))
	for i, v := range c.XValues {
		cats[i] = fmt.Sprintf("%g", v)
	}
	return cats
}

func (c BubbleChart) GetValues() []float64 {
	return c.YValues
}
