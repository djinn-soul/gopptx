package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

type (
	// Connector is one connector.
	Connector = shapes.Connector
	// ConnectorAdjustment is one connector geometry adjustment point.
	ConnectorAdjustment = shapes.ConnectorAdjustment
)

const (
	ConnectorTypeStraight = shapes.ConnectorTypeStraight
	ConnectorTypeElbow    = shapes.ConnectorTypeElbow
	ConnectorTypeCurved   = shapes.ConnectorTypeCurved

	ArrowTypeNone     = shapes.ArrowTypeNone
	ArrowTypeTriangle = shapes.ArrowTypeTriangle
	ArrowTypeStealth  = shapes.ArrowTypeStealth
	ArrowTypeDiamond  = shapes.ArrowTypeDiamond
	ArrowTypeOval     = shapes.ArrowTypeOval
	ArrowTypeOpen     = shapes.ArrowTypeOpen

	ArrowSizeSmall  = shapes.ArrowSizeSmall
	ArrowSizeMedium = shapes.ArrowSizeMedium
	ArrowSizeLarge  = shapes.ArrowSizeLarge

	ConnectionSiteTop         = shapes.ConnectionSiteTop
	ConnectionSiteRight       = shapes.ConnectionSiteRight
	ConnectionSiteBottom      = shapes.ConnectionSiteBottom
	ConnectionSiteLeft        = shapes.ConnectionSiteLeft
	ConnectionSiteTopLeft     = shapes.ConnectionSiteTopLeft
	ConnectionSiteTopRight    = shapes.ConnectionSiteTopRight
	ConnectionSiteBottomRight = shapes.ConnectionSiteBottomRight
	ConnectionSiteBottomLeft  = shapes.ConnectionSiteBottomLeft
	ConnectionSiteCenter      = shapes.ConnectionSiteCenter
)

func NewConnector(connectorType string, startX, startY, endX, endY styling.Length) Connector {
	return shapes.NewConnector(connectorType, startX, startY, endX, endY)
}

func NewStraightConnector(startX, startY, endX, endY styling.Length) Connector {
	return shapes.NewStraightConnector(startX, startY, endX, endY)
}

func NewElbowConnector(startX, startY, endX, endY styling.Length) Connector {
	return shapes.NewElbowConnector(startX, startY, endX, endY)
}

func NewCurvedConnector(startX, startY, endX, endY styling.Length) Connector {
	return shapes.NewCurvedConnector(startX, startY, endX, endY)
}

// ConnectStartAuto anchors the connector start to a shape and auto-selects the site.
func ConnectStartAuto(c Connector, shapeIndex int) Connector {
	return c.ConnectStartAuto(shapeIndex)
}

// ConnectEndAuto anchors the connector end to a shape and auto-selects the site.
func ConnectEndAuto(c Connector, shapeIndex int) Connector {
	return c.ConnectEndAuto(shapeIndex)
}

// AutoReroute recalculates connector sites from current shape positions.
func AutoReroute(c Connector, shapes []Shape) Connector {
	return c.AutoReroute(shapes)
}
