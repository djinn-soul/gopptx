package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	outputDir  = "examples/output"
	outputFile = "43_advanced_theme_management.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// Step 1: build a base presentation with several content slides.
	tmpDir, err := os.MkdirTemp("", "gopptx-theme-mgmt-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	tmpPath := filepath.Join(tmpDir, "base.pptx")
	if err := buildBase(tmpPath); err != nil {
		return fmt.Errorf("build base: %w", err)
	}

	// Step 2: open the saved file with the editor.
	e, err := pptx.OpenPresentationEditor(tmpPath)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	defer func() {
		if closeErr := e.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: close editor: %v\n", closeErr)
		}
	}()

	// Step 3: try a built-in theme preset ("integral").
	// SetGlobalThemePreset rewrites the theme XML embedded in the PPTX.
	// It may return an error when the base presentation has no theme part yet;
	// log it as a warning rather than aborting.
	if err := e.SetGlobalThemePreset("integral"); err != nil {
		log.Printf("Note: SetGlobalThemePreset: %v", err)
	}

	// Step 4: apply a direct styling.Theme (ThemeTech).
	// ApplyTheme rewrites theme colors and fonts from a Go-native Theme struct.
	// Like SetGlobalThemePreset, it requires a theme part; log errors as warnings.
	if err := e.ApplyTheme(styling.ThemeTech); err != nil {
		log.Printf("Note: ApplyTheme(ThemeTech): %v", err)
	}

	// Step 5: add a reference slide listing every available preset and theme var.
	_, err = e.AddSlide(buildReferenceSlide())
	if err != nil {
		return fmt.Errorf("add reference slide: %w", err)
	}

	// Step 6: save.
	outputPath := filepath.Join(outputDir, outputFile)
	if err := e.Save(outputPath); err != nil {
		return fmt.Errorf("save: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

// buildBase creates a presentation with a few slides for theme-application demos.
func buildBase(path string) error {
	builder := pptx.NewPresentationBuilder("Advanced Theme Management")

	builder.AddSlide(
		pptx.NewSlide("Theme Management Overview").
			AddBullet("SetGlobalThemePreset – apply an OOXML built-in preset").
			AddBullet("ApplyTheme – apply a Go-native styling.Theme struct").
			AddBullet("Both methods rewrite the embedded theme XML"),
	)

	builder.AddSlide(
		pptx.NewSlide("Available Presets").
			AddNumbered("office2013").
			AddNumbered("facet").
			AddNumbered("integral").
			AddNumbered("ion").
			AddNumbered("retrospect").
			AddNumbered("slice").
			AddNumbered("wisp"),
	)

	builder.AddSlide(
		pptx.NewSlide("Available styling.Theme Variables").
			AddNumbered("styling.ThemeCorporate – professional and trustworthy").
			AddNumbered("styling.ThemeModern   – clean and simple").
			AddNumbered("styling.ThemeTech     – modern technology feel").
			AddNumbered("styling.ThemeDark     – easy on the eyes").
			AddNumbered("styling.ThemeVibrant  – bold and colorful").
			AddNumbered("styling.ThemeNature   – fresh and organic").
			AddNumbered("styling.ThemeCarbon   – IBM Carbon design"),
	)

	return builder.WriteToFile(path)
}

// buildReferenceSlide summarises what was applied so it is visible in the deck.
func buildReferenceSlide() pptx.SlideContent {
	return pptx.NewSlide("Applied in This Run").
		AddBullet("SetGlobalThemePreset(\"integral\") – attempted").
		AddBullet("ApplyTheme(styling.ThemeTech)   – attempted").
		AddBullet("Errors logged as warnings if no theme part present")
}
