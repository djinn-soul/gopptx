package shapes

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestShape_Builders(t *testing.T) {
	s := NewShape(ShapeTypeRectangle, 0, 0, styling.Inches(1), styling.Inches(1))

	s = s.WithText("Hello").
		WithRotation(45).
		WithAltText("Alt text").
		WithDecorative(true).
		WithName("MyShape").
		WithAdjustmentValue("adj1", 50000).
		WithAdjustment("adj2", "val 10000")

	if s.Text != "Hello" {
		t.Errorf("expected text Hello, got %s", s.Text)
	}
	if *s.RotationDeg != 45 {
		t.Errorf("expected rotation 45, got %d", *s.RotationDeg)
	}
	if s.AltText != "Alt text" {
		t.Errorf("expected alt text Alt text, got %s", s.AltText)
	}
	if !s.IsDecorative {
		t.Error("expected decorative true")
	}
	if s.Name != "MyShape" {
		t.Errorf("expected name MyShape, got %s", s.Name)
	}
	if len(s.Adjustments) != 2 {
		t.Errorf("expected 2 adjustments, got %d", len(s.Adjustments))
	}

	// Text frame builders
	s = s.WithTextMargins(styling.Emu(10), styling.Emu(20), styling.Emu(30), styling.Emu(40)).
		WithVerticalAnchor(TextAnchorBottom).
		WithTextWrap(TextWrapNone).
		WithAutoFit(TextAutoFitNormal)

	if s.TextFrame.MarginLeft != styling.Emu(10) {
		t.Errorf("expected margin left 10, got %v", s.TextFrame.MarginLeft)
	}
	if s.TextFrame.Anchor != TextAnchorBottom {
		t.Errorf("expected anchor bottom, got %s", s.TextFrame.Anchor)
	}
	if s.TextFrame.Wrap != TextWrapNone {
		t.Errorf("expected wrap none, got %s", s.TextFrame.Wrap)
	}
	if s.TextFrame.AutoFit != TextAutoFitNormal {
		t.Errorf("expected autofit normal, got %s", s.TextFrame.AutoFit)
	}

	textFrame := NewTextFrame().WithRotation(45)
	if textFrame.RotationDeg == nil || *textFrame.RotationDeg != 45 {
		t.Errorf("expected text-frame rotation 45, got %#v", textFrame.RotationDeg)
	}

	// Action builders
	link := action.NewHyperlink(action.HyperlinkURL("https://example.com"))
	s = s.WithClickAction(link).WithHoverAction(link)
	if s.ClickAction.Action.URL != "https://example.com" {
		t.Error("expected click action target")
	}
	if s.HoverAction.Action.URL != "https://example.com" {
		t.Error("expected hover action target")
	}

	// Hyperlink legacy builder
	s = s.WithHyperlink(link)
	if s.ClickAction.Action.URL != "https://example.com" {
		t.Error("expected click action target from WithHyperlink")
	}

	// Rich builders
	richFill := NewSolidFill("FF0000")
	s = s.WithRichFill(richFill)
	if s.RichFill.Type != FillTypeSolid {
		t.Error("expected rich fill solid")
	}
	if s.Fill != nil || s.GradientFill != nil {
		t.Error("legacy fills should be cleared when rich fill is set")
	}

	richLine := &RichShapeLine{}
	s = s.WithRichLine(richLine)
	if s.RichLine == nil {
		t.Error("expected rich line")
	}
	if s.Line != nil {
		t.Error("legacy line should be cleared when rich line is set")
	}

	richShadow := &RichShapeShadow{}
	s = s.WithRichShadow(richShadow)
	if s.RichShadow == nil {
		t.Error("expected rich shadow")
	}
	if s.Effects == nil || !s.Effects.Shadow {
		t.Error("expected shadow effect to be true")
	}
}

