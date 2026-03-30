// examples/77-background-api demonstrates slide background fills.
//
// Shows solid color backgrounds, gradient backgrounds, and picture (image)
// backgrounds using NewSolidBackground, NewGradientBackground, and
// NewPictureBackground.
//
// Run with: go run ./examples/77-background-api/main.go
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
	outputFile = "77_background_api.pptx"

	colorWhite = "FFFFFF"
)

// bluePNGBytes returns a minimal 1x1 blue PNG for picture background demo.
func bluePNGBytes() []byte {
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xD7, 0x63, 0x60, 0x98, 0xFF, 0xFF,
		0x00, 0x00, 0x02, 0x00, 0x01, 0xE5, 0x27, 0xDE,
		0xFC, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
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

	slides := []pptx.SlideContent{
		buildDeepBlueSolidSlide(),
		buildLightSolidSlide(),
		buildLinearGradientSlide(),
		buildVerticalGradientSlide(),
		buildPictureBackgroundSlide(),
		buildDefaultBackgroundSlide(),
		buildBackgroundConstantsSlide(),
	}

	outputPath := filepath.Join(outputDir, outputFile)
	data, err := pptx.CreateWithSlides("Background API Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

func buildDeepBlueSolidSlide() pptx.SlideContent {
	bgSolid := pptx.NewSolidBackground("1565C0") // Deep blue
	slide1 := pptx.NewSlide("Solid Color Background")
	slide1.Background = &bgSolid
	slide1.TitleColor = colorWhite
	slide1.ContentColor = colorWhite
	return slide1.AddBullet("Background: NewSolidBackground(\"1565C0\")").
		AddBullet("Any 6-digit hex color is accepted.")
}

func buildLightSolidSlide() pptx.SlideContent {
	bgLight := pptx.NewSolidBackground("F0F4FF")
	slide2 := pptx.NewSlide("Light Solid Background")
	slide2.Background = &bgLight
	return slide2.AddBullet("Background: NewSolidBackground(\"F0F4FF\")").
		AddBullet("Light backgrounds work well with dark content.")
}

func buildLinearGradientSlide() pptx.SlideContent {
	gradStops := []pptx.ShapeGradientStop{
		pptx.NewShapeGradientStop(0, "4472C4"),
		pptx.NewShapeGradientStop(100, colorWhite),
	}
	gradFill := pptx.NewShapeGradientFill("linear", gradStops).WithLinearAngle(135)
	bgGrad := pptx.NewGradientBackground(gradFill)

	slide3 := pptx.NewSlide("Gradient Background")
	slide3.Background = &bgGrad
	return slide3.AddBullet("Background: NewGradientBackground(gradientFill)").
		AddBullet("Linear gradient from blue (top-left) to white (bottom-right).")
}

func buildVerticalGradientSlide() pptx.SlideContent {
	stops2 := []pptx.ShapeGradientStop{
		pptx.NewShapeGradientStop(0, "1B5E20"),
		pptx.NewShapeGradientStop(100, "E8F5E9"),
	}
	gradFill2 := pptx.NewShapeGradientFill("linear", stops2).WithLinearAngle(90)
	bgGrad2 := pptx.NewGradientBackground(gradFill2)

	slide4 := pptx.NewSlide("Vertical Gradient Background")
	slide4.Background = &bgGrad2
	return slide4.AddBullet("Linear gradient: dark green to light green, angle 90°.")
}

func buildPictureBackgroundSlide() pptx.SlideContent {
	img := pptx.NewImageFromBytes(
		bluePNGBytes(), "png",
		styling.Inches(0), styling.Inches(0),
		styling.Inches(10), styling.Inches(7.5),
	)
	bgPic := pptx.NewPictureBackground(img)

	slide5 := pptx.NewSlide("Picture Background")
	slide5.Background = &bgPic
	slide5.TitleColor = colorWhite
	return slide5.AddBullet("Background: NewPictureBackground(image)").
		AddBullet("Accepts a shapes.Image – from bytes, file, or URL.")
}

func buildDefaultBackgroundSlide() pptx.SlideContent {
	return pptx.NewSlide("Default Theme Background").
		AddBullet("No explicit background – uses the presentation theme default.").
		AddBullet("slide.Background = nil (or not set at all)")
}

func buildBackgroundConstantsSlide() pptx.SlideContent {
	return pptx.NewSlide("Background Type Constants").
		AddBullet(fmt.Sprintf("SlideBackgroundSolid    = %q", pptx.SlideBackgroundSolid)).
		AddBullet(fmt.Sprintf("SlideBackgroundGradient = %q", pptx.SlideBackgroundGradient)).
		AddBullet(fmt.Sprintf("SlideBackgroundPicture  = %q", pptx.SlideBackgroundPicture))
}
