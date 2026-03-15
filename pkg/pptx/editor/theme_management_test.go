package editor

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestPresentationEditorThemeInventory(t *testing.T) {
	path := writeThemeFixtureDeck(t)
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	inv, err := editor.ThemeInventory()
	if err != nil {
		t.Fatalf("theme inventory: %v", err)
	}
	if len(inv.ThemeParts) != 3 {
		t.Fatalf("expected 3 theme parts, got %d", len(inv.ThemeParts))
	}
	if len(inv.Bindings) != 3 {
		t.Fatalf("expected 3 bindings, got %d", len(inv.Bindings))
	}
	if !containsBinding(inv, "ppt/slideMasters/slideMaster1.xml", "ppt/theme/theme1.xml") {
		t.Fatalf("missing slideMaster1->theme1 binding: %#v", inv.Bindings)
	}
	if !containsBinding(inv, "ppt/slideMasters/slideMaster2.xml", "ppt/theme/theme3.xml") {
		t.Fatalf("missing slideMaster2->theme3 binding: %#v", inv.Bindings)
	}
	if !containsBinding(inv, "ppt/notesMasters/notesMaster1.xml", "ppt/theme/theme2.xml") {
		t.Fatalf("missing notesMaster1->theme2 binding: %#v", inv.Bindings)
	}
}

func TestPresentationEditorSetThemeData(t *testing.T) {
	path := writeThemeFixtureDeck(t)
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	const replacement = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Replaced"/>`
	if err := editor.SetThemeData("ppt/theme/theme2.xml", []byte(replacement)); err != nil {
		t.Fatalf("set theme data: %v", err)
	}
	got, ok := editor.parts.Get("ppt/theme/theme2.xml")
	if !ok {
		t.Fatalf("theme2 missing after replacement")
	}
	if !strings.Contains(string(got), `name="Replaced"`) {
		t.Fatalf("expected replacement theme, got %s", string(got))
	}
}

func TestPresentationEditorSetThemeFontAndColors(t *testing.T) {
	path := writeThemeFixtureDeck(t)
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	if err := editor.SetThemeFontScheme("Aptos Display", "Aptos"); err != nil {
		t.Fatalf("set theme font scheme: %v", err)
	}
	if err := editor.SetThemeColorScheme(ThemeColorScheme{
		Accent1: "FF1122",
		Accent2: "00AA55",
		Dk1:     "222222",
	}); err != nil {
		t.Fatalf("set theme color scheme: %v", err)
	}

	data, ok := editor.parts.Get("ppt/theme/theme1.xml")
	if !ok {
		t.Fatalf("theme1 missing")
	}
	xmlText := string(data)
	if !strings.Contains(xmlText, `typeface="Aptos Display"`) || !strings.Contains(xmlText, `typeface="Aptos"`) {
		t.Fatalf("expected updated major/minor fonts in theme1.xml: %s", xmlText)
	}
	if !strings.Contains(xmlText, `<a:accent1><a:srgbClr val="FF1122"/></a:accent1>`) {
		t.Fatalf("expected accent1 update in theme1.xml: %s", xmlText)
	}
	if !strings.Contains(xmlText, `<a:accent2><a:srgbClr val="00AA55"/></a:accent2>`) {
		t.Fatalf("expected accent2 update in theme1.xml: %s", xmlText)
	}
	if !strings.Contains(xmlText, `lastClr="222222"`) {
		t.Fatalf("expected dk1 lastClr update in theme1.xml: %s", xmlText)
	}
}

func TestResolveThemePreset(t *testing.T) {
	theme, ok := ResolveThemePreset("Office 2013")
	if !ok {
		t.Fatalf("expected office preset to resolve")
	}
	if theme.Name == "" {
		t.Fatalf("expected non-empty theme name")
	}
}

func writeThemeFixtureDeck(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "themes.pptx")
	themeXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Legacy">
<a:themeElements>
<a:clrScheme name="Legacy colors">
<a:dk1><a:sysClr val="windowText" lastClr="000000"/></a:dk1>
<a:lt1><a:sysClr val="window" lastClr="FFFFFF"/></a:lt1>
<a:dk2><a:srgbClr val="1F497D"/></a:dk2><a:lt2><a:srgbClr val="EEECE1"/></a:lt2>
<a:accent1><a:srgbClr val="4F81BD"/></a:accent1><a:accent2><a:srgbClr val="C0504D"/></a:accent2>
<a:accent3><a:srgbClr val="9BBB59"/></a:accent3><a:accent4><a:srgbClr val="8064A2"/></a:accent4>
<a:accent5><a:srgbClr val="4BACC6"/></a:accent5><a:accent6><a:srgbClr val="F79646"/></a:accent6>
<a:hlink><a:srgbClr val="0000FF"/></a:hlink><a:folHlink><a:srgbClr val="800080"/></a:folHlink>
</a:clrScheme>
<a:fontScheme name="Legacy fonts">
<a:majorFont><a:latin typeface="Calibri"/><a:ea typeface=""/><a:cs typeface=""/></a:majorFont>
<a:minorFont><a:latin typeface="Calibri"/><a:ea typeface=""/><a:cs typeface=""/></a:minorFont>
</a:fontScheme>
</a:themeElements>
</a:theme>`
	if err := writeZipFixture(path, map[string]string{
		"[Content_Types].xml":             `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/><Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/><Override PartName="/ppt/theme/theme1.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/><Override PartName="/ppt/theme/theme2.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/><Override PartName="/ppt/theme/theme3.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/></Types>`,
		"_rels/.rels":                     `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/></Relationships>`,
		"ppt/presentation.xml":            `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:sldIdLst><p:sldId id="256" r:id="rId1"/></p:sldIdLst><p:sldSz cx="9144000" cy="6858000" type="screen4x3"/><p:notesSz cx="6858000" cy="9144000"/></p:presentation>`,
		"ppt/_rels/presentation.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/></Relationships>`,
		"ppt/slides/slide1.xml":           `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`,
		"ppt/theme/theme1.xml":             themeXML,
		"ppt/theme/theme2.xml":             themeXML,
		"ppt/theme/theme3.xml":             themeXML,
		"ppt/slideMasters/_rels/slideMaster1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId7" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme1.xml"/></Relationships>`,
		"ppt/slideMasters/_rels/slideMaster2.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId7" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme3.xml"/></Relationships>`,
		"ppt/notesMasters/_rels/notesMaster1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme2.xml"/></Relationships>`,
	}); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

func containsBinding(inv ThemeInventory, owner, theme string) bool {
	for _, b := range inv.Bindings {
		if b.OwnerPart == owner && b.ThemePart == theme {
			return true
		}
	}
	return false
}
