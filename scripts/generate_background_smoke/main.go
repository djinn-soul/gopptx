package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	builder := pptx.NewPresentationBuilder("Slide Background Showcase")

	// Slide 1: Solid Background
	slide1 := pptx.NewSlide("Solid Background (Legacy compatible)").
		WithBackgroundColor("FF9900"). // Should still work
		AddShape(pptx.NewTextBox("This uses .WithBackgroundColor(\"FF9900\")", 1, 2, 8, 1))
	builder.AddSlide(slide1)

	// Slide 2: Complex Solid Background
	bgSolid := pptx.NewSolidBackground("00AAFF")
	slide2 := pptx.NewSlide("Solid Background (New API)").
		WithBackground(bgSolid).
		AddShape(pptx.NewTextBox("This uses .WithBackground(pptx.NewSolidBackground(\"00AAFF\"))", 1, 2, 8, 1))
	builder.AddSlide(slide2)

	// Slide 3: Linear Gradient Background
	gradLinear := pptx.NewShapeGradientFill(pptx.ShapeGradientTypeLinear, []pptx.ShapeGradientStop{
		pptx.NewShapeGradientStop(0, "FFFFFF"),
		pptx.NewShapeGradientStop(100, "CCCCCC"),
	})
	slide3 := pptx.NewSlide("Linear Gradient Background").
		WithGradientBackground(gradLinear).
		AddShape(pptx.NewTextBox("Linear Gradient (Top to Bottom)", 1, 2, 8, 1))
	builder.AddSlide(slide3)

	// Slide 4: Radial Gradient Background
	gradRadial := pptx.NewShapeGradientFill(pptx.ShapeGradientTypeRadial, []pptx.ShapeGradientStop{
		pptx.NewShapeGradientStop(0, "FFEE00"),
		pptx.NewShapeGradientStop(100, "FF0000"),
	})
	slide4 := pptx.NewSlide("Radial Gradient Background").
		WithGradientBackground(gradRadial).
		AddShape(pptx.NewTextBox("Radial Gradient (Center out)", 1, 2, 8, 1))
	builder.AddSlide(slide4)

	// Slide 5: Picture Background
	imgPath := "smoke_samples/sampleimage/repository-open-graph-template.png"
	if _, err := os.Stat(imgPath); err == nil {
		img := pptx.NewImage(imgPath, 0, 0, 0, 0) // Dimensions ignored for background (stretched by default)
		slide5 := pptx.NewSlide("Picture Background").
			WithPictureBackground(img).
			AddShape(pptx.NewTextBox("This slide has a picture background", 1, 2, 8, 1))
		builder.AddSlide(slide5)
	} else {
		log.Printf("Warning: background image sample not found at %s, skipping slide 5", imgPath)
	}

	// Save
	outDir := "smoke_samples"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
	outPath := filepath.Join(outDir, "gopptx_slide_background_smoke.pptx")
	if err := builder.WriteToFile(outPath); err != nil {
		log.Fatalf("Failed to save presentation: %v", err)
	}

	log.Printf("Successfully generated smoke sample: %s", outPath)
}
