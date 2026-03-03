package editor

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	maxImageBase64     = 50 * 1024 * 1024
	maxMediaBase64     = 50 * 1024 * 1024
	maxEmbeddingBase64 = 20 * 1024 * 1024
)

func handleListShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
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
	request, err := parseAddShapeRequest(e, p, v)
	if err != nil {
		return nil, err
	}

	id, err := e.AddShape(request.slideIndex, request.shapeType, request.x, request.y, request.w, request.h)
	if err != nil {
		return nil, err
	}

	if updates, hasUpdates := buildShapeUpdateForAdd(request); hasUpdates {
		if updateErr := e.UpdateShape(request.slideIndex, id, updates); updateErr != nil {
			return nil, updateErr
		}
	}
	return map[string]int{"shape_id": id}, nil
}

func handleGetImageMetadata(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

	meta, err := e.GetImageMetadata(slideIndex, shapeID)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func handleAddImage(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
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

	imagePath := v.OptionalString(p, "path")
	base64Data := v.OptionalString(p, "data")
	format := v.OptionalString(p, "format")

	var opts *common.ShapeUpdate
	if err := decodeOptionalPayloadValue(p, "options", &opts); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}

	var newID int
	if base64Data != "" {
		if strings.TrimSpace(format) == "" {
			return nil, errors.New("image format is required when image data is provided")
		}
		if len(base64Data) > maxImageBase64 {
			return nil, errors.New("image data too large")
		}
		decodedData, decodeErr := base64.StdEncoding.DecodeString(base64Data)
		if decodeErr != nil {
			return nil, fmt.Errorf("invalid base64: %w", decodeErr)
		}
		newID, err = e.AddImageFromBytes(slideIndex, decodedData, format, x, y, w, h, opts)
	} else {
		newID, err = e.AddImage(slideIndex, imagePath, x, y, w, h, opts)
	}

	if err != nil {
		return nil, err
	}
	return map[string]int{"shape_id": newID}, nil
}

func handleRemoveShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

	if err := e.RemoveShape(slideIndex, shapeID); err != nil {
		return nil, err
	}
	return map[string]bool{"removed": true}, nil
}

func handleMoveShapeToFront(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

	if err := e.MoveShapeToFront(slideIndex, shapeID); err != nil {
		return nil, err
	}
	return map[string]bool{"moved": true}, nil
}

func handleMoveShapeToBack(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

	if err := e.MoveShapeToBack(slideIndex, shapeID); err != nil {
		return nil, err
	}
	return map[string]bool{"moved": true}, nil
}

func handleUpdateShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

	var updates common.ShapeUpdate
	if err := decodeOptionalPayloadValue(p, "updates", &updates); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
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
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	notes, err := e.GetNotes(slideIndex)
	if err != nil {
		return nil, err
	}
	hasNotesSlide, err := e.HasNotesSlide(slideIndex)
	if err != nil {
		return nil, err
	}
	var notesSlide any
	if hasNotesSlide {
		notesSlide = map[string]string{"text": notes}
	}
	return map[string]any{
		"text":        notes,
		"notes_slide": notesSlide,
	}, nil
}

func handleHasNotesSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	hasNotesSlide, err := e.HasNotesSlide(slideIndex)
	if err != nil {
		return nil, err
	}
	return map[string]bool{"has_notes_slide": hasNotesSlide}, nil
}

func handleSetNotes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	text, ok := v.RequireString(p, "text")
	if !ok {
		return nil, v.Error()
	}
	if err := e.SetNotes(slideIndex, text); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleAddVideo(e *PresentationEditor, payload json.RawMessage) (any, error) {
	return handleMediaInsertCommand(e, payload, mediaInsertSpec{
		metaKey:          "mime_type",
		primaryPathKey:   "path",
		primaryDataKey:   "data",
		secondaryPathKey: "poster_path",
		secondaryDataKey: "poster_data",
		primaryMaxLen:    maxMediaBase64,
		secondaryMaxLen:  maxMediaBase64,
		primaryLabel:     "video",
		secondaryLabel:   "poster",
		insertBinary: func(
			placement mediaPlacement,
			mimeType string,
			videoData []byte,
			posterData []byte,
		) (int, error) {
			return e.AddVideo(
				placement.slideIndex,
				videoData,
				posterData,
				mimeType,
				placement.x,
				placement.y,
				placement.w,
				placement.h,
			)
		},
		insertPath: func(
			placement mediaPlacement,
			mimeType string,
			videoPath string,
			posterPath string,
		) (int, error) {
			return e.AddVideoFromFile(
				placement.slideIndex,
				videoPath,
				posterPath,
				mimeType,
				placement.x,
				placement.y,
				placement.w,
				placement.h,
			)
		},
	})
}

func handleAddOLEObject(e *PresentationEditor, payload json.RawMessage) (any, error) {
	return handleMediaInsertCommand(e, payload, mediaInsertSpec{
		metaKey:          "prog_id",
		primaryPathKey:   "path",
		primaryDataKey:   "data",
		secondaryPathKey: "icon_path",
		secondaryDataKey: "icon_data",
		primaryMaxLen:    maxEmbeddingBase64,
		secondaryMaxLen:  maxEmbeddingBase64,
		primaryLabel:     "object",
		secondaryLabel:   "icon",
		insertBinary: func(
			placement mediaPlacement,
			progID string,
			objectData []byte,
			iconData []byte,
		) (int, error) {
			return e.AddOLEObject(
				placement.slideIndex,
				objectData,
				iconData,
				progID,
				placement.x,
				placement.y,
				placement.w,
				placement.h,
			)
		},
		insertPath: func(
			placement mediaPlacement,
			progID string,
			objectPath string,
			iconPath string,
		) (int, error) {
			return e.AddOLEObjectFromFile(
				placement.slideIndex,
				objectPath,
				iconPath,
				progID,
				placement.x,
				placement.y,
				placement.w,
				placement.h,
			)
		},
	})
}

func decodeOptionalPayloadValue(payload map[string]any, key string, target any) error {
	rawValue, ok := payload[key]
	if !ok || rawValue == nil {
		return nil
	}
	raw, err := json.Marshal(rawValue)
	if err != nil {
		return fmt.Errorf("invalid %s structure: %w", key, err)
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("invalid %s payload: %w", key, err)
	}
	return nil
}
