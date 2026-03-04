package editor

import (
	tablemod "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/table"
)

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
