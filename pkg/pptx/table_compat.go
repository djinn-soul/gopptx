package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

type (
	Table     = tables.Table
	TableCell = tables.TableCell
)

func NewTable(columnWidths []int64) Table {
	return tables.NewTable(columnWidths)
}

func NewTableCell(text string) TableCell {
	return tables.NewTableCell(text)
}

// Table Constants
const (
	TableAlignLeft    = tables.TableAlignLeft
	TableAlignCenter  = tables.TableAlignCenter
	TableAlignRight   = tables.TableAlignRight
	TableAlignJustify = tables.TableAlignJustify

	TableVAlignTop    = tables.TableVAlignTop
	TableVAlignMiddle = tables.TableVAlignMiddle
	TableVAlignBottom = tables.TableVAlignBottom

	TableBorderDashSolid    = tables.TableBorderDashSolid
	TableBorderDashDash     = tables.TableBorderDashDash
	TableBorderDashDot      = tables.TableBorderDashDot
	TableBorderDashDashDot  = tables.TableBorderDashDashDot
	TableBorderDashLongDash = tables.TableBorderDashLongDash
)
