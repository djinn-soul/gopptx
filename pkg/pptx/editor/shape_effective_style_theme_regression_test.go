package editor

import (
	"strings"
	"testing"
)

func TestEffectiveStyleUsesCurrentMastersBoundTheme(t *testing.T) {
	fixture := writeInheritanceFixture(t, "effective-style-bound-theme.pptx")
	editor, err := OpenPresentationEditor(fixture)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	firstTheme := strings.Replace(inheritanceThemeXML, "4472C4", "FF0000", 1)
	boundTheme := strings.Replace(inheritanceThemeXML, "4472C4", "00AA55", 1)
	editor.parts.Set("ppt/theme/theme0.xml", []byte(firstTheme))
	editor.parts.Set("ppt/theme/theme2.xml", []byte(boundTheme))

	relsPath := "ppt/slideMasters/_rels/slideMaster1.xml.rels"
	rels, _ := editor.parts.Get(relsPath)
	editor.parts.Set(relsPath, []byte(strings.Replace(
		string(rels),
		`Target="../theme/theme1.xml"`,
		`Target="../theme/theme2.xml"`,
		1,
	)))

	style, err := editor.GetEffectiveShapeStyle(0, 2)
	if err != nil {
		t.Fatalf("effective style: %v", err)
	}
	if style.FontColor == nil || style.FontColor.RGB != "00AA55" {
		t.Fatalf("expected master-bound theme accent, got %+v", style.FontColor)
	}
}

func TestEffectiveStyleAppliesSchemeColorLuminanceTransform(t *testing.T) {
	fixture := writeInheritanceFixture(t, "effective-style-transform.pptx")
	editor, err := OpenPresentationEditor(fixture)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	layoutPart := "ppt/slideLayouts/slideLayout1.xml"
	layout, _ := editor.parts.Get(layoutPart)
	editor.parts.Set(layoutPart, []byte(strings.Replace(
		string(layout),
		`<a:schemeClr val="accent1"/>`,
		`<a:schemeClr val="accent1"><a:lumMod val="50000"/></a:schemeClr>`,
		1,
	)))

	style, err := editor.GetEffectiveShapeStyle(0, 2)
	if err != nil {
		t.Fatalf("effective style: %v", err)
	}
	if style.FontColor == nil {
		t.Fatal("expected transformed font color")
	}
	if style.FontColor.SchemeSlot != "accent1" {
		t.Fatalf("scheme slot = %q", style.FontColor.SchemeSlot)
	}
	if style.FontColor.RGB != "203864" {
		t.Fatalf("transformed RGB = %q, want 203864", style.FontColor.RGB)
	}
}
