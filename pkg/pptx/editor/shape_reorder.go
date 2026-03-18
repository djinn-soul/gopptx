package editor

import (
	"bytes"
	"errors"
	"fmt"
)

// MoveShapeToFront moves the shape with the given ID to the front of the drawing order.
func (e *PresentationEditor) MoveShapeToFront(slideIndex, shapeID int) error {
	shapes, _, _, _, _, err := e.loadSlideShapesForReorder(slideIndex, shapeID)
	if err != nil {
		return err
	}
	return e.MoveShapeToIndex(slideIndex, shapeID, len(shapes)-1)
}

// MoveShapeToBack moves the shape with the given ID to the back of the drawing order.
func (e *PresentationEditor) MoveShapeToBack(slideIndex, shapeID int) error {
	return e.MoveShapeToIndex(slideIndex, shapeID, 0)
}

// MoveShapeToIndex moves the shape with the given ID to a specific drawing-order index.
// Index 0 is back-most, and len(shapes)-1 is front-most.
func (e *PresentationEditor) MoveShapeToIndex(slideIndex, shapeID, targetIndex int) error {
	shapes, shapeIndex, content, partPath, partFound, err := e.loadSlideShapesForReorder(slideIndex, shapeID)
	if err != nil {
		return err
	}
	if !partFound {
		return fmt.Errorf("read slide part %s: not found", partPath)
	}
	shapeCount := len(shapes)
	if targetIndex < 0 || targetIndex >= shapeCount {
		return fmt.Errorf("shape index %d out of range [0,%d)", targetIndex, shapeCount)
	}
	if shapeCount <= 1 || shapeIndex == targetIndex {
		return nil
	}

	return e.reorderShapeByIndex(content, partPath, shapes, shapeIndex, targetIndex)
}

func (e *PresentationEditor) loadSlideShapesForReorder(
	slideIndex, shapeID int,
) ([]parsedShape, int, []byte, string, bool, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, 0, nil, "", false, errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return nil, 0, nil, partPath, false, nil
	}

	shapes, err := scanShapesWithOffsets(content, false)
	if err != nil {
		return nil, 0, nil, partPath, true, fmt.Errorf("parse shapes: %w", err)
	}

	shapeIndex := -1
	for i, s := range shapes {
		if s.ID == shapeID {
			shapeIndex = i
			break
		}
	}
	if shapeIndex == -1 {
		return nil, 0, nil, partPath, true, fmt.Errorf("shape with ID %d not found", shapeID)
	}
	return shapes, shapeIndex, content, partPath, true, nil
}

func (e *PresentationEditor) reorderShapeByIndex(
	content []byte,
	partPath string,
	shapes []parsedShape,
	shapeIndex, targetIndex int,
) error {
	targetShape := shapes[shapeIndex]
	shapeXML := append([]byte(nil), content[targetShape.Start:targetShape.End]...)

	prefix := append([]byte(nil), content[:shapes[0].Start]...)
	suffix := append([]byte(nil), content[shapes[len(shapes)-1].End:]...)

	gaps := make([][]byte, 0, len(shapes)-1)
	for i := range len(shapes) - 1 {
		gap := append([]byte(nil), content[shapes[i].End:shapes[i+1].Start]...)
		gaps = append(gaps, gap)
	}

	remaining := make([]int, 0, len(shapes)-1)
	for i := range shapes {
		if i == shapeIndex {
			continue
		}
		remaining = append(remaining, i)
	}
	remaining = append(remaining[:targetIndex], append([]int{shapeIndex}, remaining[targetIndex:]...)...)

	var buf bytes.Buffer
	buf.Write(prefix)
	for pos, idx := range remaining {
		if idx == shapeIndex {
			buf.Write(shapeXML)
		} else {
			s := shapes[idx]
			buf.Write(content[s.Start:s.End])
		}
		if pos < len(gaps) {
			buf.Write(gaps[pos])
		}
	}
	buf.Write(suffix)

	e.parts.Set(partPath, buf.Bytes())
	return nil
}
