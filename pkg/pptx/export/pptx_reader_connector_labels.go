package export

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	connectorLabelNamePrefix = "connector label "
	connectorLabelWidthEMU   = 914400
	connectorLabelHeightEMU  = 228600
)

func foldGeneratedConnectorLabels(sc *elements.SlideContent) {
	if sc == nil || len(sc.Connectors) == 0 || len(sc.Shapes) == 0 {
		return
	}

	filtered := make([]shapes.Shape, 0, len(sc.Shapes))
	for _, shape := range sc.Shapes {
		if !isGeneratedConnectorLabelShape(shape) {
			filtered = append(filtered, shape)
			continue
		}
		connectorIndex, ok := connectorIndexForGeneratedLabel(sc.Connectors, shape)
		if !ok {
			filtered = append(filtered, shape)
			continue
		}
		label := strings.TrimSpace(shape.Text)
		if label == "" {
			filtered = append(filtered, shape)
			continue
		}
		if strings.TrimSpace(sc.Connectors[connectorIndex].Label) != "" &&
			strings.TrimSpace(sc.Connectors[connectorIndex].Label) != label {
			filtered = append(filtered, shape)
			continue
		}
		sc.Connectors[connectorIndex].Label = label
	}
	sc.Shapes = filtered
}

func isGeneratedConnectorLabelShape(shape shapes.Shape) bool {
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(shape.Name)), connectorLabelNamePrefix) {
		return false
	}
	if strings.TrimSpace(shape.Type) != shapes.ShapeTypeRectangle {
		return false
	}
	if strings.TrimSpace(shape.Text) == "" || strings.TrimSpace(shape.AltText) != "" {
		return false
	}
	if shape.Line != nil || shape.RichLine != nil {
		return false
	}
	if shape.ClickAction != nil || shape.HoverAction != nil || shape.Hyperlink != nil {
		return false
	}
	return shape.CX == styling.Emu(connectorLabelWidthEMU) && shape.CY == styling.Emu(connectorLabelHeightEMU)
}

func connectorIndexForGeneratedLabel(connectors []shapes.Connector, shape shapes.Shape) (int, bool) {
	centerX := shape.X + shape.CX/2
	centerY := shape.Y + shape.CY/2
	matchIndex := -1
	for i, connector := range connectors {
		if connectorMidpointX(connector) != centerX || connectorMidpointY(connector) != centerY {
			continue
		}
		if matchIndex >= 0 {
			return 0, false
		}
		matchIndex = i
	}
	if matchIndex < 0 {
		return 0, false
	}
	return matchIndex, true
}

func connectorMidpointX(connector shapes.Connector) styling.Length {
	return (connector.StartX + connector.EndX) / 2
}

func connectorMidpointY(connector shapes.Connector) styling.Length {
	return (connector.StartY + connector.EndY) / 2
}
