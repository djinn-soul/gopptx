package chart

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

var reTrendlineBlock = regexp.MustCompile(`(?s)<c:trendline>.*?</c:trendline>`)

const (
	trendlineTypeLinear    = "linear"
	trendlineTypePoly      = "poly"
	trendlineTypeMovingAvg = "movingAvg"

	// trendlinePolyOrderMin/Max and trendlineMovingAvgPeriodMin are the ranges
	// PowerPoint accepts; outside them it repairs the file on open.
	trendlinePolyOrderMin       = 2
	trendlinePolyOrderMax       = 6
	trendlineMovingAvgPeriodMin = 2
)

// trendlineTypes are the ST_TrendlineType enumeration values.
//
//nolint:gochecknoglobals // Fixed enumeration, as elsewhere in this package.
var trendlineTypes = map[string]bool{
	trendlineTypeLinear: true, trendlineTypePoly: true, "exp": true,
	"log": true, trendlineTypeMovingAvg: true, "power": true,
}

func validateChartTrendlines(
	trendlines []common.ChartTrendline,
	appendTrendlines []common.ChartTrendline,
	clearSeries []int,
) error {
	for _, group := range [][]common.ChartTrendline{trendlines, appendTrendlines} {
		for _, trendline := range group {
			if err := validateChartTrendline(trendline); err != nil {
				return err
			}
		}
	}
	for _, index := range clearSeries {
		if index < 0 {
			return fmt.Errorf("clear_trendline_series index must not be negative, got %d", index)
		}
	}
	return nil
}

func validateChartTrendline(trendline common.ChartTrendline) error {
	kind := strings.TrimSpace(trendline.Type)
	if !trendlineTypes[kind] {
		return fmt.Errorf(
			"trendline type must be one of linear,poly,exp,log,movingAvg,power, got %q", trendline.Type,
		)
	}
	if trendline.SeriesIndex < 0 {
		return fmt.Errorf("trendline series_index must not be negative, got %d", trendline.SeriesIndex)
	}
	if err := validateTrendlineOrder(kind, trendline.Order); err != nil {
		return err
	}
	if err := validateTrendlinePeriod(kind, trendline.Period); err != nil {
		return err
	}
	for _, field := range []struct {
		name  string
		value *float64
	}{
		{"forward", trendline.Forward},
		{"backward", trendline.Backward},
		{"intercept", trendline.Intercept},
	} {
		if field.value != nil && !isFinite(*field.value) {
			return fmt.Errorf("trendline %s must be finite, got %v", field.name, *field.value)
		}
	}
	if trendline.LineColor != nil {
		if _, err := editorshape.NormalizeHexColor(*trendline.LineColor); err != nil {
			return fmt.Errorf("trendline line_color: %w", err)
		}
	}
	if trendline.LineDash != nil {
		if _, err := editorshape.NormalizeLineDashStyle(*trendline.LineDash); err != nil {
			return fmt.Errorf("trendline line_dash: %w", err)
		}
	}
	if trendline.LineWidthEMU != nil && *trendline.LineWidthEMU < 0 {
		return fmt.Errorf("trendline line_width_emu must not be negative, got %d", *trendline.LineWidthEMU)
	}
	return nil
}

// validateTrendlineOrder rejects c:order outside a polynomial trendline.
// PowerPoint repairs a file that carries an order on any other type.
func validateTrendlineOrder(kind string, order *int) error {
	if order == nil {
		return nil
	}
	if kind != trendlineTypePoly {
		return fmt.Errorf("trendline order is only valid for the poly type, got %q", kind)
	}
	if *order < trendlinePolyOrderMin || *order > trendlinePolyOrderMax {
		return fmt.Errorf(
			"trendline order must be between %d and %d, got %d",
			trendlinePolyOrderMin, trendlinePolyOrderMax, *order,
		)
	}
	return nil
}

// validateTrendlinePeriod rejects c:period outside a moving-average trendline.
func validateTrendlinePeriod(kind string, period *int) error {
	if period == nil {
		return nil
	}
	if kind != trendlineTypeMovingAvg {
		return fmt.Errorf("trendline period is only valid for the movingAvg type, got %q", kind)
	}
	if *period < trendlineMovingAvgPeriodMin {
		return fmt.Errorf(
			"trendline period must be greater than or equal to %d, got %d",
			trendlineMovingAvgPeriodMin, *period,
		)
	}
	return nil
}

