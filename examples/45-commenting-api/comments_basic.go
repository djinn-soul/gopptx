package main

import (
	"fmt"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
)

func main() {
	if err := runCommentsBasic(); err != nil {
		log.Fatal(err)
	}
}

func runCommentsBasic() error {
	// 1. Create base
	log.Println("Creating base presentation...")
	p := pptx.NewPresentationBuilder("Comment Test")
	p.AddTitleSlide("Slide 1")
	basePath := "comment_base.pptx"
	if err := p.WriteToFile(basePath); err != nil {
		return fmt.Errorf("failed to save base: %w", err)
	}
	// We keep base for inspection if needed, or delete it.
	// defer os.Remove(basePath)

	// 2. Open Editor
	log.Println("Opening editor...")
	ed, openErr := editor.OpenPresentationEditor(basePath)
	if openErr != nil {
		return fmt.Errorf("failed to open editor: %w", openErr)
	}
	defer func() {
		if closeErr := ed.Close(); closeErr != nil {
			log.Printf("warning: failed to close editor: %v", closeErr)
		}
	}()

	// 3. Add Authors
	log.Println("Adding authors...")
	alice, aliceErr := ed.AddAuthor("Alice Inchains", "AI")
	if aliceErr != nil {
		return fmt.Errorf("failed to add author Alice: %w", aliceErr)
	}
	bob, bobErr := ed.AddAuthor("Bob Builder", "BB")
	if bobErr != nil {
		return fmt.Errorf("failed to add author Bob: %w", bobErr)
	}

	// 4. Add Comments
	// Slide 1 (index 0)
	log.Println("Adding comments...")
	if err := ed.AddComment(0, alice.ID, "This title is too short.", 100, 100); err != nil {
		return fmt.Errorf("failed to add Alice comment: %w", err)
	}
	if err := ed.AddComment(0, bob.ID, "I think it's fine, honestly.", 200, 200); err != nil {
		return fmt.Errorf("failed to add Bob comment: %w", err)
	}

	// 5. Save
	outPath := "comment_output.pptx"
	log.Printf("Saving to %s...\n", outPath)
	if err := ed.Save(outPath); err != nil {
		return fmt.Errorf("failed to save output: %w", err)
	}
	log.Println("Done!")
	return nil
}
