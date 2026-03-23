package editor

// TableInitSpec encapsulates all table initialization configuration.
// Used with AddTableWithData to batch all table setup operations.
type TableInitSpec struct {
	// Data is a 2D array of cell text.
	// If provided, cells are populated after table creation.
	Data [][]string

	// StyledData is an alternative to Data with per-cell styling.
	// Not currently used by AddTableWithData (reserved for future).
	StyledData [][]any

	// ColumnWidths specifies width in EMU for each column (optional).
	ColumnWidths []int64

	// RowHeights specifies height in EMU for each row (optional).
	RowHeights []int64

	// FirstRow enables first-row header formatting.
	FirstRow bool

	// FirstCol enables first-column emphasis.
	FirstCol bool

	// LastRow enables last-row emphasis.
	LastRow bool

	// LastCol enables last-column emphasis.
	LastCol bool

	// BandRow enables alternating row colors.
	BandRow bool

	// BandCol enables alternating column colors.
	BandCol bool
}
