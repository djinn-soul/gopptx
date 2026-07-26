package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func sequenceArrowGeometry(
	fromX styling.Length,
	toX styling.Length,
	participantWidth styling.Length,
) (styling.Length, styling.Length, string) {
	fromCenter := fromX + participantWidth/2
	toCenter := toX + participantWidth/2
	if fromCenter < toCenter {
		return fromCenter, toCenter - fromCenter, shapes.ShapeTypeRightArrow
	}
	return toCenter, fromCenter - toCenter, shapes.ShapeTypeLeftArrow
}
