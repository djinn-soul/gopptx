package editor

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func handleListSlides(e *PresentationEditor, _ json.RawMessage) (any, error) {
	return map[string]any{"slides": e.Slides()}, nil
}

func handleFindAndReplace(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	find, ok := v.RequireString(p, "find")
	if !ok {
		return nil, v.Error()
	}
	replace, ok := v.RequireString(p, "replace")
	if !ok {
		return nil, v.Error()
	}

	count, err := e.FindAndReplaceInShapes(find, replace)
	if err != nil {
		return nil, err
	}
	return map[string]int{"replacements": count}, nil
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

func handleGetAuthors(e *PresentationEditor, _ json.RawMessage) (any, error) {
	authors, err := e.GetAuthors()
	if err != nil {
		return nil, err
	}
	return map[string]any{"authors": authors}, nil
}

func handleAddAuthor(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	name, ok := v.RequireString(p, "name")
	if !ok {
		return nil, v.Error()
	}
	initials, ok := v.RequireString(p, "initials")
	if !ok {
		return nil, v.Error()
	}

	author, err := e.AddAuthor(name, initials)
	if err != nil {
		return nil, err
	}
	return map[string]int64{"author_id": author.ID}, nil
}

func handleGetComments(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	comments, err := e.GetComments(slideIndex)
	if err != nil {
		return nil, err
	}
	return map[string]any{"comments": comments}, nil
}

func handleAddComment(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	authorID, ok := v.RequireInt64(p, "author_id")
	if !ok {
		return nil, v.Error()
	}
	text, ok := v.RequireString(p, "text")
	if !ok {
		return nil, v.Error()
	}
	x, ok := v.RequireInt64(p, "x")
	if !ok {
		return nil, v.Error()
	}
	y, ok := v.RequireInt64(p, "y")
	if !ok {
		return nil, v.Error()
	}

	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	if err := e.AddComment(slideIndex, authorID, text, x, y); err != nil {
		return nil, err
	}
	return map[string]bool{"added": true}, nil
}

func handleRemoveComment(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	authorID, ok := v.RequireInt64(p, "author_id")
	if !ok {
		return nil, v.Error()
	}
	authorIndex, ok := v.RequireInt(p, "author_index")
	if !ok {
		return nil, v.Error()
	}

	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	if err := e.RemoveComment(slideIndex, authorID, authorIndex); err != nil {
		return nil, err
	}
	return map[string]bool{"removed": true}, nil
}

func handleSetModifyPassword(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	password, ok := v.RequireString(p, "password")
	if !ok {
		return nil, v.Error()
	}

	e.Metadata().Protection.ModifyPassword = password
	return map[string]bool{"updated": true}, nil
}

func handleSetMarkAsFinal(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	final, ok := v.OptionalBool(p, "final")
	if !ok && v.HasErrors() {
		return nil, v.Error()
	}

	e.Metadata().Protection.MarkAsFinal = final
	return map[string]bool{"updated": true}, nil
}
