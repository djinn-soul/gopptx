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
	defer func() { _ = editor.Close() }()

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
	defer func() { _ = reopened.Close() }()
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
	defer func() { _ = editor.Close() }()

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
	defer func() { _ = editor.Close() }()

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
	defer func() { _ = reopened.Close() }()
	slides = reopened.Slides()
	if slides[0].Title != "B" || slides[1].Title != "C" || slides[2].Title != "A" {
		t.Fatalf("reopened order: %s, %s, %s", slides[0].Title, slides[1].Title, slides[2].Title)
	}
}

func TestDuplicateSlide_AppendsCopySuffixToTitlePlaceholder(t *testing.T) {
	path := writeDeckFixture(t, "duplicate_title_placeholder.pptx", []elements.SlideContent{
		elements.NewSlide("Main Title"),
	})

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	srcPart := editor.Slides()[0].PartName
	editor.parts.Set(srcPart, []byte(slideWithBodyAndTitlePlaceholderXML("Body Text", "Main Title")))

	if _, err := editor.DuplicateSlide(0, 1); err != nil {
		t.Fatalf("duplicate slide: %v", err)
	}

	copyPart := editor.Slides()[1].PartName
	copyData, _ := editor.parts.Get(copyPart)
	copyXML := string(copyData)
	if !strings.Contains(copyXML, "Main Title (Copy)") {
		t.Fatalf("expected title placeholder to include copy suffix")
	}
	if strings.Contains(copyXML, "Body Text (Copy)") {
		t.Fatalf("did not expect non-title text to be modified")
	}
}

func TestDuplicateSlide_AppendsCopySuffixToLastTitleRun(t *testing.T) {
	path := writeDeckFixture(t, "duplicate_title_multirun.pptx", []elements.SlideContent{
		elements.NewSlide("Main Title"),
	})

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	srcPart := editor.Slides()[0].PartName
	editor.parts.Set(srcPart, []byte(slideWithBodyAndMultiRunTitlePlaceholderXML("Body Text", "Main ", "Title")))

	if _, err := editor.DuplicateSlide(0, 1); err != nil {
		t.Fatalf("duplicate slide: %v", err)
	}

	copyPart := editor.Slides()[1].PartName
	copyData, _ := editor.parts.Get(copyPart)
	copyXML := string(copyData)
	if !strings.Contains(copyXML, "<a:t>Main </a:t></a:r><a:r><a:t>Title (Copy)</a:t>") {
		t.Fatalf("expected suffix on last title run, got XML: %s", copyXML)
	}
	if strings.Contains(copyXML, "<a:t>Main  (Copy)</a:t></a:r><a:r><a:t>Title</a:t>") {
		t.Fatalf("unexpected suffix on first title run")
	}
	if strings.Contains(copyXML, "Body Text (Copy)") {
		t.Fatalf("did not expect non-title text to be modified")
	}
}

func TestSetSlideTitle_TargetsTitlePlaceholder(t *testing.T) {
	path := writeDeckFixture(t, "set_title_placeholder.pptx", []elements.SlideContent{
		elements.NewSlide("Main Title"),
	})

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	part := editor.Slides()[0].PartName
	editor.parts.Set(part, []byte(slideWithBodyAndTitlePlaceholderXML("Body Text", "Main Title")))

	if err := editor.SetSlideTitle(0, "Renamed Title"); err != nil {
		t.Fatalf("set slide title: %v", err)
	}

	xmlData, _ := editor.parts.Get(part)
	xmlStr := string(xmlData)
	if !strings.Contains(xmlStr, "Renamed Title") {
		t.Fatalf("expected updated title text in slide XML")
	}
	if strings.Contains(xmlStr, "Body Text (Copy)") {
		t.Fatalf("unexpected body text mutation")
	}
	if !strings.Contains(xmlStr, "Body Text") {
		t.Fatalf("expected non-title text to remain unchanged")
	}
}

func slideWithBodyAndTitlePlaceholderXML(bodyText, titleText string) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
		`<p:cSld><p:spTree>` +
		`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
		`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Body"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr><p:spPr/>` +
		`<p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:t>` + bodyText + `</a:t></a:r></a:p></p:txBody></p:sp>` +
		`<p:sp><p:nvSpPr><p:cNvPr id="3" name="Title 1"/><p:cNvSpPr/><p:nvPr><p:ph type="title"/></p:nvPr></p:nvSpPr><p:spPr/>` +
		`<p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:t>` + titleText + `</a:t></a:r></a:p></p:txBody></p:sp>` +
		`</p:spTree></p:cSld></p:sld>`
}

func slideWithBodyAndMultiRunTitlePlaceholderXML(bodyText, titlePrefix, titleSuffix string) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
		`<p:cSld><p:spTree>` +
		`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
		`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Body"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr><p:spPr/>` +
		`<p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:t>` + bodyText + `</a:t></a:r></a:p></p:txBody></p:sp>` +
		`<p:sp><p:nvSpPr><p:cNvPr id="3" name="Title 1"/><p:cNvSpPr/><p:nvPr><p:ph type="title"/></p:nvPr></p:nvSpPr><p:spPr/>` +
		`<p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:t>` + titlePrefix + `</a:t></a:r><a:r><a:t>` + titleSuffix + `</a:t></a:r></a:p></p:txBody></p:sp>` +
		`</p:spTree></p:cSld></p:sld>`
}
