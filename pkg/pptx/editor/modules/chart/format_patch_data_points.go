package chart

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

var (
	reDataPointBlock = regexp.MustCompile(`(?s)<c:dPt>.*?</c:dPt>`)
	reDataPointIdx   = regexp.MustCompile(`<c:idx val="([^"]+)"`)
	reDataPointInv   = regexp.MustCompile(`<c:invertIfNegative val="([^"]+)"`)
	reDataPointBub3D = regexp.MustCompile(`<c:bubble3D val="([^"]+)"`)
	reDataPointExpl  = regexp.MustCompile(`<c:explosion val="([^"]+)"`)
	reDataPointFill  = regexp.MustCompile(`(?s)<a:solidFill>\s*<a:srgbClr val="([^"]+)"`)
	reDataPointLine  = regexp.MustCompile(`(?s)<a:ln\b[^>]*>.*?<a:srgbClr val="([^"]+)"`)
	reDataPointWidth = regexp.MustCompile(`<a:ln\b[^>]*\bw="([^"]+)"`)

	// Effects a point may carry that this API does not model, kept across a
	// rebuild so recolouring a point does not drop its shadow.
	reDataPointEffects = regexp.MustCompile(
		`(?s)<a:effectLst>.*?</a:effectLst>|<a:effectLst\s*/>` +
			`|<a:scene3d>.*?</a:scene3d>|<a:sp3d\b[^>]*/>|<a:sp3d\b[^>]*>.*?</a:sp3d>`,
	)

	reDataPointMarker     = regexp.MustCompile(`(?s)<c:marker>.*?</c:marker>`)
	reDataPointMarkerSym  = regexp.MustCompile(`<c:symbol val="([^"]+)"`)
	reDataPointMarkerSize = regexp.MustCompile(`<c:size val="([^"]+)"`)
)

const (
	explosionMax       = 400
	markerSizeMin      = 2
	markerSizeMax      = 72
	markerSymbolCircle = "circle"
)

// markerSymbols is the CT_MarkerStyle value set.
//
//nolint:gochecknoglobals // Fixed enumeration, as elsewhere in this package.
var markerSymbols = map[string]bool{
	markerSymbolCircle: true, lineDashPresetDash: true, "diamond": true, "dot": true, "none": true,
	xmlValuePlus: true, "square": true, "star": true, "triangle": true, "x": true,
	"auto": true,
}

// patchChartDataPoints merges the requested per-point formatting into each
// addressed series, leaving points and series it does not mention untouched.
func patchChartDataPoints(xml string, points []common.ChartDataPoint, clearSeries []int) string {
	if len(points) == 0 && len(clearSeries) == 0 {
		return xml
	}
	bySeries := map[int][]common.ChartDataPoint{}
	cleared := map[int]bool{}
	for _, index := range clearSeries {
		cleared[index] = true
		bySeries[index] = nil
	}
	for _, point := range points {
		bySeries[point.SeriesIndex] = append(bySeries[point.SeriesIndex], point)
	}

	seriesIndex := -1
	return reSerBlocks.ReplaceAllStringFunc(xml, func(ser string) string {
		seriesIndex++
		requested, ok := bySeries[seriesIndex]
		if !ok {
			return ser
		}
		return applySeriesDataPoints(ser, requested, cleared[seriesIndex])
	})
}

// applySeriesDataPoints rewrites the whole c:dPt run for one series so a
// repeated patch is idempotent rather than duplicating points.
func applySeriesDataPoints(ser string, requested []common.ChartDataPoint, reset bool) string {
	merged := map[int]common.ChartDataPoint{}
	effects := map[int]string{}
	if !reset {
		for _, existing := range parseSeriesDataPoints(ser) {
			merged[existing.PointIndex] = existing
		}
		// The rebuilt c:spPr keeps the effects this API does not model, so
		// recolouring a point does not drop its shadow (upstream #450).
		effects = parseSeriesDataPointEffects(ser)
	}
	for _, point := range requested {
		merged[point.PointIndex] = mergeDataPoint(merged[point.PointIndex], point)
	}

	ser = reDataPointBlock.ReplaceAllLiteralString(ser, "")
	indexes := make([]int, 0, len(merged))
	for index := range merged {
		indexes = append(indexes, index)
	}
	// CT_Ser requires the c:dPt run in ascending c:idx order.
	sort.Ints(indexes)

	var nodes strings.Builder
	for _, index := range indexes {
		nodes.WriteString(buildDataPointBlock(merged[index], effects[index]))
	}
	return insertDataPointBlock(ser, nodes.String())
}

// mergeDataPoint layers the requested fields over what the series already had,
// so patching one property does not drop the rest of a point's formatting.
func mergeDataPoint(base common.ChartDataPoint, update common.ChartDataPoint) common.ChartDataPoint {
	merged := base
	merged.SeriesIndex = update.SeriesIndex
	merged.PointIndex = update.PointIndex
	if update.FillColor != nil {
		merged.FillColor = update.FillColor
	}
	if update.LineColor != nil {
		merged.LineColor = update.LineColor
	}
	if update.LineWidthEMU != nil {
		merged.LineWidthEMU = update.LineWidthEMU
	}
	if update.InvertIfNegative != nil {
		merged.InvertIfNegative = update.InvertIfNegative
	}
	if update.Bubble3D != nil {
		merged.Bubble3D = update.Bubble3D
	}
	if update.Explosion != nil {
		merged.Explosion = update.Explosion
	}
	if update.MarkerFillColor != nil {
		merged.MarkerFillColor = update.MarkerFillColor
	}
	if update.MarkerLineColor != nil {
		merged.MarkerLineColor = update.MarkerLineColor
	}
	if update.MarkerSymbol != nil {
		merged.MarkerSymbol = update.MarkerSymbol
	}
	if update.MarkerSize != nil {
		merged.MarkerSize = update.MarkerSize
	}
	return merged
}

