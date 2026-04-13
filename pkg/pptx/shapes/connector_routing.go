package shapes

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// AutoReroute recalculates start/end connector sites from current shape positions.
// Endpoints anchored to shapes are rewritten to the nearest valid site.
func (c Connector) AutoReroute(shapes []Shape) Connector {
	if c.StartShapeIndex > 0 {
		tmp := c
		tmp.StartSite = ""
		startIdx, _ := ResolveConnectorSiteIndices(tmp, shapes)
		if startIdx != nil {
			c.StartSite = siteFromIndex(*startIdx)
			if anchorX, anchorY, ok := shapeAnchorPointForIndex(shapes, c.StartShapeIndex, c.StartSite); ok {
				c.StartX = anchorX
				c.StartY = anchorY
			}
		}
	}
	if c.EndShapeIndex > 0 {
		tmp := c
		tmp.EndSite = ""
		_, endIdx := ResolveConnectorSiteIndices(tmp, shapes)
		if endIdx != nil {
			c.EndSite = siteFromIndex(*endIdx)
			if anchorX, anchorY, ok := shapeAnchorPointForIndex(shapes, c.EndShapeIndex, c.EndSite); ok {
				c.EndX = anchorX
				c.EndY = anchorY
			}
		}
	}
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

func shapeCenterForIndex(shapes []Shape, shapeIndex int) (styling.Length, styling.Length, bool) {
	shape, ok := shapeForIndex(shapes, shapeIndex)
	if !ok {
		return 0, 0, false
	}
	return shape.X + shape.CX/2, shape.Y + shape.CY/2, true
}

func autoConnectionSite(shape Shape, targetX styling.Length, targetY styling.Length) string {
	candidates := shapeConnectionSiteCandidatesForType(shape)
	bestSite := ConnectionSiteCenter
	var bestDistance styling.Length
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
	x    styling.Length
	y    styling.Length
}

const connectionSiteCenterDivisor = 2

func shapeConnectionSiteCandidates(shape Shape) [9]connectionSiteCandidate {
	cx := shape.X + shape.CX/connectionSiteCenterDivisor
	cy := shape.Y + shape.CY/connectionSiteCenterDivisor
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

func shapeConnectionSiteCandidatesForType(shape Shape) []connectionSiteCandidate {
	all := shapeConnectionSiteCandidates(shape)
	allowed := allowedConnectionSitesForShape(shape.Type)
	out := make([]connectionSiteCandidate, 0, len(allowed))
	for _, site := range allowed {
		for _, candidate := range all {
			if candidate.site == site {
				out = append(out, candidate)
				break
			}
		}
	}
	if len(out) == 0 {
		return all[:]
	}
	return out
}

func allowedConnectionSitesForShape(shapeType string) []string {
	switch NormalizeShapeType(shapeType) {
	case ShapeTypeEllipse, ShapeTypeCloud, ShapeTypeHeart:
		return []string{
			ConnectionSiteTop,
			ConnectionSiteRight,
			ConnectionSiteBottom,
			ConnectionSiteLeft,
			ConnectionSiteCenter,
		}
	default:
		return []string{
			ConnectionSiteTop,
			ConnectionSiteRight,
			ConnectionSiteBottom,
			ConnectionSiteLeft,
			ConnectionSiteTopLeft,
			ConnectionSiteTopRight,
			ConnectionSiteBottomRight,
			ConnectionSiteBottomLeft,
			ConnectionSiteCenter,
		}
	}
}

func siteFromIndex(idx int) string {
	switch idx {
	case connectionSiteTopIndex:
		return ConnectionSiteTop
	case connectionSiteRightIndex:
		return ConnectionSiteRight
	case connectionSiteBottomIndex:
		return ConnectionSiteBottom
	case connectionSiteLeftIndex:
		return ConnectionSiteLeft
	case connectionSiteTopLeftIndex:
		return ConnectionSiteTopLeft
	case connectionSiteTopRightIndex:
		return ConnectionSiteTopRight
	case connectionSiteBottomRightIndex:
		return ConnectionSiteBottomRight
	case connectionSiteBottomLeftIndex:
		return ConnectionSiteBottomLeft
	case connectionSiteCenterIndex:
		return ConnectionSiteCenter
	default:
		return ""
	}
}

func shapeAnchorPointForIndex(shapes []Shape, shapeIndex int, site string) (styling.Length, styling.Length, bool) {
	shape, ok := shapeForIndex(shapes, shapeIndex)
	if !ok {
		return 0, 0, false
	}
	for _, candidate := range shapeConnectionSiteCandidatesForType(shape) {
		if candidate.site == site {
			return candidate.x, candidate.y, true
		}
	}
	return 0, 0, false
}

func SiteIndexPointer(site string) *int {
	if idx, ok := ConnectionSiteIndex(site); ok {
		value := idx
		return &value
	}
	return nil
}
