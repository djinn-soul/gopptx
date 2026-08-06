package export

import (
	"math"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestArrowGeometryScalesWithTheShorterSide(t *testing.T) {
	// A connector-style arrow: 158pt wide, 10.8pt tall. PowerPoint sizes the
	// head off the shorter side, so it is ~5.4pt — not half the width.
	geom := arrowGeometryFor(nil, 158, 10.8)
	if math.Abs(geom.head-5.4) > 0.01 {
		t.Fatalf("head=%v want 5.4 (half the 10.8pt height)", geom.head)
	}
	if math.Abs(geom.shaft-5.4) > 0.01 {
		t.Fatalf("shaft=%v want 5.4 (half the height)", geom.shaft)
	}
}

func TestArrowGeometryHonoursAdjustments(t *testing.T) {
	s := shapes.NewShape("rightArrow", 0, 0, 0, 0).
		WithAdjustmentValue("adj1", 25000).
		WithAdjustmentValue("adj2", 100000)

	geom := arrowGeometryFor(s.Adjustments, 100, 40)
	if math.Abs(geom.head-40) > 0.01 {
		t.Fatalf("head=%v want 40 (adj2 of the 40pt shorter side)", geom.head)
	}
	if math.Abs(geom.shaft-10) > 0.01 {
		t.Fatalf("shaft=%v want 10 (adj1 of the 40pt height)", geom.shaft)
	}
}

func TestArrowGeometryClampsToTheBox(t *testing.T) {
	// A head longer than the shape cannot stick out of it.
	wide := shapes.NewShape("rightArrow", 0, 0, 0, 0).WithAdjustmentValue("adj2", 100000)
	geom := arrowGeometryFor(wide.Adjustments, 10, 400)
	if geom.head > 10 {
		t.Fatalf("head=%v want at most the 10pt width", geom.head)
	}
}

func TestOOXMLGuideValueFallsBackWhenUnset(t *testing.T) {
	if got := ooxmlGuideValue(nil, "adj1", 12345); got != 12345 {
		t.Fatalf("missing guide=%v want the fallback", got)
	}
	junk := []shapes.ShapeAdjustment{{Name: "adj1", Formula: "*/ 1 2 3"}}
	if got := ooxmlGuideValue(junk, "adj1", 500); got != 500 {
		t.Fatalf("non-val formula=%v want the fallback", got)
	}
}

func TestRightArrowPointsPutTheHeadAtTheTip(t *testing.T) {
	geom := arrowGeometryFor(nil, 158, 10.8)
	pts := rightArrowPoints(0, 0, 158, 10.8, geom)

	// The tip is the rightmost point, on the vertical centre line.
	tip := pts[3]
	if math.Abs(tip.X-158) > 0.01 || math.Abs(tip.Y-5.4) > 0.01 {
		t.Fatalf("tip=%+v want the mid-right of the box", tip)
	}
	// The shaft runs from the left edge to where the head starts.
	if math.Abs(pts[1].X-(158-geom.head)) > 0.01 {
		t.Fatalf("shaft ends at %v want %v", pts[1].X, 158-geom.head)
	}
}
