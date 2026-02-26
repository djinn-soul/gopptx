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
	// 1. Load pre-built VBA blob, or generate a temporary runtime blob and clean it up.
	vbaData, cleanup, err := loadVBABlob()
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		log.Printf("Failed to prepare vbaProject.bin data: %v", err)
		os.Exit(1)
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

func loadVBABlob() ([]byte, func(), error) {
	const assetPath = "examples/assets/vbaProject.bin"

	if data, err := os.ReadFile(assetPath); err == nil {
		return data, nil, nil
	}

	tmpFile, err := os.CreateTemp("", "gopptx_vbaProject_*.bin")
	if err != nil {
		return nil, nil, err
	}
	tmpPath := tmpFile.Name()
	cleanup := func() { _ = os.Remove(tmpPath) }
	defer tmpFile.Close()

	// Runtime fallback for demo-only generation when no real VBA asset exists.
	if _, err := tmpFile.Write([]byte("dummy vba data")); err != nil {
		cleanup()
		return nil, nil, err
	}

	data, err := os.ReadFile(tmpPath)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	log.Printf("examples/assets/vbaProject.bin not found; using runtime-generated temporary VBA blob.")
	return data, cleanup, nil
}
