package tables

// TableCellBorder describes one side border for a table cell.
type TableCellBorder struct {
	WidthPt float64
	Color   string
	Dash    string
}

// TableCell stores text and optional style for one table cell.
type TableCell struct {
	Text            string
	Bold            bool
	BackgroundColor string
	Align           string
	VAlign          string
	MarginLeftPt    *float64
	MarginRightPt   *float64
	MarginTopPt     *float64
	MarginBottomPt  *float64
	WrapText        *bool

	// Merge model fields.
	RowSpan int
	ColSpan int
	VMerge  bool
	HMerge  bool

	// Legacy uniform border fields (kept for backward compatibility).
	BorderColor   string
	BorderWidthPt float64

	BorderLeft   *TableCellBorder
	BorderRight  *TableCellBorder
	BorderTop    *TableCellBorder
	BorderBottom *TableCellBorder
}

type tableCellBorders struct {
	Left   *TableCellBorder
	Right  *TableCellBorder
	Top    *TableCellBorder
	Bottom *TableCellBorder
}

// NewTableCell creates a styled table cell with text.
func NewTableCell(text string) TableCell {
	return TableCell{Text: text, RowSpan: 1, ColSpan: 1}
}

// WithBold sets bold text for this cell.
func (c TableCell) WithBold(enabled bool) TableCell {
	c.Bold = enabled
	return c
}

// WithBackgroundColor sets cell background fill using RGB hex.
func (c TableCell) WithBackgroundColor(color string) TableCell {
	c.BackgroundColor = NormalizeHexColor(color)
	return c
}

// WithAlign sets horizontal text alignment.
func (c TableCell) WithAlign(align string) TableCell {
	c.Align = normalizeTableAlign(align)
	return c
}

// WithAlignLeft sets horizontal text alignment to left.
func (c TableCell) WithAlignLeft() TableCell {
	return c.WithAlign(TableAlignLeft)
}

// WithAlignCenter sets horizontal text alignment to center.
func (c TableCell) WithAlignCenter() TableCell {
	return c.WithAlign(TableAlignCenter)
}

// WithAlignRight sets horizontal text alignment to right.
func (c TableCell) WithAlignRight() TableCell {
	return c.WithAlign(TableAlignRight)
}

// WithAlignJustify sets horizontal text alignment to justify.
func (c TableCell) WithAlignJustify() TableCell {
	return c.WithAlign(TableAlignJustify)
}

// WithVAlign sets vertical text alignment.
func (c TableCell) WithVAlign(vAlign string) TableCell {
	c.VAlign = normalizeTableVAlign(vAlign)
	return c
}

// WithVAlignTop sets vertical text alignment to top.
func (c TableCell) WithVAlignTop() TableCell {
	return c.WithVAlign(TableVAlignTop)
}

// WithVAlignMiddle sets vertical text alignment to middle.
func (c TableCell) WithVAlignMiddle() TableCell {
	return c.WithVAlign(TableVAlignMiddle)
}

// WithVAlignBottom sets vertical text alignment to bottom.
func (c TableCell) WithVAlignBottom() TableCell {
	return c.WithVAlign(TableVAlignBottom)
}

// WithRowSpan sets the number of rows merged downward from this anchor cell.
func (c TableCell) WithRowSpan(span int) TableCell {
	c.RowSpan = span
	return c
}

// WithColSpan sets the number of columns merged rightward from this anchor cell.
func (c TableCell) WithColSpan(span int) TableCell {
	c.ColSpan = span
	return c
}

// WithBorder sets uniform border with default dash style.
func (c TableCell) WithBorder(widthPt float64, color string) TableCell {
	return c.WithBorderStyle(widthPt, color, TableBorderDashSolid)
}

func (c TableCell) WithBorderStyle(widthPt float64, color string, dash string) TableCell {
	normalizedColor := NormalizeHexColor(color)
	normalizedDash := NormalizeTableBorderDash(dash)
	c.BorderWidthPt = widthPt
	c.BorderColor = normalizedColor
	c.BorderLeft = &TableCellBorder{WidthPt: widthPt, Color: normalizedColor, Dash: normalizedDash}
	c.BorderRight = &TableCellBorder{WidthPt: widthPt, Color: normalizedColor, Dash: normalizedDash}
	c.BorderTop = &TableCellBorder{WidthPt: widthPt, Color: normalizedColor, Dash: normalizedDash}
	c.BorderBottom = &TableCellBorder{WidthPt: widthPt, Color: normalizedColor, Dash: normalizedDash}
	return c
}

