package pptx

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

const (
	RadarStyleMarker = "marker"
	RadarStyleFilled = "filled"
)

// RadarChart is a radar chart using marker style.
type RadarChart struct {
	Title                 string
	TitleOverlay          bool
	Categories            []string
	Values                []float64
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
	RadarStyle            string
}

func NewRadarChart(categories []string, values []float64) RadarChart {
	cats, vals := copyChartData(categories, values)
	return RadarChart{
		Title:                 "Chart",
		Categories:            cats,
		Values:                vals,
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
		RadarStyle:            RadarStyleMarker,
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

// ToChartSpec converts RadarChart to internal XML spec.
func (c RadarChart) ToChartSpec() *pptxxml.ChartSpec {
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindRadar,
		Title:                 c.Title,
		TitleOverlay:          c.TitleOverlay,
		Categories:            copyStringSlice(c.Categories),
		Values:                copyFloat64Slice(c.Values),
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
		RadarStyle:            c.RadarStyle,
	}
}

// Validate checks the radar chart for consistency.
func (c RadarChart) Validate(slideIndex int) error {
	if err := validateChartCore(
		slideIndex,
		c.Title,
		c.Categories,
		c.Values,
		c.X,
		c.Y,
		c.CX,
		c.CY,
	); err != nil {
		return err
	}
	if !isHexColor(c.LineColor) {
		return fmt.Errorf("slide %d radar chart color must be 6-digit RGB hex", slideIndex)
	}
	if strings.TrimSpace(c.SeriesName) == "" {
		return fmt.Errorf("slide %d radar chart series name cannot be empty", slideIndex)
	}
	if !isLegendPosition(c.LegendPosition) {
		return fmt.Errorf("slide %d radar chart legend position must be one of r,l,t,b", slideIndex)
	}
	if strings.TrimSpace(c.ValueFormat) == "" {
		return fmt.Errorf("slide %d radar chart value format cannot be empty", slideIndex)
	}
	if !isValueAxisCrossBetween(c.ValueAxisCrossBetween) {
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

// Validate checks the radar chart for consistency.
func (c RadarFilledChart) Validate(slideIndex int) error {
	return c.RadarChart.Validate(slideIndex)
}

func validateRadarChart(chart RadarChart, slideIndex int) error {
	return chart.Validate(slideIndex)
}
