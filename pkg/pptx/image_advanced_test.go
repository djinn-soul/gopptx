package pptx_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn09/gopptx/pkg/pptx"
)

func TestImageAdvancedSources(t *testing.T) {
	// 1. Setup a simple red pixel PNG
	redPixelPNG, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKwITAAAAABJRU5ErkJggg==")

	// 2. Setup a mock HTTP server for URL testing
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		if _, err := w.Write(redPixelPNG); err != nil {
			return
		}
	}))
	defer ts.Close()

	// 3. Create images using different sources
	imgBytes := pptx.NewImageFromBytes(redPixelPNG, "png", 1000000, 1000000, 1000000, 1000000)
	imgBase64, err := pptx.NewImageFromBase64("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKwITAAAAABJRU5ErkJggg==", "png", 3000000, 1000000, 1000000, 1000000)
	if err != nil {
		t.Fatalf("failed to create image from base64: %v", err)
	}
	imgURL := pptx.NewImageFromURL(ts.URL+"/image.png", 5000000, 1000000, 1000000, 1000000)

	// 4. Test Effects
	imgEffect := pptx.NewImageFromBytes(redPixelPNG, "png", 1000000, 3000000, 1000000, 1000000).
		WithRotation(45).
		WithFlip(true, false).
		WithCrop(0.1, 0.1, 0.1, 0.1)

	slide := pptx.NewSlide("Advanced Images")
	slide.AddImage(imgBytes)
	slide.AddImage(imgBase64)
	slide.AddImage(imgURL)
	slide.AddImage(imgEffect)

	pptxFile := filepath.Join(os.TempDir(), "image_advanced_test.pptx")
	if err := pptx.WriteFile(pptxFile, "Image Test", []pptx.SlideContent{slide}); err != nil {
		t.Fatalf("failed to write pptx: %v", err)
	}
	defer func() {
		if err := os.Remove(pptxFile); err != nil && !os.IsNotExist(err) {
			t.Fatalf("failed to cleanup test pptx: %v", err)
		}
	}()

	// 5. Verify the file was created and is a valid zip (basic check)
	// In a real scenario we'd unzip and check relationships and xml content
	// using helpers similar to those in parity_fixture_test.go
}
