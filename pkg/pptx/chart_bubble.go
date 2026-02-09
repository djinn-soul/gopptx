package pptx

import (
	"fmt"
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

// BubbleChart is a bubble chart using x/y coordinates and bubble sizes.
type BubbleChart struct {
	Title                 string
	TitleOverlay          bool
	XValues               []float64
	YValues               []float64
	BubbleSizes           []float64
	X                     int64
	Y                     int64
	CX                    int64
	CY                    int64
	LineColor             string
	SeriesName            string
	ShowLegend            bool
	LegendPosition        string
	LegendOverlay         bool
	ShowDataLabels        bool
	ShowMajorGridlines    bool
	CategoryAxisTitle     string
	ValueAxisTitle        string
	ValueFormat           string
	ValueAxisCrossBetween string
	MinValue              *float64
	MaxValue              *float64
	BubbleScale           int
}

func NewBubbleChart(xValues []float64, yValues []float64, bubbleSizes []float64) BubbleChart {
	xs := make([]float64, len(xValues))
	copy(xs, xValues)
	ys := make([]float64, len(yValues))
	copy(ys, yValues)
	bs := make([]float64, len(bubbleSizes))
	copy(bs, bubbleSizes)
	return BubbleChart{
		Title:                 "Chart",
		XValues:               xs,
		YValues:               ys,
		BubbleSizes:           bs,
		X:                     685800,
		Y:                     1800000,
		CX:                    7772400,
		CY:                    4114800,
		LineColor:             "4F81BD",
		SeriesName:            "Series 1",
		ShowLegend:            false,
		LegendPosition:        LegendPositionRight,
		ShowDataLabels:        false,
		ShowMajorGridlines:    true,
		ValueFormat:           "General",
		ValueAxisCrossBetween: ValueAxisCrossBetweenBetween,
		BubbleScale:           100,
	}
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
	c.LineColor = normalizeHexColor(color)
	return c
}

// ToChartSpec converts BubbleChart to internal XML spec.
func (c BubbleChart) ToChartSpec() *pptxxml.ChartSpec {
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindBubble,
		Title:                 c.Title,
		TitleOverlay:          c.TitleOverlay,
		XValues:               copyFloat64Slice(c.XValues),
		Values:                copyFloat64Slice(c.YValues),
		BubbleSizes:           copyFloat64Slice(c.BubbleSizes),
		X:                     c.X,
		Y:                     c.Y,
		CX:                    c.CX,
		CY:                    c.CY,
		Color:                 normalizeHexColor(c.LineColor),
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
		MinValue:              copyFloat64Pointer(c.MinValue),
		MaxValue:              copyFloat64Pointer(c.MaxValue),
		BubbleScale:           c.BubbleScale,
	}
}

// Validate checks the bubble chart for consistency.
func (c BubbleChart) Validate(slideIndex int) error {
	if c.X < 0 || c.Y < 0 {
		return fmt.Errorf("slide %d bubble chart position cannot be negative", slideIndex)
	}
	if c.CX <= 0 || c.CY <= 0 {
		return fmt.Errorf("slide %d bubble chart size must be > 0", slideIndex)
	}
	if strings.TrimSpace(c.Title) == "" {
		return fmt.Errorf("slide %d bubble chart title cannot be empty", slideIndex)
	}
	if len(c.XValues) == 0 {
		return fmt.Errorf("slide %d bubble chart must define at least one point", slideIndex)
	}
	if len(c.XValues) != len(c.YValues) || len(c.XValues) != len(c.BubbleSizes) {
		return fmt.Errorf("slide %d bubble chart x/y/size lengths must match", slideIndex)
	}
	if !isHexColor(c.LineColor) {
		return fmt.Errorf("slide %d bubble chart color must be 6-digit RGB hex", slideIndex)
	}
	if strings.TrimSpace(c.SeriesName) == "" {
		return fmt.Errorf("slide %d bubble chart series name cannot be empty", slideIndex)
	}
	if !isLegendPosition(c.LegendPosition) {
		return fmt.Errorf("slide %d bubble chart legend position must be one of r,l,t,b", slideIndex)
	}
	if strings.TrimSpace(c.ValueFormat) == "" {
		return fmt.Errorf("slide %d bubble chart value format cannot be empty", slideIndex)
	}
	if !isValueAxisCrossBetween(c.ValueAxisCrossBetween) {
		return fmt.Errorf("slide %d bubble chart value-axis crossBetween must be between or midCat", slideIndex)
	}
	if c.BubbleScale < 1 || c.BubbleScale > 300 {
		return fmt.Errorf("slide %d bubble chart scale must be between 1 and 300", slideIndex)
	}
	if err := validateValueRange(c.MinValue, c.MaxValue, slideIndex); err != nil {
		return err
	}
	for i := range c.XValues {
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
	}
	return nil
}

func validateBubbleChart(chart BubbleChart, slideIndex int) error {
	return chart.Validate(slideIndex)
}
