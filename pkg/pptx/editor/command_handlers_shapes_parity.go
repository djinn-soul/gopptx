package editor

import (
	"encoding/json"
	"fmt"
	"math"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const minConnectorDimension = 1.0

func handleAddTextbox(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	left, ok := v.RequireFloat64(p, "left")
	if !ok {
		return nil, v.Error()
	}
	top, ok := v.RequireFloat64(p, "top")
	if !ok {
		return nil, v.Error()
	}
	width, ok := v.RequireFloat64(p, "width")
	if !ok {
		return nil, v.Error()
	}
	height, ok := v.RequireFloat64(p, "height")
	if !ok {
		return nil, v.Error()
	}

	addPayload := map[string]any{
		"slide_index": slideIndex,
		"type":        "rect",
		"x":           left,
		"y":           top,
		"w":           width,
		"h":           height,
	}
	copyShapeUpdateFields(p, addPayload)
	return addShapeFromPayload(e, addPayload)
}

func handleAddConnector(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	connectorType, ok := v.RequireString(p, "connector_type")
	if !ok {
		return nil, v.Error()
	}
	beginX, ok := v.RequireFloat64(p, "begin_x")
	if !ok {
		return nil, v.Error()
	}
	beginY, ok := v.RequireFloat64(p, "begin_y")
	if !ok {
		return nil, v.Error()
	}
	endX, ok := v.RequireFloat64(p, "end_x")
	if !ok {
		return nil, v.Error()
	}
	endY, ok := v.RequireFloat64(p, "end_y")
	if !ok {
		return nil, v.Error()
	}

	left := math.Min(beginX, endX)
	top := math.Min(beginY, endY)
	width := math.Max(math.Abs(endX-beginX), minConnectorDimension)
	height := math.Max(math.Abs(endY-beginY), minConnectorDimension)

	addPayload := map[string]any{
		"slide_index": slideIndex,
		"type":        connectorType,
		"x":           left,
		"y":           top,
		"w":           width,
		"h":           height,
	}
	copyShapeUpdateFields(p, addPayload)
	return addShapeFromPayload(e, addPayload)
}

func handleAddGroupShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	shapeIDs := []int{}
	if _, present := p["shapes"]; present {
		parsedIDs, idsOK := v.RequireIntSlice(p, "shapes")
		if !idsOK {
			return nil, v.Error()
		}
		shapeIDs = parsedIDs
	}
	newID, err := e.AddGroupShape(slideIndex, shapeIDs)
	if err != nil {
		return nil, err
	}
	return map[string]int{"shape_id": newID}, nil
}

func handleBuildFreeform(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	points, err := parseFreeformPointsPayload(p)
	if err != nil {
		return nil, err
	}
	closePath := false
	if rawClose, exists := p["close"]; exists {
		flag, castOK := rawClose.(bool)
		if !castOK {
			return nil, NewBridgeError(
				ErrCodeInvalidType,
				"field close must be a boolean",
			)
		}
		closePath = flag
	}
	shapeID, err := e.AddFreeformShape(slideIndex, points, closePath)
	if err != nil {
		return nil, err
	}
	if updates, hasUpdates, updateErr := parseOptionalShapeUpdates(p); updateErr != nil {
		return nil, updateErr
	} else if hasUpdates {
		if err := e.UpdateShape(slideIndex, shapeID, updates); err != nil {
			return nil, err
		}
	}
	return map[string]int{"shape_id": shapeID}, nil
}

func parseFreeformPointsPayload(payload map[string]any) ([]freeformPoint, error) {
	rawPoints, exists := payload["points"]
	if !exists {
		return nil, NewBridgeError(ErrCodeMissingField, "missing required field: points")
	}
	pointItems, ok := rawPoints.([]any)
	if !ok {
		return nil, NewBridgeError(ErrCodeInvalidType, "field points must be an array")
	}
	if len(pointItems) < 2 {
		return nil, NewBridgeError(ErrCodeInvalidValue, "field points must contain at least 2 points")
	}

	points := make([]freeformPoint, 0, len(pointItems))
	for i, item := range pointItems {
		pair, pairOK := item.([]any)
		if !pairOK || len(pair) != 2 {
			return nil, NewBridgeError(
				ErrCodeInvalidType,
				fmt.Sprintf("field points[%d] must be [x, y]", i),
			)
		}
		x, xOK := numberToFloat64(pair[0])
		y, yOK := numberToFloat64(pair[1])
		if !xOK || !yOK {
			return nil, NewBridgeError(
				ErrCodeInvalidType,
				fmt.Sprintf("field points[%d] coordinates must be numbers", i),
			)
		}
		points = append(points, freeformPoint{X: int(x), Y: int(y)})
	}
	return points, nil
}

func numberToFloat64(value any) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	default:
		return 0, false
	}
}

func parseOptionalShapeUpdates(payload map[string]any) (common.ShapeUpdate, bool, error) {
	text := ""
	runs := []common.TextRun(nil)
	var textFrame *common.TextFrame
	var paragraph *common.Paragraph
	var clickAction *common.Hyperlink
	var hoverAction *common.Hyperlink
	properties := common.ShapeUpdate{}

	if value, ok := payload["text"]; ok {
		textValue, isString := value.(string)
		if !isString {
			return common.ShapeUpdate{}, false, NewBridgeError(
				ErrCodeInvalidType,
				"field text must be a string",
			)
		}
		text = textValue
	}
	if err := decodeOptionalPayloadValue(payload, "runs", &runs); err != nil {
		return common.ShapeUpdate{}, false, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	if err := decodeOptionalPayloadValue(payload, "text_frame", &textFrame); err != nil {
		return common.ShapeUpdate{}, false, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	if err := decodeOptionalPayloadValue(payload, "paragraph", &paragraph); err != nil {
		return common.ShapeUpdate{}, false, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	if err := decodeOptionalPayloadValue(payload, "click_action", &clickAction); err != nil {
		return common.ShapeUpdate{}, false, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	if err := decodeOptionalPayloadValue(payload, "hover_action", &hoverAction); err != nil {
		return common.ShapeUpdate{}, false, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	if err := decodeOptionalPayloadValue(payload, "properties", &properties); err != nil {
		return common.ShapeUpdate{}, false, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}

	hasExplicitUpdates := text != "" || len(runs) > 0 || textFrame != nil || paragraph != nil || clickAction != nil ||
		hoverAction != nil
	hasProperties := hasAnyUpdate(properties)
	if !hasExplicitUpdates && !hasProperties {
		return common.ShapeUpdate{}, false, nil
	}

	updates := properties
	if text != "" {
		updates.Text = &text
	}
	if len(runs) > 0 {
		updates.Runs = &runs
	}
	if textFrame != nil {
		updates.TextFrame = textFrame
	}
	if paragraph != nil {
		updates.Paragraph = paragraph
	}
	if clickAction != nil {
		updates.ClickAction = clickAction
	}
	if hoverAction != nil {
		updates.HoverAction = hoverAction
	}
	return updates, true, nil
}

func addShapeFromPayload(e *PresentationEditor, payload map[string]any) (any, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return handleAddShape(e, raw)
}

func copyShapeUpdateFields(src, dst map[string]any) {
	for _, key := range []string{
		"text",
		"runs",
		"text_frame",
		"paragraph",
		"click_action",
		"hover_action",
		"properties",
	} {
		if value, ok := src[key]; ok {
			dst[key] = value
		}
	}
}
