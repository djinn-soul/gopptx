package editor

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestPresentationEditorApplyThemeAndSlideSize(t *testing.T) {
	path := filepath.Join(t.TempDir(), "with-theme-and-size.pptx")
	if err := writeZipFixture(path, map[string]string{
		"[Content_Types].xml":              `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/><Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/><Override PartName="/ppt/theme/theme1.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/></Types>`,
		"_rels/.rels":                      `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/></Relationships>`,
		"ppt/presentation.xml":             `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:sldIdLst><p:sldId id="256" r:id="rId1"/></p:sldIdLst><p:sldSz cx="9144000" cy="6858000" type="screen4x3"/><p:notesSz cx="6858000" cy="9144000"/></p:presentation>`,
		"ppt/_rels/presentation.xml.rels":  `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="theme/theme1.xml"/></Relationships>`,
		"ppt/slides/slide1.xml":            `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`,
		"ppt/theme/theme1.xml":             `<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Legacy"><a:themeElements/></a:theme>`,
	}); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()
	if applyErr := editor.ApplyTheme(styling.ThemeTech); applyErr != nil {
		t.Fatalf("apply theme: %v", applyErr)
	}
	if sizeErr := editor.SetSlideSize(SlideSize16x9()); sizeErr != nil {
		t.Fatalf("set slide size: %v", sizeErr)
	}

	outPath := filepath.Join(t.TempDir(), "edited-theme-size.pptx")
	if saveErr := editor.Save(outPath); saveErr != nil {
		t.Fatalf("save edited deck: %v", saveErr)
	}

	themeXML := string(readZipFileBytes(t, outPath, "ppt/theme/theme1.xml"))
	if !strings.Contains(themeXML, `name="Tech marketing"`) {
		t.Fatalf("expected updated theme name, got: %s", themeXML)
	}

	presentationXML := string(readZipFileBytes(t, outPath, "ppt/presentation.xml"))
	slideSize16x9 := SlideSize16x9()
	expectedSz := fmt.Sprintf(`<p:sldSz cx="%d" cy="%d" type="screen16x9"/>`, slideSize16x9.Width, slideSize16x9.Height)
	if !strings.Contains(presentationXML, expectedSz) {
		t.Fatalf("expected updated slide size %q in presentation.xml", expectedSz)
	}
}

func TestPresentationEditorSetSlideSizeInsertsWhenMissing(t *testing.T) {
	basePath := writeDeckFixture(t, "missing-size.pptx", nil)

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()
	if sizeErr := editor.SetSlideSize(SlideSize4x3()); sizeErr != nil {
		t.Fatalf("set slide size: %v", sizeErr)
	}

	outPath := filepath.Join(t.TempDir(), "inserted-size.pptx")
	if saveErr := editor.Save(outPath); saveErr != nil {
		t.Fatalf("save edited deck: %v", saveErr)
	}

	presentationXML := string(readZipFileBytes(t, outPath, "ppt/presentation.xml"))
	slideSize4x3 := SlideSize4x3()
	expectedSz := fmt.Sprintf(`<p:sldSz cx="%d" cy="%d" type="screen4x3"/>`, slideSize4x3.Width, slideSize4x3.Height)
	if !strings.Contains(presentationXML, expectedSz) {
		t.Fatalf("expected inserted slide size %q in presentation.xml", expectedSz)
	}
}

func TestPresentationEditorApplyThemeRequiresThemePart(t *testing.T) {
	basePath := writeDeckFixture(t, "missing-theme.pptx", []elements.SlideContent{elements.NewSlide("Only")})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()
	if applyErr := editor.ApplyTheme(styling.ThemeCorporate); applyErr == nil {
		t.Fatalf("expected missing theme part error")
	}
}

func TestPresentationEditorCoreProperties(t *testing.T) {
	path := filepath.Join(t.TempDir(), "props.pptx")
	coreXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:dcmitype="http://purl.org/dc/dcmitype/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<dc:title>Original Title</dc:title>
<dc:subject>Original Subject</dc:subject>
<dc:creator>Original Creator</dc:creator>
</cp:coreProperties>`

	if err := writeZipFixture(path, map[string]string{
		"[Content_Types].xml":             `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/><Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/></Types>`,
		"_rels/.rels":                     `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/></Relationships>`,
		"ppt/presentation.xml":            `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:presentation xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:sldIdLst/></p:presentation>`,
		"ppt/_rels/presentation.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`,
		"docProps/core.xml":               coreXML,
	}); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	// 1. Open and verify read
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	props := editor.GetCoreProperties()
	if props.Title != "Original Title" {
		t.Errorf("expected Title 'Original Title', got %q", props.Title)
	}
	if props.Creator != "Original Creator" {
		t.Errorf("expected Creator 'Original Creator', got %q", props.Creator)
	}

	// 2. Update properties
	newProps := common.CoreProperties{
		Title:       "Updated Title",
		Subject:     "Updated Subject",
		Creator:     "Updated Creator",
		Description: "New Description",
	}
	editor.SetCoreProperties(newProps)

	// 3. Save
	outPath := filepath.Join(t.TempDir(), "edited-props.pptx")
	if saveErr := editor.Save(outPath); saveErr != nil {
		t.Fatalf("save edited deck: %v", saveErr)
	}

	// 4. Verify output XML
	savedCoreXML := readZipFileBytes(t, outPath, "docProps/core.xml")
	savedProps, err := parseCoreProperties(savedCoreXML)
	if err != nil {
		t.Fatalf("parse saved core properties: %v", err)
	}

	if savedProps.Title != "Updated Title" {
		t.Errorf("expected Title 'Updated Title', got %q", savedProps.Title)
	}
	if savedProps.Description != "New Description" {
		t.Errorf("expected Description 'New Description', got %q", savedProps.Description)
	}
	if savedProps.Creator != "Updated Creator" {
		t.Errorf("expected Creator 'Updated Creator', got %q", savedProps.Creator)
	}
}
