package pptx

import (
	"fmt"
	"math"
	"strings"
)

const (
	ScatterStyleMarker       = "marker"
	ScatterStyleLineMarker   = "lineMarker"
	ScatterStyleSmoothMarker = "smoothMarker"
)

// ScatterChart is a simple XY scatter chart.
type ScatterChart struct {
	Title              string
	XValues            []float64
	YValues            []float64
	X                  int64
	Y                  int64
	CX                 int64
	CY                 int64
	LineColor          string
	SeriesName         string
	ScatterStyle       string
	ShowLegend         bool
	LegendPosition     string
	ShowDataLabels     bool
	ShowMajorGridlines bool
	CategoryAxisTitle  string
	ValueAxisTitle     string
	ValueFormat        string
	MinValue           *float64
	MaxValue           *float64
}

// NewScatterChart creates a scatter chart with default layout and style.
func NewScatterChart(xValues []float64, yValues []float64) ScatterChart {
	xs := make([]float64, len(xValues))
	copy(xs, xValues)
	ys := make([]float64, len(yValues))
	copy(ys, yValues)
	return ScatterChart{
		Title:              "Chart",
		XValues:            xs,
		YValues:            ys,
		X:                  685800,
		Y:                  1800000,
		CX:                 7772400,
		CY:                 4114800,
		LineColor:          "4F81BD",
		SeriesName:         "Series 1",
		ScatterStyle:       ScatterStyleMarker,
		ShowLegend:         false,
		LegendPosition:     LegendPositionRight,
		ShowDataLabels:     false,
		ShowMajorGridlines: true,
		ValueFormat:        "General",
	}
}

// Position sets chart position in EMU.
func (c ScatterChart) Position(x int64, y int64) ScatterChart {
	c.X = x
	c.Y = y
	return c
}

// Size sets chart size in EMU.
func (c ScatterChart) Size(cx int64, cy int64) ScatterChart {
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
	c.LineColor = normalizeHexColor(color)
	return c
}

func validateScatterChart(chart ScatterChart, slideIndex int) error {
	if chart.X < 0 || chart.Y < 0 {
		return fmt.Errorf("slide %d scatter chart position cannot be negative", slideIndex)
	}
	if chart.CX <= 0 || chart.CY <= 0 {
		return fmt.Errorf("slide %d scatter chart size must be > 0", slideIndex)
	}
	if strings.TrimSpace(chart.Title) == "" {
		return fmt.Errorf("slide %d scatter chart title cannot be empty", slideIndex)
	}
	if len(chart.XValues) == 0 {
		return fmt.Errorf("slide %d scatter chart must define at least one point", slideIndex)
	}
	if len(chart.XValues) != len(chart.YValues) {
		return fmt.Errorf(
			"slide %d scatter chart x/y length mismatch (%d vs %d)",
			slideIndex,
			len(chart.XValues),
			len(chart.YValues),
		)
	}
	if !isHexColor(chart.LineColor) {
		return fmt.Errorf("slide %d scatter chart color must be 6-digit RGB hex", slideIndex)
	}
	if strings.TrimSpace(chart.SeriesName) == "" {
		return fmt.Errorf("slide %d scatter chart series name cannot be empty", slideIndex)
	}
	if !isLegendPosition(chart.LegendPosition) {
		return fmt.Errorf("slide %d scatter chart legend position must be one of r,l,t,b", slideIndex)
	}
	if !isScatterStyle(chart.ScatterStyle) {
		return fmt.Errorf(
			"slide %d scatter style must be one of marker,lineMarker,smoothMarker",
			slideIndex,
		)
	}
	if strings.TrimSpace(chart.ValueFormat) == "" {
		return fmt.Errorf("slide %d scatter chart value format cannot be empty", slideIndex)
	}
	if err := validateValueRange(chart.MinValue, chart.MaxValue, slideIndex); err != nil {
		return err
	}

	for i := range chart.XValues {
		if math.IsNaN(chart.XValues[i]) || math.IsInf(chart.XValues[i], 0) {
			return fmt.Errorf("slide %d scatter x value %d must be finite", slideIndex, i+1)
		}
		if math.IsNaN(chart.YValues[i]) || math.IsInf(chart.YValues[i], 0) {
			return fmt.Errorf("slide %d scatter y value %d must be finite", slideIndex, i+1)
		}
	}
	return nil
}

func isScatterStyle(style string) bool {
	switch strings.TrimSpace(style) {
	case ScatterStyleMarker, ScatterStyleLineMarker, ScatterStyleSmoothMarker:
		return true
	default:
		return false
	}
}
