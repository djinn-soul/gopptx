package enums

import "testing"

func TestAllPPAlign(t *testing.T) {
	values := []string{"left", "center", "right", "justify", "distribute", "thaidist", "justify_low"}
	for _, v := range values {
		if _, err := ParsePPAlign(v); err != nil {
			t.Errorf("failed to parse PPAlign %q: %v", v, err)
		}
	}
	// Coverage for aliases and XMLValue
	if PPAlignLeft.XMLValue() != "l" {
		t.Error("PPAlignLeft XML mismatch")
	}
	ParsePPAlign("centre")
	ParsePPAlign("justifylow")
	ParsePPAlign("thai_distribute")
}

func TestAllMSOAnchor(t *testing.T) {
	values := []string{"top", "middle", "bottom", "justify", "distribute"}
	for _, v := range values {
		if _, err := ParseMSOAnchor(v); err != nil {
			t.Errorf("failed to parse MSOAnchor %q: %v", v, err)
		}
	}
	if MSOAnchorTop.XMLValue() != "t" {
		t.Error("MSOAnchorTop XML mismatch")
	}
	if MSOVerticalAnchorTop.XMLValue() != "t" {
		t.Error("MSOVerticalAnchorTop XML mismatch")
	}
	ParseMSOVerticalAnchor("middle")
}

func TestAllMSOThemeColor(t *testing.T) {
	values := []string{"bg1", "tx1", "bg2", "tx2", "accent1", "accent2", "accent3", "accent4", "accent5", "accent6", "hlink", "folhlink"}
	for _, v := range values {
		if _, err := ParseMSOThemeColor(v); err != nil {
			t.Errorf("failed to parse MSOThemeColor %q: %v", v, err)
		}
	}
	// Test aliases
	aliases := []string{"background1", "text1", "background2", "text2", "hyperlink", "followedhyperlink"}
	for _, v := range aliases {
		if _, err := ParseMSOThemeColor(v); err != nil {
			t.Errorf("failed to parse MSOThemeColor alias %q: %v", v, err)
		}
	}
	if MSOThemeColorAccent1.XMLValue() != "accent1" {
		t.Error("MSOThemeColorAccent1 XML mismatch")
	}
	if _, err := ParseMSOThemeColor("invalid"); err == nil {
		t.Error("expected error for invalid theme color")
	}
}

func TestAllMSOColorType(t *testing.T) {
	values := []string{"rgb", "theme", "auto", "unknown"}
	for _, v := range values {
		if _, err := ParseMSOColorType(v); err != nil {
			t.Errorf("failed to parse MSOColorType %q: %v", v, err)
		}
	}
	if MSOColorTypeRGB.XMLValue() != "rgb" {
		t.Error("MSOColorTypeRGB XML mismatch")
	}
	if _, err := ParseMSOColorType("invalid"); err == nil {
		t.Error("expected error for invalid color type")
	}
}

func TestAllMSOShape(t *testing.T) {
	values := []string{"rectangle", "rounded_rectangle", "oval", "circle", "diamond", "triangle", "right_triangle"}
	for _, v := range values {
		s, err := ParseMSOShape(v)
		if err != nil {
			t.Errorf("failed to parse MSOShape %q: %v", v, err)
		}
		if s.XMLValue() == "" {
			t.Error("MSOShape XML mismatch")
		}
	}
	if _, err := ParseMSOShape("invalid"); err == nil {
		t.Error("expected error for invalid shape")
	}
}

func TestAllXLChartType(t *testing.T) {
	values := []string{"bar", "barhorizontal", "barstacked", "barstacked100", "line", "linemarkers", "linestacked", "scatter", "xy", "area", "areastacked", "areastacked100", "pie", "doughnut", "bubble", "radar", "radarfilled", "stockhlc", "stockohlc", "combo"}
	for _, v := range values {
		if _, err := ParseXLChartType(v); err != nil {
			t.Errorf("failed to parse XLChartType %q: %v", v, err)
		}
	}
	// aliases
	ParseXLChartType("donut")
	if XLChartTypeBar.XMLValue() == "" {
		t.Error("XLChartTypeBar XML mismatch")
	}
	if _, err := ParseXLChartType("invalid"); err == nil {
		t.Error("expected error for invalid chart type")
	}
}
