package main

import (
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

func main() {
	// 1. Load the pre-built vbaProject.bin blob.
	// In a real scenario, this would be produced by Office or a compatible tool.
	vbaData, err := os.ReadFile("examples/assets/vbaProject.bin")
	if err != nil {
		log.Printf("Warning: Failed to load vbaProject.bin assets: %v. Using dummy data.", err)
		vbaData = []byte("dummy vba data")
	}

	// 2. Define the VBA project.
	project := vba.FromData(vbaData).
		AddModule(vba.NewModule("Module1", "Sub Hello()\nMsgBox \"Hello from gopptx!\"\nEnd Sub"))

	// 3. Create a presentation with the VBA project.
	meta := pptx.Metadata{
		Metadata: common.Metadata{
			Title: "VBA Macro Example",
		},
		VBA: project,
	}

	slides := []pptx.SlideContent{
		pptx.NewSlide("VBA Macro Example").AddBullet("This presentation contains a VBA macro."),
	}

	data, err := pptx.CreateWithMetadata(meta, slides)
	if err != nil {
		log.Printf("Failed to create presentation: %v", err)
		os.Exit(1)
	}

	// 4. Save as .pptm (macro-enabled).
	outDir := "examples/output"
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		log.Printf("Failed to create output directory: %v", err)
		os.Exit(1)
	}

	filename := filepath.Join(outDir, "26_vba_macros.pptm")
	if err := os.WriteFile(filename, data, 0o600); err != nil {
		log.Printf("Failed to save .pptm: %v", err)
		os.Exit(1)
	}

	log.Printf("Successfully generated %s", filename)
}
