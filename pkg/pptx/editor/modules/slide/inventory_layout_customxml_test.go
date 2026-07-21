package slide

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type mapPartLookup map[string][]byte

func (m mapPartLookup) Get(partPath string) ([]byte, bool) {
	v, ok := m[partPath]
	return v, ok
}

func TestInventoryParsing(t *testing.T) {
	ps := mapPartLookup{
		"ppt/media/image1.png":  []byte("img-1"),
		"ppt/media/image3.png":  []byte("img-3"),
		"ppt/charts/chart2.xml": []byte("<chart/>"),
		"ppt/charts/_rels/chart2.xml.rels": []byte(
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
				`<Relationship Id="rId1" Type="` + common.RelTypePackage + `" Target="../embeddings/Microsoft_Excel_Worksheet4.xlsx"/>` +
				`</Relationships>`),
		"ppt/slides/_rels/slide1.xml.rels": []byte(
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
				`<Relationship Id="rId2" Type="` + common.RelTypeNotesSlide + `" Target="../notesSlides/notesSlide6.xml"/>` +
				`</Relationships>`),
	}

	mediaInventory, nextMedia := ParseMediaInventory(ps, []string{
		"ppt/media/image1.png",
		"ppt/media/image3.png",
		"ppt/slides/slide1.xml",
	})
	if len(mediaInventory) != 2 || nextMedia != 4 {
		t.Fatalf("unexpected media inventory: inv=%+v next=%d", mediaInventory, nextMedia)
	}

	chartInv, nextChart, nextExcel := ParseChartInventory(ps, []string{
		"ppt/charts/chart2.xml",
		"ppt/charts/chartX.xml",
	})
	if len(chartInv) != 1 ||
		chartInv["ppt/charts/chart2.xml"] != "ppt/embeddings/Microsoft_Excel_Worksheet4.xlsx" {
		t.Fatalf("unexpected chart inventory: %+v", chartInv)
	}
	if nextChart != 3 || nextExcel != 5 {
		t.Fatalf("unexpected next chart/excel values: chart=%d excel=%d", nextChart, nextExcel)
	}

	notesInv, nextNotes := ParseNotesInventory(ps, []string{"ppt/slides/_rels/slide1.xml.rels"})
	if notesInv["ppt/slides/slide1.xml"] != "ppt/notesSlides/notesSlide6.xml" || nextNotes != 7 {
		t.Fatalf("unexpected notes inventory: inv=%+v next=%d", notesInv, nextNotes)
	}
}

