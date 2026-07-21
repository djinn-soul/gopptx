package mermaid

import "strings"

// renderFlowchart parses and renders a Mermaid flowchart into PowerPoint elements.
func renderFlowchart(code string, theme Theme) DiagramElements {
	flowchart := parseFlowchart(code)
	return generateFlowchartElements(flowchart, theme)
}

func parseFlowchart(code string) *FlowchartDiagram {
	lines := ParseLines(code)
	if len(lines) == 0 {
		return &FlowchartDiagram{Direction: FlowDirectionTB}
	}

	direction := ExtractDirection(lines[0])
	state := newFlowchartParseState()
	for i := 1; i < len(lines); i++ {
		state.consumeLine(lines[i])
	}

	nodeList := make([]FlowNode, 0, len(state.nodes))
	for _, node := range state.nodes {
		nodeList = append(nodeList, node)
	}

	return &FlowchartDiagram{
		Direction:   direction,
		Nodes:       nodeList,
		Connections: state.connections,
		Subgraphs:   state.subgraphs,
	}
}

type flowchartParseState struct {
	nodes           map[string]FlowNode
	connections     []FlowConnection
	subgraphs       []Subgraph
	currentSubgraph *Subgraph
}

func newFlowchartParseState() *flowchartParseState {
	return &flowchartParseState{
		nodes:       make(map[string]FlowNode),
		connections: make([]FlowConnection, 0),
		subgraphs:   make([]Subgraph, 0),
	}
}

func (s *flowchartParseState) consumeLine(line string) {
	if after, ok := strings.CutPrefix(line, "subgraph"); ok {
		s.currentSubgraph = &Subgraph{Name: strings.TrimSpace(after), Nodes: []string{}}
		return
	}
	if line == "end" {
		s.finishCurrentSubgraph()
		return
	}
	if s.consumeConnection(line) {
		return
	}

	id, label, shape := ParseNodeDef(line)
	s.addNode(id, label, shape)
}

func (s *flowchartParseState) finishCurrentSubgraph() {
	if s.currentSubgraph == nil {
		return
	}
	s.subgraphs = append(s.subgraphs, *s.currentSubgraph)
	s.currentSubgraph = nil
}

func (s *flowchartParseState) consumeConnection(line string) bool {
	fromPart, arrow, rest, found := SplitConnection(line)
	if !found {
		return false
	}

	fromID, fromLabel, fromShape := ParseNodeDef(fromPart)
	s.addNode(fromID, fromLabel, fromShape)

	arrowLabel, toPart := ExtractArrowLabel(rest)
	toID, toLabel, toShape := ParseNodeDef(toPart)
	s.addNode(toID, toLabel, toShape)

	s.connections = append(s.connections, FlowConnection{
		From:       fromID,
		To:         toID,
		Label:      arrowLabel,
		ArrowStyle: parseArrowStyle(arrow),
	})
	return true
}

func (s *flowchartParseState) addNode(id, label string, shape NodeShape) {
	if _, exists := s.nodes[id]; exists {
		return
	}
	s.nodes[id] = FlowNode{ID: id, Label: label, Shape: shape}
	if s.currentSubgraph != nil {
		s.currentSubgraph.Nodes = append(s.currentSubgraph.Nodes, id)
	}
}

func parseArrowStyle(arrow string) ArrowStyle {
	switch arrow {
	case arrowThick:
		return ArrowStyleThick
	case arrowDotted:
		return ArrowStyleDotted
	case arrowOpen:
		return ArrowStyleOpen
	default:
		return ArrowStyleArrow
	}
}
