package command

import (
	"errors"
	"fmt"
)

type ParseInt64FieldFn func(map[string]any, string) (int64, bool)

type TableAddRequest struct {
	SlideIndex int
	Rows       int
	Cols       int
	X          int64
	Y          int64
	CX         int64
	CY         int64
}

func ParseTableAddRequest(
	payload map[string]any,
	parseSlideIndex ParseSlideIndexFn,
	parseIntField ParseIntFieldFn,
	parseInt64Field ParseInt64FieldFn,
) (TableAddRequest, bool) {
	slideIndex, ok := parseSlideIndex(payload)
	if !ok {
		return TableAddRequest{}, false
	}
	rows, ok := parseIntField(payload, "rows")
	if !ok {
		return TableAddRequest{}, false
	}
	cols, ok := parseIntField(payload, "cols")
	if !ok {
		return TableAddRequest{}, false
	}
	x, ok := parseInt64Field(payload, "x")
	_ = ok
	y, ok := parseInt64Field(payload, "y")
	_ = ok
	cx, ok := parseInt64Field(payload, "cx")
	_ = ok
	cy, ok := parseInt64Field(payload, "cy")
	_ = ok
	return TableAddRequest{
		SlideIndex: slideIndex,
		Rows:       rows,
		Cols:       cols,
		X:          x,
		Y:          y,
		CX:         cx,
		CY:         cy,
	}, true
}

func ValidateTableDimensions(rows, cols, maxDimension int) error {
	if rows < 1 || rows > maxDimension {
		return fmt.Errorf("rows %d must be between 1 and %d", rows, maxDimension)
	}
	if cols < 1 || cols > maxDimension {
		return fmt.Errorf("cols %d must be between 1 and %d", cols, maxDimension)
	}
	return nil
}

type TableShapeRequest struct {
	SlideIndex int
	ShapeID    int
}

func ParseTableShapeRequest(
	payload map[string]any,
	parseIntField ParseIntFieldFn,
) (TableShapeRequest, bool) {
	slideIndex, ok := parseIntField(payload, "slide_index")
	if !ok {
		return TableShapeRequest{}, false
	}
	shapeID, ok := parseIntField(payload, "shape_id")
	if !ok {
		return TableShapeRequest{}, false
	}
	return TableShapeRequest{
		SlideIndex: slideIndex,
		ShapeID:    shapeID,
	}, true
}

type TableCellRangeRequest struct {
	SlideIndex int
	ShapeID    int
	Row1       int
	Col1       int
	Row2       int
	Col2       int
}

func ParseTableCellRangeRequest(
	payload map[string]any,
	parseIntField ParseIntFieldFn,
) (TableCellRangeRequest, bool) {
	tableShape, ok := ParseTableShapeRequest(payload, parseIntField)
	if !ok {
		return TableCellRangeRequest{}, false
	}
	row1, ok := parseIntField(payload, "row1")
	if !ok {
		return TableCellRangeRequest{}, false
	}
	col1, ok := parseIntField(payload, "col1")
	if !ok {
		return TableCellRangeRequest{}, false
	}
	row2, ok := parseIntField(payload, "row2")
	if !ok {
		return TableCellRangeRequest{}, false
	}
	col2, ok := parseIntField(payload, "col2")
	if !ok {
		return TableCellRangeRequest{}, false
	}
	return TableCellRangeRequest{
		SlideIndex: tableShape.SlideIndex,
		ShapeID:    tableShape.ShapeID,
		Row1:       row1,
		Col1:       col1,
		Row2:       row2,
		Col2:       col2,
	}, true
}

type TableCellRequest struct {
	SlideIndex int
	ShapeID    int
	Row        int
	Col        int
}

func ParseTableCellRequest(
	payload map[string]any,
	parseIntField ParseIntFieldFn,
) (TableCellRequest, bool) {
	tableShape, ok := ParseTableShapeRequest(payload, parseIntField)
	if !ok {
		return TableCellRequest{}, false
	}
	row, ok := parseIntField(payload, "row")
	if !ok {
		return TableCellRequest{}, false
	}
	col, ok := parseIntField(payload, "col")
	if !ok {
		return TableCellRequest{}, false
	}
	return TableCellRequest{
		SlideIndex: tableShape.SlideIndex,
		ShapeID:    tableShape.ShapeID,
		Row:        row,
		Col:        col,
	}, true
}

type TableStyleRequest struct {
	SlideIndex int
	ShapeID    int
	StyleGUID  string
}

func ParseTableStyleRequest(
	payload map[string]any,
	parseIntField ParseIntFieldFn,
	parseStringField ParseStringFieldFn,
) (TableStyleRequest, bool) {
	tableShape, ok := ParseTableShapeRequest(payload, parseIntField)
	if !ok {
		return TableStyleRequest{}, false
	}
	styleGUID, ok := parseStringField(payload, "style_guid")
	if !ok {
		return TableStyleRequest{}, false
	}
	return TableStyleRequest{
		SlideIndex: tableShape.SlideIndex,
		ShapeID:    tableShape.ShapeID,
		StyleGUID:  styleGUID,
	}, true
}

func ParseRequiredObjectField(
	payload map[string]any,
	key, missingErr, typeErr string,
) (map[string]any, error) {
	value, ok := payload[key]
	if !ok {
		return nil, fmt.Errorf("%s", missingErr)
	}
	objectValue, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s", typeErr)
	}
	return objectValue, nil
}

func ParseOptionalTextUpdate(updates map[string]any) (string, bool, error) {
	text, hasText := updates["text"]
	if !hasText {
		return "", false, nil
	}
	textValue, ok := text.(string)
	if !ok {
		return "", false, errors.New("text update must be a string")
	}
	return textValue, true, nil
}
