package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesEmbedsSpeakerNotesParts(t *testing.T) {
	slides := []SlideContent{
		NewSlide("Intro").AddBullet("Content").WithNotes("First line\nSecond line"),
		NewSlide("Plain").AddBullet("No notes"),
	}

	data, err := CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	required := []string{
		"ppt/notesSlides/notesSlide1.xml",
		"ppt/notesSlides/_rels/notesSlide1.xml.rels",
		"ppt/notesMasters/notesMaster1.xml",
		"ppt/notesMasters/_rels/notesMaster1.xml.rels",
	}
	for _, name := range required {
		if !zipHasFile(zr, name) {
			t.Fatalf("missing expected notes part: %s", name)
		}
	}
	if zipHasFile(zr, "ppt/notesSlides/notesSlide2.xml") {
		t.Fatalf("did not expect notes slide for slide without notes")
	}

	contentTypes := readZipFile(t, zr, "[Content_Types].xml")
	for _, needle := range []string{
		`/ppt/notesSlides/notesSlide1.xml`,
		`presentationml.notesSlide+xml`,
		`/ppt/notesMasters/notesMaster1.xml`,
		`presentationml.notesMaster+xml`,
	} {
		if !strings.Contains(contentTypes, needle) {
			t.Fatalf("expected %q in [Content_Types].xml", needle)
		}
	}

	presRels := readZipFile(t, zr, "ppt/_rels/presentation.xml.rels")
	if !strings.Contains(presRels, `relationships/notesMaster`) {
		t.Fatalf("expected notes master relationship in presentation rels")
	}

	slide1Rels := readZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(slide1Rels, `relationships/notesSlide`) {
		t.Fatalf("expected notes slide relationship on slide1")
	}
	if !strings.Contains(slide1Rels, `Target="../notesSlides/notesSlide1.xml"`) {
		t.Fatalf("expected notes slide target for slide1")
	}

	slide2Rels := readZipFile(t, zr, "ppt/slides/_rels/slide2.xml.rels")
	if strings.Contains(slide2Rels, `relationships/notesSlide`) {
		t.Fatalf("did not expect notes slide relationship on slide2")
	}

	notesXML := readZipFile(t, zr, "ppt/notesSlides/notesSlide1.xml")
	for _, needle := range []string{
		`<p:notes`,
		`<a:t>First line</a:t>`,
		`<a:t>Second line</a:t>`,
	} {
		if !strings.Contains(notesXML, needle) {
			t.Fatalf("expected %q in notes slide XML", needle)
		}
	}

	notesRels := readZipFile(t, zr, "ppt/notesSlides/_rels/notesSlide1.xml.rels")
	for _, needle := range []string{
		`Target="../slides/slide1.xml"`,
		`relationships/notesMaster`,
	} {
		if !strings.Contains(notesRels, needle) {
			t.Fatalf("expected %q in notes slide rels", needle)
		}
	}

	appXML := readZipFile(t, zr, "docProps/app.xml")
	if !strings.Contains(appXML, `<Notes>1</Notes>`) {
		t.Fatalf("expected notes count in app props")
	}
}

func TestCreateWithSlidesMarkdownNotesPersistence(t *testing.T) {
	input := `# Topic
- Bullet
> Speaker note text`
	slides, err := SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}
	if strings.TrimSpace(slides[0].Notes) == "" {
		t.Fatalf("expected parsed notes from markdown blockquote")
	}

	data, err := CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	notesXML := readZipFile(t, zr, "ppt/notesSlides/notesSlide1.xml")
	if !strings.Contains(notesXML, `<a:t>Speaker note text</a:t>`) {
		t.Fatalf("expected markdown note text persisted into notes slide xml")
	}
}
