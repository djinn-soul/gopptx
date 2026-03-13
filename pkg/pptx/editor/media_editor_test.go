package editor

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestAddVideoEmbedsMediaAndPosterWithMimeHandling(t *testing.T) {
	e := newMediaEditorFixture()
	shapeID, err := e.AddVideo(
		0,
		[]byte("video-bytes"),
		testutil.TinyPNG(),
		"video/quicktime",
		10,
		20,
		300,
		200,
	)
	if err != nil {
		t.Fatalf("AddVideo failed: %v", err)
	}
	if shapeID == 0 {
		t.Fatal("expected non-zero shape id")
	}

	mediaParts := e.parts.KeysWithPrefix("ppt/media/")
	hasPoster := false
	hasMov := false
	for _, part := range mediaParts {
		if strings.HasSuffix(part, ".png") {
			hasPoster = true
		}
		if strings.HasSuffix(part, ".mov") {
			hasMov = true
		}
	}
	if !hasPoster || !hasMov {
		t.Fatalf("expected poster png and mov media parts, got: %v", mediaParts)
	}

	slideXML, ok := e.parts.Get("ppt/slides/slide1.xml")
	if !ok {
		t.Fatal("slide xml not found")
	}
	slideText := string(slideXML)
	if !strings.Contains(slideText, `<a:videoFile r:link="`) {
		t.Fatalf("expected video file link in slide xml: %s", slideText)
	}
	if !strings.Contains(slideText, `<p14:media`) {
		t.Fatalf("expected modern media extension in slide xml: %s", slideText)
	}
	if !strings.Contains(slideText, `<a:off x="10" y="20"/>`) ||
		!strings.Contains(slideText, `<a:ext cx="300" cy="200"/>`) {
		t.Fatalf("expected deterministic sizing in slide xml: %s", slideText)
	}

	relsXML, ok := e.parts.Get("ppt/slides/_rels/slide1.xml.rels")
	if !ok {
		t.Fatal("slide rels xml not found")
	}
	relsText := string(relsXML)
	if !strings.Contains(relsText, common.RelTypeVideo) || !strings.Contains(relsText, common.RelTypeMedia) {
		t.Fatalf("expected video+media relationships in rels xml: %s", relsText)
	}
}

func TestAddOLEObjectEmbedsPackageAndEscapesProgID(t *testing.T) {
	e := newMediaEditorFixture()
	shapeID, err := e.AddOLEObject(
		0,
		[]byte("ole-object"),
		testutil.TinyPNG(),
		`Excel.Sheet.12 & <unsafe>`,
		100,
		110,
		400,
		250,
	)
	if err != nil {
		t.Fatalf("AddOLEObject failed: %v", err)
	}
	if shapeID == 0 {
		t.Fatal("expected non-zero shape id")
	}

	embeddingParts := e.parts.KeysWithPrefix("ppt/embeddings/")
	if len(embeddingParts) != 1 || !strings.HasSuffix(embeddingParts[0], ".bin") {
		t.Fatalf("expected one .bin embedding part, got: %v", embeddingParts)
	}

	slideXML, ok := e.parts.Get("ppt/slides/slide1.xml")
	if !ok {
		t.Fatal("slide xml not found")
	}
	slideText := string(slideXML)
	if !strings.Contains(slideText, `progId="Excel.Sheet.12 &amp; &lt;unsafe&gt;"`) {
		t.Fatalf("expected escaped progId in slide xml: %s", slideText)
	}
	if !strings.Contains(slideText, `<a:off x="100" y="110"/>`) ||
		!strings.Contains(slideText, `<a:ext cx="400" cy="250"/>`) {
		t.Fatalf("expected deterministic ole bounds in slide xml: %s", slideText)
	}

	relsXML, ok := e.parts.Get("ppt/slides/_rels/slide1.xml.rels")
	if !ok {
		t.Fatal("slide rels xml not found")
	}
	relsText := string(relsXML)
	if !strings.Contains(relsText, common.RelTypePackage) {
		t.Fatalf("expected package relationship in rels xml: %s", relsText)
	}
	if !strings.Contains(relsText, "../embeddings/") {
		t.Fatalf("expected embedding target in rels xml: %s", relsText)
	}
}

func TestAddAudioEmbedsModernMediaAndAudioRelationships(t *testing.T) {
	e := newMediaEditorFixture()
	shapeID, err := e.AddAudio(
		0,
		[]byte("audio-bytes"),
		"audio/mp4",
		12,
		24,
		220,
		90,
	)
	if err != nil {
		t.Fatalf("AddAudio failed: %v", err)
	}
	if shapeID == 0 {
		t.Fatal("expected non-zero shape id")
	}

	mediaParts := e.parts.KeysWithPrefix("ppt/media/")
	hasM4A := false
	for _, part := range mediaParts {
		if strings.HasSuffix(part, ".m4a") {
			hasM4A = true
			break
		}
	}
	if !hasM4A {
		t.Fatalf("expected m4a media part, got: %v", mediaParts)
	}

	slideXML := string(getFixturePart(t, e, "ppt/slides/slide1.xml"))
	if !strings.Contains(slideXML, `<a:audioFile r:link="`) || !strings.Contains(slideXML, `<p14:media`) {
		t.Fatalf("expected audio+modern media xml in slide: %s", slideXML)
	}

	relsXML := string(getFixturePart(t, e, "ppt/slides/_rels/slide1.xml.rels"))
	if !strings.Contains(relsXML, common.RelTypeAudio) || !strings.Contains(relsXML, common.RelTypeMedia) {
		t.Fatalf("expected audio/media relationships in rels xml: %s", relsXML)
	}
}

