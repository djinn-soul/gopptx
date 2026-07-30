package chart

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

// Series-level formatting, upstream #872: the c:spPr and c:marker of a whole
// c:ser. On a line, scatter or radar series the c:spPr styles the line joining
// the points, and each marker carries its own fill and outline, so recolouring
// the markers of a series needs the marker fields rather than the line ones.

var (
	reSeriesMarker   = regexp.MustCompile(`(?s)<c:marker>.*?</c:marker>|<c:marker\s*/>`)
	reSeriesSmooth   = regexp.MustCompile(`<c:smooth val="[^"]*"/>`)
	reShapeFillBlock = regexp.MustCompile(
		`(?s)<a:solidFill>.*?</a:solidFill>|<a:noFill\s*/>` +
			`|<a:gradFill\b.*?</a:gradFill>|<a:pattFill\b.*?</a:pattFill>|<a:blipFill\b.*?</a:blipFill>`,
	)
)

// seriesFormatBoundaryTags are the c:ser children that follow c:spPr and
// c:marker. The first of them bounds the search, so a c:spPr belonging to a
// c:dPt or a data label is never mistaken for the series' own.
//
//nolint:gochecknoglobals // Fixed tag set, as elsewhere in this package.
var seriesFormatBoundaryTags = []string{
	"<c:dPt>", seriesDataLabelsTag, seriesTrendlineTag, seriesErrorBarsTag,
	seriesCategoryTag, seriesXValuesTag, seriesValuesTag, seriesSmoothTag,
}

func validateChartSeriesFormats(formats []common.ChartSeriesFormat) error {
	for _, format := range formats {
		if format.SeriesIndex < 0 {
			return fmt.Errorf("series format series_index must not be negative, got %d", format.SeriesIndex)
		}
		if err := validateSeriesFormatColors(format); err != nil {
			return err
		}
		if err := validateSeriesFormatLine(format); err != nil {
			return err
		}
		if err := validateSeriesFormatMarker(format); err != nil {
			return err
		}
	}
	return nil
}

func validateSeriesFormatLine(format common.ChartSeriesFormat) error {
	if format.LineDash != nil {
		if _, err := editorshape.NormalizeLineDashStyle(*format.LineDash); err != nil {
			return fmt.Errorf("series format line_dash: %w", err)
		}
	}
	for _, pair := range []struct {
		value *int
		name  string
	}{
		{format.LineWidthEMU, "line_width_emu"},
		{format.MarkerLineWidthEMU, "marker_line_width_emu"},
	} {
		if pair.value != nil && *pair.value < 0 {
			return fmt.Errorf("series format %s must not be negative, got %d", pair.name, *pair.value)
		}
	}
	return nil
}

func validateSeriesFormatMarker(format common.ChartSeriesFormat) error {
	if format.MarkerSymbol != nil && !markerSymbols[*format.MarkerSymbol] {
		return fmt.Errorf("series format marker_symbol %q is not a CT_MarkerStyle value", *format.MarkerSymbol)
	}
	if format.MarkerSize != nil && (*format.MarkerSize < markerSizeMin || *format.MarkerSize > markerSizeMax) {
		return fmt.Errorf(
			"series format marker_size must be between %d and %d, got %d",
			markerSizeMin, markerSizeMax, *format.MarkerSize,
		)
	}
	return nil
}

func validateSeriesFormatColors(format common.ChartSeriesFormat) error {
	for _, pair := range []struct {
		value *string
		name  string
	}{
		{format.FillColor, "fill_color"},
		{format.LineColor, "line_color"},
		{format.MarkerFillColor, "marker_fill_color"},
		{format.MarkerLineColor, "marker_line_color"},
	} {
		if pair.value == nil {
			continue
		}
		if _, err := editorshape.NormalizeHexColor(*pair.value); err != nil {
			return fmt.Errorf("series format %s: %w", pair.name, err)
		}
	}
	return nil
}

func patchChartSeriesFormats(xml string, formats []common.ChartSeriesFormat) string {
	if len(formats) == 0 {
		return xml
	}
	bySeries := map[int]common.ChartSeriesFormat{}
	for _, format := range formats {
		bySeries[format.SeriesIndex] = format
	}
	seriesIndex := -1
	return reSerBlocks.ReplaceAllStringFunc(xml, func(ser string) string {
		seriesIndex++
		format, ok := bySeries[seriesIndex]
		if !ok {
			return ser
		}
		ser = applySeriesShapeProperties(ser, format)
		ser = applySeriesMarker(ser, format)
		return applySeriesSmooth(ser, format.Smooth)
	})
}

// applySeriesShapeProperties merges the fill and line into the series' own
// c:spPr, leaving whatever else it carries — a gradient, an effect — in place.
func applySeriesShapeProperties(ser string, format common.ChartSeriesFormat) string {
	fill := chartFillNode(format.FillColor, format.NoFill)
	line := buildChartLine(&common.ChartLineFormat{
		Color: format.LineColor, WidthEMU: format.LineWidthEMU,
		Dash: format.LineDash, None: format.NoLine,
	})
	if fill == "" && line == "" {
		return ser
	}
	current, start := seriesChildBlock(ser, reChartLinesShapeProps)
	if current == "" {
		node := "<c:spPr>" + fill + line + "</c:spPr>"
		return insertSeriesFormatNode(ser, node, "<c:marker>")
	}
	updated := current
	if fill != "" {
		updated = setShapePropertiesFill(updated, fill)
	}
	updated = setShapePropertiesLine(updated, line)
	return ser[:start] + updated + ser[start+len(current):]
}

