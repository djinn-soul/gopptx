package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func (b *flowchartBounds) includeShape(shape shapes.Shape) {
	b.include(shape.X, shape.Y, shape.CX, shape.CY)
}

func (b *flowchartBounds) include(x, y, cx, cy styling.Length) {
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

func createNodeShape(node *FlowNode, x, y, width, height styling.Length, theme Theme) shapes.Shape {
	shapeType := shapes.ShapeTypeRectangle
	switch node.Shape {
	case NodeShapeRectangle:
		shapeType = shapes.ShapeTypeRectangle
	case NodeShapeRoundedRect:
		shapeType = shapes.ShapeTypeRoundedRectangle
	case NodeShapeStadium:
		shapeType = shapes.ShapeTypeRoundedRectangle // Stadium is often represented as rounded rect in PPT
	case NodeShapeDiamond:
		shapeType = shapes.ShapeTypeDiamond
	case NodeShapeCircle:
		shapeType = shapes.ShapeTypeEllipse
	case NodeShapeHexagon:
		shapeType = shapes.ShapeTypeHexagon
	}

	fillColor := theme.PrimaryFill
	if node.Shape == NodeShapeDiamond {
		fillColor = theme.SecondaryFill
	}

	return shapes.NewShape(shapeType, x, y, width, height).
		WithFill(shapes.NewShapeFill(fillColor)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
		WithText(node.Label).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
}
