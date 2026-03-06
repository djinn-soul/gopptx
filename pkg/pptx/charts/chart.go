package charts

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

var hexColorPattern = regexp.MustCompile(`^[0-9A-Fa-f]{6}$`)

const (
	LegendPositionRight  = "r"
	LegendPositionLeft   = "l"
	LegendPositionTop    = "t"
	LegendPositionBottom = "b"

	ValueAxisCrossBetweenBetween     = "between"
	ValueAxisCrossBetweenMidCategory = "midCat"

	defaultChartX  = 685800
	defaultChartY  = 1800000
	defaultChartCX = 7772400
	defaultChartCY = 4114800
)

// BarChart is a simple categorical bar chart.
type BarChart struct {
	Title        string
	TitleOverlay bool
	Categories   []string
	Values       []float64
	X            styling.Length
	Y            styling.Length
	CX           styling.Length
	CY           styling.Length

	BarColor                   string
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

	// Accessibility
	AltText      string
	IsDecorative bool
}

// NewBarChart creates a bar chart with default layout and style.
func NewBarChart(categories []string, values []float64) BarChart {
	cats, vals := copyChartData(categories, values)
	return BarChart{
		Title:      "Chart",
		Categories: cats,
		Values:     vals,
		X:          styling.Emu(defaultChartX),
		Y:          styling.Emu(defaultChartY),
		CX:         styling.Emu(defaultChartCX),
		CY:         styling.Emu(defaultChartCY),

		BarColor:                   "4F81BD",
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
	}
}

// WithAltText sets the alternative text for accessibility.
func (c BarChart) WithAltText(text string) BarChart {
	c.AltText = text
	return c
}

// WithDecorative marks the chart as decorative (ignored by screen readers).
func (c BarChart) WithDecorative(enabled bool) BarChart {
	c.IsDecorative = enabled
	return c
}

// Position sets chart position in EMU.
func (c BarChart) Position(x styling.Length, y styling.Length) BarChart {
	c.X = x
	c.Y = y
	return c
}

// Size sets chart size in EMU.
func (c BarChart) Size(cx styling.Length, cy styling.Length) BarChart {
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
	c.BarColor = NormalizeHexColor(color)
	return c
}

// LineChart is a simple categorical line chart.
type LineChart struct {
	Title        string
	TitleOverlay bool
	Categories   []string
	Values       []float64
	X            styling.Length
	Y            styling.Length
	CX           styling.Length
	CY           styling.Length

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
	Smooth                     bool

	// Accessibility
	AltText      string
	IsDecorative bool
}

// NewLineChart creates a line chart with default layout and style.
func NewLineChart(categories []string, values []float64) LineChart {
	cats, vals := copyChartData(categories, values)
	return LineChart{
		Title:      "Chart",
		Categories: cats,
		Values:     vals,
		X:          styling.Emu(defaultChartX),
		Y:          styling.Emu(defaultChartY),
		CX:         styling.Emu(defaultChartCX),
		CY:         styling.Emu(defaultChartCY),

		LineColor:                  "C0504D",
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
		Smooth:                     false,
	}
}

// WithAltText sets the alternative text for accessibility.
func (c LineChart) WithAltText(text string) LineChart {
	c.AltText = text
	return c
}

// WithDecorative marks the chart as decorative (ignored by screen readers).
func (c LineChart) WithDecorative(enabled bool) LineChart {
	c.IsDecorative = enabled
	return c
}

// Position sets chart position in EMU.
func (c LineChart) Position(x styling.Length, y styling.Length) LineChart {
	c.X = x
	c.Y = y
	return c
}

// Size sets chart size in EMU.
func (c LineChart) Size(cx styling.Length, cy styling.Length) LineChart {
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
	c.LineColor = NormalizeHexColor(color)
	return c
}

// ToChartSpec converts BarChart to internal XML spec.
func (c BarChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := &pptxxml.ChartSpec{
		Kind:         pptxxml.ChartKindBar,
		Title:        c.Title,
		TitleOverlay: c.TitleOverlay,
		Categories:   CopyStringSlice(c.Categories),
		Values:       CopyFloat64Slice(c.Values),
		X:            c.X.Emu(),
		Y:            c.Y.Emu(),
		CX:           c.CX.Emu(),
		CY:           c.CY.Emu(),

		Color:                      NormalizeHexColor(c.BarColor),
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
		BarDir:                     "col",
		Grouping:                   "clustered",
		AltText:                    c.AltText,
		IsDecorative:               c.IsDecorative,
	}
	applyDataLabelSettings(spec, c.DataLabels)
	return spec
}

// Validate checks the bar chart for consistency.
func (c BarChart) Validate(slideIndex int) error {
	if !c.IsDecorative && len(c.AltText) > common.MaxAltTextLength {
		return fmt.Errorf("slide %d bar chart alt text exceeds %d characters", slideIndex, common.MaxAltTextLength)
	}
	return validateChartCommon(
		slideIndex, c.Title, c.Categories,
		c.Values, c.X, c.Y, c.CX, c.CY,
		false, c.BarColor, c.SeriesName,
		c.LegendPosition, c.ValueFormat,
		c.CategoryTickLabelPosition, c.ValueTickLabelPosition,
		c.CategoryAxisCrosses, c.ValueAxisCrosses,
		c.ValueAxisCrossBetween,
		c.DataLabels,
		c.MinValue, c.MaxValue, "bar",
	)
}

func (c BarChart) GetCategories() []string {
	return c.Categories
}

