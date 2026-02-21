package editor

import (
	"bytes"
	"fmt"
)

func setOrInsertAttr(openingTag []byte, attrName, attrValue string) []byte {
	attrStr := []byte(" " + attrName + `="`)
	idx := bytes.Index(openingTag, attrStr)
	if idx != -1 {
		valStart := idx + len(attrStr)
		valEndRel := bytes.Index(openingTag[valStart:], []byte(`"`))
		if valEndRel != -1 {
			valEnd := valStart + valEndRel
			updated := make([]byte, 0, len(openingTag)-((valEnd-valStart)-len(attrValue)))
			updated = append(updated, openingTag[:valStart]...)
			updated = append(updated, []byte(attrValue)...)
			updated = append(updated, openingTag[valEnd:]...)
			return updated
		}
	}

	insertAt := len(openingTag) - 1
	if insertAt > 0 && openingTag[insertAt-1] == '/' {
		insertAt--
	}
	updated := make([]byte, 0, len(openingTag)+len(attrName)+len(attrValue)+4)
	updated = append(updated, openingTag[:insertAt]...)
	updated = append(updated, []byte(" "+attrName+`="`+attrValue+`"`)...)
	updated = append(updated, openingTag[insertAt:]...)
	return updated
}

func setTcAttr(tcContent []byte, attrName, attrValue string) []byte {
	tagEnd := bytes.Index(tcContent, []byte(">"))
	if tagEnd == -1 {
		return tcContent
	}
	openTag := tcContent[:tagEnd+1]
	updatedTag := setOrInsertAttr(openTag, attrName, attrValue)
	updated := make([]byte, 0, len(tcContent)-((len(openTag))-len(updatedTag)))
	updated = append(updated, updatedTag...)
	updated = append(updated, tcContent[tagEnd+1:]...)
	return updated
}

func removeTcAttr(tcContent []byte, attrName string) []byte {
	tagEnd := bytes.Index(tcContent, []byte(">"))
	if tagEnd == -1 {
		return tcContent
	}
	openTag := tcContent[:tagEnd+1]
	attrStr := []byte(" " + attrName + `="`)
	idx := bytes.Index(openTag, attrStr)
	if idx == -1 {
		return tcContent
	}
	valStart := idx + len(attrStr)
	valEndRel := bytes.Index(openTag[valStart:], []byte(`"`))
	if valEndRel == -1 {
		return tcContent
	}
	valEnd := valStart + valEndRel + 1
	updatedTag := make([]byte, 0, len(openTag)-(valEnd-idx))
	updatedTag = append(updatedTag, openTag[:idx]...)
	updatedTag = append(updatedTag, openTag[valEnd:]...)

	updated := make([]byte, 0, len(tcContent)-((len(openTag))-len(updatedTag)))
	updated = append(updated, updatedTag...)
	updated = append(updated, tcContent[tagEnd+1:]...)
	return updated
}

func mutateTableRows(
	frame []byte,
	rowStart int,
	rowEnd int,
	mutator func(row int, rowContent []byte) ([]byte, error),
) ([]byte, error) {
	var out bytes.Buffer
	cursor := 0
	row := 0

	for {
		trRel := bytes.Index(frame[cursor:], []byte("<a:tr"))
		if trRel == -1 {
			out.Write(frame[cursor:])
			break
		}
		trStart := cursor + trRel
		trEndRel := bytes.Index(frame[trStart:], []byte("</a:tr>"))
		if trEndRel == -1 {
			return nil, fmt.Errorf("invalid tr xml at row %d", row)
		}
		trEnd := trStart + trEndRel + len("</a:tr>")

		out.Write(frame[cursor:trStart])
		rowContent := frame[trStart:trEnd]
		if row >= rowStart && row <= rowEnd {
			updated, err := mutator(row, rowContent)
			if err != nil {
				return nil, err
			}
			out.Write(updated)
		} else {
			out.Write(rowContent)
		}

		cursor = trEnd
		row++
	}

	return out.Bytes(), nil
}

