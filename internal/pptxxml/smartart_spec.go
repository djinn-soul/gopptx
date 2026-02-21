package pptxxml

// SmartArtNodeSpec represents a node in the internal SmartArt data model.
type SmartArtNodeSpec struct {
	ModelID  string
	Text     string
	Children []SmartArtNodeSpec
}

// SmartArtSpec describes the internal representation of a SmartArt diagram
// for XML generation.
type SmartArtSpec struct {
	LayoutURI    string
	ColorStyleID string
	QuickStyleID string
	Nodes        []SmartArtNodeSpec
	X, Y, CX, CY int64

	// Accessibility
	AltText      string
	IsDecorative bool
}

// SmartArtFrame describes SmartArt placement in slide XML.
// It holds the 4 relationship IDs that dgm:relIds requires.
type SmartArtFrame struct {
	X           int64
	Y           int64
	CX          int64
	CY          int64
	DataRelID   string // r:dm
	LayoutRelID string // r:lo
	StyleRelID  string // r:qs
	ColorRelID  string // r:cs

	// Accessibility
	AltText      string
	IsDecorative bool
}

// SmartArtRel describes one SmartArt relationship entry for slide relationships XML.
type SmartArtRel struct {
	RID    string
	Target string
	Type   string // relationship type URI
}
