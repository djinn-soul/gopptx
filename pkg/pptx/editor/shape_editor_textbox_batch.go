package editor

import (
	"bytes"
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// AddTextboxes inserts multiple textboxes on one slide in a single XML rewrite.
func (e *PresentationEditor) AddTextboxes(slideIndex int, textboxes []common.TextboxInsert) ([]int, error) {
	return e.insertBatchShapes(
		slideIndex,
		len(textboxes),
		func(partPath string, startID int) ([]byte, []int, error) {
			return e.buildTextboxBatchXML(partPath, startID, textboxes)
		},
	)
}

// ReserveShapeIDs returns the next available shape IDs on a slide without mutating XML.
func (e *PresentationEditor) ReserveShapeIDs(slideIndex int, count int) ([]int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}
	if count < 0 {
		return nil, errors.New("count must be non-negative")
	}
	if count == 0 {
		return []int{}, nil
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return nil, fmt.Errorf("read slide part %s: not found", partPath)
	}

	maxID := e.maxObjectID(partPath, content)
	// Reserved ids are handed to the caller, so they must never be handed out
	// again even though no XML is written here.
	e.reserveObjectIDs(partPath, maxID+count)
	shapeIDs := make([]int, 0, count)
	for offset := range count {
		shapeIDs = append(shapeIDs, maxID+offset+1)
	}
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
		if textbox.ShapeID != nil && *textbox.ShapeID > 0 {
			shapeID = *textbox.ShapeID
		}
		shape := parsedShape{
			ID:   shapeID,
			Name: fmt.Sprintf("rect %d", shapeID),
			Type: shapeTypeRect,
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
