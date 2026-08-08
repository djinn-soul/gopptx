package pptxxml

import (
	"strconv"
	"strings"
)

// A chart that ships an embedded workbook has to say which cells of it feed
// each series, or "Edit Data" opens a workbook whose edits change nothing: the
// chart reads its own cached literals and has no formula tying them to a sheet.
//
// The workbook the writer generates has a fixed shape — one sheet named Sheet1,
// the headers in row 1, and the data from row 2 down, with categories in column
// A and the series values in column B (see GenerateExcelForChart). These build
// the c:f formulas that address it.
//
// A chart with no embedded workbook keeps writing literals, because there is no
// workbook for a formula to point at.

const (
	// chartSheetName is the single sheet of a generated chart workbook.
	chartSheetName = "Sheet1"
	// chartCategoryColumn and chartValueColumn are where GenerateExcelForChart
	// puts the categories and the first series.
	chartCategoryColumn = "A"
	chartValueColumn    = "B"
	// chartFirstDataRow is the first row below the header.
	chartFirstDataRow = 2
)

// chartCategoryFormula addresses the category cells of a generated workbook.
func chartCategoryFormula(count int) string {
	return chartColumnFormula(chartCategoryColumn, count)
}

// chartValueFormula addresses the value cells of a generated workbook.
func chartValueFormula(count int) string {
	return chartColumnFormula(chartValueColumn, count)
}

// chartSeriesNameFormula addresses the header cell the series is named by.
func chartSeriesNameFormula() string {
	return chartSheetName + "!$" + chartValueColumn + "$1"
}

func chartColumnFormula(column string, count int) string {
	last := chartFirstDataRow + count - 1
	if count <= 0 {
		last = chartFirstDataRow
	}
	return chartSheetName + "!$" + column + "$" + strconv.Itoa(chartFirstDataRow) +
		":$" + column + "$" + strconv.Itoa(last)
}

// writeChartCategories writes the <c:cat> block, as a reference into the
// embedded workbook when there is one and as a literal otherwise. The cached
// points are written either way, so a reader that ignores the formula — and
// PowerPoint itself, before the workbook is opened — still draws the chart.
func writeChartCategories(b *strings.Builder, categories []string, bound bool) {
	if bound {
		b.WriteString("\n<c:cat><c:strRef>\n<c:f>")
		b.WriteString(Escape(chartCategoryFormula(len(categories))))
		b.WriteString("</c:f>\n<c:strCache>")
		writeChartStringPoints(b, categories)
		b.WriteString("\n</c:strCache></c:strRef></c:cat>")
		return
	}
	b.WriteString("\n<c:cat><c:strLit>")
	writeChartStringPoints(b, categories)
	b.WriteString("\n</c:strLit></c:cat>")
}

// writeChartValues writes the <c:val> block, referencing the workbook when the
// chart is bound to one.
func writeChartValues(b *strings.Builder, values []float64, bound bool) {
	if bound {
		b.WriteString("\n<c:val><c:numRef>\n<c:f>")
		b.WriteString(Escape(chartValueFormula(len(values))))
		b.WriteString("</c:f>\n<c:numCache>\n<c:formatCode>General</c:formatCode>")
		writeNumericPoints(b, values)
		b.WriteString("\n</c:numCache></c:numRef></c:val>")
		return
	}
	// The stray tab before the closing tag is what this block has always
	// emitted. It is meaningless to XML, but the chart fixtures hash the bytes,
	// and an unbound chart's output should not move at all for a change that
	// only adds a bound form.
	b.WriteString("\n<c:val><c:numLit>\n<c:formatCode>General</c:formatCode>")
	writeNumericPoints(b, values)
	b.WriteString("\n\t</c:numLit></c:val>")
}

// writeChartSeriesName writes <c:tx>, pointing at the workbook's header cell
// when the chart is bound so that renaming the series in Excel takes effect.
func writeChartSeriesName(b *strings.Builder, name string, bound bool) {
	if bound {
		b.WriteString("\n<c:tx><c:strRef>\n<c:f>")
		b.WriteString(Escape(chartSeriesNameFormula()))
		b.WriteString("</c:f>\n<c:strCache>\n<c:ptCount val=\"1\"/>\n<c:pt idx=\"0\"><c:v>")
		b.WriteString(Escape(name))
		b.WriteString("</c:v></c:pt>\n</c:strCache></c:strRef></c:tx>")
		return
	}
	b.WriteString("\n<c:tx><c:v>" + Escape(name) + "</c:v></c:tx>")
}

func writeChartStringPoints(b *strings.Builder, values []string) {
	b.WriteString("\n<c:ptCount val=\"")
	b.WriteString(strconv.Itoa(len(values)))
	b.WriteString("\"/>")
	for i, value := range values {
		b.WriteString("\n<c:pt idx=\"")
		b.WriteString(strconv.Itoa(i))
		b.WriteString("\"><c:v>")
		b.WriteString(Escape(value))
		b.WriteString("</c:v></c:pt>")
	}
}

// chartIsWorkbookBound reports whether the spec ships an embedded workbook, so
// its series can address cells rather than carrying loose literals.
func chartIsWorkbookBound(chart *ChartSpec) bool {
	return chart.ExternalDataID != ""
}
