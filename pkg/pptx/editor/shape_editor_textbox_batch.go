package editor

import (
	"bytes"
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

// AddTextboxes inserts multiple textboxes on one slide in a single XML rewrite.
func (e *PresentationEditor) AddTextboxes(slideIndex int, textboxes []common.TextboxInsert) ([]int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}
	if len(textboxes) == 0 {
		return []int{}, nil
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return nil, fmt.Errorf("read slide part %s: not found", partPath)
	}

	maxID := editorshape.MaxObjectID(content, cNvPrIDPattern, cNvPrSubmatchSize)
	shapeXML, shapeIDs, err := e.buildTextboxBatchXML(partPath, maxID, textboxes)
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

func (e *PresentationEditor) buildTextboxBatchXML(
	partPath string,
	startID int,
	textboxes []common.TextboxInsert,
) ([]byte, []int, error) {
	var xmlBuf bytes.Buffer
	shapeIDs := make([]int, 0, len(textboxes))

	for offset, textbox := range textboxes {
		shapeID := startID + offset + 1
		shape := parsedShape{
			ID:   shapeID,
			Name: fmt.Sprintf("rect %d", shapeID),
			Type: "rect",
			Text: textbox.Text,
			X:    int(textbox.Left),
			Y:    int(textbox.Top),
			W:    int(textbox.Width),
			H:    int(textbox.Height),
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
