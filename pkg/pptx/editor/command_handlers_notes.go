package editor

import (
	"encoding/json"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

func handleGetNotes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideIndexRequest(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.Error,
		func(slideIndex int) (any, error) {
			notes, err := e.GetNotes(slideIndex)
			if err != nil {
				return nil, err
			}
			hasNotesSlide, err := e.HasNotesSlide(slideIndex)
			if err != nil {
				return nil, err
			}
			return editorcommand.BuildNotesResult(notes, hasNotesSlide), nil
		},
	)
}

func handleHasNotesSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideIndexRequest(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.Error,
		func(slideIndex int) (any, error) {
			hasNotesSlide, err := e.HasNotesSlide(slideIndex)
			if err != nil {
				return nil, err
			}
			return map[string]bool{"has_notes_slide": hasNotesSlide}, nil
		},
	)
}

func handleSetNotes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, text, ok := editorcommand.ParseSetNotesRequest(
		p,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.RequireString,
	)
	if !ok {
		return nil, v.Error()
	}
	if err := e.SetNotes(slideIndex, text); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}
