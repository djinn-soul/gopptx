package editor

import (
	"encoding/json"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

func handleSetModifyPassword(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.SetModifyPasswordRequest, bool) {
			return editorcommand.ParseSetModifyPasswordRequest(p, v.RequireString)
		},
		v.Error,
		func(request editorcommand.SetModifyPasswordRequest) (any, error) {
			e.Metadata().Protection.ModifyPassword = request.Password
			return map[string]bool{"updated": true}, nil
		},
	)
}

func handleSetMarkAsFinal(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	final, ok := editorcommand.ParseSetMarkAsFinalRequest(p, v.OptionalBool)
	if !ok && v.HasErrors() {
		return nil, v.Error()
	}

	e.Metadata().Protection.MarkAsFinal = final
	return map[string]bool{"updated": true}, nil
}
