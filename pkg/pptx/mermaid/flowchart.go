package mermaid

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// FlowNode represents a node in a flowchart.
type FlowNode struct {
	ID    string
	Label string
	Shape NodeShape
}

// FlowConnection represents a connection between two nodes in a flowchart.
type FlowConnection struct {
	From       string
	To         string
	Label      string
	ArrowStyle ArrowStyle
}

// Subgraph represents a grouped set of nodes in a flowchart.
type Subgraph struct {
	Name  string
	Nodes []string
}

// FlowchartDiagram represents the parsed structure of a Mermaid flowchart.
type FlowchartDiagram struct {
	Direction   FlowDirection
	Nodes       []FlowNode
	Connections []FlowConnection
	Subgraphs   []Subgraph
}

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

	// Skip the header line (graph/flowchart ...)
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

func (s *flowchartParseState) addNode(id string, label string, shape NodeShape) {
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
	case "==>":
		return ArrowStyleThick
	case "-.->":
		return ArrowStyleDotted
	case "---":
		return ArrowStyleOpen
	default:
		return ArrowStyleArrow
	}
}

func generateFlowchartElements(flowchart *FlowchartDiagram, theme Theme) DiagramElements {
	nodeCount := len(flowchart.Nodes)
	if nodeCount == 0 {
		return DiagramElements{Grouped: true}
	}
	isHorizontal := flowchart.Direction == FlowDirectionLR || flowchart.Direction == FlowDirectionRL
	layout := defaultFlowchartLayout()
	state := newFlowchartRenderState(layout, theme, isHorizontal, flowchart.Nodes)
	state.layoutNodes(flowchart.Subgraphs, flowchart.Connections)
	state.addConnectors(flowchart.Connections)
	return state.diagramElements()
}

type flowchartLayout struct {
	baseNodeWidth  styling.Length
	nodeHeight     styling.Length
	hSpacing       styling.Length
	vSpacing       styling.Length
	subgraphStartX styling.Length
	subgraphStartY styling.Length
	subgraphGapX   styling.Length
	subgraphPadW   styling.Length
	subgraphPadH   styling.Length
	titleHeight    styling.Length
	titlePad       styling.Length
	gridStartX     styling.Length
	gridStartY     styling.Length
	gridMaxCols    int
	labelWidth     styling.Length
	labelHeight    styling.Length
}

type flowchartPoint struct {
	x styling.Length
	y styling.Length
}

type flowchartBounds struct {
	minX  styling.Length
	minY  styling.Length
	maxX  styling.Length
	maxY  styling.Length
	empty bool
}

type flowchartRenderState struct {
	layout          flowchartLayout
	theme           Theme
	isHorizontal    bool
	nodes           []FlowNode
	nodeLookup      map[string]FlowNode
	nodePositions   map[string]flowchartPoint
	nodeWidths      map[string]styling.Length
	nodeHeights     map[string]styling.Length
	nodeShapeIndex  map[string]int
	shapes          []shapes.Shape
	connectors      []shapes.Connector
	bounds          flowchartBounds
	currentSubgraph styling.Length
}

type flowConnectorGeometry struct {
	startSite string
	endSite   string
	startX    styling.Length
	startY    styling.Length
	endX      styling.Length
	endY      styling.Length
}

func defaultFlowchartLayout() flowchartLayout {
	return flowchartLayout{
		baseNodeWidth:  styling.Inches(1.6),
		nodeHeight:     styling.Inches(0.6),
		hSpacing:       styling.Inches(2.5),
		vSpacing:       styling.Inches(1.2),
		subgraphStartX: styling.Inches(0.5),
		subgraphStartY: styling.Inches(1.5),
		subgraphGapX:   styling.Inches(0.8),
		subgraphPadW:   styling.Inches(0.5),
		subgraphPadH:   styling.Inches(0.5),
		titleHeight:    styling.Inches(0.3),
		titlePad:       styling.Inches(0.1),
		gridStartX:     styling.Inches(1.0),
		gridStartY:     styling.Inches(1.5),
		gridMaxCols:    5,
		labelWidth:     styling.Inches(0.55),
		labelHeight:    styling.Inches(0.22),
	}
}