// patchChartTrendlines replaces or appends trendlines on addressed series.
// Appended entries leave the original XML intact, including children that the
// object model does not expose. A series listed in clearSeries is cleared before
// any requested replacement or append is applied.
func patchChartTrendlines(
	xml string,
	trendlines []common.ChartTrendline,
	appendTrendlines []common.ChartTrendline,
	clearSeries []int,
) string {
	if len(trendlines) == 0 && len(appendTrendlines) == 0 && len(clearSeries) == 0 {
		return xml
	}
	replacements := map[int][]common.ChartTrendline{}
	for _, index := range clearSeries {
		replacements[index] = nil
	}
	for _, trendline := range trendlines {
		replacements[trendline.SeriesIndex] = append(
			replacements[trendline.SeriesIndex], trendline,
		)
	}
	appends := map[int][]common.ChartTrendline{}
	for _, trendline := range appendTrendlines {
		appends[trendline.SeriesIndex] = append(appends[trendline.SeriesIndex], trendline)
	}

	seriesIndex := -1
	return reSerBlocks.ReplaceAllStringFunc(xml, func(ser string) string {
		seriesIndex++
		if requested, replace := replacements[seriesIndex]; replace {
			// Rewriting replacement requests from scratch keeps them idempotent.
			ser = reTrendlineBlock.ReplaceAllLiteralString(ser, "")
			ser = insertTrendlineBlock(ser, buildTrendlineBlocks(requested))
		}
		return insertTrendlineBlock(ser, buildTrendlineBlocks(appends[seriesIndex]))
	})
}

func buildTrendlineBlocks(trendlines []common.ChartTrendline) string {
	var nodes strings.Builder
	for _, trendline := range trendlines {
		nodes.WriteString(buildTrendlineBlock(trendline))
	}
	return nodes.String()
}

// buildTrendlineBlock renders one c:trendline. CT_Trendline orders its children
// name, spPr, trendlineType, order, period, forward, backward, intercept,
// dispRSqr, dispEq, trendlineLbl, extLst.
func buildTrendlineBlock(trendline common.ChartTrendline) string {
	var b strings.Builder
	b.WriteString("<c:trendline>")
	if trendline.Name != nil {
		b.WriteString("<c:name>" + xmlEscape(*trendline.Name) + "</c:name>")
	}
	b.WriteString(buildTrendlineShapeProperties(trendline))
	b.WriteString(`<c:trendlineType val="` + strings.TrimSpace(trendline.Type) + `"/>`)
	if trendline.Order != nil {
		b.WriteString(`<c:order val="` + strconv.Itoa(*trendline.Order) + `"/>`)
	}
	if trendline.Period != nil {
		b.WriteString(`<c:period val="` + strconv.Itoa(*trendline.Period) + `"/>`)
	}
	for _, node := range []struct {
		tag   string
		value *float64
	}{
		{"forward", trendline.Forward},
		{"backward", trendline.Backward},
		{"intercept", trendline.Intercept},
	} {
		if node.value != nil {
			b.WriteString(`<c:` + node.tag + ` val="` + formatAxisNumber(*node.value) + `"/>`)
		}
	}
	if trendline.DisplayRSquared != nil {
		b.WriteString(`<c:dispRSqr val="` + boolToOneZero(*trendline.DisplayRSquared) + `"/>`)
	}
	if trendline.DisplayEquation != nil {
		b.WriteString(`<c:dispEq val="` + boolToOneZero(*trendline.DisplayEquation) + `"/>`)
	}
	b.WriteString("</c:trendline>")
	return b.String()
}

// buildTrendlineShapeProperties renders the c:spPr line style. CT_LineProperties
// puts the width on a:ln, then the fill, then the dash.
func buildTrendlineShapeProperties(trendline common.ChartTrendline) string {
	if trendline.LineColor == nil && trendline.LineWidthEMU == nil && trendline.LineDash == nil {
		return ""
	}
	attrs := ""
	if trendline.LineWidthEMU != nil {
		attrs = ` w="` + strconv.Itoa(*trendline.LineWidthEMU) + `"`
	}
	var line strings.Builder
	line.WriteString("<a:ln" + attrs + ">")
	if trendline.LineColor != nil {
		color, err := editorshape.NormalizeHexColor(*trendline.LineColor)
		if err == nil {
			line.WriteString(`<a:solidFill><a:srgbClr val="` + color + `"/></a:solidFill>`)
		}
	}
	if trendline.LineDash != nil {
		dash, err := editorshape.NormalizeLineDashStyle(*trendline.LineDash)
		if err == nil {
			line.WriteString(`<a:prstDash val="` + dash + `"/>`)
		}
	}
	line.WriteString("</a:ln>")
	return "<c:spPr>" + line.String() + "</c:spPr>"
}

// insertTrendlineBlock splices the trendlines into their schema slot: CT_Ser
// puts c:trendline after c:dLbls and before c:errBars and the data references.
func insertTrendlineBlock(ser string, nodes string) string {
	if nodes == "" {
		return ser
	}
	for _, anchor := range []string{
		seriesErrorBarsTag, seriesCategoryTag, seriesXValuesTag, seriesValuesTag, seriesSmoothTag, chartExtensionListTagPrefix,
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
