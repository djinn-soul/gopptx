package charts

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

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
	return StockOHLCChart{StockHLCChart: base, OpenValues: openVals}
}

func (c StockOHLCChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := c.StockHLCChart.ToChartSpec()
	spec.Kind = pptxxml.ChartKindStockOHLC
	spec.OpenValues = CopyFloat64Slice(c.OpenValues)
	return spec
}

func (c StockOHLCChart) WithAltText(text string) StockOHLCChart {
	c.AltText = text
	return c
}

func (c StockOHLCChart) WithDecorative(enabled bool) StockOHLCChart {
	c.IsDecorative = enabled
	return c
}

func (c StockOHLCChart) Validate(slideIndex int) error {
	if err := c.StockHLCChart.Validate(slideIndex); err != nil {
		return err
	}
	return validateStockSeries(c.OpenValues, c.Categories, slideIndex, "open")
}

func (c StockOHLCChart) GetCategories() []string { return c.Categories }
func (c StockOHLCChart) GetValues() []float64    { return c.HighValues }

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
	DataLabels            DataLabelSettings
	ShowMajorGridlines    bool
	CategoryAxisTitle     string
	ValueAxisTitle        string
	ValueFormat           string
	ValueAxisCrossBetween string
	MinValue              *float64
	MaxValue              *float64

	AltText      string
	IsDecorative bool
}

func NewComboChart(categories []string, barSeries []Series, lineSeries []Series) ComboChart {
	return ComboChart{
		Title:                 defaultChartTitle,
		Categories:            append([]string(nil), categories...),
		BarSeries:             CopySeriesList(barSeries),
		LineSeries:            CopySeriesList(lineSeries),
		X:                     styling.Emu(defaultChartX),
		Y:                     styling.Emu(defaultChartY),
		CX:                    styling.Emu(defaultChartCX),
		CY:                    styling.Emu(defaultChartCY),
		ShowLegend:            true,
		LegendPosition:        LegendPositionRight,
		ShowDataLabels:        false,
		ShowMajorGridlines:    true,
		ValueFormat:           defaultChartValueFormat,
		ValueAxisCrossBetween: ValueAxisCrossBetweenBetween,
	}
}

func (c ComboChart) WithAltText(text string) ComboChart {
	c.AltText = text
	return c
}

func (c ComboChart) WithDecorative(enabled bool) ComboChart {
	c.IsDecorative = enabled
	return c
}

func (c ComboChart) ToChartSpec() *pptxxml.ChartSpec {
	spec := &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindCombo,
		Title:                 c.Title,
		TitleOverlay:          c.TitleOverlay,
		Categories:            CopyStringSlice(c.Categories),
		BarSeries:             ToXMLSeries(c.BarSeries),
		LineSeries:            ToXMLSeries(c.LineSeries),
		X:                     c.X.Emu(),
		Y:                     c.Y.Emu(),
		CX:                    c.CX.Emu(),
		CY:                    c.CY.Emu(),
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

func (c ComboChart) Validate(slideIndex int) error {
	if !c.IsDecorative && len(c.AltText) > common.MaxAltTextLength {
		return fmt.Errorf("slide %d combo chart alt text exceeds %d characters", slideIndex, common.MaxAltTextLength)
	}
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
	if !IsDataLabelPosition(c.DataLabels.Position) {
		return fmt.Errorf(
			"slide %d combo chart data-label position must be ctr,inEnd,inBase,outEnd,bestFit,l,r,t,or b",
			slideIndex,
		)
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

func (c ComboChart) GetCategories() []string { return c.Categories }
func (c ComboChart) GetValues() []float64    { return make([]float64, len(c.Categories)) }
