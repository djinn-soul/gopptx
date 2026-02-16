package tables

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

var hexColorPattern = regexp.MustCompile(`^[0-9A-F]{6}$`)

// NormalizeHexColor sanitizes hex color strings.
func NormalizeHexColor(color string) string {
	clean := strings.TrimPrefix(strings.TrimSpace(color), "#")
	return strings.ToUpper(clean)
}

// IsHexColor checks if a string is a valid 6-digit RGB hex color.
func IsHexColor(color string) bool {
	return hexColorPattern.MatchString(NormalizeHexColor(color))
}

func validateTable(table Table, slideIndex int) error {
	if err := validateTableDimensions(table, slideIndex); err != nil {
		return err
	}
	if err := validateTableColumns(table, slideIndex); err != nil {
		return err
	}

	rows := tableRowsForRender(table)
	if err := validateTableRowCounts(table, rows, slideIndex); err != nil {
		return err
	}

	if err := validateTableCells(table, rows, slideIndex); err != nil {
		return err
	}

	if _, err := TableRowsWithMerges(table, slideIndex); err != nil {
		return err
	}
	return nil
}

func validateTableDimensions(table Table, slideIndex int) error {
	if table.X < 0 || table.Y < 0 {
		return fmt.Errorf("slide %d table position cannot be negative", slideIndex)
	}
	if table.CX <= 0 || table.CY <= 0 {
		return fmt.Errorf("slide %d table size must be > 0", slideIndex)
	}
	return nil
}

func validateTableColumns(table Table, slideIndex int) error {
	if len(table.ColumnWidths) == 0 {
		return fmt.Errorf("slide %d table must define at least one column", slideIndex)
	}
	for columnIndex, width := range table.ColumnWidths {
		if width <= 0 {
			return fmt.Errorf("slide %d table column %d width must be > 0", slideIndex, columnIndex+1)
		}
	}
	return nil
}

func validateTableRowCounts(table Table, rows [][]TableCell, slideIndex int) error {
	if len(rows) == 0 {
		return fmt.Errorf("slide %d table must define at least one row", slideIndex)
	}
	if len(table.RowHeights) > 0 {
		if len(table.RowHeights) != len(rows) {
			return fmt.Errorf("slide %d table row heights count %d must match row count %d",
				slideIndex, len(table.RowHeights), len(rows))
		}
		for rowIndex, height := range table.RowHeights {
			if height <= 0 {
				return fmt.Errorf("slide %d table row %d height must be > 0", slideIndex, rowIndex+1)
			}
		}
	}
	return nil
}