func newFlowchartRenderState(
	layout flowchartLayout,
	theme Theme,
	isHorizontal bool,
	nodes []FlowNode,
) *flowchartRenderState {
	lookup := make(map[string]FlowNode, len(nodes))
	for _, node := range nodes {
		lookup[node.ID] = node
	}
	return &flowchartRenderState{
		layout:          layout,
		theme:           theme,
		isHorizontal:    isHorizontal,
		nodes:           nodes,
		nodeLookup:      lookup,
		nodePositions:   make(map[string]flowchartPoint),
		nodeWidths:      make(map[string]styling.Length),
		nodeHeights:     make(map[string]styling.Length),
		nodeShapeIndex:  make(map[string]int),
		shapes:          make([]shapes.Shape, 0),
		connectors:      make([]shapes.Connector, 0),
		bounds:          flowchartBounds{empty: true},
		currentSubgraph: layout.subgraphStartX,
	}
}

func (s *flowchartRenderState) layoutNodes(subgraphs []Subgraph, connections []FlowConnection) {
	if len(subgraphs) > 0 {
		s.layoutSubgraphs(subgraphs)
		s.layoutOrphanNodes()
		return
	}
	if s.layoutByConnections(connections) {
		return
	}
	s.layoutSimpleGrid()
}

func (s *flowchartRenderState) layoutSubgraphs(subgraphs []Subgraph) {
	for _, sg := range subgraphs {
		sgWidth := s.subgraphWidth(sg) + s.layout.subgraphPadW
		sgHeight := styling.Length(len(sg.Nodes))*s.layout.vSpacing + s.layout.subgraphPadH + s.layout.titleHeight
		sgX := s.currentSubgraph
		sgY := s.layout.subgraphStartY
		s.addSubgraphShape(sg, sgX, sgY, sgWidth, sgHeight)
		s.layoutNodesInSubgraph(sg, sgX, sgY, sgWidth)
		s.currentSubgraph += sgWidth + s.layout.subgraphGapX
	}
}

