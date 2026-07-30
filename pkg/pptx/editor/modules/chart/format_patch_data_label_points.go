package chart

import (
	"fmt"
	"html"
	"regexp"
	"sort"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

var (
	reDataLabelLayout    = regexp.MustCompile(`(?s)<c:layout>.*?</c:layout>|<c:layout\s*/>`)
	reDataLabelRichText  = regexp.MustCompile(`(?s)<c:tx>.*?</c:tx>`)
	reDataLabelNumFmt    = regexp.MustCompile(`<c:numFmt\b[^>]*/>`)
	reDataLabelShapePr   = regexp.MustCompile(`(?s)<c:spPr>.*?</c:spPr>|<c:spPr\s*/>`)
	reDataLabelTextPr    = regexp.MustCompile(`(?s)<c:txPr>.*?</c:txPr>`)
	reDataLabelPosition  = regexp.MustCompile(`<c:dLblPos\b[^>]*/>`)
	reDataLabelSeparator = regexp.MustCompile(`(?s)<c:separator>.*?</c:separator>`)
	reDataLabelDelete    = regexp.MustCompile(`<c:delete\b[^>]*/>`)
	reDataLabelFormatVal = regexp.MustCompile(`formatCode="([^"]*)"`)
	reDataLabelLinkedVal = regexp.MustCompile(`sourceLinked="([^"]*)"`)
	reDataLabelFontSize  = regexp.MustCompile(`<a:defRPr\b[^>]*\bsz="([^"]+)"`)
	reDataLabelFontBold  = regexp.MustCompile(`<a:defRPr\b[^>]*\bb="([^"]+)"`)
	reDataLabelFontColor = regexp.MustCompile(`(?s)<a:defRPr\b.*?<a:srgbClr val="([^"]+)"`)
	reDataLabelIdxValue  = regexp.MustCompile(`<c:idx val="([^"]+)"`)

	// reDataLabelFlag matches any one CT_DLbl display flag, capturing its name
	// and value, so a whole set is read or stripped in a single pass.
	reDataLabelFlag = regexp.MustCompile(
		`<c:(showLegendKey|showVal|showCatName|showSerName|showPercent|showBubbleSize)` +
			` val="([^"]*)"/>`,
	)
)

const (
	flagShowLegendKey  = "showLegendKey"
	flagShowValue      = "showVal"
	flagShowCategory   = "showCatName"
	flagShowSeriesName = "showSerName"
	flagShowPercent    = "showPercent"
	flagShowBubbleSize = "showBubbleSize"
)

// dataLabelFlagValues reads every display flag present in a dLbl or dLbls block.
func dataLabelFlagValues(xml string) map[string]string {
	values := map[string]string{}
	for _, match := range reDataLabelFlag.FindAllStringSubmatch(xml, -1) {
		values[match[1]] = match[2]
	}
	return values
}

// labelAttribute returns the first capture of re, or "" when it does not match.
func labelAttribute(xml string, re *regexp.Regexp) string {
	match := re.FindStringSubmatch(xml)
	if len(match) != expectedSingleGroupMatch {
		return ""
	}
	return match[1]
}

const (
	dataLabelFontSizeMinPt = 1
	dataLabelFontSizeMaxPt = 400
	pointsPerFontSizeUnit  = 100
)

func validateChartDataLabelPoints(points []common.ChartDataLabelPoint) error {
	for _, point := range points {
		if point.SeriesIndex < 0 {
			return fmt.Errorf(
				"data label point series_index must not be negative, got %d", point.SeriesIndex,
			)
		}
		if point.PointIndex < 0 {
			return fmt.Errorf(
				"data label point point_index must not be negative, got %d", point.PointIndex,
			)
		}
		if point.FontColor != nil {
			if _, err := editorshape.NormalizeHexColor(*point.FontColor); err != nil {
				return fmt.Errorf("data label point font_color: %w", err)
			}
		}
		if point.FontSizePt != nil &&
			(*point.FontSizePt < dataLabelFontSizeMinPt || *point.FontSizePt > dataLabelFontSizeMaxPt) {
			return fmt.Errorf(
				"data label point font_size_pt must be between %d and %d, got %d",
				dataLabelFontSizeMinPt, dataLabelFontSizeMaxPt, *point.FontSizePt,
			)
		}
		if err := validateDataLabelPointBox(point); err != nil {
			return err
		}
	}
	return nil
}

// patchChartDataLabelPoints merges per-label formatting into each addressed
// series, rewriting only the c:dLbl elements it names.
func patchChartDataLabelPoints(xml string, points []common.ChartDataLabelPoint) string {
	if len(points) == 0 {
		return xml
	}

	bySeries := map[int][]common.ChartDataLabelPoint{}
	for _, point := range points {
		bySeries[point.SeriesIndex] = append(bySeries[point.SeriesIndex], point)
	}
	plotLabels := seriesDataLabelsBlock(stripSeriesBlocks(xml))
	plotFlags := dataLabelFlagsFrom(plotLabels)
	plotNumberFormat := reDataLabelNumFmt.FindString(plotLabels)

	seriesIndex := -1
	return reSerBlocks.ReplaceAllStringFunc(xml, func(ser string) string {
		seriesIndex++
		requested, ok := bySeries[seriesIndex]
		if !ok {
			return ser
		}
		// Highest point index first: a new <c:dLbl> is inserted at the front of
		// the dLbls block, so descending order leaves them ascending.
		sorted := append([]common.ChartDataLabelPoint(nil), requested...)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].PointIndex > sorted[j].PointIndex })
		inherited := seriesNumberFormat(ser, plotNumberFormat)
		for _, point := range sorted {
			ser = applyDataLabelPoint(ser, point, plotFlags, inherited)
		}
		return seedSeriesDataLabels(ser, plotNumberFormat, plotFlags)
	})
}

