package editorcommon

// Sources reported by an effective-style lookup, in inheritance order. A
// placeholder that sets nothing of its own still draws with a colour, a size
// and a typeface; these say where each resolved value actually came from
// (upstream #1013).
const (
	StyleSourceShape  = "shape"
	StyleSourceLayout = "layout"
	StyleSourceMaster = "master"
	StyleSourceTheme  = "theme"
)

// EffectiveShapeStyle is the resolved appearance of one shape, after walking
// shape -> layout placeholder -> master placeholder / txStyles -> theme.
// A nil field means nothing in that chain defined the value.
type EffectiveShapeStyle struct {
	FillColor    *EffectiveColor    `json:"fill_color,omitempty"`
	FontColor    *EffectiveColor    `json:"font_color,omitempty"`
	FontTypeface *EffectiveString   `json:"font_typeface,omitempty"`
	FontSizePt   *EffectiveFloat    `json:"font_size_pt,omitempty"`
	Bold         *EffectiveBool     `json:"bold,omitempty"`
	Italic       *EffectiveBool     `json:"italic,omitempty"`
	Position     *EffectivePosition `json:"position,omitempty"`
	// LayoutPart and MasterPart name the parts that were consulted, so a caller
	// can inspect them directly.
	LayoutPart string `json:"layout_part,omitempty"`
	MasterPart string `json:"master_part,omitempty"`
}

// EffectiveColor is a resolved colour. RGB is always a concrete six-digit hex
// value; SchemeSlot records the theme slot when the source referenced one.
type EffectiveColor struct {
	RGB        string `json:"rgb"`
	SchemeSlot string `json:"scheme_slot,omitempty"`
	Source     string `json:"source"`
}

// EffectiveString is a resolved string value and where it came from.
type EffectiveString struct {
	Value  string `json:"value"`
	Source string `json:"source"`
}

// EffectiveFloat is a resolved numeric value and where it came from.
type EffectiveFloat struct {
	Value  float64 `json:"value"`
	Source string  `json:"source"`
}

// EffectiveBool is a resolved flag and where it came from.
type EffectiveBool struct {
	Value  bool   `json:"value"`
	Source string `json:"source"`
}

// EffectivePosition is a resolved position and extent in EMU. A placeholder
// that overrides neither inherits both from its layout or master placeholder.
type EffectivePosition struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	W      int    `json:"w"`
	H      int    `json:"h"`
	Source string `json:"source"`
}
