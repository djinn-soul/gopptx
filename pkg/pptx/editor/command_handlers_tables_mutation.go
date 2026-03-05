package editor

import (
	"encoding/json"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

func handleMergeTableCells(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.TableCellRangeRequest, bool) {
			return editorcommand.ParseTableCellRangeRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.TableCellRangeRequest) (any, error) {
			if err := e.MergeTableCells(
				request.SlideIndex,
				request.ShapeID,
				request.Row1,
				request.Col1,
				request.Row2,
				request.Col2,
			); err != nil {
				return nil, err
			}
			return map[string]bool{"success": true}, nil
		},
	)
}

func handleSplitTableCell(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.TableCellRequest, bool) {
			return editorcommand.ParseTableCellRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.TableCellRequest) (any, error) {
			if err := e.SplitTableCell(request.SlideIndex, request.ShapeID, request.Row, request.Col); err != nil {
				return nil, err
			}
			return map[string]bool{"success": true}, nil
		},
	)
}

func handleUpdateTableFlags(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.TableShapeRequest, bool) {
			return editorcommand.ParseTableShapeRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.TableShapeRequest, p map[string]any) (any, error) {
			flags, err := editorcommand.ParseRequiredObjectField(
				p,
				"flags",
				"missing flags map",
				"flags must be an object",
			)
			if err != nil {
				return nil, err
			}
			if err := e.UpdateTableFlags(request.SlideIndex, request.ShapeID, flags); err != nil {
				return nil, err
			}
			return map[string]bool{"success": true}, nil
		},
	)
}

func handleUpdateTableCell(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.TableCellRequest, bool) {
			return editorcommand.ParseTableCellRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.TableCellRequest, p map[string]any) (any, error) {
			updates, err := editorcommand.ParseRequiredObjectField(
				p,
				"updates",
				"missing updates map",
				"updates must be an object",
			)
			if err != nil {
				return nil, err
			}
			text, hasText, err := editorcommand.ParseOptionalTextUpdate(updates)
			if err != nil {
				return nil, err
			}
			if hasText {
				if err := e.UpdateTableCellText(
					request.SlideIndex,
					request.ShapeID,
					request.Row,
					request.Col,
					text,
				); err != nil {
					return nil, err
				}
			}
			return map[string]bool{"success": true}, nil
		},
	)
}

func handleSetTableStyle(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.TableStyleRequest, bool) {
			return editorcommand.ParseTableStyleRequest(p, v.RequireInt, v.RequireString)
		},
		v.Error,
		func(request editorcommand.TableStyleRequest) (any, error) {
			if err := e.SetTableStyle(request.SlideIndex, request.ShapeID, request.StyleGUID); err != nil {
				return nil, err
			}
			return map[string]bool{"success": true}, nil
		},
	)
}