func (s *flowchartRenderState) subgraphWidth(sg Subgraph) styling.Length {
	maxWidth := s.layout.baseNodeWidth
	for _, nodeID := range sg.Nodes {
		node, ok := s.nodeLookup[nodeID]
		if !ok {
			continue
		}
		width := s.calculateWidth(node.Label)
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

func (s *flowchartRenderState) addSubgraphShape(sg Subgraph, x, y, width, height styling.Length) {
	background := shapes.NewShape(shapes.ShapeTypeRoundedRectangle, x, y, width, height).
		WithFill(shapes.NewShapeFill(s.theme.SecondaryFill)).
		WithLine(shapes.NewShapeLine(s.theme.SecondaryStroke, styling.Emu(12700)))
	s.shapes = append(s.shapes, background)
	s.bounds.include(x, y, width, height)

	title := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		x+s.layout.titlePad,
		y+s.layout.titlePad,
		width-2*s.layout.titlePad,
		s.layout.titleHeight,
	).WithText(sg.Name).
		WithAutoFit(shapes.TextAutoFitNormal)
	title.Fill = nil
	title.Line = nil
	s.shapes = append(s.shapes, title)
}

func (s *flowchartRenderState) layoutNodesInSubgraph(sg Subgraph, sgX, sgY, sgWidth styling.Length) {
	for i, nodeID := range sg.Nodes {
		node, ok := s.nodeLookup[nodeID]
		if !ok {
			continue
		}
		width := s.calculateWidth(node.Label)
		x := sgX + (sgWidth-width)/2
		y := sgY + s.layout.titleHeight + styling.Inches(0.3) + styling.Length(i)*s.layout.vSpacing
		s.addNodeShape(node, x, y, width)
	}
}

func (s *flowchartRenderState) layoutOrphanNodes() {
	orphanY := s.layout.subgraphStartY
	for _, node := range s.nodes {
		if _, exists := s.nodePositions[node.ID]; exists {
			continue
		}
		width := s.calculateWidth(node.Label)
		x := s.currentSubgraph
		y := orphanY
		s.addNodeShape(node, x, y, width)
		orphanY += s.layout.vSpacing
	}
}

func (s *flowchartRenderState) layoutSimpleGrid() {
	cols := 1
	if s.isHorizontal {
		cols = min(len(s.nodes), s.layout.gridMaxCols)
	}
	for i, node := range s.nodes {
		col := i % cols
		row := i / cols
		width := s.calculateWidth(node.Label)
		x := s.layout.gridStartX + styling.Length(col)*s.layout.hSpacing
		y := s.layout.gridStartY + styling.Length(row)*s.layout.vSpacing
		s.addNodeShape(node, x, y, width)
	}
}

func (s *flowchartRenderState) calculateWidth(label string) styling.Length {
	width := s.layout.baseNodeWidth
	if len(label) > 15 {
		width += styling.Length(len(label)-15) * styling.Inches(0.08)
	}
	return width
}

func (s *flowchartRenderState) addNodeShape(node FlowNode, x, y, width styling.Length) {
	if node.Shape == NodeShapeDiamond && width < styling.Inches(3.2) {
		width = styling.Inches(3.2)
	}
	nodeHeight := s.renderedNodeHeight(node.Shape)
	s.nodePositions[node.ID] = flowchartPoint{x: x, y: y}
	s.nodeWidths[node.ID] = width
	shape := createNodeShape(&node, x, y, width, nodeHeight, s.theme)
	s.shapes = append(s.shapes, shape)
	s.nodeShapeIndex[node.ID] = len(s.shapes)
	// Track the exact rendered dimensions so connector anchors match visual geometry.
	s.nodeHeights[node.ID] = nodeHeight
	s.bounds.include(x, y, width, nodeHeight)
}

func (s *flowchartRenderState) renderedNodeHeight(shape NodeShape) styling.Length {
	if shape == NodeShapeDiamond {
		return s.layout.nodeHeight * 2
	}
	return s.layout.nodeHeight
}

func (s *flowchartRenderState) addConnectors(connections []FlowConnection) {
	for _, conn := range connections {
		connector, labelShape, hasLabel, ok := s.buildConnector(conn)
		if !ok {
			continue
		}
		s.connectors = append(s.connectors, connector)
		if hasLabel {
			s.shapes = append(s.shapes, labelShape)
			s.bounds.includeShape(labelShape)
		}
	}
}

func (s *flowchartRenderState) buildConnector(
	conn FlowConnection,
) (shapes.Connector, shapes.Shape, bool, bool) {
	fromPos, fromOK := s.nodePositions[conn.From]
	toPos, toOK := s.nodePositions[conn.To]
	if !fromOK || !toOK {
		return shapes.Connector{}, shapes.Shape{}, false, false
	}
	fromWidth := s.nodeWidths[conn.From]
	toWidth := s.nodeWidths[conn.To]
	fromHeight := s.nodeHeights[conn.From]
	toHeight := s.nodeHeights[conn.To]
	geometry := s.connectorGeometry(fromPos, toPos, fromWidth, toWidth, fromHeight, toHeight)
	connectorType := s.connectorType(geometry)
	lineWidth, lineDash := s.connectorLine(conn.ArrowStyle)
	endArrow := s.connectorEndArrow(conn.ArrowStyle)
	connector := shapes.NewConnector(
		connectorType,
		geometry.startX,
		geometry.startY,
		geometry.endX,
		geometry.endY,
	).WithLine(shapes.NewShapeLine(s.theme.PrimaryStroke, lineWidth).WithDash(lineDash)).
		WithArrows(shapes.ArrowTypeNone, endArrow)
	if idx, ok := s.nodeShapeIndex[conn.From]; ok {
		connector = connector.ConnectStart(idx, geometry.startSite)
	}
	if idx, ok := s.nodeShapeIndex[conn.To]; ok {
		connector = connector.ConnectEnd(idx, geometry.endSite)
	}
	if conn.Label == "" {
		return connector, shapes.Shape{}, false, true
	}
	labelShape := s.connectorLabelShape(conn.Label, geometry)
	return connector, labelShape, true, true
}

func (s *flowchartRenderState) connectorGeometry(
	fromPos, toPos flowchartPoint,
	fromWidth, toWidth styling.Length,
	fromHeight, toHeight styling.Length,
) flowConnectorGeometry {
	if s.isHorizontal {
		return flowConnectorGeometry{
			startSite: shapes.ConnectionSiteRight,
			endSite:   shapes.ConnectionSiteLeft,
			startX:    fromPos.x + fromWidth,
			startY:    fromPos.y + fromHeight/2,
			endX:      toPos.x,
			endY:      toPos.y + toHeight/2,
		}
	}
	return flowConnectorGeometry{
		startSite: shapes.ConnectionSiteBottom,
		endSite:   shapes.ConnectionSiteTop,
		startX:    fromPos.x + fromWidth/2,
		startY:    fromPos.y + fromHeight,
		endX:      toPos.x + toWidth/2,
		endY:      toPos.y,
	}
}

func (s *flowchartRenderState) connectorType(geometry flowConnectorGeometry) string {
	if s.isHorizontal {
		return shapes.ConnectorTypeStraight
	}
	diffX := geometry.startX - geometry.endX
	if diffX < 0 {
		diffX = -diffX
	}
	diffY := geometry.startY - geometry.endY
	if diffY < 0 {
		diffY = -diffY
	}
	if diffX < styling.Inches(0.1) || diffY < styling.Inches(0.1) {
		return shapes.ConnectorTypeStraight
	}
	return shapes.ConnectorTypeElbow
}

func (s *flowchartRenderState) connectorLine(style ArrowStyle) (styling.Length, string) {
	switch style {
	case ArrowStyleThick:
		return s.theme.LineWeight * 2, shapes.LineDashSolid
	case ArrowStyleDotted:
		return s.theme.LineWeight, shapes.LineDashDash
	default:
		return s.theme.LineWeight, shapes.LineDashSolid
	}
}

func (s *flowchartRenderState) connectorEndArrow(style ArrowStyle) string {
	if style == ArrowStyleOpen {
		return shapes.ArrowTypeNone
	}
	return shapes.ArrowTypeTriangle
}

func (s *flowchartRenderState) connectorLabelShape(
	label string,
	geometry flowConnectorGeometry,
) shapes.Shape {
	labelX := (geometry.startX + geometry.endX) / 2
	labelY := (geometry.startY + geometry.endY) / 2
	switch {
	case s.isHorizontal:
		// Keep LR labels off the connector stroke to improve legibility.
		labelY -= styling.Inches(0.16)
	case geometry.endY >= geometry.startY:
		labelY += styling.Inches(0.18)
	default:
		labelY -= styling.Inches(0.18)
	}
	labelShape := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		labelX-s.layout.labelWidth/2,
		labelY-s.layout.labelHeight/2,
		s.layout.labelWidth,
		s.layout.labelHeight,
	).WithText(label).
		WithAutoFit(shapes.TextAutoFitNormal)
	labelShape.Fill = nil
	labelShape.Line = nil
	return labelShape
}

func (s *flowchartRenderState) diagramElements() DiagramElements {
	return DiagramElements{
		Shapes:     s.shapes,
		Connectors: s.connectors,
		Grouped:    true,
		Bounds: &DiagramBounds{
			X:  s.bounds.minX,
			Y:  s.bounds.minY,
			CX: s.bounds.maxX - s.bounds.minX,
			CY: s.bounds.maxY - s.bounds.minY,
		},
	}
}
