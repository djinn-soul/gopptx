package shape

import common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"

// ParsedShapeProperties contains the normalized properties read from shape XML.
type ParsedShapeProperties struct {
	ID         int
	Name       string
	Type       string
	Text       string
	Runs       []common.TextRun
	Paragraphs []common.ShapeTextParagraph
	TextFrame  *common.TextFrame
	Paragraph  *common.Paragraph
	Fill       *common.ShapeFill
	Line       *common.ShapeLine
	Shadow     *common.ShapeShadow
	Glow       *common.ShapeGlow
	Blur       *common.ShapeBlur
	SoftEdge   *common.ShapeSoftEdge
	Reflection *common.ShapeReflection
	Rotation   *float64
	FlipH      bool
	FlipV      bool
	Connector  *common.ConnectorInfo
	X, Y       int
	W, H       int
	PhIndex    int
	PhType     string

	Adjustments []common.ShapeAdjustment
	Freeform    *common.FreeformGeometry
}
