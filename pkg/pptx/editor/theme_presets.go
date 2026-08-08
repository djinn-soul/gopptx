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

// ResolveThemePreset resolves a theme name to a concrete theme.
//
// Both the Office preset names ("facet", "ion") and the gopptx theme names
// ("Corporate", "Dark") are accepted; styling.ResolveTheme owns the mapping.
func ResolveThemePreset(name string) (styling.Theme, bool) {
	return styling.ResolveTheme(name)
}

// SetGlobalThemePreset applies a preset to the package theme part.
func (e *PresentationEditor) SetGlobalThemePreset(name string) error {
	theme, ok := ResolveThemePreset(name)
	if !ok {
		return fmt.Errorf(
			"unknown theme preset %q; accepted names are %s",
			name, strings.Join(styling.ThemeNames(), ", "),
		)
	}
	return e.ApplyTheme(theme)
}