func validateTableCells(table Table, rows [][]TableCell, slideIndex int) error {
	for rowIndex, row := range rows {
		if len(row) != len(table.ColumnWidths) {
			return fmt.Errorf("slide %d table row %d has %d cells; expected %d", slideIndex, rowIndex+1, len(row), len(table.ColumnWidths))
		}
		for cellIndex, cell := range row {
			if err := validateTableCell(cell, slideIndex, rowIndex+1, cellIndex+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateTableCell(cell TableCell, slideIndex int, rowIndex int, cellIndex int) error {
	if err := validateTableCellSpans(cell, slideIndex, rowIndex, cellIndex); err != nil {
		return err
	}
	if err := validateTableCellBasicProps(cell, slideIndex, rowIndex, cellIndex); err != nil {
		return err
	}
	if err := validateTableCellMargins(cell, slideIndex, rowIndex, cellIndex); err != nil {
		return err
	}
	return validateTableCellBorders(cell, slideIndex, rowIndex, cellIndex)
}

func validateTableCellSpans(cell TableCell, slideIndex, rowIndex, cellIndex int) error {
	if cell.RowSpan <= 0 {
		return fmt.Errorf("slide %d table row %d cell %d row span must be >= 1", slideIndex, rowIndex, cellIndex)
	}
	if cell.ColSpan <= 0 {
		return fmt.Errorf("slide %d table row %d cell %d col span must be >= 1", slideIndex, rowIndex, cellIndex)
	}
	return nil
}

func validateTableCellBasicProps(cell TableCell, slideIndex, rowIndex, cellIndex int) error {
	if color := strings.TrimSpace(cell.BackgroundColor); color != "" && !IsHexColor(color) {
		return fmt.Errorf("slide %d table row %d cell %d background color must be 6-digit RGB hex", slideIndex, rowIndex, cellIndex)
	}
	if align := strings.TrimSpace(cell.Align); align != "" && !isTableAlign(align) {
		return fmt.Errorf("slide %d table row %d cell %d align must be one of l|ctr|r|just", slideIndex, rowIndex, cellIndex)
	}
	if vAlign := strings.TrimSpace(cell.VAlign); vAlign != "" && !isTableVAlign(vAlign) {
		return fmt.Errorf("slide %d table row %d cell %d valign must be one of t|ctr|b", slideIndex, rowIndex, cellIndex)
	}
	return nil
}

func validateTableCellMargins(cell TableCell, slideIndex, rowIndex, cellIndex int) error {
	if err := validateTableCellMargin(cell.MarginLeftPt, slideIndex, rowIndex, cellIndex, "left"); err != nil {
		return err
	}
	if err := validateTableCellMargin(cell.MarginRightPt, slideIndex, rowIndex, cellIndex, "right"); err != nil {
		return err
	}
	if err := validateTableCellMargin(cell.MarginTopPt, slideIndex, rowIndex, cellIndex, "top"); err != nil {
		return err
	}
	return validateTableCellMargin(cell.MarginBottomPt, slideIndex, rowIndex, cellIndex, "bottom")
}

func validateTableCellBorders(cell TableCell, slideIndex, rowIndex, cellIndex int) error {
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
		return validateTableCellBorder(cell.BorderBottom, slideIndex, rowIndex, cellIndex, borderSideBottom)
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
		return fmt.Errorf(
			"slide %d table row %d cell %d %s width must be finite",
			slideIndex,
			rowIndex,
			cellIndex,
			fieldPrefix,
		)
	}
	if border.WidthPt < 0 {
		return fmt.Errorf(
			"slide %d table row %d cell %d %s width must be >= 0",
			slideIndex,
			rowIndex,
			cellIndex,
			fieldPrefix,
		)
	}
	if border.WidthPt > 0 && !IsHexColor(border.Color) {
		return fmt.Errorf(
			"slide %d table row %d cell %d %s color must be 6-digit RGB hex",
			slideIndex,
			rowIndex,
			cellIndex,
			fieldPrefix,
		)
	}
	if NormalizeHexColor(border.Color) != "" && border.WidthPt <= 0 {
		return fmt.Errorf(
			"slide %d table row %d cell %d %s width must be > 0 when %s color is set",
			slideIndex,
			rowIndex,
			cellIndex,
			fieldPrefix,
			fieldPrefix,
		)
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

func validateTableCellMargin(marginPt *float64, slideIndex int, rowIndex int, cellIndex int, side string) error {
	if marginPt == nil {
		return nil
	}
	if math.IsNaN(*marginPt) || math.IsInf(*marginPt, 0) {
		return fmt.Errorf(
			"slide %d table row %d cell %d %s margin must be finite",
			slideIndex,
			rowIndex,
			cellIndex,
			side,
		)
	}
	if *marginPt < 0 {
		return fmt.Errorf(
			"slide %d table row %d cell %d %s margin must be >= 0",
			slideIndex,
			rowIndex,
			cellIndex,
			side,
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

// NormalizeTableBorderDash sanitizes table border dash styles.
func NormalizeTableBorderDash(dash string) string {
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
	switch NormalizeTableBorderDash(dash) {
	case TableBorderDashSolid, TableBorderDashDash, TableBorderDashDot, TableBorderDashDashDot, TableBorderDashLongDash:
		return true
	default:
		return false
	}
}
