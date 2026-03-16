package editor

import (
	"encoding/base64"
	"encoding/json"
)

// handleMoveShapeToIndex reorders a shape to a specific z-index within its slide.
//
// Payload: {"slide_index": N, "shape_id": M, "target_index": T}.
// Response: {"moved": true}.
func handleMoveShapeToIndex(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	shapeID, ok := v.RequireInt(p, "shape_id")
	if !ok {
		return nil, v.Error()
	}
	targetIndex, ok := v.RequireInt(p, "target_index")
	if !ok {
		return nil, v.Error()
	}

	if err := e.MoveShapeToIndex(slideIndex, shapeID, targetIndex); err != nil {
		return nil, err
	}
	return map[string]bool{"moved": true}, nil
}

// handleListSlideImages lists all images embedded in a slide.
//
// Payload: {"slide_index": N}.
// Response: {"images": [...]}.
func handleListSlideImages(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	images, err := e.ListSlideImages(slideIndex)
	if err != nil {
		return nil, err
	}
	return map[string]any{"images": images}, nil
}

// handleSwapImageByIndex replaces an image at a given index within a slide.
//
// Payload: {"slide_index": N, "image_index": I, "data": "<base64>", "format": "<string>"}.
// Response: {"swapped": true}.
//
//nolint:dupl // handleSwapImageByIndex and handleSwapImageByRelID differ by selector field (image_index vs rel_id).
func handleSwapImageByIndex(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	imageIndex, ok := v.RequireInt(p, "image_index")
	if !ok {
		return nil, v.Error()
	}
	b64, ok := v.RequireString(p, "data")
	if !ok {
		return nil, v.Error()
	}
	format, ok := v.RequireString(p, "format")
	if !ok {
		return nil, v.Error()
	}

	data, decErr := base64.StdEncoding.DecodeString(b64)
	if decErr != nil {
		return nil, NewBridgeError(ErrCodeInvalidType, "data must be a valid base64 string")
	}

	if err := e.SwapImageByIndex(slideIndex, imageIndex, data, format); err != nil {
		return nil, err
	}
	return map[string]bool{"swapped": true}, nil
}

// handleSwapImageByRelID replaces an image identified by its relationship ID.
//
// Payload: {"slide_index": N, "rel_id": "<string>", "data": "<base64>", "format": "<string>"}.
// Response: {"swapped": true}.
//
//nolint:dupl // handleSwapImageByIndex and handleSwapImageByRelID differ by selector field (image_index vs rel_id).
func handleSwapImageByRelID(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	relID, ok := v.RequireString(p, "rel_id")
	if !ok {
		return nil, v.Error()
	}
	b64, ok := v.RequireString(p, "data")
	if !ok {
		return nil, v.Error()
	}
	format, ok := v.RequireString(p, "format")
	if !ok {
		return nil, v.Error()
	}

	data, decErr := base64.StdEncoding.DecodeString(b64)
	if decErr != nil {
		return nil, NewBridgeError(ErrCodeInvalidType, "data must be a valid base64 string")
	}

	if err := e.SwapImageByRelID(slideIndex, relID, data, format); err != nil {
		return nil, err
	}
	return map[string]bool{"swapped": true}, nil
}
