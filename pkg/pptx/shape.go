package pptx

import "strings"

const (
	// ShapeTypeRectangle renders a rectangle shape.
	ShapeTypeRectangle = "rect"
	// ShapeTypeRoundedRectangle renders a rounded rectangle.
	ShapeTypeRoundedRectangle = "roundRect"
	// ShapeTypeEllipse renders an ellipse shape.
	ShapeTypeEllipse = "ellipse"
	// ShapeTypeTriangle renders a triangle shape.
	ShapeTypeTriangle = "triangle"
	// ShapeTypeRightTriangle renders a right triangle shape.
	ShapeTypeRightTriangle = "rtTriangle"
	// ShapeTypeDiamond renders a diamond shape.
	ShapeTypeDiamond = "diamond"
	// ShapeTypePentagon renders a pentagon shape.
	ShapeTypePentagon = "pentagon"
	// ShapeTypeHexagon renders a hexagon shape.
	ShapeTypeHexagon = "hexagon"
	// ShapeTypeParallelogram renders a parallelogram shape.
	ShapeTypeParallelogram = "parallelogram"
	// ShapeTypeFlowChartProcess renders a flowchart process shape.
	ShapeTypeFlowChartProcess = "flowChartProcess"
	// ShapeTypeFlowChartDecision renders a flowchart decision shape.
	ShapeTypeFlowChartDecision = "flowChartDecision"
	// ShapeTypeFlowChartTerminator renders a flowchart terminator shape.
	ShapeTypeFlowChartTerminator = "flowChartTerminator"
	// ShapeTypeRightArrow renders a right arrow shape.
	ShapeTypeRightArrow = "rightArrow"
	// ShapeTypeLeftArrow renders a left arrow shape.
	ShapeTypeLeftArrow = "leftArrow"
	// ShapeTypeUpArrow renders an up arrow shape.
	ShapeTypeUpArrow = "upArrow"
	// ShapeTypeDownArrow renders a down arrow shape.
	ShapeTypeDownArrow = "downArrow"
	// ShapeTypeCloud renders a cloud shape.
	ShapeTypeCloud = "cloud"
	// ShapeTypeStar5 renders a 5-pointed star.
	ShapeTypeStar5 = "star5"
	// ShapeTypeHeart renders a heart shape.
	ShapeTypeHeart = "heart"
	// ShapeTypeFlowChartDocument renders a flowchart document shape.
	ShapeTypeFlowChartDocument = "flowChartDocument"
	// ShapeTypeFlowChartData renders a flowchart data shape (parallelogram).
	ShapeTypeFlowChartData = "flowChartInputOutput"
)

// ShapeFill configures solid fill properties for one shape.
type ShapeFill struct {
	Color           string
	TransparencyPct *int
}

// NewShapeFill creates a solid fill using a 6-digit RGB color.
func NewShapeFill(color string) ShapeFill {
	return ShapeFill{Color: normalizeHexColor(color)}
}

// WithTransparency sets fill transparency percentage in the range [0,100].
func (f ShapeFill) WithTransparency(percent int) ShapeFill {
	value := percent
	f.TransparencyPct = &value
	return f
}

// ShapeLine configures line style for one shape or connector.
type ShapeLine struct {
	Color string
	Width int64
	Dash  string
}

// NewShapeLine creates a line style with RGB color and EMU width.
func NewShapeLine(color string, width int64) ShapeLine {
	return ShapeLine{
		Color: normalizeHexColor(color),
		Width: width,
		Dash:  LineDashSolid,
	}
}

// WithDash sets line dash style (solid|dash|dot|dashDot|lgDash|lgDashDot|lgDashDotDot).
func (l ShapeLine) WithDash(dash string) ShapeLine {
	l.Dash = normalizeDrawingLineDash(dash)
	return l
}

// Shape is one auto shape rendered as p:sp in slide XML.
type Shape struct {
	Type         string
	X            int64
	Y            int64
	CX           int64
	CY           int64
	Fill         *ShapeFill
	GradientFill *ShapeGradientFill
	Line         *ShapeLine
	Text         string
	RotationDeg  *int
	Hyperlink    *Hyperlink
	AltText      string
	IsDecorative bool
}

// NewShape creates one shape with explicit preset type, position, and size.
func NewShape(shapeType string, x int64, y int64, cx int64, cy int64) Shape {
	return Shape{
		Type: normalizeShapeType(shapeType),
		X:    x,
		Y:    y,
		CX:   cx,
		CY:   cy,
	}
}

// WithFill applies solid fill to a shape.
func (s Shape) WithFill(fill ShapeFill) Shape {
	value := fill
	s.Fill = &value
	s.GradientFill = nil
	return s
}

// WithLine applies a line style to a shape.
func (s Shape) WithLine(line ShapeLine) Shape {
	value := line
	s.Line = &value
	return s
}

