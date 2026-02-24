package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func main() {
	// 1. Create a dummy image
	imgPath := "dummy_image.png"
	createDummyImage(imgPath)
	defer os.Remove(imgPath)

	// 2. Create Presentation structure
	// Note: We use elements.SlideContent directly as per export API
	slides := []elements.SlideContent{
		elements.NewSlide("HTML Export Demo").
			AddBullet("This presentation was exported to HTML via gopptx.").
			AddBullet("It mimics the ppt-rs export style."),
		elements.NewSlide("Visual Content").
			AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, 1000000, 1000000, 3000000, 1000000).
				WithText("This is a shape text.")).
			AddImage(shapes.Image{Path: imgPath, CX: 1000000, CY: 1000000}),
		elements.NewSlide("Code Simulation").
			AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, 1000000, 2000000, 3000000, 2000000).
				WithText("func main() {\n    fmt.Println(\"Hello World\")\n}")),
	}

	title := "Export Demo"
	outputHTML := "examples/output/25_export.html"
	outputPDF := "examples/output/25_export.pdf"

	// Ensure output dir exists
	if err := os.MkdirAll("examples/output", 0o755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// 3. Export HTML
	fmt.Println("Exporting to HTML...")
	htmlContent := export.HTML(title, slides)
	if err := os.WriteFile(outputHTML, []byte(htmlContent), 0o644); err != nil {
		fmt.Printf("Error writing HTML: %v\n", err)
		return
	}
	fmt.Printf("HTML exported to %s\n", outputHTML)

	// 4. Export PDF (if LibreOffice available)
	fmt.Println("Attempting PDF export (requires LibreOffice)...")
	if err := export.PDF(title, slides, outputPDF); err != nil {
		fmt.Printf("PDF export failed (expected if LibreOffice not installed): %v\n", err)
	} else {
		fmt.Printf("PDF exported to %s\n", outputPDF)
	}
}

func createDummyImage(path string) {
	// 1x1 red pixel PNG
	data := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0xF8, 0xCF, 0x00, 0x00,
		0x02, 0x03, 0x01, 0x01, 0x24, 0x95, 0x8C, 0xFC,
		0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44,
		0xAE, 0x42, 0x60, 0x82,
	}
	_ = os.WriteFile(path, data, 0o644)
	abs, _ := filepath.Abs(path)
	fmt.Printf("Created dummy image at %s\n", abs)
}
