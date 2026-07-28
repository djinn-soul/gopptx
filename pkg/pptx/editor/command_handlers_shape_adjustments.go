package editor

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	slidescatalog "github.com/djinn-soul/gopptx/pkg/pptx/editor/handlers/slidescatalog"
)

// handleSetShapeAdjustments writes preset-geometry adjustment values, the
// yellow handles in PowerPoint's UI.
//
// Payload: {"slide_index": N, "shape_id": I,
//
//	"adjustments": [{"name": "adj1", "value": 0.5}]}.
//
// Response: {"updated": true}.
func handleSetShapeAdjustments(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	shapeID, _ := v.RequireInt(p, "shape_id")
	if v.HasErrors() {
		return nil, v.Error()
	}

	var params struct {
		Adjustments []common.ShapeAdjustmentValue `json:"adjustments"`
	}
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}

	if err := e.SetShapeAdjustments(slideIndex, shapeID, params.Adjustments); err != nil {
		return nil, err
	}
	return slidescatalog.BuildUpdatedResponse(), nil
}
