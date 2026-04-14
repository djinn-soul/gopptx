package editor

import (
	"encoding/json"
	"errors"
	"fmt"

	slidesmeta "github.com/djinn-soul/gopptx/pkg/pptx/editor/handlers/slidesmeta"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

func handleGetCoreProperties(e *PresentationEditor, _ json.RawMessage) (any, error) {
	return e.GetCoreProperties(), nil
}

func handleSetCoreProperties(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	props := editorcommand.ParseCorePropertiesRequest(p, v.OptionalString)
	if v.HasErrors() {
		return nil, v.Error()
	}

	e.SetCoreProperties(props)
	return respUpdated, nil
}

func handleApplyTheme(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	themeName, ok := v.RequireString(p, "theme_name")
	if !ok {
		return nil, v.Error()
	}

	theme, err := slidesmeta.ResolveThemeByName(themeName)
	if err != nil {
		if errors.Is(err, slidesmeta.ErrUnknownThemeName) {
			return nil, NewBridgeError(ErrCodeInvalidValue, fmt.Sprintf("unknown theme name %q", themeName))
		}
		return nil, err
	}
	if err := e.ApplyTheme(theme); err != nil {
		return nil, err
	}
	return respApplied, nil
}

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
			return respAdded, nil
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
			return respRemoved, nil
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
