package editor

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestPresentationEditorPreservesCustomXML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom-xml.pptx")

	err := writeZipFixture(path, map[string]string{
		"[Content_Types].xml":              `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/><Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/><Override PartName="/customXml/itemProps1.xml" ContentType="application/vnd.openxmlformats-officedocument.customXmlProperties+xml"/><Override PartName="/customXml/itemProps2.xml" ContentType="application/vnd.openxmlformats-officedocument.customXmlProperties+xml"/></Types>`,
		"_rels/.rels":                      `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXml" Target="customXml/item1.xml"/><Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXml" Target="customXml/item2.xml"/></Relationships>`,
		"ppt/presentation.xml":             `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:sldIdLst><p:sldId id="256" r:id="rId1"/></p:sldIdLst></p:presentation>`,
		"ppt/_rels/presentation.xml.rels":  `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/></Relationships>`,
		"ppt/slides/slide1.xml":            `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`,
		"customXml/item1.xml":              `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><RootElement xmlns="http://example.com/item1"><Prop1>Fish &amp; Chips &lt;ok&gt;</Prop1></RootElement>`,
		"customXml/itemProps1.xml":         `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><ds:datastoreItem ds:itemID="{B0AF63C1-F0B0-4054-94AE-ED6EF831A21A}" xmlns:ds="http://schemas.openxmlformats.org/officeDocument/2006/customXml"><ds:schemaRefs><ds:schemaRef ds:uri="http://example.com/item1"/></ds:schemaRefs></ds:datastoreItem>`,
		"customXml/item2.xml":              `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><RawXmlContent><data>123</data></RawXmlContent>`,
		"customXml/itemProps2.xml":         `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><ds:datastoreItem ds:itemID="{F63F6E2A-9485-4240-84E3-CEB3CDD9BDC8}" xmlns:ds="http://schemas.openxmlformats.org/officeDocument/2006/customXml"><ds:schemaRefs/></ds:datastoreItem>`,
	})
	if err != nil {
		t.Fatalf("failed to write mock pptx: %v", err)
	}

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	if len(editor.metadata.CustomXML) != 2 {
		t.Fatalf("expected 2 CustomXML parts, got %d", len(editor.metadata.CustomXML))
	}

	// Order could be inverted depending on map iteration from ps.Keys()
	// Let's identify by RootElement or Content
	var strucPart, rawPart any
	for _, p := range editor.metadata.CustomXML {
		if p.RootElement == "RootElement" {
			strucPart = p
		} else if strings.Contains(p.Content, "RawXmlContent") {
			rawPart = p
		}
	}

	if strucPart == nil || rawPart == nil {
		t.Fatalf("Failed to parse both custom xml types properly")
	}

	outPath := filepath.Join(t.TempDir(), "custom-xml-edited.pptx")
	if saveErr := editor.Save(outPath); saveErr != nil {
		t.Fatalf("save edited deck: %v", saveErr)
	}

	// Verify both survived
	item1 := string(readZipFileBytes(t, outPath, "customXml/item1.xml"))
	item2 := string(readZipFileBytes(t, outPath, "customXml/item2.xml"))

	if item1 == "" || item2 == "" {
		t.Fatalf("customXml/item1.xml and item2.xml were not written out")
	}

	if !strings.Contains(item1, "RootElement") && !strings.Contains(item2, "RootElement") {
		t.Fatalf("structured item lost on save")
	}
	if strings.Contains(item1, "RootElement") && !strings.Contains(item1, "Fish &amp; Chips &lt;ok&gt;") {
		t.Fatalf("structured item text was not XML-escaped on save: %s", item1)
	}
	if strings.Contains(item2, "RootElement") && !strings.Contains(item2, "Fish &amp; Chips &lt;ok&gt;") {
		t.Fatalf("structured item text was not XML-escaped on save: %s", item2)
	}
	if !strings.Contains(item1, "RawXmlContent") && !strings.Contains(item2, "RawXmlContent") {
		t.Fatalf("raw item lost on save")
	}
}
