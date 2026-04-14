package mermaid

import (
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
