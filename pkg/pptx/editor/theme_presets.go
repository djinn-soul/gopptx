package editor

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// StandardThemePresets returns common preset names mapped to theme payloads.
func StandardThemePresets() map[string]styling.Theme {
	return map[string]styling.Theme{
		"office2013": styling.ThemeCorporate,
		"office":     styling.ThemeCorporate,
		"facet":      styling.ThemeModern,
		"integral":   styling.ThemeTech,
		"ion":        styling.ThemeDark,
		"retrospect": styling.ThemeVibrant,
		"slice":      styling.ThemeNature,
		"wisp":       styling.ThemeCarbon,
	}
}

// ResolveThemePreset resolves a preset name to a concrete theme.
func ResolveThemePreset(name string) (styling.Theme, bool) {
	switch normalizeThemePresetName(name) {
	case "office2013", "office":
		return styling.ThemeCorporate, true
	case "facet":
		return styling.ThemeModern, true
	case "integral":
		return styling.ThemeTech, true
	case "ion":
		return styling.ThemeDark, true
	case "retrospect":
		return styling.ThemeVibrant, true
	case "slice":
		return styling.ThemeNature, true
	case "wisp":
		return styling.ThemeCarbon, true
	default:
		return styling.Theme{}, false
	}
}

// SetGlobalThemePreset applies a preset to the package theme part.
func (e *PresentationEditor) SetGlobalThemePreset(name string) error {
	theme, ok := ResolveThemePreset(name)
	if !ok {
		return fmt.Errorf("unknown theme preset %q", name)
	}
	return e.ApplyTheme(theme)
}

func normalizeThemePresetName(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	normalized = strings.ReplaceAll(normalized, " ", "")
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, "_", "")
	return normalized
}
