package pptxxml

import (
	"strings"
	"testing"
)

func TestSignatureOriginXML(t *testing.T) {
	xml := SignatureOrigin()

	if !strings.Contains(
		xml,
		`<SignatureOrigin xmlns="http://schemas.openxmlformats.org/package/2006/digital-signature"/>`,
	) {
		t.Fatalf("unexpected signature origin xml: %s", xml)
	}
	if strings.Contains(xml, "<vnd.openxmlformats-package.digital-signature-origin") {
		t.Fatalf("signature origin must use SignatureOrigin element: %s", xml)
	}
}

func TestContentTypesCustomPropertiesOverride(t *testing.T) {
	withCustom := ContentTypes(1, nil, 0, 0, nil, false, 0, 1, 0, false, nil, true, false, false, false, false)
	if !strings.Contains(
		withCustom,
		`<Override PartName="/docProps/custom.xml" ContentType="application/vnd.openxmlformats-officedocument.custom-properties+xml"/>`,
	) {
		t.Fatalf("missing custom properties override in content types: %s", withCustom)
	}

	withoutCustom := ContentTypes(1, nil, 0, 0, nil, false, 0, 1, 0, false, nil, false, false, false, false, false)
	if strings.Contains(withoutCustom, `/docProps/custom.xml`) {
		t.Fatalf("unexpected custom properties override without custom props enabled: %s", withoutCustom)
	}
}

func TestRootRelationships(t *testing.T) {
	tests := []struct {
		name          string
		hasCustom     bool
		hasSignatures bool
		contains      []string
	}{
		{
			name:      "basic",
			hasCustom: false, hasSignatures: false,
			contains: []string{"rId1", "rId2", "rId3"},
		},
		{
			name:      "with custom and signatures",
			hasCustom: true, hasSignatures: true,
			contains: []string{"rId4", "rId5", "custom-properties", "digital-signature"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RootRelationships(tt.hasCustom, tt.hasSignatures)
			for _, s := range tt.contains {
				if !strings.Contains(got, s) {
					t.Errorf("RootRelationships() missing %v", s)
				}
			}
		})
	}
}

func TestPresentationRelationships(t *testing.T) {
	got := PresentationRelationships(2, true, 1, 1, true, true, true, true, 1)
	expected := []string{
		"slideMasters/slideMaster1.xml",
		"slides/slide1.xml",
		"slides/slide2.xml",
		"notesMasters/notesMaster1.xml",
		"customXml/item1.xml",
		"sectionList.xml",
		"commentAuthors.xml",
		"vbaProject.bin",
		"handoutMasters/handoutMaster1.xml",
		"fonts/font1.fntdata",
	}

	for _, s := range expected {
		if !strings.Contains(got, s) {
			t.Errorf("PresentationRelationships() missing %v", s)
		}
	}
}

func TestSectionListXML(t *testing.T) {
	sections := []Section{
		{
			Name:     "Section 1",
			GUID:     "{GUID1}",
			SlideIDs: []int64{256, 257},
		},
	}
	got := SectionListXML(sections)
	expected := []string{
		"Section 1",
		"{GUID1}",
		"sldId id=\"256\"",
		"sldId id=\"257\"",
	}

	for _, s := range expected {
		if !strings.Contains(got, s) {
			t.Errorf("SectionListXML() missing %v", s)
		}
	}
}

func TestPresentation(t *testing.T) {
	got := Presentation("Title", 1, true, 9144000, 6858000, 1, &ProtectionInfo{SpinCount: 100000}, nil, true, nil)
	expected := []string{
		"rtl=\"1\"",
		"sldMasterId id=\"2147483648\"",
		"notesMasterId",
		"sldId id=\"257\"",
		"screen4x3",
		"modifyVerifier",
		"spinCount=\"100000\"",
	}

	for _, s := range expected {
		if !strings.Contains(got, s) {
			t.Errorf("Presentation() missing %v", s)
		}
	}
}

func TestEscape(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"&", "&amp;"},
		{"<", "&lt;"},
		{">", "&gt;"},
		{"\"", "&quot;"},
		{"'", "&apos;"},
		{"hello & world", "hello &amp; world"},
	}

	for _, tt := range tests {
		if got := Escape(tt.input); got != tt.want {
			t.Errorf("Escape(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestEmbeddedFontsXML(t *testing.T) {
	fonts := []EmbeddedFontRef{
		{
			Typeface: "Arial",
			Style:    "regular",
			RelID:    "rId1",
		},
		{
			Typeface: "Arial",
			Style:    "bold",
			RelID:    "rId2",
		},
	}
	got := EmbeddedFontsXML(fonts)
	if !strings.Contains(got, "typeface=\"Arial\"") {
		t.Errorf("EmbeddedFontsXML() missing typeface: %v", got)
	}
	if !strings.Contains(got, "regular r:id=\"rId1\"") {
		t.Errorf("EmbeddedFontsXML() missing regular variant: %v", got)
	}
	if !strings.Contains(got, "bold r:id=\"rId2\"") {
		t.Errorf("EmbeddedFontsXML() missing bold variant: %v", got)
	}
}
