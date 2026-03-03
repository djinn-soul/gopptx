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
	paragraph   *common.Paragraph
	clickAction *common.Hyperlink
	hoverAction *common.Hyperlink
	runs        []common.TextRun
	properties  common.ShapeUpdate
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
	if err := decodeOptionalPayloadValue(payload, "paragraph", &request.paragraph); err != nil {
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
	if err := decodeOptionalPayloadValue(payload, "properties", &request.properties); err != nil {
		return addShapeRequest{}, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	return request, nil
}

func buildShapeUpdateForAdd(request addShapeRequest) (common.ShapeUpdate, bool) {
	hasExplicitUpdates := request.text != "" ||
		len(request.runs) > 0 ||
		request.textFrame != nil ||
		request.paragraph != nil ||
		request.clickAction != nil ||
		request.hoverAction != nil
	hasProperties := hasAnyUpdate(request.properties)

	if !hasExplicitUpdates && !hasProperties {
		return common.ShapeUpdate{}, false
	}

	updates := request.properties
	if request.text != "" {
		updates.Text = &request.text
	}
	if len(request.runs) > 0 {
		updates.Runs = &request.runs
	}
	if request.textFrame != nil {
		updates.TextFrame = request.textFrame
	}
	if request.paragraph != nil {
		updates.Paragraph = request.paragraph
	}
	if request.clickAction != nil {
		updates.ClickAction = request.clickAction
	}
	if request.hoverAction != nil {
		updates.HoverAction = request.hoverAction
	}
	return updates, true
}

func hasAnyUpdate(u common.ShapeUpdate) bool {
	return u.Text != nil || u.Runs != nil || u.TextFrame != nil ||
		u.Paragraph != nil || u.Fill != nil || u.Line != nil || u.Shadow != nil || u.Glow != nil || u.Blur != nil || u.SoftEdge != nil || u.Reflection != nil ||
		u.ClickAction != nil || u.HoverAction != nil || u.X != nil ||
		u.Y != nil || u.W != nil || u.H != nil || u.Rotation != nil ||
		u.FlipH != nil || u.FlipV != nil || u.Crop != nil
}
