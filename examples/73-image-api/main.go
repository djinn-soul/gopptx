// examples/73-image-api demonstrates image embedding in slides.
//
// Shows images from raw bytes, from a file path, from a base64 string,
// and image options: rotation, crop, flip (horizontal/vertical),
// shadow, reflection, and alt text.
//
// Run with: go run ./examples/73-image-api/main.go
package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	outputDir  = "examples/output"
	outputFile = "73_image_api.pptx"
)

// whitePNGBytes returns a minimal 1x1 white PNG.
func whitePNGBytes() []byte {
	return []byte{
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
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	slide2, err := buildImageFromFileSlide()
	if err != nil {
		return err
	}
	slide3, err := buildImageFromBase64Slide()
	if err != nil {
		return err
	}

	slides := []pptx.SlideContent{
		buildImageFromBytesSlide(),
		slide2,
		slide3,
		buildRotationFlipSlide(),
		buildCropSlide(),
		buildShadowReflectionSlide(),
		buildDecorativeSlide(),
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err := pptx.WriteFile(outputPath, "Image API Demo", slides); err != nil {
		return fmt.Errorf("write presentation: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

func buildImageFromBytesSlide() pptx.SlideContent {
	imgFromBytes := pptx.NewImageFromBytes(
		whitePNGBytes(), "png",
		styling.Inches(1), styling.Inches(1.5),
		styling.Inches(4), styling.Inches(3),
	).WithAltText("A white 1×1 PNG embedded from raw bytes")

	return pptx.NewSlide("Image from Bytes").
		AddBullet("pptx.NewImageFromBytes(data, format, x, y, cx, cy)").
		AddBullet("No file on disk required at runtime.").
		AddImage(imgFromBytes)
}

func buildImageFromFileSlide() (pptx.SlideContent, error) {
	tmpDir, err := os.MkdirTemp("", "gopptx-img-*")
	if err != nil {
		return pptx.SlideContent{}, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpPNG := filepath.Join(tmpDir, "sample.png")
	if err := os.WriteFile(tmpPNG, whitePNGBytes(), 0o600); err != nil {
		return pptx.SlideContent{}, fmt.Errorf("write temp PNG: %w", err)
	}

	imgFromFile := pptx.NewImage(
		tmpPNG,
		styling.Inches(1), styling.Inches(1.5),
		styling.Inches(4), styling.Inches(3),
	).WithAltText("Same white PNG loaded from a temp file")

	slide := pptx.NewSlide("Image from File Path").
		AddBullet("pptx.NewImage(path, x, y, cx, cy)").
		AddBullet("Reads the image from disk at creation time.").
		AddImage(imgFromFile)
	return slide, nil
}

func buildImageFromBase64Slide() (pptx.SlideContent, error) {
	b64Str := base64.StdEncoding.EncodeToString(whitePNGBytes())
	imgFromB64, err := pptx.NewImageFromBase64(
		b64Str, "png",
		styling.Inches(1), styling.Inches(1.5),
		styling.Inches(4), styling.Inches(3),
	)
	if err != nil {
		return pptx.SlideContent{}, fmt.Errorf("decode base64 image: %w", err)
	}

	slide := pptx.NewSlide("Image from Base64").
		AddBullet("pptx.NewImageFromBase64(b64, format, x, y, cx, cy)").
		AddBullet("Useful when images arrive as base64-encoded strings.").
		AddImage(imgFromB64)
	return slide, nil
}

func buildRotationFlipSlide() pptx.SlideContent {
	imgRotated := pptx.NewImageFromBytes(
		whitePNGBytes(), "png",
		styling.Inches(1), styling.Inches(1.5),
		styling.Inches(3), styling.Inches(2),
	).WithRotation(30)

	imgFlipped := pptx.NewImageFromBytes(
		whitePNGBytes(), "png",
		styling.Inches(5), styling.Inches(1.5),
		styling.Inches(3), styling.Inches(2),
	).WithFlip(true, false)

	return pptx.NewSlide("Image Rotation & Flip").
		AddBullet("WithRotation(degrees) – rotate the image").
		AddBullet("WithFlip(horizontal, vertical) – mirror the image").
		AddImage(imgRotated).
		AddImage(imgFlipped)
}

func buildCropSlide() pptx.SlideContent {
	imgCropped := pptx.NewImageFromBytes(
		whitePNGBytes(), "png",
		styling.Inches(1), styling.Inches(1.5),
		styling.Inches(6), styling.Inches(3),
	).WithCrop(0.1, 0.1, 0.1, 0.1)

	return pptx.NewSlide("Image Crop").
		AddBullet("WithCrop(left, right, top, bottom) – values are 0.0–1.0 fractions").
		AddBullet("Crops 10% from each side in this example.").
		AddImage(imgCropped)
}

func buildShadowReflectionSlide() pptx.SlideContent {
	imgShadow := pptx.NewImageFromBytes(
		whitePNGBytes(), "png",
		styling.Inches(0.5), styling.Inches(1.5),
		styling.Inches(3.5), styling.Inches(2.5),
	).WithShadow(true)

	imgReflection := pptx.NewImageFromBytes(
		whitePNGBytes(), "png",
		styling.Inches(5), styling.Inches(1.5),
		styling.Inches(3.5), styling.Inches(2.5),
	).WithReflection(true)

	return pptx.NewSlide("Shadow & Reflection").
		AddBullet("WithShadow(true) – adds an outer shadow effect").
		AddBullet("WithReflection(true) – adds a reflection effect").
		AddImage(imgShadow).
		AddImage(imgReflection)
}

func buildDecorativeSlide() pptx.SlideContent {
	imgDecorative := pptx.NewImageFromBytes(
		whitePNGBytes(), "png",
		styling.Inches(1), styling.Inches(2),
		styling.Inches(6), styling.Inches(3),
	).WithDecorative(true)

	return pptx.NewSlide("Decorative Image (Accessibility)").
		AddBullet("WithDecorative(true) marks the image as decorative.").
		AddBullet("Screen readers skip decorative images.").
		AddImage(imgDecorative)
}
