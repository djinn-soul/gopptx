package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Padding around the states nested inside a composite.
const (
	statePadX      styling.Length = 182880 // 0.2"
	stateTitleBand styling.Length = 320040 // 0.35" for the composite's name
	statePadBottom styling.Length = 182880 // 0.2"
)

// statePlacement is where one state's box goes.
type statePlacement struct {
	x styling.Length
	y styling.Length
}

// planStateLayout assigns every state a position, laying the children of a
// composite out as a block so a container can be drawn around them.
//
// Composite bodies used to be flattened into the same grid as everything else,
// which lost the nesting entirely: the children of "state First { … }" were
// scattered among the top-level states.
func planStateLayout(states []StateNode, layout stateLayout) map[string]statePlacement {
	placements := make(map[string]statePlacement, len(states))
	children := groupStateChildren(states)

	row := 0
	column := 0
	advance := func() {
		column++
		if column >= layout.cols {
			column = 0
			row++
		}
	}

	for _, state := range states {
		if state.Parent != "" {
			continue // placed with its parent
		}
		if !state.IsComposite {
			placements[state.ID] = statePlacementAt(row, column, layout)
			advance()
			continue
		}

		// A composite starts on its own row and takes as many rows as its
		// children need, so the container never overlaps a sibling.
		if column != 0 {
			column = 0
			row++
		}
		rowsUsed := placeCompositeChildren(state.ID, children, placements, layout, row)
		row += rowsUsed
	}
	return placements
}

// groupStateChildren indexes states by the composite that holds them.
func groupStateChildren(states []StateNode) map[string][]StateNode {
	children := make(map[string][]StateNode)
	for _, state := range states {
		if state.Parent == "" {
			continue
		}
		children[state.Parent] = append(children[state.Parent], state)
	}
	return children
}

// placeCompositeChildren lays a composite's children out inside its band and
// returns how many grid rows the composite consumed.
func placeCompositeChildren(
	id string,
	children map[string][]StateNode,
	placements map[string]statePlacement,
	layout stateLayout,
	startRow int,
) int {
	nested := children[id]
	if len(nested) == 0 {
		// An empty composite still occupies a row so its container has a size.
		placements[id] = statePlacementAt(startRow, 0, layout)
		return 1
	}

	innerCols := min(len(nested), layout.cols)
	rows := (len(nested) + innerCols - 1) / innerCols
	for i, child := range nested {
		spot := statePlacementAt(startRow+i/innerCols, i%innerCols, layout)
		// Nudge children down to clear the container's title band.
		spot.y += stateTitleBand
		placements[child.ID] = spot
	}
	return rows
}

func statePlacementAt(row, column int, layout stateLayout) statePlacement {
	return statePlacement{
		x: layout.startX + styling.Length(column)*layout.hSpacing,
		y: layout.startY + styling.Length(row)*layout.vSpacing,
	}
}

// stateContainerShape is the box drawn around a composite's children, labelled
// with the composite's name.
func stateContainerShape(
	state StateNode,
	children []StateNode,
	placements map[string]statePlacement,
	layout stateLayout,
	theme Theme,
) shapes.Shape {
	bounds := compositeChildBounds(state, children, placements, layout)
	return shapes.NewShape(
		shapes.ShapeTypeRoundedRectangle,
		bounds.minX,
		bounds.minY,
		bounds.maxX-bounds.minX,
		bounds.maxY-bounds.minY,
	).WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
		WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
		WithText(state.Label).
		WithVerticalAnchor(shapes.TextAnchorTop).
		WithTextMargins(styling.Inches(0.12), styling.Inches(0.06), styling.Inches(0.12), styling.Inches(0.06))
}

// compositeChildBounds is the extent of a composite's children, padded out to
// leave room for the title band above them.
func compositeChildBounds(
	state StateNode,
	children []StateNode,
	placements map[string]statePlacement,
	layout stateLayout,
) stateBounds {
	bounds := newStateBounds()
	for _, child := range children {
		spot, ok := placements[child.ID]
		if !ok {
			continue
		}
		bounds.include(spot.x, spot.y, layout.stateWidth, layout.stateHeight)
	}
	if bounds.first {
		// An empty body still needs a box: fall back to the composite's own slot.
		own := placements[state.ID]
		bounds.include(own.x, own.y+stateTitleBand, layout.stateWidth, layout.stateHeight)
	}
	bounds.minX -= statePadX
	bounds.minY -= stateTitleBand
	bounds.maxX += statePadX
	bounds.maxY += statePadBottom
	return *bounds
}
