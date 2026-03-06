package shapes

import (
	"errors"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestNewFreeform(t *testing.T) {
	points := []FreeformPoint{
		{X: styling.Emu(0), Y: styling.Emu(0)},
		{X: styling.Emu(100000), Y: styling.Emu(0)},
		{X: styling.Emu(100000), Y: styling.Emu(100000)},
	}

	f := NewFreeform(points)

	if len(f.Points) != 3 {
		t.Errorf("expected 3 points, got %d", len(f.Points))
	}
	if !f.ClosePath {
		t.Errorf("expected ClosePath=true by default")
	}
}

func TestNewFreeformCoords(t *testing.T) {
	xCoords := []int64{0, 100000, 100000}
	yCoords := []int64{0, 0, 100000}

	f, err := NewFreeformCoords(xCoords, yCoords)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(f.Points) != 3 {
		t.Errorf("expected 3 points, got %d", len(f.Points))
	}
}

func TestNewFreeformCoordsUnequalLength(t *testing.T) {
	xCoords := []int64{0, 100000}
	yCoords := []int64{0, 100000, 50000}

	_, err := NewFreeformCoords(xCoords, yCoords)
	if err == nil {
		t.Error("expected error for unequal coordinate slices")
	}
}

func TestNewFreeformCoordsTooFewPoints(t *testing.T) {
	xCoords := []int64{0}
	yCoords := []int64{0}

	_, err := NewFreeformCoords(xCoords, yCoords)
	if err == nil {
		t.Error("expected error for fewer than 2 points")
	}
}

func TestNewFreeformInches(t *testing.T) {
	points := [][2]float64{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
	}

	f, err := NewFreeformInches(points)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(f.Points) != 4 {
		t.Errorf("expected 4 points, got %d", len(f.Points))
	}
}

func TestNewFreeformInchesTooFewPoints(t *testing.T) {
	points := [][2]float64{
		{0, 0},
	}

	_, err := NewFreeformInches(points)
	if err == nil {
		t.Error("expected error for fewer than 2 points")
	}
}

func TestNewFreeformClosed(t *testing.T) {
	points := []FreeformPoint{
		{X: 0, Y: 0},
		{X: 100, Y: 0},
	}

	f := NewFreeformClosed(points)

	if !f.ClosePath {
		t.Error("expected ClosePath=true")
	}
}

func TestNewFreeformOpen(t *testing.T) {
	points := []FreeformPoint{
		{X: 0, Y: 0},
		{X: 100, Y: 0},
	}

	f := NewFreeformOpen(points)

	if f.ClosePath {
		t.Error("expected ClosePath=false")
	}
}

func TestFreeformWithClosePath(t *testing.T) {
	points := []FreeformPoint{{X: 0, Y: 0}}
	f := NewFreeform(points).WithClosePath(false)

	if f.ClosePath {
		t.Error("expected ClosePath=false after WithClosePath(false)")
	}
}

func TestFreeformWithFill(t *testing.T) {
	fill := NewShapeFill("FF0000")
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}}).WithFill(fill)

	if f.Fill == nil {
		t.Error("expected Fill to be set")
	}
	if f.RichFill != nil {
		t.Error("expected RichFill to be nil when using legacy fill")
	}
}

func TestFreeformWithRichFill(t *testing.T) {
	richFill := NewSolidFill("00FF00")
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}}).WithRichFill(richFill)

	if f.RichFill == nil {
		t.Error("expected RichFill to be set")
	}
	if f.Fill != nil {
		t.Error("expected Fill to be nil when using rich fill")
	}
}

func TestFreeformWithLine(t *testing.T) {
	line := NewShapeLine("0000FF", styling.Points(2))
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}}).WithLine(line)

	if f.Line == nil {
		t.Error("expected Line to be set")
	}
	if f.RichLine != nil {
		t.Error("expected RichLine to be nil when using legacy line")
	}
}

