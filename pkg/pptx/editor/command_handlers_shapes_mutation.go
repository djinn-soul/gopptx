package editor

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

const (
	maxImageBase64     = 50 * 1024 * 1024
	maxMediaBase64     = 50 * 1024 * 1024
	maxEmbeddingBase64 = 20 * 1024 * 1024
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
	return respShapeID(id), nil
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
	return respShapeID(newID), nil
}

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

func handleClearShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideIndexRequest(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.Error,
		func(slideIndex int) (any, error) {
			if err := e.ClearShapes(slideIndex); err != nil {
				return nil, err
			}
			return respCleared, nil
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
			return respGroupID(groupID), nil
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
			return respGroupID(shapeID), nil
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
			return respUpdated, nil
		},
	)
}

func handleAddVideo(e *PresentationEditor, payload json.RawMessage) (any, error) {
	if shouldUseMediaPlaybackCommand(payload) {
		return handleAddVideoWithPlaybackCommand(e, payload)
	}
	return handleMediaInsertCommand(
		e,
		payload,
		editorcommand.NewVideoInsertSpec(
			maxMediaBase64,
			editorcommand.AdaptVideoBinaryInsert(e.AddVideo),
			editorcommand.AdaptVideoPathInsert(e.AddVideoFromFile),
		),
	)
}

func handleAddAudio(e *PresentationEditor, payload json.RawMessage) (any, error) {
	if shouldUseMediaPlaybackCommand(payload) {
		return handleAddAudioWithPlaybackCommand(e, payload)
	}
	return handleMediaInsertCommand(
		e,
		payload,
		editorcommand.NewAudioInsertSpec(
			maxMediaBase64,
			editorcommand.AdaptAudioBinaryInsertWithOptionalIcon(e.AddAudio, e.AddAudioWithIcon),
			editorcommand.AdaptAudioPathInsertWithOptionalIcon(e.AddAudioFromFile, e.AddAudioWithIconFromFile),
		),
	)
}

func handleAddOLEObject(e *PresentationEditor, payload json.RawMessage) (any, error) {
	return handleMediaInsertCommand(
		e,
		payload,
		editorcommand.NewOLEInsertSpec(
			maxEmbeddingBase64,
			editorcommand.AdaptOLEBinaryInsert(e.AddOLEObject),
			editorcommand.AdaptOLEPathInsert(e.AddOLEObjectFromFile),
		),
	)
}

type mediaInsertSpec = editorcommand.MediaInsertSpec

func handleMediaInsertCommand(
	e *PresentationEditor,
	payload json.RawMessage,
	spec mediaInsertSpec,
) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleMediaInsertCommand(
		payload,
		e.SlideCount(),
		parseRawPayloadBytes,
		v.RequireInt,
		v.RequireFloat64,
		v.IndexBounds,
		v.OptionalString,
		v.Error,
		func(shapeID int) any { return respShapeID(shapeID) },
		spec,
	)
}
