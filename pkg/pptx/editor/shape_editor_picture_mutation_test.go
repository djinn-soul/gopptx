package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestUpdateShapePictureAppliesCropRotationAndFlip(t *testing.T) {
	pptxPath := createPictureFixturePPTX(t, testutil.TinyPNG())
	ed, err := OpenPresentationEditor(pptxPath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	rotation := 30.0
	flipH := true
	flipV := false
	err = ed.UpdateShape(0, 2, common.ShapeUpdate{
		Crop: &common.ImageCrop{
			Left:   0.1,
			Right:  0.2,
			Top:    0.3,
			Bottom: 0.4,
		},
		Rotation: &rotation,
		FlipH:    &flipH,
		FlipV:    &flipV,
	})
	if err != nil {
		t.Fatalf("update picture shape: %v", err)
	}

	slideBytes, ok := ed.parts.Get("ppt/slides/slide1.xml")
	if !ok {
		t.Fatal("slide part not found after update")
	}
	slideXML := string(slideBytes)
	if !strings.Contains(slideXML, `<a:srcRect l="10000" r="20000" t="30000" b="40000"/>`) {
		t.Fatalf("expected crop srcRect in slide xml, got: %s", slideXML)
	}
	if !strings.Contains(slideXML, `<a:xfrm rot="1800000" flipH="1">`) {
		t.Fatalf("expected updated xfrm attrs in slide xml, got: %s", slideXML)
	}
	if strings.Contains(slideXML, `flipV="1"`) {
		t.Fatalf("unexpected flipV attr in slide xml: %s", slideXML)
	}
}

func TestUpdateShapePictureRejectsCropForNonPicture(t *testing.T) {
	content := strings.Replace(
		pictureFixtureSlideXML(),
		`</p:pic></p:spTree>`,
		`</p:pic><p:sp><p:nvSpPr><p:cNvPr id="3" name="Rect 3"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr><p:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="1000" cy="1000"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr></p:sp></p:spTree>`,
		1,
	)
	pptxPath := createPictureFixtureWithSlideXML(t, []byte(content), testutil.TinyPNG())

	ed, err := OpenPresentationEditor(pptxPath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	err = ed.UpdateShape(0, 3, common.ShapeUpdate{
		Crop: &common.ImageCrop{Left: 0.2},
	})
	if err == nil {
		t.Fatal("expected non-picture crop update to fail")
	}
	if !strings.Contains(err.Error(), "not a picture shape") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func pictureFixtureSlideXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
		`<p:cSld><p:spTree>` +
		`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
		`<p:pic><p:nvPicPr><p:cNvPr id="2" name="Picture 2"/><p:cNvPicPr/><p:nvPr/></p:nvPicPr>` +
		`<p:blipFill><a:blip r:embed="rId2"/><a:stretch><a:fillRect/></a:stretch></p:blipFill>` +
		`<p:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="1000" cy="1000"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr>` +
		`</p:pic></p:spTree></p:cSld></p:sld>`
}

func createPictureFixtureWithSlideXML(t *testing.T, slideXML []byte, imageData []byte) string {
	t.Helper()
	pptxPath := createPictureFixturePPTX(t, imageData)

	ed, err := OpenPresentationEditor(pptxPath)
	if err != nil {
		t.Fatalf("open fixture editor: %v", err)
	}
	ed.parts.Set("ppt/slides/slide1.xml", slideXML)
	if err := ed.Save(pptxPath); err != nil {
		_ = ed.Close()
		t.Fatalf("save fixture editor: %v", err)
	}
	_ = ed.Close()
	return pptxPath
}
