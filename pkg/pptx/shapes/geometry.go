package shapes

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

	// LineDashSolid emits a solid line.
	LineDashSolid = "solid"
	// LineDashDash emits a dashed line.
	LineDashDash = "dash"
	// LineDashDot emits a dotted line.
	LineDashDot = "dot"
	// LineDashDashDot emits a dash-dot line.
	LineDashDashDot = "dashDot"
	// LineDashDashDotDot emits a dash-dot-dot line.
	LineDashDashDotDot = "lgDashDotDot"
	// LineDashLongDash emits a long-dash line.
	LineDashLongDash = "lgDash"
	// LineDashLongDashDot emits a long-dash-dot line.
	LineDashLongDashDot = "lgDashDot"

	// ShapeGradientTypeLinear renders a linear gradient.
	ShapeGradientTypeLinear = "linear"
	// ShapeGradientTypeRadial renders a radial gradient.
	ShapeGradientTypeRadial = "radial"
	// ShapeGradientTypeRectangular renders a rectangular gradient.
	ShapeGradientTypeRectangular = "rectangular"
	// ShapeGradientTypePath renders a shape-path gradient.
	ShapeGradientTypePath = "path"

	// ConnectorTypeStraight renders a straight connector.
	ConnectorTypeStraight = "straightConnector1"
	// ConnectorTypeElbow renders an elbow connector.
	ConnectorTypeElbow = "bentConnector3"
	// ConnectorTypeCurved renders a curved connector.
	ConnectorTypeCurved = "curvedConnector3"

	// ArrowTypeNone renders no arrowhead.
	ArrowTypeNone = "none"
	// ArrowTypeTriangle renders a triangle arrowhead.
	ArrowTypeTriangle = "triangle"
	// ArrowTypeStealth renders a stealth arrowhead.
	ArrowTypeStealth = "stealth"
	// ArrowTypeDiamond renders a diamond arrowhead.
	ArrowTypeDiamond = "diamond"
	// ArrowTypeOval renders an oval arrowhead.
	ArrowTypeOval = "oval"
	// ArrowTypeOpen renders an open arrowhead.
	ArrowTypeOpen = "arrow"

	// ArrowSizeSmall renders small arrowheads.
	ArrowSizeSmall = "sm"
	// ArrowSizeMedium renders medium arrowheads.
	ArrowSizeMedium = "med"
	// ArrowSizeLarge renders large arrowheads.
	ArrowSizeLarge = "lg"

	// ConnectionSiteTop anchors on top-center.
	ConnectionSiteTop = "top"
	// ConnectionSiteRight anchors on right-center.
	ConnectionSiteRight = "right"
	// ConnectionSiteBottom anchors on bottom-center.
	ConnectionSiteBottom = "bottom"
	// ConnectionSiteLeft anchors on left-center.
	ConnectionSiteLeft = "left"
	// ConnectionSiteTopLeft anchors on top-left.
	ConnectionSiteTopLeft = "topLeft"
	// ConnectionSiteTopRight anchors on top-right.
	ConnectionSiteTopRight = "topRight"
	// ConnectionSiteBottomRight anchors on bottom-right.
	ConnectionSiteBottomRight = "bottomRight"
	// ConnectionSiteBottomLeft anchors on bottom-left.
	ConnectionSiteBottomLeft = "bottomLeft"
	// ConnectionSiteCenter anchors on center.
	ConnectionSiteCenter = "center"
)

