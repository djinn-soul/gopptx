package tables

import (
	"testing"
)

func TestTableCell_FluentAPI(t *testing.T) {
	c := TableCell{Text: "Cell"}.
		WithBackgroundColor("FF0000").
		WithBold(true).
		WithAlign(TableAlignLeft).
		WithAlignCenter().
		WithAlignRight().
		WithAlignJustify().
		WithAlignLeft().
		WithVAlign(TableVAlignTop).
		WithVAlignMiddle().
		WithVAlignBottom().
		WithRowSpan(2).
		WithColSpan(3).
		WithBorder(1, "000000").
		WithBorderStyle(2, "00FF00", "dashed").
		WithLeftBorder(1, "111111").
		WithLeftBorderStyle(1, "222222", "dotted").
		WithRightBorder(1, "333333").
		WithRightBorderStyle(1, "444444", "dash-dot").
		WithTopBorder(1, "555555").
		WithTopBorderStyle(1, "666666", "lg-dash").
		WithBottomBorder(1, "777777").
		WithBottomBorderStyle(1, "888888", "solid")

	if c.BackgroundColor != "FF0000" || !c.Bold || c.Align != TableAlignLeft || c.VAlign != TableVAlignBottom {
		t.Error("Basic cell props failed")
	}
	if c.RowSpan != 2 || c.ColSpan != 3 { t.Error("Span failed") }
	if c.BorderLeft == nil || c.BorderLeft.Color != "222222" { t.Error("Borders failed") }
}

func TestTableCell_LayoutAPI(t *testing.T) {
	c := TableCell{Text: "Layout"}.
		WithMarginsPt(10).
		WithMarginLeftPt(5).
		WithMarginRightPt(6).
		WithMarginTopPt(7).
		WithMarginBottomPt(8).
		WithWrap(false)
	
	if *c.MarginLeftPt != 5 || *c.MarginBottomPt != 8 || *c.WrapText != false {
		t.Error("Layout props failed")
	}
}

func TestTable_Normalization_Helpers(t *testing.T) {
	// These are unexported but we are in package tables
	if normalizeTableAlign("CTR") != TableAlignCenter { t.Error("Align normalization failed") }
	if normalizeTableVAlign("CTR") != TableVAlignMiddle { t.Error("VAlign normalization failed") }
	if NormalizeTableBorderDash("DASH") != "dash" { t.Error("Dash normalization failed") }
}
