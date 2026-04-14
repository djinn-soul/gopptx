package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func stateTransitionShapes(
	trans StateTransition,
	statePositions map[string]struct{ x, y styling.Length },
	stateSizes map[string]struct{ w, h styling.Length },
	stateShapeIndices map[string]int,
	theme Theme,
) (shapes.Connector, *shapes.Shape, bool) {
	fromPos, fromExists := statePositions[trans.From]
	toPos, toExists := statePositions[trans.To]
	fromSize, fromSizeOK := stateSizes[trans.From]
	toSize, toSizeOK := stateSizes[trans.To]
	if !fromExists || !toExists || !fromSizeOK || !toSizeOK {
		return shapes.Connector{}, nil, false
	}

	geom := stateTransitionEndpoints(fromPos, toPos, fromSize.w, toSize.w, fromSize.h, toSize.h)
	connector := shapes.NewConnector(shapes.ConnectorTypeElbow, geom.startX, geom.startY, geom.endX, geom.endY).
		WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
		WithArrows(shapes.ArrowTypeNone, shapes.ArrowTypeTriangle)

	if idx, ok := stateShapeIndices[trans.From]; ok {
		connector = connector.ConnectStart(idx, geom.startSite)
	}
	if idx, ok := stateShapeIndices[trans.To]; ok {
		connector = connector.ConnectEnd(idx, geom.endSite)
	}

	if trans.Label == "" {
		return connector, nil, true
	}
	label := stateTransitionLabelShape(trans.Label, geom.startX, geom.startY, geom.endX, geom.endY, theme)
	return connector, &label, true
}
