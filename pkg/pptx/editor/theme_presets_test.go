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
