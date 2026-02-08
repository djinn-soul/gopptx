package pptx

import (
	"fmt"
	"strings"
)

func validateSlideDrawings(s SlideContent, slideIndex int) error {
	for shapeIndex, shape := range s.Shapes {
		if err := validateShape(shape, slideIndex, shapeIndex+1); err != nil {
			return err
		}
	}
	for connectorIndex, connector := range s.Connectors {
		if err := validateConnector(connector, len(s.Shapes), slideIndex, connectorIndex+1); err != nil {
			return err
		}
	}
	return nil
}

func validateShape(shape Shape, slideIndex int, shapeIndex int) error {
	if !isShapeType(shape.Type) {
		return fmt.Errorf("slide %d shape %d type %q is not supported", slideIndex, shapeIndex, shape.Type)
	}
	if shape.X < 0 || shape.Y < 0 {
		return fmt.Errorf("slide %d shape %d position cannot be negative", slideIndex, shapeIndex)
	}
	if shape.CX <= 0 || shape.CY <= 0 {
		return fmt.Errorf("slide %d shape %d size must be > 0", slideIndex, shapeIndex)
	}
	if shape.Fill != nil && shape.GradientFill != nil {
		return fmt.Errorf("slide %d shape %d cannot set both solid and gradient fill", slideIndex, shapeIndex)
	}
	if shape.Fill != nil {
		if !isHexColor(shape.Fill.Color) {
			return fmt.Errorf("slide %d shape %d fill color must be 6-digit RGB hex", slideIndex, shapeIndex)
		}
		if shape.Fill.TransparencyPct != nil {
			if *shape.Fill.TransparencyPct < 0 || *shape.Fill.TransparencyPct > 100 {
				return fmt.Errorf("slide %d shape %d fill transparency must be in [0,100]", slideIndex, shapeIndex)
			}
		}
	}
	if shape.GradientFill != nil {
		if err := validateShapeGradientFill(*shape.GradientFill, slideIndex, shapeIndex); err != nil {
			return err
		}
	}
	if shape.Line != nil {
		if !isHexColor(shape.Line.Color) {
			return fmt.Errorf("slide %d shape %d line color must be 6-digit RGB hex", slideIndex, shapeIndex)
		}
		if shape.Line.Width <= 0 {
			return fmt.Errorf("slide %d shape %d line width must be > 0", slideIndex, shapeIndex)
		}
		if !isDrawingLineDash(shape.Line.Dash) {
			return fmt.Errorf(
				"slide %d shape %d line dash must be one of solid|dash|dot|dashDot|lgDash|lgDashDot|lgDashDotDot",
				slideIndex,
				shapeIndex,
			)
		}
	}
	if shape.RotationDeg != nil {
		if *shape.RotationDeg < -360 || *shape.RotationDeg > 360 {
			return fmt.Errorf("slide %d shape %d rotation must be in [-360,360]", slideIndex, shapeIndex)
		}
	}
	if shape.Hyperlink != nil {
		if err := validateHyperlink(shape.Hyperlink, fmt.Sprintf("slide %d shape %d", slideIndex, shapeIndex)); err != nil {
			return err
		}
	}
	return nil
}

func validateShapeGradientFill(fill ShapeGradientFill, slideIndex int, shapeIndex int) error {
	if !isShapeGradientType(fill.Type) {
		return fmt.Errorf("slide %d shape %d gradient type %q is not supported", slideIndex, shapeIndex, fill.Type)
	}
	if len(fill.Stops) < 2 {
		return fmt.Errorf("slide %d shape %d gradient must contain at least 2 stops", slideIndex, shapeIndex)
	}

	lastPosition := -1
	for stopIndex, stop := range fill.Stops {
		if stop.PositionPct < 0 || stop.PositionPct > 100 {
			return fmt.Errorf(
				"slide %d shape %d gradient stop %d position must be in [0,100]",
				slideIndex,
				shapeIndex,
				stopIndex+1,
			)
		}
		if stop.PositionPct <= lastPosition {
			return fmt.Errorf(
				"slide %d shape %d gradient stop positions must be strictly increasing",
				slideIndex,
				shapeIndex,
			)
		}
		lastPosition = stop.PositionPct
		if !isHexColor(stop.Color) {
			return fmt.Errorf(
				"slide %d shape %d gradient stop %d color must be 6-digit RGB hex",
				slideIndex,
				shapeIndex,
				stopIndex+1,
			)
		}
		if stop.TransparencyPct != nil && (*stop.TransparencyPct < 0 || *stop.TransparencyPct > 100) {
			return fmt.Errorf(
				"slide %d shape %d gradient stop %d transparency must be in [0,100]",
				slideIndex,
				shapeIndex,
				stopIndex+1,
			)
		}
	}

	if fill.AngleDeg != nil {
		if normalizeShapeGradientType(fill.Type) != ShapeGradientTypeLinear {
			return fmt.Errorf(
				"slide %d shape %d gradient angle is only supported for linear gradients",
				slideIndex,
				shapeIndex,
			)
		}
		if *fill.AngleDeg < -360 || *fill.AngleDeg > 360 {
			return fmt.Errorf(
				"slide %d shape %d gradient angle must be in [-360,360]",
				slideIndex,
				shapeIndex,
			)
		}
	}

	return nil
}

func validateConnector(connector Connector, shapeCount int, slideIndex int, connectorIndex int) error {
	if !isConnectorType(connector.Type) {
		return fmt.Errorf("slide %d connector %d type %q is not supported", slideIndex, connectorIndex, connector.Type)
	}
	if connector.StartX < 0 || connector.StartY < 0 || connector.EndX < 0 || connector.EndY < 0 {
		return fmt.Errorf("slide %d connector %d coordinates cannot be negative", slideIndex, connectorIndex)
	}
	if connector.StartX == connector.EndX && connector.StartY == connector.EndY {
		return fmt.Errorf("slide %d connector %d must have distinct start and end points", slideIndex, connectorIndex)
	}
	if !isHexColor(connector.Line.Color) {
		return fmt.Errorf("slide %d connector %d line color must be 6-digit RGB hex", slideIndex, connectorIndex)
	}
	if connector.Line.Width <= 0 {
		return fmt.Errorf("slide %d connector %d line width must be > 0", slideIndex, connectorIndex)
	}
	if !isDrawingLineDash(connector.Line.Dash) {
		return fmt.Errorf(
			"slide %d connector %d line dash must be one of solid|dash|dot|dashDot|lgDash|lgDashDot|lgDashDotDot",
			slideIndex,
			connectorIndex,
		)
	}
	if !isArrowType(connector.StartArrow) {
		return fmt.Errorf("slide %d connector %d start arrow %q is invalid", slideIndex, connectorIndex, connector.StartArrow)
	}
	if !isArrowType(connector.EndArrow) {
		return fmt.Errorf("slide %d connector %d end arrow %q is invalid", slideIndex, connectorIndex, connector.EndArrow)
	}
	if !isArrowSize(connector.ArrowSize) {
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
	if _, ok := connectionSiteIndex(site); !ok {
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
