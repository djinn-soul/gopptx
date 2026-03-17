package layoutmaster

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestAddMasterRelationship(t *testing.T) {
	rels := []common.EditorRelationship{
		{ID: "rId1", Type: common.RelTypeSlide, Target: "slides/slide1.xml"},
	}
	presentationXML := `<p:presentation><p:sldMasterIdLst><p:sldMasterId id="2147483648" r:id="rId9"/></p:sldMasterIdLst><p:sldIdLst/></p:presentation>`

	updatedRels, updatedXML, nextRelIDNum, err := AddMasterRelationship(
		rels,
		presentationXML,
		10,
		"ppt/slideMasters/slideMaster2.xml",
	)
	if err != nil {
		t.Fatalf("AddMasterRelationship failed: %v", err)
	}
	if nextRelIDNum != 11 {
		t.Fatalf("expected nextRelIDNum=11, got %d", nextRelIDNum)
	}
	if len(updatedRels) != 2 || updatedRels[1].ID != "rId10" {
		t.Fatalf("unexpected relationships after add: %+v", updatedRels)
	}
	if !strings.Contains(updatedRels[1].Target, "slideMasters/slideMaster2.xml") {
		t.Fatalf("unexpected relative target for new master rel: %+v", updatedRels[1])
	}
	if !strings.Contains(updatedXML, `r:id="rId10"`) {
		t.Fatalf("expected presentation XML to reference new rel id: %s", updatedXML)
	}

	_, _, _, err = AddMasterRelationship(rels, "", 10, "ppt/slideMasters/slideMaster2.xml")
	if err == nil {
		t.Fatal("expected invalid presentation xml error")
	}
}

func TestLayoutRelationshipHelpers(t *testing.T) {
	masterRels := []common.EditorRelationship{
		{ID: "rId2", Type: common.RelTypeSlideLayout, Target: "../slideLayouts/slideLayout1.xml"},
		{ID: "rId9", Type: common.RelTypeTheme, Target: "../theme/theme1.xml"},
	}
	if got := NextLayoutRelID(masterRels); got != "rId10" {
		t.Fatalf("NextLayoutRelID=%q, want rId10", got)
	}

	appended := AppendLayoutRelationship(
		masterRels,
		"ppt/slideMasters/slideMaster1.xml",
		"ppt/slideLayouts/slideLayout5.xml",
	)
	if len(appended) != len(masterRels)+1 {
		t.Fatalf("expected one appended relationship, got %+v", appended)
	}
	last := appended[len(appended)-1]
	if last.ID != "rId10" || last.Type != common.RelTypeSlideLayout {
		t.Fatalf("unexpected appended relationship: %+v", last)
	}
	if !strings.Contains(last.Target, "slideLayout5.xml") {
		t.Fatalf("unexpected appended relative target: %+v", last)
	}

	filtered := FilterOutRelationshipTarget(appended, last.Target)
	if len(filtered) != len(appended)-1 {
		t.Fatalf("expected one filtered relationship, got %+v", filtered)
	}
}

func TestExtractMasterNumber(t *testing.T) {
	if got := ExtractMasterNumber("ppt/slideMasters/slideMaster7.xml"); got != 7 {
		t.Fatalf("ExtractMasterNumber valid=%d, want 7", got)
	}
	if got := ExtractMasterNumber("ppt/slideMasters/custom.xml"); got != 1 {
		t.Fatalf("ExtractMasterNumber fallback=%d, want 1", got)
	}
}
