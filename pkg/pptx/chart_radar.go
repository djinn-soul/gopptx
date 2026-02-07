package pptx

import (
	"fmt"
	"strings"
)

const (
	RadarStyleMarker = "marker"
	RadarStyleFilled = "filled"
)

// RadarChart is a radar chart using marker style.
type RadarChart struct {
	Title              string
	Categories         []string
	Values             []float64
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
	RadarStyle         string
}

func NewRadarChart(categories []string, values []float64) RadarChart {
	cats, vals := copyChartData(categories, values)
	return RadarChart{
		Title:              "Chart",
		Categories:         cats,
		Values:             vals,
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
		RadarStyle:         RadarStyleMarker,
	}
}

func (c RadarChart) Position(x int64, y int64) RadarChart {
	c.X = x
	c.Y = y
	return c
}

func (c RadarChart) Size(cx int64, cy int64) RadarChart {
	c.CX = cx
	c.CY = cy
	return c
}

func (c RadarChart) WithTitle(title string) RadarChart {
	c.Title = title
	return c
}

func (c RadarChart) WithLineColor(color string) RadarChart {
	c.LineColor = normalizeHexColor(color)
	return c
}

func validateRadarChart(chart RadarChart, slideIndex int) error {
	if err := validateChartCore(
		slideIndex,
		chart.Title,
		chart.Categories,
		chart.Values,
		chart.X,
		chart.Y,
		chart.CX,
		chart.CY,
	); err != nil {
		return err
	}
	if !isHexColor(chart.LineColor) {
		return fmt.Errorf("slide %d radar chart color must be 6-digit RGB hex", slideIndex)
	}
	if strings.TrimSpace(chart.SeriesName) == "" {
		return fmt.Errorf("slide %d radar chart series name cannot be empty", slideIndex)
	}
	if !isLegendPosition(chart.LegendPosition) {
		return fmt.Errorf("slide %d radar chart legend position must be one of r,l,t,b", slideIndex)
	}
	if strings.TrimSpace(chart.ValueFormat) == "" {
		return fmt.Errorf("slide %d radar chart value format cannot be empty", slideIndex)
	}
	if chart.RadarStyle != RadarStyleMarker && chart.RadarStyle != RadarStyleFilled {
		return fmt.Errorf("slide %d radar chart style must be marker or filled", slideIndex)
	}
	if err := validateValueRange(chart.MinValue, chart.MaxValue, slideIndex); err != nil {
		return err
	}
	return nil
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

func validateRadarFilledChart(chart RadarFilledChart, slideIndex int) error {
	chart.RadarStyle = RadarStyleFilled
	return validateRadarChart(chart.RadarChart, slideIndex)
}
