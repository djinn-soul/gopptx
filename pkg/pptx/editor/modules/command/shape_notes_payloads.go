package command

import common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"

type ParseSlideIndexFn func(map[string]any) (int, bool)
type ParseIntFieldFn func(map[string]any, string) (int, bool)
type ParseIntSliceFieldFn func(map[string]any, string) ([]int, bool)
type ParseStringFieldFn func(map[string]any, string) (string, bool)

type SlideShapeRequest struct {
	SlideIndex int
	ShapeID    int
}

func ParseSlideIndexRequest(
	payload map[string]any,
	parseSlideIndex ParseSlideIndexFn,
) (int, bool) {
	return parseSlideIndex(payload)
}

func ParseSlideShapeRequest(
	payload map[string]any,
	parseSlideIndex ParseSlideIndexFn,
	parseIntField ParseIntFieldFn,
) (SlideShapeRequest, bool) {
	slideIndex, ok := parseSlideIndex(payload)
	if !ok {
		return SlideShapeRequest{}, false
	}
	shapeID, ok := parseIntField(payload, "shape_id")
	if !ok {
		return SlideShapeRequest{}, false
	}
	return SlideShapeRequest{SlideIndex: slideIndex, ShapeID: shapeID}, true
}

type SlideShapeIDsRequest struct {
	SlideIndex int
	ShapeIDs   []int
}

func ParseSlideShapeIDsRequest(
	payload map[string]any,
	parseSlideIndex ParseSlideIndexFn,
	parseIntSliceField ParseIntSliceFieldFn,
) (SlideShapeIDsRequest, bool) {
	slideIndex, ok := parseSlideIndex(payload)
	if !ok {
		return SlideShapeIDsRequest{}, false
	}
	shapeIDs, ok := parseIntSliceField(payload, "shape_ids")
	if !ok {
		return SlideShapeIDsRequest{}, false
	}
	return SlideShapeIDsRequest{SlideIndex: slideIndex, ShapeIDs: shapeIDs}, true
}

func ParseSetNotesRequest(
	payload map[string]any,
	parseSlideIndex ParseSlideIndexFn,
	parseStringField ParseStringFieldFn,
) (int, string, bool) {
	slideIndex, ok := parseSlideIndex(payload)
	if !ok {
		return 0, "", false
	}
	text, ok := parseStringField(payload, "text")
	if !ok {
		return 0, "", false
	}
	return slideIndex, text, true
}

func BuildNotesResult(text string, hasNotesSlide bool) map[string]any {
	return BuildNotesResultDetailed(text, hasNotesSlide, nil)
}

func BuildNotesResultDetailed(
	text string,
	hasNotesSlide bool,
	placeholders []common.PlaceholderInfo,
) map[string]any {
	var notesSlide any
	if hasNotesSlide {
		notesSlide = map[string]any{
			"text": text,
		}
	}
	return map[string]any{
		"text":               text,
		"notes_slide":        notesSlide,
		"notes_placeholders": placeholders,
	}
}
