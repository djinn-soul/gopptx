package editor

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestUpdateChartDataCategoryByIndexPreservesStyle(t *testing.T) {
	e := newChartUpdateEditorFixture()
	e.parts.Set("ppt/charts/chart1.xml", []byte(categoryChartXML()))

	idx := 0
	err := e.UpdateChartData(0, common.ChartSelector{Index: &idx}, common.ChartDataUpdate{
		Categories: []string{"Q1", "Q2"},
		Series: []common.ChartSeriesData{
			{Values: []float64{10, 20}},
		},
	})
	if err != nil {
		t.Fatalf("UpdateChartData failed: %v", err)
	}

	chartData, _ := e.parts.Get("ppt/charts/chart1.xml")
	xml := string(chartData)
	if !strings.Contains(xml, "<c:spPr><a:solidFill/></c:spPr>") {
		t.Fatalf("expected style node preserved")
	}
	if !strings.Contains(xml, "Sheet1!$A$2:$A$3") || !strings.Contains(xml, "Sheet1!$B$2:$B$3") {
		t.Fatalf("expected formulas rewritten, got: %s", xml)
	}

	excelPath := e.chartEmbeddings["ppt/charts/chart1.xml"]
	excelData, _ := e.parts.Get(excelPath)
	if _, zipErr := zip.NewReader(bytes.NewReader(excelData), int64(len(excelData))); zipErr != nil {
		t.Fatalf("updated excel payload is invalid zip: %v", zipErr)
	}
}

func TestUpdateChartDataScatterByRelID(t *testing.T) {
	e := newChartUpdateEditorFixture()
	e.parts.Set("ppt/charts/chart1.xml", []byte(scatterChartXML()))

	err := e.UpdateChartData(0, common.ChartSelector{RelID: "rIdChart"}, common.ChartDataUpdate{
		Series: []common.ChartSeriesData{
			{
				XValues: []float64{1, 2},
				YValues: []float64{10, 20},
			},
		},
	})
	if err != nil {
		t.Fatalf("UpdateChartData failed: %v", err)
	}

	chartData, _ := e.parts.Get("ppt/charts/chart1.xml")
	xml := string(chartData)
	if !strings.Contains(xml, "Sheet1!$A$2:$A$3") || !strings.Contains(xml, "Sheet1!$B$2:$B$3") {
		t.Fatalf("expected scatter formulas rewritten")
	}
}

func TestUpdateChartDataBubbleByRelID(t *testing.T) {
	e := newChartUpdateEditorFixture()
	e.parts.Set("ppt/charts/chart1.xml", []byte(bubbleChartXML()))

	err := e.UpdateChartData(0, common.ChartSelector{RelID: "rIdChart"}, common.ChartDataUpdate{
		Series: []common.ChartSeriesData{
			{
				XValues: []float64{1, 2},
				YValues: []float64{10, 20},
				Sizes:   []float64{5, 8},
			},
		},
	})
	if err != nil {
		t.Fatalf("UpdateChartData failed: %v", err)
	}

	chartData, _ := e.parts.Get("ppt/charts/chart1.xml")
	xml := string(chartData)
	if !strings.Contains(xml, "Sheet1!$C$2:$C$3") {
		t.Fatalf("expected bubble size formula rewritten")
	}
}

func TestUpdateChartDataPreservesNumericFormatCode(t *testing.T) {
	e := newChartUpdateEditorFixture()
	e.parts.Set("ppt/charts/chart1.xml", []byte(categoryChartXMLWithNumericFormat()))

	idx := 0
	err := e.UpdateChartData(0, common.ChartSelector{Index: &idx}, common.ChartDataUpdate{
		Categories: []string{"Q1", "Q2"},
		Series: []common.ChartSeriesData{
			{Values: []float64{0.125, 0.25}},
		},
	})
	if err != nil {
		t.Fatalf("UpdateChartData failed: %v", err)
	}

	chartData, _ := e.parts.Get("ppt/charts/chart1.xml")
	xml := string(chartData)
	if !strings.Contains(xml, "<c:formatCode>0.00%</c:formatCode>") {
		t.Fatalf("expected numeric format code to be preserved, got: %s", xml)
	}
}

func TestUpdateChartDataSelectorMismatchFails(t *testing.T) {
	e := newChartUpdateEditorFixture()
	e.parts.Set("ppt/charts/chart1.xml", []byte(categoryChartXML()))
	idx := 0
	err := e.UpdateChartData(0, common.ChartSelector{Index: &idx, RelID: "rIdOther"}, common.ChartDataUpdate{
		Categories: []string{"Q1"},
		Series:     []common.ChartSeriesData{{Values: []float64{1}}},
	})
	if err == nil {
		t.Fatalf("expected selector mismatch error")
	}
}

