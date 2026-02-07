package pptx

import "strings"

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

func resolveConnectorSiteIndices(connector Connector, shapes []Shape) (*int, *int) {
	startSite := connector.StartSite
	if strings.TrimSpace(startSite) == "" && connector.StartShapeIndex > 0 {
		targetX, targetY := connector.EndX, connector.EndY
		if centerX, centerY, ok := shapeCenterForIndex(shapes, connector.EndShapeIndex); ok {
			targetX, targetY = centerX, centerY
		}
		if centerX, centerY, ok := shapeCenterForIndex(shapes, connector.StartShapeIndex); ok {
			startSite = autoConnectionSite(centerX, centerY, targetX, targetY)
		}
	}

	endSite := connector.EndSite
	if strings.TrimSpace(endSite) == "" && connector.EndShapeIndex > 0 {
		targetX, targetY := connector.StartX, connector.StartY
		if centerX, centerY, ok := shapeCenterForIndex(shapes, connector.StartShapeIndex); ok {
			targetX, targetY = centerX, centerY
		}
		if centerX, centerY, ok := shapeCenterForIndex(shapes, connector.EndShapeIndex); ok {
			endSite = autoConnectionSite(centerX, centerY, targetX, targetY)
		}
	}

	return siteIndexPointer(startSite), siteIndexPointer(endSite)
}

func shapeCenterForIndex(shapes []Shape, shapeIndex int) (int64, int64, bool) {
	if shapeIndex <= 0 || shapeIndex > len(shapes) {
		return 0, 0, false
	}
	shape := shapes[shapeIndex-1]
	return shape.X + shape.CX/2, shape.Y + shape.CY/2, true
}

func autoConnectionSite(fromX int64, fromY int64, toX int64, toY int64) string {
	dx := toX - fromX
	if dx < 0 {
		dx = -dx
	}
	dy := toY - fromY
	if dy < 0 {
		dy = -dy
	}

	if dx >= dy {
		if toX >= fromX {
			return ConnectionSiteRight
		}
		return ConnectionSiteLeft
	}
	if toY >= fromY {
		return ConnectionSiteBottom
	}
	return ConnectionSiteTop
}

func siteIndexPointer(site string) *int {
	if idx, ok := connectionSiteIndex(site); ok {
		value := idx
		return &value
	}
	return nil
}
