package shapes

import (
	"fmt"
	"slices"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

// Validate checks connector properties and anchor references.
func (c Connector) Validate(shapeCount int, slideIndex int, connectorIndex int) error {
	if !c.IsDecorative && len(c.AltText) > common.MaxAltTextLength {
		return fmt.Errorf(
			"slide %d connector %d alt text exceeds %d characters",
			slideIndex,
			connectorIndex,
			common.MaxAltTextLength,
		)
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
	if err := c.validateAdjustments(slideIndex, connectorIndex); err != nil {
		return err
	}
	return c.validateActions(slideIndex, connectorIndex)
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
	if c.StartX < 0 || c.StartY < 0 || c.EndX < 0 || c.EndY < 0 {
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
		return fmt.Errorf(
			"slide %d connector %d start arrow width %q is invalid",
			slideIndex, connectorIndex, c.StartArrowWidth,
		)
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
	if len(c.Adjustments) == 0 {
		return nil
	}
	ct := NormalizeConnectorType(c.Type)
	if ct != ConnectorTypeElbow && ct != ConnectorTypeCurved {
		return fmt.Errorf(
			"slide %d connector %d adjustments are only supported for elbow/curved connectors",
			slideIndex, connectorIndex,
		)
	}
	return nil
}

func (c Connector) validateActions(slideIndex, connectorIndex int) error {
	if c.ClickAction != nil {
		if err := c.ClickAction.Validate(); err != nil {
			return fmt.Errorf("slide %d connector %d has invalid click action: %w", slideIndex, connectorIndex, err)
		}
	}
	if c.HoverAction != nil {
		if err := c.HoverAction.Validate(); err != nil {
			return fmt.Errorf("slide %d connector %d has invalid hover action: %w", slideIndex, connectorIndex, err)
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
