package charts

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

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
			slideIndex, chartType,
		)
	}
	if strings.TrimSpace(valueFormat) == "" {
		return fmt.Errorf("slide %d %s chart value format cannot be empty", slideIndex, chartType)
	}
	if !IsAxisTickLabelPosition(categoryTickLabelPosition) {
		return fmt.Errorf(
			"slide %d %s chart category-axis tick label position must be nextTo, low, high, or none",
			slideIndex, chartType,
		)
	}
	if !IsAxisTickLabelPosition(valueTickLabelPosition) {
		return fmt.Errorf(
			"slide %d %s chart value-axis tick label position must be nextTo, low, high, or none",
			slideIndex, chartType,
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

func validateChartCore(
	slideIndex int,
	title string,
	categories []string,
	values []float64,
	x, y, cx, cy styling.Length,
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
