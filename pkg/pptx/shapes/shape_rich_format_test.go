package shapes

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestRichSolidFill(t *testing.T) {
	fill := NewSolidFill("FF0000").
		WithTransparency(0.5)

	if fill.Type != FillTypeSolid {
		t.Errorf("expected fill type solid, got %s", fill.Type)
	}
	if fill.Solid.Color != "FF0000" {
		t.Errorf("expected color FF0000, got %s", fill.Solid.Color)
	}
	if fill.Solid.Transparency != 0.5 {
		t.Errorf("expected transparency 0.5, got %f", fill.Solid.Transparency)
	}

	if err := fill.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}
}

func TestRichPatternFill(t *testing.T) {
	fill := NewPatternFill(PatternTypeDiagCross).
		WithPatternColors("000000", "FFFFFF")

	if fill.Type != FillTypePattern {
		t.Errorf("expected fill type pattern, got %s", fill.Type)
	}
	if fill.Pattern.Pattern != PatternTypeDiagCross {
		t.Errorf("expected pattern diagCross, got %s", fill.Pattern.Pattern)
	}
	if fill.Pattern.FgColor != "000000" {
		t.Errorf("expected fg color 000000, got %s", fill.Pattern.FgColor)
	}
	if fill.Pattern.BgColor != "FFFFFF" {
		t.Errorf("expected bg color FFFFFF, got %s", fill.Pattern.BgColor)
	}

	if err := fill.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}
}

func TestRichNoFill(t *testing.T) {
	fill := NewNoFill()

	if fill.Type != FillTypeNoFill {
		t.Errorf("expected fill type noFill, got %s", fill.Type)
	}

	if err := fill.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}
}

func TestRichShapeLine(t *testing.T) {
	line := NewRichShapeLine("00FF00", styling.Points(2)).
		WithDashStyle(LineDashStyleDash).
		WithCapStyle(LineCapStyleRound).
		WithJoinStyle(LineJoinStyleBevel).
		WithTransparency(0.3)

	if line.Color != "00FF00" {
		t.Errorf("expected color 00FF00, got %s", line.Color)
	}
	if line.DashStyle != LineDashStyleDash {
		t.Errorf("expected dash style dash, got %s", line.DashStyle)
	}
	if line.CapStyle != LineCapStyleRound {
		t.Errorf("expected cap style round, got %s", line.CapStyle)
	}
	if line.JoinStyle != LineJoinStyleBevel {
		t.Errorf("expected join style bevel, got %s", line.JoinStyle)
	}
	if line.Transparency != 0.3 {
		t.Errorf("expected transparency 0.3, got %f", line.Transparency)
	}

	if err := line.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}
}

func TestRichOuterShadow(t *testing.T) {
	shadow := NewOuterShadow("000000").
		WithTransparency(0.6).
		WithBlurRadius(50000).
		WithDistance(30000).
		WithAngle(45).
		WithAlignment(ShadowAlignBottomRight)

	if shadow.Type != ShadowTypeOuter {
		t.Errorf("expected shadow type outer, got %s", shadow.Type)
	}
	if shadow.Color != "000000" {
		t.Errorf("expected color 000000, got %s", shadow.Color)
	}
	if shadow.Transparency != 0.6 {
		t.Errorf("expected transparency 0.6, got %f", shadow.Transparency)
	}
	if shadow.BlurRadius != 50000 {
		t.Errorf("expected blur radius 50000, got %d", shadow.BlurRadius)
	}
	if shadow.Distance != 30000 {
		t.Errorf("expected distance 30000, got %d", shadow.Distance)
	}
	if shadow.Angle != 45 {
		t.Errorf("expected angle 45, got %f", shadow.Angle)
	}
	if shadow.Alignment != ShadowAlignBottomRight {
		t.Errorf("expected alignment bottom right, got %s", shadow.Alignment)
	}

	if err := shadow.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}
}

func TestRichInnerShadow(t *testing.T) {
	shadow := NewInnerShadow("FF0000").
		WithBlurRadius(40000).
		WithDistance(20000)

	if shadow.Type != ShadowTypeInner {
		t.Errorf("expected shadow type inner, got %s", shadow.Type)
	}
	if shadow.Color != "FF0000" {
		t.Errorf("expected color FF0000, got %s", shadow.Color)
	}

	if err := shadow.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}
}

