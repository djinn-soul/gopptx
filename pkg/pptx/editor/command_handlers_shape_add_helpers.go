package editor

import (
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type addShapeRequest struct {
	slideIndex  int
	shapeType   string
	x           float64
	y           float64
	w           float64
	h           float64
	text        string
	textFrame   *common.TextFrame
	clickAction *common.Hyperlink
	hoverAction *common.Hyperlink
	runs        []common.TextRun
}

func parseAddShapeRequest(
	e *PresentationEditor,
	payload map[string]any,
	v *PayloadValidator,
) (addShapeRequest, error) {
	slideIndex, ok := requireSlideIndex(e, payload, v)
	if !ok {
		return addShapeRequest{}, v.Error()
	}
	shapeType, ok := v.RequireString(payload, "type")
	if !ok {
		return addShapeRequest{}, v.Error()
	}
	x, ok := v.RequireFloat64(payload, "x")
	if !ok {
		return addShapeRequest{}, v.Error()
	}
	y, ok := v.RequireFloat64(payload, "y")
	if !ok {
		return addShapeRequest{}, v.Error()
	}
	w, ok := v.RequireFloat64(payload, "w")
	if !ok {
		return addShapeRequest{}, v.Error()
	}
	h, ok := v.RequireFloat64(payload, "h")
	if !ok {
		return addShapeRequest{}, v.Error()
	}

	request := addShapeRequest{
		slideIndex: slideIndex,
		shapeType:  shapeType,
		x:          x,
		y:          y,
		w:          w,
		h:          h,
		text:       v.OptionalString(payload, "text"),
	}

	if err := decodeOptionalPayloadValue(payload, "text_frame", &request.textFrame); err != nil {
		return addShapeRequest{}, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	if err := decodeOptionalPayloadValue(payload, "click_action", &request.clickAction); err != nil {
		return addShapeRequest{}, fmt.Errorf("invalid click_action: %w", err)
	}
	if err := decodeOptionalPayloadValue(payload, "hover_action", &request.hoverAction); err != nil {
		return addShapeRequest{}, fmt.Errorf("invalid hover_action: %w", err)
	}
	if err := decodeOptionalPayloadValue(payload, "runs", &request.runs); err != nil {
		return addShapeRequest{}, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	return request, nil
}

func buildShapeUpdateForAdd(request addShapeRequest) (common.ShapeUpdate, bool) {
	hasUpdates := request.text != "" ||
		len(request.runs) > 0 ||
		request.textFrame != nil ||
		request.clickAction != nil ||
		request.hoverAction != nil
	if !hasUpdates {
		return common.ShapeUpdate{}, false
	}

	updates := common.ShapeUpdate{}
	if request.text != "" {
		updates.Text = &request.text
	}
	if len(request.runs) > 0 {
		updates.Runs = &request.runs
	}
	if request.textFrame != nil {
		updates.TextFrame = request.textFrame
	}
	if request.clickAction != nil {
		updates.ClickAction = request.clickAction
	}
	if request.hoverAction != nil {
		updates.HoverAction = request.hoverAction
	}
	return updates, true
}
