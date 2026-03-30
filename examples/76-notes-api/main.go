// examples/76-notes-api demonstrates speaker notes on slides.
//
// Shows plain text notes via WithNotes, rich paragraph notes via WithRichNotes
// and AddNoteParagraph, bulleted note items, numbered note items, and updating
// notes on an existing slide via PresentationEditor.SetNotes.
//
// Run with: go run ./examples/76-notes-api/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

const (
	outputDir  = "examples/output"
	outputFile = "76_notes_api.pptx"
)

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

	// --- Slide 1: Plain text notes via WithNotes ---
	slide1 := pptx.NewSlide("Plain Text Notes").
		AddBullet("This slide has plain text speaker notes.").
		AddBullet("Open the notes panel to read them.").
		WithNotes("These are plain text speaker notes for slide 1.\nUse them to remind yourself of talking points.")

	// --- Slide 2: Multi-paragraph notes via AddNoteParagraph ---
	p1 := elements.NewParagraph()
	p1.Runs = []elements.Run{elements.NewRun("Opening paragraph – introduce the topic.")}

	p2 := elements.NewParagraph()
	p2.Style.BulletStyle = text.BulletStyleBullet
	p2.Runs = []elements.Run{elements.NewRun("Bullet point: key concept #1")}

	p3 := elements.NewParagraph()
	p3.Style.BulletStyle = text.BulletStyleBullet
	p3.Runs = []elements.Run{elements.NewRun("Bullet point: key concept #2")}

	slide2 := pptx.NewSlide("Rich Paragraph Notes").
		AddBullet("This slide uses AddNoteParagraph for structured notes.").
		AddNoteParagraph(p1).
		AddNoteParagraph(p2).
		AddNoteParagraph(p3)

	// --- Slide 3: AddNoteBullet and AddNoteNumbered ---
	slide3 := pptx.NewSlide("Notes with Bullets & Numbers").
		AddBullet("Notes have bullet and numbered styles.").
		AddNoteBullet("First bullet note item").
		AddNoteBullet("Second bullet note item").
		AddNoteNumbered("First numbered note item").
		AddNoteNumbered("Second numbered note item")

	// --- Slide 4: AddNoteSubBullet – indented notes ---
	slide4 := pptx.NewSlide("Sub-Bullet Notes").
		AddBullet("Notes support sub-bullets at different indent levels.").
		AddNoteBullet("Top-level note bullet").
		AddNoteSubBullet(1, "Indented level 1 note").
		AddNoteSubBullet(2, "Indented level 2 note").
		AddNoteSubBullet(1, "Back to level 1")

	// --- Slide 5: WithRichNotes – full paragraph array ---
	rp1 := elements.NewParagraph()
	rp1.Runs = []elements.Run{
		func() elements.Run { r := elements.NewRun("Bold intro: "); r.Bold = true; return r }(),
		elements.NewRun("remember to smile during the presentation."),
	}

	rp2 := elements.NewParagraph()
	rp2.Style.BulletStyle = text.BulletStyleNumber
	rp2.Runs = []elements.Run{elements.NewRun("Reference slide 7 for the comparison data.")}

	slide5 := pptx.NewSlide("WithRichNotes").
		AddBullet("WithRichNotes sets the full notes body in one call.").
		WithRichNotes([]elements.Paragraph{rp1, rp2})

	// --- Build base presentation ---
	tmpDir, err := os.MkdirTemp("", "gopptx-notes-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpPath := filepath.Join(tmpDir, "base.pptx")
	builder := pptx.NewPresentationBuilder("Notes API Demo")
	for _, s := range []pptx.SlideContent{slide1, slide2, slide3, slide4, slide5} {
		builder.AddSlide(s)
	}
	if err := builder.WriteToFile(tmpPath); err != nil {
		return fmt.Errorf("write base: %w", err)
	}

	// --- Part 2: Edit notes via PresentationEditor.SetNotes ---
	ed, err := pptx.OpenPresentationEditor(tmpPath)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	defer ed.Close()

	if err := ed.SetNotes(0, "Overwritten by editor.SetNotes – slide 1 notes updated programmatically."); err != nil {
		return fmt.Errorf("set notes slide 0: %w", err)
	}
	log.Printf("Updated notes on slide 0 via SetNotes.\n")

	// Add a new slide with notes set at add-time.
	_, err = ed.AddSlide(
		pptx.NewSlide("Added by Editor").
			AddBullet("This slide was added via editor.AddSlide.").
			WithNotes("Notes written directly at AddSlide time via WithNotes."),
	)
	if err != nil {
		return fmt.Errorf("add slide: %w", err)
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err := ed.Save(outputPath); err != nil {
		return fmt.Errorf("save: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
