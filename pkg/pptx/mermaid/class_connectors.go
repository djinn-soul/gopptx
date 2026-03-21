package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func classNeedsHollowInheritanceMarker(relType string) bool {
	return relType == "<|--" || relType == "<|.."
}

func classInheritanceMarker(geometry classConnectorGeometry, theme Theme) *shapes.Shape {
	size := styling.Inches(0.2)
	x, y, rotation := classInheritanceMarkerPlacement(geometry, size)
	shape := shapes.NewShape(shapes.ShapeTypeTriangle, x, y, size, size).
		WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
		WithRotation(rotation)
	return &shape
}

func classInheritanceMarkerPlacement(
	geometry classConnectorGeometry,
	size styling.Length,
) (styling.Length, styling.Length, int) {
	switch geometry.endSite {
	case shapes.ConnectionSiteLeft:
		return geometry.endX, geometry.endY - size/2, -90
	case shapes.ConnectionSiteRight:
		return geometry.endX - size, geometry.endY - size/2, 90
	case shapes.ConnectionSiteBottom:
		return geometry.endX - size/2, geometry.endY - size, 180
	default:
		return geometry.endX - size/2, geometry.endY, 0
	}
}

func classConnectorPoints(fromPos, toPos classPosition, layout classLayout) classConnectorGeometry {
	dx := fromPos.x - toPos.x
	if dx < 0 {
		dx = -dx
	}
	dy := fromPos.y - toPos.y
	if dy < 0 {
		dy = -dy
	}
	if dx > dy {
		if fromPos.x < toPos.x {
			return classConnectorGeometry{
				startX:    fromPos.x + layout.classWidth,
				startY:    fromPos.y + layout.headerHeight/2,
				endX:      toPos.x,
				endY:      toPos.y + layout.headerHeight/2,
				startSite: shapes.ConnectionSiteRight,
				endSite:   shapes.ConnectionSiteLeft,
			}
		}
		return classConnectorGeometry{
			startX:    fromPos.x,
			startY:    fromPos.y + layout.headerHeight/2,
			endX:      toPos.x + layout.classWidth,
			endY:      toPos.y + layout.headerHeight/2,
			startSite: shapes.ConnectionSiteLeft,
			endSite:   shapes.ConnectionSiteRight,
		}
	}
	if fromPos.y < toPos.y {
		return classConnectorGeometry{
			startX:    fromPos.x + layout.classWidth/2,
			startY:    fromPos.y + layout.headerHeight,
			endX:      toPos.x + layout.classWidth/2,
			endY:      toPos.y,
			startSite: shapes.ConnectionSiteBottom,
			endSite:   shapes.ConnectionSiteTop,
		}
	}
	return classConnectorGeometry{
		startX:    fromPos.x + layout.classWidth/2,
		startY:    fromPos.y,
		endX:      toPos.x + layout.classWidth/2,
		endY:      toPos.y + layout.headerHeight,
		startSite: shapes.ConnectionSiteTop,
		endSite:   shapes.ConnectionSiteBottom,
	}
}

func classArrowTypes(relType string) (string, string) {
	startArrow := shapes.ArrowTypeNone
	endArrow := shapes.ArrowTypeTriangle
	switch relType {
	case "<|--", "<|..":
		endArrow = shapes.ArrowTypeOpen
	case "*--", "*..":
		startArrow = shapes.ArrowTypeDiamond
	case "o--", "o..":
		startArrow = shapes.ArrowTypeOval
	case "--", "..":
		endArrow = shapes.ArrowTypeNone
	}
	return startArrow, endArrow
}

func (b *classBounds) include(x, y, cx, cy styling.Length) {
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

func (s *classRenderState) diagramElements() DiagramElements {
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
