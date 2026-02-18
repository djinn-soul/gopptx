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

	coverage := initMergeCoverage(len(mergedRows), columnCount)

	for rowIdx := range mergedRows {
		for cellIdx := range mergedRows[rowIdx] {
			cell := mergedRows[rowIdx][cellIdx]

			if err := validateCellSpans(cell, slideIndex, rowIdx, cellIdx); err != nil {
				return nil, err
			}

			if covered := coverage[rowIdx][cellIdx]; covered != nil {
				if err := handleCoveredCell(&cell, covered, slideIndex, rowIdx, cellIdx); err != nil {
					return nil, err
				}
				mergedRows[rowIdx][cellIdx] = cell
				continue
			}

			cell.HMerge = false
			cell.VMerge = false
			mergedRows[rowIdx][cellIdx] = cell

			if err := updateCoverage(coverage, mergedRows, cell, slideIndex, rowIdx, cellIdx, columnCount); err != nil {
				return nil, err
			}
		}
	}

	return mergedRows, nil
}

func initMergeCoverage(rowCount, columnCount int) [][]*tableMergeCoverage {
	coverage := make([][]*tableMergeCoverage, rowCount)
	for r := range coverage {
		coverage[r] = make([]*tableMergeCoverage, columnCount)
	}
	return coverage
}

func validateCellSpans(cell TableCell, slideIndex, rowIdx, cellIdx int) error {
	if cell.RowSpan <= 0 {
		return fmt.Errorf("slide %d table row %d cell %d row span must be >= 1", slideIndex, rowIdx+1, cellIdx+1)
	}
	if cell.ColSpan <= 0 {
		return fmt.Errorf("slide %d table row %d cell %d col span must be >= 1", slideIndex, rowIdx+1, cellIdx+1)
	}
	return nil
}

func handleCoveredCell(cell *TableCell, covered *tableMergeCoverage, slideIndex, rowIdx, cellIdx int) error {
	if !isTableMergePlaceholderCell(*cell) {
		return fmt.Errorf(
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
	return nil
}

func updateCoverage(
	coverage [][]*tableMergeCoverage,
	rows [][]TableCell,
	cell TableCell,
	slideIndex, rowIdx, cellIdx, columnCount int,
) error {
	rowEnd := rowIdx + cell.RowSpan
	colEnd := cellIdx + cell.ColSpan
	if rowEnd > len(rows) || colEnd > columnCount {
		return fmt.Errorf(
			"slide %d table row %d cell %d merged span (%dx%d) exceeds table bounds",
			slideIndex,
			rowIdx+1,
			cellIdx+1,
			cell.RowSpan,
			cell.ColSpan,
		)
	}

	if cell.RowSpan == 1 && cell.ColSpan == 1 {
		return nil
	}

	for rr := rowIdx; rr < rowEnd; rr++ {
		for cc := cellIdx; cc < colEnd; cc++ {
			if rr == rowIdx && cc == cellIdx {
				continue
			}
			if coverage[rr][cc] != nil {
				return fmt.Errorf(
					"slide %d table row %d cell %d has overlapping merged ranges",
					slideIndex, rowIdx+1, cellIdx+1,
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
	return nil
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
