package presentation

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// A generated chart used to be a lone chartN.xml: no rels, no workbook and no
// <c:externalData>, so "Edit Data" in PowerPoint had nothing to open.
func TestGeneratedChartCarriesItsWorkbook(t *testing.T) {
	bar := charts.NewBarChart([]string{"A", "B"}, []float64{1, 2})
	slide := elements.NewSlide("Chart")
	slide.Chart = &bar

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})

	for _, want := range []string{
		"ppt/charts/chart1.xml",
		"ppt/charts/_rels/chart1.xml.rels",
		"ppt/embeddings/Microsoft_Excel_Worksheet1.xlsx",
	} {
		if _, ok := parts[want]; !ok {
			t.Fatalf("package is missing %s", want)
		}
	}

	if !strings.Contains(parts["ppt/charts/chart1.xml"], `<c:externalData r:id="rId1">`) {
		t.Fatalf("chart1.xml has no externalData: %s", parts["ppt/charts/chart1.xml"])
	}
	if !strings.Contains(
		parts["ppt/charts/_rels/chart1.xml.rels"],
		`Target="../embeddings/Microsoft_Excel_Worksheet1.xlsx"`,
	) {
		t.Fatalf("chart rels do not point at the workbook: %s", parts["ppt/charts/_rels/chart1.xml.rels"])
	}
	if !strings.Contains(
		parts["[Content_Types].xml"],
		`PartName="/ppt/embeddings/Microsoft_Excel_Worksheet1.xlsx"`,
	) {
		t.Fatal("no content type declared for the embedded workbook")
	}
}
