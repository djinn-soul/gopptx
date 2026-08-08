package shapes_test

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// pngBytes encodes a solid image of the given pixel size.
func pngBytes(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	return buf.Bytes()
}

func TestNaturalSizeReadsTheImage(t *testing.T) {
	img := shapes.Image{Data: pngBytes(t, 192, 96), Format: "png"}
	cx, cy, err := img.NaturalSize()
	if err != nil {
		t.Fatalf("NaturalSize: %v", err)
	}
	// 96 pixels per inch.
	if cx != styling.Inches(2) || cy != styling.Inches(1) {
		t.Fatalf("natural size = (%d, %d), want (%d, %d)", cx, cy, styling.Inches(2), styling.Inches(1))
	}
}

func TestNaturalSizeReadsFromDisk(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pic.png")
	if err := os.WriteFile(path, pngBytes(t, 96, 96), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	img, err := shapes.NewImageAtNaturalSize(path, styling.Inches(1), styling.Inches(1))
	if err != nil {
		t.Fatalf("NewImageAtNaturalSize: %v", err)
	}
	if img.CX != styling.Inches(1) || img.CY != styling.Inches(1) {
		t.Fatalf("size = (%d, %d), want one inch square", img.CX, img.CY)
	}
}

func TestScaleToWidthKeepsAspectRatio(t *testing.T) {
	img := shapes.Image{Data: pngBytes(t, 400, 100), Format: "png"}
	scaled, err := img.ScaleToWidth(styling.Inches(4))
	if err != nil {
		t.Fatalf("ScaleToWidth: %v", err)
	}
	if scaled.CY != styling.Inches(1) {
		t.Fatalf("height = %d, want %d", scaled.CY, styling.Inches(1))
	}

	tall, err := img.ScaleToHeight(styling.Inches(1))
	if err != nil {
		t.Fatalf("ScaleToHeight: %v", err)
	}
	if tall.CX != styling.Inches(4) {
		t.Fatalf("width = %d, want %d", tall.CX, styling.Inches(4))
	}
}

func TestFitWithinUsesTheConstrainingAxis(t *testing.T) {
	wide := shapes.Image{Data: pngBytes(t, 400, 100), Format: "png"}
	fitted, err := wide.FitWithin(styling.Inches(4), styling.Inches(0.5))
	if err != nil {
		t.Fatalf("FitWithin: %v", err)
	}
	// Height is the tighter bound, so the width comes down with it.
	if fitted.CY != styling.Inches(0.5) || fitted.CX != styling.Inches(2) {
		t.Fatalf("fitted = (%d, %d), want (%d, %d)",
			fitted.CX, fitted.CY, styling.Inches(2), styling.Inches(0.5))
	}
}

func TestSizeOfAnUnfetchedURLIsUnknown(t *testing.T) {
	img := shapes.Image{SourceURL: "https://example.com/pic.png"}
	if _, _, err := img.NaturalSize(); !errors.Is(err, shapes.ErrUnknownImageSize) {
		t.Fatalf("error = %v, want ErrUnknownImageSize", err)
	}
}

func TestAspectRatioPrefersTheStatedExtent(t *testing.T) {
	img := shapes.Image{
		Data: pngBytes(t, 400, 100), Format: "png",
		CX: styling.Inches(2), CY: styling.Inches(2),
	}
	ratio, err := img.AspectRatio()
	if err != nil {
		t.Fatalf("AspectRatio: %v", err)
	}
	if ratio != 1 {
		t.Fatalf("ratio = %v, want 1 from the stated extent", ratio)
	}
}
