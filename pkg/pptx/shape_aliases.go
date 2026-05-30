package pptx

//go:generate go run ../../cmd/gen_shape_types shapes shape_aliases.go shape_types.go

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

	// RichShapeFill provides rich fill properties (solid, gradient, pattern, no-fill).
	RichShapeFill = shapes.RichShapeFill
	// RichShapeLine provides rich line properties with full styling control.
	RichShapeLine = shapes.RichShapeLine
	// RichShapeShadow provides rich shadow properties with detailed controls.
	RichShapeShadow = shapes.RichShapeShadow

	// PatternType represents pattern fill types.
	PatternType = shapes.PatternType
	// LineDashStyle represents line dash styles.
	LineDashStyle = shapes.LineDashStyle
	// LineCapStyle represents line cap styles.
	LineCapStyle = shapes.LineCapStyle
	// LineJoinStyle represents line join styles.
	LineJoinStyle = shapes.LineJoinStyle
	// ShadowType represents shadow types.
	ShadowType = shapes.ShadowType
	// ShadowAlignment represents shadow alignment options.
	ShadowAlignment = shapes.ShadowAlignment

	// GroupShape represents a group of shapes that move and transform together.
	GroupShape = shapes.GroupShape
	// Freeform represents a custom-geometry (freeform) shape.
	Freeform = shapes.Freeform
	// FreeformPoint represents a point in a freeform shape path.
	FreeformPoint = shapes.FreeformPoint
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
	ShapeTypeGroup               = shapes.ShapeTypeGroup

	ShapeGradientTypeLinear      = shapes.ShapeGradientTypeLinear
	ShapeGradientTypeRadial      = shapes.ShapeGradientTypeRadial
	ShapeGradientTypeRectangular = shapes.ShapeGradientTypeRectangular
	ShapeGradientTypePath        = shapes.ShapeGradientTypePath

	LineCapFlat   = shapes.LineCapFlat
	LineCapRound  = shapes.LineCapRound
	LineCapSquare = shapes.LineCapSquare

	LineJoinRound = shapes.LineJoinRound
	LineJoinBevel = shapes.LineJoinBevel
	LineJoinMiter = shapes.LineJoinMiter

	TextAnchorTop    = shapes.TextAnchorTop
	TextAnchorMiddle = shapes.TextAnchorMiddle
	TextAnchorBottom = shapes.TextAnchorBottom

	TextWrapNone   = shapes.TextWrapNone
	TextWrapSquare = shapes.TextWrapSquare

	TextAutoFitNone   = shapes.TextAutoFitNone
	TextAutoFitShape  = shapes.TextAutoFitShape
	TextAutoFitNormal = shapes.TextAutoFitNormal
)
