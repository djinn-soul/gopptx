package editor

import (
	"encoding/json"

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
	v := NewPayloadValidator()
	return editorcommand.HandleSlideShapeRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (int, bool) { return requireSlideIndex(e, p, v) },
		v.RequireInt,
		v.Error,
		func(request editorcommand.SlideShapeRequest, p map[string]any) (any, error) {
			var runs []common.TextRun
			if err := editorcommand.DecodeOptionalPayloadValue(p, "runs", &runs); err != nil {
				return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
			}
			if err := e.SetShapeRuns(request.SlideIndex, request.ShapeID, runs); err != nil {
				return nil, err
			}
			return map[string]bool{"updated": true}, nil
		},
	)
}

func handleUpdateDeckRunTexts(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	var slideUpdates []common.SlideRunTextUpdates
	if decodeErr := editorcommand.DecodeOptionalPayloadValue(p, "slides", &slideUpdates); decodeErr != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, decodeErr.Error())
	}
	if err := e.UpdateDeckRunTexts(slideUpdates); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleUpdateSlideRunTexts(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (int, bool) {
			return requireSlideIndex(e, p, v)
		},
		v.Error,
		func(slideIndex int, p map[string]any) (any, error) {
			var updates []common.ShapeRunTextUpdate
			if err := editorcommand.DecodeOptionalPayloadValue(p, "updates", &updates); err != nil {
				return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
			}
			if err := e.UpdateSlideRunTexts(slideIndex, updates); err != nil {
				return nil, err
			}
			return map[string]bool{"updated": true}, nil
		},
	)
}

func handleUpdateShapeRunText(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideShapeRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (int, bool) { return requireSlideIndex(e, p, v) },
		v.RequireInt,
		v.Error,
		func(request editorcommand.SlideShapeRequest, p map[string]any) (any, error) {
			runIndex, ok := v.RequireInt(p, "run_index")
			if !ok {
				return nil, v.Error()
			}
			text, ok := v.RequireString(p, "text")
			if !ok {
				return nil, v.Error()
			}
			if err := e.UpdateRunText(request.SlideIndex, request.ShapeID, runIndex, text); err != nil {
				return nil, err
			}
			return map[string]bool{"updated": true}, nil
		},
	)
}

func handleAppendShapeRun(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideShapeRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (int, bool) { return requireSlideIndex(e, p, v) },
		v.RequireInt,
		v.Error,
		func(request editorcommand.SlideShapeRequest, p map[string]any) (any, error) {
			var run common.TextRun
			if err := editorcommand.DecodeOptionalPayloadValue(p, "run", &run); err != nil {
				return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
			}
			if err := e.AppendShapeRun(request.SlideIndex, request.ShapeID, run); err != nil {
				return nil, err
			}
			return map[string]bool{"updated": true}, nil
		},
	)
}
