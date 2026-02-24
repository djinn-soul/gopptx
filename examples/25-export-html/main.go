package main

import (
	"log"
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
	if err := os.MkdirAll("examples/output", 0o750); err != nil {
		log.Printf("Error creating output directory: %v", err)
		return
	}

	// 3. Export HTML
	log.Println("Exporting to HTML...")
	htmlContent := export.HTML(title, slides)
	if err := os.WriteFile(outputHTML, []byte(htmlContent), 0o600); err != nil {
		log.Printf("Error writing HTML: %v", err)
		return
	}
	log.Printf("HTML exported to %s", outputHTML)

	// 4. Export PDF (if LibreOffice available)
	log.Println("Attempting PDF export (requires LibreOffice)...")
	if err := export.PDF(title, slides, outputPDF); err != nil {
		log.Printf("PDF export failed (expected if LibreOffice not installed): %v", err)
	} else {
		log.Printf("PDF exported to %s", outputPDF)
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
	_ = os.WriteFile(path, data, 0o600)
	abs, _ := filepath.Abs(path)
	log.Printf("Created dummy image at %s", abs)
}
