package pptx

import "fmt"

// Slide is a high-level wrapper around an editable slide.
type Slide struct {
	ID       int64
	PartName string
	editor   *PresentationEditor

	// Parsed content
	Placeholders []Placeholder
	// Shapes     []Shape
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

// GetPlaceholder returns the placeholder at the given index.
func (s *Slide) GetPlaceholder(idx int) (*Placeholder, error) {
	if s.Placeholders == nil {
		if err := s.loadContent(); err != nil {
			return nil, err
		}
	}
	for i := range s.Placeholders {
		if s.Placeholders[i].Index == idx {
			return &s.Placeholders[i], nil
		}
	}
	return nil, fmt.Errorf("placeholder with idx %d not found", idx)
}

func (s *Slide) loadContent() error {
	content, ok := s.editor.parts[s.PartName]
	if !ok {
		return fmt.Errorf("part %q not found", s.PartName)
	}

	placeholders, err := ParseSlidePlaceholders(content)
	if err != nil {
		return err
	}
	s.Placeholders = placeholders
	return nil
}
