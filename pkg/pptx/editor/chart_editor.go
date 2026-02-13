package editor

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"strconv"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// AddChart adds a new chart to a specific slide.
func (e *PresentationEditor) AddChart(slideIndex int, chartDef charts.ChartDefinition) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]

	// 1. Generate Excel data part
	excelData, err := generateExcelForChart(chartDef.GetCategories(), chartDef.GetValues())
	if err != nil {
		return fmt.Errorf("generate excel: %w", err)
	}

	// 2. Allocate new part paths
	// Chart part
	chartNum := e.nextChartNum
	e.nextChartNum++
	chartPartPath := fmt.Sprintf("ppt/charts/chart%d.xml", chartNum)

	// Excel part
	excelNum := e.nextExcelNum
	e.nextExcelNum++
	excelPartPath := fmt.Sprintf("ppt/embeddings/Microsoft_Excel_Worksheet%d.xlsx", excelNum)

	// 3. Register parts
	e.parts.Set(chartPartPath, nil) // Will populate content later
	e.parts.Set(excelPartPath, excelData)

	// Update [Content_Types].xml overrides
	e.addContentTypeOverride(chartPartPath, "application/vnd.openxmlformats-officedocument.drawingml.chart+xml")
	e.addContentTypeOverride(excelPartPath, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// 4. Create relationships
	// Slide -> Chart
	slideRelID := fmt.Sprintf("rId%d", e.nextRelIDNum)
	e.nextRelIDNum++
	if err := e.addSlideRelationship(slideRef.Part, slideRelID, common.RelTypeChart, "../charts/"+path.Base(chartPartPath)); err != nil {
		return fmt.Errorf("add slide rel: %w", err)
	}

	// Chart -> Excel
	chartRelID := "rId1"
	if err := e.addRelationship(chartPartPath, chartRelID, common.RelTypePackage, "../embeddings/"+path.Base(excelPartPath)); err != nil {
		return fmt.Errorf("add chart rel: %w", err)
	}

	// 5. Generate Chart XML
	chartSpec := chartDef.ToChartSpec()
	// Update external data ID to match our new rId
	chartSpec.ExternalDataID = chartRelID

	chartXMLBytes := pptxxml.RenderChart(chartSpec)
	e.parts.Set(chartPartPath, chartXMLBytes)

	// 6. Update Slide XML to include the graphic frame
	// We need to find non-conflicting ID for the graphic frame
	shapeID := e.nextShapeID(slideRef.Part)

	// Create GraphicFrame XML
	gfxFrame := e.createChartGraphicFrameXML(shapeID, chartSpec.Title, slideRelID, chartSpec.X, chartSpec.Y, chartSpec.CX, chartSpec.CY)

	if err := e.appendShapeToSlide(slideRef.Part, gfxFrame); err != nil {
		return fmt.Errorf("append chart shape: %w", err)
	}

	// Update inventory
	e.chartEmbeddings[chartPartPath] = excelPartPath

	return nil
}

// ReplaceChartData updates the data source and cached values of an existing chart.
// chartIndex is the 0-based index of the chart on the slide (order of appearance).
func (e *PresentationEditor) ReplaceChartData(slideIndex int, chartIndex int, categories []string, values []float64) error {
	idx := chartIndex
	return e.UpdateChartData(slideIndex, common.ChartSelector{Index: &idx}, common.ChartDataUpdate{
		Categories: categories,
		Series: []common.ChartSeriesData{
			{Values: values},
		},
	})
}

func extractChartRelIDs(content []byte) []string {
	re := regexp.MustCompile(`<(?:c:chart|chart)[^>]*r:id="([^"]+)"`)
	matches := re.FindAllSubmatch(content, -1)
	ids := make([]string, 0, len(matches))
	for _, m := range matches {
		ids = append(ids, string(m[1]))
	}
	return ids
}

