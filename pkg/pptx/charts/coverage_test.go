package charts_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestChartAccessibilityAndPosition(t *testing.T) {
	bar := charts.NewBarChart([]string{"A"}, []float64{10}).
		WithAltText("Alternative Text").
		WithDecorative(true).
		Position(styling.Inches(1), styling.Inches(1)).
		Size(styling.Inches(4), styling.Inches(3))

	if bar.AltText != "Alternative Text" {
		t.Errorf("Expected AltText 'Alternative Text', got %q", bar.AltText)
	}
	if !bar.IsDecorative {
		t.Error("Expected IsDecorative to be true")
	}
	if bar.X.Emu() != styling.Inches(1).Emu() || bar.Y.Emu() != styling.Inches(1).Emu() {
		t.Errorf("Unexpected position: %v, %v", bar.X, bar.Y)
	}
	if bar.CX.Emu() != styling.Inches(4).Emu() || bar.CY.Emu() != styling.Inches(3).Emu() {
		t.Errorf("Unexpected size: %v, %v", bar.CX, bar.CY)
	}

	line := charts.NewLineChart([]string{"A"}, []float64{10}).
		WithAltText("Line Alt").
		WithDecorative(false).
		Position(styling.Emu(0), styling.Emu(0)).
		Size(styling.Emu(100), styling.Emu(100))

	if line.AltText != "Line Alt" {
		t.Errorf("Expected AltText 'Line Alt', got %q", line.AltText)
	}
	if line.IsDecorative {
		t.Error("Expected IsDecorative to be false")
	}

	// Test Getters
	if len(bar.GetCategories()) != 1 || bar.GetCategories()[0] != "A" {
		t.Errorf("Unexpected Bar Categories: %v", bar.GetCategories())
	}
	if len(bar.GetValues()) != 1 || bar.GetValues()[0] != 10 {
		t.Errorf("Unexpected Bar Values: %v", bar.GetValues())
	}
	if len(line.GetCategories()) != 1 || line.GetCategories()[0] != "A" {
		t.Errorf("Unexpected Line Categories: %v", line.GetCategories())
	}
	if len(line.GetValues()) != 1 || line.GetValues()[0] != 10 {
		t.Errorf("Unexpected Line Values: %v", line.GetValues())
	}
}

func TestChartColorNormalization(t *testing.T) {
	bar := charts.NewBarChart([]string{"A"}, []float64{10}).WithBarColor("#ff0000")
	if bar.BarColor != "FF0000" {
		t.Errorf("Expected FF0000, got %s", bar.BarColor)
	}

	line := charts.NewLineChart([]string{"A"}, []float64{10}).WithLineColor("00ff00")
	if line.LineColor != "00FF00" {
		t.Errorf("Expected 00FF00, got %s", line.LineColor)
	}
}

func TestChartValidationErrors(t *testing.T) {
	// Alt text too long
	longAlt := strings.Repeat("A", 1025)
	bar := charts.NewBarChart([]string{"A"}, []float64{10}).WithAltText(longAlt)
	if err := bar.Validate(1); err == nil {
		t.Error("Expected error for long alt text")
	}

	line := charts.NewLineChart([]string{"A"}, []float64{10}).WithAltText(longAlt)
	if err := line.Validate(1); err == nil {
		t.Error("Expected error for long alt text")
	}

	// Negative position
	barPos := charts.NewBarChart([]string{"A"}, []float64{10}).Position(styling.Emu(-1), 0)
	if err := barPos.Validate(1); err == nil {
		t.Error("Expected error for negative position")
	}

	// Invalid size
	barSize := charts.NewBarChart([]string{"A"}, []float64{10}).Size(0, 100)
	if err := barSize.Validate(1); err == nil {
		t.Error("Expected error for zero size")
	}

	// Empty title
	barTitle := charts.NewBarChart([]string{"A"}, []float64{10}).WithTitle("")
	if err := barTitle.Validate(1); err == nil {
		t.Error("Expected error for empty title")
	}

	// Mismatched lengths
	barMismatch := charts.NewBarChart([]string{"A", "B"}, []float64{10})
	if err := barMismatch.Validate(1); err == nil {
		t.Error("Expected error for mismatched lengths")
	}

	// Empty category
	barEmptyCat := charts.NewBarChart([]string{""}, []float64{10})
	if err := barEmptyCat.Validate(1); err == nil {
		t.Error("Expected error for empty category")
	}

	// Invalid color
	barColor := charts.NewBarChart([]string{"A"}, []float64{10}).WithBarColor("not-a-color")
	if err := barColor.Validate(1); err == nil {
		t.Error("Expected error for invalid color")
	}

	// Value range error
	minValue, maxValue := 10.0, 5.0
	barRange := charts.NewBarChart([]string{"A"}, []float64{10})
	barRange.MinValue = &minValue
	barRange.MaxValue = &maxValue
	if err := barRange.Validate(1); err == nil {
		t.Error("Expected error for invalid value range")
	}
}

func TestChartSpecAccessibility(t *testing.T) {
	bar := charts.NewBarChart([]string{"A"}, []float64{10}).
		WithAltText("Alt").
		WithDecorative(true)
	spec := bar.ToChartSpec()
	if spec.AltText != "Alt" || !spec.IsDecorative {
		t.Errorf("Spec missing accessibility: %+v", spec)
	}

	line := charts.NewLineChart([]string{"A"}, []float64{10}).
		WithAltText("Alt Line").
		WithDecorative(true)
	lineSpec := line.ToChartSpec()
	if lineSpec.AltText != "Alt Line" || !lineSpec.IsDecorative {
		t.Errorf("Line spec missing accessibility: %+v", lineSpec)
	}
}
