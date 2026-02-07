package pptx

import (
	"fmt"
	"math"
	"strings"
)

func validateTable(table Table, slideIndex int) error {
	if table.X < 0 || table.Y < 0 {
		return fmt.Errorf("slide %d table position cannot be negative", slideIndex)
	}
	if table.CX <= 0 || table.CY <= 0 {
		return fmt.Errorf("slide %d table size must be > 0", slideIndex)
	}
	if len(table.ColumnWidths) == 0 {
		return fmt.Errorf("slide %d table must define at least one column", slideIndex)
	}
	for columnIndex, width := range table.ColumnWidths {
		if width <= 0 {
			return fmt.Errorf("slide %d table column %d width must be > 0", slideIndex, columnIndex+1)
		}
	}

	rows := tableRowsForRender(table)
	if len(rows) == 0 {
		return fmt.Errorf("slide %d table must define at least one row", slideIndex)
	}
	for rowIndex, row := range rows {
		if len(row) != len(table.ColumnWidths) {
			return fmt.Errorf(
				"slide %d table row %d has %d cells; expected %d",
				slideIndex,
				rowIndex+1,
				len(row),
				len(table.ColumnWidths),
			)
		}
		for cellIndex, cell := range row {
			if err := validateTableCell(cell, slideIndex, rowIndex+1, cellIndex+1); err != nil {
				return err
			}
		}
	}
	if _, err := tableRowsWithMerges(table, slideIndex); err != nil {
		return err
	}
	return nil
}

func validateTableCell(cell TableCell, slideIndex int, rowIndex int, cellIndex int) error {
	if cell.RowSpan <= 0 {
		return fmt.Errorf("slide %d table row %d cell %d row span must be >= 1", slideIndex, rowIndex, cellIndex)
	}
	if cell.ColSpan <= 0 {
		return fmt.Errorf("slide %d table row %d cell %d col span must be >= 1", slideIndex, rowIndex, cellIndex)
	}
	if color := strings.TrimSpace(cell.BackgroundColor); color != "" && !isHexColor(color) {
		return fmt.Errorf("slide %d table row %d cell %d background color must be 6-digit RGB hex", slideIndex, rowIndex, cellIndex)
	}
	if align := strings.TrimSpace(cell.Align); align != "" && !isTableAlign(align) {
		return fmt.Errorf("slide %d table row %d cell %d align must be one of l|ctr|r|just", slideIndex, rowIndex, cellIndex)
	}
	if vAlign := strings.TrimSpace(cell.VAlign); vAlign != "" && !isTableVAlign(vAlign) {
		return fmt.Errorf("slide %d table row %d cell %d valign must be one of t|ctr|b", slideIndex, rowIndex, cellIndex)
	}

	if cell.hasExplicitBorderSides() {
		if err := validateTableCellBorder(cell.BorderLeft, slideIndex, rowIndex, cellIndex, borderSideLeft); err != nil {
			return err
		}
		if err := validateTableCellBorder(cell.BorderRight, slideIndex, rowIndex, cellIndex, borderSideRight); err != nil {
			return err
		}
		if err := validateTableCellBorder(cell.BorderTop, slideIndex, rowIndex, cellIndex, borderSideTop); err != nil {
			return err
		}
		if err := validateTableCellBorder(cell.BorderBottom, slideIndex, rowIndex, cellIndex, borderSideBottom); err != nil {
			return err
		}
		return nil
	}

	return validateTableCellBorder(cell.uniformLegacyBorder(), slideIndex, rowIndex, cellIndex, "")
}

func validateTableCellBorder(border *TableCellBorder, slideIndex int, rowIndex int, cellIndex int, side string) error {
	if border == nil {
		return nil
	}

	fieldPrefix := "border"
	if side != "" {
		fieldPrefix = side + " border"
	}

	if math.IsNaN(border.WidthPt) || math.IsInf(border.WidthPt, 0) {
		return fmt.Errorf("slide %d table row %d cell %d %s width must be finite", slideIndex, rowIndex, cellIndex, fieldPrefix)
	}
	if border.WidthPt < 0 {
		return fmt.Errorf("slide %d table row %d cell %d %s width must be >= 0", slideIndex, rowIndex, cellIndex, fieldPrefix)
	}
	if border.WidthPt > 0 && !isHexColor(border.Color) {
		return fmt.Errorf("slide %d table row %d cell %d %s color must be 6-digit RGB hex", slideIndex, rowIndex, cellIndex, fieldPrefix)
	}
	if normalizeHexColor(border.Color) != "" && border.WidthPt <= 0 {
		return fmt.Errorf("slide %d table row %d cell %d %s width must be > 0 when %s color is set", slideIndex, rowIndex, cellIndex, fieldPrefix, fieldPrefix)
	}
	if dash := strings.TrimSpace(border.Dash); dash != "" && !isTableBorderDash(dash) {
		return fmt.Errorf(
			"slide %d table row %d cell %d %s dash style must be one of solid|dash|dot|dashDot|lgDash",
			slideIndex,
			rowIndex,
			cellIndex,
			fieldPrefix,
		)
	}
	return nil
}

func normalizeTableAlign(align string) string {
	return strings.ToLower(strings.TrimSpace(align))
}

func normalizeTableVAlign(vAlign string) string {
	return strings.ToLower(strings.TrimSpace(vAlign))
}

func isTableAlign(align string) bool {
	switch normalizeTableAlign(align) {
	case TableAlignLeft, TableAlignCenter, TableAlignRight, TableAlignJustify:
		return true
	default:
		return false
	}
}

func isTableVAlign(vAlign string) bool {
	switch normalizeTableVAlign(vAlign) {
	case TableVAlignTop, TableVAlignMiddle, TableVAlignBottom:
		return true
	default:
		return false
	}
}

func normalizeTableBorderDash(dash string) string {
	switch strings.ToLower(strings.TrimSpace(dash)) {
	case "", "solid":
		return TableBorderDashSolid
	case "dash":
		return TableBorderDashDash
	case "dot":
		return TableBorderDashDot
	case "dashdot", "dash-dot", "dash_dot":
		return TableBorderDashDashDot
	case "lgdash", "lg-dash", "longdash", "long-dash", "long_dash":
		return TableBorderDashLongDash
	default:
		return strings.TrimSpace(dash)
	}
}

func isTableBorderDash(dash string) bool {
	switch normalizeTableBorderDash(dash) {
	case TableBorderDashSolid, TableBorderDashDash, TableBorderDashDot, TableBorderDashDashDot, TableBorderDashLongDash:
		return true
	default:
		return false
	}
}
