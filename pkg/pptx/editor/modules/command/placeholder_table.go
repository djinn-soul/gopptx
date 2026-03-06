package command

import (
	"errors"
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

func ParsePlaceholderTableSpec(payload map[string]any) (*pptxxml.TableSpec, bool, error) {
	rawTable, ok := payload["table"]
	if !ok {
		return nil, false, nil
	}
	tableMap, ok := rawTable.(map[string]any)
	if !ok {
		return nil, false, errors.New("table must be an object")
	}
	rawRows, ok := tableMap["rows"]
	if !ok {
		return nil, false, errors.New("table.rows is required")
	}
	rowsSlice, ok := rawRows.([]any)
	if !ok || len(rowsSlice) == 0 {
		return nil, false, errors.New("table.rows must be a non-empty 2D array")
	}

	rows := make([][]string, 0, len(rowsSlice))
	for rowIdx, rowRaw := range rowsSlice {
		rowSlice, ok := rowRaw.([]any)
		if !ok || len(rowSlice) == 0 {
			return nil, false, fmt.Errorf("table.rows[%d] must be a non-empty array", rowIdx)
		}
		row := make([]string, 0, len(rowSlice))
		for _, cell := range rowSlice {
			if cell == nil {
				row = append(row, "")
				continue
			}
			if str, ok := cell.(string); ok {
				row = append(row, str)
				continue
			}
			row = append(row, fmt.Sprint(cell))
		}
		rows = append(rows, row)
	}

	spec := &pptxxml.TableSpec{Rows: rows}
	if altText, ok := tableMap["alt_text"].(string); ok {
		spec.AltText = altText
	}
	if decorative, ok := tableMap["decorative"].(bool); ok {
		spec.IsDecorative = decorative
	}

	return spec, true, nil
}
