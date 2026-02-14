package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	// Master 1: Blue theme with Arial
	master1 := pptx.NewMaster().
		WithBackground(pptx.SlideBackground{Type: pptx.SlideBackgroundSolid, SolidFill: &pptx.ShapeFill{Color: "E3F2FD"}}).
		WithFooter("Master 1: Professional Blue").
		WithTitleStyle([]pptx.TextLevelStyle{
			{Level: 0, Font: "Arial", SizePt: 44, Bold: true, Color: "0D47A1"},
		}).
		WithBodyStyle([]pptx.TextLevelStyle{
			{Level: 0, SizePt: 28, Color: "1A237E"},
			{Level: 1, SizePt: 24, Color: "3949AB", BulletChar: "•"},
		})

	meta := pptx.PresentationMetadata{
		Masters: []*pptx.SlideMaster{master1},
	}
	meta.Title = "Master Styling Smoke Test"

	slides := []pptx.SlideContent{
		pptx.NewSlide("Master 1 applied").
			AddBullet("Bullet Level 1").
			AddSubBullet(1, "Bullet Level 2"),
		pptx.NewSlide("Second slide"),
	}

	data, err := pptx.CreateWithMetadata(meta, slides)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir error: %v\n", err)
		os.Exit(1)
	}
	outPath := filepath.Join(outputDir, "36_multi_master_smoke.pptx")
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created %s\n", outPath)
}
