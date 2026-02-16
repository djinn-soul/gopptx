package editor

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestPresentationEditor_Sections(t *testing.T) {
	initial := []elements.SlideContent{
		elements.NewSlide("Slide 0"),
		elements.NewSlide("Slide 1"),
		elements.NewSlide("Slide 2"),
		elements.NewSlide("Slide 3"),
	}
	path := writeDeckFixture(t, "sections_base.pptx", initial)

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	// 1. Test AddSection
	if addIntroErr := editor.AddSection("Intro", []int{0}); addIntroErr != nil {
		t.Errorf("failed to add intro section: %v", addIntroErr)
	}
	if addContentErr := editor.AddSection("Content", []int{1, 2}); addContentErr != nil {
		t.Errorf("failed to add content section: %v", addContentErr)
	}

	if len(editor.Sections()) != 2 {
		t.Errorf("expected 2 sections, got %d", len(editor.Sections()))
	}

	// 2. Test RenameSection
	if renameErr := editor.RenameSection("Intro", "Introduction"); renameErr != nil {
		t.Errorf("failed to rename section: %v", renameErr)
	}
	if editor.Sections()[0].Name != "Introduction" {
		t.Errorf("expected section name 'Introduction', got %q", editor.Sections()[0].Name)
	}

	// 3. Test RemoveSection
	if removeErr := editor.RemoveSection("Content"); removeErr != nil {
		t.Errorf("failed to remove section: %v", removeErr)
	}
	if len(editor.Sections()) != 1 {
		t.Errorf("expected 1 section after removal, got %d", len(editor.Sections()))
	}

	// 4. Test Error Cases
	if addErr := editor.AddSection("", []int{0}); addErr == nil {
		t.Error("expected error for empty section name")
	}
	if addErr := editor.AddSection("Invalid", []int{99}); addErr == nil {
		t.Error("expected error for out of range slide index")
	}
	if renameErr := editor.RenameSection("NonExistent", "New"); renameErr == nil {
		t.Error("expected error for renaming non-existent section")
	}
}

func TestPresentationEditor_SectionsPersistence(t *testing.T) {
	initial := []elements.SlideContent{
		elements.NewSlide("S1"),
		elements.NewSlide("S2"),
	}
	path := writeDeckFixture(t, "sections_persistence.pptx", initial)

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}

	if addErr := editor.AddSection("Main", []int{0, 1}); addErr != nil {
		t.Fatalf("add section: %v", addErr)
	}

	outPath := filepath.Join(t.TempDir(), "sections_saved.pptx")
	if saveErr := editor.Save(outPath); saveErr != nil {
		t.Fatalf("save: %v", saveErr)
	}
	_ = editor.Close()

	// Verify persistence
	reopened, err := OpenPresentationEditor(outPath)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer func() {
		if closeErr := reopened.Close(); closeErr != nil {
			t.Errorf("close reopened editor: %v", closeErr)
		}
	}()

	sections := reopened.Sections()
	if len(sections) != 1 {
		t.Fatalf("expected 1 section in reopened deck, got %d", len(sections))
	}
	if sections[0].Name != "Main" {
		t.Errorf("expected section name 'Main', got %q", sections[0].Name)
	}
	if len(sections[0].SlideIDs) != 2 {
		t.Errorf("expected 2 slides in section, got %d", len(sections[0].SlideIDs))
	}

	// Verify XML parts exist
	contentTypes := string(readZipFileBytes(t, outPath, "[Content_Types].xml"))
	if !strings.Contains(contentTypes, "ppt/sectionList.xml") {
		t.Error("Content_Types.xml missing sectionList override")
	}

	presRels := string(readZipFileBytes(t, outPath, "ppt/_rels/presentation.xml.rels"))
	if !strings.Contains(presRels, "sectionList.xml") {
		t.Error("presentation.xml.rels missing sectionList relationship")
	}

	sectionList := string(readZipFileBytes(t, outPath, "ppt/sectionList.xml"))
	if !strings.Contains(sectionList, `name="Main"`) {
		t.Error("sectionList.xml missing section entry")
	}

	presXML := string(readZipFileBytes(t, outPath, "ppt/presentation.xml"))
	if !strings.Contains(presXML, "<p14:sectionLst") {
		t.Error("presentation.xml missing p14:sectionLst extension")
	}
}

func TestPresentationEditor_SectionsPreservedOnMove(t *testing.T) {
	initial := []elements.SlideContent{
		elements.NewSlide("A"),
		elements.NewSlide("B"),
	}
	path := writeDeckFixture(t, "sections_move.pptx", initial)

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() {
		if closeErr := editor.Close(); closeErr != nil {
			t.Errorf("close editor: %v", closeErr)
		}
	}()

	// Capture slide IDs
	slides := editor.Slides()
	idA := slides[0].SlideID
	idB := slides[1].SlideID

	if addAErr := editor.AddSection("Section A", []int{0}); addAErr != nil {
		t.Fatalf("add section A: %v", addAErr)
	}
	if addBErr := editor.AddSection("Section B", []int{1}); addBErr != nil {
		t.Fatalf("add section B: %v", addBErr)
	}

	// Move slide 0 to index 1: [B, A]
	if moveErr := editor.MoveSlide(0, 1); moveErr != nil {
		t.Fatalf("move slide: %v", moveErr)
	}

	// Sections track SlideIDs, so they should still point to the same slides
	sections := editor.Sections()
	if sections[0].SlideIDs[0] != idA {
		t.Errorf("Section A expected SlideID %d, got %d", idA, sections[0].SlideIDs[0])
	}
	if sections[1].SlideIDs[0] != idB {
		t.Errorf("Section B expected SlideID %d, got %d", idB, sections[1].SlideIDs[0])
	}
}
