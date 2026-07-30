package editorcommon

// SlideImageRef describes one image relationship on a slide.
type SlideImageRef struct {
	Index  int
	RelID  string
	Target string
}

// SlideMediaRef describes one media relationship on a slide: an image, a sound
// or a movie. Images already had a listing; audio and video did not, so an
// embedded movie could not be found without walking relationships by hand
// (upstream #1049).
type SlideMediaRef struct {
	Index int    `json:"index"`
	RelID string `json:"rel_id"`
	// Kind is "image", "audio" or "video".
	Kind string `json:"kind"`
	// Target is the raw relationship target; PartPath is that target resolved
	// to a package part, and is empty for an external link.
	Target      string `json:"target"`
	PartPath    string `json:"part_path,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	SizeBytes   int    `json:"size_bytes,omitempty"`
	External    bool   `json:"external,omitempty"`
}

// ShapeAdjustmentValue sets one preset-geometry adjustment — a yellow handle in
// PowerPoint's UI (upstream #1017). Value is a fraction in the same units the
// reader reports, so 0.5 is the halfway point; Formula overrides it when a
// caller needs a raw OOXML guide expression.
type ShapeAdjustmentValue struct {
	Name    string  `json:"name"`
	Value   float64 `json:"value,omitempty"`
	Formula string  `json:"formula,omitempty"`
}
