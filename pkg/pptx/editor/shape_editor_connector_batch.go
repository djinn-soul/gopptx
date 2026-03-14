package editor

import (
	"bytes"
	"errors"
	"fmt"
	"math"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

// AddConnectors inserts multiple connectors on one slide in a single XML rewrite.
func (e *PresentationEditor) AddConnectors(slideIndex int, connectors []common.ConnectorInsert) ([]int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}
	if len(connectors) == 0 {
		return []int{}, nil
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return nil, fmt.Errorf("read slide part %s: not found", partPath)
	}

	maxID := editorshape.MaxObjectID(content, cNvPrIDPattern, cNvPrSubmatchSize)
	shapeXML, shapeIDs, err := e.buildConnectorBatchXML(partPath, maxID, connectors)
	if err != nil {
		return nil, err
	}

	newXML, err := insertShapeXML(content, shapeXML)
	if err != nil {
		return nil, err
	}
	e.parts.Set(partPath, newXML)
	return shapeIDs, nil
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
