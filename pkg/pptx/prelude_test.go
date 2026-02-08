package pptx

import (
	"testing"
)

func TestColorHelpers(t *testing.T) {
	if ColorRed != "FF0000" {
		t.Errorf("expected RED to be FF0000, got %s", ColorRed)
	}
	if ColorCorporateBlue != "1565C0" {
		t.Errorf("expected CORPORATE_BLUE to be 1565C0, got %s", ColorCorporateBlue)
	}
	if ColorMaterialRed != "F44336" {
		t.Errorf("expected MATERIAL_RED to be F44336, got %s", ColorMaterialRed)
	}
	if ColorCarbonBlue60 != "0043CE" {
		t.Errorf("expected CARBON_BLUE_60 to be 0043CE, got %s", ColorCarbonBlue60)
	}
}

func TestFontSizeHelpers(t *testing.T) {
	if FontSizeTitle != 44 {
		t.Errorf("expected FontSizeTitle to be 44, got %d", FontSizeTitle)
	}
	if FontSizeBody != 18 {
		t.Errorf("expected FontSizeBody to be 18, got %d", FontSizeBody)
	}
}

func TestThemeHelpers(t *testing.T) {
	themes := AllThemes()
	if len(themes) != 7 {
		t.Errorf("expected 7 themes, got %d", len(themes))
	}

	corporate := ThemeCorporate
	if corporate.Name != "Corporate" {
		t.Errorf("expected Corporate theme name, got %s", corporate.Name)
	}
	if corporate.Primary != "1565C0" {
		t.Errorf("expected Corporate primary color 1565C0, got %s", corporate.Primary)
	}

	dark := ThemeDark
	if dark.Background != "121212" {
		t.Errorf("expected Dark background color 121212, got %s", dark.Background)
	}
}

func TestUnitHelpers(t *testing.T) {
	// Verify existing unit helpers are still consistent
	if Inches(1.0) != 914400 {
		t.Errorf("expected 1 inch to be 914400 EMU, got %d", Inches(1.0))
	}
	if Centimeters(1.0) != 360000 {
		t.Errorf("expected 1 cm to be 360000 EMU, got %d", Centimeters(1.0))
	}
	if Points(1.0) != 12700 {
		t.Errorf("expected 1 pt to be 12700 EMU, got %d", Points(1.0))
	}
}
