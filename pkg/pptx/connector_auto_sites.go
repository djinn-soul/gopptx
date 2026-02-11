package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// ConnectStartAuto anchors the connector start to a shape and auto-selects the site.
func ConnectStartAuto(c Connector, shapeIndex int) Connector {
	return c.ConnectStartAuto(shapeIndex)
}

// ConnectEndAuto anchors the connector end to a shape and auto-selects the site.
func ConnectEndAuto(c Connector, shapeIndex int) Connector {
	return c.ConnectEndAuto(shapeIndex)
}

func resolveConnectorSiteIndices(connector Connector, shapes []Shape) (*int, *int) {
	return elements.ResolveConnectorSiteIndices(connector, shapes)
}
