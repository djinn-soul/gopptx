package styling_test

import (
	"math"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestColorHelpers(t *testing.T) {
	if pptx.ColorRed != "FF0000" {
		t.Errorf("expected RED to be FF0000, got %s", pptx.ColorRed)
	}
	if pptx.ColorCorporateBlue != "1565C0" {
		t.Errorf("expected CORPORATE_BLUE to be 1565C0, got %s", pptx.ColorCorporateBlue)
	}
	if pptx.ColorMaterialRed != "F44336" {
		t.Errorf("expected MATERIAL_RED to be F44336, got %s", pptx.ColorMaterialRed)
	}
	if pptx.ColorCarbonBlue60 != "0043CE" {
		t.Errorf("expected CARBON_BLUE_60 to be 0043CE, got %s", pptx.ColorCarbonBlue60)
	}
}

func TestFontSizeHelpers(t *testing.T) {
	if pptx.FontSizeTitle != 44 {
		t.Errorf("expected FontSizeTitle to be 44, got %d", pptx.FontSizeTitle)
	}
	if pptx.FontSizeBody != 18 {
		t.Errorf("expected FontSizeBody to be 18, got %d", pptx.FontSizeBody)
	}
}

func TestThemeHelpers(t *testing.T) {
	themes := pptx.AllThemes()
	if len(themes) != 7 {
		t.Errorf("expected 7 themes, got %d", len(themes))
	}

	corporate := pptx.ThemeCorporate
	if corporate.Name != "Corporate" {
		t.Errorf("expected Corporate theme name, got %s", corporate.Name)
	}
	if corporate.Primary != "1565C0" {
		t.Errorf("expected Corporate primary color 1565C0, got %s", corporate.Primary)
	}

	dark := pptx.ThemeDark
	if dark.Background != "121212" {
		t.Errorf("expected Dark background color 121212, got %s", dark.Background)
	}
}

func TestUnitHelpers(t *testing.T) {
	if pptx.Inches(1.0) != 914400 {
		t.Errorf("expected 1 inch to be 914400 EMU, got %d", pptx.Inches(1.0))
	}
	if pptx.Centimeters(1.0) != 360000 {
		t.Errorf("expected 1 cm to be 360000 EMU, got %d", pptx.Centimeters(1.0))
	}
	if pptx.Points(1.0) != 12700 {
		t.Errorf("expected 1 pt to be 12700 EMU, got %d", pptx.Points(1.0))
	}
}

func TestUnitConvertersOverflow(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		fn       func(float64) pptx.Length
		expected pptx.Length
	}{
		{"Inches Overflow", 2e15, pptx.Inches, styling.MaxEMU},
		{"Inches Underflow", -2e15, pptx.Inches, -styling.MaxEMU},
		{"CM Overflow", 5e15, pptx.Centimeters, styling.MaxEMU},
		{"CM Underflow", -5e15, pptx.Centimeters, -styling.MaxEMU},
		{"Points Overflow", 1e16, pptx.Points, styling.MaxEMU},
		{"Points Underflow", -1e16, pptx.Points, -styling.MaxEMU},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.input)
			if got != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, got)
			}
		})
	}
}

func TestFontSizeConverter(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected int
	}{
		{"Standard 12pt", 12.0, 1200},
		{"Small 10.5pt", 10.5, 1050},
		{"Large 72pt", 72, 7200},
		{"Overflow", 1e10, math.MaxInt32},
		{"Underflow", -1e10, math.MinInt32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pptx.FontSize(tt.input)
			if got != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, got)
			}
		})
	}
}
