package editor

import (
	"encoding/json"
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	slidesmeta "github.com/djinn-soul/gopptx/pkg/pptx/editor/handlers/slidesmeta"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

func handleSetSlideSize(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.SlideSizeRequest, bool) {
			return editorcommand.ParseSlideSizeRequest(p, v.RequireInt64)
		},
		v.Error,
		func(request editorcommand.SlideSizeRequest) (any, error) {
			if err := e.SetSlideSize(common.SlideSize{Width: request.Width, Height: request.Height}); err != nil {
				return nil, err
			}
			return respUpdated, nil
		},
	)
}

func handleSetSlideTitle(e *PresentationEditor, payload json.RawMessage) (any, error) {
	request, err := parseSetSlideTitlePayload(payload, e.SlideCount())
	if err != nil {
		return nil, err
	}
	if err := e.SetSlideTitle(request.SlideIndex, request.Title); err != nil {
		return nil, err
	}
	return respUpdated, nil
}

func parseSetSlideTitlePayload(payload json.RawMessage, slideCount int) (editorcommand.SlideTitleRequest, error) {
	if len(payload) == 0 {
		return editorcommand.SlideTitleRequest{}, NewBridgeError(ErrCodeInvalidPayload, "empty payload")
	}

	var fast struct {
		SlideIndex *float64 `json:"slide_index"`
		Title      *string  `json:"title"`
	}
	if err := json.Unmarshal(payload, &fast); err == nil {
		if fast.SlideIndex == nil {
			return editorcommand.SlideTitleRequest{}, newPayloadValidationBridgeError(
				ErrCodeMissingField,
				"missing required field: slide_index",
			)
		}
		if fast.Title == nil {
			return editorcommand.SlideTitleRequest{}, newPayloadValidationBridgeError(
				ErrCodeMissingField,
				"missing required field: title",
			)
		}
		slideIndex := int(*fast.SlideIndex)
		if slideIndex < 0 || slideIndex >= slideCount {
			msg := fmt.Sprintf("slide_index %d out of bounds [%d, %d)", slideIndex, 0, slideCount)
			return editorcommand.SlideTitleRequest{}, newPayloadValidationBridgeError(ErrCodeInvalidIndex, msg)
		}
		return editorcommand.SlideTitleRequest{SlideIndex: slideIndex, Title: *fast.Title}, nil
	}
	return parseSetSlideTitlePayloadSlow(payload, slideCount)
}

func parseSetSlideTitlePayloadSlow(payload json.RawMessage, slideCount int) (editorcommand.SlideTitleRequest, error) {
	var raw struct {
		SlideIndex json.RawMessage `json:"slide_index"`
		Title      json.RawMessage `json:"title"`
	}
	if err := json.Unmarshal(payload, &raw); err != nil {
		msg := fmt.Sprintf("invalid JSON payload: %v", err)
		return editorcommand.SlideTitleRequest{}, NewBridgeError(ErrCodeInvalidPayload, msg)
	}
	if len(raw.SlideIndex) == 0 {
		return editorcommand.SlideTitleRequest{}, newPayloadValidationBridgeError(
			ErrCodeMissingField,
			"missing required field: slide_index",
		)
	}
	if len(raw.Title) == 0 {
		return editorcommand.SlideTitleRequest{}, newPayloadValidationBridgeError(
			ErrCodeMissingField,
			"missing required field: title",
		)
	}

	var slideNumber float64
	if err := json.Unmarshal(raw.SlideIndex, &slideNumber); err != nil {
		return editorcommand.SlideTitleRequest{}, newPayloadValidationBridgeError(
			ErrCodeInvalidType,
			"field slide_index must be an integer",
		)
	}
	slideIndex := int(slideNumber)
	if slideIndex < 0 || slideIndex >= slideCount {
		msg := fmt.Sprintf("slide_index %d out of bounds [%d, %d)", slideIndex, 0, slideCount)
		return editorcommand.SlideTitleRequest{}, newPayloadValidationBridgeError(ErrCodeInvalidIndex, msg)
	}

	var title string
	if err := json.Unmarshal(raw.Title, &title); err != nil {
		return editorcommand.SlideTitleRequest{}, newPayloadValidationBridgeError(
			ErrCodeInvalidType,
			"field title must be a string",
		)
	}
	return editorcommand.SlideTitleRequest{SlideIndex: slideIndex, Title: title}, nil
}

func newPayloadValidationBridgeError(code, detail string) error {
	return NewBridgeErrorWithDetails(code, "payload validation failed", []string{detail})
}

// handleSetSlideHidden marks or unmarks a slide as hidden.
//
// Payload: {"slide_index": N, "hidden": bool}.
// Response: {"updated": true}.
func handleSetSlideHidden(e *PresentationEditor, payload json.RawMessage) (any, error) {
	slideIndex, hidden, err := parseSetSlideHiddenPayload(payload, e.SlideCount())
	if err != nil {
		return nil, err
	}
	if setErr := e.SetSlideHidden(slideIndex, hidden); setErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, setErr.Error())
	}
	return respUpdated, nil
}

