package export

import (
	"math"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestHasImageCrop(t *testing.T) {
	if hasImageCrop(shapes.ImageCrop{}) {
		t.Fatal("empty srcRect reported as a crop")
	}
	if !hasImageCrop(shapes.ImageCrop{Left: 0.1}) {
		t.Fatal("left crop not detected")
	}
	// Negative insets pad rather than crop, but still need the crop path.
	if !hasImageCrop(shapes.ImageCrop{Bottom: -0.2}) {
		t.Fatal("negative inset not detected")
	}
}

func TestCroppedImagePlacementFillsTheBox(t *testing.T) {
	// Keep the middle half of the picture horizontally and the bottom three
	// quarters vertically.
	crop := shapes.ImageCrop{Left: 0.25, Right: 0.25, Top: 0.25}
	got, ok := croppedImagePlacement(crop, 100, 200, 50, 30)
	if !ok {
		t.Fatal("placement rejected a valid crop")
	}

	// Visible width is half the source, so the whole picture is drawn twice the
	// box width; the kept region then starts a quarter of that further left.
	wantW, wantH := 100.0, 40.0
	wantX, wantY := 100-0.25*wantW, 200-0.25*wantH
	for _, c := range []struct {
		name      string
		got, want float64
	}{
		{"x", got.X, wantX},
		{"y", got.Y, wantY},
		{"w", got.W, wantW},
		{"h", got.H, wantH},
	} {
		if math.Abs(c.got-c.want) > 1e-9 {
			t.Fatalf("%s=%v want %v", c.name, c.got, c.want)
		}
	}

	// The kept region must line up with the destination box on every edge.
	if left := got.X + crop.Left*got.W; math.Abs(left-100) > 1e-9 {
		t.Fatalf("visible left edge=%v want 100", left)
	}
	if right := got.X + (1-crop.Right)*got.W; math.Abs(right-150) > 1e-9 {
		t.Fatalf("visible right edge=%v want 150", right)
	}
	if top := got.Y + crop.Top*got.H; math.Abs(top-200) > 1e-9 {
		t.Fatalf("visible top edge=%v want 200", top)
	}
}

func TestCroppedImagePlacementRejectsEmptyCrop(t *testing.T) {
	if _, ok := croppedImagePlacement(shapes.ImageCrop{Left: 0.5, Right: 0.5}, 0, 0, 10, 10); ok {
		t.Fatal("a crop that keeps nothing was accepted")
	}
	if _, ok := croppedImagePlacement(shapes.ImageCrop{Top: 1.2}, 0, 0, 10, 10); ok {
		t.Fatal("an over-large crop was accepted")
	}
}
