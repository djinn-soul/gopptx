package editorcommon

// ShapeTextState is a snapshot of a shape's text-related state.
type ShapeTextState struct {
	Text      string
	Runs      []TextRun
	TextFrame *TextFrame
	Paragraph *Paragraph
}
