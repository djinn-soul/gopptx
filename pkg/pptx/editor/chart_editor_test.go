package editor

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestAddChart(t *testing.T) {
	// 1. Setup minimal editor with one slide
	editor := &PresentationEditor{
		parts: NewPartStore(),
		slides: []common.EditorSlideRef{
			{Part: "ppt/slides/slide1.xml"},
		},
		chartEmbeddings: make(map[string]string),
		nextChartNum:    1,
		nextExcelNum:    1,
		nextRelIDNum:    1,
	}

	// Mock slide content
	editor.parts.Set("ppt/slides/slide1.xml", []byte(`<p:sld><p:spTree></p:spTree></p:sld>`))
	editor.parts.Set("[Content_Types].xml", []byte(`<Types></Types>`))

	// 2. Define Chart
	chartDef := charts.NewBarChart(
		[]string{"A", "B"},
		[]float64{10, 20},
	).WithTitle("My Chart")

	// 3. Execute AddChart
	err := editor.AddChart(0, chartDef)
	if err != nil {
		t.Fatalf("AddChart failed: %v", err)
	}

	// 4. Verify Side Effects
	// a. Check Chart Part exists
	if !editor.parts.Has("ppt/charts/chart1.xml") {
		t.Error("missing chart part")
	}

	// b. Check Excel Part exists
	if !editor.parts.Has("ppt/embeddings/Microsoft_Excel_Worksheet1.xlsx") {
		t.Error("missing excel part")
	}

	// c. Check Relationships
	// Slide -> Chart
	slideRelsData, _ := editor.parts.Get("ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(string(slideRelsData), "http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart") {
		t.Error("missing slide->chart relationship")
	}

	// Chart -> Excel
	chartRelsData, _ := editor.parts.Get("ppt/charts/_rels/chart1.xml.rels")
	if !strings.Contains(string(chartRelsData), "http://schemas.openxmlformats.org/officeDocument/2006/relationships/package") {
		t.Error("missing chart->excel relationship")
	}

	// d. Check Inventory
	if editor.chartEmbeddings["ppt/charts/chart1.xml"] != "ppt/embeddings/Microsoft_Excel_Worksheet1.xlsx" {
		t.Error("inventory not updated correctly")
	}
}

func TestReplaceChartData(t *testing.T) {
	// 1. Setup editor with existing chart
	editor := &PresentationEditor{
		parts: NewPartStore(),
		slides: []common.EditorSlideRef{
			{Part: "ppt/slides/slide1.xml"},
		},
		chartEmbeddings: map[string]string{
			"ppt/charts/chart1.xml": "ppt/embeddings/Microsoft_Excel_Worksheet1.xlsx",
		},
	}

	// Mock Slide with GraphicFrame pointing to rId1
	editor.parts.Set("ppt/slides/slide1.xml", []byte(`
		<p:sld>
			<p:spTree>
				<p:graphicFrame>
					<a:graphic>
						<a:graphicData>
							<c:chart r:id="rId1"/>
						</a:graphicData>
					</a:graphic>
				</p:graphicFrame>
			</p:spTree>
		</p:sld>
	`))

	// Mock Slide Rels
	editor.parts.Set("ppt/slides/_rels/slide1.xml.rels", []byte(`
		<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
			<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart" Target="../charts/chart1.xml"/>
		</Relationships>
	`))
	editor.parts.Set("ppt/charts/chart1.xml", []byte(`<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart"><c:plotArea><c:barChart><c:ser><c:cat><c:strRef><c:f>Sheet1!$A$2:$A$2</c:f><c:strCache><c:ptCount val="1"/><c:pt idx="0"><c:v>A</c:v></c:pt></c:strCache></c:strRef></c:cat><c:val><c:numRef><c:f>Sheet1!$B$2:$B$2</c:f><c:numCache><c:ptCount val="1"/><c:pt idx="0"><c:v>1</c:v></c:pt></c:numCache></c:numRef></c:val></c:ser></c:barChart></c:plotArea></c:chartSpace>`))

	// Mock Excel Part
	editor.parts.Set("ppt/embeddings/Microsoft_Excel_Worksheet1.xlsx", []byte("old data"))

	// 2. Execute ReplaceChartData
	newCats := []string{"X", "Y"}
	newVals := []float64{100, 200}

	err := editor.ReplaceChartData(0, 0, newCats, newVals)
	if err != nil {
		t.Fatalf("ReplaceChartData failed: %v", err)
	}

	// 3. Verify Excel content changed (it should now be a valid zip, not "old data")
	excelData, _ := editor.parts.Get("ppt/embeddings/Microsoft_Excel_Worksheet1.xlsx")
	if string(excelData) == "old data" {
		t.Error("excel data was not updated")
	}

	// Check if it's a valid zip
	_, err = zip.NewReader(bytes.NewReader(excelData), int64(len(excelData)))
	if err != nil {
		t.Errorf("updated excel data is not a valid zip: %v", err)
	}
}
