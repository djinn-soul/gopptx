package editor

import (
	"encoding/json"
)

// handleGetLayoutShapes returns the shape names in a slide layout.
//
// Payload: {"layout_part": "<string>"}.
// Response: {"shapes": [...]}.
func handleGetLayoutShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	layoutPart, ok := v.RequireString(p, "layout_part")
	if !ok {
		return nil, v.Error()
	}

	shapes := e.GetLayoutShapes(layoutPart)
	return map[string]any{keyShapes: shapes}, nil
}

// handleGetMasterShapes returns the shape names in a slide master.
//
// Payload: {"master_part": "<string>"}.
// Response: {"shapes": [...]}.
func handleGetMasterShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	masterPart, ok := v.RequireString(p, "master_part")
	if !ok {
		return nil, v.Error()
	}

	shapes := e.GetMasterShapes(masterPart)
	return map[string]any{keyShapes: shapes}, nil
}

// handleGetLayoutPlaceholders returns the placeholders in a slide layout.
//
// Payload: {"layout_part": "<string>"}.
// Response: {"placeholders": [...]}.
func handleGetLayoutPlaceholders(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	layoutPart, ok := v.RequireString(p, "layout_part")
	if !ok {
		return nil, v.Error()
	}

	placeholders := e.GetLayoutPlaceholders(layoutPart)
	return map[string]any{keyPlaceholder: placeholders}, nil
}

// handleGetMasterPlaceholders returns the placeholders in a slide master.
//
// Payload: {"master_part": "<string>"}.
// Response: {"placeholders": [...]}.
func handleGetMasterPlaceholders(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	masterPart, ok := v.RequireString(p, "master_part")
	if !ok {
		return nil, v.Error()
	}

	placeholders := e.GetMasterPlaceholders(masterPart)
	return map[string]any{keyPlaceholder: placeholders}, nil
}

// handleSetGlobalThemePreset applies a named preset theme to the presentation.
//
// Payload: {"name": "<string>"}.
// Response: {"applied": true}.
func handleSetGlobalThemePreset(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	name, ok := v.RequireString(p, "name")
	if !ok {
		return nil, v.Error()
	}

	if err := e.SetGlobalThemePreset(name); err != nil {
		return nil, err
	}
	return respApplied, nil
}

// handleSetThemeFontScheme updates major/minor font typefaces across all themes.
//
// Payload: {"major": "<string>", "minor": "<string>"}.
// Response: {"updated": true}.
func handleSetThemeFontScheme(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	major, ok := v.RequireString(p, "major")
	if !ok {
		return nil, v.Error()
	}
	minor, ok := v.RequireString(p, "minor")
	if !ok {
		return nil, v.Error()
	}

	if err := e.SetThemeFontScheme(major, minor); err != nil {
		return nil, err
	}
	return respUpdated, nil
}

// handleSetThemeColorScheme updates the 12 standard theme color slots.
//
// Payload: {"dk1": ..., "lt1": ..., ... (all optional)}.
// Response: {"updated": true}.
func handleSetThemeColorScheme(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	scheme := ThemeColorScheme{
		Dk1:      v.OptionalString(p, themeSlotDk1),
		Lt1:      v.OptionalString(p, themeSlotLt1),
		Dk2:      v.OptionalString(p, themeSlotDk2),
		Lt2:      v.OptionalString(p, themeSlotLt2),
		Accent1:  v.OptionalString(p, themeSlotAccent1),
		Accent2:  v.OptionalString(p, themeSlotAccent2),
		Accent3:  v.OptionalString(p, themeSlotAccent3),
		Accent4:  v.OptionalString(p, themeSlotAccent4),
		Accent5:  v.OptionalString(p, themeSlotAccent5),
		Accent6:  v.OptionalString(p, themeSlotAccent6),
		Hlink:    v.OptionalString(p, themeSlotHlink),
		FolHlink: v.OptionalString(p, "fol_hlink"),
	}

	if v.HasErrors() {
		return nil, v.Error()
	}
	if err := e.SetThemeColorScheme(scheme); err != nil {
		return nil, err
	}
	return respUpdated, nil
}

// handleGetThemeInventory returns all theme part paths and owner bindings.
//
// Payload: {} (empty).
// Response: {"theme_parts": [...], "bindings": [...]}.
func handleGetThemeInventory(e *PresentationEditor, _ json.RawMessage) (any, error) {
	inv, err := e.GetThemeInventory()
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"theme_parts": inv.ThemeParts,
		"bindings":    inv.Bindings,
	}, nil
}

func handleGetThemeColorScheme(e *PresentationEditor, _ json.RawMessage) (any, error) {
	scheme, err := e.GetThemeColorScheme()
	if err != nil {
		return nil, err
	}
	return map[string]any{
		themeSlotDk1: scheme.Dk1, themeSlotLt1: scheme.Lt1,
		themeSlotDk2: scheme.Dk2, themeSlotLt2: scheme.Lt2,
		themeSlotAccent1: scheme.Accent1, themeSlotAccent2: scheme.Accent2,
		themeSlotAccent3: scheme.Accent3, themeSlotAccent4: scheme.Accent4,
		themeSlotAccent5: scheme.Accent5, themeSlotAccent6: scheme.Accent6,
		themeSlotHlink: scheme.Hlink, "fol_hlink": scheme.FolHlink,
	}, nil
}
