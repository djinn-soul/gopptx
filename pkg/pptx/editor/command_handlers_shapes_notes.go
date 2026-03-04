package editor

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
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
	if err := editorcommand.DecodeOptionalPayloadValue(p, "options", &opts); err != nil {
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

func handleGroupShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	shapeIDs, ok := v.RequireIntSlice(p, "shape_ids")
	if !ok {
		return nil, v.Error()
	}

	groupID, err := e.GroupShapes(slideIndex, shapeIDs)
	if err != nil {
		return nil, err
	}
	return map[string]int{"group_id": groupID}, nil
}

func handleUngroupShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

	shapeID, err = e.UngroupShapes(slideIndex, shapeID)
	if err != nil {
		return nil, err
	}
	return map[string]int{"group_id": shapeID}, nil
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
	if err := editorcommand.DecodeOptionalPayloadValue(p, "updates", &updates); err != nil {
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
		MetaKey:          "mime_type",
		PrimaryPathKey:   "path",
		PrimaryDataKey:   "data",
		SecondaryPathKey: "poster_path",
		SecondaryDataKey: "poster_data",
		PrimaryMaxLen:    maxMediaBase64,
		SecondaryMaxLen:  maxMediaBase64,
		PrimaryLabel:     "video",
		SecondaryLabel:   "poster",
		InsertBinary: func(
			placement editorcommand.MediaPlacement,
			mimeType string,
			videoData []byte,
			posterData []byte,
		) (int, error) {
			return e.AddVideo(
				placement.SlideIndex,
				videoData,
				posterData,
				mimeType,
				placement.X,
				placement.Y,
				placement.W,
				placement.H,
			)
		},
		InsertPath: func(
			placement editorcommand.MediaPlacement,
			mimeType string,
			videoPath string,
			posterPath string,
		) (int, error) {
			return e.AddVideoFromFile(
				placement.SlideIndex,
				videoPath,
				posterPath,
				mimeType,
				placement.X,
				placement.Y,
				placement.W,
				placement.H,
			)
		},
	})
}

func handleAddOLEObject(e *PresentationEditor, payload json.RawMessage) (any, error) {
	return handleMediaInsertCommand(e, payload, mediaInsertSpec{
		MetaKey:          "prog_id",
		PrimaryPathKey:   "path",
		PrimaryDataKey:   "data",
		SecondaryPathKey: "icon_path",
		SecondaryDataKey: "icon_data",
		PrimaryMaxLen:    maxEmbeddingBase64,
		SecondaryMaxLen:  maxEmbeddingBase64,
		PrimaryLabel:     "object",
		SecondaryLabel:   "icon",
		InsertBinary: func(
			placement editorcommand.MediaPlacement,
			progID string,
			objectData []byte,
			iconData []byte,
		) (int, error) {
			return e.AddOLEObject(
				placement.SlideIndex,
				objectData,
				iconData,
				progID,
				placement.X,
				placement.Y,
				placement.W,
				placement.H,
			)
		},
		InsertPath: func(
			placement editorcommand.MediaPlacement,
			progID string,
			objectPath string,
			iconPath string,
		) (int, error) {
			return e.AddOLEObjectFromFile(
				placement.SlideIndex,
				objectPath,
				iconPath,
				progID,
				placement.X,
				placement.Y,
				placement.W,
				placement.H,
			)
		},
	})
}

type mediaInsertSpec = editorcommand.MediaInsertSpec

func handleMediaInsertCommand(
	e *PresentationEditor,
	payload json.RawMessage,
	spec mediaInsertSpec,
) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	placement, ok := parseMediaPlacement(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	shapeID, err := editorcommand.ExecuteMediaInsert(p, placement, v.OptionalString, spec)
	if err != nil {
		return nil, err
	}
	return map[string]int{"shape_id": shapeID}, nil
}

func parseMediaPlacement(
	e *PresentationEditor,
	payload map[string]any,
	v *PayloadValidator,
) (editorcommand.MediaPlacement, bool) {
	slideIndex, ok := v.RequireInt(payload, "slide_index")
	if !ok {
		return editorcommand.MediaPlacement{}, false
	}
	x, ok := v.RequireFloat64(payload, "x")
	if !ok {
		return editorcommand.MediaPlacement{}, false
	}
	y, ok := v.RequireFloat64(payload, "y")
	if !ok {
		return editorcommand.MediaPlacement{}, false
	}
	w, ok := v.RequireFloat64(payload, "w")
	if !ok {
		return editorcommand.MediaPlacement{}, false
	}
	h, ok := v.RequireFloat64(payload, "h")
	if !ok {
		return editorcommand.MediaPlacement{}, false
	}
	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return editorcommand.MediaPlacement{}, false
	}
	return editorcommand.MediaPlacement{
		SlideIndex: slideIndex,
		X:          x,
		Y:          y,
		W:          w,
		H:          h,
	}, true
}

