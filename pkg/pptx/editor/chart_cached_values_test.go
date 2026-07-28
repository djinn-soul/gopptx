package editor

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// Upstream #115: a deck rebuilt from a linked workbook needs its cached numbers
// refreshed without each chart being opened and refreshed by hand — and without
// the link being replaced by an embedded workbook.
func TestUpdateChartCachedValuesKeepsExternalLink(t *testing.T) {
	ed, chartPart := newChartFixtureEditor(t, "chart-cached-values.pptx")
	defer func() { _ = ed.Close() }()

	linkChartToExternalWorkbook(t, ed, chartPart)

	err := ed.UpdateChartCachedValues(0, firstChartSelector(), common.ChartDataUpdate{
		Categories: []string{"Q1", "Q2"},
		Series: []common.ChartSeriesData{
			{Name: strPtr("Revenue"), Values: []float64{11, 22}},
		},
	})
	if err != nil {
		t.Fatalf("update cached values: %v", err)
	}

	chartXML, ok := ed.parts.Get(chartPart)
	if !ok {
		t.Fatalf("chart part missing after update")
	}
	text := string(chartXML)
	if !strings.Contains(text, "<c:v>11</c:v>") || !strings.Contains(text, "<c:v>22</c:v>") {
		t.Fatalf("expected the refreshed values in the cache:\n%s", text)
	}
	if !strings.Contains(text, "<c:externalData") {
		t.Fatalf("the external data link was dropped:\n%s", text)
	}

	source, err := ed.GetChartDataSource(0, firstChartSelector())
	if err != nil {
		t.Fatalf("chart data source: %v", err)
	}
	if source.Kind != ChartDataSourceExternal {
		t.Fatalf("expected an external source, got %+v", source)
	}
	if source.Target != "file:///C:/reports/weekly.xlsx" {
		t.Fatalf("unexpected link target: %q", source.Target)
	}
	if source.AutoUpdate == nil || !*source.AutoUpdate {
		t.Fatalf("expected autoUpdate to be reported as true, got %+v", source.AutoUpdate)
	}
	if source.PartPath != "" {
		t.Fatalf("an external link has no package part, got %q", source.PartPath)
	}
}

// The embedded case is the other half of the same question: a caller has to be
// able to tell which kind of chart it is holding.
func TestGetChartDataSourceReportsEmbeddedWorkbook(t *testing.T) {
	ed, _ := newChartFixtureEditor(t, "chart-embedded-source.pptx")
	defer func() { _ = ed.Close() }()

	err := ed.UpdateChartData(0, firstChartSelector(), common.ChartDataUpdate{
		Categories: []string{"Q1", "Q2"},
		Series: []common.ChartSeriesData{
			{Name: strPtr("Revenue"), Values: []float64{5, 6}},
		},
	})
	if err != nil {
		t.Fatalf("update chart data: %v", err)
	}

	source, err := ed.GetChartDataSource(0, firstChartSelector())
	if err != nil {
		t.Fatalf("chart data source: %v", err)
	}
	if source.Kind != ChartDataSourceEmbedded {
		t.Fatalf("expected an embedded source, got %+v", source)
	}
	if !strings.HasPrefix(source.PartPath, "ppt/embeddings/") {
		t.Fatalf("expected the embedded workbook part, got %q", source.PartPath)
	}
}

func firstChartSelector() common.ChartSelector {
	index := 0
	return common.ChartSelector{Index: &index}
}

func newChartFixtureEditor(t *testing.T, name string) (*PresentationEditor, string) {
	t.Helper()

	basePath := writeDeckFixture(t, name, []elements.SlideContent{
		elements.NewSlide("Chart"),
	})
	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}

	chartDef := charts.NewBarChart([]string{"Q1", "Q2"}, []float64{1, 2}).WithTitle("Revenue")
	if err := ed.AddChart(0, chartDef); err != nil {
		_ = ed.Close()
		t.Fatalf("add chart: %v", err)
	}

	charts, err := ed.ListSlideCharts(0)
	if err != nil || len(charts) == 0 {
		_ = ed.Close()
		t.Fatalf("list slide charts: %v (%d found)", err, len(charts))
	}
	return ed, charts[0].ChartPart
}

// linkChartToExternalWorkbook rewrites the chart the way PowerPoint writes one
// whose data lives in a workbook outside the package.
func linkChartToExternalWorkbook(t *testing.T, ed *PresentationEditor, chartPart string) {
	t.Helper()

	chartXML, ok := ed.parts.Get(chartPart)
	if !ok {
		t.Fatalf("chart part %s missing", chartPart)
	}
	text := string(chartXML)
	if strings.Contains(text, "<c:externalData") {
		text = externalDataPattern.ReplaceAllLiteralString(text, "")
	}
	text = replaceOnce(t, text, "</c:chartSpace>",
		`<c:externalData r:id="rIdLink"><c:autoUpdate val="1"/></c:externalData></c:chartSpace>`)
	ed.parts.Set(chartPart, []byte(text))

	relsPath := common.RelsPathFor(chartPart)
	relsData, ok := ed.parts.Get(relsPath)
	if !ok {
		t.Fatalf("chart rels %s missing", relsPath)
	}
	rels := replaceOnce(t, string(relsData), "</Relationships>",
		`<Relationship Id="rIdLink" `+
			`Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/package" `+
			`Target="file:///C:/reports/weekly.xlsx" TargetMode="External"/></Relationships>`)
	ed.parts.Set(relsPath, []byte(rels))
}
