package editor

import (
	"encoding/base64"
	"encoding/json"
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
	var textFrame *common.TextFrame
	if err := decodeOptionalPayloadValue(p, "text_frame", &textFrame); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	var clickAction *common.Hyperlink
	if err := decodeOptionalPayloadValue(p, "click_action", &clickAction); err != nil {
		return nil, fmt.Errorf("invalid click_action: %w", err)
	}

	var hoverAction *common.Hyperlink
	if err := decodeOptionalPayloadValue(p, "hover_action", &hoverAction); err != nil {
		return nil, fmt.Errorf("invalid hover_action: %w", err)
	}

	var runs []common.TextRun
	if err := decodeOptionalPayloadValue(p, "runs", &runs); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}

	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	id, err := e.AddShape(slideIndex, shapeType, x, y, w, h)
	if err != nil {
		return nil, err
	}

	// Create properties update if provided
	if text != "" || len(runs) > 0 || textFrame != nil || clickAction != nil || hoverAction != nil {
		updates := common.ShapeUpdate{}
		if text != "" {
			updates.Text = &text
		}
		if len(runs) > 0 {
			updates.Runs = &runs
		}
		if textFrame != nil {
			updates.TextFrame = textFrame
		}
		if clickAction != nil {
			updates.ClickAction = clickAction
		}
		if hoverAction != nil {
			updates.HoverAction = hoverAction
		}
		if updateErr := e.UpdateShape(slideIndex, id, updates); updateErr != nil {
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
	slideIndex, ok := v.RequireInt(p, "slide_index")
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
			return nil, fmt.Errorf("image format is required when image data is provided")
		}
		if len(base64Data) > maxImageBase64 {
			return nil, fmt.Errorf("image data too large")
		}
		data, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return nil, fmt.Errorf("invalid base64: %w", err)
		}
		newID, err = e.AddImageFromBytes(slideIndex, data, format, x, y, w, h, opts)
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

func handleMoveShapeToFront(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

func handleAddVideo(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
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

	mimeType := v.OptionalString(p, "mime_type")

	videoPath := v.OptionalString(p, "path")
	videoBase64 := v.OptionalString(p, "data")
	posterPath := v.OptionalString(p, "poster_path")
	posterBase64 := v.OptionalString(p, "poster_data")

	var videoData, posterData []byte
	if videoBase64 != "" {
		if len(videoBase64) > maxMediaBase64 {
			return nil, fmt.Errorf("video data too large")
		}
		videoData, err = base64.StdEncoding.DecodeString(videoBase64)
		if err != nil {
			return nil, fmt.Errorf("invalid video base64: %w", err)
		}
	}
	if posterBase64 != "" {
		if len(posterBase64) > maxMediaBase64 {
			return nil, fmt.Errorf("poster data too large")
		}
		posterData, err = base64.StdEncoding.DecodeString(posterBase64)
		if err != nil {
			return nil, fmt.Errorf("invalid poster base64: %w", err)
		}
	}

	var newID int
	if len(videoData) > 0 || len(posterData) > 0 {
		newID, err = e.AddVideo(slideIndex, videoData, posterData, mimeType, x, y, w, h)
	} else {
		newID, err = e.AddVideoFromFile(slideIndex, videoPath, posterPath, mimeType, x, y, w, h)
	}

	if err != nil {
		return nil, err
	}
	return map[string]int{"shape_id": newID}, nil
}

func handleAddOLEObject(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
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

	progID := v.OptionalString(p, "prog_id")

	objectPath := v.OptionalString(p, "path")
	objectBase64 := v.OptionalString(p, "data")
	iconPath := v.OptionalString(p, "icon_path")
	iconBase64 := v.OptionalString(p, "icon_data")

	var objectData, iconData []byte
	if objectBase64 != "" {
		if len(objectBase64) > maxEmbeddingBase64 {
			return nil, fmt.Errorf("object data too large")
		}
		objectData, err = base64.StdEncoding.DecodeString(objectBase64)
		if err != nil {
			return nil, fmt.Errorf("invalid object base64: %w", err)
		}
	}
	if iconBase64 != "" {
		if len(iconBase64) > maxEmbeddingBase64 {
			return nil, fmt.Errorf("icon data too large")
		}
		iconData, err = base64.StdEncoding.DecodeString(iconBase64)
		if err != nil {
			return nil, fmt.Errorf("invalid icon base64: %w", err)
		}
	}

	var newID int
	if len(objectData) > 0 || len(iconData) > 0 {
		newID, err = e.AddOLEObject(slideIndex, objectData, iconData, progID, x, y, w, h)
	} else {
		newID, err = e.AddOLEObjectFromFile(slideIndex, objectPath, iconPath, progID, x, y, w, h)
	}

	if err != nil {
		return nil, err
	}
	return map[string]int{"shape_id": newID}, nil
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
