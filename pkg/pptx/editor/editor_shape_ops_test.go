package editor

import (
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestClearShapesRemovesAllShapes(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-clear-all.pptx", []elements.SlideContent{
		elements.NewSlide("Shapes"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	if _, err := ed.AddShape(0, "rect", 100, 100, 600, 400); err != nil {
		t.Fatalf("add shape 1: %v", err)
	}
	if _, err := ed.AddShape(0, "ellipse", 900, 100, 600, 400); err != nil {
		t.Fatalf("add shape 2: %v", err)
	}

	before, err := ed.GetShapes(0)
	if err != nil {
		t.Fatalf("get shapes before clear: %v", err)
	}
	if len(before) < 2 {
		t.Fatalf("expected at least 2 shapes before clear, got %d", len(before))
	}

	if err := ed.ClearShapes(0); err != nil {
		t.Fatalf("clear shapes: %v", err)
	}

	after, err := ed.GetShapes(0)
	if err != nil {
		t.Fatalf("get shapes after clear: %v", err)
	}
	if len(after) != 0 {
		t.Fatalf("expected zero shapes after clear, got %d", len(after))
	}
}

func TestGetShapesIncludesPlaceholderMetadata(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-placeholder-metadata.pptx", []elements.SlideContent{
		elements.NewSlide("Placeholder Metadata"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	ed.parts.Set("ppt/slides/slide1.xml", []byte(
		slideWithBodyAndTitlePlaceholderXML("Body Placeholder", "Title Placeholder"),
	))

	shapes, err := ed.GetShapes(0)
	if err != nil {
		t.Fatalf("get shapes: %v", err)
	}

	for _, shape := range shapes {
		if shape.PlaceholderType != "title" {
			continue
		}
		if shape.PlaceholderIndex == nil || *shape.PlaceholderIndex != 0 {
			t.Fatalf("expected title placeholder index 0, got %#v", shape.PlaceholderIndex)
		}
		return
	}
	t.Fatalf("expected title placeholder metadata in shape listing")
}

func TestGetShapesResolvesActionsAndAltText(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shape-reader-actions.pptx")
	err := writeZipFixture(path, map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
			`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
			`<Default Extension="xml" ContentType="application/xml"/>` +
			`<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>` +
			`<Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>` +
			`<Override PartName="/ppt/slides/slide2.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>` +
			`</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>` +
			`</Relationships>`,
		"ppt/presentation.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
			`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
			`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
			`<p:sldIdLst><p:sldId id="256" r:id="rId1"/><p:sldId id="257" r:id="rId2"/></p:sldIdLst>` +
			`</p:presentation>`,
		"ppt/_rels/presentation.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/>` +
			`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide2.xml"/>` +
			`</Relationships>`,
		"ppt/slides/slide1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
			`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
			`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
			`<p:cSld><p:spTree>` +
			`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
			`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Action Shape" descr="Accessible rectangle">` +
			`<a:hlinkClick r:id="rId5" tooltip="Visit"/>` +
			`<a:hlinkHover action="ppaction://macro?name=HoverMacro"/>` +
			`</p:cNvPr><p:cNvSpPr/><p:nvPr/></p:nvSpPr>` +
			`<p:spPr><a:xfrm><a:off x="10" y="20"/><a:ext cx="300" cy="200"/></a:xfrm></p:spPr>` +
			`<p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr>` +
			`<a:hlinkClick r:id="rId6" tooltip="Run tip"/>` +
			`<a:hlinkMouseOver action="ppaction://hlinkshowjump?jump=nextslide"/>` +
			`</a:rPr><a:t>Hello</a:t></a:r></a:p></p:txBody></p:sp>` +
			`</p:spTree></p:cSld></p:sld>`,
		"ppt/slides/slide2.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
			`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
			`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
			`<p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/>` +
			`<p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId5" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink" Target="https://example.com" TargetMode="External"/>` +
			`<Relationship Id="rId6" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slide2.xml"/>` +
			`</Relationships>`,
		"ppt/slides/_rels/slide2.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`,
	})
	if err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	ed, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	shapes, err := ed.GetShapes(0)
	if err != nil {
		t.Fatalf("get shapes: %v", err)
	}
	if len(shapes) != 1 {
		t.Fatalf("expected one shape, got %d", len(shapes))
	}
	shape := shapes[0]
	if shape.AltText != "Accessible rectangle" {
		t.Fatalf("expected alt text, got %q", shape.AltText)
	}
	if shape.ClickAction == nil || shape.ClickAction.Address == nil ||
		*shape.ClickAction.Address != "https://example.com" {
		t.Fatalf("expected resolved click action address, got %+v", shape.ClickAction)
	}
	if shape.HoverAction == nil || shape.HoverAction.Macro == nil || *shape.HoverAction.Macro != "HoverMacro" {
		t.Fatalf("expected resolved hover macro, got %+v", shape.HoverAction)
	}
	if len(shape.Paragraphs) != 1 || len(shape.Paragraphs[0].Runs) != 1 {
		t.Fatalf("expected one run in one paragraph, got %+v", shape.Paragraphs)
	}
	if shape.Paragraphs[0].Runs[0].Hyperlink == nil || shape.Paragraphs[0].Runs[0].Hyperlink.TargetSlide == nil ||
		*shape.Paragraphs[0].Runs[0].Hyperlink.TargetSlide != 1 {
		t.Fatalf("expected run hyperlink to resolve to slide index 1, got %+v", shape.Paragraphs[0].Runs[0].Hyperlink)
	}
	if shape.Paragraphs[0].Runs[0].HoverAction == nil || shape.Paragraphs[0].Runs[0].HoverAction.TargetJump == nil ||
		*shape.Paragraphs[0].Runs[0].HoverAction.TargetJump != "nextslide" {
		t.Fatalf("expected run hover jump action, got %+v", shape.Paragraphs[0].Runs[0].HoverAction)
	}
}
