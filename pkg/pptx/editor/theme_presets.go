package editor

import (
	"fmt"
	"maps"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

var standardThemePresets = map[string]styling.Theme{
	"office2013": styling.ThemeCorporate,
	"office":     styling.ThemeCorporate,
	"facet":      styling.ThemeModern,
	"integral":   styling.ThemeTech,
	"ion":        styling.ThemeDark,
	"retrospect": styling.ThemeVibrant,
	"slice":      styling.ThemeNature,
	"wisp":       styling.ThemeCarbon,
}

// StandardThemePresets returns common preset names mapped to theme payloads.
func StandardThemePresets() map[string]styling.Theme {
	out := make(map[string]styling.Theme, len(standardThemePresets))
	maps.Copy(out, standardThemePresets)
	return out
}

// ResolveThemePreset resolves a preset name to a concrete theme.
func ResolveThemePreset(name string) (styling.Theme, bool) {
	key := normalizeThemePresetName(name)
	theme, ok := standardThemePresets[key]
	return theme, ok
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
