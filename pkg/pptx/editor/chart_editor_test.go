package editor

import (
	"archive/zip"
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodchart "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/chart"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
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
	if !strings.Contains(
		string(slideRelsData),
		"http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart",
	) {
		t.Error("missing slide->chart relationship")
	}

	// Chart -> Excel
	chartRelsData, _ := editor.parts.Get("ppt/charts/_rels/chart1.xml.rels")
	if !strings.Contains(
		string(chartRelsData),
		"http://schemas.openxmlformats.org/officeDocument/2006/relationships/package",
	) {
		t.Error("missing chart->excel relationship")
	}

	// d. Check Inventory
	if editor.chartEmbeddings["ppt/charts/chart1.xml"] != "ppt/embeddings/Microsoft_Excel_Worksheet1.xlsx" {
		t.Error("inventory not updated correctly")
	}
}

func TestAddChartDeduplicatesIdenticalExcelEmbeddings(t *testing.T) {
	editor := &PresentationEditor{
		parts: NewPartStore(),
		slides: []common.EditorSlideRef{
			{Part: "ppt/slides/slide1.xml"},
			{Part: "ppt/slides/slide2.xml"},
		},
		chartEmbeddings: make(map[string]string),
		nextChartNum:    1,
		nextExcelNum:    1,
		nextRelIDNum:    1,
	}
	editor.parts.Set("ppt/slides/slide1.xml", []byte(`<p:sld><p:spTree></p:spTree></p:sld>`))
	editor.parts.Set("ppt/slides/slide2.xml", []byte(`<p:sld><p:spTree></p:spTree></p:sld>`))
	editor.parts.Set("[Content_Types].xml", []byte(`<Types></Types>`))

	chartDef := charts.NewBarChart(
		[]string{"A", "B"},
		[]float64{10, 20},
	).WithTitle("Same Data")

	if err := editor.AddChart(0, chartDef); err != nil {
		t.Fatalf("add chart 1 failed: %v", err)
	}
	if err := editor.AddChart(1, chartDef); err != nil {
		t.Fatalf("add chart 2 failed: %v", err)
	}

	embeddingParts := editor.parts.KeysWithPrefix("ppt/embeddings/")
	if len(embeddingParts) != 1 {
		t.Fatalf("expected exactly one deduplicated embedding part, got %d: %v", len(embeddingParts), embeddingParts)
	}
}

func TestMergeFromFilePreservesChartEmbeddingChain(t *testing.T) {
	sourcePath := writeDeckFixture(t, "source-chart-base.pptx", []elements.SlideContent{
		elements.NewSlide("Source Chart"),
	})
	sourceEditor, err := OpenPresentationEditor(sourcePath)
	if err != nil {
		t.Fatalf("open source editor: %v", err)
	}
	chartDef := charts.NewBarChart(
		[]string{"Q1", "Q2"},
		[]float64{100, 120},
	).WithTitle("Revenue")
	if addErr := sourceEditor.AddChart(0, chartDef); addErr != nil {
		t.Fatalf("add source chart: %v", addErr)
	}
	sourceWithChart := filepath.Join(t.TempDir(), "source-with-chart.pptx")
	if saveSourceErr := sourceEditor.Save(sourceWithChart); saveSourceErr != nil {
		t.Fatalf("save source with chart: %v", saveSourceErr)
	}
	_ = sourceEditor.Close()

	destPath := writeDeckFixture(t, "dest-base.pptx", []elements.SlideContent{
		elements.NewSlide("Dest Slide"),
	})
	destEditor, err := OpenPresentationEditor(destPath)
	if err != nil {
		t.Fatalf("open dest editor: %v", err)
	}
	defer func() { _ = destEditor.Close() }()
	if mergeErr := destEditor.MergeFromFile(sourceWithChart); mergeErr != nil {
		t.Fatalf("merge from file failed: %v", mergeErr)
	}
	outPath := filepath.Join(t.TempDir(), "merged-chart.pptx")
	if saveDestErr := destEditor.Save(outPath); saveDestErr != nil {
		t.Fatalf("save merged deck: %v", saveDestErr)
	}

	merged, err := OpenPresentationEditor(outPath)
	if err != nil {
		t.Fatalf("reopen merged deck: %v", err)
	}
	defer func() { _ = merged.Close() }()
	if merged.SlideCount() != 2 {
		t.Fatalf("expected 2 slides after merge, got %d", merged.SlideCount())
	}

	chartParts := merged.parts.KeysWithPrefix("ppt/charts/chart")
	if len(chartParts) == 0 {
		t.Fatalf("expected merged chart parts")
	}
	embeddingParts := merged.parts.KeysWithPrefix("ppt/embeddings/")
	if len(embeddingParts) == 0 {
		t.Fatalf("expected merged embedding parts")
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
			<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart"
				Target="../charts/chart1.xml"/>
		</Relationships>
	`))
	editor.parts.Set(
		"ppt/charts/chart1.xml",
		[]byte(
			`<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart">`+
				`<c:plotArea><c:barChart><c:ser>`+
				`<c:cat><c:strRef><c:f>Sheet1!$A$2:$A$2</c:f><c:strCache><c:ptCount val="1"/><c:pt idx="0"><c:v>A</c:v></c:pt>`+
				`</c:strCache></c:strRef></c:cat><c:val><c:numRef><c:f>Sheet1!$B$2:$B$2</c:f><c:numCache><c:ptCount val="1"/>`+
				`<c:pt idx="0"><c:v>1</c:v></c:pt></c:numCache></c:numRef></c:val></c:ser></c:barChart></c:plotArea></c:chartSpace>`,
		),
	)
	editor.parts.Set(
		"ppt/charts/_rels/chart1.xml.rels",
		[]byte(
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`+
				`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/package" `+
				`Target="../embeddings/Microsoft_Excel_Worksheet1.xlsx"/></Relationships>`,
		),
	)

	// Mock Excel Part
	editor.parts.Set("ppt/embeddings/Microsoft_Excel_Worksheet1.xlsx", []byte("old data"))

	// 2. Execute ReplaceChartData
	newCats := []string{"X", "Y"}
	newVals := []float64{100, 200}

	err := editor.ReplaceChartData(0, 0, newCats, newVals)
	if err != nil {
		t.Fatalf("ReplaceChartData failed: %v", err)
	}

	// 3. Verify Excel content changed (it should now be a valid zip)
	excelPath := editor.chartEmbeddings["ppt/charts/chart1.xml"]
	excelData, _ := editor.parts.Get(excelPath)

	// Check if it's a valid zip
	_, err = zip.NewReader(bytes.NewReader(excelData), int64(len(excelData)))
	if err != nil {
		t.Errorf("updated excel data is not a valid zip: %v", err)
	}
}

func TestGenerateExcelForChart(t *testing.T) {
	categories := []string{"Cat 1", "Cat 2"}
	values := []float64{10.5, 20.0}

	data, err := editormodchart.GenerateExcelForChart(categories, values)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it's a valid zip
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("invalid zip archive: %v", err)
	}

	requiredFiles := []string{
		"[Content_Types].xml",
		"_rels/.rels",
		"xl/workbook.xml",
		"xl/_rels/workbook.xml.rels",
		"xl/styles.xml",
		"xl/worksheets/sheet1.xml",
	}

	for _, req := range requiredFiles {
		found := false
		for _, f := range zr.File {
			if f.Name == req {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing required file in xlsx: %s", req)
		}
	}
}
