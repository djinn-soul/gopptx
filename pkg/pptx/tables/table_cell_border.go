package tables

// TableCellBorder describes one side border for a table cell.
type TableCellBorder struct {
	WidthPt float64
	Color   string
	Dash    string
	// Cap is the line cap: TableBorderCapFlat, TableBorderCapRound or
	// TableBorderCapSquare.
	Cap string
	// Join is the line join: TableBorderJoinRound, TableBorderJoinBevel or
	// TableBorderJoinMiter.
	Join string
	// MiterLimitPct scales a miter join; ignored unless Join is miter.
	MiterLimitPct float64
	// Compound is the compound line type, e.g. TableBorderCompoundDouble.
	Compound string
	// Inset draws the pen inside the cell boundary instead of centred on it.
	Inset bool
}

// Table cell border line caps.
const (
	TableBorderCapFlat   = "flat"
	TableBorderCapRound  = "rnd"
	TableBorderCapSquare = "sq"
)

// Table cell border line joins.
const (
	TableBorderJoinRound = "round"
	TableBorderJoinBevel = "bevel"
	TableBorderJoinMiter = "miter"
)

// Table cell border compound line types.
const (
	TableBorderCompoundSingle    = "sng"
	TableBorderCompoundDouble    = "dbl"
	TableBorderCompoundThickThin = "thickThin"
	TableBorderCompoundThinThick = "thinThick"
	TableBorderCompoundTriple    = "tri"
)

// NewTableCellBorder builds a border spec that the WithSideBorderSpec setters
// accept, carrying the cap, join, compound and inset controls the
// fixed-argument helpers do not expose.
func NewTableCellBorder(widthPt float64, color string, dash string) TableCellBorder {
	return TableCellBorder{
		WidthPt: widthPt,
		Color:   NormalizeHexColor(color),
		Dash:    NormalizeTableBorderDash(dash),
	}
}

// WithCap returns the border with an explicit line cap.
func (b TableCellBorder) WithCap(lineCap string) TableCellBorder {
	b.Cap = lineCap
	return b
}

// WithJoin returns the border with an explicit line join.
func (b TableCellBorder) WithJoin(join string) TableCellBorder {
	b.Join = join
	return b
}

// WithMiterLimit returns the border with a miter limit percentage.
func (b TableCellBorder) WithMiterLimit(limitPct float64) TableCellBorder {
	b.MiterLimitPct = limitPct
	return b
}

// WithCompound returns the border with a compound line type.
func (b TableCellBorder) WithCompound(compound string) TableCellBorder {
	b.Compound = compound
	return b
}

// WithInset returns the border drawn inside the cell boundary.
func (b TableCellBorder) WithInset(inset bool) TableCellBorder {
	b.Inset = inset
	return b
}
