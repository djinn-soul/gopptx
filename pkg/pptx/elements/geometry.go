package elements

import (
	"encoding/base64"
	"fmt"
	"strings"
)

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

// ShapeFill configures solid fill properties for one shape.
type ShapeFill struct {
	Color           string
	TransparencyPct *int
}

// NewShapeFill creates a solid fill using a 6-digit RGB color.
func NewShapeFill(color string) ShapeFill {
	return ShapeFill{Color: NormalizeHexColor(color)}
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
		Color: NormalizeHexColor(color),
		Width: width,
		Dash:  LineDashSolid,
	}
}

// WithDash sets line dash style.
func (l ShapeLine) WithDash(dash string) ShapeLine {
	l.Dash = NormalizeDrawingLineDash(dash)
	return l
}

// ShapeGradientStop configures one gradient stop for a shape fill.
type ShapeGradientStop struct {
	PositionPct     int
	Color           string
	TransparencyPct *int
}

// NewShapeGradientStop creates a gradient stop at one position in [0,100].
func NewShapeGradientStop(positionPct int, color string) ShapeGradientStop {
	return ShapeGradientStop{
		PositionPct: positionPct,
		Color:       NormalizeHexColor(color),
	}
}

// WithTransparency sets stop transparency percentage in the range [0,100].
func (s ShapeGradientStop) WithTransparency(percent int) ShapeGradientStop {
	value := percent
	s.TransparencyPct = &value
	return s
}

// ShapeGradientFill configures gradient fill properties for one shape.
type ShapeGradientFill struct {
	Type     string
	Stops    []ShapeGradientStop
	AngleDeg *int
}

// NewShapeGradientFill creates one gradient fill.
func NewShapeGradientFill(gradientType string, stops []ShapeGradientStop) ShapeGradientFill {
	copiedStops := make([]ShapeGradientStop, len(stops))
	copy(copiedStops, stops)
	return ShapeGradientFill{
		Type:  NormalizeShapeGradientType(gradientType),
		Stops: copiedStops,
	}
}

// WithLinearAngle sets the linear gradient angle in degrees.
func (f ShapeGradientFill) WithLinearAngle(degrees int) ShapeGradientFill {
	value := degrees
	f.AngleDeg = &value
	return f
}

// Shape is one auto shape.
type Shape struct {
	Type         string
	X            int64
	Y            int64
	CX           int64
	CY           int64
	Fill         *ShapeFill
	Line         *ShapeLine
	GradientFill *ShapeGradientFill
	Text         string
	RotationDeg  *int
	Hyperlink    *Hyperlink
	AltText      string
	IsDecorative bool
}

// NewShape creates one shape.
func NewShape(shapeType string, x, y, cx, cy int64) Shape {
	return Shape{
		Type: NormalizeShapeType(shapeType),
		X:    x,
		Y:    y,
		CX:   cx,
		CY:   cy,
	}
}

// WithFill applies solid fill to a shape.
func (s Shape) WithFill(fill ShapeFill) Shape {
	s.Fill = &fill
	s.GradientFill = nil
	return s
}

// WithLine applies a line style to a shape.
func (s Shape) WithLine(line ShapeLine) Shape {
	s.Line = &line
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

// WithGradientFill applies gradient fill to a shape.
func (s Shape) WithGradientFill(fill ShapeGradientFill) Shape {
	s.GradientFill = &fill
	s.Fill = nil
	return s
}

// ImageCrop defines cropping details for an image.
type ImageCrop struct {
	Left   float64
	Right  float64
	Top    float64
	Bottom float64
}

// Image describes one image placement.
type Image struct {
	Path         string
	SourceURL    string
	Data         []byte
	Format       string
	X            int64
	Y            int64
	CX           int64
	CY           int64
	Rotation     float64
	Crop         ImageCrop
	FlipH        bool
	FlipV        bool
	Shadow       bool
	Reflection   bool
	AltText      string
	IsDecorative bool
	Placeholder  *Placeholder
}

// NewImage creates an image placement.
func NewImage(path string, x, y, cx, cy int64) Image {
	return Image{Path: path, X: x, Y: y, CX: cx, CY: cy}
}

// NewImageFromBytes creates an image placement from raw bytes.
func NewImageFromBytes(data []byte, format string, x, y, cx, cy int64) Image {
	return Image{Data: data, Format: format, X: x, Y: y, CX: cx, CY: cy}
}

// NewImageFromBase64 creates an image placement from a base64 string.
func NewImageFromBase64(b64 string, format string, x, y, cx, cy int64) (Image, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return Image{}, fmt.Errorf("invalid base64 image data: %w", err)
	}
	return NewImageFromBytes(data, format, x, y, cx, cy), nil
}

