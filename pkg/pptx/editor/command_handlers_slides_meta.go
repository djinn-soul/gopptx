package editor

import (
	"encoding/json"
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	slidesmeta "github.com/djinn-soul/gopptx/pkg/pptx/editor/handlers/slidesmeta"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

func handleSetSlideSize(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.SlideSizeRequest, bool) {
			return editorcommand.ParseSlideSizeRequest(p, v.RequireInt64)
		},
		v.Error,
		func(request editorcommand.SlideSizeRequest) (any, error) {
			if err := e.SetSlideSize(common.SlideSize{Width: request.Width, Height: request.Height}); err != nil {
				return nil, err
			}
			return map[string]bool{"updated": true}, nil
		},
	)
}

func handleSetSlideTitle(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.SlideTitleRequest, bool) {
			return editorcommand.ParseSlideTitleRequest(
				p,
				func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
				v.RequireString,
			)
		},
		v.Error,
		func(request editorcommand.SlideTitleRequest) (any, error) {
			if err := e.SetSlideTitle(request.SlideIndex, request.Title); err != nil {
				return nil, err
			}
			return map[string]bool{"updated": true}, nil
		},
	)
}

func handleMergeFromFile(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.MergeFromFileRequest, bool) {
			return editorcommand.ParseMergeFromFileRequest(p, v.RequireString)
		},
		v.Error,
		func(request editorcommand.MergeFromFileRequest) (any, error) {
			if err := e.MergeFromFile(request.Path); err != nil {
				return nil, err
			}
			return map[string]bool{"merged": true}, nil
		},
	)
}

func handleUpdateSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.UpdateSlideRequest, bool) {
			return editorcommand.ParseUpdateSlideRequest(
				p,
				func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
				v.OptionalString,
				v.OptionalStringSlice,
			)
		},
		v.Error,
		func(request editorcommand.UpdateSlideRequest) (any, error) {
			slide := slidesmeta.BuildSlideContent(request, e.slides[request.SlideIndex].Title)
			if err := e.UpdateSlide(request.SlideIndex, slide); err != nil {
				return nil, err
			}
			return map[string]bool{"updated": true}, nil
		},
	)
}

func handleAddChart(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.AddChartRequest, bool) {
			return editorcommand.ParseAddChartRequest(
				p,
				func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
				v.RequireString,
				v.OptionalString,
				v.RequireStringSlice,
				v.RequireFloat64Slice,
				v.OptionalInt64,
			)
		},
		v.Error,
		func(request editorcommand.AddChartRequest) (any, error) {
			chart, err := slidesmeta.BuildChartDefinition(request)
			if err != nil {
				if errors.Is(err, slidesmeta.ErrUnsupportedChartType) {
					message := fmt.Sprintf("unsupported chart type: %q", request.ChartType)
					return nil, NewBridgeError(ErrCodeInvalidValue, message)
				}
				return nil, err
			}

			if err := e.AddChart(request.SlideIndex, chart); err != nil {
				return nil, err
			}
			return map[string]bool{"added": true}, nil
		},
	)
}

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
	return map[string]bool{"updated": true}, nil
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
	return map[string]bool{"applied": true}, nil
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
