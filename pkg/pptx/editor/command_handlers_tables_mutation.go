package editor

import (
	"encoding/json"
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	tablemod "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/table"
)

func getSlideTableFrame(e *PresentationEditor, slideIndex, shapeID int) (
	string,
	[]byte,
	int,
	int,
	[]byte,
	error,
) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return "", nil, 0, 0, nil, fmt.Errorf("slide index %d out of range", slideIndex)
	}
	partPath := e.slides[slideIndex].Part
	var ok bool
	var slideContent []byte
	slideContent, ok = e.parts.Get(partPath)
	if !ok {
		return "", nil, 0, 0, nil, errors.New("slide part not found")
	}
	frameStart, frameEnd, frame, err := tablemod.FindTableFrame(slideContent, shapeID)
	if err != nil {
		return "", nil, 0, 0, nil, err
	}
	return partPath, slideContent, frameStart, frameEnd, frame, nil
}

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
			style, err := editorcommand.ParseOptionalCellStyleUpdate(updates)
			if err != nil {
				return nil, err
			}
			if style.HasStyle {
				var textPtr *string
				if hasText {
					textPtr = &text
				}
				if err := e.UpdateTableCellContent(
					request.SlideIndex,
					request.ShapeID,
					request.Row,
					request.Col,
					tablemod.CellContentUpdate{
						Text:     textPtr,
						SizePt:   style.SizePt,
						FontName: style.FontName,
					},
				); err != nil {
					return nil, err
				}
			} else if hasText {
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
