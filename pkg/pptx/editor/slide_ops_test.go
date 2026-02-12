package editor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestPresentationEditorDuplicateSlide(t *testing.T) {
	initial := []elements.SlideContent{
		elements.NewSlide("Slide 1"),
		elements.NewSlide("Slide 2"),
	}
	path := writeDeckFixture(t, "duplicate_test.pptx", initial)

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}

	// Duplicate Slide 1 to the end
	if _, err := editor.DuplicateSlide(0, 2); err != nil {
		t.Fatalf("duplicate slide 0 to 2: %v", err)
	}

	if editor.SlideCount() != 3 {
		t.Fatalf("expected 3 slides, got %d", editor.SlideCount())
	}

	slides := editor.Slides()
	if slides[2].Title != "Slide 1 (Copy)" {
		t.Fatalf("expected Slide 1 (Copy) at index 2, got %q", slides[2].Title)
	}

	// Duplicate Slide 2 between 1 and its copy
	if _, err := editor.DuplicateSlide(1, 1); err != nil {
		t.Fatalf("duplicate slide 1 to 1: %v", err)
	}

	if editor.SlideCount() != 4 {
		t.Fatalf("expected 4 slides, got %d", editor.SlideCount())
	}

	slides = editor.Slides()
	// Order: [Slide 1, Slide 2 (Copy), Slide 2, Slide 1 (Copy)]
	if slides[1].Title != "Slide 2 (Copy)" {
		t.Fatalf("expected Slide 2 (Copy) at index 1, got %q", slides[1].Title)
	}

	outPath := filepath.Join(t.TempDir(), "duplicated.pptx")
	if err := editor.Save(outPath); err != nil {
		t.Fatalf("save deck: %v", err)
	}

	// Reopen and check
	reopened, err := OpenPresentationEditor(outPath)
	if err != nil {
		t.Fatalf("reopen deck: %v", err)
	}
	if reopened.SlideCount() != 4 {
		t.Fatalf("reopened: expected 4 slides, got %d", reopened.SlideCount())
	}

	// Check XML contents of the duplicated slide to ensure it preserved titles
	// Order: [Slide 1, Slide 2 (Copy), Slide 2, Slide 1 (Copy)]
	// slides[1] (index 1) is Slide 2 (Copy)
	slide2Part := reopened.Slides()[1].PartName
	slide2XML := string(readZipFileBytes(t, outPath, slide2Part))
	if !strings.Contains(slide2XML, "Slide 2") {
		t.Errorf("duplicated slide XML does not contain expected title content (expected Slide 2)")
	}
}

func TestDuplicateSlideWithImage(t *testing.T) {
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "sample.png")
	if err := os.WriteFile(imgPath, []byte("fake content"), 0o600); err != nil {
		t.Fatalf("write image fixture: %v", err)
	}

	initial := []elements.SlideContent{
		elements.NewSlide("Image Slide").AddImage(shapes.NewImage(imgPath, 914400, 914400, 1828800, 1828800)),
		elements.NewSlide("Other"),
	}
	path := writeDeckFixture(t, "duplicate_image_test.pptx", initial)

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}

	// Duplicate Image Slide
	if _, err := editor.DuplicateSlide(0, 2); err != nil {
		t.Fatalf("duplicate image slide: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "duplicated_image.pptx")
	if err := editor.Save(outPath); err != nil {
		t.Fatalf("save: %v", err)
	}

	// Verify rels for the copy
	copyRef := editor.Slides()[2]
	relsXML := string(readZipFileBytes(t, outPath, common.SlideRelsPartName(copyRef.PartName)))
	if !strings.Contains(relsXML, "media/image1.png") {
		t.Errorf("duplicated slide relationship XML does not contain reference to image1.png")
	}
}

func TestPresentationEditorMoveSlide(t *testing.T) {
	initial := []elements.SlideContent{
		elements.NewSlide("A"),
		elements.NewSlide("B"),
		elements.NewSlide("C"),
	}
	path := writeDeckFixture(t, "move_test.pptx", initial)

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}

	// Move B to start: [B, A, C]
	if err := editor.MoveSlide(1, 0); err != nil {
		t.Fatalf("move 1 to 0: %v", err)
	}

	slides := editor.Slides()
	if slides[0].Title != "B" || slides[1].Title != "A" || slides[2].Title != "C" {
		t.Fatalf("order after move 1->0: %s, %s, %s", slides[0].Title, slides[1].Title, slides[2].Title)
	}

	// Move A to end: [B, C, A]
	if err := editor.MoveSlide(1, 2); err != nil {
		t.Fatalf("move 1 to 2: %v", err)
	}

	slides = editor.Slides()
	if slides[0].Title != "B" || slides[1].Title != "C" || slides[2].Title != "A" {
		t.Fatalf("order after move 1->2: %s, %s, %s", slides[0].Title, slides[1].Title, slides[2].Title)
	}

	outPath := filepath.Join(t.TempDir(), "moved.pptx")
	if err := editor.Save(outPath); err != nil {
		t.Fatalf("save deck: %v", err)
	}

	reopened, err := OpenPresentationEditor(outPath)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	slides = reopened.Slides()
	if slides[0].Title != "B" || slides[1].Title != "C" || slides[2].Title != "A" {
		t.Fatalf("reopened order: %s, %s, %s", slides[0].Title, slides[1].Title, slides[2].Title)
	}
}
