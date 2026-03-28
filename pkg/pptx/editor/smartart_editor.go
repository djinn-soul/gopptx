package editor

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

const (
	relTypeDiagramData       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/diagramData"
	relTypeDiagramLayout     = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/diagramLayout"
	relTypeDiagramQuickStyle = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/diagramQuickStyle"
	relTypeDiagramColors     = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/diagramColors"
	relTypeDiagramDrawing    = "http://schemas.microsoft.com/office/2007/relationships/diagramDrawing"

	contentTypeDiagramData    = "application/vnd.openxmlformats-officedocument.drawingml.diagramData+xml"
	contentTypeDiagramLayout  = "application/vnd.openxmlformats-officedocument.drawingml.diagramLayout+xml"
	contentTypeDiagramColors  = "application/vnd.openxmlformats-officedocument.drawingml.diagramColors+xml"
	contentTypeDiagramStyle   = "application/vnd.openxmlformats-officedocument.drawingml.diagramStyle+xml"
	contentTypeDiagramDrawing = "application/vnd.ms-office.drawingml.diagramDrawing+xml"

	// Relative offsets for the 5 diagram relationship IDs (dm=base, lo=base+1, ...).
	diagRelOffsetQuickStyle = 2
	diagRelOffsetColors     = 3
	diagRelOffsetDrawing    = 4
)

// AddSmartArt inserts a SmartArt diagram onto an existing slide.
// Returns the shape ID of the graphic frame inserted into the slide.
func (e *PresentationEditor) AddSmartArt(slideIndex int, sa smartart.SmartArt) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]

	num := e.nextDiagramNum
	e.nextDiagramNum++

	spec := sa.ToSpec()
	category := smartArtCategoryFromURI(spec.LayoutURI)

	// Generate the 5 diagram parts.
	dataPath := fmt.Sprintf("ppt/diagrams/data%d.xml", num)
	layoutPath := fmt.Sprintf("ppt/diagrams/layout%d.xml", num)
	colorsPath := fmt.Sprintf("ppt/diagrams/colors%d.xml", num)
	stylePath := fmt.Sprintf("ppt/diagrams/quickStyle%d.xml", num)
	drawingPath := fmt.Sprintf("ppt/diagrams/drawing%d.xml", num)
	e.parts.Set(dataPath, []byte(pptxxml.SmartArtDataXML(spec)))
	e.parts.Set(layoutPath, []byte(pptxxml.SmartArtLayoutXML(spec.LayoutURI, category)))
	e.parts.Set(colorsPath, []byte(pptxxml.SmartArtColorsXML(spec.ColorStyleID)))
	e.parts.Set(stylePath, []byte(pptxxml.SmartArtStyleXML(spec.QuickStyleID)))
	e.parts.Set(drawingPath, []byte(pptxxml.SmartArtDrawingXML(spec)))

	// Register content types.
	e.addContentTypeOverride(dataPath, contentTypeDiagramData)
	e.addContentTypeOverride(layoutPath, contentTypeDiagramLayout)
	e.addContentTypeOverride(colorsPath, contentTypeDiagramColors)
	e.addContentTypeOverride(stylePath, contentTypeDiagramStyle)
	e.addContentTypeOverride(drawingPath, contentTypeDiagramDrawing)

	// Allocate 5 consecutive relationship IDs: dm, lo, qs, cs, dr.
	baseRelID, relErr := e.nextSlideRelID(slideRef.Part)
	if relErr != nil {
		return 0, fmt.Errorf("alloc diagram rel ids: %w", relErr)
	}
	baseNum, _ := common.ParseRelationshipNumber(baseRelID)

	dmID := fmt.Sprintf("rId%d", baseNum)
	loID := fmt.Sprintf("rId%d", baseNum+1)
	qsID := fmt.Sprintf("rId%d", baseNum+diagRelOffsetQuickStyle)
	csID := fmt.Sprintf("rId%d", baseNum+diagRelOffsetColors)
	drID := fmt.Sprintf("rId%d", baseNum+diagRelOffsetDrawing)

	diagRels := []common.EditorRelationship{
		{ID: dmID, Type: relTypeDiagramData, Target: fmt.Sprintf("../diagrams/data%d.xml", num)},
		{ID: loID, Type: relTypeDiagramLayout, Target: fmt.Sprintf("../diagrams/layout%d.xml", num)},
		{ID: qsID, Type: relTypeDiagramQuickStyle, Target: fmt.Sprintf("../diagrams/quickStyle%d.xml", num)},
		{ID: csID, Type: relTypeDiagramColors, Target: fmt.Sprintf("../diagrams/colors%d.xml", num)},
		{ID: drID, Type: relTypeDiagramDrawing, Target: fmt.Sprintf("../diagrams/drawing%d.xml", num)},
	}
	if addErr := e.addSlideRelationships(slideRef.Part, diagRels); addErr != nil {
		return 0, fmt.Errorf("add diagram slide rels: %w", addErr)
	}

	// Build and insert the graphic frame into the slide spTree.
	shapeID := e.nextShapeID(slideRef.Part)
	frame := pptxxml.SmartArtFrame{
		X: spec.X, Y: spec.Y, CX: spec.CX, CY: spec.CY,
		DataRelID:    dmID,
		LayoutRelID:  loID,
		StyleRelID:   qsID,
		ColorRelID:   csID,
		AltText:      spec.AltText,
		IsDecorative: spec.IsDecorative,
	}
	gfxXML := buildSmartArtGraphicFrameXML(frame, shapeID)
	if appendErr := e.appendShapeToSlide(slideRef.Part, gfxXML); appendErr != nil {
		return 0, fmt.Errorf("append smartart frame: %w", appendErr)
	}

	return shapeID, nil
}