// WithLeftBorder sets left border with default dash style.
func (c TableCell) WithLeftBorder(widthPt float64, color string) TableCell {
	return c.WithLeftBorderStyle(widthPt, color, TableBorderDashSolid)
}

// WithLeftBorderStyle sets left border with explicit dash style.
func (c TableCell) WithLeftBorderStyle(widthPt float64, color string, dash string) TableCell {
	return c.withSideBorder(borderSideLeft, widthPt, color, dash)
}

// WithRightBorder sets right border with default dash style.
func (c TableCell) WithRightBorder(widthPt float64, color string) TableCell {
	return c.WithRightBorderStyle(widthPt, color, TableBorderDashSolid)
}

// WithRightBorderStyle sets right border with explicit dash style.
func (c TableCell) WithRightBorderStyle(widthPt float64, color string, dash string) TableCell {
	return c.withSideBorder(borderSideRight, widthPt, color, dash)
}

// WithTopBorder sets top border with default dash style.
func (c TableCell) WithTopBorder(widthPt float64, color string) TableCell {
	return c.WithTopBorderStyle(widthPt, color, TableBorderDashSolid)
}

// WithTopBorderStyle sets top border with explicit dash style.
func (c TableCell) WithTopBorderStyle(widthPt float64, color string, dash string) TableCell {
	return c.withSideBorder(borderSideTop, widthPt, color, dash)
}

// WithBottomBorder sets bottom border with default dash style.
func (c TableCell) WithBottomBorder(widthPt float64, color string) TableCell {
	return c.WithBottomBorderStyle(widthPt, color, TableBorderDashSolid)
}

// WithBottomBorderStyle sets bottom border with explicit dash style.
func (c TableCell) WithBottomBorderStyle(widthPt float64, color string, dash string) TableCell {
	return c.withSideBorder(borderSideBottom, widthPt, color, dash)
}

func (c TableCell) withSideBorder(side string, widthPt float64, color string, dash string) TableCell {
	border := &TableCellBorder{
		WidthPt: widthPt,
		Color:   NormalizeHexColor(color),
		Dash:    NormalizeTableBorderDash(dash),
	}
	switch side {
	case borderSideLeft:
		c.BorderLeft = border
	case borderSideRight:
		c.BorderRight = border
	case borderSideTop:
		c.BorderTop = border
	case borderSideBottom:
		c.BorderBottom = border
	}
	return c
}

func (c TableCell) bordersForRender() tableCellBorders {
	borders := tableCellBorders{
		Left:   cloneTableCellBorder(c.BorderLeft),
		Right:  cloneTableCellBorder(c.BorderRight),
		Top:    cloneTableCellBorder(c.BorderTop),
		Bottom: cloneTableCellBorder(c.BorderBottom),
	}
	if borders.Left == nil && borders.Right == nil && borders.Top == nil && borders.Bottom == nil {
		if legacy := c.uniformLegacyBorder(); legacy != nil {
			borders.Left = cloneTableCellBorder(legacy)
			borders.Right = cloneTableCellBorder(legacy)
			borders.Top = cloneTableCellBorder(legacy)
			borders.Bottom = cloneTableCellBorder(legacy)
		}
	}
	return borders
}

func (c TableCell) hasExplicitBorderSides() bool {
	return c.BorderLeft != nil || c.BorderRight != nil || c.BorderTop != nil || c.BorderBottom != nil
}

func (c TableCell) uniformLegacyBorder() *TableCellBorder {
	if c.BorderWidthPt <= 0 && NormalizeHexColor(c.BorderColor) == "" {
		return nil
	}
	return &TableCellBorder{
		WidthPt: c.BorderWidthPt,
		Color:   NormalizeHexColor(c.BorderColor),
		Dash:    TableBorderDashSolid,
	}
}

func cloneTableCellBorder(border *TableCellBorder) *TableCellBorder {
	if border == nil {
		return nil
	}
	clone := *border
	return &clone
}
