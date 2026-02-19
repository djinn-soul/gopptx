package editor

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func handleListShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

	shapes, err := e.GetShapes(slideIndex)
	if err != nil {
		return nil, err
	}
	return map[string]any{"shapes": shapes}, nil
}

func handleAddShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	shapeType, ok := v.RequireString(p, "type")
	if !ok {
		return nil, v.Error()
	}
	x, ok := v.RequireFloat64(p, "x")
	if !ok {
		return nil, v.Error()
	}
	y, ok := v.RequireFloat64(p, "y")
	if !ok {
		return nil, v.Error()
	}
	w, ok := v.RequireFloat64(p, "w")
	if !ok {
		return nil, v.Error()
	}
	h, ok := v.RequireFloat64(p, "h")
	if !ok {
		return nil, v.Error()
	}
	text := v.OptionalString(p, "text")

	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	id, err := e.AddShape(slideIndex, shapeType, x, y, w, h)
	if err != nil {
		return nil, err
	}

	if text != "" {
		updates := common.ShapeUpdate{Text: &text}
		if updateErr := e.UpdateShape(slideIndex, id, updates); updateErr != nil {
			return nil, updateErr
		}
	}
	return map[string]int{"shape_id": id}, nil
}

func handleAddImage(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	path, ok := v.RequireString(p, "path")
	if !ok {
		return nil, v.Error()
	}
	x, ok := v.RequireFloat64(p, "x")
	if !ok {
		return nil, v.Error()
	}
	y, ok := v.RequireFloat64(p, "y")
	if !ok {
		return nil, v.Error()
	}
	w, ok := v.RequireFloat64(p, "w")
	if !ok {
		return nil, v.Error()
	}
	h, ok := v.RequireFloat64(p, "h")
	if !ok {
		return nil, v.Error()
	}

	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	id, err := e.AddImage(slideIndex, path, x, y, w, h)
	if err != nil {
		return nil, err
	}
	return map[string]int{"shape_id": id}, nil
}

func handleRemoveShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	shapeID, ok := v.RequireInt(p, "shape_id")
	if !ok {
		return nil, v.Error()
	}

	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	if err := e.RemoveShape(slideIndex, shapeID); err != nil {
		return nil, err
	}
	return map[string]bool{"removed": true}, nil
}

func handleUpdateShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	shapeID, ok := v.RequireInt(p, "shape_id")
	if !ok {
		return nil, v.Error()
	}

	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	// For embedded updates struct, we can unmarshal just that branch if it exists
	var updates common.ShapeUpdate
	if updatesRaw, ok := p["updates"]; ok {
		raw, _ := json.Marshal(updatesRaw)
		if err := json.Unmarshal(raw, &updates); err != nil {
			return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
		}
	}

	if err := e.UpdateShape(slideIndex, shapeID, updates); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleGetNotes(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

	notes, err := e.GetNotes(slideIndex)
	if err != nil {
		return nil, err
	}
	return map[string]string{"text": notes}, nil
}

func handleSetNotes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	text, ok := v.RequireString(p, "text")
	if !ok {
		return nil, v.Error()
	}
	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	if err := e.SetNotes(slideIndex, text); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}