func TestShapeWithRichFill(t *testing.T) {
	fill := NewSolidFill("4472C4")
	shape := NewShape(ShapeTypeRectangle, 0, 0, 1000000, 1000000).
		WithRichFill(fill).
		WithText("Rich Fill Shape")

	if shape.RichFill == nil {
		t.Error("expected rich fill to be set")
	}
	if shape.Fill != nil {
		t.Error("expected legacy fill to be cleared")
	}
	if shape.RichFill.Solid.Color != "4472C4" {
		t.Errorf("expected fill color 4472C4, got %s", shape.RichFill.Solid.Color)
	}

	if err := shape.Validate(1, 1); err != nil {
		t.Errorf("shape validation failed: %v", err)
	}
}

func TestShapeWithRichLine(t *testing.T) {
	line := NewRichShapeLine("FF0000", styling.Points(3))
	shape := NewShape(ShapeTypeRectangle, 0, 0, 1000000, 1000000).
		WithRichLine(line)

	if shape.RichLine == nil {
		t.Error("expected rich line to be set")
	}
	if shape.Line != nil {
		t.Error("expected legacy line to be cleared")
	}

	if err := shape.Validate(1, 1); err != nil {
		t.Errorf("shape validation failed: %v", err)
	}
}

func TestShapeWithRichShadow(t *testing.T) {
	shadow := NewOuterShadow("000000").
		WithTransparency(0.5).
		WithDistance(20000)

	shape := NewShape(ShapeTypeRectangle, 0, 0, 1000000, 1000000).
		WithRichShadow(shadow)

	if shape.RichShadow == nil {
		t.Error("expected rich shadow to be set")
	}
	if !shape.Effects.Shadow {
		t.Error("expected effects shadow flag to be set")
	}

	if err := shape.Validate(1, 1); err != nil {
		t.Errorf("shape validation failed: %v", err)
	}
}

func TestRichFillValidation(t *testing.T) {
	// Invalid color
	fill := &RichShapeFill{
		Type: FillTypeSolid,
		Solid: &SolidFill{
			Color:        "invalid",
			Transparency: 0.0,
		},
	}
	if err := fill.Validate(); err == nil {
		t.Error("expected validation error for invalid color")
	}

	// Invalid transparency
	fill = NewSolidFill("FF0000").WithTransparency(1.5)
	if err := fill.Validate(); err == nil {
		t.Error("expected validation error for transparency > 1")
	}
}

func TestRichLineValidation(t *testing.T) {
	// Invalid color
	line := NewRichShapeLine("invalid", styling.Points(1))
	if err := line.Validate(); err == nil {
		t.Error("expected validation error for invalid color")
	}

	// Invalid transparency
	line = NewRichShapeLine("FF0000", styling.Points(1)).WithTransparency(-0.5)
	if err := line.Validate(); err == nil {
		t.Error("expected validation error for negative transparency")
	}
}

func TestRichShadowValidation(t *testing.T) {
	// Invalid color
	shadow := NewOuterShadow("invalid")
	if err := shadow.Validate(); err == nil {
		t.Error("expected validation error for invalid color")
	}

	// Invalid transparency
	shadow = NewOuterShadow("000000").WithTransparency(2.0)
	if err := shadow.Validate(); err == nil {
		t.Error("expected validation error for transparency > 1")
	}
}

func TestFillTypeString(t *testing.T) {
	if FillTypeSolid != "solid" {
		t.Errorf("expected solid fill type, got %s", FillTypeSolid)
	}
	if FillTypePattern != "pattern" {
		t.Errorf("expected pattern fill type, got %s", FillTypePattern)
	}
	if FillTypeNoFill != "noFill" {
		t.Errorf("expected noFill fill type, got %s", FillTypeNoFill)
	}
}

