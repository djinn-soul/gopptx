package slide

import (
	"errors"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestPresentationMutationHelpers(t *testing.T) {
	slides := []common.EditorSlideRef{
		{SlideID: 256, RelID: "rId2", Hidden: true},
		{SlideID: 300, RelID: "rId5"},
	}
	slideListXML := BuildPresentationSlideListXML(slides)
	if !strings.Contains(slideListXML, `id="256"`) ||
		!strings.Contains(slideListXML, `r:id="rId5"`) {
		t.Fatalf("BuildPresentationSlideListXML unexpected output: %s", slideListXML)
	}
	if strings.Contains(slideListXML, ` show="0"`) {
		t.Fatalf("BuildPresentationSlideListXML should not emit hidden show flag on p:sldId: %s", slideListXML)
	}

	source := `<p:presentation><p:sldIdLst><p:sldId id="1" r:id="rId1"/></p:sldIdLst></p:presentation>`
	rewritten, err := RewritePresentationSlideList([]byte(source), slides)
	if err != nil {
		t.Fatalf("RewritePresentationSlideList failed: %v", err)
	}
	if !strings.Contains(rewritten, `id="300"`) {
		t.Fatalf("RewritePresentationSlideList missing replacement content: %s", rewritten)
	}
	if _, err = RewritePresentationSlideList(nil, slides); err == nil {
		t.Fatal("expected missing presentation content error")
	}
	if _, err = RewritePresentationSlideList([]byte(`<p:presentation/>`), slides); err == nil {
		t.Fatal("expected missing <p:sldIdLst> error")
	}

	withNotes, err := RewritePresentationNotesMasterList(
		[]byte(`<p:presentation><p:sldMasterIdLst/><p:sldIdLst></p:sldIdLst></p:presentation>`),
		"rId10",
		true,
	)
	if err != nil || !strings.Contains(withNotes, `<p:notesMasterId r:id="rId10"/>`) {
		t.Fatalf("RewritePresentationNotesMasterList(enable) failed: %s err=%v", withNotes, err)
	}
	withoutNotes, err := RewritePresentationNotesMasterList([]byte(withNotes), "", false)
	if err != nil || strings.Contains(withoutNotes, "notesMasterIdLst") {
		t.Fatalf("RewritePresentationNotesMasterList(disable) failed: %s err=%v", withoutNotes, err)
	}
	if _, err = RewritePresentationNotesMasterList([]byte(`<p:presentation/>`), "", true); err == nil {
		t.Fatal("expected missing relID error when enabling notes master")
	}

	masterXML, err := RewritePresentationSlideMasterList(
		[]byte(
			`<p:presentation><p:sldMasterIdLst><p:sldMasterId id="2147483648" r:id="rId1"/></p:sldMasterIdLst><p:sldIdLst/></p:presentation>`,
		),
		"rId99",
	)
	if err != nil || !strings.Contains(masterXML, `r:id="rId99"`) {
		t.Fatalf("RewritePresentationSlideMasterList failed: %s err=%v", masterXML, err)
	}
	if _, err = RewritePresentationSlideMasterList([]byte(`<p:presentation/>`), ""); err == nil {
		t.Fatal("expected missing relID validation error")
	}
}

func TestSlideHiddenMutationHelpers(t *testing.T) {
	const visibleSlideXML = `<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld/></p:sld>`
	const hiddenSlidePrefix = `<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" show="0">`
	visibleXML := []byte(visibleSlideXML)
	hiddenXML, err := RewriteSlideHidden(visibleXML, true)
	if err != nil {
		t.Fatalf("RewriteSlideHidden(hidden=true) failed: %v", err)
	}
	if !strings.Contains(string(hiddenXML), hiddenSlidePrefix) {
		t.Fatalf("RewriteSlideHidden should add show=0 to p:sld root: %s", string(hiddenXML))
	}

	parsedHidden, err := ParseSlideHidden(hiddenXML)
	if err != nil {
		t.Fatalf("ParseSlideHidden(hidden) failed: %v", err)
	}
	if !parsedHidden {
		t.Fatal("ParseSlideHidden should detect hidden slide")
	}

	restoredXML, err := RewriteSlideHidden(hiddenXML, false)
	if err != nil {
		t.Fatalf("RewriteSlideHidden(hidden=false) failed: %v", err)
	}
	if strings.Contains(string(restoredXML), `show="0"`) {
		t.Fatalf("RewriteSlideHidden(hidden=false) should remove show attr: %s", string(restoredXML))
	}
	parsedVisible, err := ParseSlideHidden(restoredXML)
	if err != nil {
		t.Fatalf("ParseSlideHidden(visible) failed: %v", err)
	}
	if parsedVisible {
		t.Fatal("ParseSlideHidden should report visible slide after show attr removal")
	}
}

func TestSectionMutationHelpers(t *testing.T) {
	sections := []SectionData{{Name: "A&B", GUID: "{GUID-1}", SlideIDs: []int64{256, 300}}}
	sectionListXML := BuildSectionListXML(sections)
	if !strings.Contains(sectionListXML, "A&amp;B") ||
		!strings.Contains(sectionListXML, `id="300"`) {
		t.Fatalf("BuildSectionListXML unexpected output: %s", sectionListXML)
	}

	chartXML := []byte(`<c:externalData r:id="rId1"/>`)
	rewrittenChart := string(RewriteChartExternalData(chartXML, "rId9"))
	if !strings.Contains(rewrittenChart, `r:id="rId9"`) {
		t.Fatalf("RewriteChartExternalData failed: %s", rewrittenChart)
	}

	presXML := `<p:presentation><p:extLst></p:extLst></p:presentation>`
	withSections, err := RewritePresentationSections([]byte(presXML), sections)
	if err != nil || !strings.Contains(withSections, "sectionLst") {
		t.Fatalf("RewritePresentationSections extLst failed: %s err=%v", withSections, err)
	}
	withSections, err = RewritePresentationSections(
		[]byte(`<p:presentation></p:presentation>`),
		sections,
	)
	if err != nil || !strings.Contains(withSections, "<p:extLst>") {
		t.Fatalf("RewritePresentationSections append failed: %s err=%v", withSections, err)
	}
	if _, err = RewritePresentationSections(nil, sections); err == nil {
		t.Fatal("expected empty presentation xml error")
	}

	fontLst := `<p:embeddedFontLst><p:embeddedFont typeface="x"/></p:embeddedFontLst>`
	withFonts, err := RewritePresentationEmbeddedFonts(
		[]byte(`<p:presentation><p:extLst/></p:presentation>`),
		fontLst,
	)
	if err != nil || !strings.Contains(withFonts, "embeddedFontLst") {
		t.Fatalf("RewritePresentationEmbeddedFonts insert failed: %s err=%v", withFonts, err)
	}
	replacedFonts, err := RewritePresentationEmbeddedFonts(
		[]byte(withFonts),
		`<p:embeddedFontLst><p:embeddedFont typeface="y"/></p:embeddedFontLst>`,
	)
	if err != nil || !strings.Contains(replacedFonts, `typeface="y"`) {
		t.Fatalf("RewritePresentationEmbeddedFonts replace failed: %s err=%v", replacedFonts, err)
	}
	if got := ExtractEmbeddedFontList([]byte(replacedFonts)); got == "" {
		t.Fatal("ExtractEmbeddedFontList expected non-empty match")
	}
}

func TestTitleMutationHelpers(t *testing.T) {
	titleShapeXML := []byte(`
<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
  <p:sp><p:nvSpPr><p:nvPr><p:ph type="title"/></p:nvPr></p:nvSpPr><p:txBody><a:p><a:r><a:t>Old</a:t></a:r></a:p></p:txBody></p:sp>
  <p:sp><p:txBody><a:p><a:r><a:t>Other</a:t></a:r></a:p></p:txBody></p:sp>
</p:sld>`)
	updated, ok := ReplaceAllTitleTextRuns(titleShapeXML, "New Title")
	if !ok || !strings.Contains(string(updated), "<a:t>New Title</a:t>") {
		t.Fatalf("ReplaceAllTitleTextRuns failed: %s", string(updated))
	}
	appended := string(AppendCopySuffixToXML(titleShapeXML))
	if !strings.Contains(appended, "Old (Copy)") {
		t.Fatalf("AppendCopySuffixToXML failed: %s", appended)
	}

	noRuns := []byte(`<p:sld><p:sp/></p:sld>`)
	if out, ok := ReplaceAllTitleTextRuns(noRuns, "x"); ok || string(out) != string(noRuns) {
		t.Fatalf("ReplaceAllTitleTextRuns expected no-op, got ok=%v out=%q", ok, string(out))
	}
}

func TestValidationHelpers(t *testing.T) {
	getPart := func(path string) ([]byte, bool) {
		if path == common.SlideRelsPartName("ppt/slides/slide1.xml") {
			return []byte("rels"), true
		}
		return nil, false
	}
	parseOK := func([]byte) ([]common.EditorRelationship, error) {
		return []common.EditorRelationship{
			{Type: common.RelTypeSlideLayout, Target: "../slideLayouts/slideLayout1.xml"},
		}, nil
	}
	rels, err := Relationships("ppt/slides/slide1.xml", getPart, parseOK)
	if err != nil || len(rels) != 1 {
		t.Fatalf("Relationships failed: rels=%+v err=%v", rels, err)
	}
	if _, err = Relationships("ppt/slides/slide2.xml", getPart, parseOK); err == nil {
		t.Fatal("expected missing rels part error")
	}
	if _, err = Relationships("ppt/slides/slide1.xml", getPart, func([]byte) ([]common.EditorRelationship, error) {
		return nil, errors.New("bad xml")
	}); err == nil {
		t.Fatal("expected parse relationship wrapper error")
	}

	notesTarget, err := ScanSupportedSlideRels([]common.EditorRelationship{
		{Type: common.RelTypeSlideLayout},
		{Type: common.RelTypeNotesSlide, Target: "../notesSlides/notesSlide1.xml"},
		{Type: common.RelTypeHyperlink},
	})
	if err != nil || notesTarget != "../notesSlides/notesSlide1.xml" {
		t.Fatalf("ScanSupportedSlideRels failed: notesTarget=%q err=%v", notesTarget, err)
	}
	if _, err = ScanSupportedSlideRels([]common.EditorRelationship{{Type: "unsupported"}}); err == nil {
		t.Fatal("expected unsupported relationship type error")
	}
	if !HasSlideLayoutRelationship([]common.EditorRelationship{{Type: common.RelTypeSlideLayout}}) {
		t.Fatal("HasSlideLayoutRelationship should detect layout rel")
	}
	if HasSlideLayoutRelationship(nil) {
		t.Fatal("HasSlideLayoutRelationship(nil) should be false")
	}

	withImage := elements.NewSlide("t")
	withImage.Images = []shapes.Image{{Path: "img.png"}}
	if !HasImageContent(withImage) {
		t.Fatal("HasImageContent should be true with explicit images")
	}
}
