package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	outputDir   = "smoke_samples"
	overwriteFN = "57_editor_overwrite_existing.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	targetPath := filepath.Join(outputDir, overwriteFN)

	// 1) Create an existing PPTX file.
	baseSlides := []pptx.SlideContent{
		pptx.NewSlide("Original Title").AddBullet("This file will be edited in-place"),
		pptx.NewSlide("Original Content").AddBullet("Second slide"),
	}
	if err := pptx.WriteFile(targetPath, "Editor Overwrite Demo", baseSlides); err != nil {
		return fmt.Errorf("create base deck: %w", err)
	}
	fmt.Printf("1. Created existing deck: %s\n", targetPath)

	// 2) Open existing file with editor.
	editor, err := pptx.OpenPresentationEditor(targetPath)
	if err != nil {
		return fmt.Errorf("open existing deck: %w", err)
	}
	fmt.Println("2. Opened existing deck with PresentationEditor")

	// 3) Edit slides + presentation properties.
	theme := styling.ThemeTech

	updatedSlide := pptx.NewSlide("Edited Title").
		WithBackgroundColor(theme.Colors.Lt2).
		AddBullet("Updated in place").
		AddBullet("Visual theme colors applied")
	updatedSlide = updatedSlide.AddShape(
		pptx.NewRectangle(0.7, 1.7, 5.8, 0.65).
			WithFill(pptx.NewShapeFill(theme.Colors.Accent1)).
			WithText("Accent 1 bar"),
	)
	updatedSlide = updatedSlide.AddShape(
		pptx.NewRectangle(0.7, 2.55, 5.8, 0.65).
			WithFill(pptx.NewShapeFill(theme.Colors.Accent2)).
			WithText("Accent 2 bar"),
	)

	if err := editor.UpdateSlide(0, updatedSlide); err != nil {
		return fmt.Errorf("update slide 1: %w", err)
	}

	addedSlide := pptx.NewSlide("Added Slide").
		WithBackgroundColor("FFFFFF").
		AddBullet("Inserted before overwrite save").
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, 914400, 1828800, 5486400, 914400).
				WithFill(pptx.NewShapeFill(theme.Colors.Accent3)).
				WithText("Accent 3 highlight"),
		)
	if _, err := editor.AddSlide(addedSlide); err != nil {
		return fmt.Errorf("add slide: %w", err)
	}
	if err := editor.ApplyTheme(theme); err != nil {
		return fmt.Errorf("apply theme: %w", err)
	}
	if err := editor.SetSlideSize(pptx.SlideSize16x9); err != nil {
		return fmt.Errorf("set slide size: %w", err)
	}
	fmt.Println("3. Edited slides, theme, and slide size")

	// 4) Save back to the same path (overwrite existing PPTX).
	if err := editor.Save(targetPath); err != nil {
		return fmt.Errorf("overwrite save: %w", err)
	}
	fmt.Printf("4. Overwrote existing file successfully: %s\n", targetPath)

	// 5) Reopen to verify overwrite result is readable.
	edited, err := pptx.OpenPresentationEditor(targetPath)
	if err != nil {
		return fmt.Errorf("reopen overwritten deck: %w", err)
	}
	fmt.Printf("5. Verified overwritten file. Slide count: %d\n", edited.SlideCount())

	return nil
}