func TestNormalizePatternType(t *testing.T) {
	tests := []struct {
		input    string
		expected PatternType
	}{
		{"pct5", PatternTypePct5},
		{"pct50", PatternTypePct50},
		{"horz", PatternTypeHorz},
		{"diagCross", PatternTypeDiagCross},
		{"invalid", PatternTypePct5}, // defaults to pct5
	}

	for _, test := range tests {
		result := NormalizePatternType(test.input)
		if result != test.expected {
			t.Errorf("NormalizePatternType(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestIsValidPatternType(t *testing.T) {
	if !IsValidPatternType(PatternTypePct5) {
		t.Error("expected pct5 to be valid")
	}
	if !IsValidPatternType(PatternTypeDiagCross) {
		t.Error("expected diagCross to be valid")
	}
	if IsValidPatternType(PatternType("invalid")) {
		t.Error("expected invalid pattern type to be invalid")
	}
}

func TestLineDashStyles(t *testing.T) {
	styles := []LineDashStyle{
		LineDashStyleSolid,
		LineDashStyleDash,
		LineDashStyleDot,
		LineDashStyleDashDot,
		LineDashStyleDashDotDot,
		LineDashStyleLongDash,
		LineDashStyleLongDashDot,
	}

	for _, style := range styles {
		if !IsValidLineDashStyle(style) {
			t.Errorf("expected %s to be a valid dash style", style)
		}
	}

	// Invalid style
	if IsValidLineDashStyle(LineDashStyle("invalid")) {
		t.Error("expected invalid dash style to be invalid")
	}
}

func TestShadowTypes(t *testing.T) {
	if !IsValidShadowType(ShadowTypeOuter) {
		t.Error("expected outer shadow type to be valid")
	}
	if !IsValidShadowType(ShadowTypeInner) {
		t.Error("expected inner shadow type to be valid")
	}
	if !IsValidShadowType(ShadowTypePerspective) {
		t.Error("expected perspective shadow type to be valid")
	}
	if IsValidShadowType(ShadowType("invalid")) {
		t.Error("expected invalid shadow type to be invalid")
	}
}

func TestShadowAlignments(t *testing.T) {
	alignments := []ShadowAlignment{
		ShadowAlignTopLeft,
		ShadowAlignTop,
		ShadowAlignTopRight,
		ShadowAlignLeft,
		ShadowAlignCenter,
		ShadowAlignRight,
		ShadowAlignBottomLeft,
		ShadowAlignBottom,
		ShadowAlignBottomRight,
	}

	for _, align := range alignments {
		if !IsValidShadowAlignment(align) {
			t.Errorf("expected %s to be a valid shadow alignment", align)
		}
	}
}

func TestToXMLShapeSpecWithRichFill(t *testing.T) {
	fill := NewSolidFill("4472C4").WithTransparency(0.5)
	shape := NewShape(ShapeTypeRectangle, styling.Inches(1), styling.Inches(1), styling.Inches(2), styling.Inches(2)).
		WithRichFill(fill)

	specs := ToXMLShapeSpecs([]Shape{shape}, nil)
	if len(specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(specs))
	}

	spec := specs[0]
	if spec.RichFill == nil {
		t.Error("expected rich fill spec to be set")
	}
	if spec.RichFill.Type != "solid" {
		t.Errorf("expected fill type solid, got %s", spec.RichFill.Type)
	}
	if spec.Fill != nil {
		t.Error("expected legacy fill spec to be nil")
	}
}

func TestToXMLShapeSpecWithRichLine(t *testing.T) {
	line := NewRichShapeLine("FF0000", styling.Points(2)).WithDashStyle(LineDashStyleDash)
	shape := NewShape(ShapeTypeRectangle, 0, 0, 1000000, 1000000).
		WithRichLine(line)

	specs := ToXMLShapeSpecs([]Shape{shape}, nil)
	if len(specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(specs))
	}

	spec := specs[0]
	if spec.RichLine == nil {
		t.Error("expected rich line spec to be set")
	}
	if spec.RichLine.DashStyle != "dash" {
		t.Errorf("expected dash style dash, got %s", spec.RichLine.DashStyle)
	}
}

func TestToXMLShapeSpecWithRichShadow(t *testing.T) {
	shadow := NewOuterShadow("000000").WithTransparency(0.5)
	shape := NewShape(ShapeTypeRectangle, 0, 0, 1000000, 1000000).
		WithRichShadow(shadow)

	specs := ToXMLShapeSpecs([]Shape{shape}, nil)
	if len(specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(specs))
	}

	spec := specs[0]
	if spec.RichShadow == nil {
		t.Error("expected rich shadow spec to be set")
	}
	if spec.RichShadow.Type != "outer" {
		t.Errorf("expected shadow type outer, got %s", spec.RichShadow.Type)
	}
}

