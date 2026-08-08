package editor

import (
	"testing"
)

func TestStandardThemePresets(t *testing.T) {
	presets := StandardThemePresets()
	if len(presets) == 0 {
		t.Error("expected non-empty presets map")
	}
	if _, ok := presets["facet"]; !ok {
		t.Error("expected 'facet' preset")
	}
}

func TestSetGlobalThemePreset(t *testing.T) {
	path := writeThemeFixtureDeck(t)
	ed, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	// Test successful application
	err = ed.SetGlobalThemePreset(" FA cet ")
	if err != nil {
		t.Fatalf("failed to apply preset 'facet': %v", err)
	}

	// Test invalid preset
	err = ed.SetGlobalThemePreset("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent theme preset")
	}
}

// TestResolveThemePresetAcceptsBothVocabularies guards the defect where the
// exported THEME_* constants named themes that no operation would accept: the
// gopptx theme names resolved only through apply_theme, and the Office preset
// names only through set_global_theme_preset.
func TestResolveThemePresetAcceptsBothVocabularies(t *testing.T) {
	names := []string{
		// gopptx theme names, as carried by the Python THEME_* constants.
		"Corporate", "Modern", "Vibrant", "Dark", "Nature", "Tech", "Carbon", "Office",
		// Office preset names.
		"office", "office2013", "facet", "integral", "ion", "retrospect", "slice", "wisp",
		// Matching ignores case and separators.
		" cor-porate ", "OFFICE_2013",
	}
	for _, name := range names {
		if _, ok := ResolveThemePreset(name); !ok {
			t.Errorf("ResolveThemePreset(%q) = false, want a theme", name)
		}
	}

	if _, ok := ResolveThemePreset("NotATheme"); ok {
		t.Error("ResolveThemePreset(\"NotATheme\") = true, want false")
	}
}
