package export

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	placeholderTitle       = "title"
	placeholderCtrTitle    = "ctrtitle"
	placeholderBody        = "body"
	placeholderSubtitle    = "subtitle"
	placeholderObject      = "obj"
	placeholderContentName = "content"
	minTableBorderWidthPt  = 0.25
)

// SlidesFromPPTX reads an existing PPTX file and extracts slide content
// (title, bullets, shapes, embedded images) for the native PDF/HTML export pipeline.
func SlidesFromPPTX(pptxPath string) (string, []elements.SlideContent, error) {
	ed, err := editor.OpenPresentationEditor(pptxPath)
	if err != nil {
		return "", nil, fmt.Errorf("open PPTX: %w", err)
	}
	defer ed.Close()

	meta := ed.Metadata()
	presTitle := ""
	if meta != nil {
		presTitle = meta.Title
	}

	// Extract embedded images per slide via direct PPTX zip parsing.
	slideImages, err := extractSlideImages(pptxPath)
	if err != nil {
		// Best-effort image extraction; continue without images when parsing fails.
		slideImages = nil
	}
	slideCharts, err := extractSlideCharts(pptxPath)
	if err != nil {
		// Best-effort chart extraction; continue without charts when parsing fails.
		slideCharts = nil
	}
	slideSmartArt, err := extractSlideSmartArt(pptxPath)
	if err != nil {
		// Best-effort SmartArt extraction; continue without semantic diagrams when parsing fails.
		slideSmartArt = nil
	}

	slideMeta := ed.Slides()
	slideContents := make([]elements.SlideContent, 0, len(slideMeta))

	for _, sm := range slideMeta {
		slideContents = append(slideContents, extractSlideContent(ed, sm, slideImages, slideCharts, slideSmartArt))
	}

	if presTitle == "" && len(slideContents) > 0 {
		presTitle = slideContents[0].Title
	}

	return presTitle, slideContents, nil
}

//nolint:gocognit,funlen // Slide extraction maps multiple editor shape categories and image attachments in one pass.
func extractSlideContent(
	ed *editor.PresentationEditor,
	sm editorcommon.SlideMetadata,
	slideImages [][]SlideImage,
	slideCharts [][]parsedChart,
	slideSmartArt [][]parsedSmartArt,
) elements.SlideContent {
	editorShapes, err := ed.GetShapes(sm.Index)
	if err != nil {
		editorShapes = nil
	}

	sc := elements.SlideContent{
		Title: sm.Title,
	}
	shapeIndexByID := make(map[int]int)
	connectorShapes := make([]editorcommon.Shape, 0)
	chartShapeIDs := chartShapeIDSet(slideCharts, sm.Index)
	smartArtShapeIDs := smartArtShapeIDSet(slideSmartArt, sm.Index)

	for _, es := range editorShapes {
		lowerType := strings.ToLower(es.Type)
		lowerName := strings.ToLower(strings.TrimSpace(es.Name))
		// For shapes that carry a placeholder type but no geometry preset (e.g.
		// python-pptx generated shapes whose <p:spPr/> is empty), use the
		// placeholder type as the effective classification type so that title and
		// body placeholders are still handled correctly.
		lowerPhType := strings.ToLower(es.PlaceholderType)
		if lowerPhType != "" && (lowerType == "" || lowerType == "sp") {
			lowerType = lowerPhType
		}
		if isEditorConnector(es) {
			connectorShapes = append(connectorShapes, es)
			continue
		}

		switch lowerType {
		case "pic":
			continue
		case "graphicframe":
			if _, ok := chartShapeIDs[es.ID]; ok {
				continue
			}
			if _, ok := smartArtShapeIDs[es.ID]; ok {
				continue
			}
			if tbl := extractTableContent(ed, sm.Index, es); tbl != nil {
				sc.Table = tbl
				continue
			}
			appendReaderShape(&sc, shapeIndexByID, es)
		case placeholderTitle, placeholderCtrTitle:
			if sc.Title == "" && es.Text != "" {
				sc.Title = es.Text
			}
			applyTitleBounds(&sc, es)
			applyTitleSizeFromRuns(&sc, es)
			applyTitleAlignFromShape(&sc, es)
		case placeholderBody, placeholderSubtitle, placeholderObject:
			if consumeBodyPlaceholderAsBullets(&sc, es) {
				applyContentBounds(&sc, es)
				continue
			}
			appendReaderShape(&sc, shapeIndexByID, es)
		default:
			switch {
			case isTitlePlaceholder(lowerType, lowerName):
				if sc.Title == "" && es.Text != "" {
					sc.Title = es.Text
				}
				applyTitleBounds(&sc, es)
				applyTitleSizeFromRuns(&sc, es)
				applyTitleAlignFromShape(&sc, es)
			case isBodyPlaceholder(lowerType, lowerName):
				if consumeBodyPlaceholderAsBullets(&sc, es) {
					applyContentBounds(&sc, es)
					continue
				}
				appendReaderShape(&sc, shapeIndexByID, es)
			default:
				appendReaderShape(&sc, shapeIndexByID, es)
			}
		}
	}
	for _, es := range connectorShapes {
		connector, ok := editorShapeToConnector(es, shapeIndexByID)
		if !ok {
			appendReaderShape(&sc, shapeIndexByID, es)
			continue
		}
		sc.Connectors = append(sc.Connectors, connector)
	}

	// Attach images for this slide.
	if sm.Index < len(slideImages) {
		for _, img := range slideImages[sm.Index] {
			sc.Images = append(sc.Images, shapes.Image{
				Data:     img.Bytes,
				Format:   img.Format,
				X:        styling.Emu(img.X),
				Y:        styling.Emu(img.Y),
				CX:       styling.Emu(img.CX),
				CY:       styling.Emu(img.CY),
				Rotation: img.Rotation,
				Crop: shapes.ImageCrop{
					Left:   img.CropLeft,
					Right:  img.CropRight,
					Top:    img.CropTop,
					Bottom: img.CropBottom,
				},
				FlipH:        img.FlipH,
				FlipV:        img.FlipV,
				Shadow:       img.Shadow,
				Reflection:   img.Reflection,
				AltText:      img.AltText,
				IsDecorative: img.IsDecorative,
			})
		}
	}
	if sm.Index < len(slideCharts) {
		applyParsedCharts(&sc, slideCharts[sm.Index])
	}
	if sm.Index < len(slideSmartArt) {
		applyParsedSmartArt(&sc, slideSmartArt[sm.Index])
	}
	return sc
}

func appendReaderShape(sc *elements.SlideContent, shapeIndexByID map[int]int, es editorcommon.Shape) {
	sc.Shapes = append(sc.Shapes, editorShapeToShape(es))
	if es.ID > 0 {
		shapeIndexByID[es.ID] = len(sc.Shapes)
	}
}