// NewImageFromURL creates an image placement from a URL.
func NewImageFromURL(url string, x, y, cx, cy int64) Image {
	return Image{SourceURL: url, X: x, Y: y, CX: cx, CY: cy}
}

// WithShadow adds an outer shadow effect to the image.
func (img Image) WithShadow(enabled bool) Image {
	img.Shadow = enabled
	return img
}

// WithReflection adds a reflection effect to the image.
func (img Image) WithReflection(enabled bool) Image {
	img.Reflection = enabled
	return img
}

// WithRotation adds rotation (degrees) to the image.
func (img Image) WithRotation(degrees float64) Image {
	img.Rotation = degrees
	return img
}

// WithCrop adds cropping to the image.
func (img Image) WithCrop(left, right, top, bottom float64) Image {
	img.Crop = ImageCrop{
		Left:   left,
		Right:  right,
		Top:    top,
		Bottom: bottom,
	}
	return img
}

// WithFlip adds horizontal/vertical flip.
func (img Image) WithFlip(horizontal, vertical bool) Image {
	img.FlipH = horizontal
	img.FlipV = vertical
	return img
}

// WithAltText sets the alternative text for accessibility.
func (img Image) WithAltText(text string) Image {
	img.AltText = text
	return img
}

// WithDecorative marks the image as decorative (ignored by screen readers).
func (img Image) WithDecorative(enabled bool) Image {
	img.IsDecorative = enabled
	return img
}

// Connector is one connector.
type Connector struct {
	Type            string
	StartX          int64
	StartY          int64
	EndX            int64
	EndY            int64
	Line            ShapeLine
	StartArrow      string
	EndArrow        string
	ArrowSize       string
	StartShapeIndex int
	StartSite       string
	EndShapeIndex   int
	EndSite         string
	Label           string
	AltText         string
	IsDecorative    bool
	Placeholder     *Placeholder
}

// NewConnector creates a connector.
func NewConnector(connectorType string, startX, startY, endX, endY int64) Connector {
	return Connector{
		Type:       NormalizeConnectorType(connectorType),
		StartX:     startX,
		StartY:     startY,
		EndX:       endX,
		EndY:       endY,
		Line:       NewShapeLine("000000", 12700),
		StartArrow: ArrowTypeNone,
		EndArrow:   ArrowTypeNone,
		ArrowSize:  ArrowSizeMedium,
	}
}

func NewStraightConnector(startX, startY, endX, endY int64) Connector {
	return NewConnector(ConnectorTypeStraight, startX, startY, endX, endY)
}

func NewElbowConnector(startX, startY, endX, endY int64) Connector {
	return NewConnector(ConnectorTypeElbow, startX, startY, endX, endY)
}

func NewCurvedConnector(startX, startY, endX, endY int64) Connector {
	return NewConnector(ConnectorTypeCurved, startX, startY, endX, endY)
}

// WithLine sets connector line color and width.
func (c Connector) WithLine(line ShapeLine) Connector {
	c.Line = line
	return c
}

// WithDash sets connector dash style.
func (c Connector) WithDash(dash string) Connector {
	c.Line.Dash = NormalizeDrawingLineDash(dash)
	return c
}

// WithArrows sets start and end arrowhead types.
func (c Connector) WithArrows(startArrow string, endArrow string) Connector {
	c.StartArrow = NormalizeArrowType(startArrow)
	c.EndArrow = NormalizeArrowType(endArrow)
	return c
}

// WithArrowSize sets arrowhead size for both ends.
func (c Connector) WithArrowSize(size string) Connector {
	c.ArrowSize = NormalizeArrowSize(size)
	return c
}

// ConnectStart anchors the connector start to the indexed custom shape (1-based).
func (c Connector) ConnectStart(shapeIndex int, site string) Connector {
	c.StartShapeIndex = shapeIndex
	c.StartSite = NormalizeConnectionSite(site)
	return c
}

// ConnectEnd anchors the connector end to the indexed custom shape (1-based).
func (c Connector) ConnectEnd(shapeIndex int, site string) Connector {
	c.EndShapeIndex = shapeIndex
	c.EndSite = NormalizeConnectionSite(site)
	return c
}

// WithLabel sets connector label text.
func (c Connector) WithLabel(label string) Connector {
	c.Label = label
	return c
}

