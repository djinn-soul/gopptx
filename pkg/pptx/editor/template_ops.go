package editor

import "encoding/json"

// RenderTemplate renders all Jinja2 template expressions ({{ var }},
// {{ var | filter }}, {% if %}, {% for %}, etc.) across every slide shape
// using the provided context map. Tags that have no matching key are left
// untouched. Returns the total number of text-run replacements performed.
func (e *PresentationEditor) RenderTemplate(ctx map[string]any) (int, error) {
	payload, err := json.Marshal(map[string]any{"context": ctx})
	if err != nil {
		return 0, err
	}
	result, err := handleRenderTemplate(e, payload)
	if err != nil {
		return 0, err
	}
	if m, ok := result.(map[string]any); ok {
		if n, ok := m["replacements"].(int); ok {
			return n, nil
		}
	}
	return 0, nil
}

// FromTemplate opens a .pptx file whose shape text contains Jinja2 template
// expressions, renders them in-place with the provided context, and returns
// the ready-to-save editor. The caller is responsible for calling Close().
func FromTemplate(filePath string, ctx map[string]any) (*PresentationEditor, error) {
	e, err := OpenPresentationEditor(filePath)
	if err != nil {
		return nil, err
	}
	if _, err = e.RenderTemplate(ctx); err != nil {
		_ = e.Close()
		return nil, err
	}
	return e, nil
}