func TestRichFillBackground(t *testing.T) {
	fill := NewSolidFill("FF0000").
		Background()

	if fill.Type != FillTypeNoFill {
		t.Errorf("expected fill type noFill after Background(), got %s", fill.Type)
	}
}

func TestRichFillForeground(t *testing.T) {
	fill := NewNoFill().
		Foreground()

	if fill.Type != FillTypeSolid {
		t.Errorf("expected fill type solid after Foreground(), got %s", fill.Type)
	}
}

func TestRichShapeFillFluentWithGradient(t *testing.T) {
	gradient := NewShapeGradientFill("linear", []ShapeGradientStop{
		NewShapeGradientStop(0, "FF0000"),
		NewShapeGradientStop(100, "0000FF"),
	})

	fill := NewSolidFill("FFFFFF").WithGradient(gradient)

	if fill.Type != FillTypeGradient {
		t.Errorf("expected fill type gradient, got %s", fill.Type)
	}
	if fill.Gradient == nil {
		t.Error("expected gradient to be set")
	}
	if fill.Solid != nil {
		t.Error("expected solid fill to be cleared when switching to gradient")
	}
}

func TestRichShapeFillFluentWithPattern(t *testing.T) {
	fill := NewSolidFill("FFFFFF").WithPattern(PatternTypeHorz)

	if fill.Type != FillTypePattern {
		t.Errorf("expected fill type pattern, got %s", fill.Type)
	}
	if fill.Pattern == nil {
		t.Error("expected pattern to be set")
	}
	if fill.Pattern.Pattern != PatternTypeHorz {
		t.Errorf("expected pattern horz, got %s", fill.Pattern.Pattern)
	}
}

func TestRichShapeLineFluentAPI(t *testing.T) {
	line := NewRichShapeLine("000000", styling.Points(1))

	// Test chaining
	result := line.
		WithColor("FF0000").
		WithWidth(styling.Points(2)).
		WithDashStyle(LineDashStyleDot).
		WithCapStyle(LineCapStyleSquare).
		WithJoinStyle(LineJoinStyleMiter).
		WithTransparency(0.5)

	if result.Color != "FF0000" {
		t.Errorf("expected color FF0000, got %s", result.Color)
	}
	if result.DashStyle != LineDashStyleDot {
		t.Errorf("expected dash style dot, got %s", result.DashStyle)
	}
	if result.CapStyle != LineCapStyleSquare {
		t.Errorf("expected cap style square, got %s", result.CapStyle)
	}
	if result.JoinStyle != LineJoinStyleMiter {
		t.Errorf("expected join style miter, got %s", result.JoinStyle)
	}
	if result.Transparency != 0.5 {
		t.Errorf("expected transparency 0.5, got %f", result.Transparency)
	}
}

func TestRichShapeShadowFluentAPI(t *testing.T) {
	shadow := NewRichShapeShadow().
		WithColor("FF0000").
		WithType(ShadowTypeInner).
		WithTransparency(0.7).
		WithBlurRadius(60000).
		WithDistance(40000).
		WithAngle(90).
		WithAlignment(ShadowAlignCenter).
		WithSkew(10, 20).
		WithScale(1.2, 1.3).
		WithRotateWithShape(false)

	if shadow.Color != "FF0000" {
		t.Errorf("expected color FF0000, got %s", shadow.Color)
	}
	if shadow.Type != ShadowTypeInner {
		t.Errorf("expected type inner, got %s", shadow.Type)
	}
	if shadow.Transparency != 0.7 {
		t.Errorf("expected transparency 0.7, got %f", shadow.Transparency)
	}
	if shadow.BlurRadius != 60000 {
		t.Errorf("expected blur radius 60000, got %d", shadow.BlurRadius)
	}
	if shadow.Distance != 40000 {
		t.Errorf("expected distance 40000, got %d", shadow.Distance)
	}
	if shadow.Angle != 90 {
		t.Errorf("expected angle 90, got %f", shadow.Angle)
	}
	if shadow.Alignment != ShadowAlignCenter {
		t.Errorf("expected alignment center, got %s", shadow.Alignment)
	}
	if shadow.SkewX != 10 {
		t.Errorf("expected skewX 10, got %f", shadow.SkewX)
	}
	if shadow.SkewY != 20 {
		t.Errorf("expected skewY 20, got %f", shadow.SkewY)
	}
	if shadow.ScaleX != 1.2 {
		t.Errorf("expected scaleX 1.2, got %f", shadow.ScaleX)
	}
	if shadow.ScaleY != 1.3 {
		t.Errorf("expected scaleY 1.3, got %f", shadow.ScaleY)
	}
	if shadow.RotateWithShape {
		t.Error("expected rotate with shape to be false")
	}
}

