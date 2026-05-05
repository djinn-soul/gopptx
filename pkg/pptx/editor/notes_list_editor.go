package editor

import (
	"errors"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// ListNotesPlaceholders enumerates placeholder metadata from a slide's notes page.
func (e *PresentationEditor) ListNotesPlaceholders(slideIndex int) ([]Placeholder, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}

	notesPart, err := e.resolveSlideNotesPart(slideIndex)
	if err != nil {
		return nil, err
	}
	if notesPart == "" {
		return []Placeholder{}, nil
	}

	data, ok := e.parts.Get(notesPart)
	if !ok {
		return []Placeholder{}, nil
	}
	return parsePlaceholdersFromSlideXML(data), nil
}

// ListNotesShapes enumerates shape metadata from a slide's notes page.
func (e *PresentationEditor) ListNotesShapes(slideIndex int) ([]common.NotesShapeInfo, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}

	notesPart, err := e.resolveSlideNotesPart(slideIndex)
	if err != nil {
		return nil, err
	}
	if notesPart == "" {
		return []common.NotesShapeInfo{}, nil
	}

	data, ok := e.parts.Get(notesPart)
	if !ok {
		return []common.NotesShapeInfo{}, nil
	}
	parsed, err := scanShapesWithOffsets(data, false)
	if err != nil {
		return nil, err
	}
	shapes := make([]common.NotesShapeInfo, 0, len(parsed))
	for _, shape := range parsed {
		shapes = append(shapes, common.NotesShapeInfo{
			ID:                shape.ID,
			Name:              shape.Name,
			Type:              shape.Type,
			Text:              shape.Text,
			X:                 float64(shape.X),
			Y:                 float64(shape.Y),
			CX:                float64(shape.W),
			CY:                float64(shape.H),
			PlaceholderIndex:  shape.PhIndex,
			PlaceholderType:   shape.PhType,
			SupportsTextFrame: shape.TextFrame != nil || len(shape.Runs) > 0 || shape.Text != "",
		})
	}
	return shapes, nil
}