// seriesNumberFormat returns the number format a per-point label inherits: the
// series c:dLbls setting, else the plot-level one.
func seriesNumberFormat(ser string, plotNumberFormat string) string {
	if current := reDataLabelNumFmt.FindString(seriesDataLabelsBlock(ser)); current != "" {
		return current
	}
	return plotNumberFormat
}

// seedSeriesDataLabels copies the plot-level number format and display flags
// into the series c:dLbls when it has none.
//
// Hosting a per-point <c:dLbl> requires a series-level <c:dLbls>, and that block
// shadows the plot-level one for the whole series: without the copy, formatting
// one label silently drops the chart-wide number format from every other label
// of that series.
func seedSeriesDataLabels(ser string, plotNumberFormat string, plotFlags string) string {
	block := reDataLabelsBlock.FindString(ser)
	if block == "" {
		return ser
	}
	seriesOwn := seriesDataLabelsBlock(ser)
	updated := block
	if plotNumberFormat != "" && !reDataLabelNumFmt.MatchString(seriesOwn) {
		updated = strings.Replace(updated, "<c:dLbls>", "<c:dLbls>"+plotNumberFormat, 1)
	}
	if flags := dataLabelFlagsFrom(seriesOwn); flags == "" && plotFlags != "" {
		updated = strings.Replace(
			updated, "</c:dLbls>", completeDataLabelFlags(plotFlags)+"</c:dLbls>", 1,
		)
	}
	return strings.Replace(ser, block, updated, 1)
}

func applyDataLabelPoint(
	ser string, point common.ChartDataLabelPoint, plotFlags string, inheritedNumFmt string,
) string {
	existing, start, end := findDataLabelBlock(ser, point.PointIndex)
	label := buildDataLabelPointBlock(ser, existing, point, plotFlags, inheritedNumFmt)
	if existing != "" {
		return ser[:start] + label + ser[end:]
	}
	return insertDataLabelBlock(ser, label)
}

// buildDataLabelPointBlock renders one c:dLbl, keeping the parts of an existing
// label the request does not mention. CT_DLbl orders its children idx, then
// either delete or the layout, tx, numFmt, spPr, txPr, dLblPos, display flag
// and separator group.
func buildDataLabelPointBlock(
	ser string, existing string, point common.ChartDataLabelPoint,
	plotFlags string, inheritedNumFmt string,
) string {
	var b strings.Builder
	b.WriteString(`<c:dLbl><c:idx val="` + strconv.Itoa(point.PointIndex) + `"/>`)

	if deleted := labelIsDeleted(existing, point); deleted {
		b.WriteString(`<c:delete val="1"/></c:dLbl>`)
		return b.String()
	}

	b.WriteString(reDataLabelLayout.FindString(existing))
	b.WriteString(reDataLabelRichText.FindString(existing))
	b.WriteString(buildDataLabelNumberFormat(existing, point, inheritedNumFmt))
	b.WriteString(buildDataLabelShapeProperties(existing, point))
	b.WriteString(buildDataLabelTextProperties(existing, point))
	b.WriteString(reDataLabelPosition.FindString(existing))
	b.WriteString(applyDataLabelFlagOverrides(seriesDataLabelFlags(ser, existing, plotFlags), point))
	b.WriteString(reDataLabelSeparator.FindString(existing))
	b.WriteString("</c:dLbl>")
	return b.String()
}