func TestNormalizeLineDashStyle(t *testing.T) {
	tests := []struct {
		input    string
		expected LineDashStyle
	}{
		{"solid", LineDashStyleSolid},
		{"dash", LineDashStyleDash},
		{"dot", LineDashStyleDot},
		{"dashDot", LineDashStyleDashDot},
		{"invalid", LineDashStyleSolid}, // defaults to solid
	}

	for _, test := range tests {
		result := NormalizeLineDashStyle(test.input)
		if result != test.expected {
			t.Errorf("NormalizeLineDashStyle(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestNormalizeLineCapStyle(t *testing.T) {
	tests := []struct {
		input    string
		expected LineCapStyle
	}{
		{"flat", LineCapStyleFlat},
		{"rnd", LineCapStyleRound},
		{"sq", LineCapStyleSquare},
		{"invalid", LineCapStyleFlat}, // defaults to flat
	}

	for _, test := range tests {
		result := NormalizeLineCapStyle(test.input)
		if result != test.expected {
			t.Errorf("NormalizeLineCapStyle(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestNormalizeLineJoinStyle(t *testing.T) {
	tests := []struct {
		input    string
		expected LineJoinStyle
	}{
		{"round", LineJoinStyleRound},
		{"bevel", LineJoinStyleBevel},
		{"miter", LineJoinStyleMiter},
		{"invalid", LineJoinStyleRound}, // defaults to round
	}

	for _, test := range tests {
		result := NormalizeLineJoinStyle(test.input)
		if result != test.expected {
			t.Errorf("NormalizeLineJoinStyle(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestIsValidLineCapStyle(t *testing.T) {
	if !IsValidLineCapStyle(LineCapStyleFlat) {
		t.Error("expected flat cap style to be valid")
	}
	if !IsValidLineCapStyle(LineCapStyleRound) {
		t.Error("expected round cap style to be valid")
	}
	if IsValidLineCapStyle(LineCapStyle("invalid")) {
		t.Error("expected invalid cap style to be invalid")
	}
}

func TestIsValidLineJoinStyle(t *testing.T) {
	if !IsValidLineJoinStyle(LineJoinStyleRound) {
		t.Error("expected round join style to be valid")
	}
	if !IsValidLineJoinStyle(LineJoinStyleBevel) {
		t.Error("expected bevel join style to be valid")
	}
	if IsValidLineJoinStyle(LineJoinStyle("invalid")) {
		t.Error("expected invalid join style to be invalid")
	}
}

func TestRichFillGetType(t *testing.T) {
	fill := NewSolidFill("FF0000")
	if fill.GetType() != FillTypeSolid {
		t.Errorf("expected get type to return solid, got %s", fill.GetType())
	}

	// Test nil fill
	var nilFill *RichShapeFill
	if nilFill.GetType() != FillTypeNoFill {
		t.Errorf("expected nil fill to return noFill, got %s", nilFill.GetType())
	}
}

func TestRichShapeShadowDefaultValues(t *testing.T) {
	shadow := NewRichShapeShadow()

	if shadow.Type != ShadowTypeOuter {
		t.Errorf("expected default type outer, got %s", shadow.Type)
	}
	if shadow.Color != "000000" {
		t.Errorf("expected default color 000000, got %s", shadow.Color)
	}
	if shadow.Transparency != 0.6 {
		t.Errorf("expected default transparency 0.6, got %f", shadow.Transparency)
	}
	if shadow.BlurRadius != 40000 {
		t.Errorf("expected default blur radius 40000, got %d", shadow.BlurRadius)
	}
	if shadow.Distance != 20000 {
		t.Errorf("expected default distance 20000, got %d", shadow.Distance)
	}
	if shadow.Angle != 45 {
		t.Errorf("expected default angle 45, got %f", shadow.Angle)
	}
	if shadow.Alignment != ShadowAlignBottomRight {
		t.Errorf("expected default alignment bottom right, got %s", shadow.Alignment)
	}
	if shadow.ScaleX != 1.0 {
		t.Errorf("expected default scaleX 1.0, got %f", shadow.ScaleX)
	}
	if shadow.ScaleY != 1.0 {
		t.Errorf("expected default scaleY 1.0, got %f", shadow.ScaleY)
	}
	if !shadow.RotateWithShape {
		t.Error("expected default rotate with shape to be true")
	}
}

func TestValidationErrorMessages(t *testing.T) {
	tests := []struct {
		name string
		fill *RichShapeFill
		want string
	}{
		{
			name: "invalid hex color",
			fill: NewSolidFill("GGGGGG"),
			want: "invalid",
		},
		{
			name: "negative transparency",
			fill: NewSolidFill("FF0000").WithTransparency(-0.1),
			want: "transparency",
		},
		{
			name: "transparency over 1",
			fill: NewSolidFill("FF0000").WithTransparency(1.5),
			want: "transparency",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fill.Validate()
			if err == nil {
				t.Error("expected validation error")
				return
			}
			if !strings.Contains(strings.ToLower(err.Error()), tt.want) {
				t.Errorf("expected error to contain %q, got %q", tt.want, err.Error())
			}
		})
	}
}

func TestShapeWithRichFormattingConflict(t *testing.T) {
	// Test that setting rich fill clears legacy fill
	legacyFill := NewShapeFill("FF0000")
	richFill := NewSolidFill("00FF00")

	shape := NewShape(ShapeTypeRectangle, 0, 0, 1000000, 1000000).
		WithFill(legacyFill).
		WithRichFill(richFill)

	if shape.RichFill == nil {
		t.Error("expected rich fill to be set")
	}
	if shape.Fill != nil {
		t.Error("expected legacy fill to be cleared when setting rich fill")
	}

	// Test that setting rich line clears legacy line
	legacyLine := NewShapeLine("000000", styling.Points(1))
	richLine := NewRichShapeLine("FFFFFF", styling.Points(2))

	shape2 := NewShape(ShapeTypeRectangle, 0, 0, 1000000, 1000000).
		WithLine(legacyLine).
		WithRichLine(richLine)

	if shape2.RichLine == nil {
		t.Error("expected rich line to be set")
	}
	if shape2.Line != nil {
		t.Error("expected legacy line to be cleared when setting rich line")
	}
}

func TestShapeValidationConflicts(t *testing.T) {
	// Test validation catches rich + legacy fill conflict
	shape := NewShape(ShapeTypeRectangle, 0, 0, 1000000, 1000000)
	shape.Fill = &ShapeFill{Color: "FF0000"}
	shape.RichFill = NewSolidFill("00FF00")

	err := shape.Validate(1, 1)
	if err == nil {
		t.Error("expected validation error for conflicting fill types")
	}
	if !strings.Contains(err.Error(), "cannot set both rich fill and legacy fill") {
		t.Errorf("expected conflict error message, got: %v", err)
	}

	// Test validation catches rich + legacy line conflict
	shape2 := NewShape(ShapeTypeRectangle, 0, 0, 1000000, 1000000)
	shape2.Line = &ShapeLine{Color: "000000", Width: styling.Points(1)}
	shape2.RichLine = NewRichShapeLine("FFFFFF", styling.Points(2))

	err2 := shape2.Validate(1, 1)
	if err2 == nil {
		t.Error("expected validation error for conflicting line types")
	}
	if !strings.Contains(err2.Error(), "cannot set both rich line and legacy line") {
		t.Errorf("expected conflict error message, got: %v", err2)
	}
}
