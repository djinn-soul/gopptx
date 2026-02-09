package pptx

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

// AreaChart is a simple categorical area chart.
type AreaChart struct {
	Title                 string
	TitleOverlay          bool
	Categories            []string
	Values                []float64
	X                     int64
	Y                     int64
	CX                    int64
	CY                    int64
	AreaColor             string
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
}

// NewAreaChart creates an area chart with default layout and style.
func NewAreaChart(categories []string, values []float64) AreaChart {
	cats, vals := copyChartData(categories, values)
	return AreaChart{
		Title:                 "Chart",
		Categories:            cats,
		Values:                vals,
		X:                     685800,
		Y:                     1800000,
		CX:                    7772400,
		CY:                    4114800,
		AreaColor:             "9BBB59",
		SeriesName:            "Series 1",
		ShowLegend:            false,
		LegendPosition:        LegendPositionRight,
		ShowDataLabels:        false,
		ShowMajorGridlines:    true,
		ValueFormat:           "General",
		ValueAxisCrossBetween: ValueAxisCrossBetweenBetween,
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

// ToChartSpec converts AreaChart to internal XML spec.
func (c AreaChart) ToChartSpec() *pptxxml.ChartSpec {
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindArea,
		Title:                 c.Title,
		TitleOverlay:          c.TitleOverlay,
		Categories:            copyStringSlice(c.Categories),
		Values:                copyFloat64Slice(c.Values),
		X:                     c.X,
		Y:                     c.Y,
		CX:                    c.CX,
		CY:                    c.CY,
		Color:                 normalizeHexColor(c.AreaColor),
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
		Grouping:              "standard",
	}
}

// Validate checks the area chart for consistency.
func (c AreaChart) Validate(slideIndex int) error {
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
	if !isHexColor(c.AreaColor) {
		return fmt.Errorf("slide %d area chart color must be 6-digit RGB hex", slideIndex)
	}
	if strings.TrimSpace(c.SeriesName) == "" {
		return fmt.Errorf("slide %d area chart series name cannot be empty", slideIndex)
	}
	if !isLegendPosition(c.LegendPosition) {
		return fmt.Errorf("slide %d area chart legend position must be one of r,l,t,b", slideIndex)
	}
	if strings.TrimSpace(c.ValueFormat) == "" {
		return fmt.Errorf("slide %d area chart value format cannot be empty", slideIndex)
	}
	if !isValueAxisCrossBetween(c.ValueAxisCrossBetween) {
		return fmt.Errorf("slide %d area chart value-axis crossBetween must be between or midCat", slideIndex)
	}
	if err := validateValueRange(c.MinValue, c.MaxValue, slideIndex); err != nil {
		return err
	}
	return nil
}

func validateAreaChart(chart AreaChart, slideIndex int) error {
	return chart.Validate(slideIndex)
}
