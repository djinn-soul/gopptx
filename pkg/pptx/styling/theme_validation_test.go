package styling_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestValidateColorScheme(t *testing.T) {
	valid := styling.ThemeCorporate.Colors
	if err := styling.ValidateColorScheme(valid); err != nil {
		t.Errorf("expected valid color scheme, got error: %v", err)
	}

	tests := []struct {
		name    string
		modify  func(styling.ColorScheme) styling.ColorScheme
		wantErr string
	}{
		{
			name:    "Empty Dk1",
			modify:  func(cs styling.ColorScheme) styling.ColorScheme { cs.Dk1 = ""; return cs },
			wantErr: "Dk1 is required",
		},
		{
			name:    "Invalid Accent1 hex",
			modify:  func(cs styling.ColorScheme) styling.ColorScheme { cs.Accent1 = "GGGGGG"; return cs },
			wantErr: "Accent1 must be a 6-digit RGB hex",
		},
		{
			name:    "Short color",
			modify:  func(cs styling.ColorScheme) styling.ColorScheme { cs.Hlink = "FFF"; return cs },
			wantErr: "Hlink must be a 6-digit RGB hex",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := tt.modify(valid)
			err := styling.ValidateColorScheme(cs)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error %q should contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestValidateTheme(t *testing.T) {
	for _, theme := range styling.AllThemes() {
		t.Run(theme.Name, func(t *testing.T) {
			if err := styling.ValidateTheme(theme); err != nil {
				t.Errorf("preset theme %q should be valid: %v", theme.Name, err)
			}
		})
	}

	tests := []struct {
		name    string
		theme   styling.Theme
		wantErr string
	}{
		{
			name:    "Empty name",
			theme:   styling.Theme{Colors: styling.ThemeCorporate.Colors, Fonts: styling.ThemeCorporate.Fonts},
			wantErr: "theme name is required",
		},
		{
			name: "Missing major font",
			theme: styling.Theme{
				Name:   "Test",
				Colors: styling.ThemeCorporate.Colors,
				Fonts:  styling.FontScheme{Name: "T", MinorFont: "Arial"},
			},
			wantErr: "major font is required",
		},
		{
			name: "Missing minor font",
			theme: styling.Theme{
				Name:   "Test",
				Colors: styling.ThemeCorporate.Colors,
				Fonts:  styling.FontScheme{Name: "T", MajorFont: "Arial"},
			},
			wantErr: "minor font is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := styling.ValidateTheme(tt.theme)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error %q should contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}
