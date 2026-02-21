package editor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func handleAddSection(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	name, ok := v.RequireString(p, "name")
	if !ok {
		return nil, v.Error()
	}
	slideIndices, ok := v.RequireIntSlice(p, "slide_indices")
	if !ok {
		return nil, v.Error()
	}

	if err := e.AddSection(name, slideIndices); err != nil {
		return nil, err
	}
	return map[string]bool{"added": true}, nil
}

func handleRemoveSection(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	name, ok := v.RequireString(p, "name")
	if !ok {
		return nil, v.Error()
	}

	if err := e.RemoveSection(name); err != nil {
		return nil, err
	}
	return map[string]bool{"removed": true}, nil
}

func handleRenameSection(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	oldName, ok := v.RequireString(p, "old_name")
	if !ok {
		return nil, v.Error()
	}
	newName, ok := v.RequireString(p, "new_name")
	if !ok {
		return nil, v.Error()
	}

	if err := e.RenameSection(oldName, newName); err != nil {
		return nil, err
	}
	return map[string]bool{"renamed": true}, nil
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
	props := common.CoreProperties{
		Title:          v.OptionalString(p, "title"),
		Subject:        v.OptionalString(p, "subject"),
		Creator:        v.OptionalString(p, "creator"),
		Keywords:       v.OptionalString(p, "keywords"),
		Description:    v.OptionalString(p, "description"),
		LastModifiedBy: v.OptionalString(p, "lastModifiedBy"),
		Revision:       v.OptionalString(p, "revision"),
		Created:        v.OptionalString(p, "created"),
		Modified:       v.OptionalString(p, "modified"),
		Category:       v.OptionalString(p, "category"),
		ContentStatus:  v.OptionalString(p, "contentStatus"),
	}

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

	var theme styling.Theme
	switch themeName {
	case "Corporate":
		theme = styling.ThemeCorporate
	case "Modern":
		theme = styling.ThemeModern
	case "Vibrant":
		theme = styling.ThemeVibrant
	case "Dark":
		theme = styling.ThemeDark
	case "Nature":
		theme = styling.ThemeNature
	case "Tech":
		theme = styling.ThemeTech
	case "Carbon":
		theme = styling.ThemeCarbon
	default:
		return nil, NewBridgeError(ErrCodeInvalidValue, fmt.Sprintf("unknown theme name %q", themeName))
	}

	if err := e.ApplyTheme(theme); err != nil {
		return nil, err
	}
	return map[string]bool{"applied": true}, nil
}

func handleSetSlideSize(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	width, ok := v.RequireInt64(p, "width")
	if !ok {
		return nil, v.Error()
	}
	height, ok := v.RequireInt64(p, "height")
	if !ok {
		return nil, v.Error()
	}

	if err := e.SetSlideSize(common.SlideSize{Width: width, Height: height}); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleSetSlideTitle(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	title, ok := v.RequireString(p, "title")
	if !ok {
		return nil, v.Error()
	}

	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	if err := e.SetSlideTitle(slideIndex, title); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleMergeFromFile(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	path, ok := v.RequireString(p, "path")
	if !ok {
		return nil, v.Error()
	}

	if err := e.MergeFromFile(path); err != nil {
		return nil, err
	}
	return map[string]bool{"merged": true}, nil
}

func handleUpdateSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

	title := v.OptionalString(p, "title")
	layout := v.OptionalString(p, "layout")
	bullets, _ := v.OptionalStringSlice(p, "bullets")

	// Get current title if not provided
	currentTitle := e.slides[slideIndex].Title
	if title == "" {
		title = currentTitle
	}

	slide := elements.NewSlide(title)
	if layout != "" {
		slide = slide.WithLayout(layout)
	}
	for _, b := range bullets {
		slide = slide.AddBullet(b)
	}
	if err := e.UpdateSlide(slideIndex, slide); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleAddChart(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	chartType, ok := v.RequireString(p, "chart_type")
	if !ok {
		return nil, v.Error()
	}
	title := v.OptionalString(p, "title")
	categories, ok := v.RequireStringSlice(p, "categories")
	if !ok {
		return nil, v.Error()
	}
	values, ok := v.RequireFloat64Slice(p, "values")
	if !ok {
		return nil, v.Error()
	}

	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	x, _ := v.OptionalInt64(p, "x")
	y, _ := v.OptionalInt64(p, "y")
	w, _ := v.OptionalInt64(p, "w")
	h, _ := v.OptionalInt64(p, "h")

	var chart charts.ChartDefinition
	switch strings.ToLower(chartType) {
	case "bar":
		c := charts.NewBarChart(categories, values).WithTitle(title)
		if w > 0 {
			c = c.Size(styling.Emu(w), styling.Emu(h)).Position(styling.Emu(x), styling.Emu(y))
		}
		chart = c
	case "line":
		c := charts.NewLineChart(categories, values).WithTitle(title)
		if w > 0 {
			c = c.Size(styling.Emu(w), styling.Emu(h)).Position(styling.Emu(x), styling.Emu(y))
		}
		chart = c
	case "pie":
		c := charts.NewPieChart(categories, values).WithTitle(title)
		if w > 0 {
			c = c.Size(styling.Emu(w), styling.Emu(h)).Position(styling.Emu(x), styling.Emu(y))
		}
		chart = c
	default:
		return nil, NewBridgeError(ErrCodeInvalidValue, fmt.Sprintf("unsupported chart type: %q", chartType))
	}

	if err := e.AddChart(slideIndex, chart); err != nil {
		return nil, err
	}
	return map[string]bool{"added": true}, nil
}
