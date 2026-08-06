package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// This file lays a parsed state diagram out and turns it into shapes; the
// parsing that produces the diagram lives in state.go.

func generateStateElements(diagram *StateDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape
	var connectors []shapes.Connector

	if len(diagram.States) == 0 {
		return DiagramElements{Grouped: true}
	}

	layout := stateLayout{
		stateWidth:  styling.Inches(1.8),
		stateHeight: styling.Inches(0.8),
		hSpacing:    styling.Inches(2.8),
		vSpacing:    styling.Inches(1.8),
		startX:      styling.Inches(1.0),
		startY:      styling.Inches(1.0),
		cols:        3,
	}

	statePositions := make(map[string]struct{ x, y styling.Length })
	stateSizes := make(map[string]struct{ w, h styling.Length })
	stateShapeIndices := make(map[string]int)
	bounds := newStateBounds()

	placements := planStateLayout(diagram.States, layout)
	childrenByParent := groupStateChildren(diagram.States)

	// Containers paint first so their children sit on top of them.
	for _, state := range diagram.States {
		if !state.IsComposite {
			continue
		}
		container := stateContainerShape(state, childrenByParent[state.ID], placements, layout, theme)
		statePositions[state.ID] = struct{ x, y styling.Length }{container.X, container.Y}
		stateSizes[state.ID] = struct{ w, h styling.Length }{container.CX, container.CY}
		shapesList = append(shapesList, container)
		bounds.includeShape(container)
		stateShapeIndices[state.ID] = len(shapesList)
	}

	for _, state := range diagram.States {
		if state.IsComposite {
			continue
		}
		spot := placements[state.ID]
		x, y := spot.x, spot.y
		shape := stateNodeShape(state, x, y, layout, theme)
		statePositions[state.ID] = struct{ x, y styling.Length }{shape.X, shape.Y}
		stateSizes[state.ID] = struct{ w, h styling.Length }{shape.CX, shape.CY}
		shapesList = append(shapesList, shape)
		bounds.includeShape(shape)
		stateShapeIndices[state.ID] = len(shapesList)
	}

	for _, trans := range diagram.Transitions {
		connector, label, ok := stateTransitionShapes(trans, statePositions, stateSizes, stateShapeIndices, theme)
		if !ok {
			continue
		}
		connectors = append(connectors, connector)
		if label != nil {
			shapesList = append(shapesList, *label)
		}
	}

	// A note is tied to its state with a dashed leader, as Mermaid draws it.
	for _, state := range diagram.States {
		if state.Type != stateTypeNote || state.NoteTarget == "" {
			continue
		}
		leader, ok := stateNoteConnector(state, statePositions, stateSizes, stateShapeIndices, theme)
		if ok {
			connectors = append(connectors, leader)
		}
	}

	return DiagramElements{
		Shapes:     shapesList,
		Connectors: connectors,
		Grouped:    true,
		Bounds: &DiagramBounds{
			X:  bounds.minX,
			Y:  bounds.minY,
			CX: bounds.maxX - bounds.minX,
			CY: bounds.maxY - bounds.minY,
		},
	}
}

type stateLayout struct {
	stateWidth  styling.Length
	stateHeight styling.Length
	hSpacing    styling.Length
	vSpacing    styling.Length
	startX      styling.Length
	startY      styling.Length
	cols        int
}

type stateBounds struct {
	minX  styling.Length
	minY  styling.Length
	maxX  styling.Length
	maxY  styling.Length
	first bool
}

func newStateBounds() *stateBounds {
	return &stateBounds{first: true}
}

func (b *stateBounds) includeShape(s shapes.Shape) {
	b.include(s.X, s.Y, s.CX, s.CY)
}

func (b *stateBounds) include(x, y, cx, cy styling.Length) {
	if b.first {
		b.minX, b.minY = x, y
		b.maxX, b.maxY = x+cx, y+cy
		b.first = false
		return
	}
	if x < b.minX {
		b.minX = x
	}
	if y < b.minY {
		b.minY = y
	}
	if x+cx > b.maxX {
		b.maxX = x + cx
	}
	if y+cy > b.maxY {
		b.maxY = y + cy
	}
}

func stateNodeShape(state StateNode, x styling.Length, y styling.Length, layout stateLayout, theme Theme) shapes.Shape {
	switch state.Type {
	case stateTypeFork:
		// A fork or join is a heavy bar, not a box.
		barHeight := styling.Inches(0.12)
		return shapes.NewShape(
			shapes.ShapeTypeRectangle,
			x,
			y+(layout.stateHeight-barHeight)/2,
			layout.stateWidth,
			barHeight,
		).WithFill(shapes.NewShapeFill(theme.PrimaryStroke)).
			WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
			WithAltText(state.Label)
	case stateTypeChoice:
		size := layout.stateHeight
		return shapes.NewShape(
			shapes.ShapeTypeDiamond,
			x+(layout.stateWidth-size)/2,
			y,
			size,
			size,
		).WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
			WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
			WithAltText(state.Label)
	case stateTypeNote:
		return shapes.NewShape(
			shapes.ShapeTypeRectangle,
			x,
			y,
			layout.stateWidth,
			layout.stateHeight,
		).WithFill(shapes.NewShapeFill(stateNoteFill)).
			WithLine(shapes.NewShapeLine(stateNoteStroke, theme.LineWeight)).
			WithText(state.Label).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
	}

	if state.Type == "start" || state.Type == stateTypeEnd {
		circleSize := styling.Inches(0.36)
		lineColor := theme.PrimaryStroke
		lineWeight := theme.LineWeight
		fillColor := theme.PrimaryStroke
		if state.Type == stateTypeEnd {
			fillColor = theme.Background
			lineWeight = theme.LineWeight * 2
		}
		return shapes.NewShape(
			shapes.ShapeTypeEllipse,
			x+(layout.stateWidth-circleSize)/2,
			y+(layout.stateHeight-circleSize)/2,
			circleSize,
			circleSize,
		).WithFill(shapes.NewShapeFill(fillColor)).
			WithLine(shapes.NewShapeLine(lineColor, lineWeight))
	}
	return shapes.NewShape(
		shapes.ShapeTypeRoundedRectangle,
		x,
		y,
		layout.stateWidth,
		layout.stateHeight,
	).WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
		WithText(state.Label).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
}

type stateTransitionGeometry struct {
	startX    styling.Length
	startY    styling.Length
	endX      styling.Length
	endY      styling.Length
	startSite string
	endSite   string
}
