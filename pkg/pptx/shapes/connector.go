package shapes

import (
	"fmt"
	"strings"
)

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
	if err := connector.Line.Validate(); err != nil {
		return fmt.Errorf("slide %d connector %d invalid line: %w", slideIndex, connectorIndex, err)
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
