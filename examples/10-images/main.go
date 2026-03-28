// examples/10-images/main.go demonstrates image embedding from a local file and from bytes.
//
// Shows how to embed images using pptx.NewImageFromBytes (no file required at runtime)
// and pptx.NewImage (from a temp file on disk).
//
// Run with: go run ./examples/10-images/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "10_images.pptx"
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

	// --- Slide 1: image embedded from raw bytes ---
	imgFromBytes := pptx.NewImageFromBytes(
		whitePNG,
		"png",
		styling.Inches(1), styling.Inches(1), // x, y
		styling.Inches(4), styling.Inches(3), // width, height
	)

	slide1 := pptx.NewSlide("Image from Bytes").
		AddBullet("This slide embeds a 1x1 white PNG supplied as a raw []byte.").
		AddBullet("No file is required at runtime.").
		AddImage(imgFromBytes)

	// --- Slide 2: image loaded from a temp file on disk ---
	tmpDir, err := os.MkdirTemp("", "gopptx-images-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpPNG := filepath.Join(tmpDir, "white.png")
	if err := os.WriteFile(tmpPNG, whitePNG, 0o600); err != nil {
		return fmt.Errorf("write temp PNG: %w", err)
	}

	imgFromFile := pptx.NewImage(
		tmpPNG,
		styling.Inches(1), styling.Inches(1), // x, y
		styling.Inches(4), styling.Inches(3), // width, height
	)

	slide2 := pptx.NewSlide("Image from File Path").
		AddBullet("This slide loads the same PNG from a temporary file on disk.").
		AddBullet("Use pptx.NewImage(path, x, y, cx, cy) for file-based images.").
		AddImage(imgFromFile)

	slides := []pptx.SlideContent{slide1, slide2}

	outputPath := filepath.Join(outputDir, outputFile)
	if err := pptx.WriteFile(outputPath, "Image Embedding Demo", slides); err != nil {
		return fmt.Errorf("write presentation: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
