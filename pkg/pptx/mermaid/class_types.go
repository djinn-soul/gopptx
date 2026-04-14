package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// ClassNode represents a class in a class diagram.
type ClassNode struct {
	ID         string
	Name       string
	Attributes []string
	Methods    []string
}

// ClassRelationship represents a relationship between two classes.
type ClassRelationship struct {
	From  string
	To    string
	Type  string
	Label string
}

// ClassDiagram represents the parsed structure of a Mermaid class diagram.
type ClassDiagram struct {
	Classes       []ClassNode
	Relationships []ClassRelationship
}

type classLayout struct {
	classWidth   styling.Length
	headerHeight styling.Length
	itemHeight   styling.Length
	hSpacing     styling.Length
	vSpacing     styling.Length
	startX       styling.Length
	startY       styling.Length
	cols         int
}

type classRenderState struct {
	layout       classLayout
	shapes       []shapes.Shape
	connectors   []shapes.Connector
	positions    map[string]classPosition
	shapeIndices map[string]int
	bounds       classBounds
}

type classPosition struct {
	x styling.Length
	y styling.Length
}

type classBounds struct {
	minX  styling.Length
	minY  styling.Length
	maxX  styling.Length
	maxY  styling.Length
	empty bool
}

type classNodeShapes struct {
	position    classPosition
	totalHeight styling.Length
	header      shapes.Shape
	attrBox     shapes.Shape
	methodBox   shapes.Shape
}

type classConnectorGeometry struct {
	startX    styling.Length
	startY    styling.Length
	endX      styling.Length
	endY      styling.Length
	startSite string
	endSite   string
}

func defaultClassLayout() classLayout {
	return classLayout{
		classWidth:   styling.Inches(2.2),
		headerHeight: styling.Inches(0.5),
		itemHeight:   styling.Inches(0.35),
		hSpacing:     styling.Inches(3.0),
		vSpacing:     styling.Inches(2.5),
		startX:       styling.Inches(1.0),
		startY:       styling.Inches(1.0),
		cols:         3,
	}
}

func newClassRenderState(layout classLayout) *classRenderState {
	return &classRenderState{
		layout:       layout,
		positions:    make(map[string]classPosition),
		shapeIndices: make(map[string]int),
		bounds:       classBounds{empty: true},
	}
}

func classPositionForIndex(index int, layout classLayout) classPosition {
	col := index % layout.cols
	row := index / layout.cols
	return classPosition{
		x: layout.startX + styling.Length(col)*layout.hSpacing,
		y: layout.startY + styling.Length(row)*layout.vSpacing,
	}
}
