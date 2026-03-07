package slidesmeta

import (
	"encoding/json"
	"fmt"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
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
