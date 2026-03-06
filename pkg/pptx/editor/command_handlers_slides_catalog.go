package editor

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

func handleUpdateChartData(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
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
	v := NewPayloadValidator()
	return editorcommand.HandleSlideIndexRequest(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.Error,
		func(slideIndex int) (any, error) {
			refs, err := e.ListSlideCharts(slideIndex)
			if err != nil {
				return nil, err
			}
			return map[string]any{"charts": refs}, nil
		},
	)
}

func handleListSlideLayouts(e *PresentationEditor, _ json.RawMessage) (any, error) {
	layouts, err := e.ListSlideLayouts()
	if err != nil {
		return nil, err
	}
	return map[string]any{"layouts": layouts}, nil
}

func handleListSlideMasters(e *PresentationEditor, _ json.RawMessage) (any, error) {
	masters, err := e.ListSlideMasters()
	if err != nil {
		return nil, err
	}
	return map[string]any{"masters": masters}, nil
}

func handleListMasterLayouts(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	masterPart, ok := v.RequireString(p, "master_part")
	if !ok {
		return nil, v.Error()
	}

	layouts, err := e.ListMasterLayouts(masterPart)
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
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	layoutPart, ok := v.RequireString(p, "layout_part")
	if !ok {
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

func handleAddSlideMaster(e *PresentationEditor, _ json.RawMessage) (any, error) {
	masterPart, err := e.AddSlideMaster()
	if err != nil {
		return nil, err
	}
	return map[string]any{"master_part": masterPart}, nil
}

func handleRemoveSlideMaster(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	masterPart, ok := v.RequireString(p, "master_part")
	if !ok {
		return nil, v.Error()
	}

	if err := e.RemoveSlideMaster(masterPart); err != nil {
		return nil, err
	}
	return map[string]bool{"removed": true}, nil
}

func handleAddSlideLayout(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	masterPart, ok := v.RequireString(p, "master_part")
	if !ok {
		return nil, v.Error()
	}
	layoutName, ok := v.RequireString(p, "layout_name")
	if !ok {
		layoutName = "Custom Layout"
	}

	layoutPart, err := e.AddSlideLayout(masterPart, layoutName)
	if err != nil {
		return nil, err
	}
	return map[string]any{"layout_part": layoutPart}, nil
}

func handleRemoveSlideLayout(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	layoutPart, ok := v.RequireString(p, "layout_part")
	if !ok {
		return nil, v.Error()
	}

	if err := e.RemoveSlideLayout(layoutPart); err != nil {
		return nil, err
	}
	return map[string]bool{"removed": true}, nil
}
