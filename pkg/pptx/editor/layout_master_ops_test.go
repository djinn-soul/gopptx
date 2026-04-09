package editor

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLayoutMasterCloneAndRebindRoundTrip(t *testing.T) {
	editor := newLayoutFixtureEditor(t)

	layouts, err := editor.ListSlideLayouts()
	if err != nil {
		t.Fatalf("ListSlideLayouts failed: %v", err)
	}
	if len(layouts) != 1 {
		t.Fatalf("expected one layout, got %d", len(layouts))
	}

	cloned, err := editor.CloneLayoutMasterFamily("ppt/slideLayouts/slideLayout1.xml")
	if err != nil {
		t.Fatalf("CloneLayoutMasterFamily failed: %v", err)
	}
	if cloned.MasterPart == "" || len(cloned.LayoutMap) == 0 {
		t.Fatalf("unexpected clone result: %+v", cloned)
	}

	var clonedLayout string
	for _, v := range cloned.LayoutMap {
		clonedLayout = v
		break
	}
	if clonedLayout == "" {
		t.Fatalf("expected cloned layout")
	}
	if rebindErr := editor.RebindSlideLayout(0, clonedLayout); rebindErr != nil {
		t.Fatalf("RebindSlideLayout failed: %v", rebindErr)
	}

	relsData, _ := editor.parts.Get("ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(string(relsData), "../slideLayouts/"+filepath.Base(clonedLayout)) {
		t.Fatalf("slide rels did not rebind layout: %s", string(relsData))
	}

	out := filepath.Join(t.TempDir(), "layout-clone.pptx")
	if saveErr := editor.Save(out); saveErr != nil {
		t.Fatalf("save failed: %v", saveErr)
	}
	reopen, err := OpenPresentationEditor(out)
	if err != nil {
		t.Fatalf("reopen failed: %v", err)
	}
	defer func() { _ = reopen.Close() }()

	if !reopen.parts.Has(cloned.MasterPart) || !reopen.parts.Has(clonedLayout) {
		t.Fatalf("cloned parts missing after reopen")
	}

	contentTypes, _ := reopen.parts.Get("[Content_Types].xml")
	if !strings.Contains(string(contentTypes), cloned.MasterPart) ||
		!strings.Contains(string(contentTypes), clonedLayout) {
		t.Fatalf("content types missing cloned master/layout")
	}
}

func TestListSlideMastersUsesPresentationRelationships(t *testing.T) {
	editor := newLayoutFixtureEditor(t)
	editor.parts.Set(
		"ppt/slideMasters/slideMaster99.xml",
		[]byte(`<p:sldMaster xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"/>`),
	)

	masters, err := editor.ListSlideMasters()
	if err != nil {
		t.Fatalf("ListSlideMasters failed: %v", err)
	}
	if len(masters) != 1 {
		t.Fatalf("expected one presentation-linked master, got %d", len(masters))
	}
	if masters[0].Part != "ppt/slideMasters/slideMaster1.xml" {
		t.Fatalf("unexpected master part: %s", masters[0].Part)
	}
}

