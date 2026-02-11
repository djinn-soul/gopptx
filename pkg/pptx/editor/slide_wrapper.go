package editor

import (
	"fmt"
)

// Slide is a high-level wrapper around an editable slide.
type Slide struct {
	ID       int64
	PartName string
	editor   *PresentationEditor

	// Placeholders []Placeholder
}

// GetSlide returns a Slide object for the given index (0-based).
func (e *PresentationEditor) GetSlide(index int) (*Slide, error) {
	if index < 0 || index >= len(e.slides) {
		return nil, fmt.Errorf("slide index %d out of range", index)
	}
	ref := e.slides[index]
	return &Slide{
		ID:       ref.SlideID,
		PartName: ref.Part,
		editor:   e,
	}, nil
}

// Placeholder is a simplified version for this refactor.
// In a real implementation, this would be more complex.
type Placeholder struct {
	Index int
	Type  string
}
