package editor

import (
	tablemod "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/table"
)

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
