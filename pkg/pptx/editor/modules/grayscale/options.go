package grayscale

import "github.com/djinn-soul/gopptx/pkg/pptx/shapes"

// ShapeRef identifies one shape on one slide.
type ShapeRef struct {
	SlideIndex int `json:"slide_index"`
	ShapeID    int `json:"shape_id"`
}

// TextRef identifies one shape text target and optional run subset.
type TextRef struct {
	SlideIndex int   `json:"slide_index"`
	ShapeID    int   `json:"shape_id"`
	RunIndices []int `json:"run_indices,omitempty"`
}

// PlaceholderRef identifies one or more placeholder-backed shapes on one slide.
type PlaceholderRef struct {
	SlideIndex int                    `json:"slide_index"`
	Type       shapes.PlaceholderType `json:"type,omitempty"`
	Index      *int                   `json:"index,omitempty"`
}

// Options controls which presentation parts are converted to grayscale.
type Options struct {
	Slides       []int            `json:"slides,omitempty"`
	Shapes       []ShapeRef       `json:"shapes,omitempty"`
	Text         []TextRef        `json:"text,omitempty"`
	Placeholders []PlaceholderRef `json:"placeholders,omitempty"`
	Colors       bool             `json:"colors,omitempty"`
	Images       bool             `json:"images,omitempty"`
	Backgrounds  bool             `json:"backgrounds,omitempty"`
}
