package mermaid

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func generateERElements(diagram *ERDiagram, theme Theme) DiagramElements {
	if len(diagram.Entities) == 0 {
		return DiagramElements{Grouped: true}
	}

	layout := defaultERLayout()
	state := newERRenderState(layout)
	for index, entity := range diagram.Entities {
		state.addEntity(entity.Name, erShapesForEntity(entity, index, layout, theme))
	}
	for _, rel := range diagram.Relationships {
		connector, labelShape, hasLabel, ok := erRelationshipConnector(rel, state, layout, theme)
		if !ok {
			continue
		}
		state.connectors = append(state.connectors, connector)
		if hasLabel {
			state.shapes = append(state.shapes, labelShape)
			state.bounds.includeShape(labelShape)
		}
	}
	return state.diagramElements()
}

func erShapesForEntity(entity EREntity, index int, layout erLayout, theme Theme) erEntityShapes {
	position := erPositionForIndex(index, layout)
	attrCount := len(entity.Attributes)
	totalHeight := layout.headerHeight + styling.Length(attrCount)*layout.itemHeight
	header := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		position.x,
		position.y,
		layout.entityWidth,
		layout.headerHeight,
	).WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
		WithText(entity.Name).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
	if attrCount == 0 {
		return erEntityShapes{position: position, totalHeight: totalHeight, header: header}
	}

	attrY := position.y + layout.headerHeight
	attrHeight := styling.Length(attrCount) * layout.itemHeight
	attrBox := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		position.x,
		attrY,
		layout.entityWidth,
		attrHeight,
	).WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, layout.entityLineWgt)).
		WithText(strings.Join(entity.Attributes, "\n")).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithVerticalAnchor(shapes.TextAnchorTop).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
	return erEntityShapes{
		position:    position,
		totalHeight: totalHeight,
		header:      header,
		attrBox:     attrBox,
		hasAttrBox:  true,
	}
}

func (s *erRenderState) addEntity(entityName string, entityShapes erEntityShapes) {
	s.positions[entityName] = entityShapes.position
	s.shapes = append(s.shapes, entityShapes.header)
	s.shapeIndices[entityName] = len(s.shapes)
	s.bounds.include(entityShapes.position.x, entityShapes.position.y, s.layout.entityWidth, entityShapes.totalHeight)
	if entityShapes.hasAttrBox {
		s.shapes = append(s.shapes, entityShapes.attrBox)
	}
}

func erRelationshipConnector(
	rel ERRelationship,
	state *erRenderState,
	layout erLayout,
	theme Theme,
) (shapes.Connector, shapes.Shape, bool, bool) {
	fromPos, fromExists := state.positions[rel.From]
	toPos, toExists := state.positions[rel.To]
	if !fromExists || !toExists {
		return shapes.Connector{}, shapes.Shape{}, false, false
	}

	geometry := erConnectorPoints(fromPos, toPos, layout)
	endArrow := shapes.ArrowTypeTriangle
	if strings.Contains(rel.Type, "{") || strings.Contains(rel.Type, "}") {
		endArrow = shapes.ArrowTypeStealth
	}
	connector := shapes.NewConnector(
		shapes.ConnectorTypeElbow,
		geometry.startX,
		geometry.startY,
		geometry.endX,
		geometry.endY,
	).WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
		WithArrows(shapes.ArrowTypeNone, endArrow)
	if idx, ok := state.shapeIndices[rel.From]; ok {
		connector = connector.ConnectStart(idx, geometry.startSite)
	}
	if idx, ok := state.shapeIndices[rel.To]; ok {
		connector = connector.ConnectEnd(idx, geometry.endSite)
	}
	if rel.Label == "" {
		return connector, shapes.Shape{}, false, true
	}

	midX := (geometry.startX + geometry.endX) / 2
	midY := (geometry.startY + geometry.endY) / 2
	labelShape := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		midX-layout.labelWidth/2,
		midY-layout.labelHeight/2,
		layout.labelWidth,
		layout.labelHeight,
	).WithFill(shapes.NewShapeFill(theme.Background)).
		WithText(rel.Label).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(layout.labelMarginX, layout.labelMarginY, layout.labelMarginX, layout.labelMarginY)
	return connector, labelShape, true, true
}

func erConnectorPoints(fromPos, toPos erPosition, layout erLayout) erConnectorGeometry {
	if fromPos.y < toPos.y {
		return erConnectorGeometry{
			startX:    fromPos.x + layout.entityWidth/2,
			startY:    fromPos.y + layout.headerHeight,
			endX:      toPos.x + layout.entityWidth/2,
			endY:      toPos.y,
			startSite: shapes.ConnectionSiteBottom,
			endSite:   shapes.ConnectionSiteTop,
		}
	}
	return erConnectorGeometry{
		startX:    fromPos.x + layout.entityWidth/2,
		startY:    fromPos.y,
		endX:      toPos.x + layout.entityWidth/2,
		endY:      toPos.y + layout.headerHeight,
		startSite: shapes.ConnectionSiteTop,
		endSite:   shapes.ConnectionSiteBottom,
	}
}

func (b *erBounds) includeShape(shape shapes.Shape) {
	b.include(shape.X, shape.Y, shape.CX, shape.CY)
}

func (b *erBounds) include(x, y, cx, cy styling.Length) {
	if b.empty {
		b.minX, b.minY = x, y
		b.maxX, b.maxY = x+cx, y+cy
		b.empty = false
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

func (s *erRenderState) diagramElements() DiagramElements {
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