// WithAltText sets the alternative text for accessibility.
func (c Connector) WithAltText(text string) Connector {
	c.AltText = text
	return c
}

// WithDecorative marks the connector as decorative (ignored by screen readers).
func (c Connector) WithDecorative(enabled bool) Connector {
	c.IsDecorative = enabled
	return c
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

// ConnectStartAuto anchors the connector start to a shape and auto-selects the site.
func (c Connector) ConnectStartAuto(shapeIndex int) Connector {
	c.StartShapeIndex = shapeIndex
	c.StartSite = ""
	return c
}

// ConnectEndAuto anchors the connector end to a shape and auto-selects the site.
func (c Connector) ConnectEndAuto(shapeIndex int) Connector {
	c.EndShapeIndex = shapeIndex
	c.EndSite = ""
	return c
}

func ResolveConnectorSiteIndices(connector Connector, shapes []Shape) (*int, *int) {
	startSite := connector.StartSite
	if strings.TrimSpace(startSite) == "" && connector.StartShapeIndex > 0 {
		targetX, targetY := connector.EndX, connector.EndY
		if centerX, centerY, ok := shapeCenterForIndex(shapes, connector.EndShapeIndex); ok {
			targetX, targetY = centerX, centerY
		}
		if shape, ok := shapeForIndex(shapes, connector.StartShapeIndex); ok {
			startSite = autoConnectionSite(shape, targetX, targetY)
		}
	}

	endSite := connector.EndSite
	if strings.TrimSpace(endSite) == "" && connector.EndShapeIndex > 0 {
		targetX, targetY := connector.StartX, connector.StartY
		if centerX, centerY, ok := shapeCenterForIndex(shapes, connector.StartShapeIndex); ok {
			targetX, targetY = centerX, centerY
		}
		if shape, ok := shapeForIndex(shapes, connector.EndShapeIndex); ok {
			endSite = autoConnectionSite(shape, targetX, targetY)
		}
	}

	return SiteIndexPointer(startSite), SiteIndexPointer(endSite)
}

func shapeForIndex(shapes []Shape, shapeIndex int) (Shape, bool) {
	if shapeIndex <= 0 || shapeIndex > len(shapes) {
		return Shape{}, false
	}
	return shapes[shapeIndex-1], true
}

func shapeCenterForIndex(shapes []Shape, shapeIndex int) (int64, int64, bool) {
	shape, ok := shapeForIndex(shapes, shapeIndex)
	if !ok {
		return 0, 0, false
	}
	return shape.X + shape.CX/2, shape.Y + shape.CY/2, true
}

func autoConnectionSite(shape Shape, targetX int64, targetY int64) string {
	candidates := shapeConnectionSiteCandidates(shape)
	bestSite := ConnectionSiteCenter
	var bestDistance int64
	first := true
	for _, candidate := range candidates {
		dx := candidate.x - targetX
		dy := candidate.y - targetY
		distance := dx*dx + dy*dy
		if first || distance < bestDistance {
			first = false
			bestDistance = distance
			bestSite = candidate.site
		}
	}
	return bestSite
}

type connectionSiteCandidate struct {
	site string
	x    int64
	y    int64
}

func shapeConnectionSiteCandidates(shape Shape) [9]connectionSiteCandidate {
	cx := shape.X + shape.CX/2
	cy := shape.Y + shape.CY/2
	right := shape.X + shape.CX
	bottom := shape.Y + shape.CY

	return [9]connectionSiteCandidate{
		{site: ConnectionSiteTop, x: cx, y: shape.Y},
		{site: ConnectionSiteRight, x: right, y: cy},
		{site: ConnectionSiteBottom, x: cx, y: bottom},
		{site: ConnectionSiteLeft, x: shape.X, y: cy},
		{site: ConnectionSiteTopLeft, x: shape.X, y: shape.Y},
		{site: ConnectionSiteTopRight, x: right, y: shape.Y},
		{site: ConnectionSiteBottomRight, x: right, y: bottom},
		{site: ConnectionSiteBottomLeft, x: shape.X, y: bottom},
		{site: ConnectionSiteCenter, x: cx, y: cy},
	}
}

func SiteIndexPointer(site string) *int {
	if idx, ok := ConnectionSiteIndex(site); ok {
		value := idx
		return &value
	}
	return nil
}

func RunsToPlainText(runs []TextRun) string {
	var b strings.Builder
	for _, run := range runs {
		b.WriteString(run.Text)
	}
	return b.String()
}

// ShapeDefinition allows external shape builders to plug into slide composition.
type ShapeDefinition interface {
	ToShape() Shape
}

// ToShape returns the shape itself and satisfies ShapeDefinition.
func (s Shape) ToShape() Shape {
	return s
}

// Validate checks the image for common constraints.
func (img Image) Validate(slideIndex, imageIndex int) error {
	if img.Path == "" && len(img.Data) == 0 && img.SourceURL == "" {
		return fmt.Errorf("slide %d image %d has no source (Path, Data, or SourceURL)", slideIndex, imageIndex)
	}
	if len(img.Data) > 0 && img.Format == "" {
		return fmt.Errorf("slide %d image %d has Data but no Format", slideIndex, imageIndex)
	}
	if img.X < 0 || img.Y < 0 {
		return fmt.Errorf("slide %d image %d position cannot be negative", slideIndex, imageIndex)
	}
	if img.CX <= 0 || img.CY <= 0 {
		return fmt.Errorf("slide %d image %d size must be > 0", slideIndex, imageIndex)
	}
	return nil
}

// Validate checks shape properties for OOXML compliance and library constraints.
func (shape Shape) Validate(slideIndex int, shapeIndex int) error {
	if !IsShapeType(shape.Type) {
		return fmt.Errorf("slide %d shape %d type %q is not supported", slideIndex, shapeIndex, shape.Type)
	}
	if shape.X < 0 || shape.Y < 0 {
		return fmt.Errorf("slide %d shape %d position cannot be negative", slideIndex, shapeIndex)
	}
	if shape.CX <= 0 || shape.CY <= 0 {
		return fmt.Errorf("slide %d shape %d size must be > 0", slideIndex, shapeIndex)
	}
	if shape.Fill != nil && shape.GradientFill != nil {
		return fmt.Errorf("slide %d shape %d cannot set both solid and gradient fill", slideIndex, shapeIndex)
	}
	if shape.Fill != nil {
		if !IsHexColor(shape.Fill.Color) {
			return fmt.Errorf("slide %d shape %d fill color must be 6-digit RGB hex", slideIndex, shapeIndex)
		}
		if shape.Fill.TransparencyPct != nil {
			if *shape.Fill.TransparencyPct < 0 || *shape.Fill.TransparencyPct > 100 {
				return fmt.Errorf("slide %d shape %d fill transparency must be in [0,100]", slideIndex, shapeIndex)
			}
		}
	}
	if shape.GradientFill != nil {
		if err := validateShapeGradientFill(*shape.GradientFill, slideIndex, shapeIndex); err != nil {
			return err
		}
	}
	if shape.Line != nil {
		if !IsHexColor(shape.Line.Color) {
			return fmt.Errorf("slide %d shape %d line color must be 6-digit RGB hex", slideIndex, shapeIndex)
		}
		if shape.Line.Width <= 0 {
			return fmt.Errorf("slide %d shape %d line width must be > 0", slideIndex, shapeIndex)
		}
		if !IsDrawingLineDash(shape.Line.Dash) {
			return fmt.Errorf(
				"slide %d shape %d line dash must be one of solid|dash|dot|dashDot|lgDash|lgDashDot|lgDashDotDot",
				slideIndex,
				shapeIndex,
			)
		}
	}
	if shape.RotationDeg != nil {
		if *shape.RotationDeg < -360 || *shape.RotationDeg > 360 {
			return fmt.Errorf("slide %d shape %d rotation must be in [-360,360]", slideIndex, shapeIndex)
		}
	}
	// TODO: Move validateHyperlink to elements
	return nil
}

func validateShapeGradientFill(fill ShapeGradientFill, slideIndex int, shapeIndex int) error {
	if !IsShapeGradientType(fill.Type) {
		return fmt.Errorf("slide %d shape %d gradient type %q is not supported", slideIndex, shapeIndex, fill.Type)
	}
	if len(fill.Stops) < 2 {
		return fmt.Errorf("slide %d shape %d gradient must contain at least 2 stops", slideIndex, shapeIndex)
	}

	lastPosition := -1
	for stopIndex, stop := range fill.Stops {
		if stop.PositionPct < 0 || stop.PositionPct > 100 {
			return fmt.Errorf(
				"slide %d shape %d gradient stop %d position must be in [0,100]",
				slideIndex,
				shapeIndex,
				stopIndex+1,
			)
		}
		if stop.PositionPct <= lastPosition {
			return fmt.Errorf(
				"slide %d shape %d gradient stop positions must be strictly increasing",
				slideIndex,
				shapeIndex,
			)
		}
		lastPosition = stop.PositionPct
		if !IsHexColor(stop.Color) {
			return fmt.Errorf(
				"slide %d shape %d gradient stop %d color must be 6-digit RGB hex",
				slideIndex,
				shapeIndex,
				stopIndex+1,
			)
		}
		if stop.TransparencyPct != nil && (*stop.TransparencyPct < 0 || *stop.TransparencyPct > 100) {
			return fmt.Errorf(
				"slide %d shape %d gradient stop %d transparency must be in [0,100]",
				slideIndex,
				shapeIndex,
				stopIndex+1,
			)
		}
	}

	if fill.AngleDeg != nil {
		if NormalizeShapeGradientType(fill.Type) != ShapeGradientTypeLinear {
			return fmt.Errorf(
				"slide %d shape %d gradient angle is only supported for linear gradients",
				slideIndex,
				shapeIndex,
			)
		}
		if *fill.AngleDeg < -360 || *fill.AngleDeg > 360 {
			return fmt.Errorf(
				"slide %d shape %d gradient angle must be in [-360,360]",
				slideIndex,
				shapeIndex,
			)
		}
	}

	return nil
}

// Validate checks connector properties and anchor references.
func (connector Connector) Validate(shapeCount int, slideIndex int, connectorIndex int) error {
	if !IsConnectorType(connector.Type) {
		return fmt.Errorf("slide %d connector %d type %q is not supported", slideIndex, connectorIndex, connector.Type)
	}
	if connector.StartX < 0 || connector.StartY < 0 || connector.EndX < 0 || connector.EndY < 0 {
		return fmt.Errorf("slide %d connector %d coordinates cannot be negative", slideIndex, connectorIndex)
	}
	if connector.StartX == connector.EndX && connector.StartY == connector.EndY {
		return fmt.Errorf("slide %d connector %d must have distinct start and end points", slideIndex, connectorIndex)
	}
	if !IsHexColor(connector.Line.Color) {
		return fmt.Errorf("slide %d connector %d line color must be 6-digit RGB hex", slideIndex, connectorIndex)
	}
	if connector.Line.Width <= 0 {
		return fmt.Errorf("slide %d connector %d line width must be > 0", slideIndex, connectorIndex)
	}
	if !IsDrawingLineDash(connector.Line.Dash) {
		return fmt.Errorf(
			"slide %d connector %d line dash must be one of solid|dash|dot|dashDot|lgDash|lgDashDot|lgDashDotDot",
			slideIndex,
			connectorIndex,
		)
	}
	if !IsArrowType(connector.StartArrow) {
		return fmt.Errorf("slide %d connector %d start arrow %q is invalid", slideIndex, connectorIndex, connector.StartArrow)
	}
	if !IsArrowType(connector.EndArrow) {
		return fmt.Errorf("slide %d connector %d end arrow %q is invalid", slideIndex, connectorIndex, connector.EndArrow)
	}
	if !IsArrowSize(connector.ArrowSize) {
		return fmt.Errorf("slide %d connector %d arrow size %q is invalid", slideIndex, connectorIndex, connector.ArrowSize)
	}
	if err := validateConnectorAnchor("start", connector.StartShapeIndex, connector.StartSite, shapeCount, slideIndex, connectorIndex); err != nil {
		return err
	}
	if err := validateConnectorAnchor("end", connector.EndShapeIndex, connector.EndSite, shapeCount, slideIndex, connectorIndex); err != nil {
		return err
	}
	return nil
}

func validateConnectorAnchor(
	side string,
	shapeIndex int,
	site string,
	shapeCount int,
	slideIndex int,
	connectorIndex int,
) error {
	if shapeIndex == 0 {
		if strings.TrimSpace(site) != "" {
			return fmt.Errorf(
				"slide %d connector %d %s site requires a %s shape index",
				slideIndex,
				connectorIndex,
				side,
				side,
			)
		}
		return nil
	}
	if shapeIndex < 0 || shapeIndex > shapeCount {
		return fmt.Errorf(
			"slide %d connector %d %s shape index %d is out of range [1,%d]",
			slideIndex,
			connectorIndex,
			side,
			shapeIndex,
			shapeCount,
		)
	}
	if strings.TrimSpace(site) == "" {
		return nil
	}
	if _, ok := ConnectionSiteIndex(site); !ok {
		return fmt.Errorf(
			"slide %d connector %d %s site %q is invalid",
			slideIndex,
			connectorIndex,
			side,
			site,
		)
	}
	return nil
}
