package editor

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/noirbizarre/gonja"
)

// handleRenderTemplate renders all Jinja2 template expressions across every
// slide shape using the provided context map. It supports full Jinja2 syntax
// (variables, filters, blocks, loops) via gonja.
func handleRenderTemplate(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var req struct {
		Context map[string]any `json:"context"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, errors.New("invalid render_template payload: expected {\"context\": {...}}")
	}
	if len(req.Context) == 0 {
		return map[string]any{"replacements": 0}, nil
	}

	ctx := gonja.Context(req.Context)
	replacements := collectTemplateReplacements(e, ctx)
	total := applyTemplateReplacements(e, replacements)
	return map[string]any{"replacements": total}, nil
}

func collectTemplateReplacements(e *PresentationEditor, ctx gonja.Context) map[string]string {
	replacements := make(map[string]string) // raw shape text -> rendered text
	slideCount := e.SlideCount()
	for slideIdx := range slideCount {
		shapes, err := e.GetShapes(slideIdx)
		if err != nil {
			continue
		}
		for _, shape := range shapes {
			if _, exists := replacements[shape.Text]; exists {
				continue
			}
			if !containsTemplateSyntax(shape.Text) {
				continue
			}
			replacements[shape.Text] = renderTemplateText(shape.Text, ctx)
		}
	}
	return replacements
}

func containsTemplateSyntax(text string) bool {
	return strings.Contains(text, "{{") || strings.Contains(text, "{%")
}

func renderTemplateText(text string, ctx gonja.Context) string {
	// Render the entire shape text so multi-line blocks ({% for %}, {% if %})
	// are processed correctly.
	tpl, err := gonja.FromString(text)
	if err != nil {
		return text
	}
	rendered, err := tpl.Execute(ctx)
	if err != nil {
		return text
	}
	return rendered
}

func applyTemplateReplacements(e *PresentationEditor, replacements map[string]string) int {
	total := 0
	for raw, rendered := range replacements {
		if raw == rendered {
			continue
		}
		total += applyLineAwareReplacement(e, raw, rendered)
	}
	return total
}

func applyLineAwareReplacement(e *PresentationEditor, raw, rendered string) int {
	// When line count is unchanged, replace line-by-line to preserve per-run XML
	// formatting (bold, color, etc.) as much as possible.
	rawLines := strings.Split(raw, "\n")
	renderedLines := strings.Split(rendered, "\n")
	if len(rawLines) == len(renderedLines) {
		return replaceMatchingLines(e, rawLines, renderedLines)
	}
	// Line count changed (e.g. loop expansion); attempt whole-text replacement.
	n, _ := e.FindAndReplaceInShapes(raw, rendered)
	return n
}

func replaceMatchingLines(e *PresentationEditor, rawLines, renderedLines []string) int {
	total := 0
	for i, rawLine := range rawLines {
		if rawLine == renderedLines[i] {
			continue
		}
		n, _ := e.FindAndReplaceInShapes(rawLine, renderedLines[i])
		total += n
	}
	return total
}
