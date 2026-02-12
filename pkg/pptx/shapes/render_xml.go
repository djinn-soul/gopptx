package shapes

import (
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
			X:            shape.X,
			Y:            shape.Y,
			CX:           shape.CX,
			CY:           shape.CY,
			Text:         shape.Text,
			AltText:      shape.AltText,
			IsDecorative: shape.IsDecorative,
		}
		if shape.Fill != nil {
			spec.Fill = &pptxxml.ShapeFillSpec{
				Color:           common.NormalizeHexColor(shape.Fill.Color),
				TransparencyPct: shape.Fill.TransparencyPct,
			}
		}
		if shape.GradientFill != nil {
			stops := make([]pptxxml.ShapeGradientStopSpec, 0, len(shape.GradientFill.Stops))
			for _, stop := range shape.GradientFill.Stops {
				stops = append(stops, pptxxml.ShapeGradientStopSpec{
					PositionPct:     stop.PositionPct,
					Color:           common.NormalizeHexColor(stop.Color),
					TransparencyPct: stop.TransparencyPct,
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
				Width: shape.Line.Width,
				Dash:  NormalizeDrawingLineDash(shape.Line.Dash),
			}
		}
		spec.RotationDeg = shape.RotationDeg
		if shape.TextFrame != nil {
			spec.TextFrame = &pptxxml.TextFrameSpec{
				MarginLeft:   shape.TextFrame.MarginLeft,
				MarginRight:  shape.TextFrame.MarginRight,
				MarginTop:    shape.TextFrame.MarginTop,
				MarginBottom: shape.TextFrame.MarginBottom,
				Anchor:       string(shape.TextFrame.Anchor),
				Wrap:         string(shape.TextFrame.Wrap),
				AutoFit:      string(shape.TextFrame.AutoFit),
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
		spec := pptxxml.ConnectorSpec{
			Type:   NormalizeConnectorType(connector.Type),
			StartX: connector.StartX,
			StartY: connector.StartY,
			EndX:   connector.EndX,
			EndY:   connector.EndY,
			Line: pptxxml.ShapeLineSpec{
				Color: common.NormalizeHexColor(connector.Line.Color),
				Width: connector.Line.Width,
				Dash:  NormalizeDrawingLineDash(connector.Line.Dash),
			},
			StartArrow:      NormalizeArrowType(connector.StartArrow),
			EndArrow:        NormalizeArrowType(connector.EndArrow),
			ArrowSize:       NormalizeArrowSize(connector.ArrowSize),
			StartShapeIndex: connector.StartShapeIndex,
			EndShapeIndex:   connector.EndShapeIndex,
			Label:           connector.Label,
			AltText:         connector.AltText,
			IsDecorative:    connector.IsDecorative,
		}
		spec.StartSiteIndex, spec.EndSiteIndex = ResolveConnectorSiteIndices(connector, shapes)
		specs = append(specs, spec)
	}
	return specs
}
