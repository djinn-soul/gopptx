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

	slideMeta := ed.Slides()
	slideContents := make([]elements.SlideContent, 0, len(slideMeta))

	for _, sm := range slideMeta {
		slideContents = append(slideContents, extractSlideContent(ed, sm, slideImages, slideCharts))
	}

	if presTitle == "" && len(slideContents) > 0 {
		presTitle = slideContents[0].Title
	}

	return presTitle, slideContents, nil
}

//nolint:gocognit // Slide extraction maps multiple editor shape categories and image attachments in one pass.
func extractSlideContent(
	ed *editor.PresentationEditor,
	sm editorcommon.SlideMetadata,
	slideImages [][]SlideImage,
	slideCharts [][]parsedChart,
) elements.SlideContent {
	editorShapes, err := ed.GetShapes(sm.Index)
	if err != nil {
		editorShapes = nil
	}

	sc := elements.SlideContent{
		Title: sm.Title,
	}

	for _, es := range editorShapes {
		lowerType := strings.ToLower(es.Type)
		lowerName := strings.ToLower(strings.TrimSpace(es.Name))

		switch lowerType {
		case "graphicframe":
			if tbl := extractTableContent(ed, sm.Index, es); tbl != nil {
				sc.Table = tbl
				continue
			}
			sc.Shapes = append(sc.Shapes, editorShapeToShape(es))
		case placeholderTitle, placeholderCtrTitle:
			if sc.Title == "" && es.Text != "" {
				sc.Title = es.Text
			}
		case placeholderBody, placeholderSubtitle, placeholderObject:
			if consumeBodyPlaceholderAsBullets(&sc, es.Text) {
				continue
			}
			sc.Shapes = append(sc.Shapes, editorShapeToShape(es))
		default:
			switch {
			case isTitlePlaceholder(lowerType, lowerName):
				if sc.Title == "" && es.Text != "" {
					sc.Title = es.Text
				}
			case isBodyPlaceholder(lowerType, lowerName):
				if consumeBodyPlaceholderAsBullets(&sc, es.Text) {
					continue
				}
				sc.Shapes = append(sc.Shapes, editorShapeToShape(es))
			default:
				sc.Shapes = append(sc.Shapes, editorShapeToShape(es))
			}
		}
	}

	// Attach images for this slide.
	if sm.Index < len(slideImages) {
		for _, img := range slideImages[sm.Index] {
			sc.Images = append(sc.Images, shapes.Image{
				Data:   img.Bytes,
				Format: img.Format,
				X:      styling.Emu(img.X),
				Y:      styling.Emu(img.Y),
				CX:     styling.Emu(img.CX),
				CY:     styling.Emu(img.CY),
			})
		}
	}
	if sm.Index < len(slideCharts) {
		applyParsedCharts(&sc, slideCharts[sm.Index])
	}
	return sc
}

// editorShapeToShape maps an editor common.Shape to an export shapes.Shape.
// X, Y, W, H from the editor are in EMU (int).
func editorShapeToShape(es editorcommon.Shape) shapes.Shape {
	shapeFill := editorFillToExportFill(es.Fill)
	shapeLine := editorLineToExportLine(es.Line)
	adjustments := editorAdjustmentsToExport(es.Adjustments)
	return shapes.Shape{
		// Map OOXML preset geometry name directly — pdf_native.go uses these strings.
		Type:        editorTypeToPreset(es.Type),
		X:           styling.Emu(int64(es.X)),
		Y:           styling.Emu(int64(es.Y)),
		CX:          styling.Emu(int64(es.W)),
		CY:          styling.Emu(int64(es.H)),
		Text:        es.Text,
		Name:        es.Name,
		Fill:        shapeFill,
		Line:        shapeLine,
		Adjustments: adjustments,
	}
}

func editorFillToExportFill(fill *editorcommon.ShapeFill) *shapes.ShapeFill {
	if fill == nil || fill.Solid == nil || *fill.Solid == "" {
		return nil
	}
	return &shapes.ShapeFill{Color: *fill.Solid}
}

func editorLineToExportLine(line *editorcommon.ShapeLine) *shapes.ShapeLine {
	if line == nil || line.Color == nil || *line.Color == "" {
		return nil
	}
	width := styling.Emu(0)
	if line.WidthEmu != nil && *line.WidthEmu > 0 {
		width = styling.Emu(int64(*line.WidthEmu))
	}
	return &shapes.ShapeLine{
		Color: *line.Color,
		Width: width,
	}
}

func editorAdjustmentsToExport(
	adjustments []editorcommon.ShapeAdjustment,
) []shapes.ShapeAdjustment {
	if len(adjustments) == 0 {
		return nil
	}
	out := make([]shapes.ShapeAdjustment, 0, len(adjustments))
	for _, adjustment := range adjustments {
		if adjustment.Name == "" || adjustment.Formula == "" {
			continue
		}
		out = append(out, shapes.ShapeAdjustment{
			Name:    adjustment.Name,
			Formula: adjustment.Formula,
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func isTitlePlaceholder(shapeType, shapeName string) bool {
	return shapeType == placeholderTitle ||
		shapeType == placeholderCtrTitle ||
		shapeName == placeholderTitle ||
		strings.Contains(shapeName, "title placeholder")
}

func isBodyPlaceholder(shapeType, shapeName string) bool {
	if shapeType == placeholderBody || shapeType == placeholderSubtitle || shapeType == placeholderObject {
		return true
	}
	return shapeName == "content" ||
		shapeName == placeholderBody ||
		strings.Contains(shapeName, "content placeholder") ||
		strings.Contains(shapeName, "body placeholder")
}

func consumeBodyPlaceholderAsBullets(sc *elements.SlideContent, bodyText string) bool {
	bodyText = strings.TrimSpace(bodyText)
	if bodyText == "" {
		return false
	}
	parts := strings.Split(bodyText, "\n")
	normalized := make([]string, 0, len(parts))
	for _, line := range parts {
		line = strings.TrimSpace(strings.TrimPrefix(line, "•"))
		if line == "" {
			continue
		}
		normalized = append(normalized, line)
	}
	if len(normalized) == 0 {
		return false
	}
	// Body placeholders on title+content layouts are usually bullet paragraphs.
	sc.Bullets = append(sc.Bullets, normalized...)
	return true
}

// editorTypeToPreset normalizes the editor shape type string to the OOXML
// preset geometry name used by our shape renderers.
func editorTypeToPreset(t string) string {
	switch strings.ToLower(t) {
	case "rect", "rectangle":
		return "rect"
	case "roundrect", "roundedrectangle":
		return "roundRect"
	case "ellipse", "oval", "circle":
		return "ellipse"
	case "triangle", "rt_triangle":
		return "triangle"
	case "rightarrow":
		return "rightArrow"
	case "leftarrow":
		return "leftArrow"
	default:
		// Pass through unknown presets as-is; renderers will fall back to rect.
		return t
	}
}
