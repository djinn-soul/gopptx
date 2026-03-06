package editor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
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
			title := request.Title
			if title == "" {
				title = e.slides[request.SlideIndex].Title
			}

			slide := elements.NewSlide(title)
			if request.Layout != "" {
				slide = slide.WithLayout(request.Layout)
			}
			for _, bullet := range request.Bullets {
				slide = slide.AddBullet(bullet)
			}
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
			var chart charts.ChartDefinition
			switch strings.ToLower(request.ChartType) {
			case "bar":
				c := charts.NewBarChart(request.Categories, request.Values).WithTitle(request.Title)
				if request.W > 0 {
					c = c.Size(styling.Emu(request.W), styling.Emu(request.H)).
						Position(styling.Emu(request.X), styling.Emu(request.Y))
				}
				chart = c
			case "line":
				c := charts.NewLineChart(request.Categories, request.Values).WithTitle(request.Title)
				if request.W > 0 {
					c = c.Size(styling.Emu(request.W), styling.Emu(request.H)).
						Position(styling.Emu(request.X), styling.Emu(request.Y))
				}
				chart = c
			case "pie":
				c := charts.NewPieChart(request.Categories, request.Values).WithTitle(request.Title)
				if request.W > 0 {
					c = c.Size(styling.Emu(request.W), styling.Emu(request.H)).
						Position(styling.Emu(request.X), styling.Emu(request.Y))
				}
				chart = c
			default:
				return nil, NewBridgeError(
					ErrCodeInvalidValue,
					fmt.Sprintf("unsupported chart type: %q", request.ChartType),
				)
			}

			if err := e.AddChart(request.SlideIndex, chart); err != nil {
				return nil, err
			}
			return map[string]bool{"added": true}, nil
		},
	)
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

	theme, err := resolveThemeByName(themeName)
	if err != nil {
		return nil, err
	}
	if err := e.ApplyTheme(theme); err != nil {
		return nil, err
	}
	return map[string]bool{"applied": true}, nil
}

func resolveThemeByName(name string) (styling.Theme, error) {
	switch name {
	case "Corporate":
		return styling.ThemeCorporate, nil
	case "Modern":
		return styling.ThemeModern, nil
	case "Vibrant":
		return styling.ThemeVibrant, nil
	case "Dark":
		return styling.ThemeDark, nil
	case "Nature":
		return styling.ThemeNature, nil
	case "Tech":
		return styling.ThemeTech, nil
	case "Carbon":
		return styling.ThemeCarbon, nil
	default:
		return styling.Theme{}, NewBridgeError(ErrCodeInvalidValue, fmt.Sprintf("unknown theme name %q", name))
	}
}
