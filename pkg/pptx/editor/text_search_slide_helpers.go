package editor

import (
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

// Slide is a high-level wrapper around an editable slide.
type Slide struct {
	ID       int64
	PartName string
	editor   *PresentationEditor
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

func (e *PresentationEditor) slideRelationships(slidePart string) ([]common.EditorRelationship, error) {
	return editorslide.Relationships(slidePart, e.parts.Get, parseRelationshipsXML)
}

// Placeholder describes a discovered placeholder in an existing slide.
type Placeholder struct {
	Index int
	Type  string
	Name  string
}

// Placeholders parses the slide XML and returns all placeholder elements found.
func (s *Slide) Placeholders() ([]Placeholder, error) {
	content, ok := s.editor.parts.Get(s.PartName)
	if !ok {
		return nil, fmt.Errorf("slide part %q not found", s.PartName)
	}
	return parsePlaceholdersFromSlideXML(content), nil
}

func parsePlaceholdersFromSlideXML(content []byte) []Placeholder {
	parsed, _ := scanShapesWithOffsets(content, false)
	var result []Placeholder
	for _, s := range parsed {
		if s.PhIndex != -1 {
			result = append(result, Placeholder{
				Index: s.PhIndex,
				Type:  s.PhType,
				Name:  s.Name,
			})
		}
	}
	return result
}
