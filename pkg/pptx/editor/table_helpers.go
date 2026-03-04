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

//nolint:gocognit // Merge behavior needs explicit XML-state checks for correctness and bounds safety.
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
