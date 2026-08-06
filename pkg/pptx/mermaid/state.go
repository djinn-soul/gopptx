package mermaid

import (
	"fmt"
	"sort"
	"strings"
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