func TestUpdateChartDataMultiLevelCategoriesByIndex(t *testing.T) {
	e := newChartUpdateEditorFixture()
	e.parts.Set("ppt/charts/chart1.xml", []byte(categoryChartXMLWithMultiLevelCategories()))

	idx := 0
	err := e.UpdateChartData(0, common.ChartSelector{Index: &idx}, common.ChartDataUpdate{
		MultiLevelCategories: [][]string{
			{"2026", "2026"},
			{"Q1", "Q2"},
		},
		Series: []common.ChartSeriesData{
			{Values: []float64{100, 200}},
		},
	})
	if err != nil {
		t.Fatalf("UpdateChartData failed: %v", err)
	}

	chartData, _ := e.parts.Get("ppt/charts/chart1.xml")
	xml := string(chartData)
	if !strings.Contains(xml, "Sheet1!$A$2:$B$3") {
		t.Fatalf("expected multi-level category formula range, got: %s", xml)
	}
	if !strings.Contains(xml, "<c:multiLvlStrCache>") || !strings.Contains(xml, "<c:lvl>") {
		t.Fatalf("expected multi-level category cache nodes, got: %s", xml)
	}
	if !strings.Contains(xml, "<c:v>2026</c:v>") || !strings.Contains(xml, "<c:v>Q2</c:v>") {
		t.Fatalf("expected updated multi-level category values, got: %s", xml)
	}
}

func TestListSlideCharts(t *testing.T) {
	e := newChartUpdateEditorFixture()
	refs, err := e.ListSlideCharts(0)
	if err != nil {
		t.Fatalf("ListSlideCharts failed: %v", err)
	}
	if len(refs) != 1 || refs[0].RelID != "rIdChart" || refs[0].ChartPart != "ppt/charts/chart1.xml" {
		t.Fatalf("unexpected chart refs: %+v", refs)
	}
}

func TestUpdateChartDataSuccessWhenEmbeddingMissingButRelExists(t *testing.T) {
	e := newChartUpdateEditorFixture()
	e.parts.Set("ppt/charts/chart1.xml", []byte(categoryChartXML()))
	e.parts.Delete("ppt/embeddings/Microsoft_Excel_Worksheet1.xlsx")
	delete(e.chartEmbeddings, "ppt/charts/chart1.xml")
	e.parts.Set("ppt/charts/_rels/chart1.xml.rels", []byte(
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`+
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/package" `+
			`Target="../embeddings/Microsoft_Excel_Worksheet1.xlsx"/></Relationships>`,
	))

	idx := 0
	err := e.UpdateChartData(0, common.ChartSelector{Index: &idx}, common.ChartDataUpdate{
		Categories: []string{"Q1"},
		Series:     []common.ChartSeriesData{{Values: []float64{10}}},
	})
	if err != nil {
		t.Fatalf("UpdateChartData should succeed and create a new embedding: %v", err)
	}

	excelPath := e.chartEmbeddings["ppt/charts/chart1.xml"]
	if !e.parts.Has(excelPath) {
		t.Error("new excel part was not created")
	}
}

func newChartUpdateEditorFixture() *PresentationEditor {
	e := &PresentationEditor{
		parts: NewPartStore(),
		slides: []common.EditorSlideRef{
			{Part: "ppt/slides/slide1.xml"},
		},
		chartEmbeddings: map[string]string{
			"ppt/charts/chart1.xml": "ppt/embeddings/Microsoft_Excel_Worksheet1.xlsx",
		},
	}
	e.parts.Set(
		"ppt/slides/slide1.xml",
		[]byte(
			`<p:sld><p:spTree><p:graphicFrame><a:graphic><a:graphicData>`+
				`<c:chart r:id="rIdChart"/></a:graphicData></a:graphic></p:graphicFrame></p:spTree></p:sld>`,
		),
	)
	e.parts.Set(
		"ppt/slides/_rels/slide1.xml.rels",
		[]byte(
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`+
				`<Relationship Id="rIdChart" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart" `+
				`Target="../charts/chart1.xml"/></Relationships>`,
		),
	)
	e.parts.Set(
		"ppt/charts/_rels/chart1.xml.rels",
		[]byte(
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`+
				`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/package" `+
				`Target="../embeddings/Microsoft_Excel_Worksheet1.xlsx"/></Relationships>`,
		),
	)
	e.parts.Set("ppt/embeddings/Microsoft_Excel_Worksheet1.xlsx", []byte("old"))
	return e
}

