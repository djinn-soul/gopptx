package editor

import (
	"errors"
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
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

// SetTableStyle sets the table style for the specified table on a slide.
// The styleGUID must be a valid PowerPoint table style GUID, e.g.:
//
//	"{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}" - Medium Style 2 - Accent 1
//	"{B9AC3A68-259E-4EED-9050-4AE35E7F2B2D}" - Light Style 1
//	"{5940675A-B579-460E-94D1-54222C63F5DA}" - Medium Style 1 - Accent 1
func (e *PresentationEditor) SetTableStyle(slideIndex, shapeID int, styleGUID string) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.SetTableStyleInFrame(frame, styleGUID)
	if err != nil {
		return err
	}
	updatedSlide := tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame)
	e.parts.Set(partPath, updatedSlide)
	return nil
}

// AddTable adds a new graphic frame containing a table.
func (e *PresentationEditor) AddTable(slideIndex, rowCount, colCount int, x, y, cx, cy int64) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, fmt.Errorf("slide index %d out of range", slideIndex)
	}

	slideRef := e.slides[slideIndex]
	spec := &pptxxml.TableSpec{
		X:          x,
		Y:          y,
		CX:         cx,
		CY:         cy,
		Rows:       make([][]string, rowCount),
		StyledRows: make([][]pptxxml.TableCellSpec, rowCount),
	}
	for i := range rowCount {
		spec.Rows[i] = make([]string, colCount)
		spec.StyledRows[i] = make([]pptxxml.TableCellSpec, colCount)
		for j := range colCount {
			spec.StyledRows[i][j] = pptxxml.TableCellSpec{}
		}
	}

	shapeID := e.nextShapeID(slideRef.Part)
	shapeXML := pptxxml.RenderTable(spec, shapeID)
	if err := e.appendShapeToSlide(slideRef.Part, shapeXML); err != nil {
		return 0, fmt.Errorf("append table shape: %w", err)
	}
	return shapeID, nil
}

// GetTable reads a table's structure entirely from XML.
func (e *PresentationEditor) GetTable(slideIndex, shapeID int) (map[string]any, error) {
	_, _, _, _, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return nil, err
	}
	return tablemod.BuildTableInfo(frame)
}

// UpdateTableFlags modifies properties of the table like firstRow, bandRow, etc.
func (e *PresentationEditor) UpdateTableFlags(slideIndex, shapeID int, flags map[string]any) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.UpdateTableFlagsInFrame(frame, flags)
	if err != nil {
		return err
	}
	e.parts.Set(partPath, tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

func (e *PresentationEditor) UpdateTableCellText(slideIndex, shapeID, rowIdx, colIdx int, text string) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.UpdateTableCellTextInFrame(frame, rowIdx, colIdx, text)
	if err != nil {
		return err
	}
	e.parts.Set(partPath, tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

func (e *PresentationEditor) MergeTableCells(slideIndex, shapeID, row1, col1, row2, col2 int) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.MergeCellsInFrame(frame, row1, col1, row2, col2)
	if err != nil {
		return err
	}
	e.parts.Set(partPath, tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

func (e *PresentationEditor) SplitTableCell(slideIndex, shapeID, row, col int) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.SplitCellInFrame(frame, row, col)
	if err != nil {
		return err
	}
	e.parts.Set(partPath, tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

func (e *PresentationEditor) SetTableRowHeight(slideIndex, shapeID, row int, height int64) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.UpdateTableRowHeightInFrame(frame, row, height)
	if err != nil {
		return err
	}
	e.parts.Set(partPath, tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

func (e *PresentationEditor) SetTableColumnWidth(slideIndex, shapeID, col int, width int64) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.UpdateTableColumnWidthInFrame(frame, col, width)
	if err != nil {
		return err
	}
	e.parts.Set(partPath, tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

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

	// 1. Create empty table
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
	// Skipping this call preserves the firstRow="1" bandRow="1" defaults
	// that RenderTable writes, so AddTableWithData(..., &TableInitSpec{Data: ...})
	// does not silently strip the default header-row/banded-row styling.
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
