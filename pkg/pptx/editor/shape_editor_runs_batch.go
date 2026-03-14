package editor

import (
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

// SetSlideShapeRuns replaces runs on multiple slide shapes in one parse/rewrite pass.
func (e *PresentationEditor) SetSlideShapeRuns(slideIndex int, updates []common.ShapeRunsUpdate) error {
	if len(updates) == 0 {
		return nil
	}
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

	updatesByShape := make(map[int]common.ShapeRunsUpdate, len(updates))
	for _, update := range updates {
		updatesByShape[update.ShapeID] = update
	}
	found := make(map[int]bool, len(updatesByShape))
	var replaceErr error

	newXML := replaceShapeNodes(content, shapes, func(_ int, shape *parsedShape) ([]byte, bool) {
		targetID, ok := matchShapeRunsTarget(shape, updatesByShape)
		if !ok || replaceErr != nil {
			return nil, false
		}
		update := updatesByShape[targetID]
		shape.Runs = editorshape.CopyTextRuns(update.Runs)
		updatedXML, err := replaceShapeTextBody(e, partPath, content[shape.Start:shape.End], shape)
		if err != nil {
			replaceErr = err
			return nil, false
		}
		found[targetID] = true
		return updatedXML, true
	})
	if replaceErr != nil {
		return replaceErr
	}
	for shapeID := range updatesByShape {
		if !found[shapeID] {
			return fmt.Errorf("shape id %d not found on slide %d", shapeID, slideIndex)
		}
	}

	e.parts.Set(partPath, newXML)
	return nil
}

func matchShapeRunsTarget(shape *parsedShape, updatesByShape map[int]common.ShapeRunsUpdate) (int, bool) {
	if _, ok := updatesByShape[shape.ID]; ok {
		return shape.ID, true
	}
	if shape.PhType == placeholderTypeTitle {
		if _, ok := updatesByShape[0]; ok {
			return 0, true
		}
	}
	return 0, false
}
