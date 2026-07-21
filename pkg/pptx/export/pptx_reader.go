package export

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
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
	title, slides, _, err := slidesFromPPTXWithSize(pptxPath)
	return title, slides, err
}

// slidesFromPPTXWithSize is SlidesFromPPTX plus the deck's slide size in EMUs,
// which the native PDF renderer needs to size its pages.
func slidesFromPPTXWithSize(pptxPath string) (string, []elements.SlideContent, common.SlideSize, error) {
	var slideSize common.SlideSize

	ed, err := editor.OpenPresentationEditor(pptxPath)
	if err != nil {
		return "", nil, slideSize, fmt.Errorf("open PPTX: %w", err)
	}
	defer ed.Close()

	meta := ed.Metadata()
	presTitle := ""
	if meta != nil {
		presTitle = meta.Title
		slideSize = meta.SlideSize
	}

	// Read slide master's txStyles to resolve inherited title alignment/size.
	// Slides with placeholder references (empty <p:spPr/>) inherit these defaults
	// via the OOXML chain: slide → layout → master → txStyles.
	masterStyle := extractMasterTitleStyle(pptxPath)

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
		sc := extractSlideContent(ed, sm, slideImages, slideCharts, slideSmartArt)
		applyMasterTitleDefaults(&sc, masterStyle)
		slideContents = append(slideContents, sc)
	}

	if presTitle == "" && len(slideContents) > 0 {
		presTitle = slideContents[0].Title
	}

	return presTitle, slideContents, slideSize, nil
}

// applyMasterTitleDefaults fills in title alignment and size from master txStyles
// when they were not explicitly set by the slide or layout XML.
func applyMasterTitleDefaults(sc *elements.SlideContent, ms masterTitleStyle) {
	if sc.TitleAlign == "" && ms.Align != "" {
		sc.TitleAlign = ms.Align
	}
	if sc.TitleSize == 0 && ms.SizePt > 0 {
		sc.TitleSize = ms.SizePt
	}
}

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

	sc := elements.SlideContent{Title: sm.Title, Hidden: sm.Hidden}
	applySlideMetadata(&sc, ed, sm.Index)

	shapeIndexByID := make(map[int]int)
	connectorShapes := make([]editorcommon.Shape, 0)
	chartShapeIDs := chartShapeIDSet(slideCharts, sm.Index)
	smartArtShapeIDs := smartArtShapeIDSet(slideSmartArt, sm.Index)

	for _, es := range editorShapes {
		if isEditorConnector(es) {
			connectorShapes = append(connectorShapes, es)
			continue
		}
		processEditorShape(&sc, es, shapeIndexByID, chartShapeIDs, smartArtShapeIDs, ed, sm.Index)
	}
	for _, es := range connectorShapes {
		connector, ok := editorShapeToConnector(es, shapeIndexByID)
		if !ok {
			appendReaderShape(&sc, shapeIndexByID, es)
			continue
		}
		sc.Connectors = append(sc.Connectors, connector)
	}
	foldGeneratedConnectorLabels(&sc)

	attachSlideImages(&sc, slideImages, sm.Index)
	if sm.Index < len(slideCharts) {
		applyParsedCharts(&sc, slideCharts[sm.Index])
	}
	if sm.Index < len(slideSmartArt) {
		applyParsedSmartArt(&sc, slideSmartArt[sm.Index])
	}
	return sc
}

// applySlideMetadata reads background, header/footer, and notes from the editor.
func applySlideMetadata(sc *elements.SlideContent, ed *editor.PresentationEditor, idx int) {
	if bg, bgErr := ed.GetSlideBackground(idx); bgErr == nil && bg != nil {
		sc.Background = bg
	}
	if hf, hfErr := ed.GetSlideHeaderFooter(idx); hfErr == nil {
		if hf.ShowFooter {
			sc.FooterText = hf.Footer
		}
		if hf.ShowSlideNum {
			sc.ShowSlideNumber = true
		}
	}
	if notes, notesErr := ed.GetNotes(idx); notesErr == nil && notes != "" {
		sc.Notes = notes
	}
}