func TestFreeformWithRichLine(t *testing.T) {
	richLine := NewRichShapeLine("00FFFF", styling.Points(2))
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}}).WithRichLine(richLine)

	if f.RichLine == nil {
		t.Error("expected RichLine to be set")
	}
	if f.Line != nil {
		t.Error("expected Line to be nil when using rich line")
	}
}

func TestFreeformWithRichShadow(t *testing.T) {
	shadow := NewOuterShadow("333333")
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}}).WithRichShadow(shadow)

	if f.RichShadow == nil {
		t.Error("expected RichShadow to be set")
	}
	if f.Effects == nil || !f.Effects.Shadow {
		t.Error("expected Effects.Shadow to be true")
	}
}

func TestFreeformWithText(t *testing.T) {
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}}).WithText("Hello")

	if f.Text != "Hello" {
		t.Errorf("expected text 'Hello', got %q", f.Text)
	}
}

func TestFreeformWithRotation(t *testing.T) {
	degrees := 45
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}}).WithRotation(degrees)

	if f.RotationDeg == nil || *f.RotationDeg != 45 {
		t.Error("expected rotation to be set to 45")
	}
}

func TestFreeformWithName(t *testing.T) {
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}}).WithName("MyFreeform")

	if f.Name != "MyFreeform" {
		t.Errorf("expected name 'MyFreeform', got %q", f.Name)
	}
}

func TestFreeformWithAltText(t *testing.T) {
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}}).WithAltText("Accessible shape")

	if f.AltText != "Accessible shape" {
		t.Errorf("expected alt text 'Accessible shape', got %q", f.AltText)
	}
}

func TestFreeformWithDecorative(t *testing.T) {
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}}).WithDecorative(true)

	if !f.IsDecorative {
		t.Error("expected IsDecorative to be true")
	}
}

func TestFreeformValidateTooFewPoints(t *testing.T) {
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}})

	err := f.Validate()
	if err == nil {
		t.Error("expected error for fewer than 2 points")
	}
}

// TestFreeformValidateConflictRichAndLegacyFill tests that fluent methods clear the other type.
// This test verifies that using WithRichFill after WithFill clears the legacy fill.
func TestFreeformValidateRichFillClearsLegacy(t *testing.T) {
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}, {X: 100, Y: 100}}).
		WithFill(NewShapeFill("FF0000")).
		WithRichFill(NewSolidFill("00FF00"))

	// WithRichFill clears Fill, so there's no conflict
	if f.Fill != nil {
		t.Error("expected Fill to be cleared after WithRichFill")
	}
	if f.RichFill == nil {
		t.Error("expected RichFill to be set")
	}

	err := f.Validate()
	if err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

// TestFreeformValidateConflictSolidAndGradientFill tests that fluent methods clear the other type.
func TestFreeformValidateGradientFillClearsLegacy(t *testing.T) {
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}, {X: 100, Y: 100}}).
		WithFill(NewShapeFill("FF0000")).
		WithGradientFill(NewShapeGradientFill("linear", []ShapeGradientStop{
			NewShapeGradientStop(0, "FF0000"),
			NewShapeGradientStop(100, "00FF00"),
		}))

	// WithGradientFill clears Fill, so there's no conflict
	if f.Fill != nil {
		t.Error("expected Fill to be cleared after WithGradientFill")
	}
	if f.GradientFill == nil {
		t.Error("expected GradientFill to be set")
	}

	err := f.Validate()
	if err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

