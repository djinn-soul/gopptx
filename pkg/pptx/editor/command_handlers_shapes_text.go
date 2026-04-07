package editor

import (
	"encoding/json"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

func handleGetShapeTextState(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideShapeRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (int, bool) { return requireSlideIndex(e, p, v) },
		v.RequireInt,
		v.Error,
		func(request editorcommand.SlideShapeRequest) (any, error) {
			state, err := e.GetShapeTextState(request.SlideIndex, request.ShapeID)
			if err != nil {
				return nil, err
			}
			return map[string]any{
				"text":       state.Text,
				"runs":       state.Runs,
				"text_frame": state.TextFrame,
				"paragraph":  state.Paragraph,
			}, nil
		},
	)
}

func handleGetSlideTextStates(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideIndexRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (int, bool) { return requireSlideIndex(e, p, v) },
		v.Error,
		func(slideIndex int) (any, error) {
			states, err := e.GetSlideTextStates(slideIndex)
			if err != nil {
				return nil, err
			}
			return map[string]any{"states": states}, nil
		},
	)
}

func handleGetShapeRuns(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideShapeRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (int, bool) { return requireSlideIndex(e, p, v) },
		v.RequireInt,
		v.Error,
		func(request editorcommand.SlideShapeRequest) (any, error) {
			runs, err := e.GetShapeRuns(request.SlideIndex, request.ShapeID)
			if err != nil {
				return nil, err
			}
			return map[string]any{"runs": runs}, nil
		},
	)
}

func handleSetShapeRuns(e *PresentationEditor, payload json.RawMessage) (any, error) {
	slideIndex, shapeID, runs, err := parseSetShapeRunsPayload(payload, e.SlideCount())
	if err != nil {
		return nil, err
	}
	if err := e.SetShapeRuns(slideIndex, shapeID, runs); err != nil {
		return nil, err
	}
	return respUpdated, nil
}

func handleSetSlideShapeRuns(e *PresentationEditor, payload json.RawMessage) (any, error) {
	slideIndex, updates, err := parseSetSlideShapeRunsPayload(payload, e.SlideCount())
	if err != nil {
		return nil, err
	}
	if err := e.SetSlideShapeRuns(slideIndex, updates); err != nil {
		return nil, err
	}
	return respUpdated, nil
}

func handleUpdateDeckRunTexts(e *PresentationEditor, payload json.RawMessage) (any, error) {
	slideUpdates, err := parseUpdateDeckRunTextsPayload(payload)
	if err != nil {
		return nil, err
	}
	if err := e.UpdateDeckRunTexts(slideUpdates); err != nil {
		return nil, err
	}
	return respUpdated, nil
}

func handleUpdateSlideRunTexts(e *PresentationEditor, payload json.RawMessage) (any, error) {
	slideIndex, updates, err := parseUpdateSlideRunTextsPayload(payload, e.SlideCount())
	if err != nil {
		return nil, err
	}
	if err := e.UpdateSlideRunTexts(slideIndex, updates); err != nil {
		return nil, err
	}
	return respUpdated, nil
}

func handleUpdateShapeRunText(e *PresentationEditor, payload json.RawMessage) (any, error) {
	slideIndex, shapeID, runIndex, text, err := parseUpdateShapeRunTextPayload(payload, e.SlideCount())
	if err != nil {
		return nil, err
	}
	if err := e.UpdateRunText(slideIndex, shapeID, runIndex, text); err != nil {
		return nil, err
	}
	return respUpdated, nil
}

func handleAppendShapeRun(e *PresentationEditor, payload json.RawMessage) (any, error) {
	slideIndex, shapeID, run, err := parseAppendShapeRunPayload(payload, e.SlideCount())
	if err != nil {
		return nil, err
	}
	if err := e.AppendShapeRun(slideIndex, shapeID, run); err != nil {
		return nil, err
	}
	return respUpdated, nil
}

func parseSlideShapeBase(payload json.RawMessage, slideCount int) (int, int, error) {
	if len(payload) == 0 {
		return 0, 0, NewBridgeError(ErrCodeInvalidPayload, "empty payload")
	}
	var fast struct {
		SlideIndex *float64 `json:"slide_index"`
		ShapeID    *float64 `json:"shape_id"`
	}
	if err := json.Unmarshal(payload, &fast); err == nil {
		if fast.SlideIndex == nil {
			return 0, 0, newPayloadValidationBridgeError(ErrCodeMissingField, "missing required field: slide_index")
		}
		if fast.ShapeID == nil {
			return 0, 0, newPayloadValidationBridgeError(ErrCodeMissingField, "missing required field: shape_id")
		}
		slideIndex := int(*fast.SlideIndex)
		if slideIndex < 0 || slideIndex >= slideCount {
			msg := fmt.Sprintf("slide_index %d out of bounds [%d, %d)", slideIndex, 0, slideCount)
			return 0, 0, newPayloadValidationBridgeError(ErrCodeInvalidIndex, msg)
		}
		return slideIndex, int(*fast.ShapeID), nil
	}

	var raw struct {
		SlideIndex json.RawMessage `json:"slide_index"`
		ShapeID    json.RawMessage `json:"shape_id"`
	}
	if err := json.Unmarshal(payload, &raw); err != nil {
		return 0, 0, NewBridgeError(ErrCodeInvalidPayload, "invalid JSON payload: "+err.Error())
	}
	if len(raw.SlideIndex) == 0 {
		return 0, 0, newPayloadValidationBridgeError(ErrCodeMissingField, "missing required field: slide_index")
	}
	if len(raw.ShapeID) == 0 {
		return 0, 0, newPayloadValidationBridgeError(ErrCodeMissingField, "missing required field: shape_id")
	}
	var slideNumber float64
	if err := json.Unmarshal(raw.SlideIndex, &slideNumber); err != nil {
		return 0, 0, newPayloadValidationBridgeError(ErrCodeInvalidType, "field slide_index must be an integer")
	}
	slideIndex := int(slideNumber)
	if slideIndex < 0 || slideIndex >= slideCount {
		msg := fmt.Sprintf("slide_index %d out of bounds [%d, %d)", slideIndex, 0, slideCount)
		return 0, 0, newPayloadValidationBridgeError(ErrCodeInvalidIndex, msg)
	}
	var shapeNumber float64
	if err := json.Unmarshal(raw.ShapeID, &shapeNumber); err != nil {
		return 0, 0, newPayloadValidationBridgeError(ErrCodeInvalidType, "field shape_id must be an integer")
	}
	return slideIndex, int(shapeNumber), nil
}

