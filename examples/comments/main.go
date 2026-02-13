package main

import (
	"fmt"
	"log"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
)

func main() {
	// 1. Create base
	fmt.Println("Creating base presentation...")
	p := pptx.NewPresentationBuilder("Comment Test")
	p.AddTitleSlide("Slide 1")
	basePath := "comment_base.pptx"
	if err := p.WriteToFile(basePath); err != nil {
		log.Fatalf("failed to save base: %v", err)
	}
	// We keep base for inspection if needed, or delete it.
	// defer os.Remove(basePath)

	// 2. Open Editor
	fmt.Println("Opening editor...")
	ed, err := editor.OpenPresentationEditor(basePath)
	if err != nil {
		log.Fatalf("failed to open editor: %v", err)
	}
	defer func() {
		if closeErr := ed.Close(); closeErr != nil {
			log.Printf("warning: failed to close editor: %v", closeErr)
		}
	}()

	// 3. Add Authors
	fmt.Println("Adding authors...")
	alice, err := ed.AddAuthor("Alice Inchains", "AI")
	if err != nil {
		log.Fatal(err)
	}
	bob, err := ed.AddAuthor("Bob Builder", "BB")
	if err != nil {
		log.Fatal(err)
	}

	// 4. Add Comments
	// Slide 1 (index 0)
	fmt.Println("Adding comments...")
	if err := ed.AddComment(0, alice.ID, "This title is too short.", 100, 100); err != nil {
		log.Fatal(err)
	}
	if err := ed.AddComment(0, bob.ID, "I think it's fine, honestly.", 200, 200); err != nil {
		log.Fatal(err)
	}

	// 5. Save
	outPath := "comment_output.pptx"
	fmt.Printf("Saving to %s...\n", outPath)
	if err := ed.Save(outPath); err != nil {
		log.Fatalf("failed to save output: %v", err)
	}
	fmt.Println("Done!")
}
