package editor

import (
	"archive/zip"
	"bytes"
	"testing"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// Issue #131: OPC requires [Content_Types].xml to be the first entry in the
// package, and every entry needs a valid MS-DOS timestamp. Sorted part names
// happen to put the content types stream first, because "[" sorts below "_" and
// the lowercase directory names — but a source deck carrying a part whose name
// starts with an uppercase letter would break that, so the order is explicit.

func TestSave_ContentTypesIsFirstEntry(t *testing.T) {
	base := writeDeckFixture(t, "package-order.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1"),
		elements.NewSlide("Slide 2"),
	})

	ed, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	saved, err := ed.SaveToBytes()
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(saved), int64(len(saved)))
	if err != nil {
		t.Fatalf("read saved package: %v", err)
	}
	if len(zr.File) == 0 {
		t.Fatal("saved package is empty")
	}
	if zr.File[0].Name != pptxxml.ContentTypesPartName {
		t.Fatalf("expected %s first, got %q", pptxxml.ContentTypesPartName, zr.File[0].Name)
	}

	for _, f := range zr.File {
		if f.Modified.Month() < 1 || f.Modified.Day() < 1 {
			t.Fatalf("entry %q has an invalid MS-DOS date: %s", f.Name, f.Modified)
		}
	}
}

func TestOrderPackageNames_HoistsContentTypesAheadOfUppercaseParts(t *testing.T) {
	// "Fonts/font1.fntdata" sorts before "[Content_Types].xml": "F" is 0x46,
	// "[" is 0x5B. Sorting alone would bury the content types stream.
	names := []string{"Fonts/font1.fntdata", "[Content_Types].xml", "_rels/.rels", "ppt/presentation.xml"}

	got := orderPackageNames(names)

	want := []string{"[Content_Types].xml", "_rels/.rels", "Fonts/font1.fntdata", "ppt/presentation.xml"}
	if len(got) != len(want) {
		t.Fatalf("entry lost or duplicated: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("wrong order: got %v, want %v", got, want)
		}
	}
}

func TestOrderPackageNames_LeavesOrderedInputUntouched(t *testing.T) {
	names := []string{"[Content_Types].xml", "_rels/.rels", "docProps/app.xml", "ppt/presentation.xml"}

	got := orderPackageNames(names)

	if &got[0] != &names[0] {
		t.Fatal("expected the already-ordered slice to be returned as is")
	}
}
