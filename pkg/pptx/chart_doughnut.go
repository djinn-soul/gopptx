package pptx

import (
	"fmt"
	"strings"
)

// DoughnutChart is a simple categorical doughnut chart.
type DoughnutChart struct {
	Title          string
	Categories     []string
	Values         []float64
	X              int64
	Y              int64
	CX             int64
	CY             int64
	SeriesName     string
	ShowLegend     bool
	LegendPosition string
	ShowDataLabels bool
	HoleSize       int
}

// NewDoughnutChart creates a doughnut chart with default layout and style.
func NewDoughnutChart(categories []string, values []float64) DoughnutChart {
	cats, vals := copyChartData(categories, values)
	return DoughnutChart{
		Title:          "Chart",
		Categories:     cats,
		Values:         vals,
		X:              685800,
		Y:              1800000,
		CX:             7772400,
		CY:             4114800,
		SeriesName:     "Series 1",
		ShowLegend:     false,
		LegendPosition: LegendPositionRight,
		ShowDataLabels: false,
		HoleSize:       50,
	}
}

// Position sets chart position in EMU.
func (c DoughnutChart) Position(x int64, y int64) DoughnutChart {
	c.X = x
	c.Y = y
	return c
}

// Size sets chart size in EMU.
func (c DoughnutChart) Size(cx int64, cy int64) DoughnutChart {
	c.CX = cx
	c.CY = cy
	return c
}

// WithTitle sets the chart title.
func (c DoughnutChart) WithTitle(title string) DoughnutChart {
	c.Title = title
	return c
}

func validateDoughnutChart(chart DoughnutChart, slideIndex int) error {
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
	if strings.TrimSpace(chart.SeriesName) == "" {
		return fmt.Errorf("slide %d doughnut chart series name cannot be empty", slideIndex)
	}
	if !isLegendPosition(chart.LegendPosition) {
		return fmt.Errorf("slide %d doughnut chart legend position must be one of r,l,t,b", slideIndex)
	}
	if chart.HoleSize < 10 || chart.HoleSize > 90 {
		return fmt.Errorf("slide %d doughnut chart hole size must be between 10 and 90", slideIndex)
	}
	return nil
}
