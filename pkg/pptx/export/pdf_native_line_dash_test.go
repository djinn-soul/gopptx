package export

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestPDFDashPatternSolidAndUnknown(t *testing.T) {
	for _, dash := range []string{"", "  ", shapes.LineDashSolid, "notADash"} {
		if got := pdfDashPattern(dash, 2); got != nil {
			t.Fatalf("pdfDashPattern(%q)=%v want nil (solid)", dash, got)
		}
	}
}

func TestPDFDashPatternScalesByStrokeWidth(t *testing.T) {
	got := pdfDashPattern(shapes.LineDashDash, 2)
	want := []float64{8, 6} // 4 and 3 times the 2pt stroke
	if len(got) != len(want) {
		t.Fatalf("pattern=%v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("pattern=%v want %v", got, want)
		}
	}
}

func TestPDFDashPatternFloorsHairlineUnit(t *testing.T) {
	// A 0.1pt rule would dash at 0.1pt marks, which prints as a solid line;
	// the unit floors at minDashUnitPt instead.
	got := pdfDashPattern(shapes.LineDashDot, 0.1)
	if len(got) != 2 || got[0] != minDashUnitPt || got[1] != 3*minDashUnitPt {
		t.Fatalf("hairline pattern=%v want [%v %v]", got, minDashUnitPt, 3*minDashUnitPt)
	}
}

func TestPDFDashPatternCoversEveryPreset(t *testing.T) {
	presets := []string{
		shapes.LineDashDot,
		shapes.LineDashDash,
		shapes.LineDashDashDot,
		shapes.LineDashDashDotDot,
		shapes.LineDashLongDash,
		shapes.LineDashLongDashDot,
		string(shapes.LineDashStyleDashDotDot),
		string(shapes.LineDashStyleSystemDash),
		string(shapes.LineDashStyleSystemDot),
		string(shapes.LineDashStyleSystemDashDot),
	}
	for _, dash := range presets {
		pattern := pdfDashPattern(dash, 1)
		if len(pattern) < 2 || len(pattern)%2 != 0 {
			t.Fatalf("preset %q pattern=%v want an even, non-empty run of marks and gaps", dash, pattern)
		}
		for _, v := range pattern {
			if v <= 0 {
				t.Fatalf("preset %q pattern=%v has a non-positive length", dash, pattern)
			}
		}
	}
}

func TestPDFDashPatternIsCaseInsensitive(t *testing.T) {
	lower := pdfDashPattern("lgdashdot", 1)
	canonical := pdfDashPattern(shapes.LineDashLongDashDot, 1)
	if len(lower) != len(canonical) {
		t.Fatalf("case-insensitive lookup failed: %v vs %v", lower, canonical)
	}
}

func TestShapeLineDashPrefersRichLine(t *testing.T) {
	s := shapes.NewShape("rect", styling.Inches(0), styling.Inches(0), styling.Inches(1), styling.Inches(1))
	if got := shapeLineDash(s); got != "" {
		t.Fatalf("dash of a shape with no line=%q want empty", got)
	}

	line := shapes.NewShapeLine("FF0000", styling.Points(1))
	line.Dash = shapes.LineDashDash
	s.Line = &line
	if got := shapeLineDash(s); got != shapes.LineDashDash {
		t.Fatalf("simple line dash=%q want %q", got, shapes.LineDashDash)
	}

	// The rich line is the fuller reading of the same a:ln, so it wins.
	s.RichLine = &shapes.RichShapeLine{DashStyle: shapes.LineDashStyleSystemDot}
	if got := shapeLineDash(s); got != string(shapes.LineDashStyleSystemDot) {
		t.Fatalf("rich line dash=%q want %q", got, shapes.LineDashStyleSystemDot)
	}

	// A rich line that states no dash falls back to the simple one.
	s.RichLine = &shapes.RichShapeLine{}
	if got := shapeLineDash(s); got != shapes.LineDashDash {
		t.Fatalf("fallback dash=%q want %q", got, shapes.LineDashDash)
	}
}