func TestListMasterLayoutsUsesMasterRelationships(t *testing.T) {
	editor := newLayoutFixtureEditor(t)
	editor.parts.Set(
		"ppt/slideLayouts/slideLayout99.xml",
		[]byte(
			`<p:sldLayout xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld name="Orphan"/></p:sldLayout>`,
		),
	)
	editor.parts.Set(
		"ppt/slideLayouts/_rels/slideLayout99.xml.rels",
		[]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="../slideMasters/slideMaster1.xml"/>
</Relationships>`),
	)

	layouts, err := editor.ListMasterLayouts("ppt/slideMasters/slideMaster1.xml")
	if err != nil {
		t.Fatalf("ListMasterLayouts failed: %v", err)
	}
	if len(layouts) != 1 {
		t.Fatalf("expected one master-linked layout, got %d", len(layouts))
	}
	if layouts[0].Part != "ppt/slideLayouts/slideLayout1.xml" {
		t.Fatalf("unexpected layout part: %s", layouts[0].Part)
	}
	if layouts[0].Name != "Title and Content" {
		t.Fatalf("unexpected layout name: %q", layouts[0].Name)
	}
}

func newLayoutFixtureEditor(t *testing.T) *PresentationEditor {
	t.Helper()
	ps := NewPartStore()
	ps.Set("[Content_Types].xml", []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>
<Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>
<Override PartName="/ppt/slideLayouts/slideLayout1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>
<Override PartName="/ppt/slideMasters/slideMaster1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml"/>
<Override PartName="/ppt/theme/theme1.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>
</Types>`))
	ps.Set("ppt/presentation.xml", []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
<p:sldMasterIdLst><p:sldMasterId id="2147483648" r:id="rId1"/></p:sldMasterIdLst>
<p:sldIdLst><p:sldId id="256" r:id="rId2"/></p:sldIdLst>
</p:presentation>`))
	ps.Set("ppt/_rels/presentation.xml.rels", []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="slideMasters/slideMaster1.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/>
</Relationships>`))
	ps.Set(
		"ppt/slides/slide1.xml",
		[]byte(
			`<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`,
		),
	)
	ps.Set("ppt/slides/_rels/slide1.xml.rels", []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
</Relationships>`))
	ps.Set(
		"ppt/slideLayouts/slideLayout1.xml",
		[]byte(
			`<p:sldLayout xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" type="title"><p:cSld name="Title and Content"><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sldLayout>`,
		),
	)
	ps.Set(
		"ppt/slideLayouts/_rels/slideLayout1.xml.rels",
		[]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="../slideMasters/slideMaster1.xml"/>
</Relationships>`),
	)
	ps.Set(
		"ppt/slideMasters/slideMaster1.xml",
		[]byte(
			`<p:sldMaster xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sldMaster>`,
		),
	)
	ps.Set(
		"ppt/slideMasters/_rels/slideMaster1.xml.rels",
		[]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
<Relationship Id="rId7" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme1.xml"/>
</Relationships>`),
	)
	ps.Set(
		"ppt/theme/theme1.xml",
		[]byte(`<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Office Theme"/>`),
	)

	editor, err := newPresentationEditorFromParts(ps)
	if err != nil {
		t.Fatalf("newPresentationEditorFromParts failed: %v", err)
	}
	return editor
}

func TestAddSlideMaster(t *testing.T) {
	editor := newLayoutFixtureEditor(t)

	// Add a new slide master
	masterPart, err := editor.AddSlideMaster()
	if err != nil {
		t.Fatalf("AddSlideMaster failed: %v", err)
	}

	// Verify master was created with correct path
	if masterPart != "ppt/slideMasters/slideMaster2.xml" {
		t.Fatalf("expected masterPart ppt/slideMasters/slideMaster2.xml, got %s", masterPart)
	}

	// Verify master XML exists
	if !editor.parts.Has(masterPart) {
		t.Fatalf("master XML not found: %s", masterPart)
	}

	// Verify master relationships exist
	masterRelsPath := "ppt/slideMasters/_rels/slideMaster2.xml.rels"
	if !editor.parts.Has(masterRelsPath) {
		t.Fatalf("master relationships not found: %s", masterRelsPath)
	}

	// Verify master is listed
	masters, err := editor.ListSlideMasters()
	if err != nil {
		t.Fatalf("ListSlideMasters failed: %v", err)
	}
	if len(masters) != 2 { // Original master1 + new master2
		t.Fatalf("expected 2 masters, got %d", len(masters))
	}
}

func TestRemoveFreshlyAddedSlideMaster(t *testing.T) {
	editor := newLayoutFixtureEditor(t)

	masterPart, err := editor.AddSlideMaster()
	if err != nil {
		t.Fatalf("AddSlideMaster failed: %v", err)
	}

	if err := editor.RemoveSlideMaster(masterPart); err != nil {
		t.Fatalf("RemoveSlideMaster failed: %v", err)
	}

	if editor.parts.Has(masterPart) {
		t.Fatalf("expected removed master part to be deleted: %s", masterPart)
	}
	if !editor.parts.Has("ppt/slideLayouts/slideLayout1.xml") {
		t.Fatalf("expected existing layout to remain after removing unrelated master")
	}

	masters, err := editor.ListSlideMasters()
	if err != nil {
		t.Fatalf("ListSlideMasters failed: %v", err)
	}
	if len(masters) != 1 {
		t.Fatalf("expected 1 master after removal, got %d", len(masters))
	}
}

