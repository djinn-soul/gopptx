package pptx

import (
	"fmt"
	"math"
	"strings"
)

// StockHLCChart is a stock chart with High/Low/Close series.
type StockHLCChart struct {
	Title                 string
	TitleOverlay          bool
	Categories            []string
	HighValues            []float64
	LowValues             []float64
	CloseValues           []float64
	X                     int64
	Y                     int64
	CX                    int64
	CY                    int64
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

func NewStockHLCChart(categories []string, high []float64, low []float64, close []float64) StockHLCChart {
	cats := append([]string(nil), categories...)
	highVals := append([]float64(nil), high...)
	lowVals := append([]float64(nil), low...)
	closeVals := append([]float64(nil), close...)
	return StockHLCChart{
		Title:                 "Chart",
		Categories:            cats,
		HighValues:            highVals,
		LowValues:             lowVals,
		CloseValues:           closeVals,
		X:                     685800,
		Y:                     1800000,
		CX:                    7772400,
		CY:                    4114800,
		ShowLegend:            true,
		LegendPosition:        LegendPositionRight,
		ShowDataLabels:        false,
		ShowMajorGridlines:    true,
		ValueFormat:           "General",
		ValueAxisCrossBetween: ValueAxisCrossBetweenBetween,
	}
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
	close []float64,
) StockOHLCChart {
	base := NewStockHLCChart(categories, high, low, close)
	openVals := append([]float64(nil), open...)
	return StockOHLCChart{
		StockHLCChart: base,
		OpenValues:    openVals,
	}
}

// ComboChart mixes bar and line series on one category axis.
type ComboChart struct {
	Title                 string
	TitleOverlay          bool
	Categories            []string
	BarSeries             []Series
	LineSeries            []Series
	X                     int64
	Y                     int64
	CX                    int64
	CY                    int64
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
		Title:                 "Chart",
		Categories:            append([]string(nil), categories...),
		BarSeries:             copySeriesList(barSeries),
		LineSeries:            copySeriesList(lineSeries),
		X:                     685800,
		Y:                     1800000,
		CX:                    7772400,
		CY:                    4114800,
		ShowLegend:            true,
		LegendPosition:        LegendPositionRight,
		ShowDataLabels:        false,
		ShowMajorGridlines:    true,
		ValueFormat:           "General",
		ValueAxisCrossBetween: ValueAxisCrossBetweenBetween,
	}
}

func validateStockHLCChart(chart StockHLCChart, slideIndex int) error {
	if err := validateStockCore(
		slideIndex,
		chart.Title,
		chart.Categories,
		chart.X,
		chart.Y,
		chart.CX,
		chart.CY,
		chart.ValueFormat,
		chart.LegendPosition,
		chart.ValueAxisCrossBetween,
		chart.MinValue,
		chart.MaxValue,
	); err != nil {
		return err
	}
	if err := validateStockSeries(chart.HighValues, chart.Categories, slideIndex, "high"); err != nil {
		return err
	}
	if err := validateStockSeries(chart.LowValues, chart.Categories, slideIndex, "low"); err != nil {
		return err
	}
	if err := validateStockSeries(chart.CloseValues, chart.Categories, slideIndex, "close"); err != nil {
		return err
	}
	return nil
}

func validateStockOHLCChart(chart StockOHLCChart, slideIndex int) error {
	if err := validateStockHLCChart(chart.StockHLCChart, slideIndex); err != nil {
		return err
	}
	return validateStockSeries(chart.OpenValues, chart.Categories, slideIndex, "open")
}

func validateComboChart(chart ComboChart, slideIndex int) error {
	if err := validateChartCore(
		slideIndex,
		chart.Title,
		chart.Categories,
		chart.CategoriesToValues(),
		chart.X,
		chart.Y,
		chart.CX,
		chart.CY,
	); err != nil {
		return err
	}
	if !isLegendPosition(chart.LegendPosition) {
		return fmt.Errorf("slide %d combo chart legend position must be one of r,l,t,b", slideIndex)
	}
	if strings.TrimSpace(chart.ValueFormat) == "" {
		return fmt.Errorf("slide %d combo chart value format cannot be empty", slideIndex)
	}
	if !isValueAxisCrossBetween(chart.ValueAxisCrossBetween) {
		return fmt.Errorf("slide %d combo chart value-axis crossBetween must be between or midCat", slideIndex)
	}
	if err := validateValueRange(chart.MinValue, chart.MaxValue, slideIndex); err != nil {
		return err
	}
	if err := validateSeriesList(chart.BarSeries, len(chart.Categories), slideIndex, "combo bar"); err != nil {
		return err
	}
	if err := validateSeriesList(chart.LineSeries, len(chart.Categories), slideIndex, "combo line"); err != nil {
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
	x int64,
	y int64,
	cx int64,
	cy int64,
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
	if !isLegendPosition(legendPosition) {
		return fmt.Errorf("slide %d stock chart legend position must be one of r,l,t,b", slideIndex)
	}
	if strings.TrimSpace(valueFormat) == "" {
		return fmt.Errorf("slide %d stock chart value format cannot be empty", slideIndex)
	}
	if !isValueAxisCrossBetween(valueAxisCrossBetween) {
		return fmt.Errorf("slide %d stock chart value-axis crossBetween must be between or midCat", slideIndex)
	}
	if err := validateValueRange(minValue, maxValue, slideIndex); err != nil {
		return err
	}
	return nil
}
