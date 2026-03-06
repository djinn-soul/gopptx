package editor

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestPresentationEditorAddNotesInjectsMasterAndWiring(t *testing.T) {
	basePath := writeDeckFixture(t, "clean-no-notes.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1").AddBullet("Body"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()
	if updateErr := editor.UpdateSlide(
		0,
		elements.NewSlide("Slide 1").AddBullet("Body").WithNotes("Speaker script"),
	); updateErr != nil {
		t.Fatalf("update slide with notes: %v", updateErr)
	}

	outPath := filepath.Join(t.TempDir(), "with-notes.pptx")
	if saveErr := editor.Save(outPath); saveErr != nil {
		t.Fatalf("save edited deck: %v", saveErr)
	}

	if part := readZipFileBytes(t, outPath, "ppt/notesMasters/notesMaster1.xml"); len(part) == 0 {
		t.Fatalf("expected notes master part")
	}
	masterRels := string(readZipFileBytes(t, outPath, "ppt/notesMasters/_rels/notesMaster1.xml.rels"))
	if !strings.Contains(masterRels, `Target="../theme/theme2.xml"`) {
		t.Fatalf("expected notes master rels to target theme2")
	}

	presRels := string(readZipFileBytes(t, outPath, "ppt/_rels/presentation.xml.rels"))
	if !strings.Contains(presRels, "/relationships/notesMaster") {
		t.Fatalf("expected presentation relationship to notes master")
	}
	presXML := string(readZipFileBytes(t, outPath, "ppt/presentation.xml"))
	if !strings.Contains(presXML, "<p:notesMasterIdLst>") {
		t.Fatalf("expected notesMasterIdLst in presentation.xml")
	}

	slideRels := string(readZipFileBytes(t, outPath, "ppt/slides/_rels/slide1.xml.rels"))
	if !strings.Contains(slideRels, "/relationships/notesSlide") {
		t.Fatalf("expected slide1 notes relationship")
	}
	notesXML := string(readZipFileBytes(t, outPath, "ppt/notesSlides/notesSlide1.xml"))
	if !strings.Contains(notesXML, "Speaker script") {
		t.Fatalf("expected notes content persisted")
	}
	notesRels := string(readZipFileBytes(t, outPath, "ppt/notesSlides/_rels/notesSlide1.xml.rels"))
	if !strings.Contains(notesRels, `Target="../slides/slide1.xml"`) {
		t.Fatalf("expected notes back-reference to slide1")
	}
}

func TestPresentationEditorMoveSlidePreservesNotesAttachment(t *testing.T) {
	basePath := writeTwoSlideNotesFixture(t)

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()
	if moveErr := editor.MoveSlide(0, 1); moveErr != nil {
		t.Fatalf("move slide: %v", moveErr)
	}

	outPath := filepath.Join(t.TempDir(), "moved-with-notes.pptx")
	if saveErr := editor.Save(outPath); saveErr != nil {
		t.Fatalf("save edited deck: %v", saveErr)
	}

	slide1Rels := string(readZipFileBytes(t, outPath, "ppt/slides/_rels/slide1.xml.rels"))
	if !strings.Contains(slide1Rels, `Target="../notesSlides/notesSlide1.xml"`) {
		t.Fatalf("expected moved parent slide to keep notes target")
	}
	notesRels := string(readZipFileBytes(t, outPath, "ppt/notesSlides/_rels/notesSlide1.xml.rels"))
	if !strings.Contains(notesRels, `Target="../slides/slide1.xml"`) {
		t.Fatalf("expected notes back-reference to remain on parent slide")
	}
}

func TestPresentationEditorRemoveSlideRemovesAssociatedNotes(t *testing.T) {
	basePath := writeTwoSlideNotesFixture(t)

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()
	if removeErr := editor.RemoveSlide(0); removeErr != nil {
		t.Fatalf("remove slide: %v", removeErr)
	}

	outPath := filepath.Join(t.TempDir(), "removed-with-notes.pptx")
	if saveErr := editor.Save(outPath); saveErr != nil {
		t.Fatalf("save edited deck: %v", saveErr)
	}

	if part := readZipFileBytes(t, outPath, "ppt/notesSlides/notesSlide1.xml"); len(part) != 0 {
		t.Fatalf("expected removed slide notes part to be deleted")
	}
	if part := readZipFileBytes(t, outPath, "ppt/notesSlides/_rels/notesSlide1.xml.rels"); len(part) != 0 {
		t.Fatalf("expected removed slide notes rels part to be deleted")
	}

	contentTypes := string(readZipFileBytes(t, outPath, "[Content_Types].xml"))
	if strings.Contains(contentTypes, "/ppt/notesSlides/notesSlide1.xml") {
		t.Fatalf("did not expect removed notes slide override in content types")
	}
}

func TestPresentationEditorUpdateNotesMasterSupportsFlagsAndPictureBackground(t *testing.T) {
	basePath := writeDeckFixture(t, "notes-master-update.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1").AddBullet("Body"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	img := shapes.NewImageFromBytes(testutil.TinyPNG(), "png", 0, 0, 0, 0)
	bg := elements.NewPictureBackground(img)
	master := elements.NewNotesMaster().
		WithHeader("CONFIDENTIAL").
		WithFooter("Internal").
		WithDateTime(false).
		WithSlideNumber(false).
		WithBackground(bg)

	if err := editor.UpdateNotesMaster(master); err != nil {
		t.Fatalf("update notes master: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "notes-master-updated.pptx")
	if err := editor.Save(outPath); err != nil {
		t.Fatalf("save edited deck: %v", err)
	}

	notesMasterXML := string(readZipFileBytes(t, outPath, "ppt/notesMasters/notesMaster1.xml"))
	if strings.Contains(notesMasterXML, `type="dt"`) {
		t.Fatalf("did not expect date placeholder in notes master")
	}
	if strings.Contains(notesMasterXML, `type="sldNum"`) {
		t.Fatalf("did not expect slide number placeholder in notes master")
	}
	if !strings.Contains(notesMasterXML, `a:blip r:embed="rId2"`) {
		t.Fatalf("expected notes master picture background relationship in xml")
	}

	masterRels := string(readZipFileBytes(t, outPath, "ppt/notesMasters/_rels/notesMaster1.xml.rels"))
	if !strings.Contains(masterRels, `Target="../theme/theme2.xml"`) {
		t.Fatalf("expected notes master rels to target theme2")
	}
	if !strings.Contains(masterRels, `Id="rId2"`) || !strings.Contains(masterRels, `/relationships/image`) {
		t.Fatalf("expected image relationship in notes master rels")
	}
}

func TestListNotesPlaceholders(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notes-placeholders.pptx")
	_ = writeZipFixture(path, map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>
<Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>
<Override PartName="/ppt/notesSlides/notesSlide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.notesSlide+xml"/>
</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
</Relationships>`,
		"ppt/presentation.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:sldIdLst><p:sldId id="256" r:id="rId1"/></p:sldIdLst>
</p:presentation>`,
		"ppt/_rels/presentation.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/>
</Relationships>`,
		"ppt/slides/slide1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide" Target="../notesSlides/notesSlide1.xml"/>
</Relationships>`,
		"ppt/notesSlides/notesSlide1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:notes xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:cSld><p:spTree>
<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>
<p:sp><p:nvSpPr><p:cNvPr id="2" name="Body Placeholder"/><p:cNvSpPr/><p:nvPr><p:ph type="body" idx="1"/></p:nvPr></p:nvSpPr><p:spPr/></p:sp>
<p:sp><p:nvSpPr><p:cNvPr id="3" name="Title Placeholder"/><p:cNvSpPr/><p:nvPr><p:ph type="title" idx="0"/></p:nvPr></p:nvSpPr><p:spPr/></p:sp>
</p:spTree></p:cSld>
</p:notes>`,
		"ppt/notesSlides/_rels/notesSlide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="../slides/slide1.xml"/>
</Relationships>`,
	})

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	placeholders, err := editor.ListNotesPlaceholders(0)
	if err != nil {
		t.Fatalf("ListNotesPlaceholders failed: %v", err)
	}
	if len(placeholders) != 2 {
		t.Fatalf("expected two notes placeholders, got %d", len(placeholders))
	}
	var bodyFound, titleFound bool
	for _, ph := range placeholders {
		if ph.Type == "body" && ph.Index == 1 {
			bodyFound = true
		}
		if ph.Type == "title" && ph.Index == 0 {
			titleFound = true
		}
	}
	if !bodyFound || !titleFound {
		t.Fatalf("expected body idx=1 and title idx=0 placeholders, got %#v", placeholders)
	}
}

func writeTwoSlideNotesFixture(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "two-slide-notes.pptx")
	_ = writeZipFixture(path, map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>
<Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>
<Override PartName="/ppt/slides/slide2.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>
<Override PartName="/ppt/notesSlides/notesSlide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.notesSlide+xml"/>
<Override PartName="/ppt/notesMasters/notesMaster1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.notesMaster+xml"/>
<Override PartName="/ppt/theme/theme1.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>
</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
</Relationships>`,
		"ppt/presentation.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:sldMasterIdLst><p:sldMasterId id="2147483648" r:id="rId10"/></p:sldMasterIdLst>
<p:notesMasterIdLst><p:notesMasterId r:id="rId3"/></p:notesMasterIdLst>
<p:sldIdLst>
<p:sldId id="256" r:id="rId1"/>
<p:sldId id="257" r:id="rId2"/>
</p:sldIdLst>
</p:presentation>`,
		"ppt/_rels/presentation.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide2.xml"/>
<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesMaster" Target="notesMasters/notesMaster1.xml"/>
</Relationships>`,
		"ppt/slides/slide1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`,
		"ppt/slides/slide2.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide" Target="../notesSlides/notesSlide1.xml"/>
</Relationships>`,
		"ppt/slides/_rels/slide2.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
</Relationships>`,
		"ppt/notesSlides/notesSlide1.xml": fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:notes xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:sp><p:txBody><a:p><a:r><a:t>%s</a:t></a:r></a:p></p:txBody></p:sp></p:spTree></p:cSld></p:notes>`, "Original notes"),
		"ppt/notesSlides/_rels/notesSlide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="../slides/slide1.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesMaster" Target="../notesMasters/notesMaster1.xml"/>
</Relationships>`,
		"ppt/notesMasters/notesMaster1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:notesMaster xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:notesMaster>`,
		"ppt/notesMasters/_rels/notesMaster1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme2.xml"/>
</Relationships>`,
		"ppt/theme/theme2.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Office Theme"/>`,
		"ppt/theme/theme1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Office Theme"/>`,
	})
	return path
}