func parseSetShapeRunsPayload(payload json.RawMessage, slideCount int) (int, int, []common.TextRun, error) {
	slideIndex, shapeID, err := parseSlideShapeBase(payload, slideCount)
	if err != nil {
		return 0, 0, nil, err
	}
	var p struct {
		Runs []common.TextRun `json:"runs"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return 0, 0, nil, NewBridgeError(ErrCodeInvalidPayload, "invalid JSON payload: "+err.Error())
	}
	return slideIndex, shapeID, p.Runs, nil
}

func parseSetSlideShapeRunsPayload(payload json.RawMessage, slideCount int) (int, []common.ShapeRunsUpdate, error) {
	if len(payload) == 0 {
		return 0, nil, NewBridgeError(ErrCodeInvalidPayload, "empty payload")
	}
	var fast struct {
		SlideIndex *float64                 `json:"slide_index"`
		Updates    []common.ShapeRunsUpdate `json:"updates"`
	}
	if err := json.Unmarshal(payload, &fast); err != nil {
		return 0, nil, NewBridgeError(ErrCodeInvalidPayload, "invalid JSON payload: "+err.Error())
	}
	if fast.SlideIndex == nil {
		return 0, nil, newPayloadValidationBridgeError(ErrCodeMissingField, "missing required field: slide_index")
	}
	slideIndex := int(*fast.SlideIndex)
	if slideIndex < 0 || slideIndex >= slideCount {
		msg := fmt.Sprintf("slide_index %d out of bounds [%d, %d)", slideIndex, 0, slideCount)
		return 0, nil, newPayloadValidationBridgeError(ErrCodeInvalidIndex, msg)
	}
	return slideIndex, fast.Updates, nil
}

func parseUpdateDeckRunTextsPayload(payload json.RawMessage) ([]common.SlideRunTextUpdates, error) {
	if len(payload) == 0 {
		return nil, NewBridgeError(ErrCodeInvalidPayload, "empty payload")
	}
	var p struct {
		Slides []common.SlideRunTextUpdates `json:"slides"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, "invalid JSON payload: "+err.Error())
	}
	return p.Slides, nil
}

func parseUpdateSlideRunTextsPayload(
	payload json.RawMessage,
	slideCount int,
) (int, []common.ShapeRunTextUpdate, error) {
	if len(payload) == 0 {
		return 0, nil, NewBridgeError(ErrCodeInvalidPayload, "empty payload")
	}
	var p struct {
		SlideIndex *float64                    `json:"slide_index"`
		Updates    []common.ShapeRunTextUpdate `json:"updates"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return 0, nil, NewBridgeError(ErrCodeInvalidPayload, "invalid JSON payload: "+err.Error())
	}
	if p.SlideIndex == nil {
		return 0, nil, newPayloadValidationBridgeError(ErrCodeMissingField, "missing required field: slide_index")
	}
	slideIndex := int(*p.SlideIndex)
	if slideIndex < 0 || slideIndex >= slideCount {
		msg := fmt.Sprintf("slide_index %d out of bounds [%d, %d)", slideIndex, 0, slideCount)
		return 0, nil, newPayloadValidationBridgeError(ErrCodeInvalidIndex, msg)
	}
	return slideIndex, p.Updates, nil
}

func parseUpdateShapeRunTextPayload(payload json.RawMessage, slideCount int) (int, int, int, string, error) {
	slideIndex, shapeID, err := parseSlideShapeBase(payload, slideCount)
	if err != nil {
		return 0, 0, 0, "", err
	}
	var p struct {
		RunIndex *float64 `json:"run_index"`
		Text     *string  `json:"text"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return 0, 0, 0, "", NewBridgeError(ErrCodeInvalidPayload, "invalid JSON payload: "+err.Error())
	}
	if p.RunIndex == nil {
		return 0, 0, 0, "", newPayloadValidationBridgeError(ErrCodeMissingField, "missing required field: run_index")
	}
	if p.Text == nil {
		return 0, 0, 0, "", newPayloadValidationBridgeError(ErrCodeMissingField, "missing required field: text")
	}
	return slideIndex, shapeID, int(*p.RunIndex), *p.Text, nil
}

func parseAppendShapeRunPayload(payload json.RawMessage, slideCount int) (int, int, common.TextRun, error) {
	slideIndex, shapeID, err := parseSlideShapeBase(payload, slideCount)
	if err != nil {
		return 0, 0, common.TextRun{}, err
	}
	var p struct {
		Run common.TextRun `json:"run"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return 0, 0, common.TextRun{}, NewBridgeError(ErrCodeInvalidPayload, "invalid JSON payload: "+err.Error())
	}
	return slideIndex, shapeID, p.Run, nil
}
