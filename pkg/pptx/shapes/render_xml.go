package shapes

import (
	"math"
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
		Name:         shape.Name,
		Adjustments:  toXMLShapeAdjustments(shape.Adjustments),
	}
	if shape.Effects != nil {
		spec.Effects = &pptxxml.ShapeEffectsSpec{
			Shadow:     shape.Effects.Shadow,
			Glow:       shape.Effects.Glow,
			SoftEdges:  shape.Effects.SoftEdges,
			Reflection: shape.Effects.Reflection,
		}
	}

	// Handle rich fill first (takes precedence over legacy fill)
	if shape.RichFill != nil {
		spec.RichFill = toXMLRichFillSpec(shape.RichFill)
	} else {
		// Legacy fill handling
		if shape.Fill != nil {
			spec.Fill = &pptxxml.ShapeFillSpec{
				Color:        common.NormalizeHexColor(shape.Fill.Color),
				Transparency: shape.Fill.Transparency,
			}
		}
		if shape.GradientFill != nil {
			spec.GradientFill = toXMLGradientFillSpec(shape.GradientFill)
		}
	}

	// Handle rich line first (takes precedence over legacy line)
	if shape.RichLine != nil {
		spec.RichLine = toXMLRichLineSpec(shape.RichLine)
	} else if shape.Line != nil {
		spec.Line = &pptxxml.ShapeLineSpec{
			Color: common.NormalizeHexColor(shape.Line.Color),
			Width: shape.Line.Width.Emu(),
			Dash:  NormalizeDrawingLineDash(shape.Line.Dash),
			Cap:   shape.Line.Cap,
			Join:  shape.Line.Join,
		}
	}

	// Handle rich shadow
	if shape.RichShadow != nil {
		spec.RichShadow = toXMLRichShadowSpec(shape.RichShadow)
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

func toXMLRichFillSpec(fill *RichShapeFill) *pptxxml.RichShapeFillSpec {
	spec := &pptxxml.RichShapeFillSpec{
		Type: pptxxml.FillType(fill.Type),
	}

	switch fill.Type {
	case FillTypeSolid:
		if fill.Solid != nil {
			spec.Solid = &pptxxml.SolidFillSpec{
				Color:        fill.Solid.Color,
				Transparency: fill.Solid.Transparency,
			}
		}
	case FillTypeGradient:
		if fill.Gradient != nil {
			spec.Gradient = toXMLGradientFillSpec(fill.Gradient)
		}
	case FillTypePattern:
		if fill.Pattern != nil {
			spec.Pattern = &pptxxml.PatternFillSpec{
				Pattern: string(fill.Pattern.Pattern),
				FgColor: fill.Pattern.FgColor,
				BgColor: fill.Pattern.BgColor,
			}
		}
	case FillTypeNoFill:
		// No nested payload is required for explicit no-fill.
	}

	return spec
}

func toXMLRichLineSpec(line *RichShapeLine) *pptxxml.RichShapeLineSpec {
	return &pptxxml.RichShapeLineSpec{
		Color:        line.Color,
		Width:        line.Width.Emu(),
		DashStyle:    pptxxml.LineDashStyle(line.DashStyle),
		CapStyle:     pptxxml.LineCapStyle(line.CapStyle),
		JoinStyle:    pptxxml.LineJoinStyle(line.JoinStyle),
		Transparency: line.Transparency,
	}
}

func toXMLRichShadowSpec(shadow *RichShapeShadow) *pptxxml.RichShapeShadowSpec {
	return &pptxxml.RichShapeShadowSpec{
		Type:            pptxxml.ShadowType(shadow.Type),
		Color:           shadow.Color,
		Transparency:    shadow.Transparency,
		BlurRadius:      shadow.BlurRadius,
		Distance:        shadow.Distance,
		Angle:           shadow.Angle,
		Alignment:       string(shadow.Alignment),
		SkewX:           shadow.SkewX,
		SkewY:           shadow.SkewY,
		ScaleX:          shadow.ScaleX,
		ScaleY:          shadow.ScaleY,
		RotateWithShape: shadow.RotateWithShape,
	}
}

func toXMLTextFrameSpec(tf *TextFrame) *pptxxml.TextFrameSpec {
	var rotation *int64
	if tf.RotationDeg != nil {
		value := int64(math.Round(*tf.RotationDeg * float64(ooxmlAngleUnitsPerDegree)))
		rotation = &value
	}
	return &pptxxml.TextFrameSpec{
		MarginLeft:   tf.MarginLeft.Emu(),
		MarginRight:  tf.MarginRight.Emu(),
		MarginTop:    tf.MarginTop.Emu(),
		MarginBottom: tf.MarginBottom.Emu(),
		Anchor:       string(tf.Anchor),
		Wrap:         string(tf.Wrap),
		AutoFit:      string(tf.AutoFit),
		Rotation:     rotation,
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

func toXMLShapeAdjustments(adjustments []ShapeAdjustment) []pptxxml.ConnectorAdjustmentSpec {
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
