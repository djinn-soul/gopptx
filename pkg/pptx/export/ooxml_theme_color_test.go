package export

import "testing"

func TestResolveOOXMLColorToken_Hex(t *testing.T) {
	r, g, b, ok := resolveOOXMLColorToken("#4F81BD")
	if !ok {
		t.Fatal("expected hex color to resolve")
	}
	if r != 0x4F || g != 0x81 || b != 0xBD {
		t.Fatalf("unexpected rgb: %d,%d,%d", r, g, b)
	}
}

func TestResolveOOXMLColorToken_SchemeAlias(t *testing.T) {
	r, g, b, ok := resolveOOXMLColorToken("tx1")
	if !ok {
		t.Fatal("expected tx1 alias to resolve")
	}
	if r != 0x00 || g != 0x00 || b != 0x00 {
		t.Fatalf("expected black for tx1 default, got %d,%d,%d", r, g, b)
	}
}

func TestResolveOOXMLColorToken_WithTransforms(t *testing.T) {
	baseR, baseG, baseB, ok := resolveOOXMLColorToken(themeColorAccent1)
	if !ok {
		t.Fatal("expected accent1 to resolve")
	}
	r, g, b, ok := resolveOOXMLColorToken(themeColorAccent1 + "|tint=40000|lumMod=65000|lumOff=35000")
	if !ok {
		t.Fatal("expected transformed accent1 to resolve")
	}
	if r <= baseR || g <= baseG || b <= baseB {
		t.Fatalf("expected transformed color to lighten base (%d,%d,%d -> %d,%d,%d)", baseR, baseG, baseB, r, g, b)
	}
}

func TestResolveOOXMLColorToken_Invalid(t *testing.T) {
	_, _, _, ok := resolveOOXMLColorToken("not-a-color")
	if ok {
		t.Fatal("expected invalid color token to fail resolution")
	}
}

func TestApplyColorTransforms_Shade(t *testing.T) {
	// Shade 50% of white should be middle gray (127, 127, 127 -> 7F, 7F, 7F)
	white := rgbColor{r: 255, g: 255, b: 255}
	res := shadeColor(white, 50000)
	if res.r != 127 || res.g != 127 || res.b != 127 {
		t.Errorf("expected 127,127,127 for 50%% shade of white, got %d,%d,%d", res.r, res.g, res.b)
	}
}

func TestApplyColorTransforms_Scale(t *testing.T) {
	// Scale 50% of 255
	res := scaleColor(255, 50000)
	if res != 127 {
		t.Errorf("expected 127 for 50%% scale of 255, got %d", res)
	}
}

func TestResolveColorAlias(t *testing.T) {
	tests := []struct {
		alias    string
		expected string
	}{
		{"tx1", "dk1"},
		{"bg1", "lt1"},
		{"tx2", "dk2"},
		{"bg2", "lt2"},
		{"unknown", "unknown"},
	}
	for _, tt := range tests {
		if got := resolveColorAlias(tt.alias); got != tt.expected {
			t.Errorf("resolveColorAlias(%q) = %q, want %q", tt.alias, got, tt.expected)
		}
	}
}

func TestResolveThemeBaseColor(t *testing.T) {
	res, ok := resolveThemeBaseColor(themeColorAccent1)
	if !ok || res.r != 0x4F || res.g != 0x81 || res.b != 0xBD {
		t.Errorf("expected accent1 default, got %v, %v", res, ok)
	}
	_, ok = resolveThemeBaseColor("unknown")
	if ok {
		t.Error("expected unknown color to fail")
	}
}

func TestNormalizeColorName(t *testing.T) {
	if got := normalizeColorName("  scheme:accent1  "); got != themeColorAccent1 {
		t.Errorf("expected accent1, got %q", got)
	}
	if got := normalizeColorName("ACCENT1"); got != themeColorAccent1 {
		t.Errorf("expected accent1, got %q", got)
	}
}

func TestLumModOffColor(t *testing.T) {
	c := rgbColor{r: 100, g: 100, b: 100}
	// lumMod=100000 (no change), lumOff=0 (no change)
	res := lumModOffColor(c, 100000, 0)
	if res.r != 100 {
		t.Errorf("expected 100, got %d", res.r)
	}

	// lumMod=0 should default to 100000
	res = lumModOffColor(c, 0, 0)
	if res.r != 100 {
		t.Errorf("expected 100 for 0 lumMod, got %d", res.r)
	}
}
