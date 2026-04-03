package export

import (
	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	connSiteTop         = 0
	connSiteRight       = 1
	connSiteBottom      = 2
	connSiteLeft        = 3
	connSiteTopLeft     = 4
	connSiteTopRight    = 5
	connSiteBottomRight = 6
	connSiteBottomLeft  = 7
	connSiteCenter      = 8
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
	connector = applyConnectorLineOptions(connector, es.Line)
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

func applyConnectorLineOptions(connector shapes.Connector, line *editorcommon.ShapeLine) shapes.Connector {
	if line == nil {
		return connector
	}
	connector.Line = editorConnectorLine(line, connector.Line)
	if line.StartArrow != nil {
		connector.StartArrow = shapes.NormalizeArrowType(*line.StartArrow)
	}
	if line.StartArrowWidth != nil {
		connector.StartArrowWidth = shapes.NormalizeArrowSize(*line.StartArrowWidth)
	}
	if line.StartArrowLength != nil {
		connector.StartArrowLen = shapes.NormalizeArrowSize(*line.StartArrowLength)
	}
	if line.EndArrow != nil {
		connector.EndArrow = shapes.NormalizeArrowType(*line.EndArrow)
	}
	if line.EndArrowWidth != nil {
		connector.EndArrowWidth = shapes.NormalizeArrowSize(*line.EndArrowWidth)
	}
	if line.EndArrowLength != nil {
		connector.EndArrowLen = shapes.NormalizeArrowSize(*line.EndArrowLength)
	}
	return connector
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
	case connSiteTop:
		return shapes.ConnectionSiteTop
	case connSiteRight:
		return shapes.ConnectionSiteRight
	case connSiteBottom:
		return shapes.ConnectionSiteBottom
	case connSiteLeft:
		return shapes.ConnectionSiteLeft
	case connSiteTopLeft:
		return shapes.ConnectionSiteTopLeft
	case connSiteTopRight:
		return shapes.ConnectionSiteTopRight
	case connSiteBottomRight:
		return shapes.ConnectionSiteBottomRight
	case connSiteBottomLeft:
		return shapes.ConnectionSiteBottomLeft
	case connSiteCenter:
		return shapes.ConnectionSiteCenter
	default:
		return ""
	}
}
