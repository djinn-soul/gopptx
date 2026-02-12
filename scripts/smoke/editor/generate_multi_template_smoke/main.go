package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

const (
	inputDir  = "smoke_samples/sampleppt"
	outputDir = "smoke_samples"
	file1     = "160070-labyrinth-template-16x9.pptx"
	file2     = "162301-moneybox-template-16x9.pptx"
	outFile   = "37_multi_template_duplication.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	path1 := filepath.Join(inputDir, file1)
	path2 := filepath.Join(inputDir, file2)
	destPath := filepath.Join(outputDir, outFile)

	if _, err := os.Stat(path1); os.IsNotExist(err) {
		return fmt.Errorf("source template 1 missing: %s", path1)
	}
	if _, err := os.Stat(path2); os.IsNotExist(err) {
		return fmt.Errorf("source template 2 missing: %s", path2)
	}

	// 1. Open first template
	editor, err := pptx.OpenPresentationEditor(path1)
	if err != nil {
		return fmt.Errorf("open editor 1: %w", err)
	}

	// Mark original slides from Template 1
	for i := 0; i < editor.SlideCount(); i++ {
		origTitle := editor.Slides()[i].Title
		if origTitle == "" {
			origTitle = "Slide Title"
		}
		_ = editor.SetSlideTitle(i, "[T1] "+origTitle)
	}

	// 2. Merge slides from second template
	mergeStartIdx := editor.SlideCount()
	if err := editor.MergeFromFile(path2); err != nil {
		fmt.Printf("Warning during merge: %v\n", err)
	}

	// Mark original slides from Template 2
	for i := mergeStartIdx; i < editor.SlideCount(); i++ {
		origTitle := editor.Slides()[i].Title
		if origTitle == "" {
			origTitle = "Slide Title"
		}
		_ = editor.SetSlideTitle(i, "[T2] "+origTitle)
	}

	// 3. Perform 6 Duplications
	fmt.Println("Performing 6 duplication operations...")

	// Copy 1: T1 Title to end
	_, _ = editor.DuplicateSlide(0, editor.SlideCount())

	// Copy 2: T2 Title to start
	_, _ = editor.DuplicateSlide(mergeStartIdx, 0)

	// Copy 3: A middle slide from T1
	_, _ = editor.DuplicateSlide(3, 5)

	// Copy 4: A middle slide from T2
	_, _ = editor.DuplicateSlide(8, 2)

	// Copy 5: Another T1 slide
	_, _ = editor.DuplicateSlide(6, 10)

	// Copy 6: Another T2 slide
	_, _ = editor.DuplicateSlide(4, 1)

	// 4. Move a slide just for fun
	_ = editor.MoveSlide(editor.SlideCount()-1, 5)

	// 5. Save final result
	if err := editor.Save(destPath); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	fmt.Printf("Generated multi-template duplication smoke sample with 6 copies: %s\n", destPath)
	fmt.Printf("Final slide count: %d\n", editor.SlideCount())

	return nil
}
