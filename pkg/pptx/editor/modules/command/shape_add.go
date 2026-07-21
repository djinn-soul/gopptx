package command

import (
	"fmt"
	"reflect"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type AddShapeRequest struct {
	SlideIndex  int
	ShapeType   string
	X           float64
	Y           float64
	W           float64
	H           float64
	Text        string
	TextFrame   *common.TextFrame
	Paragraph   *common.Paragraph
	ClickAction *common.Hyperlink
	HoverAction *common.Hyperlink
	Runs        []common.TextRun
	Properties  common.ShapeUpdate
}

type AddShapeRequireSlideIndexFn func(map[string]any) (int, bool)
type AddShapeRequireStringFn func(map[string]any, string) (string, bool)
type AddShapeRequireFloatFn func(map[string]any, string) (float64, bool)
type AddShapeOptionalStringFn func(map[string]any, string) string
type AddShapeCreateFn func(slideIndex int, shapeType string, x, y, w, h float64) (int, error)
type AddShapeUpdateFn func(slideIndex int, shapeID int, update common.ShapeUpdate) error

func ParseAddShapeBase(
	payload map[string]any,
	requireSlideIndex AddShapeRequireSlideIndexFn,
	requireString AddShapeRequireStringFn,
	requireFloat AddShapeRequireFloatFn,
	optionalString AddShapeOptionalStringFn,
) (AddShapeRequest, bool) {
	slideIndex, ok := requireSlideIndex(payload)
	if !ok {
		return AddShapeRequest{}, false
	}
	shapeType, ok := requireString(payload, "type")
	if !ok {
		return AddShapeRequest{}, false
	}
	x, ok := requireFloat(payload, "x")
	if !ok {
		return AddShapeRequest{}, false
	}
	y, ok := requireFloat(payload, "y")
	if !ok {
		return AddShapeRequest{}, false
	}
	w, ok := requireFloat(payload, "w")
	if !ok {
		return AddShapeRequest{}, false
	}
	h, ok := requireFloat(payload, "h")
	if !ok {
		return AddShapeRequest{}, false
	}
	return AddShapeRequest{
		SlideIndex: slideIndex,
		ShapeType:  shapeType,
		X:          x,
		Y:          y,
		W:          w,
		H:          h,
		Text:       optionalString(payload, "text"),
	}, true
}

func BuildShapeUpdateForAdd(request AddShapeRequest) (common.ShapeUpdate, bool) {
	hasExplicitUpdates := request.Text != "" ||
		len(request.Runs) > 0 ||
		request.TextFrame != nil ||
		request.Paragraph != nil ||
		request.ClickAction != nil ||
		request.HoverAction != nil
	hasProperties := HasAnyUpdate(request.Properties)

	if !hasExplicitUpdates && !hasProperties {
		return common.ShapeUpdate{}, false
	}

	updates := request.Properties
	if request.Text != "" {
		updates.Text = &request.Text
	}
	if len(request.Runs) > 0 {
		updates.Runs = &request.Runs
	}
	if request.TextFrame != nil {
		updates.TextFrame = request.TextFrame
	}
	if request.Paragraph != nil {
		updates.Paragraph = request.Paragraph
	}
	if request.ClickAction != nil {
		updates.ClickAction = request.ClickAction
	}
	if request.HoverAction != nil {
		updates.HoverAction = request.HoverAction
	}
	return updates, true
}

func HasAnyUpdate(u common.ShapeUpdate) bool {
	v := reflect.ValueOf(u)
	for i := range v.NumField() {
		field := v.Field(i)
		if field.Kind() == reflect.Pointer && !field.IsNil() {
			return true
		}
	}
	return false
}

func DecodeAddShapeOptionals(payload map[string]any, request *AddShapeRequest) error {
	if err := DecodeOptionalPayloadValue(payload, "text_frame", &request.TextFrame); err != nil {
		return err
	}
	if err := DecodeOptionalPayloadValue(payload, "paragraph", &request.Paragraph); err != nil {
		return err
	}
	if err := DecodeOptionalPayloadValue(payload, "click_action", &request.ClickAction); err != nil {
		return fmt.Errorf("invalid click_action: %w", err)
	}
	if err := DecodeOptionalPayloadValue(payload, "hover_action", &request.HoverAction); err != nil {
		return fmt.Errorf("invalid hover_action: %w", err)
	}
	if err := DecodeOptionalPayloadValue(payload, "runs", &request.Runs); err != nil {
		return err
	}
	if err := DecodeOptionalPayloadValue(payload, "properties", &request.Properties); err != nil {
		return err
	}
	return nil
}

func ExecuteAddShapeRequest(
	request AddShapeRequest,
	createShape AddShapeCreateFn,
	updateShape AddShapeUpdateFn,
) (int, error) {
	shapeID, err := createShape(
		request.SlideIndex,
		request.ShapeType,
		request.X,
		request.Y,
		request.W,
		request.H,
	)
	if err != nil {
		return 0, err
	}
	if updates, hasUpdates := BuildShapeUpdateForAdd(request); hasUpdates {
		if err := updateShape(request.SlideIndex, shapeID, updates); err != nil {
			return 0, err
		}
	}
	return shapeID, nil
}
