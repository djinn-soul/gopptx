package pptx

import (
	"fmt"
	"strings"
)

// AreaChart is a simple categorical area chart.
type AreaChart struct {
	Title              string
	Categories         []string
	Values             []float64
	X                  int64
	Y                  int64
	CX                 int64
	CY                 int64
	AreaColor          string
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
}

// NewAreaChart creates an area chart with default layout and style.
func NewAreaChart(categories []string, values []float64) AreaChart {
	cats, vals := copyChartData(categories, values)
	return AreaChart{
		Title:              "Chart",
		Categories:         cats,
		Values:             vals,
		X:                  685800,
		Y:                  1800000,
		CX:                 7772400,
		CY:                 4114800,
		AreaColor:          "9BBB59",
		SeriesName:         "Series 1",
		ShowLegend:         false,
		LegendPosition:     LegendPositionRight,
		ShowDataLabels:     false,
		ShowMajorGridlines: true,
		ValueFormat:        "General",
	}
}

// Position sets chart position in EMU.
func (c AreaChart) Position(x int64, y int64) AreaChart {
	c.X = x
	c.Y = y
	return c
}

// Size sets chart size in EMU.
func (c AreaChart) Size(cx int64, cy int64) AreaChart {
	c.CX = cx
	c.CY = cy
	return c
}

// WithTitle sets the chart title.
func (c AreaChart) WithTitle(title string) AreaChart {
	c.Title = title
	return c
}

// WithAreaColor sets the area fill color using RGB hex.
func (c AreaChart) WithAreaColor(color string) AreaChart {
	c.AreaColor = normalizeHexColor(color)
	return c
}

func validateAreaChart(chart AreaChart, slideIndex int) error {
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
	if !isHexColor(chart.AreaColor) {
		return fmt.Errorf("slide %d area chart color must be 6-digit RGB hex", slideIndex)
	}
	if strings.TrimSpace(chart.SeriesName) == "" {
		return fmt.Errorf("slide %d area chart series name cannot be empty", slideIndex)
	}
	if !isLegendPosition(chart.LegendPosition) {
		return fmt.Errorf("slide %d area chart legend position must be one of r,l,t,b", slideIndex)
	}
	if strings.TrimSpace(chart.ValueFormat) == "" {
		return fmt.Errorf("slide %d area chart value format cannot be empty", slideIndex)
	}
	if err := validateValueRange(chart.MinValue, chart.MaxValue, slideIndex); err != nil {
		return err
	}
	return nil
}
