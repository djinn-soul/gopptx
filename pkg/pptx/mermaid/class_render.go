package mermaid

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func generateClassElements(diagram *ClassDiagram, theme Theme) DiagramElements {
	if len(diagram.Classes) == 0 {
		return DiagramElements{Grouped: true}
	}

	layout := defaultClassLayout()
	state := newClassRenderState(layout)
	for index, class := range diagram.Classes {
		state.addClass(class.ID, classShapesForNode(class, index, layout, theme))
	}
	for _, rel := range diagram.Relationships {
		connector, marker, ok := classRelationshipConnector(rel, state, layout, theme)
		if !ok {
			continue
		}
		state.connectors = append(state.connectors, connector)
		if marker != nil {
			state.shapes = append(state.shapes, *marker)
		}
	}
	return state.diagramElements()
}

func classHeights(class ClassNode, layout classLayout) (styling.Length, styling.Length, styling.Length) {
	attrCount := len(class.Attributes)
	methodCount := len(class.Methods)
	attrHeight := styling.Length(attrCount) * layout.itemHeight
	methodHeight := styling.Length(methodCount) * layout.itemHeight
	total := layout.headerHeight + attrHeight + methodHeight
	return attrHeight, methodHeight, total
}

func classShapesForNode(class ClassNode, index int, layout classLayout, theme Theme) classNodeShapes {
	position := classPositionForIndex(index, layout)
	attrHeight, methodHeight, totalHeight := classHeights(class, layout)
	header := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		position.x,
		position.y,
		layout.classWidth,
		layout.headerHeight,
	).WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
		WithText(class.Name).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))

	attrY := position.y + layout.headerHeight
	attrBox := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		position.x,
		attrY,
		layout.classWidth,
		attrHeight,
	).WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, styling.Emu(12700))).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithVerticalAnchor(shapes.TextAnchorTop).
		WithText(strings.Join(class.Attributes, "\n")).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))

	methodY := attrY + attrHeight
	methodBox := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		position.x,
		methodY,
		layout.classWidth,
		methodHeight,
	).WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, styling.Emu(12700))).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithVerticalAnchor(shapes.TextAnchorTop).
		WithText(strings.Join(class.Methods, "\n")).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))

	return classNodeShapes{
		position:    position,
		totalHeight: totalHeight,
		header:      header,
		attrBox:     attrBox,
		methodBox:   methodBox,
	}
}

func (s *classRenderState) addClass(classID string, nodeShapes classNodeShapes) {
	s.positions[classID] = nodeShapes.position
	s.shapes = append(s.shapes, nodeShapes.header)
	s.shapeIndices[classID] = len(s.shapes)
	s.bounds.include(nodeShapes.position.x, nodeShapes.position.y, s.layout.classWidth, nodeShapes.totalHeight)
	if nodeShapes.attrBox.CY > 0 {
		s.shapes = append(s.shapes, nodeShapes.attrBox)
	}
	if nodeShapes.methodBox.CY > 0 {
		s.shapes = append(s.shapes, nodeShapes.methodBox)
	}
}

func classRelationshipConnector(
	rel ClassRelationship,
	state *classRenderState,
	layout classLayout,
	theme Theme,
) (shapes.Connector, *shapes.Shape, bool) {
	fromID := rel.From
	toID := rel.To
	if strings.HasPrefix(rel.Type, "<|") {
		fromID, toID = rel.To, rel.From
	}

	fromPos, fromExists := state.positions[fromID]
	toPos, toExists := state.positions[toID]
	if !fromExists || !toExists {
		return shapes.Connector{}, nil, false
	}

	geometry := classConnectorPoints(fromPos, toPos, layout)
	line := shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)
	if strings.Contains(rel.Type, "..") {
		line = line.WithDash(shapes.LineDashDash)
	}
	startArrow, endArrow := classArrowTypes(rel.Type)
	var marker *shapes.Shape
	if classNeedsHollowInheritanceMarker(rel.Type) {
		endArrow = shapes.ArrowTypeNone
		marker = classInheritanceMarker(geometry, theme)
	}
	connector := shapes.NewConnector(
		shapes.ConnectorTypeElbow,
		geometry.startX,
		geometry.startY,
		geometry.endX,
		geometry.endY,
	).WithLine(line).WithArrows(startArrow, endArrow)
	if idx, ok := state.shapeIndices[fromID]; ok {
		connector = connector.ConnectStart(idx, geometry.startSite)
	}
	if idx, ok := state.shapeIndices[toID]; ok {
		connector = connector.ConnectEnd(idx, geometry.endSite)
	}
	return connector, marker, true
}
