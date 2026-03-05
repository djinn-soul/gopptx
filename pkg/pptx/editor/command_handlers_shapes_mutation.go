package editor

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

func executeSlideShapeMutation(
	e *PresentationEditor,
	payload json.RawMessage,
	resultKey string,
	mutate func(slideIndex, shapeID int) error,
) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideShapeRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (int, bool) { return requireSlideIndex(e, p, v) },
		v.RequireInt,
		v.Error,
		func(request editorcommand.SlideShapeRequest) (any, error) {
			if err := mutate(request.SlideIndex, request.ShapeID); err != nil {
				return nil, err
			}
			return map[string]bool{resultKey: true}, nil
		},
	)
}

func handleRemoveShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	return executeSlideShapeMutation(
		e,
		payload,
		"removed",
		func(slideIndex, shapeID int) error {
			return e.RemoveShape(slideIndex, shapeID)
		},
	)
}

func handleGroupShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideShapeIDsRequest(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.RequireIntSlice,
		v.Error,
		func(request editorcommand.SlideShapeIDsRequest) (any, error) {
			groupID, err := e.GroupShapes(request.SlideIndex, request.ShapeIDs)
			if err != nil {
				return nil, err
			}
			return map[string]int{"group_id": groupID}, nil
		},
	)
}

func handleUngroupShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideShapeRequest(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.RequireInt,
		v.Error,
		func(request editorcommand.SlideShapeRequest) (any, error) {
			shapeID, err := e.UngroupShapes(request.SlideIndex, request.ShapeID)
			if err != nil {
				return nil, err
			}
			return map[string]int{"group_id": shapeID}, nil
		},
	)
}

func handleMoveShapeToFront(e *PresentationEditor, payload json.RawMessage) (any, error) {
	return executeSlideShapeMutation(
		e,
		payload,
		"moved",
		func(slideIndex, shapeID int) error {
			return e.MoveShapeToFront(slideIndex, shapeID)
		},
	)
}

func handleMoveShapeToBack(e *PresentationEditor, payload json.RawMessage) (any, error) {
	return executeSlideShapeMutation(
		e,
		payload,
		"moved",
		func(slideIndex, shapeID int) error {
			return e.MoveShapeToBack(slideIndex, shapeID)
		},
	)
}

func handleUpdateShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideShapeRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.RequireInt,
		v.Error,
		func(request editorcommand.SlideShapeRequest, p map[string]any) (any, error) {
			var updates common.ShapeUpdate
			if err := editorcommand.DecodeOptionalPayloadValue(p, "updates", &updates); err != nil {
				return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
			}
			if err := e.UpdateShape(request.SlideIndex, request.ShapeID, updates); err != nil {
				return nil, err
			}
			return map[string]bool{"updated": true}, nil
		},
	)
}