func NormalizeShapeType(shapeType string) string {
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

func IsShapeType(shapeType string) bool {
	switch NormalizeShapeType(shapeType) {
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

func NormalizeDrawingLineDash(dash string) string {
	switch strings.ToLower(strings.TrimSpace(dash)) {
	case "", "solid":
		return LineDashSolid
	case "dash":
		return LineDashDash
	case "dot":
		return LineDashDot
	case "dashdot", "dash-dot", "dash_dot":
		return LineDashDashDot
	case "dashdotdot", "dash-dot-dot", "dash_dot_dot", "lgdashdotdot", "lg-dash-dot-dot":
		return LineDashDashDotDot
	case "lgdash", "lg-dash", "longdash", "long-dash", "long_dash":
		return LineDashLongDash
	case "lgdashdot", "lg-dash-dot", "longdashdot", "long-dash-dot", "long_dash_dot":
		return LineDashLongDashDot
	default:
		return strings.TrimSpace(dash)
	}
}

func IsDrawingLineDash(dash string) bool {
	switch NormalizeDrawingLineDash(dash) {
	case LineDashSolid, LineDashDash, LineDashDot, LineDashDashDot, LineDashDashDotDot, LineDashLongDash, LineDashLongDashDot:
		return true
	default:
		return false
	}
}

func NormalizeShapeGradientType(gradientType string) string {
	switch strings.ToLower(strings.TrimSpace(gradientType)) {
	case ShapeGradientTypeLinear:
		return ShapeGradientTypeLinear
	case ShapeGradientTypeRadial, "radial-gradient", "radial_gradient":
		return ShapeGradientTypeRadial
	case ShapeGradientTypeRectangular, "rectangular-gradient", "rectangular_gradient", "rect":
		return ShapeGradientTypeRectangular
	case ShapeGradientTypePath, "path-gradient", "path_gradient":
		return ShapeGradientTypePath
	default:
		return strings.TrimSpace(gradientType)
	}
}

func IsShapeGradientType(gradientType string) bool {
	switch NormalizeShapeGradientType(gradientType) {
	case ShapeGradientTypeLinear, ShapeGradientTypeRadial, ShapeGradientTypeRectangular, ShapeGradientTypePath:
		return true
	default:
		return false
	}
}

func NormalizeConnectorType(connectorType string) string {
	t := strings.ToLower(strings.TrimSpace(connectorType))
	switch t {
	case strings.ToLower(ConnectorTypeStraight), "straight", "s":
		return ConnectorTypeStraight
	case strings.ToLower(ConnectorTypeElbow), "elbow", "bent", "e":
		return ConnectorTypeElbow
	case strings.ToLower(ConnectorTypeCurved), "curved", "curve", "c":
		return ConnectorTypeCurved
	default:
		return strings.TrimSpace(connectorType)
	}
}

func IsConnectorType(connectorType string) bool {
	switch NormalizeConnectorType(connectorType) {
	case ConnectorTypeStraight, ConnectorTypeElbow, ConnectorTypeCurved:
		return true
	default:
		return false
	}
}

func NormalizeArrowType(arrowType string) string {
	t := strings.ToLower(strings.TrimSpace(arrowType))
	switch t {
	case strings.ToLower(ArrowTypeNone), "", "n":
		return ArrowTypeNone
	case strings.ToLower(ArrowTypeTriangle), "t":
		return ArrowTypeTriangle
	case strings.ToLower(ArrowTypeStealth), "s":
		return ArrowTypeStealth
	case strings.ToLower(ArrowTypeDiamond), "d":
		return ArrowTypeDiamond
	case strings.ToLower(ArrowTypeOval), "o":
		return ArrowTypeOval
	case strings.ToLower(ArrowTypeOpen), "open", "a":
		return ArrowTypeOpen
	default:
		return strings.TrimSpace(arrowType)
	}
}

func IsArrowType(arrowType string) bool {
	switch NormalizeArrowType(arrowType) {
	case ArrowTypeNone, ArrowTypeTriangle, ArrowTypeStealth, ArrowTypeDiamond, ArrowTypeOval, ArrowTypeOpen:
		return true
	default:
		return false
	}
}

func NormalizeArrowSize(size string) string {
	t := strings.ToLower(strings.TrimSpace(size))
	switch t {
	case strings.ToLower(ArrowSizeMedium), "", "medium", "m":
		return ArrowSizeMedium
	case strings.ToLower(ArrowSizeSmall), "small", "s":
		return ArrowSizeSmall
	case strings.ToLower(ArrowSizeLarge), "large", "l":
		return ArrowSizeLarge
	default:
		return strings.TrimSpace(size)
	}
}

func IsArrowSize(size string) bool {
	switch NormalizeArrowSize(size) {
	case ArrowSizeSmall, ArrowSizeMedium, ArrowSizeLarge:
		return true
	default:
		return false
	}
}

func NormalizeConnectionSite(site string) string {
	t := strings.ToLower(strings.TrimSpace(site))
	switch t {
	case strings.ToLower(ConnectionSiteTop), "t":
		return ConnectionSiteTop
	case strings.ToLower(ConnectionSiteRight), "r":
		return ConnectionSiteRight
	case strings.ToLower(ConnectionSiteBottom), "b":
		return ConnectionSiteBottom
	case strings.ToLower(ConnectionSiteLeft), "l":
		return ConnectionSiteLeft
	case "topleft", "top-left", "top_left", "tl":
		return ConnectionSiteTopLeft
	case "topright", "top-right", "top_right", "tr":
		return ConnectionSiteTopRight
	case "bottomright", "bottom-right", "bottom_right", "br":
		return ConnectionSiteBottomRight
	case "bottomleft", "bottom-left", "bottom_left", "bl":
		return ConnectionSiteBottomLeft
	case strings.ToLower(ConnectionSiteCenter), "ctr", "c":
		return ConnectionSiteCenter
	default:
		return strings.TrimSpace(site)
	}
}

func ConnectionSiteIndex(site string) (int, bool) {
	switch NormalizeConnectionSite(site) {
	case ConnectionSiteTop:
		return 0, true
	case ConnectionSiteRight:
		return 1, true
	case ConnectionSiteBottom:
		return 2, true
	case ConnectionSiteLeft:
		return 3, true
	case ConnectionSiteTopLeft:
		return 4, true
	case ConnectionSiteTopRight:
		return 5, true
	case ConnectionSiteBottomRight:
		return 6, true
	case ConnectionSiteBottomLeft:
		return 7, true
	case ConnectionSiteCenter:
		return 8, true
	default:
		return 0, false
	}
}
