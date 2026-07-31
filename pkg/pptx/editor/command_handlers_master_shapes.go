package editor

import (
	"encoding/json"
)

func handleGetSlideShowMasterShapes(
	e *PresentationEditor,
	payload json.RawMessage,
) (any, error) {
	slideIndex, err := parseMasterShapeSlideIndex(e, payload)
	if err != nil {
		return nil, err
	}
	visible, err := e.SlideShowMasterShapes(slideIndex)
	if err != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
	}
	return map[string]any{"visible": visible}, nil
}

func handleSetSlideShowMasterShapes(
	e *PresentationEditor,
	payload json.RawMessage,
) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	visible, ok := v.RequireBool(p, "visible")
	if !ok {
		return nil, v.Error()
	}
	if err := e.SetSlideShowMasterShapes(slideIndex, visible); err != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
	}
	return respUpdated, nil
}

func parseMasterShapeSlideIndex(
	e *PresentationEditor,
	payload json.RawMessage,
) (int, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return 0, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return 0, v.Error()
	}
	return slideIndex, nil
}
