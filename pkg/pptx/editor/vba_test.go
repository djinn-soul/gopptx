package editor

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestPresentationEditorPreservesVBA(t *testing.T) {
	path := filepath.Join(t.TempDir(), "vba-macro.pptm")

	// Synthetic minimal vbaProject.bin (just some bytes for the test)
	mockVBAData := []byte("fake_vba_bin_content")

	err := writeZipFixture(path, map[string]string{
		"[Content_Types].xml":              `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="bin" ContentType="application/vnd.ms-office.vbaProject"/><Override PartName="/ppt/presentation.xml" ContentType="application/vnd.ms-powerpoint.presentation.macroEnabled.main+xml"/><Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/></Types>`,
		"_rels/.rels":                      `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/></Relationships>`,
		"ppt/presentation.xml":             `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:sldIdLst><p:sldId id="256" r:id="rId1"/></p:sldIdLst></p:presentation>`,
		"ppt/_rels/presentation.xml.rels":  `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/><Relationship Id="rId2" Type="http://schemas.microsoft.com/office/2006/relationships/vbaProject" Target="vbaProject.bin"/></Relationships>`,
		"ppt/slides/slide1.xml":            `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`,
		"ppt/vbaProject.bin":               string(mockVBAData),
	})
	if err != nil {
		t.Fatalf("failed to write mock .pptm: %v", err)
	}

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	// Verify VBA project was extracted into memory
	if editor.metadata.VBA == nil {
		t.Fatalf("VBA project was not extracted into editor metadata")
	}

	outPath := filepath.Join(t.TempDir(), "vba-macro-edited.pptm")
	if saveErr := editor.Save(outPath); saveErr != nil {
		t.Fatalf("save edited deck: %v", saveErr)
	}

	// Verify VBA blob survived the round-trip
	actualVBAData := readZipFileBytes(t, outPath, "ppt/vbaProject.bin")
	if !bytes.Equal(actualVBAData, mockVBAData) {
		t.Fatalf("expected vbaProject.bin to survive round trip exactly")
	}

	// Double check the manifest was correctly rewritten
	contentTypes := string(readZipFileBytes(t, outPath, "[Content_Types].xml"))
	if !strings.Contains(contentTypes, "application/vnd.ms-office.vbaProject") {
		t.Fatalf("expected ContentTypes rewrite to preserve vba properties")
	}
	if !strings.Contains(contentTypes, "application/vnd.ms-powerpoint.presentation.macroEnabled.main+xml") {
		t.Fatalf("expected macro-enabled presentation content type in [Content_Types].xml")
	}
	if strings.Count(contentTypes, `/ppt/vbaProject.bin`) != 1 {
		t.Fatalf("expected exactly one /ppt/vbaProject.bin override, got content types: %s", contentTypes)
	}

	presentationRels := string(readZipFileBytes(t, outPath, "ppt/_rels/presentation.xml.rels"))
	if !strings.Contains(presentationRels, "http://schemas.microsoft.com/office/2006/relationships/vbaProject") {
		t.Fatalf("expected vbaProject relationship in presentation rels")
	}
	if !strings.Contains(presentationRels, `Target="vbaProject.bin"`) {
		t.Fatalf("expected presentation rel target to vbaProject.bin")
	}
}

func TestPresentationEditorSaveRejectsVBAToPptxExtension(t *testing.T) {
	path := filepath.Join(t.TempDir(), "vba-macro-source.pptm")
	err := writeZipFixture(path, map[string]string{
		"[Content_Types].xml":              `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="bin" ContentType="application/vnd.ms-office.vbaProject"/><Override PartName="/ppt/presentation.xml" ContentType="application/vnd.ms-powerpoint.presentation.macroEnabled.main+xml"/><Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/><Override PartName="/ppt/vbaProject.bin" ContentType="application/vnd.ms-office.vbaProject"/></Types>`,
		"_rels/.rels":                      `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/></Relationships>`,
		"ppt/presentation.xml":             `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:sldIdLst><p:sldId id="256" r:id="rId1"/></p:sldIdLst></p:presentation>`,
		"ppt/_rels/presentation.xml.rels":  `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/><Relationship Id="rId2" Type="http://schemas.microsoft.com/office/2006/relationships/vbaProject" Target="vbaProject.bin"/></Relationships>`,
		"ppt/slides/slide1.xml":            `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`,
		"ppt/vbaProject.bin":               "fake_vba_bin_content",
	})
	if err != nil {
		t.Fatalf("failed to write mock .pptm: %v", err)
	}

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	outPath := filepath.Join(t.TempDir(), "vba-macro-edited.pptx")
	saveErr := editor.Save(outPath)
	if saveErr == nil {
		t.Fatal("expected save to .pptx with VBA metadata to fail")
	}
	if !strings.Contains(saveErr.Error(), ".pptm extension") {
		t.Fatalf("expected .pptm extension guidance, got: %v", saveErr)
	}
}
