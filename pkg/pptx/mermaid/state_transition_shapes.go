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

// stateNoteConnector is the dashed, arrowless leader from a note to the state
// it annotates.
func stateNoteConnector(
	note StateNode,
	statePositions map[string]struct{ x, y styling.Length },
	stateSizes map[string]struct{ w, h styling.Length },
	stateShapeIndices map[string]int,
	theme Theme,
) (shapes.Connector, bool) {
	notePos, noteExists := statePositions[note.ID]
	targetPos, targetExists := statePositions[note.NoteTarget]
	noteSize, noteSizeOK := stateSizes[note.ID]
	targetSize, targetSizeOK := stateSizes[note.NoteTarget]
	if !noteExists || !targetExists || !noteSizeOK || !targetSizeOK {
		return shapes.Connector{}, false
	}

	geom := stateTransitionEndpoints(notePos, targetPos, noteSize.w, targetSize.w, noteSize.h, targetSize.h)
	connector := shapes.NewConnector(shapes.ConnectorTypeStraight, geom.startX, geom.startY, geom.endX, geom.endY).
		WithLine(shapes.NewShapeLine(stateNoteStroke, theme.LineWeight)).
		WithDash(shapes.LineDashDash).
		WithArrows(shapes.ArrowTypeNone, shapes.ArrowTypeNone)

	if idx, ok := stateShapeIndices[note.ID]; ok {
		connector = connector.ConnectStart(idx, geom.startSite)
	}
	if idx, ok := stateShapeIndices[note.NoteTarget]; ok {
		connector = connector.ConnectEnd(idx, geom.endSite)
	}
	return connector, true
}
