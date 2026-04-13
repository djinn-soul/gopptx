package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// EREntity represents an entity in an ER diagram.
type EREntity struct {
	Name       string
	Attributes []string
}

// ERRelationship represents a relationship between two entities.
type ERRelationship struct {
	From  string
	To    string
	Type  string
	Label string
}

// ERDiagram represents the parsed structure of a Mermaid ER diagram.
type ERDiagram struct {
	Entities      []EREntity
	Relationships []ERRelationship
}

type erLayout struct {
	entityWidth   styling.Length
	headerHeight  styling.Length
	itemHeight    styling.Length
	hSpacing      styling.Length
	vSpacing      styling.Length
	startX        styling.Length
	startY        styling.Length
	cols          int
	labelWidth    styling.Length
	labelHeight   styling.Length
	labelMarginX  styling.Length
	labelMarginY  styling.Length
	entityLineWgt styling.Length
}

type erPosition struct {
	x styling.Length
	y styling.Length
}

type erBounds struct {
	minX  styling.Length
	minY  styling.Length
	maxX  styling.Length
	maxY  styling.Length
	empty bool
}

type erRenderState struct {
	layout       erLayout
	shapes       []shapes.Shape
	connectors   []shapes.Connector
	positions    map[string]erPosition
	shapeIndices map[string]int
	bounds       erBounds
}

type erEntityShapes struct {
	position    erPosition
	totalHeight styling.Length
	header      shapes.Shape
	attrBox     shapes.Shape
	hasAttrBox  bool
}

type erConnectorGeometry struct {
	startX    styling.Length
	startY    styling.Length
	endX      styling.Length
	endY      styling.Length
	startSite string
	endSite   string
}

func defaultERLayout() erLayout {
	return erLayout{
		entityWidth:   styling.Inches(2.2),
		headerHeight:  styling.Inches(0.5),
		itemHeight:    styling.Inches(0.35),
		hSpacing:      styling.Inches(3.0),
		vSpacing:      styling.Inches(2.5),
		startX:        styling.Inches(1.0),
		startY:        styling.Inches(1.0),
		cols:          3,
		labelWidth:    styling.Inches(1.0),
		labelHeight:   styling.Inches(0.4),
		labelMarginX:  styling.Inches(0.05),
		labelMarginY:  styling.Inches(0.02),
		entityLineWgt: styling.Emu(12700),
	}
}

func newERRenderState(layout erLayout) *erRenderState {
	return &erRenderState{
		layout:       layout,
		positions:    make(map[string]erPosition),
		shapeIndices: make(map[string]int),
		bounds:       erBounds{empty: true},
	}
}

func erPositionForIndex(index int, layout erLayout) erPosition {
	col := index % layout.cols
	row := index / layout.cols
	return erPosition{
		x: layout.startX + styling.Length(col)*layout.hSpacing,
		y: layout.startY + styling.Length(row)*layout.vSpacing,
	}
}
