package editor

import (
	"encoding/json"
	"errors"
	"math"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

const minConnectorDimension = 1.0

func handleAddTextbox(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.TextboxPlacementRequest, bool) {
			return editorcommand.ParseTextboxPlacementRequest(
				p,
				func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
				v.RequireFloat64,
			)
		},
		v.Error,
		func(request editorcommand.TextboxPlacementRequest, p map[string]any) (any, error) {
			addPayload := map[string]any{
				"slide_index": request.SlideIndex,
				"type":        "rect",
				"x":           request.Left,
				"y":           request.Top,
				"w":           request.Width,
				"h":           request.Height,
			}
			editorcommand.CopyShapeUpdateFields(p, addPayload)
			return addShapeFromPayload(e, addPayload)
		},
	)
}

func handleAddConnector(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.ConnectorPlacementRequest, bool) {
			return editorcommand.ParseConnectorPlacementRequest(
				p,
				func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
				v.RequireString,
				v.RequireFloat64,
			)
		},
		v.Error,
		func(request editorcommand.ConnectorPlacementRequest, p map[string]any) (any, error) {
			left := math.Min(request.BeginX, request.EndX)
			top := math.Min(request.BeginY, request.EndY)
			width := math.Max(math.Abs(request.EndX-request.BeginX), minConnectorDimension)
			height := math.Max(math.Abs(request.EndY-request.BeginY), minConnectorDimension)

			addPayload := map[string]any{
				"slide_index": request.SlideIndex,
				"type":        request.ConnectorType,
				"x":           left,
				"y":           top,
				"w":           width,
				"h":           height,
			}
			editorcommand.CopyShapeUpdateFields(p, addPayload)
			return addShapeFromPayload(e, addPayload)
		},
	)
}

func handleAddTextboxes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (int, bool) { return requireSlideIndex(e, p, v) },
		v.Error,
		func(slideIndex int, p map[string]any) (any, error) {
			if _, ok := p["textboxes"]; !ok {
				return nil, NewBridgeError(ErrCodeInvalidPayload, "missing required field: textboxes")
			}
			var textboxes []common.TextboxInsert
			if err := editorcommand.DecodeOptionalPayloadValue(p, "textboxes", &textboxes); err != nil {
				return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
			}
			shapeIDs, err := e.AddTextboxes(slideIndex, textboxes)
			if err != nil {
				return nil, err
			}
			return map[string]any{"shape_ids": shapeIDs}, nil
		},
	)
}

func handleReserveShapeIDs(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (int, bool) { return requireSlideIndex(e, p, v) },
		v.Error,
		func(slideIndex int, p map[string]any) (any, error) {
			count, ok := v.RequireInt(p, "count")
			if !ok {
				return nil, v.Error()
			}
			shapeIDs, err := e.ReserveShapeIDs(slideIndex, count)
			if err != nil {
				return nil, err
			}
			return map[string]any{"shape_ids": shapeIDs}, nil
		},
	)
}

func handleAddGroupShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.GroupShapeRequest, bool) {
			return editorcommand.ParseGroupShapeRequest(p, v.RequireInt, v.RequireIntSlice)
		},
		v.Error,
		func(request editorcommand.GroupShapeRequest) (any, error) {
			newID, err := e.AddGroupShape(request.SlideIndex, request.ShapeIDs)
			if err != nil {
				return nil, err
			}
			return map[string]int{"shape_id": newID}, nil
		},
	)
}

func handleBuildFreeform(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	points, err := editorcommand.ParseFreeformPoints(p)
	if err != nil {
		return nil, bridgeValidationError(err)
	}
	closePath, err := editorcommand.ParseOptionalCloseFlag(p)
	if err != nil {
		return nil, bridgeValidationError(err)
	}
	freeformPoints := make([]freeformPoint, 0, len(points))
	for _, point := range points {
		freeformPoints = append(freeformPoints, freeformPoint{X: point.X, Y: point.Y})
	}
	shapeID, err := e.AddFreeformShape(slideIndex, freeformPoints, closePath)
	if err != nil {
		return nil, err
	}
	if updates, hasUpdates, updateErr := editorcommand.ParseOptionalShapeUpdates(p); updateErr != nil {
		return nil, bridgeValidationError(updateErr)
	} else if hasUpdates {
		if err := e.UpdateShape(slideIndex, shapeID, updates); err != nil {
			return nil, err
		}
	}
	return map[string]int{"shape_id": shapeID}, nil
}

func addShapeFromPayload(e *PresentationEditor, payload map[string]any) (any, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return handleAddShape(e, raw)
}

func bridgeValidationError(err error) error {
	var validationErr *editorcommand.ValidationError
	if errors.As(err, &validationErr) {
		return NewBridgeError(validationErr.Code, validationErr.Message)
	}
	return err
}
