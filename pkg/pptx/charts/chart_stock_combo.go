package charts

import (
	"fmt"
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
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
	ShowMajorGridlines    bool
	CategoryAxisTitle     string
	ValueAxisTitle        string
	ValueFormat           string
	ValueAxisCrossBetween string
	MinValue              *float64
	MaxValue              *float64
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
		X:           styling.Emu(685800),
		Y:           styling.Emu(1800000),
		CX:          styling.Emu(7772400),
		CY:          styling.Emu(4114800),

		ShowLegend:            true,
		LegendPosition:        LegendPositionRight,
		ShowDataLabels:        false,
		ShowMajorGridlines:    true,
		ValueFormat:           "General",
		ValueAxisCrossBetween: ValueAxisCrossBetweenBetween,
	}
}

// ToChartSpec converts StockHLCChart to internal XML spec.
func (c StockHLCChart) ToChartSpec() *pptxxml.ChartSpec {
	return &pptxxml.ChartSpec{
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
	}
}

// Validate checks the stock chart for consistency.
func (c StockHLCChart) Validate(slideIndex int) error {
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

// StockOHLCChart is a stock chart with Open/High/Low/Close series.
type StockOHLCChart struct {
	StockHLCChart

	OpenValues []float64
}

func NewStockOHLCChart(
	categories []string,
	open []float64,
	high []float64,
	low []float64,
	closeValues []float64,
) StockOHLCChart {
	base := NewStockHLCChart(categories, high, low, closeValues)
	openVals := append([]float64(nil), open...)
	return StockOHLCChart{
		StockHLCChart: base,
		OpenValues:    openVals,
	}
}

// ToChartSpec converts StockOHLCChart to internal XML spec.
func (c StockOHLCChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.StockHLCChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindStockOHLC
	spec.OpenValues = CopyFloat64Slice(c.OpenValues)
	return spec
}

// Validate checks the stock chart for consistency.
func (c StockOHLCChart) Validate(slideIndex int) error {
	if err := c.StockHLCChart.Validate(slideIndex); err != nil {
		return err
	}
	return validateStockSeries(c.OpenValues, c.Categories, slideIndex, "open")
}

func (c StockOHLCChart) GetCategories() []string {
	return c.Categories
}

func (c StockOHLCChart) GetValues() []float64 {
	return c.HighValues
}

// ComboChart mixes bar and line series on one category axis.
type ComboChart struct {
	Title        string
	TitleOverlay bool
	Categories   []string
	BarSeries    []Series
	LineSeries   []Series
	X            styling.Length
	Y            styling.Length
	CX           styling.Length
	CY           styling.Length

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

func NewComboChart(categories []string, barSeries []Series, lineSeries []Series) ComboChart {
	return ComboChart{
		Title:      "Chart",
		Categories: append([]string(nil), categories...),
		BarSeries:  CopySeriesList(barSeries),
		LineSeries: CopySeriesList(lineSeries),
		X:          styling.Emu(685800),
		Y:          styling.Emu(1800000),
		CX:         styling.Emu(7772400),
		CY:         styling.Emu(4114800),

		ShowLegend:            true,
		LegendPosition:        LegendPositionRight,
		ShowDataLabels:        false,
		ShowMajorGridlines:    true,
		ValueFormat:           "General",
		ValueAxisCrossBetween: ValueAxisCrossBetweenBetween,
	}
}

// ToChartSpec converts ComboChart to internal XML spec.
func (c ComboChart) ToChartSpec() *pptxxml.ChartSpec {
	return &pptxxml.ChartSpec{
		Kind:         pptxxml.ChartKindCombo,
		Title:        c.Title,
		TitleOverlay: c.TitleOverlay,
		Categories:   CopyStringSlice(c.Categories),
		BarSeries:    ToXMLSeries(c.BarSeries),
		LineSeries:   ToXMLSeries(c.LineSeries),
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
	}
}

// Validate checks the combo chart for consistency.
func (c ComboChart) Validate(slideIndex int) error {
	if err := validateChartCore(
		slideIndex,
		c.Title,
		c.Categories,
		c.CategoriesToValues(),
		c.X,
		c.Y,
		c.CX,
		c.CY,

		false,
	); err != nil {
		return err
	}
	if !IsLegendPosition(c.LegendPosition) {
		return fmt.Errorf("slide %d combo chart legend position must be one of r,l,t,b", slideIndex)
	}
	if strings.TrimSpace(c.ValueFormat) == "" {
		return fmt.Errorf("slide %d combo chart value format cannot be empty", slideIndex)
	}
	if !IsValueAxisCrossBetween(c.ValueAxisCrossBetween) {
		return fmt.Errorf("slide %d combo chart value-axis crossBetween must be between or midCat", slideIndex)
	}
	if err := validateValueRange(c.MinValue, c.MaxValue, slideIndex); err != nil {
		return err
	}
	if err := ValidateSeriesList(c.BarSeries, len(c.Categories), slideIndex, "combo bar"); err != nil {
		return err
	}
	if err := ValidateSeriesList(c.LineSeries, len(c.Categories), slideIndex, "combo line"); err != nil {
		return err
	}
	return nil
}

func (c ComboChart) CategoriesToValues() []float64 {
	values := make([]float64, len(c.Categories))
	for i := range values {
		values[i] = 1
	}
	return values
}

func (c ComboChart) GetCategories() []string {
	return c.Categories
}

func (c ComboChart) GetValues() []float64 {
	// Return zero-value slice of correct length since combo has mixed series
	return make([]float64, len(c.Categories))
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
