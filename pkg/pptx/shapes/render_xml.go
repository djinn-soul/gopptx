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
		spec := pptxxml.ShapeSpec{
			Type:         NormalizeShapeType(shape.Type),
			X:            shape.X.Emu(),
			Y:            shape.Y.Emu(),
			CX:           shape.CX.Emu(),
			CY:           shape.CY.Emu(),
			Text:         shape.Text,
			AltText:      shape.AltText,
			IsDecorative: shape.IsDecorative,
		}

		if shape.Fill != nil {
			spec.Fill = &pptxxml.ShapeFillSpec{
				Color:        common.NormalizeHexColor(shape.Fill.Color),
				Transparency: shape.Fill.Transparency,
			}
		}
		if shape.GradientFill != nil {
			stops := make([]pptxxml.ShapeGradientStopSpec, 0, len(shape.GradientFill.Stops))
			for _, stop := range shape.GradientFill.Stops {
				stops = append(stops, pptxxml.ShapeGradientStopSpec{
					PositionPct:  stop.PositionPct,
					Color:        common.NormalizeHexColor(stop.Color),
					Transparency: stop.Transparency,
				})
			}
			spec.GradientFill = &pptxxml.ShapeGradientFillSpec{
				Type:     NormalizeShapeGradientType(shape.GradientFill.Type),
				Stops:    stops,
				AngleDeg: shape.GradientFill.AngleDeg,
			}
		}
		if shape.Line != nil {
			spec.Line = &pptxxml.ShapeLineSpec{
				Color: common.NormalizeHexColor(shape.Line.Color),
				Width: shape.Line.Width.Emu(),
				Dash:  NormalizeDrawingLineDash(shape.Line.Dash),
			}
		}

		spec.RotationDeg = shape.RotationDeg
		if shape.TextFrame != nil {
			spec.TextFrame = &pptxxml.TextFrameSpec{
				MarginLeft:   shape.TextFrame.MarginLeft.Emu(),
				MarginRight:  shape.TextFrame.MarginRight.Emu(),
				MarginTop:    shape.TextFrame.MarginTop.Emu(),
				MarginBottom: shape.TextFrame.MarginBottom.Emu(),
				Anchor:       string(shape.TextFrame.Anchor),

				Wrap:    string(shape.TextFrame.Wrap),
				AutoFit: string(shape.TextFrame.AutoFit),
			}
		}
		if shape.ClickAction != nil {
			if rid, ok := hyperlinkRIDs[shape.ClickAction]; ok {
				spec.ClickAction = &pptxxml.HyperlinkSpec{
					RelID:          rid,
					Tooltip:        shape.ClickAction.Tooltip,
					HighlightClick: shape.ClickAction.HighlightClick,
					Action:         shape.ClickAction.Action.ActionType(),
				}
			}
		} else if shape.Hyperlink != nil {
			if rid, ok := hyperlinkRIDs[shape.Hyperlink]; ok {
				spec.ClickAction = &pptxxml.HyperlinkSpec{
					RelID:          rid,
					Tooltip:        shape.Hyperlink.Tooltip,
					HighlightClick: shape.Hyperlink.HighlightClick,
					Action:         shape.Hyperlink.Action.ActionType(),
				}
			}
		}

		if shape.HoverAction != nil {
			if rid, ok := hyperlinkRIDs[shape.HoverAction]; ok {
				spec.HoverAction = &pptxxml.HyperlinkSpec{
					RelID:          rid,
					Tooltip:        shape.HoverAction.Tooltip,
					HighlightClick: shape.HoverAction.HighlightClick,
					Action:         shape.HoverAction.Action.ActionType(),
				}
			}
		}
		specs = append(specs, spec)
	}
	return specs
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
