package mermaid

import (
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

	for _, line := range lines {
		if strings.HasPrefix(line, "stateDiagram") {
			continue
		}

		if from, to, label, found := splitStateTransition(line); found {
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
		if id == "[*]" {
			stateType = "start" // Can be start or end depending on context, but we'll handle it in rendering
			label = ""
		}
		states[id] = &StateNode{ID: id, Label: label, Type: stateType}
	}
}

func splitStateTransition(line string) (string, string, string, bool) {
	if idx := strings.Index(line, "-->"); idx != -1 {
		from := strings.TrimSpace(line[:idx])
		rest := strings.TrimSpace(line[idx+3:])
		to := rest
		label := ""
		if labelIdx := strings.Index(rest, ":"); labelIdx != -1 {
			to = strings.TrimSpace(rest[:labelIdx])
			label = strings.TrimSpace(rest[labelIdx+1:])
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

	// Layout parameters
	stateWidth := styling.Inches(1.8)
	stateHeight := styling.Inches(0.8)
	hSpacing := styling.Inches(2.8)
	vSpacing := styling.Inches(1.8)

	statePositions := make(map[string]struct{ x, y styling.Length })
	stateShapeIndices := make(map[string]int)

	var minX, minY, maxX, maxY styling.Length
	firstElement := true

	updateBounds := func(x, y, cx, cy styling.Length) {
		if firstElement {
			minX, minY = x, y
			maxX, maxY = x+cx, y+cy
			firstElement = false
		} else {
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}
			if x+cx > maxX {
				maxX = x + cx
			}
			if y+cy > maxY {
				maxY = y + cy
			}
		}
	}

	// Simple grid layout
	startX := styling.Inches(1.0)
	startY := styling.Inches(1.0)
	cols := 3

	for i, state := range diagram.States {
		col := i % cols
		row := i / cols

		x := startX + styling.Length(col)*hSpacing
		y := startY + styling.Length(row)*vSpacing

		statePositions[state.ID] = struct{ x, y styling.Length }{x, y}

		var shape shapes.Shape
		if state.ID == "[*]" {
			// Start/End state is a small circle
			circleSize := styling.Inches(0.3)
			shape = shapes.NewShape(shapes.ShapeTypeEllipse, x+(stateWidth-circleSize)/2, y+(stateHeight-circleSize)/2, circleSize, circleSize).
				WithFill(shapes.NewShapeFill(theme.PrimaryStroke))
			updateBounds(x+(stateWidth-circleSize)/2, y+(stateHeight-circleSize)/2, circleSize, circleSize)
		} else {
			shape = shapes.NewShape(shapes.ShapeTypeRoundedRectangle, x, y, stateWidth, stateHeight).
				WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
				WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
				WithText(state.Label).
				WithAutoFit(shapes.TextAutoFitNormal).
				WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
			updateBounds(x, y, stateWidth, stateHeight)
		}

		shapesList = append(shapesList, shape)
		stateShapeIndices[state.ID] = len(shapesList)
	}

	// Create connectors
	for _, trans := range diagram.Transitions {
		fromPos, fromExists := statePositions[trans.From]
		toPos, toExists := statePositions[trans.To]

		if fromExists && toExists {
			var startX, startY, endX, endY styling.Length
			var startSite, endSite string

			// Simple logic to decide connection points
			if fromPos.x < toPos.x {
				startX = fromPos.x + stateWidth
				startY = fromPos.y + stateHeight/2
				startSite = shapes.ConnectionSiteRight
				endSite = shapes.ConnectionSiteLeft
				endX = toPos.x
				endY = toPos.y + stateHeight/2
			} else if fromPos.x > toPos.x {
				startX = fromPos.x
				startY = fromPos.y + stateHeight/2
				startSite = shapes.ConnectionSiteLeft
				endSite = shapes.ConnectionSiteRight
				endX = toPos.x + stateWidth
				endY = toPos.y + stateHeight/2
			} else if fromPos.y < toPos.y {
				startX = fromPos.x + stateWidth/2
				startY = fromPos.y + stateHeight
				startSite = shapes.ConnectionSiteBottom
				endSite = shapes.ConnectionSiteTop
				endX = toPos.x + stateWidth/2
				endY = toPos.y
			} else {
				startX = fromPos.x + stateWidth/2
				startY = fromPos.y
				startSite = shapes.ConnectionSiteTop
				endSite = shapes.ConnectionSiteBottom
				endX = toPos.x + stateWidth/2
				endY = toPos.y + stateHeight
			}

			connector := shapes.NewConnector(shapes.ConnectorTypeElbow, startX, startY, endX, endY).
				WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
				WithArrows(shapes.ArrowTypeNone, shapes.ArrowTypeTriangle)

			if idx, ok := stateShapeIndices[trans.From]; ok {
				connector = connector.ConnectStart(idx, startSite)
			}
			if idx, ok := stateShapeIndices[trans.To]; ok {
				connector = connector.ConnectEnd(idx, endSite)
			}

			connectors = append(connectors, connector)

			if trans.Label != "" {
				labelWidth := styling.Inches(1.0)
				labelHeight := styling.Inches(0.4)
				midX := (startX + endX) / 2
				midY := (startY + endY) / 2

				labelShape := shapes.NewShape(shapes.ShapeTypeRectangle, midX-labelWidth/2, midY-labelHeight/2, labelWidth, labelHeight).
					WithFill(shapes.NewShapeFill(theme.Background)).
					WithText(trans.Label).
					WithAutoFit(shapes.TextAutoFitNormal).
					WithTextMargins(styling.Inches(0.05), styling.Inches(0.02), styling.Inches(0.05), styling.Inches(0.02))
				shapesList = append(shapesList, labelShape)
			}
		}
	}

	return DiagramElements{
		Shapes:     shapesList,
		Connectors: connectors,
		Grouped:    true,
		Bounds: &DiagramBounds{
			X:  minX,
			Y:  minY,
			CX: maxX - minX,
			CY: maxY - minY,
		},
	}
}
