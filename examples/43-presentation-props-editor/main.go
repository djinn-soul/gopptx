package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	outputDir  = "examples/output"
	outputFile = "43_presentation_props_editor.pptx"
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

	tmpDir, tempErr := os.MkdirTemp("", "gopptx-props-example-*")
	if tempErr != nil {
		return fmt.Errorf("create temp directory: %w", tempErr)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to remove temp dir %s: %v\n", tmpDir, err)
		}
	}()

	inputPath := filepath.Join(tmpDir, "base_props_input.pptx")
	outputPath := filepath.Join(outputDir, outputFile)

	slides := []pptx.SlideContent{
		pptx.NewSlide("Presentation Properties Demo").
			AddBullet("This deck is edited with PresentationEditor."),
	}
	if err := pptx.WriteFile(inputPath, "Presentation Props Base", slides); err != nil {
		return fmt.Errorf("create base presentation: %w", err)
	}

	editor, err := pptx.OpenPresentationEditor(inputPath)
	if err != nil {
		return fmt.Errorf("open base presentation: %w", err)
	}
	defer func() {
		if closeErr := editor.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close editor: %v\n", closeErr)
		}
	}()

	if err := editor.ApplyTheme(styling.ThemeCorporate); err != nil {
		return fmt.Errorf("apply theme: %w", err)
	}
	if err := editor.SetSlideSize(pptx.SlideSize16x9); err != nil {
		return fmt.Errorf("set slide size: %w", err)
	}

	props := common.CoreProperties{
		Title:          "Presentation Properties Example",
		Subject:        "Editor metadata update",
		Creator:        "gopptx example",
		Description:    "Demonstrates theme, slide size, and core properties edits.",
		Keywords:       "gopptx, editor, metadata",
		LastModifiedBy: "gopptx example",
	}
	editor.SetCoreProperties(props)

	if err := editor.Save(outputPath); err != nil {
		return fmt.Errorf("save edited presentation: %w", err)
	}

	verificationEditor, err := pptx.OpenPresentationEditor(outputPath)
	if err != nil {
		return fmt.Errorf("reopen saved presentation: %w", err)
	}
	defer func() {
		if closeErr := verificationEditor.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close verification editor: %v\n", closeErr)
		}
	}()

	got := verificationEditor.GetCoreProperties()
	if got.Title != props.Title {
		return fmt.Errorf("core title mismatch: got %q want %q", got.Title, props.Title)
	}

	log.Printf("Generated presentation properties example: %s\n", outputPath)
	return nil
}
