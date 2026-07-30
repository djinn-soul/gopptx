package editor

import (
	"errors"
	"fmt"
)

// Shape lookups used by the text-state operations. They parse the slide part
// and resolve its relationships, so run-level hyperlinks survive a read.

func (e *PresentationEditor) getShapeForTextOps(slideIndex, shapeID int) (parsedShape, error) {
	shapes, err := e.getShapesForTextOps(slideIndex)
	if err != nil {
		return parsedShape{}, err
	}

	for _, shape := range shapes {
		if shape.ID == shapeID || (shape.PhType == placeholderTypeTitle && shapeID == 0) {
			return shape, nil
		}
	}

	return parsedShape{}, fmt.Errorf("shape id %d not found on slide %d", shapeID, slideIndex)
}

func (e *PresentationEditor) getShapesForTextOps(slideIndex int) ([]parsedShape, error) {
	if e == nil {
		return nil, errors.New("editor cannot be nil")
	}
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return nil, fmt.Errorf("read slide part %s: not found", partPath)
	}

	shapes, err := parseSlideShapes(content)
	if err != nil {
		return nil, fmt.Errorf("parse shapes: %w", err)
	}
	// Resolve run-level hyperlink relationship IDs so text state round-trips
	// addresses and slide jumps instead of dropping them.
	if err := e.enrichParsedShapeRelationships(partPath, shapes); err != nil {
		return nil, fmt.Errorf("resolve shape relationships: %w", err)
	}
	return shapes, nil
}
