package pptxxml

import (
	"strings"
	"testing"
)

func boundChartSpec(relID string) *ChartSpec {
	return &ChartSpec{
		Kind:           ChartKindBar,
		BarDir:         "col",
		Grouping:       "clustered",
		Title:          "T",
		SeriesName:     "Series 1",
		Color:          "4472C4",
		Categories:     []string{"A", "B", "C"},
		Values:         []float64{1, 2, 3},
		ExternalDataID: relID,
	}
}

// A chart that ships a workbook has to address its cells, or editing them in
// Excel changes nothing on the slide.
func TestBoundChartSeriesReferencesTheWorkbook(t *testing.T) {
	xml := ChartPartXML(boundChartSpec("rId1"))

	for _, want := range []string{
		"<c:strRef>", "<c:numRef>",
		"Sheet1!$A$2:$A$4", "Sheet1!$B$2:$B$4", "Sheet1!$B$1",
		"<c:strCache>", "<c:numCache>",
	} {
		if !strings.Contains(xml, want) {
			t.Errorf("bound chart XML is missing %q", want)
		}
	}
	// The cached points stay, so the chart still draws before the workbook is
	// ever opened.
	if !strings.Contains(xml, "<c:v>A</c:v>") || !strings.Contains(xml, "<c:v>3.000000</c:v>") {
		t.Errorf("bound chart lost its cached points:\n%s", xml)
	}
	if strings.Contains(xml, "<c:strLit>") || strings.Contains(xml, "<c:numLit>") {
		t.Errorf("bound chart still writes literals:\n%s", xml)
	}
}

// Without a workbook there is nothing for a formula to point at, so the literal
// form is the correct one.
func TestUnboundChartSeriesStaysLiteral(t *testing.T) {
	xml := ChartPartXML(boundChartSpec(""))

	if !strings.Contains(xml, "<c:strLit>") || !strings.Contains(xml, "<c:numLit>") {
		t.Errorf("unbound chart should write literals:\n%s", xml)
	}
	if strings.Contains(xml, "<c:strRef>") || strings.Contains(xml, "<c:numRef>") {
		t.Errorf("unbound chart references a workbook it does not have:\n%s", xml)
	}
}

func TestChartFormulasCoverEveryPoint(t *testing.T) {
	if got := chartCategoryFormula(5); got != "Sheet1!$A$2:$A$6" {
		t.Errorf("category formula = %q, want Sheet1!$A$2:$A$6", got)
	}
	if got := chartValueFormula(1); got != "Sheet1!$B$2:$B$2" {
		t.Errorf("single-point value formula = %q, want Sheet1!$B$2:$B$2", got)
	}
	// An empty series must still address a syntactically valid range.
	if got := chartValueFormula(0); got != "Sheet1!$B$2:$B$2" {
		t.Errorf("empty value formula = %q, want Sheet1!$B$2:$B$2", got)
	}
}
