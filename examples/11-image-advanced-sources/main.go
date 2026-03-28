// examples/11-image-advanced-sources/main.go demonstrates advanced image source options.
//
// Shows how to embed images using:
//   - pptx.NewImageFromBytes  (raw []byte)
//   - pptx.NewImageFromBase64 (base64-encoded string)
//
// Run with: go run ./examples/11-image-advanced-sources/main.go
package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "11_image_advanced_sources.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Minimal 1x1 white PNG.
	whitePNG := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xFF, 0xFF, 0x3F,
		0x00, 0x05, 0xFE, 0x02, 0xFE, 0xDC, 0x44, 0x74,
		0x06, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	// --- Slide 1: image from raw bytes ---
	imgFromBytes := pptx.NewImageFromBytes(
		whitePNG,
		"png",
		styling.Inches(1), styling.Inches(1.5), // x, y
		styling.Inches(3), styling.Inches(2), // width, height
	)

	slide1 := pptx.NewSlide("Image from Bytes").
		AddBullet("pptx.NewImageFromBytes(data, format, x, y, cx, cy)").
		AddBullet("Pass any []byte PNG/JPEG/GIF directly — no file needed.").
		AddImage(imgFromBytes)

	// --- Slide 2: image from base64 string ---
	b64 := base64.StdEncoding.EncodeToString(whitePNG)

	imgFromBase64, err := pptx.NewImageFromBase64(
		b64,
		"png",
		styling.Inches(1), styling.Inches(1.5), // x, y
		styling.Inches(3), styling.Inches(2), // width, height
	)
	if err != nil {
		return fmt.Errorf("base64 image: %w", err)
	}

	slide2 := pptx.NewSlide("Image from Base64").
		AddBullet("pptx.NewImageFromBase64(b64string, format, x, y, cx, cy)").
		AddBullet("Useful when images are stored as base64 strings (e.g. JSON APIs, config files).").
		AddImage(imgFromBase64)

	slides := []pptx.SlideContent{slide1, slide2}

	outputPath := filepath.Join(outputDir, outputFile)
	if err := pptx.WriteFile(outputPath, "Advanced Image Sources Demo", slides); err != nil {
		return fmt.Errorf("write presentation: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
