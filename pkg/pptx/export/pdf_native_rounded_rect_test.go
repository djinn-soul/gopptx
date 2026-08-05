package export

import (
	"math"
	"testing"
)

// gopdf degrades a rounded rectangle to a plain polygon when it is given no
// corner points, which is what the renderer used to pass, so every roundRect
// came out square.
func TestRoundRectCornerPointsAreDrawn(t *testing.T) {
	if roundRectCornerPoints <= 0 {
		t.Fatalf("roundRectCornerPoints=%d; gopdf draws a square below 1", roundRectCornerPoints)
	}
}

// PowerPoint's roundRect preset rounds by its default adj, 16667/100000 of the
// shorter side.
func TestRoundRectRadiusMatchesPowerPoint(t *testing.T) {
	const shorter = 60.0
	got := roundRectRadius(120, shorter)
	want := shorter * 0.16667
	if math.Abs(got-want) > 0.01 {
		t.Errorf("radius=%v want %v", got, want)
	}
}

// gopdf errors out rather than clamping when the radius exceeds a side, so the
// radius must never grow past half the shorter side.
func TestRoundRectRadiusNeverExceedsHalfTheShorterSide(t *testing.T) {
	tests := []struct{ w, h float64 }{
		{100, 1},
		{1, 100},
		{4, 4},
		{0.5, 80},
	}
	for _, tt := range tests {
		shorter := math.Min(tt.w, tt.h)
		if got := roundRectRadius(tt.w, tt.h); got > shorter/2 {
			t.Errorf("%vx%v gave radius %v, over half the shorter side %v", tt.w, tt.h, got, shorter)
		}
	}
}

func TestRoundRectRadiusOfADegenerateBoxIsZero(t *testing.T) {
	if got := roundRectRadius(0, 10); got != 0 {
		t.Errorf("radius=%v want 0 for a zero-width box", got)
	}
	if got := roundRectRadius(10, -5); got != 0 {
		t.Errorf("radius=%v want 0 for a negative-height box", got)
	}
}