func TestLayoutHelpers(t *testing.T) {
	if got := NextMasterPartPath(5); got != "ppt/slideMasters/slideMaster5.xml" {
		t.Fatalf("NextMasterPartPath unexpected: %q", got)
	}

	layoutMap := BuildLayoutCloneMap([]string{"a", "b"}, 10)
	if layoutMap["a"] != "ppt/slideLayouts/slideLayout10.xml" ||
		layoutMap["b"] != "ppt/slideLayouts/slideLayout11.xml" {
		t.Fatalf("unexpected layout clone map: %+v", layoutMap)
	}
	if got := CloneResultTheme("ppt/theme/theme1.xml", ""); got != "ppt/theme/theme1.xml" {
		t.Fatalf("CloneResultTheme fallback unexpected: %q", got)
	}

	nextNum := NextPartNumber(
		[]string{"ppt/slideLayouts/slideLayout2.xml", "ppt/slideLayouts/slideLayout10.xml"},
		mustCompile(`^slideLayout([0-9]+)\.xml$`),
		2,
	)
	if nextNum != 11 {
		t.Fatalf("NextPartNumber=%d, want 11", nextNum)
	}

	layoutName := ParseLayoutName([]byte(`<p:sldLayout name="Title Slide"></p:sldLayout>`))
	if layoutName != "Title Slide" {
		t.Fatalf("ParseLayoutName unexpected: %q", layoutName)
	}

	_, _, err := CloneFamilyInputs(
		"ppt/slideLayouts/slideLayout1.xml",
		func(string) bool { return false },
		common.CanonicalPartPath,
		func(string) (string, error) { return "", nil },
		func(string) ([]string, error) { return nil, nil },
	)
	if err == nil {
		t.Fatal("expected missing layout part error")
	}
	sourceMaster, family, err := CloneFamilyInputs(
		"ppt/slideLayouts/slideLayout1.xml",
		func(string) bool { return true },
		common.CanonicalPartPath,
		func(string) (string, error) { return "ppt/slideMasters/slideMaster1.xml", nil },
		func(string) ([]string, error) { return []string{"ppt/slideLayouts/slideLayout1.xml"}, nil },
	)
	if err != nil || sourceMaster == "" || len(family) != 1 {
		t.Fatalf(
			"CloneFamilyInputs success case failed: master=%q family=%v err=%v",
			sourceMaster,
			family,
			err,
		)
	}

	getPart := func(path string) ([]byte, bool) {
		if path == "ppt/slideLayouts/_rels/slideLayout1.xml.rels" {
			return []byte(
				`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
					`<Relationship Id="rId1" Type="` + common.RelTypeSlideMaster + `" Target="../slideMasters/slideMaster1.xml"/>` +
					`</Relationships>`,
			), true
		}
		return nil, false
	}
	masterPart, err := ResolveLayoutMasterPart(
		"ppt/slideLayouts/slideLayout1.xml",
		getPart,
		ParseRelationshipsXML,
	)
	if err != nil || masterPart != "ppt/slideMasters/slideMaster1.xml" {
		t.Fatalf("ResolveLayoutMasterPart failed: master=%q err=%v", masterPart, err)
	}
	_, err = ResolveLayoutMasterPart(
		"ppt/slideLayouts/slideLayout2.xml",
		getPart,
		ParseRelationshipsXML,
	)
	if err == nil {
		t.Fatal("expected missing layout rels error")
	}

	if !strings.Contains(DefaultSlideMaster(), "<p:sldMaster") {
		t.Fatal("DefaultSlideMaster should render master xml")
	}
	if !strings.Contains(DefaultSlideMasterRelationships(), common.RelTypeTheme) {
		t.Fatal("DefaultSlideMasterRelationships should include theme relationship")
	}
	if !strings.Contains(DefaultSlideLayout("Custom"), `name="Custom"`) {
		t.Fatal("DefaultSlideLayout should include provided layout name")
	}
	if !strings.Contains(DefaultSlideLayoutRelationships(3), "slideMaster3.xml") {
		t.Fatal("DefaultSlideLayoutRelationships should point to target master number")
	}
}

func TestCustomXMLInventoryParsing(t *testing.T) {
	ps := mapPartLookup{
		"customXml/item1.xml": []byte(`<root><name>Alice</name><age>30</age></root>`),
		"customXml/itemProps1.xml": []byte(
			`<ds:datastoreItem ds:itemID="{ID-1}" xmlns:ds="http://schemas.openxmlformats.org/officeDocument/2006/customXml">` +
				`<ds:schemaRefs><ds:schemaRef ds:uri="urn:test"/></ds:schemaRefs>` +
				`</ds:datastoreItem>`,
		),
		"customXml/item2.xml": []byte(`<raw>text</raw>`),
	}
	parts := ParseCustomXMLInventory(ps, []string{
		"customXml/item1.xml",
		"customXml/itemProps1.xml",
		"customXml/item2.xml",
	})
	if len(parts) != 2 {
		t.Fatalf("expected 2 custom xml parts, got %d (%+v)", len(parts), parts)
	}
	if parts[0].Namespace != "urn:test" || parts[0].RootElement != "root" ||
		parts[0].ItemID != "{ID-1}" {
		t.Fatalf("unexpected structured custom xml parse result: %+v", parts[0])
	}
	if len(parts[0].Properties) != 2 || parts[0].Properties[0].Key != "name" {
		t.Fatalf("unexpected structured custom xml properties: %+v", parts[0].Properties)
	}
	if parts[1].Content == "" {
		t.Fatalf("expected raw custom xml fallback content: %+v", parts[1])
	}
}

func mustCompile(pattern string) *regexp.Regexp {
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(errors.New("invalid test regexp: " + err.Error()))
	}
	return re
}
