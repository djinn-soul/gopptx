package tables

import (
	"fmt"
	"math"
	"strings"
)

type tableMergeCoverage struct {
	AnchorRow int
	AnchorCol int
	HMerge    bool
	VMerge    bool
}

// TableRowsWithMerges handles merged cell calculation for table rendering.
func TableRowsWithMerges(table Table, slideIndex int) ([][]TableCell, error) {
	rows := tableRowsForRender(table)
	return applyTableCellMerges(rows, len(table.ColumnWidths), slideIndex)
}

func applyTableCellMerges(rows [][]TableCell, columnCount int, slideIndex int) ([][]TableCell, error) {
	mergedRows := copyTableRows(rows)
	if len(mergedRows) == 0 || columnCount == 0 {
		return mergedRows, nil
	}

	coverage := make([][]*tableMergeCoverage, len(mergedRows))
	for r := range coverage {
		coverage[r] = make([]*tableMergeCoverage, columnCount)
	}

	for rowIdx := range mergedRows {
		for cellIdx := range mergedRows[rowIdx] {
			cell := mergedRows[rowIdx][cellIdx]

			if cell.RowSpan <= 0 {
				return nil, fmt.Errorf(
					"slide %d table row %d cell %d row span must be >= 1",
					slideIndex,
					rowIdx+1,
					cellIdx+1,
				)
			}
			if cell.ColSpan <= 0 {
				return nil, fmt.Errorf(
					"slide %d table row %d cell %d col span must be >= 1",
					slideIndex,
					rowIdx+1,
					cellIdx+1,
				)
			}

			if covered := coverage[rowIdx][cellIdx]; covered != nil {
				if !isTableMergePlaceholderCell(cell) {
					return nil, fmt.Errorf(
						"slide %d table row %d cell %d overlaps merged range from row %d cell %d; covered cells must be empty placeholders",
						slideIndex,
						rowIdx+1,
						cellIdx+1,
						covered.AnchorRow+1,
						covered.AnchorCol+1,
					)
				}
				cell.RowSpan = 1
				cell.ColSpan = 1
				cell.HMerge = covered.HMerge
				cell.VMerge = covered.VMerge
				mergedRows[rowIdx][cellIdx] = cell
				continue
			}

			cell.HMerge = false
			cell.VMerge = false
			mergedRows[rowIdx][cellIdx] = cell

			rowEnd := rowIdx + cell.RowSpan
			colEnd := cellIdx + cell.ColSpan
			if rowEnd > len(mergedRows) || colEnd > columnCount {
				return nil, fmt.Errorf(
					"slide %d table row %d cell %d merged span (%dx%d) exceeds table bounds",
					slideIndex,
					rowIdx+1,
					cellIdx+1,
					cell.RowSpan,
					cell.ColSpan,
				)
			}

			if cell.RowSpan == 1 && cell.ColSpan == 1 {
				continue
			}

			for rr := rowIdx; rr < rowEnd; rr++ {
				for cc := cellIdx; cc < colEnd; cc++ {
					if rr == rowIdx && cc == cellIdx {
						continue
					}
					if coverage[rr][cc] != nil {
						return nil, fmt.Errorf(
							"slide %d table row %d cell %d has overlapping merged ranges",
							slideIndex,
							rowIdx+1,
							cellIdx+1,
						)
					}
					coverage[rr][cc] = &tableMergeCoverage{
						AnchorRow: rowIdx,
						AnchorCol: cellIdx,
						HMerge:    cc > cellIdx,
						VMerge:    rr > rowIdx,
					}
				}
			}
		}
	}

	return mergedRows, nil
}

func isTableMergePlaceholderCell(cell TableCell) bool {
	if strings.TrimSpace(cell.Text) != "" {
		return false
	}
	if cell.Bold || strings.TrimSpace(cell.BackgroundColor) != "" {
		return false
	}
	if strings.TrimSpace(cell.Align) != "" || strings.TrimSpace(cell.VAlign) != "" {
		return false
	}
	if strings.TrimSpace(cell.BorderColor) != "" {
		return false
	}
	if !isFiniteNonNegative(cell.BorderWidthPt) || cell.BorderWidthPt > 0 {
		return false
	}
	if cell.BorderLeft != nil || cell.BorderRight != nil || cell.BorderTop != nil || cell.BorderBottom != nil {
		return false
	}
	if cell.MarginLeftPt != nil || cell.MarginRightPt != nil || cell.MarginTopPt != nil || cell.MarginBottomPt != nil {
		return false
	}
	if cell.WrapText != nil {
		return false
	}
	if cell.RowSpan != 1 || cell.ColSpan != 1 {
		return false
	}
	if cell.HMerge || cell.VMerge {
		return false
	}
	return true
}

func isFiniteNonNegative(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0) && v >= 0
}
