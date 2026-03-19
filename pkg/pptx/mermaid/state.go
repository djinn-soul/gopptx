package mermaid

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// StateNode represents a state in a state diagram.
type StateNode struct {
	ID    string
	Label string
	Type  string // "normal", "start", "end"
}

// StateTransition represents a transition between states.
type StateTransition struct {
	From  string
	To    string
	Label string
}

// StateDiagram represents the parsed structure of a Mermaid state diagram.
type StateDiagram struct {
	States      []StateNode
	Transitions []StateTransition
}

// renderState parses and renders a Mermaid state diagram into PowerPoint elements.
func renderState(code string, theme Theme) DiagramElements {
	diagram := parseState(code)
	return generateStateElements(diagram, theme)
}

func parseState(code string) *StateDiagram {
	lines := ParseLines(code)
	states := make(map[string]*StateNode)
	var transitions []StateTransition
	startMarkerCount := 0
	endMarkerCount := 0

	for _, line := range lines {
		if strings.HasPrefix(line, "stateDiagram") {
			continue
		}

		if from, to, label, found := splitStateTransition(line); found {
			from = resolveStateEndpoint(from, &startMarkerCount, true)
			to = resolveStateEndpoint(to, &endMarkerCount, false)
			transitions = append(transitions, StateTransition{
				From:  from,
				To:    to,
				Label: label,
			})

			ensureState(states, from)
			ensureState(states, to)
			continue
		}

		// Handle state definition: state "Label" as ID
		if strings.HasPrefix(line, "state ") {
			parts := strings.Fields(line)
			if len(parts) >= 4 && parts[2] == "as" {
				id := parts[3]
				label := strings.Trim(parts[1], "\"")
				states[id] = &StateNode{ID: id, Label: label, Type: "normal"}
			} else if len(parts) >= 2 {
				id := parts[1]
				states[id] = &StateNode{ID: id, Label: id, Type: "normal"}
			}
			continue
		}

		// Simple state label: ID : Label
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			id := strings.TrimSpace(parts[0])
			label := strings.TrimSpace(parts[1])
			if s, ok := states[id]; ok {
				s.Label = label
			} else {
				states[id] = &StateNode{ID: id, Label: label, Type: "normal"}
			}
		}
	}

	stateList := make([]StateNode, 0, len(states))
	for _, s := range states {
		stateList = append(stateList, *s)
	}

	return &StateDiagram{
		States:      stateList,
		Transitions: transitions,
	}
}

func ensureState(states map[string]*StateNode, id string) {
	if _, ok := states[id]; !ok {
		stateType := "normal"
		label := id
		if strings.HasPrefix(id, "__start_") {
			stateType = "start"
			label = ""
		}
		if strings.HasPrefix(id, "__end_") {
			stateType = "end"
			label = ""
		}
		states[id] = &StateNode{ID: id, Label: label, Type: stateType}
	}
}

func resolveStateEndpoint(id string, counter *int, isFrom bool) string {
	if id != "[*]" {
		return id
	}
	(*counter)++
	if isFrom {
		return fmt.Sprintf("__start_%d", *counter)
	}
	return fmt.Sprintf("__end_%d", *counter)
}

func splitStateTransition(line string) (string, string, string, bool) {
	if before, after, ok := strings.Cut(line, "-->"); ok {
		from := strings.TrimSpace(before)
		rest := strings.TrimSpace(after)
		to := rest
		label := ""
		if before, after, ok := strings.Cut(rest, ":"); ok {
			to = strings.TrimSpace(before)
			label = strings.TrimSpace(after)
		}
		return from, to, label, true
	}
	return "", "", "", false
}

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

	for i, state := range diagram.States {
		x, y := stateGridPosition(i, layout)
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

func stateGridPosition(index int, layout stateLayout) (styling.Length, styling.Length) {
	col := index % layout.cols
	row := index / layout.cols
	x := layout.startX + styling.Length(col)*layout.hSpacing
	y := layout.startY + styling.Length(row)*layout.vSpacing
	return x, y
}

func stateNodeShape(state StateNode, x styling.Length, y styling.Length, layout stateLayout, theme Theme) shapes.Shape {
	if state.Type == "start" || state.Type == "end" {
		circleSize := styling.Inches(0.36)
		lineColor := theme.PrimaryStroke
		lineWeight := theme.LineWeight
		fillColor := theme.PrimaryStroke
		if state.Type == "end" {
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

func stateTransitionEndpoints(
	fromPos struct{ x, y styling.Length },
	toPos struct{ x, y styling.Length },
	fromWidth styling.Length,
	toWidth styling.Length,
	fromHeight styling.Length,
	toHeight styling.Length,
) stateTransitionGeometry {
	switch {
	case fromPos.x < toPos.x:
		return stateTransitionGeometry{
			startX:    fromPos.x + fromWidth,
			startY:    fromPos.y + fromHeight/2,
			endX:      toPos.x,
			endY:      toPos.y + toHeight/2,
			startSite: shapes.ConnectionSiteRight,
			endSite:   shapes.ConnectionSiteLeft,
		}
	case fromPos.x > toPos.x:
		return stateTransitionGeometry{
			startX:    fromPos.x,
			startY:    fromPos.y + fromHeight/2,
			endX:      toPos.x + toWidth,
			endY:      toPos.y + toHeight/2,
			startSite: shapes.ConnectionSiteLeft,
			endSite:   shapes.ConnectionSiteRight,
		}
	case fromPos.y < toPos.y:
		return stateTransitionGeometry{
			startX:    fromPos.x + fromWidth/2,
			startY:    fromPos.y + fromHeight,
			endX:      toPos.x + toWidth/2,
			endY:      toPos.y,
			startSite: shapes.ConnectionSiteBottom,
			endSite:   shapes.ConnectionSiteTop,
		}
	default:
		return stateTransitionGeometry{
			startX:    fromPos.x + fromWidth/2,
			startY:    fromPos.y,
			endX:      toPos.x + toWidth/2,
			endY:      toPos.y + toHeight,
			startSite: shapes.ConnectionSiteTop,
			endSite:   shapes.ConnectionSiteBottom,
		}
	}
}

func stateTransitionLabelShape(
	label string,
	startX styling.Length,
	startY styling.Length,
	endX styling.Length,
	endY styling.Length,
	theme Theme,
) shapes.Shape {
	labelWidth := styling.Inches(1.0)
	labelHeight := styling.Inches(0.4)
	midX := (startX + endX) / 2
	midY := (startY + endY) / 2
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		midX-labelWidth/2,
		midY-labelHeight/2,
		labelWidth,
		labelHeight,
	).WithText(label).
		WithFill(shapes.NewShapeFill(theme.Background)).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.05), styling.Inches(0.02), styling.Inches(0.05), styling.Inches(0.02))
}
