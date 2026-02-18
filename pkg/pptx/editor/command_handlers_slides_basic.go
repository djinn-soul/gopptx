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
	var p struct {
		Title   string   `json:"title"`
		Layout  string   `json:"layout"`
		Bullets []string `json:"bullets"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	slide := elements.NewSlide(p.Title)
	if p.Layout != "" {
		slide = slide.WithLayout(p.Layout)
	}
	for _, b := range p.Bullets {
		slide = slide.AddBullet(b)
	}
	index, err := e.AddSlide(slide)
	if err != nil {
		return nil, err
	}
	return map[string]int{"index": index}, nil
}

func handleRemoveSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		Index int `json:"index"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	return nil, e.RemoveSlide(p.Index)
}

func handleMoveSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		From int `json:"from"`
		To   int `json:"to"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	return nil, e.MoveSlide(p.From, p.To)
}

func handleDuplicateSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		Index    int `json:"index"`
		InsertAt int `json:"insert_at"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	newIdx, err := e.DuplicateSlide(p.Index, p.InsertAt)
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
	var p struct {
		SlideIndex    int                    `json:"slide_index"`
		ChartSelector common.ChartSelector   `json:"chart_selector"`
		Data          common.ChartDataUpdate `json:"data"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.UpdateChartData(p.SlideIndex, p.ChartSelector, p.Data); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleListSlideCharts(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex int `json:"slide_index"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	refs, err := e.ListSlideCharts(p.SlideIndex)
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
	var p struct {
		SlideIndex int    `json:"slide_index"`
		LayoutPart string `json:"layout_part"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.RebindSlideLayout(p.SlideIndex, p.LayoutPart); err != nil {
		return nil, err
	}
	return map[string]bool{"rebound": true}, nil
}

func handleCloneLayoutMasterFamily(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		LayoutPart string `json:"layout_part"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	result, err := e.CloneLayoutMasterFamily(p.LayoutPart)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"master_part": result.MasterPart,
		"theme_part":  result.ThemePart,
		"layout_map":  result.LayoutMap,
	}, nil
}
