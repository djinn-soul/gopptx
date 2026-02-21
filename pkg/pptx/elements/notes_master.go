package elements

import "fmt"

// NotesMaster defines the configuration for the notes master part.
type NotesMaster struct {
	HeaderText   string
	FooterText   string
	ShowDateTime bool
	ShowSlideNum bool
	Background   *SlideBackground
	BodyStyle    []TextLevelStyle
}

// NewNotesMaster creates a new NotesMaster with defaults.
func NewNotesMaster() *NotesMaster {
	return &NotesMaster{
		ShowDateTime: true,
		ShowSlideNum: true,
	}
}

// WithHeader sets the header text for the notes master.
func (n *NotesMaster) WithHeader(text string) *NotesMaster {
	n.HeaderText = text
	return n
}

// WithFooter sets the footer text for the notes master.
func (n *NotesMaster) WithFooter(text string) *NotesMaster {
	n.FooterText = text
	return n
}

// WithDateTime enables or disables the date/time placeholder.
func (n *NotesMaster) WithDateTime(show bool) *NotesMaster {
	n.ShowDateTime = show
	return n
}

// WithSlideNumber enables or disables the slide number placeholder.
func (n *NotesMaster) WithSlideNumber(show bool) *NotesMaster {
	n.ShowSlideNum = show
	return n
}

// WithBackground sets the background for the notes master.
func (n *NotesMaster) WithBackground(bg SlideBackground) *NotesMaster {
	n.Background = &bg
	return n
}

// WithBodyStyle sets the text styles for the notes body.
func (n *NotesMaster) WithBodyStyle(styles []TextLevelStyle) *NotesMaster {
	n.BodyStyle = styles
	return n
}

// Validate checks notes master configuration for unsupported values.
func (n *NotesMaster) Validate() error {
	if n == nil || n.Background == nil {
		return nil
	}
	if err := n.Background.Validate(); err != nil {
		return fmt.Errorf("invalid notes master background: %w", err)
	}
	return nil
}
