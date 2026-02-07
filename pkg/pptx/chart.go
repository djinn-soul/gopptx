package pptx

import (
	"fmt"
	"regexp"
	"strings"
)

var hexColorPattern = regexp.MustCompile(`^[0-9A-F]{6}$`)

const (
	LegendPositionRight  = "r"
	LegendPositionLeft   = "l"
	LegendPositionTop    = "t"
	LegendPositionBottom = "b"

	ValueAxisCrossBetweenBetween     = "between"
	ValueAxisCrossBetweenMidCategory = "midCat"
)

// BarChart is a simple categorical bar chart.
type BarChart struct {
	Title                 string
	TitleOverlay          bool
	Categories            []string
	Values                []float64
	X                     int64
	Y                     int64
	CX                    int64
	CY                    int64
	BarColor              string
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

// NewBarChart creates a bar chart with default layout and style.
func NewBarChart(categories []string, values []float64) BarChart {
	cats, vals := copyChartData(categories, values)
	return BarChart{
		Title:                 "Chart",
		Categories:            cats,
		Values:                vals,
		X:                     685800,
		Y:                     1800000,
		CX:                    7772400,
		CY:                    4114800,
		BarColor:              "4F81BD",
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
func (c BarChart) Position(x int64, y int64) BarChart {
	c.X = x
	c.Y = y
	return c
}

// Size sets chart size in EMU.
func (c BarChart) Size(cx int64, cy int64) BarChart {
	c.CX = cx
	c.CY = cy
	return c
}

// WithTitle sets the chart title.
func (c BarChart) WithTitle(title string) BarChart {
	c.Title = title
	return c
}

// WithBarColor sets the bar fill color using RGB hex.
func (c BarChart) WithBarColor(color string) BarChart {
	c.BarColor = normalizeHexColor(color)
	return c
}

// LineChart is a simple categorical line chart.
type LineChart struct {
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
	Smooth                bool
}

// NewLineChart creates a line chart with default layout and style.
func NewLineChart(categories []string, values []float64) LineChart {
	cats, vals := copyChartData(categories, values)
	return LineChart{
		Title:                 "Chart",
		Categories:            cats,
		Values:                vals,
		X:                     685800,
		Y:                     1800000,
		CX:                    7772400,
		CY:                    4114800,
		LineColor:             "C0504D",
		SeriesName:            "Series 1",
		ShowLegend:            false,
		LegendPosition:        LegendPositionRight,
		ShowDataLabels:        false,
		ShowMajorGridlines:    true,
		ValueFormat:           "General",
		ValueAxisCrossBetween: ValueAxisCrossBetweenBetween,
		Smooth:                false,
	}
}

// Position sets chart position in EMU.
func (c LineChart) Position(x int64, y int64) LineChart {
	c.X = x
	c.Y = y
	return c
}

// Size sets chart size in EMU.
func (c LineChart) Size(cx int64, cy int64) LineChart {
	c.CX = cx
	c.CY = cy
	return c
}

// WithTitle sets the chart title.
func (c LineChart) WithTitle(title string) LineChart {
	c.Title = title
	return c
}

// WithLineColor sets the line color using RGB hex.
func (c LineChart) WithLineColor(color string) LineChart {
	c.LineColor = normalizeHexColor(color)
	return c
}

func validateBarChart(chart BarChart, slideIndex int) error {
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
	if !isHexColor(chart.BarColor) {
		return fmt.Errorf("slide %d bar chart color must be 6-digit RGB hex", slideIndex)
	}
	if strings.TrimSpace(chart.SeriesName) == "" {
		return fmt.Errorf("slide %d bar chart series name cannot be empty", slideIndex)
	}
	if !isLegendPosition(chart.LegendPosition) {
		return fmt.Errorf("slide %d bar chart legend position must be one of r,l,t,b", slideIndex)
	}
	if strings.TrimSpace(chart.ValueFormat) == "" {
		return fmt.Errorf("slide %d bar chart value format cannot be empty", slideIndex)
	}
	if !isValueAxisCrossBetween(chart.ValueAxisCrossBetween) {
		return fmt.Errorf("slide %d bar chart value-axis crossBetween must be between or midCat", slideIndex)
	}
	if err := validateValueRange(chart.MinValue, chart.MaxValue, slideIndex); err != nil {
		return err
	}
	return nil
}

func validateLineChart(chart LineChart, slideIndex int) error {
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
		return fmt.Errorf("slide %d line chart color must be 6-digit RGB hex", slideIndex)
	}
	if strings.TrimSpace(chart.SeriesName) == "" {
		return fmt.Errorf("slide %d line chart series name cannot be empty", slideIndex)
	}
	if !isLegendPosition(chart.LegendPosition) {
		return fmt.Errorf("slide %d line chart legend position must be one of r,l,t,b", slideIndex)
	}
	if strings.TrimSpace(chart.ValueFormat) == "" {
		return fmt.Errorf("slide %d line chart value format cannot be empty", slideIndex)
	}
	if !isValueAxisCrossBetween(chart.ValueAxisCrossBetween) {
		return fmt.Errorf("slide %d line chart value-axis crossBetween must be between or midCat", slideIndex)
	}
	if err := validateValueRange(chart.MinValue, chart.MaxValue, slideIndex); err != nil {
		return err
	}
	return nil
}

func validateChartCore(
	slideIndex int,
	title string,
	categories []string,
	values []float64,
	x int64,
	y int64,
	cx int64,
	cy int64,
) error {
	if x < 0 || y < 0 {
		return fmt.Errorf("slide %d chart position cannot be negative", slideIndex)
	}
	if cx <= 0 || cy <= 0 {
		return fmt.Errorf("slide %d chart size must be > 0", slideIndex)
	}
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("slide %d chart title cannot be empty", slideIndex)
	}
	if len(categories) == 0 {
		return fmt.Errorf("slide %d chart must define at least one category", slideIndex)
	}
	if len(categories) != len(values) {
		return fmt.Errorf(
			"slide %d chart category/value length mismatch (%d vs %d)",
			slideIndex,
			len(categories),
			len(values),
		)
	}

	hasPositive := false
	for i := range categories {
		if strings.TrimSpace(categories[i]) == "" {
			return fmt.Errorf("slide %d chart category %d cannot be empty", slideIndex, i+1)
		}
		if values[i] < 0 {
			return fmt.Errorf("slide %d chart value %d cannot be negative", slideIndex, i+1)
		}
		if values[i] > 0 {
			hasPositive = true
		}
	}
	if !hasPositive {
		return fmt.Errorf("slide %d chart requires at least one positive value", slideIndex)
	}
	return nil
}

func copyChartData(categories []string, values []float64) ([]string, []float64) {
	cats := make([]string, len(categories))
	copy(cats, categories)
	vals := make([]float64, len(values))
	copy(vals, values)
	return cats, vals
}

func validateValueRange(minValue *float64, maxValue *float64, slideIndex int) error {
	if (minValue == nil) != (maxValue == nil) {
		return fmt.Errorf("slide %d chart value range requires both min and max", slideIndex)
	}
	if minValue != nil && maxValue != nil && *minValue >= *maxValue {
		return fmt.Errorf("slide %d chart value range requires min < max", slideIndex)
	}
	return nil
}