func parseSetSlideHiddenPayload(payload json.RawMessage, slideCount int) (int, bool, error) {
	if len(payload) == 0 {
		return 0, false, NewBridgeError(ErrCodeInvalidPayload, "empty payload")
	}
	var fast struct {
		SlideIndex *float64 `json:"slide_index"`
		Hidden     *bool    `json:"hidden"`
	}
	if err := json.Unmarshal(payload, &fast); err == nil {
		if fast.SlideIndex == nil {
			return 0, false, newPayloadValidationBridgeError(ErrCodeMissingField, "missing required field: slide_index")
		}
		if fast.Hidden == nil {
			return 0, false, newPayloadValidationBridgeError(ErrCodeMissingField, "missing required field: hidden")
		}
		slideIndex := int(*fast.SlideIndex)
		if slideIndex < 0 || slideIndex >= slideCount {
			msg := fmt.Sprintf("slide_index %d out of bounds [%d, %d)", slideIndex, 0, slideCount)
			return 0, false, newPayloadValidationBridgeError(ErrCodeInvalidIndex, msg)
		}
		return slideIndex, *fast.Hidden, nil
	}
	return parseSetSlideHiddenPayloadSlow(payload, slideCount)
}

func parseSetSlideHiddenPayloadSlow(payload json.RawMessage, slideCount int) (int, bool, error) {
	var raw struct {
		SlideIndex json.RawMessage `json:"slide_index"`
		Hidden     json.RawMessage `json:"hidden"`
	}
	if err := json.Unmarshal(payload, &raw); err != nil {
		msg := fmt.Sprintf("invalid JSON payload: %v", err)
		return 0, false, NewBridgeError(ErrCodeInvalidPayload, msg)
	}
	if len(raw.SlideIndex) == 0 {
		return 0, false, newPayloadValidationBridgeError(ErrCodeMissingField, "missing required field: slide_index")
	}
	if len(raw.Hidden) == 0 {
		return 0, false, newPayloadValidationBridgeError(ErrCodeMissingField, "missing required field: hidden")
	}

	var slideNumber float64
	if err := json.Unmarshal(raw.SlideIndex, &slideNumber); err != nil {
		return 0, false, newPayloadValidationBridgeError(ErrCodeInvalidType, "field slide_index must be an integer")
	}
	slideIndex := int(slideNumber)
	if slideIndex < 0 || slideIndex >= slideCount {
		msg := fmt.Sprintf("slide_index %d out of bounds [%d, %d)", slideIndex, 0, slideCount)
		return 0, false, newPayloadValidationBridgeError(ErrCodeInvalidIndex, msg)
	}

	var hidden bool
	if err := json.Unmarshal(raw.Hidden, &hidden); err != nil {
		return 0, false, newPayloadValidationBridgeError(ErrCodeInvalidType, "field hidden must be a boolean")
	}
	return slideIndex, hidden, nil
}

func handleMergeFromFile(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.MergeFromFileRequest, bool) {
			return editorcommand.ParseMergeFromFileRequest(p, v.RequireString)
		},
		v.Error,
		func(request editorcommand.MergeFromFileRequest) (any, error) {
			if err := e.MergeFromFile(request.Path); err != nil {
				return nil, err
			}
			return respMerged, nil
		},
	)
}

func handleUpdateSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.UpdateSlideRequest, bool) {
			return editorcommand.ParseUpdateSlideRequest(
				p,
				func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
				v.OptionalString,
				v.OptionalStringSlice,
			)
		},
		v.Error,
		func(request editorcommand.UpdateSlideRequest) (any, error) {
			slide := slidesmeta.BuildSlideContent(request, e.slides[request.SlideIndex].Title)
			if err := e.UpdateSlide(request.SlideIndex, slide); err != nil {
				return nil, err
			}
			return respUpdated, nil
		},
	)
}

func handleAddChart(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.AddChartRequest, bool) {
			return editorcommand.ParseAddChartRequest(
				p,
				func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
				v.RequireString,
				v.OptionalString,
				v.RequireStringSlice,
				v.RequireFloat64Slice,
				v.OptionalInt64,
			)
		},
		v.Error,
		func(request editorcommand.AddChartRequest) (any, error) {
			chart, err := slidesmeta.BuildChartDefinition(request)
			if err != nil {
				if errors.Is(err, slidesmeta.ErrUnsupportedChartType) {
					message := fmt.Sprintf("unsupported chart type: %q", request.ChartType)
					return nil, NewBridgeError(ErrCodeInvalidValue, message)
				}
				return nil, err
			}

			if err := e.AddChart(request.SlideIndex, chart); err != nil {
				return nil, err
			}
			return respAdded, nil
		},
	)
}
