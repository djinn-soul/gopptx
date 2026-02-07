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

	return siteIndexPointer(startSite), siteIndexPointer(endSite)
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

func siteIndexPointer(site string) *int {
	if idx, ok := connectionSiteIndex(site); ok {
		value := idx
		return &value
	}
	return nil
}
