package shapes

import (
	"fmt"
	"slices"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Connector is one connector.
type Connector struct {
	Type            string
	StartX          styling.Length
	StartY          styling.Length
	EndX            styling.Length
	EndY            styling.Length
	Line            ShapeLine
	StartArrow      string
	StartArrowWidth string
	StartArrowLen   string
	EndArrow        string
	EndArrowWidth   string
	EndArrowLen     string
	StartShapeIndex int
	StartSite       string
	EndShapeIndex   int
	EndSite         string
	Label           string
	AltText         string
	IsDecorative    bool
	Placeholder     *Placeholder
	Adjustments     []ConnectorAdjustment
}

// ConnectorAdjustment represents one connector geometry adjustment point (<a:gd>) entry.
type ConnectorAdjustment struct {
	Name    string
	Formula string
}

// NewConnector creates a connector.
func NewConnector(connectorType string, startX, startY, endX, endY styling.Length) Connector {
	return Connector{
		Type:            NormalizeConnectorType(connectorType),
		StartX:          startX,
		StartY:          startY,
		EndX:            endX,
		EndY:            endY,
		Line:            NewShapeLine("000000", styling.Emu(int64(styling.EmuPerPt))),
		StartArrow:      ArrowTypeNone,
		StartArrowWidth: ArrowSizeMedium,
		StartArrowLen:   ArrowSizeMedium,
		EndArrow:        ArrowTypeNone,
		EndArrowWidth:   ArrowSizeMedium,
		EndArrowLen:     ArrowSizeMedium,
	}
}

func NewStraightConnector(startX, startY, endX, endY styling.Length) Connector {
	return NewConnector(ConnectorTypeStraight, startX, startY, endX, endY)
}

func NewElbowConnector(startX, startY, endX, endY styling.Length) Connector {
	return NewConnector(ConnectorTypeElbow, startX, startY, endX, endY)
}

func NewCurvedConnector(startX, startY, endX, endY styling.Length) Connector {
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

// WithArrowSize sets arrowhead size (both width and length) for both ends.
func (c Connector) WithArrowSize(size string) Connector {
	return c.WithStartArrowSize(size, size).WithEndArrowSize(size, size)
}

// WithStartArrowSize sets start arrowhead width and length.
func (c Connector) WithStartArrowSize(width, length string) Connector {
	return c.WithStartArrowWidth(width).WithStartArrowLen(length)
}

// WithEndArrowSize sets end arrowhead width and length.
func (c Connector) WithEndArrowSize(width, length string) Connector {
	return c.WithEndArrowWidth(width).WithEndArrowLen(length)
}

// WithStartArrowWidth sets start arrowhead width.
func (c Connector) WithStartArrowWidth(width string) Connector {
	c.StartArrowWidth = NormalizeArrowSize(width)
	return c
}

// WithStartArrowLen sets start arrowhead length.
func (c Connector) WithStartArrowLen(length string) Connector {
	c.StartArrowLen = NormalizeArrowSize(length)
	return c
}

// WithEndArrowWidth sets end arrowhead width.
func (c Connector) WithEndArrowWidth(width string) Connector {
	c.EndArrowWidth = NormalizeArrowSize(width)
	return c
}

// WithEndArrowLen sets end arrowhead length.
func (c Connector) WithEndArrowLen(length string) Connector {
	c.EndArrowLen = NormalizeArrowSize(length)
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

// WithAdjustment appends one geometry adjustment point for elbow/curved connectors.
func (c Connector) WithAdjustment(name, formula string) Connector {
	c.Adjustments = append(c.Adjustments, ConnectorAdjustment{
		Name:    strings.TrimSpace(name),
		Formula: strings.TrimSpace(formula),
	})
	return c
}

// WithAdjustmentValue appends one "val" adjustment helper entry for elbow/curved connectors.
func (c Connector) WithAdjustmentValue(name string, value int) Connector {
	return c.WithAdjustment(name, fmt.Sprintf("val %d", value))
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

func shapeConnectionSiteCandidates(shape Shape) [9]connectionSiteCandidate {
	//nolint:mnd // Center point divisor
	cx := shape.X + shape.CX/2
	//nolint:mnd // Center point divisor
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
	//nolint:mnd // OOXML connection site indices
	switch idx {
	case 0:
		return ConnectionSiteTop
	case 1:
		return ConnectionSiteRight
	case 2:
		return ConnectionSiteBottom
	case 3:
		return ConnectionSiteLeft
	case 4:
		return ConnectionSiteTopLeft
	case 5:
		return ConnectionSiteTopRight
	case 6:
		return ConnectionSiteBottomRight
	case 7:
		return ConnectionSiteBottomLeft
	case 8:
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

// Validate checks connector properties and anchor references.
func (c Connector) Validate(shapeCount int, slideIndex int, connectorIndex int) error {
	if !c.IsDecorative && len(c.AltText) > common.MaxAltTextLength {
		return fmt.Errorf("slide %d connector %d alt text exceeds %d characters", slideIndex, connectorIndex, common.MaxAltTextLength)
	}

	if err := c.validateBasicProps(slideIndex, connectorIndex); err != nil {
		return err
	}
	if err := c.validateArrows(slideIndex, connectorIndex); err != nil {
		return err
	}
	if err := c.validateAnchors(shapeCount, slideIndex, connectorIndex); err != nil {
		return err
	}
	return c.validateAdjustments(slideIndex, connectorIndex)
}

func (c Connector) validateBasicProps(slideIndex, connectorIndex int) error {
	if !IsConnectorType(c.Type) {
		return fmt.Errorf(
			"slide %d connector %d type %q is not supported",
			slideIndex,
			connectorIndex,
			c.Type,
		)
	}
	if c.StartX < 0 || c.StartY < 0 ||
		c.EndX < 0 || c.EndY < 0 {
		return fmt.Errorf(
			"slide %d connector %d coordinates cannot be negative",
			slideIndex, connectorIndex,
		)
	}
	if c.StartX == c.EndX && c.StartY == c.EndY {
		return fmt.Errorf("slide %d connector %d must have distinct start and end points", slideIndex, connectorIndex)
	}
	if err := c.Line.Validate(); err != nil {
		return fmt.Errorf("slide %d connector %d invalid line: %w", slideIndex, connectorIndex, err)
	}
	return nil
}

func (c Connector) validateArrows(slideIndex, connectorIndex int) error {
	if !IsArrowType(c.StartArrow) {
		return fmt.Errorf("slide %d connector %d start arrow %q is invalid", slideIndex, connectorIndex, c.StartArrow)
	}
	if !IsArrowType(c.EndArrow) {
		return fmt.Errorf("slide %d connector %d end arrow %q is invalid", slideIndex, connectorIndex, c.EndArrow)
	}
	if !IsArrowSize(c.StartArrowWidth) {
		return fmt.Errorf("slide %d connector %d start arrow width %q is invalid",
			slideIndex, connectorIndex, c.StartArrowWidth)
	}
	if !IsArrowSize(c.StartArrowLen) {
		return fmt.Errorf(
			"slide %d connector %d start arrow length %q is invalid",
			slideIndex, connectorIndex, c.StartArrowLen,
		)
	}
	if !IsArrowSize(c.EndArrowWidth) {
		return fmt.Errorf(
			"slide %d connector %d end arrow width %q is invalid",
			slideIndex, connectorIndex, c.EndArrowWidth,
		)
	}
	if !IsArrowSize(c.EndArrowLen) {
		return fmt.Errorf(
			"slide %d connector %d end arrow length %q is invalid",
			slideIndex, connectorIndex, c.EndArrowLen,
		)
	}
	return nil
}

func (c Connector) validateAnchors(shapeCountUX int, slideIndex, connectorIndex int) error {
	if err := validateConnectorAnchor(
		"start", c.StartShapeIndex, c.StartSite, shapeCountUX, slideIndex, connectorIndex,
	); err != nil {
		return err
	}
	return validateConnectorAnchor("end", c.EndShapeIndex, c.EndSite, shapeCountUX, slideIndex, connectorIndex)
}

func (c Connector) validateAdjustments(slideIndex, connectorIndex int) error {
	for i, adj := range c.Adjustments {
		if strings.TrimSpace(adj.Name) == "" {
			return fmt.Errorf(
				"slide %d connector %d adjustment %d name cannot be empty",
				slideIndex, connectorIndex, i+1,
			)
		}
		if strings.TrimSpace(adj.Formula) == "" {
			return fmt.Errorf(
				"slide %d connector %d adjustment %d formula cannot be empty",
				slideIndex, connectorIndex, i+1,
			)
		}
	}
	if len(c.Adjustments) > 0 {
		ct := NormalizeConnectorType(c.Type)
		if ct != ConnectorTypeElbow && ct != ConnectorTypeCurved {
			return fmt.Errorf(
				"slide %d connector %d adjustments are only supported for elbow/curved connectors",
				slideIndex, connectorIndex,
			)
		}
	}
	return nil
}

// ValidateWithShapes checks connector fields and validates site compatibility with connected shape types.
func (c Connector) ValidateWithShapes(shapes []Shape, slideIndex, connectorIndex int) error {
	if err := c.Validate(len(shapes), slideIndex, connectorIndex); err != nil {
		return err
	}
	if err := validateConnectorSiteForShape(
		"start",
		c.StartShapeIndex,
		c.StartSite,
		shapes,
		slideIndex,
		connectorIndex,
	); err != nil {
		return err
	}
	if err := validateConnectorSiteForShape(
		"end",
		c.EndShapeIndex,
		c.EndSite,
		shapes,
		slideIndex,
		connectorIndex,
	); err != nil {
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

func validateConnectorSiteForShape(
	side string,
	shapeIndex int,
	site string,
	shapes []Shape,
	slideIndex int,
	connectorIndex int,
) error {
	if shapeIndex <= 0 || strings.TrimSpace(site) == "" {
		return nil
	}
	shape, ok := shapeForIndex(shapes, shapeIndex)
	if !ok {
		return nil
	}
	normalized := NormalizeConnectionSite(site)
	if slices.Contains(allowedConnectionSitesForShape(shape.Type), normalized) {
		return nil
	}
	return fmt.Errorf(
		"slide %d connector %d %s site %q is not supported for shape type %q",
		slideIndex,
		connectorIndex,
		side,
		normalized,
		shape.Type,
	)
}
