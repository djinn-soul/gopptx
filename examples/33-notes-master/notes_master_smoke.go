package main

import (
	"log"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	// 1. Create a custom Notes Master
	nm := pptx.NewNotesMaster().
		WithHeader("CONFIDENTIAL - internal use only").
		WithFooter("Notes Master Smoke Test").
		WithDateTime(true).
		WithSlideNumber(true)

	// 2. Set default text style for notes levels
	nm.WithBodyStyle([]pptx.TextLevelStyle{
		{Level: 0, SizePt: 12, Color: "0000FF", Bold: true}, // Lvl 1: Blue, Bold
		{Level: 1, SizePt: 10, Color: "333333"},             // Lvl 2: Dark Grey
	})

	// 3. Create presentation metadata
	meta := pptx.Metadata{
		Metadata:    pptx.MetadataFields{Title: "Notes Master Smoke"},
		NotesMaster: nm,
	}

	// 4. Create slides with speaker notes
	slides := []pptx.SlideContent{
		pptx.NewSlide("Notes Master Demo").
			WithNotes("This is level 1 notes text.\n\tThis is level 2 notes text."),
		pptx.NewSlide("Another Slide").
			WithNotes("Check the header and footer on the notes page!"),
	}

	// 5. Generate PPTX
	data, buildErr := pptx.CreateWithMetadata(meta, slides)
	if buildErr != nil {
		log.Fatalf("failed to create pptx: %v", buildErr)
	}

	outputPath := "examples/output/33_notes_master_smoke.pptx"
	if err := os.WriteFile(outputPath, data, 0o644); err != nil {
		log.Fatalf("failed to write file: %v", err)
	}

	log.Printf("Generated %s\n", outputPath)
}
