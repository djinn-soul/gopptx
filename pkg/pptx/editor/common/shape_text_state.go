package editorcommon

// ShapeTextState is a snapshot of a shape's text-related state.
type ShapeTextState struct {
	Text      string
	Runs      []TextRun
	TextFrame *TextFrame
	Paragraph *Paragraph
}

// SlideShapeTextState identifies one shape and its text-related state on a slide.
type SlideShapeTextState struct {
	ShapeID   int        `json:"shape_id"`
	Text      string     `json:"text"`
	Runs      []TextRun  `json:"runs"`
	TextFrame *TextFrame `json:"text_frame,omitempty"`
	Paragraph *Paragraph `json:"paragraph,omitempty"`
}

// ShapeRunTextUpdate describes one run-text mutation on a slide shape.
type ShapeRunTextUpdate struct {
	ShapeID  int    `json:"shape_id"`
	RunIndex int    `json:"run_index"`
	Text     string `json:"text"`
}

// SlideRunTextUpdates describes a set of run-text mutations for one slide.
type SlideRunTextUpdates struct {
	SlideIndex int                  `json:"slide_index"`
	Updates    []ShapeRunTextUpdate `json:"updates"`
}