// TestFreeformValidateRichLineClearsLegacy tests that fluent methods clear the other type.
func TestFreeformValidateRichLineClearsLegacy(t *testing.T) {
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}, {X: 100, Y: 100}}).
		WithLine(NewShapeLine("FF0000", styling.Points(1))).
		WithRichLine(NewRichShapeLine("00FF00", styling.Points(2)))

	// WithRichLine clears Line, so there's no conflict
	if f.Line != nil {
		t.Error("expected Line to be cleared after WithRichLine")
	}
	if f.RichLine == nil {
		t.Error("expected RichLine to be set")
	}

	err := f.Validate()
	if err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestFreeformValidateInvalidRichFill(t *testing.T) {
	// RichFill with invalid color
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}, {X: 100, Y: 100}}).
		WithRichFill(&RichShapeFill{Type: FillTypeSolid, Solid: &SolidFill{Color: "invalid"}})

	err := f.Validate()
	if err == nil {
		t.Error("expected error for invalid rich fill color")
	}
}

func TestFreeformValidateInvalidRichLine(t *testing.T) {
	// RichLine with invalid color
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}, {X: 100, Y: 100}}).
		WithRichLine(&RichShapeLine{Color: "invalid", Width: styling.Points(1)})

	err := f.Validate()
	if err == nil {
		t.Error("expected error for invalid rich line color")
	}
}

func TestFreeformValidateInvalidRichShadow(t *testing.T) {
	// RichShadow with invalid type
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}, {X: 100, Y: 100}}).
		WithRichShadow(&RichShapeShadow{Type: "invalid"})

	err := f.Validate()
	if err == nil {
		t.Error("expected error for invalid rich shadow type")
	}
}

func TestFreeformToShape(t *testing.T) {
	f := NewFreeform([]FreeformPoint{
		{X: styling.Emu(0), Y: styling.Emu(0)},
		{X: styling.Emu(100000), Y: styling.Emu(0)},
		{X: styling.Emu(100000), Y: styling.Emu(100000)},
	}).WithName("TestFreeform").
		WithFill(NewShapeFill("FF0000"))

	s := f.ToShape()

	// Should convert to a rectangle bounding the points
	if s.Name != "TestFreeform" {
		t.Errorf("expected name 'TestFreeform', got %q", s.Name)
	}
	if s.Fill == nil {
		t.Error("expected Fill to be preserved")
	}
}

func TestFreeformToShapeEmpty(t *testing.T) {
	f := NewFreeform([]FreeformPoint{})
	s := f.ToShape()

	// Empty freeform should return empty shape
	if s.X != 0 || s.Y != 0 || s.CX != 0 || s.CY != 0 {
		t.Errorf("expected empty shape bounds, got (%v,%v,%v,%v)", s.X, s.Y, s.CX, s.CY)
	}
}

// Test that validation passes when there are no conflicts.
func TestFreeformValidateNoConflict(t *testing.T) {
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}, {X: 100, Y: 100}}).
		WithFill(NewShapeFill("FF0000")).
		WithLine(NewShapeLine("0000FF", styling.Points(1)))

	err := f.Validate()
	if err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

// Test that validation passes with valid rich fill.
func TestFreeformValidateValidRichFill(t *testing.T) {
	f := NewFreeform([]FreeformPoint{{X: 0, Y: 0}, {X: 100, Y: 100}}).
		WithRichFill(NewSolidFill("FF0000"))

	err := f.Validate()
	if err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

// Test that validation correctly identifies specific error types.
func TestFreeformValidateSpecificErrors(t *testing.T) {
	tests := []struct {
		name        string
		freeform    Freeform
		validateErr error
	}{
		{
			name: "invalid rich fill color",
			freeform: NewFreeform([]FreeformPoint{{X: 0, Y: 0}, {X: 100, Y: 100}}).
				WithRichFill(&RichShapeFill{Type: FillTypeSolid, Solid: &SolidFill{Color: "xyz"}}),
			validateErr: errors.New("invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.freeform.Validate()
			if err == nil {
				t.Error("expected validation error")
			}
			if tt.validateErr != nil && err != nil && !strings.Contains(err.Error(), tt.validateErr.Error()) {
				t.Errorf("expected error containing %q, got %q", tt.validateErr.Error(), err.Error())
			}
		})
	}
}
