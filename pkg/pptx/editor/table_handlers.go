package editor

import (
	"encoding/json"
	"errors"
	"fmt"
)

const maxTableDimension = 1000

// handleAddTable handles the OP_ADD_TABLE JSON-RPC command.
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