func TestShape_Validation(t *testing.T) {
	tests := []struct {
		name    string
		shape   Shape
		wantErr bool
	}{
		{
			"Valid shape",
			NewShape(ShapeTypeRectangle, 0, 0, styling.Inches(1), styling.Inches(1)),
			false,
		},
		{
			"Negative position",
			NewShape(ShapeTypeRectangle, -1, 0, styling.Inches(1), styling.Inches(1)),
			true,
		},
		{
			"Zero size",
			NewShape(ShapeTypeRectangle, 0, 0, 0, styling.Inches(1)),
			true,
		},
		{
			"Invalid shape type",
			NewShape("invalid", 0, 0, styling.Inches(1), styling.Inches(1)),
			true,
		},
		{
			"Invalid rotation",
			NewShape(ShapeTypeRectangle, 0, 0, styling.Inches(1), styling.Inches(1)).WithRotation(500),
			true,
		},
		{
			"Conflict rich/legacy fill",
			Shape{
				Type:     ShapeTypeRectangle,
				X:        0,
				Y:        0,
				CX:       100,
				CY:       100,
				RichFill: NewSolidFill("FF0000"),
				Fill:     &ShapeFill{Color: "00FF00"},
			},
			true,
		},
		{
			"Conflict legacy solid/gradient",
			Shape{
				Type:         ShapeTypeRectangle,
				X:            0,
				Y:            0,
				CX:           100,
				CY:           100,
				Fill:         &ShapeFill{Color: "00FF00"},
				GradientFill: &ShapeGradientFill{Type: ShapeGradientTypeLinear},
			},
			true,
		},
		{
			"Invalid legacy fill color",
			NewShape(ShapeTypeRectangle, 0, 0, 100, 100).WithFill(ShapeFill{Color: "ZZZZZZ"}),
			true,
		},
		{
			"Invalid rotation negative",
			NewShape(ShapeTypeRectangle, 0, 0, 100, 100).WithRotation(-400),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.shape.Validate(0, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShapeFill_Validation(t *testing.T) {
	fill := NewShapeFill("FF0000")
	if err := fill.Validate(); err != nil {
		t.Errorf("valid fill failed: %v", err)
	}

	fill.Color = "invalid"
	if err := fill.Validate(); err == nil {
		t.Error("expected error for invalid color")
	}

	fill = NewShapeFill("FF0000").WithTransparency(1.5)
	if err := fill.Validate(); err == nil {
		t.Error("expected error for transparency > 1")
	}
}

func TestShapeLine_BuildersAndValidation(t *testing.T) {
	line := NewShapeLine("FF0000", styling.Emu(100)).
		WithDash(LineDashDot).
		WithCap(LineCapRound).
		WithJoin(LineJoinBevel)

	if line.Dash != LineDashDot {
		t.Errorf("expected dash %s, got %s", LineDashDot, line.Dash)
	}
	if line.Cap != LineCapRound {
		t.Errorf("expected cap %s, got %s", LineCapRound, line.Cap)
	}
	if line.Join != LineJoinBevel {
		t.Errorf("expected join %s, got %s", LineJoinBevel, line.Join)
	}

	if err := line.Validate(); err != nil {
		t.Errorf("valid line failed: %v", err)
	}

	line.Width = 0
	if err := line.Validate(); err == nil {
		t.Error("expected error for zero width")
	}
}

func TestShapeGradient_Validation(t *testing.T) {
	stops := []ShapeGradientStop{
		NewShapeGradientStop(0, "FF0000"),
		NewShapeGradientStop(100, "0000FF"),
	}
	grad := NewShapeGradientFill(ShapeGradientTypeLinear, stops)
	if err := grad.Validate(); err != nil {
		t.Errorf("valid gradient failed: %v", err)
	}

	// Test non-strictly increasing positions
	badStops := []ShapeGradientStop{
		NewShapeGradientStop(50, "FF0000"),
		NewShapeGradientStop(50, "0000FF"),
	}
	grad = NewShapeGradientFill(ShapeGradientTypeLinear, badStops)
	if err := grad.Validate(); err == nil {
		t.Error("expected error for non-strictly increasing positions")
	}

	// Test angle for non-linear gradient
	grad = NewShapeGradientFill(ShapeGradientTypeRadial, stops).WithLinearAngle(45)
	if err := grad.Validate(); err == nil {
		t.Error("expected error for angle on non-linear gradient")
	}
}

func TestRichShapeFill_Validation(t *testing.T) {
	tests := []struct {
		name    string
		fill    *RichShapeFill
		wantErr bool
	}{
		{"Nil fill", nil, false},
		{"Valid solid", NewSolidFill("FF0000"), false},
		{"Valid no-fill", NewNoFill(), false},
		{"Valid pattern", NewPatternFill(PatternTypeHorz), false},
		{"Invalid solid color", NewSolidFill("invalid"), true},
		{"Invalid transparency", NewSolidFill("FF0000").WithTransparency(1.5), true},
		{"Solid with nil Solid", &RichShapeFill{Type: FillTypeSolid}, true},
		{"Gradient with nil Gradient", &RichShapeFill{Type: FillTypeGradient}, true},
		{"Pattern with nil Pattern", &RichShapeFill{Type: FillTypePattern}, true},
		{
			"Invalid pattern colors",
			NewPatternFill(PatternTypeHorz).WithPatternColors("invalid", "FFFFFF"),
			true,
		},
		{"Unknown fill type", &RichShapeFill{Type: "unknown"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fill.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRichShapeFill_Builders(t *testing.T) {
	f := &RichShapeFill{}
	f = f.WithSolid("FF0000").WithTransparency(0.5)
	if f.Type != FillTypeSolid || f.Solid.Transparency != 0.5 {
		t.Error("WithSolid/WithTransparency failed")
	}

	f = f.WithGradient(NewShapeGradientFill(ShapeGradientTypeLinear, nil))
	if f.Type != FillTypeGradient || f.Gradient == nil {
		t.Error("WithGradient failed")
	}

	f = f.WithPattern(PatternTypeVert).WithPatternColors("111111", "222222")
	if f.Type != FillTypePattern || f.Pattern.FgColor != "111111" {
		t.Error("WithPattern failed")
	}

	f = f.Background()
	if f.Type != FillTypeNoFill {
		t.Error("Background failed")
	}

	f = f.Foreground()
	if f.Type != FillTypeSolid {
		t.Error("Foreground failed")
	}
}

func TestPatternType_Helpers(t *testing.T) {
	if !IsValidPatternType(PatternTypeHorz) {
		t.Error("Horz should be valid")
	}
	if IsValidPatternType("invalid") {
		t.Error("invalid should be invalid")
	}
	if NormalizePatternType("horz") != PatternTypeHorz {
		t.Error("normalize horz failed")
	}
	if NormalizePatternType("invalid") != PatternTypePct5 {
		t.Error("normalize invalid failed")
	}
}
