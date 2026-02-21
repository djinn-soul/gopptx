package editor

import (
	"bytes"
	"errors"
	"fmt"
)

// MoveShapeToFront moves the shape with the given ID to the front of the drawing order.
func (e *PresentationEditor) MoveShapeToFront(slideIndex, shapeID int) error {
	return e.moveShape(slideIndex, shapeID, true)
}

// MoveShapeToBack moves the shape with the given ID to the back of the drawing order.
func (e *PresentationEditor) MoveShapeToBack(slideIndex, shapeID int) error {
	return e.moveShape(slideIndex, shapeID, false)
}

func (e *PresentationEditor) moveShape(slideIndex, shapeID int, toFront bool) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return fmt.Errorf("read slide part %s: not found", partPath)
	}

	shapes, err := scanShapesWithOffsets(content, false)
	if err != nil {
		return fmt.Errorf("parse shapes: %w", err)
	}

	shapeIndex := -1
	for i, s := range shapes {
		if s.ID == shapeID {
			shapeIndex = i
			break
		}
	}

	if shapeIndex == -1 {
		return fmt.Errorf("shape with ID %d not found", shapeID)
	}

	if len(shapes) <= 1 {
		return nil // Nothing to move
	}

	if toFront && shapeIndex == len(shapes)-1 {
		return nil // Already at front
	}
	if !toFront && shapeIndex == 0 {
		return nil // Already at back
	}

	targetShape := shapes[shapeIndex]
	shapeXML := content[targetShape.Start:targetShape.End]

	var buf bytes.Buffer

	if toFront {
		// Move to the end of the shape list
		lastShape := shapes[len(shapes)-1]

		// 1. Everything before target shape
		buf.Write(content[:targetShape.Start])
		// 2. Everything between target shape and after last shape (excluding target shape itself)
		buf.Write(content[targetShape.End:lastShape.End])
		// 3. The target shape
		buf.Write(shapeXML)
		// 4. Everything after last shape
		buf.Write(content[lastShape.End:])
	} else {
		// Move to the start of the shape list
		firstShape := shapes[0]

		// 1. Everything before the first shape
		buf.Write(content[:firstShape.Start])
		// 2. The target shape
		buf.Write(shapeXML)
		// 3. Everything between first shape and target shape (excluding target shape itself)
		buf.Write(content[firstShape.Start:targetShape.Start])
		// 4. Everything after the target shape
		buf.Write(content[targetShape.End:])
	}

	e.parts.Set(partPath, buf.Bytes())
	return nil
}
