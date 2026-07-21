package charts

import (
	"fmt"
	"regexp"

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

	// Shared constructor defaults across every chart kind.
	defaultChartTitle       = "Chart"
	defaultChartSeriesColor = "4F81BD"
	defaultChartValueFormat = "General"
	defaultChartSeriesName  = "Series 1"

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
		Title:      defaultChartTitle,
		Categories: cats,
		Values:     vals,
		X:          styling.Emu(defaultChartX),
		Y:          styling.Emu(defaultChartY),
		CX:         styling.Emu(defaultChartCX),
		CY:         styling.Emu(defaultChartCY),

		BarColor:                   defaultChartSeriesColor,
		SeriesName:                 defaultChartSeriesName,
		ShowLegend:                 false,
		LegendPosition:             LegendPositionRight,
		ShowDataLabels:             false,
		ShowMajorGridlines:         true,
		ShowCategoryMajorGridlines: false,
		CategoryTickLabelPosition:  AxisTickLabelPositionNextTo,
		ValueTickLabelPosition:     AxisTickLabelPositionNextTo,
		CategoryAxisCrosses:        AxisCrossesAutoZero,
		ValueAxisCrosses:           AxisCrossesAutoZero,
		ValueFormat:                defaultChartValueFormat,
		ValueAxisCrossBetween:      ValueAxisCrossBetweenBetween,
	}
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
		Title:      defaultChartTitle,
		Categories: cats,
		Values:     vals,
		X:          styling.Emu(defaultChartX),
		Y:          styling.Emu(defaultChartY),
		CX:         styling.Emu(defaultChartCX),
		CY:         styling.Emu(defaultChartCY),

		LineColor:                  defaultChartSeriesColor,
		SeriesName:                 defaultChartSeriesName,
		ShowLegend:                 false,
		LegendPosition:             LegendPositionRight,
		ShowDataLabels:             false,
		ShowMajorGridlines:         true,
		ShowCategoryMajorGridlines: false,
		CategoryTickLabelPosition:  AxisTickLabelPositionNextTo,
		ValueTickLabelPosition:     AxisTickLabelPositionNextTo,
		CategoryAxisCrosses:        AxisCrossesAutoZero,
		ValueAxisCrosses:           AxisCrossesAutoZero,
		ValueFormat:                defaultChartValueFormat,
		ValueAxisCrossBetween:      ValueAxisCrossBetweenBetween,
		Smooth:                     true,
	}
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

func (c LineChart) GetCategories() []string {
	return c.Categories
}

func (c LineChart) GetValues() []float64 {
	return c.Values
}
