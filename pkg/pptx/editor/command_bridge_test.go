package editor

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestCommandShapeOps(t *testing.T) {
	basePath := writeDeckFixture(t, "bridge-shape-test.pptx", []elements.SlideContent{
		elements.NewSlide("Shape Test").AddBullet("body"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	// 1. Add Shape
	addReq := `{"api_version":1,"request_id":"r1","op":"add_shape","payload":{"slide_index":0,"type":"rect","x":100,"y":100,"w":1000,"h":500}}`
	resp := ExecuteCommand(e, addReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("add_shape failed: %s", resp)
	}

	var addOut struct {
		Result struct {
			ShapeID int `json:"shape_id"`
		} `json:"result"`
	}
	if unmarshalErr := json.Unmarshal([]byte(resp), &addOut); unmarshalErr != nil {
		t.Fatalf("unmarshal add response: %v", unmarshalErr)
	}
	shapeID := addOut.Result.ShapeID
	if shapeID == 0 {
		t.Fatalf("expected valid shape_id, got 0")
	}

	// 2. List Shapes — should contain the new shape
	listReq := `{"api_version":1,"request_id":"r2","op":"list_shapes","payload":{"slide_index":0}}`
	resp = ExecuteCommand(e, listReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("list_shapes failed: %s", resp)
	}

	// 3. Update Shape
	updateReq := fmt.Sprintf(
		`{"api_version":1,"request_id":"r3","op":"update_shape","payload":{"slide_index":0,"shape_id":%d,"updates":{"text":"Updated"}}}`,
		shapeID,
	)
	resp = ExecuteCommand(e, updateReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("update_shape failed: %s", resp)
	}

	// 4. Remove Shape
	removeReq := fmt.Sprintf(
		`{"api_version":1,"request_id":"r4","op":"remove_shape","payload":{"slide_index":0,"shape_id":%d}}`,
		shapeID,
	)
	resp = ExecuteCommand(e, removeReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("remove_shape failed: %s", resp)
	}
}

func TestCommandNotesOps(t *testing.T) {
	basePath := writeDeckFixture(t, "bridge-notes-test.pptx", []elements.SlideContent{
		elements.NewSlide("Notes Test").AddBullet("body"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	// 1. Set Notes
	setReq := `{"api_version":1,"request_id":"n1","op":"set_notes",` +
		`"payload":{"slide_index":0,"text":"Speaker notes here"}}`
	resp := ExecuteCommand(e, setReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("set_notes failed: %s", resp)
	}

	// 2. Get Notes
	getReq := `{"api_version":1,"request_id":"n2","op":"get_notes","payload":{"slide_index":0}}`
	resp = ExecuteCommand(e, getReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("get_notes failed: %s", resp)
	}
	if !strings.Contains(resp, "Speaker notes here") {
		t.Fatalf("get_notes missing expected text: %s", resp)
	}

	// 3. Update Notes
	setReq2 := `{"api_version":1,"request_id":"n3","op":"set_notes","payload":{"slide_index":0,"text":"Updated notes"}}`
	resp = ExecuteCommand(e, setReq2)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("set_notes update failed: %s", resp)
	}

	// 4. Verify Update
	resp = ExecuteCommand(e, getReq)
	if !strings.Contains(resp, "Updated notes") {
		t.Fatalf("get_notes mismatch after update: %s", resp)
	}
}

func TestCommandUpdateSlidePreservesTitleWhenOmitted(t *testing.T) {
	basePath := writeDeckFixture(t, "bridge-update-slide-title-preserve.pptx", []elements.SlideContent{
		elements.NewSlide("Keep Title").AddBullet("before"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	updateReq := `{"api_version":1,"request_id":"u1","op":"update_slide","payload":{"slide_index":0,"bullets":["after"]}}`
	resp := ExecuteCommand(e, updateReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("update_slide failed: %s", resp)
	}

	slides := e.Slides()
	if len(slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(slides))
	}
	if slides[0].Title != "Keep Title" {
		t.Fatalf("expected title to be preserved, got %q", slides[0].Title)
	}
}

func TestCommandUpdateSlidePreservesTransitionWhenOmitted(t *testing.T) {
	basePath := filepath.Join(t.TempDir(), "bridge-update-slide-transition-preserve.pptx")
	if err := writeZipFixture(basePath, map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
			`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
			`<Default Extension="xml" ContentType="application/xml"/>` +
			`<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>` +
			`<Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>` +
			`</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>` +
			`</Relationships>`,
		"ppt/presentation.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
			`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
			`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
			`<p:sldIdLst><p:sldId id="256" r:id="rId1"/></p:sldIdLst></p:presentation>`,
		"ppt/_rels/presentation.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/>` +
			`</Relationships>`,
		"ppt/slides/slide1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
			`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
			`xmlns:p15="http://schemas.microsoft.com/office/powerpoint/2015/09/main" ` +
			`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
			`<p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/>` +
			`<p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
			`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Title"/>` +
			`<p:cNvSpPr/><p:nvPr/></p:nvSpPr><p:spPr/>` +
			`<p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:t>Keep Transition</a:t></a:r></a:p></p:txBody>` +
			`</p:sp></p:spTree></p:cSld>` +
			`<p:transition><p:extLst><p:ext uri="{AE3914FA-7E93-4B9E-9A96-D1E12CAF14E6}">` +
			`<p15:morph option="obj"/></p:ext></p:extLst></p:transition>` +
			`</p:sld>`,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`,
	}); err != nil {
		t.Fatalf("write fixture deck: %v", err)
	}
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	updateReq := `{"api_version":1,"request_id":"u2","op":"update_slide","payload":{"slide_index":0,"bullets":["after"]}}`
	resp := ExecuteCommand(e, updateReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("update_slide failed: %s", resp)
	}

	partName := e.Slides()[0].PartName
	slideXML, ok := e.parts.Get(partName)
	if !ok {
		t.Fatalf("missing slide part %q", partName)
	}
	if !strings.Contains(string(slideXML), `<p:ext uri="{AE3914FA-7E93-4B9E-9A96-D1E12CAF14E6}">`) {
		t.Fatalf("expected preserved morph transition ext URI in updated slide XML")
	}
	if !strings.Contains(string(slideXML), `<p15:morph`) {
		t.Fatalf("expected preserved morph transition node in updated slide XML")
	}
}
