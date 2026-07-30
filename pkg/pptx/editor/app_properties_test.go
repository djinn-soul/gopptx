package editor

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// Issue #105 / #131: docProps/app.xml was generated when a deck was created but
// never rewritten on an editor save, so <Slides> kept the count the deck had
// when PowerPoint last touched it. Gmail and the Windows shell preview from that
// number, which is why a deck grown through the editor previewed short.

func savedPart(t *testing.T, data []byte, name string) []byte {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("read saved package: %v", err)
	}
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}
		defer func() { _ = rc.Close() }()
		content, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		return content
	}
	t.Fatalf("saved package has no %s", name)
	return nil
}

func TestSave_WritesAppPropertiesSlideCount(t *testing.T) {
	base := writeDeckFixture(t, "app-props-count.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1"),
	})

	ed, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	for i := 2; i <= 4; i++ {
		if _, err := ed.AddSlide(elements.NewSlide("Added")); err != nil {
			t.Fatalf("add slide %d: %v", i, err)
		}
	}

	saved, err := ed.SaveToBytes()
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	appXML := string(savedPart(t, saved, "docProps/app.xml"))
	if !strings.Contains(appXML, "<Slides>4</Slides>") {
		t.Fatalf("expected <Slides>4</Slides> after adding three slides, got:\n%s", appXML)
	}

	contentTypes := string(savedPart(t, saved, "[Content_Types].xml"))
	if !strings.Contains(contentTypes, `PartName="/docProps/app.xml"`) {
		t.Fatalf("app.xml written without a content type override:\n%s", contentTypes)
	}

	rootRels := string(savedPart(t, saved, "_rels/.rels"))
	if !strings.Contains(rootRels, "extended-properties") ||
		!strings.Contains(rootRels, `Target="docProps/app.xml"`) {
		t.Fatalf("app.xml written without a package relationship:\n%s", rootRels)
	}
}

func TestSave_AppPropertiesCountsNotesAndHiddenSlides(t *testing.T) {
	base := writeDeckFixture(t, "app-props-notes.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1"),
		elements.NewSlide("Slide 2"),
		elements.NewSlide("Slide 3"),
	})

	ed, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	if err := ed.SetNotes(0, "first note"); err != nil {
		t.Fatalf("set notes on slide 1: %v", err)
	}
	if err := ed.SetNotes(2, "third note"); err != nil {
		t.Fatalf("set notes on slide 3: %v", err)
	}
	if err := ed.SetSlideHidden(1, true); err != nil {
		t.Fatalf("hide slide 2: %v", err)
	}

	saved, err := ed.SaveToBytes()
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	appXML := string(savedPart(t, saved, "docProps/app.xml"))
	for _, want := range []string{"<Slides>3</Slides>", "<Notes>2</Notes>", "<HiddenSlides>1</HiddenSlides>"} {
		if !strings.Contains(appXML, want) {
			t.Fatalf("expected %s in app.xml, got:\n%s", want, appXML)
		}
	}
}

// A deck that already carries app.xml keeps the authoring metadata PowerPoint
// wrote into it. Only the counts change.
func TestSave_AppPropertiesPreservesExistingMetadata(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app-props-existing.pptx")
	fixture := writeDeckFixture(t, "app-props-source.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1"),
	})
	files := readZipEntries(t, fixture)
	files["docProps/app.xml"] = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties" ` +
		`xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes">` +
		`<Application>Microsoft Office PowerPoint</Application>` +
		`<Company>Acme $ Co</Company>` +
		`<Manager>Ada</Manager>` +
		`<Slides>1</Slides><Notes>0</Notes><HiddenSlides>0</HiddenSlides>` +
		`<TitlesOfParts><vt:vector size="1" baseType="lpstr"><vt:lpstr>Slide 1</vt:lpstr></vt:vector></TitlesOfParts>` +
		`</Properties>`
	files["[Content_Types].xml"] = strings.Replace(
		files["[Content_Types].xml"],
		"</Types>",
		`<Override PartName="/docProps/app.xml" `+
			`ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/></Types>`,
		1,
	)
	if err := writeZipFixture(path, files); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	ed, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	if _, err := ed.AddSlide(elements.NewSlide("Slide 2")); err != nil {
		t.Fatalf("add slide: %v", err)
	}

	saved, err := ed.SaveToBytes()
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	appXML := string(savedPart(t, saved, "docProps/app.xml"))
	if !strings.Contains(appXML, "<Slides>2</Slides>") {
		t.Fatalf("slide count not updated:\n%s", appXML)
	}
	// "$" in caller text must survive: regexp replacement expands $1/$name.
	if !strings.Contains(appXML, "<Company>Acme $ Co</Company>") {
		t.Fatalf("Company dropped or mangled:\n%s", appXML)
	}
	for _, want := range []string{"<Manager>Ada</Manager>", "<TitlesOfParts>"} {
		if !strings.Contains(appXML, want) {
			t.Fatalf("expected preserved %s in app.xml, got:\n%s", want, appXML)
		}
	}
	if strings.Count(appXML, "<Slides>") != 1 {
		t.Fatalf("expected a single <Slides> element:\n%s", appXML)
	}
}

func TestSetAppPropertyInt_AppendsMissingElement(t *testing.T) {
	const source = `<Properties xmlns="x"><Application>gopptx</Application></Properties>`

	got := setAppPropertyInt(source, appSlidesPattern, "Slides", 7)

	if !strings.Contains(got, "<Slides>7</Slides>") {
		t.Fatalf("missing element not appended: %s", got)
	}
	if !strings.HasSuffix(strings.TrimSpace(got), "</Properties>") {
		t.Fatalf("element appended outside the document root: %s", got)
	}
}

func readZipEntries(t *testing.T, path string) map[string]string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	files := make(map[string]string, len(zr.File))
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open %s: %v", f.Name, err)
		}
		content, err := io.ReadAll(rc)
		_ = rc.Close()
		if err != nil {
			t.Fatalf("read %s: %v", f.Name, err)
		}
		files[f.Name] = string(content)
	}
	return files
}
