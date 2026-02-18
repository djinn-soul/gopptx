package main

import (
	"fmt"
	"log"
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

	// Master 2: Warm theme with Calibri
	master2 := pptx.NewMaster().
		WithBackground(pptx.SlideBackground{Type: pptx.SlideBackgroundSolid, SolidFill: &pptx.ShapeFill{Color: "FFF3E0"}}).
		WithFooter("Master 2: Warm Orange").
		WithTitleStyle([]pptx.TextLevelStyle{
			{Level: 0, Font: "Calibri", SizePt: 44, Bold: true, Color: "BF360C"},
		}).
		WithBodyStyle([]pptx.TextLevelStyle{
			{Level: 0, SizePt: 28, Color: "E65100"},
			{Level: 1, SizePt: 24, Color: "EF6C00", BulletChar: "•"},
		})

	meta := pptx.Metadata{
		Metadata: pptx.MetadataFields{Title: "Multi-Master Smoke"},
		Masters:  []*pptx.SlideMaster{master1, master2},
	}

	slides := []pptx.SlideContent{
		pptx.NewSlide("Slide 1 (Master 1)").
			AddBullet("Bullet Level 1").
			AddSubBullet(1, "Bullet Level 2"),
		pptx.NewSlide("Slide 2 (Master 2)").
			AddBullet("Second master visual family"),
		pptx.NewSlide("Slide 3 (Master 1)"),
		pptx.NewSlide("Slide 4 (Master 2)"),
	}

	data, buildErr := pptx.CreateWithMetadata(meta, slides)
	if buildErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", buildErr)
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
	log.Printf("Created %s\n", outPath)
}
