package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

// Table aliases for backward compatibility
type (
	Table     = tables.Table
	TableCell = tables.TableCell
)

// Function aliases
var (
	NewTable     = tables.NewTable
	NewTableCell = tables.NewTableCell
)

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