// buildDataPointBlock renders one c:dPt. CT_DPt orders its children idx,
// invertIfNegative, marker, bubble3D, explosion, spPr, pictureOptions, extLst.
func buildDataPointBlock(point common.ChartDataPoint, effects string) string {
	var b strings.Builder
	b.WriteString(`<c:dPt><c:idx val="` + strconv.Itoa(point.PointIndex) + `"/>`)
	if point.InvertIfNegative != nil {
		b.WriteString(`<c:invertIfNegative val="` + boolToOneZero(*point.InvertIfNegative) + `"/>`)
	}
	b.WriteString(buildDataPointMarker(point))
	if point.Bubble3D != nil {
		b.WriteString(`<c:bubble3D val="` + boolToOneZero(*point.Bubble3D) + `"/>`)
	}
	if point.Explosion != nil {
		b.WriteString(`<c:explosion val="` + strconv.Itoa(*point.Explosion) + `"/>`)
	}
	b.WriteString(buildDataPointShapeProperties(point, effects))
	b.WriteString("</c:dPt>")
	return b.String()
}

// buildDataPointMarker renders the point's c:marker. CT_Marker orders symbol,
// size, then spPr, and on a scatter or line series this is what carries the
// point's own colour: the c:spPr formats the connecting segment instead.
func buildDataPointMarker(point common.ChartDataPoint) string {
	if point.MarkerSymbol == nil && point.MarkerSize == nil &&
		point.MarkerFillColor == nil && point.MarkerLineColor == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString("<c:marker>")
	if point.MarkerSymbol != nil {
		b.WriteString(`<c:symbol val="` + *point.MarkerSymbol + `"/>`)
	}
	if point.MarkerSize != nil {
		b.WriteString(`<c:size val="` + strconv.Itoa(*point.MarkerSize) + `"/>`)
	}
	b.WriteString(buildMarkerShapeProperties(point))
	b.WriteString("</c:marker>")
	return b.String()
}

// buildMarkerShapeProperties renders the marker's own c:spPr, fill before line.
func buildMarkerShapeProperties(point common.ChartDataPoint) string {
	if point.MarkerFillColor == nil && point.MarkerLineColor == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString("<c:spPr>")
	if color := normalizedColor(point.MarkerFillColor); color != "" {
		b.WriteString(`<a:solidFill><a:srgbClr val="` + color + `"/></a:solidFill>`)
	}
	if color := normalizedColor(point.MarkerLineColor); color != "" {
		b.WriteString(`<a:ln><a:solidFill><a:srgbClr val="` + color + `"/></a:solidFill></a:ln>`)
	}
	b.WriteString("</c:spPr>")
	return b.String()
}

// normalizedColor returns the hex form of a colour, or "" when it is unset or
// unparsable. Validation has already rejected bad values by this point.
func normalizedColor(value *string) string {
	if value == nil {
		return ""
	}
	color, err := editorshape.NormalizeHexColor(*value)
	if err != nil {
		return ""
	}
	return color
}

// buildDataPointShapeProperties renders the c:spPr. CT_ShapeProperties orders
// the fill before the line and the effects after both; effects carries whatever
// the point already had there, which this API does not otherwise model.
func buildDataPointShapeProperties(point common.ChartDataPoint, effects string) string {
	if point.FillColor == nil && point.LineColor == nil && point.LineWidthEMU == nil {
		if effects == "" {
			return ""
		}
		return "<c:spPr>" + effects + "</c:spPr>"
	}
	var b strings.Builder
	b.WriteString("<c:spPr>")
	if color := normalizedColor(point.FillColor); color != "" {
		b.WriteString(`<a:solidFill><a:srgbClr val="` + color + `"/></a:solidFill>`)
	}
	if point.LineColor != nil || point.LineWidthEMU != nil {
		attrs := ""
		if point.LineWidthEMU != nil {
			attrs = ` w="` + strconv.Itoa(*point.LineWidthEMU) + `"`
		}
		b.WriteString("<a:ln" + attrs + ">")
		if color := normalizedColor(point.LineColor); color != "" {
			b.WriteString(`<a:solidFill><a:srgbClr val="` + color + `"/></a:solidFill>`)
		}
		b.WriteString("</a:ln>")
	}
	b.WriteString(effects)
	b.WriteString("</c:spPr>")
	return b.String()
}

// insertDataPointBlock splices the c:dPt run into its schema slot: CT_Ser puts
// it after c:spPr / c:invertIfNegative and before c:dLbls.
func insertDataPointBlock(ser string, nodes string) string {
	if nodes == "" {
		return ser
	}
	for _, anchor := range []string{
		seriesDataLabelsTag, seriesTrendlineTag, seriesErrorBarsTag, seriesCategoryTag, seriesXValuesTag, seriesValuesTag,
		seriesSmoothTag, chartExtensionListTagPrefix,
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
