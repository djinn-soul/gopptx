package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

type (
	// Shape is one auto shape.
	Shape = shapes.Shape
	// ShapeDefinition allows external shape builders to plug into slide composition.
	ShapeDefinition = shapes.ShapeDefinition

	// ShapeFill configures solid fill properties for one shape.
	ShapeFill = shapes.ShapeFill
	// ShapeLine configures line style for one shape or connector.
	ShapeLine = shapes.ShapeLine

	// ShapeGradientStop configures one gradient stop for a shape fill.
	ShapeGradientStop = shapes.ShapeGradientStop
	// ShapeGradientFill configures gradient fill properties for one shape.
	ShapeGradientFill = shapes.ShapeGradientFill
)

const (
	ShapeTypeRectangle           = shapes.ShapeTypeRectangle
	ShapeTypeRoundedRectangle    = shapes.ShapeTypeRoundedRectangle
	ShapeTypeEllipse             = shapes.ShapeTypeEllipse
	ShapeTypeTriangle            = shapes.ShapeTypeTriangle
	ShapeTypeRightTriangle       = shapes.ShapeTypeRightTriangle
	ShapeTypeDiamond             = shapes.ShapeTypeDiamond
	ShapeTypePentagon            = shapes.ShapeTypePentagon
	ShapeTypeHexagon             = shapes.ShapeTypeHexagon
	ShapeTypeParallelogram       = shapes.ShapeTypeParallelogram
	ShapeTypeFlowChartProcess    = shapes.ShapeTypeFlowChartProcess
	ShapeTypeFlowChartDecision   = shapes.ShapeTypeFlowChartDecision
	ShapeTypeFlowChartTerminator = shapes.ShapeTypeFlowChartTerminator
	ShapeTypeRightArrow          = shapes.ShapeTypeRightArrow
	ShapeTypeLeftArrow           = shapes.ShapeTypeLeftArrow
	ShapeTypeUpArrow             = shapes.ShapeTypeUpArrow
	ShapeTypeDownArrow           = shapes.ShapeTypeDownArrow
	ShapeTypeCloud               = shapes.ShapeTypeCloud
	ShapeTypeStar5               = shapes.ShapeTypeStar5
	ShapeTypeHeart               = shapes.ShapeTypeHeart
	ShapeTypeFlowChartDocument   = shapes.ShapeTypeFlowChartDocument
	ShapeTypeFlowChartData       = shapes.ShapeTypeFlowChartData

	ShapeGradientTypeLinear      = shapes.ShapeGradientTypeLinear
	ShapeGradientTypeRadial      = shapes.ShapeGradientTypeRadial
	ShapeGradientTypeRectangular = shapes.ShapeGradientTypeRectangular
	ShapeGradientTypePath        = shapes.ShapeGradientTypePath
)

func NewShape(shapeType string, x, y, cx, cy int64) Shape {
	return shapes.NewShape(shapeType, x, y, cx, cy)
}

func NewShapeFill(color string) ShapeFill {
	return shapes.NewShapeFill(color)
}

func NewShapeLine(color string, width int64) ShapeLine {
	return shapes.NewShapeLine(color, width)
}

func NewShapeGradientStop(positionPct int, color string) ShapeGradientStop {
	return shapes.NewShapeGradientStop(positionPct, color)
}

func NewShapeGradientFill(gradientType string, stops []ShapeGradientStop) ShapeGradientFill {
	return shapes.NewShapeGradientFill(gradientType, stops)
}

// Fluent API Macros (Inches based)

func NewRectangle(x, y, w, h float64) Shape {
	return shapes.NewRectangle(x, y, w, h)
}

func NewEllipse(x, y, w, h float64) Shape {
	return shapes.NewEllipse(x, y, w, h)
}

func NewTextBox(text string, x, y, w, h float64) Shape {
	return shapes.NewTextBox(text, x, y, w, h)
}

func NewRoundedRectangle(x, y, w, h float64) Shape {
	return shapes.NewRoundedRectangle(x, y, w, h)
}

func NewTriangle(x, y, w, h float64) Shape {
	return shapes.NewTriangle(x, y, w, h)
}

func NewRightTriangle(x, y, w, h float64) Shape {
	return shapes.NewRightTriangle(x, y, w, h)
}

func NewDiamond(x, y, w, h float64) Shape {
	return shapes.NewDiamond(x, y, w, h)
}

func NewPentagon(x, y, w, h float64) Shape {
	return shapes.NewPentagon(x, y, w, h)
}

func NewHexagon(x, y, w, h float64) Shape {
	return shapes.NewHexagon(x, y, w, h)
}

func NewParallelogram(x, y, w, h float64) Shape {
	return shapes.NewParallelogram(x, y, w, h)
}

func NewFlowChartProcess(x, y, w, h float64) Shape {
	return shapes.NewFlowChartProcess(x, y, w, h)
}

func NewFlowChartDecision(x, y, w, h float64) Shape {
	return shapes.NewFlowChartDecision(x, y, w, h)
}

func NewFlowChartTerminator(x, y, w, h float64) Shape {
	return shapes.NewFlowChartTerminator(x, y, w, h)
}

func NewRightArrow(x, y, w, h float64) Shape {
	return shapes.NewRightArrow(x, y, w, h)
}

func NewLeftArrow(x, y, w, h float64) Shape {
	return shapes.NewLeftArrow(x, y, w, h)
}

func NewUpArrow(x, y, w, h float64) Shape {
	return shapes.NewUpArrow(x, y, w, h)
}

func NewDownArrow(x, y, w, h float64) Shape {
	return shapes.NewDownArrow(x, y, w, h)
}

func NewCloud(x, y, w, h float64) Shape {
	return shapes.NewCloud(x, y, w, h)
}

func NewCircle(x, y, diameter float64) Shape {
	return shapes.NewCircle(x, y, diameter)
}

func NewStar(x, y, size float64) Shape {
	return shapes.NewStar(x, y, size)
}

func NewHeart(x, y, size float64) Shape {
	return shapes.NewHeart(x, y, size)
}

func NewFlowChartDocument(x, y, w, h float64) Shape {
	return shapes.NewFlowChartDocument(x, y, w, h)
}

func NewFlowChartData(x, y, w, h float64) Shape {
	return shapes.NewFlowChartData(x, y, w, h)
}

func NewBadge(text string, x, y float64, color string) Shape {
	return shapes.NewBadge(text, x, y, color)
}