func TestAddSlideLayout(t *testing.T) {
	editor := newLayoutFixtureEditor(t)

	// Add a new slide layout to existing master
	layoutPart, err := editor.AddSlideLayout("ppt/slideMasters/slideMaster1.xml", "Custom Layout")
	if err != nil {
		t.Fatalf("AddSlideLayout failed: %v", err)
	}

	// Verify layout was created (should be slideLayout2 since slideLayout1 exists)
	if layoutPart != "ppt/slideLayouts/slideLayout2.xml" {
		t.Fatalf("expected layoutPart ppt/slideLayouts/slideLayout2.xml, got %s", layoutPart)
	}

	// Verify layout XML exists
	if !editor.parts.Has(layoutPart) {
		t.Fatalf("layout XML not found: %s", layoutPart)
	}

	// Verify layout relationships exist
	layoutRelsPath := "ppt/slideLayouts/_rels/slideLayout2.xml.rels"
	if !editor.parts.Has(layoutRelsPath) {
		t.Fatalf("layout relationships not found: %s", layoutRelsPath)
	}

	// Verify layout is listed under the master
	layouts, err := editor.ListMasterLayouts("ppt/slideMasters/slideMaster1.xml")
	if err != nil {
		t.Fatalf("ListMasterLayouts failed: %v", err)
	}
	if len(layouts) != 2 { // Original layout1 + new layout2
		t.Fatalf("expected 2 layouts under master, got %d", len(layouts))
	}

	// Verify layout has correct name
	found := false
	for _, l := range layouts {
		if l.Name == "Custom Layout" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected to find layout with name 'Custom Layout'")
	}
}

func TestRemoveSlideLayout(t *testing.T) {
	editor := newLayoutFixtureEditor(t)

	// Add a new layout first
	layoutPart, err := editor.AddSlideLayout("ppt/slideMasters/slideMaster1.xml", "To Be Removed")
	if err != nil {
		t.Fatalf("AddSlideLayout failed: %v", err)
	}

	// Verify layout exists
	layoutsBefore, err := editor.ListMasterLayouts("ppt/slideMasters/slideMaster1.xml")
	if err != nil {
		t.Fatalf("ListMasterLayouts failed: %v", err)
	}
	countBefore := len(layoutsBefore)

	// Remove the layout
	err = editor.RemoveSlideLayout(layoutPart)
	if err != nil {
		t.Fatalf("RemoveSlideLayout failed: %v", err)
	}

	// Verify layout is removed
	if editor.parts.Has(layoutPart) {
		t.Fatalf("layout still exists after removal: %s", layoutPart)
	}

	// Verify layout relationships are removed
	layoutRelsPath := "ppt/slideLayouts/_rels/slideLayout2.xml.rels"
	if editor.parts.Has(layoutRelsPath) {
		t.Fatalf("layout relationships still exist after removal: %s", layoutRelsPath)
	}

	// Verify it's gone from the list
	layoutsAfter, err := editor.ListMasterLayouts("ppt/slideMasters/slideMaster1.xml")
	if err != nil {
		t.Fatalf("ListMasterLayouts failed: %v", err)
	}
	if len(layoutsAfter) != countBefore-1 {
		t.Fatalf("expected %d layouts after removal, got %d", countBefore-1, len(layoutsAfter))
	}
}

func TestAddSlideMasterAndLayoutRoundTrip(t *testing.T) {
	editor := newLayoutFixtureEditor(t)

	// Add a new slide master
	masterPart, err := editor.AddSlideMaster()
	if err != nil {
		t.Fatalf("AddSlideMaster failed: %v", err)
	}

	// Add a layout to the new master
	layoutPart, err := editor.AddSlideLayout(masterPart, "Custom Layout for Master 2")
	if err != nil {
		t.Fatalf("AddSlideLayout failed: %v", err)
	}

	// Save and reopen
	out := filepath.Join(t.TempDir(), "master-layout-add.pptx")
	if err := editor.Save(out); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	reopen, err := OpenPresentationEditor(out)
	if err != nil {
		t.Fatalf("reopen failed: %v", err)
	}
	defer func() { _ = reopen.Close() }()

	// Verify master exists after reopen
	if !reopen.parts.Has(masterPart) {
		t.Fatalf("master missing after reopen: %s", masterPart)
	}

	// Verify layout exists after reopen
	if !reopen.parts.Has(layoutPart) {
		t.Fatalf("layout missing after reopen: %s", layoutPart)
	}

	// Verify content types include both
	contentTypes, _ := reopen.parts.Get("[Content_Types].xml")
	if !strings.Contains(string(contentTypes), masterPart) {
		t.Fatalf("content types missing master: %s", string(contentTypes))
	}
	if !strings.Contains(string(contentTypes), layoutPart) {
		t.Fatalf("content types missing layout: %s", string(contentTypes))
	}
}