func TestAddVideoRejectsUnsupportedMimeType(t *testing.T) {
	e := newMediaEditorFixture()
	_, err := e.AddVideo(0, []byte("video-bytes"), testutil.TinyPNG(), "video/unsupported", 1, 2, 3, 4)
	if err == nil || !strings.Contains(err.Error(), "unsupported video mime type") {
		t.Fatalf("expected unsupported video mime validation error, got: %v", err)
	}
}

func TestAddAudioRejectsUnsupportedMimeType(t *testing.T) {
	e := newMediaEditorFixture()
	_, err := e.AddAudio(0, []byte("audio-bytes"), "audio/unsupported", 1, 2, 3, 4)
	if err == nil || !strings.Contains(err.Error(), "unsupported audio mime type") {
		t.Fatalf("expected unsupported audio mime validation error, got: %v", err)
	}
}

func newMediaEditorFixture() *PresentationEditor {
	ps := NewPartStore()
	ps.Set(
		"ppt/slides/slide1.xml",
		[]byte(
			`<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">`+
				`<p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`,
		),
	)
	ps.Set(
		"ppt/slides/_rels/slide1.xml.rels",
		[]byte(
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`+
				`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`+
				`<Relationship Id="rId1" Type="`+common.RelTypeSlideLayout+`" Target="../slideLayouts/slideLayout1.xml"/>`+
				`</Relationships>`,
		),
	)
	return &PresentationEditor{
		parts: ps,
		slides: []common.EditorSlideRef{{
			SlideID: 256,
			Part:    "ppt/slides/slide1.xml",
		}},
		mediaInventory: map[string]string{},
		nextMediaNum:   1,
	}
}

func getFixturePart(t *testing.T, e *PresentationEditor, part string) []byte {
	t.Helper()
	value, ok := e.parts.Get(part)
	if !ok {
		t.Fatalf("missing part %s", part)
	}
	return value
}

func TestMediaAndOLESaveWiresPackagingAndContentTypes(t *testing.T) {
	base := writeDeckFixture(t, "media-ole-base.pptx", []elements.SlideContent{elements.NewSlide("Slide 1")})
	editor, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	if _, err := editor.AddVideo(
		0,
		[]byte("video-bytes"),
		testutil.TinyPNG(),
		"video/mp4",
		10,
		10,
		200,
		120,
	); err != nil {
		t.Fatalf("AddVideo failed: %v", err)
	}
	if _, err := editor.AddOLEObject(
		0,
		[]byte("ole-object"),
		testutil.TinyPNG(),
		"Excel.Sheet.12",
		50,
		70,
		180,
		110,
	); err != nil {
		t.Fatalf("AddOLEObject failed: %v", err)
	}

	out := filepath.Join(t.TempDir(), "media-ole-out.pptx")
	if err := editor.Save(out); err != nil {
		t.Fatalf("save output: %v", err)
	}

	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		t.Fatalf("open output zip: %v", err)
	}

	hasMP4 := false
	hasOLE := false
	for _, f := range zr.File {
		if strings.HasPrefix(f.Name, "ppt/media/media") && strings.HasSuffix(f.Name, ".mp4") {
			hasMP4 = true
		}
		if strings.HasPrefix(f.Name, "ppt/embeddings/oleObject") && strings.HasSuffix(f.Name, ".bin") {
			hasOLE = true
		}
	}
	if !hasMP4 || !hasOLE {
		t.Fatalf("expected mp4 media and bin embedding parts, hasMP4=%v hasOLE=%v", hasMP4, hasOLE)
	}

	relsXML := testutil.ReadZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(relsXML, common.RelTypeVideo) ||
		!strings.Contains(relsXML, common.RelTypeMedia) ||
		!strings.Contains(relsXML, common.RelTypePackage) {
		t.Fatalf("expected video/media/package relationships in slide rels: %s", relsXML)
	}

	contentTypes := testutil.ReadZipFile(t, zr, "[Content_Types].xml")
	if !strings.Contains(contentTypes, `Extension="mp4"`) {
		t.Fatalf("expected mp4 content-type default in [Content_Types].xml: %s", contentTypes)
	}
	if !strings.Contains(contentTypes, `Extension="bin"`) {
		t.Fatalf("expected bin content-type default in [Content_Types].xml: %s", contentTypes)
	}
}
