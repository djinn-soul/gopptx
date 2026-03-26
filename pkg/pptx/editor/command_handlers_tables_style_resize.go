package editor

import (
	"encoding/json"
	"errors"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	tablemod "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/table"
)

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

func handleDefineTableStyle(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	name, ok := v.RequireString(p, "name")
	if !ok {
		return nil, v.Error()
	}
	styleID := v.OptionalString(p, "style_id")
	id, err := e.DefineTableStyle(common.TableStyleDefinition{
		StyleID: styleID,
		Name:    name,
	})
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"style_id": id,
		"name":     name,
	}, nil
}

func handleListTableStyles(e *PresentationEditor, _ json.RawMessage) (any, error) {
	styles, err := e.ListTableStyles()
	if err != nil {
		return nil, err
	}
	out := make([]map[string]string, 0, len(styles))
	for _, style := range styles {
		out = append(out, map[string]string{
			"style_id": style.StyleID,
			"name":     style.Name,
		})
	}
	return map[string]any{"styles": out}, nil
}

func handleSetTableRowHeight(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.TableShapeRequest, bool) {
			return editorcommand.ParseTableShapeRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.TableShapeRequest, p map[string]any) (any, error) {
			row, ok := v.RequireInt(p, "row")
			if !ok {
				return nil, v.Error()
			}
			height, ok := v.RequireInt(p, "height")
			if !ok {
				return nil, v.Error()
			}
			if err := e.SetTableRowHeight(request.SlideIndex, request.ShapeID, row, int64(height)); err != nil {
				return nil, err
			}
			return map[string]bool{"success": true}, nil
		},
	)
}

func handleSetTableColumnWidth(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.TableShapeRequest, bool) {
			return editorcommand.ParseTableShapeRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.TableShapeRequest, p map[string]any) (any, error) {
			col, ok := v.RequireInt(p, "col")
			if !ok {
				return nil, v.Error()
			}
			width, ok := v.RequireInt(p, "width")
			if !ok {
				return nil, v.Error()
			}
			if err := e.SetTableColumnWidth(request.SlideIndex, request.ShapeID, col, int64(width)); err != nil {
				return nil, err
			}
			return map[string]bool{"success": true}, nil
		},
	)
}

func handleAddTableRow(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.TableShapeRequest, bool) {
			return editorcommand.ParseTableShapeRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.TableShapeRequest, p map[string]any) (any, error) {
			height, _ := v.OptionalInt64(p, "height")
			if err := e.AddTableRow(request.SlideIndex, request.ShapeID, height); err != nil {
				return nil, err
			}
			return map[string]bool{"success": true}, nil
		},
	)
}

func handleAddTableColumn(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.TableShapeRequest, bool) {
			return editorcommand.ParseTableShapeRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.TableShapeRequest, p map[string]any) (any, error) {
			width, ok := v.RequireInt(p, "width")
			if !ok {
				return nil, v.Error()
			}
			if err := e.AddTableColumn(request.SlideIndex, request.ShapeID, int64(width)); err != nil {
				return nil, err
			}
			return map[string]bool{"success": true}, nil
		},
	)
}

func parseCellBorderUpdate(p map[string]any) (*tablemod.CellBorderSideUpdate, error) {
	borderRaw, hasBorder := p["border"]
	if !hasBorder || borderRaw == nil {
		return nil, nil //nolint:nilnil // nil update means clear the border
	}
	borderMap, ok := borderRaw.(map[string]any)
	if !ok {
		return nil, errors.New("border must be an object or null")
	}
	update := &tablemod.CellBorderSideUpdate{}
	if w, exists := borderMap["width"]; exists {
		switch wv := w.(type) {
		case float64:
			update.Width = int64(wv)
		case int:
			update.Width = int64(wv)
		case int64:
			update.Width = wv
		}
	}
	if c, ok := borderMap["color"].(string); ok {
		update.Color = c
	}
	if d, ok := borderMap["dash"].(string); ok {
		update.Dash = d
	}
	return update, nil
}

func handleUpdateTableCellBorder(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.TableCellRequest, bool) {
			return editorcommand.ParseTableCellRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.TableCellRequest, p map[string]any) (any, error) {
			side, ok := v.RequireString(p, "side")
			if !ok {
				return nil, v.Error()
			}
			switch side {
			case "left", "right", "top", "bottom":
			default:
				return nil, errors.New("side must be one of: left, right, top, bottom")
			}
			update, err := parseCellBorderUpdate(p)
			if err != nil {
				return nil, err
			}
			if err := e.UpdateTableCellBorder(
				request.SlideIndex, request.ShapeID, request.Row, request.Col, side, update,
			); err != nil {
				return nil, err
			}
			return map[string]bool{"success": true}, nil
		},
	)
}
