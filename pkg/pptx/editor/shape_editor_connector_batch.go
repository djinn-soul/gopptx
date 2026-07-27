package editor

import (
	"bytes"
	"fmt"
	"math"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// AddConnectors inserts multiple connectors on one slide in a single XML rewrite.
func (e *PresentationEditor) AddConnectors(slideIndex int, connectors []common.ConnectorInsert) ([]int, error) {
	return e.insertBatchShapes(
		slideIndex,
		len(connectors),
		func(partPath string, startID int) ([]byte, []int, error) {
			return e.buildConnectorBatchXML(partPath, startID, connectors)
		},
	)
}

func (e *PresentationEditor) buildConnectorBatchXML(
	partPath string,
	startID int,
	connectors []common.ConnectorInsert,
) ([]byte, []int, error) {
	var xmlBuf bytes.Buffer
	shapeIDs := make([]int, 0, len(connectors))

	for offset, connector := range connectors {
		shapeID := startID + offset + 1
		if connector.ShapeID != nil && *connector.ShapeID > 0 {
			shapeID = *connector.ShapeID
		}
		left := math.Min(connector.BeginX, connector.EndX)
		top := math.Min(connector.BeginY, connector.EndY)
		width := math.Max(math.Abs(connector.EndX-connector.BeginX), minConnectorDimension)
		height := math.Max(math.Abs(connector.EndY-connector.BeginY), minConnectorDimension)

		shape := parsedShape{
			ID:          shapeID,
			Name:        fmt.Sprintf("%s %d", connector.ConnectorType, shapeID),
			Type:        connector.ConnectorType,
			TextFrame:   connector.TextFrame,
			Paragraph:   connector.Paragraph,
			Fill:        connector.Fill,
			Line:        connector.Line,
			Shadow:      connector.Shadow,
			Glow:        connector.Glow,
			Blur:        connector.Blur,
			SoftEdge:    connector.SoftEdge,
			Reflection:  connector.Reflection,
			X:           int(left),
			Y:           int(top),
			W:           int(width),
			H:           int(height),
			ClickAction: connector.ClickAction,
			HoverAction: connector.HoverAction,
		}
		if connector.Text != nil {
			shape.Text = *connector.Text
		}
		if connector.Runs != nil {
			shape.Runs = *connector.Runs
		}

		shapeNode, err := e.renderShapeXML(partPath, &shape)
		if err != nil {
			return nil, nil, err
		}
		xmlBuf.Write(shapeNode)
		shapeIDs = append(shapeIDs, shapeID)
	}

	return xmlBuf.Bytes(), shapeIDs, nil
}