// buildSmartArtGraphicFrameXML renders the p:graphicFrame XML for a SmartArt.
func buildSmartArtGraphicFrameXML(frame pptxxml.SmartArtFrame, shapeID int) string {
	altAttr := ""
	if frame.AltText != "" {
		altAttr = fmt.Sprintf(` descr=%q`, frame.AltText)
	}
	return fmt.Sprintf(`
<p:graphicFrame xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" `+
		`xmlns:dgm="http://schemas.openxmlformats.org/drawingml/2006/diagram" `+
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" `+
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
<p:nvGraphicFramePr>
<p:cNvPr id="%d" name="Diagram %d"%s/>
<p:cNvGraphicFramePr/>
<p:nvPr/>
</p:nvGraphicFramePr>
<p:xfrm>
<a:off x="%d" y="%d"/>
<a:ext cx="%d" cy="%d"/>
</p:xfrm>
<a:graphic>
<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/diagram">
<dgm:relIds r:dm=%q r:lo=%q r:qs=%q r:cs=%q/>
</a:graphicData>
</a:graphic>
</p:graphicFrame>`,
		shapeID, shapeID, altAttr,
		frame.X, frame.Y, frame.CX, frame.CY,
		frame.DataRelID, frame.LayoutRelID, frame.StyleRelID, frame.ColorRelID,
	)
}

// addSlideRelationships appends multiple relationships to a slide's .rels file in one pass.
func (e *PresentationEditor) addSlideRelationships(slidePart string, rels []common.EditorRelationship) error {
	relsPath := common.RelsPathFor(slidePart)
	var existing []common.EditorRelationship
	if data, ok := e.parts.Get(relsPath); ok {
		parsed, err := parseRelationshipsXML(data)
		if err != nil {
			return fmt.Errorf("parse %s: %w", relsPath, err)
		}
		existing = parsed
	}
	existing = append(existing, rels...)
	return e.writeRelationships(relsPath, existing)
}

// smartArtCategoryFromURI infers SmartArt category from the layout URI.
func smartArtCategoryFromURI(uri string) string {
	switch {
	case strings.Contains(uri, "process"):
		return "process"
	case strings.Contains(uri, "cycle"):
		return "cycle"
	case strings.Contains(uri, "hierarchy"), strings.Contains(uri, "orgChart"):
		return "hierarchy"
	case strings.Contains(uri, "venn"), strings.Contains(uri, "radial"), strings.Contains(uri, "target"):
		return "relationship"
	case strings.Contains(uri, "matrix"):
		return "matrix"
	case strings.Contains(uri, "pyramid"):
		return "pyramid"
	case strings.Contains(uri, "picture"):
		return "picture"
	default:
		return "list"
	}
}
