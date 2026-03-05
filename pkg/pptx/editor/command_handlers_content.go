package editor

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
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
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
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
	slideIndex, ok := requireSlideIndex(e, p, v)
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
	slideIndex, ok := requireSlideIndex(e, p, v)
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

func handleAddCustomXML(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	content := v.OptionalString(p, "content")
	rootElement := v.OptionalString(p, "root_element")
	namespace := v.OptionalString(p, "namespace")

	props := parseCustomXMLProperties(p, v)

	if content == "" && rootElement == "" {
		v.setCode(ErrCodeMissingField)
		v.errors = append(v.errors, "either content or root_element must be provided")
	}

	if v.HasErrors() {
		return nil, v.Error()
	}

	part := common.CustomXMLPart{
		Content:     content,
		RootElement: rootElement,
		Namespace:   namespace,
	}

	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		part.Properties = append(part.Properties, common.CustomXMLKV{Key: k, Value: props[k]})
	}

	e.metadata.CustomXML = append(e.metadata.CustomXML, part)

	return map[string]int{"index": len(e.metadata.CustomXML) - 1}, nil
}

func handleListCustomXML(e *PresentationEditor, _ json.RawMessage) (any, error) {
	type CustomXMLResp struct {
		Content     string            `json:"content,omitempty"`
		RootElement string            `json:"root_element,omitempty"`
		Namespace   string            `json:"namespace,omitempty"`
		Properties  map[string]string `json:"properties,omitempty"`
	}

	out := make([]CustomXMLResp, len(e.metadata.CustomXML))
	for i, part := range e.metadata.CustomXML {
		var props map[string]string
		if len(part.Properties) > 0 {
			props = make(map[string]string, len(part.Properties))
			for _, kv := range part.Properties {
				props[kv.Key] = kv.Value
			}
		}

		out[i] = CustomXMLResp{
			Content:     part.Content,
			RootElement: part.RootElement,
			Namespace:   part.Namespace,
			Properties:  props,
		}
	}
	return map[string]any{"custom_xml": out}, nil
}

func handleRemoveCustomXML(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	index, ok := v.RequireInt(p, "index")
	if !ok {
		return nil, v.Error()
	}

	if index < 0 || index >= len(e.metadata.CustomXML) {
		v.setCode(ErrCodeInvalidIndex)
		v.errors = append(
			v.errors,
			fmt.Sprintf("custom xml index %d out of range [0,%d)", index, len(e.metadata.CustomXML)),
		)
		return nil, v.Error()
	}

	e.metadata.CustomXML = append(e.metadata.CustomXML[:index], e.metadata.CustomXML[index+1:]...)

	return map[string]bool{"removed": true}, nil
}

func parseCustomXMLProperties(payload map[string]any, v *PayloadValidator) map[string]string {
	rawProps, ok := payload["properties"]
	if !ok || rawProps == nil {
		return nil
	}
	propMap, ok := rawProps.(map[string]any)
	if !ok {
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, "properties must be an object with string values")
		return nil
	}
	props := make(map[string]string, len(propMap))
	for k, val := range propMap {
		s, ok := val.(string)
		if !ok {
			v.setCode(ErrCodeInvalidType)
			v.errors = append(v.errors, fmt.Sprintf("property %q must be a string", k))
			continue
		}
		props[k] = s
	}
	return props
}

func handleAddVba(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	dataBase64, ok := v.RequireString(p, "data")
	if !ok {
		return nil, v.Error()
	}

	data, err := base64.StdEncoding.DecodeString(dataBase64)
	if err != nil {
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, "data must be a valid base64 string")
		return nil, v.Error()
	}

	var project *vba.VBAProject
	if e.metadata.VBA == nil {
		project = vba.New()
		e.metadata.VBA = project
	} else {
		project, ok = e.metadata.VBA.(*vba.VBAProject)
		if !ok {
			return nil, errors.New("invalid VBA metadata type")
		}
	}
	project.SetData(data)

	return map[string]bool{"added": true}, nil
}
