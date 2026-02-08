package pptx

import (
	"github.com/djinn09/goppt/internal/pptxxml"
)

func buildTableSpec(table Table, slideNumber int) (*pptxxml.TableSpec, error) {
	styledRows, err := tableRowsWithMerges(table, slideNumber)
	if err != nil {
		return nil, err
	}
	rows := make([][]string, 0, len(styledRows))
	styledSpecRows := make([][]pptxxml.TableCellSpec, 0, len(styledRows))
	for _, srcRow := range styledRows {
		row := make([]string, len(srcRow))
		specRow := make([]pptxxml.TableCellSpec, len(srcRow))
		for i, cell := range srcRow {
			borders := cell.bordersForRender()
			row[i] = cell.Text
			specRow[i] = pptxxml.TableCellSpec{
				Text:            cell.Text,
				Bold:            cell.Bold,
				BackgroundColor: cell.BackgroundColor,
				Align:           cell.Align,
				VAlign:          cell.VAlign,
				MarginLeft:      tableMarginEMU(cell.MarginLeftPt),
				MarginRight:     tableMarginEMU(cell.MarginRightPt),
				MarginTop:       tableMarginEMU(cell.MarginTopPt),
				MarginBottom:    tableMarginEMU(cell.MarginBottomPt),
				WrapText:        cloneBoolPointer(cell.WrapText),
				RowSpan:         cell.RowSpan,
				ColSpan:         cell.ColSpan,
				VMerge:          cell.VMerge,
				HMerge:          cell.HMerge,
				BorderColor:     cell.BorderColor,
				BorderWidth:     tableBorderWidthEMU(cell.BorderWidthPt),
				BorderLeft:      toXMLTableBorderSpec(borders.Left),
				BorderRight:     toXMLTableBorderSpec(borders.Right),
				BorderTop:       toXMLTableBorderSpec(borders.Top),
				BorderBottom:    toXMLTableBorderSpec(borders.Bottom),
			}
		}
		rows = append(rows, row)
		styledSpecRows = append(styledSpecRows, specRow)
	}
	columnWidths := make([]int64, len(table.ColumnWidths))
	copy(columnWidths, table.ColumnWidths)
	rowHeights := make([]int64, len(table.RowHeights))
	copy(rowHeights, table.RowHeights)

	return &pptxxml.TableSpec{
		X:            table.X,
		Y:            table.Y,
		CX:           table.CX,
		CY:           table.CY,
		ColumnWidths: columnWidths,
		RowHeights:   rowHeights,
		Rows:         rows,
		StyledRows:   styledSpecRows,
	}, nil
}
