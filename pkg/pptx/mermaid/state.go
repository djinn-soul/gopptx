package mermaid

import (
	"fmt"
	"sort"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Notes keep Mermaid's own pale yellow in every theme, the way a sticky note
// reads as an aside rather than another state.
const (
	stateNoteFill   = "FFF5AD"
	stateNoteStroke = "D6B656"
)

const (
	stateTypeEnd    = "end"
	stateTypeFork   = "fork"
	stateTypeChoice = "choice"
	stateTypeNote   = "note"
)

// StateNode represents a state in a state diagram.
type StateNode struct {
	ID    string
	Label string
	Type  string // "normal", "start", "end", "fork", "choice", "note"
	// NoteTarget is the state a note is attached to, empty for other types.
	NoteTarget string
	// Parent is the composite state this one is nested in, empty at top level.
	Parent string
	// IsComposite marks a state declared with a "state X { … }" body. It is
	// drawn as a container around its children rather than as a plain box.
	IsComposite bool
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

// stateParseState is what one line of a state diagram needs from the lines
// before it: the states and transitions so far, the counters that name the
// anonymous [*] markers, and the stack of open "state X {" bodies.
type stateParseState struct {
	states           map[string]*StateNode
	transitions      []StateTransition
	startMarkerCount int
	endMarkerCount   int
	composites       []string
}

func parseState(code string) *StateDiagram {
	lines := ParseLines(code)
	parser := &stateParseState{states: make(map[string]*StateNode)}

	for _, line := range lines {
		if strings.HasPrefix(line, "stateDiagram") {
			continue
		}
		parser.consumeLine(line)
	}

	states := parser.states
	transitions := parser.transitions

	stateList := make([]StateNode, 0, len(states))
	for _, s := range states {
		stateList = append(stateList, *s)
	}
	sortStatesForLayout(stateList)

	return &StateDiagram{
		States:      stateList,
		Transitions: transitions,
	}
}

// consumeLine reads one statement into the diagram being built.
func (p *stateParseState) consumeLine(line string) {
	if line == "}" {
		if len(p.composites) > 0 {
			p.composites = p.composites[:len(p.composites)-1]
		}
		return
	}
	if id, ok := parseCompositeStateOpen(line); ok {
		p.states[id] = &StateNode{
			ID:          id,
			Label:       id,
			Type:        stateTypeNormal,
			IsComposite: true,
			Parent:      currentComposite(p.composites),
		}
		p.composites = append(p.composites, id)
		return
	}
	if p.consumeTransition(line) {
		return
	}
	// A note is attached to a state rather than being one: without this it
	// fell through to the "ID : Label" branch and became a state of its own
	// named after the whole "note right of X" phrase.
	if note, ok := parseStateNote(line, len(p.states)); ok {
		p.states[note.ID] = note
		return
	}
	if p.consumeStateDeclaration(line) {
		return
	}
	p.consumeStateLabel(line)
}

func (p *stateParseState) consumeTransition(line string) bool {
	from, to, label, found := splitStateTransition(line)
	if !found {
		return false
	}
	from = resolveStateEndpoint(from, &p.startMarkerCount, true)
	to = resolveStateEndpoint(to, &p.endMarkerCount, false)
	p.transitions = append(p.transitions, StateTransition{From: from, To: to, Label: label})
	ensureState(p.states, from, currentComposite(p.composites))
	ensureState(p.states, to, currentComposite(p.composites))
	return true
}

// consumeStateDeclaration handles `state X`, `state "Label" as X` and the
// `<<fork>>` style markers.
func (p *stateParseState) consumeStateDeclaration(line string) bool {
	if !strings.HasPrefix(line, "state ") {
		return false
	}
	parts := strings.Fields(line)
	switch {
	case len(parts) >= 4 && parts[2] == "as":
		id := parts[3]
		p.states[id] = &StateNode{
			ID:     id,
			Label:  strings.Trim(parts[1], `"`),
			Type:   stateTypeNormal,
			Parent: currentComposite(p.composites),
		}
	case len(parts) >= 2:
		id := parts[1]
		p.states[id] = &StateNode{
			ID: id, Label: id, Type: stateTypeMarker(line), Parent: currentComposite(p.composites),
		}
	}
	return true
}

// consumeStateLabel handles the `ID : Label` form.
func (p *stateParseState) consumeStateLabel(line string) {
	id, label, found := strings.Cut(line, ":")
	if !found {
		return
	}
	id = strings.TrimSpace(id)
	label = strings.TrimSpace(label)
	if existing, ok := p.states[id]; ok {
		existing.Label = label
		return
	}
	p.states[id] = &StateNode{
		ID: id, Label: label, Type: stateTypeNormal, Parent: currentComposite(p.composites),
	}
}

// parseCompositeStateOpen matches "state X {", the header of a nested state.
func parseCompositeStateOpen(line string) (string, bool) {
	if !strings.HasPrefix(line, "state ") || !strings.HasSuffix(line, "{") {
		return "", false
	}
	body := strings.TrimSpace(strings.TrimSuffix(line[len("state "):], "{"))
	fields := strings.Fields(body)
	if len(fields) == 0 {
		return "", false
	}
	// "state "Label" as ID {" names the state by its id, which comes last.
	if len(fields) >= 3 && fields[len(fields)-2] == "as" {
		return fields[len(fields)-1], true
	}
	return fields[0], true
}

func currentComposite(composites []string) string {
	if len(composites) == 0 {
		return ""
	}
	return composites[len(composites)-1]
}

// sortStatesForLayout gives the layout a stable order — states grouped under
// their parent, top-level states first — because the parser collects them from
// a map, whose iteration order changes between runs.
func sortStatesForLayout(states []StateNode) {
	sort.SliceStable(states, func(a, b int) bool {
		if states[a].Parent != states[b].Parent {
			return states[a].Parent < states[b].Parent
		}
		return states[a].ID < states[b].ID
	})
}

// stateTypeMarker reads the "<<fork>>", "<<join>>" or "<<choice>>" marker that
// may follow a state declaration. These draw as bars and diamonds, not boxes.
func stateTypeMarker(line string) string {
	switch {
	case strings.Contains(line, "<<fork>>"), strings.Contains(line, "<<join>>"):
		return stateTypeFork
	case strings.Contains(line, "<<choice>>"):
		return stateTypeChoice
	default:
		return stateTypeNormal
	}
}

// parseStateNote reads "note right of X : text" and its left/above/below
// spellings. The sequence number keeps notes on the same state distinct.
func parseStateNote(line string, sequence int) (*StateNode, bool) {
	lower := strings.ToLower(line)
	if !strings.HasPrefix(lower, "note ") {
		return nil, false
	}
	rest, _, found := strings.Cut(line[len("note "):], ":")
	if !found {
		return nil, false
	}
	_, text, _ := strings.Cut(line, ":")

	fields := strings.Fields(rest)
	target := ""
	if len(fields) > 0 {
		target = fields[len(fields)-1]
	}
	if target == "" {
		return nil, false
	}
	return &StateNode{
		ID:         fmt.Sprintf("__note_%d_%s", sequence, target),
		Label:      strings.TrimSpace(text),
		Type:       stateTypeNote,
		NoteTarget: target,
	}, true
}

func ensureState(states map[string]*StateNode, id string, parent string) {
	if _, ok := states[id]; !ok {
		stateType := stateTypeNormal
		label := id
		if strings.HasPrefix(id, "__start_") {
			stateType = "start"
			label = ""
		}
		if strings.HasPrefix(id, "__end_") {
			stateType = stateTypeEnd
			label = ""
		}
		states[id] = &StateNode{ID: id, Label: label, Type: stateType, Parent: parent}
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
	if before, after, ok := strings.Cut(line, arrowSolid); ok {
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
