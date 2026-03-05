package editor

import (
	"encoding/json"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

const (
	maxImageBase64 = 50 * 1024 * 1024
)

func handleListShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideIndexRequest(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.Error,
		func(slideIndex int) (any, error) {
			shapes, err := e.GetShapes(slideIndex)
			if err != nil {
				return nil, err
			}
			return map[string]any{"shapes": shapes}, nil
		},
	)
}

func parseRawPayloadBytes(raw []byte) (map[string]any, error) {
	return ParseRawPayload(raw)
}

func handleAddShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	request, ok := editorcommand.ParseAddShapeBase(
		p,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.RequireString,
		v.RequireFloat64,
		v.OptionalString,
	)
	if !ok {
		return nil, v.Error()
	}
	if err := editorcommand.DecodeAddShapeOptionals(p, &request); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	id, err := editorcommand.ExecuteAddShapeRequest(request, e.AddShape, e.UpdateShape)
	if err != nil {
		return nil, err
	}
	return map[string]int{"shape_id": id}, nil
}

func handleGetImageMetadata(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideShapeRequest(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.RequireInt,
		v.Error,
		func(request editorcommand.SlideShapeRequest) (any, error) {
			return e.GetImageMetadata(request.SlideIndex, request.ShapeID)
		},
	)
}

func handleAddImage(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	request, ok, parseErr := editorcommand.ParseAddImageRequest(
		p,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.RequireFloat64,
		v.OptionalString,
	)
	if parseErr != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, parseErr.Error())
	}
	if !ok {
		return nil, v.Error()
	}
	newID, err := editorcommand.ExecuteAddImageRequest(
		request,
		maxImageBase64,
		e.AddImageFromBytes,
		e.AddImageFromURL,
		e.AddImage,
	)
	if err != nil {
		return nil, err
	}
	return map[string]int{"shape_id": newID}, nil
}
