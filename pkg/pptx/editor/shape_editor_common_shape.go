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
		Shapes:           children,
	}
}
