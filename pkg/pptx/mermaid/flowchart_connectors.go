package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

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
	return connector, s.connectorLabelShape(conn.Label, geometry), true, true
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
