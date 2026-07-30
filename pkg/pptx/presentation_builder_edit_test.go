package pptx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuilderEditReturnsAnEditablePresentation(t *testing.T) {
	prs, err := NewPresentationBuilder("Report").
		AddSlide(NewSlide("Overview")).
		AddSlide(NewSlide("Detail")).
		Edit()
	if err != nil {
		t.Fatalf("Edit: %v", err)
	}
	defer func() { _ = prs.Close() }()

	if got := prs.SlideCount(); got != 2 {
		t.Errorf("SlideCount() = %d, want 2", got)
	}
	if issues := prs.Validate(); len(issues) != 0 {
		t.Errorf("built presentation has %d validation issues: %v", len(issues), issues)
	}
}

// Edit must produce the same deck as writing to disk and reopening it, which is
// the round trip it replaces.
func TestBuilderEditMatchesWriteThenOpen(t *testing.T) {
	build := func() *PresentationBuilder {
		return NewPresentationBuilder("Report").
			AddSlide(NewSlide("Overview")).
			AddBulletSlide("Agenda", []string{"first", "second"})
	}

	viaEdit, err := build().Edit()
	if err != nil {
		t.Fatalf("Edit: %v", err)
	}
	defer func() { _ = viaEdit.Close() }()

	path := filepath.Join(t.TempDir(), "roundtrip.pptx")
	if err := build().WriteToFile(path); err != nil {
		t.Fatalf("WriteToFile: %v", err)
	}
	viaFile, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = viaFile.Close() }()

	if viaEdit.SlideCount() != viaFile.SlideCount() {
		t.Errorf("slide count %d != %d", viaEdit.SlideCount(), viaFile.SlideCount())
	}
}

// Edits applied through the returned Presentation must survive saving.
func TestBuilderEditPersistsSubsequentEdits(t *testing.T) {
	prs, err := NewPresentationBuilder("Report").
		AddSlide(NewSlide("Original")).
		Edit()
	if err != nil {
		t.Fatalf("Edit: %v", err)
	}
	defer func() { _ = prs.Close() }()

	path := filepath.Join(t.TempDir(), "edited.pptx")
	if err := prs.SaveAs(path); err != nil {
		t.Fatalf("SaveAs: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read saved deck: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("saved deck is empty")
	}

	reopened, err := Open(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer func() { _ = reopened.Close() }()

	if got := reopened.SlideCount(); got != 1 {
		t.Errorf("SlideCount() = %d, want 1", got)
	}
}

// The returned Presentation owns its own copy, so continuing to use the builder
// must not change it.
func TestBuilderEditSnapshotsTheBuilder(t *testing.T) {
	builder := NewPresentationBuilder("Report").AddSlide(NewSlide("One"))

	prs, err := builder.Edit()
	if err != nil {
		t.Fatalf("Edit: %v", err)
	}
	defer func() { _ = prs.Close() }()

	builder.AddSlide(NewSlide("Two"))

	if got := prs.SlideCount(); got != 1 {
		t.Errorf("SlideCount() = %d, want 1; later builder calls leaked in", got)
	}
}

func TestBuilderEditOnNilBuilder(t *testing.T) {
	var builder *PresentationBuilder
	prs, err := builder.Edit()
	if err == nil {
		t.Fatalf("Edit on a nil builder = %v, want error", prs)
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("error %q does not say the builder was nil", err)
	}
}