type addShapeRequest struct {
	slideIndex  int
	shapeType   string
	x           float64
	y           float64
	w           float64
	h           float64
	text        string
	textFrame   *common.TextFrame
	paragraph   *common.Paragraph
	clickAction *common.Hyperlink
	hoverAction *common.Hyperlink
	runs        []common.TextRun
	properties  common.ShapeUpdate
}

func parseAddShapeRequest(
	e *PresentationEditor,
	payload map[string]any,
	v *PayloadValidator,
) (addShapeRequest, error) {
	slideIndex, ok := requireSlideIndex(e, payload, v)
	if !ok {
		return addShapeRequest{}, v.Error()
	}
	shapeType, ok := v.RequireString(payload, "type")
	if !ok {
		return addShapeRequest{}, v.Error()
	}
	x, ok := v.RequireFloat64(payload, "x")
	if !ok {
		return addShapeRequest{}, v.Error()
	}
	y, ok := v.RequireFloat64(payload, "y")
	if !ok {
		return addShapeRequest{}, v.Error()
	}
	w, ok := v.RequireFloat64(payload, "w")
	if !ok {
		return addShapeRequest{}, v.Error()
	}
	h, ok := v.RequireFloat64(payload, "h")
	if !ok {
		return addShapeRequest{}, v.Error()
	}

	request := addShapeRequest{
		slideIndex: slideIndex,
		shapeType:  shapeType,
		x:          x,
		y:          y,
		w:          w,
		h:          h,
		text:       v.OptionalString(payload, "text"),
	}

	if err := editorcommand.DecodeOptionalPayloadValue(payload, "text_frame", &request.textFrame); err != nil {
		return addShapeRequest{}, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	if err := editorcommand.DecodeOptionalPayloadValue(payload, "paragraph", &request.paragraph); err != nil {
		return addShapeRequest{}, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	if err := editorcommand.DecodeOptionalPayloadValue(payload, "click_action", &request.clickAction); err != nil {
		return addShapeRequest{}, fmt.Errorf("invalid click_action: %w", err)
	}
	if err := editorcommand.DecodeOptionalPayloadValue(payload, "hover_action", &request.hoverAction); err != nil {
		return addShapeRequest{}, fmt.Errorf("invalid hover_action: %w", err)
	}
	if err := editorcommand.DecodeOptionalPayloadValue(payload, "runs", &request.runs); err != nil {
		return addShapeRequest{}, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	if err := editorcommand.DecodeOptionalPayloadValue(payload, "properties", &request.properties); err != nil {
		return addShapeRequest{}, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	return request, nil
}

func buildShapeUpdateForAdd(request addShapeRequest) (common.ShapeUpdate, bool) {
	hasExplicitUpdates := request.text != "" ||
		len(request.runs) > 0 ||
		request.textFrame != nil ||
		request.paragraph != nil ||
		request.clickAction != nil ||
		request.hoverAction != nil
	hasProperties := hasAnyUpdate(request.properties)

	if !hasExplicitUpdates && !hasProperties {
		return common.ShapeUpdate{}, false
	}

	updates := request.properties
	if request.text != "" {
		updates.Text = &request.text
	}
	if len(request.runs) > 0 {
		updates.Runs = &request.runs
	}
	if request.textFrame != nil {
		updates.TextFrame = request.textFrame
	}
	if request.paragraph != nil {
		updates.Paragraph = request.paragraph
	}
	if request.clickAction != nil {
		updates.ClickAction = request.clickAction
	}
	if request.hoverAction != nil {
		updates.HoverAction = request.hoverAction
	}
	return updates, true
}

func hasAnyUpdate(u common.ShapeUpdate) bool {
	return u.Text != nil || u.Runs != nil || u.TextFrame != nil ||
		u.Paragraph != nil || u.Fill != nil || u.Line != nil || u.Shadow != nil || u.Glow != nil || u.Blur != nil || u.SoftEdge != nil || u.Reflection != nil ||
		u.ClickAction != nil || u.HoverAction != nil || u.X != nil ||
		u.Y != nil || u.W != nil || u.H != nil || u.Rotation != nil ||
		u.FlipH != nil || u.FlipV != nil || u.Crop != nil
}
