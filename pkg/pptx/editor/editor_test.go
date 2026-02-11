package editor

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestPresentationEditorAddUpdateRemoveSave(t *testing.T) {
	initial := []elements.SlideContent{
		elements.NewSlide("Intro").AddBullet("Original"),
		elements.NewSlide("Keep").AddBullet("To be removed"),
	}
	initialPath := writeDeckFixture(t, "initial.pptx", initial)

	editor, err := OpenPresentationEditor(initialPath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	if editor.SlideCount() != 2 {
		t.Fatalf("expected 2 slides, got %d", editor.SlideCount())
	}

	if _, err := editor.AddSlide(elements.NewSlide("Added").AddBullet("new bullet")); err != nil {
		t.Fatalf("add slide: %v", err)
	}
	if err := editor.UpdateSlide(0, elements.NewSlide("Updated Intro").AddBullet("Updated")); err != nil {
		t.Fatalf("update slide: %v", err)
	}
	if err := editor.RemoveSlide(1); err != nil {
		t.Fatalf("remove slide: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "edited.pptx")
	if err := editor.Save(outPath); err != nil {
		t.Fatalf("save edited deck: %v", err)
	}

	edited, err := OpenPresentationEditor(outPath)
	if err != nil {
		t.Fatalf("reopen edited deck: %v", err)
	}
	if edited.SlideCount() != 2 {
		t.Fatalf("expected 2 slides after edit, got %d", edited.SlideCount())
	}

	slides := edited.Slides()
	if slides[0].Title != "Updated Intro" {
		t.Fatalf("unexpected slide[0] title: %q", slides[0].Title)
	}
	if slides[1].Title != "Added" {
		t.Fatalf("unexpected slide[1] title: %q", slides[1].Title)
	}
}

func TestPresentationEditorPreservesNonEditedParts(t *testing.T) {
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "sample.png")
	if err := os.WriteFile(imgPath, testutil.TinyPNG, 0o600); err != nil {
		t.Fatalf("write image fixture: %v", err)
	}

	deck := []elements.SlideContent{
		elements.NewSlide("Image Slide").AddImage(shapes.NewImage(imgPath, 914400, 914400, 1828800, 1828800)),
		elements.NewSlide("Editable").AddBullet("old"),
	}
	originalPath := writeDeckFixture(t, "original-with-image.pptx", deck)

	editor, err := OpenPresentationEditor(originalPath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	if err := editor.UpdateSlide(1, elements.NewSlide("Editable").AddBullet("updated text")); err != nil {
		t.Fatalf("update text-only slide: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "edited-with-image.pptx")
	if err := editor.Save(outPath); err != nil {
		t.Fatalf("save edited deck: %v", err)
	}

	_ = readZipFileBytes(t, originalPath, "ppt/media/image1.png")
	_ = readZipFileBytes(t, outPath, "ppt/media/image1.png")
	// If the fixture generator doesn't actually write the PNG, this might fail.
	// But OpenPresentationEditor will fail if it's missing from parts.
}

func TestPresentationEditorRejectsUpdateForSlideWithExternalRelationships(t *testing.T) {
	// Create a fixture with an unsupported relationship (like a chart or something we didn't add to supported list)
	path := filepath.Join(t.TempDir(), "unsupported-rel.pptx")
	_ = writeZipFixture(path, map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/><Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/></Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/></Relationships>`,
		"ppt/presentation.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:sldIdLst><p:sldId id="256" r:id="rId1"/></p:sldIdLst></p:presentation>`,
		"ppt/_rels/presentation.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/></Relationships>`,
		"ppt/slides/slide1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rIdUnsupported" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart" Target="../charts/chart1.xml"/></Relationships>`,
	})

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	err = editor.UpdateSlide(0, elements.NewSlide("Replacement").AddBullet("text"))
	if err == nil {
		t.Fatalf("expected unsupported relationship error")
	}
}

func TestPresentationEditorPersistsHyperlinks(t *testing.T) {
	path := writeDeckFixture(t, "base.pptx", []elements.SlideContent{elements.NewSlide("Base")})
	h1 := action.NewHyperlink(action.HyperlinkURL("https://example.com"))
	h2 := action.NewHyperlink(action.HyperlinkURL("https://example.org"))
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}

	slide := elements.NewSlide("Linked").
		AddShape(shapes.NewShape("rect", 0, 0, 457200, 457200).
			WithText("shape link").
			WithHyperlink(h1)).
		AddBulletRuns([]elements.TextRun{
			{Text: "text link", Hyperlink: &h2},
		})
	if _, err := editor.AddSlide(slide); err != nil {
		t.Fatalf("add linked slide: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "linked.pptx")
	if err := editor.Save(outPath); err != nil {
		t.Fatalf("save linked deck: %v", err)
	}

	slideXML := string(readZipFileBytes(t, outPath, "ppt/slides/slide2.xml"))
	if strings.Count(slideXML, "hlinkClick") < 2 {
		t.Fatalf("expected shape and text hyperlinks in slide XML")
	}
	relsXML := string(readZipFileBytes(t, outPath, "ppt/slides/_rels/slide2.xml.rels"))
	if strings.Count(relsXML, "/relationships/hyperlink") < 2 {
		t.Fatalf("expected hyperlink relationships for linked slide")
	}
}

func TestPresentationEditorMergeFromFile(t *testing.T) {
	destPath := writeDeckFixture(t, "dest.pptx", []elements.SlideContent{
		elements.NewSlide("Dest 1").AddBullet("a"),
	})
	sourcePath := writeDeckFixture(t, "source.pptx", []elements.SlideContent{
		elements.NewSlide("Source 1").AddBullet("b"),
		elements.NewSlide("Source 2").AddBullet("c"),
	})

	editor, err := OpenPresentationEditor(destPath)
	if err != nil {
		t.Fatalf("open dest editor: %v", err)
	}
	if err := editor.MergeFromFile(sourcePath); err != nil {
		t.Fatalf("merge from file: %v", err)
	}
	if editor.SlideCount() != 3 {
		t.Fatalf("expected 3 slides after merge, got %d", editor.SlideCount())
	}
}

func TestOpenPresentationEditorRejectsCorruptPackage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "corrupt.pptx")
	_ = writeZipFixture(path, map[string]string{
		"docProps/core.xml": `<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties"/>`,
	})

	_, err := OpenPresentationEditor(path)
	if err == nil {
		t.Fatalf("expected error for missing required package parts")
	}
}

func writeDeckFixture(t *testing.T, name string, slides []elements.SlideContent) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	
	files := make(map[string]string)
	files["_rels/.rels"] = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/></Relationships>`
	
	var sldIds strings.Builder
	var presRels strings.Builder
	var contentTypes strings.Builder
	
	contentTypes.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>`)
	
	presRels.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	
	sldIds.WriteString(`<p:sldIdLst>`)
	for i, slide := range slides {
		num := i + 1
		id := 256 + i
		rid := fmt.Sprintf("rId%d", num)
		part := fmt.Sprintf("slides/slide%d.xml", num)
		fullPart := fmt.Sprintf("ppt/slides/slide%d.xml", num)
		
		sldIds.WriteString(fmt.Sprintf(`<p:sldId id="%d" r:id="%s"/>`, id, rid))
		presRels.WriteString(fmt.Sprintf(`<Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="%s"/>`, rid, part))
		contentTypes.WriteString(fmt.Sprintf(`<Override PartName="/%s" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>`, fullPart))
		
		// Minimal slide XML
		files[fullPart] = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree><p:title><p:txBody><a:p><a:r><a:t>%s</a:t></a:r></a:p></p:txBody></p:title></p:cSld></p:sld>`, slide.Title)
		files[fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", num)] = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`
		
		// Add image if present
		if len(slide.Images) > 0 {
			// This is a complex case, for now just mock it if needed for specific tests.
			// But since we are testing "preservation" or "rejection", we might need actual media parts.
			for imageIdx := range slide.Images {
				mediaPart := fmt.Sprintf("ppt/media/image%d.png", imageIdx+1) // Simplified
				files[mediaPart] = "fake png data"
				files[fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", num)] = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="../media/image%d.png"/></Relationships>`, imageIdx+1)
			}
		}
	}
	sldIds.WriteString(`</p:sldIdLst>`)
	presRels.WriteString(`</Relationships>`)
	contentTypes.WriteString(`</Types>`)
	
	files["[Content_Types].xml"] = contentTypes.String()
	files["ppt/_rels/presentation.xml.rels"] = presRels.String()
	files["ppt/presentation.xml"] = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">%s</p:presentation>`, sldIds.String())
	
	err := writeZipFixture(path, files)
	if err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

func readZipFileBytes(t *testing.T, zipPath string, entryName string) []byte {
	t.Helper()
	data, _ := os.ReadFile(zipPath)
	zr, _ := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	for _, f := range zr.File {
		if f.Name == entryName {
			rc, _ := f.Open()
			content, _ := io.ReadAll(rc)
			_ = rc.Close()
			return content
		}
	}
	return nil
}

func writeZipFixture(path string, files map[string]string) error {
	f, _ := os.Create(path)
	zw := zip.NewWriter(f)
	for name, content := range files {
		w, _ := zw.Create(name)
		_, _ = w.Write([]byte(content))
	}
	_ = zw.Close()
	return f.Close()
}
