package editorcommon

import "testing"

// A group whose child space is twice its drawn size halves every child
// coordinate on the way to the slide. Reporting the raw child numbers is the
// wrong position upstream issue #925 describes.
func TestChildToSlideScalesAndTranslates(t *testing.T) {
	space := &GroupChildSpace{OffsetX: 1000, OffsetY: 2000, ExtentCx: 4000, ExtentCy: 8000}

	x, y, w, h := space.ChildToSlide(
		5000, 6000, 2000, 4000,
		2000, 4000, 1000, 2000,
	)

	if x != 5500 || y != 7000 {
		t.Fatalf("expected child origin (5500,7000), got (%d,%d)", x, y)
	}
	if w != 500 || h != 1000 {
		t.Fatalf("expected child extent (500,1000), got (%d,%d)", w, h)
	}
}

// Most groups state a child space identical to their own box, and then the
// mapping has to be an exact identity or every child moves.
func TestChildToSlideIsIdentityForMatchingSpaces(t *testing.T) {
	space := &GroupChildSpace{OffsetX: 100, OffsetY: 200, ExtentCx: 900, ExtentCy: 800}

	x, y, w, h := space.ChildToSlide(100, 200, 900, 800, 400, 500, 300, 100)

	if x != 400 || y != 500 || w != 300 || h != 100 {
		t.Fatalf("expected the child unchanged, got (%d,%d,%d,%d)", x, y, w, h)
	}
}

// A shape that is not a group has no child space, and must pass through.
func TestChildToSlideNilSpacePassesThrough(t *testing.T) {
	var space *GroupChildSpace

	x, y, w, h := space.ChildToSlide(10, 20, 30, 40, 1, 2, 3, 4)

	if x != 1 || y != 2 || w != 3 || h != 4 {
		t.Fatalf("expected the child unchanged, got (%d,%d,%d,%d)", x, y, w, h)
	}
}

// chExt="0" carries no scale; PowerPoint treats the spaces as the same size.
func TestChildToSlideZeroChildExtentOnlyTranslates(t *testing.T) {
	space := &GroupChildSpace{OffsetX: 50, OffsetY: 50}

	x, y, w, h := space.ChildToSlide(1000, 1000, 500, 500, 150, 250, 100, 200)

	if x != 1100 || y != 1200 || w != 100 || h != 200 {
		t.Fatalf("expected translation only, got (%d,%d,%d,%d)", x, y, w, h)
	}
}
