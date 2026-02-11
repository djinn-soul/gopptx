package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type (
	// Shape is one auto shape.
	Shape = elements.Shape
	// ShapeFill configures solid fill properties for one shape.
	ShapeFill = elements.ShapeFill
	// ShapeLine configures line style for one shape or connector.
	ShapeLine = elements.ShapeLine
)

const (
	ShapeTypeRectangle           = elements.ShapeTypeRectangle
	ShapeTypeRoundedRectangle    = elements.ShapeTypeRoundedRectangle
	ShapeTypeEllipse             = elements.ShapeTypeEllipse
	ShapeTypeTriangle            = elements.ShapeTypeTriangle
	ShapeTypeRightTriangle       = elements.ShapeTypeRightTriangle
	ShapeTypeDiamond             = elements.ShapeTypeDiamond
	ShapeTypePentagon            = elements.ShapeTypePentagon
	ShapeTypeHexagon             = elements.ShapeTypeHexagon
	ShapeTypeParallelogram       = elements.ShapeTypeParallelogram
	ShapeTypeFlowChartProcess    = elements.ShapeTypeFlowChartProcess
	ShapeTypeFlowChartDecision   = elements.ShapeTypeFlowChartDecision
	ShapeTypeFlowChartTerminator = elements.ShapeTypeFlowChartTerminator
	ShapeTypeRightArrow          = elements.ShapeTypeRightArrow
	ShapeTypeLeftArrow           = elements.ShapeTypeLeftArrow
	ShapeTypeUpArrow             = elements.ShapeTypeUpArrow
	ShapeTypeDownArrow           = elements.ShapeTypeDownArrow
	ShapeTypeCloud               = elements.ShapeTypeCloud
	ShapeTypeStar5               = elements.ShapeTypeStar5
	ShapeTypeHeart               = elements.ShapeTypeHeart
	ShapeTypeFlowChartDocument   = elements.ShapeTypeFlowChartDocument
	ShapeTypeFlowChartData       = elements.ShapeTypeFlowChartData
)

func NewShape(shapeType string, x, y, cx, cy int64) Shape {
	return elements.NewShape(shapeType, x, y, cx, cy)
}

func NewShapeFill(color string) ShapeFill {
	return elements.NewShapeFill(color)
}

func NewShapeLine(color string, width int64) ShapeLine {
	return elements.NewShapeLine(color, width)
}

func normalizeShapeType(shapeType string) string {
	return elements.NormalizeShapeType(shapeType)
}
