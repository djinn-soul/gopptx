package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// MermaidType represents the supported Mermaid diagram types.
type MermaidType string //nolint:revive // keeping exported name for API compatibility

const (
	// Flowchart represents a Mermaid flowchart.
	Flowchart MermaidType = "flowchart"
	// Sequence represents a Mermaid sequence diagram.
	Sequence MermaidType = "sequence"
	// Pie represents a Mermaid pie chart.
	Pie MermaidType = "pie"
	// Gantt represents a Mermaid gantt chart.
	Gantt MermaidType = "gantt"
	// Class represents a Mermaid class diagram.
	Class MermaidType = "class"
	// State represents a Mermaid state diagram.
	State MermaidType = "state"
	// ER represents a Mermaid entity relationship diagram.
	ER MermaidType = "er"
	// Mindmap represents a Mermaid mindmap.
	Mindmap MermaidType = "mindmap"
	// Timeline represents a Mermaid timeline diagram.
	Timeline MermaidType = "timeline"
	// Journey represents a Mermaid user journey diagram.
	Journey MermaidType = "journey"
	// Quadrant represents a Mermaid quadrant chart.
	Quadrant MermaidType = "quadrant"
	// GitGraph represents a Mermaid git graph.
	GitGraph MermaidType = "gitgraph"
	// Unknown represents an unsupported or unrecognized Mermaid diagram type.
	Unknown MermaidType = "unknown"
)

// DiagramBounds defines the bounding box of a diagram.
type DiagramBounds struct {
	X  styling.Length
	Y  styling.Length
	CX styling.Length
	CY styling.Length
}

// DiagramElements contains the generated PowerPoint elements for a Mermaid diagram.
type DiagramElements struct {
	Shapes     []shapes.Shape
	Connectors []shapes.Connector
	Bounds     *DiagramBounds
	Grouped    bool
}

// NodeShape represents the shape of a node in a diagram.
type NodeShape string

const (
	// NodeShapeRectangle represents a rectangular node.
	NodeShapeRectangle NodeShape = "rectangle"
	// NodeShapeRoundedRect represents a rounded rectangular node.
	NodeShapeRoundedRect NodeShape = "rounded_rect"
	// NodeShapeStadium represents a stadium-shaped node.
	NodeShapeStadium NodeShape = "stadium"
	// NodeShapeDiamond represents a diamond-shaped node.
	NodeShapeDiamond NodeShape = "diamond"
	// NodeShapeCircle represents a circular node.
	NodeShapeCircle NodeShape = "circle"
	// NodeShapeHexagon represents a hexagonal node.
	NodeShapeHexagon NodeShape = "hexagon"
)

// ArrowStyle represents the style of an arrow in a connection.
type ArrowStyle string

const (
	// ArrowStyleArrow represents a standard arrow.
	ArrowStyleArrow ArrowStyle = "arrow"
	// ArrowStyleOpen represents an open connection without an arrowhead.
	ArrowStyleOpen ArrowStyle = "open"
	// ArrowStyleDotted represents a dotted arrow connection.
	ArrowStyleDotted ArrowStyle = "dotted"
	// ArrowStyleThick represents a thick arrow connection.
	ArrowStyleThick ArrowStyle = "thick"
)

// FlowDirection represents the layout direction of a flowchart.
type FlowDirection string

const (
	// FlowDirectionLR represents Left to Right layout.
	FlowDirectionLR FlowDirection = "LR"
	// FlowDirectionRL represents Right to Left layout.
	FlowDirectionRL FlowDirection = "RL"
	// FlowDirectionTB represents Top to Bottom layout.
	FlowDirectionTB FlowDirection = "TB"
	// FlowDirectionBT represents Bottom to Top layout.
	FlowDirectionBT FlowDirection = "BT"
)
