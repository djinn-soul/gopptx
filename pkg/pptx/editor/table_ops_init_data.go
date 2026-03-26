package editor

import "fmt"

// AddTableWithData creates a table with optional data population and configuration.
// This is a convenience method that batches all table setup operations.
// Returns the shape ID of the created table.
func (e *PresentationEditor) AddTableWithData(
	slideIndex, rowCount, colCount int,
	x, y, cx, cy int64,
	spec *TableInitSpec,
) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, fmt.Errorf("slide index %d out of range", slideIndex)
	}

	shapeID, err := e.AddTable(slideIndex, rowCount, colCount, x, y, cx, cy)
	if err != nil {
		return 0, fmt.Errorf("add table: %w", err)
	}
	if spec == nil {
		return shapeID, nil
	}
	if err := e.populateTableData(slideIndex, shapeID, spec.Data); err != nil {
		return shapeID, err
	}
	if err := e.applyTableColumnWidths(slideIndex, shapeID, spec.ColumnWidths); err != nil {
		return shapeID, err
	}
	if err := e.applyTableRowHeights(slideIndex, shapeID, spec.RowHeights); err != nil {
		return shapeID, err
	}
	if err := e.applyExplicitTableFlags(slideIndex, shapeID, spec); err != nil {
		return shapeID, err
	}
	return shapeID, nil
}

func (e *PresentationEditor) populateTableData(slideIndex, shapeID int, data [][]string) error {
	for rowIdx, row := range data {
		for colIdx, text := range row {
			if err := e.UpdateTableCellText(slideIndex, shapeID, rowIdx, colIdx, text); err != nil {
				return fmt.Errorf("set cell [%d,%d]: %w", rowIdx, colIdx, err)
			}
		}
	}
	return nil
}

func (e *PresentationEditor) applyTableColumnWidths(slideIndex, shapeID int, widths []int64) error {
	for colIdx, width := range widths {
		if err := e.SetTableColumnWidth(slideIndex, shapeID, colIdx, width); err != nil {
			return fmt.Errorf("set column width %d: %w", colIdx, err)
		}
	}
	return nil
}

func (e *PresentationEditor) applyTableRowHeights(slideIndex, shapeID int, heights []int64) error {
	for rowIdx, height := range heights {
		if err := e.SetTableRowHeight(slideIndex, shapeID, rowIdx, height); err != nil {
			return fmt.Errorf("set row height %d: %w", rowIdx, err)
		}
	}
	return nil
}

func (e *PresentationEditor) applyExplicitTableFlags(
	slideIndex, shapeID int,
	spec *TableInitSpec,
) error {
	// Set flags only when the caller explicitly enables at least one.
	// Skipping this call preserves firstRow/bandRow defaults from RenderTable.
	if !spec.FirstRow && !spec.FirstCol && !spec.LastRow && !spec.LastCol && !spec.BandRow && !spec.BandCol {
		return nil
	}
	flags := map[string]any{
		"first_row": spec.FirstRow,
		"first_col": spec.FirstCol,
		"last_row":  spec.LastRow,
		"last_col":  spec.LastCol,
		"band_row":  spec.BandRow,
		"band_col":  spec.BandCol,
	}
	if err := e.UpdateTableFlags(slideIndex, shapeID, flags); err != nil {
		return fmt.Errorf("update flags: %w", err)
	}
	return nil
}
