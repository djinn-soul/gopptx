package tables

// WithMarginsPt sets all cell text margins in points.
func (c TableCell) WithMarginsPt(points float64) TableCell {
	c = c.WithMarginLeftPt(points)
	c = c.WithMarginRightPt(points)
	c = c.WithMarginTopPt(points)
	c = c.WithMarginBottomPt(points)
	return c
}

// WithMarginLeftPt sets left text margin in points.
func (c TableCell) WithMarginLeftPt(points float64) TableCell {
	return c.withMargin(borderSideLeft, points)
}

// WithMarginRightPt sets right text margin in points.
func (c TableCell) WithMarginRightPt(points float64) TableCell {
	return c.withMargin(borderSideRight, points)
}

// WithMarginTopPt sets top text margin in points.
func (c TableCell) WithMarginTopPt(points float64) TableCell {
	return c.withMargin(borderSideTop, points)
}

// WithMarginBottomPt sets bottom text margin in points.
func (c TableCell) WithMarginBottomPt(points float64) TableCell {
	return c.withMargin(borderSideBottom, points)
}

// WithWrap sets whether text wrapping is enabled for this cell.
func (c TableCell) WithWrap(enabled bool) TableCell {
	value := enabled
	c.WrapText = &value
	return c
}

func (c TableCell) withMargin(side string, points float64) TableCell {
	value := points
	switch side {
	case borderSideLeft:
		c.MarginLeftPt = &value
	case borderSideRight:
		c.MarginRightPt = &value
	case borderSideTop:
		c.MarginTopPt = &value
	case borderSideBottom:
		c.MarginBottomPt = &value
	}
	return c
}

// CloneFloat64Pointer clones a float64 pointer.
func CloneFloat64Pointer(value *float64) *float64 {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}

// CloneBoolPointer clones a bool pointer.
func CloneBoolPointer(value *bool) *bool {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}
