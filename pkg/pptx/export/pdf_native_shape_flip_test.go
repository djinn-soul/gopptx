package export

import (
	"math"
	"testing"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

func flipShape(flipH, flipV bool) shapes.Shape {
	s := shapes.NewShape(shapes.ShapeTypeRightArrow, 0, 0, 0, 0)
	s.FlipH = flipH
	s.FlipV = flipV
	return s
}

func TestFlipLeavesAnUnflippedShapeAlone(t *testing.T) {
	fl := flipFor(flipShape(false, false), 10, 20, 100, 50)
	points := []gopdf.Point{{X: 10, Y: 20}, {X: 110, Y: 70}}

	got := fl.points(points)
	if &got[0] != &points[0] {
		t.Error("an unflipped shape copied its points instead of passing them through")
	}
}

func TestFlipMirrorsAboutTheShapeCentre(t *testing.T) {
	const x, y, w, h = 10.0, 20.0, 100.0, 50.0
	// The tip of a right arrow, on the shape's right edge at mid height.
	tip := gopdf.Point{X: x + w, Y: y + h/2}

	horizontal := flipFor(flipShape(true, false), x, y, w, h).points([]gopdf.Point{tip})[0]
	if math.Abs(horizontal.X-x) > 0.001 {
		t.Errorf("flipH put the tip at x=%.3f, want the left edge %v", horizontal.X, x)
	}
	if math.Abs(horizontal.Y-tip.Y) > 0.001 {
		t.Errorf("flipH moved the tip vertically to %.3f, want %v", horizontal.Y, tip.Y)
	}

	corner := gopdf.Point{X: x, Y: y}
	vertical := flipFor(flipShape(false, true), x, y, w, h).points([]gopdf.Point{corner})[0]
	if math.Abs(vertical.Y-(y+h)) > 0.001 {
		t.Errorf("flipV put the corner at y=%.3f, want the bottom edge %v", vertical.Y, y+h)
	}
	if math.Abs(vertical.X-corner.X) > 0.001 {
		t.Errorf("flipV moved the corner horizontally to %.3f, want %v", vertical.X, corner.X)
	}

	both := flipFor(flipShape(true, true), x, y, w, h).points([]gopdf.Point{corner})[0]
	if math.Abs(both.X-(x+w)) > 0.001 || math.Abs(both.Y-(y+h)) > 0.001 {
		t.Errorf("flipping both put the corner at %.3f,%.3f, want %v,%v", both.X, both.Y, x+w, y+h)
	}
}

func TestFlippingTwiceReturnsTheOriginal(t *testing.T) {
	const x, y, w, h = 0.0, 0.0, 80.0, 40.0
	fl := flipFor(flipShape(true, true), x, y, w, h)
	points := chevronPoints(x, y, w, h)

	round := fl.points(fl.points(points))
	for i := range points {
		if math.Abs(round[i].X-points[i].X) > 0.001 || math.Abs(round[i].Y-points[i].Y) > 0.001 {
			t.Fatalf("point %d came back as %.3f,%.3f, want %.3f,%.3f",
				i, round[i].X, round[i].Y, points[i].X, points[i].Y)
		}
	}
}

func TestFlipKeepsTheShapeInsideItsBox(t *testing.T) {
	const x, y, w, h = 5.0, 7.0, 60.0, 30.0
	fl := flipFor(flipShape(true, true), x, y, w, h)
	for _, p := range fl.points(rightArrowPoints(x, y, w, h, arrowGeometryFor(nil, w, h))) {
		if p.X < x-0.001 || p.X > x+w+0.001 || p.Y < y-0.001 || p.Y > y+h+0.001 {
			t.Errorf("flipped point %.3f,%.3f left the %vx%v box at %v,%v", p.X, p.Y, w, h, x, y)
		}
	}
}

func TestCachedShapeCarriesItsFlip(t *testing.T) {
	cached := cachedNodeShape("Step")
	cached.FlipH = true
	cached.FlipV = true

	shape, ok := smartArtCachedShape(cachedDiagram([]smartart.DrawingShape{cached}, "Step"), cached)
	if !ok {
		t.Fatal("the cached shape was rejected")
	}
	if !shape.FlipH || !shape.FlipV {
		t.Errorf("flip was dropped: FlipH=%v FlipV=%v", shape.FlipH, shape.FlipV)
	}
}
