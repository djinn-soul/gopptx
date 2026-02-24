package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

const (
	outputDirComments = "examples/output"
	baseFile          = "45_comments_base.pptx"
	finalFile         = "45_comments_smoke.pptx"
)

func mainCommentingSmoke() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDirComments, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	basePath := filepath.Join(outputDirComments, baseFile)
	finalPath := filepath.Join(outputDirComments, finalFile)

	// 1. Create a base presentation
	baseSlides := []pptx.SlideContent{
		pptx.NewSlide("Slide 1 (Team Review)").AddBullet("Team member A should comment here."),
		pptx.NewSlide("Slide 2 (Client Review)").AddBullet("Client feedback wanted here."),
	}
	if err := pptx.WriteFile(basePath, "Comments API Smoke Test", baseSlides); err != nil {
		return fmt.Errorf("create base file: %w", err)
	}
	log.Printf("1. Created base: %s\n", basePath)

	// 2. Open with Editor
	editor, openErr := pptx.OpenPresentationEditor(basePath)
	if openErr != nil {
		return fmt.Errorf("open editor: %w", openErr)
	}
	log.Println("2. Opened presentation with Editor")

	// 3. Add Authors
	authorA, err := editor.AddAuthor("Alice Reviewer", "AR")
	if err != nil {
		return fmt.Errorf("add author Alice: %w", err)
	}
	authorB, err := editor.AddAuthor("Bob Architect", "BA")
	if err != nil {
		return fmt.Errorf("add author Bob: %w", err)
	}
	log.Printf("3. Registered authors: Alice (ID=%d), Bob (ID=%d)\n", authorA.ID, authorB.ID)

	// 4. Add comments to Slide 1 (Index 0)
	err = editor.AddComment(0, authorA.ID, "Looks good, but check the font size.", 500000, 500000)
	if err != nil {
		return fmt.Errorf("add alice comment on slide 1: %w", err)
	}
	err = editor.AddComment(0, authorB.ID, "Agreed. Let's make it 24pt.", 600000, 600000)
	if err != nil {
		return fmt.Errorf("add bob comment on slide 1: %w", err)
	}
	log.Println("4. Added comments to Slide 1 from Alice and Bob")

	// 5. Add comment to Slide 2 (Index 1)
	err = editor.AddComment(1, authorA.ID, "Need comparison chart here.", 1000000, 1000000)
	if err != nil {
		return fmt.Errorf("add alice comment on slide 2: %w", err)
	}
	log.Println("5. Added comment to Slide 2 from Alice")

	// 6. List and Remove logic verification
	comments, err := editor.GetComments(0)
	if err != nil {
		return fmt.Errorf("get comments slide 0: %w", err)
	}
	log.Printf("6. Verified %d comments on Slide 1\n", len(comments))

	// 7. Save final result
	if err := editor.Save(finalPath); err != nil {
		return fmt.Errorf("save modified: %w", err)
	}
	log.Printf("7. Saved final presentation with comments: %s\n", finalPath)

	return nil
}

var _ = mainCommentingSmoke
