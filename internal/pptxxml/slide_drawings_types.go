package pptxxml

// ShapeFillSpec describes solid fill properties for a custom shape.
type ShapeFillSpec struct {
	Color        string
	Transparency *float64
}

// ShapeGradientStopSpec describes one gradient stop for a custom shape.
type ShapeGradientStopSpec struct {
	PositionPct  int
	Color        string
	Transparency *float64
}

// ShapeGradientFillSpec describes gradient fill properties for a custom shape.
type ShapeGradientFillSpec struct {
	Type     string
	Stops    []ShapeGradientStopSpec
	AngleDeg *int
}

// ShapeLineSpec describes line properties for a custom shape or connector.
type ShapeLineSpec struct {
	Color string
	Width int64
	Dash  string
	Cap   string
	Join  string
}

// ShapeSpec describes one custom shape rendered as p:sp.
type ShapeSpec struct {
	Type         string
	X            int64
	Y            int64
	CX           int64
	CY           int64
	Fill         *ShapeFillSpec
	GradientFill *ShapeGradientFillSpec
	Line         *ShapeLineSpec
	Text         string
	RotationDeg  *int
	Hyperlink    *HyperlinkSpec // Legacy: mapped to ClickAction
	ClickAction  *HyperlinkSpec
	HoverAction  *HyperlinkSpec
	AltText      string
	IsDecorative bool
	TextFrame    *TextFrameSpec
	Name         string
	Adjustments  []ConnectorAdjustmentSpec
	Effects      *ShapeEffectsSpec
	RichFill     *RichShapeFillSpec
	RichLine     *RichShapeLineSpec
	RichShadow   *RichShapeShadowSpec
}

// ShapeEffectsSpec describes effects for one custom shape.
type ShapeEffectsSpec struct {
	Shadow         bool
	Glow           bool
	SoftEdges      bool
	Reflection     bool
	GlowSpec       *ShapeGlowSpec
	BlurSpec       *ShapeBlurSpec
	SoftEdgeSpec   *ShapeSoftEdgeSpec
	ReflectionSpec *ShapeReflectionSpec
}

// ShapeGlowSpec describes detailed glow effect settings.
type ShapeGlowSpec struct {
	Color     string
	RadiusEmu int
}

// ShapeBlurSpec describes detailed blur effect settings.
type ShapeBlurSpec struct {
	RadiusEmu int
}

// ShapeSoftEdgeSpec describes detailed soft-edge effect settings.
type ShapeSoftEdgeSpec struct {
	RadiusEmu int
}

// ShapeReflectionSpec describes detailed reflection effect settings.
type ShapeReflectionSpec struct {
	BlurEmu     int
	DistanceEmu int
}

// FillType represents the type of shape fill.
type FillType string

const (
	FillTypeSolid    FillType = "solid"
	FillTypeGradient FillType = "gradient"
	FillTypePattern  FillType = "pattern"
	FillTypeNoFill   FillType = "noFill"
)

// SolidFillSpec describes a solid fill.
type SolidFillSpec struct {
	Color        string
	Transparency float64
}

// PatternFillSpec describes a pattern fill.
type PatternFillSpec struct {
	Pattern string
	FgColor string
	BgColor string
}

// RichShapeFillSpec provides a unified spec for all fill types.
type RichShapeFillSpec struct {
	Type     FillType
	Solid    *SolidFillSpec
	Gradient *ShapeGradientFillSpec
	Pattern  *PatternFillSpec
}

// LineDashStyle represents line dash styles.
type LineDashStyle string

const (
	LineDashStyleSolid       LineDashStyle = "solid"
	LineDashStyleDash        LineDashStyle = "dash"
	LineDashStyleDot         LineDashStyle = "dot"
	LineDashStyleDashDot     LineDashStyle = "dashDot"
	LineDashStyleDashDotDot  LineDashStyle = "dashDotDot"
	LineDashStyleLongDash    LineDashStyle = "lgDash"
	LineDashStyleLongDashDot LineDashStyle = "lgDashDot"
)

// LineCapStyle represents line cap styles.
type LineCapStyle string

const (
	LineCapStyleFlat   LineCapStyle = "flat"
	LineCapStyleRound  LineCapStyle = "rnd"
	LineCapStyleSquare LineCapStyle = "sq"
)

// LineJoinStyle represents line join styles.
type LineJoinStyle string

const (
	LineJoinStyleRound LineJoinStyle = "round"
	LineJoinStyleBevel LineJoinStyle = "bevel"
	LineJoinStyleMiter LineJoinStyle = "miter"
)

// RichShapeLineSpec provides detailed control over shape line properties.
type RichShapeLineSpec struct {
	Color        string
	Width        int64
	DashStyle    LineDashStyle
	CapStyle     LineCapStyle
	JoinStyle    LineJoinStyle
	Transparency float64
}

// ShadowType represents shadow types.
type ShadowType string

const (
	ShadowTypeOuter       ShadowType = "outer"
	ShadowTypeInner       ShadowType = "inner"
	ShadowTypePerspective ShadowType = "perspective"
)

// RichShapeShadowSpec provides detailed control over shape shadow effects.
type RichShapeShadowSpec struct {
	Type            ShadowType
	Color           string
	Transparency    float64
	BlurRadius      int
	Distance        int
	Angle           float64
	Alignment       string
	SkewX           float64
	SkewY           float64
	ScaleX          float64
	ScaleY          float64
	RotateWithShape bool
}

// TextFrameSpec describes the text layout within a shape.
type TextFrameSpec struct {
	MarginLeft   int64
	MarginRight  int64
	MarginTop    int64
	MarginBottom int64
	Anchor       string
	Wrap         string
	AutoFit      string
	Orientation  string
	NumCol       int
	Rotation     *int64
}

// ConnectorSpec describes one custom connector rendered as p:cxnSp.
type ConnectorSpec struct {
	Type            string
	StartX, StartY  int64
	EndX, EndY      int64
	Line            ShapeLineSpec
	StartArrow      string
	StartArrowWidth string
	StartArrowLen   string
	EndArrow        string
	EndArrowWidth   string
	EndArrowLen     string
	StartShapeIndex int
	StartSiteIndex  *int
	EndShapeIndex   int
	EndSiteIndex    *int
	Label           string
	ClickAction     *HyperlinkSpec
	HoverAction     *HyperlinkSpec
	AltText         string
	IsDecorative    bool
	Adjustments     []ConnectorAdjustmentSpec
}

// ConnectorAdjustmentSpec describes one connector adjustment entry (a:gd).
type ConnectorAdjustmentSpec struct {
	Name    string
	Formula string
}
