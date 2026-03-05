package editor

import (
	"encoding/json"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

func handleAddSection(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.SectionAddRequest, bool) {
			return editorcommand.ParseSectionAddRequest(p, v.RequireString, v.RequireIntSlice)
		},
		v.Error,
		func(request editorcommand.SectionAddRequest) (any, error) {
			if err := e.AddSection(request.Name, request.SlideIndices); err != nil {
				return nil, err
			}
			return map[string]bool{"added": true}, nil
		},
	)
}

func handleRemoveSection(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.SectionNameRequest, bool) {
			return editorcommand.ParseSectionNameRequest(p, v.RequireString)
		},
		v.Error,
		func(request editorcommand.SectionNameRequest) (any, error) {
			if err := e.RemoveSection(request.Name); err != nil {
				return nil, err
			}
			return map[string]bool{"removed": true}, nil
		},
	)
}

func handleRenameSection(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.SectionRenameRequest, bool) {
			return editorcommand.ParseSectionRenameRequest(p, v.RequireString)
		},
		v.Error,
		func(request editorcommand.SectionRenameRequest) (any, error) {
			if err := e.RenameSection(request.OldName, request.NewName); err != nil {
				return nil, err
			}
			return map[string]bool{"renamed": true}, nil
		},
	)
}

func handleGetSections(e *PresentationEditor, _ json.RawMessage) (any, error) {
	return map[string]any{"sections": e.Sections()}, nil
}