func mutateTableCells(
	rowContent []byte,
	colStart int,
	colEnd int,
	mutator func(col int, cellContent []byte) ([]byte, error),
) ([]byte, error) {
	var out bytes.Buffer
	cursor := 0
	col := 0

	for {
		tcRel := bytes.Index(rowContent[cursor:], []byte("<a:tc"))
		if tcRel == -1 {
			out.Write(rowContent[cursor:])
			break
		}
		tcStart := cursor + tcRel
		tcEndRel := bytes.Index(rowContent[tcStart:], []byte("</a:tc>"))
		if tcEndRel == -1 {
			return nil, fmt.Errorf("invalid tc xml at col %d", col)
		}
		tcEnd := tcStart + tcEndRel + len("</a:tc>")

		out.Write(rowContent[cursor:tcStart])
		cellContent := rowContent[tcStart:tcEnd]
		if col >= colStart && col <= colEnd {
			updated, err := mutator(col, cellContent)
			if err != nil {
				return nil, err
			}
			out.Write(updated)
		} else {
			out.Write(cellContent)
		}

		cursor = tcEnd
		col++
	}

	return out.Bytes(), nil
}

func (e *PresentationEditor) MergeTableCells(slideIndex, shapeID, row1, col1, row2, col2 int) error {
	if row1 < 0 || col1 < 0 || row2 < 0 || col2 < 0 {
		return fmt.Errorf("merge coordinates must be non-negative")
	}
	if row1 > row2 || col1 > col2 {
		return fmt.Errorf("merge coordinates must be ordered: row1<=row2 and col1<=col2")
	}

	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	parsed, err := parseTable(frame)
	if err != nil {
		return err
	}
	rows, cols := tableDimensions(parsed)
	if row2 >= rows || col2 >= cols {
		return fmt.Errorf("merge range [%d:%d,%d:%d] out of table bounds %dx%d", row1, row2, col1, col2, rows, cols)
	}

	rowSpan := row2 - row1 + 1
	colSpan := col2 - col1 + 1

	updatedFrame, err := mutateTableRows(frame, row1, row2, func(r int, rowContent []byte) ([]byte, error) {
		return mutateTableCells(rowContent, col1, col2, func(c int, cellContent []byte) ([]byte, error) {
			cellContent = removeTcAttr(cellContent, "rowSpan")
			cellContent = removeTcAttr(cellContent, "gridSpan")
			cellContent = removeTcAttr(cellContent, "vMerge")
			cellContent = removeTcAttr(cellContent, "hMerge")

			if r == row1 && c == col1 {
				if rowSpan > 1 {
					cellContent = setTcAttr(cellContent, "rowSpan", fmt.Sprintf("%d", rowSpan))
				}
				if colSpan > 1 {
					cellContent = setTcAttr(cellContent, "gridSpan", fmt.Sprintf("%d", colSpan))
				}
				return cellContent, nil
			}
			if r > row1 {
				cellContent = setTcAttr(cellContent, "vMerge", "1")
			}
			if c > col1 {
				cellContent = setTcAttr(cellContent, "hMerge", "1")
			}
			return cellContent, nil
		})
	})
	if err != nil {
		return err
	}

	e.parts.Set(partPath, replaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

func (e *PresentationEditor) SplitTableCell(slideIndex, shapeID, row, col int) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	parsed, err := parseTable(frame)
	if err != nil {
		return err
	}
	rows, cols := tableDimensions(parsed)
	if row < 0 || row >= rows || col < 0 || col >= cols {
		return fmt.Errorf("table cell [%d,%d] out of range", row, col)
	}

	cell := parsed.Rows[row].Cells[col]
	rowSpan := cell.RowSpan
	if rowSpan <= 0 {
		rowSpan = 1
	}
	colSpan := cell.GridSpan
	if colSpan <= 0 {
		colSpan = 1
	}
	if rowSpan == 1 && colSpan == 1 {
		return fmt.Errorf("cell [%d,%d] is not a merge origin", row, col)
	}

	rowEnd := row + rowSpan - 1
	colEnd := col + colSpan - 1
	if rowEnd >= rows || colEnd >= cols {
		return fmt.Errorf("merged span at [%d,%d] exceeds table bounds", row, col)
	}

	updatedFrame, err := mutateTableRows(frame, row, rowEnd, func(r int, rowContent []byte) ([]byte, error) {
		return mutateTableCells(rowContent, col, colEnd, func(c int, cellContent []byte) ([]byte, error) {
			cellContent = removeTcAttr(cellContent, "rowSpan")
			cellContent = removeTcAttr(cellContent, "gridSpan")
			cellContent = removeTcAttr(cellContent, "vMerge")
			cellContent = removeTcAttr(cellContent, "hMerge")
			if r == row && c == col {
				return cellContent, nil
			}
			return cellContent, nil
		})
	})
	if err != nil {
		return err
	}

	e.parts.Set(partPath, replaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}
