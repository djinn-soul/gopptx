package editor

import (
	"bytes"
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

const (
	shapeTypePicture = "pic"
	minGroupShapes   = 2
	groupShapeTag    = "grpSp"
)

// GetShapes returns a list of shapes found on the specified slide (0-based index).
func (e *PresentationEditor) GetShapes(slideIndex int) ([]common.Shape, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return nil, fmt.Errorf("read slide part %s: not found", partPath)
	}

	parsed, err := parseSlideShapes(content)
	if err != nil {
		return nil, fmt.Errorf("parse shapes: %w", err)
	}

	shapes := make([]common.Shape, len(parsed))
	for i, p := range parsed {
		shapes[i] = common.Shape{
			ID:          p.ID,
			Name:        p.Name,
			Type:        p.Type,
			Text:        p.Text,
			X:           p.X,
			Y:           p.Y,
			W:           p.W,
			H:           p.H,
			Fill:        p.Fill,
			Line:        p.Line,
			Shadow:      p.Shadow,
			Adjustments: p.Adjustments,
		}
	}
	return shapes, nil
}

// RemoveShapeByIndex removes a shape from the slide by its index.
func (e *PresentationEditor) RemoveShapeByIndex(slideIndex, shapeIndex int) error {
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

	if shapeIndex < 0 || shapeIndex >= len(shapes) {
		return errors.New("shape index out of range")
	}

	return e.applyShapeRemoval(partPath, content, shapes, shapeIndex)
}

// RemoveShape removes a shape from the slide by its ID.
func (e *PresentationEditor) RemoveShape(slideIndex, shapeID int) error {
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

	return e.applyShapeRemoval(partPath, content, shapes, shapeIndex)
}

// GroupShapes groups the specified shapes on a slide into a group shape.
// Returns the ID of the created group shape.
func (e *PresentationEditor) GroupShapes(slideIndex int, shapeIDs []int) (int, error) {
	if len(shapeIDs) < minGroupShapes {
		return 0, errors.New("at least 2 shapes are required to form a group")
	}
	return e.AddGroupShape(slideIndex, shapeIDs)
}

// UngroupShapes ungroups a group shape, returning the ID of the first member shape.
func (e *PresentationEditor) UngroupShapes(slideIndex, shapeID int) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return 0, errors.New("read slide part: not found")
	}

	shapes, err := parseSlideShapes(content)
	if err != nil {
		return 0, fmt.Errorf("parse shapes: %w", err)
	}

	groupShapeIndex := -1
	for i, s := range shapes {
		if s.ID == shapeID && s.IsGroup {
			groupShapeIndex = i
			break
		}
	}

	if groupShapeIndex == -1 {
		return 0, fmt.Errorf("group shape with ID %d not found", shapeID)
	}

	groupShape := shapes[groupShapeIndex]
	groupXML := content[groupShape.Start:groupShape.End]
	children, err := editorshape.ExtractGroupChildShapeNodes(groupXML, groupShapeTag, shapeTypePicture)
	if err != nil {
		return 0, fmt.Errorf("extract grouped shapes: %w", err)
	}
	if len(children) == 0 {
		return 0, fmt.Errorf("group shape with ID %d has no child shapes", shapeID)
	}

	replacement := bytes.Join(children, nil)
	newContent := replaceShapeNodes(content, shapes, func(i int, _ *parsedShape) ([]byte, bool) {
		if i == groupShapeIndex {
			return replacement, true
		}
		return nil, false
	})
	e.parts.Set(partPath, newContent)

	firstChildID := editorshape.FirstShapeIDInXML(children, cNvPrIDPattern, cNvPrSubmatchSize)
	if firstChildID == 0 {
		return shapeID, nil
	}
	return firstChildID, nil
}
