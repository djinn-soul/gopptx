package export

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestSanitizeTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "Hello_World"},
		{"File!@#$%^Name", "File______Name"},
		{"Already-Safe_123", "Already-Safe_123"},
		{"", "presentation"},
	}
	for _, tt := range tests {
		if got := sanitizeTitle(tt.input); got != tt.expected {
			t.Errorf("sanitizeTitle(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		input   string
		r, g, b uint8
	}{
		{"FF0000", 255, 0, 0},
		{"#00FF00", 0, 255, 0},
		{"0000FF", 0, 0, 255},
		{"invalid", 0, 0, 0},
	}
	for _, tt := range tests {
		r, g, b := hexToRGB(tt.input)
		if r != tt.r || g != tt.g || b != tt.b {
			t.Errorf("hexToRGB(%q) = %d,%d,%d, want %d,%d,%d", tt.input, r, g, b, tt.r, tt.g, tt.b)
		}
	}
}

func TestRomanNumeral(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{1, "I"},
		{4, "IV"},
		{9, "IX"},
		{10, "X"},
		{40, "XL"},
		{90, "XC"},
		{100, "C"},
		{400, "CD"},
		{900, "CM"},
		{1000, "M"},
		{2024, "MMXXIV"},
		{0, ""},
		{-1, ""},
	}
	for _, tt := range tests {
		if got := romanNumeral(tt.input); got != tt.expected {
			t.Errorf("romanNumeral(%d) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestInterpolateGradient(t *testing.T) {
	stops := []gradientStop{
		{pos: 0, color: rgbColor{r: 0, g: 0, b: 0}, alpha: 1},
		{pos: 1, color: rgbColor{r: 255, g: 255, b: 255}, alpha: 1},
	}

	// At 0.5, should be middle gray (127.5 rounded to 128)
	res := interpolateGradient(stops, 0.5)
	if res.color.r != 128 || res.color.g != 128 || res.color.b != 128 {
		t.Errorf("expected 128,128,128, got %d,%d,%d", res.color.r, res.color.g, res.color.b)
	}
}

func TestBlendOverWhite(t *testing.T) {
	c := gradientStop{
		color: rgbColor{r: 255, g: 0, b: 0}, // Red
	}
	// 50% alpha over white should be pink (127.5 rounded to 128)
	res := blendOverWhite(c, 0.5)
	if res.r != 255 || res.g != 128 || res.b != 128 {
		t.Errorf("expected 255,128,128, got %d,%d,%d", res.r, res.g, res.b)
	}
}

func TestGradientStopsFromFill(t *testing.T) {
	grad := &shapes.ShapeGradientFill{
		Stops: []shapes.ShapeGradientStop{
			shapes.NewShapeGradientStop(0, "FF0000"),
			shapes.NewShapeGradientStop(100, "0000FF"),
		},
	}
	stops := gradientStopsFromFill(grad)
	if len(stops) != 2 {
		t.Errorf("expected 2 stops, got %d", len(stops))
	}
	if stops[0].pos != 0 || stops[1].pos != 1 {
		t.Errorf("expected positions 0 and 1, got %f and %f", stops[0].pos, stops[1].pos)
	}
}

func TestIsMostlyVerticalGradient(t *testing.T) {
	if !isMostlyVerticalGradient(90) {
		t.Error("90 should be vertical")
	}
	if !isMostlyVerticalGradient(270) {
		t.Error("270 should be vertical")
	}
	if isMostlyVerticalGradient(0) {
		t.Error("0 should be horizontal")
	}
	if isMostlyVerticalGradient(180) {
		t.Error("180 should be horizontal")
	}
	// Edge cases
	if isMostlyVerticalGradient(45) {
		t.Error("45 should be horizontal")
	}
	if !isMostlyVerticalGradient(46) {
		t.Error("46 should be vertical")
	}
}

func TestDrawStyle(t *testing.T) {
	if got := drawStyle(true, true); got != "DF" {
		t.Errorf("expected DF, got %q", got)
	}
	if got := drawStyle(true, false); got != "F" {
		t.Errorf("expected F, got %q", got)
	}
	if got := drawStyle(false, true); got != "D" {
		t.Errorf("expected D, got %q", got)
	}
	if got := drawStyle(false, false); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestStripHash(t *testing.T) {
	if got := stripHash("#ABC"); got != "ABC" {
		t.Errorf("expected ABC, got %q", got)
	}
	if got := stripHash("ABC"); got != "ABC" {
		t.Errorf("expected ABC, got %q", got)
	}
}

func TestNormalizePDFDriver(t *testing.T) {
	d, err := normalizePDFDriver(PDFOptions{Driver: ""})
	if err != nil || d != PDFDriverAuto {
		t.Error("Auto failed")
	}

	d, err = normalizePDFDriver(PDFOptions{Driver: "native"})
	if err != nil || d != PDFDriverNative {
		t.Error("Native failed")
	}

	_, err = normalizePDFDriver(PDFOptions{Driver: "invalid"})
	if err == nil {
		t.Error("expected error for invalid driver")
	}
}