// labelIsDeleted reports whether the rebuilt label must carry c:delete: the
// request wins, and otherwise a label already deleted stays deleted.
func labelIsDeleted(existing string, point common.ChartDataLabelPoint) bool {
	if point.Delete != nil {
		return *point.Delete
	}
	deleted := reDataLabelDelete.FindString(existing)
	return deleted != "" && !strings.Contains(deleted, `val="0"`)
}

// buildDataLabelNumberFormat renders the label's c:numFmt, falling back to the
// format inherited from the series or plot. A c:dLbl that omits it is drawn in
// the general format, so a label patched for its font alone would otherwise lose
// the chart-wide number format.
func buildDataLabelNumberFormat(
	existing string, point common.ChartDataLabelPoint, inherited string,
) string {
	format, linked := "", false
	current := reDataLabelNumFmt.FindString(existing)
	if current == "" {
		current = inherited
	}
	if current != "" {
		// The attribute is already escaped, so it is decoded here and escaped
		// again below; otherwise a second patch would double-escape a code
		// containing quotes, such as `0.0 "kg"`.
		format = html.UnescapeString(labelAttribute(current, reDataLabelFormatVal))
		linked = labelAttribute(current, reDataLabelLinkedVal) == "1"
	}
	if point.NumberFormat != nil {
		format = *point.NumberFormat
		// A label whose format is linked to the source ignores its own code.
		linked = false
	}
	if point.FormatLinked != nil {
		linked = *point.FormatLinked
	}
	if format == "" {
		return ""
	}
	return `<c:numFmt formatCode="` + common.XMLEscape(format) +
		`" sourceLinked="` + boolToOneZero(linked) + `"/>`
}

// buildDataLabelTextProperties rebuilds the label's c:txPr from the font
// properties this API models, seeded with whatever the label already carried.
func buildDataLabelTextProperties(existing string, point common.ChartDataLabelPoint) string {
	current := reDataLabelTextPr.FindString(existing)
	if point.FontColor == nil && point.FontSizePt == nil && point.FontBold == nil {
		return current
	}

	size := labelAttribute(current, reDataLabelFontSize)
	if point.FontSizePt != nil {
		size = strconv.Itoa(*point.FontSizePt * pointsPerFontSizeUnit)
	}
	bold := labelAttribute(current, reDataLabelFontBold)
	if point.FontBold != nil {
		bold = boolToOneZero(*point.FontBold)
	}
	color := labelAttribute(current, reDataLabelFontColor)
	if point.FontColor != nil {
		if normalized, err := editorshape.NormalizeHexColor(*point.FontColor); err == nil {
			color = normalized
		}
	}

	attrs := ""
	if size != "" {
		attrs += ` sz="` + size + `"`
	}
	if bold != "" {
		attrs += ` b="` + bold + `"`
	}
	fill := ""
	if color != "" {
		fill = `<a:solidFill><a:srgbClr val="` + color + `"/></a:solidFill>`
	}
	return `<c:txPr><a:bodyPr/><a:lstStyle/><a:p><a:pPr><a:defRPr` + attrs + `>` +
		fill + `</a:defRPr></a:pPr><a:endParaRPr lang="en-US"/></a:p></c:txPr>`
}

// applyDataLabelFlagOverrides layers the requested display flags over the ones
// carried across, keeping CT_DLbl's flag order. Setting a label's colour must
// not silently drop the category name it was already showing (upstream #650).
func applyDataLabelFlagOverrides(flags string, point common.ChartDataLabelPoint) string {
	overrides := map[string]*bool{
		flagShowLegendKey:  point.ShowLegendKey,
		flagShowValue:      point.ShowValue,
		flagShowCategory:   point.ShowCategory,
		flagShowSeriesName: point.ShowSeriesName,
		flagShowPercent:    point.ShowPercent,
	}

	values := dataLabelFlagValues(flags)
	var b strings.Builder
	for _, name := range dataLabelFlagNames() {
		value := values[name]
		if override, ok := overrides[name]; ok && override != nil {
			value = boolToOneZero(*override)
		}
		if value == "" {
			continue
		}
		b.WriteString(`<c:` + name + ` val="` + value + `"/>`)
	}
	return b.String()
}