// setShapePropertiesFill replaces the fill of a c:spPr. CT_ShapeProperties puts
// the fill first, after the optional a:xfrm and geometry; only the part before
// a:ln is searched, so the line's own fill is never mistaken for the shape's.
func setShapePropertiesFill(spPr string, fill string) string {
	// Cut returns the whole string when there is no line, which is the search
	// range wanted in that case too.
	head, _, _ := strings.Cut(spPr, "<a:ln")
	if current := reShapeFillBlock.FindString(head); current != "" {
		return strings.Replace(spPr, current, fill, 1)
	}
	for _, anchor := range []string{"<a:ln", "<a:effectLst", "</c:spPr>"} {
		if index := strings.Index(spPr, anchor); index >= 0 {
			return spPr[:index] + fill + spPr[index:]
		}
	}
	return spPr
}

// applySeriesMarker merges the marker fields into the series' c:marker.
// CT_Marker orders symbol, size, then spPr.
func applySeriesMarker(ser string, format common.ChartSeriesFormat) string {
	if format.MarkerSymbol == nil && format.MarkerSize == nil && format.MarkerFillColor == nil &&
		format.MarkerLineColor == nil && format.MarkerLineWidthEMU == nil &&
		format.MarkerNoFill == nil && format.MarkerNoLine == nil {
		return ser
	}
	current, start := seriesChildBlock(ser, reSeriesMarker)
	symbol, size := format.MarkerSymbol, format.MarkerSize
	fill, line := format.MarkerFillColor, format.MarkerLineColor
	if current != "" {
		symbol, size = markerCarryOver(current, symbol, size)
		if fill == nil && format.MarkerNoFill == nil {
			fill = existingMarkerColor(current, reChartLineColor, false)
		}
		if line == nil && format.MarkerNoLine == nil {
			line = existingMarkerColor(current, reChartLineColor, true)
		}
	}
	node := buildSeriesMarker(format, symbol, size, fill, line)
	if current == "" {
		return insertSeriesFormatNode(ser, node, "<c:dPt>")
	}
	return ser[:start] + node + ser[start+len(current):]
}

func buildSeriesMarker(
	format common.ChartSeriesFormat,
	symbol *string,
	size *int,
	fill *string,
	line *string,
) string {
	var b strings.Builder
	b.WriteString("<c:marker>")
	if symbol != nil {
		b.WriteString(`<c:symbol val="` + *symbol + `"/>`)
	}
	if size != nil {
		b.WriteString(`<c:size val="` + strconv.Itoa(*size) + `"/>`)
	}
	fillNode := chartFillNode(fill, format.MarkerNoFill)
	lineNode := buildChartLine(&common.ChartLineFormat{
		Color: line, WidthEMU: format.MarkerLineWidthEMU, None: format.MarkerNoLine,
	})
	if fillNode != "" || lineNode != "" {
		b.WriteString("<c:spPr>" + fillNode + lineNode + "</c:spPr>")
	}
	b.WriteString("</c:marker>")
	return b.String()
}

func applySeriesSmooth(ser string, smooth *bool) string {
	if smooth == nil {
		return ser
	}
	node := `<c:smooth val="` + boolToOneZero(*smooth) + `"/>`
	if reSeriesSmooth.MatchString(ser) {
		return reSeriesSmooth.ReplaceAllLiteralString(ser, node)
	}
	// CT_LineSer puts c:smooth last, before c:extLst.
	if index := strings.Index(ser, chartExtensionListTagPrefix); index >= 0 {
		return ser[:index] + node + ser[index:]
	}
	if index := strings.LastIndex(ser, seriesCloseTag); index >= 0 {
		return ser[:index] + node + ser[index:]
	}
	return ser
}

// seriesChildBlock finds a c:ser child that precedes the data points and
// references, and returns it with its offset. A match after the boundary
// belongs to a c:dPt or a data label, not to the series.
func seriesChildBlock(ser string, re *regexp.Regexp) (string, int) {
	location := re.FindStringIndex(ser)
	if location == nil {
		return "", -1
	}
	if boundary := seriesFormatBoundary(ser); boundary >= 0 && location[0] > boundary {
		return "", -1
	}
	return ser[location[0]:location[1]], location[0]
}

func seriesFormatBoundary(ser string) int {
	boundary := -1
	for _, tag := range seriesFormatBoundaryTags {
		index := strings.Index(ser, tag)
		if index >= 0 && (boundary < 0 || index < boundary) {
			boundary = index
		}
	}
	return boundary
}

// insertSeriesFormatNode splices a new c:spPr or c:marker into its schema slot:
// CT_Ser orders idx, order, tx, spPr, marker, dPt, dLbls, and the references.
func insertSeriesFormatNode(ser string, node string, preferred string) string {
	anchors := append([]string{preferred}, seriesFormatBoundaryTags...)
	anchors = append(anchors, "<c:invertIfNegative", "<c:pictureOptions", chartExtensionListTagPrefix)
	for _, anchor := range anchors {
		if index := strings.Index(ser, anchor); index >= 0 {
			return ser[:index] + node + ser[index:]
		}
	}
	if index := strings.LastIndex(ser, seriesCloseTag); index >= 0 {
		return ser[:index] + node + ser[index:]
	}
	return ser
}
