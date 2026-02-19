package editor

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func handleSlideCount(e *PresentationEditor, _ json.RawMessage) (any, error) {
	return map[string]int{"count": e.SlideCount()}, nil
}

func handleAddSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	// Parse and validate payload
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	title := v.OptionalString(p, "title")
	layout := v.OptionalString(p, "layout")
	bullets, _ := v.OptionalStringSlice(p, "bullets")

	if v.HasErrors() {
		return nil, v.Error()
	}

	slide := elements.NewSlide(title)
	if layout != "" {
		slide = slide.WithLayout(layout)
	}
	for _, b := range bullets {
		slide = slide.AddBullet(b)
	}
	index, err := e.AddSlide(slide)
	if err != nil {
		return nil, err
	}
	return map[string]int{"index": index}, nil
}

func handleRemoveSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	index, ok := v.RequireInt(p, "index")
	if !ok {
		return nil, v.Error()
	}
	if !v.IndexBounds(index, 0, e.SlideCount(), "index") {
		return nil, v.Error()
	}

	return nil, e.RemoveSlide(index)
}

func handleMoveSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	from, ok := v.RequireInt(p, "from")
	if !ok {
		return nil, v.Error()
	}
	to, ok := v.RequireInt(p, "to")
	if !ok {
		return nil, v.Error()
	}
	slideCount := e.SlideCount()
	if !v.IndexBounds(from, 0, slideCount, "from") {
		return nil, v.Error()
	}
	if !v.IndexBounds(to, 0, slideCount, "to") {
		return nil, v.Error()
	}

	return nil, e.MoveSlide(from, to)
}

func handleDuplicateSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	index, ok := v.RequireInt(p, "index")
	if !ok {
		return nil, v.Error()
	}
	insertAt, _ := v.OptionalInt(p, "insert_at")

	if !v.IndexBounds(index, 0, e.SlideCount(), "index") {
		return nil, v.Error()
	}

	newIdx, err := e.DuplicateSlide(index, insertAt)
	if err != nil {
		return nil, err
	}
	return map[string]int{"new_index": newIdx}, nil
}

func handleGetMetadata(e *PresentationEditor, _ json.RawMessage) (any, error) {
	m := e.Metadata()
	return map[string]any{
		"title":       m.Title,
		"slide_count": m.SlideCount,
		"size": map[string]int64{
			"width":  m.SlideSize.Width,
			"height": m.SlideSize.Height,
		},
	}, nil
}

func handleUpdateChartData(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	var params struct {
		ChartSelector common.ChartSelector   `json:"chart_selector"`
		Data          common.ChartDataUpdate `json:"data"`
	}
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}

	if err := e.UpdateChartData(slideIndex, params.ChartSelector, params.Data); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleListSlideCharts(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	refs, err := e.ListSlideCharts(slideIndex)
	if err != nil {
		return nil, err
	}
	return map[string]any{"charts": refs}, nil
}

func handleListSlideLayouts(e *PresentationEditor, _ json.RawMessage) (any, error) {
	layouts, err := e.ListSlideLayouts()
	if err != nil {
		return nil, err
	}
	return map[string]any{"layouts": layouts}, nil
}

func handleRebindSlideLayout(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	layoutPart, ok := v.RequireString(p, "layout_part")
	if !ok {
		return nil, v.Error()
	}
	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	if err := e.RebindSlideLayout(slideIndex, layoutPart); err != nil {
		return nil, err
	}
	return map[string]bool{"rebound": true}, nil
}

func handleCloneLayoutMasterFamily(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	layoutPart, ok := v.RequireString(p, "layout_part")
	if !ok {
		return nil, v.Error()
	}

	result, err := e.CloneLayoutMasterFamily(layoutPart)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"master_part": result.MasterPart,
		"theme_part":  result.ThemePart,
		"layout_map":  result.LayoutMap,
	}, nil
}
