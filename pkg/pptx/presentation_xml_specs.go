package pptx

import (
	"github.com/djinn09/gopptx/internal/pptxxml"
)

func toXMLTableBorderSpec(border *TableCellBorder) *pptxxml.TableCellBorderSpec {
	if border == nil {
		return nil
	}
	return &pptxxml.TableCellBorderSpec{
		Width: border.widthEMU(),
		Color: normalizeHexColor(border.Color),
		Dash:  normalizeTableBorderDash(border.Dash),
	}
}

func toXMLTextRunRows(rows [][]TextRun, hyperlinkRIDs map[*Hyperlink]string) [][]pptxxml.TextRunSpec {
	if len(rows) == 0 {
		return nil
	}
	out := make([][]pptxxml.TextRunSpec, len(rows))
	for i := range rows {
		if len(rows[i]) == 0 {
			continue
		}
		runs := make([]pptxxml.TextRunSpec, 0, len(rows[i]))
		for _, run := range rows[i] {
			spec := pptxxml.TextRunSpec{
				Text:          run.Text,
				Bold:          run.Bold,
				Italic:        run.Italic,
				Underline:     run.Underline,
				Strikethrough: run.Strikethrough,
				Subscript:     run.Subscript,
				Superscript:   run.Superscript,
				Color:         normalizeHexColor(run.Color),
				Highlight:     normalizeHexColor(run.Highlight),
				Font:          run.Font,
				SizePt:        run.SizePt,
				Code:          run.Code,
			}
			if run.Hyperlink != nil {
				if rid, ok := hyperlinkRIDs[run.Hyperlink]; ok {
					spec.Hyperlink = &pptxxml.HyperlinkSpec{
						RelID:          rid,
						Tooltip:        run.Hyperlink.Tooltip,
						HighlightClick: run.Hyperlink.HighlightClick,
						Action:         run.Hyperlink.Action.ActionType(),
					}
				}
			}
			runs = append(runs, spec)
		}
		out[i] = runs
	}
	return out
}

func toXMLBulletParagraphStyles(styles []TextParagraphStyle) []pptxxml.BulletParagraphSpec {
	if len(styles) == 0 {
		return nil
	}
	out := make([]pptxxml.BulletParagraphSpec, len(styles))
	for i, style := range styles {
		out[i] = pptxxml.BulletParagraphSpec{
			Align:          normalizeTextAlign(style.Align),
			SpaceBeforePt:  style.SpaceBeforePt,
			SpaceAfterPt:   style.SpaceAfterPt,
			LineSpacingPct: style.LineSpacingPct,
			BulletStyle:    normalizeBulletStyle(style.BulletStyle),
			BulletChar:     style.BulletChar,
			Level:          style.Level,
		}
	}
	return out
}

func toXMLShapeSpecs(shapes []Shape, hyperlinkRIDs map[*Hyperlink]string) []pptxxml.ShapeSpec {
	if len(shapes) == 0 {
		return nil
	}
	specs := make([]pptxxml.ShapeSpec, 0, len(shapes))
	for _, shape := range shapes {
		spec := pptxxml.ShapeSpec{
			Type: normalizeShapeType(shape.Type),
			X:    shape.X,
			Y:    shape.Y,
			CX:   shape.CX,
			CY:   shape.CY,
			Text: shape.Text,
		}
		if shape.Fill != nil {
			spec.Fill = &pptxxml.ShapeFillSpec{
				Color:           normalizeHexColor(shape.Fill.Color),
				TransparencyPct: shape.Fill.TransparencyPct,
			}
		}
		if shape.GradientFill != nil {
			stops := make([]pptxxml.ShapeGradientStopSpec, 0, len(shape.GradientFill.Stops))
			for _, stop := range shape.GradientFill.Stops {
				stops = append(stops, pptxxml.ShapeGradientStopSpec{
					PositionPct:     stop.PositionPct,
					Color:           normalizeHexColor(stop.Color),
					TransparencyPct: stop.TransparencyPct,
				})
			}
			spec.GradientFill = &pptxxml.ShapeGradientFillSpec{
				Type:     normalizeShapeGradientType(shape.GradientFill.Type),
				Stops:    stops,
				AngleDeg: shape.GradientFill.AngleDeg,
			}
		}
		if shape.Line != nil {
			spec.Line = &pptxxml.ShapeLineSpec{
				Color: normalizeHexColor(shape.Line.Color),
				Width: shape.Line.Width,
				Dash:  normalizeDrawingLineDash(shape.Line.Dash),
			}
		}
		spec.RotationDeg = shape.RotationDeg
		if shape.Hyperlink != nil {
			if rid, ok := hyperlinkRIDs[shape.Hyperlink]; ok {
				spec.Hyperlink = &pptxxml.HyperlinkSpec{
					RelID:          rid,
					Tooltip:        shape.Hyperlink.Tooltip,
					HighlightClick: shape.Hyperlink.HighlightClick,
					Action:         shape.Hyperlink.Action.ActionType(),
				}
			}
		}
		specs = append(specs, spec)
	}
	return specs
}

func toXMLConnectorSpecs(connectors []Connector, shapes []Shape) []pptxxml.ConnectorSpec {
	if len(connectors) == 0 {
		return nil
	}
	specs := make([]pptxxml.ConnectorSpec, 0, len(connectors))
	for _, connector := range connectors {
		spec := pptxxml.ConnectorSpec{
			Type:   normalizeConnectorType(connector.Type),
			StartX: connector.StartX,
			StartY: connector.StartY,
			EndX:   connector.EndX,
			EndY:   connector.EndY,
			Line: pptxxml.ShapeLineSpec{
				Color: normalizeHexColor(connector.Line.Color),
				Width: connector.Line.Width,
				Dash:  normalizeDrawingLineDash(connector.Line.Dash),
			},
			StartArrow:      normalizeArrowType(connector.StartArrow),
			EndArrow:        normalizeArrowType(connector.EndArrow),
			ArrowSize:       normalizeArrowSize(connector.ArrowSize),
			StartShapeIndex: connector.StartShapeIndex,
			EndShapeIndex:   connector.EndShapeIndex,
			Label:           connector.Label,
		}
		spec.StartSiteIndex, spec.EndSiteIndex = resolveConnectorSiteIndices(connector, shapes)
		specs = append(specs, spec)
	}
	return specs
}
