package pptx

import (
	"math"
	"testing"
)

func TestUnitConvertersOverflow(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		fn       func(float64) int64
		expected int64
	}{
		{"Inches Overflow", 2e15, Inches, int64(maxEMU)},
		{"Inches Underflow", -2e15, Inches, int64(-maxEMU)},
		{"CM Overflow", 5e15, Centimeters, int64(maxEMU)},
		{"CM Underflow", -5e15, Centimeters, int64(-maxEMU)},
		{"Points Overflow", 1e16, Points, int64(maxEMU)},
		{"Points Underflow", -1e16, Points, int64(-maxEMU)},
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
			got := FontSize(tt.input)
			if got != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, got)
			}
		})
	}
}