// WithText sets text rendered inside the shape.
func (s Shape) WithText(text string) Shape {
	s.Text = text
	return s
}

// WithRotation rotates shape geometry in degrees.
func (s Shape) WithRotation(degrees int) Shape {
	value := degrees
	s.RotationDeg = &value
	return s
}

// WithHyperlink attaches a clickable hyperlink to the shape.
func (s Shape) WithHyperlink(hyperlink Hyperlink) Shape {
	s.Hyperlink = &hyperlink
	return s
}

// WithAltText sets the alternative text for accessibility.
func (s Shape) WithAltText(text string) Shape {
	s.AltText = text
	return s
}

// WithDecorative marks the shape as decorative (ignored by screen readers).
func (s Shape) WithDecorative(enabled bool) Shape {
	s.IsDecorative = enabled
	return s
}

func normalizeShapeType(shapeType string) string {
	t := strings.ToLower(strings.TrimSpace(shapeType))
	switch t {
	case strings.ToLower(ShapeTypeRectangle), "rectangle":
		return ShapeTypeRectangle
	case strings.ToLower(ShapeTypeRoundedRectangle), "roundedrectangle", "rounded-rectangle", "rounded_rectangle":
		return ShapeTypeRoundedRectangle
	case strings.ToLower(ShapeTypeEllipse), "circle":
		return ShapeTypeEllipse
	case strings.ToLower(ShapeTypeTriangle):
		return ShapeTypeTriangle
	case strings.ToLower(ShapeTypeRightTriangle), "righttriangle", "right-triangle", "right_triangle":
		return ShapeTypeRightTriangle
	case strings.ToLower(ShapeTypeDiamond):
		return ShapeTypeDiamond
	case strings.ToLower(ShapeTypePentagon):
		return ShapeTypePentagon
	case strings.ToLower(ShapeTypeHexagon):
		return ShapeTypeHexagon
	case strings.ToLower(ShapeTypeParallelogram):
		return ShapeTypeParallelogram
	case strings.ToLower(ShapeTypeFlowChartProcess), "flowchartprocess", "flowchart-process", "flowchart_process":
		return ShapeTypeFlowChartProcess
	case strings.ToLower(ShapeTypeFlowChartDecision), "flowchartdecision", "flowchart-decision", "flowchart_decision":
		return ShapeTypeFlowChartDecision
	case strings.ToLower(ShapeTypeFlowChartTerminator), "flowchartterminator", "flowchart-terminator", "flowchart_terminator":
		return ShapeTypeFlowChartTerminator
	case strings.ToLower(ShapeTypeRightArrow), "rightarrow", "right-arrow", "right_arrow":
		return ShapeTypeRightArrow
	case strings.ToLower(ShapeTypeLeftArrow), "leftarrow", "left-arrow", "left_arrow":
		return ShapeTypeLeftArrow
	case strings.ToLower(ShapeTypeUpArrow), "uparrow", "up-arrow", "up_arrow":
		return ShapeTypeUpArrow
	case strings.ToLower(ShapeTypeDownArrow), "downarrow", "down-arrow", "down_arrow":
		return ShapeTypeDownArrow
	case strings.ToLower(ShapeTypeCloud):
		return ShapeTypeCloud
	case "star", strings.ToLower(ShapeTypeStar5):
		return ShapeTypeStar5
	case "heart", strings.ToLower(ShapeTypeHeart):
		return ShapeTypeHeart
	case "document", "flowchartdocument", strings.ToLower(ShapeTypeFlowChartDocument):
		return ShapeTypeFlowChartDocument
	case "data", "flowchartdata", "flowchartinputoutput", strings.ToLower(ShapeTypeFlowChartData):
		return ShapeTypeFlowChartData
	default:
		return strings.TrimSpace(shapeType)
	}
}

func isShapeType(shapeType string) bool {
	switch normalizeShapeType(shapeType) {
	case ShapeTypeRectangle,
		ShapeTypeRoundedRectangle,
		ShapeTypeEllipse,
		ShapeTypeTriangle,
		ShapeTypeRightTriangle,
		ShapeTypeDiamond,
		ShapeTypePentagon,
		ShapeTypeHexagon,
		ShapeTypeParallelogram,
		ShapeTypeFlowChartProcess,
		ShapeTypeFlowChartDecision,
		ShapeTypeFlowChartTerminator,
		ShapeTypeRightArrow,
		ShapeTypeLeftArrow,
		ShapeTypeUpArrow,
		ShapeTypeDownArrow,
		ShapeTypeCloud,
		ShapeTypeStar5,
		ShapeTypeHeart,
		ShapeTypeFlowChartDocument,
		ShapeTypeFlowChartData:
		return true
	default:
		return false
	}
}
