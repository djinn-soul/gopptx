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

	slideMeta := ed.Slides()
	slideContents := make([]elements.SlideContent, 0, len(slideMeta))

	for _, sm := range slideMeta {
		slideContents = append(slideContents, extractSlideContent(ed, sm, slideImages))
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
		case "title", "ctrtitle":
			if sc.Title == "" && es.Text != "" {
				sc.Title = es.Text
			}
		case "body", "subtitle", "obj":
			if es.Text != "" {
				for line := range strings.SplitSeq(es.Text, "\n") {
					if line = strings.TrimSpace(line); line != "" {
						sc.Bullets = append(sc.Bullets, line)
					}
				}
			}
		default:
			switch {
			case isTitlePlaceholder(lowerType, lowerName):
				if sc.Title == "" && es.Text != "" {
					sc.Title = es.Text
				}
			case isBodyPlaceholder(lowerType, lowerName):
				for line := range strings.SplitSeq(es.Text, "\n") {
					if line = strings.TrimSpace(line); line != "" {
						sc.Bullets = append(sc.Bullets, line)
					}
				}
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
	return sc
}

// editorShapeToShape maps an editor common.Shape to an export shapes.Shape.
// X, Y, W, H from the editor are in EMU (int).
func editorShapeToShape(es editorcommon.Shape) shapes.Shape {
	return shapes.Shape{
		// Map OOXML preset geometry name directly — pdf_native.go uses these strings.
		Type: editorTypeToPreset(es.Type),
		X:    styling.Emu(int64(es.X)),
		Y:    styling.Emu(int64(es.Y)),
		CX:   styling.Emu(int64(es.W)),
		CY:   styling.Emu(int64(es.H)),
		Text: es.Text,
		Name: es.Name,
	}
}

func isTitlePlaceholder(shapeType, shapeName string) bool {
	return shapeType == "title" ||
		shapeType == "ctrtitle" ||
		shapeName == "title" ||
		strings.Contains(shapeName, "title placeholder")
}

func isBodyPlaceholder(shapeType, shapeName string) bool {
	if shapeType == "body" || shapeType == "subtitle" || shapeType == "obj" {
		return true
	}
	return shapeName == "content" ||
		shapeName == "body" ||
		strings.Contains(shapeName, "content placeholder") ||
		strings.Contains(shapeName, "body placeholder")
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
