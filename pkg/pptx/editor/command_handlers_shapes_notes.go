package editor

import (
	"encoding/json"
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

const (
	maxImageBase64     = 50 * 1024 * 1024
	maxMediaBase64     = 50 * 1024 * 1024
	maxEmbeddingBase64 = 20 * 1024 * 1024
	maxTableDimension  = 1000
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
	newID, err := editorcommand.ExecuteAddImageRequest(request, maxImageBase64, e.AddImageFromBytes, e.AddImage)
	if err != nil {
		return nil, err
	}
	return map[string]int{"shape_id": newID}, nil
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

func handleGetNotes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideIndexRequest(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.Error,
		func(slideIndex int) (any, error) {
			notes, err := e.GetNotes(slideIndex)
			if err != nil {
				return nil, err
			}
			hasNotesSlide, err := e.HasNotesSlide(slideIndex)
			if err != nil {
				return nil, err
			}
			return editorcommand.BuildNotesResult(notes, hasNotesSlide), nil
		},
	)
}

func handleHasNotesSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideIndexRequest(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.Error,
		func(slideIndex int) (any, error) {
			hasNotesSlide, err := e.HasNotesSlide(slideIndex)
			if err != nil {
				return nil, err
			}
			return map[string]bool{"has_notes_slide": hasNotesSlide}, nil
		},
	)
}

func handleSetNotes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, text, ok := editorcommand.ParseSetNotesRequest(
		p,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.RequireString,
	)
	if !ok {
		return nil, v.Error()
	}
	if err := e.SetNotes(slideIndex, text); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleAddTable(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	rowCount, ok := v.RequireInt(p, "rows")
	if !ok {
		return nil, v.Error()
	}
	colCount, ok := v.RequireInt(p, "cols")
	if !ok {
		return nil, v.Error()
	}

	if rowCount < 1 || rowCount > maxTableDimension {
		return nil, fmt.Errorf("rows %d must be between 1 and %d", rowCount, maxTableDimension)
	}
	if colCount < 1 || colCount > maxTableDimension {
		return nil, fmt.Errorf("cols %d must be between 1 and %d", colCount, maxTableDimension)
	}

	x, _ := v.OptionalInt64(p, "x")
	y, _ := v.OptionalInt64(p, "y")
	cx, _ := v.OptionalInt64(p, "cx")
	cy, _ := v.OptionalInt64(p, "cy")

	if v.HasErrors() {
		return nil, v.Error()
	}

	shapeID, err := e.AddTable(slideIndex, rowCount, colCount, x, y, cx, cy)
	if err != nil {
		return nil, err
	}

	return map[string]int{"shape_id": shapeID}, nil
}

func handleGetTable(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

	return e.GetTable(slideIndex, shapeID)
}

func handleMergeTableCells(e *PresentationEditor, payload json.RawMessage) (any, error) {
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
	row1, ok := v.RequireInt(p, "row1")
	if !ok {
		return nil, v.Error()
	}
	col1, ok := v.RequireInt(p, "col1")
	if !ok {
		return nil, v.Error()
	}
	row2, ok := v.RequireInt(p, "row2")
	if !ok {
		return nil, v.Error()
	}
	col2, ok := v.RequireInt(p, "col2")
	if !ok {
		return nil, v.Error()
	}

	if err := e.MergeTableCells(slideIndex, shapeID, row1, col1, row2, col2); err != nil {
		return nil, err
	}
	return map[string]bool{"success": true}, nil
}

func handleSplitTableCell(e *PresentationEditor, payload json.RawMessage) (any, error) {
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
	row, ok := v.RequireInt(p, "row")
	if !ok {
		return nil, v.Error()
	}
	col, ok := v.RequireInt(p, "col")
	if !ok {
		return nil, v.Error()
	}

	if err := e.SplitTableCell(slideIndex, shapeID, row, col); err != nil {
		return nil, err
	}
	return map[string]bool{"success": true}, nil
}

func handleUpdateTableFlags(e *PresentationEditor, payload json.RawMessage) (any, error) {
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
	flagsVal, ok := p["flags"]
	if !ok {
		return nil, errors.New("missing flags map")
	}
	flags, ok := flagsVal.(map[string]any)
	if !ok {
		return nil, errors.New("flags must be an object")
	}

	if err := e.UpdateTableFlags(slideIndex, shapeID, flags); err != nil {
		return nil, err
	}

	return map[string]bool{"success": true}, nil
}

func handleUpdateTableCell(e *PresentationEditor, payload json.RawMessage) (any, error) {
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
	row, ok := v.RequireInt(p, "row")
	if !ok {
		return nil, v.Error()
	}
	col, ok := v.RequireInt(p, "col")
	if !ok {
		return nil, v.Error()
	}
	updatesVal, ok := p["updates"]
	if !ok {
		return nil, errors.New("missing updates map")
	}
	updates, ok := updatesVal.(map[string]any)
	if !ok {
		return nil, errors.New("updates must be an object")
	}

	if text, hasText := updates["text"]; hasText {
		textStr, isStr := text.(string)
		if !isStr {
			return nil, errors.New("text update must be a string")
		}
		if err := e.UpdateTableCellText(slideIndex, shapeID, row, col, textStr); err != nil {
			return nil, err
		}
	}

	return map[string]bool{"success": true}, nil
}

func handleSetTableStyle(e *PresentationEditor, payload json.RawMessage) (any, error) {
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
	styleGUID, ok := v.RequireString(p, "style_guid")
	if !ok {
		return nil, v.Error()
	}

	if err := e.SetTableStyle(slideIndex, shapeID, styleGUID); err != nil {
		return nil, err
	}

	return map[string]bool{"success": true}, nil
}

func handleAddVideo(e *PresentationEditor, payload json.RawMessage) (any, error) {
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
		func(shapeID int) any { return map[string]int{"shape_id": shapeID} },
		spec,
	)
}
