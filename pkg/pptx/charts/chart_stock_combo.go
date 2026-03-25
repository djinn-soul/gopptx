package charts

import (
	"fmt"
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// StockHLCChart is a stock chart with High/Low/Close series.
type StockHLCChart struct {
	Title        string
	TitleOverlay bool
	Categories   []string
	HighValues   []float64
	LowValues    []float64
	CloseValues  []float64
	X            styling.Length
	Y            styling.Length
	CX           styling.Length
	CY           styling.Length

	ShowLegend            bool
	LegendPosition        string
	LegendOverlay         bool
	ShowDataLabels        bool
	DataLabels            DataLabelSettings
	ShowMajorGridlines    bool
	CategoryAxisTitle     string
	ValueAxisTitle        string
	ValueFormat           string
	ValueAxisCrossBetween string
	MinValue              *float64
	MaxValue              *float64

	// Accessibility
	AltText      string
	IsDecorative bool
}

func NewStockHLCChart(categories []string, high []float64, low []float64, closeValues []float64) StockHLCChart {
	cats := append([]string(nil), categories...)
	highVals := append([]float64(nil), high...)
	lowVals := append([]float64(nil), low...)
	closeVals := append([]float64(nil), closeValues...)
	return StockHLCChart{
		Title:       "Chart",
		Categories:  cats,
		HighValues:  highVals,
		LowValues:   lowVals,
		CloseValues: closeVals,
		X:           styling.Emu(defaultChartX),
		Y:           styling.Emu(defaultChartY),
		CX:          styling.Emu(defaultChartCX),
		CY:          styling.Emu(defaultChartCY),

		ShowLegend:            true,
		LegendPosition:        LegendPositionRight,
		ShowDataLabels:        false,
		ShowMajorGridlines:    true,
		ValueFormat:           "General",
		ValueAxisCrossBetween: ValueAxisCrossBetweenBetween,
	}
}

// WithAltText sets the alternative text for accessibility.
func (c StockHLCChart) WithAltText(text string) StockHLCChart {
	c.AltText = text
	return c
}

// WithDecorative marks the chart as decorative (ignored by screen readers).
func (c StockHLCChart) WithDecorative(enabled bool) StockHLCChart {
	c.IsDecorative = enabled
	return c
}

// ToChartSpec converts StockHLCChart to internal XML spec.
func (c StockHLCChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := &pptxxml.ChartSpec{
		Kind:         pptxxml.ChartKindStockHLC,
		Title:        c.Title,
		TitleOverlay: c.TitleOverlay,
		Categories:   CopyStringSlice(c.Categories),
		HighValues:   CopyFloat64Slice(c.HighValues),
		LowValues:    CopyFloat64Slice(c.LowValues),
		CloseValues:  CopyFloat64Slice(c.CloseValues),
		X:            c.X.Emu(),
		Y:            c.Y.Emu(),
		CX:           c.CX.Emu(),
		CY:           c.CY.Emu(),

		ShowLegend:            c.ShowLegend,
		LegendPosition:        c.LegendPosition,
		LegendOverlay:         c.LegendOverlay,
		ShowDataLabels:        c.ShowDataLabels,
		ShowMajorGridlines:    c.ShowMajorGridlines,
		CategoryAxisTitle:     c.CategoryAxisTitle,
		ValueAxisTitle:        c.ValueAxisTitle,
		ValueFormat:           c.ValueFormat,
		ValueAxisCrossBetween: c.ValueAxisCrossBetween,
		MinValue:              CopyFloat64Pointer(c.MinValue),
		MaxValue:              CopyFloat64Pointer(c.MaxValue),
		AltText:               c.AltText,
		IsDecorative:          c.IsDecorative,
	}
	applyDataLabelSettings(spec, c.DataLabels)
	return spec
}

// Validate checks the stock chart for consistency.
func (c StockHLCChart) Validate(slideIndex int) error {
	if !c.IsDecorative && len(c.AltText) > common.MaxAltTextLength {
		return fmt.Errorf("slide %d stock chart alt text exceeds %d characters", slideIndex, common.MaxAltTextLength)
	}
	if err := validateStockCore(
		slideIndex,
		c.Title,
		c.Categories,
		c.X,
		c.Y,
		c.CX,
		c.CY,

		c.ValueFormat,
		c.LegendPosition,
		c.ValueAxisCrossBetween,
		c.DataLabels,
		c.MinValue,
		c.MaxValue,
	); err != nil {
		return err
	}
	if err := validateStockSeries(c.HighValues, c.Categories, slideIndex, "high"); err != nil {
		return err
	}
	if err := validateStockSeries(c.LowValues, c.Categories, slideIndex, "low"); err != nil {
		return err
	}
	if err := validateStockSeries(c.CloseValues, c.Categories, slideIndex, "close"); err != nil {
		return err
	}
	return nil
}

func (c StockHLCChart) GetCategories() []string {
	return c.Categories
}

func (c StockHLCChart) GetValues() []float64 {
	return c.HighValues
}

func validateStockSeries(values []float64, categories []string, slideIndex int, label string) error {
	if len(values) != len(categories) {
		return fmt.Errorf("slide %d stock %s series length mismatch", slideIndex, label)
	}
	for i := range values {
		if math.IsNaN(values[i]) || math.IsInf(values[i], 0) {
			return fmt.Errorf("slide %d stock %s value %d must be finite", slideIndex, label, i+1)
		}
	}
	return nil
}

func validateStockCore(
	slideIndex int,
	title string,
	categories []string,
	x styling.Length,
	y styling.Length,
	cx styling.Length,
	cy styling.Length,

	valueFormat string,
	legendPosition string,
	valueAxisCrossBetween string,
	dataLabels DataLabelSettings,
	minValue *float64,
	maxValue *float64,
) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("slide %d stock chart title cannot be empty", slideIndex)
	}
	if x < 0 || y < 0 || cx <= 0 || cy <= 0 {
		return fmt.Errorf("slide %d stock chart geometry is invalid", slideIndex)
	}
	if len(categories) == 0 {
		return fmt.Errorf("slide %d stock chart must define categories", slideIndex)
	}
	if !IsLegendPosition(legendPosition) {
		return fmt.Errorf("slide %d stock chart legend position must be one of r,l,t,b", slideIndex)
	}
	if !IsDataLabelPosition(dataLabels.Position) {
		return fmt.Errorf(
			"slide %d stock chart data-label position must be ctr,inEnd,inBase,outEnd,bestFit,l,r,t,or b",
			slideIndex,
		)
	}
	if strings.TrimSpace(valueFormat) == "" {
		return fmt.Errorf("slide %d stock chart value format cannot be empty", slideIndex)
	}
	if !IsValueAxisCrossBetween(valueAxisCrossBetween) {
		return fmt.Errorf("slide %d stock chart value-axis crossBetween must be between or midCat", slideIndex)
	}
	if err := validateValueRange(minValue, maxValue, slideIndex); err != nil {
		return err
	}
	return nil
}
