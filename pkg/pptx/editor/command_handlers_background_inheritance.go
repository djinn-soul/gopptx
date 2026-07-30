package editor

import "encoding/json"

func handleSetSlideFollowMasterBackground(
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
	follow, ok := v.RequireBool(p, "follow")
	if !ok {
		return nil, v.Error()
	}
	if err := e.SetSlideFollowMasterBackground(slideIndex, follow); err != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
	}
	return respUpdated, nil
}