func (e *PresentationEditor) addSlideRelationship(slidePart, id, relType, target string) error {
	relsPath := common.RelsPathFor(slidePart)

	// Load or create rels
	var rels []common.EditorRelationship
	if data, ok := e.parts.Get(relsPath); ok {
		var err error
		rels, err = parseRelationshipsXML(data)
		if err != nil {
			return err
		}
	}

	rels = append(rels, common.EditorRelationship{
		ID:     id,
		Type:   relType,
		Target: target,
	})

	return e.writeRelationships(relsPath, rels)
}

func (e *PresentationEditor) addRelationship(partPath, id, relType, target string) error {
	relsPath := common.RelsPathFor(partPath)
	var rels []common.EditorRelationship
	if data, ok := e.parts.Get(relsPath); ok {
		var err error
		rels, err = parseRelationshipsXML(data)
		if err != nil {
			return err
		}
	}
	rels = append(rels, common.EditorRelationship{
		ID:     id,
		Type:   relType,
		Target: target,
	})
	return e.writeRelationships(relsPath, rels)
}

func (e *PresentationEditor) writeRelationships(path string, rels []common.EditorRelationship) error {
	output := `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`
	for _, r := range rels {
		output += fmt.Sprintf(`<Relationship Id="%s" Type="%s" Target="%s"`, r.ID, r.Type, r.Target)
		if r.TargetMode != "" {
			output += fmt.Sprintf(` TargetMode="%s"`, r.TargetMode)
		}
		output += `/>`
	}
	output += `</Relationships>`
	e.parts.Set(path, []byte(output))
	return nil
}

func (e *PresentationEditor) addContentTypeOverride(partName, contentType string) {
	ctPath := "[Content_Types].xml"
	data, ok := e.parts.Get(ctPath)
	if !ok {
		return
	}

	partNameRooted := "/" + partName
	if bytes.Contains(data, []byte(fmt.Sprintf(`PartName="%s"`, partNameRooted))) {
		return
	}

	override := fmt.Sprintf(`<Override PartName="%s" ContentType="%s"/>`, partNameRooted, contentType)
	replaced := bytes.Replace(data, []byte("</Types>"), []byte(override+"</Types>"), 1)
	e.parts.Set(ctPath, replaced)
}

func (e *PresentationEditor) nextShapeID(slidePart string) int {
	data, _ := e.parts.Get(slidePart)
	re := regexp.MustCompile(`id="(\d+)"`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	maxID := 0
	for _, m := range matches {
		id, _ := strconv.Atoi(m[1])
		if id > maxID {
			maxID = id
		}
	}
	return maxID + 1
}

func (e *PresentationEditor) createChartGraphicFrameXML(id int, name, rId string, x, y, cx, cy int64) string {
	return fmt.Sprintf(`
	<p:graphicFrame>
		<p:nvGraphicFramePr>
			<p:cNvPr id="%d" name="%s"/>
			<p:cNvGraphicFramePr/>
			<p:nvPr/>
		</p:nvGraphicFramePr>
		<p:xf>
			<a:off x="%d" y="%d"/>
			<a:ext cx="%d" cy="%d"/>
		</p:xf>
		<a:graphic>
			<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/chart">
				<c:chart xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" r:id="%s"/>
			</a:graphicData>
		</a:graphic>
	</p:graphicFrame>`, id, name, x, y, cx, cy, rId)
}

func (e *PresentationEditor) appendShapeToSlide(slidePart, shapeXML string) error {
	data, ok := e.parts.Get(slidePart)
	if !ok {
		return fmt.Errorf("part not found")
	}

	if !bytes.Contains(data, []byte("</p:spTree>")) {
		return fmt.Errorf("invalid slide xml: missing spTree end")
	}

	replaced := bytes.Replace(data, []byte("</p:spTree>"), []byte(shapeXML+"</p:spTree>"), 1)
	e.parts.Set(slidePart, replaced)
	return nil
}

func resolveRelationshipTarget(basePart, target string) string {
	baseDir := path.Dir(basePart)
	joined := path.Join(baseDir, target)
	return path.Clean(joined)
}
