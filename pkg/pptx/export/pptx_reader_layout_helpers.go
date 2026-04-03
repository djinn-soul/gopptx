package export

import (
	"strings"

	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

const (
	editorTypeTriangle = "triangle"
)

func chartShapeIDSet(slideCharts [][]parsedChart, slideIndex int) map[int]struct{} {
	out := make(map[int]struct{})
	if slideIndex < 0 || slideIndex >= len(slideCharts) {
		return out
	}
	for _, chart := range slideCharts[slideIndex] {
		if chart.ShapeID <= 0 {
			continue
		}
		out[chart.ShapeID] = struct{}{}
	}
	return out
}

func smartArtShapeIDSet(slideSmartArt [][]parsedSmartArt, slideIndex int) map[int]struct{} {
	out := make(map[int]struct{})
	if slideIndex < 0 || slideIndex >= len(slideSmartArt) {
		return out
	}
	for _, diagram := range slideSmartArt[slideIndex] {
		if diagram.ShapeID <= 0 {
			continue
		}
		out[diagram.ShapeID] = struct{}{}
	}
	return out
}

// applyTitleBounds records title placeholder geometry for downstream renderers.
func applyTitleBounds(sc *elements.SlideContent, es editorcommon.Shape) {
	if es.W == 0 && es.H == 0 {
		return
	}
	sc.TitleBoundsEMU = [4]int64{int64(es.X), int64(es.Y), int64(es.W), int64(es.H)}
}

// applyTitleSizeFromRuns captures the first explicit run font size on title text.
func applyTitleSizeFromRuns(sc *elements.SlideContent, es editorcommon.Shape) {
	if sc.TitleSize > 0 {
		return
	}
	for _, run := range es.Runs {
		if run.SizePt != nil && *run.SizePt > 0 {
			sc.TitleSize = *run.SizePt
			return
		}
	}
}

// applyTitleAlignFromShape captures the first explicit paragraph alignment.
func applyTitleAlignFromShape(sc *elements.SlideContent, es editorcommon.Shape) {
	if sc.TitleAlign != "" {
		return
	}
	for _, paragraph := range es.Paragraphs {
		if paragraph.Paragraph != nil && paragraph.Paragraph.Alignment != nil && *paragraph.Paragraph.Alignment != "" {
			sc.TitleAlign = *paragraph.Paragraph.Alignment
			return
		}
	}
	if es.Paragraph != nil && es.Paragraph.Alignment != nil && *es.Paragraph.Alignment != "" {
		sc.TitleAlign = *es.Paragraph.Alignment
	}
}

// applyContentBounds records body placeholder geometry for bullet layout.
func applyContentBounds(sc *elements.SlideContent, es editorcommon.Shape) {
	if es.W == 0 && es.H == 0 {
		return
	}
	sc.ContentBoundsEMU = [4]int64{int64(es.X), int64(es.Y), int64(es.W), int64(es.H)}
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
	return shapeName == placeholderContentName ||
		shapeName == placeholderBody ||
		strings.Contains(shapeName, "content placeholder") ||
		strings.Contains(shapeName, "body placeholder")
}

// editorTypeToPreset normalizes editor shape types to OOXML preset names.
func editorTypeToPreset(value string) string {
	switch strings.ToLower(value) {
	case "rect", "rectangle":
		return "rect"
	case "roundrect", "roundedrectangle":
		return "roundRect"
	case "ellipse", "oval", "circle":
		return "ellipse"
	case editorTypeTriangle, "rt_triangle":
		return editorTypeTriangle
	case "rightarrow":
		return "rightArrow"
	case "leftarrow":
		return "leftArrow"
	default:
		return value
	}
}
