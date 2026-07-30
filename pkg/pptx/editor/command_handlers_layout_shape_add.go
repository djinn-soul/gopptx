package editor

import (
	"encoding/json"
)

// layoutShapeRequest is the shared payload for adding a shape to a layout or
// master: the part path, the bounds, and either a shape type or text.
type layoutShapeRequest struct {
	partPath   string
	shapeType  string
	text       string
	x, y, w, h float64
}

func parseLayoutShapeRequest(
	payload json.RawMessage,
	partKey string,
	wantShapeType bool,
	wantText bool,
) (layoutShapeRequest, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return layoutShapeRequest{}, err
	}

	v := NewPayloadValidator()
	req := layoutShapeRequest{}
	partPath, ok := v.RequireString(p, partKey)
	if !ok {
		return layoutShapeRequest{}, v.Error()
	}
	req.partPath = partPath

	if wantShapeType {
		shapeType, okType := v.RequireString(p, "shape_type")
		if !okType {
			return layoutShapeRequest{}, v.Error()
		}
		req.shapeType = shapeType
	}
	if wantText {
		if text, okText := p["text"].(string); okText {
			req.text = text
		}
	}

	bounds := []struct {
		key string
		dst *float64
	}{
		{keyLeft, &req.x}, {keyTop, &req.y}, {keyWidth, &req.w}, {keyHeight, &req.h},
	}
	for _, bound := range bounds {
		value, okNum := v.RequireFloat64(p, bound.key)
		if !okNum {
			return layoutShapeRequest{}, v.Error()
		}
		*bound.dst = value
	}
	return req, nil
}

// handleAddLayoutShape adds an autoshape to a slide layout.
//
// Payload: {"layout_part", "shape_type", "left", "top", "width", "height"}.
func handleAddLayoutShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	req, err := parseLayoutShapeRequest(payload, "layout_part", true, false)
	if err != nil {
		return nil, err
	}
	shapeID, err := e.AddLayoutShape(req.partPath, req.shapeType, req.x, req.y, req.w, req.h)
	if err != nil {
		return nil, err
	}
	return map[string]any{keyShapeID: shapeID}, nil
}

// handleAddLayoutTextbox adds a text box to a slide layout.
//
// Payload: {"layout_part", "text", "left", "top", "width", "height"}.
func handleAddLayoutTextbox(e *PresentationEditor, payload json.RawMessage) (any, error) {
	req, err := parseLayoutShapeRequest(payload, "layout_part", false, true)
	if err != nil {
		return nil, err
	}
	shapeID, err := e.AddLayoutTextBox(req.partPath, req.text, req.x, req.y, req.w, req.h)
	if err != nil {
		return nil, err
	}
	return map[string]any{keyShapeID: shapeID}, nil
}

// handleAddMasterShape adds an autoshape to a slide master.
//
// Payload: {"master_part", "shape_type", "left", "top", "width", "height"}.
func handleAddMasterShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	req, err := parseLayoutShapeRequest(payload, "master_part", true, false)
	if err != nil {
		return nil, err
	}
	shapeID, err := e.AddMasterShape(req.partPath, req.shapeType, req.x, req.y, req.w, req.h)
	if err != nil {
		return nil, err
	}
	return map[string]any{keyShapeID: shapeID}, nil
}

// handleAddMasterTextbox adds a text box to a slide master.
//
// Payload: {"master_part", "text", "left", "top", "width", "height"}.
func handleAddMasterTextbox(e *PresentationEditor, payload json.RawMessage) (any, error) {
	req, err := parseLayoutShapeRequest(payload, "master_part", false, true)
	if err != nil {
		return nil, err
	}
	shapeID, err := e.AddMasterTextBox(req.partPath, req.text, req.x, req.y, req.w, req.h)
	if err != nil {
		return nil, err
	}
	return map[string]any{keyShapeID: shapeID}, nil
}