func categoryChartXML() string {
	return `<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
		`<c:plotArea><c:barChart><c:ser><c:idx val="0"/><c:order val="0"/><c:tx><c:v>Series 1</c:v></c:tx>` +
		`<c:spPr><a:solidFill/></c:spPr><c:cat><c:strRef><c:f>Sheet1!$A$2:$A$2</c:f>` +
		`<c:strCache><c:ptCount val="1"/><c:pt idx="0"><c:v>Old</c:v></c:pt></c:strCache></c:strRef></c:cat>` +
		`<c:val><c:numRef><c:f>Sheet1!$B$2:$B$2</c:f><c:numCache><c:ptCount val="1"/><c:pt idx="0"><c:v>1</c:v></c:pt>` +
		`</c:numCache></c:numRef></c:val></c:ser></c:barChart></c:plotArea></c:chartSpace>`
}

func categoryChartXMLWithNumericFormat() string {
	return `<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart">` +
		`<c:plotArea><c:barChart><c:ser>` +
		`<c:cat><c:strRef><c:f>Sheet1!$A$2:$A$2</c:f><c:strCache><c:ptCount val="1"/>` +
		`<c:pt idx="0"><c:v>Old</c:v></c:pt></c:strCache></c:strRef></c:cat>` +
		`<c:val><c:numRef><c:f>Sheet1!$B$2:$B$2</c:f><c:numCache>` +
		`<c:formatCode>0.00%</c:formatCode><c:ptCount val="1"/><c:pt idx="0"><c:v>0.10</c:v></c:pt>` +
		`</c:numCache></c:numRef></c:val>` +
		`</c:ser></c:barChart></c:plotArea></c:chartSpace>`
}

func categoryChartXMLWithMultiLevelCategories() string {
	return `<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart">` +
		`<c:plotArea><c:barChart><c:ser><c:idx val="0"/><c:order val="0"/>` +
		`<c:cat><c:multiLvlStrRef><c:f>Sheet1!$A$2:$B$2</c:f><c:multiLvlStrCache><c:ptCount val="1"/>` +
		`<c:lvl><c:pt idx="0"><c:v>OLD-YEAR</c:v></c:pt></c:lvl>` +
		`<c:lvl><c:pt idx="0"><c:v>OLD-Q</c:v></c:pt></c:lvl>` +
		`</c:multiLvlStrCache></c:multiLvlStrRef></c:cat>` +
		`<c:val><c:numRef><c:f>Sheet1!$C$2:$C$2</c:f><c:numCache><c:ptCount val="1"/>` +
		`<c:pt idx="0"><c:v>1</c:v></c:pt></c:numCache></c:numRef></c:val>` +
		`</c:ser></c:barChart></c:plotArea></c:chartSpace>`
}

func scatterChartXML() string {
	return `<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart"><c:plotArea><c:scatterChart><c:ser>` +
		`<c:xVal><c:numRef><c:f>Sheet1!$A$2:$A$2</c:f><c:numCache><c:ptCount val="1"/><c:pt idx="0"><c:v>1</c:v></c:pt>` +
		`</c:numCache></c:numRef></c:xVal><c:yVal><c:numRef><c:f>Sheet1!$B$2:$B$2</c:f>` +
		`<c:numCache><c:ptCount val="1"/><c:pt idx="0"><c:v>2</c:v></c:pt></c:numCache></c:numRef></c:yVal></c:ser>` +
		`</c:scatterChart></c:plotArea></c:chartSpace>`
}

func bubbleChartXML() string {
	return `<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart"><c:plotArea>` +
		`<c:bubbleChart><c:ser>` +
		`<c:xVal><c:numRef><c:f>Sheet1!$A$2:$A$2</c:f><c:numCache><c:ptCount val="1"/><c:pt idx="0"><c:v>1</c:v></c:pt>` +
		`</c:numCache></c:numRef></c:xVal><c:yVal><c:numRef><c:f>Sheet1!$B$2:$B$2</c:f>` +
		`<c:numCache><c:ptCount val="1"/><c:pt idx="0"><c:v>2</c:v></c:pt></c:numCache></c:numRef></c:yVal>` +
		`<c:bubbleSize><c:numRef><c:f>Sheet1!$C$2:$C$2</c:f><c:numCache><c:ptCount val="1"/><c:pt idx="0"><c:v>3</c:v></c:pt>` +
		`</c:numCache></c:numRef></c:bubbleSize></c:ser></c:bubbleChart></c:plotArea></c:chartSpace>`
}
