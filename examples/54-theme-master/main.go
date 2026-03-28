// examples/54-theme-master/main.go demonstrates theme and master customization
// using PresentationBuilder.WithTheme() with built-in theme presets.
//
// Run with: go run ./examples/54-theme-master/main.go
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
	outputFile = "54_theme_master.pptx"
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

	// Apply the ThemeTech preset (Integral-inspired blue/teal palette).
	builder := pptx.NewPresentationBuilder("Theme & Master Demo").
		WithTheme(styling.ThemeTech)

	builder.AddSlide(pptx.NewSlide("Tech Theme Applied").
		AddBullet("Using ThemeTech (Integral-style preset)").
		AddBullet("Colors: tech-focused blue/teal palette").
		AddBullet("Fonts: modern technical typeface"))

	builder.AddSlide(pptx.NewSlide("Available Theme Presets").
		AddBullet("ThemeCorporate — professional blue/grey (Office-style)").
		AddBullet("ThemeModern    — clean monochrome with cyan accents").
		AddBullet("ThemeTech      — tech-focused blue/teal (Integral-style)").
		AddBullet("ThemeDark      — dark background with bold accents (Ion-style)"))

	builder.AddSlide(pptx.NewSlide("Theme Color Slots").
		AddNumbered("Dk1 / Lt1  — primary dark and light (text / background)").
		AddNumbered("Dk2 / Lt2  — secondary dark and light").
		AddNumbered("Accent1-6  — six highlight colors used by charts and shapes").
		AddNumbered("Hlink      — hyperlink color").
		AddNumbered("FolHlink   — followed-hyperlink color"))

	builder.AddSlide(pptx.NewSlide("Theme Font Slots").
		AddNumbered("MajorFont — heading typeface").
		AddNumbered("MinorFont — body/paragraph typeface").
		AddBullet("PowerPoint resolves +mj-lt and +mn-lt against these slots"))

	builder.AddSlide(pptx.NewSlide("Custom Themes").
		AddBullet("Build a styling.Theme{} struct with your own ColorScheme and FontScheme").
		AddBullet("Pass it to PresentationBuilder.WithTheme() or Metadata.Theme").
		AddBullet("Combine WithTheme() and WithMaster() for full brand identity control"))

	outputPath := filepath.Join(outputDir, outputFile)
	if err := builder.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