// processEditorShape classifies a single editor shape and applies it to the slide content.
func processEditorShape(
	sc *elements.SlideContent,
	es editorcommon.Shape,
	shapeIndexByID map[int]int,
	chartShapeIDs, smartArtShapeIDs map[int]struct{},
	ed *editor.PresentationEditor,
	slideIdx int,
) {
	lowerType := strings.ToLower(es.Type)
	lowerName := strings.ToLower(strings.TrimSpace(es.Name))
	// For shapes that carry a placeholder type but no geometry preset (e.g.
	// python-pptx generated shapes whose <p:spPr/> is empty), use the
	// placeholder type as the effective classification type so that title and
	// body placeholders are still handled correctly.
	if lowerPhType := strings.ToLower(es.PlaceholderType); lowerPhType != "" && (lowerType == "" || lowerType == "sp") {
		lowerType = lowerPhType
	}
	switch lowerType {
	case "pic":
		return
	case "graphicframe":
		applyGraphicFrame(sc, es, shapeIndexByID, chartShapeIDs, smartArtShapeIDs, ed, slideIdx)
	case placeholderTitle, placeholderCtrTitle:
		applyTitleShape(sc, es, lowerType)
	case placeholderBody, placeholderSubtitle, placeholderObject:
		applyBodyShape(sc, es, shapeIndexByID)
	default:
		applyDefaultShape(sc, es, shapeIndexByID, lowerType, lowerName)
	}
}

func applyGraphicFrame(
	sc *elements.SlideContent,
	es editorcommon.Shape,
	shapeIndexByID map[int]int,
	chartShapeIDs, smartArtShapeIDs map[int]struct{},
	ed *editor.PresentationEditor,
	slideIdx int,
) {
	if _, ok := chartShapeIDs[es.ID]; ok {
		return
	}
	if _, ok := smartArtShapeIDs[es.ID]; ok {
		return
	}
	if tbl := extractTableContent(ed, slideIdx, es); tbl != nil {
		// Keep the first table in Table; later ones would otherwise overwrite it.
		if sc.Table == nil {
			sc.Table = tbl
		} else {
			sc.Tables = append(sc.Tables, *tbl)
		}
		return
	}
	appendReaderShape(sc, shapeIndexByID, es)
}

func applyTitleShape(sc *elements.SlideContent, es editorcommon.Shape, lowerType string) {
	if sc.Title == "" && es.Text != "" {
		sc.Title = es.Text
	}
	applyTitleBounds(sc, es)
	applyTitleSizeFromRuns(sc, es)
	applyTitleAlignFromShape(sc, es)
	if lowerType == placeholderCtrTitle && sc.Layout == "" {
		sc.Layout = elements.SlideLayoutCenteredTitle
	}
}

func applyBodyShape(sc *elements.SlideContent, es editorcommon.Shape, shapeIndexByID map[int]int) {
	if consumeBodyPlaceholderAsBullets(sc, es) {
		applyContentBounds(sc, es)
		return
	}
	appendReaderShape(sc, shapeIndexByID, es)
}

func applyDefaultShape(
	sc *elements.SlideContent, es editorcommon.Shape,
	shapeIndexByID map[int]int, lowerType, lowerName string,
) {
	switch {
	case isTitlePlaceholder(lowerType, lowerName):
		if sc.Title == "" && es.Text != "" {
			sc.Title = es.Text
		}
		applyTitleBounds(sc, es)
		applyTitleSizeFromRuns(sc, es)
		applyTitleAlignFromShape(sc, es)
	case isBodyPlaceholder(lowerType, lowerName):
		applyBodyShape(sc, es, shapeIndexByID)
	default:
		appendReaderShape(sc, shapeIndexByID, es)
	}
}

// attachSlideImages appends all embedded images for the given slide index.
func attachSlideImages(sc *elements.SlideContent, slideImages [][]SlideImage, idx int) {
	if idx >= len(slideImages) {
		return
	}
	for _, img := range slideImages[idx] {
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

func appendReaderShape(sc *elements.SlideContent, shapeIndexByID map[int]int, es editorcommon.Shape) {
	sc.Shapes = append(sc.Shapes, editorShapeToShape(es))
	if es.ID > 0 {
		shapeIndexByID[es.ID] = len(sc.Shapes)
	}
}