func (c BarChart) GetValues() []float64 {
	return c.Values
}

// ToChartSpec converts LineChart to internal XML spec.
func (c LineChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := &pptxxml.ChartSpec{
		Kind:         pptxxml.ChartKindLine,
		Title:        c.Title,
		TitleOverlay: c.TitleOverlay,
		Categories:   CopyStringSlice(c.Categories),
		Values:       CopyFloat64Slice(c.Values),
		X:            c.X.Emu(),
		Y:            c.Y.Emu(),
		CX:           c.CX.Emu(),
		CY:           c.CY.Emu(),

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
		Grouping:                   "standard",
		Smooth:                     c.Smooth,
		AltText:                    c.AltText,
		IsDecorative:               c.IsDecorative,
	}
	applyDataLabelSettings(spec, c.DataLabels)
	return spec
}

// Validate checks the line chart for consistency.
func (c LineChart) Validate(slideIndex int) error {
	if !c.IsDecorative && len(c.AltText) > common.MaxAltTextLength {
		return fmt.Errorf("slide %d line chart alt text exceeds %d characters", slideIndex, common.MaxAltTextLength)
	}
	return validateChartCommon(
		slideIndex, c.Title, c.Categories,
		c.Values, c.X, c.Y, c.CX, c.CY,
		true, c.LineColor, c.SeriesName,
		c.LegendPosition, c.ValueFormat,
		c.CategoryTickLabelPosition, c.ValueTickLabelPosition,
		c.CategoryAxisCrosses, c.ValueAxisCrosses,
		c.ValueAxisCrossBetween,
		c.DataLabels,
		c.MinValue, c.MaxValue, "line",
	)
}

func validateChartCommon(
	slideIndex int,
	title string,
	categories []string,
	values []float64,
	x, y, cx, cy styling.Length,
	allowNegative bool,
	color string,
	seriesName string,
	legendPosition string,
	valueFormat string,
	categoryTickLabelPosition string,
	valueTickLabelPosition string,
	categoryAxisCrosses string,
	valueAxisCrosses string,
	crossBetween string,
	dataLabels DataLabelSettings,
	minValue *float64,
	maxValue *float64,
	chartType string,
) error {
	if err := validateChartCore(slideIndex, title, categories, values, x, y, cx, cy, allowNegative); err != nil {
		return err
	}
	if !IsHexColor(color) {
		return fmt.Errorf("slide %d %s chart color must be 6-digit RGB hex", slideIndex, chartType)
	}
	if strings.TrimSpace(seriesName) == "" {
		return fmt.Errorf("slide %d %s chart series name cannot be empty", slideIndex, chartType)
	}
	if !IsLegendPosition(legendPosition) {
		return fmt.Errorf("slide %d %s chart legend position must be one of r,l,t,b", slideIndex, chartType)
	}
	if !IsDataLabelPosition(dataLabels.Position) {
		return fmt.Errorf(
			"slide %d %s chart data-label position must be ctr,inEnd,inBase,outEnd,bestFit,l,r,t,or b",
			slideIndex,
			chartType,
		)
	}
	if strings.TrimSpace(valueFormat) == "" {
		return fmt.Errorf("slide %d %s chart value format cannot be empty", slideIndex, chartType)
	}
	if !IsAxisTickLabelPosition(categoryTickLabelPosition) {
		return fmt.Errorf(
			"slide %d %s chart category-axis tick label position must be nextTo, low, high, or none",
			slideIndex,
			chartType,
		)
	}
	if !IsAxisTickLabelPosition(valueTickLabelPosition) {
		return fmt.Errorf(
			"slide %d %s chart value-axis tick label position must be nextTo, low, high, or none",
			slideIndex,
			chartType,
		)
	}
	if !IsAxisCrosses(categoryAxisCrosses) {
		return fmt.Errorf(
			"slide %d %s chart category-axis crosses must be autoZero, min, or max",
			slideIndex,
			chartType,
		)
	}
	if !IsAxisCrosses(valueAxisCrosses) {
		return fmt.Errorf(
			"slide %d %s chart value-axis crosses must be autoZero, min, or max",
			slideIndex,
			chartType,
		)
	}
	if !IsValueAxisCrossBetween(crossBetween) {
		return fmt.Errorf("slide %d %s chart value-axis crossBetween must be between or midCat", slideIndex, chartType)
	}
	return validateValueRange(minValue, maxValue, slideIndex)
}

func (c LineChart) GetCategories() []string {
	return c.Categories
}

func (c LineChart) GetValues() []float64 {
	return c.Values
}

func validateChartCore(
	slideIndex int,
	title string,
	categories []string,
	values []float64,
	x styling.Length,
	y styling.Length,
	cx styling.Length,
	cy styling.Length,
	allowNegative bool,
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

	hasNonZero := false
	for i := range categories {
		if strings.TrimSpace(categories[i]) == "" {
			return fmt.Errorf("slide %d chart category %d cannot be empty", slideIndex, i+1)
		}
		if !allowNegative && values[i] < 0 {
			return fmt.Errorf("slide %d chart value %d cannot be negative", slideIndex, i+1)
		}
		if values[i] != 0 {
			hasNonZero = true
		}
	}
	if !hasNonZero {
		return fmt.Errorf("slide %d chart requires at least one non-zero value", slideIndex)
	}
	return nil
}

func copyChartData(categories []string, values []float64) ([]string, []float64) {
	return CopyStringSlice(categories), CopyFloat64Slice(values)
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
