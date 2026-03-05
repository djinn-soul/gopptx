package editor

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

func handleListSlides(e *PresentationEditor, _ json.RawMessage) (any, error) {
	return map[string]any{"slides": e.Slides()}, nil
}

func handleFindAndReplace(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.FindReplaceRequest, bool) {
			return editorcommand.ParseFindReplaceRequest(p, v.RequireString)
		},
		v.Error,
		func(request editorcommand.FindReplaceRequest) (any, error) {
			count, err := e.FindAndReplaceInShapes(request.Find, request.Replace)
			if err != nil {
				return nil, err
			}
			return map[string]int{"replacements": count}, nil
		},
	)
}

func handleSearchShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	query := common.ShapeSearchQuery{
		NameContains: v.OptionalString(p, "name_contains"),
		TypeEquals:   v.OptionalString(p, "type_equals"),
		TextContains: v.OptionalString(p, "text_contains"),
	}
	query.CaseSensitive, _ = v.OptionalBool(p, "case_sensitive")

	if v.HasErrors() {
		return nil, v.Error()
	}

	results, err := e.SearchShapes(query)
	if err != nil {
		return nil, err
	}
	return map[string]any{"results": results}, nil
}
