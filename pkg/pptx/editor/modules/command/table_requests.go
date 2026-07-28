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
	// Geometry is optional: parseInt64Field is the caller's OptionalInt64, which
	// yields a zero value when the field is absent. Callers substitute defaults.
	x, _ := parseInt64Field(payload, "x")
	y, _ := parseInt64Field(payload, "y")
	cx, _ := parseInt64Field(payload, "cx")
	cy, _ := parseInt64Field(payload, "cy")
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

// CellStyleUpdate holds optional style fields parsed from an update_table_cell payload.
type CellStyleUpdate struct {
	SizePt          float64
	FontName        string
	Bold            *bool
	Italic          *bool
	Underline       *bool
	Color           string
	BackgroundColor string
	HasStyle        bool
}

// ParseOptionalCellStyleUpdate reads the run and cell style fields from an
// updates map: size_pt, font_name, bold, italic, underline, color and
// background_color.
func ParseOptionalCellStyleUpdate(updates map[string]any) (CellStyleUpdate, error) {
	var out CellStyleUpdate
	if raw, ok := updates["size_pt"]; ok {
		switch v := raw.(type) {
		case float64:
			out.SizePt = v
			out.HasStyle = true
		case int:
			out.SizePt = float64(v)
			out.HasStyle = true
		default:
			return CellStyleUpdate{}, fmt.Errorf("size_pt must be a number, got %T", raw)
		}
	}

	stringFields := map[string]*string{
		"font_name":        &out.FontName,
		"color":            &out.Color,
		"background_color": &out.BackgroundColor,
	}
	for key, dst := range stringFields {
		raw, ok := updates[key]
		if !ok {
			continue
		}
		s, isString := raw.(string)
		if !isString {
			return CellStyleUpdate{}, fmt.Errorf("%s must be a string, got %T", key, raw)
		}
		*dst = s
		out.HasStyle = true
	}

	boolFields := map[string]**bool{
		"bold":      &out.Bold,
		"italic":    &out.Italic,
		"underline": &out.Underline,
	}
	for key, dst := range boolFields {
		raw, ok := updates[key]
		if !ok {
			continue
		}
		b, isBool := raw.(bool)
		if !isBool {
			return CellStyleUpdate{}, fmt.Errorf("%s must be a boolean, got %T", key, raw)
		}
		value := b
		*dst = &value
		out.HasStyle = true
	}
	return out, nil
}
