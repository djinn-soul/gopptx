package editor

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	editorgrayscale "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/grayscale"
)

func (e *PresentationEditor) addGrayscalePlaceholderTargets(
	targets *grayscaleTargets,
	refs []editorgrayscale.PlaceholderRef,
) error {
	for _, ref := range refs {
		if err := e.validateSlideIndex(ref.SlideIndex); err != nil {
			return err
		}
		if ref.Type == "" && ref.Index == nil {
			return fmt.Errorf(
				"placeholder target on slide %d must include at least one of type or index",
				ref.SlideIndex,
			)
		}
		matches, err := e.resolveGrayscalePlaceholderShapeIDs(ref)
		if err != nil {
			return err
		}
		for _, shapeID := range matches {
			targets.shapes[shapeTargetKey(ref.SlideIndex, shapeID)] = struct{}{}
		}
		targets.slides[ref.SlideIndex] = struct{}{}
	}
	return nil
}

func (e *PresentationEditor) resolveGrayscalePlaceholderShapeIDs(
	ref editorgrayscale.PlaceholderRef,
) ([]int, error) {
	shapes, err := e.getShapesForTextOps(ref.SlideIndex)
	if err != nil {
		return nil, err
	}

	shapeIDs := make([]int, 0, len(shapes))
	for _, shape := range shapes {
		if !matchesGrayscalePlaceholderRef(shape, ref) {
			continue
		}
		shapeIDs = append(shapeIDs, shape.ID)
	}
	if len(shapeIDs) == 0 {
		return nil, fmt.Errorf(
			"placeholder target not found on slide %d (type=%q index=%s)",
			ref.SlideIndex,
			string(ref.Type),
			formatGrayscalePlaceholderIndex(ref.Index),
		)
	}
	return shapeIDs, nil
}

func matchesGrayscalePlaceholderRef(shape parsedShape, ref editorgrayscale.PlaceholderRef) bool {
	if shape.PhType == "" {
		return false
	}
	if ref.Index != nil && shape.PhIndex != *ref.Index {
		return false
	}
	if ref.Type == "" {
		return true
	}
	requested := pptxxml.NormalizePlaceholderType(string(ref.Type))
	actual := pptxxml.NormalizePlaceholderType(shape.PhType)
	return requested == actual || (requested == placeholderTypeTitle && actual == "ctrTitle")
}

func formatGrayscalePlaceholderIndex(index *int) string {
	if index == nil {
		return "*"
	}
	return fmt.Sprintf("%d", *index)
}
