package pptx

import "strings"

const (
	// ConnectorTypeStraight renders a straight connector.
	ConnectorTypeStraight = "straightConnector1"
	// ConnectorTypeElbow renders an elbow connector.
	ConnectorTypeElbow = "bentConnector3"
	// ConnectorTypeCurved renders a curved connector.
	ConnectorTypeCurved = "curvedConnector3"
)

const (
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
)

const (
	// ArrowSizeSmall renders small arrowheads.
	ArrowSizeSmall = "sm"
	// ArrowSizeMedium renders medium arrowheads.
	ArrowSizeMedium = "med"
	// ArrowSizeLarge renders large arrowheads.
	ArrowSizeLarge = "lg"
)

const (
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

// Connector is one connector rendered as p:cxnSp in slide XML.
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
}

// NewConnector creates a connector with explicit geometry and type.
func NewConnector(connectorType string, startX int64, startY int64, endX int64, endY int64) Connector {
	return Connector{
		Type:       normalizeConnectorType(connectorType),
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

// NewStraightConnector creates a straight connector.
func NewStraightConnector(startX int64, startY int64, endX int64, endY int64) Connector {
	return NewConnector(ConnectorTypeStraight, startX, startY, endX, endY)
}

// NewElbowConnector creates an elbow connector.
func NewElbowConnector(startX int64, startY int64, endX int64, endY int64) Connector {
	return NewConnector(ConnectorTypeElbow, startX, startY, endX, endY)
}

// NewCurvedConnector creates a curved connector.
func NewCurvedConnector(startX int64, startY int64, endX int64, endY int64) Connector {
	return NewConnector(ConnectorTypeCurved, startX, startY, endX, endY)
}

// WithLine sets connector line color and width.
func (c Connector) WithLine(line ShapeLine) Connector {
	c.Line = line
	return c
}

// WithDash sets connector dash style.
func (c Connector) WithDash(dash string) Connector {
	c.Line = c.Line.WithDash(dash)
	return c
}

// WithArrows sets start and end arrowhead types.
func (c Connector) WithArrows(startArrow string, endArrow string) Connector {
	c.StartArrow = normalizeArrowType(startArrow)
	c.EndArrow = normalizeArrowType(endArrow)
	return c
}

// WithArrowSize sets arrowhead size for both ends.
func (c Connector) WithArrowSize(size string) Connector {
	c.ArrowSize = normalizeArrowSize(size)
	return c
}

// ConnectStart anchors the connector start to the indexed custom shape (1-based).
func (c Connector) ConnectStart(shapeIndex int, site string) Connector {
	c.StartShapeIndex = shapeIndex
	c.StartSite = normalizeConnectionSite(site)
	return c
}

// ConnectEnd anchors the connector end to the indexed custom shape (1-based).
func (c Connector) ConnectEnd(shapeIndex int, site string) Connector {
	c.EndShapeIndex = shapeIndex
	c.EndSite = normalizeConnectionSite(site)
	return c
}

// WithLabel sets connector label text.
func (c Connector) WithLabel(label string) Connector {
	c.Label = label
	return c
}

func normalizeConnectorType(connectorType string) string {
	switch strings.ToLower(strings.TrimSpace(connectorType)) {
	case ConnectorTypeStraight, "straight":
		return ConnectorTypeStraight
	case ConnectorTypeElbow, "elbow", "bent":
		return ConnectorTypeElbow
	case ConnectorTypeCurved, "curved", "curve":
		return ConnectorTypeCurved
	default:
		return strings.TrimSpace(connectorType)
	}
}

func isConnectorType(connectorType string) bool {
	switch normalizeConnectorType(connectorType) {
	case ConnectorTypeStraight, ConnectorTypeElbow, ConnectorTypeCurved:
		return true
	default:
		return false
	}
}

func normalizeArrowType(arrowType string) string {
	switch strings.ToLower(strings.TrimSpace(arrowType)) {
	case "", ArrowTypeNone:
		return ArrowTypeNone
	case ArrowTypeTriangle:
		return ArrowTypeTriangle
	case ArrowTypeStealth:
		return ArrowTypeStealth
	case ArrowTypeDiamond:
		return ArrowTypeDiamond
	case ArrowTypeOval:
		return ArrowTypeOval
	case ArrowTypeOpen, "open":
		return ArrowTypeOpen
	default:
		return strings.TrimSpace(arrowType)
	}
}

func isArrowType(arrowType string) bool {
	switch normalizeArrowType(arrowType) {
	case ArrowTypeNone, ArrowTypeTriangle, ArrowTypeStealth, ArrowTypeDiamond, ArrowTypeOval, ArrowTypeOpen:
		return true
	default:
		return false
	}
}

func normalizeArrowSize(size string) string {
	switch strings.ToLower(strings.TrimSpace(size)) {
	case "", ArrowSizeMedium, "medium":
		return ArrowSizeMedium
	case ArrowSizeSmall, "small":
		return ArrowSizeSmall
	case ArrowSizeLarge, "large":
		return ArrowSizeLarge
	default:
		return strings.TrimSpace(size)
	}
}

func isArrowSize(size string) bool {
	switch normalizeArrowSize(size) {
	case ArrowSizeSmall, ArrowSizeMedium, ArrowSizeLarge:
		return true
	default:
		return false
	}
}

func normalizeConnectionSite(site string) string {
	switch strings.ToLower(strings.TrimSpace(site)) {
	case ConnectionSiteTop:
		return ConnectionSiteTop
	case ConnectionSiteRight:
		return ConnectionSiteRight
	case ConnectionSiteBottom:
		return ConnectionSiteBottom
	case ConnectionSiteLeft:
		return ConnectionSiteLeft
	case "topleft", "top-left", "top_left":
		return ConnectionSiteTopLeft
	case "topright", "top-right", "top_right":
		return ConnectionSiteTopRight
	case "bottomright", "bottom-right", "bottom_right":
		return ConnectionSiteBottomRight
	case "bottomleft", "bottom-left", "bottom_left":
		return ConnectionSiteBottomLeft
	case ConnectionSiteCenter:
		return ConnectionSiteCenter
	default:
		return strings.TrimSpace(site)
	}
}

func connectionSiteIndex(site string) (int, bool) {
	switch normalizeConnectionSite(site) {
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
