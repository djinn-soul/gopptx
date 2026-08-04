package editor

import common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"

func commonShapeFromParsed(shape parsedShape) common.Shape {
	var placeholderIndex *int
	if shape.PhType != "" {
		idx := shape.PhIndex
		placeholderIndex = &idx
	}
	children := make([]common.Shape, len(shape.Shapes))
	for idx, child := range shape.Shapes {
		children[idx] = commonShapeFromParsed(child)
		// A child states its geometry in the group's child space. Reporting
		// those numbers as they stand puts the child in the wrong place on the
		// slide whenever chOff/chExt differ from the group's own off/ext
		// (upstream issue #925).
		mapSubtreeToSlideSpace(&children[idx], shape.GroupChild, shape.X, shape.Y, shape.W, shape.H)
	}
	return common.Shape{
		ID:               shape.ID,
		Name:             shape.Name,
		Type:             shape.Type,
		Text:             shape.Text,
		X:                shape.X,
		Y:                shape.Y,
		W:                shape.W,
		H:                shape.H,
		Runs:             shape.Runs,
		Paragraphs:       shape.Paragraphs,
		TextFrame:        shape.TextFrame,
		Paragraph:        shape.Paragraph,
		Rotation:         shape.Rotation,
		FlipH:            shape.FlipH,
		FlipV:            shape.FlipV,
		Hidden:           shape.Hidden,
		PlaceholderIndex: placeholderIndex,
		PlaceholderType:  shape.PhType,
		Fill:             shape.Fill,
		Line:             shape.Line,
		Shadow:           shape.Shadow,
		Glow:             shape.Glow,
		Blur:             shape.Blur,
		SoftEdge:         shape.SoftEdge,
		Reflection:       shape.Reflection,
		ClickAction:      shape.ClickAction,
		HoverAction:      shape.HoverAction,
		AltText:          shape.AltText,
		Title:            shape.Title,
		IsDecorative:     shape.IsDecorative,
		Connector:        shape.Connector,
		Adjustments:      shape.Adjustments,
		Freeform:         shape.Freeform,
		GroupChild:       shape.GroupChild,
		Shapes:           children,
	}
}

// mapSubtreeToSlideSpace applies one group's child-to-slide transform to a whole
// subtree, not just its immediate child.
//
// A nested group's descendants come back from the recursion already reduced to
// the outer group's child space, and the transform is a scale and a translate,
// so the same map applies to every one of them. Transforming only the immediate
// child left grandchildren behind in the intermediate space.
func mapSubtreeToSlideSpace(
	shape *common.Shape,
	space *common.GroupChildSpace,
	groupX, groupY, groupW, groupH int,
) {
	shape.X, shape.Y, shape.W, shape.H = space.ChildToSlide(
		groupX, groupY, groupW, groupH,
		shape.X, shape.Y, shape.W, shape.H,
	)
	for i := range shape.Shapes {
		mapSubtreeToSlideSpace(&shape.Shapes[i], space, groupX, groupY, groupW, groupH)
	}
}
