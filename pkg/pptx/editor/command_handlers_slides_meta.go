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
	var p struct {
		Name         string `json:"name"`
		SlideIndices []int  `json:"slide_indices"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.AddSection(p.Name, p.SlideIndices); err != nil {
		return nil, err
	}
	return map[string]bool{"added": true}, nil
}

func handleRemoveSection(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.RemoveSection(p.Name); err != nil {
		return nil, err
	}
	return map[string]bool{"removed": true}, nil
}

func handleRenameSection(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		OldName string `json:"old_name"`
		NewName string `json:"new_name"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.RenameSection(p.OldName, p.NewName); err != nil {
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
	var p common.CoreProperties
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	e.SetCoreProperties(p)
	return map[string]bool{"updated": true}, nil
}

func handleApplyTheme(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		ThemeName string `json:"theme_name"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}

	var theme styling.Theme
	switch p.ThemeName {
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
		return nil, fmt.Errorf("unknown theme name %q", p.ThemeName)
	}

	if err := e.ApplyTheme(theme); err != nil {
		return nil, err
	}
	return map[string]bool{"applied": true}, nil
}

func handleSetSlideSize(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p common.SlideSize
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.SetSlideSize(p); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleSetSlideTitle(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex int    `json:"slide_index"`
		Title      string `json:"title"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.SetSlideTitle(p.SlideIndex, p.Title); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleMergeFromFile(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.MergeFromFile(p.Path); err != nil {
		return nil, err
	}
	return map[string]bool{"merged": true}, nil
}

func handleUpdateSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex int       `json:"slide_index"`
		Title      *string   `json:"title"`
		Layout     *string   `json:"layout"`
		Bullets    *[]string `json:"bullets"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if p.SlideIndex < 0 || p.SlideIndex >= len(e.slides) {
		return nil, fmt.Errorf("slide index %d out of range [0,%d)", p.SlideIndex, len(e.slides))
	}

	title := e.slides[p.SlideIndex].Title
	if p.Title != nil {
		title = *p.Title
	}

	slide := elements.NewSlide(title)
	if p.Layout != nil && *p.Layout != "" {
		slide = slide.WithLayout(*p.Layout)
	}
	if p.Bullets != nil {
		for _, b := range *p.Bullets {
			slide = slide.AddBullet(b)
		}
	}
	if err := e.UpdateSlide(p.SlideIndex, slide); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleAddChart(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex int       `json:"slide_index"`
		ChartType  string    `json:"chart_type"`
		Title      string    `json:"title"`
		Categories []string  `json:"categories"`
		Values     []float64 `json:"values"`
		X          int64     `json:"x"`
		Y          int64     `json:"y"`
		W          int64     `json:"w"`
		H          int64     `json:"h"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}

	var chart charts.ChartDefinition
	switch strings.ToLower(p.ChartType) {
	case "bar":
		c := charts.NewBarChart(p.Categories, p.Values).WithTitle(p.Title)
		if p.W > 0 {
			c = c.Size(styling.Emu(p.W), styling.Emu(p.H)).Position(styling.Emu(p.X), styling.Emu(p.Y))
		}
		chart = c
	case "line":
		c := charts.NewLineChart(p.Categories, p.Values).WithTitle(p.Title)
		if p.W > 0 {
			c = c.Size(styling.Emu(p.W), styling.Emu(p.H)).Position(styling.Emu(p.X), styling.Emu(p.Y))
		}
		chart = c
	case "pie":
		c := charts.NewPieChart(p.Categories, p.Values).WithTitle(p.Title)
		if p.W > 0 {
			c = c.Size(styling.Emu(p.W), styling.Emu(p.H)).Position(styling.Emu(p.X), styling.Emu(p.Y))
		}
		chart = c
	default:
		return nil, fmt.Errorf("unsupported chart type: %q", p.ChartType)
	}

	if err := e.AddChart(p.SlideIndex, chart); err != nil {
		return nil, err
	}
	return map[string]bool{"added": true}, nil
}
