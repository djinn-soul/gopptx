package editor

import "testing"

const masterWithLayoutIDList = `<p:sldMaster xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">` +
	`<p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld>` +
	`<p:sldLayoutIdLst><p:sldLayoutId id="2147483649" r:id="rId1"/></p:sldLayoutIdLst>` +
	`</p:sldMaster>`

// SlideMaster.get_layout(slide_layout_id) needs the p:sldLayoutId/@id the master
// lists each layout under, so ListMasterLayouts has to carry it through.
func TestListMasterLayoutsReportsSlideLayoutID(t *testing.T) {
	e := newLayoutFixtureEditor(t)
	e.parts.Set("ppt/slideMasters/slideMaster1.xml", []byte(masterWithLayoutIDList))

	layouts, err := e.ListMasterLayouts("ppt/slideMasters/slideMaster1.xml")
	if err != nil {
		t.Fatalf("ListMasterLayouts: %v", err)
	}
	if len(layouts) != 1 {
		t.Fatalf("expected one layout, got %d", len(layouts))
	}
	if layouts[0].LayoutID != 2147483649 {
		t.Errorf("LayoutID = %d, want 2147483649", layouts[0].LayoutID)
	}
}

// A master with no sldLayoutIdLst still lists its layouts through the
// relationship fallback; those carry no id, which the API reports as zero.
func TestListMasterLayoutsReportsZeroWhenNotInLayoutIDList(t *testing.T) {
	e := newLayoutFixtureEditor(t)

	layouts, err := e.ListMasterLayouts("ppt/slideMasters/slideMaster1.xml")
	if err != nil {
		t.Fatalf("ListMasterLayouts: %v", err)
	}
	if len(layouts) != 1 {
		t.Fatalf("expected one layout, got %d", len(layouts))
	}
	if layouts[0].LayoutID != 0 {
		t.Errorf("LayoutID = %d, want 0 for a layout absent from sldLayoutIdLst", layouts[0].LayoutID)
	}
}

func TestSlideLayoutIDPatternReadsIDAndRelID(t *testing.T) {
	matches := slideLayoutIDPattern.FindAllStringSubmatch(
		`<p:sldLayoutIdLst>`+
			`<p:sldLayoutId id="2147483649" r:id="rId1"/>`+
			`<p:sldLayoutId id="2147483650" r:id="rId2"/>`+
			`</p:sldLayoutIdLst>`,
		-1,
	)

	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	for i, want := range [][2]string{
		{"2147483649", "rId1"},
		{"2147483650", "rId2"},
	} {
		if matches[i][1] != want[0] || matches[i][2] != want[1] {
			t.Errorf("match %d = (%s, %s), want %v", i, matches[i][1], matches[i][2], want)
		}
	}
}
