package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
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

	// TextFrame configures the text layout within a shape.
	TextFrame = shapes.TextFrame
	// TextFrameAnchor specifies the vertical alignment of text within its shape.
	TextFrameAnchor = shapes.TextFrameAnchor
	// TextFrameWrap specifies how text wraps within the shape's text frame.
	TextFrameWrap = shapes.TextFrameWrap
	// TextFrameAutoFit specifies how text is automatically resized or how the shape is resized.
	TextFrameAutoFit = shapes.TextFrameAutoFit

	// Length represents a distance in English Metric Units (EMU).
	Length = styling.Length
)

const (
	ShapeTypeRectangle        = shapes.ShapeTypeRectangle
	ShapeTypeRoundedRectangle = shapes.ShapeTypeRoundedRectangle
	ShapeTypeEllipse          = shapes.ShapeTypeEllipse
	ShapeTypeTriangle         = shapes.ShapeTypeTriangle
	ShapeTypeRightTriangle    = shapes.ShapeTypeRightTriangle
	ShapeTypeDiamond          = shapes.ShapeTypeDiamond
	ShapeTypePentagon         = shapes.ShapeTypePentagon
	ShapeTypeHexagon          = shapes.ShapeTypeHexagon
	ShapeTypeParallelogram    = shapes.ShapeTypeParallelogram
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

	TextAnchorTop    = shapes.TextAnchorTop
	TextAnchorMiddle = shapes.TextAnchorMiddle
	TextAnchorBottom = shapes.TextAnchorBottom

	TextWrapNone   = shapes.TextWrapNone
	TextWrapSquare = shapes.TextWrapSquare

	TextAutoFitNone   = shapes.TextAutoFitNone
	TextAutoFitShape  = shapes.TextAutoFitShape
	TextAutoFitNormal = shapes.TextAutoFitNormal
)


func NewShape(shapeType string, x, y, cx, cy Length) Shape {
	return shapes.NewShape(shapeType, x, y, cx, cy)
}

func NewShapeFill(color string) ShapeFill {
	return shapes.NewShapeFill(color)
}

func NewShapeLine(color string, width Length) ShapeLine {
	return shapes.NewShapeLine(color, width)
}

func NewShapeGradientStop(positionPct int, color string) ShapeGradientStop {
	return shapes.NewShapeGradientStop(positionPct, color)
}

func NewShapeGradientFill(gradientType string, stops []ShapeGradientStop) ShapeGradientFill {
	return shapes.NewShapeGradientFill(gradientType, stops)
}

