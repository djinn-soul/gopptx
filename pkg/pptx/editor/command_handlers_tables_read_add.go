package editor

import (
	"encoding/json"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

const (
	maxTableDimension = 1000
)

func handleAddTable(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.TableAddRequest, bool) {
			return editorcommand.ParseTableAddRequest(
				p,
				func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
				v.RequireInt,
				v.OptionalInt64,
			)
		},
		v.Error,
		func(request editorcommand.TableAddRequest) (any, error) {
			if err := editorcommand.ValidateTableDimensions(request.Rows, request.Cols, maxTableDimension); err != nil {
				return nil, err
			}
			shapeID, err := e.AddTable(
				request.SlideIndex,
				request.Rows,
				request.Cols,
				request.X,
				request.Y,
				request.CX,
				request.CY,
			)
			if err != nil {
				return nil, err
			}
			return map[string]int{"shape_id": shapeID}, nil
		},
	)
}

func handleGetTable(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.TableShapeRequest, bool) {
			return editorcommand.ParseTableShapeRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.TableShapeRequest) (any, error) {
			return e.GetTable(request.SlideIndex, request.ShapeID)
		},
	)
}
