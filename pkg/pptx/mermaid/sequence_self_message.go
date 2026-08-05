package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Geometry of the loop drawn for a message a participant sends to itself.
// Mermaid draws it as a bracket leaving the lifeline, dropping, and returning.
// EMU rather than styling.Inches so these stay constants: 914400 EMU per inch.
const (
	sequenceArrowThickness styling.Length = 137160 // 0.15"
	selfMessageLoopWidth   styling.Length = 640080 // 0.7"
	selfMessageLoopDrop    styling.Length = 320040 // 0.35"
	selfMessageStroke      styling.Length = 20000  // the lifeline's own width
	selfMessageLabelGap    styling.Length = 91440  // 0.1"
	selfMessageLabelWidth  styling.Length = 1463040
	selfMessageLabelHalf   styling.Length = 91440 // centres the 0.2" label
)

// sequenceSelfMessageShapes draws `A ->> A` as a loop hanging off A's lifeline.
//
// The three bars form the bracket: out from the lifeline, down, then back with
// the arrowhead. Drawing it as one arrow is impossible — source and target sit
// at the same x, so the arrow would have zero width.
func sequenceSelfMessageShapes(
	msg Message,
	participantLeftX styling.Length,
	y styling.Length,
	layout sequenceLayout,
	theme Theme,
) sequenceRenderedMessage {
	lifelineX := participantLeftX + layout.participantWidth/2
	stroke := shapes.NewShapeFill(theme.PrimaryStroke)

	outbound := shapes.NewShape(
		shapes.ShapeTypeRectangle, lifelineX, y, selfMessageLoopWidth, selfMessageStroke,
	).WithFill(stroke)

	descent := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		lifelineX+selfMessageLoopWidth-selfMessageStroke,
		y,
		selfMessageStroke,
		selfMessageLoopDrop,
	).WithFill(stroke)

	// The return leg carries the arrowhead, pointing back at the lifeline.
	returnLeg := shapes.NewShape(
		shapes.ShapeTypeLeftArrow,
		lifelineX,
		y+selfMessageLoopDrop-sequenceArrowThickness/2,
		selfMessageLoopWidth,
		sequenceArrowThickness,
	).WithFill(stroke)

	label := sequenceMessageLabel(
		msg.Text,
		lifelineX+selfMessageLoopWidth+selfMessageLabelGap,
		y+selfMessageLoopDrop/2-selfMessageLabelHalf,
		selfMessageLabelWidth,
	)

	return sequenceRenderedMessage{
		shapes: []shapes.Shape{outbound, descent, returnLeg, label},
		height: selfMessageLoopDrop + layout.messageSpacing/2,
	}
}
