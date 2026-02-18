package media_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func TestImageAdvancedSources(t *testing.T) {
	// 1. Setup a simple red pixel PNG
	redPixelPNG, _ := base64.StdEncoding.DecodeString(
		"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKwITAAAAABJRU5ErkJggg==",
	)

	// 2. Setup a mock HTTP server for URL testing
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		if _, err := w.Write(redPixelPNG); err != nil {
			return
		}
	}))
	defer ts.Close()

	// 3. Create images using different sources
	imgBytes := pptx.NewImageFromBytes(redPixelPNG, "png", 1000000, 1000000, 1000000, 1000000)
	imgBase64, err := pptx.NewImageFromBase64(
		"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKwITAAAAABJRU5ErkJggg==",
		"png",
		3000000,
		1000000,
		1000000,
		1000000,
	)
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

	pptxFile := filepath.Join(t.TempDir(), "image_advanced_test.pptx")
	if writeErr := pptx.WriteFile(pptxFile, "Image Test", []pptx.SlideContent{slide}); writeErr != nil {
		t.Fatalf("failed to write pptx: %v", writeErr)
	}
}
