package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func main() {
	// Define a custom Neon Theme
	neonTheme := styling.Theme{
		Name: "NeonStream",
		Colors: styling.ColorScheme{
			Name:     "NeonStream Colors",
			Dk1:      "000000", // Background
			Lt1:      "FFFFFF", // Text
			Dk2:      "1A1A1A", // Darker accent
			Lt2:      "E0E0E0", // Lighter accent
			Accent1:  "00FFFF", // Cyan
			Accent2:  "FF00FF", // Magenta
			Accent3:  "FFFF00", // Yellow
			Accent4:  "00FF00", // Green
			Accent5:  "FF8000", // Orange
			Accent6:  "4B0082", // Indigo
			Hlink:    "00BFFF", // Deep Sky Blue
			FolHlink: "9932CC", // Dark Orchid
		},
		Fonts: styling.FontScheme{
			Name:      "Modern Tech",
			MajorFont: "Inter",
			MinorFont: "Roboto",
		},
	}

	// Define a custom Slide Master
	// Use the correct API for gradients
	neonGradient := pptx.NewShapeGradientFill(pptx.ShapeGradientTypeLinear, []pptx.ShapeGradientStop{
		pptx.NewShapeGradientStop(0, "000000"),
		pptx.NewShapeGradientStop(100, "1A1A1A"),
	})

	master := elements.NewMaster().
		WithBackground(elements.NewGradientBackground(neonGradient)).
		WithFooter("© 2026 NeonStream Technology - Confidential").
		WithColorMapping("dk1", "lt1") // Map slide bg to theme dk1 (black) and text to theme lt1 (white)

	// Add decorative master shapes
	// Top neon bar
	master.AddShape(pptx.NewRectangle(0, 0, 13.33, 0.05).
		WithFill(pptx.NewShapeFill("00FFFF")))

	// Bottom neon bar
	master.AddShape(pptx.NewRectangle(0, 7.45, 13.33, 0.05).
		WithFill(pptx.NewShapeFill("FF00FF")))

	// Placeholder "Logo" shape in corner
	master.AddShape(pptx.NewStar(12.5, 0.2, 0.6).
		WithFill(pptx.NewShapeFill("00FFFF")).
		WithLine(pptx.NewShapeLine("FFFFFF", 10000)))

	// Create presentation with custom theme and master
	builder := pptx.NewPresentationBuilder("NeonStream Corporate Identity").
		WithTheme(neonTheme).
		WithMaster(master).
		WithSlideSize(pptx.SlideSize16x9)

	// Add slides
	builder.AddTitleSlide("NeonStream: The Future of Streaming")

	builder.AddBulletSlide("Key Competitive Advantages", []string{
		"Hyper-converged edge networking",
		"Zero-latency transcoding pipelines",
		"AI-driven audience engagement",
		"Quantum-secure content distribution",
	})

	builder.AddShapesSlide("Visual Design Language",
		pptx.NewRectangle(1, 2, 3, 1).WithFill(pptx.NewShapeFill("00FFFF")).WithText("Primary: Cyan"),
		pptx.NewRectangle(5, 2, 3, 1).WithFill(pptx.NewShapeFill("FF00FF")).WithText("Secondary: Magenta"),
		pptx.NewRectangle(9, 2, 3, 1).WithFill(pptx.NewShapeFill("FFFF00")).WithText("Tertiary: Yellow"),
	)

	// Build and save
	data, err := builder.Build()
	if err != nil {
		log.Fatalf("failed to build presentation: %v", err)
	}

	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("failed to create output dir: %v", err)
	}

	filename := filepath.Join(outputDir, "54_theme_master.pptx")
	if err := os.WriteFile(filename, data, 0o644); err != nil {
		log.Fatalf("failed to write file: %v", err)
	}

	log.Printf("Successfully generated smoke sample: %s\n", filename)
}
