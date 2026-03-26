package editor

import (
	"bytes"
	"errors"
	"fmt"
)

// BringShapeToFront moves the specified shape to the top of the z-order
// (rendered last, appearing above all other shapes).
func (e *PresentationEditor) BringShapeToFront(slideIndex, shapeID int) error {
	return e.moveShape(slideIndex, shapeID, func(_ int, count int) int {
		return count - 1
	})
}

// SendShapeToBack moves the specified shape to the bottom of the z-order
// (rendered first, appearing behind all other shapes).
func (e *PresentationEditor) SendShapeToBack(slideIndex, shapeID int) error {
	return e.moveShape(slideIndex, shapeID, func(_, _ int) int {
		return 0
	})
}

// BringShapeForward moves the specified shape one position forward in the z-order.
func (e *PresentationEditor) BringShapeForward(slideIndex, shapeID int) error {
	return e.moveShape(slideIndex, shapeID, func(idx, count int) int {
		if idx < count-1 {
			return idx + 1
		}
		return idx
	})
}

// SendShapeBackward moves the specified shape one position backward in the z-order.
func (e *PresentationEditor) SendShapeBackward(slideIndex, shapeID int) error {
	return e.moveShape(slideIndex, shapeID, func(idx, _ int) int {
		if idx > 0 {
			return idx - 1
		}
		return idx
	})
}

// ReorderShape moves a shape to the specified z-order position.
// direction must be one of: "front", "back", "forward", "backward".
func (e *PresentationEditor) ReorderShape(slideIndex, shapeID int, direction string) error {
	switch direction {
	case "front":
		return e.BringShapeToFront(slideIndex, shapeID)
	case "back":
		return e.SendShapeToBack(slideIndex, shapeID)
	case "forward":
		return e.BringShapeForward(slideIndex, shapeID)
	case "backward":
		return e.SendShapeBackward(slideIndex, shapeID)
	default:
		return fmt.Errorf("invalid direction %q: must be front, back, forward, or backward", direction)
	}
}

// moveShape is the shared implementation for all z-order operations.
// targetFn receives the current index and total shape count, and returns the desired target index.
func (e *PresentationEditor) moveShape(slideIndex, shapeID int, targetFn func(idx, count int) int) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return fmt.Errorf("read slide part %s: not found", partPath)
	}

	shapes, err := parseSlideShapes(content)
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

	target := targetFn(shapeIndex, len(shapes))
	if target == shapeIndex {
		return nil // already in position
	}

	newOrder := buildReorderedIndices(len(shapes), shapeIndex, target)
	newContent := rebuildWithReorderedShapes(content, shapes, newOrder)
	e.parts.Set(partPath, newContent)
	return nil
}

// buildReorderedIndices returns a new index order with the element at fromIdx moved to toIdx.
func buildReorderedIndices(count, fromIdx, toIdx int) []int {
	order := make([]int, 0, count)
	for i := range count {
		if i != fromIdx {
			order = append(order, i)
		}
	}
	// Insert fromIdx at toIdx position
	result := make([]int, count)
	copy(result[:toIdx], order[:toIdx])
	result[toIdx] = fromIdx
	copy(result[toIdx+1:], order[toIdx:])
	return result
}

// rebuildWithReorderedShapes reconstructs the slide XML with shapes in a new order.
// Content between shape nodes (typically whitespace) is discarded; shape XML itself is preserved.
func rebuildWithReorderedShapes(content []byte, shapes []parsedShape, newOrder []int) []byte {
	if len(shapes) == 0 || len(newOrder) == 0 {
		return content
	}

	firstStart := shapes[0].Start
	lastEnd := shapes[len(shapes)-1].End

	var buf bytes.Buffer
	buf.Write(content[:firstStart])
	for _, idx := range newOrder {
		buf.Write(content[shapes[idx].Start:shapes[idx].End])
	}
	buf.Write(content[lastEnd:])
	return buf.Bytes()
}
