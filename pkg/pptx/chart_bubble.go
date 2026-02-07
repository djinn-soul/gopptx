package pptx

import (
	"fmt"
	"math"
	"strings"
)

// BubbleChart is a bubble chart using x/y coordinates and bubble sizes.
type BubbleChart struct {
	Title              string
	XValues            []float64
	YValues            []float64
	BubbleSizes        []float64
	X                  int64
	Y                  int64
	CX                 int64
	CY                 int64
	LineColor          string
	SeriesName         string
	ShowLegend         bool
	LegendPosition     string
	ShowDataLabels     bool
	ShowMajorGridlines bool
	CategoryAxisTitle  string
	ValueAxisTitle     string
	ValueFormat        string
	MinValue           *float64
	MaxValue           *float64
	BubbleScale        int
}

func NewBubbleChart(xValues []float64, yValues []float64, bubbleSizes []float64) BubbleChart {
	xs := make([]float64, len(xValues))
	copy(xs, xValues)
	ys := make([]float64, len(yValues))
	copy(ys, yValues)
	bs := make([]float64, len(bubbleSizes))
	copy(bs, bubbleSizes)
	return BubbleChart{
		Title:              "Chart",
		XValues:            xs,
		YValues:            ys,
		BubbleSizes:        bs,
		X:                  685800,
		Y:                  1800000,
		CX:                 7772400,
		CY:                 4114800,
		LineColor:          "4F81BD",
		SeriesName:         "Series 1",
		ShowLegend:         false,
		LegendPosition:     LegendPositionRight,
		ShowDataLabels:     false,
		ShowMajorGridlines: true,
		ValueFormat:        "General",
		BubbleScale:        100,
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

func validateBubbleChart(chart BubbleChart, slideIndex int) error {
	if chart.X < 0 || chart.Y < 0 {
		return fmt.Errorf("slide %d bubble chart position cannot be negative", slideIndex)
	}
	if chart.CX <= 0 || chart.CY <= 0 {
		return fmt.Errorf("slide %d bubble chart size must be > 0", slideIndex)
	}
	if strings.TrimSpace(chart.Title) == "" {
		return fmt.Errorf("slide %d bubble chart title cannot be empty", slideIndex)
	}
	if len(chart.XValues) == 0 {
		return fmt.Errorf("slide %d bubble chart must define at least one point", slideIndex)
	}
	if len(chart.XValues) != len(chart.YValues) || len(chart.XValues) != len(chart.BubbleSizes) {
		return fmt.Errorf("slide %d bubble chart x/y/size lengths must match", slideIndex)
	}
	if !isHexColor(chart.LineColor) {
		return fmt.Errorf("slide %d bubble chart color must be 6-digit RGB hex", slideIndex)
	}
	if strings.TrimSpace(chart.SeriesName) == "" {
		return fmt.Errorf("slide %d bubble chart series name cannot be empty", slideIndex)
	}
	if !isLegendPosition(chart.LegendPosition) {
		return fmt.Errorf("slide %d bubble chart legend position must be one of r,l,t,b", slideIndex)
	}
	if strings.TrimSpace(chart.ValueFormat) == "" {
		return fmt.Errorf("slide %d bubble chart value format cannot be empty", slideIndex)
	}
	if chart.BubbleScale < 1 || chart.BubbleScale > 300 {
		return fmt.Errorf("slide %d bubble chart scale must be between 1 and 300", slideIndex)
	}
	if err := validateValueRange(chart.MinValue, chart.MaxValue, slideIndex); err != nil {
		return err
	}
	for i := range chart.XValues {
		if math.IsNaN(chart.XValues[i]) || math.IsInf(chart.XValues[i], 0) {
			return fmt.Errorf("slide %d bubble x value %d must be finite", slideIndex, i+1)
		}
		if math.IsNaN(chart.YValues[i]) || math.IsInf(chart.YValues[i], 0) {
			return fmt.Errorf("slide %d bubble y value %d must be finite", slideIndex, i+1)
		}
		if math.IsNaN(chart.BubbleSizes[i]) || math.IsInf(chart.BubbleSizes[i], 0) {
			return fmt.Errorf("slide %d bubble size %d must be finite", slideIndex, i+1)
		}
		if chart.BubbleSizes[i] < 0 {
			return fmt.Errorf("slide %d bubble size %d cannot be negative", slideIndex, i+1)
		}
	}
	return nil
}
