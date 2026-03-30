// examples/66-placeholders-api demonstrates working with slide placeholders.
//
// Shows how to override title/body placeholders with text, images, and tables
// using WithPlaceholderText, WithPlaceholderImage, and WithPlaceholderTable.
// Also demonstrates WithPlaceholderTextAs with explicit placeholder types.
//
// Run with: go run ./examples/66-placeholders-api/main.go
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
	outputFile = "66_placeholders_api.pptx"
)

// whitePNGBytes returns a minimal 1x1 white PNG for placeholder image demo.
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

	// --- Slide 1: WithPlaceholderText (default types by index) ---
	// Index 0 => title placeholder, index 1 => body placeholder.
	slide1 := pptx.NewSlide("Placeholder Text Override").
		WithPlaceholderText(0, "Title Set via Placeholder Index 0").
		WithPlaceholderText(1, "Body content set via placeholder index 1.\nSecond line of body text.")

	// --- Slide 2: WithPlaceholderTextAs with explicit type names ---
	slide2 := pptx.NewSlide("Explicit Placeholder Types").
		WithPlaceholderTextAs(0, string(pptx.PlaceholderTypeTitle), "Explicit Title Placeholder").
		WithPlaceholderTextAs(1, string(pptx.PlaceholderTypeBody), "Explicit body placeholder text")

	// --- Slide 3: Placeholder with image ---
	img := pptx.NewImageFromBytes(
		whitePNGBytes(), "png",
		styling.Inches(0), styling.Inches(0),
		styling.Inches(2), styling.Inches(2),
	)
	slide3 := pptx.NewSlide("Image Placeholder Demo").
		WithPlaceholderText(0, "Slide with Image Placeholder").
		WithPlaceholderImage(1, img)

	// --- Slide 4: Placeholder with table ---
	table := pptx.NewTable([]styling.Length{
		styling.Inches(2.5),
		styling.Inches(2.5),
		styling.Inches(2),
	}).
		AddRow([]string{"Column A", "Column B", "Column C"}).
		AddRow([]string{"Value 1", "Value 2", "Value 3"}).
		AddRow([]string{"Value 4", "Value 5", "Value 6"})

	slide4 := pptx.NewSlide("Table Placeholder Demo").
		WithPlaceholderText(0, "Slide with Table in Placeholder").
		WithPlaceholderTable(1, table)

	// --- Slide 5: Available placeholder type constants ---
	slide5 := pptx.NewSlide("Placeholder Type Constants").
		AddBullet(fmt.Sprintf("PlaceholderTypeTitle   = %q", pptx.PlaceholderTypeTitle)).
		AddBullet(fmt.Sprintf("PlaceholderTypeBody    = %q", pptx.PlaceholderTypeBody)).
		AddBullet(fmt.Sprintf("PlaceholderTypeSubTitle= %q", pptx.PlaceholderTypeSubTitle)).
		AddBullet(fmt.Sprintf("PlaceholderTypeChart   = %q", pptx.PlaceholderTypeChart)).
		AddBullet(fmt.Sprintf("PlaceholderTypeTbl     = %q", pptx.PlaceholderTypeTbl)).
		AddBullet(fmt.Sprintf("PlaceholderTypePic     = %q", pptx.PlaceholderTypePic)).
		AddBullet(fmt.Sprintf("PlaceholderTypeObj     = %q", pptx.PlaceholderTypeObj))

	slides := []pptx.SlideContent{slide1, slide2, slide3, slide4, slide5}

	outputPath := filepath.Join(outputDir, outputFile)
	data, err := pptx.CreateWithSlides("Placeholders API Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
