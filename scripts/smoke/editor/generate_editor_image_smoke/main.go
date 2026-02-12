package main

import (
	"fmt"
	"log"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func main() {
	templatePath := "smoke_samples/sampleppt/160070-labyrinth-template-16x9.pptx"
	outputPath := "smoke_samples/38_editor_image_stamping.pptx"

	fmt.Printf("Opening %s...\n", templatePath)
	e, err := pptx.OpenPresentationEditor(templatePath)
	if err != nil {
		log.Fatalf("Failed to open template: %v", err)
	}

	// Create a dummy image (small 1x1 white PNG)
	// Actually, I'll just use a real image if available or generate one.
	// For testing, I'll use the logo if I can find one, or just a placeholder.
	// But wait, I can just use a red square data.
	redSquare := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, 0x54, 0x08, 0xD7, 0x63, 0xF8, 0xFF, 0xFF, 0x3F,
		0x00, 0x05, 0xFE, 0x02, 0xFE, 0xDC, 0x44, 0x74, 0x06, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	// 1. Stamping the same image on multiple slides to test deduplication
	logo := shapes.NewImageFromBytes(redSquare, "png", styling.Inches(0.5), styling.Inches(0.5), styling.Inches(1), styling.Inches(1))

	fmt.Println("Stamping images on multiple slides...")
	for i := 0; i < 3; i++ {
		slide := elements.NewSlide(fmt.Sprintf("Image Stamp Test %d", i+1)).
			AddBullet("This slide has a stamped logo.").
			AddBullet("Deduplication should ensure only one image file is added to the package.").
			AddImage(logo)

		_, err = e.AddSlide(slide)
		if err != nil {
			log.Fatalf("Failed to add slide %d: %v", i+1, err)
		}
	}

	// 2. Picture Background
	fmt.Println("Adding slide with picture background...")
	bgSlide := elements.NewSlide("Picture Background Test").
		WithPictureBackground(logo).
		AddBullet("This slide uses the same logo as a background.").
		AddBullet("Still should be deduplicated.")

	_, err = e.AddSlide(bgSlide)
	if err != nil {
		log.Fatalf("Failed to add background slide: %v", err)
	}

	fmt.Printf("Saving to %s...\n", outputPath)
	if err := pptx.Save(e, outputPath); err != nil {
		log.Fatalf("Failed to save: %v", err)
	}

	fmt.Println("Done!")
}
