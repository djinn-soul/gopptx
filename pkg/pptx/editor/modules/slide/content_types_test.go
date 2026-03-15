package slide

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestContentTypeHelpers(t *testing.T) {
	if got := contentTypeForExtension("jpg"); got != "image/jpeg" {
		t.Fatalf("contentTypeForExtension(jpg)=%q", got)
	}
	if got := contentTypeForExtension(".mp4"); got != "video/mp4" {
		t.Fatalf("contentTypeForExtension(.mp4)=%q", got)
	}
	if got := contentTypeForExtension("unknown"); got != "" {
		t.Fatalf("contentTypeForExtension(unknown)=%q, want empty", got)
	}

	if !isSlidePartOverride("/ppt/slides/slide2.xml") {
		t.Fatal("expected slide override detection to match slide2.xml")
	}
	if isSlidePartOverride("/ppt/notesSlides/notesSlide1.xml") {
		t.Fatal("notes slide should not be treated as normal slide part")
	}

	if !shouldSkipOverridePart("ppt/charts/chart1.xml") {
		t.Fatal("chart overrides should be skipped during rewrite")
	}
	if shouldSkipOverridePart("ppt/media/image1.png") {
		t.Fatal("media paths should not be skipped by override filter")
	}
}

func TestDedupeAndFilterOverrides(t *testing.T) {
	filtered := filterDynamicOverrides([]contentTypeOverride{
		{PartName: "/ppt/slides/slide1.xml", ContentType: common.SlideContentType},
		{PartName: "/ppt/comments/comment1.xml", ContentType: commentsPartType},
		{PartName: "/ppt/customXml/item1.xml", ContentType: "application/xml"},
	}, 1)
	if len(filtered) != 1 || filtered[0].PartName != "/ppt/customXml/item1.xml" {
		t.Fatalf("unexpected filtered overrides: %+v", filtered)
	}

	deduped := dedupeContentTypeOverrides([]contentTypeOverride{
		{PartName: "ppt/slides/slide1.xml", ContentType: "A"},
		{PartName: "/ppt/slides/slide1.xml", ContentType: "B"},
		{PartName: "/ppt/slides/slide2.xml", ContentType: "C"},
	})
	if len(deduped) != 2 {
		t.Fatalf("expected 2 deduped overrides, got %d (%+v)", len(deduped), deduped)
	}
	if deduped[0].PartName != "/ppt/slides/slide1.xml" || deduped[0].ContentType != "B" {
		t.Fatalf("expected last content type to win for slide1: %+v", deduped[0])
	}
}

func TestRewriteContentTypes(t *testing.T) {
	current := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>
  <Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>
</Types>`)

	slides := []common.EditorSlideRef{
		{Part: "ppt/slides/slide2.xml"},
		{Part: "ppt/slides/slide4.xml"},
	}
	out, err := RewriteContentTypes(
		current,
		slides,
		[]string{"ppt/media/video1.mp4", "ppt/media/audio1.wav"},
		true,
		[]string{"ppt/charts/chart1.xml"},
		[]string{"ppt/notesSlides/notesSlide1.xml"},
		[]string{"ppt/theme/theme2.xml"},
		[]string{"ppt/slideLayouts/slideLayout2.xml"},
		[]string{"ppt/slideMasters/slideMaster2.xml"},
		true,
		true,
		[]string{"ppt/comments/comment1.xml"},
		true,
		true,
		[]string{"customXml/itemProps1.xml"},
	)
	if err != nil {
		t.Fatalf("RewriteContentTypes failed: %v", err)
	}

	checks := []string{
		`PartName="/ppt/slides/slide2.xml"`,
		`PartName="/ppt/slides/slide4.xml"`,
		`PartName="/ppt/sectionList.xml"`,
		`PartName="/ppt/vbaProject.bin"`,
		`PartName="/ppt/commentAuthors.xml"`,
		`PartName="/ppt/comments/comment1.xml"`,
		`PartName="/ppt/charts/chart1.xml"`,
		`PartName="/ppt/notesSlides/notesSlide1.xml"`,
		`PartName="/ppt/notesMasters/notesMaster1.xml"`,
		`PartName="/ppt/handoutMasters/handoutMaster1.xml"`,
		`PartName="/customXml/itemProps1.xml"`,
		`Extension="mp4"`,
		`Extension="wav"`,
		`Extension="bin"`,
	}
	for _, c := range checks {
		if !strings.Contains(out, c) {
			t.Fatalf("missing expected rewrite output fragment %q\n%s", c, out)
		}
	}
}

func TestParseContentTypesErrors(t *testing.T) {
	if _, err := parseContentTypesDocument(nil); err == nil {
		t.Fatal("expected missing content types content error")
	}
	if _, err := parseContentTypesDocument([]byte("<Types")); err == nil {
		t.Fatal("expected malformed xml parse error")
	}
}
