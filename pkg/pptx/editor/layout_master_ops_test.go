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
	if err := editor.RebindSlideLayout(0, clonedLayout); err != nil {
		t.Fatalf("RebindSlideLayout failed: %v", err)
	}

	relsData, _ := editor.parts.Get("ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(string(relsData), "../slideLayouts/"+filepath.Base(clonedLayout)) {
		t.Fatalf("slide rels did not rebind layout: %s", string(relsData))
	}

	out := filepath.Join(t.TempDir(), "layout-clone.pptx")
	if err := editor.Save(out); err != nil {
		t.Fatalf("save failed: %v", err)
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
	if !strings.Contains(string(contentTypes), cloned.MasterPart) || !strings.Contains(string(contentTypes), clonedLayout) {
		t.Fatalf("content types missing cloned master/layout")
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
	ps.Set("ppt/slides/slide1.xml", []byte(`<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`))
	ps.Set("ppt/slides/_rels/slide1.xml.rels", []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
</Relationships>`))
	ps.Set("ppt/slideLayouts/slideLayout1.xml", []byte(`<p:sldLayout xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" type="title"><p:cSld name="Title and Content"><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sldLayout>`))
	ps.Set("ppt/slideLayouts/_rels/slideLayout1.xml.rels", []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="../slideMasters/slideMaster1.xml"/>
</Relationships>`))
	ps.Set("ppt/slideMasters/slideMaster1.xml", []byte(`<p:sldMaster xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sldMaster>`))
	ps.Set("ppt/slideMasters/_rels/slideMaster1.xml.rels", []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
<Relationship Id="rId7" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme1.xml"/>
</Relationships>`))
	ps.Set("ppt/theme/theme1.xml", []byte(`<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Office Theme"/>`))

	editor, err := newPresentationEditorFromParts(ps)
	if err != nil {
		t.Fatalf("newPresentationEditorFromParts failed: %v", err)
	}
	return editor
}
