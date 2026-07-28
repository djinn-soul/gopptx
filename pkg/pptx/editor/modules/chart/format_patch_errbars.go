package chart

import (
	"fmt"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

var reErrBarsBlock = regexp.MustCompile(`(?s)<c:errBars>.*?</c:errBars>`)

const (
	errBarDirectionX = "x"
	errBarDirectionY = "y"

	errValTypeCustom = "cust"
	errValTypeFixed  = "fixedVal"
	errValTypeStdDev = "stdDev"
	errBarTypeBoth   = "both"
)

// errBarTypes are the ST_ErrBarType values; errValTypes are ST_ErrValType.
//
//nolint:gochecknoglobals // Fixed enumerations, as elsewhere in this package.
var (
	errBarTypes = map[string]bool{errBarTypeBoth: true, "minus": true, xmlValuePlus: true}
	errValTypes = map[string]bool{
		errValTypeCustom: true, errValTypeFixed: true,
		"percentage": true, errValTypeStdDev: true, "stdErr": true,
	}
	errBarDirections = map[string]bool{errBarDirectionX: true, errBarDirectionY: true}
)

func validateChartErrorBars(errorBars []common.ChartErrorBars, clearSeries []int) error {
	for _, bars := range errorBars {
		if err := validateChartErrorBar(bars); err != nil {
			return err
		}
	}
	for _, index := range clearSeries {
		if index < 0 {
			return fmt.Errorf("clear_error_bar_series index must not be negative, got %d", index)
		}
	}
	return nil
}

func validateChartErrorBar(bars common.ChartErrorBars) error {
	if bars.SeriesIndex < 0 {
		return fmt.Errorf("error bar series_index must not be negative, got %d", bars.SeriesIndex)
	}
	if !errBarTypes[strings.TrimSpace(bars.BarType)] {
		return fmt.Errorf("error bar type must be one of both,minus,plus, got %q", bars.BarType)
	}
	valueType := strings.TrimSpace(bars.ValueType)
	if !errValTypes[valueType] {
		return fmt.Errorf(
			"error bar value_type must be one of cust,fixedVal,percentage,stdDev,stdErr, got %q",
			bars.ValueType,
		)
	}
	if bars.Direction != nil && !errBarDirections[strings.TrimSpace(*bars.Direction)] {
		return fmt.Errorf("error bar direction must be one of x,y, got %q", *bars.Direction)
	}
	if err := validateErrorBarValues(valueType, bars); err != nil {
		return err
	}
	if bars.LineColor != nil {
		if _, err := editorshape.NormalizeHexColor(*bars.LineColor); err != nil {
			return fmt.Errorf("error bar line_color: %w", err)
		}
	}
	return nil
}

// validateErrorBarValues enforces the c:val / c:plus-c:minus split: a custom
// error bar carries formula references, every other type carries a scalar.
func validateErrorBarValues(valueType string, bars common.ChartErrorBars) error {
	if valueType == errValTypeCustom {
		if bars.PlusReference == nil && bars.MinusReference == nil {
			return fmt.Errorf(
				"error bar value_type %q requires plus_reference and/or minus_reference", errValTypeCustom,
			)
		}
		if bars.Value != nil {
			return fmt.Errorf("error bar value is not valid with value_type %q", errValTypeCustom)
		}
		return nil
	}
	if bars.PlusReference != nil || bars.MinusReference != nil {
		return fmt.Errorf(
			"error bar plus_reference/minus_reference require value_type %q, got %q",
			errValTypeCustom, valueType,
		)
	}
	if bars.Value != nil && (!isFinite(*bars.Value) || *bars.Value < 0) {
		return fmt.Errorf("error bar value must be a finite non-negative number, got %v", *bars.Value)
	}
	return nil
}

// patchChartErrorBars replaces every c:errBars on each addressed series with the
// requested set, leaving series absent from the request untouched. A series
// listed in clearSeries ends up with no error bars at all.
func patchChartErrorBars(xml string, errorBars []common.ChartErrorBars, clearSeries []int) string {
	if len(errorBars) == 0 && len(clearSeries) == 0 {
		return xml
	}
	bySeries := map[int][]common.ChartErrorBars{}
	for _, index := range clearSeries {
		bySeries[index] = nil
	}
	for _, bars := range errorBars {
		bySeries[bars.SeriesIndex] = append(bySeries[bars.SeriesIndex], bars)
	}

	seriesIndex := -1
	return reSerBlocks.ReplaceAllStringFunc(xml, func(ser string) string {
		seriesIndex++
		requested, ok := bySeries[seriesIndex]
		if !ok {
			return ser
		}
		// Rewriting from scratch keeps a repeated patch idempotent instead of
		// stacking duplicate error bars on the series.
		ser = reErrBarsBlock.ReplaceAllLiteralString(ser, "")
		var nodes strings.Builder
		for _, bars := range requested {
			nodes.WriteString(buildErrorBarsBlock(bars))
		}
		return insertErrorBarsBlock(ser, nodes.String())
	})
}

// buildErrorBarsBlock renders one c:errBars. CT_ErrBars orders its children
// errDir, errBarType, errValType, noEndCap, plus, minus, val, spPr.
func buildErrorBarsBlock(bars common.ChartErrorBars) string {
	var b strings.Builder
	b.WriteString("<c:errBars>")
	if bars.Direction != nil {
		b.WriteString(`<c:errDir val="` + strings.TrimSpace(*bars.Direction) + `"/>`)
	}
	b.WriteString(`<c:errBarType val="` + strings.TrimSpace(bars.BarType) + `"/>`)
	b.WriteString(`<c:errValType val="` + strings.TrimSpace(bars.ValueType) + `"/>`)
	// c:noEndCap is always written, even when it is the schema default. It is
	// optional per CT_ErrBars, but PowerPoint's renderer only applies the
	// c:spPr line style to every point when the element is present; omit it and
	// the colour lands on the last point alone.
	noEndCap := false
	if bars.NoEndCap != nil {
		noEndCap = *bars.NoEndCap
	}
	b.WriteString(`<c:noEndCap val="` + boolToOneZero(noEndCap) + `"/>`)
	if bars.PlusReference != nil {
		b.WriteString(buildErrorBarReference("plus", *bars.PlusReference))
	}
	if bars.MinusReference != nil {
		b.WriteString(buildErrorBarReference("minus", *bars.MinusReference))
	}
	if bars.Value != nil {
		b.WriteString(`<c:val val="` + formatAxisNumber(*bars.Value) + `"/>`)
	}
	b.WriteString(buildErrorBarShapeProperties(bars))
	b.WriteString("</c:errBars>")
	return b.String()
}

// buildErrorBarReference renders a c:plus or c:minus formula reference, which
// CT_NumDataSource expresses as a c:numRef around the sheet formula.
func buildErrorBarReference(tag string, formula string) string {
	return `<c:` + tag + `><c:numRef><c:f>` + xmlEscape(strings.TrimSpace(formula)) +
		`</c:f></c:numRef></c:` + tag + `>`
}

func buildErrorBarShapeProperties(bars common.ChartErrorBars) string {
	if bars.LineColor == nil {
		return ""
	}
	color, err := editorshape.NormalizeHexColor(*bars.LineColor)
	if err != nil {
		return ""
	}
	return `<c:spPr><a:ln><a:solidFill><a:srgbClr val="` + color + `"/></a:solidFill></a:ln></c:spPr>`
}

// insertErrorBarsBlock splices the error bars into their schema slot: CT_Ser
// puts c:errBars after c:trendline and before the data references.
func insertErrorBarsBlock(ser string, nodes string) string {
	if nodes == "" {
		return ser
	}
	for _, anchor := range []string{
		seriesCategoryTag, seriesXValuesTag, seriesValuesTag, seriesSmoothTag, chartExtensionListTagPrefix,
	} {
		if index := strings.Index(ser, anchor); index >= 0 {
			return ser[:index] + nodes + ser[index:]
		}
	}
	if index := strings.LastIndex(ser, seriesCloseTag); index >= 0 {
		return ser[:index] + nodes + ser[index:]
	}
	return ser
}
