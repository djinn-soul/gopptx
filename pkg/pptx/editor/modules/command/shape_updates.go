package command

import common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"

func ParseOptionalShapeUpdates(payload map[string]any) (common.ShapeUpdate, bool, error) {
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
			return common.ShapeUpdate{}, false, NewValidationError("INVALID_TYPE", "field text must be a string")
		}
		text = textValue
	}
	if err := DecodeOptionalPayloadValue(payload, "runs", &runs); err != nil {
		return common.ShapeUpdate{}, false, NewValidationError("INVALID_PAYLOAD", err.Error())
	}
	if err := DecodeOptionalPayloadValue(payload, "text_frame", &textFrame); err != nil {
		return common.ShapeUpdate{}, false, NewValidationError("INVALID_PAYLOAD", err.Error())
	}
	if err := DecodeOptionalPayloadValue(payload, "paragraph", &paragraph); err != nil {
		return common.ShapeUpdate{}, false, NewValidationError("INVALID_PAYLOAD", err.Error())
	}
	if err := DecodeOptionalPayloadValue(payload, "click_action", &clickAction); err != nil {
		return common.ShapeUpdate{}, false, NewValidationError("INVALID_PAYLOAD", err.Error())
	}
	if err := DecodeOptionalPayloadValue(payload, "hover_action", &hoverAction); err != nil {
		return common.ShapeUpdate{}, false, NewValidationError("INVALID_PAYLOAD", err.Error())
	}
	if err := DecodeOptionalPayloadValue(payload, "properties", &properties); err != nil {
		return common.ShapeUpdate{}, false, NewValidationError("INVALID_PAYLOAD", err.Error())
	}

	hasExplicitUpdates := text != "" || len(runs) > 0 || textFrame != nil || paragraph != nil || clickAction != nil ||
		hoverAction != nil
	hasProperties := HasAnyUpdate(properties)
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

func CopyShapeUpdateFields(src, dst map[string]any) {
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
