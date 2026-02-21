package editor

import (
	"encoding/json"
	"errors"
)

// handleAddTable handles the OP_ADD_TABLE JSON-RPC command.
func handleAddTable(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
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

	x, _ := v.OptionalInt64(p, "x")
	y, _ := v.OptionalInt64(p, "y")
	cx, _ := v.OptionalInt64(p, "cx")
	cy, _ := v.OptionalInt64(p, "cy")

	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
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
	slideIndex, _ := v.RequireInt(p, "slide_index")
	shapeID, _ := v.RequireInt(p, "shape_id")
	row1, _ := v.RequireInt(p, "row1")
	col1, _ := v.RequireInt(p, "col1")
	row2, _ := v.RequireInt(p, "row2")
	col2, _ := v.RequireInt(p, "col2")

	if v.HasErrors() {
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
	slideIndex, _ := v.RequireInt(p, "slide_index")
	shapeID, _ := v.RequireInt(p, "shape_id")
	row, _ := v.RequireInt(p, "row")
	col, _ := v.RequireInt(p, "col")

	if v.HasErrors() {
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
