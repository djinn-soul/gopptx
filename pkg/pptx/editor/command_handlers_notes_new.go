package editor

import (
	"encoding/json"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// handleListNotesShapes lists all shapes in a slide's notes pane.
//
// Payload: {"slide_index": N}.
// Response: {"shapes": [...]}.
func handleListNotesShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideIndexRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (int, bool) { return requireSlideIndex(e, p, v) },
		v.Error,
		func(slideIndex int) (any, error) {
			shapes, err := e.ListNotesShapes(slideIndex)
			if err != nil {
				return nil, err
			}
			return map[string]any{keyShapes: shapes}, nil
		},
	)
}

// handleListNotesPlaceholders lists all placeholders in a slide's notes pane.
//
// Payload: {"slide_index": N}.
// Response: {"placeholders": [...]}.
func handleListNotesPlaceholders(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideIndexRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (int, bool) { return requireSlideIndex(e, p, v) },
		v.Error,
		func(slideIndex int) (any, error) {
			placeholders, err := e.ListNotesPlaceholders(slideIndex)
			if err != nil {
				return nil, err
			}
			out := make([]map[string]any, 0, len(placeholders))
			for _, ph := range placeholders {
				out = append(out, map[string]any{
					keyType:  ph.Type,
					keyIndex: ph.Index,
					keyName:  ph.Name,
				})
			}
			return map[string]any{keyPlaceholder: out}, nil
		},
	)
}

// handleUpdateNotesMaster configures the global notes master.
//
// Payload: {"header": "<string>", "footer": "<string>", "show_date_time": bool, "show_slide_num": bool}.
// Response: {"updated": true}.
func handleUpdateNotesMaster(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	master := elements.NewNotesMaster()
	master.HeaderText = v.OptionalString(p, "header")
	master.FooterText = v.OptionalString(p, "footer")

	if showDT, ok := v.OptionalBool(p, "show_date_time"); ok {
		master.ShowDateTime = showDT
	}
	if showSN, ok := v.OptionalBool(p, "show_slide_num"); ok {
		master.ShowSlideNum = showSN
	}

	if v.HasErrors() {
		return nil, v.Error()
	}

	if err := e.UpdateNotesMaster(master); err != nil {
		return nil, err
	}
	return respUpdated, nil
}
