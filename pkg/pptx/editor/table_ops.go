package editor

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	tablemod "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/table"
)

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
// Use styleGUID to specify a PowerPoint table style GUID; pass "" to use the default style.
func (e *PresentationEditor) AddTable(
	slideIndex, rowCount, colCount int,
	x, y, cx, cy int64,
	styleGUID ...string,
) (int, error) {
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
	if len(styleGUID) > 0 && styleGUID[0] != "" {
		spec.StyleGUID = styleGUID[0]
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

func (e *PresentationEditor) UpdateTableCellContent(
	slideIndex, shapeID, rowIdx, colIdx int,
	update tablemod.CellContentUpdate,
) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.UpdateTableCellContentInFrame(frame, rowIdx, colIdx, update)
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

// AddTableRow appends a new empty row to an existing table.
// height is in EMU; pass 0 to let PowerPoint auto-size.
func (e *PresentationEditor) AddTableRow(slideIndex, shapeID int, height int64) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.AddTableRowInFrame(frame, height)
	if err != nil {
		return err
	}
	e.parts.Set(partPath, tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

// UpdateTableCellBorder sets or clears a border on a single cell side.
// side must be "left", "right", "top", or "bottom".
// update=nil removes the border; a non-nil update sets width/color/dash.
func (e *PresentationEditor) UpdateTableCellBorder(
	slideIndex, shapeID, row, col int,
	side string,
	update *tablemod.CellBorderSideUpdate,
) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.UpdateTableCellBordersInFrame(frame, row, col, side, update)
	if err != nil {
		return err
	}
	e.parts.Set(partPath, tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

// AddTableColumn appends a new empty column to an existing table.
// width is in EMU.
func (e *PresentationEditor) AddTableColumn(slideIndex, shapeID int, width int64) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.AddTableColumnInFrame(frame, width)
	if err != nil {
		return err
	}
	e.parts.Set(partPath, tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

// InsertTableRow inserts a new empty row before the row at atIndex.
// Use atIndex == current row count to append at the end.
// height is in EMU; pass 0 to let PowerPoint auto-size.
func (e *PresentationEditor) InsertTableRow(slideIndex, shapeID, atIndex int, height int64) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.InsertTableRowInFrame(frame, atIndex, height)
	if err != nil {
		return err
	}
	e.parts.Set(partPath, tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

// RemoveTableRow removes the row at atIndex from an existing table.
func (e *PresentationEditor) RemoveTableRow(slideIndex, shapeID, atIndex int) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.RemoveTableRowInFrame(frame, atIndex)
	if err != nil {
		return err
	}
	e.parts.Set(partPath, tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

// InsertTableColumn inserts a new empty column before the column at atIndex.
// Use atIndex == current column count to append at the end.
// width is in EMU.
func (e *PresentationEditor) InsertTableColumn(slideIndex, shapeID, atIndex int, width int64) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.InsertTableColumnInFrame(frame, atIndex, width)
	if err != nil {
		return err
	}
	e.parts.Set(partPath, tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

// RemoveTableColumn removes the column at atIndex from an existing table.
func (e *PresentationEditor) RemoveTableColumn(slideIndex, shapeID, atIndex int) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedFrame, err := tablemod.RemoveTableColumnInFrame(frame, atIndex)
	if err != nil {
		return err
	}
	e.parts.Set(partPath, tablemod.ReplaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}
