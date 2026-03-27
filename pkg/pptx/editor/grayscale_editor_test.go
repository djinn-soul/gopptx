package editor

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorgrayscale "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/grayscale"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	pptxshapes "github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestConvertToGrayscaleTargetsShapeAndRunSelection(t *testing.T) {
	basePath := writeDeckFixture(t, "grayscale-shapes.pptx", []elements.SlideContent{
		elements.NewSlide("Base"),
	})
	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	shapeID, err := ed.AddShape(0, "rect", 120, 120, 2000, 1000)
	if err != nil {
		t.Fatalf("add shape: %v", err)
	}
	red := "FF0000"
	green := "00FF00"
	blue := "0000FF"
	if err := ed.UpdateShape(0, shapeID, common.ShapeUpdate{
		Fill: &common.ShapeFill{Solid: &red},
		Line: &common.ShapeLine{Color: &green},
		Runs: &[]common.TextRun{
			{Text: "A", Color: &red},
			{Text: "B", Color: &blue},
		},
	}); err != nil {
		t.Fatalf("update shape: %v", err)
	}

	if err := ed.ConvertToGrayscale(editorgrayscale.Options{
		Shapes: []editorgrayscale.ShapeRef{{SlideIndex: 0, ShapeID: shapeID}},
		Text:   []editorgrayscale.TextRef{{SlideIndex: 0, ShapeID: shapeID, RunIndices: []int{1}}},
		Colors: true,
	}); err != nil {
		t.Fatalf("convert to grayscale: %v", err)
	}

	part, _ := ed.parts.Get(ed.slides[0].Part)
	xml := string(part)
	for _, token := range []string{`val="4C4C4C"`, `val="969696"`, `val="1D1D1D"`} {
		if !strings.Contains(xml, token) {
			t.Fatalf("expected grayscale token %q in slide xml: %s", token, xml)
		}
	}
}

func TestConvertToGrayscaleConvertsPictureAndBackground(t *testing.T) {
	pptxPath := createPictureFixturePPTX(t, validFixturePNG(t))
	ed, err := OpenPresentationEditor(pptxPath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	ed.parts.Set("ppt/slides/slide1.xml", []byte(
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`+
			`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`+
			`<p:cSld><p:bg><p:bgPr><a:solidFill><a:srgbClr val="FF0000"/></a:solidFill><a:effectLst/></p:bgPr></p:bg><p:spTree>`+
			`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>`+
			`<p:pic><p:nvPicPr><p:cNvPr id="2" name="Picture 2"/><p:cNvPicPr/><p:nvPr/></p:nvPicPr>`+
			`<p:blipFill><a:blip r:embed="rId2"/><a:stretch><a:fillRect/></a:stretch></p:blipFill>`+
			`<p:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="1000" cy="1000"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr>`+
			`</p:pic></p:spTree></p:cSld></p:sld>`,
	))

	if err := ed.ConvertToGrayscale(editorgrayscale.Options{
		Slides:      []int{0},
		Images:      true,
		Backgrounds: true,
	}); err != nil {
		t.Fatalf("convert to grayscale: %v", err)
	}

	slideXML, _ := ed.parts.Get("ppt/slides/slide1.xml")
	if !strings.Contains(string(slideXML), `val="4C4C4C"`) {
		t.Fatalf("expected grayscale background color, got: %s", slideXML)
	}
	meta, err := ed.GetImageMetadata(0, 2)
	if err != nil {
		t.Fatalf("get image metadata: %v", err)
	}
	if meta.Format != "png" {
		t.Fatalf("expected png grayscale image, got %q", meta.Format)
	}
}

func TestConvertToGrayscaleTargetsPlaceholderTypes(t *testing.T) {
	pptxPath := createPictureFixturePPTX(t, validFixturePNG(t))
	ed, err := OpenPresentationEditor(pptxPath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	ed.parts.Set("ppt/slides/slide1.xml", []byte(
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`+
			`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`+
			`<p:cSld><p:spTree>`+
			`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>`+
			`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Title 1"/><p:cNvSpPr/><p:nvPr><p:ph type="title" idx="0"/></p:nvPr></p:nvSpPr><p:spPr/><p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr><a:solidFill><a:srgbClr val="FF0000"/></a:solidFill></a:rPr><a:t>Title</a:t></a:r></a:p></p:txBody></p:sp>`+
			`<p:sp><p:nvSpPr><p:cNvPr id="3" name="Footer Placeholder"/><p:cNvSpPr/><p:nvPr><p:ph type="ftr" idx="11"/></p:nvPr></p:nvSpPr><p:spPr/><p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr><a:solidFill><a:srgbClr val="0000FF"/></a:solidFill></a:rPr><a:t>Footer</a:t></a:r></a:p></p:txBody></p:sp>`+
			`</p:spTree></p:cSld></p:sld>`,
	))

	if err := ed.ConvertToGrayscale(editorgrayscale.Options{
		Placeholders: []editorgrayscale.PlaceholderRef{
			{SlideIndex: 0, Type: pptxshapes.PlaceholderTypeTitle},
			{SlideIndex: 0, Type: pptxshapes.PlaceholderTypeFtr},
		},
		Colors: true,
	}); err != nil {
		t.Fatalf("convert to grayscale: %v", err)
	}

	slideXML, _ := ed.parts.Get("ppt/slides/slide1.xml")
	xml := string(slideXML)
	for _, token := range []string{`val="4C4C4C"`, `val="1D1D1D"`} {
		if !strings.Contains(xml, token) {
			t.Fatalf("expected grayscale placeholder token %q in slide xml: %s", token, xml)
		}
	}
}

func validFixturePNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	img.Set(1, 0, color.NRGBA{R: 0, G: 255, B: 0, A: 255})
	img.Set(0, 1, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
	img.Set(1, 1, color.NRGBA{R: 255, G: 255, B: 0, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png fixture: %v", err)
	}
	return buf.Bytes()
}
