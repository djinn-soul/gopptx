package export

import (
	"math"
	"testing"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// newlyDrawnPresets are the presets these two tables added. Each used to fall
// through to a plain rectangle.
func newlyDrawnPresets() []string {
	return []string{
		shapes.ShapeTypeRound1Rect,
		shapes.ShapeTypeRound2SameRect,
		shapes.ShapeTypeRound2DiagRect,
		shapes.ShapeTypeSnip1Rect,
		shapes.ShapeTypeSnip2SameRect,
		shapes.ShapeTypeSnip2DiagRect,
		shapes.ShapeTypeSnipRoundRect,
		shapes.ShapeTypeCube,
		shapes.ShapeTypeFoldedCorner,
		shapes.ShapeTypeCorner,
		shapes.ShapeTypeDiagStripe,
		shapes.ShapeTypeHalfFrame,
		shapes.ShapeTypePlaque,
		shapes.ShapeTypeChevron,
		shapes.ShapeTypeNotchedRightArrow,
		shapes.ShapeTypeTeardrop,
		shapes.ShapeTypeMoon,
		shapes.ShapeTypeRibbon,
		shapes.ShapeTypeRibbon2,
	}
}

func TestNewPresetsAreRecognised(t *testing.T) {
	pdf := newTestPDF(t)
	fl := flipState{unflippedShape: true}

	for _, preset := range newlyDrawnPresets() {
		corner := drawPDFCornerGeometry(pdf, fl, preset, 0, 0, 100, 60, "D")
		solid := drawPDFSolidGeometry(pdf, fl, preset, 0, 0, 100, 60, "D")
		if !corner && !solid {
			t.Errorf("%s still falls through to the rectangle fallback", preset)
		}
		if corner && solid {
			t.Errorf("%s is claimed by both geometry tables", preset)
		}
	}

	if drawPDFCornerGeometry(pdf, fl, "notAPreset", 0, 0, 10, 10, "D") {
		t.Error("the corner table claimed an unknown preset")
	}
	if drawPDFSolidGeometry(pdf, fl, "notAPreset", 0, 0, 10, 10, "D") {
		t.Error("the solid table claimed an unknown preset")
	}
}

func TestNewPresetOutlinesStayInsideTheirBox(t *testing.T) {
	const x, y, w, h = 12.0, 8.0, 140.0, 90.0
	cases := map[string][]gopdf.Point{
		"cube":          cubePoints(x, y, w, h),
		"folded corner": foldedCornerPoints(x, y, w, h),
		"corner":        cornerPoints(x, y, w, h),
		"diag stripe":   diagStripePoints(x, y, w, h),
		"half frame":    halfFramePoints(x, y, w, h),
		"notched arrow": notchedRightArrowPoints(x, y, w, h),
		"teardrop":      teardropPoints(x, y, w, h),
		"ribbon":        ribbonPoints(x, y, w, h, false),
		"ribbon2":       ribbonPoints(x, y, w, h, true),
	}
	for name, points := range cases {
		assertPointsInsideBox(t, name, points, x, y, w, h)
	}

	// Every rounded/snipped corner combination stays in the box too.
	for _, preset := range []string{
		shapes.ShapeTypeRound1Rect, shapes.ShapeTypeRound2SameRect,
		shapes.ShapeTypeRound2DiagRect, shapes.ShapeTypeSnip1Rect,
		shapes.ShapeTypeSnip2SameRect, shapes.ShapeTypeSnip2DiagRect,
		shapes.ShapeTypeSnipRoundRect,
	} {
		corners, ok := cornersForPreset(preset)
		if !ok {
			t.Fatalf("%s has no corner treatments", preset)
		}
		assertPointsInsideBox(t, preset, roundedRectPoints(x, y, w, h, corners), x, y, w, h)
	}
}

func TestRibbonTailsFoldOppositeWays(t *testing.T) {
	const x, y, w, h = 0.0, 0.0, 100.0, 60.0
	down := ribbonPoints(x, y, w, h, false)
	up := ribbonPoints(x, y, w, h, true)

	lowest := func(points []gopdf.Point) float64 {
		out := points[0].Y
		for _, p := range points {
			out = math.Max(out, p.Y)
		}
		return out
	}
	if math.Abs(lowest(down)-(y+h)) > 0.001 {
		t.Errorf("ribbon's tails reach y=%.3f, want the bottom edge %v", lowest(down), y+h)
	}
	if lowest(up) >= y+h {
		t.Errorf("ribbon2's tails reach y=%.3f, want them above the bottom edge", lowest(up))
	}
}

func TestSnippedAndRoundedCornersDiffer(t *testing.T) {
	const x, y, w, h = 0.0, 0.0, 100.0, 100.0
	snipped, _ := cornersForPreset(shapes.ShapeTypeSnip1Rect)
	rounded, _ := cornersForPreset(shapes.ShapeTypeRound1Rect)

	snipPoints := roundedRectPoints(x, y, w, h, snipped)
	roundPoints := roundedRectPoints(x, y, w, h, rounded)
	if len(snipPoints) == len(roundPoints) {
		t.Errorf("a snipped corner traced %d points and a rounded one %d — they should differ",
			len(snipPoints), len(roundPoints))
	}
}

func TestArcPointsWalkTheRequestedSweep(t *testing.T) {
	centre := gopdf.Point{X: 10, Y: 20}
	points := arcPoints(centre, 5, 0, 90, 4)
	if len(points) != 5 {
		t.Fatalf("got %d points, want segments+1", len(points))
	}
	if math.Abs(points[0].X-15) > 0.001 || math.Abs(points[0].Y-20) > 0.001 {
		t.Errorf("arc starts at %.3f,%.3f, want 15,20", points[0].X, points[0].Y)
	}
	if math.Abs(points[4].X-10) > 0.001 || math.Abs(points[4].Y-25) > 0.001 {
		t.Errorf("arc ends at %.3f,%.3f, want 10,25", points[4].X, points[4].Y)
	}
}
