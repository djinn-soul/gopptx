package chart

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var (
	reDataTableBlock = regexp.MustCompile(`(?s)<c:dTable>.*?</c:dTable>|<c:dTable/>`)
	reDTableHorz     = regexp.MustCompile(`<c:showHorzBorder val="([^"]+)"`)
	reDTableVert     = regexp.MustCompile(`<c:showVertBorder val="([^"]+)"`)
	reDTableOutline  = regexp.MustCompile(`<c:showOutline val="([^"]+)"`)
	reDTableKeys     = regexp.MustCompile(`<c:showKeys val="([^"]+)"`)
	reDTableFontSize = regexp.MustCompile(`<a:defRPr\b[^>]*\bsz="([^"]+)"`)
	reDTableSpPr     = regexp.MustCompile(`(?s)<c:spPr\b[^>]*>.*?</c:spPr>|<c:spPr\b[^>]*/>`)
	reDTableTxPr     = regexp.MustCompile(`(?s)<c:txPr\b[^>]*>.*?</c:txPr>|<c:txPr\b[^>]*/>`)
	reDTableDefRPr   = regexp.MustCompile(`<a:defRPr\b[^>]*>`)
)

const (
	dataTableFontMinPt = 1
	dataTableFontMaxPt = 400
	pointsToHundredths = 100
)

func validateChartDataTable(table *common.ChartDataTable) error {
	if table == nil || table.FontSizePt == nil {
		return nil
	}
	if *table.FontSizePt < dataTableFontMinPt || *table.FontSizePt > dataTableFontMaxPt {
		return errors.New("data_table font_size_pt must be between 1 and 400")
	}
	return nil
}

// patchChartDataTable writes the c:dTable under the plot area. CT_PlotArea puts
// it after the chart groups and axes and before the plot-area c:spPr.
func patchChartDataTable(xml string, table *common.ChartDataTable) string {
	if table == nil {
		return xml
	}
	if !table.Show {
		return reDataTableBlock.ReplaceAllLiteralString(xml, "")
	}
	node := buildDataTableBlock(xml, *table)
	if reDataTableBlock.MatchString(xml) {
		return reDataTableBlock.ReplaceAllLiteralString(xml, node)
	}
	return insertDataTableBlock(xml, node)
}

// buildDataTableBlock renders the c:dTable, carrying over whatever the chart
// already had so a patch of one flag leaves the others alone. CT_DTable orders
// showHorzBorder, showVertBorder, showOutline, showKeys, spPr, txPr.
func buildDataTableBlock(xml string, table common.ChartDataTable) string {
	existing := reDataTableBlock.FindString(xml)
	var b strings.Builder
	b.WriteString("<c:dTable>")
	for _, flag := range []struct {
		tag      string
		value    *bool
		pattern  *regexp.Regexp
		fallback bool
	}{
		{"showHorzBorder", table.ShowHorizontalBorder, reDTableHorz, true},
		{"showVertBorder", table.ShowVerticalBorder, reDTableVert, true},
		{"showOutline", table.ShowOutline, reDTableOutline, true},
		{"showKeys", table.ShowKeys, reDTableKeys, true},
	} {
		value := flag.fallback
		if match := flag.pattern.FindStringSubmatch(existing); len(match) == expectedSingleGroupMatch {
			value = strings.TrimSpace(match[1]) == "1"
		}
		if flag.value != nil {
			value = *flag.value
		}
		b.WriteString(`<c:` + flag.tag + ` val="` + boolToOneZero(value) + `"/>`)
	}
	b.WriteString(reDTableSpPr.FindString(existing))
	b.WriteString(buildDataTableTextProperties(existing, table.FontSizePt))
	b.WriteString("</c:dTable>")
	return b.String()
}

// buildDataTableTextProperties renders the c:txPr font size, which OOXML
// records in hundredths of a point.
func buildDataTableTextProperties(existing string, fontSizePt *int) string {
	existingTxPr := reDTableTxPr.FindString(existing)
	if fontSizePt == nil {
		return existingTxPr
	}
	size := *fontSizePt * pointsToHundredths
	if existingTxPr != "" && reDTableDefRPr.MatchString(existingTxPr) {
		return reDTableDefRPr.ReplaceAllStringFunc(existingTxPr, func(tag string) string {
			return setXMLAttribute(tag, "sz", strconv.Itoa(size))
		})
	}

	if size <= 0 {
		return ""
	}
	return `<c:txPr><a:bodyPr/><a:lstStyle/><a:p><a:pPr><a:defRPr sz="` +
		strconv.Itoa(size) + `"/></a:pPr><a:endParaRPr lang="en-US"/></a:p></c:txPr>`
}

func insertDataTableBlock(xml string, node string) string {
	start, end := plotAreaBounds(xml)
	if start < 0 || end <= start {
		return xml
	}
	block := xml[start:end]
	// The plot-area c:spPr is the only sibling that follows c:dTable. It has to
	// be matched at the plot area's own nesting level: every chart group and
	// series carries a c:spPr of its own.
	insertAt := plotAreaShapePropertiesIndex(block)
	if insertAt < 0 {
		insertAt = len(block)
	}
	patched := block[:insertAt] + node + block[insertAt:]
	return xml[:start] + patched + xml[end:]
}

// plotAreaShapePropertiesIndex returns the offset of the plot area's own
// c:spPr, or -1 when it has none. Nested copies inside chart groups, series,
// and axes are skipped by tracking element depth.
func plotAreaShapePropertiesIndex(block string) int {
	depth := 0
	for offset := 0; offset < len(block); {
		open := strings.IndexByte(block[offset:], '<')
		if open < 0 {
			return -1
		}
		open += offset
		closing := strings.IndexByte(block[open:], '>')
		if closing < 0 {
			return -1
		}
		closing += open
		tag := block[open : closing+1]
		switch {
		case strings.HasPrefix(tag, "</"):
			depth--
		case strings.HasSuffix(tag, "/>"):
			// Self-closing: no depth change.
		default:
			if depth == 0 && strings.HasPrefix(tag, "<c:spPr") {
				return open
			}
			depth++
		}
		offset = closing + 1
	}
	return -1
}

// plotAreaBounds returns the range of the plot area's children.
func plotAreaBounds(xml string) (int, int) {
	const openTag, closeTag = "<c:plotArea>", "</c:plotArea>"
	start := strings.Index(xml, openTag)
	if start < 0 {
		return -1, -1
	}
	start += len(openTag)
	endRel := strings.Index(xml[start:], closeTag)
	if endRel < 0 {
		return -1, -1
	}
	return start, start + endRel
}

// parseDataTableState reads the c:dTable back into the update payload shape.
func parseDataTableState(xml string) *common.ChartDataTable {
	block := reDataTableBlock.FindString(xml)
	if block == "" {
		return nil
	}
	table := common.ChartDataTable{Show: true}
	table.ShowHorizontalBorder = trendlineBoolValue(block, reDTableHorz)
	table.ShowVerticalBorder = trendlineBoolValue(block, reDTableVert)
	table.ShowOutline = trendlineBoolValue(block, reDTableOutline)
	table.ShowKeys = trendlineBoolValue(block, reDTableKeys)
	if size := trendlineIntValue(block, reDTableFontSize); size != nil && *size > 0 {
		points := *size / pointsToHundredths
		table.FontSizePt = &points
	}
	return &table
}
