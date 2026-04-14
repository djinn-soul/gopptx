package mermaid

import "github.com/djinn-soul/gopptx/pkg/pptx/styling"

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
	background := createSubgraphBackground(sg.Name, x, y, width, height, s.theme, s.layout)
	title := createSubgraphTitle(sg.Name, x, y, width, s.layout)
	s.shapes = append(s.shapes, background, title)
	s.bounds.include(x, y, width, height)
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
		s.addNodeShape(node, s.currentSubgraph, orphanY, width)
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
	s.shapes = append(s.shapes, createNodeShape(&node, x, y, width, nodeHeight, s.theme))
	s.nodeShapeIndex[node.ID] = len(s.shapes)
	s.nodeHeights[node.ID] = nodeHeight
	s.bounds.include(x, y, width, nodeHeight)
}

func (s *flowchartRenderState) renderedNodeHeight(shape NodeShape) styling.Length {
	if shape == NodeShapeDiamond {
		return s.layout.nodeHeight * 2
	}
	return s.layout.nodeHeight
}
