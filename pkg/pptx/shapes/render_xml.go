package shapes

import (
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

func ToXMLShapeSpecs(shapes []Shape, hyperlinkRIDs map[*action.Hyperlink]string) []pptxxml.ShapeSpec {
	if len(shapes) == 0 {
		return nil
	}
	specs := make([]pptxxml.ShapeSpec, 0, len(shapes))
	for _, shape := range shapes {
		spec := toXMLShapeSpec(shape, hyperlinkRIDs)
		specs = append(specs, spec)
	}
	return specs
}

func toXMLShapeSpec(shape Shape, hyperlinkRIDs map[*action.Hyperlink]string) pptxxml.ShapeSpec {
	spec := pptxxml.ShapeSpec{
		Type:         NormalizeShapeType(shape.Type),
		X:            shape.X.Emu(),
		Y:            shape.Y.Emu(),
		CX:           shape.CX.Emu(),
		CY:           shape.CY.Emu(),
		Text:         shape.Text,
		AltText:      shape.AltText,
		IsDecorative: shape.IsDecorative,
		RotationDeg:  shape.RotationDeg,
	}

	if shape.Fill != nil {
		spec.Fill = &pptxxml.ShapeFillSpec{
			Color:        common.NormalizeHexColor(shape.Fill.Color),
			Transparency: shape.Fill.Transparency,
		}
	}
	if shape.GradientFill != nil {
		spec.GradientFill = toXMLGradientFillSpec(shape.GradientFill)
	}
	if shape.Line != nil {
		spec.Line = &pptxxml.ShapeLineSpec{
			Color: common.NormalizeHexColor(shape.Line.Color),
			Width: shape.Line.Width.Emu(),
			Dash:  NormalizeDrawingLineDash(shape.Line.Dash),
		}
	}

	if shape.TextFrame != nil {
		spec.TextFrame = toXMLTextFrameSpec(shape.TextFrame)
	}

	spec.ClickAction = resolveActionSpec(shape.ClickAction, shape.Hyperlink, hyperlinkRIDs)
	spec.HoverAction = resolveActionSpec(shape.HoverAction, nil, hyperlinkRIDs)

	return spec
}

func toXMLGradientFillSpec(fill *ShapeGradientFill) *pptxxml.ShapeGradientFillSpec {
	stops := make([]pptxxml.ShapeGradientStopSpec, 0, len(fill.Stops))
	for _, stop := range fill.Stops {
		stops = append(stops, pptxxml.ShapeGradientStopSpec{
			PositionPct:  stop.PositionPct,
			Color:        common.NormalizeHexColor(stop.Color),
			Transparency: stop.Transparency,
		})
	}
	return &pptxxml.ShapeGradientFillSpec{
		Type:     NormalizeShapeGradientType(fill.Type),
		Stops:    stops,
		AngleDeg: fill.AngleDeg,
	}
}

func toXMLTextFrameSpec(tf *TextFrame) *pptxxml.TextFrameSpec {
	return &pptxxml.TextFrameSpec{
		MarginLeft:   tf.MarginLeft.Emu(),
		MarginRight:  tf.MarginRight.Emu(),
		MarginTop:    tf.MarginTop.Emu(),
		MarginBottom: tf.MarginBottom.Emu(),
		Anchor:       string(tf.Anchor),
		Wrap:         string(tf.Wrap),
		AutoFit:      string(tf.AutoFit),
	}
}

func resolveActionSpec(primary, secondary *action.Hyperlink, rids map[*action.Hyperlink]string) *pptxxml.HyperlinkSpec {
	h := primary
	if h == nil {
		h = secondary
	}
	if h == nil {
		return nil
	}

	if rid, ok := rids[h]; ok {
		return &pptxxml.HyperlinkSpec{
			RelID:          rid,
			Tooltip:        h.Tooltip,
			HighlightClick: h.HighlightClick,
			Action:         h.Action.ActionType(),
		}
	}
	return nil
}

func ToXMLConnectorSpecs(connectors []Connector, shapes []Shape) []pptxxml.ConnectorSpec {
	if len(connectors) == 0 {
		return nil
	}
	specs := make([]pptxxml.ConnectorSpec, 0, len(connectors))
	for _, connector := range connectors {
		startSiteIndex, endSiteIndex := ResolveConnectorSiteIndices(connector, shapes)
		spec := ToXMLConnectorSpec(connector, startSiteIndex, endSiteIndex)
		specs = append(specs, spec)
	}
	return specs
}

func ToXMLConnectorSpec(connector Connector, startSiteIndex *int, endSiteIndex *int) pptxxml.ConnectorSpec {
	return pptxxml.ConnectorSpec{
		Type:            NormalizeConnectorType(connector.Type),
		StartX:          connector.StartX.Emu(),
		StartY:          connector.StartY.Emu(),
		EndX:            connector.EndX.Emu(),
		EndY:            connector.EndY.Emu(),
		Line:            ToXMLShapeLineSpec(connector.Line),
		StartArrow:      NormalizeArrowType(connector.StartArrow),
		StartArrowWidth: NormalizeArrowSize(connector.StartArrowWidth),
		StartArrowLen:   NormalizeArrowSize(connector.StartArrowLen),
		EndArrow:        NormalizeArrowType(connector.EndArrow),
		EndArrowWidth:   NormalizeArrowSize(connector.EndArrowWidth),
		EndArrowLen:     NormalizeArrowSize(connector.EndArrowLen),
		StartShapeIndex: connector.StartShapeIndex,
		StartSiteIndex:  startSiteIndex,
		EndShapeIndex:   connector.EndShapeIndex,
		EndSiteIndex:    endSiteIndex,
		Label:           connector.Label,
		AltText:         connector.AltText,
		IsDecorative:    connector.IsDecorative,
		Adjustments:     toXMLConnectorAdjustments(connector.Adjustments),
	}
}

func ToXMLShapeLineSpec(line ShapeLine) pptxxml.ShapeLineSpec {
	lineCap := strings.TrimSpace(line.Cap)
	if lineCap != "" {
		lineCap = NormalizeLineCap(lineCap)
	}
	join := strings.TrimSpace(line.Join)
	if join != "" {
		join = NormalizeLineJoin(join)
	}
	return pptxxml.ShapeLineSpec{
		Color: common.NormalizeHexColor(line.Color),
		Width: line.Width.Emu(),
		Dash:  NormalizeDrawingLineDash(line.Dash),
		Cap:   lineCap,
		Join:  join,
	}
}

func toXMLConnectorAdjustments(adjustments []ConnectorAdjustment) []pptxxml.ConnectorAdjustmentSpec {
	if len(adjustments) == 0 {
		return nil
	}
	specs := make([]pptxxml.ConnectorAdjustmentSpec, 0, len(adjustments))
	for _, adj := range adjustments {
		specs = append(specs, pptxxml.ConnectorAdjustmentSpec{
			Name:    adj.Name,
			Formula: adj.Formula,
		})
	}
	return specs
}
