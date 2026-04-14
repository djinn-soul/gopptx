package shapes

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
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
	ClickAction     *action.Hyperlink
	HoverAction     *action.Hyperlink
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
	newAdj := make([]ConnectorAdjustment, len(c.Adjustments), len(c.Adjustments)+1)
	copy(newAdj, c.Adjustments)
	c.Adjustments = append(newAdj, ConnectorAdjustment{
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
