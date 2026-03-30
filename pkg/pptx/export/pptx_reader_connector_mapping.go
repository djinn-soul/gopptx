package export

import (
	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func isEditorConnector(es editorcommon.Shape) bool {
	return es.Connector != nil || shapes.IsConnectorType(es.Type)
}

func editorShapeToConnector(es editorcommon.Shape, shapeIndexByID map[int]int) (shapes.Connector, bool) {
	if !isEditorConnector(es) {
		return shapes.Connector{}, false
	}
	startX, startY, endX, endY := editorConnectorEndpoints(es)
	connector := shapes.NewConnector(es.Type, startX, startY, endX, endY)
	connector.AltText = es.AltText
	connector.IsDecorative = es.IsDecorative
	connector.Adjustments = editorAdjustmentsToExportConnector(es.Adjustments)
	if es.Line != nil {
		connector.Line = editorConnectorLine(es.Line, connector.Line)
		if es.Line.StartArrow != nil {
			connector.StartArrow = shapes.NormalizeArrowType(*es.Line.StartArrow)
		}
		if es.Line.StartArrowWidth != nil {
			connector.StartArrowWidth = shapes.NormalizeArrowSize(*es.Line.StartArrowWidth)
		}
		if es.Line.StartArrowLength != nil {
			connector.StartArrowLen = shapes.NormalizeArrowSize(*es.Line.StartArrowLength)
		}
		if es.Line.EndArrow != nil {
			connector.EndArrow = shapes.NormalizeArrowType(*es.Line.EndArrow)
		}
		if es.Line.EndArrowWidth != nil {
			connector.EndArrowWidth = shapes.NormalizeArrowSize(*es.Line.EndArrowWidth)
		}
		if es.Line.EndArrowLength != nil {
			connector.EndArrowLen = shapes.NormalizeArrowSize(*es.Line.EndArrowLength)
		}
	}
	if es.Connector != nil {
		if es.Connector.StartShapeID != nil {
			connector.StartShapeIndex = shapeIndexByID[*es.Connector.StartShapeID]
		}
		if es.Connector.StartSiteIndex != nil {
			connector.StartSite = connectorSiteName(*es.Connector.StartSiteIndex)
		}
		if es.Connector.EndShapeID != nil {
			connector.EndShapeIndex = shapeIndexByID[*es.Connector.EndShapeID]
		}
		if es.Connector.EndSiteIndex != nil {
			connector.EndSite = connectorSiteName(*es.Connector.EndSiteIndex)
		}
	}
	return connector, true
}

func editorConnectorEndpoints(es editorcommon.Shape) (
	styling.Length,
	styling.Length,
	styling.Length,
	styling.Length,
) {
	left := styling.Emu(int64(es.X))
	top := styling.Emu(int64(es.Y))
	right := styling.Emu(int64(es.X + es.W))
	bottom := styling.Emu(int64(es.Y + es.H))
	startX, endX := left, right
	startY, endY := top, bottom
	if es.Connector != nil && es.Connector.FlipH {
		startX, endX = right, left
	}
	if es.Connector != nil && es.Connector.FlipV {
		startY, endY = bottom, top
	}
	return startX, startY, endX, endY
}

func editorConnectorLine(src *editorcommon.ShapeLine, fallback shapes.ShapeLine) shapes.ShapeLine {
	line := fallback
	if src.Color != nil && *src.Color != "" {
		line.Color = *src.Color
	}
	if src.WidthEmu != nil && *src.WidthEmu > 0 {
		line.Width = styling.Emu(int64(*src.WidthEmu))
	}
	if src.DashStyle != nil {
		line.Dash = shapes.NormalizeDrawingLineDash(*src.DashStyle)
	}
	return line
}

func editorAdjustmentsToExportConnector(
	adjustments []editorcommon.ShapeAdjustment,
) []shapes.ConnectorAdjustment {
	if len(adjustments) == 0 {
		return nil
	}
	out := make([]shapes.ConnectorAdjustment, 0, len(adjustments))
	for _, adjustment := range adjustments {
		if adjustment.Name == "" || adjustment.Formula == "" {
			continue
		}
		out = append(out, shapes.ConnectorAdjustment{
			Name:    adjustment.Name,
			Formula: adjustment.Formula,
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func connectorSiteName(idx int) string {
	switch idx {
	case 0:
		return shapes.ConnectionSiteTop
	case 1:
		return shapes.ConnectionSiteRight
	case 2:
		return shapes.ConnectionSiteBottom
	case 3:
		return shapes.ConnectionSiteLeft
	case 4:
		return shapes.ConnectionSiteTopLeft
	case 5:
		return shapes.ConnectionSiteTopRight
	case 6:
		return shapes.ConnectionSiteBottomRight
	case 7:
		return shapes.ConnectionSiteBottomLeft
	case 8:
		return shapes.ConnectionSiteCenter
	default:
		return ""
	}
}