func NewTextFrame() TextFrame {
	return shapes.NewTextFrame()
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

func NewLeftRightArrow(x, y, w, h float64) Shape {
	return shapes.NewLeftRightArrow(x, y, w, h)
}

func NewUpDownArrow(x, y, w, h float64) Shape {
	return shapes.NewUpDownArrow(x, y, w, h)
}

func NewQuadArrow(x, y, w, h float64) Shape {
	return shapes.NewQuadArrow(x, y, w, h)
}

func NewBentArrow(x, y, w, h float64) Shape {
	return shapes.NewBentArrow(x, y, w, h)
}

func NewUturnArrow(x, y, w, h float64) Shape {
	return shapes.NewUturnArrow(x, y, w, h)
}

func NewCircularArrow(x, y, w, h float64) Shape {
	return shapes.NewCircularArrow(x, y, w, h)
}

func NewChevron(x, y, w, h float64) Shape {
	return shapes.NewChevron(x, y, w, h)
}

func NewWedgeRectCallout(x, y, w, h float64) Shape {
	return shapes.NewWedgeRectCallout(x, y, w, h)
}

func NewWedgeEllipseCallout(x, y, w, h float64) Shape {
	return shapes.NewWedgeEllipseCallout(x, y, w, h)
}

func NewCloudCallout(x, y, w, h float64) Shape {
	return shapes.NewCloudCallout(x, y, w, h)
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

func NewStar4(x, y, size float64) Shape {
	return shapes.NewStar4(x, y, size)
}

func NewStar6(x, y, size float64) Shape {
	return shapes.NewStar6(x, y, size)
}

func NewStar8(x, y, size float64) Shape {
	return shapes.NewStar8(x, y, size)
}

func NewStar12(x, y, size float64) Shape {
	return shapes.NewStar12(x, y, size)
}

func NewStar16(x, y, size float64) Shape {
	return shapes.NewStar16(x, y, size)
}

func NewStar24(x, y, size float64) Shape {
	return shapes.NewStar24(x, y, size)
}

func NewStar32(x, y, size float64) Shape {
	return shapes.NewStar32(x, y, size)
}

func NewRibbon(x, y, w, h float64) Shape {
	return shapes.NewRibbon(x, y, w, h)
}

func NewWave(x, y, w, h float64) Shape {
	return shapes.NewWave(x, y, w, h)
}

func NewSeal(x, y, size float64) Shape {
	return shapes.NewSeal(x, y, size)
}

func NewActionButtonHome(x, y, size float64) Shape {
	return shapes.NewActionButtonHome(x, y, size)
}

func NewActionButtonHelp(x, y, size float64) Shape {
	return shapes.NewActionButtonHelp(x, y, size)
}

func NewActionButtonInformation(x, y, size float64) Shape {
	return shapes.NewActionButtonInformation(x, y, size)
}

func NewActionButtonBack(x, y, size float64) Shape {
	return shapes.NewActionButtonBack(x, y, size)
}

func NewActionButtonForward(x, y, size float64) Shape {
	return shapes.NewActionButtonForward(x, y, size)
}


func NewFlowChartDocument(x, y, w, h float64) Shape {
	return shapes.NewFlowChartDocument(x, y, w, h)
}

func NewFlowChartData(x, y, w, h float64) Shape {
	return shapes.NewFlowChartData(x, y, w, h)
}

func NewOctagon(x, y, w, h float64) Shape {
	return shapes.NewOctagon(x, y, w, h)
}

func NewTrapezoid(x, y, w, h float64) Shape {
	return shapes.NewTrapezoid(x, y, w, h)
}

func NewCube(x, y, w, h float64) Shape {
	return shapes.NewCube(x, y, w, h)
}

func NewFlowChartPredefinedProcess(x, y, w, h float64) Shape {
	return shapes.NewFlowChartPredefinedProcess(x, y, w, h)
}

func NewFlowChartInternalStorage(x, y, w, h float64) Shape {
	return shapes.NewFlowChartInternalStorage(x, y, w, h)
}

func NewFlowChartManualInput(x, y, w, h float64) Shape {
	return shapes.NewFlowChartManualInput(x, y, w, h)
}

func NewFlowChartManualOperation(x, y, w, h float64) Shape {
	return shapes.NewFlowChartManualOperation(x, y, w, h)
}

func NewFlowChartConnector(x, y, w, h float64) Shape {
	return shapes.NewFlowChartConnector(x, y, w, h)
}

func NewFlowChartOffpageConnector(x, y, w, h float64) Shape {
	return shapes.NewFlowChartOffpageConnector(x, y, w, h)
}

func NewFlowChartPunchedCard(x, y, w, h float64) Shape {
	return shapes.NewFlowChartPunchedCard(x, y, w, h)
}

func NewFlowChartPunchedTape(x, y, w, h float64) Shape {
	return shapes.NewFlowChartPunchedTape(x, y, w, h)
}

func NewFlowChartSummingJunction(x, y, w, h float64) Shape {
	return shapes.NewFlowChartSummingJunction(x, y, w, h)
}

func NewFlowChartOr(x, y, w, h float64) Shape {
	return shapes.NewFlowChartOr(x, y, w, h)
}

func NewFlowChartCollate(x, y, w, h float64) Shape {
	return shapes.NewFlowChartCollate(x, y, w, h)
}

func NewFlowChartSort(x, y, w, h float64) Shape {
	return shapes.NewFlowChartSort(x, y, w, h)
}

func NewFlowChartExtract(x, y, w, h float64) Shape {
	return shapes.NewFlowChartExtract(x, y, w, h)
}

func NewFlowChartMerge(x, y, w, h float64) Shape {
	return shapes.NewFlowChartMerge(x, y, w, h)
}

func NewFlowChartOnlineStorage(x, y, w, h float64) Shape {
	return shapes.NewFlowChartOnlineStorage(x, y, w, h)
}

func NewFlowChartDelay(x, y, w, h float64) Shape {
	return shapes.NewFlowChartDelay(x, y, w, h)
}

func NewFlowChartMagneticTape(x, y, w, h float64) Shape {
	return shapes.NewFlowChartMagneticTape(x, y, w, h)
}

func NewFlowChartMagneticDisk(x, y, w, h float64) Shape {
	return shapes.NewFlowChartMagneticDisk(x, y, w, h)
}

func NewFlowChartMagneticDrum(x, y, w, h float64) Shape {
	return shapes.NewFlowChartMagneticDrum(x, y, w, h)
}

func NewFlowChartDisplay(x, y, w, h float64) Shape {
	return shapes.NewFlowChartDisplay(x, y, w, h)
}

func NewFlowChartPreparation(x, y, w, h float64) Shape {
	return shapes.NewFlowChartPreparation(x, y, w, h)
}


func NewBadge(text string, x, y float64, color string) Shape {
	return shapes.NewBadge(text, x, y, color)
}
