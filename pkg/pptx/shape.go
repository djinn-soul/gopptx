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

func normalizeShapeType(shapeType string) string {
	switch strings.ToLower(strings.TrimSpace(shapeType)) {
	case ShapeTypeRectangle, "rectangle":
		return ShapeTypeRectangle
	case ShapeTypeRoundedRectangle, "roundedrectangle", "rounded-rectangle", "rounded_rectangle":
		return ShapeTypeRoundedRectangle
	case ShapeTypeEllipse, "circle":
		return ShapeTypeEllipse
	case ShapeTypeTriangle:
		return ShapeTypeTriangle
	case ShapeTypeRightTriangle, "righttriangle", "right-triangle", "right_triangle":
		return ShapeTypeRightTriangle
	case ShapeTypeDiamond:
		return ShapeTypeDiamond
	case ShapeTypePentagon:
		return ShapeTypePentagon
	case ShapeTypeHexagon:
		return ShapeTypeHexagon
	case ShapeTypeParallelogram:
		return ShapeTypeParallelogram
	case ShapeTypeFlowChartProcess, "flowchartprocess", "flowchart-process", "flowchart_process":
		return ShapeTypeFlowChartProcess
	case ShapeTypeFlowChartDecision, "flowchartdecision", "flowchart-decision", "flowchart_decision":
		return ShapeTypeFlowChartDecision
	case ShapeTypeFlowChartTerminator, "flowchartterminator", "flowchart-terminator", "flowchart_terminator":
		return ShapeTypeFlowChartTerminator
	case ShapeTypeRightArrow, "rightarrow", "right-arrow", "right_arrow":
		return ShapeTypeRightArrow
	case ShapeTypeLeftArrow, "leftarrow", "left-arrow", "left_arrow":
		return ShapeTypeLeftArrow
	case ShapeTypeUpArrow, "uparrow", "up-arrow", "up_arrow":
		return ShapeTypeUpArrow
	case ShapeTypeDownArrow, "downarrow", "down-arrow", "down_arrow":
		return ShapeTypeDownArrow
	case ShapeTypeCloud:
		return ShapeTypeCloud
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
		ShapeTypeCloud:
		return true
	default:
		return false
	}
}
