package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type (
	// Connector is one connector.
	Connector = elements.Connector
)

const (
	ConnectorTypeStraight = elements.ConnectorTypeStraight
	ConnectorTypeElbow    = elements.ConnectorTypeElbow
	ConnectorTypeCurved   = elements.ConnectorTypeCurved
)

const (
	ArrowTypeNone     = elements.ArrowTypeNone
	ArrowTypeTriangle = elements.ArrowTypeTriangle
	ArrowTypeStealth  = elements.ArrowTypeStealth
	ArrowTypeDiamond  = elements.ArrowTypeDiamond
	ArrowTypeOval     = elements.ArrowTypeOval
	ArrowTypeOpen     = elements.ArrowTypeOpen
)

const (
	ArrowSizeSmall  = elements.ArrowSizeSmall
	ArrowSizeMedium = elements.ArrowSizeMedium
	ArrowSizeLarge  = elements.ArrowSizeLarge
)

const (
	ConnectionSiteTop         = elements.ConnectionSiteTop
	ConnectionSiteRight       = elements.ConnectionSiteRight
	ConnectionSiteBottom      = elements.ConnectionSiteBottom
	ConnectionSiteLeft        = elements.ConnectionSiteLeft
	ConnectionSiteTopLeft     = elements.ConnectionSiteTopLeft
	ConnectionSiteTopRight    = elements.ConnectionSiteTopRight
	ConnectionSiteBottomRight = elements.ConnectionSiteBottomRight
	ConnectionSiteBottomLeft  = elements.ConnectionSiteBottomLeft
	ConnectionSiteCenter      = elements.ConnectionSiteCenter
)

func NewConnector(connectorType string, startX, startY, endX, endY int64) Connector {
	return elements.NewConnector(connectorType, startX, startY, endX, endY)
}

func NewStraightConnector(startX, startY, endX, endY int64) Connector {
	return elements.NewStraightConnector(startX, startY, endX, endY)
}

func NewElbowConnector(startX, startY, endX, endY int64) Connector {
	return elements.NewElbowConnector(startX, startY, endX, endY)
}

func NewCurvedConnector(startX, startY, endX, endY int64) Connector {
	return elements.NewCurvedConnector(startX, startY, endX, endY)
}

func normalizeConnectorType(connectorType string) string {
	return elements.NormalizeConnectorType(connectorType)
}

func normalizeArrowType(arrowType string) string {
	return elements.NormalizeArrowType(arrowType)
}

func normalizeArrowSize(size string) string {
	return elements.NormalizeArrowSize(size)
}

func normalizeConnectionSite(site string) string {
	return elements.NormalizeConnectionSite(site)
}
