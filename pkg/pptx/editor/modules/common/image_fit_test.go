package common

import "testing"

func TestFitImageToBoxStretchKeepsTheBox(t *testing.T) {
	got, err := FitImageToBox(ImageFitStretch, 100, 200, 4000, 3000, 800, 200)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.X != 100 || got.Y != 200 || got.CX != 4000 || got.CY != 3000 {
		t.Fatalf("stretch changed the box: %+v", got)
	}
	if got.Crop != nil {
		t.Fatalf("stretch must not crop: %+v", got.Crop)
	}
}

func TestFitImageToBoxEmptyModeIsStretch(t *testing.T) {
	got, err := FitImageToBox("", 0, 0, 4000, 3000, 800, 200)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.CX != 4000 || got.CY != 3000 {
		t.Fatalf("empty fit mode must behave like stretch: %+v", got)
	}
}

func TestFitImageToBoxContainLetterboxesAWideImage(t *testing.T) {
	// 4:1 image in a 1:1 box: keep the width, shrink the height to a quarter,
	// and center the result vertically.
	got, err := FitImageToBox(ImageFitContain, 0, 0, 4000, 4000, 800, 200)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.CX != 4000 || got.CY != 1000 {
		t.Fatalf("expected 4000x1000, got %dx%d", got.CX, got.CY)
	}
	if got.X != 0 || got.Y != 1500 {
		t.Fatalf("expected offset (0,1500), got (%d,%d)", got.X, got.Y)
	}
	if got.Crop != nil {
		t.Fatalf("contain must not crop: %+v", got.Crop)
	}
}

func TestFitImageToBoxContainPillarboxesATallImage(t *testing.T) {
	// 1:4 image in a 1:1 box: keep the height, shrink the width.
	got, err := FitImageToBox(ImageFitContain, 500, 500, 4000, 4000, 200, 800)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.CX != 1000 || got.CY != 4000 {
		t.Fatalf("expected 1000x4000, got %dx%d", got.CX, got.CY)
	}
	if got.X != 2000 || got.Y != 500 {
		t.Fatalf("expected offset (2000,500), got (%d,%d)", got.X, got.Y)
	}
}

func TestFitImageToBoxContainKeepsAMatchingAspectRatio(t *testing.T) {
	got, err := FitImageToBox(ImageFitContain, 0, 0, 4000, 2000, 800, 400)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.CX != 4000 || got.CY != 2000 || got.X != 0 || got.Y != 0 {
		t.Fatalf("matching aspect ratio should fill the box exactly: %+v", got)
	}
}

func TestFitImageToBoxCoverCropsAWideImageHorizontally(t *testing.T) {
	// 4:1 image in a 1:1 box: only a quarter of the width is visible, so 75%
	// is cropped, split evenly as 37.5% per side.
	got, err := FitImageToBox(ImageFitCover, 0, 0, 4000, 4000, 800, 200)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.CX != 4000 || got.CY != 4000 || got.X != 0 || got.Y != 0 {
		t.Fatalf("cover must fill the box: %+v", got)
	}
	if got.Crop == nil {
		t.Fatal("cover must crop a wide image")
	}
	if got.Crop.Left != 37500 || got.Crop.Right != 37500 {
		t.Fatalf("expected 37500 left/right, got %+v", got.Crop)
	}
	if got.Crop.Top != 0 || got.Crop.Bottom != 0 {
		t.Fatalf("a wide image must not be cropped vertically: %+v", got.Crop)
	}
}

func TestFitImageToBoxCoverCropsATallImageVertically(t *testing.T) {
	got, err := FitImageToBox(ImageFitCover, 0, 0, 4000, 4000, 200, 800)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Crop == nil {
		t.Fatal("cover must crop a tall image")
	}
	if got.Crop.Top != 37500 || got.Crop.Bottom != 37500 {
		t.Fatalf("expected 37500 top/bottom, got %+v", got.Crop)
	}
	if got.Crop.Left != 0 || got.Crop.Right != 0 {
		t.Fatalf("a tall image must not be cropped horizontally: %+v", got.Crop)
	}
}

func TestFitImageToBoxCoverSkipsCropWhenAspectRatiosMatch(t *testing.T) {
	got, err := FitImageToBox(ImageFitCover, 0, 0, 4000, 2000, 800, 400)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Crop != nil {
		t.Fatalf("a matching aspect ratio needs no crop: %+v", got.Crop)
	}
}

func TestFitImageToBoxRejectsUnknownMode(t *testing.T) {
	if _, err := FitImageToBox("squish", 0, 0, 100, 100, 10, 10); err == nil {
		t.Fatal("expected an error for an unknown fit mode")
	}
}

func TestFitImageToBoxRejectsMissingDimensions(t *testing.T) {
	if _, err := FitImageToBox(ImageFitContain, 0, 0, 100, 100, 0, 10); err == nil {
		t.Fatal("expected an error when the image width is unknown")
	}
	if _, err := FitImageToBox(ImageFitCover, 0, 0, 0, 100, 10, 10); err == nil {
		t.Fatal("expected an error when the box has no width")
	}
}

func TestNormalizeImageFit(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"", ImageFitStretch},
		{ImageFitStretch, ImageFitStretch},
		{ImageFitContain, ImageFitContain},
		{ImageFitCover, ImageFitCover},
	} {
		got, err := NormalizeImageFit(tc.in)
		if err != nil {
			t.Fatalf("NormalizeImageFit(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("NormalizeImageFit(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
	if _, err := NormalizeImageFit("COVER"); err == nil {
		t.Fatal("fit modes are case sensitive; expected an error")
	}
}
